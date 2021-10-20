package requests

import (
	"errors"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

// GetContainersUsageMetricsRequestSpec represents request spec
type GetContainersUsageMetricsRequestSpec struct {
	BeginAt       string   `json:"begin_at" example:"2020-14-07T15:058:05Z07:00"`
	EndAt         string   `json:"end_at" example:"2020-14-07T16:58:05Z07:00"`
	DeploymentIDs []string `json:"deployment_ids[]"`
}

// GetContainersUsageMetricsRequestParsed represents parsed request
type GetContainersUsageMetricsRequestParsed struct {
	BeginAt       time.Time
	EndAt         time.Time
	DeploymentIDs []uint64
}

// Parse parses raw request to typed one
func (rawRequest *GetContainersUsageMetricsRequestSpec) Parse() (*GetContainersUsageMetricsRequestParsed, error) {
	if len(rawRequest.BeginAt) == 0 {
		rawRequest.BeginAt = time.Now().AddDate(0, 0, -1).String()
	}

	if len(rawRequest.EndAt) == 0 {
		rawRequest.EndAt = time.Now().String()
	}

	if len(rawRequest.DeploymentIDs) == 0 {
		return nil, errors.New("empty deployment_ids[]")
	}

	beginAt, err := dateparse.ParseAny(rawRequest.BeginAt)
	if err != nil {
		return nil, errors.New("wrong begin_at format")
	}

	endAt, err := dateparse.ParseAny(rawRequest.EndAt)
	if err != nil {
		return nil, errors.New("wrong end_at format")
	}

	var deploymentIDs []uint64

	for _, deploymentIDStr := range rawRequest.DeploymentIDs {
		deploymentID, err := strconv.ParseUint(deploymentIDStr, 10, 64)
		if err != nil {
			return nil, err
		}

		deploymentIDs = append(deploymentIDs, deploymentID)
	}

	result := &GetContainersUsageMetricsRequestParsed{
		BeginAt:       beginAt,
		EndAt:         endAt,
		DeploymentIDs: deploymentIDs,
	}

	return result, nil
}
