package server

import (
	"net/http"

	"backend-pos-colegio/internal/handler"

	"github.com/go-chi/chi/v5"
)

func New(healthHandler *handler.HealthHandler) *http.Server {
	router := chi.NewRouter()

	router.Get("/health", healthHandler.Check)

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
