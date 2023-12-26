/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\schema\schema_test.go
 * @Description: schema.go 测试，覆盖 JSON Schema 子集校验、SchemaBuilder 和 QuickSchema
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package schema

import (
	"encoding/json"
	"testing"

	"github.com/kamalyes/go-argus/validate"
)

func TestValidateJSONSchemaObjectSuccess(t *testing.T) {
	min := 2
	max := 8
	s := JSONSchema{
		Type:     "object",
		Required: []string{"name"},
		Properties: map[string]JSONSchema{
			"name": {Type: "string", MinLength: &min, MaxLength: &max},
		},
	}
	result := ValidateJSONSchema(map[string]interface{}{"name": "argus"}, s)
	if !result.Success {
		t.Fatalf("expected schema validation to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaObjectFailMinLength(t *testing.T) {
	min := 2
	max := 8
	s := JSONSchema{
		Type:     "object",
		Required: []string{"name"},
		Properties: map[string]JSONSchema{
			"name": {Type: "string", MinLength: &min, MaxLength: &max},
		},
	}
	result := ValidateJSONSchema(map[string]interface{}{"name": "a"}, s)
	if result.Success {
		t.Fatal("expected short name to fail")
	}
}

func TestValidateJSONSchemaObjectFailMaxLength(t *testing.T) {
	min := 1
	max := 3
	s := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"name": {Type: "string", MinLength: &min, MaxLength: &max},
		},
	}
	result := ValidateJSONSchema(map[string]interface{}{"name": "toolong"}, s)
	if result.Success {
		t.Fatal("expected long name to fail")
	}
}

func TestValidateJSONSchemaMissingRequired(t *testing.T) {
	s := JSONSchema{
		Type:     "object",
		Required: []string{"email"},
		Properties: map[string]JSONSchema{
			"email": {Type: "string"},
		},
	}
	result := ValidateJSONSchema(map[string]interface{}{}, s)
	if result.Success {
		t.Fatal("expected missing required field to fail")
	}
}

func TestValidateJSONSchemaTypeMismatch(t *testing.T) {
	s := JSONSchema{Type: "string"}
	result := ValidateJSONSchema(42, s)
	if result.Success {
		t.Fatal("expected type mismatch to fail")
	}
}

func TestValidateJSONSchemaEnumMismatch(t *testing.T) {
	s := JSONSchema{Type: "string", Enum: []interface{}{"red", "green", "blue"}}
	result := ValidateJSONSchema("yellow", s)
	if result.Success {
		t.Fatal("expected enum mismatch to fail")
	}
}

func TestValidateJSONSchemaEnumMatch(t *testing.T) {
	s := JSONSchema{Type: "string", Enum: []interface{}{"red", "green", "blue"}}
	result := ValidateJSONSchema("red", s)
	if !result.Success {
		t.Fatalf("expected enum match to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaNumberBelowMinimum(t *testing.T) {
	min := 10.0
	s := JSONSchema{Type: "number", Minimum: &min}
	result := ValidateJSONSchema(5.0, s)
	if result.Success {
		t.Fatal("expected below minimum to fail")
	}
}

func TestValidateJSONSchemaNumberAboveMaximum(t *testing.T) {
	max := 100.0
	s := JSONSchema{Type: "number", Maximum: &max}
	result := ValidateJSONSchema(200.0, s)
	if result.Success {
		t.Fatal("expected above maximum to fail")
	}
}

func TestValidateJSONSchemaNumberInRange(t *testing.T) {
	min := 10.0
	max := 100.0
	s := JSONSchema{Type: "number", Minimum: &min, Maximum: &max}
	result := ValidateJSONSchema(50.0, s)
	if !result.Success {
		t.Fatalf("expected number in range to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaInteger(t *testing.T) {
	s := JSONSchema{Type: "integer"}
	result := ValidateJSONSchema(42.0, s)
	if !result.Success {
		t.Fatalf("expected integer to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaIntegerFail(t *testing.T) {
	s := JSONSchema{Type: "integer"}
	result := ValidateJSONSchema(3.14, s)
	if result.Success {
		t.Fatal("expected non-integer float to fail")
	}
}

func TestValidateJSONSchemaArray(t *testing.T) {
	s := JSONSchema{
		Type:  "array",
		Items: &JSONSchema{Type: "string"},
	}
	result := ValidateJSONSchema([]interface{}{"a", "b"}, s)
	if !result.Success {
		t.Fatalf("expected array to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaArrayItemFail(t *testing.T) {
	s := JSONSchema{
		Type:  "array",
		Items: &JSONSchema{Type: "string"},
	}
	result := ValidateJSONSchema([]interface{}{"a", 123}, s)
	if result.Success {
		t.Fatal("expected array item type mismatch to fail")
	}
}

func TestValidateJSONSchemaArrayNoItems(t *testing.T) {
	s := JSONSchema{Type: "array"}
	result := ValidateJSONSchema([]interface{}{"a", 123}, s)
	if !result.Success {
		t.Fatalf("expected array without items to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaBoolean(t *testing.T) {
	s := JSONSchema{Type: "boolean"}
	result := ValidateJSONSchema(true, s)
	if !result.Success {
		t.Fatalf("expected boolean to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaNull(t *testing.T) {
	s := JSONSchema{Type: "null"}
	result := ValidateJSONSchema(nil, s)
	if !result.Success {
		t.Fatalf("expected null to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaUnknownType(t *testing.T) {
	s := JSONSchema{Type: "custom"}
	result := ValidateJSONSchema("anything", s)
	if !result.Success {
		t.Fatalf("expected unknown type to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaObjectWithProperties(t *testing.T) {
	s := JSONSchema{
		Properties: map[string]JSONSchema{
			"name": {Type: "string"},
		},
	}
	result := ValidateJSONSchema(map[string]interface{}{"name": "argus"}, s)
	if !result.Success {
		t.Fatalf("expected object with properties to pass: %s", result.Message)
	}
}

func TestValidateJSONSchemaObjectPropertyFail(t *testing.T) {
	s := JSONSchema{
		Properties: map[string]JSONSchema{
			"age": {Type: "number"},
		},
	}
	result := ValidateJSONSchema(map[string]interface{}{"age": "not-a-number"}, s)
	if result.Success {
		t.Fatal("expected object property type mismatch to fail")
	}
}

func TestValidateJSONSchemaNonObjectForProperties(t *testing.T) {
	s := JSONSchema{
		Type: "string",
		Properties: map[string]JSONSchema{
			"name": {Type: "string"},
		},
	}
	result := ValidateJSONSchema("hello", s)
	if !result.Success {
		t.Fatalf("expected non-object with properties to pass type check: %s", result.Message)
	}
}

func TestValidateJSONSchemaNonStringForStringRules(t *testing.T) {
	min := 1
	s := JSONSchema{Type: "string", MinLength: &min}
	result := ValidateJSONSchema(123, s)
	if result.Success {
		t.Fatal("expected type mismatch to fail before string rules")
	}
}

func TestValidateJSONSchemaNonNumberForNumberRules(t *testing.T) {
	min := 10.0
	s := JSONSchema{Type: "number", Minimum: &min}
	result := ValidateJSONSchema("hello", s)
	if result.Success {
		t.Fatal("expected type mismatch to fail before number rules")
	}
}

func TestValidateJSONSchemaNonArrayForArrayRules(t *testing.T) {
	s := JSONSchema{Type: "array", Items: &JSONSchema{Type: "string"}}
	result := ValidateJSONSchema("not-array", s)
	if result.Success {
		t.Fatal("expected type mismatch to fail before array rules")
	}
}

func TestValidateJSONSchemaNonObjectForObjectRules(t *testing.T) {
	s := JSONSchema{Type: "object", Required: []string{"name"}}
	result := ValidateJSONSchema(42, s)
	if result.Success {
		t.Fatal("expected type mismatch to fail before object rules")
	}
}

func TestValidateStructWithSchema(t *testing.T) {
	type user struct {
		Name string `json:"name"`
	}
	min := 1
	s := JSONSchema{
		Type:     "object",
		Required: []string{"name"},
		Properties: map[string]JSONSchema{
			"name": {Type: "string", MinLength: &min},
		},
	}
	result := ValidateStructWithSchema(user{Name: "argus"}, s)
	if !result.Success {
		t.Fatalf("expected struct validation to pass: %s", result.Message)
	}
}

func TestValidateStructWithSchemaFail(t *testing.T) {
	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	min := 1
	s := JSONSchema{
		Type:     "object",
		Required: []string{"name"},
		Properties: map[string]JSONSchema{
			"name": {Type: "string", MinLength: &min},
		},
	}
	result := ValidateStructWithSchema(user{Name: "", Age: 0}, s)
	if result.Success {
		t.Fatal("expected struct validation to fail for short name")
	}
}

func TestValidateStructWithSchemaMarshalError(t *testing.T) {
	result := ValidateStructWithSchema(make(chan int), JSONSchema{})
	if result.Success {
		t.Fatal("expected marshal error")
	}
}

func TestNormalizeSchemaFromJSONSchema(t *testing.T) {
	s := JSONSchema{Type: "string"}
	result, err := normalizeSchema(s)
	if err != nil || result.Type != "string" {
		t.Fatalf("unexpected: %v %v", result, err)
	}
}

func TestNormalizeSchemaFromPtr(t *testing.T) {
	s := JSONSchema{Type: "number"}
	result, err := normalizeSchema(&s)
	if err != nil || result.Type != "number" {
		t.Fatalf("unexpected: %v %v", result, err)
	}
}

func TestNormalizeSchemaFromNilPtr(t *testing.T) {
	var s *JSONSchema
	_, err := normalizeSchema(s)
	if err == nil {
		t.Fatal("expected error for nil pointer")
	}
}

func TestNormalizeSchemaFromBytes(t *testing.T) {
	raw := `{"type":"string"}`
	result, err := normalizeSchema([]byte(raw))
	if err != nil || result.Type != "string" {
		t.Fatalf("unexpected: %v %v", result, err)
	}
}

func TestNormalizeSchemaFromString(t *testing.T) {
	raw := `{"type":"integer"}`
	result, err := normalizeSchema(raw)
	if err != nil || result.Type != "integer" {
		t.Fatalf("unexpected: %v %v", result, err)
	}
}

func TestNormalizeSchemaFromInvalidString(t *testing.T) {
	_, err := normalizeSchema("not-json")
	if err == nil {
		t.Fatal("expected error for invalid JSON string")
	}
}

func TestNormalizeSchemaFromMap(t *testing.T) {
	m := map[string]interface{}{"type": "boolean"}
	result, err := normalizeSchema(m)
	if err != nil || result.Type != "boolean" {
		t.Fatalf("unexpected: %v %v", result, err)
	}
}

func TestNormalizeSchemaFromInvalidType(t *testing.T) {
	_, err := normalizeSchema(make(chan int))
	if err == nil {
		t.Fatal("expected error for unmarshalable type")
	}
}

func TestSchemaBuilder(t *testing.T) {
	min := 1
	max := 100
	minLen := 2.0
	maxLen := 200.0
	b := NewSchemaBuilder().
		Type("object").
		Required("name", "age").
		StringProperty("name", min, max).
		NumberProperty("age", &minLen, &maxLen).
		ArrayProperty("tags", JSONSchema{Type: "string"}).
		Enum("admin", "member")
	s := b.Build()
	if s.Type != "object" {
		t.Fatalf("expected object type, got %s", s.Type)
	}
	if len(s.Required) != 2 {
		t.Fatalf("expected 2 required fields, got %d", len(s.Required))
	}
	if len(s.Properties) != 3 {
		t.Fatalf("expected 3 properties, got %d", len(s.Properties))
	}
	if len(s.Enum) != 2 {
		t.Fatalf("expected 2 enum values, got %d", len(s.Enum))
	}
}

func TestSchemaBuilderPropertyNilMap(t *testing.T) {
	b := &SchemaBuilder{}
	b.Property("test", JSONSchema{Type: "string"})
	if b.schema.Properties == nil || len(b.schema.Properties) != 1 {
		t.Fatal("expected property map to be initialized")
	}
}

func TestSchemaBuilderStringPropertyNegative(t *testing.T) {
	b := NewSchemaBuilder()
	b.StringProperty("name", -1, -1)
	s := b.Build()
	prop := s.Properties["name"]
	if prop.MinLength != nil || prop.MaxLength != nil {
		t.Fatal("expected nil min/maxLength for negative values")
	}
}

func TestSchemaBuilderBuildJSON(t *testing.T) {
	b := NewSchemaBuilder().Type("string")
	jsonStr := b.BuildJSON()
	var parsed JSONSchema
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}
	if parsed.Type != "string" {
		t.Fatalf("expected string type, got %s", parsed.Type)
	}
}

func TestQuickSchema(t *testing.T) {
	s := QuickSchema(map[string]string{"name": "string", "age": "number"}, "name")
	if s.Type != "object" {
		t.Fatalf("expected object type, got %s", s.Type)
	}
	if len(s.Required) != 1 || s.Required[0] != "name" {
		t.Fatalf("unexpected required: %v", s.Required)
	}
	if len(s.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(s.Properties))
	}
}

func TestFormatSchemaErrorSuccess(t *testing.T) {
	result := ValidateCompareResult(true)
	msg := FormatSchemaError(result)
	if msg != "" {
		t.Fatalf("expected empty message for success, got %q", msg)
	}
}

func TestFormatSchemaErrorFailure(t *testing.T) {
	result := ValidateJSONSchema(42, JSONSchema{Type: "string"})
	msg := FormatSchemaError(result)
	if msg == "" {
		t.Fatal("expected non-empty message for failure")
	}
}

func ValidateCompareResult(success bool) validate.CompareResult {
	return validate.CompareResult{Success: success}
}

func TestValidateJSONSchemaNormalizeError(t *testing.T) {
	result := ValidateJSONSchema("data", "not-json")
	if result.Success {
		t.Fatal("expected normalize error to fail")
	}
}

func TestValidateStructWithSchemaUnmarshalError(t *testing.T) {
	result := ValidateStructWithSchema(map[string]interface{}{"a": make(chan int)}, JSONSchema{})
	if result.Success {
		t.Fatal("expected unmarshal error to fail")
	}
}

func TestValidateStringNonString(t *testing.T) {
	min := 1
	s := JSONSchema{MinLength: &min}
	result := ValidateJSONSchema(42, s)
	if !result.Success {
		t.Fatalf("expected non-string without type to pass: %s", result.Message)
	}
}

func TestValidateNumberNonNumber(t *testing.T) {
	min := 10.0
	s := JSONSchema{Minimum: &min}
	result := ValidateJSONSchema("hello", s)
	if !result.Success {
		t.Fatalf("expected non-number without type to pass: %s", result.Message)
	}
}

func TestValidateObjectNonObject(t *testing.T) {
	s := JSONSchema{Required: []string{"name"}}
	result := ValidateJSONSchema(42, s)
	if !result.Success {
		t.Fatalf("expected non-object without type to pass: %s", result.Message)
	}
}
