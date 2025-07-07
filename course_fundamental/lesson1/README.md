# Lesson 1: Basic UN/EDIFACT Syntax & Delimiters

## 📚 Learning Objectives

By the end of this lesson, you will understand:
- ✅ UN/EDIFACT syntax rules and structure
- ✅ The role and importance of delimiters
- ✅ Character encoding in EDIFACT
- ✅ Basic syntax validation concepts

## 🔍 What is UN/EDIFACT Syntax?

UN/EDIFACT syntax is a set of strict rules that govern how electronic data interchange messages are structured and formatted. It ensures that messages can be reliably transmitted, received, and processed by different systems worldwide.

### Key Syntax Principles

1. **Hierarchical Structure**: Messages are organized in a specific hierarchy
2. **Delimiter-Based Separation**: Special characters separate different data components
3. **Positional Elements**: Data elements have specific positions within segments
4. **Standardized Format**: All messages follow the same basic structure

## 🎯 EDIFACT Delimiters

Delimiters are special characters that separate different parts of an EDIFACT message. They are crucial for parsing and understanding the message structure.

### Standard EDIFACT Delimiters

| Delimiter | Name | Purpose | Example |
|-----------|------|---------|---------|
| `'` | Segment Terminator | Ends each segment | `UNH+1+INVOIC:D:97A:UN'` |
| `+` | Data Element Separator | Separates data elements within a segment | `UNH+1+INVOIC:D:97A:UN'` |
| `:` | Component Data Element Separator | Separates components within composite elements | `DTM+137:20231201:102'` |
| `?` | Release Character | Escapes special characters in data | `FTX+AAA++?+This is a ?+plus sign'` |

### Delimiter Examples

```edifact
# Simple segment with data elements
UNH+1+INVOIC:D:97A:UN'

# Segment with composite element
DTM+137:20231201:102'

# Segment with released character
FTX+AAA++?+This contains a ?+plus sign'
```

## 📖 Syntax Rules

### 1. Segment Structure
- Each segment starts with a 3-character segment tag
- Data elements are separated by `+`
- Segment ends with segment terminator `'`
- Empty elements are represented by consecutive delimiters

### 2. Data Element Rules
- Simple elements contain single values
- Composite elements contain multiple components separated by `:`
- Empty elements are allowed but must be properly delimited

### 3. Character Encoding
- EDIFACT uses specific character sets (UNOA, UNOB, UNOC, etc.)
- Special characters must be escaped using the release character `?`
- Line breaks and formatting are not allowed in data

## 🔧 Running the Examples

### Prerequisites
```bash
# Ensure you're in the lesson directory
cd examples/fundamental_un_edifact/lesson1
```

### Basic Examples
```bash
# Run the main lesson
go run main.go
```

### What You'll See
The examples demonstrate:
- Basic delimiter usage
- Segment structure
- Composite element handling
- Character escaping
- Syntax validation

## 💡 Key Concepts Explained

### 1. Segment Tags
Every EDIFACT segment begins with a 3-character tag that identifies the segment type:
- `UNH` - Message Header
- `UNT` - Message Trailer
- `DTM` - Date/Time/Period
- `FTX` - Free Text

### 2. Data Element Positioning
Data elements have specific positions within segments:
- Position 1: Segment tag (always present)
- Position 2+: Data elements (may be empty)

### 3. Composite Elements
Composite elements contain multiple related pieces of information:
```
DTM+137:20231201:102'
     ^   ^         ^
     |   |         └── Format code (102 = CCYYMMDD)
     |   └── Date value
     └── Qualifier (137 = Document/message date/time)
```

## 🧪 Practice Exercises

### Exercise 1: Identify Delimiters
Look at this segment and identify each delimiter:
```
UNH+1+INVOIC:D:97A:UN'
```

**Answer**: 
- `+` separates data elements
- `:` separates components in the composite element
- `'` terminates the segment

### Exercise 2: Count Elements
How many data elements are in this segment?
```
DTM+137:20231201:102'
```

**Answer**: 2 data elements
- Element 1: `137:20231201:102` (composite element with 3 components)
- Element 2: (empty)

### Exercise 3: Escape Characters
What does this segment contain?
```
FTX+AAA++?+This contains a ?+plus sign'
```

**Answer**: The text "This contains a +plus sign" (the `?+` is escaped to represent a literal `+`)

## ⚠️ Common Mistakes

1. **Missing Delimiters**: Forgetting to include required delimiters
2. **Incorrect Escaping**: Not properly escaping special characters
3. **Empty Element Handling**: Not accounting for empty elements
4. **Character Encoding**: Using unsupported characters

## 🔍 Troubleshooting

### Syntax Errors
- Check delimiter placement
- Verify segment termination
- Ensure proper character escaping
- Validate element positioning

### Common Issues
- **"Invalid segment"**: Check segment tag and structure
- **"Missing delimiter"**: Verify all required delimiters are present
- **"Invalid character"**: Check character encoding and escaping

## 📚 Next Steps

After completing this lesson:
1. Practice with the provided examples
2. Try modifying the delimiter values
3. Experiment with different segment structures
4. Move to Lesson 2: Segments & Elements

## 🎯 Key Takeaways

- ✅ EDIFACT uses specific delimiters to separate data components
- ✅ Delimiters are crucial for parsing and validation
- ✅ Character escaping is essential for special characters
- ✅ Syntax rules must be followed precisely
- ✅ Understanding delimiters is fundamental to EDIFACT

---

*Ready for the next lesson? Let's explore segments and elements in detail! 🚀* 