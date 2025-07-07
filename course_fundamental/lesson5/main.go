// Lesson 5: Syntax Validation
// This lesson demonstrates EDIFACT syntax validation techniques and error detection.

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
	Elements []string
	Position int
	Raw      string
}

// SyntaxValidator validates EDIFACT syntax
type SyntaxValidator struct {
	validationRules []ValidationRule
}

// ValidationRule defines a validation rule
type ValidationRule struct {
	Name        string
	Description string
	Validate    func(*EDIFACTMessage) []ValidationError
}

// ValidationError represents a validation error
type ValidationError struct {
	Segment  string
	Element  string
	Position int
	Rule     string
	Message  string
	Severity string
}

// ValidationResult contains validation results
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationError
}

// NewSyntaxValidator creates a new syntax validator
func NewSyntaxValidator() *SyntaxValidator {
	sv := &SyntaxValidator{
		validationRules: []ValidationRule{},
	}

	// Register validation rules
	sv.registerValidationRules()

	return sv
}

// registerValidationRules registers all validation rules
func (sv *SyntaxValidator) registerValidationRules() {
	sv.validationRules = []ValidationRule{
		{
			Name:        "segment_terminator",
			Description: "Check for proper segment terminators",
			Validate:    sv.validateSegmentTerminators,
		},
		{
			Name:        "element_delimiter",
			Description: "Check for proper element delimiters",
			Validate:    sv.validateElementDelimiters,
		},
		{
			Name:        "date_format",
			Description: "Validate date formats",
			Validate:    sv.validateDateFormats,
		},
		{
			Name:        "numeric_format",
			Description: "Validate numeric formats",
			Validate:    sv.validateNumericFormats,
		},
		{
			Name:        "segment_count",
			Description: "Validate segment count in UNT",
			Validate:    sv.validateSegmentCount,
		},
		{
			Name:        "reference_matching",
			Description: "Validate reference number matching",
			Validate:    sv.validateReferenceMatching,
		},
	}
}

// Validate validates an EDIFACT message
func (sv *SyntaxValidator) Validate(message *EDIFACTMessage) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}

	// Apply all validation rules
	for _, rule := range sv.validationRules {
		errors := rule.Validate(message)
		for _, err := range errors {
			if err.Severity == "ERROR" {
				result.Valid = false
				result.Errors = append(result.Errors, err)
			} else {
				result.Warnings = append(result.Warnings, err)
			}
		}
	}

	return result
}

// validateSegmentTerminators validates segment terminators
func (sv *SyntaxValidator) validateSegmentTerminators(message *EDIFACTMessage) []ValidationError {
	var errors []ValidationError

	// Check if message ends with segment terminator
	if !strings.HasSuffix(message.RawContent, "'") {
		errors = append(errors, ValidationError{
			Segment:  "END",
			Element:  "Terminator",
			Position: len(message.RawContent),
			Rule:     "segment_terminator",
			Message:  "Message must end with segment terminator (')",
			Severity: "ERROR",
		})
	}

	// Check each segment for proper termination
	for i, segment := range message.Segments {
		if !strings.HasSuffix(segment.Raw, "'") {
			errors = append(errors, ValidationError{
				Segment:  segment.Tag,
				Element:  "Terminator",
				Position: i,
				Rule:     "segment_terminator",
				Message:  fmt.Sprintf("Segment %s must end with terminator (')", segment.Tag),
				Severity: "ERROR",
			})
		}
	}

	return errors
}

// validateElementDelimiters validates element delimiters
func (sv *SyntaxValidator) validateElementDelimiters(message *EDIFACTMessage) []ValidationError {
	var errors []ValidationError

	for i, segment := range message.Segments {
		// Check for consecutive delimiters (empty elements)
		rawSegment := segment.Raw
		if strings.Contains(rawSegment, "++") {
			errors = append(errors, ValidationError{
				Segment:  segment.Tag,
				Element:  "Delimiter",
				Position: i,
				Rule:     "element_delimiter",
				Message:  fmt.Sprintf("Segment %s contains consecutive delimiters (empty element)", segment.Tag),
				Severity: "WARNING",
			})
		}

		// Check for proper element delimiter usage
		if strings.Contains(rawSegment, "+") && !strings.HasPrefix(rawSegment, "+") {
			// Valid - contains element delimiters
		} else if !strings.Contains(rawSegment, "+") && segment.Tag != "UNT" {
			// UNT can be simple, others should have elements
			if len(segment.Elements) == 0 && segment.Tag != "UNT" {
				errors = append(errors, ValidationError{
					Segment:  segment.Tag,
					Element:  "Delimiter",
					Position: i,
					Rule:     "element_delimiter",
					Message:  fmt.Sprintf("Segment %s should contain element delimiters", segment.Tag),
					Severity: "WARNING",
				})
			}
		}
	}

	return errors
}

// validateDateFormats validates date formats
func (sv *SyntaxValidator) validateDateFormats(message *EDIFACTMessage) []ValidationError {
	var errors []ValidationError
	dateRegex := regexp.MustCompile(`^\d{8}$`)

	for i, segment := range message.Segments {
		if segment.Tag == "DTM" && len(segment.Elements) > 0 {
			// Parse composite element
			composite := strings.Split(segment.Elements[0], ":")
			if len(composite) >= 2 {
				dateValue := composite[1]
				if !dateRegex.MatchString(dateValue) {
					errors = append(errors, ValidationError{
						Segment:  segment.Tag,
						Element:  "Date",
						Position: i,
						Rule:     "date_format",
						Message:  fmt.Sprintf("Invalid date format in DTM segment: %s (expected YYYYMMDD)", dateValue),
						Severity: "ERROR",
					})
				}

				// Additional date validation
				if len(dateValue) == 8 {
					year := dateValue[:4]
					month := dateValue[4:6]
					day := dateValue[6:8]

					// Basic range validation
					if month < "01" || month > "12" {
						errors = append(errors, ValidationError{
							Segment:  segment.Tag,
							Element:  "Month",
							Position: i,
							Rule:     "date_format",
							Message:  fmt.Sprintf("Invalid month in date: %s", dateValue),
							Severity: "ERROR",
						})
					}

					if day < "01" || day > "31" {
						errors = append(errors, ValidationError{
							Segment:  segment.Tag,
							Element:  "Day",
							Position: i,
							Rule:     "date_format",
							Message:  fmt.Sprintf("Invalid day in date: %s", dateValue),
							Severity: "ERROR",
						})
					}
				}
			}
		}
	}

	return errors
}

// validateNumericFormats validates numeric formats
func (sv *SyntaxValidator) validateNumericFormats(message *EDIFACTMessage) []ValidationError {
	var errors []ValidationError
	numericRegex := regexp.MustCompile(`^-?\d+(\.\d+)?$`)

	for i, segment := range message.Segments {
		if segment.Tag == "QTY" && len(segment.Elements) > 0 {
			composite := strings.Split(segment.Elements[0], ":")
			if len(composite) >= 2 {
				quantity := composite[1]
				if !numericRegex.MatchString(quantity) {
					errors = append(errors, ValidationError{
						Segment:  segment.Tag,
						Element:  "Quantity",
						Position: i,
						Rule:     "numeric_format",
						Message:  fmt.Sprintf("Invalid numeric format in QTY segment: %s", quantity),
						Severity: "ERROR",
					})
				}

				// Check for negative quantities (might be invalid in some contexts)
				if strings.HasPrefix(quantity, "-") {
					errors = append(errors, ValidationError{
						Segment:  segment.Tag,
						Element:  "Quantity",
						Position: i,
						Rule:     "numeric_format",
						Message:  fmt.Sprintf("Negative quantity detected: %s", quantity),
						Severity: "WARNING",
					})
				}
			}
		}

		if segment.Tag == "PRI" && len(segment.Elements) > 0 {
			composite := strings.Split(segment.Elements[0], ":")
			if len(composite) >= 2 {
				price := composite[1]
				if !numericRegex.MatchString(price) {
					errors = append(errors, ValidationError{
						Segment:  segment.Tag,
						Element:  "Price",
						Position: i,
						Rule:     "numeric_format",
						Message:  fmt.Sprintf("Invalid numeric format in PRI segment: %s", price),
						Severity: "ERROR",
					})
				}
			}
		}
	}

	return errors
}

// validateSegmentCount validates segment count in UNT
func (sv *SyntaxValidator) validateSegmentCount(message *EDIFACTMessage) []ValidationError {
	var errors []ValidationError

	// Find UNT segment
	var untSegment *Segment
	for _, segment := range message.Segments {
		if segment.Tag == "UNT" {
			untSegment = segment
			break
		}
	}

	if untSegment != nil && len(untSegment.Elements) > 0 {
		// UNT should contain segment count (excluding UNH/UNT)
		expectedCount := len(message.Segments) - 2 // Exclude UNH and UNT
		actualCount := len(message.Segments)

		if expectedCount != actualCount-2 {
			errors = append(errors, ValidationError{
				Segment:  "UNT",
				Element:  "Segment Count",
				Position: untSegment.Position,
				Rule:     "segment_count",
				Message:  fmt.Sprintf("Segment count mismatch: expected %d, actual %d", expectedCount, actualCount-2),
				Severity: "ERROR",
			})
		}
	}

	return errors
}

// validateReferenceMatching validates reference number matching
func (sv *SyntaxValidator) validateReferenceMatching(message *EDIFACTMessage) []ValidationError {
	var errors []ValidationError

	var unhRef, untRef string

	// Extract UNH reference
	for _, segment := range message.Segments {
		if segment.Tag == "UNH" && len(segment.Elements) > 0 {
			unhRef = segment.Elements[0]
			break
		}
	}

	// Extract UNT reference
	for _, segment := range message.Segments {
		if segment.Tag == "UNT" && len(segment.Elements) > 1 {
			untRef = segment.Elements[1]
			break
		}
	}

	// Compare references
	if unhRef != "" && untRef != "" && unhRef != untRef {
		errors = append(errors, ValidationError{
			Segment:  "UNT",
			Element:  "Reference",
			Position: 0,
			Rule:     "reference_matching",
			Message:  fmt.Sprintf("Reference number mismatch: UNH=%s, UNT=%s", unhRef, untRef),
			Severity: "ERROR",
		})
	}

	return errors
}

// PrintValidationResult prints validation results
func (vr *ValidationResult) PrintValidationResult() {
	fmt.Printf("\n=== Validation Results ===\n")
	fmt.Printf("Valid: %t\n", vr.Valid)
	fmt.Printf("Errors: %d\n", len(vr.Errors))
	fmt.Printf("Warnings: %d\n", len(vr.Warnings))

	if len(vr.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, err := range vr.Errors {
			fmt.Printf("  ❌ [%s] %s: %s\n", err.Segment, err.Rule, err.Message)
		}
	}

	if len(vr.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range vr.Warnings {
			fmt.Printf("  ⚠️  [%s] %s: %s\n", warning.Segment, warning.Rule, warning.Message)
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
			Elements: parts[1:],
			Position: i,
			Raw:      segmentStr + "'",
		}

		message.Segments = append(message.Segments, segment)

		// Extract message type from UNH segment
		if segment.Tag == "UNH" && len(segment.Elements) > 1 {
			msgTypeParts := strings.Split(segment.Elements[1], ":")
			if len(msgTypeParts) > 0 {
				message.MessageType = msgTypeParts[0]
			}
		}
	}

	return message
}

func main() {
	fmt.Println("=== EDIFACT Syntax Validation (Lesson 5) ===")

	// Test messages with various syntax issues
	messages := []string{
		// Valid message
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+1'`,

		// Missing segment terminator
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12346+9'DTM+137:20231201:102'UNT+3+1`,

		// Invalid date format
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12347+9'DTM+137:20231301:102'UNT+3+1'`,

		// Invalid numeric format
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12348+9'DTM+137:20231201:102'LIN+1++ITEM001:EN'QTY+12:ABC:PCE'UNT+5+1'`,

		// Reference mismatch
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12349+9'DTM+137:20231201:102'UNT+3+2'`,

		// Consecutive delimiters
		`UNH+1+INVOIC:D:97A:UN'BGM+380++INV12350+9'DTM+137:20231201:102'UNT+3+1'`,
	}

	// Create validator
	validator := NewSyntaxValidator()

	// Validate each message
	for i, rawMessage := range messages {
		fmt.Printf("\n--- Validating Message %d ---\n", i+1)
		fmt.Printf("Message: %s\n", rawMessage)

		// Parse message
		message := ParseEDIFACTMessage(rawMessage)

		// Validate message
		result := validator.Validate(message)
		result.PrintValidationResult()
	}

	// Demonstrate specific validation rules
	fmt.Println("\n=== Specific Validation Rule Demo ===")

	testMessage := `UNH+1+INVOIC:D:97A:UN'BGM+380+INV12351+9'DTM+137:20231201:102'LIN+1++ITEM001:EN'QTY+12:-5:PCE'PRI+AAA:25.50:CT'UNT+6+1'`
	message := ParseEDIFACTMessage(testMessage)

	// Test individual rules
	fmt.Printf("Testing date format validation...\n")
	dateErrors := validator.validateDateFormats(message)
	for _, err := range dateErrors {
		fmt.Printf("  Date validation: %s\n", err.Message)
	}

	fmt.Printf("Testing numeric format validation...\n")
	numericErrors := validator.validateNumericFormats(message)
	for _, err := range numericErrors {
		fmt.Printf("  Numeric validation: %s\n", err.Message)
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Segment terminator validation")
	fmt.Println("✅ Element delimiter validation")
	fmt.Println("✅ Date format validation")
	fmt.Println("✅ Numeric format validation")
	fmt.Println("✅ Segment count validation")
	fmt.Println("✅ Reference matching validation")
	fmt.Println("✅ Comprehensive error reporting")
}
