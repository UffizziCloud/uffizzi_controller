package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

func TestContainerList_GetPublicContainerList(t *testing.T) {
	var port int32 = 80

	containerName := "gadkins/counter-app" // nice ref
	containerTag := "1.0"
	publicContainer := domainTypes.Container{Image: containerName, Tag: &containerTag, Public: true, Port: &port}
	privateContainer := domainTypes.Container{Image: containerName, Tag: &containerTag, Public: false}
	containersCount := 1

	containerList := MakeContainerList(containersCount, publicContainer)
	containerList.AddContainer(privateContainer)

	publicContainerList := containerList.GetPublicContainerList()
	expectedPublicContainerList := MakeContainerList(containersCount, publicContainer)

	assert.Equal(
		t,
		publicContainerList,
		expectedPublicContainerList,
		"The two containerLists should be the same.",
	)
}

func MakeContainerList(count int, container domainTypes.Container) domainTypes.ContainerList {
	containerList := domainTypes.ContainerList{}

	for i := 0; i < count; i++ {
		containerList.AddContainer(container)
	}

	return containerList
}
