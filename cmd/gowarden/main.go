package main

import (
	"log"
	"net/http"

	"github.com/QVedant/GoWarden/internal/api"
	"github.com/QVedant/GoWarden/internal/registry"
)

func main() {
	reg, err := registry.Load("config/languages.yaml")
	if err != nil {
		log.Fatalf("failed to load language registry: %v", err)
	}
	log.Printf("loaded %d language(s): %v", len(reg.Names()), reg.Names())

	srv := api.NewServer(reg)

	addr := ":8080"
	log.Printf("GoWarden listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.Router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
