package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/michael-grace/cin-2021-tweets/pkg/logging"
)

type authHandler struct {
	handler http.Handler
}

func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()

	if ok {
		if username == strings.TrimSuffix(os.Getenv("AUTH_USER"), "\r") &&
			password == strings.TrimSuffix(os.Getenv("AUTH_PASS"), "\r") {
			h.handler.ServeHTTP(w, r)
			return
		}
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="controller", charset="UTF-8"`)
	http.Error(w, "Unauthorised", http.StatusUnauthorized)
}

type wsAuthHandler struct {
	authToken uuid.UUID
}

func (h wsAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(struct {
		WSAuth string `json:"wsAuthToken"`
	}{
		WSAuth: h.authToken.String(),
	})

	if err != nil {
		logging.Error(err)
	}

	fmt.Fprint(w, string(json))
}
