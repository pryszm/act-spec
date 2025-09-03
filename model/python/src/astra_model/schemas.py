"""
ASTRA JSON Schemas

This module provides JSON Schema definitions for all ASTRA types.
These schemas can be used for validation, documentation, and API specification.
"""

from typing import Dict, Any

from .types import (
    Act, Ask, Fact, Confirm, Commit, Error,
    Entity, Participant, Constraint, Conversation,
    ActMetadata, ParticipantPreferences, RangeConstraint,
    CommitError, ConversationContext, ConversationMetadata
)


def _get_schema_with_id(model_class, schema_id: str, title: str) -> Dict[str, Any]:
    """Get JSON schema for a Pydantic model with custom $id and title"""
    schema = model_class.model_json_schema()
    schema["$schema"] = "https://json-schema.org/draft/2020-12/schema"
    schema["$id"] = f"https://schemas.astra.dev/v1/{schema_id}.json"
    schema["title"] = title
    return schema


# Generate JSON schemas for all ASTRA types
SCHEMAS: Dict[str, Dict[str, Any]] = {
    # Base types
    "act": _get_schema_with_id(Act, "act", "Act"),
    "act_metadata": _get_schema_with_id(ActMetadata, "act_metadata", "Act Metadata"),
    
    # Entity types
    "entity": _get_schema_with_id(Entity, "entity", "Entity"),
    
    # Constraint types
    "constraint": _get_schema_with_id(Constraint, "constraint", "Constraint"),
    "range_constraint": _get_schema_with_id(RangeConstraint, "range_constraint", "Range Constraint"),
    
    # Participant types
    "participant": _get_schema_with_id(Participant, "participant", "Participant"),
    "participant_preferences": _get_schema_with_id(ParticipantPreferences, "participant_preferences", "Participant Preferences"),
    
    # Act types
    "ask": _get_schema_with_id(Ask, "ask", "Ask"),
    "fact": _get_schema_with_id(Fact, "fact", "Fact"),
    "confirm": _get_schema_with_id(Confirm, "confirm", "Confirm"),
    "commit": _get_schema_with_id(Commit, "commit", "Commit"),
    "error": _get_schema_with_id(Error, "error", "Error"),
    "commit_error": _get_schema_with_id(CommitError, "commit_error", "Commit Error"),
    
    # Conversation types
    "conversation": _get_schema_with_id(Conversation, "conversation", "Conversation"),
    "conversation_context": _get_schema_with_id(ConversationContext, "conversation_context", "Conversation Context"),
    "conversation_metadata": _get_schema_with_id(ConversationMetadata, "conversation_metadata", "Conversation Metadata"),
}
