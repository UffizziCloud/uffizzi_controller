package domain

import (
	"encoding/base64"
	"log"

	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	types "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

func (l *Logic) getClusterNameBy(
	namespaceName string,
) string {
	return namespaceName + "-" + "vcluster"
}

func (l *Logic) mapUffizziClusterToCluster(
	ufizziCluster *v1alpha1.UffizziCluster,
) *types.Cluster {
	cluster := &types.Cluster{
		Name:      ufizziCluster.ObjectMeta.Name,
		Namespace: ufizziCluster.ObjectMeta.Namespace,
		UID:       string(ufizziCluster.ObjectMeta.UID),
	}

	cluster.Status.Ready = ufizziCluster.Status.Ready

	if !ufizziCluster.Status.Ready {
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
	namespaceName string,
) (*types.Cluster, error) {
	namespace, err := l.KuberClient.FindNamespace(namespaceName)

	if err != nil {
		return nil, err
	}

	log.Printf("namespace/%s found", namespace.Name)

	clusterName := l.getClusterNameBy(namespaceName)

	ufizziCluster, err := l.KuberClient.CreateCluster(
		namespace.Name,
		clusterName,
	)

	if err != nil {
		log.Printf("ClusterError: %s", err)
		return nil, err
	}

	cluster := l.mapUffizziClusterToCluster(ufizziCluster)

	return cluster, err
}

func (l *Logic) GetCluster(
	namespaceName string,
) (*types.Cluster, error) {
	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	log.Printf("namespace/%s found", namespace.Name)

	clusterName := l.getClusterNameBy(namespaceName)

	ufizziCluster, err := l.KuberClient.GetCluster(
		namespace.Name,
		clusterName,
	)
	if err != nil {
		log.Printf("ClusterError: %s", err)
		return nil, err
	}

	cluster := l.mapUffizziClusterToCluster(ufizziCluster)

	return cluster, err
}
