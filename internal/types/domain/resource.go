package types

type ResourceKind string

const (
	ResourceKindConfigMap ResourceKind = "config_map"
	ResourceKindSecret    ResourceKind = "secret"
)

type ResourceVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Resource struct {
	ID        uint64             `json:"id"`
	Kind      ResourceKind       `json:"kind"`
	Variables []ResourceVariable `json:"variables,omitempty"`
}
