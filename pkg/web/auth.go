/**
URY Tweet Board
Copyright (C) 2022 Michael Grace <michael.grace@ury.org.uk>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
