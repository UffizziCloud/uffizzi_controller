package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetPodsMetrics(t *testing.T) {
	// TODO fix this
	t.Skip("Skipping broken test")

	// Create request
	req, err := http.NewRequest("GET", "/deployments/insertVar/pods/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create Handlers
	h := &Handlers{}

	// Create handler
	handler := http.HandlerFunc(h.handleGetContainersMetrics)

	// Serve
	handler.ServeHTTP(rr, req)

	// Evaluate status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("error: expected: %v | received: %v", http.StatusOK, status)
	}
}
