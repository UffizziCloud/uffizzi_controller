package types

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
)

type Cluster struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
	Status    struct {
		Ready      bool   `json:"ready"`
		KubeConfig string `json:"kubeConfig"`
		Host       string `json:"host"`
	} `json:"status"`
}

type ClusterParams struct {
	Name             string                  `json:"name"`
	Manifest         string                  `json:"manifest"`
	BaseIngressHost  string                  `json:"base_ingress_host"`
	ResourceSettings ClusterResourceSettings `json:"resource_settings"`
}

type ClusterResourceSettings struct {
	ResourceQuota v1alpha1.UffizziClusterResourceQuota `json:"resourceQuota,omitempty"`
	LimitRange    v1alpha1.UffizziClusterLimitRange    `json:"limitRange,omitempty"`
}
