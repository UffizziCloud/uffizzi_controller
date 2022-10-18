package init_containers

import (
	"fmt"
	"strings"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
)

func BuildLocalHostVolumeInitContainer(
	containerList domainTypes.ContainerList,
	composeFile domainTypes.ComposeFile,
	hostVolumeFileList *domainTypes.HostVolumeFileList,
) (corev1.Container, error) {
	container := corev1.Container{}

	if hostVolumeFileList.IsEmpty() {
		return container, nil
	}

	hostVolumes := containerList.GetUniqHostVolumes()
	volumeMounts := prepareVolumeMountsForHostVolumes(hostVolumes)
	volumeMounts = append(volumeMounts, prepareVolumeMountsForHostVolumeFiles(hostVolumeFileList)...)
	commands, err := prepareLocalHostVolumeCommand(hostVolumes, hostVolumeFileList)

	if err != nil {
		return container, nil
	}

	container = corev1.Container{
		Name:         "init-container-for-host-volumes",
		Image:        "alpine",
		VolumeMounts: volumeMounts,
		Command:      commands,
	}

	return container, nil
}

func prepareVolumeMountsForHostVolumeFiles(hostVolumeFileList *domainTypes.HostVolumeFileList) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}

	for _, hostVolumeFile := range hostVolumeFileList.Items {
		volumeName := hostVolumeFile.VolumeName()
		volumeMount := corev1.VolumeMount{
			Name:      volumeName,
			MountPath: buildMountPathForHostVolumeFile(volumeName),
		}

		volumeMounts = append(volumeMounts, volumeMount)
	}

	return volumeMounts
}

func prepareLocalHostVolumeCommand(
	hostVolumes []domainTypes.DeploymentVolume,
	hostVolumeFileList *domainTypes.HostVolumeFileList,
) ([]string, error) {
	commands := []string{}

	for _, hostVolume := range hostVolumes {
		volumeSource := hostVolume.Volume.Source
		hostVolumeName := global.Settings.ResourceName.VolumeName(hostVolume.UniqName)
		containerHostVolumeFile, err := findContainerHostVolumeFileBySourcePath(volumeSource, hostVolume.Container)

		if err != nil {
			return commands, err
		}

		hostVolumeFile, err := hostVolumeFileList.GetHostVolumeFileById(containerHostVolumeFile.HostVolumeFileId)

		if err != nil {
			return commands, err
		}

		tarSourceDir := buildMountPathForHostVolumeFile(hostVolumeFile.VolumeName())
		tarFileName := hostVolumeFile.ConfigMapPath()
		tarPath := fmt.Sprintf("%s/%s", tarSourceDir, tarFileName)
		unpackDirPath := buildUnpackPathForHostVolumeFile(tarSourceDir)
		targetDir := buildMountPathForHostVolume(hostVolumeName)
		sourceDir := fmt.Sprintf("%s%s", unpackDirPath, hostVolumeFile.Path)

		createUnpackDirCommand := fmt.Sprintf("mkdir -p %v", unpackDirPath)
		unpackCommand := fmt.Sprintf("tar -xzf %v -C %v", tarPath, unpackDirPath)

		var copyCommand string
		if hostVolumeFile.IsFile {
			copyCommand = fmt.Sprintf("cp %v %v", sourceDir, targetDir)
		} else {
			copyCommand = fmt.Sprintf("cp -a %v/. %v", sourceDir, targetDir)
		}

		hostVolumeCommands := []string{createUnpackDirCommand, unpackCommand, copyCommand}
		commands = append(commands, strings.Join(hostVolumeCommands, " && "))
	}

	return []string{"sh", "-c", strings.Join(commands, " && ")}, nil
}

func buildMountPathForHostVolumeFile(volumeName string) string {
	return fmt.Sprintf("/tmp_host_volume_file_volumes/%v", volumeName)
}

func buildUnpackPathForHostVolumeFile(path string) string {
	return fmt.Sprintf("/tmp_unpack%v/unpack", path)
}

func findContainerHostVolumeFileBySourcePath(
	sourcePath string,
	container *domainTypes.Container,
) (*domainTypes.ContainerHostVolumeFile, error) {
	for _, containerHostVolumeFile := range container.ContainerHostVolumeFiles {
		if containerHostVolumeFile.SourcePath == sourcePath {
			return containerHostVolumeFile, nil
		}
	}

	err := fmt.Errorf("can't find containerHostVolumeFile with source_path=%v", sourcePath)

	return &domainTypes.ContainerHostVolumeFile{}, err
}
