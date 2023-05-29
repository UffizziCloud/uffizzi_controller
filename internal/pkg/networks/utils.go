package networks

import (
	networkingV1 "k8s.io/api/networking/v1"
)

func GetIngresEntrypoint(ingress networkingV1.IngressLoadBalancerIngress) string {
	if len(ingress.Hostname) > 0 {
		return ingress.Hostname
	}

	return ingress.IP
}
