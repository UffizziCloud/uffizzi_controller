package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

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

	var clusterParams domainTypes.ClusterParams

	err := json.NewDecoder(r.Body).Decode(&clusterParams)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", clusterParams)

	namespaceName := vars["namespace"]

	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("namespace", fmt.Sprint(namespaceName))
	})

	cluster, err := domainLogic.CreateCluster(namespaceName, clusterParams)

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
// @Router /namespaces/{namespace}/cluster/{name} [get]
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

// @Description Update a virtual cluster within a Namespace.
// @Param namespace path string true "unique Uffizzi Namespace"
// @Param name path string true "Uffizzi Virtual Cluster Name"
// @Success 200 "OK"
// @Failure 500 "most errors including Not Found"
// @Response 403 "Incorrect Token for HTTP Basic Auth"
// @Security BasicAuth
// @Produce plain
// @Router /namespaces/{namespace}/cluster/{name} [put]
func (h *Handlers) handlePatchCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	namespaceName := vars["namespace"]
	name := vars["name"]

	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("namespace", fmt.Sprint(namespaceName))
	})

	var patchClusterParams domainTypes.PatchClusterParams

	err := json.NewDecoder(r.Body).Decode(&patchClusterParams)
	if err != nil {
		handleError(err, w, r)
		return
	}

	log.Printf("Decoded HTTP Request: %+v", patchClusterParams)

	err = domainLogic.PatchCluster(name, namespaceName, patchClusterParams)

	if err != nil {
		handleDomainError("domainLogic.PatchCluster", err, localHub)
	}

	respondWithJSON(w, r, http.StatusOK, nil)
}
