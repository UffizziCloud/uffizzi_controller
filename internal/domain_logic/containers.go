package domain

import (
	"time"
)

type ContainersUsageMetrics struct {
	ContainersMemory float64 `json:"containers_memory"`
}

func (l *Logic) GetDeploymentsContainersUsageMetrics(deploymentIDs []uint64, beginAt time.Time, endAt time.Time) (ContainersUsageMetrics, error) {
	return ContainersUsageMetrics{
		ContainersMemory: 0,
	}, nil
}
