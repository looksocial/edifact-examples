package main

import (
	"fmt"
	"strings"
)

// EDIFACTDelimiters represents the standard EDIFACT delimiters
type EDIFACTDelimiters struct {
	SegmentTerminator    string
	DataElementSeparator string
	ComponentSeparator   string
	ReleaseCharacter     string
}

// EDIFACTSegment represents a basic EDIFACT segment
type EDIFACTSegment struct {
	Tag        string
	Elements   []string
	Delimiters EDIFACTDelimiters
}

// NewEDIFACTSegment creates a new EDIFACT segment
func NewEDIFACTSegment(tag string, delimiters EDIFACTDelimiters) *EDIFACTSegment {
	return &EDIFACTSegment{
		Tag:        tag,
		Elements:   make([]string, 0),
		Delimiters: delimiters,
	}
}

// AddElement adds a data element to the segment
func (s *EDIFACTSegment) AddElement(element string) {
	s.Elements = append(s.Elements, element)
}

// AddEmptyElement adds an empty element to the segment
func (s *EDIFACTSegment) AddEmptyElement() {
	s.Elements = append(s.Elements, "")
}

// ToString converts the segment to EDIFACT string format
func (s *EDIFACTSegment) ToString() string {
	var parts []string
	parts = append(parts, s.Tag)

	for _, element := range s.Elements {
		parts = append(parts, element)
	}

	return strings.Join(parts, s.Delimiters.DataElementSeparator) + s.Delimiters.SegmentTerminator
}

// ParseSegment parses an EDIFACT segment string
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
		segment.Elements = append(segment.Elements, parts[i])
	}

	return segment, nil
}

// EscapeSpecialCharacters escapes special characters in EDIFACT data
func EscapeSpecialCharacters(data string, delimiters EDIFACTDelimiters) string {
	result := data
	result = strings.ReplaceAll(result, delimiters.ReleaseCharacter, delimiters.ReleaseCharacter+delimiters.ReleaseCharacter)
	result = strings.ReplaceAll(result, delimiters.DataElementSeparator, delimiters.ReleaseCharacter+delimiters.DataElementSeparator)
	result = strings.ReplaceAll(result, delimiters.ComponentSeparator, delimiters.ReleaseCharacter+delimiters.ComponentSeparator)
	result = strings.ReplaceAll(result, delimiters.SegmentTerminator, delimiters.ReleaseCharacter+delimiters.SegmentTerminator)
	return result
}

// UnescapeSpecialCharacters removes escape sequences from EDIFACT data
func UnescapeSpecialCharacters(data string, delimiters EDIFACTDelimiters) string {
	result := data
	result = strings.ReplaceAll(result, delimiters.ReleaseCharacter+delimiters.SegmentTerminator, delimiters.SegmentTerminator)
	result = strings.ReplaceAll(result, delimiters.ReleaseCharacter+delimiters.ComponentSeparator, delimiters.ComponentSeparator)
	result = strings.ReplaceAll(result, delimiters.ReleaseCharacter+delimiters.DataElementSeparator, delimiters.DataElementSeparator)
	result = strings.ReplaceAll(result, delimiters.ReleaseCharacter+delimiters.ReleaseCharacter, delimiters.ReleaseCharacter)
	return result
}

func main() {
	fmt.Println("🎓 Lesson 1: Basic UN/EDIFACT Syntax & Delimiters")
	fmt.Println("=" * 60)

	// Standard EDIFACT delimiters
	delimiters := EDIFACTDelimiters{
		SegmentTerminator:    "'",
		DataElementSeparator: "+",
		ComponentSeparator:   ":",
		ReleaseCharacter:     "?",
	}

	fmt.Println("\n📋 Standard EDIFACT Delimiters:")
	fmt.Printf("Segment Terminator: '%s'\n", delimiters.SegmentTerminator)
	fmt.Printf("Data Element Separator: '%s'\n", delimiters.DataElementSeparator)
	fmt.Printf("Component Separator: '%s'\n", delimiters.ComponentSeparator)
	fmt.Printf("Release Character: '%s'\n", delimiters.ReleaseCharacter)

	// Example 1: Basic segment construction
	fmt.Println("\n🔧 Example 1: Basic Segment Construction")
	segment1 := NewEDIFACTSegment("UNH", delimiters)
	segment1.AddElement("1")
	segment1.AddElement("INVOIC:D:97A:UN")
	fmt.Printf("Constructed: %s", segment1.ToString())

	// Example 2: Segment with empty elements
	fmt.Println("\n🔧 Example 2: Segment with Empty Elements")
	segment2 := NewEDIFACTSegment("DTM", delimiters)
	segment2.AddElement("137:20231201:102")
	segment2.AddEmptyElement() // Empty element
	fmt.Printf("With empty element: %s", segment2.ToString())

	// Example 3: Segment parsing
	fmt.Println("\n🔧 Example 3: Segment Parsing")
	segmentStr := "UNH+1+INVOIC:D:97A:UN'"
	parsedSegment, err := ParseSegment(segmentStr, delimiters)
	if err != nil {
		fmt.Printf("Error parsing segment: %v\n", err)
	} else {
		fmt.Printf("Parsed segment tag: %s\n", parsedSegment.Tag)
		fmt.Printf("Parsed elements: %v\n", parsedSegment.Elements)
		fmt.Printf("Reconstructed: %s", parsedSegment.ToString())
	}

	// Example 4: Character escaping
	fmt.Println("\n🔧 Example 4: Character Escaping")
	textWithSpecialChars := "This contains a +plus sign and :colon"
	escaped := EscapeSpecialCharacters(textWithSpecialChars, delimiters)
	fmt.Printf("Original: %s\n", textWithSpecialChars)
	fmt.Printf("Escaped: %s\n", escaped)

	// Example 5: Character unescaping
	fmt.Println("\n🔧 Example 5: Character Unescaping")
	unescaped := UnescapeSpecialCharacters(escaped, delimiters)
	fmt.Printf("Unescaped: %s\n", unescaped)

	// Example 6: Composite elements
	fmt.Println("\n🔧 Example 6: Composite Elements")
	compositeSegment := NewEDIFACTSegment("DTM", delimiters)
	compositeSegment.AddElement("137:20231201:102") // Composite element with 3 components
	fmt.Printf("Composite element: %s", compositeSegment.ToString())

	// Parse composite element
	parts := strings.Split(compositeSegment.Elements[0], delimiters.ComponentSeparator)
	fmt.Printf("Components: %v\n", parts)

	// Example 7: Complex segment with escaping
	fmt.Println("\n🔧 Example 7: Complex Segment with Escaping")
	complexSegment := NewEDIFACTSegment("FTX", delimiters)
	complexSegment.AddElement("AAA")
	complexSegment.AddEmptyElement()
	escapedText := EscapeSpecialCharacters("This contains a +plus sign", delimiters)
	complexSegment.AddElement(escapedText)
	fmt.Printf("Complex segment: %s", complexSegment.ToString())

	// Example 8: Syntax validation
	fmt.Println("\n🔧 Example 8: Syntax Validation")
	testSegments := []string{
		"UNH+1+INVOIC:D:97A:UN'",    // Valid
		"UNH+1+INVOIC:D:97A:UN",     // Missing terminator
		"UNH+1+INVOIC:D:97A:UN++",   // Valid with empty elements
		"UNH+1+INVOIC:D:97A:UN+++'", // Valid with multiple empty elements
	}

	for i, segmentStr := range testSegments {
		_, err := ParseSegment(segmentStr, delimiters)
		if err != nil {
			fmt.Printf("Segment %d (%s): INVALID - %v\n", i+1, segmentStr, err)
		} else {
			fmt.Printf("Segment %d (%s): VALID\n", i+1, segmentStr)
		}
	}

	// Example 9: Delimiter analysis
	fmt.Println("\n🔧 Example 9: Delimiter Analysis")
	analyzeSegment := "UNH+1+INVOIC:D:97A:UN'"
	fmt.Printf("Analyzing: %s\n", analyzeSegment)

	// Count delimiters
	plusCount := strings.Count(analyzeSegment, delimiters.DataElementSeparator)
	colonCount := strings.Count(analyzeSegment, delimiters.ComponentSeparator)
	terminatorCount := strings.Count(analyzeSegment, delimiters.SegmentTerminator)

	fmt.Printf("Data element separators (+): %d\n", plusCount)
	fmt.Printf("Component separators (:): %d\n", colonCount)
	fmt.Printf("Segment terminators ('): %d\n", terminatorCount)

	// Example 10: Real-world segment examples
	fmt.Println("\n🔧 Example 10: Real-world Segment Examples")
	realSegments := []string{
		"UNH+1+INVOIC:D:97A:UN'",   // Message header
		"BGM+380+12345678+9'",      // Beginning of message
		"DTM+137:20231201:102'",    // Date/time
		"NAD+BY+++ACME CORP'",      // Name and address
		"LIN+1++1234567890123:EN'", // Line item
		"QTY+12:100:PCE'",          // Quantity
		"PRI+AAA:25.50:CT'",        // Price
		"UNT+8+1'",                 // Message trailer
	}

	for i, segmentStr := range realSegments {
		parsed, err := ParseSegment(segmentStr, delimiters)
		if err == nil {
			fmt.Printf("%d. %s -> Tag: %s, Elements: %d\n",
				i+1, segmentStr, parsed.Tag, len(parsed.Elements))
		}
	}

	fmt.Println("\n🎉 Lesson 1 Complete!")
	fmt.Println("Key takeaways:")
	fmt.Println("- EDIFACT uses specific delimiters to separate data")
	fmt.Println("- Delimiters are crucial for parsing and validation")
	fmt.Println("- Character escaping is essential for special characters")
	fmt.Println("- Understanding delimiters is fundamental to EDIFACT")
}
