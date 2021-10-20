package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/requests"
)

// @Description Fetch metrics for all containers within a Deployment.
// @Param deploymentId path int true "unique Deployment ID"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Produce json
// @Security BasicAuth
// @Router /deployments/{deploymentId}/containers/metrics [get]
func (h *Handlers) handleGetContainersMetrics(w http.ResponseWriter, r *http.Request) {
	// Get path vars
	vars := mux.Vars(r)

	// Get deployment id
	deploymentID, err := strconv.ParseUint(vars["deploymentId"], 10, 64)
	if err != nil {
		handleError(err, w, r)
		return
	}

	// Configure scope
	localHub := h.getLocalHub(deploymentID)

	// Get pod metrics
	metrics, err := domainLogic.GetContainersMetrics(deploymentID)
	if err != nil {
		handleDomainError("domainLogic.GetContainersMetrics", err, localHub)
		return
	}

	// Handle response
	respondWithJSON(w, r, http.StatusOK, metrics)
}

// @Description Fetch memory usage for all containers within a Deployment.
// nolint:lll
// @Param deployment_ids[] query requests.GetContainersUsageMetricsRequestSpec.DeploymentIDs true "array of Uffizzi Deployment ID's"
// @Param begin_at query requests.GetContainersUsageMetricsRequestSpec.BeginAt false "time range start"
// @Param end_at query requests.GetContainersUsageMetricsRequestSpec.EndAt false "time range finish"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Produce json
// @Security BasicAuth
// @Router /deployments/usage_metrics/containers [get]
func (h *Handlers) handleGetContainersUsageMetrics(w http.ResponseWriter, r *http.Request) {
	// Configure scope
	localHub := sentry.CurrentHub().Clone()

	var rawRequest requests.GetContainersUsageMetricsRequestSpec

	queries := r.URL.Query()

	log.Printf("handleGetContainersUsageMetrics: %+#v", queries)

	rawRequest.DeploymentIDs = queries["deployment_ids[]"]
	rawRequest.BeginAt = queries.Get("begin_at")
	rawRequest.EndAt = queries.Get("end_at")

	parsedRequest, err := rawRequest.Parse()
	if err != nil {
		handleError(err, w, r)
		return
	}

	usageMetrics, err := domainLogic.GetDeploymentsContainersUsageMetrics(
		parsedRequest.DeploymentIDs,
		parsedRequest.BeginAt,
		parsedRequest.EndAt,
	)
	if err != nil {
		handleDomainError("domainLogic.GetDeploymentsContainersUsageMetrics", err, localHub)
		return
	}

	// Handle response
	respondWithJSON(w, r, http.StatusOK, usageMetrics)
}
