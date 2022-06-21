package kuber

import (
	"fmt"
	"log"
	"strings"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
)

func prepareContainerEnvironmentVariables(container domainTypes.Container) []corev1.EnvVar {
	var variables []corev1.EnvVar

	for _, variable := range container.Variables {
		newVariable := corev1.EnvVar{Name: variable.Name, Value: variable.Value}
		variables = append(variables, newVariable)
	}

	return variables
}

func prepareContainerSecrets(container domainTypes.Container, secret *corev1.Secret) []corev1.EnvVar {
	envVariables := []corev1.EnvVar{}

	if secret == nil {
		log.Printf("No Secret Variables (containerId=%d)", container.ID)
		return envVariables
	}

	name := global.Settings.ResourceName.ContainerSecret(container.ID)
	log.Printf("Container Secret Name - %s", name)

	for envVariableName := range secret.Data {
		envVariable := corev1.EnvVar{
			Name: envVariableName,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: envVariableName,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name,
					},
				},
			},
		}

		envVariables = append(envVariables, envVariable)
	}

	return envVariables
}

func prepareDeploymentVolumes(containerList domainTypes.ContainerList) []corev1.Volume {
	volumes := []corev1.Volume{}
	configFileVolumes := prepareDeploymentConfigFileVolumes(containerList)
	volumes = append(volumes, configFileVolumes...)
	pvcVolumes := prepareDeploymentPvcVolumes(containerList)
	volumes = append(volumes, pvcVolumes...)

	return volumes
}

func prepareContainerVolumeMounts(container domainTypes.Container) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}
	configVolumeMounts := prepareContainerConfigFileVolumeMounts(container)
	volumeMounts = append(volumeMounts, configVolumeMounts...)
	namedVolumeMounts := prepareContainerNamedVolumeMounts(container)
	volumeMounts = append(volumeMounts, namedVolumeMounts...)

	return volumeMounts
}

func prepareConfigFileMountPath(containerConfigFile *domainTypes.ContainerConfigFile) (string, string) {
	mountPath := containerConfigFile.MountPath

	if len(mountPath) == 0 {
		mountPath = "/"
	}

	if mountPath[0] != '/' {
		mountPath = fmt.Sprintf("/%v", containerConfigFile.MountPath)
	}

	if mountPath == "/" {
		mountPath = fmt.Sprintf("/%v", containerConfigFile.ConfigFile.Filename)
	}

	mountPathParts := strings.Split(mountPath, "/")
	mountFileName := mountPathParts[len(mountPathParts)-1]

	return mountPath, mountFileName
}

func prepareCredentialsDeployment(credentials []domainTypes.Credential) []corev1.LocalObjectReference {
	references := []corev1.LocalObjectReference{}

	for _, credential := range credentials {
		credentialName := global.Settings.ResourceName.Credential(credential.ID)

		references = append(references, corev1.LocalObjectReference{Name: credentialName})
	}

	return references
}

func prepareDeploymentConfigFileVolumes(containerList domainTypes.ContainerList) []corev1.Volume {
	volumes := []corev1.Volume{}

	for _, container := range containerList.Items {
		for _, containerConfigFile := range container.ContainerConfigFiles {
			configFile := containerConfigFile.ConfigFile

			configFileName := global.Settings.ResourceName.ConfigFile(configFile.ID)

			volumeName := global.Settings.ResourceName.ContainerVolume(container.ID, configFile.ID)

			volume := corev1.Volume{Name: volumeName}

			_, mountFileName := prepareConfigFileMountPath(containerConfigFile)

			items := []corev1.KeyToPath{
				{Key: configFile.Filename, Path: mountFileName},
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

func prepareDeploymentPvcVolumes(containerList domainTypes.ContainerList) []corev1.Volume {
	volumes := []corev1.Volume{}

	for _, namedVolume := range containerList.GetUniqNamedVolumes() {
		volume := corev1.Volume{
			Name: global.Settings.ResourceName.VolumeName(namedVolume.Source),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: global.Settings.ResourceName.PvcName(namedVolume.Source),
				},
			},
		}

		volumes = append(volumes, volume)
	}

	return volumes
}

func prepareContainerConfigFileVolumeMounts(container domainTypes.Container) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}

	for _, containerConfigFile := range container.ContainerConfigFiles {
		configFile := containerConfigFile.ConfigFile

		volumeName := global.Settings.ResourceName.ContainerVolume(container.ID, configFile.ID)

		mountPath, mountFileName := prepareConfigFileMountPath(containerConfigFile)

		volumeMount := corev1.VolumeMount{
			Name:      volumeName,
			MountPath: mountPath,
			SubPath:   mountFileName,
			ReadOnly:  true,
		}

		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumeMounts
}

func prepareContainerNamedVolumeMounts(container domainTypes.Container) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}

	for _, containerVolume := range container.ContainerVolumes {
		if containerVolume.Type != domainTypes.ContainerVolumeTypeNamed {
			continue
		}

		volumeMount := corev1.VolumeMount{
			Name:      global.Settings.ResourceName.VolumeName(containerVolume.Source),
			MountPath: containerVolume.Target,
			ReadOnly:  containerVolume.ReadOnly,
		}

		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumeMounts
}

func prepareContainerHealthcheck(container domainTypes.Container) *corev1.Probe {
	if container.Healthcheck == nil {
		return nil
	}

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
