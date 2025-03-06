package batch

import (
	"context"
	"sync"
	"time"
)

// BatchProcessor processes a batch of requests
type BatchProcessor[T any, R any] struct {
	process       func(ctx context.Context, items []T) ([]R, []error)
	maxBatchSize  int
	maxWaitTime   time.Duration
	batch         []T
	results       []ResultPromise[R]
	mutex         sync.Mutex
	timer         *time.Timer
	processingCtx context.Context
	cancelCtx     context.CancelFunc
}

// ResultPromise represents a promise for an operation result
type ResultPromise[R any] struct {
	result     R
	err        error
	done       bool
	resultChan chan struct{}
}

// NewResultPromise creates a new result promise
func NewResultPromise[R any]() ResultPromise[R] {
	return ResultPromise[R]{
		resultChan: make(chan struct{}, 1),
	}
}

// Wait waits for the result to be available
func (rp *ResultPromise[R]) Wait() (R, error) {
	<-rp.resultChan
	return rp.result, rp.err
}

// Set sets the result
func (rp *ResultPromise[R]) Set(result R, err error) {
	rp.result = result
	rp.err = err
	rp.done = true
	close(rp.resultChan)
}

// Options configures a batch processor
type Options struct {
	// MaxBatchSize is the maximum number of items in a batch
	MaxBatchSize int
	// MaxWaitTime is the maximum time to wait before processing a batch
	MaxWaitTime time.Duration
}

// New creates a new batch processor
func New[T any, R any](
	process func(ctx context.Context, items []T) ([]R, []error),
	options *Options,
) *BatchProcessor[T, R] {
	if options == nil {
		options = &Options{
			MaxBatchSize: 100,
			MaxWaitTime:  50 * time.Millisecond,
		}
	}

	// Ensure minimum values
	if options.MaxBatchSize < 1 {
		options.MaxBatchSize = 1
	}
	if options.MaxWaitTime < time.Millisecond {
		options.MaxWaitTime = time.Millisecond
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &BatchProcessor[T, R]{
		process:       process,
		maxBatchSize:  options.MaxBatchSize,
		maxWaitTime:   options.MaxWaitTime,
		batch:         make([]T, 0, options.MaxBatchSize),
		results:       make([]ResultPromise[R], 0, options.MaxBatchSize),
		processingCtx: ctx,
		cancelCtx:     cancel,
	}
}

// Process adds an item to the batch and returns a promise for the result
func (bp *BatchProcessor[T, R]) Process(ctx context.Context, item T) (*ResultPromise[R], error) {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	// Check if context is already canceled
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Create result promise
	resultPromise := NewResultPromise[R]()

	// Add item to batch
	bp.batch = append(bp.batch, item)
	bp.results = append(bp.results, resultPromise)

	// Start timer if this is the first item
	if len(bp.batch) == 1 {
		bp.timer = time.AfterFunc(bp.maxWaitTime, bp.processBatch)
	}

	// Process batch immediately if it's full
	if len(bp.batch) >= bp.maxBatchSize {
		bp.timer.Stop()
		go bp.processBatch()
	}

	return &resultPromise, nil
}

// processBatch processes the current batch
func (bp *BatchProcessor[T, R]) processBatch() {
	bp.mutex.Lock()

	// Check if there's anything to process
	if len(bp.batch) == 0 {
		bp.mutex.Unlock()
		return
	}

	// Take the current batch and results
	batch := bp.batch
	results := bp.results

	// Reset batch and results
	bp.batch = make([]T, 0, bp.maxBatchSize)
	bp.results = make([]ResultPromise[R], 0, bp.maxBatchSize)

	// Stop timer if it's running
	if bp.timer != nil {
		bp.timer.Stop()
		bp.timer = nil
	}

	bp.mutex.Unlock()

	// Process the batch
	batchResults, batchErrors := bp.process(bp.processingCtx, batch)

	// Distribute results to promises
	for i := range results {
		if i < len(batchResults) {
			var err error
			if i < len(batchErrors) {
				err = batchErrors[i]
			}
			results[i].Set(batchResults[i], err)
		} else {
			// This shouldn't happen if the processor is implemented correctly
			var zero R
			results[i].Set(zero, ErrProcessorIncomplete)
		}
	}
}

// Close closes the batch processor
func (bp *BatchProcessor[T, R]) Close() {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	// Cancel context
	bp.cancelCtx()

	// Process any pending items
	if len(bp.batch) > 0 {
		// Take the current batch and results
		batch := bp.batch
		results := bp.results

		// Reset batch and results
		bp.batch = nil
		bp.results = nil

		// Stop timer if it's running
		if bp.timer != nil {
			bp.timer.Stop()
			bp.timer = nil
		}

		// Process the batch
		go func() {
			batchResults, batchErrors := bp.process(context.Background(), batch)

			// Distribute results to promises
			for i := range results {
				if i < len(batchResults) {
					var err error
					if i < len(batchErrors) {
						err = batchErrors[i]
					}
					results[i].Set(batchResults[i], err)
				} else {
					var zero R
					results[i].Set(zero, ErrProcessorClosed)
				}
			}
		}()
	}
}

// Errors
var (
	ErrProcessorClosed     = Error("batch processor closed")
	ErrProcessorIncomplete = Error("batch processor returned incomplete results")
)

// Error is a simple string error
type Error string

func (e Error) Error() string {
	return string(e)
}
