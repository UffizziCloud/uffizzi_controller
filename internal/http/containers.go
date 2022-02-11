package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

// @Description Fetch metadata on all containers specified by a Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce json
// @Router /deployments/{deploymentId}/containers [get]
func (h *Handlers) handleGetContainers(w http.ResponseWriter, r *http.Request) {
	// Get path vars
	vars := mux.Vars(r)

	// Get deployment id
	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	// Configure scope
	localHub := h.getLocalHub(deploymentId)

	// Get pods
	pods, err := domainLogic.GetContainers(deploymentId)
	if err != nil {
		handleDomainError("domainLogic.GetContainers", err, localHub)
		return
	}

	// Handle response
	respondWithJSON(w, r, http.StatusOK, pods)
}

type applyContainersRequest struct {
	Containers     []domainTypes.Container  `json:"containers"`
	Credentials    []domainTypes.Credential `json:"credentials,omitempty"`
	Resources      []domainTypes.Resource   `json:"resources"`
	DeploymentHost string                   `json:"deployment_url"`
}

// @Description Create or Update containers within a Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Param spec body applyContainersRequest true "container specification"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /deployments/{deploymentId}/containers [post]
func (h *Handlers) handleApplyContainers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	var request applyContainersRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", request)

	go func(localHub *sentry.Hub) {
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("deploymentId", fmt.Sprint(deploymentId))
		})

		containers := request.Containers
		containerList := domainTypes.ContainerList{Items: containers}
		credentials := request.Credentials
		deploymentHost := request.DeploymentHost
		resources := request.Resources

		err = domainLogic.ApplyContainers(deploymentId, containerList, credentials, deploymentHost, resources)
		if err != nil {
			handleDomainError("domainLogic.ApplyContainers", err, localHub)
		}
	}(sentry.CurrentHub().Clone())
}
