package main

import (
	"fmt"
	"log"
	"os"

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

	fmt.Println("Backend POS Colegio iniciado")
	fmt.Println("SUPABASE_DB_URL configurada correctamente")
}