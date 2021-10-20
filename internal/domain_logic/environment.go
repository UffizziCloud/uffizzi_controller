package domain

import (
	"log"

	// "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	// "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/jobs"
	corev1 "k8s.io/api/core/v1"
)

func (l *Logic) GetNamespace(deploymentID uint64) (*corev1.Namespace, error) {
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (l *Logic) CreateNamespace(deploymentID uint64, kind string) (*corev1.Namespace, error) {
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	namespace, err := l.KuberClient.CreateNamespace(namespaceName, kind)
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (l *Logic) UpdateNamespace(deploymentID uint64, kind string) (*corev1.Namespace, error) {
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	namespace.Labels["kind"] = kind

	namespace, err = l.KuberClient.UpdateNamespace(namespace)
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (l *Logic) DeleteNamespace(deploymentID uint64) error {
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	log.Printf("Clear all deployment resources namespace: %v\n", namespaceName)

	err := l.KuberClient.RemoveNamespace(namespaceName)
	if err != nil {
		return err
	}

	return nil
}
