# Act Specification

An open standard for representing conversational state as structured, auditable sequences of actions.

## What is an Act?

Acts are atomic conversational events that capture the semantic meaning of business conversations. Instead of storing conversations as unstructured text, acts represent the implicit state changes that occur when people communicate.

### Four Act Types

**`Ask`** - Request missing information
```json
{
  "type": "ask",
  "field": "delivery_address", 
  "prompt": "What's your delivery address?"
}
```

**`State`** - Declare facts or information
```json
{
  "type": "state",
  "entity": "order",
  "field": "items",
  "value": [{"type": "pizza", "size": "large"}]
}
```

**`Confirm`** - Verify understanding before commitment
```json
{
  "type": "confirm",
  "entity": "order",
  "summary": "Two large pizzas for delivery to 123 Main St"
}
```

**`Commit`** - Execute business processes
```json
{
  "type": "commit",
  "entity": "order",
  "action": "create",
  "transaction_id": "ord_1234567890"
}
```

## Repository Contents

- **`/spec`** - Technical specification documents
- **`/schema`** - JSON Schema for validation  
- **`/examples`** - Reference implementations and sample conversations
- **`/tools`** - Validation and development utilities

## Quick Start

1. Read the [full specification](./spec/act-spec.md)
2. Validate acts using the [JSON Schema](./schema/act.json)
3. See [examples](./examples/) for common conversation patterns
4. Use [validation tools](./tools/) to test your implementations

## Validation

```bash
npm install -g @pryszm/act-validator
act-validate my-conversation.json
```
