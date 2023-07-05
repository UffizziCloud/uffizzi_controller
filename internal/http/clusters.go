package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

type createClusterRequest struct {
	Name            string `json:"name"`
	Manifest        string `json:"manifest"`
	BaseIngressHost string `json:"base_ingress_host"`
}

// @Description Create a cluster within a Namespace.
// @Param namespace path string true "unique Uffizzi Namespace"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /namespaces/{namespace}/cluster [post]
func (h *Handlers) handleCreateCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var request createClusterRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", request)

	namespaceName := vars["namespace"]
	name := request.Name
	manifest := request.Manifest
	baseIngressHost := request.BaseIngressHost

	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("namespace", fmt.Sprint(namespaceName))
	})

	cluster, err := domainLogic.CreateCluster(name, namespaceName, manifest, baseIngressHost)

	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusOK, cluster)
}

// @Description Get a virtual cluster within a Namespace.
// @Param namespace path string true "unique Uffizzi Namespace"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /namespaces/{namespace}/cluster [get]
func (h *Handlers) handleGetCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespaceName := vars["namespace"]
	name := vars["name"]

	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("namespace", fmt.Sprint(namespaceName))
	})

	cluster, err := domainLogic.GetCluster(name, namespaceName)

	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusOK, cluster)
}
