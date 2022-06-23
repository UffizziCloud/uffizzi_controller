package kuber

import (
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client Client) FindOrInitializePersistentVolumeClaim(
	namespace, name string) (*corev1.PersistentVolumeClaim, error) {
	persistentVolumeClaim, err := client.GetPersistentVolumeClaim(namespace, name)
	if err != nil && !errors.IsNotFound(err) {
		return persistentVolumeClaim, err
	}

	if persistentVolumeClaim != nil && len(persistentVolumeClaim.UID) > 0 {
		return persistentVolumeClaim, nil
	}

	var storageClassName string = global.Settings.PvcStorageClassName

	persistentVolumeClaimDraft := &corev1.PersistentVolumeClaim{
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
			},
		},
	}

	return persistentVolumeClaimDraft, nil
}

func (client *Client) GetPersistentVolumeClaim(namespace, name string) (*corev1.PersistentVolumeClaim, error) {
	persistentVolumeClaimClient := client.clientset.CoreV1().PersistentVolumeClaims(namespace)

	persistentVolumeClaim, err := persistentVolumeClaimClient.Get(client.context, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return persistentVolumeClaim, nil
}

func (client *Client) CreatePersistentVolumeClaim(
	namespace string,
	persistentVolumeClaimDraft *corev1.PersistentVolumeClaim,
) (*corev1.PersistentVolumeClaim, error) {
	persistentVolumeClaimClient := client.clientset.CoreV1().PersistentVolumeClaims(namespace)

	persistentVolumeClaim, err := persistentVolumeClaimClient.Create(
		client.context, persistentVolumeClaimDraft, metav1.CreateOptions{})

	if err != nil {
		return nil, err
	}

	return persistentVolumeClaim, nil
}
