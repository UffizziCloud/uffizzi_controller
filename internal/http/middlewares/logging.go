package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		formattedRequest, err := formatRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println(formattedRequest)
		next.ServeHTTP(w, r)
	})
}

func formatRequest(r *http.Request) (string, error) {
	request := []string{fmt.Sprintf("%v %v %v Host: %v", r.Method, r.URL, r.Proto, r.Host)}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			return "", err
		}

		request = append(request, r.Form.Encode())
	}

	return strings.Join(request, "\n"), nil
}
