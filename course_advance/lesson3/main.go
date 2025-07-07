// Lesson 3: Performance Optimization
// This lesson demonstrates performance optimization techniques for high-throughput EDI processing.

package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Segment represents an EDIFACT segment
type Segment struct {
	Tag      string
	Elements []string
}

// SegmentPool for memory optimization
type SegmentPool struct {
	pool sync.Pool
}

// NewSegmentPool creates a new segment pool
func NewSegmentPool() *SegmentPool {
	return &SegmentPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Segment{
					Elements: make([]string, 0, 10), // Pre-allocate slice
				}
			},
		},
	}
}

// Get retrieves a segment from the pool
func (sp *SegmentPool) Get() *Segment {
	return sp.pool.Get().(*Segment)
}

// Put returns a segment to the pool
func (sp *SegmentPool) Put(segment *Segment) {
	// Reset for reuse
	segment.Tag = ""
	segment.Elements = segment.Elements[:0]
	sp.pool.Put(segment)
}

// ProcessingJob represents a job to be processed
type ProcessingJob struct {
	ID      int
	Message string
}

// ProcessingResult represents the result of processing
type ProcessingResult struct {
	JobID    int
	Segments int
	Duration time.Duration
	Error    error
}

// ConcurrentProcessor handles concurrent message processing
type ConcurrentProcessor struct {
	workers    int
	jobQueue   chan *ProcessingJob
	resultChan chan *ProcessingResult
	pool       *SegmentPool
}

// NewConcurrentProcessor creates a new concurrent processor
func NewConcurrentProcessor(workers int) *ConcurrentProcessor {
	return &ConcurrentProcessor{
		workers:    workers,
		jobQueue:   make(chan *ProcessingJob, workers*2),
		resultChan: make(chan *ProcessingResult, workers*2),
		pool:       NewSegmentPool(),
	}
}

// worker processes jobs from the queue
func (cp *ConcurrentProcessor) worker() {
	for job := range cp.jobQueue {
		start := time.Now()

		// Process the message
		segments, err := cp.processMessage(job.Message)

		result := &ProcessingResult{
			JobID:    job.ID,
			Segments: segments,
			Duration: time.Since(start),
			Error:    err,
		}

		cp.resultChan <- result
	}
}

// processMessage processes a single message using the pool
func (cp *ConcurrentProcessor) processMessage(message string) (int, error) {
	segments := strings.Split(message, "'")
	segmentCount := 0

	for _, segmentStr := range segments {
		if segmentStr == "" {
			continue
		}

		// Get segment from pool
		segment := cp.pool.Get()
		defer cp.pool.Put(segment)

		// Parse segment
		parts := strings.Split(segmentStr, "+")
		if len(parts) == 0 {
			continue
		}

		segment.Tag = parts[0]
		segment.Elements = append(segment.Elements, parts[1:]...)

		// Process segment (simulate work)
		time.Sleep(1 * time.Millisecond)
		segmentCount++
	}

	return segmentCount, nil
}

// Process processes multiple messages concurrently
func (cp *ConcurrentProcessor) Process(messages []string) []*ProcessingResult {
	// Start workers
	for i := 0; i < cp.workers; i++ {
		go cp.worker()
	}

	// Submit jobs
	for i, msg := range messages {
		cp.jobQueue <- &ProcessingJob{
			ID:      i + 1,
			Message: msg,
		}
	}

	// Close job queue after all jobs are submitted
	close(cp.jobQueue)

	// Collect results
	results := make([]*ProcessingResult, 0, len(messages))
	for i := 0; i < len(messages); i++ {
		result := <-cp.resultChan
		results = append(results, result)
	}

	return results
}

// MessageCache provides caching for processed messages
type MessageCache struct {
	cache map[string]*CachedMessage
	mu    sync.RWMutex
	ttl   time.Duration
}

// CachedMessage represents a cached message
type CachedMessage struct {
	Segments   int
	Processed  time.Time
	Expiration time.Time
}

// NewMessageCache creates a new message cache
func NewMessageCache(ttl time.Duration) *MessageCache {
	return &MessageCache{
		cache: make(map[string]*CachedMessage),
		ttl:   ttl,
	}
}

// Get retrieves a message from cache
func (mc *MessageCache) Get(key string) (*CachedMessage, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if cached, exists := mc.cache[key]; exists && !cached.Expired() {
		return cached, true
	}
	return nil, false
}

// Set stores a message in cache
func (mc *MessageCache) Set(key string, segments int) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	now := time.Now()
	mc.cache[key] = &CachedMessage{
		Segments:   segments,
		Processed:  now,
		Expiration: now.Add(mc.ttl),
	}
}

// Expired checks if a cached message has expired
func (cm *CachedMessage) Expired() bool {
	return time.Now().After(cm.Expiration)
}

// OptimizedProcessor combines all optimization techniques
type OptimizedProcessor struct {
	pool       *SegmentPool
	cache      *MessageCache
	concurrent *ConcurrentProcessor
}

// NewOptimizedProcessor creates a new optimized processor
func NewOptimizedProcessor(workers int) *OptimizedProcessor {
	return &OptimizedProcessor{
		pool:       NewSegmentPool(),
		cache:      NewMessageCache(5 * time.Minute),
		concurrent: NewConcurrentProcessor(workers),
	}
}

// ProcessMessage processes a message with all optimizations
func (op *OptimizedProcessor) ProcessMessage(message string) (int, error) {
	// Check cache first
	if cached, exists := op.cache.Get(message); exists {
		fmt.Printf("Cache hit! Returning cached result: %d segments\n", cached.Segments)
		return cached.Segments, nil
	}

	// Process message
	segments, err := op.concurrent.processMessage(message)
	if err != nil {
		return 0, err
	}

	// Cache result
	op.cache.Set(message, segments)

	return segments, nil
}

// ProcessBatch processes multiple messages with optimizations
func (op *OptimizedProcessor) ProcessBatch(messages []string) []*ProcessingResult {
	return op.concurrent.Process(messages)
}

// BenchmarkProcessor provides benchmarking capabilities
type BenchmarkProcessor struct {
	processor *OptimizedProcessor
}

// NewBenchmarkProcessor creates a new benchmark processor
func NewBenchmarkProcessor(workers int) *BenchmarkProcessor {
	return &BenchmarkProcessor{
		processor: NewOptimizedProcessor(workers),
	}
}

// Benchmark runs performance benchmarks
func (bp *BenchmarkProcessor) Benchmark(messages []string) *BenchmarkResult {
	// Warm up
	for i := 0; i < 10; i++ {
		bp.processor.ProcessMessage(messages[0])
	}

	// Force GC before benchmark
	runtime.GC()

	// Run benchmark
	start := time.Now()
	results := bp.processor.ProcessBatch(messages)
	duration := time.Since(start)

	// Calculate statistics
	totalSegments := 0
	errors := 0
	for _, result := range results {
		totalSegments += result.Segments
		if result.Error != nil {
			errors++
		}
	}

	return &BenchmarkResult{
		TotalMessages:  len(messages),
		TotalSegments:  totalSegments,
		TotalDuration:  duration,
		MessagesPerSec: float64(len(messages)) / duration.Seconds(),
		SegmentsPerSec: float64(totalSegments) / duration.Seconds(),
		Errors:         errors,
		AverageLatency: duration / time.Duration(len(messages)),
	}
}

// BenchmarkResult contains benchmark statistics
type BenchmarkResult struct {
	TotalMessages  int
	TotalSegments  int
	TotalDuration  time.Duration
	MessagesPerSec float64
	SegmentsPerSec float64
	Errors         int
	AverageLatency time.Duration
}

func (br *BenchmarkResult) Print() {
	fmt.Printf("\n=== Benchmark Results ===\n")
	fmt.Printf("Total Messages: %d\n", br.TotalMessages)
	fmt.Printf("Total Segments: %d\n", br.TotalSegments)
	fmt.Printf("Total Duration: %v\n", br.TotalDuration)
	fmt.Printf("Messages/sec: %.2f\n", br.MessagesPerSec)
	fmt.Printf("Segments/sec: %.2f\n", br.SegmentsPerSec)
	fmt.Printf("Errors: %d\n", br.Errors)
	fmt.Printf("Average Latency: %v\n", br.AverageLatency)
}

func main() {
	fmt.Println("=== Performance Optimization (Lesson 3) ===")

	// Generate test messages
	messages := generateTestMessages(100)

	// Create benchmark processor
	benchmarker := NewBenchmarkProcessor(runtime.NumCPU())

	// Run benchmarks with different configurations
	fmt.Println("\n--- Benchmarking with different worker counts ---")

	workerCounts := []int{1, 2, 4, 8}
	for _, workers := range workerCounts {
		fmt.Printf("\nTesting with %d workers:\n", workers)
		benchmarker = NewBenchmarkProcessor(workers)
		result := benchmarker.Benchmark(messages)
		result.Print()
	}

	// Demonstrate memory pooling
	fmt.Println("\n--- Memory Pooling Demo ---")
	pool := NewSegmentPool()

	// Simulate high-frequency segment creation
	start := time.Now()
	for i := 0; i < 10000; i++ {
		segment := pool.Get()
		segment.Tag = "TEST"
		segment.Elements = []string{"element1", "element2"}
		pool.Put(segment)
	}
	duration := time.Since(start)
	fmt.Printf("Processed 10,000 segments with pooling in %v\n", duration)

	// Demonstrate caching
	fmt.Println("\n--- Caching Demo ---")
	processor := NewOptimizedProcessor(4)

	// Process same message multiple times
	testMessage := messages[0]
	for i := 0; i < 5; i++ {
		start := time.Now()
		segments, err := processor.ProcessMessage(testMessage)
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Attempt %d: %d segments in %v\n", i+1, segments, duration)
		}
	}

	// Memory usage demonstration
	fmt.Println("\n--- Memory Usage Demo ---")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Heap Alloc: %d bytes\n", m.HeapAlloc)
	fmt.Printf("Heap Sys: %d bytes\n", m.HeapSys)
	fmt.Printf("Num GC: %d\n", m.NumGC)

	// Force GC and show memory after
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Printf("After GC - Heap Alloc: %d bytes\n", m.HeapAlloc)

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Memory pooling for object reuse")
	fmt.Println("✅ Concurrent processing with worker pools")
	fmt.Println("✅ Caching for repeated operations")
	fmt.Println("✅ Performance benchmarking")
	fmt.Println("✅ Memory usage monitoring")
	fmt.Println("✅ Scalability testing")
}

// generateTestMessages creates test EDIFACT messages
func generateTestMessages(count int) []string {
	messages := make([]string, count)

	for i := 0; i < count; i++ {
		message := fmt.Sprintf(
			"UNH+%d+INVOIC:D:97A:UN'BGM+380+INV%d+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM%d:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+%d'",
			i+1, 1000+i, 1000+i, i+1,
		)
		messages[i] = message
	}

	return messages
}
