package config

import "testing"

func TestConfigValidate(t *testing.T) {
	cfg := Config{
		HTTP: HTTPConfig{
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Database: "go_backend_foundation",
		},
		Redis: RedisConfig{
			Address: "localhost:6379",
		},
		Logging: LoggingConfig{
			Level: "debug",
		},
	}

	err := cfg.Validate()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestConfigValidateInvalidHTTPPort(t *testing.T) {
	cfg := Config{
		HTTP: HTTPConfig{
			Port: 70000,
		},
	}

	err := cfg.Validate()

	if err == nil {
		t.Fatal("expected validation error for invalid HTTP port")
	}
}
