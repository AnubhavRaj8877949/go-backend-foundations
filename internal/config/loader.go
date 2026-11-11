package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// loadDotEnv loads local environment variables from .env.
func loadDotEnv() error {
	// 1. Pehle current directory mein check karein (.env)
	err := godotenv.Load(".env")
	if err == nil {
		return nil
	}

	// 2. Agar current dir mein nahi mili, toh root directory (../../.env) check karein
	if errors.Is(err, os.ErrNotExist) {
		if rootErr := godotenv.Load("../../.env"); rootErr == nil {
			return nil
		}
	}

	// 3. Agar koi real syntax/parsing error ho toh log return karein
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env file: %w", err)
	}

	return nil
}

// getEnv returns an environment variable or its default value.
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

// getInt reads an environment variable and converts it to an integer.
func getInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer for %s: %w", key, err)
	}

	return parsedValue, nil
}

// Load reads environment variables and creates application configuration.
func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		return Config{}, err
	}

	httpPort, err := getInt("HTTP_PORT", 8080)

	if err != nil {
		return Config{}, err
	}

	databasePort, err := getInt("DB_PORT", 5432)
	if err != nil {
		return Config{}, err
	}

	redisDatabase, err := getInt("REDIS_DATABASE", 0)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "Go Backend Foundation"),
			Environment: getEnv("APP_ENV", "development"),
		},
		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST", "0.0.0.0"),
			Port: httpPort,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     databasePort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "go_backend_foundation"),
		},
		Redis: RedisConfig{
			Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			Database: redisDatabase,
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", ""),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "console"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}
