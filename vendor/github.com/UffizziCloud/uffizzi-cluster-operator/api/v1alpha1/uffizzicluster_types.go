/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/fluxcd/pkg/apis/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HelmReleaseInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type HelmChartInfo struct {
	Name    string `json:"name"`
	Repo    string `json:"repo"`
	Version string `json:"version,omitempty"`
}

type HelmChart struct {
	Chart   HelmChartInfo   `json:"chart"`
	Values  string          `json:"values,omitempty"`
	Release HelmReleaseInfo `json:"release"`
}

type VClusterIngressSpec struct {
	IngressAnnotations    map[string]string `json:"ingressAnnotations,omitempty"`
	CertManagerTLSEnabled bool              `json:"certManagerTLSEnabled,omitempty"`
}

// UffiClusterIngress defines the ingress capabilities of the cluster,
// the basic host can be setup for all
type UffizziClusterIngress struct {
	Host  string `json:"host,omitempty"`
	Class string `json:"class,omitempty"`
}

// UffizziClusterAPIServer defines the API server capabilities of the cluster
type UffizziClusterAPIServer struct {
	Image string `json:"image,omitempty"`
}

// UffizziClusterResourceQuota defines the resource quota which defines the
// quota of resources a namespace has access to
type UffizziClusterResourceQuota struct {
	//+kubebuilder:default:=true
	Enabled bool `json:"enabled"`
	//+kubebuilder:default:={}
	Requests UffizziClusterRequestsQuota `json:"requests,omitempty"`
	//+kubebuilder:default:={}
	Limits UffizziClusterResourceQuotaLimits `json:"limits,omitempty"`
	//+kubebuilder:default:={}
	Services UffizziClusterServicesQuota `json:"services,omitempty"`
	//+kubebuilder:default:={}
	Count UffizziClusterResourceCount `json:"count,omitempty"`
}

type UffizziClusterLimitRange struct {
	//+kubebuilder:default:=true
	Enabled bool `json:"enabled"`
	//+kubebuilder:default:={}
	Default UffizziClusterLimitRangeDefault `json:"default,omitempty"`
	//+kubebuilder:default:={}
	DefaultRequest UffizziClusterLimitRangeDefaultRequest `json:"defaultRequest,omitempty"`
}

type UffizziClusterRequestsQuota struct {
	//+kubebuilder:default:="0.5"
	CPU string `json:"cpu,omitempty"`
	//+kubebuilder:default:="1Gi"
	Memory string `json:"memory,omitempty"`
	//+kubebuilder:default:="5Gi"
	EphemeralStorage string `json:"ephemeralStorage,omitempty"`
	//+kubebuilder:default:="5Gi"
	Storage string `json:"storage,omitempty"`
}

type UffizziClusterResourceQuotaLimits struct {
	//+kubebuilder:default:="0.5"
	CPU string `json:"cpu,omitempty"`
	//+kubebuilder:default:="8Gi"
	Memory string `json:"memory,omitempty"`
	//+kubebuilder:default:="5Gi"
	EphemeralStorage string `json:"ephemeralStorage,omitempty"`
}

type UffizziClusterLimitRangeDefault struct {
	//+kubebuilder:default:="0.5"
	CPU string `json:"cpu,omitempty"`
	//+kubebuilder:default:="1Gi"
	Memory string `json:"memory,omitempty"`
	//+kubebuilder:default:="8Gi"
	EphemeralStorage string `json:"ephemeralStorage,omitempty"`
}

type UffizziClusterLimitRangeDefaultRequest struct {
	//+kubebuilder:default:="0.1"
	CPU string `json:"cpu,omitempty"`
	//+kubebuilder:default:="128Mi"
	Memory string `json:"memory,omitempty"`
	//+kubebuilder:default:="1Gi"
	EphemeralStorage string `json:"ephemeralStorage,omitempty"`
}

type UffizziClusterServicesQuota struct {
	//+kubebuilder:default:=0
	NodePorts int `json:"nodePorts,omitempty"`
	//+kubebuilder:default:=3
	LoadBalancers int `json:"loadBalancers,omitempty"`
}

type UffizziClusterResourceCount struct {
	//+kubebuilder:default:=20
	Pods int `json:"pods,omitempty"`
	//+kubebuilder:default:=10
	Services int `json:"services,omitempty"`
	//+kubebuilder:default:=20
	ConfigMaps int `json:"configMaps,omitempty"`
	//+kubebuilder:default:=20
	Secrets int `json:"secrets,omitempty"`
	//+kubebuilder:default:=10
	PersistentVolumeClaims int `json:"persistentVolumeClaims,omitempty"`
	//+kubebuilder:default:=10
	Endpoints int `json:"endpoints,omitempty"`
}

// UffizziClusterSpec defines the desired state of UffizziCluster
type UffizziClusterSpec struct {
	//+kubebuilder:default:="k3s"
	//+kubebuilder:validation:Enum=k3s;k8s
	Distro        string                       `json:"distro,omitempty"`
	APIServer     UffizziClusterAPIServer      `json:"apiServer,omitempty"`
	Ingress       UffizziClusterIngress        `json:"ingress,omitempty"`
	TTL           string                       `json:"ttl,omitempty"`
	Helm          []HelmChart                  `json:"helm,omitempty"`
	Manifests     *string                      `json:"manifests,omitempty"`
	ResourceQuota *UffizziClusterResourceQuota `json:"resourceQuota,omitempty"`
	LimitRange    *UffizziClusterLimitRange    `json:"limitRange,omitempty"`
	Sleep         bool                         `json:"sleep,omitempty"`
	//+kubebuilder:default:="sqlite"
	//+kubebuilder:validation:Enum=etcd;sqlite
	ExternalDatastore string `json:"externalDatastore,omitempty"`
}

// UffizziClusterStatus defines the observed state of UffizziCluster
type UffizziClusterStatus struct {
	Conditions                 []metav1.Condition `json:"conditions,omitempty"`
	HelmReleaseRef             *string            `json:"helmReleaseRef,omitempty"`
	KubeConfig                 VClusterKubeConfig `json:"kubeConfig,omitempty"`
	Host                       *string            `json:"host,omitempty"`
	LastAppliedConfiguration   *string            `json:"lastAppliedConfiguration,omitempty"`
	LastAppliedHelmReleaseSpec *string            `json:"lastAppliedHelmReleaseSpec,omitempty"`
	LastAwakeTime              metav1.Time        `json:"lastAwakeTime,omitempty"`
}

// VClusterKubeConfig is the KubeConfig SecretReference of the related VCluster
type VClusterKubeConfig struct {
	SecretRef *meta.SecretKeyReference `json:"secretRef,omitempty"`
}

type UffizziClusterDistro struct {
	//+kubebuilder:default:="k3s"
	Type string `json:"type,omitempty"`
	//+kubebuilder:default:="v1.27.3-k3s1"
	Version string `json:"version,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=uc;ucluster
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="APIReady",type=string,JSONPath=`.status.conditions[?(@.type=='APIReady')].status`
//+kubebuilder:printcolumn:name="DataStoreReady",type=string,JSONPath=`.status.conditions[?(@.type=='DataStoreReady')].status`
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=='Ready')].status`
//+kubebuilder:printcolumn:name="Sleep",type=string,JSONPath=`.status.conditions[?(@.type=='Sleep')].status`
//+kubebuilder:printcolumn:name="Host",type=string,JSONPath=`.status.host`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="UptimeSinceLastAwake",type=date,JSONPath=`.status.lastAwakeTime`

// UffizziCluster is the Schema for the UffizziClusters API
type UffizziCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UffizziClusterSpec   `json:"spec,omitempty"`
	Status UffizziClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UffizziClusterList contains a list of UffizziCluster
type UffizziClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UffizziCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UffizziCluster{}, &UffizziClusterList{})
}
