package astra

import (
	"encoding/json"
	"fmt"
)

// Schema represents a JSON Schema definition
type Schema map[string]interface{}

// Schemas contains all embedded JSON Schema definitions for ASTRA types
var Schemas = struct {
	Act          Schema
	Ask          Schema
	Fact         Schema
	Confirm      Schema
	Commit       Schema
	Error        Schema
	Entity       Schema
	Participant  Schema
	Constraint   Schema
	Conversation Schema
}{
	Act: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/act.json",
		"title":       "Act",
		"description": "Base type for all conversational actions in ASTRA",
		"type":        "object",
		"required":    []string{"id", "timestamp", "speaker", "type"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"pattern":     "^act_[a-zA-Z0-9_-]+$",
				"description": "Unique identifier for this act within the conversation",
			},
			"timestamp": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"description": "ISO 8601 timestamp when the act occurred",
			},
			"speaker": map[string]interface{}{
				"type":        "string",
				"description": "Identifier of the conversation participant who performed this act",
			},
			"type": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"ask", "fact", "confirm", "commit", "error"},
				"description": "Type of conversational act being performed",
			},
			"confidence": map[string]interface{}{
				"type":        "number",
				"minimum":     0.0,
				"maximum":     1.0,
				"description": "Confidence score for automated act extraction (0.0 to 1.0)",
			},
			"source": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"human", "speech_recognition", "text_analysis", "system", "ai"},
				"description": "Source that generated this act",
			},
			"metadata": map[string]interface{}{
				"type":        "object",
				"description": "Additional context-specific metadata",
				"properties": map[string]interface{}{
					"channel": map[string]interface{}{
						"type":        "string",
						"description": "Communication channel (voice, text, email, etc.)",
					},
					"language": map[string]interface{}{
						"type":        "string",
						"pattern":     "^[a-z]{2}(-[A-Z]{2})?$",
						"description": "Language code (ISO 639-1, optional region)",
					},
					"original_text": map[string]interface{}{
						"type":        "string",
						"description": "Original utterance that generated this act",
					},
					"processing_time_ms": map[string]interface{}{
						"type":        "number",
						"minimum":     0,
						"description": "Time taken to process this act in milliseconds",
					},
				},
				"additionalProperties": true,
			},
		},
		"additionalProperties": false,
	},

	Ask: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/ask.json",
		"title":       "Ask",
		"description": "Act that requests missing information required to complete a business process",
		"type":        "object",
		"required":    []string{"id", "timestamp", "speaker", "type", "field", "prompt"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"pattern":     "^act_[a-zA-Z0-9_-]+$",
				"description": "Unique identifier for this act within the conversation",
			},
			"timestamp": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"description": "ISO 8601 timestamp when the act occurred",
			},
			"speaker": map[string]interface{}{
				"type":        "string",
				"description": "Identifier of the conversation participant who performed this act",
			},
			"type": map[string]interface{}{
				"const": "ask",
			},
			"confidence": map[string]interface{}{
				"type":        "number",
				"minimum":     0.0,
				"maximum":     1.0,
				"description": "Confidence score for automated act extraction (0.0 to 1.0)",
			},
			"source": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"human", "speech_recognition", "text_analysis", "system", "ai"},
				"description": "Source that generated this act",
			},
			"metadata": map[string]interface{}{
				"type":                 "object",
				"description":          "Additional context-specific metadata",
				"additionalProperties": true,
			},
			"field": map[string]interface{}{
				"type":        "string",
				"description": "Field or information being requested",
			},
			"prompt": map[string]interface{}{
				"type":        "string",
				"description": "Question or request presented to obtain the information",
			},
			"constraints": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"type"},
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type": "string",
							"enum": []string{"required", "optional", "min_length", "max_length", "pattern", "format", "range", "enum", "custom"},
						},
						"value":   map[string]interface{}{},
						"message": map[string]interface{}{
				"type":        "string",
				"description": "Human-readable error message when constraint is violated",
			},
			"code": map[string]interface{}{
				"type":        "string",
				"description": "Machine-readable error code for constraint violations",
			},
		},
		"additionalProperties": false,
	},

	Conversation: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/conversation.json",
		"title":       "Conversation",
		"description": "Complete ASTRA conversation container with acts and metadata",
		"type":        "object",
		"required":    []string{"id", "participants", "acts"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"pattern":     "^conv_[a-zA-Z0-9_-]+$",
				"description": "Unique identifier for this conversation",
			},
			"participants": map[string]interface{}{
				"type":     "array",
				"minItems": 1,
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"id", "type"},
					"properties": map[string]interface{}{
						"id":   map[string]interface{}{"type": "string"},
						"type": map[string]interface{}{"type": "string", "enum": []string{"human", "ai", "system", "bot"}},
						"role": map[string]interface{}{"type": "string"},
						"name": map[string]interface{}{"type": "string"},
						"email": map[string]interface{}{
							"type":   "string",
							"format": "email",
						},
						"phone":       map[string]interface{}{"type": "string"},
						"external_id": map[string]interface{}{"type": "string"},
						"system":      map[string]interface{}{"type": "string"},
						"capabilities": map[string]interface{}{
							"type":  "array",
							"items": map[string]interface{}{"type": "string"},
						},
						"permissions": map[string]interface{}{
							"type":  "array",
							"items": map[string]interface{}{"type": "string"},
						},
						"preferences": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": true,
						},
						"metadata": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": true,
						},
					},
				},
				"description": "List of conversation participants",
			},
			"acts": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"oneOf": []map[string]interface{}{
						{"$ref": "#/definitions/ask"},
						{"$ref": "#/definitions/fact"},
						{"$ref": "#/definitions/confirm"},
						{"$ref": "#/definitions/commit"},
						{"$ref": "#/definitions/error"},
					},
				},
				"description": "Ordered sequence of acts in this conversation",
			},
			"started_at": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"description": "When the conversation started",
			},
			"ended_at": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"description": "When the conversation ended",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"active", "paused", "completed", "failed", "cancelled"},
				"default":     "active",
				"description": "Current status of the conversation",
			},
			"channel": map[string]interface{}{
				"type":        "string",
				"description": "Primary communication channel for this conversation",
			},
			"schema": map[string]interface{}{
				"type":        "string",
				"description": "Business schema identifier used for this conversation",
			},
			"context": map[string]interface{}{
				"type":        "object",
				"description": "Conversation context and session information",
				"properties": map[string]interface{}{
					"session_id": map[string]interface{}{
						"type":        "string",
						"description": "Session identifier",
					},
					"user_agent": map[string]interface{}{
						"type":        "string",
						"description": "User agent or client information",
					},
					"ip_address": map[string]interface{}{
						"type":        "string",
						"description": "Client IP address",
					},
					"referrer": map[string]interface{}{
						"type":        "string",
						"description": "How the conversation was initiated",
					},
				},
				"additionalProperties": true,
			},
			"final_state": map[string]interface{}{
				"type":                 "object",
				"description":          "Final computed state of all entities after processing all acts",
				"additionalProperties": true,
			},
			"metadata": map[string]interface{}{
				"type":        "object",
				"description": "Additional conversation metadata",
				"properties": map[string]interface{}{
					"total_duration_ms": map[string]interface{}{
						"type":        "integer",
						"minimum":     0,
						"description": "Total conversation duration in milliseconds",
					},
					"act_count": map[string]interface{}{
						"type":        "integer",
						"minimum":     0,
						"description": "Total number of acts in the conversation",
					},
					"error_count": map[string]interface{}{
						"type":        "integer",
						"minimum":     0,
						"description": "Number of errors that occurred",
					},
					"commit_count": map[string]interface{}{
						"type":        "integer",
						"minimum":     0,
						"description": "Number of successful commits",
					},
					"avg_confidence": map[string]interface{}{
						"type":        "number",
						"minimum":     0,
						"maximum":     1,
						"description": "Average confidence score across all acts",
					},
				},
				"additionalProperties": true,
			},
		},
		"additionalProperties": false,
	},
}

// GetSchema returns a specific schema by name
func GetSchema(name string) (Schema, error) {
	switch name {
	case "act":
		return Schemas.Act, nil
	case "ask":
		return Schemas.Ask, nil
	case "fact":
		return Schemas.Fact, nil
	case "confirm":
		return Schemas.Confirm, nil
	case "commit":
		return Schemas.Commit, nil
	case "error":
		return Schemas.Error, nil
	case "entity":
		return Schemas.Entity, nil
	case "participant":
		return Schemas.Participant, nil
	case "constraint":
		return Schemas.Constraint, nil
	case "conversation":
		return Schemas.Conversation, nil
	default:
		return nil, fmt.Errorf("unknown schema: %s", name)
	}
}

// ValidateJSON validates a JSON byte slice against a named schema
func ValidateJSON(data []byte, schemaName string) error {
	schema, err := GetSchema(schemaName)
	if err != nil {
		return err
	}
	
	// Parse JSON data
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	
	// Basic validation (full JSON Schema validation would require a dedicated library)
	return validateAgainstSchema(jsonData, schema)
}

// validateAgainstSchema performs basic validation against a schema
// Note: This is a simplified validator. For full JSON Schema validation,
// consider using a dedicated library like github.com/xeipuuv/gojsonschema
func validateAgainstSchema(data interface{}, schema Schema) error {
	// Check if data is an object when schema expects object
	if schemaType, ok := schema["type"].(string); ok && schemaType == "object" {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected object, got %T", data)
		}
		
		// Check required fields
		if required, ok := schema["required"].([]string); ok {
			for _, field := range required {
				if _, exists := dataMap[field]; !exists {
					return fmt.Errorf("required field missing: %s", field)
				}
			}
		}
		
		// Check properties if they exist in schema
		if properties, ok := schema["properties"].(map[string]interface{}); ok {
			for key, value := range dataMap {
				if propSchema, exists := properties[key]; exists {
					if propMap, ok := propSchema.(map[string]interface{}); ok {
						if err := validateProperty(value, propMap); err != nil {
							return fmt.Errorf("validation failed for property %s: %w", key, err)
						}
					}
				}
			}
		}
	}
	
	return nil
}

// validateProperty validates a single property against its schema
func validateProperty(value interface{}, propSchema map[string]interface{}) error {
	// Check type constraints
	if expectedType, ok := propSchema["type"].(string); ok {
		if !validateType(value, expectedType) {
			return fmt.Errorf("expected type %s, got %T", expectedType, value)
		}
	}
	
	// Check const constraints
	if constValue, ok := propSchema["const"]; ok {
		if value != constValue {
			return fmt.Errorf("expected const value %v, got %v", constValue, value)
		}
	}
	
	// Check enum constraints
	if enumValues, ok := propSchema["enum"].([]interface{}); ok {
		found := false
		for _, enumValue := range enumValues {
			if value == enumValue {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("value %v not in enum %v", value, enumValues)
		}
	}
	
	// Check numeric constraints
	if num, ok := value.(float64); ok {
		if min, exists := propSchema["minimum"].(float64); exists && num < min {
			return fmt.Errorf("value %f is below minimum %f", num, min)
		}
		if max, exists := propSchema["maximum"].(float64); exists && num > max {
			return fmt.Errorf("value %f is above maximum %f", num, max)
		}
	}
	
	// Check string constraints
	if str, ok := value.(string); ok {
		if pattern, exists := propSchema["pattern"].(string); exists {
			// Note: Full regex validation would require regexp package
			if pattern == "^act_[a-zA-Z0-9_-]+$" && !IsValidActID(str) {
				return fmt.Errorf("string %s does not match pattern %s", str, pattern)
			}
			if pattern == "^conv_[a-zA-Z0-9_-]+$" && !IsValidConversationID(str) {
				return fmt.Errorf("string %s does not match pattern %s", str, pattern)
			}
		}
	}
	
	return nil
}

// validateType checks if a value matches the expected JSON Schema type
func validateType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		_, ok := value.(float64)
		return ok
	case "integer":
		if f, ok := value.(float64); ok {
			return f == float64(int64(f))
		}
		return false
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "object":
		_, ok := value.(map[string]interface{})
		return ok
	case "array":
		_, ok := value.([]interface{})
		return ok
	case "null":
		return value == nil
	default:
		return false
	}
}

// SchemaVersion returns the schema version for a given schema
func SchemaVersion(schemaName string) string {
	return "v1" // All current schemas are v1
}

// ListSchemas returns all available schema names
func ListSchemas() []string {
	return []string{
		"act",
		"ask", 
		"fact",
		"confirm",
		"commit",
		"error",
		"entity",
		"participant",
		"constraint",
		"conversation",
	}
}"type": "string"},
						"code":    map[string]interface{}{"type": "string"},
					},
				},
				"description": "Validation constraints for the requested information",
			},
			"required": map[string]interface{}{
				"type":        "boolean",
				"default":     true,
				"description": "Whether this information is required to proceed",
			},
			"expected_type": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"string", "number", "boolean", "object", "array", "date", "email", "phone", "address"},
				"description": "Expected data type of the response",
			},
			"retry_count": map[string]interface{}{
				"type":        "integer",
				"minimum":     0,
				"default":     0,
				"description": "Number of times this question has been asked",
			},
			"max_retries": map[string]interface{}{
				"type":        "integer",
				"minimum":     0,
				"default":     3,
				"description": "Maximum number of retry attempts before escalation",
			},
		},
		"additionalProperties": false,
	},

	Fact: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/fact.json",
		"title":       "Fact",
		"description": "Act that declares facts or information provided during conversation",
		"type":        "object",
		"required":    []string{"id", "timestamp", "speaker", "type", "entity", "field", "value"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"pattern":     "^act_[a-zA-Z0-9_-]+$",
				"description": "Unique identifier for this act within the conversation",
			},
			"timestamp": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"description": "ISO 8601 timestamp when the act occurred",
			},
			"speaker": map[string]interface{}{
				"type":        "string",
				"description": "Identifier of the conversation participant who performed this act",
			},
			"type": map[string]interface{}{
				"const": "fact",
			},
			"confidence": map[string]interface{}{
				"type":        "number",
				"minimum":     0.0,
				"maximum":     1.0,
				"description": "Confidence score for automated act extraction (0.0 to 1.0)",
			},
			"source": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"human", "speech_recognition", "text_analysis", "system", "ai"},
				"description": "Source that generated this act",
			},
			"metadata": map[string]interface{}{
				"type":                 "object",
				"description":          "Additional context-specific metadata",
				"additionalProperties": true,
			},
			"entity": map[string]interface{}{
				"oneOf": []map[string]interface{}{
					{
						"type":        "string",
						"description": "Entity identifier as string",
					},
					{
						"type":     "object",
						"required": []string{"id", "type"},
						"properties": map[string]interface{}{
							"id":          map[string]interface{}{"type": "string"},
							"type":        map[string]interface{}{"type": "string"},
							"external_id": map[string]interface{}{"type": "string"},
							"system":      map[string]interface{}{"type": "string"},
							"version":     map[string]interface{}{"type": "string"},
							"schema_url": map[string]interface{}{
								"type":   "string",
								"format": "uri",
							},
							"metadata": map[string]interface{}{
								"type":                 "object",
								"additionalProperties": true,
							},
						},
						"additionalProperties": false,
						"description":          "Structured entity reference",
					},
				},
				"description": "Business entity being modified (order, customer, appointment, etc.)",
			},
			"field": map[string]interface{}{
				"type":        "string",
				"description": "Specific field or property being set",
			},
			"value": map[string]interface{}{
				"description": "Value being assigned to the field (any JSON type)",
			},
			"operation": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"set", "append", "increment", "decrement", "delete", "merge"},
				"default":     "set",
				"description": "Operation being performed on the field",
			},
			"previous_value": map[string]interface{}{
				"description": "Previous value of the field (for audit trail)",
			},
			"validation_status": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"pending", "valid", "invalid", "partial"},
				"default":     "pending",
				"description": "Validation status of this fact",
			},
			"validation_errors": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of validation errors if validation_status is invalid",
			},
		},
		"additionalProperties": false,
	},

	Confirm: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/confirm.json",
		"title":       "Confirm",
		"description": "Act that verifies understanding of information before commitment",
		"type":        "object",
		"required":    []string{"id", "timestamp", "speaker", "type", "entity", "summary"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"pattern":     "^act_[a-zA-Z0-9_-]+$",
				"description": "Unique identifier for this act within the conversation",
			},
			"timestamp": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"description": "ISO 8601 timestamp when the act occurred",
			},
			"speaker": map[string]interface{}{
				"type":        "string",
				"description": "Identifier of the conversation participant who performed this act",
			},
			"type": map[string]interface{}{
				"const": "confirm",
			},
			"confidence": map[string]interface{}{
				"type":        "number",
				"minimum":     0.0,
				"maximum":     1.0,
				"description": "Confidence score for automated act extraction (0.0 to 1.0)",
			},
			"source": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"human", "speech_recognition", "text_analysis", "system", "ai"},
				"description": "Source that generated this act",
			},
			"metadata": map[string]interface{}{
				"type":                 "object",
				"description":          "Additional context-specific metadata",
				"additionalProperties": true,
			},
			"entity": map[string]interface{}{
				"oneOf": []map[string]interface{}{
					{
						"type":        "string",
						"description": "Entity identifier as string",
					},
					{
						"type":     "object",
						"required": []string{"id", "type"},
						"properties": map[string]interface{}{
							"id":          map[string]interface{}{"type": "string"},
							"type":        map[string]interface{}{"type": "string"},
							"external_id": map[string]interface{}{"type": "string"},
							"system":      map[string]interface{}{"type": "string"},
							"version":     map[string]interface{}{"type": "string"},
							"schema_url": map[string]interface{}{
								"type":   "string",
								"format": "uri",
							},
							"metadata": map[string]interface{}{
								"type":                 "object",
								"additionalProperties": true,
							},
						},
						"additionalProperties": false,
						"description":          "Structured entity reference",
					},
				},
				"description": "Business entity being confirmed",
			},
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "Human-readable summary of what is being confirmed",
			},
			"awaiting": map[string]interface{}{
				"type":        "boolean",
				"default":     true,
				"description": "Whether confirmation is still pending",
			},
			"confirmed": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the confirmation was accepted (true) or rejected (false)",
			},
			"confirmation_method": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"verbal", "explicit", "implicit", "timeout", "system"},
				"description": "How the confirmation was obtained",
			},
			"fields_confirmed": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Specific fields or aspects being confirmed",
			},
			"rejection_reason": map[string]interface{}{
				"type":        "string",
				"description": "Reason provided if confirmation was rejected",
			},
			"timeout_ms": map[string]interface{}{
				"type":        "integer",
				"minimum":     0,
				"description": "Timeout for awaiting confirmation in milliseconds",
			},
		},
		"additionalProperties": false,
	},

	Commit: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/commit.json",
		"title":       "Commit",
		"description": "Act that executes business processes and triggers system integrations",
		"type":        "object",
		"required":    []string{"id", "timestamp", "speaker", "type", "entity", "action"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"pattern":     "^act_[a-zA-Z0-9_-]+$",
				"description": "Unique identifier for this act within the conversation",
			},
			"timestamp": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"description": "ISO 8601 timestamp when the act occurred",
			},
			"speaker": map[string]interface{}{
				"type":        "string",
				"description": "Identifier of the conversation participant who performed this act",
			},
			"type": map[string]interface{}{
				"const": "commit",
			},
			"confidence": map[string]interface{}{
				"type":        "number",
				"minimum":     0.0,
				"maximum":     1.0,
				"description": "Confidence score for automated act extraction (0.0 to 1.0)",
			},
			"source": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"human", "speech_recognition", "text_analysis", "system", "ai"},
				"description": "Source that generated this act",
			},
			"metadata": map[string]interface{}{
				"type":                 "object",
				"description":          "Additional context-specific metadata",
				"additionalProperties": true,
			},
			"entity": map[string]interface{}{
				"oneOf": []map[string]interface{}{
					{
						"type":        "string",
						"description": "Entity identifier as string",
					},
					{
						"type":     "object",
						"required": []string{"id", "type"},
						"properties": map[string]interface{}{
							"id":          map[string]interface{}{"type": "string"},
							"type":        map[string]interface{}{"type": "string"},
							"external_id": map[string]interface{}{"type": "string"},
							"system":      map[string]interface{}{"type": "string"},
							"version":     map[string]interface{}{"type": "string"},
							"schema_url": map[string]interface{}{
								"type":   "string",
								"format": "uri",
							},
							"metadata": map[string]interface{}{
								"type":                 "object",
								"additionalProperties": true,
							},
						},
						"additionalProperties": false,
						"description":          "Structured entity reference",
					},
				},
				"description": "Business entity being committed to external systems",
			},
			"action": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"create", "update", "delete", "execute", "cancel", "pause", "resume"},
				"description": "Action being performed in the target system",
			},
			"system": map[string]interface{}{
				"type":        "string",
				"description": "Target system identifier (CRM, order_management, etc.)",
			},
			"transaction_id": map[string]interface{}{
				"type":        "string",
				"description": "External system transaction or record identifier",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"pending", "in_progress", "success", "failed", "retrying", "cancelled"},
				"default":     "pending",
				"description": "Status of the commit operation",
			},
			"error": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]interface{}{
						"type":        "string",
						"description": "Error code from the target system",
					},
					"message": map[string]interface{}{
						"type":        "string",
						"description": "Human-readable error message",
					},
					"details": map[string]interface{}{
						"type":        "object",
						"description": "Additional error context",
					},
					"recoverable": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether the error can be recovered from",
					},
				},
				"required":    []string{"code", "message"},
				"description": "Error information if commit failed",
			},
			"retry_count": map[string]interface{}{
				"type":        "integer",
				"minimum":     0,
				"default":     0,
				"description": "Number of retry attempts made",
			},
			"max_retries": map[string]interface{}{
				"type":        "integer",
				"minimum":     0,
				"default":     3,
				"description": "Maximum number of retry attempts",
			},
			"idempotency_key": map[string]interface{}{
				"type":        "string",
				"description": "Key to ensure idempotent operations",
			},
			"rollback_info": map[string]interface{}{
				"type":        "object",
				"description": "Information needed to rollback this commit if necessary",
			},
		},
		"additionalProperties": false,
	},

	Error: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/error.json",
		"title":       "Error",
		"description": "Act that handles failures and exceptions in conversational processing",
		"type":        "object",
		"required":    []string{"id", "timestamp", "speaker", "type", "code", "message", "recoverable"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"pattern":     "^act_[a-zA-Z0-9_-]+$",
				"description": "Unique identifier for this act within the conversation",
			},
			"timestamp": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"description": "ISO 8601 timestamp when the act occurred",
			},
			"speaker": map[string]interface{}{
				"type":        "string",
				"description": "Identifier of the conversation participant who performed this act",
			},
			"type": map[string]interface{}{
				"const": "error",
			},
			"confidence": map[string]interface{}{
				"type":        "number",
				"minimum":     0.0,
				"maximum":     1.0,
				"description": "Confidence score for automated act extraction (0.0 to 1.0)",
			},
			"source": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"human", "speech_recognition", "text_analysis", "system", "ai"},
				"description": "Source that generated this act",
			},
			"metadata": map[string]interface{}{
				"type":                 "object",
				"description":          "Additional context-specific metadata",
				"additionalProperties": true,
			},
			"code": map[string]interface{}{
				"type":        "string",
				"description": "Machine-readable error code",
			},
			"message": map[string]interface{}{
				"type":        "string",
				"description": "Human-readable error message",
			},
			"recoverable": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the conversation can continue after this error",
			},
			"severity": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"info", "warning", "error", "critical"},
				"default":     "error",
				"description": "Severity level of the error",
			},
			"category": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"validation", "processing", "integration", "timeout", "permission", "system", "user_input", "business_rule"},
				"description": "Category of error for classification",
			},
			"details": map[string]interface{}{
				"type":        "object",
				"description": "Additional error context and debugging information",
			},
			"related_act_id": map[string]interface{}{
				"type":        "string",
				"pattern":     "^act_[a-zA-Z0-9_-]+$",
				"description": "ID of the act that caused this error",
			},
			"suggested_action": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"retry", "escalate", "ignore", "clarify", "fallback", "terminate"},
				"description": "Suggested recovery action",
			},
			"user_message": map[string]interface{}{
				"type":        "string",
				"description": "User-friendly message to display to conversation participants",
			},
			"stack_trace": map[string]interface{}{
				"type":        "string",
				"description": "Technical stack trace for debugging (not shown to users)",
			},
		},
		"additionalProperties": false,
	},

	Entity: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/entity.json",
		"title":       "Entity",
		"description": "Reference to a business entity in ASTRA conversations",
		"type":        "object",
		"required":    []string{"id", "type"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Unique identifier for this entity within the conversation scope",
			},
			"type": map[string]interface{}{
				"type":        "string",
				"description": "Type of business entity (order, customer, appointment, ticket, etc.)",
			},
			"external_id": map[string]interface{}{
				"type":        "string",
				"description": "External system identifier for this entity",
			},
			"system": map[string]interface{}{
				"type":        "string",
				"description": "External system that owns this entity",
			},
			"version": map[string]interface{}{
				"type":        "string",
				"description": "Version or revision of this entity",
			},
			"schema_url": map[string]interface{}{
				"type":        "string",
				"format":      "uri",
				"description": "URL to the schema definition for this entity type",
			},
			"metadata": map[string]interface{}{
				"type":                 "object",
				"description":          "Additional entity-specific metadata",
				"additionalProperties": true,
			},
		},
		"additionalProperties": false,
	},

	Participant: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/participant.json",
		"title":       "Participant",
		"description": "Conversation participant in ASTRA conversations",
		"type":        "object",
		"required":    []string{"id", "type"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Unique identifier for this participant",
			},
			"type": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"human", "ai", "system", "bot"},
				"description": "Type of participant",
			},
			"role": map[string]interface{}{
				"type":        "string",
				"description": "Business role of the participant (customer, agent, manager, etc.)",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Display name of the participant",
			},
			"email": map[string]interface{}{
				"type":        "string",
				"format":      "email",
				"description": "Email address of the participant",
			},
			"phone": map[string]interface{}{
				"type":        "string",
				"description": "Phone number of the participant",
			},
			"external_id": map[string]interface{}{
				"type":        "string",
				"description": "External system identifier for this participant",
			},
			"system": map[string]interface{}{
				"type":        "string",
				"description": "External system that manages this participant",
			},
			"capabilities": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of capabilities this participant has",
			},
			"permissions": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of permissions granted to this participant",
			},
			"preferences": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"language": map[string]interface{}{
						"type":        "string",
						"pattern":     "^[a-z]{2}(-[A-Z]{2})?$",
						"description": "Preferred language code",
					},
					"timezone": map[string]interface{}{
						"type":        "string",
						"description": "Preferred timezone (IANA timezone identifier)",
					},
					"communication_channels": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
						"description": "Preferred communication channels in order of preference",
					},
				},
				"additionalProperties": true,
				"description":          "Participant preferences",
			},
			"metadata": map[string]interface{}{
				"type":                 "object",
				"description":          "Additional participant metadata",
				"additionalProperties": true,
			},
		},
		"additionalProperties": false,
	},

	Constraint: Schema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://schemas.astra.dev/v1/constraint.json",
		"title":       "Constraint",
		"description": "Validation constraint for ASTRA fields and values",
		"type":        "object",
		"required":    []string{"type"},
		"properties": map[string]interface{}{
			"type": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"required", "optional", "min_length", "max_length", "pattern", "format", "range", "enum", "custom"},
				"description": "Type of constraint being applied",
			},
			"value": map[string]interface{}{
				"description": "Constraint value (varies by constraint type)",
			},
			"message": map[string]interface{
