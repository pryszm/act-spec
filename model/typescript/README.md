# @astra/model-ts

TypeScript types and JSON schemas for ASTRA (Act State Representation Architecture) conversations.

## Installation

```bash
npm install @astra/model-ts
```

## Usage

### Import Types

```typescript
import { 
  Act, 
  Ask, 
  Fact, 
  Confirm, 
  Commit, 
  Error,
  Conversation,
  Participant,
  Entity,
  Constraint 
} from '@astra/model-ts';
```

### Create Acts

```typescript
const ask: Ask = {
  id: 'act_001',
  timestamp: '2025-01-15T14:30:00Z',
  speaker: 'agent_123',
  type: 'ask',
  field: 'delivery_address',
  prompt: 'What is your delivery address?',
  required: true
};

const fact: Fact = {
  id: 'act_002',
  timestamp: '2025-01-15T14:31:00Z',
  speaker: 'customer_456',
  type: 'fact',
  entity: 'order_789',
  field: 'delivery_address',
  value: '123 Main St, Anytown, USA',
  operation: 'set'
};
```

### Access JSON Schemas

```typescript
import { schemas } from '@astra/model-ts';

// Validate acts against JSON Schema
const isValid = validate(schemas.act, someActData);
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

## Schema Evolution

This package follows semantic versioning for schema compatibility:

- **Patch versions (1.0.x)** - Bug fixes, no breaking changes
- **Minor versions (1.x.0)** - New optional fields, backward compatible
- **Major versions (x.0.0)** - Breaking changes, migration required

## License

Licensed under the Apache License 2.0. See the [main repository](../../LICENSE) for details.

## Contributing

This package is part of the [ASTRA project](https://github.com/pryszm/astra). Please see the main repository for contribution guidelines.

## Links

- [ASTRA Repository](https://github.com/pryszm/astra)
- [Documentation](https://docs.pryszm.com)
- [Pryszm Platform](https://pryszm.com)
