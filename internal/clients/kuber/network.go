package kuber

import (
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (client *Client) FindOrCreateNetworkPolicy( // nolint:funlen
	namespaceName string,
	name string,
) (*v1.NetworkPolicy, error) {
	policies := client.clientset.NetworkingV1().NetworkPolicies(namespaceName)
	policy, err := policies.Get(client.context, name, metav1.GetOptions{})

	if err != nil && !errors.IsNotFound(err) || len(policy.UID) > 0 {
		return policy, err
	}

	// Define variables from contstants so we can make pointers later.
	protocolUDP := corev1.ProtocolUDP
	portDNS := intstr.FromInt(53) // nolint: gomnd

	policy, err = policies.Create(client.context, &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
			},
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{},
			},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					From: []v1.NetworkPolicyPeer{
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"kubernetes.io/metadata.name": "uffizzi-controller"},
							},
						},
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"kubernetes.io/metadata.name": global.Settings.KubernetesNamespace},
							},
						},
					},
				},
			},
			Egress: []v1.NetworkPolicyEgressRule{
				{
					To: []v1.NetworkPolicyPeer{
						{
							IPBlock: &v1.IPBlock{
								CIDR: "0.0.0.0/0",
								Except: []string{
									"10.0.0.0/8",
									"172.16.0.0/12",
									"192.168.0.0/16",
								},
							},
						},
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"kubernetes.io/metadata.name": "kube-system"},
							},
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"k8s-app": "kube-dns"},
							},
						},
					},
					Ports: []v1.NetworkPolicyPort{
						{
							Protocol: &protocolUDP,
							Port:     &portDNS,
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})

	return policy, err
}
