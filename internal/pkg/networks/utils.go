package networks

import (
	corev1 "k8s.io/api/core/v1"
)

func GetIngresEntrypoint(ingress corev1.LoadBalancerIngress) string {
	if len(ingress.Hostname) > 0 {
		return ingress.Hostname
	}

	return ingress.IP
}
