# Go Mutexes Learning Repository

A comprehensive guide to understanding and using mutexes in Go, covering low-level implementation details, best practices, common use cases, and frequent pitfalls.

## 📚 Contents

### Main Documentation
- **[go-mutexes.md](go-mutexes.md)** - Comprehensive guide covering:
  - What mutexes are and why we need them
  - Low-level implementation details (normal mode, starvation mode, fast/slow paths)
  - Basic usage patterns
  - RWMutex (read-write locks)
  - Best practices
  - Common use cases
  - Frequent pitfalls
  - Performance considerations
  - Advanced patterns

### Practical Examples

All examples are runnable Go programs that demonstrate key concepts:

1. **[examples/basic_mutex.go](examples/basic_mutex.go)**
   - Unsafe vs safe counter (demonstrates race conditions)
   - Bank account with concurrent deposits/withdrawals
   - Basic mutex usage patterns
   - Run with: `go run examples/basic_mutex.go`
   - Detect races: `go run -race examples/basic_mutex.go`

2. **[examples/rwmutex.go](examples/rwmutex.go)**
   - Read-write mutex for cache implementation
   - Performance comparison: Mutex vs RWMutex
   - Statistics tracker with RWMutex
   - Run with: `go run examples/rwmutex.go`

3. **[examples/pitfalls.go](examples/pitfalls.go)**
   - Deadlock prevention (lock ordering)
   - Missing unlock (why to use defer)
   - Blocking operations while holding locks
   - Copying mutexes (pointer receivers)
   - Lock contention and sharding
   - Run with: `go run examples/pitfalls.go`

4. **[examples/advanced.go](examples/advanced.go)**
   - sync.Once for lazy initialization
   - Try-lock pattern
   - sync.Cond for producer-consumer
   - Mutex vs atomic operations
   - Double-checked locking vs sync.Once
   - Read-copy-update pattern
   - Run with: `go run examples/advanced.go`

## 🚀 Quick Start

### Prerequisites
- Go 1.16 or later

### Running Examples

```bash
# Run all examples
go run examples/basic_mutex.go
go run examples/rwmutex.go
go run examples/pitfalls.go
go run examples/advanced.go

# Detect race conditions
go run -race examples/basic_mutex.go

# Run with verbose output
go run -v examples/basic_mutex.go
```

## 📖 Key Concepts

### When to Use Mutexes

✅ **Use mutexes when:**
- Protecting shared mutable state
- Multiple goroutines access the same data structure
- Operations involve multiple variables
- Complex state transitions

❌ **Don't use mutexes when:**
- Single atomic variable (use `sync/atomic`)
- Passing data between goroutines (use channels)
- Data is immutable
- Using already synchronized types (`sync.Map`, `sync.Pool`)

### Mutex vs RWMutex

**Use `sync.Mutex`:**
- Simple mutual exclusion
- Reads and writes are roughly equal
- Critical sections are very short

**Use `sync.RWMutex`:**
- Many more reads than writes
- Critical sections are long enough to amortize overhead
- Need to allow concurrent reads

### Performance Tips

1. **Keep critical sections small** - only protect what needs protection
2. **Use atomic operations** for simple counters and flags
3. **Shard data structures** to reduce contention
4. **Profile before optimizing** - measure with `-mutexprofile`
5. **Consider lock-free alternatives** for very hot paths

## 🔍 Common Pitfalls

1. **Deadlocks** - Always acquire locks in consistent order
2. **Forgetting to unlock** - Always use `defer` for unlocking
3. **Copying mutexes** - Use pointer receivers, run `go vet`
4. **Holding locks too long** - Don't do I/O or blocking ops while locked
5. **Reentrant locking** - Go mutexes are NOT reentrant

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
```

## 📊 Learning Path

1. **Start with basics** - Read [go-mutexes.md](go-mutexes.md) introduction
2. **Run basic examples** - Execute `basic_mutex.go` with and without `-race`
3. **Understand RWMutex** - Study `rwmutex.go` for read-heavy workloads
4. **Learn common pitfalls** - Review `pitfalls.go` to avoid mistakes
5. **Explore advanced patterns** - Study `advanced.go` for optimization techniques
6. **Practice** - Implement your own concurrent data structures

## 📚 Additional Resources

- [Go sync package documentation](https://pkg.go.dev/sync)
- [Go Memory Model](https://go.dev/ref/mem)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- Source code: `src/sync/mutex.go` in Go repository

## 🎯 Summary

### Key Takeaways

1. Mutexes protect shared state from concurrent access
2. Go's implementation balances throughput and fairness (normal/starvation modes)
3. Always use `defer` to unlock mutexes
4. Keep critical sections as small as possible
5. Never copy mutexes (use pointers, run `go vet`)
6. Beware of deadlocks from inconsistent lock ordering
7. Use RWMutex when reads greatly outnumber writes
8. Consider alternatives: atomics, channels, immutability
9. Profile before optimizing

### Quick Reference

| Operation | When to Use |
|-----------|-------------|
| `sync.Mutex` | Simple mutual exclusion |
| `sync.RWMutex` | Many readers, few writers |
| `sync.Once` | One-time initialization |
| `sync.Cond` | Complex condition waiting |
| `sync/atomic` | Single-variable operations |
| Channels | Passing data between goroutines |

## 🤝 Contributing

This is a learning repository. Feel free to:
- Report issues or inaccuracies
- Suggest improvements
- Add more examples
- Enhance documentation

## 📝 License

This repository is for educational purposes.

---

*Happy learning! Remember: "Don't communicate by sharing memory; share memory by communicating." - Rob Pike*
