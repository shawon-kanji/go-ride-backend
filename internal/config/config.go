package config

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/shawon-kanji/go-ride-utils/awssecrets"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret         string
	ExpiryMinutes  int
	Issuer         string
	Audience       string
	DriverAudience string
}

// Load reads configuration from the environment. In staging/production,
// DB_CREDENTIALS_SECRET_NAME and/or JWT_SECRET_NAME select AWS Secrets
// Manager as the source for those values instead (fetched here via IRSA);
// when unset, the plain env vars below are used exactly as before.
func Load(ctx context.Context) (*Config, error) {
	dbPort, err := getIntEnv("DB_PORT", 5432)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	jwtExpiry, err := getIntEnv("JWT_EXPIRY_MINUTES", 60)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_MINUTES: %w", err)
	}

	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	if secretName := getEnv("DB_CREDENTIALS_SECRET_NAME", ""); secretName != "" {
		values, err := awssecrets.FetchJSON(ctx, secretName)
		if err != nil {
			return nil, fmt.Errorf("fetch db credentials secret: %w", err)
		}
		if v, ok := values["DB_USER"]; ok {
			dbUser = v
		}
		if v, ok := values["DB_PASSWORD"]; ok {
			dbPassword = v
		}
	}

	jwtSecret := getEnv("JWT_SECRET", "change-me-in-production")
	if secretName := getEnv("JWT_SECRET_NAME", ""); secretName != "" {
		values, err := awssecrets.FetchJSON(ctx, secretName)
		if err != nil {
			return nil, fmt.Errorf("fetch jwt secret: %w", err)
		}
		if v, ok := values["JWT_SECRET"]; ok {
			jwtSecret = v
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("APP_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     getEnv("DB_NAME", "go_ride"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:         jwtSecret,
			ExpiryMinutes:  jwtExpiry,
			Issuer:         getEnv("JWT_ISSUER", "go-ride-backend"),
			Audience:       getEnv("JWT_AUDIENCE", "go-ride-clients"),
			DriverAudience: getEnv("JWT_DRIVER_AUDIENCE", "go-ride-drivers"),
		},
	}

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	return strconv.Atoi(value)
}
