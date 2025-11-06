# Go Atomic Package (sync/atomic)

## Overview

The `sync/atomic` package provides low-level atomic memory primitives useful for implementing synchronization algorithms. It offers atomic operations on integers and pointers without requiring locks, making concurrent operations faster and more efficient than mutex-based synchronization.

## Purpose

The atomic package is designed for:
- **Lock-free synchronization**: Perform operations without mutex overhead
- **Memory safety**: Ensure thread-safe operations on shared variables
- **Performance-critical code**: Optimize concurrent access to counters, flags, and pointers
- **Building higher-level synchronization primitives**: Foundation for more complex concurrent data structures

## Key Strengths

### 1. **Lock-Free Operations**
Atomic operations complete without acquiring locks, avoiding contention and context switching overhead.

### 2. **Hardware-Level Support**
Atomic operations are implemented using CPU-level instructions (like CMPXCHG on x86, LDREX/STREX on ARM), providing the fastest possible synchronization.

### 3. **Memory Ordering Guarantees**
All atomic operations have sequential consistency, meaning:
- Operations appear to occur in a single, global order
- All goroutines see operations in the same sequence
- No reordering of atomic operations

### 4. **Simple API**
Clear function names and straightforward usage patterns make atomic operations accessible.

### 5. **Zero Allocation**
Atomic operations don't allocate memory, making them ideal for hot paths.

## Low-Level Mechanics

### Memory Model

Go's atomic operations provide **sequential consistency**, the strongest memory ordering guarantee:

```
Thread 1:  atomic.Store(x, 1)  →  atomic.Store(y, 1)
Thread 2:  r1 = atomic.Load(y)  →  r2 = atomic.Load(x)

If r1 == 1, then r2 must also == 1
```

### CPU Instructions

Under the hood, atomic operations map to CPU-specific instructions:

**x86/x64:**
- `LOCK ADD` - Atomic addition
- `LOCK CMPXCHG` - Compare-and-swap
- `LOCK XCHG` - Atomic exchange
- `MFENCE` - Memory barrier

**ARM:**
- `LDREX/STREX` - Load/Store exclusive (for CAS loops)
- `DMB` - Data Memory Barrier
- `LDADD` (ARMv8.1+) - Atomic add

**Compiler Integration:**
The Go compiler recognizes atomic operations and:
1. Prevents reordering of atomic operations
2. Ensures proper memory barriers are inserted
3. Uses the most efficient CPU instructions available

### Cache Coherency

Atomic operations interact with CPU cache coherency protocols (like MESI):
1. **Modified (M)**: Cache line modified in one core
2. **Exclusive (E)**: Cache line in one core, matches memory
3. **Shared (S)**: Cache line in multiple cores
4. **Invalid (I)**: Cache line is invalid

Atomic operations force cache synchronization across all cores, ensuring consistency.

## Core Functions

### Integer Operations

```go
// Add operations
func AddInt32(addr *int32, delta int32) (new int32)
func AddInt64(addr *int64, delta int64) (new int64)
func AddUint32(addr *uint32, delta uint32) (new uint32)
func AddUint64(addr *uint64, delta uint64) (new uint64)

// Load operations
func LoadInt32(addr *int32) (val int32)
func LoadInt64(addr *int64) (val int64)
func LoadUint32(addr *uint32) (val uint32)
func LoadUint64(addr *uint64) (val uint64)

// Store operations
func StoreInt32(addr *int32, val int32)
func StoreInt64(addr *int64, val int64)
func StoreUint32(addr *uint32, val uint32)
func StoreUint64(addr *uint64, val uint64)

// Swap operations
func SwapInt32(addr *int32, new int32) (old int32)
func SwapInt64(addr *int64, new int64) (old int64)
func SwapUint32(addr *uint32, new uint32) (old uint32)
func SwapUint64(addr *uint64, new uint64) (old uint64)

// Compare-and-swap operations
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)
```

### Pointer Operations

```go
func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
```

### Value Type (Go 1.4+)

```go
type Value struct {
    // contains filtered or unexported fields
}

func (v *Value) Load() (val any)
func (v *Value) Store(val any)
func (v *Value) Swap(new any) (old any)
func (v *Value) CompareAndSwap(old, new any) (swapped bool)
```

## Usage Examples

### Example 1: Atomic Counter

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type Counter struct {
    value int64
}

func (c *Counter) Increment() {
    atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Decrement() {
    atomic.AddInt64(&c.value, -1)
}

func (c *Counter) Get() int64 {
    return atomic.LoadInt64(&c.value)
}

func main() {
    var counter Counter
    var wg sync.WaitGroup

    // Spawn 100 goroutines, each incrementing 1000 times
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                counter.Increment()
            }
        }()
    }

    wg.Wait()
    fmt.Printf("Final count: %d\n", counter.Get()) // Output: Final count: 100000
}
```

### Example 2: Compare-and-Swap Loop

```go
package main

import (
    "fmt"
    "sync/atomic"
)

// Thread-safe stack using CAS
type Stack struct {
    head unsafe.Pointer
}

type Node struct {
    value int
    next  unsafe.Pointer
}

func (s *Stack) Push(value int) {
    node := &Node{value: value}
    for {
        old := atomic.LoadPointer(&s.head)
        node.next = old
        if atomic.CompareAndSwapPointer(&s.head, old, unsafe.Pointer(node)) {
            return
        }
        // CAS failed, retry
    }
}

func (s *Stack) Pop() (int, bool) {
    for {
        old := atomic.LoadPointer(&s.head)
        if old == nil {
            return 0, false
        }
        node := (*Node)(old)
        if atomic.CompareAndSwapPointer(&s.head, old, node.next) {
            return node.value, true
        }
        // CAS failed, retry
    }
}
```

### Example 3: Atomic Value for Configuration

```go
package main

import (
    "fmt"
    "sync/atomic"
    "time"
)

type Config struct {
    Timeout time.Duration
    MaxRetries int
}

type Server struct {
    config atomic.Value
}

func (s *Server) UpdateConfig(cfg Config) {
    s.config.Store(cfg)
}

func (s *Server) GetConfig() Config {
    return s.config.Load().(Config)
}

func main() {
    server := &Server{}
    
    // Initial configuration
    server.UpdateConfig(Config{
        Timeout: 5 * time.Second,
        MaxRetries: 3,
    })
    
    // Configuration can be read by multiple goroutines safely
    for i := 0; i < 10; i++ {
        go func() {
            cfg := server.GetConfig()
            fmt.Printf("Timeout: %v, MaxRetries: %d\n", cfg.Timeout, cfg.MaxRetries)
        }()
    }
    
    time.Sleep(100 * time.Millisecond)
    
    // Update configuration atomically
    server.UpdateConfig(Config{
        Timeout: 10 * time.Second,
        MaxRetries: 5,
    })
}
```

### Example 4: Spin Lock Implementation

```go
package main

import (
    "runtime"
    "sync/atomic"
)

type SpinLock struct {
    state int32
}

func (s *SpinLock) Lock() {
    for !atomic.CompareAndSwapInt32(&s.state, 0, 1) {
        // Spin: yield to other goroutines
        runtime.Gosched()
    }
}

func (s *SpinLock) Unlock() {
    atomic.StoreInt32(&s.state, 0)
}
```

### Example 5: Flag Synchronization

```go
package main

import (
    "fmt"
    "sync/atomic"
    "time"
)

type Worker struct {
    shutdown int32
}

func (w *Worker) Shutdown() {
    atomic.StoreInt32(&w.shutdown, 1)
}

func (w *Worker) IsShutdown() bool {
    return atomic.LoadInt32(&w.shutdown) == 1
}

func (w *Worker) Run() {
    for !w.IsShutdown() {
        // Do work
        fmt.Println("Working...")
        time.Sleep(100 * time.Millisecond)
    }
    fmt.Println("Worker stopped")
}

func main() {
    worker := &Worker{}
    
    go worker.Run()
    
    time.Sleep(500 * time.Millisecond)
    worker.Shutdown()
    
    time.Sleep(200 * time.Millisecond)
}
```

## Common Scenarios

### 1. **Shared Counters**
Use atomic operations for metrics, statistics, and rate counters accessed by multiple goroutines.

### 2. **Reference Counting**
Implement reference counting for resource management without locks.

### 3. **Configuration Hot-Reload**
Use `atomic.Value` to update configuration that's read frequently but updated rarely.

### 4. **Lock-Free Data Structures**
Build stacks, queues, and other concurrent data structures using CAS operations.

### 5. **State Machines**
Implement state transitions with atomic operations to ensure consistency.

### 6. **Resource Pooling**
Manage available/in-use resource counts atomically.

## Performance Considerations

### When to Use Atomic Operations

✅ **Good Use Cases:**
- Simple counters and flags
- Single-value updates
- Read-heavy, write-light scenarios
- Performance-critical paths

❌ **Avoid When:**
- Multiple related values need to update together
- Complex state transitions
- Operations require conditional logic based on multiple variables

### Atomic vs. Mutex

| Aspect | Atomic | Mutex |
|--------|--------|-------|
| Performance | Faster (no context switch) | Slower (potential context switch) |
| Complexity | Simple operations only | Complex operations supported |
| Memory Usage | Zero allocation | Small overhead |
| Fairness | No fairness guarantees | Fair with sync.Mutex |
| Use Case | Single variable operations | Multi-variable critical sections |

## Common Pitfalls

### 1. **Mixing Atomic and Non-Atomic Access**

```go
// BAD: Race condition
var counter int64
atomic.AddInt64(&counter, 1)  // Atomic
counter++                      // Non-atomic - RACE!
```

### 2. **False Sharing**

```go
// BAD: False sharing - counters share cache line
type Counters struct {
    a int64
    b int64
}

// GOOD: Pad to separate cache lines
type Counters struct {
    a int64
    _ [7]int64  // Padding
    b int64
}
```

### 3. **ABA Problem**

```go
// CAS can succeed even if value changed and changed back
// Solution: Use version numbers or generation counters
type VersionedPointer struct {
    ptr     unsafe.Pointer
    version uint64
}
```

### 4. **Store/Load Without Proper Initialization**

```go
// BAD: Loading before first Store
var val atomic.Value
val.Load()  // Panics if never stored

// GOOD: Always Store before Load
var val atomic.Value
val.Store(initialValue)
val.Load()  // Safe
```

## Best Practices

1. **Always use atomic operations consistently** - Never mix atomic and non-atomic access
2. **Use atomic.Value for complex types** - Safer than unsafe.Pointer manipulation
3. **Document atomicity requirements** - Make it clear which fields require atomic access
4. **Consider cache line padding** - Prevent false sharing in hot paths
5. **Profile before optimizing** - Use mutexes unless profiling shows contention
6. **Keep critical sections small** - Minimize the scope of atomic variables

## References

- [Go sync/atomic package documentation](https://pkg.go.dev/sync/atomic)
- [The Go Memory Model](https://go.dev/ref/mem)
- [Atomic Operations in Go](https://go101.org/article/atomic.html)
