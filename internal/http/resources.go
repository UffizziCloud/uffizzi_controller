package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

type ApplyResourceRequest struct {
	Resource domainTypes.Resource `json:"resource"`
}

// @Description Create or update Uffizzi Resource.
// @Param resourceId path int true "unique Uffizzi Resource ID"
// @Param spec body ApplyResourceRequest true "Uffizzi Resource specification"
// @Success 200 "OK"
// @Router /resources/{resourceId} [post]
func (h *Handlers) handleApplyResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentID, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	var request ApplyResourceRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", request)

	err = domainLogic.ApplyResource(deploymentID, request.Resource)
	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusOK, nil)
}
