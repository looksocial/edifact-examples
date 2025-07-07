// Lesson 7: Data Element Types
// This lesson demonstrates EDIFACT data element types and their classification.

package main

import (
	"fmt"
	"regexp"
	"strings"
)

// EDIFACTMessage represents a parsed EDIFACT message
type EDIFACTMessage struct {
	MessageType string
	Segments    []*Segment
	RawContent  string
}

// Segment represents an EDIFACT segment
type Segment struct {
	Tag      string
	Elements []*DataElement
	Position int
}

// DataElement represents an EDIFACT data element
type DataElement struct {
	Value       string
	Type        ElementType
	Qualifier   string
	Components  []string
	Position    int
	Description string
}

// ElementType represents the type of a data element
type ElementType int

const (
	TypeQualifier ElementType = iota
	TypeCode
	TypeMeasurement
	TypeDate
	TypeText
	TypeIdentifier
	TypeComposite
)

// DataElementAnalyzer analyzes data elements in EDIFACT messages
type DataElementAnalyzer struct {
	qualifiers map[string]string
	codes      map[string]map[string]string
	units      map[string]string
	patterns   map[ElementType]*regexp.Regexp
}

// NewDataElementAnalyzer creates a new data element analyzer
func NewDataElementAnalyzer() *DataElementAnalyzer {
	dea := &DataElementAnalyzer{
		qualifiers: make(map[string]string),
		codes:      make(map[string]map[string]string),
		units:      make(map[string]string),
		patterns:   make(map[ElementType]*regexp.Regexp),
	}

	// Initialize patterns and code lists
	dea.initializePatterns()
	dea.initializeCodeLists()

	return dea
}

// initializePatterns initializes regex patterns for element types
func (dea *DataElementAnalyzer) initializePatterns() {
	dea.patterns[TypeDate] = regexp.MustCompile(`^\d{8}$`)
	dea.patterns[TypeMeasurement] = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	dea.patterns[TypeCode] = regexp.MustCompile(`^[A-Z0-9]{2,6}$`)
	dea.patterns[TypeIdentifier] = regexp.MustCompile(`^[A-Z0-9]{1,35}$`)
}

// initializeCodeLists initializes code lists and qualifiers
func (dea *DataElementAnalyzer) initializeCodeLists() {
	// Party qualifiers
	dea.qualifiers["BY"] = "Buyer"
	dea.qualifiers["SE"] = "Seller"
	dea.qualifiers["SU"] = "Supplier"
	dea.qualifiers["CA"] = "Carrier"
	dea.qualifiers["XX"] = "Unknown"

	// Date/time qualifiers
	dea.qualifiers["137"] = "Document/message date/time"
	dea.qualifiers["2"] = "Delivery date/time"
	dea.qualifiers["35"] = "Time"

	// Quantity qualifiers
	dea.qualifiers["12"] = "Number of packages"
	dea.qualifiers["145"] = "Gross weight"
	dea.qualifiers["146"] = "Volume"

	// Reference qualifiers
	dea.qualifiers["CT"] = "Contract number"
	dea.qualifiers["IV"] = "Invoice number"
	dea.qualifiers["PO"] = "Purchase order number"

	// Message type codes
	dea.codes["message_type"] = map[string]string{
		"INVOIC": "Invoice",
		"ORDERS": "Purchase Order",
		"DESADV": "Despatch Advice",
		"INVRPT": "Inventory Report",
	}

	// Document type codes
	dea.codes["document_type"] = map[string]string{
		"380": "Commercial Invoice",
		"325": "Pro-forma Invoice",
		"384": "Consignment Invoice",
		"220": "Purchase Order",
	}

	// Unit codes
	dea.units["PCE"] = "Pieces"
	dea.units["KGM"] = "Kilograms"
	dea.units["LTR"] = "Liters"
	dea.units["MTR"] = "Meters"
	dea.units["CT"] = "Cartons"
}

// AnalyzeMessage analyzes data elements in an EDIFACT message
func (dea *DataElementAnalyzer) AnalyzeMessage(message *EDIFACTMessage) *ElementAnalysisResult {
	result := &ElementAnalysisResult{
		MessageType: message.MessageType,
		Elements:    []*ElementAnalysis{},
		Statistics:  make(map[ElementType]int),
		Errors:      []string{},
		Warnings:    []string{},
	}

	// Analyze each segment
	for _, segment := range message.Segments {
		for i, element := range segment.Elements {
			analysis := dea.analyzeElement(element, segment.Tag, i)
			result.Elements = append(result.Elements, analysis)
			result.Statistics[analysis.Type]++

			// Check for validation issues
			if analysis.ValidationError != "" {
				result.Errors = append(result.Errors, analysis.ValidationError)
			}
			if analysis.Warning != "" {
				result.Warnings = append(result.Warnings, analysis.Warning)
			}
		}
	}

	return result
}

// analyzeElement analyzes a single data element
func (dea *DataElementAnalyzer) analyzeElement(element *DataElement, segmentTag string, position int) *ElementAnalysis {
	analysis := &ElementAnalysis{
		Segment:    segmentTag,
		Position:   position,
		Value:      element.Value,
		Type:       element.Type,
		Qualifier:  element.Qualifier,
		Components: element.Components,
	}

	// Determine element type based on context and content
	analysis.Type = dea.determineElementType(element, segmentTag, position)

	// Extract qualifier if present
	analysis.Qualifier = dea.extractQualifier(element, segmentTag, position)

	// Validate element based on type
	analysis.ValidationError = dea.validateElement(element, analysis.Type, segmentTag)

	// Check for warnings
	analysis.Warning = dea.checkWarnings(element, analysis.Type, segmentTag)

	// Get description
	analysis.Description = dea.getElementDescription(element, analysis.Type, segmentTag)

	return analysis
}

// determineElementType determines the type of a data element
func (dea *DataElementAnalyzer) determineElementType(element *DataElement, segmentTag string, position int) ElementType {
	value := element.Value

	// Check for composite elements (containing ':')
	if strings.Contains(value, ":") {
		return TypeComposite
	}

	// Check for qualifiers based on segment context
	if dea.isQualifierContext(segmentTag, position) {
		return TypeQualifier
	}

	// Check for codes
	if dea.patterns[TypeCode].MatchString(value) {
		return TypeCode
	}

	// Check for dates
	if dea.patterns[TypeDate].MatchString(value) {
		return TypeDate
	}

	// Check for measurements
	if dea.patterns[TypeMeasurement].MatchString(value) {
		return TypeMeasurement
	}

	// Check for identifiers
	if dea.patterns[TypeIdentifier].MatchString(value) {
		return TypeIdentifier
	}

	// Default to text
	return TypeText
}

// isQualifierContext checks if the element is in a qualifier context
func (dea *DataElementAnalyzer) isQualifierContext(segmentTag string, position int) bool {
	switch segmentTag {
	case "NAD":
		return position == 0 // First element is party qualifier
	case "DTM":
		return position == 0 // First element is date qualifier
	case "QTY":
		return position == 0 // First element is quantity qualifier
	case "RFF":
		return position == 0 // First element is reference qualifier
	}
	return false
}

// extractQualifier extracts qualifier information from an element
func (dea *DataElementAnalyzer) extractQualifier(element *DataElement, segmentTag string, position int) string {
	if dea.isQualifierContext(segmentTag, position) {
		// For composite elements, qualifier is the first component
		if len(element.Components) > 0 {
			return element.Components[0]
		}
		return element.Value
	}
	return ""
}

// validateElement validates an element based on its type
func (dea *DataElementAnalyzer) validateElement(element *DataElement, elementType ElementType, segmentTag string) string {
	value := element.Value

	switch elementType {
	case TypeQualifier:
		if qualifier, exists := dea.qualifiers[value]; !exists {
			return fmt.Sprintf("Unknown qualifier: %s", value)
		}

	case TypeCode:
		if segmentTag == "BGM" && position == 1 {
			if docType, exists := dea.codes["document_type"][value]; !exists {
				return fmt.Sprintf("Unknown document type code: %s", value)
			}
		}

	case TypeDate:
		if !dea.patterns[TypeDate].MatchString(value) {
			return fmt.Sprintf("Invalid date format: %s (expected YYYYMMDD)", value)
		}

	case TypeMeasurement:
		if !dea.patterns[TypeMeasurement].MatchString(value) {
			return fmt.Sprintf("Invalid numeric format: %s", value)
		}

	case TypeComposite:
		if len(element.Components) == 0 {
			return "Composite element should have components"
		}
	}

	return ""
}

// checkWarnings checks for warnings on an element
func (dea *DataElementAnalyzer) checkWarnings(element *DataElement, elementType ElementType, segmentTag string) string {
	value := element.Value

	switch elementType {
	case TypeMeasurement:
		if strings.HasPrefix(value, "-") {
			return fmt.Sprintf("Negative value detected: %s", value)
		}

	case TypeText:
		if len(value) > 35 {
			return fmt.Sprintf("Text element exceeds recommended length: %d characters", len(value))
		}
	}

	return ""
}

// getElementDescription gets a description for an element
func (dea *DataElementAnalyzer) getElementDescription(element *DataElement, elementType ElementType, segmentTag string) string {
	value := element.Value

	switch elementType {
	case TypeQualifier:
		if description, exists := dea.qualifiers[value]; exists {
			return description
		}
		return "Unknown qualifier"

	case TypeCode:
		if segmentTag == "BGM" {
			if description, exists := dea.codes["document_type"][value]; exists {
				return description
			}
		}
		if description, exists := dea.codes["message_type"][value]; exists {
			return description
		}
		return "Code"

	case TypeDate:
		return "Date (YYYYMMDD format)"

	case TypeMeasurement:
		return "Numeric measurement"

	case TypeIdentifier:
		return "Identifier"

	case TypeComposite:
		return "Composite element"

	case TypeText:
		return "Text element"
	}

	return "Unknown type"
}

// ElementAnalysis contains analysis results for a single element
type ElementAnalysis struct {
	Segment         string
	Position        int
	Value           string
	Type            ElementType
	Qualifier       string
	Components      []string
	Description     string
	ValidationError string
	Warning         string
}

// ElementAnalysisResult contains overall analysis results
type ElementAnalysisResult struct {
	MessageType string
	Elements    []*ElementAnalysis
	Statistics  map[ElementType]int
	Errors      []string
	Warnings    []string
}

// PrintAnalysis prints the analysis results
func (ear *ElementAnalysisResult) PrintAnalysis() {
	fmt.Printf("\n=== Data Element Analysis ===\n")
	fmt.Printf("Message Type: %s\n", ear.MessageType)
	fmt.Printf("Total Elements: %d\n", len(ear.Elements))

	// Print statistics
	fmt.Printf("\nElement Type Statistics:\n")
	typeNames := map[ElementType]string{
		TypeQualifier:   "Qualifiers",
		TypeCode:        "Codes",
		TypeMeasurement: "Measurements",
		TypeDate:        "Dates",
		TypeText:        "Text",
		TypeIdentifier:  "Identifiers",
		TypeComposite:   "Composite",
	}

	for elementType, count := range ear.Statistics {
		if name, exists := typeNames[elementType]; exists {
			fmt.Printf("  %s: %d\n", name, count)
		}
	}

	// Print element details
	fmt.Printf("\nElement Details:\n")
	for _, element := range ear.Elements {
		fmt.Printf("  [%s:%d] %s (%s)", element.Segment, element.Position, element.Value, typeNames[element.Type])
		if element.Qualifier != "" {
			fmt.Printf(" - Qualifier: %s", element.Qualifier)
		}
		if element.Description != "" {
			fmt.Printf(" - %s", element.Description)
		}
		fmt.Printf("\n")

		if element.ValidationError != "" {
			fmt.Printf("    ❌ Error: %s\n", element.ValidationError)
		}
		if element.Warning != "" {
			fmt.Printf("    ⚠️  Warning: %s\n", element.Warning)
		}
	}

	// Print summary
	if len(ear.Errors) > 0 {
		fmt.Printf("\nValidation Errors:\n")
		for _, err := range ear.Errors {
			fmt.Printf("  ❌ %s\n", err)
		}
	}

	if len(ear.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range ear.Warnings {
			fmt.Printf("  ⚠️  %s\n", warning)
		}
	}
}

// ParseEDIFACTMessage parses a raw EDIFACT message
func ParseEDIFACTMessage(rawContent string) *EDIFACTMessage {
	message := &EDIFACTMessage{
		RawContent: rawContent,
		Segments:   []*Segment{},
	}

	segments := strings.Split(rawContent, "'")
	for i, segmentStr := range segments {
		if segmentStr == "" {
			continue
		}

		parts := strings.Split(segmentStr, "+")
		if len(parts) == 0 {
			continue
		}

		segment := &Segment{
			Tag:      parts[0],
			Elements: []*DataElement{},
			Position: i,
		}

		// Parse elements
		for j, elementStr := range parts[1:] {
			element := &DataElement{
				Value:    elementStr,
				Position: j,
			}

			// Check if element is composite (contains ':')
			if strings.Contains(elementStr, ":") {
				element.Type = TypeComposite
				element.Components = strings.Split(elementStr, ":")
			} else {
				element.Value = elementStr
			}

			segment.Elements = append(segment.Elements, element)
		}

		message.Segments = append(message.Segments, segment)

		// Extract message type from UNH segment
		if segment.Tag == "UNH" && len(segment.Elements) > 1 {
			msgTypeParts := strings.Split(segment.Elements[1].Value, ":")
			if len(msgTypeParts) > 0 {
				message.MessageType = msgTypeParts[0]
			}
		}
	}

	return message
}

func main() {
	fmt.Println("=== EDIFACT Data Element Types (Lesson 7) ===")

	// Test messages with various element types
	messages := []string{
		// Valid message with various element types
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+1'`,

		// Message with unknown qualifier
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12346+9'DTM+137:20231201:102'NAD+XX+++UNKNOWN COMPANY'LIN+1++ITEM001:EN'QTY+12:100:PCE'UNT+6+1'`,

		// Message with invalid date
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12347+9'DTM+137:20231301:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'UNT+6+1'`,

		// Message with negative quantity
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12348+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:-5:PCE'PRI+AAA:25.50:CT'UNT+7+1'`,
	}

	// Create analyzer
	analyzer := NewDataElementAnalyzer()

	// Analyze each message
	for i, rawMessage := range messages {
		fmt.Printf("\n--- Analyzing Message %d ---\n", i+1)
		fmt.Printf("Message: %s\n", rawMessage)

		// Parse message
		message := ParseEDIFACTMessage(rawMessage)

		// Analyze elements
		result := analyzer.AnalyzeMessage(message)
		result.PrintAnalysis()
	}

	// Demonstrate element type classification
	fmt.Println("\n=== Element Type Classification Demo ===")

	testMessage := `UNH+1+INVOIC:D:97A:UN'BGM+380+INV12349+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+1'`
	message := ParseEDIFACTMessage(testMessage)

	// Show element type classification
	fmt.Printf("Element Type Classification:\n")
	for _, segment := range message.Segments {
		fmt.Printf("  %s segment:\n", segment.Tag)
		for i, element := range segment.Elements {
			analysis := analyzer.analyzeElement(element, segment.Tag, i)
			typeNames := map[ElementType]string{
				TypeQualifier:   "Qualifier",
				TypeCode:        "Code",
				TypeMeasurement: "Measurement",
				TypeDate:        "Date",
				TypeText:        "Text",
				TypeIdentifier:  "Identifier",
				TypeComposite:   "Composite",
			}
			fmt.Printf("    Element %d: %s (%s)\n", i+1, element.Value, typeNames[analysis.Type])
		}
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Data element type classification")
	fmt.Println("✅ Qualifier identification and validation")
	fmt.Println("✅ Code list validation")
	fmt.Println("✅ Date and numeric format validation")
	fmt.Println("✅ Composite element handling")
	fmt.Println("✅ Element context analysis")
	fmt.Println("✅ Comprehensive error reporting")
}
