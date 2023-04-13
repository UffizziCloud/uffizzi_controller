package http

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "gitlab.com/dualbootpartners/idyl/uffizzi_controller/docs"
)

func drawRoutes(r *mux.Router, h *Handlers) {
	r.HandleFunc("/", h.handleRoot)
	r.HandleFunc("/nodes", h.handleGetNodes).Methods(http.MethodGet)
	r.HandleFunc("/default_ingress/service", h.handleGetDefaultIngressService).Methods(http.MethodGet)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}", h.handleGetNamespace).Methods(http.MethodGet)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}", h.handleCreateNamespace).Methods(http.MethodPost)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}", h.handleDeleteNamespace).Methods(http.MethodDelete)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/config_files/{configFileId:[0-9]+}", h.handleApplyConfigFile).Methods(http.MethodPost)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/credentials", h.handleApplyCredential).Methods(http.MethodPost)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/credentials/{credentialId}", h.handleDeleteCredential).Methods(http.MethodDelete)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/containers", h.handleGetContainers).Methods(http.MethodGet)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/containers", h.handleApplyContainers).Methods(http.MethodPost)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/containers/metrics", h.handleGetContainersMetrics).Methods(http.MethodGet)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/containers/events", h.handleGetPodEvents).Methods(http.MethodGet)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/services", h.handleGetServices).Methods(http.MethodGet)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/containers/{containerName}/logs", h.handleGetContainerLogs).Methods(http.MethodGet)
	r.HandleFunc("/deployments/usage_metrics/containers", h.handleGetContainersUsageMetrics).Methods(http.MethodGet)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/ingress/basic_auth", h.handleApplyIngressBasciAuth).Methods(http.MethodPost)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/ingress/basic_auth", h.handleDeleteIngressBasciAuth).Methods(http.MethodDelete)
	r.HandleFunc("/deployments/{deploymentId:[0-9]+}/scale", h.handleUpdateScale).Methods(http.MethodPut)
	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
}
