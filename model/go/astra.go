// Package astra provides Go types and JSON schemas for ASTRA (Act State Representation Architecture) conversations.
//
// ASTRA defines a canonical data model for representing conversational state as typed,
// auditable sequences of actions. This package provides:
//   - Type-safe Go structs for all ASTRA types
//   - JSON schemas for runtime validation
//   - Type guards and utility functions
//   - ID generation and validation utilities
//
// Usage:
//
//	import "github.com/pryszm/astra-model-go"
//
//	// Create an Ask act
//	ask := astra.Ask{
//	    Act: astra.Act{
//	        ID:        astra.GenerateActID(),
//	        Timestamp: time.Now(),
//	        Speaker:   "agent_123",
//	        Type:      astra.ActTypeAsk,
//	    },
//	    Field:    "email",
//	    Prompt:   "What is your email address?",
//	    Required: true,
//	}
//
//	// Validate act type
//	if astra.IsAsk(ask) {
//	    fmt.Printf("Field requested: %s\n", ask.Field)
//	}
package astra

import (
	"encoding/json"
	"fmt"
)

// Version information
const (
	// Version is the current version of this package
	Version = "1.0.0"
	// SchemaVersion is the ASTRA schema version this package implements
	SchemaVersion = "v1"
)

// ConversationAct represents any valid ASTRA act type.
// This interface is implemented by Ask, Fact, Confirm, Commit, and Error.
type ConversationAct interface {
	// GetAct returns the base Act properties
	GetAct() Act
	// GetType returns the act type
	GetType() ActType
	// Validate validates the act structure
	Validate() error
}

// Ensure all act types implement ConversationAct interface
var _ ConversationAct = (*Ask)(nil)
var _ ConversationAct = (*Fact)(nil)
var _ ConversationAct = (*Confirm)(nil)
var _ ConversationAct = (*Commit)(nil)
var _ ConversationAct = (*Error)(nil)

// ActUnion provides a way to work with different act types in a type-safe manner.
// Use this when you need to handle multiple act types dynamically.
type ActUnion struct {
	Type    ActType `json:"type"`
	Ask     *Ask    `json:"ask,omitempty"`
	Fact    *Fact   `json:"fact,omitempty"`
	Confirm *Confirm `json:"confirm,omitempty"`
	Commit  *Commit  `json:"commit,omitempty"`
	Error   *Error   `json:"error,omitempty"`
}

// GetAct returns the underlying act, regardless of type
func (u ActUnion) GetAct() (ConversationAct, error) {
	switch u.Type {
	case ActTypeAsk:
		if u.Ask != nil {
			return *u.Ask, nil
		}
	case ActTypeFact:
		if u.Fact != nil {
			return *u.Fact, nil
		}
	case ActTypeConfirm:
		if u.Confirm != nil {
			return *u.Confirm, nil
		}
	case ActTypeCommit:
		if u.Commit != nil {
			return *u.Commit, nil
		}
	case ActTypeError:
		if u.Error != nil {
			return *u.Error, nil
		}
	}
	return nil, fmt.Errorf("invalid act union: type %s has no corresponding act", u.Type)
}

// NewActUnion creates a new ActUnion from a ConversationAct
func NewActUnion(act ConversationAct) ActUnion {
	union := ActUnion{Type: act.GetType()}
	
	switch a := act.(type) {
	case Ask:
		union.Ask = &a
	case Fact:
		union.Fact = &a
	case Confirm:
		union.Confirm = &a
	case Commit:
		union.Commit = &a
	case Error:
		union.Error = &a
	}
	
	return union
}

// Type Guards - Runtime type checking functions

// IsAct checks if an interface{} is a valid ASTRA Act
func IsAct(v interface{}) bool {
	_, ok := v.(ConversationAct)
	return ok
}

// IsAsk checks if an interface{} is an Ask act
func IsAsk(v interface{}) bool {
	_, ok := v.(Ask)
	return ok
}

// IsFact checks if an interface{} is a Fact act
func IsFact(v interface{}) bool {
	_, ok := v.(Fact)
	return ok
}

// IsConfirm checks if an interface{} is a Confirm act
func IsConfirm(v interface{}) bool {
	_, ok := v.(Confirm)
	return ok
}

// IsCommit checks if an interface{} is a Commit act
func IsCommit(v interface{}) bool {
	_, ok := v.(Commit)
	return ok
}

// IsError checks if an interface{} is an Error act
func IsError(v interface{}) bool {
	_, ok := v.(Error)
	return ok
}

// IsParticipant checks if an interface{} is a valid Participant
func IsParticipant(v interface{}) bool {
	p, ok := v.(Participant)
	if !ok {
		return false
	}
	return p.ID != "" && p.Type != ""
}

// IsEntity checks if an interface{} is a valid Entity
func IsEntity(v interface{}) bool {
	switch e := v.(type) {
	case Entity:
		return e.ID != "" && e.Type != ""
	case string:
		return e != ""
	default:
		return false
	}
}

// IsConversation checks if an interface{} is a valid Conversation
func IsConversation(v interface{}) bool {
	c, ok := v.(Conversation)
	if !ok {
		return false
	}
	return c.ID != "" && len(c.Participants) > 0 && c.Acts != nil
}

// JSON Marshaling and Unmarshaling helpers

// MarshalAct marshals any ConversationAct to JSON
func MarshalAct(act ConversationAct) ([]byte, error) {
	return json.Marshal(act)
}

// UnmarshalAct unmarshals JSON to the appropriate ConversationAct type
func UnmarshalAct(data []byte) (ConversationAct, error) {
	// First, determine the type
	var typeCheck struct {
		Type ActType `json:"type"`
	}
	
	if err := json.Unmarshal(data, &typeCheck); err != nil {
		return nil, fmt.Errorf("failed to determine act type: %w", err)
	}
	
	// Unmarshal to the specific type
	switch typeCheck.Type {
	case ActTypeAsk:
		var ask Ask
		err := json.Unmarshal(data, &ask)
		return ask, err
	case ActTypeFact:
		var fact Fact
		err := json.Unmarshal(data, &fact)
		return fact, err
	case ActTypeConfirm:
		var confirm Confirm
		err := json.Unmarshal(data, &confirm)
		return confirm, err
	case ActTypeCommit:
		var commit Commit
		err := json.Unmarshal(data, &commit)
		return commit, err
	case ActTypeError:
		var errorAct Error
		err := json.Unmarshal(data, &errorAct)
		return errorAct, err
	default:
		return nil, fmt.Errorf("unknown act type: %s", typeCheck.Type)
	}
}

// Validation functions

// ValidateAct validates a ConversationAct against its schema requirements
func ValidateAct(act ConversationAct) error {
	if act == nil {
		return fmt.Errorf("act cannot be nil")
	}
	
	// Validate base act properties
	baseAct := act.GetAct()
	if err := validateBaseAct(baseAct); err != nil {
		return fmt.Errorf("invalid base act: %w", err)
	}
	
	// Validate specific act type
	return act.Validate()
}

// validateBaseAct validates the base Act properties
func validateBaseAct(act Act) error {
	if act.ID == "" {
		return fmt.Errorf("act ID is required")
	}
	
	if !IsValidActID(act.ID) {
		return fmt.Errorf("invalid act ID format: %s", act.ID)
	}
	
	if act.Speaker == "" {
		return fmt.Errorf("act speaker is required")
	}
	
	if act.Type == "" {
		return fmt.Errorf("act type is required")
	}
	
	if !isValidActType(act.Type) {
		return fmt.Errorf("invalid act type: %s", act.Type)
	}
	
	if act.Confidence != nil && (*act.Confidence < 0.0 || *act.Confidence > 1.0) {
		return fmt.Errorf("confidence must be between 0.0 and 1.0, got: %f", *act.Confidence)
	}
	
	return nil
}

// isValidActType checks if an ActType is valid
func isValidActType(actType ActType) bool {
	switch actType {
	case ActTypeAsk, ActTypeFact, ActTypeConfirm, ActTypeCommit, ActTypeError:
		return true
	default:
		return false
	}
}

// Error types for better error handling

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

// InvalidActTypeError represents an error when an invalid act type is encountered
type InvalidActTypeError struct {
	ActType string
}

func (e InvalidActTypeError) Error() string {
	return fmt.Sprintf("invalid act type: %s", e.ActType)
}

// ActNotFoundError represents an error when an act is not found
type ActNotFoundError struct {
	ActID string
}

func (e ActNotFoundError) Error() string {
	return fmt.Sprintf("act not found: %s", e.ActID)
}
