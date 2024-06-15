package main

import (
	"log"
	"net/http"

	"github.com/anditakaesar/uwa-server-checker/internal/env"
)

func main() {
	mux := http.NewServeMux()
	env := env.New()
	server := &http.Server{
		Addr:    env.AppPort(),
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't start server with err: %v", err)
	}
}
