package main

import (
	"fmt"
	"strings"
	"unicode"
)

func isValidDate(date string) bool {
	if len(date) != 8 {
		return false
	}
	for _, r := range date {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	month := date[4:6]
	return month >= "01" && month <= "12"
}

func main() {
	fmt.Println("\n=== การตรวจสอบความถูกต้องของข้อความ (บทที่ 6) ===")
	msg := "BGM+380+INV12345+9'DTM+137:20231301:102'QTY+47:1000:PCE'"

	if !strings.HasSuffix(msg, "'") {
		fmt.Println("Error: ข้อความไม่มี segment terminator (' )")
	}

	for _, seg := range strings.Split(msg, "'") {
		if seg == "" { continue }
		parts := strings.Split(seg, "+")
		if parts[0] == "DTM" && len(parts) > 1 {
			comps := strings.Split(parts[1], ":")
			if len(comps) > 1 && !isValidDate(comps[1]) {
				fmt.Printf("Error: วันที่ไม่ถูกต้อง: %s\n", comps[1])
			}
		}
		if parts[0] == "BGM" && len(parts) > 1 && parts[1] != "380" {
			fmt.Printf("Error: BGM code ไม่ถูกต้อง: %s\n", parts[1])
		}
	}

	fmt.Println("\nตรวจสอบเสร็จสิ้น")
} 