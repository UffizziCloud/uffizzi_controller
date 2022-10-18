package kuber

import (
	"fmt"

	initContainers "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/clients/kuber/init_containers"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
)

func buildInitContainerForHostVolumes(
	containerList domainTypes.ContainerList,
	composeFile domainTypes.ComposeFile,
	hostVolumeFileList *domainTypes.HostVolumeFileList,
) (corev1.Container, error) {
	container := corev1.Container{}

	if !containerList.IsHostVolumesPresent() {
		return container, nil
	}

	if composeFile.IsGithubSourceKind() {
		container, err := initContainers.BuildGithubHostVolumeInitContainer(containerList, composeFile)

		if err != nil {
			return container, err
		}

		return container, nil
	}

	if composeFile.IsLocalSourceKind() {
		container, err := initContainers.BuildLocalHostVolumeInitContainer(containerList, composeFile, hostVolumeFileList)

		if err != nil {
			return container, err
		}

		return container, nil
	}

	return container, fmt.Errorf("unknown compose file source kind")
}
