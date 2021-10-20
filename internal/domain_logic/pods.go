package domain

import (
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (l *Logic) GetContainers(deploymentID uint64) ([]v1.Pod, error) {
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	pods, err := l.KuberClient.GetPods(namespace.Name)

	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func (l *Logic) GetContainersMetrics(deploymentID uint64) ([]v1beta1.PodMetrics, error) {
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	metrics, err := l.KuberClient.GetPodsMetrics(namespace.Name)

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
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	events, err := l.KuberClient.ListEvents(namespace.Name)

	return events, err
}
