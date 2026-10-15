package server

import (
	"net/http"

	"backend-pos-colegio/internal/handler"
	"backend-pos-colegio/internal/middleware"
	"backend-pos-colegio/internal/service"

	"github.com/go-chi/chi/v5"
)

func New(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
	productHandler *handler.ProductHandler,
	variantHandler *handler.VariantHandler,
) *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)

	router.Get("/health", healthHandler.Check)
	router.Post("/auth/login", authHandler.Login)

	router.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authService))
		r.Get("/auth/session", authHandler.Session)
		r.Post("/auth/logout", authHandler.Logout)

		r.With(middleware.RequirePermission("CONSULTAR_PRODUCTO")).Get("/productos", productHandler.List)
		r.With(middleware.RequirePermission("CONSULTAR_VARIANTE")).Get("/variantes", variantHandler.List)
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
