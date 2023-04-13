package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

func TestKubernetesContainerName(t *testing.T) {
	containerName := "gadkins/counter-app"
	containerTag := "1.0"
	container := domainTypes.Container{
		Image: containerName,
		Tag:   &containerTag,
	}

	actualName, err := container.KubernetesName()
	if err != nil {
		t.Error(err)
	}

	expectedName := "gadkins-counter-app-1-0"

	assert.Equal(t, expectedName, actualName, "The two names should be the same.")
}

func TestKubernetesContainerNameWithTag(t *testing.T) {
	containerName := "uffizzi/counter-app"
	containerTag := "2.0"
	container := domainTypes.Container{
		Image: containerName,
		Tag:   &containerTag,
	}

	actualName, err := container.NameWithTag()
	if err != nil {
		t.Error(err)
	}

	expectedName := "uffizzi/counter-app:2.0"

	assert.Equal(t, expectedName, actualName, "The two names should be the same.")
}

func TestKubernetesContainerNameWithTagWhenFullImageNameExists(t *testing.T) {
	containerName := "uffizzi/counter-app"
	containerTag := "3.0"
	expectedName := "image/name:tag"
	container := domainTypes.Container{
		Image:         containerName,
		Tag:           &containerTag,
		FullImageName: expectedName,
	}

	actualName, err := container.NameWithTag()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expectedName, actualName, "The two names should be the same.")
}
