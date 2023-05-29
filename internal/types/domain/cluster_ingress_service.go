package types

type ClusterIngressService struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Port      *int32 `json:"port"`
}
