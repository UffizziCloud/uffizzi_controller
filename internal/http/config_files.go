package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

type ApplyConfigFileRequest struct {
	ConfigFile domainTypes.ConfigFile `json:"config_file"`
}

// @Description create config file
// @Param configFileId path int true "Config file ID"
// @Param spec body ApplyConfigFileRequest true "Specification"
// @Success 200
// @Router /deployments/{deploymentId}/config_files/{configFileId} [post]
func (h *Handlers) handleApplyConfigFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentID, err := strconv.ParseUint(vars["deploymentId"], 10, 64)
	if err != nil {
		handleError(err, w, r)
		return
	}

	var request ApplyConfigFileRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", request)

	err = domainLogic.ApplyConfigFile(deploymentID, request.ConfigFile)
	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusOK, nil)
}
