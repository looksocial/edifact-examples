package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/looksocial/edifact/internal/model"
	"github.com/looksocial/edifact/pkg/edifact"
)

// Custom message handler for ORDERS messages
type ORDERSHandler struct{}

func NewORDERSHandler() *ORDERSHandler {
	return &ORDERSHandler{}
}

func (h *ORDERSHandler) CanHandle(messageType string) bool {
	return messageType == "ORDERS"
}

func (h *ORDERSHandler) Handle(message *model.Message) (interface{}, error) {
	result := &ORDERSMessage{
		MessageReference: "",
		OrderNumber:      "",
		OrderDate:        "",
		Buyer:            "",
		Supplier:         "",
		Items:            []OrderItem{},
	}

	// Extract information from segments
	for _, segment := range message.Segments {
		switch segment.Tag {
		case "UNH":
			if len(segment.Elements) > 0 {
				result.MessageReference = segment.Elements[0].Value
			}
		case "BGM":
			if len(segment.Elements) > 1 {
				result.OrderNumber = segment.Elements[1].Value
			}
		case "DTM":
			if len(segment.Elements) > 0 {
				dateType := segment.Elements[0].Value
				if dateType == "4" { // Order date
					if len(segment.Elements) > 1 {
						result.OrderDate = segment.Elements[1].Value
					}
				}
			}
		case "NAD":
			if len(segment.Elements) > 0 {
				partyType := segment.Elements[0].Value
				switch partyType {
				case "BY": // Buyer
					if len(segment.Elements) > 2 {
						result.Buyer = segment.Elements[2].Value
					}
				case "SU": // Supplier
					if len(segment.Elements) > 2 {
						result.Supplier = segment.Elements[2].Value
					}
				}
			}
		case "LIN":
			item := OrderItem{}
			if len(segment.Elements) > 0 {
				item.LineNumber = segment.Elements[0].Value
			}
			if len(segment.Elements) > 2 && segment.Elements[2].IsComposite {
				item.ProductCode = segment.Elements[2].Components[0]
			}
			result.Items = append(result.Items, item)
		case "QTY":
			if len(segment.Elements) > 0 {
				quantityType := segment.Elements[0].Value
				if quantityType == "21" { // Ordered quantity
					if len(segment.Elements) > 1 {
						if len(result.Items) > 0 {
							result.Items[len(result.Items)-1].Quantity = segment.Elements[1].Value
						}
					}
				}
			}
		}
	}

	return result, nil
}

type ORDERSMessage struct {
	MessageReference string      `json:"message_reference"`
	OrderNumber      string      `json:"order_number"`
	OrderDate        string      `json:"order_date"`
	Buyer            string      `json:"buyer"`
	Supplier         string      `json:"supplier"`
	Items            []OrderItem `json:"items"`
}

type OrderItem struct {
	LineNumber  string `json:"line_number"`
	ProductCode string `json:"product_code"`
	Quantity    string `json:"quantity"`
}

func main() {
	fmt.Println("=== EDIFACT Library Advanced Usage Example ===\n")

	// 1. Custom configuration
	fmt.Println("1. Using custom configuration:")
	config := model.EDIFACTConfig{
		SegmentTerminator:  "'",
		ElementSeparator:   "+",
		ComponentSeparator: ":",
		ReleaseCharacter:   "?",
	}

	reader := edifact.NewReaderWithConfig(config)
	detector := edifact.NewDetectorWithConfig(config)

	// 2. Message detection
	fmt.Println("2. Message type detection:")
	edifactData := `UNH+1+ORDERS:D:93A:UN'
BGM+220+ORD123456'
DTM+4:20231201:102'
NAD+BY+++BUYER COMPANY'
NAD+SU+++SUPPLIER COMPANY'
LIN+1++ABC123:EN'
QTY+21:5'
UNT+7+1'`

	messageType, err := detector.DetectMessageType(edifactData)
	if err != nil {
		log.Fatal("Error detecting message type:", err)
	}
	fmt.Printf("Detected message type: %s\n", messageType)

	messageInfo, err := detector.DetectMessageInfo(edifactData)
	if err != nil {
		log.Fatal("Error detecting message info:", err)
	}
	fmt.Printf("Message info: Type=%s, Version=%s\n\n", messageInfo.Type, messageInfo.Version)

	// 3. Reading and parsing
	fmt.Println("3. Reading and parsing EDIFACT message:")
	message, err := reader.ReadString(edifactData)
	if err != nil {
		log.Fatal("Error reading message:", err)
	}

	fmt.Printf("Message type: %s\n", message.Type)
	fmt.Printf("Number of segments: %d\n", len(message.Segments))

	// Access specific segments
	if unhSegment := message.GetSegmentByTag("UNH"); unhSegment != nil {
		fmt.Printf("Message reference: %s\n", unhSegment.GetElement(1).Value)
	}

	if bgmSegment := message.GetSegmentByTag("BGM"); bgmSegment != nil {
		fmt.Printf("Order number: %s\n", bgmSegment.GetElement(2).Value)
	}
	fmt.Println()

	// 4. Custom handler registration
	fmt.Println("4. Using custom message handler:")
	converter := edifact.NewConverter()
	converter.RegisterHandler("ORDERS", NewORDERSHandler())

	// Convert with custom handler
	result, err := converter.ConvertToStructured(edifactData)
	if err != nil {
		log.Fatal("Error converting with custom handler:", err)
	}

	// Type assert to our custom type
	if ordersMessage, ok := result.(*ORDERSMessage); ok {
		jsonData, _ := json.MarshalIndent(ordersMessage, "", "  ")
		fmt.Println("Converted ORDERS message:")
		fmt.Println(string(jsonData))
	}
	fmt.Println()

	// 5. Segment-level operations
	fmt.Println("5. Segment-level operations:")
	segmentStr := "NAD+BY+++BUYER COMPANY'"
	segment, err := reader.ReadSegment(segmentStr)
	if err != nil {
		log.Fatal("Error reading segment:", err)
	}

	fmt.Printf("Segment tag: %s\n", segment.Tag)
	fmt.Printf("Number of elements: %d\n", len(segment.Elements))
	for i, element := range segment.Elements {
		fmt.Printf("Element %d: %s\n", i+1, element.String())
	}
	fmt.Println()

	// 6. Error handling demonstration
	fmt.Println("6. Error handling demonstration:")
	invalidData := "INVALID+EDIFACT+DATA"
	_, err = reader.ReadString(invalidData)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}
