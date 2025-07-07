// Lesson 6: Service Segments
// This lesson demonstrates service segment handling and validation.

package main

import (
	"fmt"
	"strings"
	"time"
)

// EDIFACTMessage represents a parsed EDIFACT message
type EDIFACTMessage struct {
	MessageType string
	Segments    []*Segment
	RawContent  string
	Interchange *InterchangeInfo
}

// Segment represents an EDIFACT segment
type Segment struct {
	Tag      string
	Elements []string
	Position int
	Raw      string
}

// InterchangeInfo contains interchange-level information
type InterchangeInfo struct {
	SenderID    string
	ReceiverID  string
	Date        string
	Time        string
	RefNumber   string
	MessageType string
}

// ServiceSegmentHandler handles service segments
type ServiceSegmentHandler struct {
	validators map[string]func(*Segment) error
}

// NewServiceSegmentHandler creates a new service segment handler
func NewServiceSegmentHandler() *ServiceSegmentHandler {
	ssh := &ServiceSegmentHandler{
		validators: make(map[string]func(*Segment) error),
	}

	// Register validators for different service segments
	ssh.registerValidators()

	return ssh
}

// registerValidators registers validation functions for service segments
func (ssh *ServiceSegmentHandler) registerValidators() {
	ssh.validators["UNH"] = ssh.validateUNH
	ssh.validators["UNT"] = ssh.validateUNT
	ssh.validators["UNB"] = ssh.validateUNB
	ssh.validators["UNZ"] = ssh.validateUNZ
	ssh.validators["UNG"] = ssh.validateUNG
	ssh.validators["UNE"] = ssh.validateUNE
}

// ProcessMessage processes a message and extracts service segment information
func (ssh *ServiceSegmentHandler) ProcessMessage(message *EDIFACTMessage) *ServiceSegmentResult {
	result := &ServiceSegmentResult{
		MessageType: message.MessageType,
		Valid:       true,
		Errors:      []string{},
		Warnings:    []string{},
		ServiceSegs: make(map[string]*ServiceSegmentInfo),
	}

	// Process each segment
	for _, segment := range message.Segments {
		if validator, exists := ssh.validators[segment.Tag]; exists {
			// Validate service segment
			if err := validator(segment); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", segment.Tag, err))
			}

			// Extract service segment information
			info := ssh.extractServiceSegmentInfo(segment)
			result.ServiceSegs[segment.Tag] = info
		}
	}

	// Validate message envelope (UNH/UNT pair)
	ssh.validateMessageEnvelope(message, result)

	// Extract interchange information
	message.Interchange = ssh.extractInterchangeInfo(message)

	return result
}

// validateUNH validates UNH (Message Header) segment
func (ssh *ServiceSegmentHandler) validateUNH(segment *Segment) error {
	if len(segment.Elements) < 2 {
		return fmt.Errorf("UNH must have at least 2 elements")
	}

	// Validate message type format
	msgTypeParts := strings.Split(segment.Elements[1], ":")
	if len(msgTypeParts) < 4 {
		return fmt.Errorf("invalid message type format in UNH")
	}

	// Validate version and release
	if len(msgTypeParts) >= 4 {
		release := msgTypeParts[2]
		version := msgTypeParts[3]
		if release == "" || version == "" {
			return fmt.Errorf("invalid release/version in UNH")
		}
	}

	return nil
}

// validateUNT validates UNT (Message Trailer) segment
func (ssh *ServiceSegmentHandler) validateUNT(segment *Segment) error {
	if len(segment.Elements) < 2 {
		return fmt.Errorf("UNT must have at least 2 elements")
	}

	// Validate segment count
	segmentCount := segment.Elements[0]
	if segmentCount == "" {
		return fmt.Errorf("segment count missing in UNT")
	}

	// Validate reference number
	refNumber := segment.Elements[1]
	if refNumber == "" {
		return fmt.Errorf("reference number missing in UNT")
	}

	return nil
}

// validateUNB validates UNB (Interchange Header) segment
func (ssh *ServiceSegmentHandler) validateUNB(segment *Segment) error {
	if len(segment.Elements) < 5 {
		return fmt.Errorf("UNB must have at least 5 elements")
	}

	// Validate syntax identifier
	syntaxID := segment.Elements[0]
	if syntaxID == "" {
		return fmt.Errorf("syntax identifier missing in UNB")
	}

	// Validate sender and receiver IDs
	senderID := segment.Elements[2]
	receiverID := segment.Elements[3]
	if senderID == "" || receiverID == "" {
		return fmt.Errorf("sender or receiver ID missing in UNB")
	}

	// Validate date and time
	dateTime := segment.Elements[4]
	if dateTime == "" {
		return fmt.Errorf("date/time missing in UNB")
	}

	// Validate interchange control reference
	if len(segment.Elements) >= 6 {
		refNumber := segment.Elements[5]
		if refNumber == "" {
			return fmt.Errorf("interchange control reference missing in UNB")
		}
	}

	return nil
}

// validateUNZ validates UNZ (Interchange Trailer) segment
func (ssh *ServiceSegmentHandler) validateUNZ(segment *Segment) error {
	if len(segment.Elements) < 2 {
		return fmt.Errorf("UNZ must have at least 2 elements")
	}

	// Validate message count
	messageCount := segment.Elements[0]
	if messageCount == "" {
		return fmt.Errorf("message count missing in UNZ")
	}

	// Validate interchange control reference
	refNumber := segment.Elements[1]
	if refNumber == "" {
		return fmt.Errorf("interchange control reference missing in UNZ")
	}

	return nil
}

// validateUNG validates UNG (Functional Group Header) segment
func (ssh *ServiceSegmentHandler) validateUNG(segment *Segment) error {
	if len(segment.Elements) < 4 {
		return fmt.Errorf("UNG must have at least 4 elements")
	}

	// Validate functional group ID
	groupID := segment.Elements[0]
	if groupID == "" {
		return fmt.Errorf("functional group ID missing in UNG")
	}

	// Validate sender and receiver IDs
	senderID := segment.Elements[2]
	receiverID := segment.Elements[3]
	if senderID == "" || receiverID == "" {
		return fmt.Errorf("sender or receiver ID missing in UNG")
	}

	return nil
}

// validateUNE validates UNE (Functional Group Trailer) segment
func (ssh *ServiceSegmentHandler) validateUNE(segment *Segment) error {
	if len(segment.Elements) < 2 {
		return fmt.Errorf("UNE must have at least 2 elements")
	}

	// Validate message count
	messageCount := segment.Elements[0]
	if messageCount == "" {
		return fmt.Errorf("message count missing in UNE")
	}

	// Validate functional group reference
	refNumber := segment.Elements[1]
	if refNumber == "" {
		return fmt.Errorf("functional group reference missing in UNE")
	}

	return nil
}

// extractServiceSegmentInfo extracts information from a service segment
func (ssh *ServiceSegmentHandler) extractServiceSegmentInfo(segment *Segment) *ServiceSegmentInfo {
	info := &ServiceSegmentInfo{
		Tag:       segment.Tag,
		Elements:  segment.Elements,
		Position:  segment.Position,
		Timestamp: time.Now(),
	}

	switch segment.Tag {
	case "UNH":
		if len(segment.Elements) >= 2 {
			msgTypeParts := strings.Split(segment.Elements[1], ":")
			if len(msgTypeParts) >= 1 {
				info.MessageType = msgTypeParts[0]
			}
			if len(msgTypeParts) >= 4 {
				info.Release = msgTypeParts[2]
				info.Version = msgTypeParts[3]
			}
		}
		if len(segment.Elements) >= 1 {
			info.ReferenceNumber = segment.Elements[0]
		}

	case "UNT":
		if len(segment.Elements) >= 1 {
			info.SegmentCount = segment.Elements[0]
		}
		if len(segment.Elements) >= 2 {
			info.ReferenceNumber = segment.Elements[1]
		}

	case "UNB":
		if len(segment.Elements) >= 1 {
			info.SyntaxIdentifier = segment.Elements[0]
		}
		if len(segment.Elements) >= 3 {
			info.SenderID = segment.Elements[2]
		}
		if len(segment.Elements) >= 4 {
			info.ReceiverID = segment.Elements[3]
		}
		if len(segment.Elements) >= 5 {
			info.DateTime = segment.Elements[4]
		}
		if len(segment.Elements) >= 6 {
			info.ControlReference = segment.Elements[5]
		}

	case "UNZ":
		if len(segment.Elements) >= 1 {
			info.MessageCount = segment.Elements[0]
		}
		if len(segment.Elements) >= 2 {
			info.ControlReference = segment.Elements[1]
		}
	}

	return info
}

// validateMessageEnvelope validates the message envelope (UNH/UNT pair)
func (ssh *ServiceSegmentHandler) validateMessageEnvelope(message *EDIFACTMessage, result *ServiceSegmentResult) {
	var unhRef, untRef string

	// Find UNH and UNT segments
	for _, segment := range message.Segments {
		if segment.Tag == "UNH" && len(segment.Elements) > 0 {
			unhRef = segment.Elements[0]
		}
		if segment.Tag == "UNT" && len(segment.Elements) > 1 {
			untRef = segment.Elements[1]
		}
	}

	// Check reference number matching
	if unhRef != "" && untRef != "" && unhRef != untRef {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Reference number mismatch: UNH=%s, UNT=%s", unhRef, untRef))
	}

	// Check segment count
	if untInfo, exists := result.ServiceSegs["UNT"]; exists {
		expectedCount := len(message.Segments) - 2 // Exclude UNH and UNT
		if untInfo.SegmentCount != fmt.Sprintf("%d", expectedCount) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Segment count mismatch: expected %d, found %s", expectedCount, untInfo.SegmentCount))
		}
	}
}

// extractInterchangeInfo extracts interchange-level information
func (ssh *ServiceSegmentHandler) extractInterchangeInfo(message *EDIFACTMessage) *InterchangeInfo {
	info := &InterchangeInfo{}

	for _, segment := range message.Segments {
		if segment.Tag == "UNB" {
			if len(segment.Elements) >= 3 {
				info.SenderID = segment.Elements[2]
			}
			if len(segment.Elements) >= 4 {
				info.ReceiverID = segment.Elements[3]
			}
			if len(segment.Elements) >= 5 {
				dateTime := segment.Elements[4]
				parts := strings.Split(dateTime, ":")
				if len(parts) >= 2 {
					info.Date = parts[0]
					info.Time = parts[1]
				}
			}
			if len(segment.Elements) >= 6 {
				info.RefNumber = segment.Elements[5]
			}
		}
		if segment.Tag == "UNH" && len(segment.Elements) >= 2 {
			msgTypeParts := strings.Split(segment.Elements[1], ":")
			if len(msgTypeParts) >= 1 {
				info.MessageType = msgTypeParts[0]
			}
		}
	}

	return info
}

// ServiceSegmentResult contains service segment processing results
type ServiceSegmentResult struct {
	MessageType string
	Valid       bool
	Errors      []string
	Warnings    []string
	ServiceSegs map[string]*ServiceSegmentInfo
}

// ServiceSegmentInfo contains information about a service segment
type ServiceSegmentInfo struct {
	Tag              string
	Elements         []string
	Position         int
	Timestamp        time.Time
	MessageType      string
	ReferenceNumber  string
	SegmentCount     string
	SyntaxIdentifier string
	SenderID         string
	ReceiverID       string
	DateTime         string
	ControlReference string
	MessageCount     string
	Release          string
	Version          string
}

// PrintServiceSegmentResult prints service segment processing results
func (ssr *ServiceSegmentResult) PrintServiceSegmentResult() {
	fmt.Printf("\n=== Service Segment Analysis ===\n")
	fmt.Printf("Message Type: %s\n", ssr.MessageType)
	fmt.Printf("Valid: %t\n", ssr.Valid)
	fmt.Printf("Service Segments Found: %d\n", len(ssr.ServiceSegs))

	if len(ssr.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, err := range ssr.Errors {
			fmt.Printf("  ❌ %s\n", err)
		}
	}

	if len(ssr.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range ssr.Warnings {
			fmt.Printf("  ⚠️  %s\n", warning)
		}
	}

	if len(ssr.ServiceSegs) > 0 {
		fmt.Printf("\nService Segments:\n")
		for tag, info := range ssr.ServiceSegs {
			fmt.Printf("  📋 %s:\n", tag)
			fmt.Printf("    Position: %d\n", info.Position)
			if info.MessageType != "" {
				fmt.Printf("    Message Type: %s\n", info.MessageType)
			}
			if info.ReferenceNumber != "" {
				fmt.Printf("    Reference: %s\n", info.ReferenceNumber)
			}
			if info.SegmentCount != "" {
				fmt.Printf("    Segment Count: %s\n", info.SegmentCount)
			}
			if info.SenderID != "" {
				fmt.Printf("    Sender: %s\n", info.SenderID)
			}
			if info.ReceiverID != "" {
				fmt.Printf("    Receiver: %s\n", info.ReceiverID)
			}
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
	fmt.Println("=== EDIFACT Service Segments (Lesson 6) ===")

	// Test messages with different service segment configurations
	messages := []string{
		// Simple message with UNH/UNT
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+1'`,

		// Message with interchange envelope
		`UNB+UNOA:2+SENDER+RECEIVER+231201:1430+12345+++INVOIC'UNH+1+INVOIC:D:97A:UN'BGM+380+INV12346+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+1'UNZ+1+12345'`,

		// Message with functional group
		`UNB+UNOA:2+SENDER+RECEIVER+231201:1430+12345+++INVOIC'UNG+INVOIC+SENDER+RECEIVER+231201:1430+12345+UN+D:97A'UNH+1+INVOIC:D:97A:UN'BGM+380+INV12347+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+1'UNE+1+12345'UNZ+1+12345'`,

		// Invalid message (reference mismatch)
		`UNH+1+INVOIC:D:97A:UN'BGM+380+INV12348+9'DTM+137:20231201:102'UNT+3+2'`,
	}

	// Create service segment handler
	handler := NewServiceSegmentHandler()

	// Process each message
	for i, rawMessage := range messages {
		fmt.Printf("\n--- Processing Message %d ---\n", i+1)
		fmt.Printf("Message: %s\n", rawMessage)

		// Parse message
		message := ParseEDIFACTMessage(rawMessage)

		// Process service segments
		result := handler.ProcessMessage(message)
		result.PrintServiceSegmentResult()

		// Print interchange information if available
		if message.Interchange != nil {
			fmt.Printf("\nInterchange Information:\n")
			fmt.Printf("  Sender: %s\n", message.Interchange.SenderID)
			fmt.Printf("  Receiver: %s\n", message.Interchange.ReceiverID)
			fmt.Printf("  Date: %s\n", message.Interchange.Date)
			fmt.Printf("  Time: %s\n", message.Interchange.Time)
			fmt.Printf("  Reference: %s\n", message.Interchange.RefNumber)
			fmt.Printf("  Message Type: %s\n", message.Interchange.MessageType)
		}
	}

	// Demonstrate service segment validation
	fmt.Println("\n=== Service Segment Validation Demo ===")

	testMessage := `UNH+1+INVOIC:D:97A:UN'BGM+380+INV12349+9'DTM+137:20231201:102'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+1'`
	message := ParseEDIFACTMessage(testMessage)

	// Test individual validations
	fmt.Printf("Testing UNH validation...\n")
	if err := handler.validators["UNH"](message.Segments[0]); err != nil {
		fmt.Printf("  UNH validation error: %v\n", err)
	} else {
		fmt.Printf("  UNH validation passed\n")
	}

	fmt.Printf("Testing UNT validation...\n")
	if err := handler.validators["UNT"](message.Segments[len(message.Segments)-1]); err != nil {
		fmt.Printf("  UNT validation error: %v\n", err)
	} else {
		fmt.Printf("  UNT validation passed\n")
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Service segment parsing and validation")
	fmt.Println("✅ Message envelope validation (UNH/UNT)")
	fmt.Println("✅ Interchange envelope handling (UNB/UNZ)")
	fmt.Println("✅ Functional group handling (UNG/UNE)")
	fmt.Println("✅ Reference number matching")
	fmt.Println("✅ Segment count validation")
	fmt.Println("✅ Comprehensive error reporting")
}
