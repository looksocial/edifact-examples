package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/looksocial/edifact/internal/model"
	"github.com/looksocial/edifact/pkg/edifact"
)

// SimpleAdapter demonstrates the exact structure you requested
// BGM+220+ORD123456 becomes:
// 1. BGM
// 2. 220
// 3. ORD123456
type SimpleAdapter struct{}

func NewSimpleAdapter() *SimpleAdapter {
	return &SimpleAdapter{}
}

func (a *SimpleAdapter) CanHandle(messageType string) bool {
	return true // Can handle any message type
}

func (a *SimpleAdapter) Handle(message *model.Message) (interface{}, error) {
	result := &SimpleAdapterResult{
		MessageType: message.Type,
		MessageRef:  "",
		Elements:    make([]string, 0),
	}

	// Extract message reference from UNH segment
	if unhSegment := message.GetSegmentByTag("UNH"); unhSegment != nil {
		if len(unhSegment.Elements) > 0 {
			result.MessageRef = unhSegment.Elements[0].Value
		}
	}

	// Process each segment and flatten all elements
	for _, segment := range message.Segments {
		// Add segment tag as first element
		result.Elements = append(result.Elements, segment.Tag)

		// Add each element value
		for _, element := range segment.Elements {
			if element.IsComposite {
				// For composite elements, join components with colon
				result.Elements = append(result.Elements,
					strings.Join(element.Components, ":"))
			} else {
				result.Elements = append(result.Elements, element.Value)
			}
		}
	}

	return result, nil
}

type SimpleAdapterResult struct {
	MessageType string   `json:"message_type"`
	MessageRef  string   `json:"message_ref"`
	Elements    []string `json:"elements"`
}

func (r *SimpleAdapterResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r *SimpleAdapterResult) ToJSONString() (string, error) {
	jsonData, err := r.ToJSON()
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (r *SimpleAdapterResult) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Message Type: %s\n", r.MessageType))
	sb.WriteString(fmt.Sprintf("Message Ref: %s\n", r.MessageRef))
	sb.WriteString("Elements:\n")

	for i, element := range r.Elements {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, element))
	}

	return sb.String()
}

// GetElement returns the element at the specified position (1-based)
func (r *SimpleAdapterResult) GetElement(position int) string {
	if position <= 0 || position > len(r.Elements) {
		return ""
	}
	return r.Elements[position-1]
}

// GetElementRange returns elements in the specified range
func (r *SimpleAdapterResult) GetElementRange(start, end int) []string {
	if start <= 0 || end > len(r.Elements) || start > end {
		return nil
	}
	return r.Elements[start-1 : end]
}

func main() {
	fmt.Println("=== Simple EDIFACT Adapter Example ===\n")

	// Sample EDIFACT messages
	messages := map[string]string{
		"BGM Example": `BGM+220+ORD123456'`,

		"UNH Example": `UNH+1+IFTMIN:D:93A:UN'`,

		"Complete IFTMIN": `UNH+1+IFTMIN:D:93A:UN'
BGM+380+123456789'
DTM+137:20231201:102'
NAD+CZ+++SENDER COMPANY'
NAD+CN+++RECEIVER COMPANY'
UNT+6+1'`,

		"Complete INVOIC": `UNH+1+INVOIC:D:93A:UN'
BGM+380+INV123456'
DTM+137:20231201:102'
NAD+SE+++SELLER COMPANY'
NAD+BY+++BUYER COMPANY'
LIN+1++123456:EN'
QTY+47:10'
UNT+8+1'`,

		"Complete ORDERS": `UNH+1+ORDERS:D:93A:UN'
BGM+220+ORD123456'
DTM+4:20231201:102'
NAD+BY+++BUYER COMPANY'
NAD+SU+++SUPPLIER COMPANY'
LIN+1++ABC123:EN'
QTY+21:5'
UNT+7+1'`,
	}

	// Create converter and register simple adapter
	converter := edifact.NewConverter()
	converter.RegisterHandler("SIMPLE_ADAPTER", NewSimpleAdapter())

	// Test each message
	for name, edifactData := range messages {
		fmt.Printf("=== %s ===\n", name)
		fmt.Printf("Input EDIFACT:\n%s\n\n", edifactData)

		// Convert using simple adapter
		result, err := converter.ConvertToStructured(edifactData)
		if err != nil {
			log.Printf("Error converting %s: %v\n", name, err)
			continue
		}

		// Type assert to our simple adapter result
		if simpleResult, ok := result.(*SimpleAdapterResult); ok {
			fmt.Printf("Output (Numbered Elements):\n")
			for i, element := range simpleResult.Elements {
				fmt.Printf("%d. %s\n", i+1, element)
			}
			fmt.Println()

			// Demonstrate element access
			fmt.Printf("Element Access Examples:\n")
			fmt.Printf("- Element 1: %s\n", simpleResult.GetElement(1))
			fmt.Printf("- Element 2: %s\n", simpleResult.GetElement(2))
			fmt.Printf("- Element 3: %s\n", simpleResult.GetElement(3))

			if len(simpleResult.Elements) > 5 {
				fmt.Printf("- Elements 4-6: %v\n", simpleResult.GetElementRange(4, 6))
			}
			fmt.Println()

			// Show JSON output
			jsonData, err := simpleResult.ToJSONString()
			if err != nil {
				log.Printf("Error converting to JSON: %v\n", err)
			} else {
				fmt.Printf("JSON Output:\n%s\n", jsonData)
			}
		}

		fmt.Printf("\n" + strings.Repeat("-", 80) + "\n\n")
	}

	// Demonstrate how to create a custom adapter for specific message types
	fmt.Println("=== Custom Adapter for Specific Message Types ===")

	// Create a custom adapter that only handles ORDERS messages
	ordersAdapter := &OrdersSpecificAdapter{}
	converter.RegisterHandler("ORDERS", ordersAdapter)

	// Test with ORDERS message
	ordersData := `UNH+1+ORDERS:D:93A:UN'
BGM+220+ORD123456'
DTM+4:20231201:102'
NAD+BY+++BUYER COMPANY'
NAD+SU+++SUPPLIER COMPANY'
UNT+5+1'`

	fmt.Printf("Testing ORDERS-specific adapter:\n")
	fmt.Printf("Input: %s\n\n", ordersData)

	result, err := converter.ConvertToStructured(ordersData)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		if ordersResult, ok := result.(*OrdersResult); ok {
			fmt.Printf("Orders Result:\n")
			fmt.Printf("- Order Number: %s\n", ordersResult.OrderNumber)
			fmt.Printf("- Order Date: %s\n", ordersResult.OrderDate)
			fmt.Printf("- Buyer: %s\n", ordersResult.Buyer)
			fmt.Printf("- Supplier: %s\n", ordersResult.Supplier)
		}
	}

	fmt.Println("\n=== Benefits of Simple Adapter ===")
	fmt.Println("1. Universal - works with any EDIFACT message type")
	fmt.Println("2. Simple - just numbered elements, easy to understand")
	fmt.Println("3. Flexible - can access any element by position")
	fmt.Println("4. Extensible - easy to create custom adapters for specific needs")
	fmt.Println("5. Future-proof - no need to update for new message types")
}

// OrdersSpecificAdapter demonstrates a custom adapter for ORDERS messages
type OrdersSpecificAdapter struct{}

func (a *OrdersSpecificAdapter) CanHandle(messageType string) bool {
	return messageType == "ORDERS"
}

func (a *OrdersSpecificAdapter) Handle(message *model.Message) (interface{}, error) {
	result := &OrdersResult{}

	// Extract order number from BGM segment
	if bgmSegment := message.GetSegmentByTag("BGM"); bgmSegment != nil {
		if len(bgmSegment.Elements) > 1 {
			result.OrderNumber = bgmSegment.Elements[1].Value
		}
	}

	// Extract order date from DTM segment
	if dtmSegment := message.GetSegmentByTag("DTM"); dtmSegment != nil {
		if len(dtmSegment.Elements) > 0 {
			dateType := dtmSegment.Elements[0].Value
			if dateType == "4" && len(dtmSegment.Elements) > 1 {
				result.OrderDate = dtmSegment.Elements[1].Value
			}
		}
	}

	// Extract buyer and supplier from NAD segments
	for _, segment := range message.Segments {
		if segment.Tag == "NAD" && len(segment.Elements) > 0 {
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
	}

	return result, nil
}

type OrdersResult struct {
	OrderNumber string `json:"order_number"`
	OrderDate   string `json:"order_date"`
	Buyer       string `json:"buyer"`
	Supplier    string `json:"supplier"`
}
