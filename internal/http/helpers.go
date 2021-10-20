package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
)

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	sentry.GetHubFromContext(r.Context()).CaptureException(err)
	log.Printf("HTTP error. err=%v", err)

	body := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	respondWithJSON(w, r, http.StatusInternalServerError, body)
}

func handleDomainError(useCaseName string, err error, localHub *sentry.Hub) {
	log.Printf("DomainError: %v. err=%v\n", useCaseName, err)
	localHub.CaptureException(err)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	// Marshal JSON
	response, err := json.Marshal(payload)
	if err != nil {
		sentry.GetHubFromContext(r.Context()).CaptureException(err)
		log.Printf("error marshal http response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if code == http.StatusNoContent {
		return
	}

	// Write response
	_, err = w.Write(response)
	if err != nil {
		sentry.GetHubFromContext(r.Context()).CaptureException(err)
		log.Printf("error write http response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}
