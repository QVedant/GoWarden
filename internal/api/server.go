package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/QVedant/GoWarden/internal/registry"
)

type Server struct {
	reg *registry.Registry
	mux *http.ServeMux
}

func NewServer(reg *registry.Registry) *Server {
	s := &Server{
		reg: reg,
		mux: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealth)
	s.mux.HandleFunc("GET /languages", s.handleLanguages)
	// POST /run comes once the executor + nsjail integration exists.
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
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
