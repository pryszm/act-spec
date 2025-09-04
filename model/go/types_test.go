package astra

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// ID Generation and Validation Tests
// ============================================================================

func TestGenerateActID(t *testing.T) {
	id1 := GenerateActID()
	id2 := GenerateActID()

	// Should generate valid act IDs
	assert.True(t, IsValidActID(id1), "Generated act ID should be valid")
	assert.True(t, IsValidActID(id2), "Generated act ID should be valid")

	// Should be unique
	assert.NotEqual(t, id1, id2, "Generated act IDs should be unique")

	// Should match expected pattern
	assert.Regexp(t, `^act_[a-zA-Z0-9_-]+$`, id1)
	assert.Regexp(t, `^act_[a-zA-Z0-9_-]+$`, id2)
}

func TestGenerateConversationID(t *testing.T) {
	id1 := GenerateConversationID()
	id2 := GenerateConversationID()

	// Should generate valid conversation IDs
	assert.True(t, IsValidConversationID(id1), "Generated conversation ID should be valid")
	assert.True(t, IsValidConversationID(id2), "Generated conversation ID should be valid")

	// Should be unique
	assert.NotEqual(t, id1, id2, "Generated conversation IDs should be unique")

	// Should match expected pattern
	assert.Regexp(t, `^conv_[a-zA-Z0-9_-]+$`, id1)
	assert.Regexp(t, `^conv_[a-zA-Z0-9_-]+$`, id2)
}

func TestIsValidActID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"Valid act ID", "act_123abc", true},
		{"Valid act ID with underscores", "act_123_abc_def", true},
		{"Valid act ID with dashes", "act_123-abc-def", true},
		{"Invalid - no prefix", "123abc", false},
		{"Invalid - wrong prefix", "action_123abc", false},
		{"Invalid - spaces", "act_123 abc", false},
		{"Invalid - special chars", "act_123@abc", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidActID(tt.id)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidConversationID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"Valid conversation ID", "conv_123abc", true},
		{"Valid conversation ID with underscores", "conv_123_abc_def", true},
		{"Valid conversation ID with dashes", "conv_123-abc-def", true},
		{"Invalid - no prefix", "123abc", false},
		{"Invalid - wrong prefix", "conversation_123abc", false},
		{"Invalid - spaces", "conv_123 abc", false},
		{"Invalid - special chars", "conv_123@abc", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidConversationID(tt.id)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// Act Creation and Validation Tests
// ============================================================================

func TestCreateBaseAct(t *testing.T) {
	speaker := "agent_123"
	actType := ActTypeAsk
	confidence := 0.95
	source := SourceAI

	act := CreateBaseAct(speaker, actType,
		WithConfidence(confidence),
		WithSource(source))

	assert.Equal(t, speaker, act.Speaker)
	assert.Equal(t, actType, act.Type)
	assert.True(t, IsValidActID(act.ID))
	assert.False(t, act.Timestamp.IsZero())
	assert.NotNil(t, act.Confidence)
	assert.Equal(t, confidence, *act.Confidence)
	assert.NotNil(t, act.Source)
	assert.Equal(t, source, *act.Source)
}

func TestNewAsk(t *testing.T) {
	speaker := "agent_123"
	field := "email"
	prompt := "What is your email address?"
	required := true
	expectedType := ExpectedTypeEmail

	ask := NewAsk(speaker, field, prompt,
		WithRequired(required),
		WithExpectedType(expectedType))

	assert.Equal(t, ActTypeAsk, ask.Type)
	assert.Equal(t, speaker, ask.Speaker)
	assert.Equal(t, field, ask.Field)
	assert.Equal(t, prompt, ask.Prompt)
	assert.NotNil(t, ask.Required)
	assert.Equal(t, required, *ask.Required)
	assert.NotNil(t, ask.ExpectedType)
	assert.Equal(t, expectedType, *ask.ExpectedType)
	assert.True(t, IsValidActID(ask.ID))

	// Test validation
	err := ask.Validate()
	assert.NoError(t, err)
}

func TestNewFact(t *testing.T) {
	speaker := "customer_456"
	entity := "order_789"
	field := "email"
	value := "user@example.com"
	operation := FieldOperationSet

	fact := NewFact(speaker, entity, field, value,
		WithOperation(operation))

	assert.Equal(t, ActTypeFact, fact.Type)
	assert.Equal(t, speaker, fact.Speaker)
	assert.Equal(t, entity, fact.Entity)
	assert.Equal(t, field, fact.Field)
	assert.Equal(t, value, fact.Value)
	assert.NotNil(t, fact.Operation)
	assert.Equal(t, operation, *fact.Operation)
	assert.True(t, IsValidActID(fact.ID))

	// Test validation
	err := fact.Validate()
	assert.NoError(t, err)
}

func TestNewConfirm(t *testing.T) {
	speaker := "agent_123"
	entity := "order_456"
	summary := "Order for 2 pizzas to be delivered at 6 PM"
	awaiting := true

	confirm := NewConfirm(speaker, entity, summary,
		WithAwaiting(awaiting))

	assert.Equal(t, ActTypeConfirm, confirm.Type)
	assert.Equal(t, speaker, confirm.Speaker)
	assert.Equal(t, entity, confirm.Entity)
	assert.Equal(t, summary, confirm.Summary)
	assert.NotNil(t, confirm.Awaiting)
	assert.Equal(t, awaiting, *confirm.Awaiting)
	assert.True(t, IsValidActID(confirm.ID))

	// Test validation
	err := confirm.Validate()
	assert.NoError(t, err)
}

func TestNewCommit(t *testing.T) {
	speaker := "system_001"
	entity := "order_456"
	action := CommitActionCreate
	system := "order_management"

	commit := NewCommit(speaker, entity, action,
		WithSystem(system))

	assert.Equal(t, ActTypeCommit, commit.Type)
	assert.Equal(t, speaker, commit.Speaker)
	assert.Equal(t, entity, commit.Entity)
	assert.Equal(t, action, commit.Action)
	assert.NotNil(t, commit.System)
	assert.Equal(t, system, *commit.System)
	assert.True(t, IsValidActID(commit.ID))

	// Test validation
	err := commit.Validate()
	assert.NoError(t, err)
}

func TestNewError(t *testing.T) {
	speaker := "system_001"
	code := "VALIDATION_ERROR"
	message := "Invalid email format"
	recoverable := true
	severity := ErrorSeverityWarning

	errorAct := NewError(speaker, code, message, recoverable,
		WithSeverity(severity))

	assert.Equal(t, ActTypeError, errorAct.Type)
	assert.Equal(t, speaker, errorAct.Speaker)
	assert.Equal(t, code, errorAct.Code)
	assert.Equal(t, message, errorAct.Message)
	assert.Equal(t, recoverable, errorAct.Recoverable)
	assert.NotNil(t, errorAct.Severity)
	assert.Equal(t, severity, *errorAct.Severity)
	assert.True(t, IsValidActID(errorAct.ID))

	// Test validation
	err := errorAct.Validate()
	assert.NoError(t, err)
}

// ============================================================================
// Act Validation Tests
// ============================================================================

func TestAskValidation(t *testing.T) {
	tests := []struct {
		name      string
		ask       Ask
		wantError bool
	}{
		{
			name: "Valid ask",
			ask: Ask{
				Act:    CreateBaseAct("agent_123", ActTypeAsk),
				Field:  "email",
				Prompt: "What is your email?",
			},
			wantError: false,
		},
		{
			name: "Missing field",
			ask: Ask{
				Act:    CreateBaseAct("agent_123", ActTypeAsk),
				Prompt: "What is your email?",
			},
			wantError: true,
		},
		{
			name: "Missing prompt",
			ask: Ask{
				Act:   CreateBaseAct("agent_123", ActTypeAsk),
				Field: "email",
			},
			wantError: true,
		},
		{
			name: "Negative retry count",
			ask: Ask{
				Act:        CreateBaseAct("agent_123", ActTypeAsk),
				Field:      "email",
				Prompt:     "What is your email?",
				RetryCount: func() *int { i := -1; return &i }(),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ask.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFactValidation(t *testing.T) {
	tests := []struct {
		name      string
		fact      Fact
		wantError bool
	}{
		{
			name: "Valid fact",
			fact: Fact{
				Act:    CreateBaseAct("customer_456", ActTypeFact),
				Entity: "order_789",
				Field:  "email",
				Value:  "user@example.com",
			},
			wantError: false,
		},
		{
			name: "Missing entity",
			fact: Fact{
				Act:   CreateBaseAct("customer_456", ActTypeFact),
				Field: "email",
				Value: "user@example.com",
			},
			wantError: true,
		},
		{
			name: "Missing field",
			fact: Fact{
				Act:    CreateBaseAct("customer_456", ActTypeFact),
				Entity: "order_789",
				Value:  "user@example.com",
			},
			wantError: true,
		},
		{
			name: "Missing value",
			fact: Fact{
				Act:    CreateBaseAct("customer_456", ActTypeFact),
				Entity: "order_789",
				Field:  "email",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fact.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ============================================================================
// Type Guard Tests
// ============================================================================

func TestTypeGuards(t *testing.T) {
	ask := NewAsk("agent_123", "email", "What's your email?")
	fact := NewFact("customer_456", "order_789", "email", "user@example.com")
	confirm := NewConfirm("agent_123", "order_456", "Order confirmation")
	commit := NewCommit("system_001", "order_456", CommitActionCreate)
	errorAct := NewError("system_001", "ERROR", "Something went wrong", true)

	// Test IsAct
	assert.True(t, IsAct(ask))
	assert.True(t, IsAct(fact))
	assert.True(t, IsAct(confirm))
	assert.True(t, IsAct(commit))
	assert.True(t, IsAct(errorAct))
	assert.False(t, IsAct("not an act"))
	assert.False(t, IsAct(nil))

	// Test specific type guards
	assert.True(t, IsAsk(ask))
	assert.False(t, IsAsk(fact))

	assert.True(t, IsFact(fact))
	assert.False(t, IsFact(ask))

	assert.True(t, IsConfirm(confirm))
	assert.False(t, IsConfirm(fact))

	assert.True(t, IsCommit(commit))
	assert.False(t, IsCommit(ask))

	assert.True(t, IsError(errorAct))
	assert.False(t, IsError(fact))
}

func TestEntityTypeGuards(t *testing.T) {
	entity := Entity{ID: "entity_123", Type: "order"}
	participant := Participant{ID: "participant_456", Type: ParticipantTypeHuman}

	// Test IsEntity
	assert.True(t, IsEntity(entity))
	assert.True(t, IsEntity("string_entity_ref"))
	assert.False(t, IsEntity(123))
	assert.False(t, IsEntity(nil))

	// Test IsParticipant  
	assert.True(t, IsParticipant(participant))
	assert.False(t, IsParticipant(entity))
	assert.False(t, IsParticipant(Participant{ID: "test"})) // missing type
	assert.False(t, IsParticipant(Participant{Type: ParticipantTypeHuman})) // missing ID
}

// ============================================================================
// JSON Marshaling/Unmarshaling Tests
// ============================================================================

func TestActJSONSerialization(t *testing.T) {
	ask := NewAsk("agent_123", "email", "What's your email?",
		WithRequired(true),
		WithConfidence(0.95))

	// Marshal to JSON
	jsonData, err := json.Marshal(ask)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaledAsk Ask
	err = json.Unmarshal(jsonData, &unmarshaledAsk)
	require.NoError(t, err)

	// Compare
	assert.Equal(t, ask.ID, unmarshaledAsk.ID)
	assert.Equal(t, ask.Speaker, unmarshaledAsk.Speaker)
	assert.Equal(t, ask.Type, unmarshaledAsk.Type)
	assert.Equal(t, ask.Field, unmarshaledAsk.Field)
	assert.Equal(t, ask.Prompt, unmarshaledAsk.Prompt)
	assert.Equal(t, *ask.Required, *unmarshaledAsk.Required)
	assert.Equal(t, *ask.Confidence, *unmarshaledAsk.Confidence)
}

func TestMarshalAct(t *testing.T) {
	ask := NewAsk("agent_123", "email", "What's your email?")

	jsonData, err := MarshalAct(ask)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Should contain expected fields
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	require.NoError(t, err)

	assert.Equal(t, "ask", jsonMap["type"])
	assert.Equal(t, "agent_123", jsonMap["speaker"])
	assert.Equal(t, "email", jsonMap["field"])
	assert.Equal(t, "What's your email?", jsonMap["prompt"])
}

func TestUnmarshalAct(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectedType ActType
		shouldError bool
	}{
		{
			name: "Valid Ask",
			jsonData: `{
				"id": "act_123",
				"timestamp": "2025-01-15T14:30:00Z",
				"speaker": "agent_123",
				"type": "ask",
				"field": "email",
				"prompt": "What's your email?"
			}`,
			expectedType: ActTypeAsk,
			shouldError:  false,
		},
		{
			name: "Valid Fact",
			jsonData: `{
				"id": "act_456",
				"timestamp": "2025-01-15T14:31:00Z",
				"speaker": "customer_456",
				"type": "fact",
				"entity": "order_789",
				"field": "email",
				"value": "user@example.com"
			}`,
			expectedType: ActTypeFact,
			shouldError:  false,
		},
		{
			name: "Invalid JSON",
			jsonData: `{
				"id": "act_123",
				"timestamp": "invalid-timestamp"
			`,
			shouldError: true,
		},
		{
			name: "Unknown act type",
			jsonData: `{
				"id": "act_123",
				"timestamp": "2025-01-15T14:30:00Z",
				"speaker": "agent_123",
				"type": "unknown"
			}`,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act, err := UnmarshalAct([]byte(tt.jsonData))
			
			if tt.shouldError {
				assert.Error(t, err)
				assert.Nil(t, act)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, act)
				assert.Equal(t, tt.expectedType, act.GetType())
			}
		})
	}
}

// ============================================================================
// Entity and Participant Tests
// ============================================================================

func TestNewEntity(t *testing.T) {
	id := "order_123"
	entityType := "order"
	externalID := "ext_456"
	system := "order_management"

	entity := NewEntity(id, entityType,
		WithExternalID(externalID),
		WithEntitySystem(system))

	assert.Equal(t, id, entity.ID)
	assert.Equal(t, entityType, entity.Type)
	assert.NotNil(t, entity.ExternalID)
	assert.Equal(t, externalID, *entity.ExternalID)
	assert.NotNil(t, entity.System)
	assert.Equal(t, system, *entity.System)
}

func TestNewParticipant(t *testing.T) {
	id := "participant_123"
	participantType := ParticipantTypeHuman
	role := "customer"
	name := "John Doe"
	email := "john@example.com"

	participant := NewParticipant(id, participantType,
		WithRole(role),
		WithName(name),
		WithEmail(email))

	assert.Equal(t, id, participant.ID)
	assert.Equal(t, participantType, participant.Type)
	assert.NotNil(t, participant.Role)
	assert.Equal(t, role, *participant.Role)
	assert.NotNil(t, participant.Name)
	assert.Equal(t, name, *participant.Name)
	assert.NotNil(t, participant.Email)
	assert.Equal(t, email, *participant.Email)
}

func TestGetEntityID(t *testing.T) {
	tests := []struct {
		name        string
		entityRef   EntityRef
		expectedID  string
		shouldError bool
	}{
		{
			name:        "String entity ref",
			entityRef:   "order_123",
			expectedID:  "order_123",
			shouldError: false,
		},
		{
			name: "Entity struct ref",
			entityRef: Entity{
				ID:   "order_456",
				Type: "order",
			},
			expectedID:  "order_456",
			shouldError: false,
		},
		{
			name: "Entity pointer ref",
			entityRef: &Entity{
				ID:   "order_789",
				Type: "order",
			},
			expectedID:  "order_789",
			shouldError: false,
		},
		{
			name:        "Nil entity pointer",
			entityRef:   (*Entity)(nil),
			expectedID:  "",
			shouldError: true,
		},
		{
			name:        "Invalid type",
			entityRef:   123,
			expectedID:  "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GetEntityID(tt.entityRef)
			
			if tt.shouldError {
				assert.Error(t, err)
				assert.Empty(t, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

// ============================================================================
// Conversation Tests
// ============================================================================

func TestNewConversation(t *testing.T) {
	participants := []Participant{
		NewParticipant("agent_123", ParticipantTypeAI, WithRole("agent")),
		NewParticipant("customer_456", ParticipantTypeHuman, WithRole("customer")),
	}
	channel := "voice"

	conv := NewConversation(participants,
		WithConversationChannel(channel))

	assert.True(t, IsValidConversationID(conv.ID))
	assert.Equal(t, len(participants), len(conv.Participants))
	assert.NotNil(t, conv.Acts)
	assert.Empty(t, conv.Acts)
	assert.NotNil(t, conv.StartedAt)
	assert.NotNil(t, conv.Status)
	assert.Equal(t, ConversationStatusActive, *conv.Status)
	assert.NotNil(t, conv.Channel)
	assert.Equal(t, channel, *conv.Channel)
}

func TestConversationAddAct(t *testing.T) {
	participants := []Participant{
		NewParticipant("agent_123", ParticipantTypeAI),
	}
	conv := NewConversation(participants)

	ask := NewAsk("agent_123", "email", "What's your email?")
	
	err := conv.AddAct(ask)
	assert.NoError(t, err)
	assert.Len(t, conv.Acts, 1)
	
	// Metadata should be updated
	assert.NotNil(t, conv.Metadata)
	assert.NotNil(t, conv.Metadata.ActCount)
	assert.Equal(t, 1, *conv.Metadata.ActCount)
}

func TestConversationGetMethods(t *testing.T) {
	participants := []Participant{
		NewParticipant("agent_123", ParticipantTypeAI),
		NewParticipant("customer_456", ParticipantTypeHuman),
	}
	conv := NewConversation(participants)

	// Add various acts
	ask := NewAsk("agent_123", "email", "What's your email?")
	fact := NewFact("customer_456", "order_789", "email", "user@example.com")
	errorAct := NewError("system_001", "ERROR", "Something went wrong", true)

	conv.AddAct(ask)
	conv.AddAct(fact)
	conv.AddAct(errorAct)

	// Test GetActsByType
	asks := conv.GetActsByType(ActTypeAsk)
	assert.Len(t, asks, 1)

	facts := conv.GetActsByType(ActTypeFact)
	assert.Len(t, facts, 1)

	errors := conv.GetActsByType(ActTypeError)
	assert.Len(t, errors, 1)

	// Test GetActsBySpeaker
	agentActs := conv.GetActsBySpeaker("agent_123")
	assert.Len(t, agentActs, 1)

	customerActs := conv.GetActsBySpeaker("customer_456")
	assert.Len(t, customerActs, 1)

	// Test GetParticipantByID
	agent := conv.GetParticipantByID("agent_123")
	assert.NotNil(t, agent)
	assert.Equal(t, "agent_123", agent.ID)

	nonExistent := conv.GetParticipantByID("nonexistent")
	assert.Nil(t, nonExistent)
}

func TestConversationEndConversation(t *testing.T) {
	participants := []Participant{
		NewParticipant("agent_123", ParticipantTypeAI),
	}
	conv := NewConversation(participants)

	// Add some time delay to test duration calculation
	time.Sleep(10 * time.Millisecond)

	conv.EndConversation(ConversationStatusCompleted)

	assert.NotNil(t, conv.EndedAt)
	assert.NotNil(t, conv.Status)
	assert.Equal(t, ConversationStatusCompleted, *conv.Status)
	assert.NotNil(t, conv.Metadata)
	assert.NotNil(t, conv.Metadata.TotalDurationMs)
	assert.Greater(t, *conv.Metadata.TotalDurationMs, int64(0))
}

// ============================================================================
// Constraint Tests
// ============================================================================

func TestConstraintBuilders(t *testing.T) {
	// Test RequiredConstraint
	required := RequiredConstraint()
	assert.Equal(t, ConstraintTypeRequired, required.Type)
	assert.NotNil(t, required.Message)

	// Test MinLengthConstraint
	minLength := MinLengthConstraint(5)
	assert.Equal(t, ConstraintTypeMinLength, minLength.Type)
	assert.Equal(t, 5, minLength.Value)
	assert.NotNil(t, minLength.Message)

	// Test MaxLengthConstraint
	maxLength := MaxLengthConstraint(100)
	assert.Equal(t, ConstraintTypeMaxLength, maxLength.Type)
	assert.Equal(t, 100, maxLength.Value)

	// Test EmailFormatConstraint
	email := EmailFormatConstraint()
	assert.Equal(t, ConstraintTypeFormat, email.Type)
	assert.Equal(t, FormatTypeEmail, email.Value)

	// Test EnumConstraint
	enum := EnumConstraint([]string{"small", "medium", "large"})
	assert.Equal(t, ConstraintTypeEnum, enum.Type)
	assert.Equal(t, []string{"small", "medium", "large"}, enum.Value)

	// Test RangeConstraint
	min := 0.0
	max := 100.0
	rangeConstraint := RangeConstraint(&min, &max, true)
	assert.Equal(t, ConstraintTypeRange, rangeConstraint.Type)
	
	rangeValue, ok := rangeConstraint.Value.(RangeConstraint)
	assert.True(t, ok)
	assert.NotNil(t, rangeValue.Min)
	assert.Equal(t, min, *rangeValue.Min)
	assert.NotNil(t, rangeValue.Max)
	assert.Equal(t, max, *rangeValue.Max)
	assert.NotNil(t, rangeValue.Inclusive)
	assert.True(t, *rangeValue.Inclusive)
}

// ============================================================================
// Schema Validation Tests
// ============================================================================

func TestGetSchema(t *testing.T) {
	tests := []struct {
		name        string
		schemaName  string
		shouldError bool
	}{
		{"Valid schema - act", "act", false},
		{"Valid schema - ask", "ask", false},
		{"Valid schema - fact", "fact", false},
		{"Valid schema - confirm", "confirm", false},
		{"Valid schema - commit", "commit", false},
		{"Valid schema - error", "error", false},
		{"Valid schema - entity", "entity", false},
		{"Valid schema - participant", "participant", false},
		{"Valid schema - constraint", "constraint", false},
		{"Valid schema - conversation", "conversation", false},
		{"Invalid schema", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := GetSchema(tt.schemaName)
			
			if tt.shouldError {
				assert.Error(t, err)
				assert.Nil(t, schema)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schema)
				assert.Contains(t, schema, "$schema")
				assert.Contains(t, schema, "$id")
				assert.Contains(t, schema, "title")
			}
		})
	}
}

func TestValidateJSON(t *testing.T) {
	validAskJSON := `{
		"id": "act_123",
		"timestamp": "2025-01-15T14:30:00Z",
		"speaker": "agent_123",
		"type": "ask",
		"field": "email",
		"prompt": "What's your email?"
	}`

	invalidAskJSON := `{
		"id": "act_123",
		"timestamp": "2025-01-15T14:30:00Z",
		"speaker": "agent_123",
		"type": "ask"
	}`

	err := ValidateJSON([]byte(validAskJSON), "ask")
	assert.NoError(t, err)

	err = ValidateJSON([]byte(invalidAskJSON), "ask")
	assert.Error(t, err)

	err = ValidateJSON([]byte("invalid json"), "ask")
	assert.Error(t, err)

	err = ValidateJSON([]byte(validAskJSON), "invalid_schema")
	assert.Error(t, err)
}

func TestListSchemas(t *testing.T) {
	schemas := ListSchemas()
	expectedSchemas := []string{
		"act", "ask", "fact", "confirm", "commit", "error",
		"entity", "participant", "constraint", "conversation",
	}

	assert.Equal(t, len(expectedSchemas), len(schemas))
	
	for _, expected := range expectedSchemas {
		assert.Contains(t, schemas, expected)
	}
}

func TestSchemaVersion(t *testing.T) {
	version := SchemaVersion("act")
	assert.Equal(t, "v1", version)
	
	version = SchemaVersion("nonexistent")
	assert.Equal(t, "v1", version)
}

// ============================================================================
// ActUnion Tests
// ============================================================================

func TestActUnion(t *testing.T) {
	ask := NewAsk("agent_123", "email", "What's your email?")
	fact := NewFact("customer_456", "order_789", "email", "user@example.com")

	// Test NewActUnion
	askUnion := NewActUnion(ask)
	assert.Equal(t, ActTypeAsk, askUnion.Type)
	assert.NotNil(t, askUnion.Ask)
	assert.Nil(t, askUnion.Fact)
	assert.Nil(t, askUnion.Confirm)
	assert.Nil(t, askUnion.Commit)
	assert.Nil(t, askUnion.Error)

	factUnion := NewActUnion(fact)
	assert.Equal(t, ActTypeFact, factUnion.Type)
	assert.Nil(t, factUnion.Ask)
	assert.NotNil(t, factUnion.Fact)
	assert.Nil(t, factUnion.Confirm)
	assert.Nil(t, factUnion.Commit)
	assert.Nil(t, factUnion.Error)

	// Test GetAct
	retrievedAsk, err := askUnion.GetAct()
	assert.NoError(t, err)
	assert.Equal(t, ask, retrievedAsk)

	retrievedFact, err := factUnion.GetAct()
	assert.NoError(t, err)
	assert.Equal(t, fact, retrievedFact)

	// Test invalid union
	invalidUnion := ActUnion{Type: ActTypeAsk} // No corresponding ask field
	_, err = invalidUnion.GetAct()
	assert.Error(t, err)
}

// ============================================================================
// Metadata JSON Marshaling Tests
// ============================================================================

func TestActMetadataJSONMarshaling(t *testing.T) {
	channel := "voice"
	language := "en-US"
	originalText := "What's your email?"
	processingTime := 150.0

	metadata := ActMetadata{
		Channel:          &channel,
		Language:         &language,
		OriginalText:     &originalText,
		ProcessingTimeMs: &processingTime,
		AdditionalProperties: map[string]interface{}{
			"custom_field": "custom_value",
			"another_field": 123,
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(metadata)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaledMetadata ActMetadata
	err = json.Unmarshal(jsonData, &unmarshaledMetadata)
	require.NoError(t, err)

	// Compare
	assert.Equal(t, *metadata.Channel, *unmarshaledMetadata.Channel)
	assert.Equal(t, *metadata.Language, *unmarshaledMetadata.Language)
	assert.Equal(t, *metadata.OriginalText, *unmarshaledMetadata.OriginalText)
	assert.Equal(t, *metadata.ProcessingTimeMs, *unmarshaledMetadata.ProcessingTimeMs)
	assert.Equal(t, metadata.AdditionalProperties["custom_field"], unmarshaledMetadata.AdditionalProperties["custom_field"])
	assert.Equal(t, metadata.AdditionalProperties["another_field"], unmarshaledMetadata.AdditionalProperties["another_field"])
}

func TestParticipantPreferencesJSONMarshaling(t *testing.T) {
	language := "en-US"
	timezone := "America/New_York"
	channels := []string{"voice", "text", "email"}

	prefs := ParticipantPreferences{
		Language:              &language,
		Timezone:              &timezone,
		CommunicationChannels: channels,
		AdditionalProperties: map[string]interface{}{
			"notification_style": "immediate",
			"max_response_time":  300,
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(prefs)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaledPrefs ParticipantPreferences
	err = json.Unmarshal(jsonData, &unmarshaledPrefs)
	require.NoError(t, err)

	// Compare
	assert.Equal(t, *prefs.Language, *unmarshaledPrefs.Language)
	assert.Equal(t, *prefs.Timezone, *unmarshaledPrefs.Timezone)
	assert.Equal(t, prefs.CommunicationChannels, unmarshaledPrefs.CommunicationChannels)
	assert.Equal(t, prefs.AdditionalProperties["notification_style"], unmarshaledPrefs.AdditionalProperties["notification_style"])
	assert.Equal(t, prefs.AdditionalProperties["max_response_time"], unmarshaledPrefs.AdditionalProperties["max_response_time"])
}

func TestConversationJSONMarshaling(t *testing.T) {
	participants := []Participant{
		NewParticipant("agent_123", ParticipantTypeAI, WithRole("agent")),
		NewParticipant("customer_456", ParticipantTypeHuman, WithRole("customer")),
	}

	conv := NewConversation(participants)
	
	// Add some acts
	ask := NewAsk("agent_123", "email", "What's your email?")
	fact := NewFact("customer_456", "order_789", "email", "user@example.com")
	
	conv.AddAct(ask)
	conv.AddAct(fact)

	// Marshal to JSON
	jsonData, err := json.Marshal(conv)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaledConv Conversation
	err = json.Unmarshal(jsonData, &unmarshaledConv)
	require.NoError(t, err)

	// Compare basic fields
	assert.Equal(t, conv.ID, unmarshaledConv.ID)
	assert.Equal(t, len(conv.Participants), len(unmarshaledConv.Participants))
	assert.Equal(t, len(conv.Acts), len(unmarshaledConv.Acts))
	assert.Equal(t, *conv.Status, *unmarshaledConv.Status)

	// Compare acts
	for i, originalAct := range conv.Acts {
		unmarshaledAct := unmarshaledConv.Acts[i]
		assert.Equal(t, originalAct.GetType(), unmarshaledAct.GetType())
		assert.Equal(t, originalAct.GetAct().ID, unmarshaledAct.GetAct().ID)
		assert.Equal(t, originalAct.GetAct().Speaker, unmarshaledAct.GetAct().Speaker)
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkGenerateActID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateActID()
	}
}

func BenchmarkGenerateConversationID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateConversationID()
	}
}

func BenchmarkIsValidActID(b *testing.B) {
	id := GenerateActID()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		IsValidActID(id)
	}
}

func BenchmarkNewAsk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewAsk("agent_123", "email", "What's your email?")
	}
}

func BenchmarkNewFact(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewFact("customer_456", "order_789", "email", "user@example.com")
	}
}

func BenchmarkActValidation(b *testing.B) {
	ask := NewAsk("agent_123", "email", "What's your email?")
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		ask.Validate()
	}
}

func BenchmarkActJSONMarshal(b *testing.B) {
	ask := NewAsk("agent_123", "email", "What's your email?", WithConfidence(0.95))
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		json.Marshal(ask)
	}
}

func BenchmarkActJSONUnmarshal(b *testing.B) {
	ask := NewAsk("agent_123", "email", "What's your email?")
	jsonData, _ := json.Marshal(ask)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		UnmarshalAct(jsonData)
	}
}

func BenchmarkConversationAddAct(b *testing.B) {
	participants := []Participant{
		NewParticipant("agent_123", ParticipantTypeAI),
	}
	conv := NewConversation(participants)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		ask := NewAsk("agent_123", "email", "What's your email?")
		conv.AddAct(ask)
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestFullConversationWorkflow(t *testing.T) {
	// Create participants
	agent := NewParticipant("agent_123", ParticipantTypeAI,
		WithRole("customer_service"),
		WithName("Support Agent"),
		WithCapabilities([]string{"order_lookup", "payment_processing"}))

	customer := NewParticipant("customer_456", ParticipantTypeHuman,
		WithRole("customer"),
		WithName("John Doe"),
		WithEmail("john@example.com"))

	participants := []Participant{agent, customer}

	// Create conversation
	conv := NewConversation(participants,
		WithConversationChannel("voice"),
		WithConversationSchema("customer_service_v1"))

	assert.Equal(t, ConversationStatusActive, *conv.Status)
	assert.Equal(t, 0, len(conv.Acts))

	// Agent asks for customer information
	askEmail := NewAsk("agent_123", "email", "What's your email address?",
		WithRequired(true),
		WithExpectedType(ExpectedTypeEmail),
		WithConstraints([]Constraint{EmailFormatConstraint()}))

	err := conv.AddAct(askEmail)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(conv.Acts))

	// Customer provides email
	factEmail := NewFact("customer_456", "customer_456", "email", "john@example.com",
		WithOperation(FieldOperationSet),
		WithValidationStatus(ValidationStatusValid))

	err = conv.AddAct(factEmail)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(conv.Acts))

	// Agent confirms the information
	confirm := NewConfirm("agent_123", "customer_456", "Your email is john@example.com, is that correct?",
		WithAwaiting(true))

	err = conv.AddAct(confirm)
	assert.NoError(t, err)

	// Customer confirms
	confirmResponse := NewFact("customer_456", confirm.ID, "confirmed", true,
		WithOperation(FieldOperationSet))

	err = conv.AddAct(confirmResponse)
	assert.NoError(t, err)

	// System commits the information
	commit := NewCommit("system_001", "customer_456", CommitActionUpdate,
		WithSystem("customer_management"),
		WithCommitStatus(CommitStatusSuccess))

	err = conv.AddAct(commit)
	assert.NoError(t, err)

	// End the conversation
	conv.EndConversation(ConversationStatusCompleted)

	// Verify final state
	assert.Equal(t, ConversationStatusCompleted, *conv.Status)
	assert.Equal(t, 5, len(conv.Acts))
	assert.NotNil(t, conv.EndedAt)
	assert.NotNil(t, conv.Metadata)
	assert.Equal(t, 5, *conv.Metadata.ActCount)
	assert.Equal(t, 1, *conv.Metadata.CommitCount)
	assert.Equal(t, 0, *conv.Metadata.ErrorCount)

	// Test querying methods
	asks := conv.GetActsByType(ActTypeAsk)
	assert.Equal(t, 1, len(asks))

	facts := conv.GetActsByType(ActTypeFact)
	assert.Equal(t, 2, len(facts))

	agentActs := conv.GetActsBySpeaker("agent_123")
	assert.Equal(t, 2, len(agentActs))

	customerActs := conv.GetActsBySpeaker("customer_456")
	assert.Equal(t, 2, len(customerActs))

	// Test JSON serialization of the complete conversation
	jsonData, err := json.Marshal(conv)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var unmarshaledConv Conversation
	err = json.Unmarshal(jsonData, &unmarshaledConv)
	assert.NoError(t, err)
	assert.Equal(t, conv.ID, unmarshaledConv.ID)
	assert.Equal(t, len(conv.Acts), len(unmarshaledConv.Acts))
}

func TestErrorHandlingWorkflow(t *testing.T) {
	participants := []Participant{
		NewParticipant("agent_123", ParticipantTypeAI),
		NewParticipant("customer_456", ParticipantTypeHuman),
	}
	conv := NewConversation(participants)

	// Agent asks for invalid information (missing required fields)
	invalidAsk := Ask{
		Act: CreateBaseAct("agent_123", ActTypeAsk),
		// Missing required Field and Prompt
	}

	err := conv.AddAct(invalidAsk)
	assert.Error(t, err)
	assert.Equal(t, 0, len(conv.Acts))

	// Add a valid ask
	validAsk := NewAsk("agent_123", "email", "What's your email?")
	err = conv.AddAct(validAsk)
	assert.NoError(t, err)

	// System encounters an error
	errorAct := NewError("system_001", "VALIDATION_ERROR", "Invalid email format", true,
		WithSeverity(ErrorSeverityError),
		WithCategory(ErrorCategoryValidation),
		WithSuggestedAction(SuggestedActionClarify),
		WithRelatedActID(validAsk.ID))

	err = conv.AddAct(errorAct)
	assert.NoError(t, err)

	// Verify error tracking
	assert.Equal(t, 2, len(conv.Acts))
	assert.NotNil(t, conv.Metadata)
	assert.Equal(t, 1, *conv.Metadata.ErrorCount)

	errors := conv.GetActsByType(ActTypeError)
	assert.Equal(t, 1, len(errors))

	errorAct2 := errors[0].(Error)
	assert.Equal(t, "VALIDATION_ERROR", errorAct2.Code)
	assert.Equal(t, true, errorAct2.Recoverable)
	assert.Equal(t, ErrorSeverityError, *errorAct2.Severity)
}

// ============================================================================
// Edge Cases and Error Conditions
// ============================================================================

func TestInvalidConfidenceValues(t *testing.T) {
	// Test invalid confidence values are ignored
	act1 := CreateBaseAct("agent_123", ActTypeAsk, WithConfidence(-0.1))
	assert.Nil(t, act1.Confidence)

	act2 := CreateBaseAct("agent_123", ActTypeAsk, WithConfidence(1.1))
	assert.Nil(t, act2.Confidence)

	act3 := CreateBaseAct("agent_123", ActTypeAsk, WithConfidence(0.5))
	assert.NotNil(t, act3.Confidence)
	assert.Equal(t, 0.5, *act3.Confidence)
}

func TestNegativeRetryValues(t *testing.T) {
	// Test that negative retry values are ignored
	ask := NewAsk("agent_123", "email", "What's your email?", WithMaxRetries(-1))
	assert.Nil(t, ask.MaxRetries)

	ask2 := NewAsk("agent_123", "email", "What's your email?", WithMaxRetries(3))
	assert.NotNil(t, ask2.MaxRetries)
	assert.Equal(t, 3, *ask2.MaxRetries)
}

func TestInvalidActIDInError(t *testing.T) {
	// Test that invalid act IDs are ignored in error related_act_id
	errorAct := NewError("system_001", "ERROR", "Something went wrong", true,
		WithRelatedActID("invalid_id_format"))
	assert.Nil(t, errorAct.RelatedActID)

	errorAct2 := NewError("system_001", "ERROR", "Something went wrong", true,
		WithRelatedActID("act_valid_123"))
	assert.NotNil(t, errorAct2.RelatedActID)
	assert.Equal(t, "act_valid_123", *errorAct2.RelatedActID)
}

func TestEmptyConversationValidation(t *testing.T) {
	// Test that conversations require participants
	emptyConv := Conversation{
		ID:           GenerateConversationID(),
		Participants: []Participant{},
		Acts:         []ConversationAct{},
	}

	assert.False(t, IsConversation(emptyConv))

	// Test with participants
	validConv := Conversation{
		ID: GenerateConversationID(),
		Participants: []Participant{
			NewParticipant("agent_123", ParticipantTypeAI),
		},
		Acts: []ConversationAct{},
	}

	assert.True(t, IsConversation(validConv))
}
