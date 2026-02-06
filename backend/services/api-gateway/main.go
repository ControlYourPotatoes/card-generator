package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Logger setup (structured JSON for Docker, console for dev)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API Versioning: /api/v1/*
	r.Route("/api/v1", func(r chi.Router) {
		// Health endpoint - no auth required (used by Docker healthcheck)
		r.Get("/health", healthHandler)

		// Stub endpoints (501 Not Implemented until services ready)
		r.Post("/cards/generate", stubHandler("card-generator"))
		r.Get("/cards/{id}", stubHandler("card-generator"))
		r.Get("/cards/{id}/render", stubImageHandler("image-renderer"))
		r.Post("/cards/analyze", stubHandler("card-generator"))
		r.Post("/import/csv", stubHandler("importer"))
		r.Get("/import/{jobId}/status", stubHandler("importer"))
		r.Delete("/admin/cards", stubHandler("card-generator"))
	})

	log.Info().Msg("API Gateway starting on :8080")
	log.Info().Msg("Endpoints ready (stubs until Phase 2)")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"status": "ok",
		"services": map[string]string{
			"api-gateway": "ok",
			"postgres":    "available", // Can ping later
			// Others added in Phase 2
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func stubHandler(service string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		log.Info().Str("request_id", reqID).Str("service", service).Msg("Routing stub")

		w.Header().Set("X-Request-ID", reqID)
		http.Error(w, `{"error":"Service `+service+` not implemented (Phase 2)","code":"NOT_IMPLEMENTED"}`, http.StatusNotImplemented)
	}
}

func stubImageHandler(service string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		log.Info().Str("request_id", reqID).Str("service", service).Msg("Image stub")

		w.Header().Set("X-Request-ID", reqID)
		w.Header().Set("Content-Type", "image/png")
		http.Error(w, "Image renderer not implemented", http.StatusNotImplemented)
	}
}