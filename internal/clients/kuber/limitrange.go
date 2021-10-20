package kuber

import (
	"fmt"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FindLimitRange - find limit range by name
func (client *Client) FindLimitRange(namespace string) (*corev1.LimitRange, error) {
	limitRanges := client.clientset.CoreV1().LimitRanges(namespace)

	limitRange, err := limitRanges.Get(client.context, client.LimitRangeName(namespace), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return limitRange, nil
}

// UpdateLimitRangeAttributes - update limit range attributes
func (client *Client) UpdateLimitRangeAttributes(
	limitRange *corev1.LimitRange,
	memoryLimit string,
) (*corev1.LimitRange, error) {
	resourceMemory, err := resource.ParseQuantity(memoryLimit)
	if err != nil {
		return nil, err
	}

	ephemeralStorage := ephemeralStorageProportion(resourceMemory)

	resourceCPU := *cpuProportion(resourceMemory)

	resourceStorage, err := resource.ParseQuantity("0")
	if err != nil {
		return nil, err
	}

	limitRange.Spec = corev1.LimitRangeSpec{
		Limits: []corev1.LimitRangeItem{{
			Type: corev1.LimitTypePersistentVolumeClaim,
			Max: corev1.ResourceList{
				corev1.ResourceStorage:          resourceStorage,
				corev1.ResourceMemory:           resourceMemory,
				corev1.ResourceCPU:              resourceCPU,
				corev1.ResourceEphemeralStorage: *ephemeralStorage,
			},
		}},
	}

	return limitRange, nil
}

// CreateLimitRange - create limit range
func (client *Client) CreateLimitRange(namespace, name, memoryLimit string) (*corev1.LimitRange, error) {
	limitRanges := client.clientset.CoreV1().LimitRanges(namespace)

	draftLimitRange := &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
			},
		},
	}

	draftLimitRange, err := client.UpdateLimitRangeAttributes(draftLimitRange, memoryLimit)
	if err != nil {
		return nil, err
	}

	limitRange, err := limitRanges.Create(client.context, draftLimitRange, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return limitRange, nil
}

// UpdateLimitRange - update limit range
func (client *Client) UpdateLimitRange(limitRange *corev1.LimitRange, memoryLimit string) (*corev1.LimitRange, error) {
	limitRanges := client.clientset.CoreV1().LimitRanges(limitRange.Namespace)

	limitRange, err := client.UpdateLimitRangeAttributes(limitRange, memoryLimit)
	if err != nil {
		return nil, err
	}

	limitRange, err = limitRanges.Update(client.context, limitRange, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return limitRange, nil
}

func (client *Client) LimitRangeName(namespaceName string) string {
	return fmt.Sprintf("limit-range-%v", namespaceName)
}

// CreateOrUpdateLimitRange - create or update limit range
func (client *Client) CreateOrUpdateLimitRange(namespaceName, memoryLimit string) (*corev1.LimitRange, error) {
	limitRange, err := client.FindLimitRange(namespaceName)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	if limitRange != nil {
		limitRange, err = client.UpdateLimitRange(limitRange, memoryLimit)
		if err != nil {
			return nil, err
		}
	} else {
		limitRange, err = client.CreateLimitRange(namespaceName, client.LimitRangeName(namespaceName), memoryLimit)
		if err != nil {
			return nil, err
		}
	}

	return limitRange, nil
}
