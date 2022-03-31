package domain

import (
	"time"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
)

type ContainersUsageMetrics struct {
	ContainersMemory float64 `json:"containers_memory"`
}

func (l *Logic) GetDeploymentsContainersUsageMetrics(deploymentIDs []uint64, beginAt time.Time, endAt time.Time) (ContainersUsageMetrics, error) {
	return ContainersUsageMetrics{
		ContainersMemory: 0,
	}, nil
}

func (l *Logic) ApplyContainerSecrets(namespace string, containerList domainTypes.ContainerList) error {
	for _, container := range containerList.Items {
		name := global.Settings.ResourceName.ContainerSecret(container.ID)

		secret, err := l.KuberClient.FindOrInitializeSecret(namespace, name)
		if err != nil {
			return err
		}

		secretVariables := map[string]string{}

		for _, secretVariable := range container.SecretVariables {
			secretVariables[secretVariable.Name] = secretVariable.Value
		}

		secret.StringData = secretVariables

		if len(secret.UID) > 0 {
			if len(secret.StringData) > 0 {
				_, err = l.KuberClient.UpdateSecret(namespace, secret)
				if err != nil {
					return err
				}
			} else {
				err = l.KuberClient.DeleteSecret(namespace, secret.Name)
				if err != nil {
					return err
				}
			}
		} else {
			if len(secret.StringData) > 0 {
				_, err = l.KuberClient.CreateSecret(namespace, secret)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
