// Package main demonstrates basic mutex usage in Go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Counter demonstrates unsafe concurrent access
type UnsafeCounter struct {
	value int
}

func (c *UnsafeCounter) Increment() {
	c.value++
}

func (c *UnsafeCounter) Value() int {
	return c.value
}

// SafeCounter demonstrates proper mutex usage
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func demonstrateUnsafe() {
	fmt.Println("=== Unsafe Counter (Race Condition) ===")
	counter := &UnsafeCounter{}
	var wg sync.WaitGroup

	// Launch 1000 goroutines
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Expected: 1000, Got: %d (likely incorrect due to race)\n", counter.Value())
	fmt.Println("Run with 'go run -race basic_mutex.go' to detect the race condition")
}

func demonstrateSafe() {
	fmt.Println("\n=== Safe Counter (With Mutex) ===")
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// Launch 1000 goroutines
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Expected: 1000, Got: %d (correct!)\n", counter.Value())
}

// BankAccount demonstrates mutex protecting multiple fields
type BankAccount struct {
	mu      sync.Mutex
	balance int
	holder  string
}

func (a *BankAccount) Deposit(amount int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	fmt.Printf("[%s] Depositing %d, balance before: %d\n", a.holder, amount, a.balance)
	time.Sleep(10 * time.Millisecond) // Simulate processing
	a.balance += amount
	fmt.Printf("[%s] Balance after: %d\n", a.holder, a.balance)
}

func (a *BankAccount) Withdraw(amount int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if a.balance >= amount {
		fmt.Printf("[%s] Withdrawing %d, balance before: %d\n", a.holder, amount, a.balance)
		time.Sleep(10 * time.Millisecond) // Simulate processing
		a.balance -= amount
		fmt.Printf("[%s] Balance after: %d\n", a.holder, a.balance)
		return true
	}
	fmt.Printf("[%s] Insufficient funds to withdraw %d\n", a.holder, amount)
	return false
}

func (a *BankAccount) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func demonstrateBankAccount() {
	fmt.Println("\n=== Bank Account Example ===")
	account := &BankAccount{holder: "Alice", balance: 1000}
	
	var wg sync.WaitGroup
	
	// Concurrent deposits
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(amount int) {
			defer wg.Done()
			account.Deposit(amount)
		}(100)
	}
	
	// Concurrent withdrawals
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(amount int) {
			defer wg.Done()
			account.Withdraw(amount)
		}(150)
	}
	
	wg.Wait()
	fmt.Printf("\nFinal balance: %d\n", account.Balance())
}

func main() {
	fmt.Println("Basic Mutex Examples in Go")
	fmt.Println("===========================\n")
	
	demonstrateUnsafe()
	demonstrateSafe()
	demonstrateBankAccount()
	
	fmt.Println("✓ All examples completed successfully!")
}
