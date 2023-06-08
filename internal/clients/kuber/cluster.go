package kuber

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiUffizziClusterV1 "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	clientsetUffizziClusterV1 "github.com/UffizziCloud/uffizzi-cluster-operator/clientset/v1alpha1"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

func (client *Client) CreateCluster(
	namespace string,
	name string,
	helm []apiUffizziClusterV1.HelmChart,
	ingressService domainTypes.ClusterIngressService,
) error {
	cluster := clientsetUffizziClusterV1.UffizziClusterProps{
		Name: name,
		Spec: apiUffizziClusterV1.UffizziClusterSpec{
			Ingress: apiUffizziClusterV1.UffizziClusterIngress{
				Host:  "deleteme2023.uffizzi.cloud",
				Class: "nginx",
				Services: []apiUffizziClusterV1.ExposedVClusterService{
					{
						Name:      ingressService.Name,
						Namespace: ingressService.Namespace,
						Port:      ingressService.Port,
					},
				},
			},
			Helm: helm,
		},
	}

	_, err := client.uffizziClusterClient.UffizziClusterV1(namespace).Create(cluster)
	if err != nil {
		log.Println("Vcluster ERROR")
		log.Println(err)
	}

	return nil
}

func (client *Client) GetCluster(namespace, name string) {
	_, err := client.uffizziClusterClient.UffizziClusterV1(namespace).Get(name, metav1.GetOptions{})

	if err != nil {
		log.Println("Vcluster ERROR")
		log.Println(err)
	}
}

func (client *Client) GetClusters(namespace string) {
	_, err := client.uffizziClusterClient.UffizziClusterV1(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println("Vcluster ERROR")
		log.Println(err)
	}
}

func (client *Client) DeleteCluster(namespace, name string) {
	err := client.uffizziClusterClient.UffizziClusterV1(namespace).Delete(name)
	if err != nil {
		log.Println("Vcluster ERROR")
		log.Println(err)
	}
}
