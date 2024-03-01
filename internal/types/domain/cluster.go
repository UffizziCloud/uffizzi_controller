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
		Sleep      bool   `json:"sleep"`
		KubeConfig string `json:"kubeConfig"`
		Host       string `json:"host"`
	} `json:"status"`
}

type ClusterParams struct {
	Name             string                  `json:"name"`
	Manifest         string                  `json:"manifest"`
	BaseIngressHost  string                  `json:"base_ingress_host"`
	ResourceSettings ClusterResourceSettings `json:"resource_settings"`
	Distro           string                  `json:"distro"`
	Image            string                  `json:"image"`
	AutoSleep        bool                    `json:"auto_sleep,omitempty"`
	Provider         string                  `json:"provider,omitempty"`
}

type PatchClusterParams struct {
	Manifest         string                  `json:"manifest"`
	BaseIngressHost  string                  `json:"base_ingress_host"`
	ResourceSettings ClusterResourceSettings `json:"resource_settings"`
	Sleep            bool                    `json:"sleep"`
}

type ClusterResourceSettings struct {
	ResourceQuota v1alpha1.UffizziClusterResourceQuota `json:"resourceQuota,omitempty"`
	LimitRange    v1alpha1.UffizziClusterLimitRange    `json:"limitRange,omitempty"`
}
