package clients

import (
	"k8s.io/client-go/rest"
)

func InitializeKubeConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	config.QPS = 20
	config.Burst = 40

	if err != nil {
		panic(err.Error())
	}

	return config, err
}
