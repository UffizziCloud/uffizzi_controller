package http

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

// @Description Get the Default Ingress that handles most incoming requests.
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce json
// @Router /default_ingress/service [get]
func (h *Handlers) handleGetDefaultIngressService(w http.ResponseWriter, r *http.Request) {
	localHub := sentry.CurrentHub().Clone()

	service, err := domainLogic.GetDefaultIngressService()
	if err != nil {
		handleDomainError("domainLogic.GetDefaultIngressService", err, localHub)
		return
	}

	// Handle response
	respondWithJSON(w, r, http.StatusOK, service)
}

// @Description Get all Ingresses for a namespace.
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce json
// @Router /namespaces/{namespace}/ingresses [get]
func (h *Handlers) handleGetNamespaceIngresses(w http.ResponseWriter, r *http.Request) {
	localHub := sentry.CurrentHub().Clone()
	vars := mux.Vars(r)
	namespaceName := vars["namespace"]

	ingresses, err := domainLogic.GetIngresses(namespaceName)
	if err != nil {
		handleDomainError("domainLogic.GetIngresses", err, localHub)
		return
	}

	// Handle response
	respondWithJSON(w, r, http.StatusOK, ingresses)
}
