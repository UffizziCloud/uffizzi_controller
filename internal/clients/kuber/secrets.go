package kuber

import (
	"fmt"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client Client) FindOrInitializeSecret(namespace, name string) (*corev1.Secret, error) {
	secret, err := client.GetSecret(namespace, name)
	if err != nil && !errors.IsNotFound(err) {
		return secret, err
	}

	if secret != nil && len(secret.UID) > 0 {
		return secret, nil
	}

	secretDraft := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
			},
		},
		StringData: map[string]string{},
	}

	return secretDraft, nil
}

func (client *Client) GetSecret(namespace, name string) (*corev1.Secret, error) {
	secretClient := client.clientset.CoreV1().Secrets(namespace)

	secret, err := secretClient.Get(client.context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (client *Client) GetSecrets(namespace string) ([]corev1.Secret, error) {
	secretClient := client.clientset.CoreV1().Secrets(namespace)

	secrets, err := secretClient.List(client.context, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/managed-by=%v", global.Settings.ManagedApplication),
	})
	if err != nil {
		return nil, err
	}

	return secrets.Items, nil
}

func (client *Client) DeleteSecret(namespace, name string) error {
	secretClient := client.clientset.CoreV1().Secrets(namespace)

	err := secretClient.Delete(client.context, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) CreateSecret(namespace string, secretDraft *corev1.Secret) (*corev1.Secret, error) {
	secretClient := client.clientset.CoreV1().Secrets(namespace)

	secret, err := secretClient.Create(client.context, secretDraft, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (client *Client) UpdateSecret(namespace string, secretDraft *corev1.Secret) (*corev1.Secret, error) {
	secretClient := client.clientset.CoreV1().Secrets(namespace)

	secret, err := secretClient.Update(client.context, secretDraft, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return secret, nil
}
