package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandleCreateNamespace(t *testing.T) {
	t.Skip("Skipping broken test")

	var deploymentID = 1

	var requestBody struct {
		Kind string `json:"kind"`
	}

	requestBody.Kind = "development"

	requestBodyJson, err := json.Marshal(&requestBody)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("/deployments/%v/namespace", deploymentID)

	req, err := http.NewRequest("POST", path, bytes.NewReader(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	h := Handlers{}
	router.HandleFunc("/deployments/{deploymentId:[0-9]+}/namespace", h.handleCreateNamespace)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandleDeleteNamespace(t *testing.T) {
	t.Skip("Skipping broken test")

	var deploymentID = 1

	path := fmt.Sprintf("/deployments/%v/namespace", deploymentID)

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	h := Handlers{}
	router.HandleFunc("/deployments/{deploymentId:[0-9]+}/namespace", h.handleDeleteNamespace)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
