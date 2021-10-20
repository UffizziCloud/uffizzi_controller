package kuber

import (
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Client) FindNamespace(name string) (*corev1.Namespace, error) {
	namespaces := client.clientset.CoreV1().Namespaces()

	namespace, err := namespaces.Get(client.context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (client *Client) CreateNamespace(name, kind string) (*corev1.Namespace, error) {
	namespaces := client.clientset.CoreV1().Namespaces()

	draftNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"name":                         name,
				"kind":                         kind,
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
			},
			Annotations: map[string]string{},
		},
	}

	namespace, err := namespaces.Create(client.context, draftNamespace, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (client *Client) UpdateNamespace(namespace *corev1.Namespace) (*corev1.Namespace, error) {
	namespaces := client.clientset.CoreV1().Namespaces()

	namespace, err := namespaces.Update(client.context, namespace, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (client *Client) UpdateAnnotationNamespace(namespaceName, name, value string) (*corev1.Namespace, error) {
	namespace, err := client.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	namespace.Annotations[name] = value

	namespace, err = client.UpdateNamespace(namespace)
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (client *Client) RemoveNamespace(name string) error {
	namespaces := client.clientset.CoreV1().Namespaces()
	err := namespaces.Delete(client.context, name, metav1.DeleteOptions{})

	return err
}
