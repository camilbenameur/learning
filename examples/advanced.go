// Package main demonstrates advanced mutex patterns in Go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Example 1: sync.Once for lazy initialization
type HeavyResource struct {
	data string
}

func NewHeavyResource() *HeavyResource {
	fmt.Println("  Initializing heavy resource (expensive operation)...")
	time.Sleep(100 * time.Millisecond) // Simulate expensive initialization
	return &HeavyResource{data: "initialized"}
}

type ServiceWithOnce struct {
	once     sync.Once
	resource *HeavyResource
}

func (s *ServiceWithOnce) GetResource() *HeavyResource {
	s.once.Do(func() {
		s.resource = NewHeavyResource()
	})
	return s.resource
}

func demonstrateSyncOnce() {
	fmt.Println("=== sync.Once for Lazy Initialization ===")
	service := &ServiceWithOnce{}
	var wg sync.WaitGroup
	
	// Multiple goroutines try to get resource
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			resource := service.GetResource()
			fmt.Printf("Goroutine %d got resource: %s\n", id, resource.data)
		}(i)
	}
	
	wg.Wait()
	fmt.Println("✓ Resource initialized exactly once, even with concurrent access")
}

// Example 2: Try-lock pattern
type TryMutex struct {
	ch chan struct{}
}

func NewTryMutex() *TryMutex {
	ch := make(chan struct{}, 1)
	ch <- struct{}{} // Initialize as unlocked
	return &TryMutex{ch: ch}
}

func (m *TryMutex) Lock() {
	<-m.ch
}

func (m *TryMutex) Unlock() {
	m.ch <- struct{}{}
}

func (m *TryMutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
		return false
	}
}

func demonstrateTryLock() {
	fmt.Println("\n=== Try-Lock Pattern ===")
	mutex := NewTryMutex()
	
	// Lock it
	mutex.Lock()
	fmt.Println("Lock acquired")
	
	// Try to lock (should fail)
	if mutex.TryLock() {
		fmt.Println("TryLock succeeded")
		mutex.Unlock()
	} else {
		fmt.Println("TryLock failed - mutex already locked")
	}
	
	// Unlock
	mutex.Unlock()
	fmt.Println("Lock released")
	
	// Try to lock again (should succeed)
	if mutex.TryLock() {
		fmt.Println("TryLock succeeded after unlock")
		mutex.Unlock()
	}
	
	fmt.Println("✓ Try-lock allows non-blocking lock attempts")
}

// Example 3: Conditional variables with sync.Cond
type Queue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []int
}

func NewQueue() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Enqueue(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	q.items = append(q.items, item)
	fmt.Printf("Enqueued %d, queue size: %d\n", item, len(q.items))
	q.cond.Signal() // Wake one waiting goroutine
}

func (q *Queue) Dequeue() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	for len(q.items) == 0 {
		fmt.Println("Queue empty, waiting...")
		q.cond.Wait() // Atomically unlocks mu and waits
	}
	
	item := q.items[0]
	q.items = q.items[1:]
	fmt.Printf("Dequeued %d, queue size: %d\n", item, len(q.items))
	return item
}

func demonstrateSyncCond() {
	fmt.Println("\n=== sync.Cond for Producer-Consumer ===")
	queue := NewQueue()
	var wg sync.WaitGroup
	
	// Start consumers first
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			item := queue.Dequeue()
			fmt.Printf("Consumer %d received: %d\n", id, item)
		}(i)
	}
	
	// Give consumers time to start waiting
	time.Sleep(100 * time.Millisecond)
	
	// Start producer
	for i := 0; i < 3; i++ {
		queue.Enqueue(i * 10)
		time.Sleep(50 * time.Millisecond)
	}
	
	wg.Wait()
	fmt.Println("✓ sync.Cond enables efficient waiting for conditions")
}

// Example 4: Mutex vs Atomic operations
type MutexCounter struct {
	mu    sync.Mutex
	value int64
}

func (c *MutexCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *MutexCounter) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func demonstrateMutexVsAtomic() {
	fmt.Println("\n=== Mutex vs Atomic Operations ===")
	
	const iterations = 100000
	const goroutines = 10
	
	// Mutex version
	mutexCounter := &MutexCounter{}
	start := time.Now()
	var wg sync.WaitGroup
	
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				mutexCounter.Increment()
			}
		}()
	}
	wg.Wait()
	mutexDuration := time.Since(start)
	
	// Atomic version
	atomicCounter := &AtomicCounter{}
	start = time.Now()
	
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				atomicCounter.Increment()
			}
		}()
	}
	wg.Wait()
	atomicDuration := time.Since(start)
	
	fmt.Printf("Mutex counter:  %v (value: %d)\n", mutexDuration, mutexCounter.Value())
	fmt.Printf("Atomic counter: %v (value: %d)\n", atomicDuration, atomicCounter.Value())
	
	if atomicDuration < mutexDuration {
		speedup := float64(mutexDuration) / float64(atomicDuration)
		fmt.Printf("Atomic operations are %.2fx faster for simple counters\n", speedup)
	}
	
	fmt.Println("✓ Use atomics for simple single-variable operations")
}

// Example 5: Double-checked locking vs sync.Once
type DoubleCheckedService struct {
	mu       sync.Mutex
	instance *HeavyResource
}

func (s *DoubleCheckedService) GetInstance() *HeavyResource {
	if s.instance != nil { // First check (no lock)
		return s.instance
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.instance == nil { // Second check (with lock)
		s.instance = NewHeavyResource()
	}
	
	return s.instance
}

func compareLazyInit() {
	fmt.Println("\n=== Double-Checked Locking vs sync.Once ===")
	
	// Double-checked locking
	dcService := &DoubleCheckedService{}
	start := time.Now()
	var wg sync.WaitGroup
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dcService.GetInstance()
		}()
	}
	wg.Wait()
	dcDuration := time.Since(start)
	
	// sync.Once
	onceService := &ServiceWithOnce{}
	start = time.Now()
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			onceService.GetResource()
		}()
	}
	wg.Wait()
	onceDuration := time.Since(start)
	
	fmt.Printf("Double-checked locking: %v\n", dcDuration)
	fmt.Printf("sync.Once:              %v\n", onceDuration)
	fmt.Println("✓ sync.Once is simpler and often more efficient than double-checked locking")
}

// Example 6: Read-copy-update pattern
type Config struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *Config) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

// Read-copy-update: copy map, modify copy, swap pointer
func (c *Config) Update(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Create a copy
	newData := make(map[string]string, len(c.data)+1)
	for k, v := range c.data {
		newData[k] = v
	}
	newData[key] = value
	
	// Atomically swap
	c.data = newData
}

func demonstrateRCU() {
	fmt.Println("\n=== Read-Copy-Update Pattern ===")
	config := &Config{data: make(map[string]string)}
	config.data["initial"] = "value"
	
	var wg sync.WaitGroup
	
	// Many readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				config.Get("initial")
			}
		}()
	}
	
	// Few writers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			config.Update(fmt.Sprintf("key%d", id), fmt.Sprintf("value%d", id))
		}(i)
	}
	
	wg.Wait()
	
	fmt.Println("✓ RCU allows lock-free reads at the cost of copying on writes")
	fmt.Println("  Best for read-heavy workloads with infrequent updates")
}

func main() {
	fmt.Println("Advanced Mutex Patterns in Go")
	fmt.Println("==============================\n")
	
	demonstrateSyncOnce()
	demonstrateTryLock()
	demonstrateSyncCond()
	demonstrateMutexVsAtomic()
	compareLazyInit()
	demonstrateRCU()
	
	fmt.Println("✓ All advanced pattern examples completed successfully!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("1. sync.Once: Simplest way for one-time initialization")
	fmt.Println("2. Try-lock: Non-blocking lock attempts")
	fmt.Println("3. sync.Cond: Efficient waiting for conditions")
	fmt.Println("4. Atomics: Faster than mutexes for simple operations")
	fmt.Println("5. RCU: Lock-free reads for read-heavy workloads")
}
