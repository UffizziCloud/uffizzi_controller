package kuber

import (
	"log"
	"time"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (client *Client) GetServices(namespace string) (*corev1.ServiceList, error) {
	services := client.clientset.CoreV1().Services(namespace)

	serviceList, err := services.List(client.context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return serviceList, nil
}

func (client *Client) GetServicesByLabel(namespace string, labelSelector string) (*corev1.ServiceList, error) {
	services := client.clientset.CoreV1().Services(namespace)

	serviceList, err := services.List(client.context, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}

	return serviceList, nil
}

func (client *Client) GetService(namespace, name string) (*corev1.Service, error) {
	services := client.clientset.CoreV1().Services(namespace)

	service, err := services.Get(client.context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (client *Client) FindSingleService(namespace *corev1.Namespace) (*corev1.Service, error) {
	serviceName, err := client.GetOrPrepareServiceName(namespace)
	if err != nil {
		return nil, err
	}

	return client.GetService(namespace.Name, serviceName)
}

func (client *Client) findOrInitializeService(namespace, serviceName string) (*corev1.Service, error) {
	services := client.clientset.CoreV1().Services(namespace)

	service, err := services.Get(client.context, serviceName, metav1.GetOptions{})

	if err != nil && !errors.IsNotFound(err) {
		return service, err
	}

	if len(service.UID) == 0 {
		service = initializeService(namespace, serviceName)
	}

	return service, nil
}

func (client *Client) updateServiceAttributes(service *corev1.Service, deploymentSelectorName string, containerList domainTypes.ContainerList) *corev1.Service {
	var servicePorts []corev1.ServicePort

	for _, draftContainer := range containerList.Items {
		port := *draftContainer.Port

		targetPort := *draftContainer.Port
		if draftContainer.TargetPort != nil {
			targetPort = *draftContainer.TargetPort
		}

		servicePort := corev1.ServicePort{}

		for _, existingServicePort := range service.Spec.Ports {
			if existingServicePort.Name == draftContainer.ControllerName && existingServicePort.Port == port {
				servicePort = *existingServicePort.DeepCopy()
				break
			}
		}

		if (corev1.ServicePort{}) == servicePort {
			servicePort = corev1.ServicePort{Name: draftContainer.ControllerName, Port: port, TargetPort: intstr.FromInt(int(targetPort))}
		}

		servicePorts = append(servicePorts, servicePort)
	}

	service.Spec.Ports = servicePorts
	service.Spec.Selector = map[string]string{"app": deploymentSelectorName}

	return service
}

func (client *Client) GetOrPrepareServiceName(namespace *corev1.Namespace) (string, error) {
	var annotationKey = "serviceName"

	serviceName := namespace.Annotations[annotationKey]

	if len(serviceName) > 0 {
		return serviceName, nil
	}

	serviceName = generateServiceName()

	_, err := client.UpdateAnnotationNamespace(namespace.Name, annotationKey, serviceName)

	return serviceName, err
}

func (client *Client) CreateOrUpdateService(namespace *corev1.Namespace, deploymentSelectorName string, containerList domainTypes.ContainerList) (*corev1.Service, error) {
	serviceName, err := client.GetOrPrepareServiceName(namespace)
	if err != nil {
		return nil, err
	}

	services := client.clientset.CoreV1().Services(namespace.Name)

	service, err := client.findOrInitializeService(namespace.Name, serviceName)
	if err != nil {
		return service, err
	}

	service = client.updateServiceAttributes(service, deploymentSelectorName, containerList)

	if len(service.UID) > 0 {
		service, err = services.Update(client.context, service, metav1.UpdateOptions{})
	} else {
		service, err = services.Create(client.context, service, metav1.CreateOptions{})
	}

	return service, err
}

func (client *Client) AwaitServiceStatus(inputService *corev1.Service) (*corev1.Service, error) {
	services := client.clientset.CoreV1().Services(inputService.Namespace)

	for {
		service, err := services.Get(client.context, inputService.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if len(service.Spec.ClusterIP) > 0 {
			return service, nil
		}

		select {
		case <-time.After(global.Settings.ServiceChecks.AwaitStatusTimeout):
			return nil, nil
		default:
			log.Printf("waiting %v seconds for service status\n", global.Settings.ResourceRequestBackOffPeriod)

			time.Sleep(global.Settings.ResourceRequestBackOffPeriod)
		}
	}
}

func (client *Client) RemoveService(namespace, name string) error {
	services := client.clientset.CoreV1().Services(namespace)

	err := services.Delete(client.context, name, metav1.DeleteOptions{})

	switch {
	case errors.IsNotFound(err):
		return nil
	default:
		return err
	}
}
