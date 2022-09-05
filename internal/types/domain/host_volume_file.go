package types

import (
	"encoding/base64"
	"fmt"
)

type HostVolumeFile struct {
	ID      uint64 `json:"id"`
	Source  string `json:"source"`
	Path    string `json:"path"`
	Payload string `json:"payload"`
	IsFile  bool   `json:"is_file"`
}

func (h HostVolumeFile) BinaryPayload() ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(h.Payload)

	return data, err
}

func (h HostVolumeFile) ConfigMapName() string {
	return fmt.Sprintf("host-volume-file-configmap-%v", h.ID)
}

func (h HostVolumeFile) VolumeName() string {
	return fmt.Sprintf("host-volume-file-volume-%v", h.ID)
}

func (h HostVolumeFile) ConfigMapKey() string {
	return "tar_file"
}

func (h HostVolumeFile) ConfigMapPath() string {
	return h.ConfigMapKey()
}
