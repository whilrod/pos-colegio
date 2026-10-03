package main

import (
	"fmt"
	"log"
	"os"

	"backend-pos-colegio/internal/database"

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

	fmt.Println("Backend POS Colegio iniciado")
	fmt.Println("Conexión a PostgreSQL establecida correctamente")
}
