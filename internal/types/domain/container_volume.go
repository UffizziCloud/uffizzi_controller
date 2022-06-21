package types

type ContainerVolume struct {
	Source   string              `json:"source"`
	Target   string              `json:"target"`
	Type     ContainerVolumeType `json:"type"`
	ReadOnly bool                `json:"read_only"`
}

type ContainerVolumeType string

const (
	ContainerVolumeTypeNamed     ContainerVolumeType = "named"
	ContainerVolumeTypeAnonymous ContainerVolumeType = "anonymous"
	ContainerVolumeTypeHost      ContainerVolumeType = "host"
)

// func (volume ContainerVolume) IsHostTypeContainerVolume() bool {
// 	return volume.Type == "host"
// }

// func (volume ContainerVolume) IsNamedTypeContainerVolume() bool {
// 	return volume.Type == "named"
// }

// func (volume ContainerVolume) IsAnonymousTypeContainerVolume() bool {
// 	return volume.Type == "anonymous"
// }
