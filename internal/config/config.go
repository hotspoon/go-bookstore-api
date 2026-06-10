package config

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	APIVersion     string
	DatabasePath   string
	LogPath        string
	GinMode        string
	AllowedOrigins []string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		Port:         envOrDefault("PORT", "8080"),
		APIVersion:   envOrDefault("API_VERSION", "v1"),
		DatabasePath: envOrDefault("DB_PATH", "bookstore.db"),
		LogPath:      envOrDefault("LOG_PATH", "app.log"),
		GinMode:      envOrDefault("GIN_MODE", gin.DebugMode),
	}

	origins := envOrDefault("FRONTEND_WEB_URL", "http://localhost:3000")
	cfg.AllowedOrigins = splitAndTrim(origins)

	gin.SetMode(cfg.GinMode)

	return cfg
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
