package retry

import (
    "context"
    "fmt"
    "time"
)

type Config struct {
    Attempts int
    Delay    time.Duration
}

func Do(
    ctx context.Context,
    cfg Config,
    fn func() error,
) error {
    if cfg.Attempts < 1 {
        return fmt.Errorf("retry attempts must be greater than zero")
    }

    if cfg.Delay <= 0 {
        return fmt.Errorf("retry delay must be greater than zero")
    }

    var lastErr error
    delay := cfg.Delay

    for attempt := 1; attempt <= cfg.Attempts; attempt++ {
        if err := fn(); err == nil {
            return nil
        } else {
            lastErr = err
        }

        if attempt == cfg.Attempts {
            break
        }

        timer := time.NewTimer(delay)

        select {
        case <-ctx.Done():
            timer.Stop()
            return ctx.Err()
        case <-timer.C:
        }

        delay *= 2
    }

    return fmt.Errorf("operation failed after %d attempts: %w", cfg.Attempts, lastErr)
}