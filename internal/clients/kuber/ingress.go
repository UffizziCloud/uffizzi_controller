package kuber

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const NAMESPACE_ANNOTATION_KEY = "ingressName"
const BASIC_AUTH_SECRET_NAME = "basic-auth"

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
	ingressName := client.GetIngressName(namespace)

	if len(ingressName) > 0 {
		return ingressName, nil
	}

	ingressName = generateIngressName()

	_, err := client.UpdateAnnotationNamespace(namespace.Name, NAMESPACE_ANNOTATION_KEY, ingressName)

	return ingressName, err
}

func (client *Client) GetIngressName(namespace *corev1.Namespace) string {
	ingressName := namespace.Annotations[NAMESPACE_ANNOTATION_KEY]

	return ingressName
}

func (client *Client) UpdateIngressAttributes(
	ingress *networkingV1.Ingress,
	namespace *corev1.Namespace,
	container domainTypes.Container,
	serviceName, deploymentHost string,
	project domainTypes.Project) (*networkingV1.Ingress, error) {
	containerPort := *container.Port
	tls := []networkingV1.IngressTLS{
		{Hosts: []string{deploymentHost}},
	}

	var err error

	ingress.ObjectMeta.Annotations["kubernetes.io/ingress.class"] = "nginx"

	if project.IsPreviewsProtected {
		ingress, err = client.AddBasicAuthToIngress(ingress, project, namespace.Name)
	}

	if err != nil {
		return nil, err
	}

	// if a `ClusterIssuer` is specified, use it.
	if len(global.Settings.CertManagerClusterIssuer) > 0 {
		ingress.ObjectMeta.Annotations["cert-manager.io/cluster-issuer"] = global.Settings.CertManagerClusterIssuer
		tls = []networkingV1.IngressTLS{
			{Hosts: []string{deploymentHost}, SecretName: deploymentHost},
		}
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

	return ingress, nil
}

func (client *Client) CreateOrUpdateIngress(namespace *corev1.Namespace,
	serviceName string,
	container domainTypes.Container,
	deploymentHost string,
	project domainTypes.Project) (*networkingV1.Ingress, error) {
	ingressName, err := client.GetOrPrepareIngressName(namespace)
	if err != nil {
		return nil, err
	}

	ingresses := client.clientset.NetworkingV1().Ingresses(namespace.Name)

	ingress, err := client.findOrInitializeIngress(namespace.Name, ingressName)
	if err != nil {
		return ingress, err
	}

	ingress, err = client.UpdateIngressAttributes(ingress, namespace, container, serviceName, deploymentHost, project)

	if len(ingress.UID) > 0 {
		ingress, err = ingresses.Update(client.context, ingress, metav1.UpdateOptions{})
	} else {
		ingress, err = ingresses.Create(client.context, ingress, metav1.CreateOptions{})
	}

	return ingress, err
}

func (client *Client) AwaitIngressStatus(inputIngress *networkingV1.Ingress) (*networkingV1.Ingress, error) {
	namespace := inputIngress.Namespace
	ctx := context.TODO()
	ctx, cancelContext := context.WithCancel(ctx)
	ingressChan := make(chan *networkingV1.Ingress, 1)
	errorChan := make(chan error, 1)

	defer cancelContext()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("namespace/%s AwaitIngressStatus goroutine was stopped\n", namespace)
				return
			default:
				ingress, err := client.clientset.NetworkingV1().Ingresses(namespace).Get(
					client.context,
					inputIngress.Name,
					metav1.GetOptions{},
				)
				if err != nil {
					errorChan <- err
					return
				}

				if len(ingress.Status.LoadBalancer.Ingress) > 0 {
					ingressChan <- ingress
					return
				}

				resourceRequestBackOffPeriod := global.Settings.ResourceRequestBackOffPeriod

				log.Printf("namespace/%s waiting %v seconds for ingress status\n", namespace, resourceRequestBackOffPeriod)

				time.Sleep(resourceRequestBackOffPeriod)
			}
		}
	}()

	select {
	case <-time.After(global.Settings.ServiceChecks.AwaitStatusTimeout):
		return nil, fmt.Errorf("ingress check timeout")
	case ingress := <-ingressChan:
		return ingress, nil
	case err := <-errorChan:
		return nil, err
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

func (client *Client) AddBasicAuthToIngress(
	ingress *networkingV1.Ingress,
	project domainTypes.Project,
	namespaceName string) (*networkingV1.Ingress, error) {

	err := client.DeleteSecret(namespaceName, BASIC_AUTH_SECRET_NAME)

	if err != nil {
		return nil, err
	}

	secretDraft := client.BuildSecretBasicAuth(project.PreviewsUserName, project.PreviewsPassword, namespaceName, BASIC_AUTH_SECRET_NAME)
	_, err = client.CreateSecret(namespaceName, secretDraft)

	if err != nil {
		return nil, err
	}

	basicAuthAnnotations := getBasicAuthAnnotation(BASIC_AUTH_SECRET_NAME)

	for key, val := range basicAuthAnnotations {
		ingress.ObjectMeta.Annotations[key] = val
	}

	return ingress, nil
}

func (client *Client) DeleteBasicAuthFromIngress(
	ingress *networkingV1.Ingress,
	namespaceName string) (*networkingV1.Ingress, error) {

	err := client.DeleteSecret(namespaceName, BASIC_AUTH_SECRET_NAME)
	if err != nil {
		return nil, err
	}

	basicAuthAnnotationKeys := getBasicAuthAnnotation(BASIC_AUTH_SECRET_NAME)

	for key, _ := range basicAuthAnnotationKeys {
		delete(ingress.ObjectMeta.Annotations, key)
	}

	return ingress, nil
}

func (client *Client) UpdateIngress(ingress *networkingV1.Ingress, namespaceName string) (*networkingV1.Ingress, error) {

	ingresses := client.clientset.NetworkingV1().Ingresses(namespaceName)
	ingress, err := ingresses.Update(client.context, ingress, metav1.UpdateOptions{})

	return ingress, err
}

func getBasicAuthAnnotation(secretName string) map[string]string {
	return map[string]string{
		"nginx.ingress.kubernetes.io/auth-type":   "basic",
		"nginx.ingress.kubernetes.io/auth-secret": BASIC_AUTH_SECRET_NAME,
		"nginx.ingress.kubernetes.io/auth-realm":  "Authentication Required",
	}
}
