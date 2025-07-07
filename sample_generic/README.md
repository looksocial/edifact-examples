# Generic Example: Custom Adapter Flow

This example demonstrates how to implement and register a custom adapter for a specific EDIFACT message type (e.g., ORDERS, VERMAS, IFTMBF).

## Flow Diagram

```mermaid
flowchart TD
    A["EDIFACT Message (ORDERS, VERMAS, IFTMBF)"] -->|Read/Parse| B["edifact.Reader / edifact.Converter"]
    B -->|Parse| C["Tokenizer & Parser"]
    C -->|Build| D["model.Message"]
    D -->|Dispatch| E["Dispatcher/Router"]
    E -->|CustomAdapter (ORDERS, VERMAS, IFTMBF)| F["Custom Struct (Booking, etc.)"]
    F -->|Marshal| G["JSON Output / DB Save"]
    style A fill:#f9f,stroke:#333,stroke-width:2px
    style G fill:#bbf,stroke:#333,stroke-width:2px
```

## Usage

See `main.go` for runnable code and custom adapter examples. 