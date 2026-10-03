package main

import (
	"log"
	"net/http"

	"backend-pos-colegio/internal/config"
	"backend-pos-colegio/internal/database"
	"backend-pos-colegio/internal/handler"
	"backend-pos-colegio/internal/repository"
	"backend-pos-colegio/internal/server"
	"backend-pos-colegio/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repository := repository.NewRepository(db)

	healthService := service.NewHealthService(repository)
	healthHandler := handler.NewHealthHandler(healthService)

	httpServer := server.New(healthHandler)

	log.Printf("Backend POS Colegio iniciado en %s", httpServer.Addr)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error en el servidor HTTP: %v", err)
	}
}
