package http

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type Handlers struct{}

// @Description welcome page and heartbeat
// @Router / [get]
// @Success 200 "OK"
// @produce html
func (h *Handlers) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	fmt.Fprintf(w, "Uffizzi Controller API documentation at "+
		"<a href=\"/docs/\">/docs/</a><br/>HTTP Basic Authentication "+
		"required for all API methods.")
}

func (h *Handlers) getLocalHub(deploymentId uint64) *sentry.Hub {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("deployment_id", fmt.Sprint(deploymentId))
	})

	return localHub
}
