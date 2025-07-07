package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/looksocial/edifact/internal/message"
	"github.com/looksocial/edifact/internal/model"
	"github.com/looksocial/edifact/pkg/edifact"
)

// CustomHandler demonstrates how to create a custom handler for any message type
type CustomHandler struct {
	messageType string
}

func NewCustomHandler(messageType string) *CustomHandler {
	return &CustomHandler{
		messageType: messageType,
	}
}

func (h *CustomHandler) CanHandle(messageType string) bool {
	return messageType == h.messageType
}

func (h *CustomHandler) Handle(message *model.Message) (interface{}, error) {
	// Create a custom structure for any message type
	result := &CustomMessage{
		MessageType: message.Type,
		MessageRef:  "",
		Data:        make(map[string]interface{}),
	}

	// Extract message reference from UNH segment
	if unhSegment := message.GetSegmentByTag("UNH"); unhSegment != nil {
		if len(unhSegment.Elements) > 0 {
			result.MessageRef = unhSegment.Elements[0].Value
		}
	}

	// Process each segment and create a custom structure
	for _, segment := range message.Segments {
		segmentData := make(map[string]interface{})

		// Add elements as numbered fields
		for i, element := range segment.Elements {
			key := fmt.Sprintf("element_%d", i+1)
			if element.IsComposite {
				segmentData[key] = element.Components
			} else {
				segmentData[key] = element.Value
			}
		}

		// Add element count
		segmentData["element_count"] = len(segment.Elements)

		// Store segment data
		result.Data[segment.Tag] = segmentData
	}

	return result, nil
}

type CustomMessage struct {
	MessageType string                 `json:"message_type"`
	MessageRef  string                 `json:"message_ref"`
	Data        map[string]interface{} `json:"data"`
}

func (m *CustomMessage) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

func (m *CustomMessage) ToJSONString() (string, error) {
	jsonData, err := m.ToJSON()
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func main() {
	fmt.Println("=== EDIFACT Generic Adapters Example ===")

	// Sample EDIFACT messages of different types
	messages := map[string]string{
		"IFTMIN": `UNH+1+IFTMIN:D:93A:UN'
BGM+380+123456789'
DTM+137:20231201:102'
NAD+CZ+++SENDER COMPANY LTD'
NAD+CN+++RECEIVER COMPANY INC'
TDT+20++1++ROAD'
LOC+5+PLACE:12345:9'
UNT+8+1'`,

		"INVOIC": `UNH+1+INVOIC:D:93A:UN'
BGM+380+INV123456'
DTM+137:20231201:102'
NAD+SE+++SELLER COMPANY LTD'
NAD+BY+++BUYER COMPANY INC'
LIN+1++123456:EN'
QTY+47:10'
PRI+AAA:100.00'
MOA+77:1000.00'
UNT+9+1'`,

		"ORDERS": `UNH+1+ORDERS:D:93A:UN'
BGM+220+ORD123456'
DTM+4:20231201:102'
NAD+BY+++BUYER COMPANY'
NAD+SU+++SUPPLIER COMPANY'
LIN+1++ABC123:EN'
QTY+21:5'
UNT+7+1'`,

		"VERMAS": `UNH+1+VERMAS:D:93A:UN'
BGM+745+VER123456'
DTM+137:20231201:102'
NAD+SE+++SELLER COMPANY'
NAD+BY+++BUYER COMPANY'
LIN+1++PROD123:EN'
QTY+47:10'
UNT+7+1'`,

		"IFTMBC": `UNH+1+IFTMBC:D:93A:UN'
BGM+380+BC123456'
DTM+137:20231201:102'
NAD+CZ+++SENDER COMPANY'
NAD+CN+++RECEIVER COMPANY'
TDT+20++1++ROAD'
UNT+6+1'`,
	}

	// Create converter with generic handlers
	converter := edifact.NewConverter()

	// Register custom handlers for specific message types
	converter.RegisterHandler("ORDERS", NewCustomHandler("ORDERS"))
	converter.RegisterHandler("VERMAS", NewCustomHandler("VERMAS"))
	converter.RegisterHandler("IFTMBC", NewCustomHandler("IFTMBC"))

	// Test each message type
	for msgType, edifactData := range messages {
		fmt.Printf("=== Processing %s Message ===\n", msgType)
		fmt.Printf("Input EDIFACT:\n%s\n\n", edifactData)

		// Convert to structured data
		result, err := converter.ConvertToStructured(edifactData)
		if err != nil {
			log.Printf("Error converting %s: %v\n", msgType, err)
			continue
		}

		// Convert to JSON
		jsonData, err := converter.ConvertToJSON(edifactData)
		if err != nil {
			log.Printf("Error converting %s to JSON: %v\n", msgType, err)
			continue
		}

		fmt.Printf("Output JSON:\n%s\n\n", string(jsonData))

		// Demonstrate different ways to access the data
		fmt.Printf("Data Access Examples:\n")

		switch typedResult := result.(type) {
		case *message.SimpleMessage:
			fmt.Printf("- Message Type: %s\n", typedResult.MessageType)
			fmt.Printf("- Message Ref: %s\n", typedResult.MessageRef)

			// Access BGM segment
			if bgmSegment := typedResult.GetSegmentByTag("BGM"); bgmSegment != nil {
				fmt.Printf("- BGM Element 1: %s\n", bgmSegment.GetElementValue(1))
				fmt.Printf("- BGM Element 2: %s\n", bgmSegment.GetElementValue(2))
			}

			// Access UNH segment
			if unhSegment := typedResult.GetSegmentByTag("UNH"); unhSegment != nil {
				fmt.Printf("- UNH Element 1: %s\n", unhSegment.GetElementValue(1))
				fmt.Printf("- UNH Element 2: %s\n", unhSegment.GetElementValue(2))
			}

		case *message.IFTMINMessage:
			fmt.Printf("- Message Type: %s\n", typedResult.MessageReference)
			fmt.Printf("- Document Number: %s\n", typedResult.DocumentNumber)

		case *message.INVOICMessage:
			fmt.Printf("- Message Type: %s\n", typedResult.MessageReference)
			fmt.Printf("- Invoice Number: %s\n", typedResult.InvoiceNumber)

		case *CustomMessage:
			fmt.Printf("- Message Type: %s\n", typedResult.MessageType)
			fmt.Printf("- Message Ref: %s\n", typedResult.MessageRef)

			// Access BGM segment data
			if bgmData, exists := typedResult.Data["BGM"]; exists {
				if bgmMap, ok := bgmData.(map[string]interface{}); ok {
					fmt.Printf("- BGM Element 1: %v\n", bgmMap["element_1"])
					fmt.Printf("- BGM Element 2: %v\n", bgmMap["element_2"])
				}
			}

		default:
			fmt.Printf("- Unknown result type: %T\n", result)
		}

		fmt.Printf("\n" + strings.Repeat("-", 80) + "\n\n")
	}

	// Demonstrate how to create a simple reader for any message type
	fmt.Println("=== Simple Reader Example ===")
	reader := edifact.NewReader()

	// Read any EDIFACT message
	message, err := reader.ReadString(messages["ORDERS"])
	if err != nil {
		log.Fatal("Error reading message:", err)
	}

	fmt.Printf("Parsed Message Type: %s\n", message.Type)
	fmt.Printf("Number of Segments: %d\n", len(message.Segments))

	// Access segments directly
	for i, segment := range message.Segments {
		fmt.Printf("Segment %d: %s\n", i+1, segment.Tag)
		for j, element := range segment.Elements {
			if element.IsComposite {
				fmt.Printf("  Element %d: %s\n", j+1, strings.Join(element.Components, ":"))
			} else {
				fmt.Printf("  Element %d: %s\n", j+1, element.Value)
			}
		}
	}

	fmt.Println("\n=== Generic Handler Benefits ===")
	fmt.Println("1. No need to create specific handlers for each message type")
	fmt.Println("2. Can handle any EDIFACT message type (IFTMIN, INVOIC, ORDERS, VERMAS, IFTMBC, etc.)")
	fmt.Println("3. Simple, flat structure that's easy to work with")
	fmt.Println("4. Extensible - users can create their own custom handlers")
	fmt.Println("5. Consistent API across all message types")
}
