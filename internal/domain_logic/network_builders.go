package domain

import (
	"fmt"
	"log"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/pkg/networks"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
)

type INetworkBuilder interface {
	GetDeploymentSelectorName() (string, error)
	Create() error
	AwaitNetworkCreation() error
}

type NetworkDependencies struct {
	DomainLogic    *Logic
	Namespace      *corev1.Namespace
	AppName        string
	ContainerList  domainTypes.ContainerList
	Deployment     *appsv1.Deployment
	DeploymentHost string
	Project        domainTypes.Project
}

func NewNetworkDependencies(
	domainLogic *Logic,
	namespace *corev1.Namespace,
	appName string,
	containerList domainTypes.ContainerList,
	deployment *appsv1.Deployment,
	deploymentHost string,
	project domainTypes.Project,
) *NetworkDependencies {
	return &NetworkDependencies{
		DomainLogic:    domainLogic,
		Namespace:      namespace,
		AppName:        appName,
		ContainerList:  containerList,
		Deployment:     deployment,
		DeploymentHost: deploymentHost,
		Project:        project,
	}
}

type IngressNetworkBuilder struct {
	Network *NetworkDependencies
	Service *corev1.Service
	Ingress *v1.Ingress
}

func NewIngressNetworkBuilder(network *NetworkDependencies) INetworkBuilder {
	return &IngressNetworkBuilder{
		Network: network,
	}
}

func (builder *IngressNetworkBuilder) GetDeploymentSelectorName() (string, error) {
	network := builder.Network
	deployment := network.Deployment

	selector, ok := deployment.Spec.Selector.MatchLabels["app"]
	if !ok {
		err := fmt.Errorf(
			"deployment/%s : Cannot find 'app' in Spec.Selector.MatchLabels",
			deployment.Name,
		)

		return "", err
	}

	return selector, nil
}

func (builder *IngressNetworkBuilder) Create() error {
	network := builder.Network
	namespace := network.Namespace
	domainLogic := network.DomainLogic
	deploymentHost := network.DeploymentHost
	project := network.Project

	deploymentSelector, err := builder.GetDeploymentSelectorName()
	if err != nil {
		return err
	}

	publicContainerList := network.ContainerList.GetPublicContainerList()

	service, err := domainLogic.KuberClient.CreateOrUpdateService(namespace, deploymentSelector, publicContainerList)
	if err != nil {
		return err
	}

	log.Printf("service/%s configured", service.Name)

	service, err = domainLogic.KuberClient.AwaitServiceStatus(service)
	if err != nil {
		return err
	}

	log.Printf("service/%s got load balancer", service.Name)

	generalPublicContainer, err := publicContainerList.GetPublicContainer()
	if err != nil {
		return err
	}

	ingress, err := domainLogic.KuberClient.CreateOrUpdateIngress(
		namespace,
		service.Name,
		generalPublicContainer,
		deploymentHost,
		project,
	)
	if err != nil {
		return err
	}

	log.Printf("ingress/%s configured", ingress.Name)

	ingress, err = domainLogic.KuberClient.AwaitIngressStatus(ingress)
	if err != nil {
		return err
	}

	log.Printf("ingress/%s got load balancer", ingress.Name)

	namespace, err = domainLogic.KuberClient.UpdateAnnotationNamespace(
		namespace.Name,
		"ingress_ip",
		networks.GetIngresEntrypoint(ingress.Status.LoadBalancer.Ingress[0]),
	)
	if err != nil {
		return err
	}

	log.Printf("namespace/%s annotations.ingress_ip configured", namespace.Name)

	builder.Service = service
	builder.Ingress = ingress

	return nil
}

func (builder *IngressNetworkBuilder) AwaitNetworkCreation() error {
	domainLogic := builder.Network.DomainLogic
	resources := &networkingResources{
		builder.Ingress,
		builder.Service,
		builder.Network.ContainerList,
	}

	err := domainLogic.RunNetworkConnectivityChecks(
		builder.Network.Namespace,
		resources,
		builder.Network.Deployment,
	)
	if err != nil {
		return err
	}

	return nil
}
