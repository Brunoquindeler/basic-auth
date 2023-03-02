package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
)

type application struct {
	auth struct {
		username string
		password string
	}
}

func NewApplication(username, password string) *application {
	return &application{
		auth: struct {
			username string
			password string
		}{username: username, password: password},
	}
}

func (a application) Validate() {
	if a.auth.username == "" {
		log.Fatal("basic auth username must be provided ")
	}

	if a.auth.password == "" {
		log.Fatal("basic auth password must be provided ")
	}
}

func (a application) protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the protected handler")
}

func (a application) unprotectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the unprotected handler")
}

func (a *application) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(a.auth.username))
			expectedPasswordHash := sha256.Sum256([]byte(a.auth.password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
