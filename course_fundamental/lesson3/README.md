# Lesson 3: Composite Elements

## 📚 Learning Objectives

By the end of this lesson, you will understand:
- ✅ Composite data element structure and purpose
- ✅ Component separation and positioning
- ✅ Qualifiers and their role in composite elements
- ✅ Complex composite element scenarios

## 🔍 What are Composite Elements?

Composite elements are data elements that contain multiple related pieces of information within a single element. They use the component separator (`:`) to divide the different components.

### Composite Element Structure

```
DTM+137:20231201:102'
     ^   ^         ^
     |   |         └── Component 3 (Format code)
     |   └── Component 2 (Date value)
     └── Component 1 (Qualifier)
```

## 🎯 Composite Element Types

### 1. Qualifier + Value + Format
Most common pattern:
- Component 1: Qualifier (defines the type)
- Component 2: Value (the actual data)
- Component 3: Format code (how to interpret the value)

### 2. Identification + Qualifier
For identifying objects:
- Component 1: Identification code
- Component 2: Code list qualifier

### 3. Multiple Values
For complex data:
- Multiple components with different meanings
- Each component serves a specific purpose

## 📖 Common Composite Elements

### Date/Time Elements
```
DTM+137:20231201:102'    # Document date
DTM+2:20231201:102'      # Delivery date
DTM+35:1430:201'         # Time (14:30)
```

### Identification Elements
```
LIN+1++1234567890123:EN' # Item identification
NAD+BY+++ACME CORP'      # Party identification
RFF+CT:123456'           # Reference number
```

### Quantity Elements
```
QTY+12:100:PCE'          # Quantity (100 pieces)
QTY+145:25.5:KGM'        # Weight (25.5 kg)
QTY+146:10:LTR'          # Volume (10 liters)
```

## 🔧 Running the Examples

### Prerequisites
```bash
# Ensure you're in the lesson directory
cd examples/fundamental_un_edifact/lesson3
```

### Basic Examples
```bash
# Run the main lesson
go run main.go
```

### What You'll See
The examples demonstrate:
- Composite element parsing
- Component extraction and validation
- Qualifier interpretation
- Complex composite scenarios
- Error handling for malformed composites

## 💡 Key Concepts Explained

### 1. Component Positioning
Components have specific positions within composite elements:
- Position 1: Usually a qualifier or identifier
- Position 2+: Values, codes, or additional qualifiers

### 2. Qualifiers
Qualifiers define the meaning of the composite element:
- `137` - Document/message date/time
- `2` - Delivery date/time
- `12` - Number of packages
- `AAA` - Free text qualifier

### 3. Format Codes
Format codes specify how to interpret values:
- `102` - CCYYMMDD (Century, Year, Month, Day)
- `201` - HHMM (Hour, Minute)
- `EN` - EAN (European Article Number)

## 🧪 Practice Exercises

### Exercise 1: Parse Date Element
Break down this composite element:
```
DTM+137:20231201:102'
```

**Answer**:
- Component 1: `137` (qualifier - document/message date/time)
- Component 2: `20231201` (date value - December 1, 2023)
- Component 3: `102` (format - CCYYMMDD)

### Exercise 2: Identify Components
What are the components in this element?
```
QTY+12:100:PCE'
```

**Answer**:
- Component 1: `12` (qualifier - number of packages)
- Component 2: `100` (quantity value)
- Component 3: `PCE` (unit qualifier - pieces)

### Exercise 3: Complex Composite
Analyze this complex composite element:
```
NAD+BY+++ACME CORP+123 MAIN ST+CITY+ST+12345+US'
```

**Answer**: This is actually multiple simple elements, not a composite:
- Element 1: `BY` (party qualifier)
- Element 2: (empty)
- Element 3: (empty)
- Element 4: `ACME CORP` (name)
- Element 5: `123 MAIN ST` (address)
- Element 6: `CITY` (city)
- Element 7: `ST` (state)
- Element 8: `12345` (postal code)
- Element 9: `US` (country code)

## ⚠️ Common Mistakes

1. **Confusing Elements and Components**: Not distinguishing between segment elements and composite components
2. **Wrong Separator**: Using `+` instead of `:` for component separation
3. **Missing Components**: Not accounting for empty components
4. **Qualifier Confusion**: Not understanding what qualifiers mean

## 🔍 Troubleshooting

### Composite Parsing Issues
- Verify component separator usage (`:`)
- Check for empty components (consecutive `:`)
- Validate qualifier meanings
- Ensure proper component count

### Validation Problems
- Component count validation
- Qualifier code validation
- Format code validation
- Value format validation

## 📚 Next Steps

After completing this lesson:
1. Practice parsing different composite elements
2. Learn common qualifier codes
3. Understand format codes
4. Move to Lesson 4: Message Structure

## 🎯 Key Takeaways

- ✅ Composite elements contain multiple related components
- ✅ Components are separated by the `:` delimiter
- ✅ Qualifiers define the meaning of composite elements
- ✅ Format codes specify how to interpret values
- ✅ Understanding composites is essential for data extraction

---

*Ready for the next lesson? Let's explore message structure and hierarchy! 🚀* 