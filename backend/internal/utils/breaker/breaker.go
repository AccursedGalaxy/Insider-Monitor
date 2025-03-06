package breaker

import (
	"errors"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	// StateIsClosed indicates the circuit is closed and operating normally
	StateIsClosed CircuitState = iota
	// StateIsOpen indicates the circuit is open (failing)
	StateIsOpen
	// StateIsHalfOpen indicates the circuit is half-open (testing recovery)
	StateIsHalfOpen
)

var (
	// ErrCircuitOpen is returned when a circuit is open
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// Options configures a circuit breaker
type Options struct {
	// Threshold is the number of failures allowed before opening
	Threshold int
	// Timeout is the duration the circuit remains open before testing recovery
	Timeout time.Duration
	// OnStateChange is called when circuit state changes
	OnStateChange func(from, to CircuitState)
	// HalfOpenMaxRequests is the max number of requests allowed in half-open state
	HalfOpenMaxRequests int
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mutex          sync.RWMutex
	state          CircuitState
	failCount      int
	lastStateTime  time.Time
	threshold      int
	timeout        time.Duration
	onStateChange  func(from, to CircuitState)
	halfOpenMax    int
	halfOpenCount  int
	halfOpenResult sync.Once
}

// New creates a new circuit breaker
func New(options *Options) *CircuitBreaker {
	if options == nil {
		options = &Options{
			Threshold:           5,
			Timeout:             10 * time.Second,
			HalfOpenMaxRequests: 1,
		}
	}

	// Ensure minimum values
	if options.Threshold < 1 {
		options.Threshold = 1
	}
	if options.Timeout < time.Second {
		options.Timeout = time.Second
	}
	if options.HalfOpenMaxRequests < 1 {
		options.HalfOpenMaxRequests = 1
	}

	return &CircuitBreaker{
		state:          StateIsClosed,
		failCount:      0,
		lastStateTime:  time.Now(),
		threshold:      options.Threshold,
		timeout:        options.Timeout,
		onStateChange:  options.OnStateChange,
		halfOpenMax:    options.HalfOpenMaxRequests,
		halfOpenResult: sync.Once{},
	}
}

// Execute runs the given function guarded by the circuit breaker
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.AllowRequest() {
		return ErrCircuitOpen
	}

	err := fn()
	cb.ReportResult(err == nil)
	return err
}

// AllowRequest checks if a request should be permitted
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case StateIsClosed:
		return true
	case StateIsOpen:
		// Check if we should enter half-open state
		if time.Since(cb.lastStateTime) > cb.timeout {
			// Only transition in the write lock
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			defer cb.mutex.Unlock()

			// Need to double-check after acquiring the write lock
			if cb.state == StateIsOpen && time.Since(cb.lastStateTime) > cb.timeout {
				cb.setState(StateIsHalfOpen)
				cb.halfOpenResult = sync.Once{}
				cb.halfOpenCount = 0
				return true
			}
			return false
		}
		return false
	case StateIsHalfOpen:
		// Allow limited requests in half-open state
		return cb.halfOpenCount < cb.halfOpenMax
	default:
		return false
	}
}

// ReportResult reports the result of a request execution
func (cb *CircuitBreaker) ReportResult(success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateIsClosed:
		if success {
			// Reset fail count on success
			cb.failCount = 0
		} else {
			// Increment fail count on failure
			cb.failCount++
			if cb.failCount >= cb.threshold {
				cb.setState(StateIsOpen)
			}
		}
	case StateIsHalfOpen:
		cb.halfOpenCount++

		// Use Once to ensure we only make the transition once, even if multiple
		// goroutines report results simultaneously
		if success {
			cb.halfOpenResult.Do(func() {
				// Success in half-open state closes the circuit
				cb.setState(StateIsClosed)
				cb.failCount = 0
			})
		} else {
			cb.halfOpenResult.Do(func() {
				// Failure in half-open state re-opens the circuit
				cb.setState(StateIsOpen)
			})
		}
	}
}

// Reset forces the circuit breaker back to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.setState(StateIsClosed)
	cb.failCount = 0
}

// Trip forces the circuit breaker to open (useful for testing)
func (cb *CircuitBreaker) Trip() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.setState(StateIsOpen)
}

// State returns the current circuit state
func (cb *CircuitBreaker) State() CircuitState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	return cb.state
}

// setState changes the circuit state and calls the state change handler
func (cb *CircuitBreaker) setState(newState CircuitState) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState
	cb.lastStateTime = time.Now()

	if cb.onStateChange != nil {
		go cb.onStateChange(oldState, newState)
	}
}
