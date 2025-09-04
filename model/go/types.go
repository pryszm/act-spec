package astra

import (
	"encoding/json"
	"fmt"
	"time"
)

// ============================================================================
// Base Types
// ============================================================================

// Source represents the source that generated an act
type Source string

const (
	SourceHuman              Source = "human"
	SourceSpeechRecognition  Source = "speech_recognition"
	SourceTextAnalysis       Source = "text_analysis"
	SourceSystem             Source = "system"
	SourceAI                 Source = "ai"
)

// ActType represents the type of conversational act being performed
type ActType string

const (
	ActTypeAsk     ActType = "ask"
	ActTypeFact    ActType = "fact"
	ActTypeConfirm ActType = "confirm"
	ActTypeCommit  ActType = "commit"
	ActTypeError   ActType = "error"
)

// ActMetadata contains additional metadata for acts
type ActMetadata struct {
	// Communication channel (voice, text, email, etc.)
	Channel *string `json:"channel,omitempty"`
	// Language code (ISO 639-1, optional region)
	Language *string `json:"language,omitempty"`
	// Original utterance that generated this act
	OriginalText *string `json:"original_text,omitempty"`
	// Time taken to process this act in milliseconds
	ProcessingTimeMs *float64 `json:"processing_time_ms,omitempty"`
	// Additional context-specific metadata
	AdditionalProperties map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for ActMetadata
func (m ActMetadata) MarshalJSON() ([]byte, error) {
	type Alias ActMetadata
	base := map[string]interface{}{
		"channel":             m.Channel,
		"language":            m.Language,
		"original_text":       m.OriginalText,
		"processing_time_ms":  m.ProcessingTimeMs,
	}
	
	// Add additional properties
	for k, v := range m.AdditionalProperties {
		base[k] = v
	}
	
	return json.Marshal(base)
}

// UnmarshalJSON implements custom JSON unmarshaling for ActMetadata
func (m *ActMetadata) UnmarshalJSON(data []byte) error {
	type Alias ActMetadata
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// Extract additional properties
	m.AdditionalProperties = make(map[string]interface{})
	for k, v := range raw {
		switch k {
		case "channel", "language", "original_text", "processing_time_ms":
			// Skip known fields
		default:
			m.AdditionalProperties[k] = v
		}
	}
	
	return nil
}

// Act represents the base type for all conversational actions in ASTRA
type Act struct {
	// Unique identifier for this act within the conversation
	ID string `json:"id"`
	// Timestamp when the act occurred
	Timestamp time.Time `json:"timestamp"`
	// Identifier of the conversation participant who performed this act
	Speaker string `json:"speaker"`
	// Type of conversational act being performed
	Type ActType `json:"type"`
	// Confidence score for automated act extraction (0.0 to 1.0)
	Confidence *float64 `json:"confidence,omitempty"`
	// Source that generated this act
	Source *Source `json:"source,omitempty"`
	// Additional context-specific metadata
	Metadata *ActMetadata `json:"metadata,omitempty"`
}

// GetAct implements ConversationAct interface
func (a Act) GetAct() Act {
	return a
}

// GetType implements ConversationAct interface
func (a Act) GetType() ActType {
	return a.Type
}

// ============================================================================
// Entity Types
// ============================================================================

// Entity represents a reference to a business entity in ASTRA conversations
type Entity struct {
	// Unique identifier for this entity within the conversation scope
	ID string `json:"id"`
	// Type of business entity (order, customer, appointment, ticket, etc.)
	Type string `json:"type"`
	// External system identifier for this entity
	ExternalID *string `json:"external_id,omitempty"`
	// External system that owns this entity
	System *string `json:"system,omitempty"`
	// Version or revision of this entity
	Version *string `json:"version,omitempty"`
	// URL to the schema definition for this entity type
	SchemaURL *string `json:"schema_url,omitempty"`
	// Additional entity-specific metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// EntityRef represents an entity reference that can be either a string ID or structured Entity
type EntityRef interface{}

// GetEntityID extracts the entity ID from an EntityRef
func GetEntityID(ref EntityRef) (string, error) {
	switch e := ref.(type) {
	case string:
		return e, nil
	case Entity:
		return e.ID, nil
	case *Entity:
		if e == nil {
			return "", fmt.Errorf("entity reference is nil")
		}
		return e.ID, nil
	default:
		return "", fmt.Errorf("invalid entity reference type: %T", ref)
	}
}

// ============================================================================
// Constraint Types
// ============================================================================

// ConstraintType represents the type of constraint being applied
type ConstraintType string

const (
	ConstraintTypeRequired   ConstraintType = "required"
	ConstraintTypeOptional   ConstraintType = "optional"
	ConstraintTypeMinLength  ConstraintType = "min_length"
	ConstraintTypeMaxLength  ConstraintType = "max_length"
	ConstraintTypePattern    ConstraintType = "pattern"
	ConstraintTypeFormat     ConstraintType = "format"
	ConstraintTypeRange      ConstraintType = "range"
	ConstraintTypeEnum       ConstraintType = "enum"
	ConstraintTypeCustom     ConstraintType = "custom"
)

// FormatType represents format validation types
type FormatType string

const (
	FormatTypeEmail    FormatType = "email"
	FormatTypePhone    FormatType = "phone"
	FormatTypeURL      FormatType = "url"
	FormatTypeDate     FormatType = "date"
	FormatTypeTime     FormatType = "time"
	FormatTypeDateTime FormatType = "datetime"
	FormatTypeUUID     FormatType = "uuid"
	FormatTypeIPv4     FormatType = "ipv4"
	FormatTypeIPv6     FormatType = "ipv6"
)

// RangeConstraint represents a range constraint value
type RangeConstraint struct {
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	Inclusive *bool    `json:"inclusive,omitempty"`
}

// Constraint represents a validation constraint for ASTRA fields and values
type Constraint struct {
	// Type of constraint being applied
	Type ConstraintType `json:"type"`
	// Constraint value (varies by constraint type)
	Value interface{} `json:"value,omitempty"`
	// Human-readable error message when constraint is violated
	Message *string `json:"message,omitempty"`
	// Machine-readable error code for constraint violations
	Code *string `json:"code,omitempty"`
}

// ============================================================================
// Participant Types
// ============================================================================

// ParticipantType represents the type of participant
type ParticipantType string

const (
	ParticipantTypeHuman  ParticipantType = "human"
	ParticipantTypeAI     ParticipantType = "ai"
	ParticipantTypeSystem ParticipantType = "system"
	ParticipantTypeBot    ParticipantType = "bot"
)

// ParticipantPreferences contains participant preferences
type ParticipantPreferences struct {
	// Preferred language code (ISO 639-1 with optional region)
	Language *string `json:"language,omitempty"`
	// Preferred timezone (IANA timezone identifier)
	Timezone *string `json:"timezone,omitempty"`
	// Preferred communication channels in order of preference
	CommunicationChannels []string `json:"communication_channels,omitempty"`
	// Additional preferences
	AdditionalProperties map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for ParticipantPreferences
func (p ParticipantPreferences) MarshalJSON() ([]byte, error) {
	base := map[string]interface{}{
		"language":               p.Language,
		"timezone":               p.Timezone,
		"communication_channels": p.CommunicationChannels,
	}
	
	for k, v := range p.AdditionalProperties {
		base[k] = v
	}
	
	return json.Marshal(base)
}

// UnmarshalJSON implements custom JSON unmarshaling for ParticipantPreferences
func (p *ParticipantPreferences) UnmarshalJSON(data []byte) error {
	type Alias ParticipantPreferences
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// Extract additional properties
	p.AdditionalProperties = make(map[string]interface{})
	for k, v := range raw {
		switch k {
		case "language", "timezone", "communication_channels":
			// Skip known fields
		default:
			p.AdditionalProperties[k] = v
		}
	}
	
	return nil
}

// Participant represents a conversation participant in ASTRA conversations
type Participant struct {
	// Unique identifier for this participant
	ID string `json:"id"`
	// Type of participant
	Type ParticipantType `json:"type"`
	// Business role of the participant (customer, agent, manager, etc.)
	Role *string `json:"role,omitempty"`
	// Display name of the participant
	Name *string `json:"name,omitempty"`
	// Email address of the participant
	Email *string `json:"email,omitempty"`
	// Phone number of the participant
	Phone *string `json:"phone,omitempty"`
	// External system identifier for this participant
	ExternalID *string `json:"external_id,omitempty"`
	// External system that manages this participant
	System *string `json:"system,omitempty"`
	// List of capabilities this participant has
	Capabilities []string `json:"capabilities,omitempty"`
	// List of permissions granted to this participant
	Permissions []string `json:"permissions,omitempty"`
	// Participant preferences
	Preferences *ParticipantPreferences `json:"preferences,omitempty"`
	// Additional participant metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ============================================================================
// Act Types
// ============================================================================

// ExpectedType represents the expected data type of a response
type ExpectedType string

const (
	ExpectedTypeString  ExpectedType = "string"
	ExpectedTypeNumber  ExpectedType = "number"
	ExpectedTypeBoolean ExpectedType = "boolean"
	ExpectedTypeObject  ExpectedType = "object"
	ExpectedTypeArray   ExpectedType = "array"
	ExpectedTypeDate    ExpectedType = "date"
	ExpectedTypeEmail   ExpectedType = "email"
	ExpectedTypePhone   ExpectedType = "phone"
	ExpectedTypeAddress ExpectedType = "address"
)

// Ask represents an act that requests missing information required to complete a business process
type Ask struct {
	Act
	// Field or information being requested
	Field string `json:"field"`
	// Question or request presented to obtain the information
	Prompt string `json:"prompt"`
	// Validation constraints for the requested information
	Constraints []Constraint `json:"constraints,omitempty"`
	// Whether this information is required to proceed
	Required *bool `json:"required,omitempty"`
	// Expected data type of the response
	ExpectedType *ExpectedType `json:"expected_type,omitempty"`
	// Number of times this question has been asked
	RetryCount *int `json:"retry_count,omitempty"`
	// Maximum number of retry attempts before escalation
	MaxRetries *int `json:"max_retries,omitempty"`
}

// Validate implements ConversationAct interface
func (a Ask) Validate() error {
	if a.Field == "" {
		return ValidationError{Field: "field", Message: "field is required", Value: a.Field}
	}
	if a.Prompt == "" {
		return ValidationError{Field: "prompt", Message: "prompt is required", Value: a.Prompt}
	}
	if a.RetryCount != nil && *a.RetryCount < 0 {
		return ValidationError{Field: "retry_count", Message: "retry_count cannot be negative", Value: *a.RetryCount}
	}
	if a.MaxRetries != nil && *a.MaxRetries < 0 {
		return ValidationError{Field: "max_retries", Message: "max_retries cannot be negative", Value: *a.MaxRetries}
	}
	return nil
}

// FieldOperation represents the operation being performed on a field
type FieldOperation string

const (
	FieldOperationSet       FieldOperation = "set"
	FieldOperationAppend    FieldOperation = "append"
	FieldOperationIncrement FieldOperation = "increment"
	FieldOperationDecrement FieldOperation = "decrement"
	FieldOperationDelete    FieldOperation = "delete"
	FieldOperationMerge     FieldOperation = "merge"
)

// ValidationStatus represents the validation status of a fact
type ValidationStatus string

const (
	ValidationStatusPending ValidationStatus = "pending"
	ValidationStatusValid   ValidationStatus = "valid"
	ValidationStatusInvalid ValidationStatus = "invalid"
	ValidationStatusPartial ValidationStatus = "partial"
)

// Fact represents an act that declares facts or information provided during conversation
type Fact struct {
	Act
	// Business entity being modified (order, customer, appointment, etc.)
	Entity EntityRef `json:"entity"`
	// Specific field or property being set
	Field string `json:"field"`
	// Value being assigned to the field
	Value interface{} `json:"value"`
	// Operation being performed on the field
	Operation *FieldOperation `json:"operation,omitempty"`
	// Previous value of the field (for audit trail)
	PreviousValue interface{} `json:"previous_value,omitempty"`
	// Validation status of this fact
	ValidationStatus *ValidationStatus `json:"validation_status,omitempty"`
	// List of validation errors if validation_status is invalid
	ValidationErrors []string `json:"validation_errors,omitempty"`
}

// Validate implements ConversationAct interface
func (f Fact) Validate() error {
	if f.Entity == nil {
		return ValidationError{Field: "entity", Message: "entity is required", Value: f.Entity}
	}
	if f.Field == "" {
		return ValidationError{Field: "field", Message: "field is required", Value: f.Field}
	}
	if f.Value == nil {
		return ValidationError{Field: "value", Message: "value is required", Value: f.Value}
	}
	return nil
}

// ConfirmationMethod represents how a confirmation was obtained
type ConfirmationMethod string

const (
	ConfirmationMethodVerbal   ConfirmationMethod = "verbal"
	ConfirmationMethodExplicit ConfirmationMethod = "explicit"
	ConfirmationMethodImplicit ConfirmationMethod = "implicit"
	ConfirmationMethodTimeout  ConfirmationMethod = "timeout"
	ConfirmationMethodSystem   ConfirmationMethod = "system"
)

// Confirm represents an act that verifies understanding of information before commitment
type Confirm struct {
	Act
	// Business entity being confirmed
	Entity EntityRef `json:"entity"`
	// Human-readable summary of what is being confirmed
	Summary string `json:"summary"`
	// Whether confirmation is still pending
	Awaiting *bool `json:"awaiting,omitempty"`
	// Whether the confirmation was accepted (true) or rejected (false)
	Confirmed *bool `json:"confirmed,omitempty"`
	// How the confirmation was obtained
	ConfirmationMethod *ConfirmationMethod `json:"confirmation_method,omitempty"`
	// Specific fields or aspects being confirmed
	FieldsConfirmed []string `json:"fields_confirmed,omitempty"`
	// Reason provided if confirmation was rejected
	RejectionReason *string `json:"rejection_reason,omitempty"`
	// Timeout for awaiting confirmation in milliseconds
	TimeoutMs *int64 `json:"timeout_ms,omitempty"`
}

// Validate implements ConversationAct interface
func (c Confirm) Validate() error {
	if c.Entity == nil {
		return ValidationError{Field: "entity", Message: "entity is required", Value: c.Entity}
	}
	if c.Summary == "" {
		return ValidationError{Field: "summary", Message: "summary is required", Value: c.Summary}
	}
	if c.TimeoutMs != nil && *c.TimeoutMs < 0 {
		return ValidationError{Field: "timeout_ms", Message: "timeout_ms cannot be negative", Value: *c.TimeoutMs}
	}
	return nil
}

// CommitAction represents the action being performed in the target system
type CommitAction string

const (
	CommitActionCreate  CommitAction = "create"
	CommitActionUpdate  CommitAction = "update"
	CommitActionDelete  CommitAction = "delete"
	CommitActionExecute CommitAction = "execute"
	CommitActionCancel  CommitAction = "cancel"
	CommitActionPause   CommitAction = "pause"
	CommitActionResume  CommitAction = "resume"
)

// CommitStatus represents the status of a commit operation
type CommitStatus string

const (
	CommitStatusPending    CommitStatus = "pending"
	CommitStatusInProgress CommitStatus = "in_progress"
	CommitStatusSuccess    CommitStatus = "success"
	CommitStatusFailed     CommitStatus = "failed"
	CommitStatusRetrying   CommitStatus = "retrying"
	CommitStatusCancelled  CommitStatus = "cancelled"
)

// CommitError represents error information for failed commits
type CommitError struct {
	// Error code from the target system
	Code string `json:"code"`
	// Human-readable error message
	Message string `json:"message"`
	// Additional error context
	Details map[string]interface{} `json:"details,omitempty"`
	// Whether the error can be recovered from
	Recoverable bool `json:"recoverable"`
}

// Commit represents an act that executes business processes and triggers system integrations
type Commit struct {
	Act
	// Business entity being committed to external systems
	Entity EntityRef `json:"entity"`
	// Action being performed in the target system
	Action CommitAction `json:"action"`
	// Target system identifier (CRM, order_management, etc.)
	System *string `json:"system,omitempty"`
	// External system transaction or record identifier
	TransactionID *string `json:"transaction_id,omitempty"`
	// Status of the commit operation
	Status *CommitStatus `json:"status,omitempty"`
	// Error information if commit failed
	Error *CommitError `json:"error,omitempty"`
	// Number of retry attempts made
	RetryCount *int `json:"retry_count,omitempty"`
	// Maximum number of retry attempts
	MaxRetries *int `json:"max_retries,omitempty"`
	// Key to ensure idempotent operations
	IdempotencyKey *string `json:"idempotency_key,omitempty"`
	// Information needed to rollback this commit if necessary
	RollbackInfo map[string]interface{} `json:"rollback_info,omitempty"`
}

// Validate implements ConversationAct interface
func (c Commit) Validate() error {
	if c.Entity == nil {
		return ValidationError{Field: "entity", Message: "entity is required", Value: c.Entity}
	}
	if c.Action == "" {
		return ValidationError{Field: "action", Message: "action is required", Value: c.Action}
	}
	if c.RetryCount != nil && *c.RetryCount < 0 {
		return ValidationError{Field: "retry_count", Message: "retry_count cannot be negative", Value: *c.RetryCount}
	}
	if c.MaxRetries != nil && *c.MaxRetries < 0 {
		return ValidationError{Field: "max_retries", Message: "max_retries cannot be negative", Value: *c.MaxRetries}
	}
	return nil
}

// ErrorSeverity represents the severity level of an error
type ErrorSeverity string

const (
	ErrorSeverityInfo     ErrorSeverity = "info"
	ErrorSeverityWarning  ErrorSeverity = "warning"
	ErrorSeverityError    ErrorSeverity = "error"
	ErrorSeverityCritical ErrorSeverity = "critical"
)

// ErrorCategory represents the category of error for classification
type ErrorCategory string

const (
	ErrorCategoryValidation   ErrorCategory = "validation"
	ErrorCategoryProcessing   ErrorCategory = "processing"
	ErrorCategoryIntegration  ErrorCategory = "integration"
	ErrorCategoryTimeout      ErrorCategory = "timeout"
	ErrorCategoryPermission   ErrorCategory = "permission"
	ErrorCategorySystem       ErrorCategory = "system"
	ErrorCategoryUserInput    ErrorCategory = "user_input"
	ErrorCategoryBusinessRule ErrorCategory = "business_rule"
)

// SuggestedAction represents a suggested recovery action
type SuggestedAction string

const (
	SuggestedActionRetry     SuggestedAction = "retry"
	SuggestedActionEscalate  SuggestedAction = "escalate"
	SuggestedActionIgnore    SuggestedAction = "ignore"
	SuggestedActionClarify   SuggestedAction = "clarify"
	SuggestedActionFallback  SuggestedAction = "fallback"
	SuggestedActionTerminate SuggestedAction = "terminate"
)

// Error represents an act that handles failures and exceptions in conversational processing
type Error struct {
	Act
	// Machine-readable error code
	Code string `json:"code"`
	// Human-readable error message
	Message string `json:"message"`
	// Whether the conversation can continue after this error
	Recoverable bool `json:"recoverable"`
	// Severity level of the error
	Severity *ErrorSeverity `json:"severity,omitempty"`
	// Category of error for classification
	Category *ErrorCategory `json:"category,omitempty"`
	// Additional error context and debugging information
	Details map[string]interface{} `json:"details,omitempty"`
	// ID of the act that caused this error
	RelatedActID *string `json:"related_act_id,omitempty"`
	// Suggested recovery action
	SuggestedAction *SuggestedAction `json:"suggested_action,omitempty"`
	// User-friendly message to display to conversation participants
	UserMessage *string `json:"user_message,omitempty"`
	// Technical stack trace for debugging (not shown to users)
	StackTrace *string `json:"stack_trace,omitempty"`
}

// Validate implements ConversationAct interface
func (e Error) Validate() error {
	if e.Code == "" {
		return ValidationError{Field: "code", Message: "code is required", Value: e.Code}
	}
	if e.Message == "" {
		return ValidationError{Field: "message", Message: "message is required", Value: e.Message}
	}
	if e.RelatedActID != nil && !IsValidActID(*e.RelatedActID) {
		return ValidationError{Field: "related_act_id", Message: "invalid act ID format", Value: *e.RelatedActID}
	}
	return nil
}

// ============================================================================
// Conversation Types
// ============================================================================

// ConversationStatus represents the current status of a conversation
type ConversationStatus string

const (
	ConversationStatusActive    ConversationStatus = "active"
	ConversationStatusPaused    ConversationStatus = "paused"
	ConversationStatusCompleted ConversationStatus = "completed"
	ConversationStatusFailed    ConversationStatus = "failed"
	ConversationStatusCancelled ConversationStatus = "cancelled"
)

// ConversationContext contains conversation context and session information
type ConversationContext struct {
	// Session identifier
	SessionID *string `json:"session_id,omitempty"`
	// User agent or client information
	UserAgent *string `json:"user_agent,omitempty"`
	// Client IP address
	IPAddress *string `json:"ip_address,omitempty"`
	// How the conversation was initiated
	Referrer *string `json:"referrer,omitempty"`
	// Additional context properties
	AdditionalProperties map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for ConversationContext
func (c ConversationContext) MarshalJSON() ([]byte, error) {
	base := map[string]interface{}{
		"session_id": c.SessionID,
		"user_agent": c.UserAgent,
		"ip_address": c.IPAddress,
		"referrer":   c.Referrer,
	}
	
	for k, v := range c.AdditionalProperties {
		base[k] = v
	}
	
	return json.Marshal(base)
}

// UnmarshalJSON implements custom JSON unmarshaling for ConversationContext
func (c *ConversationContext) UnmarshalJSON(data []byte) error {
	type Alias ConversationContext
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// Extract additional properties
	c.AdditionalProperties = make(map[string]interface{})
	for k, v := range raw {
		switch k {
		case "session_id", "user_agent", "ip_address", "referrer":
			// Skip known fields
		default:
			c.AdditionalProperties[k] = v
		}
	}
	
	return nil
}

// ConversationMetadata contains additional conversation metadata
type ConversationMetadata struct {
	// Total conversation duration in milliseconds
	TotalDurationMs *int64 `json:"total_duration_ms,omitempty"`
	// Total number of acts in the conversation
	ActCount *int `json:"act_count,omitempty"`
	// Number of errors that occurred
	ErrorCount *int `json:"error_count,omitempty"`
	// Number of successful commits
	CommitCount *int `json:"commit_count,omitempty"`
	// Average confidence score across all acts
	AvgConfidence *float64 `json:"avg_confidence,omitempty"`
	// Additional metadata properties
	AdditionalProperties map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for ConversationMetadata
func (m ConversationMetadata) MarshalJSON() ([]byte, error) {
	base := map[string]interface{}{
		"total_duration_ms": m.TotalDurationMs,
		"act_count":         m.ActCount,
		"error_count":       m.ErrorCount,
		"commit_count":      m.CommitCount,
		"avg_confidence":    m.AvgConfidence,
	}
	
	for k, v := range m.AdditionalProperties {
		base[k] = v
	}
	
	return json.Marshal(base)
}

// UnmarshalJSON implements custom JSON unmarshaling for ConversationMetadata
func (m *ConversationMetadata) UnmarshalJSON(data []byte) error {
	type Alias ConversationMetadata
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// Extract additional properties
	m.AdditionalProperties = make(map[string]interface{})
	for k, v := range raw {
		switch k {
		case "total_duration_ms", "act_count", "error_count", "commit_count", "avg_confidence":
			// Skip known fields
		default:
			m.AdditionalProperties[k] = v
		}
	}
	
	return nil
}

// Conversation represents a complete ASTRA conversation container with acts and metadata
type Conversation struct {
	// Unique identifier for this conversation
	ID string `json:"id"`
	// List of conversation participants
	Participants []Participant `json:"participants"`
	// Ordered sequence of acts in this conversation
	Acts []ConversationAct `json:"acts"`
	// When the conversation started
	StartedAt *time.Time `json:"started_at,omitempty"`
	// When the conversation ended
	EndedAt *time.Time `json:"ended_at,omitempty"`
	// Current status of the conversation
	Status *ConversationStatus `json:"status,omitempty"`
	// Primary communication channel for this conversation
	Channel *string `json:"channel,omitempty"`
	// Business schema identifier used for this conversation
	Schema *string `json:"schema,omitempty"`
	// Conversation context and session information
	Context *ConversationContext `json:"context,omitempty"`
	// Final computed state of all entities after processing all acts
	FinalState map[string]interface{} `json:"final_state,omitempty"`
	// Additional conversation metadata
	Metadata *ConversationMetadata `json:"metadata,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for Conversation
func (c Conversation) MarshalJSON() ([]byte, error) {
	type Alias Conversation
	
	// Convert acts to ActUnion for JSON serialization
	acts := make([]ActUnion, len(c.Acts))
	for i, act := range c.Acts {
		acts[i] = NewActUnion(act)
	}
	
	return json.Marshal(&struct {
		Acts []ActUnion `json:"acts"`
		*Alias
	}{
		Acts:  acts,
		Alias: (*Alias)(&c),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for Conversation
func (c *Conversation) UnmarshalJSON(data []byte) error {
	type Alias Conversation
	aux := &struct {
		Acts []json.RawMessage `json:"acts"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// Unmarshal acts
	c.Acts = make([]ConversationAct, len(aux.Acts))
	for i, actData := range aux.Acts {
		act, err := UnmarshalAct(actData)
		if err != nil {
			return fmt.Errorf("failed to unmarshal act at index %d: %w", i, err)
		}
		c.Acts[i] = act
	}
	
	return nil
}
