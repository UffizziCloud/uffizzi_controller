package init_containers

import (
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
)

func BuildGithubHostVolumeInitContainer(
	containerList domainTypes.ContainerList,
	composeFile domainTypes.ComposeFile,
) (corev1.Container, error) {
	hostVolumes := containerList.GetUniqHostVolumes()
	container := corev1.Container{}

	if len(hostVolumes) == 0 {
		return container, nil
	}

	container = corev1.Container{
		Name:         "init-container-for-host-volumes",
		Image:        "bitnami/git",
		VolumeMounts: prepareVolumeMountsForHostVolumes(hostVolumes),
		Command:      prepareGithubHostVolumeCommand(hostVolumes, composeFile),
	}

	return container, nil
}

func prepareVolumeMountsForHostVolumes(volumes []domainTypes.DeploymentVolume) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}

	for _, volume := range volumes {
		volumeName := global.Settings.ResourceName.VolumeName(volume.UniqName)
		volumeMount := corev1.VolumeMount{
			Name:      volumeName,
			MountPath: buildMountPathForHostVolume(volumeName),
		}

		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumeMounts
}

func prepareGithubHostVolumeCommand(
	hostVolumes []domainTypes.DeploymentVolume,
	composeFile domainTypes.ComposeFile,
) []string {
	githubUrl := fmt.Sprintf(
		"https://%v:%v@github.com/%v/%v.git",
		composeFile.RepoUsername,
		composeFile.RepoPassword,
		composeFile.RepoUsername,
		composeFile.RepoName)

	gitCloneCommand := fmt.Sprintf("git clone --branch %v --single-branch %v %v", composeFile.Branch, githubUrl, "repo")
	commands := []string{gitCloneCommand, "cd repo"}

	for _, hostVolume := range hostVolumes {
		volumeName := global.Settings.ResourceName.VolumeName(hostVolume.UniqName)
		copySource := buildCopySourceForGithubHostVolume(hostVolume.Volume.Source)
		targetDir := buildMountPathForHostVolume(volumeName)
		copyCommand := fmt.Sprintf("cp -a %v %v", copySource, targetDir)
		commands = append(commands, copyCommand)
	}

	return []string{"bash", "-c", strings.Join(commands, " && ")}
}

func buildMountPathForHostVolume(volumeName string) string {
	return fmt.Sprintf("/tmp_host_volumes/%v", volumeName)
}

func buildCopySourceForGithubHostVolume(path string) string {
	regexDoubleSlashes := regexp.MustCompile(`\/\/`)

	newPath := path + "/"
	newPath = regexDoubleSlashes.ReplaceAllString(newPath, "/") + "."

	return newPath
}
