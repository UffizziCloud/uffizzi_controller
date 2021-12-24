package domain

import (
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (l *Logic) GetContainers(deploymentID uint64) ([]v1.Pod, error) {
	namespace := l.KubernetesNamespaceName(deploymentID)

	pods, err := l.KuberClient.GetPods(namespace)

	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func (l *Logic) GetContainersMetrics(deploymentID uint64) ([]v1beta1.PodMetrics, error) {
	namespace := l.KubernetesNamespaceName(deploymentID)

	metrics, err := l.KuberClient.GetPodsMetrics(namespace)

	if err != nil {
		return nil, err
	}

	return metrics.Items, nil
}

func (l *Logic) GetPodLogs(
	deploymentID uint64,
	podName string,
	containerName string,
	limit int64,
) ([]string, error) {
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	namespace, _ := l.KuberClient.FindNamespace(namespaceName)
	logs, err := l.KuberClient.GetPodLogs(namespace.Name, podName, containerName, limit)

	return logs, err
}

func (l *Logic) GetPodEvents(
	deploymentID uint64,
) (*v1.EventList, error) {
	namespace := l.KubernetesNamespaceName(deploymentID)

	events, err := l.KuberClient.ListEvents(namespace)

	return events, err
}
