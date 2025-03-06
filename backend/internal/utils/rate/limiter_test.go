package rate

import (
	"context"
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	limit := Every(100 * time.Millisecond)
	limiter := NewLimiter(limit, 1)

	if limiter == nil {
		t.Fatal("Expected a non-nil limiter")
	}
}

func TestLimiterAllow(t *testing.T) {
	// Create a limiter that allows 1 event per second with a burst of 1
	limiter := NewLimiter(1, 1)

	// First request should be allowed
	if !limiter.Allow() {
		t.Error("Expected first request to be allowed")
	}

	// Second immediate request should be denied
	if limiter.Allow() {
		t.Error("Expected second immediate request to be denied")
	}

	// After waiting, another request should be allowed
	time.Sleep(1100 * time.Millisecond)
	if !limiter.Allow() {
		t.Error("Expected request after waiting to be allowed")
	}
}

func TestLimiterWait(t *testing.T) {
	// Create a fast limiter for testing
	limiter := NewLimiter(10, 1) // 10 events per second with burst of 1

	// First request should not block
	start := time.Now()
	err := limiter.Wait(context.Background())
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if elapsed > 50*time.Millisecond {
		t.Errorf("Expected first request to be fast, took %v", elapsed)
	}

	// Second request might block briefly but should succeed
	err = limiter.Wait(context.Background())
	if err != nil {
		t.Errorf("Expected no error on second request, got %v", err)
	}
}

func TestLimiterWithCancelledContext(t *testing.T) {
	// Create a very slow limiter
	limiter := NewLimiter(0.1, 1) // Only 1 event per 10 seconds

	// Allow first request
	limiter.Allow()

	// Create a context with a short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// This should timeout since we're trying to make a second request immediately
	err := limiter.Wait(ctx)
	if err == nil {
		t.Error("Expected context deadline exceeded error, got nil")
	}
}

func TestEvery(t *testing.T) {
	// Test different intervals
	intervals := []time.Duration{
		time.Millisecond,
		100 * time.Millisecond,
		time.Second,
		time.Minute,
	}

	for _, interval := range intervals {
		limit := Every(interval)
		if limit <= 0 {
			t.Errorf("Expected positive limit for interval %v, got %v", interval, limit)
		}

		// The rate should be the inverse of the interval
		expected := 1.0 / interval.Seconds()
		if float64(limit) < expected*0.9 || float64(limit) > expected*1.1 {
			t.Errorf("Expected limit around %v for interval %v, got %v", expected, interval, limit)
		}
	}
}
