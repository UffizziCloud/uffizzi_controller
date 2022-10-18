package types

type ContainerHostVolumeFile struct {
	SourcePath       string `json:"source_path"`
	HostVolumeFileId uint64 `json:"host_volume_file_id"`
}
