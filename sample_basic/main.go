package main

import (
	"fmt"
	"log"

	"github.com/looksocial/edifact/pkg/edifact"
)

func main() {
	fmt.Println("=== EDIFACT Library Basic Usage Example ===")

	// Create a converter
	converter := edifact.NewConverter()

	// Sample EDIFACT IFTMIN message
	iftminData := `UNH+1+IFTMIN:D:93A:UN'
BGM+380+123456789'
DTM+137:20231201:102'
NAD+CZ+++SENDER COMPANY LTD'
NAD+CN+++RECEIVER COMPANY INC'
TDT+20++1++ROAD'
LOC+5+PLACE:12345:9'
UNT+8+1'`

	fmt.Println("1. Converting IFTMIN message to JSON:")
	fmt.Println("Input EDIFACT:")
	fmt.Println(iftminData)
	fmt.Println()

	// Convert to JSON
	jsonData, err := converter.ConvertToJSONString(iftminData)
	if err != nil {
		log.Fatal("Error converting IFTMIN:", err)
	}

	fmt.Println("Output JSON:")
	fmt.Println(jsonData)
	fmt.Println()

	// Sample EDIFACT INVOIC message
	invoicData := `UNH+1+INVOIC:D:93A:UN'
BGM+380+INV123456'
DTM+137:20231201:102'
NAD+SE+++SELLER COMPANY LTD'
NAD+BY+++BUYER COMPANY INC'
LIN+1++123456:EN'
QTY+47:10'
PRI+AAA:100.00'
MOA+77:1000.00'
TAX+7+VAT+++::20'
MOA+125:200.00'
MOA+79:1200.00'
UNT+12+1'`

	fmt.Println("2. Converting INVOIC message to JSON:")
	fmt.Println("Input EDIFACT:")
	fmt.Println(invoicData)
	fmt.Println()

	// Convert to JSON
	jsonData2, err := converter.ConvertToJSONString(invoicData)
	if err != nil {
		log.Fatal("Error converting INVOIC:", err)
	}

	fmt.Println("Output JSON:")
	fmt.Println(jsonData2)
	fmt.Println()

	// Show supported message types
	fmt.Println("3. Supported message types:")
	supportedTypes := converter.GetSupportedMessageTypes()
	for _, msgType := range supportedTypes {
		fmt.Printf("- %s\n", msgType)
	}
}
