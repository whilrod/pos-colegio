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
	inventoryHandler *handler.InventoryHandler,
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
		r.With(middleware.RequirePermission("CREAR_PRODUCTO")).Post("/productos", productHandler.Create)
		r.With(middleware.RequirePermission("CONSULTAR_VARIANTE")).Get("/variantes", variantHandler.List)
		r.With(middleware.RequirePermission("CREAR_VARIANTE")).Post("/productos/{idproducto}/variantes", variantHandler.Create)
		r.With(middleware.RequirePermission("CONSULTAR_PRODUCTO")).Get("/productos/{idproducto}", productHandler.GetByID)
		r.With(middleware.RequirePermission("EDITAR_PRODUCTO")).Put("/productos/{idproducto}", productHandler.Update)
		r.With(middleware.RequirePermission("DESACTIVAR_PRODUCTO")).Post("/productos/{idproducto}/desactivar", productHandler.Deactivate)
		r.With(middleware.RequirePermission("CONSULTAR_INVENTARIO")).Get("/inventario", inventoryHandler.List)
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
