package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New() *http.Server {
	router := chi.NewRouter()
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
