package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

type updateScaleRequest struct {
	ScaleEvent domainTypes.DeploymentScaleEvent `json:"scale_event"`
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

	deploymentID, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
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

	err = domainLogic.UpdateScale(deploymentID, request.ScaleEvent)

	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusOK, nil)
}
