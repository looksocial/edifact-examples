# Bookings Example: IFTMBF Adapter

This example demonstrates how to use the `github.com/looksocial/edifact` package as a dependency and implement a custom adapter for the IFTMBF EDIFACT message type, converting it to a Booking model suitable for database storage.

## Flow Diagram

```mermaid
flowchart TD
    A["EDIFACT Message (IFTMBF)"] -->|Read/Parse| B["edifact.Reader / edifact.Converter"]
    B -->|Parse| C["Tokenizer & Parser"]
    C -->|Build| D["model.Message"]
    D -->|Dispatch| E["Dispatcher/Router"]
    E -->|IFTMBFAdapter (Custom Handler)| F["Booking Model (struct)"]
    F -->|Marshal| G["JSON Output / DB Save"]
    style A fill:#f9f,stroke:#333,stroke-width:2px
    style G fill:#bbf,stroke:#333,stroke-width:2px
```

## Usage

1. Run `go mod tidy` in this directory to fetch dependencies.
2. Run the example:

```sh
go run main.go
```

You should see the parsed Booking struct as JSON output.

## Custom Adapter

- See `edifact_adapter/iftmbf.go` for the custom adapter implementation.
- See `models/booking.go` for the Booking struct.

You can extend this pattern for other EDIFACT message types and models as needed. 