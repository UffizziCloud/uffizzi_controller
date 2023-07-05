package types

type Cluster struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
	Status    struct {
		Ready      bool   `json:"ready"`
		KubeConfig string `json:"kubeConfig"`
	} `json:"status"`
}
