package domain

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/clients/kuber"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
)

type Logic struct {
	KuberClient *kuber.Client
}

func NewLogic(kuberClient *kuber.Client) *Logic {
	return &Logic{
		KuberClient: kuberClient,
	}
}

func (l Logic) KubernetesNamespaceName(deploymentID uint64) string {
	prefix := global.Settings.NamespaceNamePrefix
	return fmt.Sprintf("%s-%d", prefix, int(deploymentID))
}

func (l Logic) GetDeploymentIDFromKubernetesNamespaceName(namespace string) (uint64, error) {
	prefix := global.Settings.NamespaceNamePrefix + "-"
	rawID := strings.TrimPrefix(namespace, prefix)

	return strconv.ParseUint(rawID, 10, 64)
}
