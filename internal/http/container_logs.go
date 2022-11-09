package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
)

// @Description Fetch logs for a specific container.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Param containerName path string true "container name"
// @Param limit query int false "maximum number of lines to return"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Produce json
// @Security BasicAuth
// @Router /deployments/{deploymentId}/containers/{containerName}/logs [get]
func (h *Handlers) handleGetContainerLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("handleGetContainerLogs: %+v\n", vars)

	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint:gomnd

	if err != nil {
		handleError(err, w, r)
		return
	}

	queries := r.URL.Query()

	limit := global.Settings.CountDisplayedEntriesForLogsOutput
	limitRaw := queries.Get("limit")
	previousValue := queries.Get("previous")
	previous, err := strconv.ParseBool(previousValue)

	if err != nil {
		handleError(err, w, r)
		return
	}

	if limitRaw != "" {
		limit, err = strconv.ParseInt(limitRaw, 10, 64) //nolint:gomnd
		if err != nil {
			handleError(err, w, r)
			return
		}
	}

	pods, err := domainLogic.GetContainers(deploymentId)
	if err != nil {
		handleError(err, w, r)
		return
	}

	response := struct {
		Logs []string
	}{
		Logs: []string{},
	}

	if len(pods) > 0 {
		logs, err := domainLogic.GetPodLogs(deploymentId, pods[0].Name, vars["containerName"], limit, previous)
		if err != nil {
			localHub := h.getLocalHub(deploymentId)
			handleDomainError("domainLogic.handleGetContainerLogs", err, localHub)

			return
		}

		response.Logs = logs
	}

	respondWithJSON(w, r, http.StatusOK, response)
}
