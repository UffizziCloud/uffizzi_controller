package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

type namespaceRequest struct {
	Namespace string `json:"namespace"`
}

// @Description Fetch the Kubernetes Namespace for a specified Uffizzi Deployment of Uffizzi Cluster.
// @Param namespace path string true "prefix plus unique Uffizzi Deployment/Cluster ID"
// @Success 200 "OK"
// @Failure 500 "most errors"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Response 404 "namespace not found"
// @Security BasicAuth
// @Produce json
// @Router /namespaces/{namespace} [get]
func (h *Handlers) handleGetNamespaceV2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespaceName := vars["namespace"]

	namespace, err := domainLogic.GetNamespaceV2(namespaceName)

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

// @Description Create Kubernetes Namespace.
// @Param spec body namespaceRequest true "New prefix plus unique Uffizzi Deployment/Cluster ID"
// @Success 201 "created successfully"
// @Failure 500 "most internal errors"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce json
// @Router /namespaces [post]
func (h *Handlers) handleCreateNamespaceV2(w http.ResponseWriter, r *http.Request) {
	var request namespaceRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", request)

	namespaceName := request.Namespace

	namespace, err := domainLogic.CreateNamespaceV2(namespaceName)
	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusCreated, namespace)
}

// @Description Delete Kubernetes Namespace and all Resources within.
// @Param namespace path string true "prefix plus unique Uffizzi Deployment/Cluster ID"
// @Success 204 "No Content (success)"
// @Failure 500 "most internal errors"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /namespaces/{namespace} [delete]
func (h *Handlers) handleDeleteNamespaceV2(w http.ResponseWriter, r *http.Request) {
	// Get path vars
	vars := mux.Vars(r)
	namespaceName := vars["namespace"]

	log.Printf("Namespace Name: %+v", namespaceName)

	go func(localHub *sentry.Hub) {
		// Configure scope
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("namespace", namespaceName)
		})
		// DeleteNamespace
		err := domainLogic.DeleteNamespaceV2(namespaceName)

		if err != nil && !isNotFoundNamespaceError(err) {
			handleDomainError("domainLogic.DeleteNamespace", err, localHub)
		}
	}(sentry.CurrentHub().Clone())

	respondWithJSON(w, r, http.StatusNoContent, nil)
}
