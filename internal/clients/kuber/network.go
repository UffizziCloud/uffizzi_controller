package kuber

import (
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Client) FindOrCreateNetworkPolicy(namespaceName string, name string) (*v1.NetworkPolicy, error) {
	policies := client.clientset.NetworkingV1().NetworkPolicies(namespaceName)
	policy, err := policies.Get(client.context, name, metav1.GetOptions{})

	if err != nil && !errors.IsNotFound(err) || len(policy.UID) > 0 {
		return policy, err
	}

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
							PodSelector: &metav1.LabelSelector{},
						},
						{
							IPBlock: &v1.IPBlock{
								CIDR:   "0.0.0.0/0",
								Except: []string{global.Settings.PodCidr},
							},
						},
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app.kubernetes.io/name": "ingress-nginx"},
							},
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})

	return policy, err
}
