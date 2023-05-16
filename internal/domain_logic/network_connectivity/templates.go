package network_connectivity

import (
	"fmt"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

const StatusPending ConnectivityStatus = "pending"
const StatusSuccess ConnectivityStatus = "success"
const StatusSuccessTcp ConnectivityStatus = "success_tcp"
const StatusFailed ConnectivityStatus = "failed"

type ConnectivityStatus string

type ConnectivityContainerStatus struct {
	Entrypoint string             `json:"ip,omitempty"`
	DomainName string             `json:"domain_name,omitempty"`
	Port       int                `json:"port"`
	Status     ConnectivityStatus `json:"status"`
}

type ConnectivityContainer struct {
	Service     *ConnectivityContainerStatus `json:"service,omitempty"`
	Ingress     *ConnectivityContainerStatus `json:"ingress,omitempty"`
	IngressHttp *ConnectivityContainerStatus `json:"ingress_http,omitempty"`
}

func (c *ConnectivityContainer) SetIngressStatus(status ConnectivityStatus, entrypoint string) {
	c.Ingress.Status = status
	c.Ingress.Entrypoint = entrypoint
}

func (c *ConnectivityContainer) SetIngressHttpStatus(status ConnectivityStatus, entrypoint string) {
	c.IngressHttp.Status = status
	c.IngressHttp.Entrypoint = entrypoint
}

func (c *ConnectivityContainer) SetLoadBalancerStatus(status ConnectivityStatus, entrypoint string) {
	c.Service.Status = status
	c.Service.Entrypoint = entrypoint
}

type ConnectivityResponse struct {
	Containers map[string]ConnectivityContainer `json:"containers"`
}

func NewNetworkConnectivityTemplate(containerList domainTypes.ContainerList) (*ConnectivityResponse, error) {
	networkConnectivity := &ConnectivityResponse{}
	networkConnectivity.Containers = make(map[string]ConnectivityContainer, containerList.Count())

	for _, container := range containerList.Items {
		containerPort := int(*container.Port)

		domainName, err := container.KubernetesName()
		if err != nil {
			return nil, err
		}

		networkConnectivityContainer := ConnectivityContainer{
			Service: &ConnectivityContainerStatus{
				DomainName: domainName,
				Port:       containerPort,
				Status:     StatusPending,
			},
		}

		networkConnectivity.Containers[fmt.Sprint(container.ID)] = networkConnectivityContainer
	}

	return networkConnectivity, nil
}

func (response *ConnectivityResponse) AddIngressContainer(container *domainTypes.Container) {
	containerID := fmt.Sprint(container.ID)
	networkConnectivityContainer := response.Containers[containerID]
	networkConnectivityContainer.Ingress = &ConnectivityContainerStatus{
		Port:   global.Settings.IngressDefaultPort,
		Status: StatusPending,
	}

	response.Containers[containerID] = networkConnectivityContainer
}

func (response *ConnectivityResponse) AddIngressHttpStatus(container *domainTypes.Container) {
	containerID := fmt.Sprint(container.ID)
	networkConnectivityContainer := response.Containers[containerID]
	networkConnectivityContainer.IngressHttp = &ConnectivityContainerStatus{
		Status: StatusPending,
	}

	response.Containers[containerID] = networkConnectivityContainer
}

func (response *ConnectivityResponse) SetIngressStatus(containerID string,
	status ConnectivityStatus, entrypoint string) {
	container := response.Containers[containerID]
	container.SetIngressStatus(status, entrypoint)
	response.Containers[containerID] = container
}

func (response *ConnectivityResponse) SetLoadBalancerStatus(containerID string,
	status ConnectivityStatus, entrypoint string) {
	container := response.Containers[containerID]
	container.SetLoadBalancerStatus(status, entrypoint)
	response.Containers[containerID] = container
}
