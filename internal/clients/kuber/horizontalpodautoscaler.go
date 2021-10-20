package kuber

import (
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Client) findHorizontalPodAutoscaler(
	namespaceName, deploymentName string,
) (*autoscalingV1.HorizontalPodAutoscaler, error) {
	autoscalers := client.clientset.AutoscalingV1().HorizontalPodAutoscalers(namespaceName)
	return autoscalers.Get(client.context, deploymentName, metav1.GetOptions{})
}

func (client *Client) findOrInitializeHorizontalPodAutoscaler(
	namespaceName, deploymentName string,
	minReplicas, maxReplicas int32,
) (*autoscalingV1.HorizontalPodAutoscaler, error) {
	autoscaler, err := client.findHorizontalPodAutoscaler(namespaceName, deploymentName)

	if err != nil && !errors.IsNotFound(err) {
		return autoscaler, err
	}

	if len(autoscaler.UID) == 0 {
		autoscaler = initializeHorizontalPodAutoscaler(
			namespaceName,
			deploymentName,
			deploymentName,
			minReplicas,
			maxReplicas,
		)
	}

	return autoscaler, nil
}

func (client *Client) DeleteHorizontalPodAutoscalerIfExists(
	namespace *v1.Namespace, deploymentName string,
) error {
	autoscaler, err := client.findHorizontalPodAutoscaler(
		namespace.Name,
		deploymentName,
	)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if len(autoscaler.UID) != 0 {
		autoscalers := client.clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace.Name)
		err = autoscalers.Delete(client.context, autoscaler.Name, metav1.DeleteOptions{})

		return err
	}

	return nil
}

func (client *Client) CreateOrUpdateHorizontalPodAutoscaler(
	namespace *v1.Namespace, deploymentName string,
	minReplicas, maxReplicas int32,
) (*autoscalingV1.HorizontalPodAutoscaler, error) {
	autoscaler, err := client.findOrInitializeHorizontalPodAutoscaler(
		namespace.Name,
		deploymentName,
		minReplicas,
		maxReplicas,
	)

	if err != nil {
		return nil, err
	}

	autoscalers := client.clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace.Name)

	if len(autoscaler.UID) > 0 {
		autoscaler, err = autoscalers.Update(client.context, autoscaler, metav1.UpdateOptions{})
	} else {
		autoscaler, err = autoscalers.Create(client.context, autoscaler, metav1.CreateOptions{})
	}

	return autoscaler, err
}
