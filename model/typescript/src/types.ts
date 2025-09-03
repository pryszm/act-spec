/**
 * ASTRA (Act State Representation Architecture) TypeScript Types
 * 
 * Core types for representing conversational state as typed, auditable sequences of actions.
 */

// ============================================================================
// Base Types
// ============================================================================

/**
 * Source that generated this act
 */
export type Source = 'human' | 'speech_recognition' | 'text_analysis' | 'system' | 'ai';

/**
 * Type of conversational act being performed
 */
export type ActType = 'ask' | 'fact' | 'confirm' | 'commit' | 'error';

/**
 * Additional metadata for acts
 */
export interface ActMetadata {
  /** Communication channel (voice, text, email, etc.) */
  channel?: string;
  /** Language code (ISO 639-1, optional region) */
  language?: string;
  /** Original utterance that generated this act */
  original_text?: string;
  /** Time taken to process this act in milliseconds */
  processing_time_ms?: number;
  /** Additional context-specific metadata */
  [key: string]: any;
}

/**
 * Base type for all conversational actions in ASTRA
 */
export interface Act {
  /** Unique identifier for this act within the conversation */
  id: string;
  /** ISO 8601 timestamp when the act occurred */
  timestamp: string;
  /** Identifier of the conversation participant who performed this act */
  speaker: string;
  /** Type of conversational act being performed */
  type: ActType;
  /** Confidence score for automated act extraction (0.0 to 1.0) */
  confidence?: number;
  /** Source that generated this act */
  source?: Source;
  /** Additional context-specific metadata */
  metadata?: ActMetadata;
}

// ============================================================================
// Entity Types
// ============================================================================

/**
 * Reference to a business entity in ASTRA conversations
 */
export interface Entity {
  /** Unique identifier for this entity within the conversation scope */
  id: string;
  /** Type of business entity (order, customer, appointment, ticket, etc.) */
  type: string;
  /** External system identifier for this entity */
  external_id?: string;
  /** External system that owns this entity */
  system?: string;
  /** Version or revision of this entity */
  version?: string;
  /** URL to the schema definition for this entity type */
  schema_url?: string;
  /** Additional entity-specific metadata */
  metadata?: Record<string, any>;
}

/**
 * Entity reference that can be either a string ID or structured Entity
 */
export type EntityRef = string | Entity;

// ============================================================================
// Constraint Types
// ============================================================================

/**
 * Type of constraint being applied
 */
export type ConstraintType = 
  | 'required' 
  | 'optional' 
  | 'min_length' 
  | 'max_length' 
  | 'pattern' 
  | 'format' 
  | 'range' 
  | 'enum' 
  | 'custom';

/**
 * Format validation types
 */
export type FormatType = 
  | 'email' 
  | 'phone' 
  | 'url' 
  | 'date' 
  | 'time' 
  | 'datetime' 
  | 'uuid' 
  | 'ipv4' 
  | 'ipv6';

/**
 * Range constraint value
 */
export interface RangeConstraint {
  min?: number;
  max?: number;
  inclusive?: boolean;
}

/**
 * Validation constraint for ASTRA fields and values
 */
export interface Constraint {
  /** Type of constraint being applied */
  type: ConstraintType;
  /** Constraint value (varies by constraint type) */
  value?: number | string | FormatType | RangeConstraint | string[];
  /** Human-readable error message when constraint is violated */
  message?: string;
  /** Machine-readable error code for constraint violations */
  code?: string;
}

// ============================================================================
// Participant Types
// ============================================================================

/**
 * Type of participant
 */
export type ParticipantType = 'human' | 'ai' | 'system' | 'bot';

/**
 * Participant preferences
 */
export interface ParticipantPreferences {
  /** Preferred language code (ISO 639-1 with optional region) */
  language?: string;
  /** Preferred timezone (IANA timezone identifier) */
  timezone?: string;
  /** Preferred communication channels in order of preference */
  communication_channels?: string[];
  /** Additional preferences */
  [key: string]: any;
}

/**
 * Conversation participant in ASTRA conversations
 */
export interface Participant {
  /** Unique identifier for this participant */
  id: string;
  /** Type of participant */
  type: ParticipantType;
  /** Business role of the participant (customer, agent, manager, etc.) */
  role?: string;
  /** Display name of the participant */
  name?: string;
  /** Email address of the participant */
  email?: string;
  /** Phone number of the participant */
  phone?: string;
  /** External system identifier for this participant */
  external_id?: string;
  /** External system that manages this participant */
  system?: string;
  /** List of capabilities this participant has */
  capabilities?: string[];
  /** List of permissions granted to this participant */
  permissions?: string[];
  /** Participant preferences */
  preferences?: ParticipantPreferences;
  /** Additional participant metadata */
  metadata?: Record<string, any>;
}

// ============================================================================
// Act Types
// ============================================================================

/**
 * Expected data type of the response
 */
export type ExpectedType = 
  | 'string' 
  | 'number' 
  | 'boolean' 
  | 'object' 
  | 'array' 
  | 'date' 
  | 'email' 
  | 'phone' 
  | 'address';

/**
 * Act that requests missing information required to complete a business process
 */
export interface Ask extends Act {
  type: 'ask';
  /** Field or information being requested */
  field: string;
  /** Question or request presented to obtain the information */
  prompt: string;
  /** Validation constraints for the requested information */
  constraints?: Constraint[];
  /** Whether this information is required to proceed */
  required?: boolean;
  /** Expected data type of the response */
  expected_type?: ExpectedType;
  /** Number of times this question has been asked */
  retry_count?: number;
  /** Maximum number of retry attempts before escalation */
  max_retries?: number;
}

/**
 * Operation being performed on the field
 */
export type FieldOperation = 'set' | 'append' | 'increment' | 'decrement' | 'delete' | 'merge';

/**
 * Validation status of this fact
 */
export type ValidationStatus = 'pending' | 'valid' | 'invalid' | 'partial';

/**
 * Act that declares facts or information provided during conversation
 */
export interface Fact extends Act {
  type: 'fact';
  /** Business entity being modified (order, customer, appointment, etc.) */
  entity: EntityRef;
  /** Specific field or property being set */
  field: string;
  /** Value being assigned to the field */
  value: any;
  /** Operation being performed on the field */
  operation?: FieldOperation;
  /** Previous value of the field (for audit trail) */
  previous_value?: any;
  /** Validation status of this fact */
  validation_status?: ValidationStatus;
  /** List of validation errors if validation_status is invalid */
  validation_errors?: string[];
}

/**
 * How the confirmation was obtained
 */
export type ConfirmationMethod = 'verbal' | 'explicit' | 'implicit' | 'timeout' | 'system';

/**
 * Act that verifies understanding of information before commitment
 */
export interface Confirm extends Act {
  type: 'confirm';
  /** Business entity being confirmed */
  entity: EntityRef;
  /** Human-readable summary of what is being confirmed */
  summary: string;
  /** Whether confirmation is still pending */
  awaiting?: boolean;
  /** Whether the confirmation was accepted (true) or rejected (false) */
  confirmed?: boolean;
  /** How the confirmation was obtained */
  confirmation_method?: ConfirmationMethod;
  /** Specific fields or aspects being confirmed */
  fields_confirmed?: string[];
  /** Reason provided if confirmation was rejected */
  rejection_reason?: string;
  /** Timeout for awaiting confirmation in milliseconds */
  timeout_ms?: number;
}

/**
 * Action being performed in the target system
 */
export type CommitAction = 'create' | 'update' | 'delete' | 'execute' | 'cancel' | 'pause' | 'resume';

/**
 * Status of the commit operation
 */
export type CommitStatus = 'pending' | 'in_progress' | 'success' | 'failed' | 'retrying' | 'cancelled';

/**
 * Error information for failed commits
 */
export interface CommitError {
  /** Error code from the target system */
  code: string;
  /** Human-readable error message */
  message: string;
  /** Additional error context */
  details?: Record<string, any>;
  /** Whether the error can be recovered from */
  recoverable: boolean;
}

/**
 * Act that executes business processes and triggers system integrations
 */
export interface Commit extends Act {
  type: 'commit';
  /** Business entity being committed to external systems */
  entity: EntityRef;
  /** Action being performed in the target system */
  action: CommitAction;
  /** Target system identifier (CRM, order_management, etc.) */
  system?: string;
  /** External system transaction or record identifier */
  transaction_id?: string;
  /** Status of the commit operation */
  status?: CommitStatus;
  /** Error information if commit failed */
  error?: CommitError;
  /** Number of retry attempts made */
  retry_count?: number;
  /** Maximum number of retry attempts */
  max_retries?: number;
  /** Key to ensure idempotent operations */
  idempotency_key?: string;
  /** Information needed to rollback this commit if necessary */
  rollback_info?: Record<string, any>;
}

/**
 * Severity level of the error
 */
export type ErrorSeverity = 'info' | 'warning' | 'error' | 'critical';

/**
 * Category of error for classification
 */
export type ErrorCategory = 
  | 'validation' 
  | 'processing' 
  | 'integration' 
  | 'timeout' 
  | 'permission' 
  | 'system' 
  | 'user_input' 
  | 'business_rule';

/**
 * Suggested recovery action
 */
export type SuggestedAction = 'retry' | 'escalate' | 'ignore' | 'clarify' | 'fallback' | 'terminate';

/**
 * Act that handles failures and exceptions in conversational processing
 */
export interface Error extends Act {
  type: 'error';
  /** Machine-readable error code */
  code: string;
  /** Human-readable error message */
  message: string;
  /** Whether the conversation can continue after this error */
  recoverable: boolean;
  /** Severity level of the error */
  severity?: ErrorSeverity;
  /** Category of error for classification */
  category?: ErrorCategory;
  /** Additional error context and debugging information */
  details?: Record<string, any>;
  /** ID of the act that caused this error */
  related_act_id?: string;
  /** Suggested recovery action */
  suggested_action?: SuggestedAction;
  /** User-friendly message to display to conversation participants */
  user_message?: string;
  /** Technical stack trace for debugging (not shown to users) */
  stack_trace?: string;
}

// ============================================================================
// Conversation Types
// ============================================================================

/**
 * Union type for all possible acts
 */
export type ConversationAct = Ask | Fact | Confirm | Commit | Error;

/**
 * Current status of the conversation
 */
export type ConversationStatus = 'active' | 'paused' | 'completed' | 'failed' | 'cancelled';

/**
 * Conversation context and session information
 */
export interface ConversationContext {
  /** Session identifier */
  session_id?: string;
  /** User agent or client information */
  user_agent?: string;
  /** Client IP address */
  ip_address?: string;
  /** How the conversation was initiated */
  referrer?: string;
  /** Additional context properties */
  [key: string]: any;
}

/**
 * Additional conversation metadata
 */
export interface ConversationMetadata {
  /** Total conversation duration in milliseconds */
  total_duration_ms?: number;
  /** Total number of acts in the conversation */
  act_count?: number;
  /** Number of errors that occurred */
  error_count?: number;
  /** Number of successful commits */
  commit_count?: number;
  /** Average confidence score across all acts */
  avg_confidence?: number;
  /** Additional metadata properties */
  [key: string]: any;
}

/**
 * Complete ASTRA conversation container with acts and metadata
 */
export interface Conversation {
  /** Unique identifier for this conversation */
  id: string;
  /** List of conversation participants */
  participants: Participant[];
  /** Ordered sequence of acts in this conversation */
  acts: ConversationAct[];
  /** When the conversation started */
  started_at?: string;
  /** When the conversation ended */
  ended_at?: string;
  /** Current status of the conversation */
  status?: ConversationStatus;
  /** Primary communication channel for this conversation */
  channel?: string;
  /** Business schema identifier used for this conversation */
  schema?: string;
  /** Conversation context and session information */
  context?: ConversationContext;
  /** Final computed state of all entities after processing all acts */
  final_state?: Record<string, any>;
  /** Additional conversation metadata */
  metadata?: ConversationMetadata;
}
