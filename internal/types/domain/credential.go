package types

import (
	"encoding/base64"
	"fmt"
)

// Credential example
type Credential struct {
	ID          uint64 `json:"id"`
	RegistryUrl string `json:"registry_url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}
type CredentialAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

type CredentialConfig struct {
	CredentialAuths map[string]CredentialAuth `json:"auths"`
}

func NewCredentialAuth(username, password string) CredentialAuth {
	msg := fmt.Sprintf("%v:%v", username, password)
	authBase64 := base64.StdEncoding.EncodeToString([]byte(msg))

	auths := CredentialAuth{
		Username: username,
		Password: password,
		Auth:     authBase64,
	}

	return auths
}
