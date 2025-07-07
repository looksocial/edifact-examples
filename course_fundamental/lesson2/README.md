# Lesson 2: Segments & Elements

## 📚 Learning Objectives

By the end of this lesson, you will understand:
- ✅ EDIFACT segment structure and types
- ✅ Simple and composite data elements
- ✅ Element positioning and mandatory/conditional elements
- ✅ Segment hierarchy and relationships

## 🔍 What are EDIFACT Segments?

Segments are the building blocks of EDIFACT messages. Each segment represents a logical group of related data elements and serves a specific purpose in the message structure.

### Segment Characteristics

1. **Fixed Tag**: Every segment starts with a 3-character tag
2. **Data Elements**: Contains one or more data elements
3. **Positional Structure**: Elements have specific positions
4. **Conditional/Mandatory**: Elements can be required or optional

## 🎯 Segment Types

### Service Segments
Control the message structure and transmission:
- `UNH` - Message Header
- `UNT` - Message Trailer
- `UNS` - Section Control
- `UNZ` - Interchange Trailer

### Data Segments
Contain the actual business data:
- `BGM` - Beginning of Message
- `DTM` - Date/Time/Period
- `NAD` - Name and Address
- `LIN` - Line Item
- `QTY` - Quantity
- `PRI` - Price Details

## 📖 Data Elements

### Simple Elements
Contain a single value:
```
UNH+1+INVOIC:D:97A:UN'
     ^
     └── Simple element (message reference number)
```

### Composite Elements
Contain multiple related components:
```
DTM+137:20231201:102'
     ^   ^         ^
     |   |         └── Format code
     |   └── Date value
     └── Qualifier
```

## 🔧 Running the Examples

### Prerequisites
```bash
# Ensure you're in the lesson directory
cd examples/fundamental_un_edifact/lesson2
```

### Basic Examples
```bash
# Run the main lesson
go run main.go
```

### What You'll See
The examples demonstrate:
- Segment structure analysis
- Element parsing and validation
- Composite element handling
- Position-based element access
- Segment type identification

## 💡 Key Concepts Explained

### 1. Element Positioning
Elements have specific positions within segments:
- Position 1: Segment tag (always present)
- Position 2+: Data elements (may be empty)

### 2. Mandatory vs Conditional Elements
- **Mandatory**: Must be present for valid message
- **Conditional**: May be omitted based on business rules

### 3. Element Types
- **Simple**: Single value
- **Composite**: Multiple components
- **Empty**: No value (represented by consecutive delimiters)

## 🧪 Practice Exercises

### Exercise 1: Identify Element Types
Analyze this segment and identify element types:
```
NAD+BY+++ACME CORP+123 MAIN ST+CITY+ST+12345+US'
```

**Answer**:
- Element 1: `BY` (simple - party qualifier)
- Element 2: (empty)
- Element 3: (empty)
- Element 4: `ACME CORP` (simple - name)
- Element 5: `123 MAIN ST` (simple - address)
- Element 6: `CITY` (simple - city)
- Element 7: `ST` (simple - state)
- Element 8: `12345` (simple - postal code)
- Element 9: `US` (simple - country code)

### Exercise 2: Count Elements
How many data elements are in this segment?
```
LIN+1++1234567890123:EN'
```

**Answer**: 3 data elements
- Element 1: `1` (line item number)
- Element 2: (empty)
- Element 3: `1234567890123:EN` (composite - item identification)

### Exercise 3: Parse Composite Element
Break down this composite element:
```
DTM+137:20231201:102'
```

**Answer**:
- Component 1: `137` (qualifier - document/message date/time)
- Component 2: `20231201` (date value)
- Component 3: `102` (format code - CCYYMMDD)

## ⚠️ Common Mistakes

1. **Ignoring Empty Elements**: Empty elements must be accounted for
2. **Wrong Element Count**: Not counting the segment tag as position 1
3. **Composite Confusion**: Treating composite elements as simple
4. **Position Errors**: Accessing elements by wrong position

## 🔍 Troubleshooting

### Element Access Issues
- Verify element positions (start from 1, not 0)
- Account for empty elements
- Check composite element structure

### Segment Validation
- Ensure segment tag is valid
- Verify required elements are present
- Check element data types

## 📚 Next Steps

After completing this lesson:
1. Practice with different segment types
2. Experiment with element positioning
3. Try creating your own segments
4. Move to Lesson 3: Composite Elements

## 🎯 Key Takeaways

- ✅ Segments are the building blocks of EDIFACT messages
- ✅ Elements have specific positions and types
- ✅ Empty elements must be properly handled
- ✅ Composite elements contain multiple components
- ✅ Understanding segments is crucial for message processing

---

*Ready for the next lesson? Let's dive deeper into composite elements! 🚀* 