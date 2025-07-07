# Simple Adapter Example: Flat Elements Output

This example demonstrates how to use a simple adapter to convert any EDIFACT message to a flat, numbered list of elements.

## Flow Diagram

```mermaid
flowchart TD
    A["EDIFACT Message"] -->|Read/Parse| B["edifact.Reader / edifact.Converter"]
    B -->|Parse| C["Tokenizer & Parser"]
    C -->|Build| D["model.Message"]
    D -->|Dispatch| E["SimpleAdapter"]
    E -->|Flatten| F["[]string (Numbered Elements)"]
    F -->|Marshal| G["JSON Output / UI"]
    style A fill:#f9f,stroke:#333,stroke-width:2px
    style G fill:#bbf,stroke:#333,stroke-width:2px
```

## Usage

See `main.go` for runnable code and element access examples. 