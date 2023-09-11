package kuber

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	clientsetUffizziClusterV1 "github.com/UffizziCloud/uffizzi-cluster-operator/clientset/v1alpha1"

	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Client) CreateCluster(
	namespace string,
	clusterParams domainTypes.ClusterParams,
) (*v1alpha1.UffizziCluster, error) {
	clusterSpec := clientsetUffizziClusterV1.UffizziClusterProps{
		Name: clusterParams.Name,
		Spec: v1alpha1.UffizziClusterSpec{
			Manifests: &clusterParams.Manifest,
			Ingress: v1alpha1.UffizziClusterIngress{
				Host: clusterParams.BaseIngressHost,
			},
			ResourceQuota: &clusterParams.ResourceSettings.ResourceQuota,
			LimitRange:    &clusterParams.ResourceSettings.LimitRange,
		},
	}

	return client.uffizziClusterClient.UffizziClusterV1(namespace).Create(clusterSpec)
}

func (client *Client) GetCluster(
	name string,
	namespace string,
) (*v1alpha1.UffizziCluster, error) {
	return client.uffizziClusterClient.UffizziClusterV1(namespace).Get(name, metav1.GetOptions{})
}
