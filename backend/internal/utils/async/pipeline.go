package async

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// Pipeline represents a data processing pipeline
type Pipeline[T any, R any] struct {
	input         chan T
	output        chan R
	errorOutput   chan error
	workers       int
	processing    atomic.Int32
	processorFunc func(ctx context.Context, item T) (R, error)
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
}

// Options configures a pipeline
type Options struct {
	// Workers is the number of worker goroutines
	Workers int
	// InputBuffer is the size of the input channel buffer
	InputBuffer int
	// OutputBuffer is the size of the output channel buffer
	OutputBuffer int
	// ErrorBuffer is the size of the error channel buffer
	ErrorBuffer int
}

// NewPipeline creates a new processing pipeline
func NewPipeline[T any, R any](
	process func(ctx context.Context, item T) (R, error),
	options *Options,
) *Pipeline[T, R] {
	if options == nil {
		options = &Options{
			Workers:      10,
			InputBuffer:  100,
			OutputBuffer: 100,
			ErrorBuffer:  100,
		}
	}

	// Ensure minimum values
	if options.Workers < 1 {
		options.Workers = 1
	}
	if options.InputBuffer < 1 {
		options.InputBuffer = 1
	}
	if options.OutputBuffer < 1 {
		options.OutputBuffer = 1
	}
	if options.ErrorBuffer < 1 {
		options.ErrorBuffer = 1
	}

	ctx, cancel := context.WithCancel(context.Background())

	pipeline := &Pipeline[T, R]{
		input:         make(chan T, options.InputBuffer),
		output:        make(chan R, options.OutputBuffer),
		errorOutput:   make(chan error, options.ErrorBuffer),
		workers:       options.Workers,
		processorFunc: process,
		ctx:           ctx,
		cancel:        cancel,
	}

	// Start worker goroutines
	for i := 0; i < options.Workers; i++ {
		pipeline.wg.Add(1)
		go pipeline.worker()
	}

	return pipeline
}

// Submit adds an item to the pipeline for processing
func (p *Pipeline[T, R]) Submit(item T) error {
	select {
	case p.input <- item:
		p.processing.Add(1)
		return nil
	case <-p.ctx.Done():
		return ErrPipelineClosed
	}
}

// Result returns a channel that receives processing results
func (p *Pipeline[T, R]) Result() <-chan R {
	return p.output
}

// Errors returns a channel that receives processing errors
func (p *Pipeline[T, R]) Errors() <-chan error {
	return p.errorOutput
}

// Close closes the pipeline and waits for all workers to finish
func (p *Pipeline[T, R]) Close() error {
	// Signal cancellation
	p.cancel()

	// Close input channel to signal workers to stop
	close(p.input)

	// Wait for workers to finish
	p.wg.Wait()

	// Close output channels
	close(p.output)
	close(p.errorOutput)

	return nil
}

// worker processes items from the input channel
func (p *Pipeline[T, R]) worker() {
	defer p.wg.Done()

	for {
		select {
		case item, ok := <-p.input:
			if !ok {
				// Input channel closed, exit worker
				return
			}

			// Process the item
			result, err := p.processorFunc(p.ctx, item)

			// Decrement processing counter
			p.processing.Add(-1)

			// Check if context was canceled during processing
			if errors.Is(err, context.Canceled) {
				// Discard result if context was canceled
				continue
			}

			// Send result or error
			if err != nil {
				select {
				case p.errorOutput <- err:
					// Error sent successfully
				case <-p.ctx.Done():
					// Context canceled, exit worker
					return
				}
			} else {
				select {
				case p.output <- result:
					// Result sent successfully
				case <-p.ctx.Done():
					// Context canceled, exit worker
					return
				}
			}

		case <-p.ctx.Done():
			// Context canceled, exit worker
			return
		}
	}
}

// ProcessingCount returns the number of items currently being processed
func (p *Pipeline[T, R]) ProcessingCount() int {
	return int(p.processing.Load())
}

// ErrPipelineClosed indicates the pipeline is closed
var ErrPipelineClosed = errors.New("pipeline is closed")

// Worker executes a task with automatic retry
type Worker[T any] struct {
	tasks        chan func(ctx context.Context) (T, error)
	results      chan T
	errors       chan error
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	workerCount  int
	maxQueueSize int
}

// NewWorker creates a worker pool with the specified number of workers
func NewWorker[T any](workerCount, maxQueueSize int) *Worker[T] {
	if workerCount <= 0 {
		workerCount = 5
	}
	if maxQueueSize <= 0 {
		maxQueueSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	w := &Worker[T]{
		tasks:        make(chan func(ctx context.Context) (T, error), maxQueueSize),
		results:      make(chan T, maxQueueSize),
		errors:       make(chan error, maxQueueSize),
		ctx:          ctx,
		cancel:       cancel,
		workerCount:  workerCount,
		maxQueueSize: maxQueueSize,
	}

	// Start workers
	for i := 0; i < workerCount; i++ {
		w.wg.Add(1)
		go w.runWorker()
	}

	return w
}

// Submit adds a task to the worker pool
func (w *Worker[T]) Submit(task func(ctx context.Context) (T, error)) error {
	select {
	case w.tasks <- task:
		return nil
	case <-w.ctx.Done():
		return ErrWorkerClosed
	default:
		return ErrQueueFull
	}
}

// Results returns the results channel
func (w *Worker[T]) Results() <-chan T {
	return w.results
}

// Errors returns the errors channel
func (w *Worker[T]) Errors() <-chan error {
	return w.errors
}

// Close shuts down the worker pool
func (w *Worker[T]) Close() {
	w.cancel()
	close(w.tasks)
	w.wg.Wait()
	close(w.results)
	close(w.errors)
}

// runWorker processes tasks from the queue
func (w *Worker[T]) runWorker() {
	defer w.wg.Done()

	for {
		select {
		case task, ok := <-w.tasks:
			if !ok {
				return
			}

			result, err := task(w.ctx)

			if err != nil {
				// Don't send errors if context is canceled
				if errors.Is(err, context.Canceled) {
					continue
				}

				select {
				case w.errors <- err:
					// Sent error
				case <-w.ctx.Done():
					return
				}
			} else {
				select {
				case w.results <- result:
					// Sent result
				case <-w.ctx.Done():
					return
				}
			}

		case <-w.ctx.Done():
			return
		}
	}
}

// Errors
var (
	ErrQueueFull    = errors.New("task queue is full")
	ErrWorkerClosed = errors.New("worker is closed")
)
