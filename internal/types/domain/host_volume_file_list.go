package types

import "fmt"

type HostVolumeFileList struct {
	Items []HostVolumeFile `json:"items"`
}

func (list HostVolumeFileList) IsEmpty() bool {
	return list.Count() == 0
}

func (list HostVolumeFileList) Count() int {
	return len(list.Items)
}

func (list HostVolumeFileList) GetHostVolumeFileById(id uint64) (HostVolumeFile, error) {
	hostVolumeFile := HostVolumeFile{}

	for _, hostVolumeFile := range list.Items {
		if hostVolumeFile.ID == id {
			return hostVolumeFile, nil
		}
	}

	return hostVolumeFile, fmt.Errorf("can't find hostVolumeFile with id=%v", id)
}
