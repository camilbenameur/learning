# Gomega - Go Matcher Library

## Overview

Gomega is a matcher/assertion library for Go, designed to work seamlessly with testing frameworks like Ginkgo and Go's built-in testing package. It provides expressive, readable assertions that make tests easier to write and understand.

## Purpose

Gomega is designed for:
- **Expressive Testing**: Write tests that read like natural language
- **Rich Matchers**: Comprehensive set of matchers for various data types
- **Async Testing**: Built-in support for testing asynchronous behavior
- **Flexible Integration**: Works with any testing framework
- **Better Error Messages**: Clear, actionable failure messages

## Key Strengths

### 1. **Readable Assertions**
Gomega uses a fluent API that makes tests self-documenting:

```go
Expect(response.StatusCode).To(Equal(200))
Expect(user.Name).NotTo(BeEmpty())
Expect(items).To(ContainElement("apple"))
```

### 2. **Rich Matcher Library**
Over 50 built-in matchers covering:
- Equality and identity
- Numeric comparisons
- String operations
- Collection operations
- Type assertions
- Error handling
- Channel operations
- HTTP responses

### 3. **Asynchronous Testing**
First-class support for testing concurrent code:

```go
Eventually(func() int {
    return counter.Value()
}).Should(Equal(100))

Consistently(func() bool {
    return server.IsHealthy()
}).Should(BeTrue())
```

### 4. **Composable Matchers**
Combine matchers for complex assertions:

```go
Expect(users).To(ContainElement(And(
    HaveField("Name", "Alice"),
    HaveField("Age", BeNumerically(">", 18)),
)))
```

### 5. **Custom Matchers**
Extend Gomega with domain-specific matchers:

```go
func BeValidEmail() types.GomegaMatcher {
    return &emailMatcher{}
}
```

## Low-Level Mechanics

### Matcher Interface

At its core, Gomega uses the `GomegaMatcher` interface:

```go
type GomegaMatcher interface {
    Match(actual interface{}) (success bool, err error)
    FailureMessage(actual interface{}) (message string)
    NegatedFailureMessage(actual interface{}) (message string)
}
```

### Assertion Flow

1. **Expectation Creation**: `Expect(actual)` wraps the value
2. **Matcher Application**: `.To(matcher)` applies the matcher
3. **Match Evaluation**: Matcher's `Match()` is called
4. **Result Handling**: On failure, calls appropriate message method
5. **Test Failure**: Reports to testing framework

### Internal Architecture

```
┌─────────────┐
│   Expect()  │ ← Entry point
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Assertion  │ ← Holds actual value
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Matcher   │ ← Implements matching logic
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Result    │ ← Success/Failure + Messages
└─────────────┘
```

### Gomega Context

Gomega uses a context to:
- Store the current testing T
- Track timeout configurations
- Manage polling intervals for Eventually/Consistently
- Handle custom failure handlers

```go
type Gomega interface {
    Expect(actual interface{}, extra ...interface{}) Assertion
    Eventually(args ...interface{}) AsyncAssertion
    Consistently(args ...interface{}) AsyncAssertion
}
```

### Async Assertion Mechanics

`Eventually` and `Consistently` use goroutines and timers:

```go
// Pseudo-code for Eventually
func (a *AsyncAssertion) Should(matcher GomegaMatcher) {
    timeout := time.After(a.timeout)
    ticker := time.NewTicker(a.pollingInterval)
    
    for {
        select {
        case <-timeout:
            // Fail with timeout message
            return
        case <-ticker.C:
            result, err := matcher.Match(a.actualFunc())
            if result && err == nil {
                return // Success
            }
        }
    }
}
```

## Core Matchers

### Equality Matchers

```go
// Equal - Deep equality
Expect(actual).To(Equal(expected))

// BeIdenticalTo - Pointer equality
Expect(actual).To(BeIdenticalTo(expected))

// BeEquivalentTo - Type-insensitive equality
Expect(int32(5)).To(BeEquivalentTo(5))
```

### Numeric Matchers

```go
Expect(value).To(BeNumerically("==", 100))
Expect(value).To(BeNumerically(">", 50))
Expect(value).To(BeNumerically("~", 100, 5)) // Within 5 of 100
```

### String Matchers

```go
Expect(str).To(ContainSubstring("hello"))
Expect(str).To(HavePrefix("http://"))
Expect(str).To(HaveSuffix(".com"))
Expect(str).To(MatchRegexp(`^\d{3}-\d{4}$`))
```

### Collection Matchers

```go
Expect(slice).To(HaveLen(5))
Expect(slice).To(BeEmpty())
Expect(slice).To(ContainElement("apple"))
Expect(slice).To(ContainElements("apple", "banana"))
Expect(slice).To(ConsistOf("apple", "banana", "cherry"))
Expect(map).To(HaveKey("username"))
Expect(map).To(HaveKeyWithValue("age", 25))
```

### Type Matchers

```go
Expect(value).To(BeNil())
Expect(value).To(BeAssignableToTypeOf(&User{}))
Expect(func() { panic("oh no") }).To(Panic())
```

### Error Matchers

```go
Expect(err).To(HaveOccurred())
Expect(err).NotTo(HaveOccurred())
Expect(err).To(MatchError("file not found"))
Expect(err).To(MatchError(ContainSubstring("timeout")))
```

### Boolean Matchers

```go
Expect(condition).To(BeTrue())
Expect(condition).To(BeFalse())
```

### Channel Matchers

```go
Expect(ch).To(BeClosed())
Expect(ch).To(Receive())
Expect(ch).To(Receive(&result))
Expect(ch).To(BeSent(value))
```

## Usage Examples

### Example 1: Basic Testing with Go's testing Package

```go
package user_test

import (
    "testing"
    . "github.com/onsi/gomega"
)

func TestUserValidation(t *testing.T) {
    g := NewWithT(t)
    
    user := &User{
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   25,
    }
    
    g.Expect(user.Name).To(Equal("Alice"))
    g.Expect(user.Email).To(MatchRegexp(`^[a-z]+@[a-z]+\.[a-z]+$`))
    g.Expect(user.Age).To(BeNumerically(">=", 18))
}
```

### Example 2: Ginkgo Integration

```go
package user_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
    var user *User
    
    BeforeEach(func() {
        user = NewUser("Alice", "alice@example.com")
    })
    
    Context("when creating a new user", func() {
        It("should have a valid name", func() {
            Expect(user.Name).NotTo(BeEmpty())
            Expect(user.Name).To(HaveLen(5))
        })
        
        It("should have a valid email", func() {
            Expect(user.Email).To(ContainSubstring("@"))
            Expect(user.Email).To(HaveSuffix(".com"))
        })
    })
    
    Context("when validating user", func() {
        It("should return no error for valid user", func() {
            err := user.Validate()
            Expect(err).NotTo(HaveOccurred())
        })
        
        It("should return error for invalid email", func() {
            user.Email = "invalid"
            err := user.Validate()
            Expect(err).To(HaveOccurred())
            Expect(err).To(MatchError(ContainSubstring("email")))
        })
    })
})
```

### Example 3: Asynchronous Testing

```go
package worker_test

import (
    "testing"
    "time"
    . "github.com/onsi/gomega"
)

func TestWorkerProcessing(t *testing.T) {
    g := NewWithT(t)
    
    worker := NewWorker()
    worker.Start()
    
    // Eventually - wait for condition to become true
    g.Eventually(func() int {
        return worker.ProcessedCount()
    }, "5s", "100ms").Should(BeNumerically(">", 10))
    
    // Consistently - verify condition stays true
    g.Consistently(func() bool {
        return worker.IsHealthy()
    }, "2s", "200ms").Should(BeTrue())
    
    worker.Stop()
    
    // Wait for graceful shutdown
    g.Eventually(worker.IsRunning, "3s").Should(BeFalse())
}
```

### Example 4: Complex Matchers

```go
package api_test

import (
    "testing"
    . "github.com/onsi/gomega"
)

func TestAPIResponse(t *testing.T) {
    g := NewWithT(t)
    
    response := &APIResponse{
        Status: 200,
        Data: map[string]interface{}{
            "users": []User{
                {Name: "Alice", Age: 25},
                {Name: "Bob", Age: 30},
            },
            "total": 2,
        },
    }
    
    // Nested matchers
    g.Expect(response).To(And(
        HaveField("Status", Equal(200)),
        HaveField("Data", HaveKey("users")),
    ))
    
    // Collection with element matchers
    g.Expect(response.Data["users"]).To(ContainElement(
        HaveField("Name", "Alice"),
    ))
    
    // Multiple conditions
    g.Expect(response.Data["users"]).To(And(
        HaveLen(2),
        ContainElement(HaveField("Age", BeNumerically(">", 20))),
    ))
}
```

### Example 5: Testing HTTP Handlers

```go
package handler_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    . "github.com/onsi/gomega"
)

func TestUserHandler(t *testing.T) {
    g := NewWithT(t)
    
    handler := NewUserHandler()
    req := httptest.NewRequest("GET", "/users/123", nil)
    rec := httptest.NewRecorder()
    
    handler.ServeHTTP(rec, req)
    
    g.Expect(rec.Code).To(Equal(http.StatusOK))
    g.Expect(rec.Header().Get("Content-Type")).To(Equal("application/json"))
    g.Expect(rec.Body.String()).To(ContainSubstring(`"name"`))
    g.Expect(rec.Body.String()).To(MatchJSON(`{"id": 123, "name": "Alice"}`))
}
```

### Example 6: Custom Matchers

```go
package matchers

import (
    "fmt"
    "strings"
    "github.com/onsi/gomega/types"
)

type validEmailMatcher struct {
    expectedDomain string
}

func BeValidEmail(domain string) types.GomegaMatcher {
    return &validEmailMatcher{
        expectedDomain: domain,
    }
}

func (m *validEmailMatcher) Match(actual interface{}) (bool, error) {
    email, ok := actual.(string)
    if !ok {
        return false, fmt.Errorf("BeValidEmail expects a string")
    }
    
    if !strings.Contains(email, "@") {
        return false, nil
    }
    
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false, nil
    }
    
    if m.expectedDomain != "" && parts[1] != m.expectedDomain {
        return false, nil
    }
    
    return true, nil
}

func (m *validEmailMatcher) FailureMessage(actual interface{}) string {
    return fmt.Sprintf("Expected\n\t%v\nto be a valid email address", actual)
}

func (m *validEmailMatcher) NegatedFailureMessage(actual interface{}) string {
    return fmt.Sprintf("Expected\n\t%v\nnot to be a valid email address", actual)
}

// Usage
func TestCustomMatcher(t *testing.T) {
    g := NewWithT(t)
    
    g.Expect("alice@example.com").To(BeValidEmail("example.com"))
    g.Expect("invalid").NotTo(BeValidEmail(""))
}
```

### Example 7: Table-Driven Tests

```go
package calculator_test

import (
    "testing"
    . "github.com/onsi/gomega"
)

func TestCalculator(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        op       string
        expected int
        shouldErr bool
    }{
        {"addition", 2, 3, "+", 5, false},
        {"subtraction", 5, 3, "-", 2, false},
        {"multiplication", 4, 5, "*", 20, false},
        {"division", 10, 2, "/", 5, false},
        {"division by zero", 10, 0, "/", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            g := NewWithT(t)
            
            result, err := Calculate(tt.a, tt.b, tt.op)
            
            if tt.shouldErr {
                g.Expect(err).To(HaveOccurred())
            } else {
                g.Expect(err).NotTo(HaveOccurred())
                g.Expect(result).To(Equal(tt.expected))
            }
        })
    }
}
```

### Example 8: Testing Goroutines

```go
package concurrent_test

import (
    "sync"
    "testing"
    "time"
    . "github.com/onsi/gomega"
)

func TestConcurrentMap(t *testing.T) {
    g := NewWithT(t)
    
    m := NewSafeMap()
    var wg sync.WaitGroup
    
    // Spawn multiple goroutines writing to map
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(val int) {
            defer wg.Done()
            m.Set(fmt.Sprintf("key%d", val), val)
        }(i)
    }
    
    // Wait for all writes to complete
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()
    
    g.Eventually(done, "5s").Should(BeClosed())
    g.Expect(m.Len()).To(Equal(100))
    
    // Verify all keys exist
    for i := 0; i < 100; i++ {
        key := fmt.Sprintf("key%d", i)
        g.Expect(m.Has(key)).To(BeTrue())
        
        val, exists := m.Get(key)
        g.Expect(exists).To(BeTrue())
        g.Expect(val).To(Equal(i))
    }
}
```

## Common Scenarios

### 1. **Unit Testing**
Test individual functions and methods with clear, expressive assertions.

### 2. **Integration Testing**
Verify interactions between components with async matchers.

### 3. **API Testing**
Test HTTP handlers and responses with JSON and status code matchers.

### 4. **Concurrent Code Testing**
Use Eventually/Consistently for goroutines and channels.

### 5. **Error Handling Testing**
Verify error conditions and error messages.

### 6. **Data Validation**
Test data structures and business logic validation.

## Advanced Features

### Transform Functions

Apply transformations before matching:

```go
Expect(users).To(WithTransform(func(u []User) []string {
    names := make([]string, len(u))
    for i, user := range u {
        names[i] = user.Name
    }
    return names
}, ContainElement("Alice")))
```

### Pointer Matchers

```go
Expect(ptr).To(PointTo(Equal(42)))
Expect(ptr).To(PointTo(MatchFields(IgnoreExtras, Fields{
    "Name": Equal("Alice"),
})))
```

### Struct Field Matching

```go
Expect(user).To(MatchFields(IgnoreExtras, Fields{
    "Name": Equal("Alice"),
    "Age":  BeNumerically(">", 18),
}))

Expect(user).To(MatchAllFields(Fields{
    "ID":    Not(BeZero()),
    "Name":  Equal("Alice"),
    "Email": ContainSubstring("@"),
    "Age":   Equal(25),
}))
```

### Stop Trying

Control when async assertions stop:

```go
Eventually(fn).WithTimeout(5*time.Second).
    WithPolling(100*time.Millisecond).
    Should(Equal(expected))
```

## Best Practices

### 1. **Use Descriptive Test Names**

```go
// Good
It("should return error when email is invalid", func() {
    Expect(user.Validate()).To(HaveOccurred())
})

// Avoid
It("test 1", func() {
    Expect(user.Validate()).To(HaveOccurred())
})
```

### 2. **Choose Appropriate Matchers**

```go
// Good - specific matcher
Expect(slice).To(HaveLen(5))

// Avoid - less clear
Expect(len(slice)).To(Equal(5))
```

### 3. **Use Eventually for Async Operations**

```go
// Good
Eventually(func() bool {
    return server.IsReady()
}).Should(BeTrue())

// Avoid - flaky
time.Sleep(1 * time.Second)
Expect(server.IsReady()).To(BeTrue())
```

### 4. **Provide Context with Extra Arguments**

```go
Expect(result).To(Equal(expected), "Database query should return expected results")
```

### 5. **Combine Matchers Logically**

```go
Expect(response).To(And(
    HaveHTTPStatus(200),
    HaveHTTPHeaderWithValue("Content-Type", "application/json"),
))
```

### 6. **Create Custom Matchers for Domain Logic**

```go
// Better than complex inline assertions
Expect(order).To(BeValidOrder())
Expect(payment).To(BeSuccessfulPayment())
```

### 7. **Use Table-Driven Tests**

Organize multiple test cases systematically with Gomega assertions.

### 8. **Set Reasonable Timeouts**

```go
// Good - reasonable timeout
Eventually(fn, "2s", "100ms").Should(Succeed())

// Avoid - timeout too long or too short
Eventually(fn, "30s").Should(Succeed()) // Too long
Eventually(fn, "10ms").Should(Succeed()) // Too short
```

## Performance Considerations

### Matcher Performance

- **Simple matchers** (Equal, BeTrue) are very fast
- **Complex matchers** (ContainElement, MatchFields) iterate over data
- **Async matchers** (Eventually) poll repeatedly

### Optimization Tips

1. **Use specific matchers** - They're optimized for their use case
2. **Avoid unnecessary Eventually** - Use synchronous matchers when possible
3. **Tune polling intervals** - Balance responsiveness vs. CPU usage
4. **Cache expensive computations** - Don't recalculate in matcher functions

## Common Pitfalls

### 1. **Not Using NewWithT in Tests**

```go
// BAD - panic on failure
func TestSomething(t *testing.T) {
    Expect(value).To(Equal(expected)) // Panics
}

// GOOD - proper integration
func TestSomething(t *testing.T) {
    g := NewWithT(t)
    g.Expect(value).To(Equal(expected)) // Fails test properly
}
```

### 2. **Incorrect Eventually Usage**

```go
// BAD - function called immediately
Eventually(expensiveFunction()).Should(Equal(expected))

// GOOD - function passed as reference
Eventually(expensiveFunction).Should(Equal(expected))

// GOOD - wrapped in closure
Eventually(func() int {
    return expensiveFunction()
}).Should(Equal(expected))
```

### 3. **Comparing Uncomparable Types**

```go
// BAD - slices can't be compared with Equal
Expect(slice1).To(Equal(slice2))

// GOOD - use ConsistOf
Expect(slice1).To(ConsistOf(slice2))
```

### 4. **Ignoring Error Returns**

```go
// BAD - ignores error
result := SomeFunction()
Expect(result).To(Equal(expected))

// GOOD - checks error
result, err := SomeFunction()
Expect(err).NotTo(HaveOccurred())
Expect(result).To(Equal(expected))
```

## Integration Patterns

### With Ginkgo

```go
var _ = Describe("Service", func() {
    var service *Service
    
    BeforeEach(func() {
        service = NewService()
    })
    
    It("should process requests", func() {
        Expect(service.Process()).To(Succeed())
    })
})
```

### With Standard Testing

```go
func TestService(t *testing.T) {
    g := NewWithT(t)
    service := NewService()
    g.Expect(service.Process()).To(Succeed())
}
```

### With Testify

```go
func TestService(t *testing.T) {
    g := NewWithT(t)
    suite := &ServiceTestSuite{}
    
    suite.Run(t)
    g.Expect(suite.Service).NotTo(BeNil())
}
```

## References

- [Gomega Official Documentation](https://onsi.github.io/gomega/)
- [Gomega Matcher Reference](https://onsi.github.io/gomega/#provided-matchers)
- [Ginkgo + Gomega Guide](https://onsi.github.io/ginkgo/)
- [GitHub Repository](https://github.com/onsi/gomega)

## Summary

Gomega provides:
- **Expressive syntax** for readable tests
- **Rich matcher library** for various scenarios
- **Async testing support** for concurrent code
- **Flexible integration** with any testing framework
- **Custom matchers** for domain-specific assertions

Use Gomega when you want clear, maintainable tests that communicate intent effectively.
