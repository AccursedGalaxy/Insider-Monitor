package rate

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// Limiter wraps the standard rate limiter
type Limiter struct {
	limiter *rate.Limiter
}

// NewLimiter creates a new rate limiter
func NewLimiter(r rate.Limit, b int) *Limiter {
	return &Limiter{
		limiter: rate.NewLimiter(r, b),
	}
}

// Wait blocks until the limiter allows an event to happen
func (l *Limiter) Wait(ctx context.Context) error {
	return l.limiter.Wait(ctx)
}

// Allow checks if an operation is allowed
func (l *Limiter) Allow() bool {
	return l.limiter.Allow()
}

// Every returns a rate limiter that allows events at the specified interval
func Every(interval time.Duration) rate.Limit {
	return rate.Every(interval)
}
