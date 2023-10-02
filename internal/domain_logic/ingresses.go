package domain

import (
	networkingV1 "k8s.io/api/networking/v1"
)

func (l *Logic) GetIngresses(namespaceName string) (*networkingV1.IngressList, error) {
	ingresses, err := l.KuberClient.GetIngresses(namespaceName)
	if err != nil {
		return nil, err
	}

	return ingresses, nil
}
