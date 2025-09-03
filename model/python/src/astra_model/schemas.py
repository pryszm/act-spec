"""
JSON Schemas for ASTRA types

These schemas are embedded versions of the JSON Schema definitions
from the main ASTRA repository for runtime validation.
"""

SCHEMAS = {
    "act": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/act.json",
        "title": "Act",
        "description": "Base type for all conversational actions in ASTRA",
        "type": "object",
        "required": ["id", "timestamp", "speaker", "type"],
        "properties": {
            "id": {
                "type": "string",
                "pattern": "^act_[a-zA-Z0-9_-]+$",
                "description": "Unique identifier for this act within the conversation"
            },
            "timestamp": {
                "type": "string",
                "format": "date-time",
                "description": "ISO 8601 timestamp when the act occurred"
            },
            "speaker": {
                "type": "string",
                "description": "Identifier of the conversation participant who performed this act"
            },
            "type": {
                "type": "string",
                "enum": ["ask", "fact", "confirm", "commit", "error"],
                "description": "Type of conversational act being performed"
            },
            "confidence": {
                "type": "number",
                "minimum": 0.0,
                "maximum": 1.0,
                "description": "Confidence score for automated act extraction (0.0 to 1.0)"
            },
            "source": {
                "type": "string",
                "enum": ["human", "speech_recognition", "text_analysis", "system", "ai"],
                "description": "Source that generated this act"
            },
            "metadata": {
                "type": "object",
                "description": "Additional context-specific metadata",
                "properties": {
                    "channel": {
                        "type": "string",
                        "description": "Communication channel (voice, text, email, etc.)"
                    },
                    "language": {
                        "type": "string",
                        "pattern": "^[a-z]{2}(-[A-Z]{2})?$",
                        "description": "Language code (ISO 639-1, optional region)"
                    },
                    "original_text": {
                        "type": "string",
                        "description": "Original utterance that generated this act"
                    },
                    "processing_time_ms": {
                        "type": "number",
                        "minimum": 0,
                        "description": "Time taken to process this act in milliseconds"
                    }
                },
                "additionalProperties": True
            }
        },
        "additionalProperties": False
    },

    "ask": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/ask.json",
        "title": "Ask",
        "description": "Act that requests missing information required to complete a business process",
        "type": "object",
        "required": ["id", "timestamp", "speaker", "type", "field", "prompt"],
        "properties": {
            "id": {
                "type": "string",
                "pattern": "^act_[a-zA-Z0-9_-]+$",
                "description": "Unique identifier for this act within the conversation"
            },
            "timestamp": {
                "type": "string",
                "format": "date-time",
                "description": "ISO 8601 timestamp when the act occurred"
            },
            "speaker": {
                "type": "string",
                "description": "Identifier of the conversation participant who performed this act"
            },
            "type": {
                "const": "ask"
            },
            "confidence": {
                "type": "number",
                "minimum": 0.0,
                "maximum": 1.0,
                "description": "Confidence score for automated act extraction (0.0 to 1.0)"
            },
            "source": {
                "type": "string",
                "enum": ["human", "speech_recognition", "text_analysis", "system", "ai"],
                "description": "Source that generated this act"
            },
            "metadata": {
                "type": "object",
                "description": "Additional context-specific metadata",
                "additionalProperties": True
            },
            "field": {
                "type": "string",
                "description": "Field or information being requested"
            },
            "prompt": {
                "type": "string",
                "description": "Question or request presented to obtain the information"
            },
            "constraints": {
                "type": "array",
                "items": {
                    "type": "object",
                    "required": ["type"],
                    "properties": {
                        "type": {
                            "type": "string",
                            "enum": ["required", "optional", "min_length", "max_length", "pattern", "format", "range", "enum", "custom"]
                        },
                        "value": {},
                        "message": {"type": "string"},
                        "code": {"type": "string"}
                    }
                },
                "description": "Validation constraints for the requested information"
            },
            "required": {
                "type": "boolean",
                "default": True,
                "description": "Whether this information is required to proceed"
            },
            "expected_type": {
                "type": "string",
                "enum": ["string", "number", "boolean", "object", "array", "date", "email", "phone", "address"],
                "description": "Expected data type of the response"
            },
            "retry_count": {
                "type": "integer",
                "minimum": 0,
                "default": 0,
                "description": "Number of times this question has been asked"
            },
            "max_retries": {
                "type": "integer",
                "minimum": 0,
                "default": 3,
                "description": "Maximum number of retry attempts before escalation"
            }
        },
        "additionalProperties": False
    },

    "fact": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/fact.json",
        "title": "Fact",
        "description": "Act that declares facts or information provided during conversation",
        "type": "object",
        "required": ["id", "timestamp", "speaker", "type", "entity", "field", "value"],
        "properties": {
            "id": {
                "type": "string",
                "pattern": "^act_[a-zA-Z0-9_-]+$",
                "description": "Unique identifier for this act within the conversation"
            },
            "timestamp": {
                "type": "string",
                "format": "date-time",
                "description": "ISO 8601 timestamp when the act occurred"
            },
            "speaker": {
                "type": "string",
                "description": "Identifier of the conversation participant who performed this act"
            },
            "type": {
                "const": "fact"
            },
            "confidence": {
                "type": "number",
                "minimum": 0.0,
                "maximum": 1.0,
                "description": "Confidence score for automated act extraction (0.0 to 1.0)"
            },
            "source": {
                "type": "string",
                "enum": ["human", "speech_recognition", "text_analysis", "system", "ai"],
                "description": "Source that generated this act"
            },
            "metadata": {
                "type": "object",
                "description": "Additional context-specific metadata",
                "additionalProperties": True
            },
            "entity": {
                "oneOf": [
                    {
                        "type": "string",
                        "description": "Entity identifier as string"
                    },
                    {
                        "type": "object",
                        "required": ["id", "type"],
                        "properties": {
                            "id": {"type": "string"},
                            "type": {"type": "string"},
                            "external_id": {"type": "string"},
                            "system": {"type": "string"},
                            "version": {"type": "string"},
                            "schema_url": {"type": "string", "format": "uri"},
                            "metadata": {"type": "object", "additionalProperties": True}
                        },
                        "additionalProperties": False,
                        "description": "Structured entity reference"
                    }
                ],
                "description": "Business entity being modified (order, customer, appointment, etc.)"
            },
            "field": {
                "type": "string",
                "description": "Specific field or property being set"
            },
            "value": {
                "description": "Value being assigned to the field (any JSON type)"
            },
            "operation": {
                "type": "string",
                "enum": ["set", "append", "increment", "decrement", "delete", "merge"],
                "default": "set",
                "description": "Operation being performed on the field"
            },
            "previous_value": {
                "description": "Previous value of the field (for audit trail)"
            },
            "validation_status": {
                "type": "string",
                "enum": ["pending", "valid", "invalid", "partial"],
                "default": "pending",
                "description": "Validation status of this fact"
            },
            "validation_errors": {
                "type": "array",
                "items": {"type": "string"},
                "description": "List of validation errors if validation_status is invalid"
            }
        },
        "additionalProperties": False
    },

    "confirm": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/confirm.json",
        "title": "Confirm",
        "description": "Act that verifies understanding of information before commitment",
        "type": "object",
        "required": ["id", "timestamp", "speaker", "type", "entity", "summary"],
        "properties": {
            "id": {
                "type": "string",
                "pattern": "^act_[a-zA-Z0-9_-]+$",
                "description": "Unique identifier for this act within the conversation"
            },
            "timestamp": {
                "type": "string",
                "format": "date-time",
                "description": "ISO 8601 timestamp when the act occurred"
            },
            "speaker": {
                "type": "string",
                "description": "Identifier of the conversation participant who performed this act"
            },
            "type": {
                "const": "confirm"
            },
            "confidence": {
                "type": "number",
                "minimum": 0.0,
                "maximum": 1.0,
                "description": "Confidence score for automated act extraction (0.0 to 1.0)"
            },
            "source": {
                "type": "string",
                "enum": ["human", "speech_recognition", "text_analysis", "system", "ai"],
                "description": "Source that generated this act"
            },
            "metadata": {
                "type": "object",
                "description": "Additional context-specific metadata",
                "additionalProperties": True
            },
            "entity": {
                "oneOf": [
                    {
                        "type": "string",
                        "description": "Entity identifier as string"
                    },
                    {
                        "type": "object",
                        "required": ["id", "type"],
                        "properties": {
                            "id": {"type": "string"},
                            "type": {"type": "string"},
                            "external_id": {"type": "string"},
                            "system": {"type": "string"},
                            "version": {"type": "string"},
                            "schema_url": {"type": "string", "format": "uri"},
                            "metadata": {"type": "object", "additionalProperties": True}
                        },
                        "additionalProperties": False,
                        "description": "Structured entity reference"
                    }
                ],
                "description": "Business entity being confirmed"
            },
            "summary": {
                "type": "string",
                "description": "Human-readable summary of what is being confirmed"
            },
            "awaiting": {
                "type": "boolean",
                "default": True,
                "description": "Whether confirmation is still pending"
            },
            "confirmed": {
                "type": "boolean",
                "description": "Whether the confirmation was accepted (true) or rejected (false)"
            },
            "confirmation_method": {
                "type": "string",
                "enum": ["verbal", "explicit", "implicit", "timeout", "system"],
                "description": "How the confirmation was obtained"
            },
            "fields_confirmed": {
                "type": "array",
                "items": {"type": "string"},
                "description": "Specific fields or aspects being confirmed"
            },
            "rejection_reason": {
                "type": "string",
                "description": "Reason provided if confirmation was rejected"
            },
            "timeout_ms": {
                "type": "integer",
                "minimum": 0,
                "description": "Timeout for awaiting confirmation in milliseconds"
            }
        },
        "additionalProperties": False
    },

    "commit": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/commit.json",
        "title": "Commit",
        "description": "Act that executes business processes and triggers system integrations",
        "type": "object",
        "required": ["id", "timestamp", "speaker", "type", "entity", "action"],
        "properties": {
            "id": {
                "type": "string",
                "pattern": "^act_[a-zA-Z0-9_-]+$",
                "description": "Unique identifier for this act within the conversation"
            },
            "timestamp": {
                "type": "string",
                "format": "date-time",
                "description": "ISO 8601 timestamp when the act occurred"
            },
            "speaker": {
                "type": "string",
                "description": "Identifier of the conversation participant who performed this act"
            },
            "type": {
                "const": "commit"
            },
            "confidence": {
                "type": "number",
                "minimum": 0.0,
                "maximum": 1.0,
                "description": "Confidence score for automated act extraction (0.0 to 1.0)"
            },
            "source": {
                "type": "string",
                "enum": ["human", "speech_recognition", "text_analysis", "system", "ai"],
                "description": "Source that generated this act"
            },
            "metadata": {
                "type": "object",
                "description": "Additional context-specific metadata",
                "additionalProperties": True
            },
            "entity": {
                "oneOf": [
                    {
                        "type": "string",
                        "description": "Entity identifier as string"
                    },
                    {
                        "type": "object",
                        "required": ["id", "type"],
                        "properties": {
                            "id": {"type": "string"},
                            "type": {"type": "string"},
                            "external_id": {"type": "string"},
                            "system": {"type": "string"},
                            "version": {"type": "string"},
                            "schema_url": {"type": "string", "format": "uri"},
                            "metadata": {"type": "object", "additionalProperties": True}
                        },
                        "additionalProperties": False,
                        "description": "Structured entity reference"
                    }
                ],
                "description": "Business entity being committed to external systems"
            },
            "action": {
                "type": "string",
                "enum": ["create", "update", "delete", "execute", "cancel", "pause", "resume"],
                "description": "Action being performed in the target system"
            },
            "system": {
                "type": "string",
                "description": "Target system identifier (CRM, order_management, etc.)"
            },
            "transaction_id": {
                "type": "string",
                "description": "External system transaction or record identifier"
            },
            "status": {
                "type": "string",
                "enum": ["pending", "in_progress", "success", "failed", "retrying", "cancelled"],
                "default": "pending",
                "description": "Status of the commit operation"
            },
            "error": {
                "type": "object",
                "properties": {
                    "code": {
                        "type": "string",
                        "description": "Error code from the target system"
                    },
                    "message": {
                        "type": "string",
                        "description": "Human-readable error message"
                    },
                    "details": {
                        "type": "object",
                        "description": "Additional error context"
                    },
                    "recoverable": {
                        "type": "boolean",
                        "description": "Whether the error can be recovered from"
                    }
                },
                "required": ["code", "message"],
                "description": "Error information if commit failed"
            },
            "retry_count": {
                "type": "integer",
                "minimum": 0,
                "default": 0,
                "description": "Number of retry attempts made"
            },
            "max_retries": {
                "type": "integer",
                "minimum": 0,
                "default": 3,
                "description": "Maximum number of retry attempts"
            },
            "idempotency_key": {
                "type": "string",
                "description": "Key to ensure idempotent operations"
            },
            "rollback_info": {
                "type": "object",
                "description": "Information needed to rollback this commit if necessary"
            }
        },
        "additionalProperties": False
    },

    "error": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/error.json",
        "title": "Error",
        "description": "Act that handles failures and exceptions in conversational processing",
        "type": "object",
        "required": ["id", "timestamp", "speaker", "type", "code", "message", "recoverable"],
        "properties": {
            "id": {
                "type": "string",
                "pattern": "^act_[a-zA-Z0-9_-]+$",
                "description": "Unique identifier for this act within the conversation"
            },
            "timestamp": {
                "type": "string",
                "format": "date-time",
                "description": "ISO 8601 timestamp when the act occurred"
            },
            "speaker": {
                "type": "string",
                "description": "Identifier of the conversation participant who performed this act"
            },
            "type": {
                "const": "error"
            },
            "confidence": {
                "type": "number",
                "minimum": 0.0,
                "maximum": 1.0,
                "description": "Confidence score for automated act extraction (0.0 to 1.0)"
            },
            "source": {
                "type": "string",
                "enum": ["human", "speech_recognition", "text_analysis", "system", "ai"],
                "description": "Source that generated this act"
            },
            "metadata": {
                "type": "object",
                "description": "Additional context-specific metadata",
                "additionalProperties": True
            },
            "code": {
                "type": "string",
                "description": "Machine-readable error code"
            },
            "message": {
                "type": "string",
                "description": "Human-readable error message"
            },
            "recoverable": {
                "type": "boolean",
                "description": "Whether the conversation can continue after this error"
            },
            "severity": {
                "type": "string",
                "enum": ["info", "warning", "error", "critical"],
                "default": "error",
                "description": "Severity level of the error"
            },
            "category": {
                "type": "string",
                "enum": ["validation", "processing", "integration", "timeout", "permission", "system", "user_input", "business_rule"],
                "description": "Category of error for classification"
            },
            "details": {
                "type": "object",
                "description": "Additional error context and debugging information"
            },
            "related_act_id": {
                "type": "string",
                "pattern": "^act_[a-zA-Z0-9_-]+$",
                "description": "ID of the act that caused this error"
            },
            "suggested_action": {
                "type": "string",
                "enum": ["retry", "escalate", "ignore", "clarify", "fallback", "terminate"],
                "description": "Suggested recovery action"
            },
            "user_message": {
                "type": "string",
                "description": "User-friendly message to display to conversation participants"
            },
            "stack_trace": {
                "type": "string",
                "description": "Technical stack trace for debugging (not shown to users)"
            }
        },
        "additionalProperties": False
    },

    "entity": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/entity.json",
        "title": "Entity",
        "description": "Reference to a business entity in ASTRA conversations",
        "type": "object",
        "required": ["id", "type"],
        "properties": {
            "id": {
                "type": "string",
                "description": "Unique identifier for this entity within the conversation scope"
            },
            "type": {
                "type": "string",
                "description": "Type of business entity (order, customer, appointment, ticket, etc.)"
            },
            "external_id": {
                "type": "string",
                "description": "External system identifier for this entity"
            },
            "system": {
                "type": "string",
                "description": "External system that owns this entity"
            },
            "version": {
                "type": "string",
                "description": "Version or revision of this entity"
            },
            "schema_url": {
                "type": "string",
                "format": "uri",
                "description": "URL to the schema definition for this entity type"
            },
            "metadata": {
                "type": "object",
                "description": "Additional entity-specific metadata",
                "additionalProperties": True
            }
        },
        "additionalProperties": False
    },

    "participant": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/participant.json",
        "title": "Participant",
        "description": "Conversation participant in ASTRA conversations",
        "type": "object",
        "required": ["id", "type"],
        "properties": {
            "id": {
                "type": "string",
                "description": "Unique identifier for this participant"
            },
            "type": {
                "type": "string",
                "enum": ["human", "ai", "system", "bot"],
                "description": "Type of participant"
            },
            "role": {
                "type": "string",
                "description": "Business role of the participant (customer, agent, manager, etc.)"
            },
            "name": {
                "type": "string",
                "description": "Display name of the participant"
            },
            "email": {
                "type": "string",
                "format": "email",
                "description": "Email address of the participant"
            },
            "phone": {
                "type": "string",
                "description": "Phone number of the participant"
            },
            "external_id": {
                "type": "string",
                "description": "External system identifier for this participant"
            },
            "system": {
                "type": "string",
                "description": "External system that manages this participant"
            },
            "capabilities": {
                "type": "array",
                "items": {"type": "string"},
                "description": "List of capabilities this participant has"
            },
            "permissions": {
                "type": "array",
                "items": {"type": "string"},
                "description": "List of permissions granted to this participant"
            },
            "preferences": {
                "type": "object",
                "properties": {
                    "language": {
                        "type": "string",
                        "pattern": "^[a-z]{2}(-[A-Z]{2})?$",
                        "description": "Preferred language code"
                    },
                    "timezone": {
                        "type": "string",
                        "description": "Preferred timezone (IANA timezone identifier)"
                    },
                    "communication_channels": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "Preferred communication channels in order of preference"
                    }
                },
                "additionalProperties": True,
                "description": "Participant preferences"
            },
            "metadata": {
                "type": "object",
                "description": "Additional participant metadata",
                "additionalProperties": True
            }
        },
        "additionalProperties": False
    },

    "constraint": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/constraint.json",
        "title": "Constraint",
        "description": "Validation constraint for ASTRA fields and values",
        "type": "object",
        "required": ["type"],
        "properties": {
            "type": {
                "type": "string",
                "enum": ["required", "optional", "min_length", "max_length", "pattern", "format", "range", "enum", "custom"],
                "description": "Type of constraint being applied"
            },
            "value": {
                "description": "Constraint value (varies by constraint type)"
            },
            "message": {
                "type": "string",
                "description": "Human-readable error message when constraint is violated"
            },
            "code": {
                "type": "string",
                "description": "Machine-readable error code for constraint violations"
            }
        },
        "additionalProperties": False
    },

    "conversation": {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://schemas.astra.dev/v1/conversation.json",
        "title": "Conversation",
        "description": "Complete ASTRA conversation container with acts and metadata",
        "type": "object",
        "required": ["id", "participants", "acts"],
        "properties": {
            "id": {
                "type": "string",
                "pattern": "^conv_[a-zA-Z0-9_-]+$",
                "description": "Unique identifier for this conversation"
            },
            "participants": {
                "type": "array",
                "minItems": 1,
                "items": {
                    "type": "object",
                    "required": ["id", "type"],
                    "properties": {
                        "id": {"type": "string"},
                        "type": {"type": "string", "enum": ["human", "ai", "system", "bot"]},
                        "role": {"type": "string"},
                        "name": {"type": "string"},
                        "email": {"type": "string", "format": "email"},
                        "phone": {"type": "string"},
                        "external_id": {"type": "string"},
                        "system": {"type": "string"},
                        "capabilities": {"type": "array", "items": {"type": "string"}},
                        "permissions": {"type": "array", "items": {"type": "string"}},
                        "preferences": {"type": "object", "additionalProperties": True},
                        "metadata": {"type": "object", "additionalProperties": True}
                    }
                },
                "description": "List of conversation participants"
            },
            "acts": {
                "type": "array",
                "items": {
                    "oneOf": [
                        {"$ref": "#/ask"},
                        {"$ref": "#/fact"},
                        {"$ref": "#/confirm"},
                        {"$ref": "#/commit"},
                        {"$ref": "#/error"}
                    ]
                },
                "description": "Ordered sequence of acts in this conversation"
            },
            "started_at": {
                "type": "string",
                "format": "date-time",
                "description": "When the conversation started"
            },
            "ended_at": {
                "type": "string",
                "format": "date-time",
                "description": "When the conversation ended"
            },
            "status": {
                "type": "string",
                "enum": ["active", "paused", "completed", "failed", "cancelled"],
                "default": "active",
                "description": "Current status of the conversation"
            },
            "channel": {
                "type": "string",
                "description": "Primary communication channel for this conversation"
            },
            "schema": {
                "type": "string",
                "description": "Business schema identifier used for this conversation"
            },
            "context": {
                "type": "object",
                "description": "Conversation context and session information",
                "properties": {
                    "session_id": {
                        "type": "string",
                        "description": "Session identifier"
                    },
                    "user_agent": {
                        "type": "string",
                        "description": "User agent or client information"
                    },
                    "ip_address": {
                        "type": "string",
                        "description": "Client IP address"
                    },
                    "referrer": {
                        "type": "string",
                        "description": "How the conversation was initiated"
                    }
                },
                "additionalProperties": True
            },
            "final_state": {
                "type": "object",
                "description": "Final computed state of all entities after processing all acts",
                "additionalProperties": True
            },
            "metadata": {
                "type": "object",
                "description": "Additional conversation metadata",
                "properties": {
                    "total_duration_ms": {
                        "type": "integer",
                        "minimum": 0,
                        "description": "Total conversation duration in milliseconds"
                    },
                    "act_count": {
                        "type": "integer",
                        "minimum": 0,
                        "description": "Total number of acts in the conversation"
                    },
                    "error_count": {
                        "type": "integer",
                        "minimum": 0,
                        "description": "Number of errors that occurred"
                    },
                    "commit_count": {
                        "type": "integer",
                        "minimum": 0,
                        "description": "Number of successful commits"
                    },
                    "avg_confidence": {
                        "type": "number",
                        "minimum": 0,
                        "maximum": 1,
                        "description": "Average confidence score across all acts"
                    }
                },
                "additionalProperties": True
            }
        },
        "additionalProperties": False
    }
}
