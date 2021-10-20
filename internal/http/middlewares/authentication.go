package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/docs") {
			next.ServeHTTP(w, r)
			return
		}

		canAccess, err := checkAuthorization(r)
		if canAccess && err == nil {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
	})
}

func checkAuthorization(r *http.Request) (bool, error) {
	controllerLogin := global.Settings.ControllerLogin
	controllerPassword := global.Settings.ControllerPassword

	if controllerLogin == "" || controllerPassword == "" {
		return true, nil
	}

	login, password, ok := r.BasicAuth()

	if !ok {
		return false, errors.New("incorrect token for basic auth")
	}

	if controllerLogin == login && controllerPassword == password {
		return true, nil
	}

	return false, errors.New("incorrect login or password")
}
