package kuber

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Client) ListEvents(namespace string) (*v1.EventList, error) {
	events, err := client.clientset.CoreV1().Events(namespace).List(client.context, metav1.ListOptions{})

	return events, err
}
