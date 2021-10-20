package resource_name_utils

import "fmt"

type ResouceNameUtils struct{}

func (resouceNameUtils *ResouceNameUtils) Resource(id interface{}) string {
	return fmt.Sprintf("resource-%v", id)
}

func (resouceNameUtils *ResouceNameUtils) ConfigFile(id interface{}) string {
	return fmt.Sprintf("config-file-%v", id)
}

func (resouceNameUtils *ResouceNameUtils) ContainerVolume(containerID, configFileID interface{}) string {
	return fmt.Sprintf("container-%v-config-file-%v", containerID, configFileID)
}

func (resouceNameUtils *ResouceNameUtils) Credential(credentialId uint64) string {
	return fmt.Sprintf("credential-%v", credentialId)
}
