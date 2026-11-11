package logger

import "testing"

func TestNewLogger(t *testing.T) {
	logger, err := New("debug", "console")

	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	if logger == nil {
		t.Fatal("expected logger, got nil")
	}

	logger.Info("Logger test successful")
}
