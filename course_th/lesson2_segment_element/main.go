package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("\n=== ตัวอย่าง Segment และ Element (บทที่ 2) ===")
	msg := "BGM+380+12345678+9'DTM+137:20231201:102'"

	fmt.Println("\nข้อความ EDIFACT ตัวอย่าง:")
	fmt.Println(msg)

	fmt.Println("\n1. แยก Segment (ใช้ ')")
	segments := strings.Split(msg, "'")
	for i, seg := range segments {
		if seg == "" { continue }
		fmt.Printf("  %d: %s\n", i+1, seg)
	}

	fmt.Println("\n2. แยก Element ในแต่ละ Segment (ใช้ +)")
	for _, seg := range segments {
		if seg == "" { continue }
		elements := strings.Split(seg, "+")
		fmt.Printf("  Segment: %s\n", seg)
		for i, el := range elements {
			fmt.Printf("    Element %d: %s\n", i+1, el)
		}
	}

	fmt.Println("\n3. แยก Composite ใน DTM (ใช้ :)")
	for _, seg := range segments {
		if strings.HasPrefix(seg, "DTM+") {
			parts := strings.Split(seg, "+")
			if len(parts) > 1 {
				composite := strings.Split(parts[1], ":")
				for i, c := range composite {
					fmt.Printf("  Component %d: %s\n", i+1, c)
				}
			}
		}
	}
} 