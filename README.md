# ASTRA — Act State Representation Architecture

Canonical types for Acts, Entities, Facts, Confirms, Commits, Errors with constraints, typing & extensibility.

## What is ASTRA?

ASTRA provides shared structure so different conversational runtimes and applications can interoperate. It defines the canonical data model for representing conversational state as typed, auditable sequences of actions.

## Core Types

**`Act`** - Base type for all conversational actions
```typescript
interface Act {
  id: string;
  timestamp: ISO8601;
  speaker: ParticipantId;
  type: "ask" | "fact" | "confirm" | "commit" | "error";
}
```

**`Ask`** - Request missing information
```typescript
interface Ask extends Act {
  type: "ask";
  field: string;
  prompt: string;
  constraints?: Constraint[];
}
```

**`Fact`** - Declare state information
```typescript
interface Fact extends Act {
  type: "fact";
  entity: EntityId;
  field: string;
  value: any;
  confidence?: number;
}
```

**`Confirm`** - Verify understanding before commitment
```typescript
interface Confirm extends Act {
  type: "confirm";
  entity: EntityId;
  summary: string;
  awaiting?: boolean;
  confirmed?: boolean;
}
```

**`Commit`** - Execute business processes
```typescript
interface Commit extends Act {
  type: "commit";
  entity: EntityId;
  action: "create" | "update" | "delete" | "execute";
  system?: string;
  transaction_id?: string;
}
```

**`Error`** - Handle failures and exceptions
```typescript
interface Error extends Act {
  type: "error";
  code: string;
  message: string;
  recoverable: boolean;
}
```

## Repository Structure

```
astra/
├── idl/                  # Interface Definition Languages
│   ├── json-schema/      # JSON Schema definitions
│   ├── protobuf/         # Protocol Buffer schemas  
│   └── avro/             # Apache Avro schemas
├── model/                # Language implementations
│   ├── typescript/       # astra-model-ts package
│   ├── python/           # astra-model-py package
│   └── go/               # astra-model-go package
├── docs/                 # Architecture documentation
├── examples/             # Reference implementations
├── tools/                # Validation and migration tools
└── compatibility/        # Version compatibility matrix
```

## Success Criteria

- **Schema-first evolution** - Forward/backward compatible versioning
- **Extensible type system** - Custom Acts & entity kinds
- **Static validation** - Compile-time type safety
- **Runtime guards** - Runtime validation and constraints

## Quick Start

### Install Model Library
```bash
npm install @astra/model-ts
pip install astra-model-py
go get github.com/pryszm/astra-model-go
```

### Validate Acts
```typescript
import { validateAct, Ask } from '@astra/model-ts';

const ask: Ask = {
  id: 'act_001',
  timestamp: '2025-01-15T14:30:00Z',
  speaker: 'agent_123',
  type: 'ask',
  field: 'delivery_address',
  prompt: 'What is your delivery address?'
};

const result = validateAct(ask); // ✅ Valid
```

### Extend Types
```typescript
interface CustomEntityAct extends Act {
  type: 'custom_entity';
  entity_type: 'my_business_object';
  custom_field: string;
}

// Register custom type
registerActType('custom_entity', CustomEntityActSchema);
```

## Deliverables

- **`astra/idl/`** - Multi-format schemas (JSON Schema, Protobuf, Avro)
- **`astra-model` libraries** - TypeScript, Python implementations with validators & builders  
- **Compatibility matrix** - Version compatibility and migration playbook
- **Migration tools** - Automated schema evolution utilities

## Schema Evolution

ASTRA supports forward and backward compatible evolution:

```yaml
# v1.0 -> v1.1 (backward compatible)
- Add optional fields
- Extend enums  
- Add new act types

# v1.x -> v2.0 (breaking changes)  
- Remove required fields
- Change field types
- Remove act types
```

See [compatibility matrix](./compatibility/matrix.md) for detailed evolution rules.
