# Basic Example: EDIFACT to JSON

This example demonstrates how to use the `edifact` package to convert an EDIFACT message to JSON using the generic handler.

## Flow Diagram

```mermaid
flowchart TD
    A["EDIFACT Message (string)"] -->|Read/Parse| B["edifact.Reader / edifact.Converter"]
    B -->|Parse| C["Tokenizer & Parser"]
    C -->|Build| D["model.Message"]
    D -->|Dispatch| E["Dispatcher/Router"]
    E -->|GenericHandler| F["Generic Structure (map/struct)"]
    F -->|Marshal| G["JSON Output"]
    style A fill:#f9f,stroke:#333,stroke-width:2px
    style G fill:#bbf,stroke:#333,stroke-width:2px
```

## Usage

See `main.go` for runnable code. 