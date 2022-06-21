package resource_name_utils

import (
	"fmt"
	"strings"
)

type ResouceNameUtils struct{}

func (resouceNameUtils *ResouceNameUtils) ConfigFile(id interface{}) string {
	return fmt.Sprintf("config-file-%v", id)
}

func (resouceNameUtils *ResouceNameUtils) ContainerVolume(containerID, configFileID interface{}) string {
	return fmt.Sprintf("container-%v-config-file-%v", containerID, configFileID)
}

func (resouceNameUtils *ResouceNameUtils) Credential(credentialId uint64) string {
	return fmt.Sprintf("credential-%v", credentialId)
}

func (resouceNameUtils *ResouceNameUtils) ContainerSecret(containerID uint64) string {
	return fmt.Sprintf("container-%v-secret", containerID)
}

func (resouceNameUtils *ResouceNameUtils) Deployment(namespace string) string {
	return fmt.Sprintf("app-%v", namespace)
}

func (resouceNameUtils *ResouceNameUtils) Policy(namespace string) string {
	return fmt.Sprintf("app-%v", namespace)
}

func (resouceNameUtils *ResouceNameUtils) PvcName(name string) string {
	rfcName := strings.ReplaceAll(name, "_", "-")
	return fmt.Sprintf("pvc-%v", rfcName)
}

func (resouceNameUtils *ResouceNameUtils) VolumeName(name string) string {
	rfcName := strings.ReplaceAll(name, "_", "-")
	return fmt.Sprintf("volume-%v", rfcName)
}
