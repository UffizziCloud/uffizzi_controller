package domain

import (
	"encoding/json"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (l *Logic) ApplyCredential(deploymentID uint64, credential domainTypes.Credential) (*corev1.Secret, error) {
	namespaceName := l.KubernetesNamespaceName(deploymentID)
	credentialName := global.Settings.ResourceName.Credential(credential.ID)

	secret, err := l.KuberClient.FindOrInitializeSecret(namespaceName, credentialName)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	secret.Type = corev1.SecretTypeDockerConfigJson

	config := domainTypes.CredentialConfig{
		CredentialAuths: map[string]domainTypes.CredentialAuth{
			credential.RegistryUrl: domainTypes.NewCredentialAuth(credential.Username, credential.Password),
		},
	}

	dockerconfigjson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	secret.Data = map[string][]byte{
		".dockerconfigjson": dockerconfigjson,
	}

	if len(secret.UID) == 0 {
		secret, err := l.KuberClient.CreateSecret(namespaceName, secret)
		if err != nil {
			return nil, err
		}

		return secret, nil
	}

	secret, err = l.KuberClient.UpdateSecret(namespaceName, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (l *Logic) DeleteCredential(deploymentID, credentialId uint64) error {
	namespaceName := l.KubernetesNamespaceName(deploymentID)
	secretName := global.Settings.ResourceName.Credential(credentialId)

	err := l.KuberClient.DeleteSecret(namespaceName, secretName)
	if err != nil {
		return err
	}

	return nil
}
