package examples

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestAtomicCounter(t *testing.T) {
	g := NewWithT(t)

	counter := &AtomicCounter{}

	// Test initial value
	g.Expect(counter.Get()).To(Equal(int64(0)))

	// Test increment
	result := counter.Increment()
	g.Expect(result).To(Equal(int64(1)))
	g.Expect(counter.Get()).To(Equal(int64(1)))

	// Test multiple increments
	counter.Increment()
	counter.Increment()
	g.Expect(counter.Get()).To(Equal(int64(3)))

	// Test decrement
	result = counter.Decrement()
	g.Expect(result).To(Equal(int64(2)))

	// Test set
	counter.Set(100)
	g.Expect(counter.Get()).To(Equal(int64(100)))

	// Test CompareAndSet
	success := counter.CompareAndSet(100, 200)
	g.Expect(success).To(BeTrue())
	g.Expect(counter.Get()).To(Equal(int64(200)))

	// Failed CompareAndSet
	success = counter.CompareAndSet(100, 300)
	g.Expect(success).To(BeFalse())
	g.Expect(counter.Get()).To(Equal(int64(200)))
}

func TestAtomicCounterConcurrency(t *testing.T) {
	g := NewWithT(t)

	counter := &AtomicCounter{}
	done := make(chan bool)

	// Start 100 goroutines, each incrementing 1000 times
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				counter.Increment()
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify final count
	g.Expect(counter.Get()).To(Equal(int64(100000)))
}

func TestAtomicConfig(t *testing.T) {
	g := NewWithT(t)

	initial := Config{
		MaxConnections: 100,
		Timeout:        5,
		Debug:          false,
	}

	ac := NewAtomicConfig(initial)

	// Test initial config
	cfg := ac.Get()
	g.Expect(cfg.MaxConnections).To(Equal(100))
	g.Expect(cfg.Timeout).To(Equal(5))
	g.Expect(cfg.Debug).To(BeFalse())

	// Test update
	newCfg := Config{
		MaxConnections: 200,
		Timeout:        10,
		Debug:          true,
	}
	ac.Update(newCfg)

	cfg = ac.Get()
	g.Expect(cfg.MaxConnections).To(Equal(200))
	g.Expect(cfg.Timeout).To(Equal(10))
	g.Expect(cfg.Debug).To(BeTrue())
}

func TestAtomicFlag(t *testing.T) {
	g := NewWithT(t)

	flag := &AtomicFlag{}

	// Test initial state
	g.Expect(flag.IsSet()).To(BeFalse())

	// Test set
	flag.Set()
	g.Expect(flag.IsSet()).To(BeTrue())

	// Test clear
	flag.Clear()
	g.Expect(flag.IsSet()).To(BeFalse())

	// Test toggle
	result := flag.Toggle()
	g.Expect(result).To(BeTrue())
	g.Expect(flag.IsSet()).To(BeTrue())

	result = flag.Toggle()
	g.Expect(result).To(BeFalse())
	g.Expect(flag.IsSet()).To(BeFalse())
}

func TestReferenceCounter(t *testing.T) {
	g := NewWithT(t)

	called := false
	rc := NewReferenceCounter(func() {
		called = true
	})

	// Initial count should be 1
	g.Expect(rc.Count()).To(Equal(int32(1)))

	// Acquire increases count
	rc.Acquire()
	g.Expect(rc.Count()).To(Equal(int32(2)))

	rc.Acquire()
	g.Expect(rc.Count()).To(Equal(int32(3)))

	// Release decreases count
	rc.Release()
	g.Expect(rc.Count()).To(Equal(int32(2)))
	g.Expect(called).To(BeFalse())

	rc.Release()
	g.Expect(rc.Count()).To(Equal(int32(1)))
	g.Expect(called).To(BeFalse())

	// Final release triggers callback
	rc.Release()
	g.Expect(rc.Count()).To(Equal(int32(0)))
	g.Expect(called).To(BeTrue())
}

func TestSpinLock(t *testing.T) {
	g := NewWithT(t)

	lock := &SpinLock{}

	// Test TryLock on unlocked lock
	g.Expect(lock.TryLock()).To(BeTrue())

	// Test TryLock on locked lock
	g.Expect(lock.TryLock()).To(BeFalse())

	// Unlock
	lock.Unlock()

	// Test Lock/Unlock
	lock.Lock()
	g.Expect(lock.TryLock()).To(BeFalse())
	lock.Unlock()
	g.Expect(lock.TryLock()).To(BeTrue())
	lock.Unlock()
}

func TestMetrics(t *testing.T) {
	g := NewWithT(t)

	metrics := &Metrics{}

	// Test initial state
	req, err, bytes := metrics.GetSnapshot()
	g.Expect(req).To(Equal(int64(0)))
	g.Expect(err).To(Equal(int64(0)))
	g.Expect(bytes).To(Equal(int64(0)))

	// Record some metrics
	metrics.RecordRequest()
	metrics.RecordRequest()
	metrics.RecordError()
	metrics.RecordBytes(1024)

	req, err, bytes = metrics.GetSnapshot()
	g.Expect(req).To(Equal(int64(2)))
	g.Expect(err).To(Equal(int64(1)))
	g.Expect(bytes).To(Equal(int64(1024)))

	// Record more
	metrics.RecordRequest()
	metrics.RecordBytes(512)

	req, err, bytes = metrics.GetSnapshot()
	g.Expect(req).To(Equal(int64(3)))
	g.Expect(err).To(Equal(int64(1)))
	g.Expect(bytes).To(Equal(int64(1536)))

	// Reset
	metrics.Reset()
	req, err, bytes = metrics.GetSnapshot()
	g.Expect(req).To(Equal(int64(0)))
	g.Expect(err).To(Equal(int64(0)))
	g.Expect(bytes).To(Equal(int64(0)))
}

func TestMetricsConcurrency(t *testing.T) {
	g := NewWithT(t)

	metrics := &Metrics{}
	done := make(chan bool)

	// Start multiple goroutines recording metrics
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				metrics.RecordRequest()
				if j%10 == 0 {
					metrics.RecordError()
				}
				metrics.RecordBytes(100)
			}
			done <- true
		}()
	}

	// Wait for completion
	for i := 0; i < 100; i++ {
		<-done
	}

	req, err, bytes := metrics.GetSnapshot()
	g.Expect(req).To(Equal(int64(10000)))
	g.Expect(err).To(Equal(int64(1000)))
	g.Expect(bytes).To(Equal(int64(1000000)))
}

func TestWorker(t *testing.T) {
	g := NewWithT(t)

	worker := NewWorker()

	// Initially not running
	g.Expect(worker.IsRunning()).To(BeFalse())

	// Start worker
	worker.Start()
	g.Expect(worker.IsRunning()).To(BeTrue())

	// Submit some work
	for i := 0; i < 10; i++ {
		worker.Submit(i)
	}

	// Eventually should process work
	g.Eventually(func() int64 {
		return worker.ProcessedCount()
	}, "2s", "50ms").Should(BeNumerically(">=", int64(10)))

	// Stop worker
	worker.Stop()
	g.Eventually(worker.IsRunning, "2s").Should(BeFalse())
}

func TestWorkerMultipleStarts(t *testing.T) {
	g := NewWithT(t)

	worker := NewWorker()

	worker.Start()
	g.Expect(worker.IsRunning()).To(BeTrue())

	// Starting again should be no-op
	worker.Start()
	g.Expect(worker.IsRunning()).To(BeTrue())

	worker.Stop()
	g.Eventually(worker.IsRunning, "2s").Should(BeFalse())
}

func TestWorkerProcessing(t *testing.T) {
	g := NewWithT(t)

	worker := NewWorker()
	worker.Start()

	// Submit work continuously
	go func() {
		for i := 0; i < 100; i++ {
			worker.Submit(i)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Verify work is being processed
	g.Eventually(func() int64 {
		return worker.ProcessedCount()
	}, "5s", "100ms").Should(BeNumerically(">", int64(50)))

	// Consistently verify it keeps processing
	initialCount := worker.ProcessedCount()
	g.Consistently(func() int64 {
		return worker.ProcessedCount()
	}, "1s", "100ms").Should(BeNumerically(">=", initialCount))

	worker.Stop()
}

func TestSafeMap(t *testing.T) {
	g := NewWithT(t)

	sm := &SafeMap{}

	// Test initial size
	g.Expect(sm.Size()).To(Equal(int64(0)))

	// Test set and get
	sm.Set("key1", "value1")
	g.Expect(sm.Size()).To(Equal(int64(1)))

	val, exists := sm.Get("key1")
	g.Expect(exists).To(BeTrue())
	g.Expect(val).To(Equal("value1"))

	// Test multiple sets
	sm.Set("key2", "value2")
	sm.Set("key3", "value3")
	g.Expect(sm.Size()).To(Equal(int64(3)))

	// Test overwrite (size shouldn't change)
	sm.Set("key1", "newvalue1")
	g.Expect(sm.Size()).To(Equal(int64(3)))

	val, exists = sm.Get("key1")
	g.Expect(exists).To(BeTrue())
	g.Expect(val).To(Equal("newvalue1"))

	// Test delete
	sm.Delete("key2")
	g.Expect(sm.Size()).To(Equal(int64(2)))

	_, exists = sm.Get("key2")
	g.Expect(exists).To(BeFalse())

	// Test non-existent key
	_, exists = sm.Get("nonexistent")
	g.Expect(exists).To(BeFalse())
}

func TestSafeMapConcurrency(t *testing.T) {
	g := NewWithT(t)

	sm := &SafeMap{}
	done := make(chan bool)

	// Start multiple goroutines setting values
	for i := 0; i < 100; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key_%d", id)
				sm.Set(key, id*10+j)
			}
			done <- true
		}(i)
	}

	// Wait for completion
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify size
	g.Expect(sm.Size()).To(BeNumerically(">", int64(0)))
	g.Expect(sm.Size()).To(BeNumerically("<=", int64(100)))
}

func TestGomegaMatcherExamples(t *testing.T) {
	g := NewWithT(t)

	// Numeric matchers
	g.Expect(42).To(Equal(42))
	g.Expect(42).To(BeNumerically("==", 42))
	g.Expect(10).To(BeNumerically(">", 5))
	g.Expect(10).To(BeNumerically(">=", 10))
	g.Expect(10).To(BeNumerically("<", 20))
	g.Expect(10).To(BeNumerically("~", 12, 3)) // Within 3 of 12

	// String matchers
	g.Expect("hello world").To(ContainSubstring("world"))
	g.Expect("hello").To(HavePrefix("hel"))
	g.Expect("hello").To(HaveSuffix("llo"))

	// Collection matchers
	slice := []string{"apple", "banana", "cherry"}
	g.Expect(slice).To(HaveLen(3))
	g.Expect(slice).To(ContainElement("banana"))
	g.Expect(slice).To(ContainElements("apple", "cherry"))
	g.Expect(slice).NotTo(BeEmpty())

	// Boolean matchers
	g.Expect(true).To(BeTrue())
	g.Expect(false).To(BeFalse())

	// Nil matchers
	var ptr *int
	g.Expect(ptr).To(BeNil())
	value := 42
	ptr = &value
	g.Expect(ptr).NotTo(BeNil())
}
