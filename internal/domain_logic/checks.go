package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	networkConnectivity "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/domain_logic/network_connectivity"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"

	//"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/pkg/networks"
	availabilityManager "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/pkg/resource_availability_manager"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type networkingResources struct {
	Ingress       *networkingV1.Ingress
	Service       *corev1.Service
	ContainerList domainTypes.ContainerList
}

func (nR *networkingResources) String() string {
	return fmt.Sprintf("networkingResources{Ingress:%+#v, Service:%+#v, ContainerList:%+#v}",
		nR.Ingress,
		nR.Service,
		nR.ContainerList)
}

func (l *Logic) BuildResourceAvailabilityRequests(containerList domainTypes.ContainerList,
	service *corev1.Service,
	ingress *networkingV1.Ingress,
) ([]availabilityManager.ResourceAvailabilityRequest, int, error) {
	loadBalancerRequest := availabilityManager.ResourceAvailabilityRequest{}
	loadBalancerRequest.Entrypoint = service.Spec.ClusterIP
	loadBalancerRequest.Points = []availabilityManager.ResourceAvailabilityPoint{}

	for _, container := range containerList.Items {
		containerPort := *container.Port

		loadBalancerPoint := availabilityManager.ResourceAvailabilityPoint{
			Port: int(containerPort),
			Kind: availabilityManager.NetworkPointService,
			Payload: map[string]string{
				"id": fmt.Sprint(container.ID),
			},
		}

		loadBalancerRequest.Points = append(loadBalancerRequest.Points, loadBalancerPoint)
	}

	requests := []availabilityManager.ResourceAvailabilityRequest{loadBalancerRequest}

	if ingress != nil {
		ingressRequest := availabilityManager.ResourceAvailabilityRequest{}
		ingressRequest.Entrypoint = ingress.Spec.Rules[0].Host

		generalPublicContainer, err := containerList.GetPublicContainer()
		if err != nil {
			return requests, 0, err
		}

		ingressPoint := availabilityManager.ResourceAvailabilityPoint{
			Port: global.Settings.IngressDefaultPort,
			Kind: availabilityManager.NetworkPointIngress,
			Payload: map[string]string{
				"id": fmt.Sprint(generalPublicContainer.ID),
			},
		}

		ingressRequest.Points = append(ingressRequest.Points, ingressPoint)

		requests = append(requests, ingressRequest)
	}

	countPoints := 0
	for _, request := range requests {
		countPoints += len(request.Points)
	}

	return requests, countPoints, nil
}

func (l *Logic) UpdateNetworkConnectivity(namespace *corev1.Namespace,
	networkConnectivity *networkConnectivity.ConnectivityResponse) error {
	networkConnectivityJson, _ := json.Marshal(networkConnectivity)

	_, err := l.KuberClient.UpdateAnnotationNamespace(
		namespace.Name,
		"network_connectivity",
		string(networkConnectivityJson),
	)

	return err
}

func (l *Logic) ResetNetworkConnectivityTemplateForIngress(
	namespace *corev1.Namespace,
	containerList domainTypes.ContainerList,
) error {
	publicContainerList := containerList.GetPublicContainerList()

	networkConnectivityTemplate, err := networkConnectivity.NewNetworkConnectivityTemplate(publicContainerList)
	if err != nil {
		return err
	}

	ingressContainer := containerList.GetIngressContainer()

	networkConnectivityTemplate.AddIngressContainer(ingressContainer)

	err = l.UpdateNetworkConnectivity(namespace, networkConnectivityTemplate)

	return err
}

func (l *Logic) AddNetworkConnectivityTemplateForIngress(
	namespace *corev1.Namespace,
	containerList domainTypes.ContainerList,
) error {
	networkConnectivityTemplate, err := l.DecodeNetworkConnectivityJson(namespace.Annotations["network_connectivity"])
	if err != nil {
		return err
	}

	ingressContainer := containerList.GetIngressContainer()

	networkConnectivityTemplate.AddIngressContainer(ingressContainer)

	err = l.UpdateNetworkConnectivity(namespace, networkConnectivityTemplate)

	return err
}

func (l *Logic) ResetNetworkConnectivityTemplate(

	namespace *corev1.Namespace,
	containerList domainTypes.ContainerList,
) error {
	return l.ResetNetworkConnectivityTemplateForIngress(namespace, containerList)
}

func (l *Logic) DecodeNetworkConnectivityJson(
	networkConnectivityJson string,
) (*networkConnectivity.ConnectivityResponse, error) {
	networkConnectivityTemplate := networkConnectivity.ConnectivityResponse{}

	err := json.Unmarshal([]byte(networkConnectivityJson), &networkConnectivityTemplate)
	if err != nil {
		return nil, err
	}

	return &networkConnectivityTemplate, nil
}

func (l *Logic) MarkUnresponsiveContainersAsFailed(namespaceName string) error {
	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return err
	}

	networkConnectivityTemplate, err := l.DecodeNetworkConnectivityJson(namespace.Annotations["network_connectivity"])
	if err != nil {
		return err
	}

	for key := range networkConnectivityTemplate.Containers {
		if networkConnectivityTemplate.Containers[key].Ingress != nil {
			if networkConnectivityTemplate.Containers[key].Ingress.Status == networkConnectivity.StatusPending {
				networkConnectivityTemplate.Containers[key].Ingress.Status = networkConnectivity.StatusFailed
			}
		}

		if networkConnectivityTemplate.Containers[key].Service.Status == networkConnectivity.StatusPending {
			networkConnectivityTemplate.Containers[key].Service.Status = networkConnectivity.StatusFailed
		}
	}

	err = l.UpdateNetworkConnectivity(namespace, networkConnectivityTemplate)

	return err
}

func (l *Logic) UpdateContainerInNetworkConnectivity(
	namespaceName string,
	kind availabilityManager.NetworkPointType,
	entrypoint, containerID string,
	status networkConnectivity.ConnectivityStatus,
) error {
	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return err
	}

	networkConnectivityObj, err := l.DecodeNetworkConnectivityJson(namespace.Annotations["network_connectivity"])
	if err != nil {
		return err
	}

	if kind == availabilityManager.NetworkPointIngress {
		networkConnectivityObj.SetIngressStatus(containerID, status, entrypoint)
	} else {
		networkConnectivityObj.SetLoadBalancerStatus(containerID, status, entrypoint)
	}

	err = l.UpdateNetworkConnectivity(namespace, networkConnectivityObj)

	return err
}

func (l *Logic) WatchChangeDeployment(
	ctx context.Context,
	cancelContext context.CancelFunc,
	deployment *appsv1.Deployment,
) {
	oldGeneration := deployment.Generation

	for {
		deployment, err := l.KuberClient.FindDeployment(deployment.Namespace, deployment.Name)

		if errors.IsNotFound(err) || deployment.Generation != oldGeneration {
			cancelContext()
		}

		var delayOnGettingDeployment = 5 * time.Second

		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(delayOnGettingDeployment)
		}
	}
}

func (l *Logic) CheckResourcesAvailability(
	ctx context.Context,
	namespaceName string,
	resources *networkingResources,
) error {
	ingress := resources.Ingress

	service, err := l.KuberClient.AwaitServiceStatus(resources.Service)
	if err != nil {
		return err
	}

	if ingress != nil {
		ingress, err = l.KuberClient.AwaitIngressStatus(resources.Ingress)
		if err != nil {
			return err
		}
	}

	publicContainerList := resources.ContainerList.GetPublicContainerList()

	requests, countPoints, err := l.BuildResourceAvailabilityRequests(publicContainerList, service, ingress)
	if err != nil {
		return err
	}

	settings := availabilityManager.ResourceAvailabilitySettings{
		IPPingTimeout:                global.Settings.ServiceChecks.IPPingTimeout,
		PerAddressTimeout:            global.Settings.ServiceChecks.PerAddressTimeout,
		PerAddressAttempts:           global.Settings.ServiceChecks.PerAddressAttempts,
		ResourceRequestBackOffPeriod: global.Settings.ResourceRequestBackOffPeriod,
	}

	resourceAvailability := availabilityManager.NewResourceAvailabilityManager(settings)

	ch := make(chan availabilityManager.ResourceAvailabilityPointResponse, countPoints)

	for _, request := range requests {
		go func(request availabilityManager.ResourceAvailabilityRequest) {
			resourceAvailability.CheckResourceAvailability(ctx, ch, &request)
		}(request)
	}

	for i := 0; i < countPoints; i++ {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(global.Settings.ServiceChecks.AvailabilityTimeout):
			return fmt.Errorf("resources check timeout")
		case point := <-ch:
			status := networkConnectivity.StatusSuccess

			if !point.Status {
				status = networkConnectivity.StatusFailed
			}

			err = l.UpdateContainerInNetworkConnectivity(
				namespaceName, point.Kind, point.Entrypoint, point.Payload["id"], status)
			if err != nil {
				log.Printf("error UpdateContainerInNetworkConnectivity %+#v\n", err)
			}
		}
	}

	return nil
}

func (l *Logic) RunNetworkConnectivityChecks(namespace *corev1.Namespace,
	resources *networkingResources, deployment *appsv1.Deployment) error {
	ctx := context.TODO()
	ctx, cancelContext := context.WithCancel(ctx)

	defer cancelContext()

	go l.WatchChangeDeployment(ctx, cancelContext, deployment)

	errorChanel := make(chan error)

	go func() {
		errorChanel <- l.CheckResourcesAvailability(ctx, namespace.Name, resources)
	}()

	select {
	case <-ctx.Done():
		log.Println("Cancel to check resources availability")
		return nil

	case value := <-errorChanel:
		if value != nil {
			return l.handleDomainDeploymentError(namespace.Name, value)
		}

		return nil
	}
}
