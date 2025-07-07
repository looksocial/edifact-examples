// Lesson 2: Custom Message Handlers
// This lesson demonstrates creating custom message handlers for specific business requirements.

package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// MessageHandler interface for processing EDIFACT messages
type MessageHandler interface {
	Process(message *EDIFACTMessage) error
	GetMessageType() string
}

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
}

// MessageDetector routes messages to appropriate handlers
type MessageDetector struct {
	handlers map[string]MessageHandler
	fallback MessageHandler
	mu       sync.RWMutex
}

// NewMessageDetector creates a new message detector
func NewMessageDetector() *MessageDetector {
	return &MessageDetector{
		handlers: make(map[string]MessageHandler),
		fallback: &GenericHandler{},
	}
}

// RegisterHandler registers a handler for a specific message type
func (md *MessageDetector) RegisterHandler(msgType string, handler MessageHandler) {
	md.mu.Lock()
	defer md.mu.Unlock()
	md.handlers[msgType] = handler
}

// DetectAndRoute routes a message to the appropriate handler
func (md *MessageDetector) DetectAndRoute(message *EDIFACTMessage) error {
	md.mu.RLock()
	handler, exists := md.handlers[message.MessageType]
	md.mu.RUnlock()

	if exists {
		fmt.Printf("Routing %s message to %s handler\n", message.MessageType, handler.GetMessageType())
		return handler.Process(message)
	}

	fmt.Printf("No specific handler for %s, using fallback\n", message.MessageType)
	return md.fallback.Process(message)
}

// InvoiceHandler processes INVOIC messages
type InvoiceHandler struct {
	validator  *InvoiceValidator
	processor  *InvoiceProcessor
	notifier   *NotificationService
	repository *InvoiceRepository
}

// NewInvoiceHandler creates a new invoice handler
func NewInvoiceHandler() *InvoiceHandler {
	return &InvoiceHandler{
		validator:  &InvoiceValidator{},
		processor:  &InvoiceProcessor{},
		notifier:   &NotificationService{},
		repository: &InvoiceRepository{},
	}
}

func (h *InvoiceHandler) GetMessageType() string {
	return "INVOIC"
}

func (h *InvoiceHandler) Process(message *EDIFACTMessage) error {
	fmt.Println("Processing invoice message...")

	// Extract invoice data
	invoice := h.extractInvoice(message)
	fmt.Printf("Extracted invoice: %s\n", invoice.Number)

	// Validate business rules
	if err := h.validator.Validate(invoice); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	// Process invoice
	result := h.processor.Process(invoice)
	fmt.Printf("Processed invoice with total: $%.2f\n", result.Total)

	// Store in database
	if err := h.repository.Save(result); err != nil {
		return fmt.Errorf("failed to save: %v", err)
	}

	// Send notifications
	return h.notifier.Notify(result)
}

func (h *InvoiceHandler) extractInvoice(message *EDIFACTMessage) *Invoice {
	invoice := &Invoice{}

	for _, segment := range message.Segments {
		switch segment.Tag {
		case "BGM":
			if len(segment.Elements) > 1 {
				invoice.Number = segment.Elements[1]
			}
		case "DTM":
			if len(segment.Elements) > 0 {
				parts := strings.Split(segment.Elements[0], ":")
				if len(parts) > 1 {
					invoice.Date = parts[1]
				}
			}
		case "LIN":
			if len(segment.Elements) > 2 {
				parts := strings.Split(segment.Elements[2], ":")
				if len(parts) > 0 {
					invoice.Items = append(invoice.Items, &InvoiceItem{
						ProductCode: parts[0],
					})
				}
			}
		}
	}

	return invoice
}

// OrderHandler processes ORDERS messages
type OrderHandler struct {
	validator  *OrderValidator
	processor  *OrderProcessor
	notifier   *NotificationService
	repository *OrderRepository
}

// NewOrderHandler creates a new order handler
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		validator:  &OrderValidator{},
		processor:  &OrderProcessor{},
		notifier:   &NotificationService{},
		repository: &OrderRepository{},
	}
}

func (h *OrderHandler) GetMessageType() string {
	return "ORDERS"
}

func (h *OrderHandler) Process(message *EDIFACTMessage) error {
	fmt.Println("Processing order message...")

	// Extract order data
	order := h.extractOrder(message)
	fmt.Printf("Extracted order: %s\n", order.Number)

	// Validate business rules
	if err := h.validator.Validate(order); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	// Process order
	result := h.processor.Process(order)
	fmt.Printf("Processed order with status: %s\n", result.Status)

	// Store in database
	if err := h.repository.Save(result); err != nil {
		return fmt.Errorf("failed to save: %v", err)
	}

	// Send notifications
	return h.notifier.Notify(result)
}

func (h *OrderHandler) extractOrder(message *EDIFACTMessage) *Order {
	order := &Order{}

	for _, segment := range message.Segments {
		switch segment.Tag {
		case "BGM":
			if len(segment.Elements) > 1 {
				order.Number = segment.Elements[1]
			}
		case "DTM":
			if len(segment.Elements) > 0 {
				parts := strings.Split(segment.Elements[0], ":")
				if len(parts) > 1 {
					order.Date = parts[1]
				}
			}
		}
	}

	return order
}

// GenericHandler processes unknown message types
type GenericHandler struct{}

func (h *GenericHandler) GetMessageType() string {
	return "GENERIC"
}

func (h *GenericHandler) Process(message *EDIFACTMessage) error {
	fmt.Printf("Processing generic message of type: %s\n", message.MessageType)
	fmt.Printf("Message contains %d segments\n", len(message.Segments))
	return nil
}

// CompositeHandler applies multiple handlers in sequence
type CompositeHandler struct {
	handlers []MessageHandler
}

// NewCompositeHandler creates a new composite handler
func NewCompositeHandler(handlers ...MessageHandler) *CompositeHandler {
	return &CompositeHandler{handlers: handlers}
}

func (ch *CompositeHandler) Process(message *EDIFACTMessage) error {
	fmt.Println("Processing with composite handler...")

	for i, handler := range ch.handlers {
		fmt.Printf("Step %d: Applying %s handler\n", i+1, handler.GetMessageType())
		if err := handler.Process(message); err != nil {
			return fmt.Errorf("handler %s failed: %v", handler.GetMessageType(), err)
		}
	}

	return nil
}

// Business logic components (simplified implementations)
type Invoice struct {
	Number string
	Date   string
	Items  []*InvoiceItem
}

type InvoiceItem struct {
	ProductCode string
}

type Order struct {
	Number string
	Date   string
}

type InvoiceValidator struct{}
type InvoiceProcessor struct{}
type NotificationService struct{}
type InvoiceRepository struct{}
type OrderValidator struct{}
type OrderProcessor struct{}
type OrderRepository struct{}

func (v *InvoiceValidator) Validate(invoice *Invoice) error {
	if invoice.Number == "" {
		return fmt.Errorf("invoice number is required")
	}
	return nil
}

func (p *InvoiceProcessor) Process(invoice *Invoice) *ProcessedInvoice {
	return &ProcessedInvoice{
		Number: invoice.Number,
		Total:  100.50,
	}
}

func (n *NotificationService) Notify(result interface{}) error {
	fmt.Printf("Sending notification for: %v\n", result)
	return nil
}

func (r *InvoiceRepository) Save(invoice *ProcessedInvoice) error {
	fmt.Printf("Saving invoice to database: %s\n", invoice.Number)
	return nil
}

func (v *OrderValidator) Validate(order *Order) error {
	if order.Number == "" {
		return fmt.Errorf("order number is required")
	}
	return nil
}

func (p *OrderProcessor) Process(order *Order) *ProcessedOrder {
	return &ProcessedOrder{
		Number: order.Number,
		Status: "PROCESSED",
	}
}

func (r *OrderRepository) Save(order *ProcessedOrder) error {
	fmt.Printf("Saving order to database: %s\n", order.Number)
	return nil
}

type ProcessedInvoice struct {
	Number string
	Total  float64
}

type ProcessedOrder struct {
	Number string
	Status string
}

// ParseEDIFACTMessage parses a raw EDIFACT message
func ParseEDIFACTMessage(rawContent string) *EDIFACTMessage {
	message := &EDIFACTMessage{
		RawContent: rawContent,
		Segments:   []*Segment{},
	}

	segments := strings.Split(rawContent, "'")
	for _, segmentStr := range segments {
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
	fmt.Println("=== Custom Message Handlers (Lesson 2) ===")

	// Create message detector
	detector := NewMessageDetector()

	// Register handlers
	detector.RegisterHandler("INVOIC", NewInvoiceHandler())
	detector.RegisterHandler("ORDERS", NewOrderHandler())

	// Test messages
	invoiceMessage := `UNH+1+INVOIC:D:97A:UN'BGM+380+INV12345+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM001:EN'QTY+12:100:PCE'PRI+AAA:25.50:CT'UNT+7+1'`
	orderMessage := `UNH+1+ORDERS:D:97A:UN'BGM+220+PO67890+9'DTM+137:20231201:102'NAD+BY+++ACME CORP'LIN+1++ITEM002:EN'QTY+12:50:PCE'UNT+6+1'`
	unknownMessage := `UNH+1+UNKNOWN:D:97A:UN'BGM+999+TEST123+9'UNT+2+1'`

	// Process messages
	messages := []string{invoiceMessage, orderMessage, unknownMessage}

	for i, rawMessage := range messages {
		fmt.Printf("\n--- Processing Message %d ---\n", i+1)

		// Parse message
		message := ParseEDIFACTMessage(rawMessage)
		fmt.Printf("Detected message type: %s\n", message.MessageType)

		// Route to appropriate handler
		start := time.Now()
		if err := detector.DetectAndRoute(message); err != nil {
			fmt.Printf("Error processing message: %v\n", err)
		}
		duration := time.Since(start)
		fmt.Printf("Processing completed in %v\n", duration)
	}

	// Demonstrate composite handler
	fmt.Println("\n=== Composite Handler Demo ===")

	composite := NewCompositeHandler(
		&GenericHandler{},
		&InvoiceHandler{},
	)

	testMessage := ParseEDIFACTMessage(invoiceMessage)
	if err := composite.Process(testMessage); err != nil {
		fmt.Printf("Composite handler error: %v\n", err)
	}

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✅ Custom message handlers for different message types")
	fmt.Println("✅ Message routing and detection")
	fmt.Println("✅ Business logic integration")
	fmt.Println("✅ Handler composition patterns")
	fmt.Println("✅ Error handling and logging")
}
