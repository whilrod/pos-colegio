package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env")
	}

	dbURL := os.Getenv("SUPABASE_DB_URL")

	if dbURL == "" {
		log.Fatal("SUPABASE_DB_URL no está configurada")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("No se pudo abrir la conexión a PostgreSQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}

	fmt.Println("Backend POS Colegio iniciado")
	fmt.Println("Conexión a PostgreSQL establecida correctamente")
}
