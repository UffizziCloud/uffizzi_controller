package types

import "fmt"

type ContainerList struct {
	Items []Container `json:"items"`
}

type DeploymentVolume struct {
	Volume    *ContainerVolume
	Container *Container
	UniqName  string
}

func (list ContainerList) IsEmpty() bool {
	return list.Count() == 0
}

func (list ContainerList) Count() int {
	return len(list.Items)
}

func (list ContainerList) GetPublicContainerList() ContainerList {
	var publicContainers ContainerList

	for _, container := range list.Items {
		if container.IsPublic() {
			publicContainers.Items = append(publicContainers.Items, container)
		}
	}

	return publicContainers
}

func (list ContainerList) GetIngressContainer() *Container {
	for _, container := range list.Items {
		if container.ReceiveIncomingRequests {
			return &container
		}
	}

	return nil
}

func (list ContainerList) GetPublicContainer() (Container, error) {
	for _, container := range list.Items {
		if container.IsPublic() {
			return container, nil
		}
	}

	return Container{}, fmt.Errorf("can't find public container with port")
}

func (list ContainerList) GetUserContainerList() ContainerList {
	var userContainers ContainerList

	for _, container := range list.Items {
		if !container.IsInternal() {
			userContainers.Items = append(userContainers.Items, container)
		}
	}

	return userContainers
}

func (list *ContainerList) AddContainer(container Container) {
	list.Items = append(list.Items, container)
}

func (list ContainerList) GetUniqNamedVolumes() []DeploymentVolume {
	return getUniqVolumesByType(list, ContainerVolumeTypeNamed)
}

func (list ContainerList) GetUniqAnonymousVolumes() []DeploymentVolume {
	return getUniqVolumesByType(list, ContainerVolumeTypeAnonymous)
}

func (list ContainerList) GetUniqHostVolumes() []DeploymentVolume {
	return getUniqVolumesByType(list, ContainerVolumeTypeHost)
}

func getUniqVolumesByType(list ContainerList, volumeType ContainerVolumeType) []DeploymentVolume {
	volumes := []DeploymentVolume{}

	for i, container := range list.Items {
		for _, containerVolume := range container.ContainerVolumes {
			if containerVolume.Type != volumeType {
				continue
			}

			isVolumeExists := false
			uniqName := containerVolume.BuildUniqName(&list.Items[i])

			for _, existsVolume := range volumes {
				if existsVolume.UniqName == uniqName {
					isVolumeExists = true
					break
				}
			}

			if isVolumeExists {
				continue
			}

			volume := DeploymentVolume{
				Volume:    containerVolume,
				Container: &list.Items[i],
				UniqName:  uniqName,
			}

			volumes = append(volumes, volume)
		}
	}

	return volumes
}

func (list ContainerList) IsAnyVolumeExists() bool {
	isExists := false

	for _, container := range list.Items {
		if len(container.ContainerVolumes) > 0 {
			isExists = true
		}

		if isExists {
			continue
		}
	}

	return isExists
}
