package kuber

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Client) GetNodes() (*v1.NodeList, error) {
	nodes, err := client.clientset.CoreV1().Nodes().List(client.context, metav1.ListOptions{})

	return nodes, err
}
