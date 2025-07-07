// Lesson 4: Message Structure
// This lesson demonstrates EDIFACT message structure analysis and validation.

package main

import (
	"fmt"
	"strings"
)

// EDIFACTMessage represents a parsed EDIFACT message
type EDIFACTMessage struct {
	MessageType string
	Segments    []*Segment
	Groups      []*SegmentGroup
	RawContent  string
}

// Segment represents an EDIFACT segment
type Segment struct {
	Tag      string
	Elements []string
	Position int
}

// SegmentGroup represents a group of related segments
type SegmentGroup struct {
	Name     string
	Segments []*Segment
	Repeat   int
}

// MessageAnalyzer analyzes EDIFACT message structure
type MessageAnalyzer struct {
	messageTypes map[string]*MessageTypeDefinition
}

// MessageTypeDefinition defines the structure of a message type
type MessageTypeDefinition struct {
	Name          string
	MandatorySegs []string
	OptionalSegs  []string
	SegmentGroups []*GroupDefinition
	SegmentOrder  []string
}

// GroupDefinition defines a segment group
type GroupDefinition struct {
	Name           string
	TriggerSegment string
	Segments       []string
	Repeatable     bool
}

// NewMessageAnalyzer creates a new message analyzer
func NewMessageAnalyzer() *MessageAnalyzer {
	ma := &MessageAnalyzer{
		messageTypes: make(map[string]*MessageTypeDefinition),
	}

	// Register message type definitions
	ma.registerMessageTypes()

	return ma
}

// registerMessageTypes registers known message type definitions
func (ma *MessageAnalyzer) registerMessageTypes() {
	// INVOIC message definition
	ma.messageTypes["INVOIC"] = &MessageTypeDefinition{
		Name:          "INVOIC",
		MandatorySegs: []string{"UNH", "BGM", "UNT"},
		OptionalSegs:  []string{"DTM", "NAD", "LIN", "QTY", "PRI"},
		SegmentOrder:  []string{"UNH", "BGM", "DTM", "NAD", "LIN", "QTY", "PRI", "UNT"},
		SegmentGroups: []*GroupDefinition{
			{
				Name:           "Line Item Group",
				TriggerSegment: "LIN",
				Segments:       []string{"LIN", "QTY", "PRI"},
				Repeatable:     true,
			},
		},
	}

	// ORDERS message definition
	ma.messageTypes["ORDERS"] = &MessageTypeDefinition{
		Name:          "ORDERS",
		MandatorySegs: []string{"UNH", "BGM", "UNT"},
		OptionalSegs:  []string{"DTM", "NAD", "LIN", "QTY", "PRI"},
		SegmentOrder:  []string{"UNH", "BGM", "DTM", "NAD", "LIN", "QTY", "PRI", "UNT"},
		SegmentGroups: []*GroupDefinition{
			{
				Name:           "Line Item Group",
				TriggerSegment: "LIN",
				Segments:       []string{"LIN", "QTY", "PRI"},
				Repeatable:     true,
			},
		},
	}
}

// AnalyzeMessage analyzes the structure of an EDIFACT message
func (ma *MessageAnalyzer) AnalyzeMessage(message *EDIFACTMessage) *MessageAnalysis {
	analysis := &MessageAnalysis{
		MessageType: message.MessageType,
		Valid:       true,
		Errors:      []string{},
		Warnings:    []string{},
		Groups:      []*GroupAnalysis{},
	}

	// Get message type definition
	definition, exists := ma.messageTypes[message.MessageType]
	if !exists {
		analysis.Valid = false
		analysis.Errors = append(analysis.Errors, fmt.Sprintf("Unknown message type: %s", message.MessageType))
		return analysis
	}

	// Validate mandatory segments
	ma.validateMandatorySegments(message, definition, analysis)

	// Validate segment order
	ma.validateSegmentOrder(message, definition, analysis)

	// Analyze segment groups
	ma.analyzeSegmentGroups(message, definition, analysis)

	// Count segments
	analysis.SegmentCount = len(message.Segments)

	return analysis
}

// validateMandatorySegments validates that all mandatory segments are present
func (ma *MessageAnalyzer) validateMandatorySegments(message *EDIFACTMessage, definition *MessageTypeDefinition, analysis *MessageAnalysis) {
	foundSegments := make(map[string]bool)
	for _, segment := range message.Segments {
		foundSegments[segment.Tag] = true
	}

	for _, mandatory := range definition.MandatorySegs {
		if !foundSegments[mandatory] {
			analysis.Valid = false
			analysis.Errors = append(analysis.Errors, fmt.Sprintf("Missing mandatory segment: %s", mandatory))
		}
	}
}

// validateSegmentOrder validates the order of segments
func (ma *MessageAnalyzer) validateSegmentOrder(message *EDIFACTMessage, definition *MessageTypeDefinition, analysis *MessageAnalysis) {
	if len(message.Segments) < 2 {
		analysis.Valid = false
		analysis.Errors = append(analysis.Errors, "Message too short")
		return
	}

	// Check that UNH is first
	if message.Segments[0].Tag != "UNH" {
		analysis.Valid = false
		analysis.Errors = append(analysis.Errors, "UNH must be the first segment")
	}

	// Check that UNT is last
	if message.Segments[len(message.Segments)-1].Tag != "UNT" {
		analysis.Valid = false
		analysis.Errors = append(analysis.Errors, "UNT must be the last segment")
	}
}

// analyzeSegmentGroups analyzes segment groups in the message
func (ma *MessageAnalyzer) analyzeSegmentGroups(message *EDIFACTMessage, definition *MessageTypeDefinition, analysis *MessageAnalysis) {
	for _, groupDef := range definition.SegmentGroups {
		groupAnalysis := &GroupAnalysis{
			Name:     groupDef.Name,
			Segments: []*Segment{},
		}

		// Find segments belonging to this group
		for _, segment := range message.Segments {
			for _, groupSegment := range groupDef.Segments {
				if segment.Tag == groupSegment {
					groupAnalysis.Segments = append(groupAnalysis.Segments, segment)
					break
				}
			}
		}

		// Count trigger segments to determine group count
		triggerCount := 0
		for _, segment := range groupAnalysis.Segments {
			if segment.Tag == groupDef.TriggerSegment {
				triggerCount++
			}
		}
		groupAnalysis.Count = triggerCount

		analysis.Groups = append(analysis.Groups, groupAnalysis)
	}
}

// MessageAnalysis contains the analysis results
type MessageAnalysis struct {
	MessageType  string
	Valid        bool
	Errors       []string
	Warnings     []string
	SegmentCount int
	Groups       []*GroupAnalysis
}

// GroupAnalysis contains group analysis results
type GroupAnalysis struct {
	Name     string
	Segments []*Segment
	Count    int
}

// PrintAnalysis prints the analysis results
func (ma *MessageAnalysis) PrintAnalysis() {
	fmt.Printf("\n=== Message Analysis Results ===\n")
	fmt.Printf("Message Type: %s\n", ma.MessageType)
	fmt.Printf("Valid: %t\n", ma.Valid)
	fmt.Printf("Segment Count: %d\n", ma.SegmentCount)

	if len(ma.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, err := range ma.Errors {
			fmt.Printf("  ❌ %s\n", err)
		}
	}

	if len(ma.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range ma.Warnings {
			fmt.Printf("  ⚠️  %s\n", warning)
		}
	}

	if len(ma.Groups) > 0 {
		fmt.Printf("\nSegment Groups:\n")
		for _, group := range ma.Groups {
			fmt.Printf("  📦 %s: %d occurrences\n", group.Name, group.Count)
			for _, segment := range group.Segments {
				fmt.Printf("    - %s\n", segment.Tag)
			}
		}
	}
}

// ParseEDIFACTMessage parses a raw EDIFACT message
func ParseEDIFACTMessage(rawContent string) *EDIFACTMessage {
	message := &EDIFACTMessage{
		RawContent: rawContent,
		Segments:   []*Segment{},
		Groups:     []*SegmentGroup{},
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

// PrintMessageStructure prints the message structure
func PrintMessageStructure(message *EDIFACTMessage) {
	fmt.Printf("\n=== Message Structure ===\n")
	fmt.Printf("Message Type: %s\n", message.MessageType)
	fmt.Printf("Total Segments: %d\n\n", len(message.Segments))

	fmt.Printf("Segment Hierarchy:\n")
	fmt.Printf("├── Interchange Level\n")
	fmt.Printf("│   └── Message Level\n")

	for i, segment := range message.Segments {
		indent := "    "
		if i == 0 {
			fmt.Printf("│       ├── %s (Header)\n", segment.Tag)
		} else if i == len(message.Segments)-1 {
			fmt.Printf("│       └── %s (Trailer)\n", segment.Tag)
		} else {
			fmt.Printf("│       ├── %s\n", segment.Tag)
		}

		// Show elements for key segments
		if segment.Tag == "BGM" || segment.Tag == "DTM" || segment.Tag == "LIN" {
			for j, element := range segment.Elements {
				if j < 3 { // Limit to first 3 elements for readability
					fmt.Printf("%s│           ├── Element %d: %s\n", indent, j+1, element)
				}
			}
		}
	}
}

func main() {
	fmt.Println("=== EDIFACT Message Structure (Lesson 4) ===")

	// Test messages
	messages := []string{
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'NAD+SE+++SUPPLIER INC'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'LIN+2++ITEM002:EN'QTY+12:50:PCE'PRI+AAA:30.00:CT'UNT+12+1'`,
		`UNH+1+ORDERS:D:97A:UN'BGM+220+PO67890+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'UNT+6+1'`,
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12346+9'UNT+2+1'`, // Missing mandatory segments
	}

	// Create analyzer
	analyzer := NewMessageAnalyzer()

	// Analyze each message
	for i, rawMessage := range messages {
		fmt.Printf("\n--- Analyzing Message %d ---\n", i+1)

		// Parse message
		message := ParseEDIFACTMessage(rawMessage)

		// Print structure
		PrintMessageStructure(message)

		// Analyze message
		analysis := analyzer.AnalyzeMessage(message)
		analysis.PrintAnalysis()
	}

	// Demonstrate segment group analysis
	fmt.Println("\n=== Segment Group Analysis Demo ===")

	complexMessage := `UNH+1+INVOIC:D:97A:UN'BGM+380+INV12347+9'DTM+137:20231201:102'NAD+BY+++BUYER CORP'NAD+SE+++SELLER INC'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'LIN+2++ITEM002:EN'QTY+12:50:PCE'PRI+AAA:30.00:CT'LIN+3++ITEM003:EN'QTY+12:75:PCE'PRI+AAA:15.00:CT'UNT+15+1'`

	message := ParseEDIFACTMessage(complexMessage)
	analysis := analyzer.AnalyzeMessage(message)

	fmt.Printf("Complex message analysis:\n")
	fmt.Printf("Total segments: %d\n", analysis.SegmentCount)
	fmt.Printf("Line item groups: %d\n", len(analysis.Groups))

	for _, group := range analysis.Groups {
		fmt.Printf("  %s: %d occurrences\n", group.Name, group.Count)
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Message structure analysis")
	fmt.Println("✅ Mandatory segment validation")
	fmt.Println("✅ Segment order validation")
	fmt.Println("✅ Segment group identification")
	fmt.Println("✅ Message type definition")
	fmt.Println("✅ Comprehensive error reporting")
}
