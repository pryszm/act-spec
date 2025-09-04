# Getting Started with ASTRA

This guide will help you get up and running with ASTRA (Act State Representation Architecture) in your conversational applications. Whether you're building a customer service bot, sales assistant, or any other conversation-driven system, ASTRA provides the foundational types and patterns you need for reliable, structured conversations.

## Quick Overview

ASTRA transforms conversations from unstructured transcripts into structured, auditable sequences of actions. Instead of parsing free-form text, you work with typed acts that represent the essential patterns of business conversation:

- **Ask** - Request missing information ("What's your email address?")
- **Fact** - Declare state information ("My email is user@example.com")
- **Confirm** - Verify understanding ("So you want 2 pizzas delivered to 123 Main St?")
- **Commit** - Execute business processes ("Order created with ID #12345")
- **Error** - Handle failures ("Payment method declined")

## Installation

Choose your language implementation:

### TypeScript/JavaScript
```bash
npm install @astra/model-ts
```

### Python
```bash
pip install astra-model-py
```

### Go
```bash
go get github.com/pryszm/astra-model-go
```

## Your First ASTRA Conversation

Let's build a simple pizza ordering conversation to demonstrate core concepts:

### TypeScript Example

```typescript
import { 
  Ask, 
  Fact, 
  Confirm, 
  Commit, 
  Conversation,
  Participant
} from '@astra/model-ts';

// Create participants
const customer: Participant = {
  id: 'customer_123',
  type: 'human',
  role: 'customer'
};

const agent: Participant = {
  id: 'agent_456', 
  type: 'ai',
  role: 'sales_assistant'
};

// Build conversation acts
const askSize: Ask = {
  id: 'act_001',
  timestamp: '2025-01-15T14:30:00Z',
  speaker: 'agent_456',
  type: 'ask',
  field: 'pizza_size',
  prompt: 'What size pizza would you like?',
  constraints: [{
    type: 'enum',
    value: ['small', 'medium', 'large'],
    message: 'Please choose small, medium, or large'
  }]
};

const provideSize: Fact = {
  id: 'act_002', 
  timestamp: '2025-01-15T14:31:00Z',
  speaker: 'customer_123',
  type: 'fact',
  entity: 'order_789',
  field: 'pizza_size',
  value: 'large'
};

const confirmOrder: Confirm = {
  id: 'act_003',
  timestamp: '2025-01-15T14:32:00Z',
  speaker: 'agent_456',
  type: 'confirm',
  entity: 'order_789',
  summary: 'One large pepperoni pizza for delivery to 123 Main St',
  awaiting: true
};

const createOrder: Commit = {
  id: 'act_004',
  timestamp: '2025-01-15T14:33:00Z', 
  speaker: 'agent_456',
  type: 'commit',
  entity: 'order_789',
  action: 'create',
  system: 'order_management'
};

// Create the conversation
const conversation: Conversation = {
  id: 'conv_pizza_order',
  participants: [customer, agent],
  acts: [askSize, provideSize, confirmOrder, createOrder],
  status: 'active'
};
```

### Python Example

```python
from astra_model import Ask, Fact, Confirm, Commit, Conversation, Participant
from datetime import datetime

# Create participants
customer = Participant(
    id="customer_123",
    type="human", 
    role="customer"
)

agent = Participant(
    id="agent_456",
    type="ai",
    role="sales_assistant"
)

# Build conversation acts
ask_size = Ask(
    id="act_001",
    timestamp=datetime.now().isoformat() + "Z",
    speaker="agent_456",
    type="ask",
    field="pizza_size",
    prompt="What size pizza would you like?",
    constraints=[{
        "type": "enum",
        "value": ["small", "medium", "large"],
        "message": "Please choose small, medium, or large"
    }]
)

provide_size = Fact(
    id="act_002",
    timestamp=datetime.now().isoformat() + "Z",
    speaker="customer_123", 
    type="fact",
    entity="order_789",
    field="pizza_size",
    value="large"
)

confirm_order = Confirm(
    id="act_003",
    timestamp=datetime.now().isoformat() + "Z",
    speaker="agent_456",
    type="confirm", 
    entity="order_789",
    summary="One large pepperoni pizza for delivery to 123 Main St",
    awaiting=True
)

create_order = Commit(
    id="act_004",
    timestamp=datetime.now().isoformat() + "Z",
    speaker="agent_456",
    type="commit",
    entity="order_789", 
    action="create",
    system="order_management"
)

# Create the conversation
conversation = Conversation(
    id="conv_pizza_order",
    participants=[customer, agent],
    acts=[ask_size, provide_size, confirm_order, create_order],
    status="active"
)

# Convert to JSON for storage/transmission
conversation_json = conversation.model_dump_json()
```

### Go Example

```go
package main

import (
    "fmt"
    "time"
    
    astra "github.com/pryszm/astra-model-go"
)

func main() {
    // Create participants
    customer := astra.Participant{
        ID:   "customer_123",
        Type: astra.ParticipantTypeHuman,
        Role: astra.StringPtr("customer"),
    }
    
    agent := astra.Participant{
        ID:   "agent_456", 
        Type: astra.ParticipantTypeAI,
        Role: astra.StringPtr("sales_assistant"),
    }
    
    // Build conversation acts
    askSize := astra.Ask{
        Act: astra.Act{
            ID:        "act_001",
            Timestamp: time.Now(),
            Speaker:   "agent_456",
            Type:      astra.ActTypeAsk,
        },
        Field:  "pizza_size",
        Prompt: "What size pizza would you like?",
        Constraints: []astra.Constraint{{
            Type: astra.ConstraintTypeEnum,
            Value: []string{"small", "medium", "large"},
            Message: astra.StringPtr("Please choose small, medium, or large"),
        }},
    }
    
    provideSize := astra.Fact{
        Act: astra.Act{
            ID:        "act_002",
            Timestamp: time.Now(),
            Speaker:   "customer_123", 
            Type:      astra.ActTypeFact,
        },
        Entity: "order_789",
        Field:  "pizza_size",
        Value:  "large",
    }
    
    confirmOrder := astra.Confirm{
        Act: astra.Act{
            ID:        "act_003",
            Timestamp: time.Now(),
            Speaker:   "agent_456",
            Type:      astra.ActTypeConfirm,
        },
        Entity:   "order_789",
        Summary:  "One large pepperoni pizza for delivery to 123 Main St",
        Awaiting: astra.BoolPtr(true),
    }
    
    createOrder := astra.Commit{
        Act: astra.Act{
            ID:        "act_004", 
            Timestamp: time.Now(),
            Speaker:   "agent_456",
            Type:      astra.ActTypeCommit,
        },
        Entity: "order_789",
        Action: astra.CommitActionCreate,
        System: astra.StringPtr("order_management"),
    }
    
    // Create the conversation
    conversation := astra.Conversation{
        ID:           "conv_pizza_order",
        Participants: []astra.Participant{customer, agent},
        Acts:         []astra.ConversationAct{askSize, provideSize, confirmOrder, createOrder},
        Status:       astra.ConversationStatusActive,
    }
    
    fmt.Printf("Created conversation with %d acts\n", len(conversation.Acts))
}
```

## Core Concepts in Action

### 1. Act Sequencing

Notice how acts flow naturally in the examples above:
1. **Ask** - Agent requests information
2. **Fact** - Customer provides information  
3. **Confirm** - Agent verifies understanding
4. **Commit** - System executes the business process

This sequence represents the fundamental pattern of business conversations: information gathering, validation, and execution.

### 2. Entity Management

The `entity` field in Facts, Confirms, and Commits ties acts to business objects. In our pizza example, `order_789` represents the order being built throughout the conversation. This enables:

- **State tracking** - All information about the order is tied to one entity
- **Validation** - Business rules can be applied per entity type
- **System integration** - External systems receive structured entity data

### 3. Constraint Validation

The `constraints` array on Ask acts defines what constitutes valid input:

```typescript
constraints: [{
  type: 'enum',
  value: ['small', 'medium', 'large'],
  message: 'Please choose small, medium, or large'
}]
```

This enables automatic validation of user responses before they're committed to business systems.

## Validation and Type Safety

ASTRA provides multiple layers of validation:

### Compile-Time Safety (TypeScript/Go)

Static typing catches structural errors during development:

```typescript
// TypeScript will catch this error - missing required fields
const invalidAsk: Ask = {
  id: 'act_001',
  // Missing: timestamp, speaker, type, field, prompt
};
```

### Runtime Validation

JSON Schema validation ensures data integrity at system boundaries:

```typescript
import { schemas, isAsk } from '@astra/model-ts';

// Type guard validation
if (isAsk(someAct)) {
  // TypeScript now knows someAct is an Ask
  console.log(`Requesting field: ${someAct.field}`);
}

// Schema validation
const validationResult = validate(schemas.ask, someActData);
```

### Business Rule Validation

Constraints enable domain-specific validation:

```python
from astra_model import Constraint

# Email format validation
email_constraint = Constraint(
    type="format",
    value="email",
    message="Please provide a valid email address"
)

# Custom business rule validation
credit_limit_constraint = Constraint(
    type="custom", 
    value={"max_amount": 5000, "check_credit_score": True},
    message="Credit limit exceeded"
)
```

## Integration Patterns

### Event-Driven Architecture

ASTRA acts work naturally as event payloads:

```typescript
// Publish acts to message queue
await messageQueue.publish('conversation.fact', {
  conversationId: 'conv_123',
  act: provideSize
});

// Subscribe to commits for system integration
messageQueue.subscribe('conversation.commit', async (event) => {
  if (event.act.system === 'order_management') {
    await orderSystem.createOrder(event.act.entity);
  }
});
```

### Stream Processing

Process conversations as real-time streams:

```python
from astra_model import ConversationAct

def process_act_stream(act: ConversationAct):
    if act.type == "fact":
        # Update entity state
        entity_store.update(act.entity, act.field, act.value)
    elif act.type == "commit": 
        # Trigger business process
        business_process.execute(act.entity, act.action)
```

### API Integration

Use ASTRA types for consistent API contracts:

```go
// HTTP endpoint that accepts ASTRA acts
func handleAct(w http.ResponseWriter, r *http.Request) {
    var act astra.ConversationAct
    if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Process act based on type
    switch act.GetType() {
    case astra.ActTypeFact:
        handleFact(act.(astra.Fact))
    case astra.ActTypeCommit:
        handleCommit(act.(astra.Commit))
    }
}
```

## Building Your Business Schema

Define entities and constraints that match your business domain:

### E-commerce Example

```typescript
interface Product {
  id: string;
  name: string;
  price: number;
  inventory: number;
}

interface Order {
  id: string;
  customer_id: string;
  items: OrderItem[];
  shipping_address: Address;
  status: 'pending' | 'confirmed' | 'shipped' | 'delivered';
}

// ASTRA acts for product catalog
const askProduct: Ask = {
  id: 'act_001',
  timestamp: '2025-01-15T14:30:00Z',
  speaker: 'bot',
  type: 'ask',
  field: 'product_name',
  prompt: 'What product are you looking for?',
  constraints: [{
    type: 'custom',
    value: { search_catalog: true },
    message: 'Product not found in catalog'
  }]
};
```

### Healthcare Example

```python
# Patient appointment scheduling
ask_appointment = Ask(
    id="act_001",
    timestamp=datetime.now().isoformat() + "Z",
    speaker="scheduler_bot",
    type="ask",
    field="preferred_date",
    prompt="What date would you prefer for your appointment?",
    constraints=[
        {
            "type": "format",
            "value": "date",
            "message": "Please provide a valid date"
        },
        {
            "type": "custom",
            "value": {"business_days_only": True, "min_advance_hours": 24},
            "message": "Appointments must be scheduled at least 24 hours in advance on business days"
        }
    ]
)
```

## Error Handling and Recovery

ASTRA's Error acts enable systematic failure handling:

```typescript
const paymentError: Error = {
  id: 'act_error_001',
  timestamp: '2025-01-15T14:35:00Z',
  speaker: 'payment_system',
  type: 'error',
  code: 'PAYMENT_DECLINED',
  message: 'Credit card payment was declined',
  recoverable: true,
  context: {
    error_code: '51',
    bank_message: 'Insufficient funds',
    suggested_actions: ['try_different_card', 'use_bank_transfer']
  }
};

// Recovery pattern
const askAlternatePayment: Ask = {
  id: 'act_recovery_001', 
  timestamp: '2025-01-15T14:36:00Z',
  speaker: 'payment_assistant',
  type: 'ask',
  field: 'payment_method',
  prompt: 'Your card was declined. Would you like to try a different payment method?',
  constraints: [{
    type: 'enum',
    value: ['different_card', 'bank_transfer', 'paypal'],
    message: 'Please choose an alternative payment method'
  }]
};
```

## Testing ASTRA Applications

### Unit Testing Acts

```typescript
import { Ask, isAsk, schemas } from '@astra/model-ts';

describe('Pizza Ordering', () => {
  test('should create valid Ask act', () => {
    const ask: Ask = {
      id: 'act_test_001',
      timestamp: '2025-01-15T14:30:00Z',
      speaker: 'test_agent',
      type: 'ask',
      field: 'pizza_size',
      prompt: 'What size pizza would you like?'
    };
    
    expect(isAsk(ask)).toBe(true);
    expect(ask.field).toBe('pizza_size');
  });
  
  test('should validate against JSON schema', () => {
    const askData = { /* act data */ };
    const result = validate(schemas.ask, askData);
    expect(result.valid).toBe(true);
  });
});
```

### Integration Testing

```python
import pytest
from astra_model import Conversation, Ask, Fact, Commit

def test_pizza_ordering_flow():
    conversation = Conversation(
        id="test_conv_001",
        participants=[],  # Add test participants
        acts=[],
        status="active"
    )
    
    # Test act sequence
    ask_size = Ask(...)  # Create test ask
    conversation.acts.append(ask_size)
    
    provide_size = Fact(...)  # Create test fact  
    conversation.acts.append(provide_size)
    
    create_order = Commit(...)  # Create test commit
    conversation.acts.append(create_order)
    
    # Assert conversation state
    assert len(conversation.acts) == 3
    assert conversation.acts[-1].type == "commit"
```

## Next Steps

Now that you understand the basics:

1. **Read the [Architecture Guide](./architecture.md)** to understand ASTRA's design principles
2. **Study the [Act Types Reference](./act-types.md)** for complete specifications  
3. **Check out [examples/](../examples/)** for more complex implementations
4. **Review [Integration Guide](./integration-guide.md)** for production deployment patterns

### Common Questions

**Q: How do I handle multi-turn conversations?**
A: Use the `Conversation` container to maintain state across multiple acts. Each act references the conversation ID and builds on previous acts.

**Q: Can I extend ASTRA with custom act types?**
A: Yes! ASTRA is designed for extensibility. You can define custom act types while maintaining compatibility with the base specification.

**Q: How do I integrate with my existing systems?**
A: Use Commit acts to trigger integrations. The `system` and `transaction_id` fields help track external operations.

**Q: What about privacy and compliance?**
A: ASTRA acts can include privacy metadata and support redaction/encryption patterns. See the Integration Guide for compliance strategies.

Ready to build structured, reliable conversations? The next step is exploring how different act types work together in complex business scenarios.
