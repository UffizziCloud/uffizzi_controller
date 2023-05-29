package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	apiUffizziClusterV1 "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

type createClusterRequest struct {
	Helm           []apiUffizziClusterV1.HelmChart   `json:"helm"`
	IngressService domainTypes.ClusterIngressService `json:"ingress_service"`
	Name           string                            `json:"name"`
	DeploymentHost string                            `json:"deployment_url"`
}

// @Description Create or Update containers within a Deployment.
// @Param deploymentId path int true "unique Uffizzi Deployment ID"
// @Param spec body createClusterRequest true "container specification"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /deployments/{deploymentId}/clusters [post]
func (h *Handlers) handleCreateCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deploymentId, err := strconv.ParseUint(vars["deploymentId"], 10, 64) //nolint: gomnd
	if err != nil {
		handleError(err, w, r)
		return
	}

	var request createClusterRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", request)

	go func(localHub *sentry.Hub) {
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("deploymentId", fmt.Sprint(deploymentId))
		})

		name := request.Name
		helm := request.Helm
		ingressService := request.IngressService
		deploymentHost := request.DeploymentHost

		err = domainLogic.CreateCluster(
			deploymentId,
			name,
			helm,
			ingressService,
			deploymentHost,
		)
		if err != nil {
			handleDomainError("domainLogic.ApplyContainers", err, localHub)
		}
	}(sentry.CurrentHub().Clone())
}
