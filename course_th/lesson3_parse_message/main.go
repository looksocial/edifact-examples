package main

import (
	"fmt"
	"strings"
)

type Segment struct {
	Tag      string
	Elements []string
}

func parseEDIFACT(msg string) []Segment {
	segments := []Segment{}
	for _, seg := range strings.Split(msg, "'") {
		if seg == "" {
			continue
		}
		parts := strings.Split(seg, "+")
		if len(parts) == 0 {
			continue
		}
		segments = append(segments, Segment{
			Tag:      parts[0],
			Elements: parts[1:],
		})
	}
	return segments
}

func main() {
	fmt.Println("\n=== การแยกวิเคราะห์ข้อความ EDIFACT (บทที่ 3) ===")
	msg := "UNH+1+INVOIC:D:97A:UN'BGM+380+12345678+9'DTM+137:20231201:102'"

	segments := parseEDIFACT(msg)
	for i, seg := range segments {
		fmt.Printf("Segment %d: %s\n", i+1, seg.Tag)
		for j, el := range seg.Elements {
			fmt.Printf("  Element %d: %s\n", j+1, el)
			if strings.Contains(el, ":") {
				comps := strings.Split(el, ":")
				for k, c := range comps {
					fmt.Printf("    Component %d: %s\n", k+1, c)
				}
			}
		}
	}
}
