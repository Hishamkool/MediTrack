// Package config reads settings from environment variables.
// In local development they come from backend/.env (loaded by godotenv).
package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	TokenTTL       time.Duration
	AllowedOrigins []string // frontend URLs allowed by CORS
	AutoMigrate    bool     // run DB migrations when the API starts

	// Used once, to create the first admin if no admin exists yet.
	AdminName     string
	AdminEmail    string
	AdminPassword string
}

// Load reads the config. It returns an error if a required value is missing.
func Load() (*Config, error) {
	// Ignore the error: on Render/Fly there is no .env file, the
	// variables are set in the dashboard instead.
	_ = godotenv.Load()

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		TokenTTL:       12 * time.Hour,
		AllowedOrigins: splitList(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),
		AutoMigrate:    getEnv("AUTO_MIGRATE", "true") == "true",
		AdminName:      getEnv("ADMIN_NAME", "Administrator"),
		AdminEmail:     os.Getenv("ADMIN_EMAIL"),
		AdminPassword:  os.Getenv("ADMIN_PASSWORD"),
	}

	if ttl := os.Getenv("TOKEN_TTL"); ttl != "" {
		d, err := time.ParseDuration(ttl) // e.g. "12h"
		if err != nil {
			return nil, errors.New("TOKEN_TTL must be a duration like 12h")
		}
		cfg.TokenTTL = d
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, errors.New("JWT_SECRET is required and must be at least 32 characters")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// splitList turns "a, b,c" into ["a","b","c"].
func splitList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
