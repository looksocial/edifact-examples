package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/looksocial/edifact/examples/bookings/edifact_adapter"
	"github.com/looksocial/edifact/pkg/edifact"
)

func main() {
	// Example IFTMBF EDIFACT message (simplified)
	edifactData := `UNH+1+IFTMBF:D:93A:UN'
BGM+335+BOOK123456'
NAD+CA+++BOOKING PARTY LTD'
NAD+CN+++CONSIGNEE INC'
TDT+20++VESSELNAME+VOY123'
LOC+9+SGSIN:139:6'
LOC+11+NLRTM:139:6'
EQD+CN+CONT1234567'
UNT+9+1'`

	// Create a converter and register the IFTMBF adapter
	converter := edifact.NewConverter()
	converter.RegisterHandler("IFTMBF", edifact_adapter.NewIFTMBFAdapter())

	// Convert the EDIFACT message to a Booking model
	result, err := converter.ConvertToStructured(edifactData)
	if err != nil {
		log.Fatalf("Failed to convert EDIFACT: %v", err)
	}

	// Marshal to JSON for demonstration
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal to JSON: %v", err)
	}

	fmt.Println("Booking JSON:")
	fmt.Println(string(jsonData))
}
