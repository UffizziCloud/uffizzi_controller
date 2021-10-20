package types

type ConfigFileKind string

const (
	ConfigFileKindConfigMap ConfigFileKind = "config_map"
	ConfigFileKindSecret    ConfigFileKind = "secret"
)

type ConfigFile struct {
	ID       uint64         `json:"id"`
	Filename string         `json:"filename"`
	Kind     ConfigFileKind `json:"kind"`
	Payload  string         `json:"payload"`
}
