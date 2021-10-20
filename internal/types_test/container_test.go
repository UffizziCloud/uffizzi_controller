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

	assert.Equal(t, actualName, expectedName, "The two names should be the same.")
}
