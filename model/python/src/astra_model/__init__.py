"""
ASTRA Model Python Package

This package provides Python types and JSON schemas for ASTRA (Act State Representation Architecture).

ASTRA provides shared structure so different conversational runtimes and applications can interoperate. 
It defines the canonical data model for representing conversational state as typed, auditable sequences of actions.
"""

from .types import (
    # Base types
    Act,
    ActType,
    Source,
    ActMetadata,
    
    # Entity types
    Entity,
    EntityRef,
    
    # Constraint types
    Constraint,
    ConstraintType,
    FormatType,
    RangeConstraint,
    
    # Participant types
    Participant,
    ParticipantType,
    ParticipantPreferences,
    
    # Act types
    Ask,
    ExpectedType,
    Fact,
    FieldOperation,
    ValidationStatus,
    Confirm,
    ConfirmationMethod,
    Commit,
    CommitAction,
    CommitStatus,
    CommitError,
    Error,
    ErrorSeverity,
    ErrorCategory,
    SuggestedAction,
    
    # Conversation types
    ConversationAct,
    ConversationStatus,
    ConversationContext,
    ConversationMetadata,
    Conversation,
    
    # Utility functions
    generate_act_id,
    generate_conversation_id,
    create_base_act,
)

from .schemas import SCHEMAS

# Package metadata
__version__ = "1.0.0"
__schema_version__ = "v1"
__author__ = "Pryszm"
__license__ = "Apache-2.0"
__description__ = "Python types and JSON schemas for ASTRA conversations"

# Export all public APIs
__all__ = [
    # Base types
    "Act",
    "ActType", 
    "Source",
    "ActMetadata",
    
    # Entity types
    "Entity",
    "EntityRef",
    
    # Constraint types
    "Constraint",
    "ConstraintType",
    "FormatType", 
    "RangeConstraint",
    
    # Participant types
    "Participant",
    "ParticipantType",
    "ParticipantPreferences",
    
    # Act types
    "Ask",
    "ExpectedType",
    "Fact",
    "FieldOperation",
    "ValidationStatus",
    "Confirm",
    "ConfirmationMethod",
    "Commit",
    "CommitAction",
    "CommitStatus",
    "CommitError", 
    "Error",
    "ErrorSeverity",
    "ErrorCategory",
    "SuggestedAction",
    
    # Conversation types
    "ConversationAct",
    "ConversationStatus",
    "ConversationContext",
    "ConversationMetadata",
    "Conversation",
    
    # Utility functions
    "generate_act_id",
    "generate_conversation_id",
    "create_base_act",
    
    # Schemas
    "SCHEMAS",
    
    # Package metadata
    "__version__",
    "__schema_version__",
    "__author__",
    "__license__",
    "__description__",
]
