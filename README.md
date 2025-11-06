# Learning Repository

A comprehensive collection of technical learning resources covering programming languages, system fundamentals, and best practices. Each topic includes detailed documentation, practical examples, and hands-on code.

## 📚 Contents

### Go Programming

Comprehensive Go language resources covering concurrency primitives and testing frameworks:

- **[Go Documentation](go/README.md)** - Complete guide to Go packages and patterns
  - [Mutexes](go/mutexes.md) - Understanding mutexes, RWMutex, and synchronization primitives
  - [Atomic Operations](go/packages/atomic.md) - Lock-free synchronization with `sync/atomic`
  - [Gomega Testing](go/packages/gomega.md) - Expressive matcher library for testing
  - [Mutex Examples](go/examples-mutexes/) - Practical mutex usage patterns
  - [Atomic Examples](go/examples/) - Working code using atomic operations

### Linux System Fundamentals

Deep dive into Linux kernel internals and system programming:

- **[Linux Documentation](linux/README.md)** - System-level concepts and implementation details
  - [Kernel Buffers](linux/kernel-buffers.md) - Memory management, buffering mechanisms, and I/O optimization

## 🚀 Quick Start

### Go Topics

```bash
# Navigate to Go directory
cd go

# Run mutex examples
go run examples-mutexes/basic_mutex.go
go run -race examples-mutexes/pitfalls.go

# Run atomic examples and tests
go test -v ./examples/
go test -race ./examples/
```

### Linux Topics

Browse the Linux documentation to understand kernel-level concepts:

```bash
cd linux
# Read through kernel-buffers.md
```

## 🎯 What's Inside

Each topic includes:

- 📖 **Comprehensive Documentation** - Detailed explanations with theory and implementation details
- 💻 **Working Code Examples** - Production-quality, runnable code
- ✅ **Test Suites** - Real-world testing scenarios and best practices
- 🔍 **Low-Level Details** - Under-the-hood implementation and hardware interactions
- 📊 **Performance Considerations** - Optimization tips and benchmarking guidance
- ⚠️ **Common Pitfalls** - Known issues and how to avoid them

## 📂 Repository Structure

```
learning/
├── README.md           # This file
├── go/                 # Go programming resources
│   ├── README.md       # Go topics overview
│   ├── mutexes.md      # Mutex comprehensive guide
│   ├── packages/       # Package-specific documentation
│   ├── examples/       # Atomic operations examples
│   └── examples-mutexes/ # Mutex usage examples
└── linux/              # Linux system resources
    ├── README.md       # Linux topics overview
    └── kernel-buffers.md # Kernel buffering guide
```

## 📖 Learning Path

### For Go Developers

1. **Start with Mutexes** - Understand basic synchronization
   - Read [mutexes.md](go/mutexes.md)
   - Run examples in [examples-mutexes/](go/examples-mutexes/)
   
2. **Explore Atomic Operations** - Learn lock-free programming
   - Study [atomic.md](go/packages/atomic.md)
   - Experiment with [atomic examples](go/examples/)
   
3. **Master Testing** - Write better tests
   - Review [gomega.md](go/packages/gomega.md)
   - Study test patterns in [atomic_examples_test.go](go/examples/atomic_examples_test.go)

### For Systems Programmers

1. **Linux Fundamentals** - Understand kernel internals
   - Read [kernel-buffers.md](linux/kernel-buffers.md)
   - Study buffer types and memory management

## 🛠️ Tools and Prerequisites

### Go Development
- Go 1.16 or later
- Race detector: `go run -race`
- Vet tool: `go vet ./...`
- Format: `go fmt ./...`

### Linux Topics
- Basic understanding of C and systems programming
- Familiarity with Unix/Linux command line
- Optional: Linux kernel source code for reference

## 🤝 Contributing

This is an educational repository. Contributions welcome:

- 📝 Improve documentation clarity
- 💡 Add more examples
- 🐛 Fix inaccuracies
- ✨ Suggest new topics
- 📚 Add references and resources

## 📝 License

Educational purposes. Free to use and share.

---

*Happy learning! Build deep understanding through theory, practice, and experimentation.* 🚀
