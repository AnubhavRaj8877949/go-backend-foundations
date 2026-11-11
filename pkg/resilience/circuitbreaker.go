package resilience

import (
    "sync"
    "time"
)
type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

type CircuitBreakerConfig struct {
    FailureThreshold int
    OpenTimeout      time.Duration
}
type CircuitBreaker struct {
	 mu sync.Mutex
    config CircuitBreakerConfig
    state  State

    failures  int
    openedAt  time.Time
}

func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
    return &CircuitBreaker{
        config: config,
        state:  StateClosed,
    }
}


func (cb *CircuitBreaker) Allow() bool {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    switch cb.state {
    case StateClosed:
        return true

    case StateOpen:
        if time.Since(cb.openedAt) < cb.config.OpenTimeout {
            return false
        }

        cb.state = StateHalfOpen
        return true

    case StateHalfOpen:
        return false

    default:
        return false
    }
}

func (cb *CircuitBreaker) RecordSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.failures = 0
    cb.state = StateClosed
    cb.openedAt = time.Time{}
}



func (cb *CircuitBreaker) RecordFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.failures++

    if cb.failures >= cb.config.FailureThreshold {
        cb.state = StateOpen
        cb.openedAt = time.Now()
    }
}


func (cb *CircuitBreaker) Execute(operation func() error) error {
    if operation == nil {
        return errors.New("operation is nil")
    }

    if !cb.Allow() {
        return ErrCircuitOpen
    }

    err := operation()
    if err != nil {
        cb.RecordFailure()
        return err
    }

    cb.RecordSuccess()
    return nil
}