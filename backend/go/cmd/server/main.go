package main

import (
	"log"
	"net/http"
	"os"

	"backend-pos-colegio/internal/database"
	"backend-pos-colegio/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env")
	}

	dbURL := os.Getenv("SUPABASE_DB_URL")

	if dbURL == "" {
		log.Fatal("SUPABASE_DB_URL no está configurada")
	}

	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	httpServer := server.New()

	log.Printf("Backend POS Colegio iniciado en %s", httpServer.Addr)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error en el servidor HTTP: %v", err)
	}
}
