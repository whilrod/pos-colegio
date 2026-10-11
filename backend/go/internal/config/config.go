package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL            string
	SessionExpirationHours int
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

	sessionExpirationHours, err := strconv.Atoi(
		os.Getenv("SESSION_EXPIRATION_HOURS"),
	)
	if err != nil || sessionExpirationHours <= 0 {
		return nil, fmt.Errorf(
			"SESSION_EXPIRATION_HOURS debe ser un número entero mayor que cero",
		)
	}

	return &Config{
		DatabaseURL:            databaseURL,
		SessionExpirationHours: sessionExpirationHours,
	}, nil
}
