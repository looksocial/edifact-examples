package main

import (
	"fmt"
	"strings"
)

type Invoice struct {
	Number string
	Date   string
	Amount string
}

func ParseInvoice(msg string) Invoice {
	inv := Invoice{}
	for _, seg := range strings.Split(msg, "'") {
		if seg == "" {
			continue
		}
		parts := strings.Split(seg, "+")
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "BGM":
			if len(parts) > 2 {
				inv.Number = parts[2]
			}
		case "DTM":
			if len(parts) > 1 {
				comps := strings.Split(parts[1], ":")
				if len(comps) > 1 {
					inv.Date = comps[1]
				}
			}
		case "QTY":
			if len(parts) > 1 {
				comps := strings.Split(parts[1], ":")
				if len(comps) > 1 {
					inv.Amount = comps[1]
				}
			}
		}
	}
	return inv
}

func main() {
	fmt.Println("\n=== การสร้าง Custom Adapter (บทที่ 5) ===")
	msg := "BGM+380+INV12345+9'DTM+137:20231201:102'QTY+47:1000:PCE'"

	inv := ParseInvoice(msg)
	fmt.Println("Invoice struct ที่ได้:")
	fmt.Printf("  Number: %s\n", inv.Number)
	fmt.Printf("  Date:   %s\n", inv.Date)
	fmt.Printf("  Amount: %s\n", inv.Amount)
}
