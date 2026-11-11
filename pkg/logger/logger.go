package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func New(level string, format string) (*slog.Logger, error) {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		return nil, fmt.Errorf("unsupported log level: %s", level)
	}

	options := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler

	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, options)
	case "console":
		handler = slog.NewTextHandler(os.Stdout, options)
	default:
		return nil, fmt.Errorf("unsupported log format: %s", format)
	}

	return slog.New(handler), nil
}
