# Go Concurrency and Testing

Comprehensive documentation and examples for Go concurrency primitives, synchronization patterns, and testing frameworks. This collection covers mutexes, atomic operations, and the Gomega testing library with detailed explanations, working code, and best practices.

## 📚 Contents

### Core Documentation

#### Synchronization Primitives
- **[Mutexes (`sync.Mutex`, `sync.RWMutex`)](mutexes.md)** - Complete guide to mutual exclusion
  - What mutexes are and why we need them
  - Low-level implementation (normal mode, starvation mode, fast/slow paths)
  - Basic usage patterns and RWMutex
  - Best practices and common pitfalls
  - Performance considerations and advanced patterns

- **[Atomic Operations (`sync/atomic`)](packages/atomic.md)** - Lock-free synchronization primitives
  - Hardware-level atomic operations
  - CPU instruction usage (LOCK, CMPXCHG, MFENCE)
  - Memory model and sequential consistency
  - Lock-free data structures
  - Performance characteristics

#### Testing Framework
- **[Gomega](packages/gomega.md)** - Expressive matcher/assertion library
  - 50+ built-in matchers
  - Async testing support (Eventually, Consistently)
  - Custom matcher implementation
  - Integration with Go testing and Ginkgo

### 💻 Working Examples

#### Mutex Examples (`examples-mutexes/`)
- **[basic_mutex.go](examples-mutexes/basic_mutex.go)** - Counter examples and basic patterns
- **[rwmutex.go](examples-mutexes/rwmutex.go)** - Read-write locks for cache and statistics
- **[pitfalls.go](examples-mutexes/pitfalls.go)** - Deadlocks, missing unlocks, lock contention
- **[advanced.go](examples-mutexes/advanced.go)** - sync.Once, try-lock, sync.Cond patterns

#### Atomic Examples (`examples/`)
- **[atomic_examples.go](examples/atomic_examples.go)** - Production-ready atomic implementations
- **[atomic_examples_test.go](examples/atomic_examples_test.go)** - Test suite using Gomega matchers

## 🚀 Quick Start

### Running Mutex Examples

```bash
# Navigate to Go directory
cd go

# Run mutex examples
go run examples-mutexes/basic_mutex.go
go run examples-mutexes/rwmutex.go
go run examples-mutexes/pitfalls.go
go run examples-mutexes/advanced.go

# Detect race conditions
go run -race examples-mutexes/basic_mutex.go
```

### Running Atomic Examples and Tests

```bash
# Install dependencies
go mod tidy

# Run all tests
go test -v ./examples/

# Run specific tests
go test -v ./examples/ -run TestAtomicCounter
go test -v ./examples/ -run TestWorker

# Run with race detector
go test -race ./examples/
```

## 🎯 Topic Overview

### Mutexes - Mutual Exclusion

**Purpose:** Protect shared state from concurrent access by multiple goroutines.

**When to Use:**
- ✅ Protecting shared mutable state
- ✅ Multiple goroutines access the same data structure
- ✅ Operations involve multiple variables
- ✅ Complex state transitions

**When NOT to Use:**
- ❌ Single atomic variable (use `sync/atomic`)
- ❌ Passing data between goroutines (use channels)
- ❌ Data is immutable
- ❌ Already synchronized types (`sync.Map`, `sync.Pool`)

**Key Takeaways:**
- Always use `defer` to unlock mutexes
- Keep critical sections as small as possible
- Never copy mutexes (use pointers, run `go vet`)
- Beware of deadlocks from inconsistent lock ordering
- Use RWMutex when reads greatly outnumber writes

### Atomic Operations - Lock-Free Synchronization

**Purpose:** Lock-free synchronization primitives using hardware-level CPU instructions.

**When to Use:**
- ✅ Simple counters and flags
- ✅ Performance-critical paths
- ✅ Lock-free data structures
- ✅ Configuration hot-reload

**Key Strengths:**
- Hardware-level performance (1-10 ns operations)
- Zero allocation
- Sequential consistency guarantees
- Simple, clear API

**Best Practices:**
- Always use atomic operations consistently
- Consider cache line padding to prevent false sharing
- Use atomic.Value for complex types
- Profile before optimizing

### Gomega - Expressive Testing

**Purpose:** Matcher/assertion library for writing readable, expressive tests.

**Key Features:**
- 50+ built-in matchers
- First-class async testing support
- Composable matcher system
- Custom matcher extensibility
- Works with Go's testing package and Ginkgo

**Best Practices:**
- Use NewWithT for proper test integration
- Choose specific matchers over generic ones
- Use Eventually for async operations
- Set reasonable timeouts

## 📊 Comparison: Mutex vs Atomic vs Channels

| Scenario | Best Choice | Reasoning |
|----------|-------------|-----------|
| Simple counter | Atomic | Fastest, simplest |
| Multiple related variables | Mutex | Atomic updates across fields |
| Passing data between goroutines | Channels | Communicates intent |
| Read-heavy workload | RWMutex | Allows concurrent reads |
| Lock-free stack/queue | Atomic (CAS) | No blocking |
| Complex state machine | Mutex | Easier to reason about |

## 📖 Learning Path

### Beginner Path

1. **Understand the basics**
   - Read [mutexes.md](mutexes.md) introduction
   - Run [basic_mutex.go](examples-mutexes/basic_mutex.go) with `-race`
   - Study race conditions and how mutexes fix them

2. **Explore RWMutex**
   - Study [rwmutex.go](examples-mutexes/rwmutex.go)
   - Understand read vs write locks
   - Learn when to use RWMutex over Mutex

3. **Learn common pitfalls**
   - Review [pitfalls.go](examples-mutexes/pitfalls.go)
   - Understand deadlocks, missing unlocks
   - Practice with `go vet` and `-race`

### Intermediate Path

4. **Master atomic operations**
   - Read [atomic.md](packages/atomic.md)
   - Study hardware-level implementation
   - Run [atomic_examples.go](examples/atomic_examples.go)

5. **Advanced mutex patterns**
   - Study [advanced.go](examples-mutexes/advanced.go)
   - Learn sync.Once, sync.Cond
   - Compare double-checked locking vs alternatives

6. **Write better tests**
   - Read [gomega.md](packages/gomega.md)
   - Study [atomic_examples_test.go](examples/atomic_examples_test.go)
   - Practice writing async tests with Eventually

### Advanced Path

7. **Implement lock-free data structures**
   - Use CAS loops for stacks/queues
   - Apply atomic operations to real problems
   - Benchmark against mutex-based implementations

8. **Optimize concurrent code**
   - Profile with `-mutexprofile`
   - Apply sharding techniques
   - Reduce lock contention

9. **Build concurrent systems**
   - Combine mutexes, atomics, and channels
   - Apply patterns from examples
   - Write comprehensive tests

## 🛠️ Tools and Commands

```bash
# Detect race conditions
go run -race main.go
go test -race ./...

# Detect common mistakes (including mutex copying)
go vet ./...

# Profile lock contention
go test -bench=. -mutexprofile=mutex.out
go tool pprof mutex.out

# Format code
go fmt ./...

# Run tests with coverage
go test -cover ./...
```

## 📚 Additional Resources

### Official Documentation
- [Go Memory Model](https://go.dev/ref/mem)
- [sync package](https://pkg.go.dev/sync)
- [sync/atomic package](https://pkg.go.dev/sync/atomic)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)

### Articles and Guides
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [Go 101: Atomic Operations](https://go101.org/article/atomic.html)

### Testing Resources
- [Gomega Documentation](https://onsi.github.io/gomega/)
- [Ginkgo + Gomega Guide](https://onsi.github.io/ginkgo/)

## ✅ Test Coverage

All examples include comprehensive testing:

### Mutex Examples
- ✅ Race condition detection
- ✅ Deadlock prevention examples
- ✅ RWMutex performance comparisons
- ✅ Advanced pattern demonstrations

### Atomic Examples
- ✅ 14 test functions covering all operations
- ✅ Concurrency tests (100 goroutines)
- ✅ Async worker behavior validation
- ✅ Race detector clean
- ✅ Gomega matcher examples

## 🎓 Summary

### Quick Reference

| Primitive | Use Case | Pros | Cons |
|-----------|----------|------|------|
| `sync.Mutex` | Simple mutual exclusion | Easy to use, safe | Can block, overhead |
| `sync.RWMutex` | Read-heavy workloads | Concurrent reads | More complex, overhead |
| `sync.Once` | One-time initialization | Thread-safe, efficient | Single use only |
| `sync.Cond` | Complex waiting conditions | Flexible signaling | Complex to use |
| `sync/atomic` | Simple counters/flags | Fastest, lock-free | Limited to simple ops |
| Channels | Communication | Natural Go idiom | Memory allocation |

### Key Principles

1. **"Don't communicate by sharing memory; share memory by communicating"** - Prefer channels when appropriate
2. **Always use `defer` to unlock** - Prevents missing unlocks
3. **Keep critical sections small** - Minimize lock holding time
4. **Never copy mutexes** - Use pointer receivers
5. **Profile before optimizing** - Measure, don't guess
6. **Use the right tool** - Mutexes, atomics, and channels each have their place

## 🤝 Contributing

This is a learning repository. Contributions welcome:
- 📝 Improve documentation
- 💡 Add more examples
- 🐛 Fix inaccuracies
- ✅ Add test cases
- 📊 Add benchmarks

## 📝 License

Educational purposes. Free to use and share.

---

*Master Go concurrency through understanding, practice, and proper testing.* 🚀
