package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/clients"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/clients/kuber"
	internalDomainLogic "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/domain_logic"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
	"k8s.io/client-go/rest"
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

	path := fmt.Sprintf("/deployments/%v/containers", deploymentID)

	req, err := http.NewRequest("POST", path, bytes.NewReader(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	h := Handlers{}
	router.HandleFunc("/deployments/{deploymentId:[0-9]+}/containers", h.handleApplyContainers)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func setup(testHttpClient *http.Client) *rest.Config {
	if err := global.Init(os.Getenv("ENV")); err != nil {
		panic(err)
	}

	config, err := clients.InitializeKubeConfig()
	if err != nil {
		panic(err)
	}

	kuberClient, err := kuber.NewClient2(config, testHttpClient)
	if err != nil {
		panic(err)
	}

	logic := internalDomainLogic.NewLogic(kuberClient)
	domainLogic = logic

	return config
}

func TestStartApplyContainers(t *testing.T) {
	r, err := recorder.New("fixtures/iana-reserved-domains")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	if r.Mode() != recorder.ModeRecordOnce {
		t.Fatal("Recorder should be in ModeRecordOnce")
	}

	r.SetReplayableInteractions(true)
	client := r.GetDefaultClient()
	client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	_ = setup(client)

	// resp, err := client.Get("https://104.154.101.254/api/v1/namespaces/agafonov-env-1")
	// resp, err := client.Get("https://google.com")
	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Printf("%s", string(body))
	requestBody := ApplyContainersRequest{}

	var containerID uint64 = 1
	var deploymentID uint64 = 57
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

	requestBody.Credentials = []domainTypes.Credential{}
	requestBody.DeploymentHost = "deployment-1.app.uffizzi.com"
	requestBody.Project = domainTypes.Project{}
	requestBody.ComposeFile = domainTypes.ComposeFile{}

	err = startApplyContainers(requestBody, deploymentID)

	if err != nil {
		t.Errorf("ERR %v", err)
	}
}
