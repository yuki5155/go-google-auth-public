package config

import (
	"os"
	"strings"
)

type Config struct {
	Port              string
	Environment       string
	AllowedOrigins    []string
	FrontendURL       string
	GoogleClientID    string
	GoogleSecret      string
	GoogleRedirectURL string
}

func Load() *Config {
	// CORS Allowed Origins - comma separated
	allowedOriginsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:5173")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	return &Config{
		Port:              getEnv("PORT", "8080"),
		Environment:       getEnv("GO_ENV", "development"),
		AllowedOrigins:    allowedOrigins,
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:5173"),
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleSecret:      getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL: getEnv("GOOGLE_REDIRECT_URL", ""),
	}
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
