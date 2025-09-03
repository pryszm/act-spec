/**
 * @astra/model-ts
 * 
 * TypeScript types and JSON schemas for ASTRA (Act State Representation Architecture)
 * 
 * This package provides:
 * - TypeScript interfaces for all ASTRA types
 * - JSON schemas for runtime validation
 * - Type guards and utility functions
 */

// Export all types
export * from './types';

// Export schemas
export { schemas } from './schemas';

// Re-export commonly used types for convenience
export type {
  Act,
  Ask,
  Fact,
  Confirm,
  Commit,
  Error,
  ConversationAct,
  Conversation,
  Participant,
  Entity,
  EntityRef,
  Constraint
} from './types';

/**
 * Package version information
 */
export const VERSION = '1.0.0';

/**
 * ASTRA schema version this package implements
 */
export const SCHEMA_VERSION = 'v1';

/**
 * Type guard functions for runtime type checking
 */

/**
 * Type guard to check if an object is a valid Act
 */
export function isAct(obj: any): obj is Act {
  return obj && 
         typeof obj.id === 'string' && 
         typeof obj.timestamp === 'string' && 
         typeof obj.speaker === 'string' && 
         ['ask', 'fact', 'confirm', 'commit', 'error'].includes(obj.type);
}

/**
 * Type guard to check if an object is an Ask act
 */
export function isAsk(obj: any): obj is Ask {
  return isAct(obj) && obj.type === 'ask' && typeof obj.field === 'string' && typeof obj.prompt === 'string';
}

/**
 * Type guard to check if an object is a Fact act
 */
export function isFact(obj: any): obj is Fact {
  return isAct(obj) && obj.type === 'fact' && typeof obj.field === 'string' && obj.value !== undefined;
}

/**
 * Type guard to check if an object is a Confirm act
 */
export function isConfirm(obj: any): obj is Confirm {
  return isAct(obj) && obj.type === 'confirm' && typeof obj.summary === 'string';
}

/**
 * Type guard to check if an object is a Commit act
 */
export function isCommit(obj: any): obj is Commit {
  return isAct(obj) && obj.type === 'commit' && typeof obj.action === 'string';
}

/**
 * Type guard to check if an object is an Error act
 */
export function isError(obj: any): obj is Error {
  return isAct(obj) && 
         obj.type === 'error' && 
         typeof obj.code === 'string' && 
         typeof obj.message === 'string' && 
         typeof obj.recoverable === 'boolean';
}

/**
 * Type guard to check if an object is a valid Participant
 */
export function isParticipant(obj: any): obj is Participant {
  return obj && 
         typeof obj.id === 'string' && 
         ['human', 'ai', 'system', 'bot'].includes(obj.type);
}

/**
 * Type guard to check if an object is a valid Entity
 */
export function isEntity(obj: any): obj is Entity {
  return obj && typeof obj.id === 'string' && typeof obj.type === 'string';
}

/**
 * Type guard to check if an object is a valid Conversation
 */
export function isConversation(obj: any): obj is Conversation {
  return obj && 
         typeof obj.id === 'string' && 
         Array.isArray(obj.participants) && 
         Array.isArray(obj.acts);
}

/**
 * Utility function to generate ASTRA-compliant act IDs
 */
export function generateActId(): string {
  const timestamp = Date.now().toString(36);
  const random = Math.random().toString(36).substring(2, 8);
  return `act_${timestamp}_${random}`;
}

/**
 * Utility function to generate ASTRA-compliant conversation IDs
 */
export function generateConversationId(): string {
  const timestamp = Date.now().toString(36);
  const random = Math.random().toString(36).substring(2, 8);
  return `conv_${timestamp}_${random}`;
}

/**
 * Utility function to create a basic Act structure with required fields
 */
export function createBaseAct(
  speaker: string,
  type: ActType,
  additionalFields: Partial<Act> = {}
): Omit<Act, keyof typeof additionalFields> & typeof additionalFields {
  return {
    id: generateActId(),
    timestamp: new Date().toISOString(),
    speaker,
    type,
    ...additionalFields
  };
}

// Import types for the utility function
import type { Act, Ask, Fact, Confirm, Commit, Error, ActType, Participant, Entity, Conversation } from './types';
