// Package main demonstrates common mutex pitfalls and how to avoid them
package main

import (
	"fmt"
	"sync"
	"time"
)

// Example 1: Deadlock from circular lock ordering
type Account struct {
	mu      sync.Mutex
	id      int
	balance int
}

func badTransfer(from, to *Account, amount int) {
	// BAD: Can cause deadlock!
	from.mu.Lock()
	defer from.mu.Unlock()
	
	time.Sleep(10 * time.Millisecond) // Simulate work
	
	to.mu.Lock()
	defer to.mu.Unlock()
	
	from.balance -= amount
	to.balance += amount
}

func goodTransfer(from, to *Account, amount int) {
	// GOOD: Consistent lock ordering
	first, second := from, to
	if from.id > to.id {
		first, second = to, from
	}
	
	first.mu.Lock()
	defer first.mu.Unlock()
	
	second.mu.Lock()
	defer second.mu.Unlock()
	
	from.balance -= amount
	to.balance += amount
}

func demonstrateDeadlock() {
	fmt.Println("=== Deadlock Prevention ===")
	acc1 := &Account{id: 1, balance: 1000}
	acc2 := &Account{id: 2, balance: 1000}
	
	fmt.Println("Attempting transfers with proper lock ordering...")
	var wg sync.WaitGroup
	
	// These won't deadlock due to consistent ordering
	wg.Add(2)
	go func() {
		defer wg.Done()
		goodTransfer(acc1, acc2, 100)
		fmt.Println("Transfer 1: acc1 -> acc2 completed")
	}()
	
	go func() {
		defer wg.Done()
		goodTransfer(acc2, acc1, 50)
		fmt.Println("Transfer 2: acc2 -> acc1 completed")
	}()
	
	wg.Wait()
	fmt.Printf("Final balances: acc1=%d, acc2=%d\n", acc1.balance, acc2.balance)
	fmt.Println("✓ No deadlock occurred due to consistent lock ordering")
	
	fmt.Println("\nNote: badTransfer() is commented out to prevent actual deadlock")
	fmt.Println("Uncommenting it would cause the program to hang!")
}

// Example 2: Forgetting to unlock
type BadCounter struct {
	mu    sync.Mutex
	value int
}

// BAD: Easy to forget unlock in error paths
func (c *BadCounter) IncrementBad() error {
	c.mu.Lock()
	// Missing defer unlock!
	
	if c.value >= 100 {
		return fmt.Errorf("max reached") // BUG: Mutex never unlocked!
	}
	
	c.value++
	c.mu.Unlock() // Only unlocked on success path
	return nil
}

// GOOD: Always use defer
func (c *BadCounter) IncrementGood() error {
	c.mu.Lock()
	defer c.mu.Unlock() // Always unlocked
	
	if c.value >= 100 {
		return fmt.Errorf("max reached")
	}
	
	c.value++
	return nil
}

func demonstrateMissingUnlock() {
	fmt.Println("\n=== Missing Unlock Pitfall ===")
	fmt.Println("Always use 'defer' to ensure unlock happens")
	
	counter := &BadCounter{}
	
	// Use good version
	for i := 0; i < 10; i++ {
		counter.IncrementGood()
	}
	
	fmt.Printf("Counter value: %d\n", counter.value)
	fmt.Println("✓ Using defer ensures mutex is always unlocked")
}

// Example 3: Holding lock during blocking operation
type BadCache struct {
	mu   sync.Mutex
	data map[string]string
}

func (c *BadCache) GetOrFetchBad(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if val, ok := c.data[key]; ok {
		return val
	}
	
	// BAD: Blocking operation while holding lock!
	time.Sleep(100 * time.Millisecond) // Simulates network call
	val := "fetched-" + key
	c.data[key] = val
	return val
}

type GoodCache struct {
	mu   sync.Mutex
	data map[string]string
}

func (c *GoodCache) GetOrFetchGood(key string) string {
	c.mu.Lock()
	if val, ok := c.data[key]; ok {
		c.mu.Unlock()
		return val
	}
	c.mu.Unlock()
	
	// GOOD: Blocking operation without lock
	time.Sleep(100 * time.Millisecond) // Simulates network call
	val := "fetched-" + key
	
	c.mu.Lock()
	c.data[key] = val
	c.mu.Unlock()
	
	return val
}

func demonstrateBlockingWithLock() {
	fmt.Println("\n=== Blocking Operation with Lock ===")
	
	badCache := &BadCache{data: make(map[string]string)}
	goodCache := &GoodCache{data: make(map[string]string)}
	
	// Bad version - serializes all operations
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			badCache.GetOrFetchBad(fmt.Sprintf("key%d", id))
		}(i)
	}
	wg.Wait()
	badDuration := time.Since(start)
	fmt.Printf("Bad cache (lock during fetch): %v\n", badDuration)
	
	// Good version - allows concurrent fetches
	start = time.Now()
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			goodCache.GetOrFetchGood(fmt.Sprintf("key%d", id))
		}(i)
	}
	wg.Wait()
	goodDuration := time.Since(start)
	fmt.Printf("Good cache (no lock during fetch): %v\n", goodDuration)
	
	fmt.Println("✓ Minimize critical sections - don't hold locks during blocking operations")
}

// Example 4: Copying mutexes
type CopyableBad struct {
	mu    sync.Mutex
	value int
}

// BAD: Value receiver copies the mutex!
func (c CopyableBad) IncrementBad() {
	c.mu.Lock() // Locks a COPY of the mutex
	defer c.mu.Unlock()
	c.value++ // Increments a COPY of value
}

// GOOD: Pointer receiver
func (c *CopyableBad) IncrementGood() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func demonstrateCopying() {
	fmt.Println("\n=== Copying Mutex Pitfall ===")
	
	counter := &CopyableBad{}
	
	// Use good version (pointer receiver)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.IncrementGood()
		}()
	}
	wg.Wait()
	
	fmt.Printf("Counter value with pointer receiver: %d\n", counter.value)
	fmt.Println("✓ Always use pointer receivers for types with mutexes")
	fmt.Println("Run 'go vet' to detect mutex copying issues")
}

// Example 5: Lock contention
type HighContentionCounter struct {
	mu    sync.Mutex
	value int64
}

func (c *HighContentionCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

type ShardedCounter struct {
	shards [16]struct {
		mu    sync.Mutex
		value int64
	}
}

func (c *ShardedCounter) Increment(id int) {
	shard := &c.shards[id%16]
	shard.mu.Lock()
	shard.value++
	shard.mu.Unlock()
}

func (c *ShardedCounter) Total() int64 {
	var total int64
	for i := range c.shards {
		c.shards[i].mu.Lock()
		total += c.shards[i].value
		c.shards[i].mu.Unlock()
	}
	return total
}

func demonstrateContention() {
	fmt.Println("\n=== Lock Contention and Sharding ===")
	
	// High contention
	highContention := &HighContentionCounter{}
	start := time.Now()
	var wg sync.WaitGroup
	
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				highContention.Increment()
			}
		}()
	}
	wg.Wait()
	highDuration := time.Since(start)
	
	// Low contention (sharded)
	sharded := &ShardedCounter{}
	start = time.Now()
	
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				sharded.Increment(id)
			}
		}(i)
	}
	wg.Wait()
	shardedDuration := time.Since(start)
	
	fmt.Printf("High contention (single mutex): %v\n", highDuration)
	fmt.Printf("Low contention (sharded):       %v\n", shardedDuration)
	
	if shardedDuration < highDuration {
		speedup := float64(highDuration) / float64(shardedDuration)
		fmt.Printf("Sharding provides %.2fx speedup\n", speedup)
	}
	
	fmt.Printf("Total count: %d\n", sharded.Total())
	fmt.Println("✓ Sharding reduces lock contention and improves performance")
}

func main() {
	fmt.Println("Common Mutex Pitfalls in Go")
	fmt.Println("============================\n")
	
	demonstrateDeadlock()
	demonstrateMissingUnlock()
	demonstrateBlockingWithLock()
	demonstrateCopying()
	demonstrateContention()
	
	fmt.Println("✓ All pitfall examples completed successfully!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("1. Use consistent lock ordering to prevent deadlocks")
	fmt.Println("2. Always use 'defer' to unlock mutexes")
	fmt.Println("3. Don't hold locks during blocking operations")
	fmt.Println("4. Never copy mutexes - use pointer receivers")
	fmt.Println("5. Use sharding to reduce lock contention")
}
