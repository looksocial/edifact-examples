package main

import (
	"fmt"
	"strings"
)

// CompositeElement represents a composite data element
type CompositeElement struct {
	Components []string
	Qualifier  string
	Value      string
	Format     string
}

// EDIFACTDelimiters represents the standard EDIFACT delimiters
type EDIFACTDelimiters struct {
	SegmentTerminator    string
	DataElementSeparator string
	ComponentSeparator   string
	ReleaseCharacter     string
}

// ParseCompositeElement parses a composite element string
func ParseCompositeElement(elementStr string, delimiters EDIFACTDelimiters) *CompositeElement {
	components := strings.Split(elementStr, delimiters.ComponentSeparator)

	composite := &CompositeElement{
		Components: components,
	}

	// Extract common patterns
	if len(components) >= 1 {
		composite.Qualifier = components[0]
	}
	if len(components) >= 2 {
		composite.Value = components[1]
	}
	if len(components) >= 3 {
		composite.Format = components[2]
	}

	return composite
}

// GetComponent returns a component by position (1-based)
func (c *CompositeElement) GetComponent(position int) (string, error) {
	if position < 1 || position > len(c.Components) {
		return "", fmt.Errorf("component position %d out of range", position)
	}
	return c.Components[position-1], nil
}

// GetComponentCount returns the number of components
func (c *CompositeElement) GetComponentCount() int {
	return len(c.Components)
}

// IsEmptyComponent checks if a component is empty
func (c *CompositeElement) IsEmptyComponent(position int) bool {
	if position < 1 || position > len(c.Components) {
		return true
	}
	return c.Components[position-1] == ""
}

// ToString converts the composite element back to string
func (c *CompositeElement) ToString(delimiters EDIFACTDelimiters) string {
	return strings.Join(c.Components, delimiters.ComponentSeparator)
}

// Common qualifier definitions
var qualifierDefinitions = map[string]string{
	"137": "Document/message date/time",
	"2":   "Delivery date/time",
	"35":  "Time",
	"12":  "Number of packages",
	"145": "Gross weight",
	"146": "Volume",
	"AAA": "Free text",
	"CT":  "Contract number",
	"EN":  "EAN (European Article Number)",
	"PCE": "Pieces",
	"KGM": "Kilograms",
	"LTR": "Liters",
}

// Common format codes
var formatCodes = map[string]string{
	"102": "CCYYMMDD (Century, Year, Month, Day)",
	"201": "HHMM (Hour, Minute)",
	"203": "HHMMSS (Hour, Minute, Second)",
	"EN":  "EAN (European Article Number)",
	"PCE": "Pieces",
	"KGM": "Kilograms",
	"LTR": "Liters",
}

func main() {
	fmt.Println("🎓 Lesson 3: Composite Elements")
	fmt.Println("=" * 60)

	// Standard EDIFACT delimiters
	delimiters := EDIFACTDelimiters{
		SegmentTerminator:    "'",
		DataElementSeparator: "+",
		ComponentSeparator:   ":",
		ReleaseCharacter:     "?",
	}

	// Example 1: Basic composite element parsing
	fmt.Println("\n🔧 Example 1: Basic Composite Element Parsing")
	basicComposite := "137:20231201:102"
	parsed := ParseCompositeElement(basicComposite, delimiters)

	fmt.Printf("Composite: %s\n", basicComposite)
	fmt.Printf("Components: %v\n", parsed.Components)
	fmt.Printf("Qualifier: %s (%s)\n", parsed.Qualifier, qualifierDefinitions[parsed.Qualifier])
	fmt.Printf("Value: %s\n", parsed.Value)
	fmt.Printf("Format: %s (%s)\n", parsed.Format, formatCodes[parsed.Format])

	// Example 2: Date/Time composite elements
	fmt.Println("\n🔧 Example 2: Date/Time Composite Elements")
	dateTimeExamples := []string{
		"137:20231201:102", // Document date
		"2:20231201:102",   // Delivery date
		"35:1430:201",      // Time
		"35:143045:203",    // Time with seconds
	}

	for i, example := range dateTimeExamples {
		parsed := ParseCompositeElement(example, delimiters)
		fmt.Printf("%d. %s\n", i+1, example)
		fmt.Printf("   Qualifier: %s (%s)\n", parsed.Qualifier, qualifierDefinitions[parsed.Qualifier])
		fmt.Printf("   Value: %s\n", parsed.Value)
		fmt.Printf("   Format: %s (%s)\n", parsed.Format, formatCodes[parsed.Format])
	}

	// Example 3: Quantity composite elements
	fmt.Println("\n🔧 Example 3: Quantity Composite Elements")
	quantityExamples := []string{
		"12:100:PCE",   // Number of packages
		"145:25.5:KGM", // Weight
		"146:10:LTR",   // Volume
	}

	for i, example := range quantityExamples {
		parsed := ParseCompositeElement(example, delimiters)
		fmt.Printf("%d. %s\n", i+1, example)
		fmt.Printf("   Qualifier: %s (%s)\n", parsed.Qualifier, qualifierDefinitions[parsed.Qualifier])
		fmt.Printf("   Value: %s\n", parsed.Value)
		fmt.Printf("   Unit: %s (%s)\n", parsed.Format, formatCodes[parsed.Format])
	}

	// Example 4: Identification composite elements
	fmt.Println("\n🔧 Example 4: Identification Composite Elements")
	identificationExamples := []string{
		"1234567890123:EN", // EAN code
		"CT:123456",        // Contract number
		"AAA:Free text",    // Free text
	}

	for i, example := range identificationExamples {
		parsed := ParseCompositeElement(example, delimiters)
		fmt.Printf("%d. %s\n", i+1, example)
		fmt.Printf("   Identifier: %s\n", parsed.Qualifier)
		fmt.Printf("   Value: %s\n", parsed.Value)
		if parsed.Format != "" {
			fmt.Printf("   Type: %s (%s)\n", parsed.Format, formatCodes[parsed.Format])
		}
	}

	// Example 5: Component access by position
	fmt.Println("\n🔧 Example 5: Component Access by Position")
	complexComposite := "137:20231201:102:EXTRA"
	parsed = ParseCompositeElement(complexComposite, delimiters)

	fmt.Printf("Composite: %s\n", complexComposite)
	fmt.Printf("Total components: %d\n", parsed.GetComponentCount())

	for pos := 1; pos <= parsed.GetComponentCount(); pos++ {
		component, err := parsed.GetComponent(pos)
		if err != nil {
			fmt.Printf("Error accessing component %d: %v\n", pos, err)
		} else {
			emptyStatus := ""
			if parsed.IsEmptyComponent(pos) {
				emptyStatus = " (empty)"
			}
			fmt.Printf("Component %d: %s%s\n", pos, component, emptyStatus)
		}
	}

	// Example 6: Empty components handling
	fmt.Println("\n🔧 Example 6: Empty Components Handling")
	emptyComponentExamples := []string{
		"137:20231201:102", // No empty components
		"137::102",         // Empty value component
		":20231201:102",    // Empty qualifier
		"137:20231201:",    // Empty format
		"137::",            // Multiple empty components
	}

	for i, example := range emptyComponentExamples {
		parsed := ParseCompositeElement(example, delimiters)
		fmt.Printf("%d. %s\n", i+1, example)
		fmt.Printf("   Components: %v\n", parsed.Components)

		for pos := 1; pos <= parsed.GetComponentCount(); pos++ {
			if parsed.IsEmptyComponent(pos) {
				fmt.Printf("   Component %d is empty\n", pos)
			}
		}
	}

	// Example 7: Composite element validation
	fmt.Println("\n🔧 Example 7: Composite Element Validation")
	validationExamples := []string{
		"137:20231201:102", // Valid date
		"137:20231301:102", // Invalid date (month 13)
		"12:100:PCE",       // Valid quantity
		"12:-5:PCE",        // Invalid quantity (negative)
		"999:value:format", // Unknown qualifier
	}

	for i, example := range validationExamples {
		parsed := ParseCompositeElement(example, delimiters)
		fmt.Printf("%d. %s\n", i+1, example)

		// Validate qualifier
		if qualifierDef, exists := qualifierDefinitions[parsed.Qualifier]; exists {
			fmt.Printf("   Qualifier: %s (%s) - VALID\n", parsed.Qualifier, qualifierDef)
		} else {
			fmt.Printf("   Qualifier: %s - UNKNOWN\n", parsed.Qualifier)
		}

		// Validate format
		if formatDef, exists := formatCodes[parsed.Format]; exists {
			fmt.Printf("   Format: %s (%s) - VALID\n", parsed.Format, formatDef)
		} else {
			fmt.Printf("   Format: %s - UNKNOWN\n", parsed.Format)
		}
	}

	// Example 8: Real-world composite elements
	fmt.Println("\n🔧 Example 8: Real-world Composite Elements")
	realWorldExamples := []string{
		"DTM+137:20231201:102'",    // Date/time
		"QTY+12:100:PCE'",          // Quantity
		"PRI+AAA:25.50:CT'",        // Price
		"LIN+1++1234567890123:EN'", // Line item with composite
		"RFF+CT:123456'",           // Reference
	}

	for i, example := range realWorldExamples {
		fmt.Printf("%d. %s\n", i+1, example)

		// Extract composite element from segment
		parts := strings.Split(example, "+")
		if len(parts) > 1 {
			// Find composite elements (containing :)
			for j, part := range parts {
				if strings.Contains(part, ":") {
					composite := ParseCompositeElement(part, delimiters)
					fmt.Printf("   Composite element %d: %s\n", j, part)
					fmt.Printf("   Components: %v\n", composite.Components)
				}
			}
		}
	}

	// Example 9: Composite element reconstruction
	fmt.Println("\n🔧 Example 9: Composite Element Reconstruction")
	reconstructionExamples := []string{
		"137:20231201:102",
		"12:100:PCE",
		"1234567890123:EN",
	}

	for i, example := range reconstructionExamples {
		parsed := ParseCompositeElement(example, delimiters)
		reconstructed := parsed.ToString(delimiters)
		fmt.Printf("%d. Original: %s\n", i+1, example)
		fmt.Printf("   Reconstructed: %s\n", reconstructed)
		fmt.Printf("   Match: %t\n", example == reconstructed)
	}

	// Example 10: Complex composite scenarios
	fmt.Println("\n🔧 Example 10: Complex Composite Scenarios")
	complexScenarios := []string{
		"137:20231201:102:EXTRA:INFO", // Multiple components
		"AAA:This is:complex:text",    // Text with colons
		"12:100:PCE:BOX:10",           // Quantity with packaging info
	}

	for i, example := range complexScenarios {
		parsed := ParseCompositeElement(example, delimiters)
		fmt.Printf("%d. %s\n", i+1, example)
		fmt.Printf("   Component count: %d\n", parsed.GetComponentCount())
		fmt.Printf("   Components: %v\n", parsed.Components)

		// Analyze each component
		for pos := 1; pos <= parsed.GetComponentCount(); pos++ {
			component, _ := parsed.GetComponent(pos)
			fmt.Printf("   Component %d: '%s' (length: %d)\n", pos, component, len(component))
		}
	}

	fmt.Println("\n🎉 Lesson 3 Complete!")
	fmt.Println("Key takeaways:")
	fmt.Println("- Composite elements contain multiple related components")
	fmt.Println("- Components are separated by the ':' delimiter")
	fmt.Println("- Qualifiers define the meaning of composite elements")
	fmt.Println("- Format codes specify how to interpret values")
	fmt.Println("- Understanding composites is essential for data extraction")
}
