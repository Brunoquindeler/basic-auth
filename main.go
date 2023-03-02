package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	username := flag.String("username", "", "Define the username")
	password := flag.String("password", "", "Define the password")

	flag.Parse()

	app := NewApplication(*username, *password)

	app.Validate()

	mux := http.NewServeMux()
	mux.HandleFunc("/protected", app.basicAuth(app.protectedHandler))
	mux.HandleFunc("/unprotected", app.unprotectedHandler)

	srv := &http.Server{
		Addr:         ":4000",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServeTLS("./localhost.pem", "./localhost-key.pem")
	log.Fatal(err)
}
