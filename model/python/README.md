# astra-model-py

Python types and JSON schemas for ASTRA (Act State Representation Architecture) conversations.

## Installation

```bash
pip install astra-model-py
```

## Usage

### Import Types

```python
from astra_model import (
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
)
```

### Create Acts

```python
from datetime import datetime
from astra_model import Ask, Fact

# Create an Ask act
ask = Ask(
    id="act_001",
    timestamp=datetime.now().isoformat() + "Z",
    speaker="agent_123",
    type="ask",
    field="delivery_address",
    prompt="What is your delivery address?",
    required=True
)

# Create a Fact act
fact = Fact(
    id="act_002", 
    timestamp=datetime.now().isoformat() + "Z",
    speaker="customer_456",
    type="fact",
    entity="order_789",
    field="delivery_address",
    value="123 Main St, Anytown, USA",
    operation="set"
)
```

### Validate Acts using Pydantic

```python
from astra_model import Ask

# This will raise ValidationError if invalid
ask = Ask(
    id="act_001",
    timestamp="2025-01-15T14:30:00Z",
    speaker="agent_123",
    type="ask",
    field="email",
    prompt="What is your email address?"
)

# Validate data from dict
data = {
    "id": "act_001",
    "timestamp": "2025-01-15T14:30:00Z", 
    "speaker": "agent_123",
    "type": "ask",
    "field": "email",
    "prompt": "What is your email address?"
}
ask = Ask.model_validate(data)
```

### Access JSON Schemas

```python
from astra_model.schemas import SCHEMAS

# Get schema for validation
act_schema = SCHEMAS["act"]
ask_schema = SCHEMAS["ask"]

# Validate against JSON Schema
import jsonschema
jsonschema.validate(some_act_data, act_schema)
```

### Generate IDs

```python
from astra_model import generate_act_id, generate_conversation_id

act_id = generate_act_id()  # e.g., "act_lx2z8q_3k9p1r"
conversation_id = generate_conversation_id()  # e.g., "conv_lx2z8q_7m4n2s"
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

This library uses [Pydantic](https://pydantic.dev/) for runtime type validation and serialization. All ASTRA types are defined as Pydantic models, providing:

- **Runtime validation** - Automatic validation of field types and constraints
- **JSON serialization** - Easy conversion to/from JSON
- **IDE support** - Full type hints and autocompletion
- **Documentation** - Built-in field documentation and examples

## Schema Evolution

This package follows semantic versioning for schema compatibility:

- **Patch versions (1.0.x)** - Bug fixes, no breaking changes
- **Minor versions (1.x.0)** - New optional fields, backward compatible
- **Major versions (x.0.0)** - Breaking changes, migration required

## Examples

### Complete Conversation

```python
from astra_model import Conversation, Participant, Ask, Fact
from datetime import datetime

# Create participants
agent = Participant(id="agent_001", type="ai", role="customer_service")
customer = Participant(id="customer_123", type="human", role="customer")

# Create acts
ask_email = Ask(
    id="act_001",
    timestamp=datetime.now().isoformat() + "Z",
    speaker="agent_001",
    type="ask",
    field="email",
    prompt="What is your email address?"
)

provide_email = Fact(
    id="act_002",
    timestamp=datetime.now().isoformat() + "Z", 
    speaker="customer_123",
    type="fact",
    entity="customer_123",
    field="email",
    value="user@example.com"
)

# Create conversation
conversation = Conversation(
    id="conv_001",
    participants=[agent, customer],
    acts=[ask_email, provide_email],
    status="active"
)

# Convert to JSON
conversation_json = conversation.model_dump_json()

# Load from JSON
conversation_loaded = Conversation.model_validate_json(conversation_json)
```

## License

Licensed under the Apache License 2.0. See the [main repository](../../LICENSE) for details.

## Contributing

This package is part of the [ASTRA project](https://github.com/pryszm/astra). Please see the main repository for contribution guidelines.

## Links

- [ASTRA Repository](https://github.com/pryszm/astra)
- [Documentation](https://docs.pryszm.com)
- [Pryszm Platform](https://pryszm.com)
