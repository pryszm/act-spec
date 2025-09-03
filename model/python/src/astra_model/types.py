"""
ASTRA Type Definitions

This module provides Python type definitions for ASTRA conversations.
"""

from typing import Any, Dict, List, Optional, Union
from enum import Enum

# Base types
class ActType(str, Enum):
    ASK = "ask"
    FACT = "fact"
    CONFIRM = "confirm"
    COMMIT = "commit"
    ERROR = "error"

class Source(str, Enum):
    HUMAN = "human"
    SPEECH_RECOGNITION = "speech_recognition"
    TEXT_ANALYSIS = "text_analysis"
    SYSTEM = "system"
    AI = "ai"

# Placeholder classes for types that would be fully implemented
class Act:
    """Base type for all conversational actions"""
    pass

class ActMetadata:
    """Metadata for acts"""
    pass

class Entity:
    """Business entity reference"""
    pass

class EntityRef:
    """Entity reference"""
    pass

class Constraint:
    """Validation constraint"""
    pass

class ConstraintType(str, Enum):
    REQUIRED = "required"
    OPTIONAL = "optional"
    MIN_LENGTH = "min_length"
    MAX_LENGTH = "max_length"
    PATTERN = "pattern"
    FORMAT = "format"
    RANGE = "range"
    ENUM = "enum"
    CUSTOM = "custom"

class FormatType:
    """Format type for constraints"""
    pass

class RangeConstraint:
    """Range constraint"""
    pass

class Participant:
    """Conversation participant"""
    pass

class ParticipantType(str, Enum):
    HUMAN = "human"
    AI = "ai"
    SYSTEM = "system"
    BOT = "bot"

class ParticipantPreferences:
    """Participant preferences"""
    pass

class Ask:
    """Ask act type"""
    pass

class ExpectedType(str, Enum):
    STRING = "string"
    NUMBER = "number"
    BOOLEAN = "boolean"
    OBJECT = "object"
    ARRAY = "array"
    DATE = "date"
    EMAIL = "email"
    PHONE = "phone"
    ADDRESS = "address"

class Fact:
    """Fact act type"""
    pass

class FieldOperation(str, Enum):
    SET = "set"
    APPEND = "append"
    INCREMENT = "increment"
    DECREMENT = "decrement"
    DELETE = "delete"
    MERGE = "merge"

class ValidationStatus(str, Enum):
    PENDING = "pending"
    VALID = "valid"
    INVALID = "invalid"
    PARTIAL = "partial"

class Confirm:
    """Confirm act type"""
    pass

class ConfirmationMethod(str, Enum):
    VERBAL = "verbal"
    EXPLICIT = "explicit"
    IMPLICIT = "implicit"
    TIMEOUT = "timeout"
    SYSTEM = "system"

class Commit:
    """Commit act type"""
    pass

class CommitAction(str, Enum):
    CREATE = "create"
    UPDATE = "update"
    DELETE = "delete"
    EXECUTE = "execute"
    CANCEL = "cancel"
    PAUSE = "pause"
    RESUME = "resume"

class CommitStatus(str, Enum):
    PENDING = "pending"
    IN_PROGRESS = "in_progress"
    SUCCESS = "success"
    FAILED = "failed"
    RETRYING = "retrying"
    CANCELLED = "cancelled"

class CommitError:
    """Commit error information"""
    pass

class Error:
    """Error act type"""
    pass

class ErrorSeverity(str, Enum):
    INFO = "info"
    WARNING = "warning"
    ERROR = "error"
    CRITICAL = "critical"

class ErrorCategory(str, Enum):
    VALIDATION = "validation"
    PROCESSING = "processing"
    INTEGRATION = "integration"
    TIMEOUT = "timeout"
    PERMISSION = "permission"
    SYSTEM = "system"
    USER_INPUT = "user_input"
    BUSINESS_RULE = "business_rule"

class SuggestedAction(str, Enum):
    RETRY = "retry"
    ESCALATE = "escalate"
    IGNORE = "ignore"
    CLARIFY = "clarify"
    FALLBACK = "fallback"
    TERMINATE = "terminate"

class ConversationAct:
    """Conversation act wrapper"""
    pass

class ConversationStatus(str, Enum):
    ACTIVE = "active"
    PAUSED = "paused"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"

class ConversationContext:
    """Conversation context"""
    pass

class ConversationMetadata:
    """Conversation metadata"""
    pass

class Conversation:
    """Complete conversation container"""
    pass

# Utility functions
def generate_act_id() -> str:
    """Generate a unique act ID"""
    import uuid
    return f"act_{uuid.uuid4().hex[:12]}"

def generate_conversation_id() -> str:
    """Generate a unique conversation ID"""
    import uuid
    return f"conv_{uuid.uuid4().hex[:12]}"

def create_base_act(speaker: str, act_type: str) -> Dict[str, Any]:
    """Create a base act structure"""
    from datetime import datetime
    return {
        "id": generate_act_id(),
        "timestamp": datetime.utcnow().isoformat() + "Z",
        "speaker": speaker,
        "type": act_type
    }
