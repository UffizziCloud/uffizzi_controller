package types

import "fmt"

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

func (volume ContainerVolume) BuildUniqName(container *Container) string {
	switch volume.Type {
	case ContainerVolumeTypeAnonymous:
		return fmt.Sprintf("anonymous-%s-%s", container.ServiceName, volume.Source)
	case ContainerVolumeTypeNamed:
		return volume.Source
	case ContainerVolumeTypeHost:
		return fmt.Sprintf("host-%s-%s", container.ServiceName, volume.Source)
	default:
		return ""
	}
}

func (volume ContainerVolume) IsHostType() bool {
	return volume.Type == ContainerVolumeTypeHost
}

func (volume ContainerVolume) IsNamedType() bool {
	return volume.Type == ContainerVolumeTypeNamed
}

func (volume ContainerVolume) IsAnonymousType() bool {
	return volume.Type == ContainerVolumeTypeAnonymous
}
