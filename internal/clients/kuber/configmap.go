package kuber

import (
	"fmt"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client Client) FindOrInitializeConfigMap(namespace, name string) (*corev1.ConfigMap, error) {
	configMap, err := client.GetConfigMap(namespace, name)
	if err != nil && !errors.IsNotFound(err) {
		return configMap, err
	}

	if configMap != nil && len(configMap.UID) > 0 {
		return configMap, nil
	}

	configMapDraft := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
			},
		},
		Data:       map[string]string{},
		BinaryData: map[string][]byte{},
	}

	return configMapDraft, nil
}

func (client *Client) GetConfigMap(namespace, name string) (*corev1.ConfigMap, error) {
	configMapClient := client.clientset.CoreV1().ConfigMaps(namespace)

	configMap, err := configMapClient.Get(client.context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

func (client *Client) GetConfigMaps(namespace string) ([]corev1.ConfigMap, error) {
	configMapClient := client.clientset.CoreV1().ConfigMaps(namespace)

	configMaps, err := configMapClient.List(client.context, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/managed-by=%v", global.Settings.ManagedApplication),
	})
	if err != nil {
		return nil, err
	}

	return configMaps.Items, nil
}

func (client *Client) DeleteConfigMap(namespace, name string) error {
	configMapClient := client.clientset.CoreV1().ConfigMaps(namespace)

	err := configMapClient.Delete(client.context, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) CreateConfigMap(namespace string, configMapDraft *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	configMapClient := client.clientset.CoreV1().ConfigMaps(namespace)

	configMap, err := configMapClient.Create(client.context, configMapDraft, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

func (client *Client) UpdateConfigMap(namespace string, configMapDraft *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	configMapClient := client.clientset.CoreV1().ConfigMaps(namespace)

	configMap, err := configMapClient.Update(client.context, configMapDraft, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}
