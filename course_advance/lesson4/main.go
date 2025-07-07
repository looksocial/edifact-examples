// Lesson 4: Detecting Message Types
// This lesson covers how to detect and validate EDIFACT message types.

package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// ErrorType represents different types of processing errors
type ErrorType int

const (
	ErrorTypeSyntax ErrorType = iota
	ErrorTypeValidation
	ErrorTypeSystem
	ErrorTypeTransient
)

func (et ErrorType) String() string {
	switch et {
	case ErrorTypeSyntax:
		return "SYNTAX"
	case ErrorTypeValidation:
		return "VALIDATION"
	case ErrorTypeSystem:
		return "SYSTEM"
	case ErrorTypeTransient:
		return "TRANSIENT"
	default:
		return "UNKNOWN"
	}
}

// ProcessingError represents a detailed error with context
type ProcessingError struct {
	Type      ErrorType
	Message   string
	Segment   string
	Position  int
	Retryable bool
	Timestamp time.Time
}

func (e *ProcessingError) Error() string {
	return fmt.Sprintf("[%s] %s at position %d: %s",
		e.Type.String(), e.Segment, e.Position, e.Message)
}

// ErrorClassifier determines error types and retryability
type ErrorClassifier struct {
	rules map[string]ErrorType
}

func NewErrorClassifier() *ErrorClassifier {
	return &ErrorClassifier{
		rules: map[string]ErrorType{
			"UNH": ErrorTypeSyntax,
			"BGM": ErrorTypeValidation,
			"DTM": ErrorTypeValidation,
			"UNT": ErrorTypeSyntax,
		},
	}
}

func (ec *ErrorClassifier) Classify(err error, segment string) *ProcessingError {
	errorType := ErrorTypeSystem // Default to system error

	if segmentType, exists := ec.rules[segment]; exists {
		errorType = segmentType
	}

	return &ProcessingError{
		Type:      errorType,
		Message:   err.Error(),
		Segment:   segment,
		Position:  0,
		Retryable: ec.isRetryable(errorType),
		Timestamp: time.Now(),
	}
}

func (ec *ErrorClassifier) isRetryable(errorType ErrorType) bool {
	switch errorType {
	case ErrorTypeTransient:
		return true
	case ErrorTypeSystem:
		return true
	case ErrorTypeValidation:
		return false
	case ErrorTypeSyntax:
		return false
	default:
		return false
	}
}

// RetryHandler implements exponential backoff retry logic
type RetryHandler struct {
	maxRetries    int
	baseDelay     time.Duration
	maxDelay      time.Duration
	backoffFactor float64
}

func NewRetryHandler(maxRetries int, baseDelay, maxDelay time.Duration) *RetryHandler {
	return &RetryHandler{
		maxRetries:    maxRetries,
		baseDelay:     baseDelay,
		maxDelay:      maxDelay,
		backoffFactor: 2.0,
	}
}

func (rh *RetryHandler) ExecuteWithRetry(operation func() error) error {
	var lastErr error

	for attempt := 0; attempt <= rh.maxRetries; attempt++ {
		if err := operation(); err == nil {
			if attempt > 0 {
				fmt.Printf("Operation succeeded after %d retries\n", attempt)
			}
			return nil
		} else {
			lastErr = err
			fmt.Printf("Attempt %d failed: %v\n", attempt+1, err)
		}

		if attempt < rh.maxRetries {
			delay := rh.calculateDelay(attempt)
			fmt.Printf("Waiting %v before retry...\n", delay)
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w",
		rh.maxRetries+1, lastErr)
}

func (rh *RetryHandler) calculateDelay(attempt int) time.Duration {
	delay := time.Duration(float64(rh.baseDelay) *
		math.Pow(rh.backoffFactor, float64(attempt)))
	if delay > rh.maxDelay {
		delay = rh.maxDelay
	}
	return delay
}

// CircuitBreaker implements circuit breaker pattern
type CircuitBreaker struct {
	failureThreshold int
	timeout          time.Duration
	state            CircuitState
	failureCount     int
	lastFailureTime  time.Time
	mu               sync.RWMutex
}

type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

func (cs CircuitState) String() string {
	switch cs {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

func NewCircuitBreaker(failureThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		timeout:          timeout,
		state:            StateClosed,
	}
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
	if !cb.canExecute() {
		return fmt.Errorf("circuit breaker is %s", cb.state.String())
	}

	err := operation()
	cb.recordResult(err)
	return err
}

func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.state = StateHalfOpen
			cb.mu.Unlock()
			cb.mu.RLock()
			fmt.Println("Circuit breaker transitioning to HALF_OPEN")
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		if cb.state == StateHalfOpen ||
			(cb.state == StateClosed && cb.failureCount >= cb.failureThreshold) {
			cb.state = StateOpen
			fmt.Printf("Circuit breaker opened after %d failures\n", cb.failureCount)
		}
	} else {
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
			cb.failureCount = 0
			fmt.Println("Circuit breaker closed - operation succeeded")
		}
	}
}

// ErrorLogger provides comprehensive error logging
type ErrorLogger struct {
	errors []*ProcessingError
	mu     sync.Mutex
}

func NewErrorLogger() *ErrorLogger {
	return &ErrorLogger{
		errors: make([]*ProcessingError, 0),
	}
}

func (el *ErrorLogger) LogError(err *ProcessingError) {
	el.mu.Lock()
	defer el.mu.Unlock()

	el.errors = append(el.errors, err)
	fmt.Printf("[ERROR] %s: %s\n", err.Type.String(), err.Message)
}

func (el *ErrorLogger) GetErrorStats() map[ErrorType]int {
	el.mu.Lock()
	defer el.mu.Unlock()

	stats := make(map[ErrorType]int)
	for _, err := range el.errors {
		stats[err.Type]++
	}
	return stats
}

// Mock operations for demonstration
func simulateTransientError() error {
	// Simulate a transient error (e.g., network timeout)
	return fmt.Errorf("network timeout")
}

func simulateSystemError() error {
	// Simulate a system error (e.g., database connection)
	return fmt.Errorf("database connection failed")
}

func simulateValidationError() error {
	// Simulate a validation error (e.g., invalid data)
	return fmt.Errorf("invalid invoice amount")
}

func main() {
	fmt.Println("=== Error Handling & Recovery (Lesson 4) ===")

	// Create components
	classifier := NewErrorClassifier()
	retryHandler := NewRetryHandler(3, 100*time.Millisecond, 1*time.Second)
	circuitBreaker := NewCircuitBreaker(2, 5*time.Second)
	errorLogger := NewErrorLogger()

	// Test error classification
	fmt.Println("\n=== Error Classification Demo ===")

	testErrors := []struct {
		message string
		segment string
	}{
		{"Invalid UNH segment", "UNH"},
		{"Invalid invoice amount", "BGM"},
		{"Invalid date format", "DTM"},
		{"Database connection failed", "UNT"},
	}

	for _, test := range testErrors {
		err := fmt.Errorf(test.message)
		processingErr := classifier.Classify(err, test.segment)
		fmt.Printf("Error: %s -> Type: %s, Retryable: %t\n",
			test.message, processingErr.Type.String(), processingErr.Retryable)
		errorLogger.LogError(processingErr)
	}

	// Test retry logic
	fmt.Println("\n=== Retry Logic Demo ===")

	fmt.Println("Testing retry with transient error:")
	err := retryHandler.ExecuteWithRetry(simulateTransientError)
	if err != nil {
		fmt.Printf("Final error: %v\n", err)
	}

	// Test circuit breaker
	fmt.Println("\n=== Circuit Breaker Demo ===")

	fmt.Println("Testing circuit breaker with system errors:")
	for i := 0; i < 5; i++ {
		err := circuitBreaker.Execute(simulateSystemError)
		if err != nil {
			fmt.Printf("Attempt %d: %v\n", i+1, err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Wait for circuit breaker to reset
	fmt.Println("Waiting for circuit breaker timeout...")
	time.Sleep(6 * time.Second)

	// Test successful operation
	fmt.Println("Testing successful operation:")
	err = circuitBreaker.Execute(func() error { return nil })
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Operation succeeded")
	}

	// Display error statistics
	fmt.Println("\n=== Error Statistics ===")
	stats := errorLogger.GetErrorStats()
	for errorType, count := range stats {
		fmt.Printf("%s errors: %d\n", errorType.String(), count)
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Error classification and categorization")
	fmt.Println("✅ Retry logic with exponential backoff")
	fmt.Println("✅ Circuit breaker pattern implementation")
	fmt.Println("✅ Comprehensive error logging")
	fmt.Println("✅ Error statistics and monitoring")
}
