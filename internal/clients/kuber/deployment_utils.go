package kuber

import (
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
)

func prepareEnvFromResources(resources []domainTypes.Resource) []corev1.EnvFromSource {
	environmentsFrom := []corev1.EnvFromSource{}

	for _, resource := range resources {
		envFrom := prepareEnvFromResource(resource)

		environmentsFrom = append(environmentsFrom, *envFrom)
	}

	return environmentsFrom
}

func prepareEnvFromResource(resource domainTypes.Resource) *corev1.EnvFromSource {
	resourceName := global.Settings.ResourceName.Resource(resource.ID)

	if resource.Kind == domainTypes.ResourceKindConfigMap {
		return &corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: resourceName,
				},
			},
		}
	}

	if resource.Kind == domainTypes.ResourceKindSecret {
		return &corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: resourceName,
				},
			},
		}
	}

	return nil
}

func prepareDeploymentVolumes(containerList domainTypes.ContainerList) []corev1.Volume {
	volumes := []corev1.Volume{}

	for _, container := range containerList.Items {
		for _, containerConfigFile := range container.ContainerConfigFiles {
			configFile := containerConfigFile.ConfigFile

			configFileName := global.Settings.ResourceName.ConfigFile(configFile.ID)

			volumeName := global.Settings.ResourceName.ContainerVolume(container.ID, configFile.ID)

			volume := corev1.Volume{Name: volumeName}

			items := []corev1.KeyToPath{
				{Key: configFile.Filename, Path: configFile.Filename},
			}

			if configFile.Kind == domainTypes.ConfigFileKindConfigMap {
				volume.ConfigMap = &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configFileName,
					},
					Items: items,
				}
			}

			if configFile.Kind == domainTypes.ConfigFileKindSecret {
				volume.Secret = &corev1.SecretVolumeSource{
					SecretName: configFileName,
					Items:      items,
				}
			}

			volumes = append(volumes, volume)
		}
	}

	return volumes
}

func prepareContainerVolumeMounts(container domainTypes.Container) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}

	for _, containerConfigFile := range container.ContainerConfigFiles {
		configFile := containerConfigFile.ConfigFile

		volumeName := global.Settings.ResourceName.ContainerVolume(container.ID, configFile.ID)

		volumeMount := corev1.VolumeMount{
			Name:      volumeName,
			MountPath: containerConfigFile.MountPath,
			ReadOnly:  true,
		}

		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumeMounts
}

func prepareCredentialsDeployment(credentials []domainTypes.Credential) []corev1.LocalObjectReference {
	references := []corev1.LocalObjectReference{}

	for _, credential := range credentials {
		credentialName := global.Settings.ResourceName.Credential(credential.ID)

		references = append(references, corev1.LocalObjectReference{Name: credentialName})
	}

	return references
}

func prepareContainerHealthcheck(container domainTypes.Container) *corev1.Probe {
	healthcheck := *container.Healthcheck

	if len(healthcheck.Test) == 0 || healthcheck.Disable {
		return nil
	}

	probe := &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: healthcheck.Test,
			},
		},
		InitialDelaySeconds: healthcheck.StartPeriod,
		TimeoutSeconds:      healthcheck.Timeout,
		PeriodSeconds:       healthcheck.Interval,
		FailureThreshold:    healthcheck.Retries,
	}

	return probe
}
