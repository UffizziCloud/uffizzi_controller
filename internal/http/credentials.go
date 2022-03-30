package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	types "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

// @Description Add or Update credential within a Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Param spec body types.Credential true "credential specification"
// @Success 201 "created successfully"
// @Failure 500 "most errors including Not Found"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Security BasicAuth
// @Accept json
// @Produce json
// @Router /deployments/{deploymentId}/credentials [post]
func (h *Handlers) handleApplyCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentID, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	var credential types.Credential

	err = json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", credential)

	namespace, err := domainLogic.ApplyCredential(deploymentID, credential)
	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusCreated, namespace)
}

// @Description Delete credential from a Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Param credentialId path int true "сredential ID"
// @Success 204 "no content (success)"
// @Failure 500 "most errors including Not Found"
// @Response 403 "incorrect token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce json
// @Router /deployments/{deploymentId}/credentials/{credentialId} [delete]
func (h *Handlers) handleDeleteCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentID, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	credentialID, err := strconv.ParseUint(vars["credentialId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	err = domainLogic.DeleteCredential(deploymentID, credentialID)
	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusNoContent, nil)
}
