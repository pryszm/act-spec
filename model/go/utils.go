package astra

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// ============================================================================
// ID Generation and Validation
// ============================================================================

// actIDPattern is the regex pattern for valid ASTRA act IDs
var actIDPattern = regexp.MustCompile(`^act_[a-zA-Z0-9_-]+$`)

// conversationIDPattern is the regex pattern for valid ASTRA conversation IDs
var conversationIDPattern = regexp.MustCompile(`^conv_[a-zA-Z0-9_-]+$`)

// GenerateActID generates a new ASTRA-compliant act ID
func GenerateActID() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 36)
	random := generateRandomString(8)
	return fmt.Sprintf("act_%s_%s", timestamp, random)
}

// GenerateConversationID generates a new ASTRA-compliant conversation ID
func GenerateConversationID() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 36)
	random := generateRandomString(8)
	return fmt.Sprintf("conv_%s_%s", timestamp, random)
}

// GenerateParticipantID generates a new participant ID
func GenerateParticipantID() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 36)
	random := generateRandomString(6)
	return fmt.Sprintf("participant_%s_%s", timestamp, random)
}

// GenerateEntityID generates a new entity ID
func GenerateEntityID(entityType string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 36)
	random := generateRandomString(6)
	if entityType == "" {
		entityType = "entity"
	}
	return fmt.Sprintf("%s_%s_%s", entityType, timestamp, random)
}

// generateRandomString generates a cryptographically secure random string
func generateRandomString(length int) string {
	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based generation if crypto/rand fails
		return strconv.FormatInt(time.Now().UnixNano(), 36)[:length]
	}
	return hex.EncodeToString(bytes)[:length]
}

// IsValidActID validates an act ID format
func IsValidActID(id string) bool {
	return actIDPattern.MatchString(id)
}

// IsValidConversationID validates a conversation ID format
func IsValidConversationID(id string) bool {
	return conversationIDPattern.MatchString(id)
}

// IsValidTimestamp validates if a time.Time is not zero
func IsValidTimestamp(t time.Time) bool {
	return !t.IsZero()
}

// ============================================================================
// Act Creation Utilities
// ============================================================================

// CreateBaseAct creates a basic Act with required fields filled in
func CreateBaseAct(speaker string, actType ActType, options ...ActOption) Act {
	act := Act{
		ID:        GenerateActID(),
		Timestamp: time.Now(),
		Speaker:   speaker,
		Type:      actType,
	}
	
	// Apply options
	for _, option := range options {
		option(&act)
	}
	
	return act
}

// ActOption is a function type for configuring Act creation
type ActOption func(*Act)

// WithConfidence sets the confidence score for an act
func WithConfidence(confidence float64) ActOption {
	return func(a *Act) {
		if confidence >= 0.0 && confidence <= 1.0 {
			a.Confidence = &confidence
		}
	}
}

// WithSource sets the source for an act
func WithSource(source Source) ActOption {
	return func(a *Act) {
		a.Source = &source
	}
}

// WithMetadata sets metadata for an act
func WithMetadata(metadata ActMetadata) ActOption {
	return func(a *Act) {
		a.Metadata = &metadata
	}
}

// WithChannel sets the communication channel in metadata
func WithChannel(channel string) ActOption {
	return func(a *Act) {
		if a.Metadata == nil {
			a.Metadata = &ActMetadata{}
		}
		a.Metadata.Channel = &channel
	}
}

// WithLanguage sets the language in metadata
func WithLanguage(language string) ActOption {
	return func(a *Act) {
		if a.Metadata == nil {
			a.Metadata = &ActMetadata{}
		}
		a.Metadata.Language = &language
	}
}

// WithOriginalText sets the original text in metadata
func WithOriginalText(originalText string) ActOption {
	return func(a *Act) {
		if a.Metadata == nil {
			a.Metadata = &ActMetadata{}
		}
		a.Metadata.OriginalText = &originalText
	}
}

// ============================================================================
// Specialized Act Builders
// ============================================================================

// NewAsk creates a new Ask act with required fields
func NewAsk(speaker, field, prompt string, options ...AskOption) Ask {
	ask := Ask{
		Act:    CreateBaseAct(speaker, ActTypeAsk),
		Field:  field,
		Prompt: prompt,
	}
	
	// Apply Ask-specific options
	for _, option := range options {
		option(&ask)
	}
	
	return ask
}

// AskOption is a function type for configuring Ask creation
type AskOption func(*Ask)

// WithRequired sets whether the ask is required
func WithRequired(required bool) AskOption {
	return func(a *Ask) {
		a.Required = &required
	}
}

// WithExpectedType sets the expected response type
func WithExpectedType(expectedType ExpectedType) AskOption {
	return func(a *Ask) {
		a.ExpectedType = &expectedType
	}
}

// WithConstraints sets validation constraints
func WithConstraints(constraints []Constraint) AskOption {
	return func(a *Ask) {
		a.Constraints = constraints
	}
}

// WithMaxRetries sets the maximum retry count
func WithMaxRetries(maxRetries int) AskOption {
	return func(a *Ask) {
		if maxRetries >= 0 {
			a.MaxRetries = &maxRetries
		}
	}
}

// NewFact creates a new Fact act with required fields
func NewFact(speaker string, entity EntityRef, field string, value interface{}, options ...FactOption) Fact {
	fact := Fact{
		Act:    CreateBaseAct(speaker, ActTypeFact),
		Entity: entity,
		Field:  field,
		Value:  value,
	}
	
	// Apply Fact-specific options
	for _, option := range options {
		option(&fact)
	}
	
	return fact
}

// FactOption is a function type for configuring Fact creation
type FactOption func(*Fact)

// WithOperation sets the field operation
func WithOperation(operation FieldOperation) FactOption {
	return func(f *Fact) {
		f.Operation = &operation
	}
}

// WithPreviousValue sets the previous value for audit trail
func WithPreviousValue(previousValue interface{}) FactOption {
	return func(f *Fact) {
		f.PreviousValue = previousValue
	}
}

// WithValidationStatus sets the validation status
func WithValidationStatus(status ValidationStatus) FactOption {
	return func(f *Fact) {
		f.ValidationStatus = &status
	}
}

// NewConfirm creates a new Confirm act with required fields
func NewConfirm(speaker string, entity EntityRef, summary string, options ...ConfirmOption) Confirm {
	confirm := Confirm{
		Act:     CreateBaseAct(speaker, ActTypeConfirm),
		Entity:  entity,
		Summary: summary,
	}
	
	// Apply Confirm-specific options
	for _, option := range options {
		option(&confirm)
	}
	
	return confirm
}

// ConfirmOption is a function type for configuring Confirm creation
type ConfirmOption func(*Confirm)

// WithAwaiting sets the awaiting status
func WithAwaiting(awaiting bool) ConfirmOption {
	return func(c *Confirm) {
		c.Awaiting = &awaiting
	}
}

// WithConfirmed sets the confirmed status
func WithConfirmed(confirmed bool) ConfirmOption {
	return func(c *Confirm) {
		c.Confirmed = &confirmed
	}
}

// WithConfirmationMethod sets the confirmation method
func WithConfirmationMethod(method ConfirmationMethod) ConfirmOption {
	return func(c *Confirm) {
		c.ConfirmationMethod = &method
	}
}

// WithTimeoutMs sets the confirmation timeout
func WithTimeoutMs(timeoutMs int64) ConfirmOption {
	return func(c *Confirm) {
		if timeoutMs >= 0 {
			c.TimeoutMs = &timeoutMs
		}
	}
}

// NewCommit creates a new Commit act with required fields
func NewCommit(speaker string, entity EntityRef, action CommitAction, options ...CommitOption) Commit {
	commit := Commit{
		Act:    CreateBaseAct(speaker, ActTypeCommit),
		Entity: entity,
		Action: action,
	}
	
	// Apply Commit-specific options
	for _, option := range options {
		option(&commit)
	}
	
	return commit
}

// CommitOption is a function type for configuring Commit creation
type CommitOption func(*Commit)

// WithSystem sets the target system
func WithSystem(system string) CommitOption {
	return func(c *Commit) {
		c.System = &system
	}
}

// WithTransactionID sets the transaction ID
func WithTransactionID(transactionID string) CommitOption {
	return func(c *Commit) {
		c.TransactionID = &transactionID
	}
}

// WithCommitStatus sets the commit status
func WithCommitStatus(status CommitStatus) CommitOption {
	return func(c *Commit) {
		c.Status = &status
	}
}

// WithIdempotencyKey sets the idempotency key
func WithIdempotencyKey(key string) CommitOption {
	return func(c *Commit) {
		c.IdempotencyKey = &key
	}
}

// NewError creates a new Error act with required fields
func NewError(speaker, code, message string, recoverable bool, options ...ErrorOption) Error {
	errorAct := Error{
		Act:         CreateBaseAct(speaker, ActTypeError),
		Code:        code,
		Message:     message,
		Recoverable: recoverable,
	}
	
	// Apply Error-specific options
	for _, option := range options {
		option(&errorAct)
	}
	
	return errorAct
}

// ErrorOption is a function type for configuring Error creation
type ErrorOption func(*Error)

// WithSeverity sets the error severity
func WithSeverity(severity ErrorSeverity) ErrorOption {
	return func(e *Error) {
		e.Severity = &severity
	}
}

// WithCategory sets the error category
func WithCategory(category ErrorCategory) ErrorOption {
	return func(e *Error) {
		e.Category = &category
	}
}

// WithSuggestedAction sets the suggested recovery action
func WithSuggestedAction(action SuggestedAction) ErrorOption {
	return func(e *Error) {
		e.SuggestedAction = &action
	}
}

// WithUserMessage sets a user-friendly error message
func WithUserMessage(userMessage string) ErrorOption {
	return func(e *Error) {
		e.UserMessage = &userMessage
	}
}

// WithRelatedActID sets the ID of the act that caused this error
func WithRelatedActID(actID string) ErrorOption {
	return func(e *Error) {
		if IsValidActID(actID) {
			e.RelatedActID = &actID
		}
	}
}

// ============================================================================
// Entity and Participant Utilities
// ============================================================================

// NewEntity creates a new Entity with required fields
func NewEntity(id, entityType string, options ...EntityOption) Entity {
	entity := Entity{
		ID:   id,
		Type: entityType,
	}
	
	// Apply options
	for _, option := range options {
		option(&entity)
	}
	
	return entity
}

// EntityOption is a function type for configuring Entity creation
type EntityOption func(*Entity)

// WithExternalID sets the external system ID
func WithExternalID(externalID string) EntityOption {
	return func(e *Entity) {
		e.ExternalID = &externalID
	}
}

// WithEntitySystem sets the owning system
func WithEntitySystem(system string) EntityOption {
	return func(e *Entity) {
		e.System = &system
	}
}

// WithVersion sets the entity version
func WithVersion(version string) EntityOption {
	return func(e *Entity) {
		e.Version = &version
	}
}

// WithSchemaURL sets the schema URL
func WithSchemaURL(schemaURL string) EntityOption {
	return func(e *Entity) {
		e.SchemaURL = &schemaURL
	}
}

// WithEntityMetadata sets entity metadata
func WithEntityMetadata(metadata map[string]interface{}) EntityOption {
	return func(e *Entity) {
		e.Metadata = metadata
	}
}

// NewParticipant creates a new Participant with required fields
func NewParticipant(id string, participantType ParticipantType, options ...ParticipantOption) Participant {
	participant := Participant{
		ID:   id,
		Type: participantType,
	}
	
	// Apply options
	for _, option := range options {
		option(&participant)
	}
	
	return participant
}

// ParticipantOption is a function type for configuring Participant creation
type ParticipantOption func(*Participant)

// WithRole sets the participant's business role
func WithRole(role string) ParticipantOption {
	return func(p *Participant) {
		p.Role = &role
	}
}

// WithName sets the participant's display name
func WithName(name string) ParticipantOption {
	return func(p *Participant) {
		p.Name = &name
	}
}

// WithEmail sets the participant's email address
func WithEmail(email string) ParticipantOption {
	return func(p *Participant) {
		p.Email = &email
	}
}

// WithPhone sets the participant's phone number
func WithPhone(phone string) ParticipantOption {
	return func(p *Participant) {
		p.Phone = &phone
	}
}

// WithCapabilities sets the participant's capabilities
func WithCapabilities(capabilities []string) ParticipantOption {
	return func(p *Participant) {
		p.Capabilities = capabilities
	}
}

// WithPermissions sets the participant's permissions
func WithPermissions(permissions []string) ParticipantOption {
	return func(p *Participant) {
		p.Permissions = permissions
	}
}

// WithPreferences sets the participant's preferences
func WithPreferences(preferences ParticipantPreferences) ParticipantOption {
	return func(p *Participant) {
		p.Preferences = &preferences
	}
}

// ============================================================================
// Conversation Utilities
// ============================================================================

// NewConversation creates a new Conversation with required fields
func NewConversation(participants []Participant, options ...ConversationOption) Conversation {
	conversation := Conversation{
		ID:           GenerateConversationID(),
		Participants: participants,
		Acts:         make([]ConversationAct, 0),
	}
	
	// Set started time
	now := time.Now()
	conversation.StartedAt = &now
	
	// Set default status
	status := ConversationStatusActive
	conversation.Status = &status
	
	// Apply options
	for _, option := range options {
		option(&conversation)
	}
	
	return conversation
}

// ConversationOption is a function type for configuring Conversation creation
type ConversationOption func(*Conversation)

// WithChannel sets the primary communication channel
func WithConversationChannel(channel string) ConversationOption {
	return func(c *Conversation) {
		c.Channel = &channel
	}
}

// WithSchema sets the business schema identifier
func WithConversationSchema(schema string) ConversationOption {
	return func(c *Conversation) {
		c.Schema = &schema
	}
}

// WithContext sets the conversation context
func WithContext(context ConversationContext) ConversationOption {
	return func(c *Conversation) {
		c.Context = &context
	}
}

// WithConversationStatus sets the conversation status
func WithConversationStatus(status ConversationStatus) ConversationOption {
	return func(c *Conversation) {
		c.Status = &status
	}
}

// AddAct adds an act to a conversation and returns the updated conversation
func (c *Conversation) AddAct(act ConversationAct) error {
	// Validate the act
	if err := ValidateAct(act); err != nil {
		return fmt.Errorf("invalid act: %w", err)
	}
	
	// Add to acts slice
	c.Acts = append(c.Acts, act)
	
	// Update metadata
	c.updateMetadata()
	
	return nil
}

// EndConversation marks a conversation as ended
func (c *Conversation) EndConversation(status ConversationStatus) {
	now := time.Now()
	c.EndedAt = &now
	c.Status = &status
	c.updateMetadata()
}

// updateMetadata updates the conversation metadata based on current state
func (c *Conversation) updateMetadata() {
	if c.Metadata == nil {
		c.Metadata = &ConversationMetadata{}
	}
	
	// Update act count
	actCount := len(c.Acts)
	c.Metadata.ActCount = &actCount
	
	// Count errors and commits
	errorCount := 0
	commitCount := 0
	var totalConfidence float64
	confidenceCount := 0
	
	for _, act := range c.Acts {
		baseAct := act.GetAct()
		
		switch act.GetType() {
		case ActTypeError:
			errorCount++
		case ActTypeCommit:
			if commitAct, ok := act.(Commit); ok {
				if commitAct.Status != nil && *commitAct.Status == CommitStatusSuccess {
					commitCount++
				}
			}
		}
		
		// Accumulate confidence scores
		if baseAct.Confidence != nil {
			totalConfidence += *baseAct.Confidence
			confidenceCount++
		}
	}
	
	c.Metadata.ErrorCount = &errorCount
	c.Metadata.CommitCount = &commitCount
	
	// Calculate average confidence
	if confidenceCount > 0 {
		avgConfidence := totalConfidence / float64(confidenceCount)
		c.Metadata.AvgConfidence = &avgConfidence
	}
	
	// Calculate duration if conversation has ended
	if c.StartedAt != nil && c.EndedAt != nil {
		duration := c.EndedAt.Sub(*c.StartedAt).Milliseconds()
		c.Metadata.TotalDurationMs = &duration
	}
}

// GetActsByType returns all acts of a specific type from the conversation
func (c *Conversation) GetActsByType(actType ActType) []ConversationAct {
	var acts []ConversationAct
	for _, act := range c.Acts {
		if act.GetType() == actType {
			acts = append(acts, act)
		}
	}
	return acts
}

// GetActsBySpeaker returns all acts from a specific speaker
func (c *Conversation) GetActsBySpeaker(speaker string) []ConversationAct {
	var acts []ConversationAct
	for _, act := range c.Acts {
		if act.GetAct().Speaker == speaker {
			acts = append(acts, act)
		}
	}
	return acts
}

// GetParticipantByID finds a participant by their ID
func (c *Conversation) GetParticipantByID(id string) *Participant {
	for i := range c.Participants {
		if c.Participants[i].ID == id {
			return &c.Participants[i]
		}
	}
	return nil
}

// ============================================================================
// Constraint Utilities
// ============================================================================

// NewConstraint creates a new validation constraint
func NewConstraint(constraintType ConstraintType, options ...ConstraintOption) Constraint {
	constraint := Constraint{
		Type: constraintType,
	}
	
	// Apply options
	for _, option := range options {
		option(&constraint)
	}
	
	return constraint
}

// ConstraintOption is a function type for configuring Constraint creation
type ConstraintOption func(*Constraint)

// WithValue sets the constraint value
func WithConstraintValue(value interface{}) ConstraintOption {
	return func(c *Constraint) {
		c.Value = value
	}
}

// WithMessage sets the constraint error message
func WithConstraintMessage(message string) ConstraintOption {
	return func(c *Constraint) {
		c.Message = &message
	}
}

// WithCode sets the constraint error code
func WithConstraintCode(code string) ConstraintOption {
	return func(c *Constraint) {
		c.Code = &code
	}
}

// Common constraint builders

// RequiredConstraint creates a required field constraint
func RequiredConstraint() Constraint {
	return NewConstraint(ConstraintTypeRequired,
		WithConstraintMessage("This field is required"))
}

// MinLengthConstraint creates a minimum length constraint
func MinLengthConstraint(minLength int) Constraint {
	return NewConstraint(ConstraintTypeMinLength,
		WithConstraintValue(minLength),
		WithConstraintMessage(fmt.Sprintf("Minimum length is %d characters", minLength)))
}

// MaxLengthConstraint creates a maximum length constraint
func MaxLengthConstraint(maxLength int) Constraint {
	return NewConstraint(ConstraintTypeMaxLength,
		WithConstraintValue(maxLength),
		WithConstraintMessage(fmt.Sprintf("Maximum length is %d characters", maxLength)))
}

// EmailFormatConstraint creates an email format constraint
func EmailFormatConstraint() Constraint {
	return NewConstraint(ConstraintTypeFormat,
		WithConstraintValue(FormatTypeEmail),
		WithConstraintMessage("Must be a valid email address"))
}

// PhoneFormatConstraint creates a phone format constraint
func PhoneFormatConstraint() Constraint {
	return NewConstraint(ConstraintTypeFormat,
		WithConstraintValue(FormatTypePhone),
		WithConstraintMessage("Must be a valid phone number"))
}

// EnumConstraint creates an enumeration constraint
func EnumConstraint(values []string) Constraint {
	return NewConstraint(ConstraintTypeEnum,
		WithConstraintValue(values),
		WithConstraintMessage(fmt.Sprintf("Must be one of: %v", values)))
}

// RangeConstraint creates a numeric range constraint
func RangeConstraint(min, max *float64, inclusive bool) Constraint {
	rangeValue := RangeConstraint{
		Min:       min,
		Max:       max,
		Inclusive: &inclusive,
	}
	
	return NewConstraint(ConstraintTypeRange,
		WithConstraintValue(rangeValue),
		WithConstraintMessage("Value must be within the specified range"))
}
