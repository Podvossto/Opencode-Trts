// Purpose: App configuration — loads env vars, exposes typed Config struct
// Ref: ADR — all config from environment, no config files
package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all runtime configuration for the API server.
type Config struct {
	Port         string
	DatabaseURL  string
	RedisURL     string
	JWTSecret    string
	JWTExpiry    int // minutes
	AIServiceURL string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPass     string
	SMTPFrom     string
	CORSOrigins  []string // allowed origins for CORS
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://ats_user:change_me_in_production@localhost:5432/ats_db?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:    getEnv("JWT_SECRET", "change_me_minimum_32_chars_secret_key"),
		JWTExpiry:    getEnvInt("JWT_EXPIRY_MINUTES", 30),
		AIServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:8000"),
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnvInt("SMTP_PORT", 587),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPass:     getEnv("SMTP_PASS", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@tigersoft.com"),
		CORSOrigins:  getEnvSlice("CORS_ORIGINS", []string{"http://localhost:3000", "http://localhost:3001"}),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvSlice(key string, fallback []string) []string {
	if v := os.Getenv(key); v != "" {
		parts := strings.Split(v, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				out = append(out, trimmed)
			}
		}
		return out
	}
	return fallback
}
