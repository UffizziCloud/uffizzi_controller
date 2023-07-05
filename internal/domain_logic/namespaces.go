package domain

import (
	"log"

	// "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	// "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/jobs"
	corev1 "k8s.io/api/core/v1"
)

func (l *Logic) GetNamespaceV2(namespaceName string) (*corev1.Namespace, error) {
	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (l *Logic) CreateNamespaceV2(namespaceName string) (*corev1.Namespace, error) {
	namespace, err := l.KuberClient.CreateNamespace(namespaceName)

	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (l *Logic) DeleteNamespaceV2(namespaceName string) error {
	log.Printf("Clear all resources namespace: %v\n", namespaceName)

	err := l.KuberClient.RemoveNamespace(namespaceName)
	if err != nil {
		return err
	}

	return nil
}
