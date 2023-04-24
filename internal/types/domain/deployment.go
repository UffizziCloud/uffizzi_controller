package types

type DeploymentScaleEvent string

const (
	DeploymentScaleEventScaleDown DeploymentScaleEvent = "scale_down"
	DeploymentScaleEventScaleUp   DeploymentScaleEvent = "scale_up"
)
