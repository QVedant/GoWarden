package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/QVedant/GoWarden/internal/registry"
)

type Server struct {
	reg    *registry.Registry
	Router *chi.Mux
}

func NewServer(reg *registry.Registry) *Server {
	s := &Server{
		reg:    reg,
		Router: chi.NewRouter(),
	}

	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(10 * time.Second))

	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.Get("/healthz", s.handleHealth)
	s.Router.Get("/languages", s.handleLanguages)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleLanguages(w http.ResponseWriter, r *http.Request) {
	names := s.reg.Names()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]string{"languages": names}); err != nil {
		log.Printf("failed to encode languages response: %v", err)
	}
}
