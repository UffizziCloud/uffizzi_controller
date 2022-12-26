package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

type deploymentRequest struct {
	Kind string `json:"kind"`
}

// @Description Fetch the Kubernetes Namespace for a specified Uffizzi Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Success 200 "OK"
// @Failure 500 "most errors"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Response 404 "namespace not found"
// @Security BasicAuth
// @Produce json
// @Router /deployments/{deploymentId} [get]
func (h *Handlers) handleGetNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentID, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	namespace, err := domainLogic.GetNamespace(deploymentID)

	if err != nil && isNotFoundNamespaceError(err) {
		respondWithJSON(w, r, http.StatusNotFound, namespace)
		return
	}

	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusOK, namespace)
}

// @Description Create Kubernetes Namespace for a new Uffizzi Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Param spec body deploymentRequest true "Uffizzi Deployment Specification"
// @Success 201 "created successfully"
// @Failure 500 "most internal errors"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce json
// @Router /deployments/{deploymentId} [post]
func (h *Handlers) handleCreateNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentID, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	var request deploymentRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", request)

	namespace, err := domainLogic.CreateNamespace(
		deploymentID,
		request.Kind,
	)
	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusCreated, namespace)
}

// @Description Update Kubernetes Namespace.
// @Param deployment Id path int true "unique Uffizzi Deployment ID"
// @Param spec body deploymentRequest true "Uffizzi Deployment specification"
// @Success 200 "OK"
// @Failure 500 "most internal errors"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Response 404 "namespace not found"
// @Security BasicAuth
// @Accept json
// @Produce json
// @Router /deployments/{deploymentId} [put]
func (h *Handlers) handleUpdateNamespace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentID, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	var request deploymentRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", request)

	namespace, err := domainLogic.UpdateNamespace(deploymentID, request.Kind)
	if err != nil && isNotFoundNamespaceError(err) {
		respondWithJSON(w, r, http.StatusNotFound, namespace)
		return
	}

	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusOK, namespace)
}

// @Description Delete Kubernetes Namespace and all Resources within.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Success 204 "No Content (success)"
// @Failure 500 "most internal errors"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /deployments/{deploymentId} [delete]
func (h *Handlers) handleDeleteNamespace(w http.ResponseWriter, r *http.Request) {
	// Get path vars
	vars := mux.Vars(r)

	// Get deployment id
	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	go func(localHub *sentry.Hub) {
		// Configure scope
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("deployment_id", fmt.Sprint(deploymentId))
		})
		// DeleteNamespace
		err = domainLogic.DeleteNamespace(deploymentId)

		if err != nil && !isNotFoundNamespaceError(err) {
			handleDomainError("domainLogic.DeleteNamespace", err, localHub)
		}
	}(sentry.CurrentHub().Clone())

	respondWithJSON(w, r, http.StatusNoContent, nil)
}

func isNotFoundNamespaceError(err error) bool {
	if err == nil {
		return false
	}

	notFound, _ := regexp.MatchString(`namespaces.*?not found`, err.Error())

	return notFound
}
