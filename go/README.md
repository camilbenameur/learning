# Go Packages Documentation

This directory contains comprehensive documentation and working examples for important Go packages.

## Contents

### 📚 Documentation

- **[Atomic Package (`sync/atomic`)](packages/atomic.md)** - Low-level atomic memory primitives for lock-free synchronization
- **[Gomega](packages/gomega.md)** - Matcher/assertion library for expressive testing

### 💻 Working Examples

The `examples/` directory contains fully functional Go code demonstrating both packages:

- **[atomic_examples.go](examples/atomic_examples.go)** - Production-ready implementations using atomic operations
- **[atomic_examples_test.go](examples/atomic_examples_test.go)** - Comprehensive test suite using Gomega matchers

## Quick Start

### Running the Examples

```bash
# Navigate to the Go directory
cd go

# Install dependencies
go mod tidy

# Run all tests
go test -v ./examples/

# Run specific tests
go test -v ./examples/ -run TestAtomicCounter
go test -v ./examples/ -run TestWorker

# Run tests with race detector
go test -race ./examples/
```

### Atomic Package Overview

The `sync/atomic` package provides lock-free operations for:
- **Counters**: Thread-safe increment/decrement
- **Flags**: Atomic boolean operations
- **Configuration**: Hot-reload with atomic.Value
- **Reference counting**: Memory management without locks
- **Synchronization primitives**: Spin locks, barriers

**Key Strengths:**
- ✅ Hardware-level performance
- ✅ Zero allocation
- ✅ Sequential consistency guarantees
- ✅ Simple, clear API

**Learn more:** [atomic.md](packages/atomic.md)

### Gomega Overview

Gomega is a rich matcher library for writing expressive, readable tests:
- **Fluent API**: Tests that read like natural language
- **50+ Built-in matchers**: Cover all common testing scenarios
- **Async testing**: First-class support for `Eventually` and `Consistently`
- **Custom matchers**: Extend with domain-specific assertions
- **Framework agnostic**: Works with any testing framework

**Key Strengths:**
- ✅ Readable assertions
- ✅ Clear failure messages
- ✅ Composable matchers
- ✅ Async/concurrent testing support

**Learn more:** [gomega.md](packages/gomega.md)

## Example Implementations

### 1. Atomic Counter

```go
type AtomicCounter struct {
    value int64
}

func (c *AtomicCounter) Increment() int64 {
    return atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Get() int64 {
    return atomic.LoadInt64(&c.value)
}
```

### 2. Configuration Hot-Reload

```go
type AtomicConfig struct {
    config atomic.Value
}

func (ac *AtomicConfig) Update(cfg Config) {
    ac.config.Store(cfg)
}

func (ac *AtomicConfig) Get() Config {
    return ac.config.Load().(Config)
}
```

### 3. Worker Coordination

```go
type Worker struct {
    running   int32
    processed int64
}

func (w *Worker) Start() {
    if atomic.CompareAndSwapInt32(&w.running, 0, 1) {
        go w.run()
    }
}
```

### 4. Testing with Gomega

```go
func TestWorker(t *testing.T) {
    g := NewWithT(t)
    
    worker := NewWorker()
    worker.Start()
    
    // Wait for async processing
    g.Eventually(func() int64 {
        return worker.ProcessedCount()
    }, "2s", "50ms").Should(BeNumerically(">=", int64(10)))
    
    // Verify consistent state
    g.Consistently(func() bool {
        return worker.IsHealthy()
    }, "2s", "200ms").Should(BeTrue())
}
```

## Test Coverage

The test suite includes:
- ✅ Unit tests for all atomic operations
- ✅ Concurrency tests (100 goroutines)
- ✅ Async worker testing
- ✅ Reference counting validation
- ✅ Thread-safe map operations
- ✅ Comprehensive Gomega matcher examples

All tests pass with `-race` detector enabled.

## Use Cases

### Atomic Operations

1. **Performance-critical counters** - Metrics, statistics, rate limiting
2. **Lock-free data structures** - Stacks, queues, lists
3. **Configuration management** - Hot-reload without service restart
4. **Resource pooling** - Track available resources
5. **State machines** - Atomic state transitions
6. **Shutdown coordination** - Graceful shutdown flags

### Gomega Testing

1. **Unit testing** - Clear, expressive assertions
2. **Integration testing** - Test component interactions
3. **API testing** - HTTP response validation
4. **Concurrent code** - Async behavior verification
5. **Error handling** - Validate error conditions
6. **Data validation** - Business logic verification

## Performance Characteristics

### Atomic Operations

| Operation | Typical Latency | Notes |
|-----------|----------------|-------|
| Load | ~1-2 ns | Single CPU instruction |
| Store | ~1-2 ns | May require memory barrier |
| Add | ~2-5 ns | Read-modify-write cycle |
| CAS | ~5-10 ns | May retry on contention |

### Gomega Matchers

| Matcher Type | Performance | Best For |
|--------------|-------------|----------|
| Simple (Equal, BeTrue) | Very fast | Basic assertions |
| Collection (Contains) | O(n) | Small-medium collections |
| Eventually | Polling overhead | Async operations |
| Complex (MatchFields) | Reflection overhead | Detailed struct validation |

## Best Practices

### Atomic Operations

1. ✅ **Always use consistently** - Never mix atomic and non-atomic access
2. ✅ **Use atomic.Value for complex types** - Safer than unsafe.Pointer
3. ✅ **Consider cache line padding** - Prevent false sharing
4. ✅ **Profile before optimizing** - Mutexes are often sufficient
5. ❌ **Avoid for complex operations** - Use mutexes for multi-variable updates

### Gomega Testing

1. ✅ **Use descriptive test names** - Communicate intent clearly
2. ✅ **Choose specific matchers** - More readable and better errors
3. ✅ **Use Eventually for async** - Avoid flaky sleep-based tests
4. ✅ **Set reasonable timeouts** - Balance responsiveness vs reliability
5. ❌ **Don't ignore errors** - Always check function errors

## Further Reading

### Documentation
- [Go Memory Model](https://go.dev/ref/mem)
- [sync/atomic Package](https://pkg.go.dev/sync/atomic)
- [Gomega Documentation](https://onsi.github.io/gomega/)
- [Ginkgo + Gomega Guide](https://onsi.github.io/ginkgo/)

### Articles
- [Go 101: Atomic Operations](https://go101.org/article/atomic.html)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

## Contributing

This is a learning repository. Feel free to:
- Add more examples
- Improve documentation
- Add test cases
- Share insights and best practices

## License

This documentation and code examples are provided for educational purposes.
