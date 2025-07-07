# คอร์สเรียนรู้ EDIFACT และการประยุกต์ใช้กับ looksocial/edifact (ภาษาไทย)

คอร์สนี้เหมาะสำหรับผู้เริ่มต้นที่ต้องการเข้าใจไวยากรณ์ (Syntax) ของ UN/EDIFACT และนำไปประยุกต์ใช้กับแพ็กเกจ Go: [looksocial/edifact](https://github.com/looksocial/edifact) เพื่อแยกวิเคราะห์และแปลงข้อมูล EDI ในงานจริง

## 🏁 สิ่งที่คุณจะได้เรียนรู้
- เข้าใจโครงสร้างและไวยากรณ์ของ EDIFACT
- วิเคราะห์และแยกส่วน Segment, Element, Composite
- ใช้งานแพ็กเกจ looksocial/edifact เพื่ออ่าน/แปลงข้อความ
- สร้าง Adapter สำหรับข้อมูลเฉพาะ
- ตรวจสอบความถูกต้องของข้อความ EDI
- ประยุกต์ใช้กับกรณีจริงในธุรกิจ

## 🗂️ โครงสร้างคอร์ส

```
examples/course_th/
├── README.md                # ภาพรวมคอร์ส
├── lesson1_edifact_syntax/  # พื้นฐานไวยากรณ์ EDIFACT
├── lesson2_segment_element/ # ส่วนประกอบ Segment และ Element
├── lesson3_parse_message/   # การแยกวิเคราะห์ข้อความ EDIFACT
├── lesson4_detect_convert/  # การตรวจจับและแปลงข้อความด้วย looksocial/edifact
├── lesson5_custom_adapter/  # การสร้าง Adapter สำหรับข้อมูลเฉพาะ
├── lesson6_validation/      # การตรวจสอบความถูกต้องของข้อความ
└── lesson7_apply_realworld/ # ประยุกต์ใช้กับกรณีจริง
```

## 🚦 วิธีเรียน
- อ่าน README ของแต่ละบทเรียน
- ทดลองรันโค้ด Go (`main.go`) ในแต่ละบท
- ทำแบบฝึกหัดท้ายบทเพื่อทบทวนความเข้าใจ
- สามารถนำความรู้ไปประยุกต์ใช้กับงานจริงได้ทันที

## 🔗 แหล่งข้อมูลเพิ่มเติม
- [UN/EDIFACT Overview (TH)](https://www.gs1th.org/knowledge/edi/edifact)
- [looksocial/edifact GitHub](https://github.com/looksocial/edifact)
- [ตัวอย่าง EDI Message](https://www.unece.org/trade/untdid/welcome.html)

---

> **หมายเหตุ:** คอร์สนี้เหมาะสำหรับผู้มีพื้นฐาน Go เล็กน้อย หากยังไม่เคยเขียน Go สามารถศึกษาจาก [Go by Example (TH)](https://gobyexample.com/) ก่อนเริ่มคอร์สนี้ 