package config

import (
	"fmt"
	"strings"
)

func (c Config) Validate() error {
	// Validate HTTP port.
	if c.HTTP.Port < 1 || c.HTTP.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.HTTP.Port)
	}

	// Validate database configuration.
	if strings.TrimSpace(c.Database.Host) == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.Port < 1 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}

	if strings.TrimSpace(c.Database.User) == "" {
		return fmt.Errorf("database user is required")
	}

	if strings.TrimSpace(c.Database.Database) == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate Redis configuration.
	if strings.TrimSpace(c.Redis.Address) == "" {
		return fmt.Errorf("redis address is required")
	}

	if c.Redis.Database < 0 {
		return fmt.Errorf("redis database cannot be negative")
	}

	// Validate logging configuration.
	level := strings.ToLower(c.Logging.Level)

	switch level {
	case "debug", "info", "warn", "error":
		// Valid logging level.
	default:
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	return nil
}
