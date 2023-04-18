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

type Deployment struct {
	ScaleEvent domainTypes.DeploymentScaleEvent `json:"scale_event"`
}

type updateScaleRequest struct {
	Deployment     Deployment              `json:"deployment"`
	Containers     []domainTypes.Container `json:"containers"`
	DeploymentHost string                  `json:"deployment_url"`
	Project        domainTypes.Project     `json:"project"`
}

// @Description Update Kubernetes Deployment Scale.
// @Param deployment Id path int true "unique Uffizzi Deployment ID"
// @Param spec body updateScaleRequest true "Uffizzi Deployment specification"
// @Success 200 "OK"
// @Failure 500 "most internal errors"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Response 404 "namespace not found"
// @Security BasicAuth
// @Accept json
// @Produce json
// @Router /deployments/{deploymentId}/replicas [put]
func (h *Handlers) handleUpdateScale(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	var request updateScaleRequest

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
		deploymentHost := request.DeploymentHost
		project := request.Project
		scaleEvent := request.Deployment.ScaleEvent

		err = domainLogic.UpdateScale(
			scaleEvent,
			deploymentId,
			containerList,
			deploymentHost,
			project,
		)

		if err != nil {
			handleDomainError("domainLogic.ApplyContainers", err, localHub)
		}
	}(sentry.CurrentHub().Clone())
}
