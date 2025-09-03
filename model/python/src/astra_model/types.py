"""
ASTRA (Act State Representation Architecture) Python Types

This module provides Python types and Pydantic models for ASTRA conversations.

ASTRA provides shared structure so different conversational runtimes and applications can interoperate. 
It defines the canonical data model for representing conversational state as typed, auditable sequences of actions.
"""

import uuid
from datetime import datetime
from typing import Any, Dict, List, Optional, Union
from enum import Enum

from pydantic import BaseModel, Field, ConfigDict


# ============================================================================
# Base Types
# ============================================================================

class Source(str, Enum):
    """Source that generated this act"""
    HUMAN = "human"
    SPEECH_RECOGNITION = "speech_recognition"
    TEXT_ANALYSIS = "text_analysis"
    SYSTEM = "system"
    AI = "ai"


class ActType(str, Enum):
    """Type of conversational act being performed"""
    ASK = "ask"
    FACT = "fact"
    CONFIRM = "confirm"
    COMMIT = "commit"
    ERROR = "error"


class ActMetadata(BaseModel):
    """Additional metadata for acts"""
    model_config = ConfigDict(extra="allow")
    
    channel: Optional[str] = Field(None, description="Communication channel (voice, text, email, etc.)")
    language: Optional[str] = Field(None, description="Language code (ISO 639-1, optional region)")
    original_text: Optional[str] = Field(None, description="Original utterance that generated this act")
    processing_time_ms: Optional[int] = Field(None, description="Time taken to process this act in milliseconds")


class Act(BaseModel):
    """Base type for all conversational actions in ASTRA"""
    model_config = ConfigDict(extra="forbid")
    
    id: str = Field(..., description="Unique identifier for this act within the conversation")
    timestamp: str = Field(..., description="ISO 8601 timestamp when the act occurred")
    speaker: str = Field(..., description="Identifier of the conversation participant who performed this act")
    type: ActType = Field(..., description="Type of conversational act being performed")
    confidence: Optional[float] = Field(None, ge=0.0, le=1.0, description="Confidence score for automated act extraction (0.0 to 1.0)")
    source: Optional[Source] = Field(None, description="Source that generated this act")
    metadata: Optional[ActMetadata] = Field(None, description="Additional context-specific metadata")


# ============================================================================
# Entity Types
# ============================================================================

class Entity(BaseModel):
    """Reference to a business entity in ASTRA conversations"""
    model_config = ConfigDict(extra="allow")
    
    id: str = Field(..., description="Unique identifier for this entity within the conversation scope")
    type: str = Field(..., description="Type of business entity (order, customer, appointment, ticket, etc.)")
    external_id: Optional[str] = Field(None, description="External system identifier for this entity")
    system: Optional[str] = Field(None, description="External system that owns this entity")
    version: Optional[str] = Field(None, description="Version or revision of this entity")
    schema_url: Optional[str] = Field(None, description="URL to the schema definition for this entity type")
    metadata: Optional[Dict[str, Any]] = Field(None, description="Additional entity-specific metadata")


# Entity reference that can be either a string ID or structured Entity
EntityRef = Union[str, Entity]


# ============================================================================
# Constraint Types
# ============================================================================

class ConstraintType(str, Enum):
    """Type of constraint being applied"""
    REQUIRED = "required"
    OPTIONAL = "optional"
    MIN_LENGTH = "min_length"
    MAX_LENGTH = "max_length"
    PATTERN = "pattern"
    FORMAT = "format"
    RANGE = "range"
    ENUM = "enum"
    CUSTOM = "custom"


class FormatType(str, Enum):
    """Format validation types"""
    EMAIL = "email"
    PHONE = "phone"
    URL = "url"
    DATE = "date"
    TIME = "time"
    DATETIME = "datetime"
    UUID = "uuid"
    IPV4 = "ipv4"
    IPV6 = "ipv6"


class RangeConstraint(BaseModel):
    """Range constraint value"""
    min: Optional[float] = Field(None, description="Minimum value")
    max: Optional[float] = Field(None, description="Maximum value")
    inclusive: Optional[bool] = Field(True, description="Whether the range is inclusive")


class Constraint(BaseModel):
    """Validation constraint for ASTRA fields and values"""
    type: ConstraintType = Field(..., description="Type of constraint being applied")
    value: Optional[Union[int, float, str, FormatType, RangeConstraint, List[str]]] = Field(
        None, description="Constraint value (varies by constraint type)"
    )
    message: Optional[str] = Field(None, description="Human-readable error message when constraint is violated")
    code: Optional[str] = Field(None, description="Machine-readable error code for constraint violations")


# ============================================================================
# Participant Types
# ============================================================================

class ParticipantType(str, Enum):
    """Type of participant"""
    HUMAN = "human"
    AI = "ai"
    SYSTEM = "system"
    BOT = "bot"


class ParticipantPreferences(BaseModel):
    """Participant preferences"""
    model_config = ConfigDict(extra="allow")
    
    language: Optional[str] = Field(None, description="Preferred language code (ISO 639-1 with optional region)")
    timezone: Optional[str] = Field(None, description="Preferred timezone (IANA timezone identifier)")
    communication_channels: Optional[List[str]] = Field(
        None, description="Preferred communication channels in order of preference"
    )


class Participant(BaseModel):
    """Conversation participant in ASTRA conversations"""
    model_config = ConfigDict(extra="allow")
    
    id: str = Field(..., description="Unique identifier for this participant")
    type: ParticipantType = Field(..., description="Type of participant")
    role: Optional[str] = Field(None, description="Business role of the participant (customer, agent, manager, etc.)")
    name: Optional[str] = Field(None, description="Display name of the participant")
    email: Optional[str] = Field(None, description="Email address of the participant")
    phone: Optional[str] = Field(None, description="Phone number of the participant")
    external_id: Optional[str] = Field(None, description="External system identifier for this participant")
    system: Optional[str] = Field(None, description="External system that manages this participant")
    capabilities: Optional[List[str]] = Field(None, description="List of capabilities this participant has")
    permissions: Optional[List[str]] = Field(None, description="List of permissions granted to this participant")
    preferences: Optional[ParticipantPreferences] = Field(None, description="Participant preferences")
    metadata: Optional[Dict[str, Any]] = Field(None, description="Additional participant metadata")


# ============================================================================
# Act Types
# ============================================================================

class ExpectedType(str, Enum):
    """Expected data type of the response"""
    STRING = "string"
    NUMBER = "number"
    BOOLEAN = "boolean"
    OBJECT = "object"
    ARRAY = "array"
    DATE = "date"
    EMAIL = "email"
    PHONE = "phone"
    ADDRESS = "address"


class Ask(Act):
    """Act that requests missing information required to complete a business process"""
    type: ActType = Field(ActType.ASK, description="Type of conversational act")
    field: str = Field(..., description="Field or information being requested")
    prompt: str = Field(..., description="Question or request presented to obtain the information")
    constraints: Optional[List[Constraint]] = Field(None, description="Validation constraints for the requested information")
    required: Optional[bool] = Field(None, description="Whether this information is required to proceed")
    expected_type: Optional[ExpectedType] = Field(None, description="Expected data type of the response")
    retry_count: Optional[int] = Field(None, ge=0, description="Number of times this question has been asked")
    max_retries: Optional[int] = Field(None, ge=0, description="Maximum number of retry attempts before escalation")


class FieldOperation(str, Enum):
    """Operation being performed on the field"""
    SET = "set"
    APPEND = "append"
    INCREMENT = "increment"
    DECREMENT = "decrement"
    DELETE = "delete"
    MERGE = "merge"


class ValidationStatus(str, Enum):
    """Validation status of this fact"""
    PENDING = "pending"
    VALID = "valid"
    INVALID = "invalid"
    PARTIAL = "partial"


class Fact(Act):
    """Act that declares facts or information provided during conversation"""
    type: ActType = Field(ActType.FACT, description="Type of conversational act")
    entity: EntityRef = Field(..., description="Business entity being modified (order, customer, appointment, etc.)")
    field: str = Field(..., description="Specific field or property being set")
    value: Any = Field(..., description="Value being assigned to the field")
    operation: Optional[FieldOperation] = Field(None, description="Operation being performed on the field")
    previous_value: Optional[Any] = Field(None, description="Previous value of the field (for audit trail)")
    validation_status: Optional[ValidationStatus] = Field(None, description="Validation status of this fact")
    validation_errors: Optional[List[str]] = Field(
        None, description="List of validation errors if validation_status is invalid"
    )


class ConfirmationMethod(str, Enum):
    """How the confirmation was obtained"""
    VERBAL = "verbal"
    EXPLICIT = "explicit"
    IMPLICIT = "implicit"
    TIMEOUT = "timeout"
    SYSTEM = "system"


class Confirm(Act):
    """Act that verifies understanding of information before commitment"""
    type: ActType = Field(ActType.CONFIRM, description="Type of conversational act")
    entity: EntityRef = Field(..., description="Business entity being confirmed")
    summary: str = Field(..., description="Human-readable summary of what is being confirmed")
    awaiting: Optional[bool] = Field(None, description="Whether confirmation is still pending")
    confirmed: Optional[bool] = Field(None, description="Whether the confirmation was accepted (true) or rejected (false)")
    confirmation_method: Optional[ConfirmationMethod] = Field(None, description="How the confirmation was obtained")
    fields_confirmed: Optional[List[str]] = Field(None, description="Specific fields or aspects being confirmed")
    rejection_reason: Optional[str] = Field(None, description="Reason provided if confirmation was rejected")
    timeout_ms: Optional[int] = Field(None, ge=0, description="Timeout for awaiting confirmation in milliseconds")


class CommitAction(str, Enum):
    """Action being performed in the target system"""
    CREATE = "create"
    UPDATE = "update"
    DELETE = "delete"
    EXECUTE = "execute"
    CANCEL = "cancel"
    PAUSE = "pause"
    RESUME = "resume"


class CommitStatus(str, Enum):
    """Status of the commit operation"""
    PENDING = "pending"
    IN_PROGRESS = "in_progress"
    SUCCESS = "success"
    FAILED = "failed"
    RETRYING = "retrying"
    CANCELLED = "cancelled"


class CommitError(BaseModel):
    """Error information for failed commits"""
    code: str = Field(..., description="Error code from the target system")
    message: str = Field(..., description="Human-readable error message")
    details: Optional[Dict[str, Any]] = Field(None, description="Additional error context")
    recoverable: bool = Field(..., description="Whether the error can be recovered from")


class Commit(Act):
    """Act that executes business processes and triggers system integrations"""
    type: ActType = Field(ActType.COMMIT, description="Type of conversational act")
    entity: EntityRef = Field(..., description="Business entity being committed to external systems")
    action: CommitAction = Field(..., description="Action being performed in the target system")
    system: Optional[str] = Field(None, description="Target system identifier (CRM, order_management, etc.)")
    transaction_id: Optional[str] = Field(None, description="External system transaction or record identifier")
    status: Optional[CommitStatus] = Field(None, description="Status of the commit operation")
    error: Optional[CommitError] = Field(None, description="Error information if commit failed")
    retry_count: Optional[int] = Field(None, ge=0, description="Number of retry attempts made")
    max_retries: Optional[int] = Field(None, ge=0, description="Maximum number of retry attempts")
    idempotency_key: Optional[str] = Field(None, description="Key to ensure idempotent operations")
    rollback_info: Optional[Dict[str, Any]] = Field(
        None, description="Information needed to rollback this commit if necessary"
    )


class ErrorSeverity(str, Enum):
    """Severity level of the error"""
    INFO = "info"
    WARNING = "warning"
    ERROR = "error"
    CRITICAL = "critical"


class ErrorCategory(str, Enum):
    """Category of error for classification"""
    VALIDATION = "validation"
    PROCESSING = "processing"
    INTEGRATION = "integration"
    TIMEOUT = "timeout"
    PERMISSION = "permission"
    SYSTEM = "system"
    USER_INPUT = "user_input"
    BUSINESS_RULE = "business_rule"


class SuggestedAction(str, Enum):
    """Suggested recovery action"""
    RETRY = "retry"
    ESCALATE = "escalate"
    IGNORE = "ignore"
    CLARIFY = "clarify"
    FALLBACK = "fallback"
    TERMINATE = "terminate"


class Error(Act):
    """Act that handles failures and exceptions in conversational processing"""
    type: ActType = Field(ActType.ERROR, description="Type of conversational act")
    code: str = Field(..., description="Machine-readable error code")
    message: str = Field(..., description="Human-readable error message")
    recoverable: bool = Field(..., description="Whether the conversation can continue after this error")
    severity: Optional[ErrorSeverity] = Field(None, description="Severity level of the error")
    category: Optional[ErrorCategory] = Field(None, description="Category of error for classification")
    details: Optional[Dict[str, Any]] = Field(None, description="Additional error context and debugging information")
    related_act_id: Optional[str] = Field(None, description="ID of the act that caused this error")
    suggested_action: Optional[SuggestedAction] = Field(None, description="Suggested recovery action")
    user_message: Optional[str] = Field(None, description="User-friendly message to display to conversation participants")
    stack_trace: Optional[str] = Field(None, description="Technical stack trace for debugging (not shown to users)")


# ============================================================================
# Conversation Types
# ============================================================================

# Union type for all possible acts
ConversationAct = Union[Ask, Fact, Confirm, Commit, Error]


class ConversationStatus(str, Enum):
    """Current status of the conversation"""
    ACTIVE = "active"
    PAUSED = "paused"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"


class ConversationContext(BaseModel):
    """Conversation context and session information"""
    model_config = ConfigDict(extra="allow")
    
    session_id: Optional[str] = Field(None, description="Session identifier")
    user_agent: Optional[str] = Field(None, description="User agent or client information")
    ip_address: Optional[str] = Field(None, description="Client IP address")
    referrer: Optional[str] = Field(None, description="How the conversation was initiated")


class ConversationMetadata(BaseModel):
    """Additional conversation metadata"""
    model_config = ConfigDict(extra="allow")
    
    total_duration_ms: Optional[int] = Field(None, ge=0, description="Total conversation duration in milliseconds")
    act_count: Optional[int] = Field(None, ge=0, description="Total number of acts in the conversation")
    error_count: Optional[int] = Field(None, ge=0, description="Number of errors that occurred")
    commit_count: Optional[int] = Field(None, ge=0, description="Number of successful commits")
    avg_confidence: Optional[float] = Field(None, ge=0.0, le=1.0, description="Average confidence score across all acts")


class Conversation(BaseModel):
    """Complete ASTRA conversation container with acts and metadata"""
    model_config = ConfigDict(extra="forbid")
    
    id: str = Field(..., description="Unique identifier for this conversation")
    participants: List[Participant] = Field(..., description="List of conversation participants")
    acts: List[ConversationAct] = Field(..., description="Ordered sequence of acts in this conversation")
    started_at: Optional[str] = Field(None, description="When the conversation started")
    ended_at: Optional[str] = Field(None, description="When the conversation ended")
    status: Optional[ConversationStatus] = Field(None, description="Current status of the conversation")
    channel: Optional[str] = Field(None, description="Primary communication channel for this conversation")
    business_schema: Optional[str] = Field(None, description="Business schema identifier used for this conversation")
    context: Optional[ConversationContext] = Field(None, description="Conversation context and session information")
    final_state: Optional[Dict[str, Any]] = Field(
        None, description="Final computed state of all entities after processing all acts"
    )
    metadata: Optional[ConversationMetadata] = Field(None, description="Additional conversation metadata")


# ============================================================================
# Utility Functions
# ============================================================================

def generate_act_id() -> str:
    """Generate a unique act ID"""
    # Generate a UUID and format it as act_<short_uuid>
    uid = str(uuid.uuid4()).replace('-', '')[:12]
    return f"act_{uid}"


def generate_conversation_id() -> str:
    """Generate a unique conversation ID"""
    # Generate a UUID and format it as conv_<short_uuid>
    uid = str(uuid.uuid4()).replace('-', '')[:12]
    return f"conv_{uid}"


def create_base_act(
    speaker: str,
    act_type: ActType,
    additional_fields: Optional[Dict[str, Any]] = None
) -> Dict[str, Any]:
    """Create a base act with required fields and optional additional fields"""
    base_act = {
        "id": generate_act_id(),
        "timestamp": datetime.now().isoformat() + "Z",
        "speaker": speaker,
        "type": act_type,
    }
    
    if additional_fields:
        base_act.update(additional_fields)
    
    return base_act
