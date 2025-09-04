# ASTRA Act Types Reference

This document provides the complete specification for all five ASTRA act types. Each act type represents a fundamental pattern of business conversation and includes detailed field descriptions, validation rules, usage patterns, and implementation examples.

## Base Act Structure

All ASTRA acts share a common base structure that provides essential metadata:

```typescript
interface Act {
  id: string;              // Unique identifier within conversation scope
  timestamp: string;       // ISO 8601 timestamp when act occurred
  speaker: string;         // Participant who performed this act
  type: ActType;           // What kind of act this is
  confidence?: number;     // Confidence score (0.0-1.0) for automated extraction
  source?: Source;         // Origin: 'human' | 'ai' | 'system' | 'speech_recognition' | 'text_analysis'
  metadata?: ActMetadata;  // Additional context-specific information
}
```

### Base Fields

**`id`** (required)
- Unique identifier for this act within the conversation scope
- Format: `act_[alphanumeric]` (e.g., `act_001`, `act_pizza_size_request`)
- Must be unique within the conversation

**`timestamp`** (required)
- ISO 8601 formatted timestamp indicating when the act occurred
- Must include timezone information (Z for UTC recommended)
- Used for conversation sequencing and audit trails

**`speaker`** (required)
- Identifier of the conversation participant who performed this act
- References a participant ID defined in the conversation's participant list
- Enables proper attribution and role-based processing

**`type`** (required)
- Specifies which act type this is: `ask`, `fact`, `confirm`, `commit`, or `error`
- Determines which additional fields are required and how the act should be processed

**`confidence`** (optional)
- Floating point value between 0.0 and 1.0
- Indicates confidence level for automated act extraction from natural language
- Higher values indicate greater certainty in the interpretation

**`source`** (optional)
- Indicates how this act was generated or extracted
- Values: `human`, `ai`, `system`, `speech_recognition`, `text_analysis`
- Useful for debugging and quality assurance

**`metadata`** (optional)
- Additional context-specific information
- Common fields: `channel`, `language`, `original_text`, `processing_time_ms`
- Extensible for domain-specific requirements

---

## Ask Acts

Ask acts request missing information required to complete a business process. They establish the conversational contract for what information is needed and how it should be structured.

```typescript
interface Ask extends Act {
  type: "ask";
  field: string;                    // What information is being requested
  prompt: string;                   // Human-readable question or instruction
  constraints?: Constraint[];       // Validation rules for expected responses
  required?: boolean;               // Whether this information is mandatory
  expected_type?: ExpectedType;     // Expected data type of response
  retry_count?: number;             // Number of times this has been asked
  max_retries?: number;             // Maximum retry attempts before escalation
}
```

### Ask-Specific Fields

**`field`** (required)
- Identifies what specific information is being requested
- Should be a clear, consistent identifier (e.g., `email`, `delivery_address`, `payment_method`)
- Used to map responses to the correct entity fields

**`prompt`** (required)
- Human-readable question or instruction presented to obtain the information
- Should be clear, concise, and actionable
- Examples: "What's your email address?", "Please provide your delivery address"

**`constraints`** (optional)
- Array of validation constraints that define acceptable responses
- Enables automatic validation before accepting user input
- See [Constraint System](#constraint-system) for detailed specification

**`required`** (optional, default: true)
- Boolean indicating whether this information is mandatory to proceed
- Optional fields can be skipped without blocking conversation progress
- Affects validation and conversation flow logic

**`expected_type`** (optional)
- Hints about the expected data type of the response
- Values: `string`, `number`, `boolean`, `object`, `array`, `date`, `email`, `phone`, `address`
- Helps with input parsing and validation

**`retry_count`** (optional, default: 0)
- Number of times this specific question has been asked
- Increments automatically on repeated asks
- Used for escalation and alternative prompting strategies

**`max_retries`** (optional, default: 3)
- Maximum number of times to ask before escalating or failing
- Prevents infinite retry loops
- Triggers alternative conversation flows when exceeded

### Usage Patterns

**Basic Information Request**
```typescript
const askEmail: Ask = {
  id: 'act_001',
  timestamp: '2025-01-15T14:30:00Z',
  speaker: 'agent_123',
  type: 'ask',
  field: 'email',
  prompt: 'What is your email address?',
  constraints: [{
    type: 'format',
    value: 'email',
    message: 'Please provide a valid email address'
  }]
};
```

**Complex Validation**
```typescript
const askCreditCard: Ask = {
  id: 'act_002',
  timestamp: '2025-01-15T14:31:00Z',
  speaker: 'payment_bot',
  type: 'ask',
  field: 'credit_card_number',
  prompt: 'Please enter your credit card number',
  required: true,
  expected_type: 'string',
  constraints: [
    {
      type: 'pattern',
      value: '^[0-9]{13,19}$',
      message: 'Credit card number must be 13-19 digits'
    },
    {
      type: 'custom',
      value: { luhn_check: true },
      message: 'Invalid credit card number'
    }
  ]
};
```

**Optional Information with Retry Logic**
```typescript
const askPhoneOptional: Ask = {
  id: 'act_003', 
  timestamp: '2025-01-15T14:32:00Z',
  speaker: 'agent_123',
  type: 'ask',
  field: 'phone_number',
  prompt: 'Would you like to provide a phone number for order updates? (Optional)',
  required: false,
  retry_count: 1,
  max_retries: 2,
  constraints: [{
    type: 'format',
    value: 'phone',
    message: 'Please provide a valid phone number'
  }]
};
```

---

## Fact Acts

Fact acts declare information provided during conversation. They capture state changes and form the factual foundation of conversational state.

```typescript
interface Fact extends Act {
  type: "fact";
  entity: EntityRef;                    // Which business entity is being modified
  field: string;                        // What property is being set
  value: any;                           // The value being assigned
  operation?: FieldOperation;           // How to apply this value
  previous_value?: any;                 // Previous value (for audit trail)
  validation_status?: ValidationStatus; // Validation state
  validation_errors?: string[];         // Errors if validation failed
}
```

### Fact-Specific Fields

**`entity`** (required)
- References the business entity being modified
- Can be a simple string ID or a structured entity object
- Examples: `order_789`, `customer_123`, `appointment_456`

**`field`** (required)
- Specific property or field of the entity being set
- Should match the field names used in Ask acts
- Examples: `email`, `delivery_address`, `total_amount`

**`value`** (required)
- The actual value being assigned to the field
- Can be any JSON-serializable type (string, number, boolean, object, array)
- Should conform to the constraints defined in corresponding Ask acts

**`operation`** (optional, default: "set")
- How the value should be applied to the field
- Values: `set`, `append`, `increment`, `decrement`, `delete`, `merge`
- Enables sophisticated state manipulation patterns

**`previous_value`** (optional)
- The value that was previously stored in this field
- Provides audit trail and enables rollback operations
- Automatically populated by some implementations

**`validation_status`** (optional, default: "pending")
- Current validation state of this fact
- Values: `pending`, `valid`, `invalid`, `partial`
- Updated as validation processes complete

**`validation_errors`** (optional)
- Array of validation error messages if validation failed
- Provides specific feedback about what constraints were violated
- Used for error reporting and conversation repair

### Field Operations

**`set`** (default) - Replace the current value entirely
```typescript
const setEmail: Fact = {
  id: 'act_004',
  timestamp: '2025-01-15T14:33:00Z',
  speaker: 'customer_456',
  type: 'fact',
  entity: 'customer_123',
  field: 'email',
  value: 'user@example.com',
  operation: 'set'
};
```

**`append`** - Add to existing array or string value
```typescript
const addItem: Fact = {
  id: 'act_005',
  timestamp: '2025-01-15T14:34:00Z',
  speaker: 'customer_456',
  type: 'fact', 
  entity: 'order_789',
  field: 'items',
  value: { product: 'pizza', size: 'large', quantity: 1 },
  operation: 'append'
};
```

**`increment`/`decrement`** - Mathematical operations on numeric values
```typescript
const updateQuantity: Fact = {
  id: 'act_006',
  timestamp: '2025-01-15T14:35:00Z',
  speaker: 'customer_456', 
  type: 'fact',
  entity: 'order_789',
  field: 'total_quantity',
  value: 2,
  operation: 'increment'
};
```

**`merge`** - Merge object values (shallow or deep)
```typescript
const updateAddress: Fact = {
  id: 'act_007',
  timestamp: '2025-01-15T14:36:00Z',
  speaker: 'customer_456',
  type: 'fact',
  entity: 'customer_123', 
  field: 'address',
  value: { apartment: '4B' },
  operation: 'merge' // Merges with existing address object
};
```

### Validation Status Flow

Facts typically progress through validation states:

1. **`pending`** - Fact received but not yet validated
2. **`valid`** - Passed all validation constraints
3. **`invalid`** - Failed validation (check validation_errors)
4. **`partial`** - Some constraints passed, others failed

---

## Confirm Acts

Confirm acts verify understanding of information before commitment. They represent the verification phase where systems or participants validate their understanding before taking irreversible actions.

```typescript
interface Confirm extends Act {
  type: "confirm";
  entity: EntityRef;                        // Entity being confirmed
  summary: string;                          // Human-readable summary
  awaiting?: boolean;                       // Whether confirmation is pending
  confirmed?: boolean;                      // Whether confirmation was received
  confirmation_method?: ConfirmationMethod; // How confirmation was obtained
  fields_confirmed?: string[];              // Specific fields being confirmed
  rejection_reason?: string;                // Reason if confirmation rejected
  timeout_ms?: number;                      // Timeout for awaiting confirmation
}
```

### Confirm-Specific Fields

**`entity`** (required)
- References the business entity being confirmed
- Should match entities referenced in related Facts
- The entity whose state is being verified

**`summary`** (required)
- Human-readable summary of what is being confirmed
- Should clearly describe the current understanding
- Examples: "Delivering 2 large pizzas to 123 Main St at 6 PM", "Charging $29.99 to card ending in 1234"

**`awaiting`** (optional, default: false)
- Boolean indicating whether confirmation is still pending
- True when waiting for user response
- False when confirmation process is complete

**`confirmed`** (optional)
- Boolean result of the confirmation process
- True if user confirmed, false if rejected
- Undefined while still awaiting confirmation

**`confirmation_method`** (optional)
- How the confirmation was obtained
- Values: `verbal`, `explicit`, `implicit`, `timeout`, `system`
- Useful for audit trails and compliance

**`fields_confirmed`** (optional)
- Array of specific field names being confirmed
- Enables partial confirmations of complex entities
- Examples: `['delivery_address', 'delivery_time']`

**`rejection_reason`** (optional)
- User-provided reason if confirmation was rejected
- Helps identify what needs to be corrected
- Examples: "Wrong address", "Price too high"

**`timeout_ms`** (optional)
- Timeout in milliseconds for awaiting confirmation
- After timeout, confirmation may auto-proceed or escalate
- Default timeout behavior is implementation-specific

### Confirmation Methods

**`verbal`** - Spoken confirmation ("Yes", "That's correct")
```typescript
const verbalConfirm: Confirm = {
  id: 'act_008',
  timestamp: '2025-01-15T14:37:00Z',
  speaker: 'customer_456',
  type: 'confirm',
  entity: 'order_789',
  summary: 'One large pepperoni pizza for $16.99',
  confirmed: true,
  confirmation_method: 'verbal'
};
```

**`explicit`** - Direct UI interaction (button click, checkbox)
```typescript
const explicitConfirm: Confirm = {
  id: 'act_009',
  timestamp: '2025-01-15T14:38:00Z',
  speaker: 'customer_456',
  type: 'confirm',
  entity: 'payment_info',
  summary: 'Charge $16.99 to Visa ending in 1234',
  confirmed: true,
  confirmation_method: 'explicit',
  fields_confirmed: ['payment_method', 'total_amount']
};
```

**`implicit`** - Inferred from context or continued interaction
```typescript
const implicitConfirm: Confirm = {
  id: 'act_010', 
  timestamp: '2025-01-15T14:39:00Z',
  speaker: 'system',
  type: 'confirm',
  entity: 'user_preferences',
  summary: 'Using saved delivery address',
  confirmed: true,
  confirmation_method: 'implicit'
};
```

**`timeout`** - Automatic confirmation after no response
```typescript
const timeoutConfirm: Confirm = {
  id: 'act_011',
  timestamp: '2025-01-15T14:40:00Z', 
  speaker: 'system',
  type: 'confirm',
  entity: 'order_789',
  summary: 'Proceeding with order after confirmation timeout',
  confirmed: true,
  confirmation_method: 'timeout',
  timeout_ms: 30000
};
```

### Confirmation Patterns

**Pre-Commit Verification**
```typescript
const preCommitConfirm: Confirm = {
  id: 'act_012',
  timestamp: '2025-01-15T14:41:00Z',
  speaker: 'agent_123',
  type: 'confirm',
  entity: 'order_789',
  summary: 'Ready to place your order: 2 large pizzas, total $29.98, deliver to 123 Main St',
  awaiting: true,
  timeout_ms: 60000
};
```

---

## Commit Acts

Commit acts execute business processes and trigger system integrations. They represent the transition from conversational state to business system execution.

```typescript
interface Commit extends Act {
  type: "commit";
  entity: EntityRef;                  // Entity being acted upon
  action: CommitAction;               // Type of operation to perform
  system?: string;                    // Target system identifier
  transaction_id?: string;            // External transaction ID
  status?: CommitStatus;              // Current operation status
  error?: CommitError;                // Error info if commit failed
  retry_count?: number;               // Number of retry attempts
  max_retries?: number;               // Maximum retry attempts
  idempotency_key?: string;           // Key for idempotent operations
  rollback_info?: object;             // Information for rollback if needed
}
```

### Commit-Specific Fields

**`entity`** (required)
- Business entity being acted upon by external systems
- Must reference an entity that has been populated with Facts
- The subject of the business process execution

**`action`** (required)
- Type of operation to perform on the entity
- Values: `create`, `update`, `delete`, `execute`, `cancel`, `pause`, `resume`
- Determines what business process logic to trigger

**`system`** (optional)
- Identifier of the target system for execution
- Examples: `order_management`, `crm`, `payment_processor`, `inventory`
- Enables routing to appropriate integration handlers

**`transaction_id`** (optional)
- External system transaction or record identifier
- Populated by integration handlers after successful execution
- Used for tracking and correlation with external systems

**`status`** (optional, default: "pending")
- Current status of the commit operation
- Values: `pending`, `in_progress`, `success`, `failed`, `retrying`, `cancelled`
- Updated as the commit progresses through execution

**`error`** (optional)
- Detailed error information if the commit failed
- Includes error code, message, details, and recoverability information
- Used for error handling and retry logic

**`retry_count`** (optional, default: 0)
- Number of retry attempts made for this commit
- Incremented automatically on retries
- Used with max_retries for retry limit enforcement

**`max_retries`** (optional, default: 3)
- Maximum number of retry attempts before giving up
- Prevents infinite retry loops
- Can be configured per commit based on operation criticality

**`idempotency_key`** (optional)
- Unique key to ensure idempotent operations
- Prevents duplicate execution if commit is retried
- Critical for financial transactions and other non-repeatable operations

**`rollback_info`** (optional)
- Information needed to rollback this commit if necessary
- Stored before execution to enable compensation transactions
- Used in saga patterns and distributed transaction management

### Commit Actions

**`create`** - Create new entity in target system
```typescript
const createOrder: Commit = {
  id: 'act_013',
  timestamp: '2025-01-15T14:42:00Z',
  speaker: 'order_system',
  type: 'commit',
  entity: 'order_789',
  action: 'create',
  system: 'order_management',
  idempotency_key: 'order_789_create_20250115'
};
```

**`update`** - Modify existing entity
```typescript
const updateCustomer: Commit = {
  id: 'act_014',
  timestamp: '2025-01-15T14:43:00Z',
  speaker: 'crm_system',
  type: 'commit',
  entity: 'customer_123',
  action: 'update',
  system: 'crm',
  transaction_id: 'crm_txn_456',
  status: 'success'
};
```

**`execute`** - Trigger business process
```typescript
const processPayment: Commit = {
  id: 'act_015',
  timestamp: '2025-01-15T14:44:00Z',
  speaker: 'payment_system',
  type: 'commit',
  entity: 'payment_request_789',
  action: 'execute',
  system: 'payment_processor',
  status: 'in_progress',
  max_retries: 1 // Financial operations typically have low retry limits
};
```

### Commit Status Progression

Commits typically follow this status progression:

1. **`pending`** - Commit created but not yet started
2. **`in_progress`** - Execution has begun
3. **`success`** - Completed successfully
4. **`failed`** - Execution failed (check error field)
5. **`retrying`** - Automatic retry in progress
6. **`cancelled`** - Manually cancelled before completion

### Error Handling

```typescript
const failedCommit: Commit = {
  id: 'act_016',
  timestamp: '2025-01-15T14:45:00Z',
  speaker: 'payment_system',
  type: 'commit',
  entity: 'payment_request_789',
  action: 'execute',
  system: 'payment_processor',
  status: 'failed',
  retry_count: 3,
  max_retries: 3,
  error: {
    code: 'PAYMENT_DECLINED',
    message: 'Credit card payment declined by issuer',
    details: { decline_code: '51', bank_message: 'Insufficient funds' },
    recoverable: true
  }
};
```

---

## Error Acts

Error acts handle failures and exceptions in conversational processing. They ensure that problems are represented as first-class conversational elements that can be addressed and resolved.

```typescript
interface Error extends Act {
  type: "error";
  code: string;                       // Machine-readable error code
  message: string;                    // Human-readable error description
  recoverable: boolean;               // Whether conversation can continue
  severity?: ErrorSeverity;           // Error severity level
  category?: ErrorCategory;           // Error classification
  details?: object;                   // Additional error context
  related_act_id?: string;            // ID of act that caused error
  suggested_action?: SuggestedAction; // Recommended recovery action
  user_message?: string;              // User-friendly error message
  stack_trace?: string;               // Technical debug information
}
```

### Error-Specific Fields

**`code`** (required)
- Machine-readable error identifier
- Should be consistent and parseable by automated systems
- Examples: `VALIDATION_FAILED`, `PAYMENT_DECLINED`, `TIMEOUT`, `SYSTEM_UNAVAILABLE`

**`message`** (required)
- Human-readable error description for technical users
- Should provide clear information about what went wrong
- Examples: "Field validation failed for email address", "Payment processor timeout"

**`recoverable`** (required)
- Boolean indicating whether the conversation can continue after this error
- True for errors that can be handled through conversation repair
- False for fatal errors that require conversation termination

**`severity`** (optional, default: "error")
- Severity level for error classification and handling
- Values: `info`, `warning`, `error`, `critical`
- Used for logging, alerting, and escalation decisions

**`category`** (optional)
- Classification category for error analysis and handling
- Values: `validation`, `processing`, `integration`, `timeout`, `permission`, `system`, `user_input`, `business_rule`
- Enables category-specific error handling strategies

**`details`** (optional)
- Additional context and debugging information
- Can include technical details, request parameters, system state
- Not typically shown to end users

**`related_act_id`** (optional)
- ID of the act that caused or triggered this error
- Enables tracing error back to root cause
- Useful for debugging and conversation repair

**`suggested_action`** (optional)
- Recommended recovery action for automated or human handlers
- Values: `retry`, `escalate`, `ignore`, `clarify`, `fallback`, `terminate`
- Guides error recovery strategies

**`user_message`** (optional)
- User-friendly error message appropriate for end users
- Should be clear, helpful, and actionable
- Differs from technical message field

**`stack_trace`** (optional)
- Technical stack trace for debugging purposes
- Should never be shown to end users
- Useful for development and production debugging

### Error Severity Levels

**`info`** - Informational messages, no action required
```typescript
const infoError: Error = {
  id: 'act_017',
  timestamp: '2025-01-15T14:46:00Z',
  speaker: 'system',
  type: 'error',
  code: 'FEATURE_UNAVAILABLE',
  message: 'Voice ordering temporarily unavailable, using text input',
  recoverable: true,
  severity: 'info',
  category: 'system',
  user_message: 'Voice ordering is temporarily unavailable. You can still place your order by typing.'
};
```

**`warning`** - Potential issues that don't stop processing
```typescript
const warningError: Error = {
  id: 'act_018',
  timestamp: '2025-01-15T14:47:00Z',
  speaker: 'validation_system',
  type: 'error',
  code: 'DATA_QUALITY_WARNING',
  message: 'Phone number format unusual but accepted',
  recoverable: true,
  severity: 'warning',
  category: 'validation',
  details: { phone_number: '+1-555-PIZZA', format_score: 0.6 }
};
```

**`error`** (default) - Standard errors requiring attention
```typescript
const standardError: Error = {
  id: 'act_019',
  timestamp: '2025-01-15T14:48:00Z',
  speaker: 'payment_system',
  type: 'error',
  code: 'PAYMENT_DECLINED',
  message: 'Credit card payment declined by issuer',
  recoverable: true,
  severity: 'error',
  category: 'integration',
  related_act_id: 'act_015',
  suggested_action: 'clarify',
  user_message: 'Your payment was declined. Please try a different payment method.',
  details: { decline_code: '05', bank_message: 'Do not honor' }
};
```

**`critical`** - Severe errors requiring immediate attention
```typescript
const criticalError: Error = {
  id: 'act_020',
  timestamp: '2025-01-15T14:49:00Z',
  speaker: 'system',
  type: 'error',
  code: 'SYSTEM_FAILURE', 
  message: 'Database connection failed after all retries',
  recoverable: false,
  severity: 'critical',
  category: 'system',
  suggested_action: 'escalate',
  user_message: 'We are experiencing technical difficulties. Please try again later.'
};
```

### Error Categories

**`validation`** - Input validation failures
- Field format errors, constraint violations
- Usually recoverable with clarification

**`processing`** - Internal processing errors  
- Parsing failures, computation errors
- May be recoverable depending on cause

**`integration`** - External system failures
- API errors, network timeouts, service unavailable
- Often temporary and recoverable with retry

**`timeout`** - Operation timeout errors
- Network timeouts, processing timeouts
- Usually recoverable with retry or escalation

**`permission`** - Authorization/authentication failures
- Access denied, insufficient permissions
- May require different user or escalation

**`system`** - Infrastructure and platform errors
- Database errors, service failures
- Often require technical intervention

**`user_input`** - User input problems
- Unclear requests, invalid choices
- Recoverable with clarification

**`business_rule`** - Business logic violations
- Policy violations, constraint failures
- May require exception handling or escalation

### Suggested Recovery Actions

**`retry`** - Attempt the operation again
```typescript
suggested_action: 'retry',
details: { retry_delay_ms: 1000, max_retries: 3 }
```

**`escalate`** - Transfer to human agent or supervisor
```typescript
suggested_action: 'escalate',
details: { escalation_reason: 'payment_issue', priority: 'high' }
```

**`clarify`** - Ask user for clarification or correction
```typescript
suggested_action: 'clarify',
details: { clarification_needed: 'payment_method' }
```

**`fallback`** - Use alternative approach or default
```typescript
suggested_action: 'fallback',
details: { fallback_method: 'manual_processing' }
```

**`terminate`** - End conversation due to unrecoverable error
```typescript
suggested_action: 'terminate',
details: { termination_reason: 'system_failure' }
```

---

## Constraint System

Constraints define validation rules for Ask acts and can be used to validate Facts. They enable sophisticated business rule enforcement within conversational flows.

```typescript
interface Constraint {
  type: ConstraintType;     // Type of constraint
  value?: any;              // Constraint-specific parameter
  message?: string;         // Human-readable error message
  code?: string;            // Machine-readable error code
}
```

### Constraint Types

**`required`** - Field must be provided
```typescript
{
  type: 'required',
  message: 'This field is required'
}
```

**`min_length`** / **`max_length`** - String/array length constraints
```typescript
{
  type: 'min_length',
  value: 8,
  message: 'Password must be at least 8 characters'
}
```

**`pattern`** - Regular expression validation
```typescript
{
  type: 'pattern',
  value: '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$',
  message: 'Please enter a valid email address'
}
```

**`format`** - Pre-defined format validation
```typescript
{
  type: 'format',
  value: 'phone',  // email, phone, url, date, time, datetime, uuid, ipv4, ipv6
  message: 'Please enter a valid phone number'
}
```

**`range`** - Numeric range validation
```typescript
{
  type: 'range',
  value: { min: 1, max: 100, inclusive: true },
  message: 'Quantity must be between 1 and 100'
}
```

**`enum`** - Must be one of specified values
```typescript
{
  type: 'enum',
  value: ['small', 'medium', 'large'],
  message: 'Size must be small, medium, or large'
}
```

**`custom`** - Domain-specific validation logic
```typescript
{
  type: 'custom',
  value: { 
    business_rule: 'credit_check',
    parameters: { min_score: 650 }
  },
  message: 'Credit check required for this purchase amount'
}
```

## Best Practices

### Act Sequencing
- Follow the natural conversation flow: Ask → Fact → Confirm → Commit
- Use Error acts for exception handling at any point
- Maintain entity references consistently across related acts

### Validation Strategy
- Apply constraints at Ask time to set expectations
- Validate Facts against constraints before acceptance
- Use Error acts to report validation failures with specific guidance

### Error Recovery
- Make errors recoverable when possible
- Provide specific, actionable error messages
- Use suggested actions to guide automated recovery

### Entity Management
- Use consistent entity identifiers across conversations
- Include entity type information for proper validation
- Maintain entity version information for audit trails

### Performance Considerations
- Keep act payloads reasonable in size
- Use references rather than embedding large objects
- Consider batch operations for multiple related commits

This reference provides the foundation for implementing robust, type-safe conversational applications using the ASTRA specification. Each act type serves a specific purpose in the conversational state machine while maintaining compatibility and extensibility across different implementations and business domains.
