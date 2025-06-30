package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db     DbConfig
	Server ServerConfig
}

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}
type ServerConfig struct {
	Port string
}

func Load() (*Config, error) {
	slog.Info("config | Loading config")
	if err := godotenv.Load(); err != nil {
		slog.Warn("config | No .env file found")
	}

	cfg := &Config{
		Db: DbConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Database: getEnv("DB_DATABASE", "postgres"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
