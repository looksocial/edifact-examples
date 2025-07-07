package main

import (
	"fmt"
	"strings"
)

// EDIFACTElement represents a data element within a segment
type EDIFACTElement struct {
	Position    int
	Value       string
	IsEmpty     bool
	IsComposite bool
	Components  []string
}

// EDIFACTSegment represents an EDIFACT segment with detailed element analysis
type EDIFACTSegment struct {
	Tag        string
	Elements   []EDIFACTElement
	Delimiters EDIFACTDelimiters
}

// EDIFACTDelimiters represents the standard EDIFACT delimiters
type EDIFACTDelimiters struct {
	SegmentTerminator    string
	DataElementSeparator string
	ComponentSeparator   string
	ReleaseCharacter     string
}

// NewEDIFACTSegment creates a new EDIFACT segment
func NewEDIFACTSegment(tag string, delimiters EDIFACTDelimiters) *EDIFACTSegment {
	return &EDIFACTSegment{
		Tag:        tag,
		Elements:   make([]EDIFACTElement, 0),
		Delimiters: delimiters,
	}
}

// AddElement adds a data element to the segment
func (s *EDIFACTSegment) AddElement(value string) {
	element := EDIFACTElement{
		Position:    len(s.Elements) + 1,
		Value:       value,
		IsEmpty:     value == "",
		IsComposite: strings.Contains(value, s.Delimiters.ComponentSeparator),
	}

	if element.IsComposite {
		element.Components = strings.Split(value, s.Delimiters.ComponentSeparator)
	}

	s.Elements = append(s.Elements, element)
}

// GetElement returns an element by position (1-based)
func (s *EDIFACTSegment) GetElement(position int) (EDIFACTElement, error) {
	if position < 1 || position > len(s.Elements) {
		return EDIFACTElement{}, fmt.Errorf("element position %d out of range", position)
	}
	return s.Elements[position-1], nil
}

// GetElementCount returns the number of data elements (excluding tag)
func (s *EDIFACTSegment) GetElementCount() int {
	return len(s.Elements)
}

// IsServiceSegment checks if this is a service segment
func (s *EDIFACTSegment) IsServiceSegment() bool {
	serviceTags := []string{"UNH", "UNT", "UNS", "UNZ", "UNB", "UNZ", "UNG", "UNE"}
	for _, tag := range serviceTags {
		if s.Tag == tag {
			return true
		}
	}
	return false
}

// ToString converts the segment to EDIFACT string format
func (s *EDIFACTSegment) ToString() string {
	var parts []string
	parts = append(parts, s.Tag)

	for _, element := range s.Elements {
		parts = append(parts, element.Value)
	}

	return strings.Join(parts, s.Delimiters.DataElementSeparator) + s.Delimiters.SegmentTerminator
}

// ParseSegment parses an EDIFACT segment string with detailed element analysis
func ParseSegment(segmentStr string, delimiters EDIFACTDelimiters) (*EDIFACTSegment, error) {
	// Remove segment terminator
	segmentStr = strings.TrimSuffix(segmentStr, delimiters.SegmentTerminator)

	// Split by data element separator
	parts := strings.Split(segmentStr, delimiters.DataElementSeparator)

	if len(parts) == 0 {
		return nil, fmt.Errorf("empty segment")
	}

	segment := NewEDIFACTSegment(parts[0], delimiters)

	// Add remaining elements
	for i := 1; i < len(parts); i++ {
		segment.AddElement(parts[i])
	}

	return segment, nil
}

// AnalyzeSegment provides detailed analysis of a segment
func AnalyzeSegment(segment *EDIFACTSegment) {
	fmt.Printf("\n📊 Segment Analysis: %s\n", segment.Tag)
	fmt.Printf("Type: %s\n", func() string {
		if segment.IsServiceSegment() {
			return "Service Segment"
		}
		return "Data Segment"
	}())
	fmt.Printf("Total Elements: %d\n", segment.GetElementCount())

	fmt.Println("\nElements:")
	for i, element := range segment.Elements {
		fmt.Printf("  %d. Position %d: ", i+1, element.Position)
		if element.IsEmpty {
			fmt.Printf("[EMPTY]\n")
		} else if element.IsComposite {
			fmt.Printf("Composite: %s (Components: %v)\n", element.Value, element.Components)
		} else {
			fmt.Printf("Simple: %s\n", element.Value)
		}
	}
}

func main() {
	fmt.Println("🎓 Lesson 2: Segments & Elements")
	fmt.Println("=" * 60)

	// Standard EDIFACT delimiters
	delimiters := EDIFACTDelimiters{
		SegmentTerminator:    "'",
		DataElementSeparator: "+",
		ComponentSeparator:   ":",
		ReleaseCharacter:     "?",
	}

	// Example 1: Service segment analysis
	fmt.Println("\n🔧 Example 1: Service Segment Analysis")
	serviceSegmentStr := "UNH+1+INVOIC:D:97A:UN'"
	serviceSegment, _ := ParseSegment(serviceSegmentStr, delimiters)
	AnalyzeSegment(serviceSegment)

	// Example 2: Data segment with simple elements
	fmt.Println("\n🔧 Example 2: Data Segment with Simple Elements")
	dataSegmentStr := "BGM+380+12345678+9'"
	dataSegment, _ := ParseSegment(dataSegmentStr, delimiters)
	AnalyzeSegment(dataSegment)

	// Example 3: Segment with composite elements
	fmt.Println("\n🔧 Example 3: Segment with Composite Elements")
	compositeSegmentStr := "DTM+137:20231201:102'"
	compositeSegment, _ := ParseSegment(compositeSegmentStr, delimiters)
	AnalyzeSegment(compositeSegment)

	// Example 4: Complex segment with mixed elements
	fmt.Println("\n🔧 Example 4: Complex Segment with Mixed Elements")
	complexSegmentStr := "NAD+BY+++ACME CORP+123 MAIN ST+CITY+ST+12345+US'"
	complexSegment, _ := ParseSegment(complexSegmentStr, delimiters)
	AnalyzeSegment(complexSegment)

	// Example 5: Line item segment
	fmt.Println("\n🔧 Example 5: Line Item Segment")
	lineSegmentStr := "LIN+1++1234567890123:EN'"
	lineSegment, _ := ParseSegment(lineSegmentStr, delimiters)
	AnalyzeSegment(lineSegment)

	// Example 6: Element positioning demonstration
	fmt.Println("\n🔧 Example 6: Element Positioning")
	positionSegmentStr := "QTY+12:100:PCE'"
	positionSegment, _ := ParseSegment(positionSegmentStr, delimiters)

	fmt.Printf("Segment: %s\n", positionSegmentStr)
	for i, element := range positionSegment.Elements {
		fmt.Printf("Element %d (Position %d): %s\n", i+1, element.Position, element.Value)
	}

	// Example 7: Accessing elements by position
	fmt.Println("\n🔧 Example 7: Accessing Elements by Position")
	accessSegmentStr := "PRI+AAA:25.50:CT'"
	accessSegment, _ := ParseSegment(accessSegmentStr, delimiters)

	for pos := 1; pos <= accessSegment.GetElementCount(); pos++ {
		element, err := accessSegment.GetElement(pos)
		if err != nil {
			fmt.Printf("Error accessing element %d: %v\n", pos, err)
		} else {
			fmt.Printf("Position %d: %s (Composite: %t)\n", pos, element.Value, element.IsComposite)
		}
	}

	// Example 8: Segment type identification
	fmt.Println("\n🔧 Example 8: Segment Type Identification")
	testSegments := []string{
		"UNH+1+INVOIC:D:97A:UN'", // Service
		"UNT+8+1'",               // Service
		"BGM+380+12345678+9'",    // Data
		"DTM+137:20231201:102'",  // Data
		"NAD+BY+++ACME CORP'",    // Data
	}

	for _, segmentStr := range testSegments {
		segment, _ := ParseSegment(segmentStr, delimiters)
		segmentType := "Data"
		if segment.IsServiceSegment() {
			segmentType = "Service"
		}
		fmt.Printf("%s -> %s Segment\n", segment.Tag, segmentType)
	}

	// Example 9: Element counting and validation
	fmt.Println("\n🔧 Example 9: Element Counting and Validation")
	countingSegments := []string{
		"UNH+1+INVOIC:D:97A:UN'",          // 2 elements
		"BGM+380+12345678+9'",             // 3 elements
		"DTM+137:20231201:102'",           // 1 composite element
		"NAD+BY+++ACME CORP+123 MAIN ST'", // 5 elements (2 empty)
	}

	for _, segmentStr := range countingSegments {
		segment, _ := ParseSegment(segmentStr, delimiters)
		emptyCount := 0
		compositeCount := 0

		for _, element := range segment.Elements {
			if element.IsEmpty {
				emptyCount++
			}
			if element.IsComposite {
				compositeCount++
			}
		}

		fmt.Printf("%s: %d total elements, %d empty, %d composite\n",
			segment.Tag, segment.GetElementCount(), emptyCount, compositeCount)
	}

	// Example 10: Real-world segment examples
	fmt.Println("\n🔧 Example 10: Real-world Segment Examples")
	realSegments := []string{
		"UNH+1+INVOIC:D:97A:UN'",                           // Message header
		"BGM+380+12345678+9'",                              // Beginning of message
		"DTM+137:20231201:102'",                            // Date/time
		"NAD+BY+++ACME CORP+123 MAIN ST+CITY+ST+12345+US'", // Name and address
		"LIN+1++1234567890123:EN'",                         // Line item
		"QTY+12:100:PCE'",                                  // Quantity
		"PRI+AAA:25.50:CT'",                                // Price
		"UNT+8+1'",                                         // Message trailer
	}

	fmt.Println("\nDetailed Analysis:")
	for i, segmentStr := range realSegments {
		segment, _ := ParseSegment(segmentStr, delimiters)
		fmt.Printf("\n%d. %s\n", i+1, segmentStr)
		fmt.Printf("   Tag: %s, Elements: %d\n", segment.Tag, segment.GetElementCount())

		for j, element := range segment.Elements {
			elementType := "Simple"
			if element.IsComposite {
				elementType = "Composite"
			} else if element.IsEmpty {
				elementType = "Empty"
			}
			fmt.Printf("   Element %d: %s (%s)\n", j+1, element.Value, elementType)
		}
	}

	fmt.Println("\n🎉 Lesson 2 Complete!")
	fmt.Println("Key takeaways:")
	fmt.Println("- Segments are the building blocks of EDIFACT messages")
	fmt.Println("- Elements have specific positions and types")
	fmt.Println("- Empty elements must be properly handled")
	fmt.Println("- Composite elements contain multiple components")
	fmt.Println("- Understanding segments is crucial for message processing")
}
