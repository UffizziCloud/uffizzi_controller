package domain

import corev1 "k8s.io/api/core/v1"

func isPodsReady(pods []corev1.Pod) bool {
	for _, pod := range pods {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if !containerStatus.Ready {
				return false
			}
		}
	}

	return true
}
