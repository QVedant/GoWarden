package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/QVedant/GoWarden/internal/executor"
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
	s.Router.Use(middleware.Timeout(15 * time.Second))

	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.Get("/healthz", s.handleHealth)
	s.Router.Get("/languages", s.handleLanguages)
	s.Router.Post("/run", s.handleRun)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleLanguages(w http.ResponseWriter, r *http.Request) {
	names := s.reg.Names()
	writeJSON(w, http.StatusOK, map[string][]string{"languages": names})
}

type runRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type runResponse struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	TimedOut   bool   `json:"timed_out"`
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	if req.Language == "" || req.Code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "language and code are required"})
		return
	}

	lang, ok := s.reg.Get(req.Language)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unsupported language: " + req.Language})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(lang.TimeoutSeconds+2)*time.Second)
	defer cancel()

	result, err := executor.Run(ctx, lang, req.Code)
	if err != nil {
		log.Printf("execution error for language %q: %v", req.Language, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "execution failed"})
		return
	}

	writeJSON(w, http.StatusOK, runResponse{
		Stdout:     result.Stdout,
		Stderr:     result.Stderr,
		ExitCode:   result.ExitCode,
		DurationMs: result.Duration.Milliseconds(),
		TimedOut:   result.TimedOut,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
