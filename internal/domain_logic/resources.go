package domain

import (
	"log"
	"regexp"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/pkg/string_utils"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
)

func (l *Logic) ApplyResource(deploymentID uint64, resource domainTypes.Resource) error {
	namespace := l.KubernetesNamespaceName(deploymentID)

	switch resource.Kind {
	case domainTypes.ResourceKindConfigMap:
		{
			err := l.ApplyResourceAsConfigMap(namespace, resource)
			if err != nil {
				return err
			}

			log.Printf("%v/configMap resource-%v was configured", namespace, resource.ID)
		}
	case domainTypes.ResourceKindSecret:
		{
			err := l.ApplyResourceAsSecret(namespace, resource)
			if err != nil {
				return err
			}

			log.Printf("%v/secret resource-%v was configured", namespace, resource.ID)
		}
	default:
		{
			log.Printf("%v/configMap resource-%v wasn't configured because of the unknown kind", namespace, resource.ID)
		}
	}

	return nil
}

func (l *Logic) ApplyResourceAsConfigMap(namespace string, resource domainTypes.Resource) error {
	name := global.Settings.ResourceName.Resource(resource.ID)

	configMap, err := l.KuberClient.FindOrInitializeConfigMap(namespace, name)
	if err != nil {
		return err
	}

	variables := map[string]string{}

	for _, variable := range resource.Variables {
		variables[variable.Name] = variable.Value
	}

	configMap.Data = variables

	if len(configMap.UID) > 0 {
		_, err = l.KuberClient.UpdateConfigMap(namespace, configMap)
		if err != nil {
			return err
		}
	} else {
		_, err = l.KuberClient.CreateConfigMap(namespace, configMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Logic) ApplyResourceAsSecret(namespace string, resource domainTypes.Resource) error {
	name := global.Settings.ResourceName.Resource(resource.ID)

	secret, err := l.KuberClient.FindOrInitializeSecret(namespace, name)
	if err != nil {
		return err
	}

	secretVariables := map[string]string{}

	for _, secretVariable := range resource.Variables {
		secretVariables[secretVariable.Name] = secretVariable.Value
	}

	secret.StringData = secretVariables

	if len(secret.UID) > 0 {
		_, err = l.KuberClient.UpdateSecret(namespace, secret)
		if err != nil {
			return err
		}
	} else {
		_, err = l.KuberClient.CreateSecret(namespace, secret)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Logic) ClearOldResources(namespace *corev1.Namespace, resources []domainTypes.Resource) error {
	err := l.ClearOldSecretResources(namespace, resources)
	if err != nil {
		return nil
	}

	err = l.ClearOldConfigMapResources(namespace, resources)
	if err != nil {
		return nil
	}

	return nil
}

func (l *Logic) ClearOldSecretResources(namespace *corev1.Namespace, resources []domainTypes.Resource) error {
	resourcesInstalled := []string{}

	for _, resource := range resources {
		resourceName := global.Settings.ResourceName.Resource(resource.ID)

		if resource.Kind == domainTypes.ResourceKindSecret {
			resourcesInstalled = append(resourcesInstalled, resourceName)
		}
	}

	pattern := global.Settings.ResourceName.Resource("")

	secrets, err := l.KuberClient.GetSecrets(namespace.Name)
	if err != nil {
		return nil
	}

	for _, secret := range secrets {
		matched, err := regexp.MatchString(pattern, secret.Name)
		if err != nil {
			return nil
		}

		if !matched || string_utils.Contains(resourcesInstalled, secret.Name) {
			continue
		}

		err = l.KuberClient.DeleteSecret(namespace.Name, secret.Name)
		if err != nil {
			return nil
		}

		log.Printf("%v/secret %v was deleted\n", namespace.Name, secret.Name)
	}

	return nil
}

func (l *Logic) ClearOldConfigMapResources(namespace *corev1.Namespace, resources []domainTypes.Resource) error {
	resourcesInstalled := []string{}

	for _, resource := range resources {
		resourceName := global.Settings.ResourceName.Resource(resource.ID)

		if resource.Kind == domainTypes.ResourceKindConfigMap {
			resourcesInstalled = append(resourcesInstalled, resourceName)
		}
	}

	pattern := global.Settings.ResourceName.Resource("")

	configMaps, err := l.KuberClient.GetConfigMaps(namespace.Name)
	if err != nil {
		return nil
	}

	for _, configMap := range configMaps {
		matched, err := regexp.MatchString(pattern, configMap.Name)
		if err != nil {
			return nil
		}

		if !matched || string_utils.Contains(resourcesInstalled, configMap.Name) {
			continue
		}

		err = l.KuberClient.DeleteConfigMap(namespace.Name, configMap.Name)
		if err != nil {
			return nil
		}

		log.Printf("%v/configMap %v was deleted\n", namespace.Name, configMap.Name)
	}

	return nil
}
