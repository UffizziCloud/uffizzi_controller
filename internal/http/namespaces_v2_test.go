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

func TestHandleCreateNamespaceV2(t *testing.T) {
	t.Skip("Skipping broken test")

	var requestBody struct {
		Namespace string `json:"namespace"`
	}

	requestBody.Namespace = "cluster-11"

	requestBodyJson, err := json.Marshal(&requestBody)
	if err != nil {
		t.Fatal(err)
	}

	path := "/namespaces"

	req, err := http.NewRequest("POST", path, bytes.NewReader(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	h := Handlers{}
	router.HandleFunc("/namespaces", h.handleCreateNamespaceV2)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandleDeleteNamespaceV2(t *testing.T) {
	t.Skip("Skipping broken test")

	var namespace = "cluster-11"

	path := fmt.Sprintf("/deployments/%v", namespace)

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	h := Handlers{}
	router.HandleFunc("/namespaces/{namespace}", h.handleDeleteNamespaceV2)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
