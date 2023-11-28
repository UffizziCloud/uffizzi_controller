package v1alpha1

import (
	"context"
	"encoding/json"

	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type UffizziClusterInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.UffizziClusterList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.UffizziCluster, error)
	Create(UffizziClusterProps) (*v1alpha1.UffizziCluster, error)
	Patch(name string, patchClusterProps PatchUffizziClusterProps) error
	Delete(name string) error
}

type UffizziClusterClient struct {
	restClient rest.Interface
	ns         string
}

type UffizziClusterProps struct {
	Name string
	Spec v1alpha1.UffizziClusterSpec
}

type PatchUffizziClusterProps struct {
	Spec v1alpha1.UffizziClusterSpec
}

type JSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func (c *UffizziClusterClient) List(opts metav1.ListOptions) (*v1alpha1.UffizziClusterList, error) {
	result := v1alpha1.UffizziClusterList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("UffizziClusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *UffizziClusterClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.UffizziCluster, error) {
	result := v1alpha1.UffizziCluster{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Name(name).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *UffizziClusterClient) Create(clusterProps UffizziClusterProps) (*v1alpha1.UffizziCluster, error) {
	uffizziCluster := v1alpha1.UffizziCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "UffizziCluster",
			APIVersion: "uffizzi.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterProps.Name,
		},
		Spec: clusterProps.Spec,
	}

	result := v1alpha1.UffizziCluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Body(&uffizziCluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

// The sleep mode only works if a patch is applied (not "update")
func (c *UffizziClusterClient) Patch(
	name string,
	patchClusterProps PatchUffizziClusterProps,
) error {

	uffizziCluster := &v1alpha1.UffizziCluster{}

	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Name(name).
		Do(context.TODO()).
		Into(uffizziCluster)

	if err != nil {
		return err
	}

	resourceVersion := uffizziCluster.ObjectMeta.ResourceVersion

	patchOps := []JSONPatchOperation{
		JSONPatchOperation{
			Op:    "replace",
			Path:  "/spec/sleep",
			Value: patchClusterProps.Spec.Sleep,
		},
	}

	patchType := types.JSONPatchType
	patchBytes, err := json.Marshal(patchOps)
	if err != nil {
		return err
	}

	err = c.restClient.
		Patch(patchType).
		Namespace(c.ns).
		Resource("UffizziClusters").
		Name(name).
		Body(patchBytes).
		Param("resourceVersion", resourceVersion).
		Do(context.TODO()).
		Error()

	if err != nil {
		return err
	}

	return nil
}

func (c *UffizziClusterClient) Delete(name string) error {
	result := v1alpha1.UffizziCluster{}
	err := c.restClient.
		Delete().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Name(name).
		Do(context.TODO()).
		Into(&result)

	return err
}
