package examples

import (
	"sync"
	"sync/atomic"
)

// AtomicCounter demonstrates atomic operations for thread-safe counting
type AtomicCounter struct {
	value int64
}

// Increment atomically increments the counter
func (c *AtomicCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Decrement atomically decrements the counter
func (c *AtomicCounter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Get atomically reads the counter value
func (c *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Set atomically sets the counter value
func (c *AtomicCounter) Set(val int64) {
	atomic.StoreInt64(&c.value, val)
}

// CompareAndSet atomically compares and sets the value
func (c *AtomicCounter) CompareAndSet(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.value, old, new)
}

// AtomicConfig demonstrates atomic.Value for configuration hot-reload
type AtomicConfig struct {
	config atomic.Value
}

type Config struct {
	MaxConnections int
	Timeout        int
	Debug          bool
}

// NewAtomicConfig creates a new atomic config with initial values
func NewAtomicConfig(initial Config) *AtomicConfig {
	ac := &AtomicConfig{}
	ac.config.Store(initial)
	return ac
}

// Get returns the current configuration
func (ac *AtomicConfig) Get() Config {
	return ac.config.Load().(Config)
}

// Update atomically updates the configuration
func (ac *AtomicConfig) Update(cfg Config) {
	ac.config.Store(cfg)
}

// AtomicFlag demonstrates a simple atomic boolean flag
type AtomicFlag struct {
	flag int32
}

// Set sets the flag to true
func (f *AtomicFlag) Set() {
	atomic.StoreInt32(&f.flag, 1)
}

// Clear sets the flag to false
func (f *AtomicFlag) Clear() {
	atomic.StoreInt32(&f.flag, 0)
}

// IsSet returns true if the flag is set
func (f *AtomicFlag) IsSet() bool {
	return atomic.LoadInt32(&f.flag) == 1
}

// Toggle atomically toggles the flag and returns the new state
func (f *AtomicFlag) Toggle() bool {
	for {
		old := atomic.LoadInt32(&f.flag)
		new := int32(1) - old
		if atomic.CompareAndSwapInt32(&f.flag, old, new) {
			return new == 1
		}
	}
}

// ReferenceCounter demonstrates atomic reference counting
type ReferenceCounter struct {
	refs  int32
	onZero func()
}

// NewReferenceCounter creates a reference counter with initial count of 1
func NewReferenceCounter(onZero func()) *ReferenceCounter {
	return &ReferenceCounter{
		refs:   1,
		onZero: onZero,
	}
}

// Acquire increments the reference count
func (rc *ReferenceCounter) Acquire() {
	atomic.AddInt32(&rc.refs, 1)
}

// Release decrements the reference count and calls onZero if it reaches 0
func (rc *ReferenceCounter) Release() {
	if atomic.AddInt32(&rc.refs, -1) == 0 {
		if rc.onZero != nil {
			rc.onZero()
		}
	}
}

// Count returns the current reference count
func (rc *ReferenceCounter) Count() int32 {
	return atomic.LoadInt32(&rc.refs)
}

// SpinLock demonstrates a simple spin lock using atomic operations
type SpinLock struct {
	state int32
}

// Lock acquires the spin lock
func (sl *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&sl.state, 0, 1) {
		// Busy-wait (spin)
		// In production, you might want to add runtime.Gosched() here
	}
}

// Unlock releases the spin lock
func (sl *SpinLock) Unlock() {
	atomic.StoreInt32(&sl.state, 0)
}

// TryLock attempts to acquire the lock without blocking
func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapInt32(&sl.state, 0, 1)
}

// Metrics demonstrates concurrent metrics collection using atomic operations
type Metrics struct {
	requests   int64
	errors     int64
	totalBytes int64
}

// RecordRequest atomically increments the request counter
func (m *Metrics) RecordRequest() {
	atomic.AddInt64(&m.requests, 1)
}

// RecordError atomically increments the error counter
func (m *Metrics) RecordError() {
	atomic.AddInt64(&m.errors, 1)
}

// RecordBytes atomically adds to the total bytes counter
func (m *Metrics) RecordBytes(bytes int64) {
	atomic.AddInt64(&m.totalBytes, bytes)
}

// GetSnapshot returns a snapshot of current metrics
func (m *Metrics) GetSnapshot() (requests, errors, totalBytes int64) {
	requests = atomic.LoadInt64(&m.requests)
	errors = atomic.LoadInt64(&m.errors)
	totalBytes = atomic.LoadInt64(&m.totalBytes)
	return
}

// Reset atomically resets all metrics to zero
func (m *Metrics) Reset() {
	atomic.StoreInt64(&m.requests, 0)
	atomic.StoreInt64(&m.errors, 0)
	atomic.StoreInt64(&m.totalBytes, 0)
}

// Worker demonstrates using atomic operations for worker coordination
type Worker struct {
	running    int32
	processed  int64
	workQueue  chan int
	stopSignal chan struct{}
}

// NewWorker creates a new worker
func NewWorker() *Worker {
	return &Worker{
		workQueue:  make(chan int, 100),
		stopSignal: make(chan struct{}),
	}
}

// Start starts the worker
func (w *Worker) Start() {
	if atomic.CompareAndSwapInt32(&w.running, 0, 1) {
		go w.run()
	}
}

// Stop stops the worker
func (w *Worker) Stop() {
	if atomic.CompareAndSwapInt32(&w.running, 1, 0) {
		close(w.stopSignal)
	}
}

// IsRunning returns true if the worker is running
func (w *Worker) IsRunning() bool {
	return atomic.LoadInt32(&w.running) == 1
}

// ProcessedCount returns the number of processed items
func (w *Worker) ProcessedCount() int64 {
	return atomic.LoadInt64(&w.processed)
}

// Submit submits work to the worker
func (w *Worker) Submit(work int) {
	if w.IsRunning() {
		select {
		case w.workQueue <- work:
		default:
			// Queue full, drop work
		}
	}
}

func (w *Worker) run() {
	for {
		select {
		case work := <-w.workQueue:
			// Process work
			_ = work
			atomic.AddInt64(&w.processed, 1)
		case <-w.stopSignal:
			return
		}
	}
}

// SafeMap demonstrates atomic operations with sync.Map for thread-safe map access
type SafeMap struct {
	m    sync.Map
	size int64
}

// Set stores a key-value pair
func (sm *SafeMap) Set(key string, value interface{}) {
	// Check if key exists first
	_, exists := sm.m.Load(key)
	// Store the value (overwrites if exists)
	sm.m.Store(key, value)
	// Only increment size if key didn't exist before
	if !exists {
		atomic.AddInt64(&sm.size, 1)
	}
}

// Get retrieves a value by key
func (sm *SafeMap) Get(key string) (interface{}, bool) {
	return sm.m.Load(key)
}

// Delete removes a key
func (sm *SafeMap) Delete(key string) {
	_, loaded := sm.m.LoadAndDelete(key)
	if loaded {
		atomic.AddInt64(&sm.size, -1)
	}
}

// Size returns the approximate size of the map
func (sm *SafeMap) Size() int64 {
	return atomic.LoadInt64(&sm.size)
}
