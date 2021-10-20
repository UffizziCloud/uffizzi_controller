package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @Description Get Kubernetes Service Resources for a Uffizzi Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Produce json
// @Security BasicAuth
// @Router /deployments/{deploymentId}/services [get]
func (h *Handlers) handleGetServices(w http.ResponseWriter, r *http.Request) {
	// Get path vars
	vars := mux.Vars(r)

	// Get deployment id
	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64)
	if err != nil {
		handleError(err, w, r)
		return
	}

	// Configure scope
	localHub := h.getLocalHub(deploymentId)

	// Get services
	services, err := domainLogic.GetServices(deploymentId)
	if err != nil {
		handleDomainError("domainLogic.GetServices", err, localHub)
		return
	}

	// Handle response
	respondWithJSON(w, r, http.StatusOK, services)
}
