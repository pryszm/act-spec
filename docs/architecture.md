# ASTRA Architecture

**Act State Representation Architecture**

ASTRA (Act State Representation Architecture) is a foundational specification for representing conversational state as typed, auditable sequences of actions. It provides the canonical data model that enables different conversational systems, runtimes, and applications to interoperate through shared structural semantics.

## Core Principles

### Conversational State as Data

ASTRA treats every meaningful business conversation as a sequence of state transformations over structured entities. Rather than capturing conversations as unstructured transcripts, ASTRA models them as append-only ledgers of typed actions that can be validated, queried, and integrated with business systems.

This paradigm shift moves us from "conversation as transcript" to "conversation as state transformation," enabling unprecedented reliability, integration capabilities, and continuity across communication channels.

### Schema-First Evolution

All ASTRA types are defined through formal schemas (JSON Schema, Protocol Buffers, Avro) that support forward and backward compatible evolution. This schema-first approach ensures that conversational applications can evolve independently while maintaining interoperability.

Version compatibility follows semantic versioning principles:
- **Patch versions (1.0.x)** - Bug fixes, no breaking changes
- **Minor versions (1.x.0)** - New optional fields, backward compatible additions
- **Major versions (x.0.0)** - Breaking changes requiring migration

### Universal Interoperability

ASTRA serves as a universal interchange format for conversational systems. Any system that can produce or consume ASTRA acts can interoperate with any other ASTRA-compliant system, regardless of underlying technology, programming language, or business domain.

## The Act Ledger Model

At ASTRA's core is the **Act Ledger** - an append-only sequence of atomic conversational actions. Every business conversation becomes a verifiable, auditable sequence of acts that capture both the informational content and the structural semantics of human communication.

### Base Act Structure

Every act in ASTRA shares a common base structure:

```typescript
interface Act {
  id: string;           // Unique identifier within conversation scope
  timestamp: ISO8601;   // When the act occurred
  speaker: string;      // Who performed this act
  type: ActType;        // What kind of act this is
  confidence?: number;  // Confidence score (0.0-1.0) for automated extraction
  source?: Source;      // Origin of this act (human, AI, system, etc.)
  metadata?: object;    // Additional context-specific information
}
```

This base structure ensures that every conversational action carries complete provenance information - who said what, when they said it, and how confident the system is in the interpretation.

### Act Type System

ASTRA defines five fundamental act types that capture the essential patterns of business conversation:

#### Ask - Request Missing Information

```typescript
interface Ask extends Act {
  type: "ask";
  field: string;                // What information is being requested
  prompt: string;               // Human-readable question or instruction
  constraints?: Constraint[];   // Validation rules for expected responses
  required?: boolean;           // Whether this information is mandatory
}
```

Ask acts represent information requests - questions, prompts, or instructions that solicit specific data from conversation participants. They establish the conversational contract for what information is needed and how it should be structured.

#### Fact - Declare State Information

```typescript
interface Fact extends Act {
  type: "fact";
  entity: EntityId;             // Which business entity is being modified
  field: string;                // What property of the entity is being set
  value: any;                   // The value being assigned
  operation?: "set" | "append"; // How to apply this value
  confidence?: number;          // Confidence in this fact extraction
}
```

Fact acts capture state declarations - information provided by conversation participants that establishes or modifies the state of business entities. They form the factual foundation of conversational state.

#### Confirm - Verify Understanding

```typescript
interface Confirm extends Act {
  type: "confirm";
  entity: EntityId;        // Entity being confirmed
  summary: string;         // Human-readable summary of understanding
  awaiting?: boolean;      // Whether confirmation is still pending
  confirmed?: boolean;     // Whether confirmation was received
}
```

Confirm acts ensure accuracy before commitment. They represent the verification phase where systems or participants validate their understanding of the current state before taking irreversible actions.

#### Commit - Execute Business Processes

```typescript
interface Commit extends Act {
  type: "commit";
  entity: EntityId;                        // Entity being acted upon
  action: "create" | "update" | "delete";  // Type of operation
  system?: string;                         // Target system for execution
  transaction_id?: string;                 // External system transaction ID
}
```

Commit acts trigger business process execution. They represent the transition from conversational state to business system integration, executing the intentions captured in the preceding conversation.

#### Error - Handle Failures

```typescript
interface Error extends Act {
  type: "error";
  code: string;            // Machine-readable error code
  message: string;         // Human-readable error description
  recoverable: boolean;    // Whether the error can be recovered from
  context?: object;        // Additional error context
}
```

Error acts capture failures and exceptions that occur during conversation processing. They ensure that problems are represented as first-class conversational elements that can be addressed and resolved.

## Entity Management

ASTRA's entity system provides structured representation for business objects that conversations manipulate. Entities serve as the nouns in conversational sentences, while acts serve as the verbs.

### Entity Structure

```typescript
interface Entity {
  id: string;              // Unique identifier within conversation scope
  type: string;            // Business entity type (order, customer, appointment)
  external_id?: string;    // Identifier in external business system
  system?: string;         // Which external system owns this entity
  version?: string;        // Version or revision of entity
  schema_url?: string;     // URL to schema definition
  metadata?: object;       // Additional entity-specific data
}
```

### Entity References

Acts can reference entities either by simple string identifier or through structured entity references that carry additional context:

```typescript
type EntityRef = string | Entity;
```

This flexibility allows acts to reference entities at different levels of specificity - from simple IDs for lightweight references to full entity objects for complex business contexts.

## Conversation Structure

ASTRA conversations serve as containers for acts and provide the broader context within which conversational state evolves.

### Conversation Components

```typescript
interface Conversation {
  id: string;                    // Unique conversation identifier
  participants: Participant[];   // Who is involved in this conversation
  acts: Act[];                   // Sequence of conversational actions
  status: ConversationStatus;    // Current conversation state
  created_at: ISO8601;          // When conversation began
  updated_at: ISO8601;          // When conversation last changed
  metadata?: object;            // Additional conversation context
}
```

### Participant Management

```typescript
interface Participant {
  id: string;                   // Unique participant identifier
  type: "human" | "ai" | "system" | "bot"; // What kind of participant
  role?: string;                // Business role (customer, agent, manager)
  name?: string;                // Display name
  metadata?: object;            // Additional participant context
}
```

Participants establish who is involved in conversations and what roles they play. This information proves essential for business process routing, security controls, and audit requirements.

## Constraint System

ASTRA's constraint system enables sophisticated validation and business rule enforcement within conversational flows.

### Constraint Types

```typescript
interface Constraint {
  type: "required" | "optional" | "min_length" | "max_length" | 
        "pattern" | "format" | "range" | "enum" | "custom";
  value?: any;                  // Constraint-specific parameter
  message?: string;             // Human-readable error message
  code?: string;                // Machine-readable error code
}
```

Constraints can be applied to Ask acts to specify validation requirements for expected responses, ensuring that conversational input meets business requirements before being committed to systems.

## Schema Evolution Strategy

ASTRA's architecture anticipates continuous evolution of conversational requirements and business processes. The schema evolution strategy ensures that systems can upgrade independently while maintaining backward compatibility.

### Compatibility Guarantees

**Within Major Versions:**
- All existing acts remain valid
- All existing fields maintain their semantics
- All existing validations continue to work
- New optional fields can be added
- New act types can be introduced

**Across Major Versions:**
- Migration tools provide automated upgrade paths
- Compatibility matrices document breaking changes
- Deprecation warnings provide advance notice
- Fallback mechanisms handle unsupported features

### Extension Mechanisms

ASTRA supports extension through several mechanisms:

**Custom Act Types:**
Systems can define domain-specific act types that extend the base Act interface while maintaining compatibility with ASTRA-compliant processors.

**Metadata Fields:**
All ASTRA types include metadata fields that can carry additional context without requiring schema changes.

**Entity Type System:**
Business domains can define custom entity types with domain-specific validation and processing rules.

## Implementation Patterns

### Multi-Language Support

ASTRA provides canonical implementations in multiple programming languages:

- **TypeScript/JavaScript** - Web applications and Node.js services
- **Python** - AI/ML processing and data analysis workflows  
- **Go** - High-performance server applications and microservices

Each implementation provides:
- Type-safe interfaces for all ASTRA types
- JSON schema validation for runtime safety
- Utility functions for common operations
- Migration tools for schema evolution

### Validation Architecture

ASTRA implements a multi-layered validation approach:

**Compile-Time Validation:**
Type systems catch structural errors during development, ensuring that applications can only construct valid ASTRA types.

**Runtime Validation:**
JSON schemas validate data at system boundaries, ensuring that external inputs conform to ASTRA specifications.

**Business Rule Validation:**
Constraint systems enable domain-specific validation that goes beyond structural correctness to business rule compliance.

### Integration Patterns

**Event-Driven Architecture:**
ASTRA acts serve as natural event payloads in event-driven systems, enabling loose coupling between conversational processing and business system integration.

**Stream Processing:**
Act sequences can be processed as streams, enabling real-time conversation analysis and business process orchestration.

**Message Queue Integration:**
Acts provide structured payloads for message queue systems, enabling reliable, asynchronous processing of conversational state changes.

## Ecosystem Vision

ASTRA aims to establish industry-standard protocols for conversational business applications. The architecture anticipates an ecosystem where:

**Conversational Applications** can focus on user experience while relying on ASTRA for reliable state management and business system integration.

**Business Systems** can consume conversational state through standardized ASTRA interfaces without needing to understand the complexities of natural language processing.

**Channel Providers** can produce ASTRA acts from any communication medium - voice, text, email, chat, or video - enabling universal conversational processing.

**Analytics Platforms** can process ASTRA streams to provide insights into conversational business processes without requiring custom integrations.

This ecosystem approach transforms conversation from an isolated application feature into a composable business capability that can be mixed, matched, and evolved independently across different systems and vendors.

## Future Considerations

ASTRA's architecture anticipates several evolving requirements:

**Multimodal Conversations:**
Support for conversations that blend text, voice, images, and other media through structured metadata and entity references.

**Distributed Conversations:**
Protocols for conversations that span multiple organizations, systems, and security boundaries while maintaining state consistency.

**Real-Time Processing:**
Optimizations for ultra-low-latency conversational applications that require immediate response to state changes.

**Privacy and Compliance:**
Enhanced support for privacy-preserving conversational processing and regulatory compliance requirements.

The modular, schema-first architecture ensures that these capabilities can be added incrementally without disrupting existing implementations or requiring coordinated upgrades across the ecosystem.

## Conclusion

ASTRA represents a foundational shift toward treating conversation as a structured, reliable business capability rather than an unstructured communication medium. By providing canonical types, evolution mechanisms, and interoperability standards, ASTRA enables the next generation of conversational business applications to be built on solid, interoperable foundations.

The architecture's success will be measured not by its technical sophistication but by its ability to enable seamless interoperability between conversational systems while maintaining the flexibility to evolve with changing business requirements. Through careful attention to schema evolution, type safety, and extensibility, ASTRA aims to provide the conversational infrastructure that will power the next decade of business communication innovation.
