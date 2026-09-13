// Package config loads application configuration from environment variables.
package config

import "os"

// Config holds all runtime configuration for the application.
type Config struct {
	// DatabaseURL is the full PostgreSQL connection string.
	// Defaults to a local development database when not set.
	DatabaseURL string

	// Port is the HTTP port the server will listen on.
	// Defaults to "8080".
	Port string
}

// Load reads configuration from environment variables and returns a Config
// populated with sensible defaults for local development.
func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
