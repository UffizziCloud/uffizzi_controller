package domain

import (
	v1 "k8s.io/api/core/v1"
)

func (l *Logic) GetNodes() ([]v1.Node, error) {
	nodes, err := l.KuberClient.GetNodes()

	if err != nil {
		return nil, err
	}

	return nodes.Items, nil
}
