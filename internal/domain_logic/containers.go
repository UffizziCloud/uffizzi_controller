package domain

import (
	"log"
	"time"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/pkg/string_utils"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

type ContainersUsageMetrics struct {
	ContainersMemory float64 `json:"containers_memory"`
}

func (l *Logic) GetDeploymentsContainersUsageMetrics(deploymentIDs []uint64, beginAt time.Time, endAt time.Time) (ContainersUsageMetrics, error) {
	return ContainersUsageMetrics{
		ContainersMemory: 0,
	}, nil
}

func (l *Logic) ApplyContainerSecrets(namespace string, containerList domainTypes.ContainerList) error {
	for _, container := range containerList.Items {
		name := global.Settings.ResourceName.ContainerSecret(container.ID)

		secret, err := l.KuberClient.FindOrInitializeSecret(namespace, name)
		if err != nil {
			return err
		}

		secretVariables := map[string]string{}

		for _, secretVariable := range container.SecretVariables {
			secretVariables[secretVariable.Name] = secretVariable.Value
		}

		secret.StringData = secretVariables

		if len(secret.UID) > 0 {
			if len(secret.StringData) > 0 {
				_, err = l.KuberClient.UpdateSecret(namespace, secret)
				if err != nil {
					return err
				}
			} else {
				err = l.KuberClient.DeleteSecret(namespace, secret.Name)
				if err != nil {
					return err
				}
			}
		} else {
			if len(secret.StringData) > 0 {
				_, err = l.KuberClient.CreateSecret(namespace, secret)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (l *Logic) ApplyContainersVolumes(
	namespace string,
	containerList domainTypes.ContainerList,
	hostVolumeFileList *domainTypes.HostVolumeFileList,
) error {
	composeFileVolumes := []domainTypes.DeploymentVolume{}
	composeFileVolumes = append(composeFileVolumes, containerList.GetUniqNamedVolumes()...)
	composeFileVolumes = append(composeFileVolumes, containerList.GetUniqAnonymousVolumes()...)
	composeFileVolumes = append(composeFileVolumes, containerList.GetUniqHostVolumes()...)

	for _, composeFileVolume := range composeFileVolumes {
		pvcName := global.Settings.ResourceName.PvcName(composeFileVolume.UniqName)
		pvc, err := l.KuberClient.FindOrInitializePersistentVolumeClaim(namespace, pvcName)

		if err != nil {
			return err
		}

		if len(pvc.UID) == 0 {
			_, err = l.KuberClient.CreatePersistentVolumeClaim(namespace, pvc)

			log.Printf("%v/pvc %v was created\n", namespace, pvcName)
		}

		if err != nil {
			return err
		}
	}

	for _, hostVolumeFile := range hostVolumeFileList.Items {
		err := l.ApplyHostVolumeFileAsConfigMap(namespace, hostVolumeFile)

		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Logic) RemoveUnusedContainersVolumes(namespace string, containerList domainTypes.ContainerList) error {
	uniqVolumes := containerList.GetUniqNamedVolumes()
	uniqVolumes = append(uniqVolumes, containerList.GetUniqAnonymousVolumes()...)
	uniqVolumes = append(uniqVolumes, containerList.GetUniqHostVolumes()...)
	newPersistentVolumeClaimNames := []string{}

	for _, volume := range uniqVolumes {
		newPersistentVolumeClaimNames = append(newPersistentVolumeClaimNames, global.Settings.ResourceName.PvcName(volume.UniqName))
	}

	existsPersistentVolumeClaims, err := l.KuberClient.GetPersistentVolumeClaims(namespace)
	if err != nil {
		return err
	}

	for _, existsPersistentVolumeClaim := range existsPersistentVolumeClaims {
		if string_utils.Contains(newPersistentVolumeClaimNames, existsPersistentVolumeClaim.Name) {
			continue
		}

		err := l.KuberClient.DeletePersistentVolumeClaim(namespace, existsPersistentVolumeClaim.Name)
		if err != nil {
			return nil
		}

		log.Printf("%v/pvc %v was deleted\n", namespace, existsPersistentVolumeClaim.Name)
	}

	return nil
}
