# ASTRA Validation Guide

This guide covers ASTRA's multi-layered validation architecture, constraint system, and business rule enforcement mechanisms. ASTRA implements validation at three levels: compile-time type safety, runtime structural validation, and business rule validation, ensuring both correctness and business logic compliance.

## Validation Architecture Overview

ASTRA's validation strategy follows a defense-in-depth approach with multiple validation layers:

1. **Compile-Time Validation** - Static type checking catches structural errors during development
2. **Runtime Schema Validation** - JSON Schema validates data structure and basic constraints at system boundaries  
3. **Business Rule Validation** - Custom constraints enforce domain-specific logic and business requirements
4. **Entity State Validation** - Cross-field validation and complex business invariants

## Layer 1: Compile-Time Validation

### TypeScript Type Safety

TypeScript's type system provides the first line of defense against invalid ASTRA data:

```typescript
// TypeScript catches structural errors at compile time
interface Ask extends Act {
  type: 'ask';
  field: string;        // Required
  prompt: string;       // Required
  constraints?: Constraint[];  // Optional
}

// Compiler error - missing required fields
const invalidAsk: Ask = {
  id: 'act_001',
  timestamp: '2025-01-15T14:30:00Z',
  speaker: 'agent',
  // type: 'ask',     // Missing required field - compile error
  // field: 'email',  // Missing required field - compile error
  // prompt: '...',   // Missing required field - compile error
};

// Valid - all required fields present
const validAsk: Ask = {
  id: 'act_001',
  timestamp: '2025-01-15T14:30:00Z',
  speaker: 'agent',
  type: 'ask',
  field: 'email',
  prompt: 'What is your email address?'
};
```

### Discriminated Unions

TypeScript's discriminated unions ensure type safety across act types:

```typescript
type ConversationAct = Ask | Fact | Confirm | Commit | Error;

function processAct(act: ConversationAct): void {
  switch (act.type) {
    case 'ask':
      // TypeScript knows this is an Ask
      console.log(`Requesting field: ${act.field}`);
      break;
    case 'fact':
      // TypeScript knows this is a Fact
      console.log(`Setting ${act.field} = ${act.value}`);
      break;
    case 'confirm':
      // TypeScript knows this is a Confirm
      console.log(`Confirming: ${act.summary}`);
      break;
    // ... other cases
  }
}
```

### Python Type Hints with Pydantic

Python implementations use Pydantic for runtime type validation:

```python
from pydantic import BaseModel, ValidationError
from typing import Literal, Optional, List

class Ask(BaseModel):
    id: str
    timestamp: str
    speaker: str
    type: Literal['ask']
    field: str
    prompt: str
    constraints: Optional[List[Constraint]] = None

# Pydantic validates at runtime
try:
    ask = Ask(
        id='act_001',
        timestamp='2025-01-15T14:30:00Z',
        speaker='agent',
        type='ask',
        field='email',
        prompt='What is your email address?'
    )
except ValidationError as e:
    print(f"Validation failed: {e}")
```

### Go Struct Validation

Go uses struct tags and explicit validation:

```go
type Ask struct {
    Act
    Field       string       `json:"field" validate:"required"`
    Prompt      string       `json:"prompt" validate:"required"`
    Constraints []Constraint `json:"constraints,omitempty"`
    Required    *bool        `json:"required,omitempty"`
}

// Validate method ensures required fields
func (a Ask) Validate() error {
    if a.Field == "" {
        return fmt.Errorf("field is required")
    }
    if a.Prompt == "" {
        return fmt.Errorf("prompt is required")
    }
    return nil
}
```

## Layer 2: Runtime Schema Validation

### JSON Schema Validation

ASTRA uses JSON Schema for runtime validation at system boundaries:

```typescript
import Ajv from 'ajv';
import addFormats from 'ajv-formats';
import { schemas } from '@astra/model-ts';

// Set up JSON Schema validator
const ajv = new Ajv({ allErrors: true });
addFormats(ajv);

// Compile schemas for performance
const validateAsk = ajv.compile(schemas.ask);
const validateFact = ajv.compile(schemas.fact);
const validateCommit = ajv.compile(schemas.commit);

class SchemaValidator {
  validate(data: unknown, actType: string): ValidationResult {
    const validator = this.getValidator(actType);
    const isValid = validator(data);
    
    return {
      valid: isValid,
      errors: isValid ? [] : validator.errors || [],
      data: isValid ? data : null
    };
  }

  private getValidator(actType: string) {
    switch (actType) {
      case 'ask': return validateAsk;
      case 'fact': return validateFact;
      case 'commit': return validateCommit;
      default: throw new Error(`Unknown act type: ${actType}`);
    }
  }
}

// Usage
const validator = new SchemaValidator();
const result = validator.validate(actData, 'ask');

if (!result.valid) {
  console.error('Validation errors:', result.errors);
}
```

### Multi-Format Schema Validation

ASTRA supports validation across different schema formats:

```python
import jsonschema
from google.protobuf.message import Message
from astra_model.schemas import SCHEMAS

class MultiFormatValidator:
    def __init__(self):
        self.json_schemas = {
            name: jsonschema.validators.validator_for(schema)(schema)
            for name, schema in SCHEMAS.items()
        }
    
    def validate_json_schema(self, data: dict, act_type: str) -> ValidationResult:
        """Validate using JSON Schema"""
        validator = self.json_schemas.get(act_type)
        if not validator:
            raise ValueError(f"No JSON schema for act type: {act_type}")
        
        errors = list(validator.iter_errors(data))
        return ValidationResult(
            valid=len(errors) == 0,
            errors=[self._format_error(error) for error in errors],
            format='json_schema'
        )
    
    def validate_protobuf(self, message: Message) -> ValidationResult:
        """Validate Protocol Buffer message"""
        try:
            # Protobuf validation happens during parsing
            message.SerializeToString()  # Validates required fields
            return ValidationResult(valid=True, errors=[], format='protobuf')
        except Exception as e:
            return ValidationResult(
                valid=False, 
                errors=[str(e)], 
                format='protobuf'
            )
    
    def validate_avro(self, data: dict, schema: dict) -> ValidationResult:
        """Validate using Avro schema"""
        try:
            import avro.schema
            import avro.io
            import io
            
            avro_schema = avro.schema.parse(json.dumps(schema))
            writer = avro.io.DatumWriter(avro_schema)
            bytes_writer = io.BytesIO()
            encoder = avro.io.BinaryEncoder(bytes_writer)
            writer.write(data, encoder)
            
            return ValidationResult(valid=True, errors=[], format='avro')
        except Exception as e:
            return ValidationResult(
                valid=False,
                errors=[str(e)], 
                format='avro'
            )

    def _format_error(self, error) -> str:
        return f"Field '{'.'.join(error.path)}': {error.message}"
```

### Performance-Optimized Validation

For high-throughput scenarios, ASTRA provides optimized validation:

```go
package validation

import (
    "sync"
    "github.com/gojsonschema/gojsonschema"
)

type CachedValidator struct {
    schemaCache map[string]*gojsonschema.Schema
    mutex       sync.RWMutex
}

func NewCachedValidator() *CachedValidator {
    return &CachedValidator{
        schemaCache: make(map[string]*gojsonschema.Schema),
    }
}

func (v *CachedValidator) Validate(data interface{}, schemaName string) (*ValidationResult, error) {
    schema, err := v.getSchema(schemaName)
    if err != nil {
        return nil, err
    }
    
    documentLoader := gojsonschema.NewGoLoader(data)
    result, err := schema.Validate(documentLoader)
    if err != nil {
        return nil, err
    }
    
    return &ValidationResult{
        Valid:  result.Valid(),
        Errors: formatErrors(result.Errors()),
    }, nil
}

func (v *CachedValidator) getSchema(schemaName string) (*gojsonschema.Schema, error) {
    // Read lock for cache lookup
    v.mutex.RLock()
    schema, exists := v.schemaCache[schemaName]
    v.mutex.RUnlock()
    
    if exists {
        return schema, nil
    }
    
    // Write lock for cache update
    v.mutex.Lock()
    defer v.mutex.Unlock()
    
    // Double-check after acquiring write lock
    if schema, exists := v.schemaCache[schemaName]; exists {
        return schema, nil
    }
    
    // Load and compile schema
    schemaJSON := GetSchemaJSON(schemaName)
    schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
    
    compiledSchema, err := gojsonschema.NewSchema(schemaLoader)
    if err != nil {
        return nil, err
    }
    
    v.schemaCache[schemaName] = compiledSchema
    return compiledSchema, nil
}
```

## Layer 3: Constraint System

### Built-in Constraint Types

ASTRA provides comprehensive built-in constraints:

```typescript
// Required/Optional constraints
const requiredConstraint: Constraint = {
  type: 'required',
  message: 'This field is required'
};

// Length constraints
const minLengthConstraint: Constraint = {
  type: 'min_length',
  value: 8,
  message: 'Password must be at least 8 characters'
};

const maxLengthConstraint: Constraint = {
  type: 'max_length',
  value: 255,
  message: 'Field cannot exceed 255 characters'
};

// Pattern matching
const emailPatternConstraint: Constraint = {
  type: 'pattern',
  value: '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$',
  message: 'Please enter a valid email address'
};

// Format validation
const emailFormatConstraint: Constraint = {
  type: 'format',
  value: 'email',
  message: 'Please enter a valid email address'
};

// Range constraints
const ageRangeConstraint: Constraint = {
  type: 'range',
  value: { min: 18, max: 120, inclusive: true },
  message: 'Age must be between 18 and 120'
};

// Enumeration constraints
const sizeEnumConstraint: Constraint = {
  type: 'enum',
  value: ['small', 'medium', 'large'],
  message: 'Size must be small, medium, or large'
};
```

### Constraint Validation Engine

```python
from typing import Any, List, Dict
from enum import Enum
import re
from datetime import datetime

class ConstraintValidator:
    def __init__(self):
        self.validators = {
            'required': self._validate_required,
            'min_length': self._validate_min_length,
            'max_length': self._validate_max_length,
            'pattern': self._validate_pattern,
            'format': self._validate_format,
            'range': self._validate_range,
            'enum': self._validate_enum,
            'custom': self._validate_custom
        }
        
        self.format_validators = {
            'email': self._validate_email_format,
            'phone': self._validate_phone_format,
            'url': self._validate_url_format,
            'date': self._validate_date_format,
            'datetime': self._validate_datetime_format,
            'uuid': self._validate_uuid_format
        }
    
    def validate_constraints(
        self, 
        value: Any, 
        constraints: List[Dict[str, Any]]
    ) -> ValidationResult:
        """Validate a value against multiple constraints"""
        errors = []
        
        for constraint in constraints:
            result = self.validate_constraint(value, constraint)
            if not result.valid:
                errors.extend(result.errors)
        
        return ValidationResult(
            valid=len(errors) == 0,
            errors=errors
        )
    
    def validate_constraint(self, value: Any, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate a single constraint"""
        constraint_type = constraint.get('type')
        validator = self.validators.get(constraint_type)
        
        if not validator:
            return ValidationResult(
                valid=False,
                errors=[f"Unknown constraint type: {constraint_type}"]
            )
        
        try:
            return validator(value, constraint)
        except Exception as e:
            return ValidationResult(
                valid=False,
                errors=[f"Constraint validation error: {str(e)}"]
            )
    
    def _validate_required(self, value: Any, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate required constraint"""
        if value is None or (isinstance(value, str) and value.strip() == ''):
            return ValidationResult(
                valid=False,
                errors=[constraint.get('message', 'This field is required')]
            )
        return ValidationResult(valid=True, errors=[])
    
    def _validate_min_length(self, value: Any, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate minimum length constraint"""
        if not isinstance(value, (str, list, dict)):
            return ValidationResult(valid=True, errors=[])  # Skip non-measurable types
        
        min_length = constraint.get('value')
        if len(value) < min_length:
            message = constraint.get('message', f'Minimum length is {min_length}')
            return ValidationResult(valid=False, errors=[message])
        
        return ValidationResult(valid=True, errors=[])
    
    def _validate_max_length(self, value: Any, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate maximum length constraint"""
        if not isinstance(value, (str, list, dict)):
            return ValidationResult(valid=True, errors=[])
        
        max_length = constraint.get('value')
        if len(value) > max_length:
            message = constraint.get('message', f'Maximum length is {max_length}')
            return ValidationResult(valid=False, errors=[message])
        
        return ValidationResult(valid=True, errors=[])
    
    def _validate_pattern(self, value: Any, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate regex pattern constraint"""
        if not isinstance(value, str):
            return ValidationResult(valid=True, errors=[])
        
        pattern = constraint.get('value')
        if not re.match(pattern, value):
            message = constraint.get('message', f'Value does not match pattern: {pattern}')
            return ValidationResult(valid=False, errors=[message])
        
        return ValidationResult(valid=True, errors=[])
    
    def _validate_format(self, value: Any, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate format constraint"""
        if not isinstance(value, str):
            return ValidationResult(valid=True, errors=[])
        
        format_type = constraint.get('value')
        format_validator = self.format_validators.get(format_type)
        
        if not format_validator:
            return ValidationResult(
                valid=False,
                errors=[f"Unknown format type: {format_type}"]
            )
        
        return format_validator(value, constraint)
    
    def _validate_range(self, value: Any, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate numeric range constraint"""
        if not isinstance(value, (int, float)):
            return ValidationResult(valid=True, errors=[])
        
        range_config = constraint.get('value', {})
        min_val = range_config.get('min')
        max_val = range_config.get('max')
        inclusive = range_config.get('inclusive', True)
        
        if min_val is not None:
            if inclusive and value < min_val:
                message = constraint.get('message', f'Value must be >= {min_val}')
                return ValidationResult(valid=False, errors=[message])
            elif not inclusive and value <= min_val:
                message = constraint.get('message', f'Value must be > {min_val}')
                return ValidationResult(valid=False, errors=[message])
        
        if max_val is not None:
            if inclusive and value > max_val:
                message = constraint.get('message', f'Value must be <= {max_val}')
                return ValidationResult(valid=False, errors=[message])
            elif not inclusive and value >= max_val:
                message = constraint.get('message', f'Value must be < {max_val}')
                return ValidationResult(valid=False, errors=[message])
        
        return ValidationResult(valid=True, errors=[])
    
    def _validate_enum(self, value: Any, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate enumeration constraint"""
        allowed_values = constraint.get('value', [])
        if value not in allowed_values:
            message = constraint.get('message', f'Value must be one of: {allowed_values}')
            return ValidationResult(valid=False, errors=[message])
        
        return ValidationResult(valid=True, errors=[])
    
    def _validate_email_format(self, value: str, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate email format"""
        email_pattern = r'^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
        if not re.match(email_pattern, value):
            message = constraint.get('message', 'Please enter a valid email address')
            return ValidationResult(valid=False, errors=[message])
        return ValidationResult(valid=True, errors=[])
    
    def _validate_phone_format(self, value: str, constraint: Dict[str, Any]) -> ValidationResult:
        """Validate phone number format"""
        # Simple phone validation - can be enhanced based on requirements
        phone_pattern = r'^\+?[\d\s\-\(\)]{10,}$'
        if not re.match(phone_pattern, value):
            message = constraint.get('message', 'Please enter a valid phone number')
            return ValidationResult(valid=False, errors=[message])
        return ValidationResult(valid=True, errors=[])
    
    def _validate_custom(self, value: Any, constraint: Dict[str, Any]) -> ValidationResult:
        """Handle custom constraint validation - to be implemented by business logic"""
        # This is a hook for custom business rule validation
        # Implementation would delegate to business rule engine
        return ValidationResult(valid=True, errors=[])
```

## Layer 4: Business Rule Validation

### Custom Business Rules

ASTRA supports complex business rule validation through custom constraints:

```typescript
interface BusinessRuleEngine {
  validateRule(rule: string, context: ValidationContext): Promise<ValidationResult>;
}

class ASTRABusinessRules implements BusinessRuleEngine {
  constructor(
    private customerService: CustomerService,
    private inventoryService: InventoryService,
    private pricingEngine: PricingEngine
  ) {}

  async validateRule(rule: string, context: ValidationContext): Promise<ValidationResult> {
    switch (rule) {
      case 'customer_exists':
        return this.validateCustomerExists(context);
      case 'inventory_available':
        return this.validateInventoryAvailable(context);
      case 'credit_limit_check':
        return this.validateCreditLimit(context);
      case 'business_hours':
        return this.validateBusinessHours(context);
      default:
        return { valid: false, errors: [`Unknown business rule: ${rule}`] };
    }
  }

  private async validateCustomerExists(context: ValidationContext): Promise<ValidationResult> {
    const customerId = context.entityData.customer_id;
    if (!customerId) {
      return { valid: false, errors: ['Customer ID is required'] };
    }

    const customer = await this.customerService.getById(customerId);
    if (!customer) {
      return { valid: false, errors: [`Customer ${customerId} not found`] };
    }

    if (customer.status === 'inactive') {
      return { valid: false, errors: ['Customer account is inactive'] };
    }

    return { valid: true, errors: [] };
  }

  private async validateInventoryAvailable(context: ValidationContext): Promise<ValidationResult> {
    const items = context.entityData.items || [];
    const errors: string[] = [];

    for (const item of items) {
      const available = await this.inventoryService.checkAvailability(
        item.product_id,
        item.quantity
      );
      
      if (!available) {
        errors.push(`Insufficient inventory for ${item.product_id}`);
      }
    }

    return {
      valid: errors.length === 0,
      errors
    };
  }

  private async validateCreditLimit(context: ValidationContext): Promise<ValidationResult> {
    const customerId = context.entityData.customer_id;
    const orderAmount = context.entityData.total_amount;
    
    if (!customerId || !orderAmount) {
      return { valid: false, errors: ['Customer ID and order amount required for credit check'] };
    }

    const creditInfo = await this.customerService.getCreditInfo(customerId);
    const availableCredit = creditInfo.limit - creditInfo.used;

    if (orderAmount > availableCredit) {
      return {
        valid: false,
        errors: [`Order amount $${orderAmount} exceeds available credit $${availableCredit}`]
      };
    }

    return { valid: true, errors: [] };
  }

  private validateBusinessHours(context: ValidationContext): ValidationResult {
    const requestTime = new Date(context.act.timestamp);
    const hour = requestTime.getHours();
    const day = requestTime.getDay();

    // Business hours: Monday-Friday 9 AM - 5 PM
    if (day === 0 || day === 6) { // Sunday or Saturday
      return {
        valid: false,
        errors: ['Orders cannot be placed on weekends']
      };
    }

    if (hour < 9 || hour >= 17) {
      return {
        valid: false,
        errors: ['Orders can only be placed during business hours (9 AM - 5 PM)']
      };
    }

    return { valid: true, errors: [] };
  }
}
```

### Entity State Validation

Cross-field validation ensures entity consistency:

```python
from typing import Dict, Any, List
from dataclasses import dataclass

@dataclass
class EntityValidationRule:
    name: str
    fields: List[str]
    validator: callable
    message: str

class EntityStateValidator:
    def __init__(self):
        self.rules = {
            'order': [
                EntityValidationRule(
                    name='delivery_address_required',
                    fields=['delivery_method', 'delivery_address'],
                    validator=self._validate_delivery_address_required,
                    message='Delivery address is required for delivery orders'
                ),
                EntityValidationRule(
                    name='payment_method_required',
                    fields=['payment_required', 'payment_method'],
                    validator=self._validate_payment_method_required,
                    message='Payment method is required when payment is needed'
                ),
                EntityValidationRule(
                    name='total_amount_consistency',
                    fields=['items', 'total_amount', 'tax_amount'],
                    validator=self._validate_total_amount_consistency,
                    message='Total amount must match sum of item prices plus tax'
                )
            ]
        }
    
    def validate_entity_state(self, entity_type: str, entity_data: Dict[str, Any]) -> ValidationResult:
        """Validate complete entity state against business rules"""
        rules = self.rules.get(entity_type, [])
        errors = []
        
        for rule in rules:
            # Check if all required fields are present
            if not all(field in entity_data for field in rule.fields):
                continue  # Skip validation if required fields missing
            
            # Extract field values
            field_values = {field: entity_data[field] for field in rule.fields}
            
            # Run validation
            if not rule.validator(field_values):
                errors.append(rule.message)
        
        return ValidationResult(
            valid=len(errors) == 0,
            errors=errors
        )
    
    def _validate_delivery_address_required(self, fields: Dict[str, Any]) -> bool:
        """Validate delivery address is provided when needed"""
        delivery_method = fields.get('delivery_method')
        delivery_address = fields.get('delivery_address')
        
        if delivery_method == 'delivery':
            return bool(delivery_address and delivery_address.strip())
        return True
    
    def _validate_payment_method_required(self, fields: Dict[str, Any]) -> bool:
        """Validate payment method is provided when payment required"""
        payment_required = fields.get('payment_required', True)
        payment_method = fields.get('payment_method')
        
        if payment_required:
            return bool(payment_method and payment_method.strip())
        return True
    
    def _validate_total_amount_consistency(self, fields: Dict[str, Any]) -> bool:
        """Validate total amount matches item prices plus tax"""
        items = fields.get('items', [])
        total_amount = fields.get('total_amount', 0)
        tax_amount = fields.get('tax_amount', 0)
        
        # Calculate expected total
        items_total = sum(
            item.get('price', 0) * item.get('quantity', 0) 
            for item in items
        )
        expected_total = items_total + tax_amount
        
        # Allow small floating point differences
        return abs(expected_total - total_amount) < 0.01
```

## Validation Performance and Optimization

### Caching and Memoization

```go
package validation

import (
    "sync"
    "time"
)

type CacheEntry struct {
    Result    ValidationResult
    ExpiresAt time.Time
}

type CachedConstraintValidator struct {
    validator BusinessRuleValidator
    cache     map[string]*CacheEntry
    mutex     sync.RWMutex
    ttl       time.Duration
}

func NewCachedValidator(validator BusinessRuleValidator, ttl time.Duration) *CachedConstraintValidator {
    return &CachedConstraintValidator{
        validator: validator,
        cache:     make(map[string]*CacheEntry),
        ttl:       ttl,
    }
}

func (v *CachedConstraintValidator) ValidateBusinessRule(
    rule string,
    context ValidationContext,
) (ValidationResult, error) {
    // Create cache key from rule and relevant context
    cacheKey := v.createCacheKey(rule, context)
    
    // Check cache first
    v.mutex.RLock()
    if entry, exists := v.cache[cacheKey]; exists {
        if time.Now().Before(entry.ExpiresAt) {
            v.mutex.RUnlock()
            return entry.Result, nil
        }
    }
    v.mutex.RUnlock()
    
    // Validate using underlying validator
    result, err := v.validator.ValidateRule(rule, context)
    if err != nil {
        return ValidationResult{}, err
    }
    
    // Cache the result
    v.mutex.Lock()
    v.cache[cacheKey] = &CacheEntry{
        Result:    result,
        ExpiresAt: time.Now().Add(v.ttl),
    }
    v.mutex.Unlock()
    
    return result, nil
}

func (v *CachedConstraintValidator) createCacheKey(rule string, context ValidationContext) string {
    // Create deterministic cache key based on rule and context
    // Implementation would hash relevant context fields
    return fmt.Sprintf("%s_%s_%d", rule, context.EntityID, context.EntityVersion)
}

// Background cleanup of expired cache entries
func (v *CachedConstraintValidator) startCacheCleanup() {
    ticker := time.NewTicker(5 * time.Minute)
    go func() {
        for range ticker.C {
            v.cleanupExpiredEntries()
        }
    }()
}

func (v *CachedConstraintValidator) cleanupExpiredEntries() {
    v.mutex.Lock()
    defer v.mutex.Unlock()
    
    now := time.Now()
    for key, entry := range v.cache {
        if now.After(entry.ExpiresAt) {
            delete(v.cache, key)
        }
    }
}
```

### Batch Validation

For high-throughput scenarios, ASTRA supports batch validation:

```python
import asyncio
from typing import List, Dict, Any
from concurrent.futures import ThreadPoolExecutor

class BatchValidator:
    def __init__(self, max_workers: int = 10):
        self.constraint_validator = ConstraintValidator()
        self.business_rule_engine = BusinessRuleEngine()
        self.entity_validator = EntityStateValidator()
        self.max_workers = max_workers
    
    async def validate_acts_batch(
        self, 
        acts: List[Dict[str, Any]]
    ) -> List[ValidationResult]:
        """Validate multiple acts concurrently"""
        
        # Create tasks for concurrent validation
        tasks = [self._validate_single_act(act) for act in acts]
        
        # Execute with controlled concurrency
        results = []
        semaphore = asyncio.Semaphore(self.max_workers)
        
        async def limited_validate(task):
            async with semaphore:
                return await task
        
        limited_tasks = [limited_validate(task) for task in tasks]
        results = await asyncio.gather(*limited_tasks)
        
        return results
    
    async def _validate_single_act(self, act: Dict[str, Any]) -> ValidationResult:
        """Validate a single act with all validation layers"""
        errors = []
        
        # Schema validation
        schema_result = self._validate_schema(act)
        if not schema_result.valid:
            errors.extend(schema_result.errors)
            # Return early on schema validation failure
            return ValidationResult(valid=False, errors=errors)
