# บทที่ 1: พื้นฐานไวยากรณ์ EDIFACT

## 🎯 เป้าหมายการเรียนรู้
- เข้าใจโครงสร้างข้อความ EDIFACT
- รู้จัก delimiter (ตัวแบ่งข้อมูล) แต่ละประเภท
- อ่านและแยก segment, element, composite ได้

## 🔍 EDIFACT คืออะไร?
EDIFACT (Electronic Data Interchange For Administration, Commerce and Transport) คือมาตรฐานสากลสำหรับการแลกเปลี่ยนข้อมูลธุรกิจระหว่างองค์กร

## 🧩 โครงสร้างพื้นฐาน
- **Segment**: กลุ่มข้อมูล เช่น UNH, BGM, DTM
- **Element**: ข้อมูลย่อยใน segment
- **Composite**: element ที่มีหลายค่าคั่นด้วย `:`
- **Delimiter**: ตัวแบ่ง เช่น `'` `+` `:` `?`

### ตัวอย่างข้อความ EDIFACT
```
UNH+1+INVOIC:D:97A:UN'
BGM+380+12345678+9'
DTM+137:20231201:102'
```
- `'` = จบ segment
- `+` = แบ่ง element
- `:` = แบ่ง component ใน composite
- `?` = escape character

## 🛠️ ทดลองรันโค้ด
ดูตัวอย่างใน `main.go` แล้วรัน:
```bash
cd lesson1_edifact_syntax
go run main.go
```

## 📝 แบบฝึกหัด
1. ข้อความนี้มี segment อะไรบ้าง?
   ```
   UNH+1+ORDERS:D:97A:UN'
   BGM+220+PO12345+9'
   DTM+4:20230401:102'
   ```
2. อธิบายหน้าที่ของ delimiter แต่ละตัว
3. ลองเปลี่ยนข้อความใน main.go แล้วสังเกตผลลัพธ์

## 🔑 สรุป
- EDIFACT ใช้ delimiter เพื่อแยกข้อมูลแต่ละส่วน
- เข้าใจ segment, element, composite คือพื้นฐานสำคัญ
- พร้อมเรียนรู้การแยกวิเคราะห์ในบทถัดไป 