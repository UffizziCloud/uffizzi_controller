package basic_auth_utils

import (
	"fmt"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/apr1_crypt"
)

func GenerateAuthPair(login string, password string) string {
	crypter := crypt.APR1.New()
	hash, _ := crypter.Generate([]byte(password), []byte(""))

	return fmt.Sprintf("%s:%s", login, hash)
}
