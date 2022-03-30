package kuber

import (
	"log"
	"time"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Client) GetIngresses(namespace string) (*networkingV1.IngressList, error) {
	ingresses := client.clientset.NetworkingV1().Ingresses(namespace)

	ingressList, err := ingresses.List(client.context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return ingressList, nil
}

func (client *Client) FindIngress(namespace, name string) (*networkingV1.Ingress, error) {
	ingresses := client.clientset.NetworkingV1().Ingresses(namespace)
	ingress, err := ingresses.Get(client.context, name, metav1.GetOptions{})

	return ingress, err
}

func (client *Client) findOrInitializeIngress(namespace, ingressName string) (*networkingV1.Ingress, error) {
	ingress, err := client.FindIngress(namespace, ingressName)

	if err != nil && !errors.IsNotFound(err) {
		return ingress, err
	}

	if len(ingress.UID) == 0 {
		ingress = initializeIngress(namespace, ingressName)
	}

	return ingress, nil
}

func (client *Client) GetOrPrepareIngressName(namespace *corev1.Namespace) (string, error) {
	var annotationKey = "ingressName"

	ingressName := namespace.Annotations[annotationKey]

	if len(ingressName) > 0 {
		return ingressName, nil
	}

	ingressName = generateIngressName()

	_, err := client.UpdateAnnotationNamespace(namespace.Name, annotationKey, ingressName)

	return ingressName, err
}

func (client *Client) UpdateIngressAttributes(
	ingress *networkingV1.Ingress,
	namespace *corev1.Namespace,
	container domainTypes.Container,
	serviceName, deploymentHost string) *networkingV1.Ingress {
	containerPort := *container.Port

	ingress.ObjectMeta.Annotations["kubernetes.io/ingress.class"] = "nginx"
	ingress.ObjectMeta.Annotations["cert-manager.io/cluster-issuer"] = global.Settings.CertManagerClusterIssuer

	tls := []networkingV1.IngressTLS{
		{Hosts: []string{deploymentHost}, SecretName: deploymentHost},
	}

	ingressBackend := networkingV1.IngressBackend{
		Service: &networkingV1.IngressServiceBackend{
			Name: serviceName,
			Port: networkingV1.ServiceBackendPort{
				Number: containerPort,
			},
		},
	}

	// constants are not addressable.
	pathTypePrefix := networkingV1.PathTypePrefix

	paths := []networkingV1.HTTPIngressPath{
		{
			Path:     "/",
			PathType: &pathTypePrefix,
			Backend:  ingressBackend,
		},
	}

	ingressRuleValue := networkingV1.IngressRuleValue{
		HTTP: &networkingV1.HTTPIngressRuleValue{Paths: paths},
	}

	rules := []networkingV1.IngressRule{
		{Host: deploymentHost, IngressRuleValue: ingressRuleValue},
	}

	ingress.Spec = networkingV1.IngressSpec{TLS: tls, Rules: rules}

	return ingress
}

func (client *Client) CreateOrUpdateIngress(namespace *corev1.Namespace,
	serviceName string,
	container domainTypes.Container,
	deploymentHost string) (*networkingV1.Ingress, error) {
	ingressName, err := client.GetOrPrepareIngressName(namespace)
	if err != nil {
		return nil, err
	}

	ingresses := client.clientset.NetworkingV1().Ingresses(namespace.Name)

	ingress, err := client.findOrInitializeIngress(namespace.Name, ingressName)
	if err != nil {
		return ingress, err
	}

	ingress = client.UpdateIngressAttributes(ingress, namespace, container, serviceName, deploymentHost)

	if len(ingress.UID) > 0 {
		ingress, err = ingresses.Update(client.context, ingress, metav1.UpdateOptions{})
	} else {
		ingress, err = ingresses.Create(client.context, ingress, metav1.CreateOptions{})
	}

	return ingress, err
}

func (client *Client) AwaitIngressStatus(inputIngress *networkingV1.Ingress) (*networkingV1.Ingress, error) {
	ingresses := client.clientset.NetworkingV1().Ingresses(inputIngress.Namespace)

	for {
		ingress, err := ingresses.Get(client.context, inputIngress.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if len(ingress.Status.LoadBalancer.Ingress) > 0 {
			return ingress, nil
		}

		select {
		case <-time.After(global.Settings.ServiceChecks.AwaitStatusTimeout):
			return nil, nil
		default:
			log.Printf("waiting %v seconds for ingress status\n", global.Settings.ResourceRequestBackOffPeriod)

			time.Sleep(global.Settings.ResourceRequestBackOffPeriod)
		}
	}
}

func (client *Client) RemoveIngress(namespace, name string) error {
	ingresses := client.clientset.NetworkingV1().Ingresses(namespace)

	err := ingresses.Delete(client.context, name, metav1.DeleteOptions{})

	switch {
	case errors.IsNotFound(err):
		return nil
	default:
		return err
	}
}
