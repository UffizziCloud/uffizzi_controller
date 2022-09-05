package domain

import (
	"log"
	"regexp"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/pkg/string_utils"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
)

func (l *Logic) ApplyConfigFile(deploymentID uint64, configFile domainTypes.ConfigFile) error {
	namespace := l.KubernetesNamespaceName(deploymentID)

	switch configFile.Kind {
	case domainTypes.ConfigFileKindConfigMap:
		{
			err := l.ApplyConfigFileAsConfigMap(namespace, configFile)
			if err != nil {
				return err
			}

			log.Printf("%v/configMap config-file-%v was configured", namespace, configFile.ID)
		}
	case domainTypes.ConfigFileKindSecret:
		{
			err := l.ApplyConfigFileAsSecret(namespace, configFile)
			if err != nil {
				return err
			}

			log.Printf("%v/secret config-file-%v was configured", namespace, configFile.ID)
		}
	default:
		{
			log.Printf("%v/configMap config-file-%v wasn't configured because of the unknown kind", namespace, configFile.ID)
		}
	}

	return nil
}

func (l *Logic) ApplyConfigFileAsConfigMap(namespace string, configFile domainTypes.ConfigFile) error {
	name := global.Settings.ResourceName.ConfigFile(configFile.ID)

	configMap, err := l.KuberClient.FindOrInitializeConfigMap(namespace, name)
	if err != nil {
		return err
	}

	configMap.Data = map[string]string{
		configFile.Filename: configFile.Payload,
	}

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

func (l *Logic) ApplyConfigFileAsSecret(namespace string, configFile domainTypes.ConfigFile) error {
	name := global.Settings.ResourceName.ConfigFile(configFile.ID)

	secret, err := l.KuberClient.FindOrInitializeSecret(namespace, name)
	if err != nil {
		return err
	}

	secret.StringData = map[string]string{
		configFile.Filename: configFile.Payload,
	}

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

func (l *Logic) ClearOldConfigurationFiles(
	namespace *corev1.Namespace,
	containerList domainTypes.ContainerList,
) error {
	err := l.ClearOldSecretConfigurationFiles(namespace, containerList)
	if err != nil {
		return nil
	}

	err = l.ClearOldConfigMapConfigurationFiles(namespace, containerList)
	if err != nil {
		return nil
	}

	return nil
}

func (l *Logic) ClearOldSecretConfigurationFiles(
	namespace *corev1.Namespace,
	containerList domainTypes.ContainerList,
) error {
	configFilesInstalled := []string{}

	for _, container := range containerList.Items {
		for _, containerConfigFile := range container.ContainerConfigFiles {
			configFile := containerConfigFile.ConfigFile
			configFileName := global.Settings.ResourceName.ConfigFile(configFile.ID)

			if configFile.Kind == domainTypes.ConfigFileKindSecret {
				configFilesInstalled = append(configFilesInstalled, configFileName)
			}
		}
	}

	pattern := global.Settings.ResourceName.ConfigFile("")

	secrets, err := l.KuberClient.GetSecrets(namespace.Name)
	if err != nil {
		return nil
	}

	for _, secret := range secrets {
		matched, err := regexp.MatchString(pattern, secret.Name)
		if err != nil {
			return nil
		}

		if !matched || string_utils.Contains(configFilesInstalled, secret.Name) {
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

func (l *Logic) ClearOldConfigMapConfigurationFiles(
	namespace *corev1.Namespace,
	containerList domainTypes.ContainerList,
) error {
	configFilesInstalled := []string{}

	for _, container := range containerList.Items {
		for _, containerConfigFile := range container.ContainerConfigFiles {
			configFile := containerConfigFile.ConfigFile
			configFileName := global.Settings.ResourceName.ConfigFile(configFile.ID)

			if configFile.Kind == domainTypes.ConfigFileKindConfigMap {
				configFilesInstalled = append(configFilesInstalled, configFileName)
			}
		}
	}

	pattern := global.Settings.ResourceName.ConfigFile("")

	configMaps, err := l.KuberClient.GetConfigMaps(namespace.Name)
	if err != nil {
		return nil
	}

	for _, configMap := range configMaps {
		matched, err := regexp.MatchString(pattern, configMap.Name)
		if err != nil {
			return nil
		}

		if !matched || string_utils.Contains(configFilesInstalled, configMap.Name) {
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

func (l *Logic) ApplyHostVolumeFileAsConfigMap(namespace string, hostVolumeFile domainTypes.HostVolumeFile) error {
	name := hostVolumeFile.ConfigMapName()

	configMap, err := l.KuberClient.FindOrInitializeConfigMap(namespace, name)
	if err != nil {
		return err
	}

	binaryPayload, err := hostVolumeFile.BinaryPayload()

	if err != nil {
		return err
	}

	configMap.BinaryData = map[string][]byte{
		hostVolumeFile.ConfigMapKey(): binaryPayload,
	}

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
