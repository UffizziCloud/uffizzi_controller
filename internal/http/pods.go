package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @Description Fetch Pod events TODO
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Produce json
// @Security BasicAuth
// @Router /deployments/{deploymentId}/containers/events [get]
func (h *Handlers) handleGetPodEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("handleGetPodEvents: %+v\n", vars)

	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64)
	if err != nil {
		handleError(err, w, r)
		return
	}

	events, err := domainLogic.GetPodEvents(deploymentId)
	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusOK, events)
}
