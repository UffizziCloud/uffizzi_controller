package kuber

import (
	"context"

	"k8s.io/client-go/rest"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	uffizziClusterOperator "github.com/UffizziCloud/uffizzi-cluster-operator/clientset/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type Client struct {
	clientset            *kubernetes.Clientset
	uffizziClusterClient *uffizziClusterOperator.V1Alpha1Clientset
	metricClient         *metrics.Clientset
	context              context.Context
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

	uffizziClusterClient, err := uffizziClusterOperator.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	client := &Client{
		clientset:            clientset,
		uffizziClusterClient: uffizziClusterClient,
		metricClient:         metricClient,
		context:              context.Background(),
	}

	return client, nil
}
