package http

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
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

	namespaceName := vars["namespace"]

	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("namespace", fmt.Sprint(namespaceName))
	})

	cluster, err := domainLogic.CreateCluster(namespaceName)

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

	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("namespace", fmt.Sprint(namespaceName))
	})

	cluster, err := domainLogic.GetCluster(namespaceName)

	if err != nil {
		handleError(err, w, r)
		return
	}

	respondWithJSON(w, r, http.StatusOK, cluster)
}
