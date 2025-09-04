# astra-model-go

Go types and JSON schemas for ASTRA (Act State Representation Architecture) conversations.

## Installation

```bash
go get github.com/pryszm/astra-model-go
```

## Usage

### Import Package

```go
import "github.com/pryszm/astra-model-go"
```

### Create Acts

```go
package main

import (
    "fmt"
    "time"
    
    astra "github.com/pryszm/astra-model-go"
)

func main() {
    // Create an Ask act
    ask := astra.Ask{
        Act: astra.Act{
            ID:        astra.GenerateActID(),
            Timestamp: time.Now(),
            Speaker:   "agent_123",
            Type:      astra.ActTypeAsk,
        },
        Field:    "delivery_address",
        Prompt:   "What is your delivery address?",
        Required: true,
    }

    // Create a Fact act
    fact := astra.Fact{
        Act: astra.Act{
            ID:        astra.GenerateActID(),
            Timestamp: time.Now(),
            Speaker:   "customer_456",
            Type:      astra.ActTypeFact,
        },
        Entity:    "order_789",
        Field:     "delivery_address",
        Value:     "123 Main St, Anytown, USA",
        Operation: astra.FieldOperationSet,
    }

    fmt.Printf("Ask: %+v\n", ask)
    fmt.Printf("Fact: %+v\n", fact)
}
```

### Type Guards

```go
import astra "github.com/pryszm/astra-model-go"

// Check act types
if astra.IsAsk(act) {
    askAct, _ := act.(astra.Ask)
    fmt.Printf("Field requested: %s\n", askAct.Field)
}

if astra.IsFact(act) {
    factAct, _ := act.(astra.Fact)
    fmt.Printf("Value set: %v\n", factAct.Value)
}
```

### JSON Schema Validation

```go
import (
    "encoding/json"
    
    astra "github.com/pryszm/astra-model-go"
)

// Validate acts against JSON Schema
actJSON, _ := json.Marshal(ask)
isValid := astra.ValidateAct(actJSON)
fmt.Printf("Act is valid: %t\n", isValid)
```

## Core Types

- **`Act`** - Base type for all conversational actions
- **`Ask`** - Request missing information
- **`Fact`** - Declare state information  
- **`Confirm`** - Verify understanding before commitment
- **`Commit`** - Execute business processes
- **`Error`** - Handle failures and exceptions
- **`Conversation`** - Container for acts and metadata
- **`Participant`** - Conversation participant
- **`Entity`** - Business entity reference
- **`Constraint`** - Validation constraint

## Type System

ASTRA Go implements a type-safe approach to conversational acts:

```go
type ActType string

const (
    ActTypeAsk     ActType = "ask"
    ActTypeFact    ActType = "fact"
    ActTypeConfirm ActType = "confirm"
    ActTypeCommit  ActType = "commit"
    ActTypeError   ActType = "error"
)
```

## Utilities

### ID Generation

```go
// Generate ASTRA-compliant IDs
actID := astra.GenerateActID()           // "act_1a2b3c4d5e"
convID := astra.GenerateConversationID() // "conv_1a2b3c4d5e"
```

### Act Creation

```go
// Create base act with required fields
baseAct := astra.CreateBaseAct("speaker_123", astra.ActTypeAsk)
```

### Validation

```go
// Validate individual fields
isValidID := astra.IsValidActID("act_123")
isValidTimestamp := astra.IsValidTimestamp(time.Now())

// Type guards for runtime checking
if astra.IsParticipant(obj) {
    participant := obj.(astra.Participant)
    fmt.Printf("Participant: %s\n", participant.Name)
}
```

## Schema Evolution

This package follows semantic versioning for schema compatibility:

- **Patch versions (1.0.x)** - Bug fixes, no breaking changes
- **Minor versions (1.x.0)** - New optional fields, backward compatible
- **Major versions (x.0.0)** - Breaking changes, migration required

## Performance

The Go implementation is optimized for:

- **Zero-allocation JSON marshaling** where possible
- **Fast type assertions** for act type checking
- **Efficient ID generation** using optimized algorithms
- **Minimal memory footprint** for high-throughput scenarios

## Testing

Run the test suite:

```bash
go test ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## License

Licensed under the Apache License 2.0. See the [main repository](../../LICENSE) for details.

## Contributing

This package is part of the [ASTRA project](https://github.com/pryszm/astra). Please see the main repository for contribution guidelines.

## Links

- [ASTRA Repository](https://github.com/pryszm/astra)
- [Documentation](https://docs.pryszm.com)
- [Pryszm Platform](https://pryszm.com)
- [Go Package Documentation](https://pkg.go.dev/github.com/pryszm/astra-model-go)
