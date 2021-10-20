package http

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

// @Description Fetch list of Kubernetes nodes.
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Router /nodes [get]
// @Produce json
// @Security BasicAuth
func (h *Handlers) handleGetNodes(w http.ResponseWriter, r *http.Request) {
	localHub := sentry.CurrentHub().Clone()

	// Get nodes
	nodes, err := domainLogic.GetNodes()
	if err != nil {
		handleDomainError("domainLogic.GetNodes", err, localHub)
		return
	}

	respondWithJSON(w, r, http.StatusOK, nodes)
}
