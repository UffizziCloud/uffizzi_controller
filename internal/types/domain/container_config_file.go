package types

type ContainerConfigFile struct {
	MountPath  string     `json:"mount_path"`
	ConfigFile ConfigFile `json:"config_file"`
}
