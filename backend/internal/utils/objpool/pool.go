package objpool

import (
	"sync"
)

// Pool is a generic object pool based on sync.Pool
type Pool[T any] struct {
	pool       sync.Pool
	resetFunc  func(obj T)
	createFunc func() T
}

// NewPool creates a new object pool
func NewPool[T any](createFunc func() T, resetFunc func(obj T)) *Pool[T] {
	p := &Pool[T]{
		createFunc: createFunc,
		resetFunc:  resetFunc,
	}

	p.pool.New = func() any {
		return createFunc()
	}

	return p
}

// Get retrieves an object from the pool or creates a new one
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put returns an object to the pool after resetting it
func (p *Pool[T]) Put(obj T) {
	if p.resetFunc != nil {
		p.resetFunc(obj)
	}
	p.pool.Put(obj)
}

// MessagePool is a specialized pool for message objects
type MessagePool[T any] struct {
	pool *Pool[T]
}

// NewMessagePool creates a message pool with specified initial and max capacity
func NewMessagePool[T any](
	createFunc func() T,
	resetFunc func(obj T),
) *MessagePool[T] {
	return &MessagePool[T]{
		pool: NewPool(createFunc, resetFunc),
	}
}

// Get retrieves a message object from the pool
func (mp *MessagePool[T]) Get() T {
	return mp.pool.Get()
}

// Put returns a message object to the pool
func (mp *MessagePool[T]) Put(msg T) {
	mp.pool.Put(msg)
}

// BufferPool is a pool of byte buffers
type BufferPool struct {
	pool *Pool[[]byte]
}

// NewBufferPool creates a new buffer pool
func NewBufferPool(bufferSize int) *BufferPool {
	createFunc := func() []byte {
		return make([]byte, 0, bufferSize)
	}

	resetFunc := func(buf []byte) {
		buf = buf[:0] // Reset length but keep capacity
	}

	return &BufferPool{
		pool: NewPool(createFunc, resetFunc),
	}
}

// Get retrieves a buffer from the pool
func (bp *BufferPool) Get() []byte {
	return bp.pool.Get()
}

// Put returns a buffer to the pool
func (bp *BufferPool) Put(buf []byte) {
	bp.pool.Put(buf)
}

// MapPool is a pool for maps
type MapPool[K comparable, V any] struct {
	pool *Pool[map[K]V]
}

// NewMapPool creates a new map pool
func NewMapPool[K comparable, V any]() *MapPool[K, V] {
	createFunc := func() map[K]V {
		return make(map[K]V)
	}

	resetFunc := func(m map[K]V) {
		// Clear the map
		for k := range m {
			delete(m, k)
		}
	}

	return &MapPool[K, V]{
		pool: NewPool(createFunc, resetFunc),
	}
}

// Get retrieves a map from the pool
func (mp *MapPool[K, V]) Get() map[K]V {
	return mp.pool.Get()
}

// Put returns a map to the pool
func (mp *MapPool[K, V]) Put(m map[K]V) {
	mp.pool.Put(m)
}
