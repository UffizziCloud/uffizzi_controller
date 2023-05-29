package domain

import (
	"log"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

func (l *Logic) CreateCluster(
	deploymentID uint64,
	clusterName string,
	helm []domainTypes.Helm,
	ingressService domainTypes.ClusterIngressService,
	deploymentHost string,
) error {
	namespaceName := l.KubernetesNamespaceName(deploymentID)

	namespace, err := l.KuberClient.FindNamespace(namespaceName)
	if err != nil {
		return err
	}

	log.Printf("namespace/%s found", namespace.Name)

	namespace, err = l.ResetNamespaceErrors(namespace)
	if err != nil {
		return err
	}

	policyName := global.Settings.ResourceName.Policy(namespace.Name)
	policy, err := l.KuberClient.FindOrCreateNetworkPolicy(namespace.Name, policyName)
	if err != nil {
		return err
	}

	log.Printf("networkPolicy/%s configured", policy.Name)

	deployment, err := l.KuberClient.CreateCluster(
		namespace,
		clusterName,
		helm,
		ingressService,
	)
	if err != nil {
		return l.handleDomainDeploymentError(namespace.Name, err)
	}

	return nil
}
