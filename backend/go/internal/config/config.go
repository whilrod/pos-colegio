package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("no se pudo cargar el archivo .env: %w", err)
		}
	}

	databaseURL := os.Getenv("SUPABASE_DB_URL")

	if databaseURL == "" {
		return nil, fmt.Errorf("SUPABASE_DB_URL no está configurada")
	}

	return &Config{
		DatabaseURL: databaseURL,
	}, nil
}
