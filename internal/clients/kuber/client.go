package kuber

import (
	"context"
	"net/http"

	"k8s.io/client-go/rest"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	"k8s.io/client-go/kubernetes"
)

type Client struct {
	clientset    *kubernetes.Clientset
	metricClient *metrics.Clientset
	context      context.Context
}

func NewClient(config *rest.Config) (*Client, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	metricClient, err := metrics.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	client := &Client{
		clientset:    clientset,
		metricClient: metricClient,
		context:      context.Background(),
	}

	return client, nil
}

func NewClient2(config *rest.Config, httpClient *http.Client) (*Client, error) {
	clientset, err := kubernetes.NewForConfigAndClient(config, httpClient)
	if err != nil {
		return nil, err
	}

	metricClient, err := metrics.NewForConfigAndClient(config, httpClient)
	if err != nil {
		return nil, err
	}

	client := &Client{
		clientset:    clientset,
		metricClient: metricClient,
		context:      context.Background(),
	}

	return client, nil
}
