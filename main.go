package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/looksocial/edifact"
	"github.com/looksocial/edifact-examples/edifact_adapter"
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

	// Parse the EDIFACT message using the new simplified API
	message, err := edifact.Parse(edifactData)
	if err != nil {
		log.Fatalf("Failed to parse EDIFACT: %v", err)
	}

	// Create the IFTMBF adapter and convert to Booking model
	adapter := edifact_adapter.NewIFTMBFAdapter()
	result, err := adapter.Handle(message)
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

	// Demonstrate the new simplified API functions
	fmt.Println("\n=== Using New Simplified API ===")

	// Get message info
	info, err := edifact.GetMessageInfo(edifactData)
	if err != nil {
		log.Printf("Error getting message info: %v", err)
	} else {
		fmt.Printf("Message Info: %s\n", info.String())
	}

	// Get booking number
	bookingNumber, err := edifact.GetElementValue(edifactData, "BGM", 2)
	if err != nil {
		log.Printf("Error getting booking number: %v", err)
	} else {
		fmt.Printf("Booking Number: %s\n", bookingNumber)
	}

	// Validate message
	if edifact.IsValid(edifactData) {
		fmt.Println("Message is valid")
	} else {
		fmt.Println("Message is invalid")
	}
}
