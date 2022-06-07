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

type applyIngressBasicAuthRequest struct {
	Project domainTypes.Project `json:"project"`
}

// @Description Create or Update containers within a Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Param spec body applyIngressBasicAuthRequest true "container specification"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /deployments/{deploymentId}/containers [post]
func (h *Handlers) handleApplyIngressBasciAuth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	var request applyIngressBasicAuthRequest

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

		project := request.Project

		err = domainLogic.ApplyIngressBasciAuth(deploymentId, project)
		if err != nil {
			handleDomainError("domainLogic.ApplyIngressBasciAuth", err, localHub)
		}
	}(sentry.CurrentHub().Clone())
}

// @Description Create or Update containers within a Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /deployments/{deploymentId}/containers [post]
func (h *Handlers) handleDeleteIngressBasciAuth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	go func(localHub *sentry.Hub) {
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("deploymentId", fmt.Sprint(deploymentId))
		})

		err = domainLogic.DeleteIngressBasciAuth(deploymentId)
		if err != nil {
			handleDomainError("domainLogic.DeleteIngressBasciAuth", err, localHub)
		}
	}(sentry.CurrentHub().Clone())
}
