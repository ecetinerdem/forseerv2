package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unAuthorizedBasicError(w, r, fmt.Errorf("missing authorization header"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unAuthorizedBasicError(w, r, fmt.Errorf("melformed authorization header"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unAuthorizedBasicError(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)

			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unAuthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return

			}

			next.ServeHTTP(w, r)
		})
	}
}
