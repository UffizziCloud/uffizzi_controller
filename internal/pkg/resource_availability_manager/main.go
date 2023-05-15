package resource_availability_manager

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/avast/retry-go"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
)

func NewResourceAvailabilityManager(settings ResourceAvailabilitySettings) *ResourceAvailabilityManager {
	return &ResourceAvailabilityManager{
		settings: settings,
	}
}

type NetworkPointType string
type ResourceAvailabilityRequestType string

const (
	NetworkPointIngress                     NetworkPointType                = "ingress"
	NetworkPointService                     NetworkPointType                = "service"
	ResourceAvailabilityRequestIngress      ResourceAvailabilityRequestType = "ingress"
	ResourceAvailabilityRequestLoadbalancer ResourceAvailabilityRequestType = "loadbalancer"
	FailedHttpCode                          int                             = 503
	HttpsPortNumber                         int                             = 443
)

type ResourceAvailabilityManager struct {
	settings ResourceAvailabilitySettings
}

type ResourceAvailabilityRequest struct {
	Entrypoint string
	Points     []ResourceAvailabilityPoint
	Kind       ResourceAvailabilityRequestType
}

type ResourceAvailabilityPoint struct {
	Port    int
	Kind    NetworkPointType
	Payload map[string]string
}

type ResourceAvailabilitySettings struct {
	IPPingTimeout                time.Duration
	PerAddressTimeout            time.Duration
	PerAddressAttempts           uint
	ResourceRequestBackOffPeriod time.Duration
}

type ResourceAvailabilityPointResponse struct {
	Entrypoint string
	Kind       NetworkPointType
	Status     bool
	Payload    map[string]string
}

func (r *ResourceAvailabilityManager) CheckResourceAvailabilityByTcp(ctx context.Context,
	ch chan ResourceAvailabilityPointResponse,
	request *ResourceAvailabilityRequest) {
	log.Printf("start naive TCP checks for: %v", request)

	dialer := net.Dialer{
		Timeout: r.settings.IPPingTimeout,
	}

	for _, point := range request.Points {
		go func(point ResourceAvailabilityPoint) {
			err := r.checkAddressConnectivityWithLoggerByTcp(ctx, &dialer, request.Entrypoint, point.Port)

			pointResponse := ResourceAvailabilityPointResponse{
				Entrypoint: request.Entrypoint,
				Kind:       point.Kind,
				Status:     err == nil,
				Payload:    point.Payload,
			}

			ch <- pointResponse
		}(point)
	}
}

func (r *ResourceAvailabilityManager) checkAddressConnectivityWithLoggerByTcp(ctx context.Context,
	dialer *net.Dialer,
	entrypoint string, port int) error {
	err := r.checkAddressConnectivityByTcp(ctx, dialer, entrypoint, port)
	if err != nil {
		log.Printf("[tcp] checks failed for %s:%v with err: %s", entrypoint, port, err)
		return err
	}

	log.Printf("[tcp] checks finished for %s:%v", entrypoint, port)

	return nil
}

func (r *ResourceAvailabilityManager) checkAddressConnectivityByTcp(ctx context.Context,
	dialer *net.Dialer, entrypoint string, port int) error {
	address := fmt.Sprintf("%s:%v", entrypoint, port)

	retryableFunc := func() error {
		network := "tcp"

		log.Printf("[tcp] dial %s", address)

		var err error

		var conn net.Conn

		if port == global.Settings.IngressDefaultPort {
			conn, err = tls.DialWithDialer(dialer, network, address, nil)
		} else {
			conn, err = dialer.DialContext(ctx, network, address)
		}

		if err == nil {
			defer conn.Close()
		}

		select {
		case <-ctx.Done():
			return retry.Unrecoverable(ctx.Err())
		case <-time.After(r.settings.PerAddressTimeout):
			return nil
		default:
			return err
		}
	}

	err := retry.Do(
		retryableFunc,
		retry.LastErrorOnly(true),
		retry.Attempts(r.settings.PerAddressAttempts),
		retry.Delay(r.settings.ResourceRequestBackOffPeriod),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("[tcp] dial failed for %s with err %s, retry attempt: %d", address, err, n)
		}),
	)

	return err
}

func (r *ResourceAvailabilityManager) CheckResourceAvailabilityByHttp(ctx context.Context,
	ch chan ResourceAvailabilityPointResponse,
	request *ResourceAvailabilityRequest) {
	log.Printf("start naive HTTP checks for: %v", request)

	for _, point := range request.Points {
		go func(point ResourceAvailabilityPoint) {
			err := r.checkAddressConnectivityWithLoggerByHttp(ctx, request.Entrypoint, point.Port)

			pointResponse := ResourceAvailabilityPointResponse{
				Entrypoint: request.Entrypoint,
				Kind:       point.Kind,
				Status:     err == nil,
				Payload:    point.Payload,
			}

			ch <- pointResponse
		}(point)
	}
}

func (r *ResourceAvailabilityManager) checkAddressConnectivityWithLoggerByHttp(
	ctx context.Context,
	entrypoint string,
	port int) error {
	err := r.CheckHttpStatus(ctx, entrypoint, port)
	if err != nil {
		log.Printf("[http] checks failed for %s with err: %s", entrypoint, err)
		return err
	}

	log.Printf("[http] checks finished for %s", entrypoint)

	return nil
}

func (r *ResourceAvailabilityManager) CheckHttpStatus(ctx context.Context, address string, port int) error {
	var protocol string

	if port == HttpsPortNumber {
		protocol = "https"
	} else {
		protocol = "http"
	}

	retryableFunc := func() error {
		log.Printf("[%s] dial %s", protocol, address)

		var err error

		resp, err := http.Get(fmt.Sprintf("%s://%s", protocol, address))

		if err == nil {
			defer resp.Body.Close()
		}

		if err != nil {
			return err
		}

		if resp.StatusCode == FailedHttpCode {
			err = fmt.Errorf("Status is %s", resp.Status)
		}

		select {
		case <-ctx.Done():
			return retry.Unrecoverable(ctx.Err())
		case <-time.After(r.settings.PerAddressTimeout):
			return nil
		default:
			return err
		}
	}

	err := retry.Do(
		retryableFunc,
		retry.LastErrorOnly(true),
		retry.Attempts(r.settings.PerAddressAttempts),
		retry.Delay(r.settings.ResourceRequestBackOffPeriod),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("[%s] dial failed for %s with err %s, retry attempt: %d", protocol, address, err, n)
		}),
	)

	return err
}
