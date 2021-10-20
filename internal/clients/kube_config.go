package clients

import (
	"k8s.io/client-go/rest"
)

func InitializeKubeConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	return config, err
}
