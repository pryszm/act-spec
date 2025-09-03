/**
 * Tests for ASTRA TypeScript types and utilities
 */

import {
  // Types
  Act,
  Ask,
  Fact,
  Confirm,
  Commit,
  Error,
  Participant,
  Entity,
  Conversation,
  Constraint,
  
  // Type guards
  isAct,
  isAsk,
  isFact,
  isConfirm,
  isCommit,
  isError,
  isParticipant,
  isEntity,
  isConversation,
  
  // Utilities
  generateActId,
  generateConversationId,
  createBaseAct,
  
  // Constants
  VERSION,
  SCHEMA_VERSION,
  
  // Schemas
  schemas
} from '../src';

describe('ASTRA Types', () => {
  describe('Type Guards', () => {
    const validAct: Act = {
      id: 'act_123',
      timestamp: '2025-01-15T14:30:00Z',
      speaker: 'agent_001',
      type: 'ask'
    };

    const validAsk: Ask = {
      ...validAct,
      type: 'ask',
      field: 'email',
      prompt: 'What is your email address?'
    };

    const validFact: Fact = {
      ...validAct,
      type: 'fact',
      entity: 'customer_123',
      field: 'email',
      value: 'user@example.com'
    };

    const validConfirm: Confirm = {
      ...validAct,
      type: 'confirm',
      entity: 'order_456',
      summary: 'Order for 2 pizzas to be delivered at 6 PM'
    };

    const validCommit: Commit = {
      ...validAct,
      type: 'commit',
      entity: 'order_456',
      action: 'create'
    };

    const validError: Error = {
      ...validAct,
      type: 'error',
      code: 'VALIDATION_ERROR',
      message: 'Invalid email format',
      recoverable: true
    };

    describe('isAct', () => {
      it('should return true for valid acts', () => {
        expect(isAct(validAct)).toBe(true);
        expect(isAct(validAsk)).toBe(true);
        expect(isAct(validFact)).toBe(true);
        expect(isAct(validConfirm)).toBe(true);
        expect(isAct(validCommit)).toBe(true);
        expect(isAct(validError)).toBe(true);
      });

      it('should return false for invalid acts', () => {
        expect(isAct(null)).toBe(false);
        expect(isAct(undefined)).toBe(false);
        expect(isAct({})).toBe(false);
        expect(isAct({ id: 'test' })).toBe(false);
        expect(isAct({ ...validAct, type: 'invalid' })).toBe(false);
      });
    });

    describe('isAsk', () => {
      it('should return true for valid Ask acts', () => {
        expect(isAsk(validAsk)).toBe(true);
      });

      it('should return false for other act types', () => {
        expect(isAsk(validFact)).toBe(false);
        expect(isAsk(validConfirm)).toBe(false);
        expect(isAsk(validCommit)).toBe(false);
        expect(isAsk(validError)).toBe(false);
      });

      it('should return false for invalid Ask acts', () => {
        expect(isAsk({ ...validAsk, field: undefined })).toBe(false);
        expect(isAsk({ ...validAsk, prompt: undefined })).toBe(false);
      });
    });

    describe('isFact', () => {
      it('should return true for valid Fact acts', () => {
        expect(isFact(validFact)).toBe(true);
      });

      it('should return false for other act types', () => {
        expect(isFact(validAsk)).toBe(false);
        expect(isFact(validConfirm)).toBe(false);
        expect(isFact(validCommit)).toBe(false);
        expect(isFact(validError)).toBe(false);
      });

      it('should return false for invalid Fact acts', () => {
        expect(isFact({ ...validFact, field: undefined })).toBe(false);
        expect(isFact({ ...validFact, value: undefined })).toBe(false);
      });
    });

    describe('isConfirm', () => {
      it('should return true for valid Confirm acts', () => {
        expect(isConfirm(validConfirm)).toBe(true);
      });

      it('should return false for other act types', () => {
        expect(isConfirm(validAsk)).toBe(false);
        expect(isConfirm(validFact)).toBe(false);
        expect(isConfirm(validCommit)).toBe(false);
        expect(isConfirm(validError)).toBe(false);
      });

      it('should return false for invalid Confirm acts', () => {
        expect(isConfirm({ ...validConfirm, summary: undefined })).toBe(false);
      });
    });

    describe('isCommit', () => {
      it('should return true for valid Commit acts', () => {
        expect(isCommit(validCommit)).toBe(true);
      });

      it('should return false for other act types', () => {
        expect(isCommit(validAsk)).toBe(false);
        expect(isCommit(validFact)).toBe(false);
        expect(isCommit(validConfirm)).toBe(false);
        expect(isCommit(validError)).toBe(false);
      });

      it('should return false for invalid Commit acts', () => {
        expect(isCommit({ ...validCommit, action: undefined })).toBe(false);
      });
    });

    describe('isError', () => {
      it('should return true for valid Error acts', () => {
        expect(isError(validError)).toBe(true);
      });

      it('should return false for other act types', () => {
        expect(isError(validAsk)).toBe(false);
        expect(isError(validFact)).toBe(false);
        expect(isError(validConfirm)).toBe(false);
        expect(isError(validCommit)).toBe(false);
      });

      it('should return false for invalid Error acts', () => {
        expect(isError({ ...validError, code: undefined })).toBe(false);
        expect(isError({ ...validError, message: undefined })).toBe(false);
        expect(isError({ ...validError, recoverable: undefined })).toBe(false);
      });
    });
  });

  describe('Participant Type Guard', () => {
    const validParticipant: Participant = {
      id: 'participant_123',
      type: 'human'
    };

    it('should return true for valid participants', () => {
      expect(isParticipant(validParticipant)).toBe(true);
      expect(isParticipant({ ...validParticipant, type: 'ai' })).toBe(true);
      expect(isParticipant({ ...validParticipant, type: 'system' })).toBe(true);
      expect(isParticipant({ ...validParticipant, type: 'bot' })).toBe(true);
    });

    it('should return false for invalid participants', () => {
      expect(isParticipant(null)).toBe(false);
      expect(isParticipant({})).toBe(false);
      expect(isParticipant({ id: 'test' })).toBe(false);
      expect(isParticipant({ type: 'human' })).toBe(false);
      expect(isParticipant({ ...validParticipant, type: 'invalid' })).toBe(false);
    });
  });

  describe('Entity Type Guard', () => {
    const validEntity: Entity = {
      id: 'entity_123',
      type: 'order'
    };

    it('should return true for valid entities', () => {
      expect(isEntity(validEntity)).toBe(true);
      expect(isEntity({ ...validEntity, external_id: 'ext_123' })).toBe(true);
    });

    it('should return false for invalid entities', () => {
      expect(isEntity(null)).toBe(false);
      expect(isEntity({})).toBe(false);
      expect(isEntity({ id: 'test' })).toBe(false);
      expect(isEntity({ type: 'order' })).toBe(false);
    });
  });

  describe('Conversation Type Guard', () => {
    const validConversation: Conversation = {
      id: 'conv_123',
      participants: [{ id: 'participant_1', type: 'human' }],
      acts: []
    };

    it('should return true for valid conversations', () => {
      expect(isConversation(validConversation)).toBe(true);
    });

    it('should return false for invalid conversations', () => {
      expect(isConversation(null)).toBe(false);
      expect(isConversation({})).toBe(false);
      expect(isConversation({ id: 'test' })).toBe(false);
      expect(isConversation({ ...validConversation, participants: null })).toBe(false);
      expect(isConversation({ ...validConversation, acts: null })).toBe(false);
    });
  });

  describe('Utility Functions', () => {
    describe('generateActId', () => {
      it('should generate valid act IDs', () => {
        const id1 = generateActId();
        const id2 = generateActId();
        
        expect(id1).toMatch(/^act_[a-zA-Z0-9_]+$/);
        expect(id2).toMatch(/^act_[a-zA-Z0-9_]+$/);
        expect(id1).not.toBe(id2); // Should be unique
      });
    });

    describe('generateConversationId', () => {
      it('should generate valid conversation IDs', () => {
        const id1 = generateConversationId();
        const id2 = generateConversationId();
        
        expect(id1).toMatch(/^conv_[a-zA-Z0-9_]+$/);
        expect(id2).toMatch(/^conv_[a-zA-Z0-9_]+$/);
        expect(id1).not.toBe(id2); // Should be unique
      });
    });

    describe('createBaseAct', () => {
      it('should create valid base acts', () => {
        const baseAct = createBaseAct('speaker_123', 'ask');
        
        expect(baseAct.speaker).toBe('speaker_123');
        expect(baseAct.type).toBe('ask');
        expect(baseAct.id).toMatch(/^act_[a-zA-Z0-9_]+$/);
        expect(baseAct.timestamp).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
      });

      it('should merge additional fields', () => {
        const baseAct = createBaseAct('speaker_123', 'fact', {
          confidence: 0.9,
          source: 'system'
        });
        
        expect(baseAct.confidence).toBe(0.9);
        expect(baseAct.source).toBe('system');
      });
    });
  });

  describe('Constants', () => {
    it('should export version information', () => {
      expect(VERSION).toBe('1.0.0');
      expect(SCHEMA_VERSION).toBe('v1');
    });
  });

  describe('Schemas', () => {
    it('should export all required schemas', () => {
      expect(schemas.act).toBeDefined();
      expect(schemas.ask).toBeDefined();
      expect(schemas.fact).toBeDefined();
      expect(schemas.confirm).toBeDefined();
      expect(schemas.commit).toBeDefined();
      expect(schemas.error).toBeDefined();
      expect(schemas.entity).toBeDefined();
      expect(schemas.participant).toBeDefined();
      expect(schemas.constraint).toBeDefined();
      expect(schemas.conversation).toBeDefined();
    });

    it('should have valid schema structure', () => {
      expect(schemas.act.$schema).toBe('https://json-schema.org/draft/2020-12/schema');
      expect(schemas.act.$id).toBe('https://schemas.astra.dev/v1/act.json');
      expect(schemas.act.title).toBe('Act');
      expect(schemas.act.type).toBe('object');
      expect(schemas.act.required).toEqual(['id', 'timestamp', 'speaker', 'type']);
    });
  });

  describe('Type Definitions', () => {
    it('should properly define Act types', () => {
      const act: Act = {
        id: 'act_001',
        timestamp: '2025-01-15T14:30:00Z',
        speaker: 'agent_123',
        type: 'ask'
      };
      
      expect(act.id).toBe('act_001');
      expect(act.type).toBe('ask');
    });

    it('should properly define specific act types', () => {
      const ask: Ask = {
        id: 'act_001',
        timestamp: '2025-01-15T14:30:00Z',
        speaker: 'agent_123',
        type: 'ask',
        field: 'email',
        prompt: 'What is your email?'
      };
      
      const fact: Fact = {
        id: 'act_002',
        timestamp: '2025-01-15T14:30:00Z',
        speaker: 'customer_456',
        type: 'fact',
        entity: 'customer_456',
        field: 'email',
        value: 'user@example.com'
      };
      
      expect(ask.field).toBe('email');
      expect(fact.value).toBe('user@example.com');
    });

    it('should support EntityRef as string or Entity', () => {
      const factWithStringEntity: Fact = {
        id: 'act_001',
        timestamp: '2025-01-15T14:30:00Z',
        speaker: 'agent_123',
        type: 'fact',
        entity: 'customer_123',
        field: 'name',
        value: 'John Doe'
      };

      const factWithObjectEntity: Fact = {
        id: 'act_002',
        timestamp: '2025-01-15T14:30:00Z',
        speaker: 'agent_123',
        type: 'fact',
        entity: {
          id: 'customer_123',
          type: 'customer',
          external_id: 'cust_ext_456'
        },
        field: 'name',
        value: 'John Doe'
      };

      expect(typeof factWithStringEntity.entity).toBe('string');
      expect(typeof factWithObjectEntity.entity).toBe('object');
    });
  });
});
