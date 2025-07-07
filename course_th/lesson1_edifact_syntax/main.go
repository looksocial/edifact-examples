package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("\n=== ตัวอย่างพื้นฐาน EDIFACT (บทที่ 1) ===")
	msg := "UNH+1+INVOIC:D:97A:UN'BGM+380+12345678+9'DTM+137:20231201:102'"

	fmt.Println("\nข้อความ EDIFACT ตัวอย่าง:")
	fmt.Println(msg)

	fmt.Println("\n1. แยก Segment (ใช้ ')")
	segments := strings.Split(msg, "'")
	for i, seg := range segments {
		if seg == "" {
			continue
		}
		fmt.Printf("  %d: %s\n", i+1, seg)
	}

	fmt.Println("\n2. แยก Element ใน Segment แรก (ใช้ +)")
	if len(segments) > 0 {
		elements := strings.Split(segments[0], "+")
		for i, el := range elements {
			fmt.Printf("  %d: %s\n", i+1, el)
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

	fmt.Println("\n--- สรุป ---")
	fmt.Println("- ' ใช้แบ่ง Segment")
	fmt.Println("- + ใช้แบ่ง Element")
	fmt.Println("- : ใช้แบ่ง Component ใน Composite")
	fmt.Println("- ? ใช้ escape ตัวอักษรพิเศษ")
}
