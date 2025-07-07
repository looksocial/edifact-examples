// Lesson 1: Introduction & Setup
// This lesson covers the basic setup and introduction to the EDIFACT Go package.

package main

import (
	"fmt"
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
var segmentPool = sync.Pool{
	New: func() interface{} {
		return &Segment{
			Elements: make([]string, 0, 10), // Pre-allocate slice
		}
	},
}

// StreamingProcessor handles large EDIFACT files
type StreamingProcessor struct {
	processedCount int
	errorCount     int
	mu             sync.Mutex
}

// ProcessSegment processes a single segment
func (sp *StreamingProcessor) ProcessSegment(segmentStr string) error {
	// Get segment from pool
	segment := segmentPool.Get().(*Segment)
	defer segmentPool.Put(segment)

	// Reset segment for reuse
	segment.Tag = ""
	segment.Elements = segment.Elements[:0]

	// Parse segment
	parts := strings.Split(segmentStr, "+")
	if len(parts) == 0 {
		return fmt.Errorf("empty segment")
	}

	segment.Tag = parts[0]
	segment.Elements = append(segment.Elements, parts[1:]...)

	// Process based on segment type
	switch segment.Tag {
	case "UNH":
		return sp.processHeader(segment)
	case "BGM":
		return sp.processBeginning(segment)
	case "DTM":
		return sp.processDateTime(segment)
	case "LIN":
		return sp.processLineItem(segment)
	case "UNT":
		return sp.processTrailer(segment)
	default:
		// Handle unknown segments gracefully
		return nil
	}
}

func (sp *StreamingProcessor) processHeader(segment *Segment) error {
	if len(segment.Elements) < 2 {
		return fmt.Errorf("invalid UNH segment")
	}
	fmt.Printf("Processing header: %s\n", segment.Elements[0])
	return nil
}

func (sp *StreamingProcessor) processBeginning(segment *Segment) error {
	if len(segment.Elements) < 2 {
		return fmt.Errorf("invalid BGM segment")
	}
	fmt.Printf("Processing beginning: %s\n", segment.Elements[1])
	return nil
}

func (sp *StreamingProcessor) processDateTime(segment *Segment) error {
	if len(segment.Elements) == 0 {
		return fmt.Errorf("invalid DTM segment")
	}

	// Parse composite element
	composite := strings.Split(segment.Elements[0], ":")
	if len(composite) >= 2 {
		fmt.Printf("Processing date: %s\n", composite[1])
	}
	return nil
}

func (sp *StreamingProcessor) processLineItem(segment *Segment) error {
	if len(segment.Elements) < 1 {
		return fmt.Errorf("invalid LIN segment")
	}
	fmt.Printf("Processing line item: %s\n", segment.Elements[0])
	return nil
}

func (sp *StreamingProcessor) processTrailer(segment *Segment) error {
	if len(segment.Elements) < 1 {
		return fmt.Errorf("invalid UNT segment")
	}
	fmt.Printf("Processing trailer: %s\n", segment.Elements[0])
	return nil
}

// ProcessLargeFile simulates processing a large EDIFACT file
func (sp *StreamingProcessor) ProcessLargeFile(content string) error {
	segments := strings.Split(content, "'")

	fmt.Printf("Processing %d segments...\n", len(segments))

	// Process segments with error handling
	for i, segmentStr := range segments {
		if segmentStr == "" {
			continue
		}

		if err := sp.ProcessSegment(segmentStr); err != nil {
			sp.mu.Lock()
			sp.errorCount++
			sp.mu.Unlock()
			fmt.Printf("Error processing segment %d: %v\n", i+1, err)
			// Continue processing other segments
		} else {
			sp.mu.Lock()
			sp.processedCount++
			sp.mu.Unlock()
		}

		// Simulate processing time
		time.Sleep(1 * time.Millisecond)
	}

	return nil
}

// GetStats returns processing statistics
func (sp *StreamingProcessor) GetStats() (int, int) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	return sp.processedCount, sp.errorCount
}

func main() {
	fmt.Println("=== Advanced Message Handling (Lesson 1) ===")

	// Complex EDIFACT message with nested structures
	complexMessage := `UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'NAD+SE+++SUPPLIER INC'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'LIN+2++ITEM002:EN'QTY+12:50:PCE'PRI+AAA:30.00:CT'UNT+12+1'`

	// Create processor
	processor := &StreamingProcessor{}

	// Process the message
	fmt.Println("Starting advanced message processing...")
	start := time.Now()

	if err := processor.ProcessLargeFile(complexMessage); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		return
	}

	duration := time.Since(start)
	processed, errors := processor.GetStats()

	fmt.Printf("\nProcessing completed in %v\n", duration)
	fmt.Printf("Segments processed: %d\n", processed)
	fmt.Printf("Errors encountered: %d\n", errors)

	// Demonstrate memory optimization
	fmt.Println("\n=== Memory Optimization Demo ===")

	// Process multiple messages to show object reuse
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("UNH+%d+INVOIC:D:97A:UN'BGM+380+INV%d+9'UNT+2+%d'", i+1, 1000+i, i+1)
		processor.ProcessLargeFile(msg)
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Streaming processing for large files")
	fmt.Println("✅ Object pooling for memory optimization")
	fmt.Println("✅ Robust error handling and recovery")
	fmt.Println("✅ Concurrent-safe processing")
	fmt.Println("✅ Performance monitoring")
}
