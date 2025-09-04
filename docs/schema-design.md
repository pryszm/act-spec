# ASTRA Schema Design Guide

This document outlines the principles, patterns, and best practices for designing robust, extensible, and interoperable ASTRA schemas. Whether you're extending ASTRA for custom business domains or contributing to the core specification, these guidelines ensure consistency and long-term maintainability.

## Design Philosophy

### Schema-First Development

ASTRA follows a schema-first approach where all types are formally defined before implementation. This ensures:

- **Contract-driven development** - APIs and integrations are built against stable contracts
- **Multi-language consistency** - All language implementations derive from the same source schemas
- **Automated validation** - Runtime validation is generated from schema definitions
- **Documentation generation** - Schema annotations become documentation automatically

### Interoperability Through Standards

ASTRA uses multiple Interface Definition Languages (IDLs) to ensure broad compatibility:

- **JSON Schema** - Web APIs, validation, and documentation
- **Protocol Buffers** - High-performance serialization and gRPC services
- **Apache Avro** - Schema evolution and data analytics pipelines

Each IDL serves specific use cases while maintaining semantic consistency across all representations.

### Backward Compatibility by Design

ASTRA schemas are designed for evolution without breaking existing implementations:

- **Additive changes only** within major versions
- **Optional fields** for extensions
- **Enum expansion** for new categories
- **Deprecation pathways** for obsolete fields

## Core Schema Principles

### 1. Semantic Clarity

Every field and type should have clear, unambiguous semantics:

```typescript
// Good: Clear semantic meaning
interface Ask extends Act {
  field: string;        // What information is being requested
  prompt: string;       // Human-readable question
  required?: boolean;   // Whether information is mandatory
}

// Avoid: Ambiguous or overly generic names
interface Ask extends Act {
  data: string;         // What data? Too generic
  text: string;         // What kind of text? Unclear
  flag?: boolean;       // What does this flag control?
}
```

### 2. Consistent Naming Conventions

ASTRA follows consistent naming across all schemas:

**Field Names:**
- `snake_case` for JSON Schema and Avro
- `camelCase` for TypeScript interfaces
- `UpperCamelCase` for Protocol Buffer messages

**Type Names:**
- `PascalCase` for all type names across all IDLs
- Descriptive names that indicate purpose: `ValidationStatus`, `CommitAction`

**Enum Values:**
- `UPPER_CASE` with underscores for Protocol Buffers
- `lowercase` with underscores for JSON Schema enums
- Consistent semantic mapping across IDLs

### 3. Extensibility Patterns

ASTRA schemas support extension through several mechanisms:

**Metadata Fields:**
```typescript
interface Act {
  id: string;
  timestamp: string;
  speaker: string;
  type: ActType;
  metadata?: ActMetadata;  // Extension point
}

interface ActMetadata {
  channel?: string;
  language?: string;
  [key: string]: any;      // Additional custom properties
}
```

**Union Types:**
```typescript
type EntityRef = string | Entity;  // Simple ID or full object

type ConversationAct = Ask | Fact | Confirm | Commit | Error;  // Extensible union
```

**Custom Constraints:**
```typescript
interface Constraint {
  type: ConstraintType;
  value?: any;              // Flexible value type
  message?: string;
  code?: string;
}

// Enables domain-specific validation
{
  type: 'custom',
  value: { business_rule: 'credit_check', min_score: 650 },
  message: 'Credit check required'
}
```

### 4. Validation-Friendly Design

Schemas should enable both structural and semantic validation:

**Required vs Optional Fields:**
```typescript
interface Ask extends Act {
  // Required fields establish contract
  field: string;            // Required: Must specify what's being requested
  prompt: string;           // Required: Must provide user-facing text
  
  // Optional fields enable flexibility
  constraints?: Constraint[]; // Optional: Not all asks need validation
  required?: boolean;        // Optional: Defaults to true
}
```

**Constrained Types:**
```typescript
interface Act {
  id: string;               // Pattern: ^act_[a-zA-Z0-9_-]+$
  timestamp: string;        // Format: ISO 8601 date-time
  confidence?: number;      // Range: 0.0 to 1.0
}
```

**Enum Constraints:**
```typescript
type ActType = 'ask' | 'fact' | 'confirm' | 'commit' | 'error';
type Source = 'human' | 'ai' | 'system' | 'speech_recognition' | 'text_analysis';
```

## Multi-IDL Schema Design

### JSON Schema Design

JSON Schema serves as the canonical definition for ASTRA types:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://schemas.astra.dev/v1/ask.json",
  "title": "Ask",
  "description": "Act that requests missing information",
  "allOf": [
    { "$ref": "act.json" },
    {
      "type": "object",
      "properties": {
        "type": { "const": "ask" },
        "field": {
          "type": "string",
          "description": "Field or information being requested"
        },
        "prompt": {
          "type": "string", 
          "description": "Question or request presented to user"
        },
        "constraints": {
          "type": "array",
          "items": { "$ref": "constraint.json" },
          "description": "Validation constraints"
        }
      },
      "required": ["field", "prompt"],
      "additionalProperties": false
    }
  ]
}
```

**JSON Schema Best Practices:**

- Use `$ref` for type composition and reuse
- Include comprehensive `description` fields
- Use `allOf` for inheritance patterns
- Set `additionalProperties: false` for strict validation
- Use `const` for discriminator fields
- Include format constraints where appropriate

### Protocol Buffer Design

Protocol Buffers provide efficient serialization and strong typing:

```protobuf
syntax = "proto3";

package astra.v1;

import "act.proto";
import "constraint.proto";

// Act that requests missing information
message Ask {
  // Base act properties
  Act act = 1;
  
  // Field or information being requested
  string field = 2;
  
  // Question or request presented to user
  string prompt = 3;
  
  // Validation constraints for expected responses
  repeated Constraint constraints = 4;
  
  // Whether this information is required to proceed
  optional bool required = 5;
  
  // Expected data type of response
  optional ExpectedType expected_type = 6;
}

// Expected response data types
enum ExpectedType {
  EXPECTED_TYPE_UNSPECIFIED = 0;
  EXPECTED_TYPE_STRING = 1;
  EXPECTED_TYPE_NUMBER = 2;
  EXPECTED_TYPE_BOOLEAN = 3;
  EXPECTED_TYPE_OBJECT = 4;
  EXPECTED_TYPE_ARRAY = 5;
  EXPECTED_TYPE_DATE = 6;
  EXPECTED_TYPE_EMAIL = 7;
  EXPECTED_TYPE_PHONE = 8;
  EXPECTED_TYPE_ADDRESS = 9;
}
```

**Protocol Buffer Best Practices:**

- Use `optional` for nullable fields
- Reserve field numbers for future use
- Use descriptive enum prefixes
- Include comprehensive comments
- Use `repeated` for arrays
- Embed related messages for composition

### Apache Avro Design

Avro provides schema evolution capabilities for data processing:

```json
{
  "type": "record",
  "name": "Ask",
  "namespace": "dev.astra.v1",
  "doc": "Act that requests missing information",
  "fields": [
    {
      "name": "act",
      "type": "dev.astra.v1.Act",
      "doc": "Base act properties"
    },
    {
      "name": "field", 
      "type": "string",
      "doc": "Field or information being requested"
    },
    {
      "name": "prompt",
      "type": "string", 
      "doc": "Question or request presented to user"
    },
    {
      "name": "constraints",
      "type": {
        "type": "array",
        "items": "dev.astra.v1.Constraint"
      },
      "default": [],
      "doc": "Validation constraints for expected responses"
    },
    {
      "name": "required",
      "type": ["null", "boolean"],
      "default": null,
      "doc": "Whether this information is required to proceed"
    }
  ]
}
```

**Avro Best Practices:**

- Use `default` values for optional fields
- Use union types (`["null", "type"]`) for nullable fields
- Include comprehensive `doc` fields
- Use namespaces for type organization
- Design for schema evolution from the start

## Extension Patterns

### Custom Act Types

ASTRA supports custom act types for domain-specific requirements:

```typescript
// Define custom act type
interface ScheduleAct extends Act {
  type: 'schedule';
  appointment_id: string;
  datetime: string;
  duration_minutes: number;
  participants: string[];
  location?: string;
}

// Register in type system
type ExtendedConversationAct = ConversationAct | ScheduleAct;
```

**Custom Act Guidelines:**

- Extend the base `Act` interface
- Use descriptive type identifiers
- Follow existing field naming conventions
- Include comprehensive documentation
- Provide JSON Schema definition
- Consider backward compatibility

### Entity Schema Extensions

Entities can be extended with domain-specific fields:

```typescript
// Base entity interface
interface Entity {
  id: string;
  type: string;
  external_id?: string;
  system?: string;
  metadata?: Record<string, any>;
}

// Domain-specific entity
interface OrderEntity extends Entity {
  type: 'order';
  customer_id: string;
  items: OrderItem[];
  total_amount: number;
  currency: string;
  status: OrderStatus;
}

// Use in acts
interface OrderFact extends Fact {
  entity: OrderEntity;
  // ... additional order-specific fields
}
```

### Constraint System Extensions

Custom constraints enable domain-specific validation:

```typescript
// Custom constraint types
type CustomConstraintType = 
  | 'credit_check'
  | 'inventory_available'  
  | 'business_hours'
  | 'geolocation_valid';

interface CustomConstraint extends Constraint {
  type: CustomConstraintType;
  value: {
    business_rule: string;
    parameters: Record<string, any>;
  };
}

// Usage in Ask acts
const askCreditCard: Ask = {
  // ... base fields
  constraints: [
    {
      type: 'custom',
      value: {
        business_rule: 'credit_limit_check',
        parameters: { 
          amount: 5000,
          customer_tier: 'premium'
        }
      },
      message: 'Credit limit verification required'
    }
  ]
};
```

## Schema Evolution Patterns

### Backward Compatible Changes

These changes can be made within the same major version:

**Adding Optional Fields:**
```typescript
// v1.0
interface Ask extends Act {
  field: string;
  prompt: string;
}

// v1.1 - Backward compatible
interface Ask extends Act {
  field: string;
  prompt: string;
  expected_type?: ExpectedType;  // New optional field
  retry_count?: number;          // New optional field
}
```

**Extending Enums:**
```typescript
// v1.0
type ActType = 'ask' | 'fact' | 'confirm' | 'commit' | 'error';

// v1.1 - Backward compatible  
type ActType = 'ask' | 'fact' | 'confirm' | 'commit' | 'error' | 'schedule';
```

**Adding Union Type Members:**
```typescript
// v1.0
type ConversationAct = Ask | Fact | Confirm | Commit | Error;

// v1.1 - Backward compatible
type ConversationAct = Ask | Fact | Confirm | Commit | Error | Schedule;
```

### Breaking Changes (Major Version)

These changes require a new major version:

**Removing Required Fields:**
```typescript
// v1.x - BREAKING: Would break existing implementations
interface Ask extends Act {
  // field: string;  // Removed required field
  prompt: string;
}
```

**Changing Field Types:**
```typescript  
// v1.x - BREAKING: Incompatible type change
interface Ask extends Act {
  field: string[];   // Changed from string to string[]
  prompt: string;
}
```

**Removing Enum Values:**
```typescript
// v1.x - BREAKING: Existing data would be invalid
type ActType = 'ask' | 'fact' | 'confirm';  // Removed 'commit' | 'error'
```

### Migration Strategies

**Field Deprecation:**
```typescript
interface Ask extends Act {
  field: string;
  prompt: string;
  
  /** @deprecated Use expected_type instead */
  response_type?: string;
  
  expected_type?: ExpectedType;  // Replacement field
}
```

**Gradual Type Migration:**
```typescript
// Support both old and new formats during transition
type EntityRef = string | Entity | LegacyEntity;

// Provide conversion utilities
function migrateEntity(legacy: LegacyEntity): Entity {
  return {
    id: legacy.entityId,
    type: legacy.entityType,
    external_id: legacy.externalRef,
    // ... field mapping
  };
}
```

**Schema Versioning:**
```typescript
// Include version in schema URLs
"$id": "https://schemas.astra.dev/v1/ask.json"
"$id": "https://schemas.astra.dev/v2/ask.json"

// Version-aware type definitions
interface ActV1 { /* v1 definition */ }
interface ActV2 { /* v2 definition */ }

type Act = ActV1 | ActV2;  // Support multiple versions
```

## Validation Architecture

### Multi-Layer Validation

ASTRA implements validation at multiple levels:

**1. Compile-Time Validation (TypeScript)**
```typescript
// Type system catches structural errors
const invalidAsk: Ask = {
  id: 'act_001',
  timestamp: '2025-01-15T14:30:00Z',
  speaker: 'agent',
  type: 'ask',
  // Missing required 'field' and 'prompt' - compile error
};
```

**2. Runtime Structural Validation (JSON Schema)**
```typescript
import { validate } from 'ajv';
import { schemas } from '@astra/model-ts';

const result = validate(schemas.ask, actData);
if (!result.valid) {
  console.error('Schema validation errors:', result.errors);
}
```

**3. Business Rule Validation (Constraints)**
```typescript
const askWithValidation: Ask = {
  // ... base fields
  constraints: [
    {
      type: 'format',
      value: 'email',
      message: 'Must be valid email address'
    },
    {
      type: 'custom', 
      value: { business_rule: 'customer_exists' },
      message: 'Customer must exist in system'
    }
  ]
};
```

### Constraint Design Patterns

**Reusable Constraint Factories:**
```typescript
// Common constraint builders
export const Constraints = {
  required: (message?: string) => ({
    type: 'required',
    message: message || 'This field is required'
  }),
  
  email: (message?: string) => ({
    type: 'format',
    value: 'email',
    message: message || 'Must be valid email address'
  }),
  
  range: (min: number, max: number, message?: string) => ({
    type: 'range',
    value: { min, max, inclusive: true },
    message: message || `Value must be between ${min} and ${max}`
  }),
  
  businessRule: (rule: string, params: any, message: string) => ({
    type: 'custom',
    value: { business_rule: rule, parameters: params },
    message
  })
};

// Usage
const askCreditCard: Ask = {
  // ... base fields
  constraints: [
    Constraints.required(),
    Constraints.range(13, 19, 'Credit card must be 13-19 digits'),
    Constraints.businessRule('luhn_check', {}, 'Invalid credit card number')
  ]
};
```

**Hierarchical Validation:**
```typescript
// Entity-level constraints
interface EntityConstraints {
  [entityType: string]: {
    [field: string]: Constraint[];
  };
}

const orderConstraints: EntityConstraints = {
  order: {
    total_amount: [
      Constraints.required(),
      Constraints.range(0.01, 10000, 'Order total must be between $0.01 and $10,000')
    ],
    customer_email: [
      Constraints.required(),
      Constraints.email(),
      Constraints.businessRule('customer_exists', {}, 'Customer not found')
    ]
  }
};
```

## Performance Considerations

### Schema Size Optimization

**Minimize Schema Overhead:**
```typescript
// Prefer simple types over complex objects when possible
interface OptimizedAct {
  id: string;              // Simple string vs complex ID object
  timestamp: string;       // ISO string vs Date object  
  metadata?: ActMetadata;  // Optional heavy objects
}

// Use references instead of embedding
interface Fact extends Act {
  entity: string;          // Entity ID reference
  // vs entity: Entity     // Full entity object
}
```

**Efficient Serialization:**
```typescript
// Design for efficient JSON serialization
interface StreamlinedAct {
  // Required fields first
  id: string;
  type: ActType;
  timestamp: string;
  
  // Optional fields last
  confidence?: number;
  metadata?: Record<string, any>;
}
```

### Validation Performance

**Precompiled Validation:**
```typescript
// Compile schemas once, use many times
import Ajv from 'ajv';

const ajv = new Ajv();
const validateAsk = ajv.compile(schemas.ask);
const validateFact = ajv.compile(schemas.fact);

// Fast validation calls
const isValidAsk = validateAsk(actData);
```

**Caching Strategies:**
```typescript
// Cache validation results for repeated validations
const validationCache = new Map<string, boolean>();

function cachedValidate(act: Act): boolean {
  const key = `${act.type}_${JSON.stringify(act)}`;
  
  if (validationCache.has(key)) {
    return validationCache.get(key)!;
  }
  
  const isValid = validateAct(act);
  validationCache.set(key, isValid);
  return isValid;
}
```

## Testing Schema Design

### Schema Validation Testing

**Test Valid Cases:**
```typescript
describe('Ask Schema Validation', () => {
  it('validates minimal valid ask', () => {
    const ask: Ask = {
      id: 'act_001',
      timestamp: '2025-01-15T14:30:00Z',
      speaker: 'agent',
      type: 'ask',
      field: 'email',
      prompt: 'What is your email?'
    };
    
    expect(validateAsk(ask)).toBe(true);
  });
  
  it('validates ask with all optional fields', () => {
    const ask: Ask = {
      id: 'act_001',
      timestamp: '2025-01-15T14:30:00Z',
      speaker: 'agent',
      type: 'ask',
      field: 'email',
      prompt: 'What is your email?',
      constraints: [Constraints.email()],
      required: true,
      expected_type: 'email',
      retry_count: 0,
      max_retries: 3
    };
    
    expect(validateAsk(ask)).toBe(true);
  });
});
```

**Test Invalid Cases:**
```typescript
describe('Ask Schema Validation Errors', () => {
  it('rejects ask missing required fields', () => {
    const invalidAsk = {
      id: 'act_001',
      timestamp: '2025-01-15T14:30:00Z',
      speaker: 'agent',
      type: 'ask'
      // Missing 'field' and 'prompt'
    };
    
    expect(validateAsk(invalidAsk)).toBe(false);
  });
  
  it('rejects ask with invalid field types', () => {
    const invalidAsk = {
      id: 'act_001',
      timestamp: '2025-01-15T14:30:00Z',
      speaker: 'agent',
      type: 'ask',
      field: 123,          // Should be string
      prompt: 'What is your email?'
    };
    
    expect(validateAsk(invalidAsk)).toBe(false);
  });
});
```

### Constraint Testing

**Test Built-in Constraints:**
```typescript
describe('Constraint Validation', () => {
  it('validates email format constraint', () => {
    const constraint = Constraints.email();
    
    expect(validateConstraint(constraint, 'user@example.com')).toBe(true);
    expect(validateConstraint(constraint, 'invalid-email')).toBe(false);
  });
  
  it('validates range constraint', () => {
    const constraint = Constraints.range(1, 10);
    
    expect(validateConstraint(constraint, 5)).toBe(true);
    expect(validateConstraint(constraint, 0)).toBe(false);
    expect(validateConstraint(constraint, 11)).toBe(false);
  });
});
```

**Test Custom Constraints:**
```typescript
describe('Custom Business Rule Constraints', () => {
  it('validates credit check constraint', () => {
    const constraint = Constraints.businessRule(
      'credit_check',
      { min_score: 650 },
      'Credit score too low'
    );
    
    const validCustomer = { credit_score: 700 };
    const invalidCustomer = { credit_score: 600 };
    
    expect(validateConstraint(constraint, validCustomer)).toBe(true);
    expect(validateConstraint(constraint, invalidCustomer)).toBe(false);
  });
});
```

## Best Practices Summary

### Schema Definition
1. **Start with JSON Schema** as the canonical definition
2. **Use clear, descriptive names** for all types and fields
3. **Include comprehensive documentation** in schema annotations
4. **Design for extensibility** with metadata fields and union types
5. **Minimize required fields** to maximize flexibility

### Validation Strategy
1. **Layer validation** - structural, semantic, and business rules
2. **Fail fast** with clear error messages
3. **Cache compiled schemas** for performance
4. **Provide constraint factories** for common validation patterns
5. **Test both valid and invalid cases** comprehensively

### Evolution Planning
1. **Version schemas explicitly** with URLs and namespaces
2. **Only make additive changes** within major versions
3. **Deprecate before removal** with clear migration paths
4. **Provide migration utilities** for breaking changes
5. **Document compatibility** requirements clearly

### Performance Optimization
1. **Keep schemas lightweight** - prefer references over embedding
2. **Compile validation once** and reuse
3. **Design for efficient serialization** 
4. **Consider caching strategies** for repeated validation
5. **Profile schema validation performance** in production

By following these schema design principles and patterns, ASTRA implementations can maintain interoperability, extensibility, and performance while supporting diverse business requirements and evolving over time.
