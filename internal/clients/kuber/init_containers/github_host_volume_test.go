package init_containers

import (
	"os"
	"testing"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
)

func TestBuildCopySourceForGithubHostVolume(t *testing.T) {
	if err := global.Init(os.Getenv("ENV")); err != nil {
		panic(err)
	}

	type testItem struct {
		composeFilePath  string
		hostVolumeSource string
		expected         string
	}

	tests := []testItem{
		{
			composeFilePath:  "docker-compose.yml",
			hostVolumeSource: "./test1",
			expected:         "./test1/.",
		},
		{
			composeFilePath:  "docker-compose.yml",
			hostVolumeSource: "./",
			expected:         "./.",
		},
		{
			composeFilePath:  "subdir/docker-compose.yml",
			hostVolumeSource: "./test1",
			expected:         "./subdir/test1/.",
		},
		{
			composeFilePath:  "subdir/docker-compose.yml",
			hostVolumeSource: "./",
			expected:         "./subdir/.",
		},
		{
			composeFilePath:  "main-dir/subdir/docker-compose.yml",
			hostVolumeSource: "./",
			expected:         "./main-dir/subdir/.",
		},
		{
			composeFilePath:  "main_dir/subdir/docker-compose.yml",
			hostVolumeSource: "./test1",
			expected:         "./main_dir/subdir/test1/.",
		},
	}

	for _, value := range tests {
		buildedCopySource := buildCopySourceForGithubHostVolume(value.composeFilePath, value.hostVolumeSource)

		if buildedCopySource != value.expected {
			t.Errorf("buildCopySourceForGithubHostVolume(%v, %v) = %v; want %v",
				value.composeFilePath,
				value.hostVolumeSource,
				buildedCopySource,
				value.expected,
			)
		}
	}
}
