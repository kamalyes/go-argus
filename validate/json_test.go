/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\json_test.go
 * @Description: json.go 测试，覆盖 JSON 有效性、字段读取和轻量路径匹配
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validate

import (
	"testing"
)

func TestValidateJSONValid(t *testing.T) {
	if err := ValidateJSON([]byte(`{"key":"value"}`)); err != nil {
		t.Fatal("expected valid JSON")
	}
}

func TestValidateJSONInvalid(t *testing.T) {
	if err := ValidateJSON([]byte(`{invalid}`)); err == nil {
		t.Fatal("expected invalid JSON to fail")
	}
}

func TestIsJSONNull(t *testing.T) {
	if !IsJSONNull([]byte("null")) {
		t.Fatal("expected null")
	}
	if IsJSONNull([]byte(`{"a":1}`)) {
		t.Fatal("expected not null")
	}
}

func TestIsJSONColumnType(t *testing.T) {
	if !IsJSONColumnType("json") {
		t.Fatal("expected json type")
	}
	if !IsJSONColumnType("JSONB") {
		t.Fatal("expected jsonb type")
	}
	if IsJSONColumnType("text") {
		t.Fatal("expected text not to be json type")
	}
}

func TestValidateJSONWithData(t *testing.T) {
	data, err := ValidateJSONWithData([]byte(`{"key":"value"}`))
	if err != nil {
		t.Fatal("expected valid JSON")
	}
	m, ok := data.(map[string]interface{})
	if !ok || m["key"] != "value" {
		t.Fatal("expected key=value")
	}
}

func TestValidateJSONWithDataInvalid(t *testing.T) {
	_, err := ValidateJSONWithData([]byte(`{invalid}`))
	if err == nil {
		t.Fatal("expected invalid JSON to fail")
	}
}

func TestValidateJSONField(t *testing.T) {
	r := ValidateJSONField([]byte(`{"name":"argus"}`), "name", "argus")
	if !r.Success {
		t.Fatal("expected field match")
	}
}

func TestValidateJSONFieldInvalidJSON(t *testing.T) {
	r := ValidateJSONField([]byte(`{invalid}`), "name", "argus")
	if r.Success || r.Message == "" {
		t.Fatal("expected invalid JSON to fail")
	}
}

func TestValidateJSONFieldNotObject(t *testing.T) {
	r := ValidateJSONField([]byte(`[1,2,3]`), "name", "argus")
	if r.Success || r.Message == "" {
		t.Fatal("expected root not object to fail")
	}
}

func TestValidateJSONFieldNotFound(t *testing.T) {
	r := ValidateJSONField([]byte(`{"name":"argus"}`), "age", 25)
	if r.Success || r.Message == "" {
		t.Fatal("expected field not found to fail")
	}
}

func TestValidateJSONFieldValueMismatch(t *testing.T) {
	r := ValidateJSONField([]byte(`{"name":"argus"}`), "name", "other")
	if r.Success {
		t.Fatal("expected value mismatch")
	}
}

func TestValidateJSONFields(t *testing.T) {
	results := ValidateJSONFields([]byte(`{"name":"argus","ver":1}`), map[string]any{
		"name": "argus",
		"ver":  1,
	})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestLookupJSONPath(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "argus",
		},
	}
	v, ok := LookupJSONPath(data, "$.user.name")
	if !ok || v != "argus" {
		t.Fatal("expected to find user.name=argus")
	}
}

func TestLookupJSONPathEmpty(t *testing.T) {
	data := "hello"
	v, ok := LookupJSONPath(data, "")
	if !ok || v != "hello" {
		t.Fatal("expected empty path to return data")
	}
}

func TestLookupJSONPathNotFound(t *testing.T) {
	data := map[string]any{"name": "argus"}
	_, ok := LookupJSONPath(data, "missing")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestLookupJSONPathNotMap(t *testing.T) {
	data := "hello"
	_, ok := LookupJSONPath(data, "name")
	if ok {
		t.Fatal("expected not map to fail")
	}
}

func TestLookupJSONPathArray(t *testing.T) {
	data := map[string]any{
		"items": []any{"a", "b", "c"},
	}
	v, ok := LookupJSONPath(data, "items[1]")
	if !ok || v != "b" {
		t.Fatal("expected items[1]=b")
	}
}

func TestLookupJSONPathArrayOutOfRange(t *testing.T) {
	data := map[string]any{
		"items": []any{"a"},
	}
	_, ok := LookupJSONPath(data, "items[5]")
	if ok {
		t.Fatal("expected out of range to fail")
	}
}

func TestLookupJSONPathNotArray(t *testing.T) {
	data := map[string]any{
		"items": "not-array",
	}
	_, ok := LookupJSONPath(data, "items[0]")
	if ok {
		t.Fatal("expected not array to fail")
	}
}

func TestLookupJSONPathNameAndIndex(t *testing.T) {
	data := map[string]any{
		"items": []any{map[string]any{"name": "first"}, map[string]any{"name": "second"}},
	}
	v, ok := LookupJSONPath(data, "items[1].name")
	if !ok || v != "second" {
		t.Fatal("expected items[1].name=second")
	}
}

func TestLookupJSONPathDollar(t *testing.T) {
	data := map[string]any{"name": "argus"}
	v, ok := LookupJSONPath(data, "$name")
	if !ok || v != "argus" {
		t.Fatal("expected $name=argus")
	}
}

func TestValidateJSONPath(t *testing.T) {
	r := ValidateJSONPath([]byte(`{"name":"argus"}`), "name", "argus", OpEqual)
	if !r.Success {
		t.Fatal("expected path match")
	}
}

func TestValidateJSONPathInvalidJSON(t *testing.T) {
	r := ValidateJSONPath([]byte(`{invalid}`), "name", "argus", OpEqual)
	if r.Success || r.Message == "" {
		t.Fatal("expected invalid JSON to fail")
	}
}

func TestValidateJSONPathNotFound(t *testing.T) {
	r := ValidateJSONPath([]byte(`{"name":"argus"}`), "missing", "argus", OpEqual)
	if r.Success {
		t.Fatal("expected path not found to fail")
	}
}

func TestValidateJSONPathExists(t *testing.T) {
	r := ValidateJSONPathExists([]byte(`{"name":"argus"}`), "name")
	if !r.Success {
		t.Fatal("expected path to exist")
	}
}

func TestValidateJSONPathExistsInvalidJSON(t *testing.T) {
	r := ValidateJSONPathExists([]byte(`{invalid}`), "name")
	if r.Success || r.Message == "" {
		t.Fatal("expected invalid JSON to fail")
	}
}

func TestValidateJSONPathExistsNotFound(t *testing.T) {
	r := ValidateJSONPathExists([]byte(`{"name":"argus"}`), "missing")
	if r.Success {
		t.Fatal("expected path not found")
	}
}

func TestParsePathPartNoIndex(t *testing.T) {
	name, indexes := parsePathPart("name")
	if name != "name" || len(indexes) != 0 {
		t.Fatal("expected name without indexes")
	}
}

func TestParsePathPartInvalidIndex(t *testing.T) {
	name, indexes := parsePathPart("items[abc]")
	if name != "items" || len(indexes) != 0 {
		t.Fatal("expected name with invalid index to have no indexes")
	}
}

func TestParsePathPartEmptyBracket(t *testing.T) {
	name, indexes := parsePathPart("items[]")
	if name != "items" || len(indexes) != 0 {
		t.Fatal("expected empty bracket to have no indexes")
	}
}

func TestParsePathPartSingleCharBracket(t *testing.T) {
	name, indexes := parsePathPart("items[a]")
	if name != "items" || len(indexes) != 0 {
		t.Fatal("expected single char bracket to have no indexes")
	}
}

func TestSplitJSONPathEmpty(t *testing.T) {
	parts := splitJSONPath("a..b")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
}
