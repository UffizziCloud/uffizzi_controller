package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

func TestHandleGetPods(t *testing.T) {
	// TODO fix this
	t.Skip("Skipping broken test")

	// Create request
	req, err := http.NewRequest("GET", "/deployments/1/pods", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create handlers
	h := &Handlers{}

	// Create handler
	handler := http.HandlerFunc(h.handleGetContainers)

	// -----
	// TEST SUCCESS
	// -----

	// Serve
	handler.ServeHTTP(rr, req)

	// Evaluate status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("error: expected: %v | received: %v", http.StatusOK, status)
	}

	// -----
	// TEST ERRORS
	// -----

	// Set varian
	handler = http.HandlerFunc(h.handleGetContainers)

	// Serve
	handler.ServeHTTP(rr, req)

	// Evaluate status code
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("error: expected: %v | received: %v", http.StatusInternalServerError, status)
	}
}

func TestHandleApplyDeployment(t *testing.T) {
	t.Skip("Skipping broken test")

	var requestBody struct {
		Containers []domainTypes.Container
	}

	var containerID uint64 = 1

	var deploymentID = 1

	var tag = "latest"

	requestBody.Containers = []domainTypes.Container{
		{
			ID:    containerID,
			Image: "nginx",
			Tag:   &tag,
			Variables: []*domainTypes.ContainerVariable{
				{
					Name:  "PORT",
					Value: "80",
				},
			},
		},
	}

	requestBodyJson, err := json.Marshal(&requestBody)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("/deployments/%v/deployments", deploymentID)

	req, err := http.NewRequest("POST", path, bytes.NewReader(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	h := Handlers{}
	router.HandleFunc("/deployments/{deploymentId:[0-9]+}/deployments", h.handleApplyContainers)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
