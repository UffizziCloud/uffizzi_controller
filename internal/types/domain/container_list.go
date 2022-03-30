package types

import "fmt"

type ContainerList struct {
	Items []Container `json:"items"`
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
