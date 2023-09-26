package domain

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	types "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (l *Logic) mapUffizziClusterToCluster(
	ufizziCluster *v1alpha1.UffizziCluster,
) *types.Cluster {
	cluster := &types.Cluster{
		Name:      ufizziCluster.ObjectMeta.Name,
		Namespace: ufizziCluster.ObjectMeta.Namespace,
		UID:       string(ufizziCluster.ObjectMeta.UID),
	}

	readyStatus := isClusterReady(ufizziCluster)
	cluster.Status.Ready = readyStatus

	sleepStatus := isClusterAsleep(ufizziCluster)
	cluster.Status.Sleep = sleepStatus

	if !readyStatus {
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
	cluster.Status.Host = *ufizziCluster.Status.Host

	return cluster
}

func (l *Logic) CreateCluster(
	namespaceName string,
	clusterParams types.ClusterParams,
) (*types.Cluster, error) {
	namespace, err := l.KuberClient.FindNamespace(namespaceName)

	if err != nil {
		return nil, err
	}

	log.Printf("namespace/%s found", namespace.Name)

	ufizziCluster, err := l.KuberClient.CreateCluster(namespace.Name, clusterParams)

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

func (l *Logic) PatchCluster(
	clusterName string,
	namespaceName string,
	patchClusterParams types.PatchClusterParams,
) error {
	namespace, err := l.KuberClient.FindNamespace(namespaceName)

	if err != nil {
		return err
	}

	log.Printf("namespace/%s found", namespace.Name)

	err = l.KuberClient.PatchCluster(
		clusterName,
		namespaceName,
		patchClusterParams,
	)

	if err != nil {
		log.Printf("ClusterError: %s", err)
		return err
	}

	return nil
}

func isClusterReady(ufizziCluster *v1alpha1.UffizziCluster) bool {
	if len(ufizziCluster.Status.Conditions) == 0 {
		return false
	}

	readyCondition, err := getConditionByName(ufizziCluster.Status.Conditions, "APIReady")

	if err != nil {
		log.Printf("ClusterError: %s", err)
		return false
	}

	return readyCondition.Status == "True"
}

func isClusterAsleep(ufizziCluster *v1alpha1.UffizziCluster) bool {
	if len(ufizziCluster.Status.Conditions) == 0 {
		return false
	}

	sleepCondition, err := getConditionByName(ufizziCluster.Status.Conditions, "Sleep")

	if err != nil {
		log.Printf("ClusterError: %s", err)
		return false
	}

	return sleepCondition.Status == "True"
}

func getConditionByName(conditions []metav1.Condition, conditionType string) (metav1.Condition, error) {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return condition, nil
		}
	}

	return conditions[0], fmt.Errorf("Status with type %v not found", conditionType)
}
