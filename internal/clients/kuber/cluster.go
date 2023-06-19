package kuber

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	clientsetUffizziClusterV1 "github.com/UffizziCloud/uffizzi-cluster-operator/clientset/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Client) CreateCluster(
	namespace string,
	name string,
	manifest string,
) (*v1alpha1.UffizziCluster, error) {
	clusterSpec := clientsetUffizziClusterV1.UffizziClusterProps{
		Name: name,
		Spec: v1alpha1.UffizziClusterSpec{
			Manifests: &manifest,
		},
	}

	return client.uffizziClusterClient.UffizziClusterV1(namespace).Create(clusterSpec)
}

func (client *Client) GetCluster(
	namespace string,
	name string,
) (*v1alpha1.UffizziCluster, error) {
	return client.uffizziClusterClient.UffizziClusterV1(namespace).Get(name, metav1.GetOptions{})
}
