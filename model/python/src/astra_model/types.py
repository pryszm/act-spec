"""
Tests for ASTRA Python types and utilities
"""

import json
import pytest
from datetime import datetime
from pydantic import ValidationError

from astra_model import (
    # Types
    Act,
    Ask,
    Fact,
    Confirm,
    Commit,
    Error,
    Participant,
    Entity,
    Conversation,
    Constraint,
    
    # Enums
    ActType,
    Source,
    ParticipantType,
    FieldOperation,
    ValidationStatus,
    ConfirmationMethod,
    CommitAction,
    CommitStatus,
    ErrorSeverity,
    ErrorCategory,
    SuggestedAction,
    ConversationStatus,
    
    # Utilities
    generate_act_id,
    generate_conversation_id,
    create_base_act,
    
    # Constants
    __version__,
    __schema_version__,
    
    # Schemas
    SCHEMAS
)


class TestActTypes:
    """Tests for core Act types"""
    
    def test_ask_creation(self):
        """Test creating Ask acts"""
        ask = Ask(
            id="act_001",
            timestamp="2025-01-15T14:30:00Z",
            speaker="agent_123",
            type=ActType.ASK,
            field="email",
            prompt="What is your email address?"
        )
        
        assert ask.id == "act_001"
        assert ask.type == ActType.ASK
        assert ask.field == "email"
        assert ask.prompt == "What is your email address?"
        assert ask.required is True  # default value
        
    def test_ask_validation(self):
        """Test Ask validation"""
        # Missing required field
        with pytest.raises(ValidationError):
            Ask(
                id="act_001",
                timestamp="2025-01-15T14:30:00Z",
                speaker="agent_123",
                type=ActType.ASK,
                # missing field
                prompt="What is your email?"
            )
            
        # Invalid act ID pattern
        with pytest.raises(ValidationError):
            Ask(
                id="invalid_id",  # doesn't match pattern
                timestamp="2025-01-15T14:30:00Z",
                speaker="agent_123",
                type=ActType.ASK,
                field="email",
                prompt="What is your email?"
            )
    
    def test_fact_creation(self):
        """Test creating Fact acts"""
        # With string entity
        fact1 = Fact(
            id="act_002",
            timestamp="2025-01-15T14:30:00Z",
            speaker="customer_456",
            type=ActType.FACT,
            entity="customer_456",
            field="email",
            value="user@example.com"
        )
        
        assert fact1.entity == "customer_456"
        assert fact1.value == "user@example.com"
        assert fact1.operation == FieldOperation.SET  # default
        
        # With Entity object
        entity = Entity(id="customer_456", type="customer")
        fact2 = Fact(
            id="act_003",
            timestamp="2025-01-15T14:30:00Z",
            speaker="customer_456",
            type=ActType.FACT,
            entity=entity,
            field="name",
            value="John Doe",
            operation=FieldOperation.APPEND
        )
        
        assert isinstance(fact2.entity, Entity)
        assert fact2.entity.id == "customer_456"
        assert fact2.operation == FieldOperation.APPEND
    
    def test_confirm_creation(self):
        """Test creating Confirm acts"""
        confirm = Confirm(
            id="act_004",
            timestamp="2025-01-15T14:30:00Z",
            speaker="agent_123",
            type=ActType.CONFIRM,
            entity="order_789",
            summary="Order for 2 pizzas to be delivered at 6 PM"
        )
        
        assert confirm.entity == "order_789"
        assert confirm.summary == "Order for 2 pizzas to be delivered at 6 PM"
        assert confirm.awaiting is True  # default
        assert confirm.confirmed is None  # not yet confirmed
    
    def test_commit_creation(self):
        """Test creating Commit acts"""
        commit = Commit(
            id="act_005",
            timestamp="2025-01-15T14:30:00Z",
            speaker="system_001",
            type=ActType.COMMIT,
            entity="order_789",
            action=CommitAction.CREATE,
            system="order_management"
        )
        
        assert commit.action == CommitAction.CREATE
        assert commit.system == "order_management"
        assert commit.status == CommitStatus.PENDING  # default
        assert commit.retry_count == 0  # default
    
    def test_error_creation(self):
        """Test creating Error acts"""
        error = Error(
            id="act_006",
            timestamp="2025-01-15T14:30:00Z",
            speaker="system_001",
            type=ActType.ERROR,
            code="VALIDATION_ERROR",
            message="Invalid email format",
            recoverable=True
        )
        
        assert error.code == "VALIDATION_ERROR"
        assert error.message == "Invalid email format"
        assert error.recoverable is True
        assert error.severity == ErrorSeverity.ERROR  # default


class TestEntityTypes:
    """Tests for Entity types"""
    
    def test_entity_creation(self):
        """Test creating Entity objects"""
        entity = Entity(
            id="customer_123",
            type="customer",
            external_id="cust_ext_456",
            system="crm",
            metadata={"priority": "high"}
        )
        
        assert entity.id == "customer_123"
        assert entity.type == "customer"
        assert entity.external_id == "cust_ext_456"
        assert entity.system == "crm"
        assert entity.metadata == {"priority": "high"}
    
    def test_entity_validation(self):
        """Test Entity validation"""
        # Missing required fields
        with pytest.raises(ValidationError):
            Entity(id="test")  # missing type
            
        with pytest.raises(ValidationError):
            Entity(type="customer")  # missing id


class TestParticipantTypes:
    """Tests for Participant types"""
    
    def test_participant_creation(self):
        """Test creating Participant objects"""
        participant = Participant(
            id="agent_001",
            type=ParticipantType.AI,
            role="customer_service",
            name="Support Agent",
            capabilities=["text", "voice"],
            preferences={
                "language": "en-US",
                "timezone": "America/New_York"
            }
        )
        
        assert participant.id == "agent_001"
        assert participant.type == ParticipantType.AI
        assert participant.role == "customer_service"
        assert "text" in participant.capabilities
        assert participant.preferences.language == "en-US"


class TestConversationTypes:
    """Tests for Conversation types"""
    
    def test_conversation_creation(self):
        """Test creating Conversation objects"""
        agent = Participant(id="agent_001", type=ParticipantType.AI)
        customer = Participant(id="customer_123", type=ParticipantType.HUMAN)
        
        ask = Ask(
            id="act_001",
            timestamp="2025-01-15T14:30:00Z",
            speaker="agent_001",
            type=ActType.ASK,
            field="email",
            prompt="What is your email?"
        )
        
        conversation = Conversation(
            id="conv_001",
            participants=[agent, customer],
            acts=[ask],
            status=ConversationStatus.ACTIVE
        )
        
        assert conversation.id == "conv_001"
        assert len(conversation.participants) == 2
        assert len(conversation.acts) == 1
        assert conversation.status == ConversationStatus.ACTIVE
    
    def test_conversation_validation(self):
        """Test Conversation validation"""
        # Missing participants
        with pytest.raises(ValidationError):
            Conversation(
                id="conv_001",
                participants=[],  # must have at least one
                acts=[]
            )
            
        # Invalid conversation ID pattern
        with pytest.raises(ValidationError):
            Conversation(
                id="invalid_id",  # doesn't match pattern
                participants=[Participant(id="p1", type=ParticipantType.HUMAN)],
                acts=[]
            )


class TestUtilityFunctions:
    """Tests for utility functions"""
    
    def test_generate_act_id(self):
        """Test act ID generation"""
        id1 = generate_act_id()
        id2 = generate_act_id()
        
        # Check format
        assert id1.startswith("act_")
        assert id2.startswith("act_")
        
        # Should be unique
        assert id1 != id2
        
        # Should match pattern
        import re
        pattern = r"^act_[a-zA-Z0-9_-]+$"
        assert re.match(pattern, id1)
        assert re.match(pattern, id2)
    
    def test_generate_conversation_id(self):
        """Test conversation ID generation"""
        id1 = generate_conversation_id()
        id2 = generate_conversation_id()
        
        # Check format
        assert id1.startswith("conv_")
        assert id2.startswith("conv_")
        
        # Should be unique
        assert id1 != id2
        
        # Should match pattern
        import re
        pattern = r"^conv_[a-zA-Z0-9_-]+$"
        assert re.match(pattern, id1)
        assert re.match(pattern, id2)
    
    def test_create_base_act(self):
        """Test base act creation"""
        base_act = create_base_act(
            speaker="test_speaker",
            act_type=ActType.ASK,
            confidence=0.9,
            source=Source.SYSTEM
        )
        
        assert base_act["speaker"] == "test_speaker"
        assert base_act["type"] == "ask"
        assert base_act["confidence"] == 0.9
        assert base_act["source"] == "system"
        assert "id" in base_act
        assert "timestamp" in base_act


class TestSerialization:
    """Tests for JSON serialization/deserialization"""
    
    def test_act_serialization(self):
        """Test Act serialization to/from JSON"""
        ask = Ask(
            id="act_001",
            timestamp="2025-01-15T14:30:00Z",
            speaker="agent_123",
            type=ActType.ASK,
            field="email",
            prompt="What is your email address?",
            confidence=0.95
        )
        
        # Serialize to JSON
        json_str = ask.model_dump_json()
        data = json.loads(json_str)
        
        assert data["id"] == "act_001"
        assert data["type"] == "ask"
        assert data["confidence"] == 0.95
        
        # Deserialize from JSON
        ask_loaded = Ask.model_validate_json(json_str)
        assert ask_loaded.id == ask.id
        assert ask_loaded.type == ask.type
        assert ask_loaded.confidence == ask.confidence
    
    def test_conversation_serialization(self):
        """Test complete Conversation serialization"""
        agent = Participant(id="agent_001", type=ParticipantType.AI)
        customer = Participant(id="customer_123", type=ParticipantType.HUMAN)
        
        ask = Ask(
            id="act_001",
            timestamp="2025-01-15T14:30:00Z",
            speaker="agent_001",
            type=ActType.ASK,
            field="email",
            prompt="What is your email?"
        )
        
        fact = Fact(
            id="act_002",
            timestamp="2025-01-15T14:31:00Z",
            speaker="customer_123",
            type=ActType.FACT,
            entity="customer_123",
            field="email",
            value="user@example.com"
        )
        
        conversation = Conversation(
            id="conv_001",
            participants=[agent, customer],
            acts=[ask, fact],
            status=ConversationStatus.ACTIVE
        )
        
        # Serialize
        json_str = conversation.model_dump_json()
        
        # Deserialize
        conversation_loaded = Conversation.model_validate_json(json_str)
        
        assert conversation_loaded.id == conversation.id
        assert len(conversation_loaded.participants) == 2
        assert len(conversation_loaded.acts) == 2
        assert isinstance(conversation_loaded.acts[0], Ask)
        assert isinstance(conversation_loaded.acts[1], Fact)


class TestSchemas:
    """Tests for JSON schemas"""
    
    def test_schemas_available(self):
        """Test that all schemas are available"""
        expected_schemas = [
            "act", "ask", "fact", "confirm", "commit", "error",
            "entity", "participant", "constraint", "conversation"
        ]
        
        for schema_name in expected_schemas:
            assert schema_name in SCHEMAS
            schema = SCHEMAS[schema_name]
            assert "$schema" in schema
            assert "$id" in schema
            assert "title" in schema
            assert "description" in schema
    
    def test_schema_structure(self):
        """Test schema structure"""
        act_schema = SCHEMAS["act"]
        
        assert act_schema["$schema"] == "https://json-schema.org/draft/2020-12/schema"
        assert act_schema["$id"] == "https://schemas.astra.dev/v1/act.json"
        assert act_schema["title"] == "Act"
        assert act_schema["type"] == "object"
        assert "id" in act_schema["required"]
        assert "timestamp" in act_schema["required"]
        assert "speaker" in act_schema["required"]
        assert "type" in act_schema["required"]


class TestConstants:
    """Tests for package constants"""
    
    def test_version_constants(self):
        """Test version information"""
        assert __version__ == "1.0.0"
        assert __schema_version__ == "v1"


class TestPydanticFeatures:
    """Tests for Pydantic-specific features"""
    
    def test_model_validation(self):
        """Test Pydantic model validation"""
        # Valid data
        data = {
            "id": "act_001",
            "timestamp": "2025-01-15T14:30:00Z",
            "speaker": "agent_123",
            "type": "ask",
            "field": "email",
            "prompt": "What is your email?"
        }
        
        ask = Ask.model_validate(data)
        assert ask.id == "act_001"
        assert ask.type == ActType.ASK
        
        # Invalid data - should raise ValidationError
        invalid_data = {
            "id": "act_001",
            "timestamp": "2025-01-15T14:30:00Z",
            "speaker": "agent_123",
            "type": "invalid_type",  # invalid enum value
            "field": "email",
            "prompt": "What is your email?"
        }
        
        with pytest.raises(ValidationError) as exc_info:
            Ask.model_validate(invalid_data)
        
        # Should mention the invalid enum value
        assert "invalid_type" in str(exc_info.value)
    
    def test_model_dump(self):
        """Test model dumping to dict/JSON"""
        ask = Ask(
            id="act_001",
            timestamp="2025-01-15T14:30:00Z",
            speaker="agent_123",
            type=ActType.ASK,
            field="email",
            prompt="What is your email?"
        )
        
        # Dump to dict
        data = ask.model_dump()
        assert isinstance(data, dict)
        assert data["type"] == "ask"
        
        # Dump to JSON
        json_str = ask.model_dump_json()
        assert isinstance(json_str, str)
        
        # Should be valid JSON
        parsed = json.loads(json_str)
        assert parsed["id"] == "act_001"
    
    def test_model_copy(self):
        """Test model copying with updates"""
        ask = Ask(
            id="act_001",
            timestamp="2025-01-15T14:30:00Z",
            speaker="agent_123",
            type=ActType.ASK,
            field="email",
            prompt="What is your email?"
        )
        
        # Copy with updates
        ask_copy = ask.model_copy(update={"prompt": "Please provide your email address"})
        
        # Original unchanged
        assert ask.prompt == "What is your email?"
        
        # Copy has updated value
        assert ask_copy.prompt == "Please provide your email address"
        
        # Other fields same
        assert ask_copy.id == ask.id
        assert ask_copy.speaker == ask.speaker


if __name__ == "__main__":
    pytest.main([__file__])
