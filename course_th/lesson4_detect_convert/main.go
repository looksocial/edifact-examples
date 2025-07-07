package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// mock ฟังก์ชันแปลง (ถ้าไม่มี package จริง)
func ConvertEDIFACTToJSON(msg string) (string, error) {
	segments := []map[string]interface{}{}
	for _, seg := range strings.Split(msg, "'") {
		if seg == "" {
			continue
		}
		parts := strings.Split(seg, "+")
		if len(parts) == 0 {
			continue
		}
		segMap := map[string]interface{}{
			"tag":      parts[0],
			"elements": parts[1:],
		}
		segments = append(segments, segMap)
	}
	b, err := json.MarshalIndent(segments, "", "  ")
	return string(b), err
}

func main() {
	fmt.Println("\n=== การตรวจจับและแปลงข้อความ EDIFACT (บทที่ 4) ===")
	msg := "UNH+1+INVOIC:D:97A:UN'BGM+380+12345678+9'DTM+137:20231201:102'"

	jsonStr, err := ConvertEDIFACTToJSON(msg)
	if err != nil {
		fmt.Println("แปลงเป็น JSON ไม่สำเร็จ:", err)
		return
	}
	fmt.Println("\nผลลัพธ์ JSON:")
	fmt.Println(jsonStr)
}
