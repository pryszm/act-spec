# ASTRA Versioning and Schema Evolution

This guide defines ASTRA's versioning strategy, compatibility requirements, and migration approaches. ASTRA follows semantic versioning principles with specific rules for schema evolution that ensure long-term stability while enabling continuous improvement.

## Versioning Strategy

### Semantic Versioning (SemVer)

ASTRA follows [Semantic Versioning 2.0.0](https://semver.org/) with domain-specific interpretations for schema compatibility:

**Version Format: `MAJOR.MINOR.PATCH`**

- **MAJOR** - Incompatible schema changes that require migration
- **MINOR** - Backward-compatible additions and extensions  
- **PATCH** - Bug fixes and clarifications without schema changes

### Version Scope

ASTRA versioning applies to:

- **IDL Schemas** - JSON Schema, Protocol Buffers, Apache Avro definitions
- **Model Libraries** - TypeScript, Python, Go implementations
- **Validation Rules** - Constraint definitions and business logic
- **API Contracts** - REST API specifications and message formats

## Compatibility Matrix

### Version Compatibility Rules

| Consumer Version | Producer Version | Compatibility | Notes |
|-----------------|------------------|---------------|--------|
| 1.0.x | 1.0.x | ✅ Full | Identical versions |
| 1.0.x | 1.1.x | ✅ Forward | Consumer can read producer data |
| 1.1.x | 1.0.x | ✅ Backward | Producer can write data consumer understands |
| 1.x.x | 2.y.y | ❌ Breaking | Migration required |
| 2.x.x | 1.y.y | ❌ Breaking | Migration required |

### Schema Evolution Rules

#### Patch Version Changes (1.0.0 → 1.0.1)

**Allowed:**
- Documentation improvements and clarifications
- Bug fixes in validation logic (stricter validation only)
- Performance optimizations in libraries
- Non-functional improvements

**Example:**
```json
// v1.0.0 - Original schema
{
  "field": {
    "type": "string",
    "description": "User email address"  // Vague description
  }
}

// v1.0.1 - Improved documentation
{
  "field": {
    "type": "string", 
    "description": "User email address in RFC 5322 format"  // Clearer description
  }
}
```

#### Minor Version Changes (1.0.0 → 1.1.0)

**Allowed:**
- Add optional fields to existing types
- Add new act types to the union
- Extend enums with new values
- Add new constraint types
- Add new utility functions to libraries
- Relax existing validation rules (less strict)

**Forbidden:**
- Remove or rename existing fields
- Change required fields to optional or vice versa
- Remove enum values
- Change field types
- Add required fields without default values

**Examples:**

*Adding Optional Fields:*
```typescript
// v1.0.0
interface Ask extends Act {
  field: string;
  prompt: string;
}

// v1.1.0 - Backward compatible
interface Ask extends Act {
  field: string;
  prompt: string;
  expected_type?: ExpectedType;  // New optional field
  retry_count?: number;          // New optional field
}
```

*Extending Enums:*
```typescript
// v1.0.0
type ActType = 'ask' | 'fact' | 'confirm' | 'commit' | 'error';

// v1.1.0 - Backward compatible
type ActType = 'ask' | 'fact' | 'confirm' | 'commit' | 'error' | 'schedule';
```

*Adding New Act Types:*
```typescript
// v1.0.0
type ConversationAct = Ask | Fact | Confirm | Commit | Error;

// v1.1.0 - Backward compatible
type ConversationAct = Ask | Fact | Confirm | Commit | Error | Schedule;
```

#### Major Version Changes (1.x.x → 2.0.0)

**Examples of Breaking Changes:**

*Removing Required Fields:*
```typescript
// v1.x.x
interface Ask extends Act {
  field: string;        // Required field
  prompt: string;       // Required field
  deprecated_field: string;  // Will be removed
}

// v2.0.0 - BREAKING CHANGE
interface Ask extends Act {
  field: string;
  prompt: string;
  // deprecated_field removed - BREAKING
}
```

*Changing Field Types:*
```typescript
// v1.x.x
interface Fact extends Act {
  confidence?: number;  // Number between 0-1
}

// v2.0.0 - BREAKING CHANGE  
interface Fact extends Act {
  confidence?: ConfidenceLevel;  // Changed to enum - BREAKING
}

type ConfidenceLevel = 'low' | 'medium' | 'high';
```

*Renaming Fields:*
```typescript
// v1.x.x
interface Commit extends Act {
  transaction_id?: string;
}

// v2.0.0 - BREAKING CHANGE
interface Commit extends Act {
  external_transaction_id?: string;  // Field renamed - BREAKING
}
```

## Migration Strategies

### Automated Migration Tools

ASTRA provides migration tools for each breaking change:

```typescript
// Migration tool example
class SchemaV1ToV2Migrator {
  migrate(v1Data: any): any {
    const v2Data = { ...v1Data };
    
    // Handle field renames
    if ('transaction_id' in v2Data) {
      v2Data.external_transaction_id = v2Data.transaction_id;
      delete v2Data.transaction_id;
    }
    
    // Handle type changes
    if (typeof v2Data.confidence === 'number') {
      v2Data.confidence = this.numberToConfidenceLevel(v2Data.confidence);
    }
    
    // Remove deprecated fields
    delete v2Data.deprecated_field;
    
    return v2Data;
  }
  
  private numberToConfidenceLevel(value: number): ConfidenceLevel {
    if (value < 0.33) return 'low';
    if (value < 0.67) return 'medium';
    return 'high';
  }
}
```

### Gradual Migration Pattern

For production systems, ASTRA supports gradual migration using version-aware processors:

```python
class VersionAwareProcessor:
    def __init__(self):
        self.v1_processor = V1ActProcessor()
        self.v2_processor = V2ActProcessor()
        self.migrator = SchemaV1ToV2Migrator()
    
    async def process_act(self, act_data: dict) -> dict:
        version = self.detect_version(act_data)
        
        if version == "1.x.x":
            # Process with v1 logic and optionally migrate
            if self.should_migrate_to_v2():
                migrated_data = self.migrator.migrate(act_data)
                return await self.v2_processor.process(migrated_data)
            else:
                return await self.v1_processor.process(act_data)
        
        elif version == "2.x.x":
            return await self.v2_processor.process(act_data)
        
        else:
            raise UnsupportedVersionError(f"Version {version} not supported")
    
    def detect_version(self, act_data: dict) -> str:
        # Version detection logic based on schema markers
        if 'transaction_id' in act_data:
            return "1.x.x"
        elif 'external_transaction_id' in act_data:
            return "2.x.x"
        else:
            # Use schema validation to determine version
            return self.validate_and_detect_version(act_data)
```

### Migration Testing Strategy

Every migration includes comprehensive test coverage:

```typescript
describe('Schema Migration v1 → v2', () => {
  const migrator = new SchemaV1ToV2Migrator();
  
  test('migrates transaction_id field correctly', () => {
    const v1Act = {
      id: 'act_001',
      type: 'commit',
      transaction_id: 'txn_123'
    };
    
    const v2Act = migrator.migrate(v1Act);
    
    expect(v2Act.external_transaction_id).toBe('txn_123');
    expect(v2Act.transaction_id).toBeUndefined();
  });
  
  test('converts confidence number to enum', () => {
    const v1Act = {
      id: 'act_002', 
      type: 'fact',
      confidence: 0.8
    };
    
    const v2Act = migrator.migrate(v1Act);
    
    expect(v2Act.confidence).toBe('high');
  });
  
  test('removes deprecated fields', () => {
    const v1Act = {
      id: 'act_003',
      type: 'ask',
      deprecated_field: 'old_value'
    };
    
    const v2Act = migrator.migrate(v1Act);
    
    expect(v2Act.deprecated_field).toBeUndefined();
  });
  
  test('preserves all other fields unchanged', () => {
    const v1Act = {
      id: 'act_004',
      timestamp: '2025-01-15T14:30:00Z',
      speaker: 'agent_123',
      type: 'ask',
      field: 'email',
      prompt: 'What is your email?'
    };
    
    const v2Act = migrator.migrate(v1Act);
    
    expect(v2Act.id).toBe(v1Act.id);
    expect(v2Act.timestamp).toBe(v1Act.timestamp);
    expect(v2Act.speaker).toBe(v1Act.speaker);
    expect(v2Act.type).toBe(v1Act.type);
    expect(v2Act.field).toBe(v1Act.field);
    expect(v2Act.prompt).toBe(v1Act.prompt);
  });
});
```

## Deprecation Policy

### Deprecation Timeline

ASTRA follows a structured deprecation process:

1. **Announcement** (Minor version) - Field/feature marked as deprecated with replacement guidance
2. **Warning Period** (Minimum 6 months) - Deprecation warnings in logs and documentation
3. **Removal** (Next major version) - Deprecated field/feature removed

### Deprecation Marking

**In Schemas:**
```json
{
  "field": {
    "type": "string",
    "deprecated": true,
    "description": "DEPRECATED: Use new_field instead. Will be removed in v2.0.0"
  },
  "new_field": {
    "type": "string", 
    "description": "Replacement for deprecated field"
  }
}
```

**In TypeScript:**
```typescript
interface Ask extends Act {
  field: string;
  prompt: string;
  
  /** @deprecated Use expected_type instead. Will be removed in v2.0.0 */
  response_type?: string;
  
  expected_type?: ExpectedType;
}
```

**In Python:**
```python
from pydantic import BaseModel, Field
from warnings import warn

class Ask(Act):
    field: str
    prompt: str
    response_type: Optional[str] = Field(
        None, 
        deprecated=True,
        description="DEPRECATED: Use expected_type instead"
    )
    expected_type: Optional[ExpectedType] = None
    
    def __init__(self, **data):
        if 'response_type' in data:
            warn(
                "response_type is deprecated. Use expected_type instead.",
                DeprecationWarning,
                stacklevel=2
            )
        super().__init__(**data)
```

### Migration Guidance

Each deprecation includes clear migration guidance:

```markdown
## Migration Guide: response_type → expected_type

### What's Changing
The `response_type` field in Ask acts is deprecated in favor of the more structured `expected_type` field.

### Timeline  
- **v1.5.0** - `expected_type` introduced, `response_type` deprecated
- **v1.6.0+** - Deprecation warnings logged when `response_type` used
- **v2.0.0** - `response_type` removed

### Migration Steps

1. **Update act creation:**
   ```typescript
   // Old (deprecated)
   const ask: Ask = {
     // ... other fields
     response_type: 'email'
   };
   
   // New (recommended)
   const ask: Ask = {
     // ... other fields  
     expected_type: 'email'
   };
   ```

2. **Update validation logic:**
   ```typescript
   // Old
   if (act.response_type === 'email') { /* ... */ }
   
   // New
   if (act.expected_type === 'email') { /* ... */ }
   ```

3. **Test thoroughly** with new field before removing old field usage
```

## Version Detection and Validation

### Schema Versioning

All ASTRA schemas include explicit version metadata:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://schemas.astra.dev/v1.5.0/act.json",
  "title": "Act",
  "version": "1.5.0",
  "description": "Base type for all conversational actions in ASTRA v1.5.0"
}
```

### Runtime Version Detection

```typescript
class VersionDetector {
  static detectSchemaVersion(data: any): string {
    // Check for version-specific markers
    const versionMarkers = [
      { version: '2.0.0', marker: (d: any) => 'external_transaction_id' in d },
      { version: '1.5.0', marker: (d: any) => 'expected_type' in d },
      { version: '1.0.0', marker: (d: any) => true } // Fallback
    ];
    
    for (const { version, marker } of versionMarkers) {
      if (marker(data)) {
        return version;
      }
    }
    
    throw new Error('Unable to detect schema version');
  }
  
  static validateVersion(data: any, expectedVersion: string): boolean {
    const detectedVersion = this.detectSchemaVersion(data);
    return semver.satisfies(detectedVersion, expectedVersion);
  }
}
```

### Multi-Version Support

Production systems can support multiple ASTRA versions simultaneously:

```python
class MultiVersionProcessor:
    def __init__(self):
        self.processors = {
            '1.0.0': V1_0_Processor(),
            '1.5.0': V1_5_Processor(), 
            '2.0.0': V2_0_Processor()
        }
        self.migrators = {
            ('1.0.0', '1.5.0'): V1_0_to_V1_5_Migrator(),
            ('1.5.0', '2.0.0'): V1_5_to_V2_0_Migrator(),
            ('1.0.0', '2.0.0'): ChainedMigrator(['1.0.0', '1.5.0', '2.0.0'])
        }
    
    async def process(self, act_data: dict, target_version: str = None):
        current_version = VersionDetector.detect_version(act_data)
        target_version = target_version or self.get_latest_version()
        
        # Migrate if needed
        if current_version != target_version:
            act_data = await self.migrate(act_data, current_version, target_version)
        
        # Process with appropriate version processor
        processor = self.processors[target_version]
        return await processor.process(act_data)
    
    async def migrate(self, data: dict, from_version: str, to_version: str):
        migrator_key = (from_version, to_version)
        
        if migrator_key in self.migrators:
            return self.migrators[migrator_key].migrate(data)
        else:
            # Find migration path through intermediate versions
            return await self.find_migration_path(data, from_version, to_version)
```

## Testing Version Compatibility

### Compatibility Test Suite

```typescript
describe('Version Compatibility', () => {
  const testCases = [
    {
      name: 'v1.0.0 consumer reads v1.1.0 data',
      consumer: '1.0.0',
      producer: '1.1.0', 
      expected: 'compatible'
    },
    {
      name: 'v1.1.0 consumer reads v1.0.0 data',
      consumer: '1.1.0',
      producer: '1.0.0',
      expected: 'compatible'  
    },
    {
      name: 'v1.x.x consumer reads v2.0.0 data',
      consumer: '1.5.0',
      producer: '2.0.0',
      expected: 'incompatible'
    }
  ];
  
  testCases.forEach(({ name, consumer, producer, expected }) => {
    test(name, async () => {
      const producerProcessor = createProcessor(producer);
      const consumerProcessor = createProcessor(consumer);
      
      // Generate test data with producer
      const testData = await producerProcessor.generateTestActs();
      
      // Attempt to process with consumer
      const results = await Promise.allSettled(
        testData.map(data => consumerProcessor.process(data))
      );
      
      if (expected === 'compatible') {
        expect(results.every(r => r.status === 'fulfilled')).toBe(true);
      } else {
        expect(results.some(r => r.status === 'rejected')).toBe(true);
      }
    });
  });
});
```

### Breaking Change Detection

Automated tools detect potential breaking changes in schema evolution:

```python
class BreakingChangeDetector:
    def analyze_changes(self, old_schema: dict, new_schema: dict) -> List[BreakingChange]:
        breaking_changes = []
        
        # Check for removed required fields
        old_required = set(old_schema.get('required', []))
        new_required = set(new_schema.get('required', []))
        
        for field in old_required - new_required:
            breaking_changes.append(BreakingChange(
                type='required_field_removed',
                field=field,
                severity='major',
                description=f'Required field "{field}" was removed'
            ))
        
        # Check for type changes
        old_properties = old_schema.get('properties', {})
        new_properties = new_schema.get('properties', {})
        
        for field, old_def in old_properties.items():
            if field in new_properties:
                new_def = new_properties[field]
                if old_def.get('type') != new_def.get('type'):
                    breaking_changes.append(BreakingChange(
                        type='type_change',
                        field=field,
                        severity='major',
                        description=f'Field "{field}" type changed from {old_def.get("type")} to {new_def.get("type")}'
                    ))
        
        # Check for enum value removals
        for field, old_def in old_properties.items():
            if field in new_properties:
                new_def = new_properties[field]
                old_enum = set(old_def.get('enum', []))
                new_enum = set(new_def.get('enum', []))
                
                removed_values = old_enum - new_enum
                if removed_values:
                    breaking_changes.append(BreakingChange(
                        type='enum_values_removed',
                        field=field,
                        severity='major',
                        description=f'Enum values removed from "{field}": {removed_values}'
                    ))
        
        return breaking_changes
```

## Version Management Best Practices

### Release Planning

1. **Feature Planning** - Group related features into minor releases
2. **Breaking Change Batching** - Accumulate breaking changes for major releases
3. **Deprecation Timeline** - Plan deprecations at least one minor version before removal
4. **Migration Tool Development** - Build migration tools before releasing breaking changes

### Documentation Requirements

Each version must include:

- **CHANGELOG.md** - Detailed changes with examples
- **MIGRATION.md** - Step-by-step migration instructions  
- **COMPATIBILITY.md** - Version compatibility matrix
- **DEPRECATION.md** - List of deprecated features and timelines

### CI/CD Integration

Version management is integrated into the development workflow:

```yaml
# GitHub Actions workflow
name: Schema Version Validation

on:
  pull_request:
    paths: ['idl/**', 'model/**']

jobs:
  detect-breaking-changes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      
      - name: Detect breaking changes
        run: |
          python tools/breaking-change-detector.py \
            --base-ref origin/main \
            --head-ref HEAD \
            --schemas idl/json-schema/
      
      - name: Validate version bump
        run: |
          python tools/validate-version-bump.py \
            --breaking-changes breaking-changes.json \
            --version-bump $(git diff origin/main..HEAD package.json | grep version)
      
      - name: Generate migration tools
        if: contains(github.event.pull_request.labels.*.name, 'breaking-change')
        run: |
          python tools/generate-migration.py \
            --from-version $(git show origin/main:package.json | jq .version) \
            --to-version $(jq .version package.json)
```

## Conclusion

ASTRA's versioning strategy balances stability with evolution, enabling long-term compatibility while supporting continuous improvement. By following semantic versioning principles with domain-specific rules, providing comprehensive migration tools, and maintaining strict compatibility testing, ASTRA ensures that conversational applications can evolve reliably over time.

The key to successful schema evolution is planning for change from the beginning - designing extensible schemas, following deprecation processes, and providing migration support for every breaking change. This approach enables the ASTRA ecosystem to grow and adapt while maintaining the reliability that production systems require.
