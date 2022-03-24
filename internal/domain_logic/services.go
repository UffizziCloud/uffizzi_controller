package domain

import (
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	v1 "k8s.io/api/core/v1"
)

func (l *Logic) GetServices(deploymentID uint64) ([]v1.Service, error) {
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	services, err := l.KuberClient.GetServices(namespace.Name)
	if err != nil {
		return nil, err
	}

	return services.Items, nil
}

func (l *Logic) GetDefaultIngressService() (*v1.Service, error) {
	services, err := l.KuberClient.GetServices(global.Settings.KubernetesNamespace)
	if err != nil {
		return nil, err
	}

	return &services.Items[0], nil
}
