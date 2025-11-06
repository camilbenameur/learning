// Package main demonstrates RWMutex usage in Go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Cache demonstrates RWMutex for read-heavy workloads
type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	val, ok := c.data[key]
	return val, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.data[key] = value
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return len(c.data)
}

func demonstrateRWMutex() {
	fmt.Println("=== RWMutex Cache Example ===")
	cache := NewCache()
	var wg sync.WaitGroup
	
	// Pre-populate cache
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	
	// Start 50 readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key%d", j%10)
				if val, ok := cache.Get(key); ok {
					if id == 0 && j == 0 {
						fmt.Printf("Reader %d: Got %s=%s\n", id, key, val)
					}
				}
				time.Sleep(1 * time.Millisecond)
			}
		}(i)
	}
	
	// Start 5 writers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				key := fmt.Sprintf("key%d", j)
				value := fmt.Sprintf("new-value%d-%d", id, j)
				cache.Set(key, value)
				if id == 0 {
					fmt.Printf("Writer %d: Set %s=%s\n", id, key, value)
				}
				time.Sleep(5 * time.Millisecond)
			}
		}(i)
	}
	
	wg.Wait()
	fmt.Printf("Cache size: %d\n", cache.Size())
	fmt.Println("✓ RWMutex allows multiple readers but exclusive writer access")
}

// Comparison of Mutex vs RWMutex performance
type MutexCache struct {
	mu   sync.Mutex
	data map[string]int
}

type RWMutexCache struct {
	mu   sync.RWMutex
	data map[string]int
}

func (c *MutexCache) Get(key string) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *RWMutexCache) Get(key string) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func benchmarkCache(name string, reads, writes int, getFunc func(string) (int, bool)) time.Duration {
	var wg sync.WaitGroup
	start := time.Now()
	
	// Readers
	for i := 0; i < reads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				getFunc("key1")
			}
		}()
	}
	
	wg.Wait()
	duration := time.Since(start)
	
	return duration
}

func demonstratePerformance() {
	fmt.Println("\n=== Performance Comparison ===")
	
	mutexCache := &MutexCache{data: map[string]int{"key1": 42}}
	rwMutexCache := &RWMutexCache{data: map[string]int{"key1": 42}}
	
	// Read-heavy workload
	fmt.Println("Read-heavy workload (100 readers, minimal writes):")
	
	mutexTime := benchmarkCache("Mutex", 100, 0, mutexCache.Get)
	fmt.Printf("  Mutex:    %v\n", mutexTime)
	
	rwMutexTime := benchmarkCache("RWMutex", 100, 0, rwMutexCache.Get)
	fmt.Printf("  RWMutex:  %v\n", rwMutexTime)
	
	if rwMutexTime < mutexTime {
		speedup := float64(mutexTime) / float64(rwMutexTime)
		fmt.Printf("  RWMutex is %.2fx faster for read-heavy workloads\n", speedup)
	}
}

// StatsTracker demonstrates RWMutex for statistics
type StatsTracker struct {
	mu       sync.RWMutex
	requests int64
	errors   int64
	latency  []time.Duration
}

func (s *StatsTracker) RecordRequest(duration time.Duration, isError bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.requests++
	if isError {
		s.errors++
	}
	s.latency = append(s.latency, duration)
}

func (s *StatsTracker) GetStats() (requests, errors int64, avgLatency time.Duration) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	requests = s.requests
	errors = s.errors
	
	if len(s.latency) > 0 {
		var total time.Duration
		for _, lat := range s.latency {
			total += lat
		}
		avgLatency = total / time.Duration(len(s.latency))
	}
	
	return
}

func demonstrateStatsTracker() {
	fmt.Println("\n=== Stats Tracker with RWMutex ===")
	stats := &StatsTracker{}
	var wg sync.WaitGroup
	
	// Simulate requests
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			duration := time.Duration(id%50) * time.Millisecond
			isError := id%10 == 0
			stats.RecordRequest(duration, isError)
		}(i)
	}
	
	// Periodic stats reading
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			requests, errors, avgLatency := stats.GetStats()
			fmt.Printf("Stats snapshot: %d requests, %d errors, avg latency: %v\n", 
				requests, errors, avgLatency)
		}()
	}
	
	wg.Wait()
	
	requests, errors, avgLatency := stats.GetStats()
	fmt.Printf("\nFinal stats: %d requests, %d errors, avg latency: %v\n", 
		requests, errors, avgLatency)
}

func main() {
	fmt.Println("RWMutex Examples in Go")
	fmt.Println("======================\n")
	
	demonstrateRWMutex()
	demonstratePerformance()
	demonstrateStatsTracker()
	
	fmt.Println("✓ All RWMutex examples completed successfully!")
}
