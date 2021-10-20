package http

import (
	"net/http"

	"github.com/getsentry/sentry-go"
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
		handleDomainError("domainLogic.GetIngresses", err, localHub)
		return
	}

	// Handle response
	respondWithJSON(w, r, http.StatusOK, service)
}
