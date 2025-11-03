package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Application struct {
}

type envelope map[string]any

func main() {
	router := chi.NewRouter()
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(envelope{"status": "ok"})
	})

	mux := chi.NewRouter()
	mux.Mount("/v1", router)

	srv := http.Server{
		Addr:    "0.0.0.0:6000",
		Handler: mux,
	}
	log.Println("start api server on ", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
