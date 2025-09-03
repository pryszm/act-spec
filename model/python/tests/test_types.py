"""
Tests for ASTRA model schemas
"""
try:
    import pytest
    HAS_PYTEST = True
except ImportError:
    HAS_PYTEST = False

from astra_model.schemas import SCHEMAS


def test_schemas_exist():
    """Test that all expected schemas are present"""
    expected_schemas = {
        'act', 'ask', 'fact', 'confirm', 'commit', 
        'error', 'entity', 'participant', 'constraint', 'conversation'
    }
    actual_schemas = set(SCHEMAS.keys())
    assert expected_schemas == actual_schemas


def test_schema_structure():
    """Test that each schema has required JSON Schema properties"""
    for schema_name, schema in SCHEMAS.items():
        assert '$schema' in schema, f"{schema_name} missing $schema"
        assert '$id' in schema, f"{schema_name} missing $id"
        assert 'title' in schema, f"{schema_name} missing title"
        assert 'type' in schema, f"{schema_name} missing type"
        assert schema['type'] == 'object', f"{schema_name} should be object type"


def test_act_schemas_have_const_types():
    """Test that specific act schemas have correct const types"""
    act_types = ['ask', 'fact', 'confirm', 'commit', 'error']
    for act_type in act_types:
        schema = SCHEMAS[act_type]
        type_prop = schema['properties']['type']
        assert 'const' in type_prop, f"{act_type} missing const type"
        assert type_prop['const'] == act_type, f"{act_type} has wrong const value"


def test_required_fields():
    """Test that schemas have expected required fields"""
    # Base act should have these required fields
    base_required = ['id', 'timestamp', 'speaker', 'type']
    for field in base_required:
        assert field in SCHEMAS['act']['required']
    
    # Entity should require id and type
    assert 'id' in SCHEMAS['entity']['required']
    assert 'type' in SCHEMAS['entity']['required']
    
    # Conversation should require id, participants, acts
    conv_required = ['id', 'participants', 'acts']
    for field in conv_required:
        assert field in SCHEMAS['conversation']['required']


def test_schemas_serializable():
    """Test that schemas can be serialized to JSON"""
    import json
    for schema_name, schema in SCHEMAS.items():
        try:
            json.dumps(schema)
        except Exception as e:
            if HAS_PYTEST:
                pytest.fail(f"Schema {schema_name} not JSON serializable: {e}")
            else:
                raise AssertionError(f"Schema {schema_name} not JSON serializable: {e}")


if __name__ == "__main__":
    # Run basic tests when script is executed directly
    test_schemas_exist()
    test_schema_structure()
    test_act_schemas_have_const_types()
    test_required_fields()
    test_schemas_serializable()
    print("All tests passed!")
