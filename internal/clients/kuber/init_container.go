package kuber

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
)

func buildInitContainerForHostVolumes(
	containerList domainTypes.ContainerList,
	composeFile domainTypes.ComposeFile,
) (corev1.Container, error) {
	hostVolumes := containerList.GetUniqHostVolumes()
	container := corev1.Container{}

	if len(hostVolumes) == 0 {
		return container, nil
	}

	if len(composeFile.RepoName) == 0 {
		return container, errors.New("host volumes supported only for compose file from github")
	}

	container = corev1.Container{
		Name:         "init-container-for-host-volumes",
		Image:        "bitnami/git",
		VolumeMounts: prepareInitContainerVolumeMounts(hostVolumes),
		Command:      prepareInitContainerCommand(hostVolumes, composeFile),
	}

	return container, nil
}

func prepareInitContainerVolumeMounts(volumes []domainTypes.DeploymentVolume) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}

	for _, volume := range volumes {
		volumeName := global.Settings.ResourceName.VolumeName(volume.UniqName)
		volumeMount := corev1.VolumeMount{
			Name:      volumeName,
			MountPath: buildMountPath(volumeName),
		}

		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumeMounts
}

func prepareInitContainerCommand(volumes []domainTypes.DeploymentVolume, composeFile domainTypes.ComposeFile) []string {
	githubUrl := fmt.Sprintf(
		"https://%v:%v@github.com/%v/%v.git",
		composeFile.RepoUsername,
		composeFile.RepoPassword,
		composeFile.RepoUsername,
		composeFile.RepoName)

	gitCloneCommand := fmt.Sprintf("git clone --branch %v --single-branch %v %v", composeFile.Branch, githubUrl, "repo")
	bashCommands := []string{gitCloneCommand, "cd repo"}

	for _, volume := range volumes {
		volumeName := global.Settings.ResourceName.VolumeName(volume.UniqName)
		sourceDir := volume.Volume.Source
		targetDir := buildMountPath(volumeName)
		copyCommand := fmt.Sprintf("cp -a %v/. %v", sourceDir, targetDir)

		bashCommands = append(bashCommands, copyCommand)
	}

	return []string{"bash", "-c", strings.Join(bashCommands, " && ")}
}

func buildMountPath(volumeName string) string {
	return fmt.Sprintf("/tmp/%v", volumeName)
}
