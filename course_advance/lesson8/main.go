// Lesson 8: Advanced EDIFACT Processing
// This lesson demonstrates advanced EDIFACT processing techniques including transformation, validation, and integration.

package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// EDIFACTMessage represents a parsed EDIFACT message
type EDIFACTMessage struct {
	MessageType string
	Segments    []*Segment
	RawContent  string
	Metadata    map[string]interface{}
}

// Segment represents an EDIFACT segment
type Segment struct {
	Tag      string
	Elements []string
	Position int
}

// MessageTransformer transforms EDIFACT messages
type MessageTransformer struct {
	transformers map[string]TransformFunc
}

// TransformFunc defines a transformation function
type TransformFunc func(*EDIFACTMessage) (*TransformedMessage, error)

// TransformedMessage represents a transformed message
type TransformedMessage struct {
	OriginalType string
	NewType      string
	Data         interface{}
	Metadata     map[string]interface{}
}

// NewMessageTransformer creates a new message transformer
func NewMessageTransformer() *MessageTransformer {
	mt := &MessageTransformer{
		transformers: make(map[string]TransformFunc),
	}

	// Register default transformers
	mt.RegisterTransformer("INVOIC", mt.transformInvoiceToOrder)
	mt.RegisterTransformer("ORDERS", mt.transformOrderToInvoice)

	return mt
}

// RegisterTransformer registers a transformation function
func (mt *MessageTransformer) RegisterTransformer(msgType string, transformer TransformFunc) {
	mt.transformers[msgType] = transformer
}

// Transform transforms a message using registered transformers
func (mt *MessageTransformer) Transform(message *EDIFACTMessage) (*TransformedMessage, error) {
	if transformer, exists := mt.transformers[message.MessageType]; exists {
		return transformer(message)
	}
	return nil, fmt.Errorf("no transformer found for message type: %s", message.MessageType)
}

// transformInvoiceToOrder transforms INVOIC to ORDERS
func (mt *MessageTransformer) transformInvoiceToOrder(message *EDIFACTMessage) (*TransformedMessage, error) {
	fmt.Println("Transforming INVOIC to ORDERS...")

	// Extract invoice data
	invoiceData := mt.extractInvoiceData(message)

	// Transform to order format
	orderData := &OrderData{
		OrderNumber: "PO-" + invoiceData.InvoiceNumber,
		OrderDate:   invoiceData.InvoiceDate,
		Items:       make([]OrderItem, len(invoiceData.Items)),
	}

	// Transform items
	for i, item := range invoiceData.Items {
		orderData.Items[i] = OrderItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}

	return &TransformedMessage{
		OriginalType: "INVOIC",
		NewType:      "ORDERS",
		Data:         orderData,
		Metadata: map[string]interface{}{
			"transformation_time": time.Now(),
			"original_invoice":    invoiceData.InvoiceNumber,
		},
	}, nil
}

// transformOrderToInvoice transforms ORDERS to INVOIC
func (mt *MessageTransformer) transformOrderToInvoice(message *EDIFACTMessage) (*TransformedMessage, error) {
	fmt.Println("Transforming ORDERS to INVOIC...")

	// Extract order data
	orderData := mt.extractOrderData(message)

	// Transform to invoice format
	invoiceData := &InvoiceData{
		InvoiceNumber: "INV-" + orderData.OrderNumber,
		InvoiceDate:   orderData.OrderDate,
		Items:         make([]InvoiceItem, len(orderData.Items)),
	}

	// Transform items
	for i, item := range orderData.Items {
		invoiceData.Items[i] = InvoiceItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}

	return &TransformedMessage{
		OriginalType: "ORDERS",
		NewType:      "INVOIC",
		Data:         invoiceData,
		Metadata: map[string]interface{}{
			"transformation_time": time.Now(),
			"original_order":      orderData.OrderNumber,
		},
	}, nil
}

// extractInvoiceData extracts invoice data from message
func (mt *MessageTransformer) extractInvoiceData(message *EDIFACTMessage) *InvoiceData {
	invoice := &InvoiceData{
		Items: []InvoiceItem{},
	}

	for _, segment := range message.Segments {
		switch segment.Tag {
		case "BGM":
			if len(segment.Elements) > 1 {
				invoice.InvoiceNumber = segment.Elements[1]
			}
		case "DTM":
			if len(segment.Elements) > 0 {
				parts := strings.Split(segment.Elements[0], ":")
				if len(parts) > 1 {
					invoice.InvoiceDate = parts[1]
				}
			}
		case "LIN":
			if len(segment.Elements) > 2 {
				parts := strings.Split(segment.Elements[2], ":")
				if len(parts) > 0 {
					item := InvoiceItem{ProductCode: parts[0]}
					invoice.Items = append(invoice.Items, item)
				}
			}
		case "QTY":
			if len(segment.Elements) > 0 {
				parts := strings.Split(segment.Elements[0], ":")
				if len(parts) > 1 && len(invoice.Items) > 0 {
					invoice.Items[len(invoice.Items)-1].Quantity = parts[1]
				}
			}
		case "PRI":
			if len(segment.Elements) > 0 {
				parts := strings.Split(segment.Elements[0], ":")
				if len(parts) > 1 && len(invoice.Items) > 0 {
					invoice.Items[len(invoice.Items)-1].UnitPrice = parts[1]
				}
			}
		}
	}

	return invoice
}

// extractOrderData extracts order data from message
func (mt *MessageTransformer) extractOrderData(message *EDIFACTMessage) *OrderData {
	order := &OrderData{
		Items: []OrderItem{},
	}

	for _, segment := range message.Segments {
		switch segment.Tag {
		case "BGM":
			if len(segment.Elements) > 1 {
				order.OrderNumber = segment.Elements[1]
			}
		case "DTM":
			if len(segment.Elements) > 0 {
				parts := strings.Split(segment.Elements[0], ":")
				if len(parts) > 1 {
					order.OrderDate = parts[1]
				}
			}
		case "LIN":
			if len(segment.Elements) > 2 {
				parts := strings.Split(segment.Elements[2], ":")
				if len(parts) > 0 {
					item := OrderItem{ProductCode: parts[0]}
					order.Items = append(order.Items, item)
				}
			}
		case "QTY":
			if len(segment.Elements) > 0 {
				parts := strings.Split(segment.Elements[0], ":")
				if len(parts) > 1 && len(order.Items) > 0 {
					order.Items[len(order.Items)-1].Quantity = parts[1]
				}
			}
		case "PRI":
			if len(segment.Elements) > 0 {
				parts := strings.Split(segment.Elements[0], ":")
				if len(parts) > 1 && len(order.Items) > 0 {
					order.Items[len(order.Items)-1].UnitPrice = parts[1]
				}
			}
		}
	}

	return order
}

// AdvancedValidator provides comprehensive validation
type AdvancedValidator struct {
	rules map[string][]ValidationRule
}

// ValidationRule defines a validation rule
type ValidationRule struct {
	Name        string
	Description string
	Validate    func(*EDIFACTMessage) error
}

// ValidationResult contains validation results
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationWarning
}

// ValidationError represents a validation error
type ValidationError struct {
	Segment  string
	Element  string
	Rule     string
	Message  string
	Severity string
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Segment string
	Element string
	Rule    string
	Message string
}

// NewAdvancedValidator creates a new advanced validator
func NewAdvancedValidator() *AdvancedValidator {
	av := &AdvancedValidator{
		rules: make(map[string][]ValidationRule),
	}

	// Register default validation rules
	av.registerDefaultRules()

	return av
}

// registerDefaultRules registers default validation rules
func (av *AdvancedValidator) registerDefaultRules() {
	// INVOIC validation rules
	av.rules["INVOIC"] = []ValidationRule{
		{
			Name:        "mandatory_segments",
			Description: "Check for mandatory segments",
			Validate:    av.validateMandatorySegments,
		},
		{
			Name:        "segment_order",
			Description: "Validate segment order",
			Validate:    av.validateSegmentOrder,
		},
		{
			Name:        "data_format",
			Description: "Validate data formats",
			Validate:    av.validateDataFormats,
		},
	}

	// ORDERS validation rules
	av.rules["ORDERS"] = []ValidationRule{
		{
			Name:        "mandatory_segments",
			Description: "Check for mandatory segments",
			Validate:    av.validateMandatorySegments,
		},
		{
			Name:        "segment_order",
			Description: "Validate segment order",
			Validate:    av.validateSegmentOrder,
		},
		{
			Name:        "data_format",
			Description: "Validate data formats",
			Validate:    av.validateDataFormats,
		},
	}
}

// Validate validates a message using registered rules
func (av *AdvancedValidator) Validate(message *EDIFACTMessage) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	if rules, exists := av.rules[message.MessageType]; exists {
		for _, rule := range rules {
			if err := rule.Validate(message); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Rule:     rule.Name,
					Message:  err.Error(),
					Severity: "ERROR",
				})
			}
		}
	}

	return result
}

// validateMandatorySegments validates mandatory segments
func (av *AdvancedValidator) validateMandatorySegments(message *EDIFACTMessage) error {
	requiredSegments := map[string]bool{
		"UNH": true,
		"BGM": true,
		"UNT": true,
	}

	foundSegments := make(map[string]bool)
	for _, segment := range message.Segments {
		foundSegments[segment.Tag] = true
	}

	for required := range requiredSegments {
		if !foundSegments[required] {
			return fmt.Errorf("missing mandatory segment: %s", required)
		}
	}

	return nil
}

// validateSegmentOrder validates segment order
func (av *AdvancedValidator) validateSegmentOrder(message *EDIFACTMessage) error {
	if len(message.Segments) < 2 {
		return fmt.Errorf("message too short")
	}

	// Check that UNH is first
	if message.Segments[0].Tag != "UNH" {
		return fmt.Errorf("UNH must be the first segment")
	}

	// Check that UNT is last
	if message.Segments[len(message.Segments)-1].Tag != "UNT" {
		return fmt.Errorf("UNT must be the last segment")
	}

	return nil
}

// validateDataFormats validates data formats
func (av *AdvancedValidator) validateDataFormats(message *EDIFACTMessage) error {
	dateRegex := regexp.MustCompile(`^\d{8}$`)

	for _, segment := range message.Segments {
		if segment.Tag == "DTM" && len(segment.Elements) > 0 {
			parts := strings.Split(segment.Elements[0], ":")
			if len(parts) > 1 {
				if !dateRegex.MatchString(parts[1]) {
					return fmt.Errorf("invalid date format in DTM segment: %s", parts[1])
				}
			}
		}
	}

	return nil
}

// MessageIntegrator provides integration capabilities
type MessageIntegrator struct {
	transformers map[string]TransformFunc
	validators   map[string][]ValidationRule
}

// NewMessageIntegrator creates a new message integrator
func NewMessageIntegrator() *MessageIntegrator {
	return &MessageIntegrator{
		transformers: make(map[string]TransformFunc),
		validators:   make(map[string][]ValidationRule),
	}
}

// ProcessMessage processes a message with transformation and validation
func (mi *MessageIntegrator) ProcessMessage(message *EDIFACTMessage) (*ProcessedMessage, error) {
	fmt.Printf("Processing %s message...\n", message.MessageType)

	// Validate message
	validator := NewAdvancedValidator()
	validationResult := validator.Validate(message)

	if !validationResult.Valid {
		return nil, fmt.Errorf("validation failed: %v", validationResult.Errors)
	}

	// Transform message if needed
	transformer := NewMessageTransformer()
	transformed, err := transformer.Transform(message)
	if err != nil {
		return nil, fmt.Errorf("transformation failed: %v", err)
	}

	return &ProcessedMessage{
		Original:    message,
		Transformed: transformed,
		Validation:  validationResult,
		ProcessedAt: time.Now(),
	}, nil
}

// ProcessedMessage represents a processed message
type ProcessedMessage struct {
	Original    *EDIFACTMessage
	Transformed *TransformedMessage
	Validation  *ValidationResult
	ProcessedAt time.Time
}

// Data structures for transformation
type InvoiceData struct {
	InvoiceNumber string
	InvoiceDate   string
	Items         []InvoiceItem
}

type InvoiceItem struct {
	ProductCode string
	Quantity    string
	UnitPrice   string
}

type OrderData struct {
	OrderNumber string
	OrderDate   string
	Items       []OrderItem
}

type OrderItem struct {
	ProductCode string
	Quantity    string
	UnitPrice   string
}

// ParseEDIFACTMessage parses a raw EDIFACT message
func ParseEDIFACTMessage(rawContent string) *EDIFACTMessage {
	message := &EDIFACTMessage{
		RawContent: rawContent,
		Segments:   []*Segment{},
		Metadata:   make(map[string]interface{}),
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

func main() {
	fmt.Println("=== Advanced EDIFACT Processing (Lesson 8) ===")

	// Test messages
	invoiceMessage := `UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+1'`
	orderMessage := `UNH+1+ORDERS:D:97A:UN'BGM+220+PO67890+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM002:EN'QTY+12:50:PCE'UNT+6+1'`

	// Create integrator
	integrator := NewMessageIntegrator()

	// Process messages
	messages := []string{invoiceMessage, orderMessage}

	for i, rawMessage := range messages {
		fmt.Printf("\n--- Processing Message %d ---\n", i+1)

		// Parse message
		message := ParseEDIFACTMessage(rawMessage)
		fmt.Printf("Message type: %s\n", message.MessageType)

		// Process with integrator
		processed, err := integrator.ProcessMessage(message)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Display results
		fmt.Printf("Validation: %t\n", processed.Validation.Valid)
		if processed.Transformed != nil {
			fmt.Printf("Transformed to: %s\n", processed.Transformed.NewType)

			// Convert to JSON for display
			if jsonData, err := json.MarshalIndent(processed.Transformed.Data, "", "  "); err == nil {
				fmt.Printf("Transformed data:\n%s\n", string(jsonData))
			}
		}
	}

	// Demonstrate advanced validation
	fmt.Println("\n=== Advanced Validation Demo ===")

	invalidMessage := `UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'DTM+137:20231301:102'UNT+3+1'` // Invalid date
	message := ParseEDIFACTMessage(invalidMessage)

	validator := NewAdvancedValidator()
	result := validator.Validate(message)

	fmt.Printf("Validation result: %t\n", result.Valid)
	for _, err := range result.Errors {
		fmt.Printf("Error: %s - %s\n", err.Rule, err.Message)
	}

	// Demonstrate transformation
	fmt.Println("\n=== Transformation Demo ===")

	transformer := NewMessageTransformer()
	transformed, err := transformer.Transform(message)
	if err != nil {
		fmt.Printf("Transformation error: %v\n", err)
	} else {
		fmt.Printf("Transformed from %s to %s\n", transformed.OriginalType, transformed.NewType)
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Advanced message transformation")
	fmt.Println("✅ Comprehensive validation rules")
	fmt.Println("✅ Message integration patterns")
	fmt.Println("✅ Data extraction and mapping")
	fmt.Println("✅ Error handling and reporting")
	fmt.Println("✅ Metadata management")
}
