package domain

import (
	"encoding/base64"

	"log"

	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	types "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

func (l *Logic) mapUffizziClusterToCluster(
	ufizziCluster *v1alpha1.UffizziCluster,
) *types.Cluster {
	cluster := &types.Cluster{
		Name:      ufizziCluster.ObjectMeta.Name,
		Namespace: ufizziCluster.ObjectMeta.Namespace,
		UID:       string(ufizziCluster.ObjectMeta.UID),
	}

	status := clusterStatus(ufizziCluster)
	cluster.Status.Ready = status

	if !status {
		return cluster
	}

	secret, err := l.KuberClient.GetSecret(ufizziCluster.ObjectMeta.Namespace,
		ufizziCluster.Status.KubeConfig.SecretRef.Name,
	)

	if err != nil {
		return cluster
	}

	kubeConfigData, ok := secret.Data["config"]
	if !ok {
		return cluster
	}

	cluster.Status.KubeConfig = base64.StdEncoding.EncodeToString(kubeConfigData)

	return cluster
}

func (l *Logic) CreateCluster(
	clusterName string,
	namespaceName string,
	manifest string,
	baseIngressHost string,
) (*types.Cluster, error) {
	namespace, err := l.KuberClient.FindNamespace(namespaceName)

	if err != nil {
		return nil, err
	}

	log.Printf("namespace/%s found", namespace.Name)

	ufizziCluster, err := l.KuberClient.CreateCluster(
		clusterName,
		namespace.Name,
		manifest,
		baseIngressHost,
	)

	if err != nil {
		log.Printf("ClusterError: %s", err)
		return nil, err
	}

	cluster := l.mapUffizziClusterToCluster(ufizziCluster)

	return cluster, err
}

func (l *Logic) GetCluster(
	clusterName string,
	namespaceName string,
) (*types.Cluster, error) {
	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	log.Printf("namespace/%s found", namespace.Name)

	ufizziCluster, err := l.KuberClient.GetCluster(
		clusterName,
		namespace.Name,
	)
	if err != nil {
		log.Printf("ClusterError: %s", err)
		return nil, err
	}

	cluster := l.mapUffizziClusterToCluster(ufizziCluster)

	return cluster, err
}

func clusterStatus(ufizziCluster *v1alpha1.UffizziCluster) bool {
	if len(ufizziCluster.Status.Conditions) == 0 {
		return false
	}

	if ufizziCluster.Status.Conditions[0].Status == "False" {
		return false
	} else {
		return true
	}
}
