package types

import (
	"errors"
	"fmt"
	"unicode"
)

// Container example
type Container struct {
	ID                       uint64                     `json:"id"`
	Image                    string                     `json:"image"`
	Kind                     string                     `json:"kind"`
	Entrypoint               []string                   `json:"entrypoint,omitempty"`
	Command                  []string                   `json:"command,omitempty"`
	Tag                      *string                    `json:"tag"`
	Port                     *int32                     `json:"port"`
	TargetPort               *int32                     `json:"target_port"`
	Public                   bool                       `json:"public"`
	HasReachedRestartsLimit  bool                       `json:"has_reached_restarts_limit"`
	ReceiveIncomingRequests  bool                       `json:"receive_incoming_requests"`
	Variables                []*ContainerVariable       `json:"variables"`
	SecretVariables          []*ContainerVariable       `json:"secret_variables"`
	ControllerName           string                     `json:"controller_name"`
	MemoryLimit              uint                       `json:"memory_limit"`
	MemoryRequest            uint                       `json:"memory_request"`
	ContainerConfigFiles     []*ContainerConfigFile     `json:"container_config_files"`
	Healthcheck              *Healthcheck               `json:"healthcheck"`
	ContainerVolumes         []*ContainerVolume         `json:"volumes"`
	ServiceName              string                     `json:"service_name"`
	AdditionalSubdomains     []string                   `json:"additional_subdomains"`
	ContainerHostVolumeFiles []*ContainerHostVolumeFile `json:"container_host_volume_files"`
	Version                  string                     `json:"version"`
}

func (c Container) IsPublic() bool {
	return c.Public && c.Port != nil
}

func (c Container) IsInternal() bool {
	return c.Kind == "internal"
}

func (c Container) NameWithTag() (string, error) {
	if len(c.Image) > 0 {
		tag := "latest"
		if c.Tag != nil {
			tag = *c.Tag
		}

		return fmt.Sprintf("%v:%v", c.Image, tag), nil
	} else {
		return "", errors.New("image cannot be blank")
	}
}

func (c Container) KubernetesName() (string, error) {
	name, err := c.NameWithTag()
	if err != nil {
		return "", err
	}

	kubernetesName := ""
	unicodeRanges := []*unicode.RangeTable{unicode.Lower, unicode.Upper, unicode.Digit}
	dashRune := '-'

	if !unicode.IsOneOf(unicodeRanges, rune(name[0])) {
		return "", errors.New("container name must start and end with an alphanumeric character")
	}

	if !unicode.IsOneOf(unicodeRanges, rune(name[len(name)-1])) {
		return "", errors.New("container name must end with an alphanumeric character")
	}

	for _, symbol := range name {
		if unicode.IsOneOf(unicodeRanges, symbol) || symbol == dashRune {
			kubernetesName += string(symbol)
		} else {
			kubernetesName += "-"
		}
	}

	return kubernetesName, nil
}
