"""
Tests for ASTRA Python types

Basic tests to validate the type implementations work correctly.
"""

import json
from datetime import datetime
from astra_model import (
    Act, Ask, Fact, Confirm, Commit, Error,
    ActType, Source, ParticipantType, ConversationStatus,
    Participant, Entity, Conversation,
    generate_act_id, generate_conversation_id, create_base_act,
    SCHEMAS
)


def test_basic_imports():
    """Test that all types can be imported"""
    assert Act is not None
    assert Ask is not None
    assert Fact is not None
    assert ActType is not None
    assert Source is not None
    print("✓ All imports successful")


def test_id_generation():
    """Test ID generation functions"""
    act_id = generate_act_id()
    conv_id = generate_conversation_id()
    
    assert act_id.startswith("act_")
    assert conv_id.startswith("conv_")
    assert len(act_id) > 4
    assert len(conv_id) > 5
    print(f"✓ Generated IDs: {act_id}, {conv_id}")


def test_ask_creation():
    """Test Ask act creation"""
    ask = Ask(
        id=generate_act_id(),
        timestamp=datetime.now().isoformat() + "Z",
        speaker="agent_123",
        type=ActType.ASK,
        field="email",
        prompt="What is your email address?",
        required=True
    )
    
    assert ask.type == ActType.ASK
    assert ask.field == "email"
    assert ask.required is True
    print(f"✓ Ask created: {ask.id}")


def test_fact_creation():
    """Test Fact act creation"""
    fact = Fact(
        id=generate_act_id(),
        timestamp=datetime.now().isoformat() + "Z",
        speaker="user_456",
        type=ActType.FACT,
        entity="customer_789",
        field="email",
        value="user@example.com"
    )
    
    assert fact.type == ActType.FACT
    assert fact.entity == "customer_789"
    assert fact.value == "user@example.com"
    print(f"✓ Fact created: {fact.id}")


def test_conversation_creation():
    """Test Conversation creation"""
    participants = [
        Participant(
            id="p1",
            type=ParticipantType.HUMAN,
            name="John Doe",
            email="john@example.com"
        ),
        Participant(
            id="p2",
            type=ParticipantType.AI,
            name="Assistant"
        )
    ]
    
    conversation = Conversation(
        id=generate_conversation_id(),
        participants=participants,
        acts=[],
        status=ConversationStatus.ACTIVE
    )
    
    assert len(conversation.participants) == 2
    assert conversation.status == ConversationStatus.ACTIVE
    print(f"✓ Conversation created: {conversation.id}")


def test_json_serialization():
    """Test JSON serialization/deserialization"""
    ask = Ask(
        id="ask_123",
        timestamp="2025-01-15T14:30:00Z",
        speaker="agent",
        type=ActType.ASK,
        field="name",
        prompt="What is your name?"
    )
    
    # Serialize to JSON
    json_str = ask.model_dump_json()
    json_data = json.loads(json_str)
    
    assert json_data["id"] == "ask_123"
    assert json_data["type"] == "ask"
    assert json_data["field"] == "name"
    
    # Deserialize from JSON
    ask_restored = Ask.model_validate_json(json_str)
    assert ask_restored.id == ask.id
    assert ask_restored.field == ask.field
    print("✓ JSON serialization/deserialization works")


def test_schemas():
    """Test that schemas are available"""
    assert "act" in SCHEMAS
    assert "ask" in SCHEMAS
    assert "fact" in SCHEMAS
    assert "conversation" in SCHEMAS
    
    act_schema = SCHEMAS["act"]
    assert "$schema" in act_schema
    assert "$id" in act_schema
    assert "title" in act_schema
    print(f"✓ Found {len(SCHEMAS)} schemas")


def test_base_act_utility():
    """Test create_base_act utility"""
    base_act = create_base_act("speaker_123", ActType.ASK, {"confidence": 0.9})
    
    assert base_act["speaker"] == "speaker_123"
    assert base_act["type"] == ActType.ASK
    assert base_act["confidence"] == 0.9
    assert "id" in base_act
    assert "timestamp" in base_act
    print("✓ Base act utility works")


if __name__ == "__main__":
    print("Running ASTRA Python types tests...")
    
    test_basic_imports()
    test_id_generation()
    test_ask_creation()
    test_fact_creation()
    test_conversation_creation()
    test_json_serialization()
    test_schemas()
    test_base_act_utility()
    
    print("\n🎉 All tests passed!")
