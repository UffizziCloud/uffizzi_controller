package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	apiUffizziClusterV1 "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

type createClusterRequest struct {
	Helm           []apiUffizziClusterV1.HelmChart   `json:"helm"`
	IngressService domainTypes.ClusterIngressService `json:"ingress_service"`
	Name           string                            `json:"name"`
}

// @Description Create a clusters within a Namespace.
// @Param namespace in path string true "unique Uffizzi Namespace"
// @Param spec body createClusterRequest true "cluster specification"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /clusters [post]
func (h *Handlers) handleCreateCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespaceName := vars["namespace"]

	var request createClusterRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", request)

	go func(localHub *sentry.Hub) {
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("namespace", fmt.Sprint(namespaceName))
		})

		name := request.Name
		helm := request.Helm
		ingressService := request.IngressService

		err = domainLogic.CreateCluster(
			namespaceName,
			name,
			helm,
			ingressService,
		)
		if err != nil {
			handleDomainError("domainLogic.CreateCluster", err, localHub)
		}
	}(sentry.CurrentHub().Clone())
}
