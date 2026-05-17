/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-17 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 11:54:16
 * @FilePath: \go-argus\reexport_test.go
 * @Description: reexport_test.go 测试，覆盖所有导出函数
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"
)

func TestValidateRegex(t *testing.T) {
	result := ValidateRegex([]byte("hello world"), `^hello`)
	if !result.Success {
		t.Fatal("expected regex to match")
	}
	result = ValidateRegex([]byte("hello world"), `^world`)
	if result.Success {
		t.Fatal("expected regex not to match")
	}
}

func TestValidateEmail(t *testing.T) {
	result := ValidateEmail("user@example.com")
	if !result.Success {
		t.Fatal("expected valid email")
	}
	result = ValidateEmail("not-email")
	if result.Success {
		t.Fatal("expected invalid email")
	}
}

func TestValidateIPAddress(t *testing.T) {
	result := ValidateIPAddress("192.168.1.1")
	if !result.Success {
		t.Fatal("expected valid IP")
	}
	result = ValidateIPAddress("not-an-ip")
	if result.Success {
		t.Fatal("expected invalid IP")
	}
}

func TestValidateProtocol(t *testing.T) {
	result := ValidateProtocol("https://example.com", "http", "https")
	if !result.Success {
		t.Fatal("expected valid protocol")
	}
	result = ValidateProtocol("ftp://example.com", "http", "https")
	if result.Success {
		t.Fatal("expected invalid protocol")
	}
}

func TestValidateHTTP(t *testing.T) {
	result := ValidateHTTP("https://example.com")
	if !result.Success {
		t.Fatal("expected valid HTTP URL")
	}
	result = ValidateHTTP("ftp://example.com")
	if result.Success {
		t.Fatal("expected invalid HTTP URL")
	}
}

func TestValidateWebSocket(t *testing.T) {
	result := ValidateWebSocket("wss://example.com/ws")
	if !result.Success {
		t.Fatal("expected valid WebSocket URL")
	}
	result = ValidateWebSocket("http://example.com")
	if result.Success {
		t.Fatal("expected invalid WebSocket URL")
	}
}

func TestValidateUUID(t *testing.T) {
	result := ValidateUUID("550e8400-e29b-41d4-a716-446655440000")
	if !result.Success {
		t.Fatal("expected valid UUID")
	}
	result = ValidateUUID("not-a-uuid")
	if result.Success {
		t.Fatal("expected invalid UUID")
	}
}

func TestValidateBase64(t *testing.T) {
	result := ValidateBase64("YXJndXM=")
	if !result.Success {
		t.Fatal("expected valid Base64")
	}
	result = ValidateBase64("!!!invalid!!!")
	if result.Success {
		t.Fatal("expected invalid Base64")
	}
}

func TestIsEmail(t *testing.T) {
	if !IsEmail("user@example.com") {
		t.Fatal("expected IsEmail to return true")
	}
	if IsEmail("not-email") {
		t.Fatal("expected IsEmail to return false")
	}
}

func TestIsIP(t *testing.T) {
	if !IsIP("192.168.1.1") {
		t.Fatal("expected IsIP to return true")
	}
	if IsIP("not-an-ip") {
		t.Fatal("expected IsIP to return false")
	}
}

func TestIsUUID(t *testing.T) {
	if !IsUUID("550e8400-e29b-41d4-a716-446655440000") {
		t.Fatal("expected IsUUID to return true")
	}
	if IsUUID("not-a-uuid") {
		t.Fatal("expected IsUUID to return false")
	}
}

func TestIsBase64(t *testing.T) {
	if !IsBase64("YXJndXM=") {
		t.Fatal("expected IsBase64 to return true")
	}
	if IsBase64("!!!invalid!!!") {
		t.Fatal("expected IsBase64 to return false")
	}
}

func TestValidateJSON(t *testing.T) {
	if err := ValidateJSON([]byte(`{"key":"value"}`)); err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}
	if err := ValidateJSON([]byte(`{invalid}`)); err == nil {
		t.Fatal("expected invalid JSON to fail")
	}
}

func TestIsJSONNull(t *testing.T) {
	if !IsJSONNull([]byte("null")) {
		t.Fatal("expected null to be JSON null")
	}
	if IsJSONNull([]byte(`{"key":"value"}`)) {
		t.Fatal("expected object not to be JSON null")
	}
}

func TestIsJSONColumnType(t *testing.T) {
	if !IsJSONColumnType("json") {
		t.Fatal("expected 'json' to be JSON column type")
	}
	if !IsJSONColumnType("jsonb") {
		t.Fatal("expected 'jsonb' to be JSON column type")
	}
	if IsJSONColumnType("varchar") {
		t.Fatal("expected 'varchar' not to be JSON column type")
	}
}

func TestValidateJSONWithData(t *testing.T) {
	data, err := ValidateJSONWithData([]byte(`{"key":"value"}`))
	if err != nil {
		t.Fatalf("expected valid JSON: %v", err)
	}
	if data == nil {
		t.Fatal("expected non-nil data")
	}
	_, err = ValidateJSONWithData([]byte(`{invalid}`))
	if err == nil {
		t.Fatal("expected invalid JSON to fail")
	}
}

func TestValidateJSONField(t *testing.T) {
	result := ValidateJSONField([]byte(`{"name":"argus"}`), "name", "argus")
	if !result.Success {
		t.Fatal("expected JSON field to match")
	}
}

func TestValidateJSONFields(t *testing.T) {
	results := ValidateJSONFields([]byte(`{"name":"argus","version":"1"}`), map[string]any{
		"name": "argus",
	})
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
}

func TestLookupJSONPath(t *testing.T) {
	data := map[string]any{"key": "value"}
	val, ok := LookupJSONPath(data, "key")
	if !ok || val != "value" {
		t.Fatal("expected to find key")
	}
	_, ok = LookupJSONPath(data, "missing")
	if ok {
		t.Fatal("expected not to find missing key")
	}
}

func TestValidateJSONPath(t *testing.T) {
	result := ValidateJSONPath([]byte(`{"name":"argus"}`), "name", "argus", OpEqual)
	if !result.Success {
		t.Fatal("expected JSON path to match")
	}
}

func TestValidateJSONPathExists(t *testing.T) {
	result := ValidateJSONPathExists([]byte(`{"name":"argus"}`), "name")
	if !result.Success {
		t.Fatal("expected JSON path to exist")
	}
	result = ValidateJSONPathExists([]byte(`{"name":"argus"}`), "missing")
	if result.Success {
		t.Fatal("expected missing JSON path not to exist")
	}
}

func TestNewEnumValidator(t *testing.T) {
	ev := NewEnumValidator("admin", "member", "guest")
	if ev == nil {
		t.Fatal("expected non-nil enum validator")
	}
}

func TestSkipJSONSpaces(t *testing.T) {
	data := []byte("  hello")
	idx := SkipJSONSpaces(data, 0)
	if idx != 2 {
		t.Fatalf("expected index 2, got %d", idx)
	}
}

func TestScanJSONString(t *testing.T) {
	data := []byte(`"hello"world`)
	end, err := ScanJSONString(data, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if end != 7 {
		t.Fatalf("expected end at 7, got %d", end)
	}
}

func TestScanJSONValueEnd(t *testing.T) {
	data := []byte(`{"key":"value"}rest`)
	end, err := ScanJSONValueEnd(data, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if end >= len(data) {
		t.Fatal("expected end before data length")
	}
}

func TestGetCompiledRegex(t *testing.T) {
	re, err := GetCompiledRegex(`^\d+$`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !re.MatchString("123") {
		t.Fatal("expected regex to match digits")
	}
}

func TestClearRegexCache(t *testing.T) {
	ClearRegexCache()
	_, _ = GetCompiledRegex(`^test$`)
	ClearRegexCache()
}

func TestValidateIPCompat(t *testing.T) {
	result := ValidateIP("192.168.1.1")
	if !result.Success {
		t.Fatal("expected valid IP")
	}
}

func TestValidateJSONSchema(t *testing.T) {
	schema := QuickSchema(map[string]string{
		"name": "string",
	}, "name")
	result := ValidateJSONSchema(map[string]any{"name": "argus"}, schema)
	if !result.Success {
		t.Fatal("expected schema validation to pass")
	}
}

func TestValidateStructWithSchema(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}
	schema := QuickSchema(map[string]string{
		"name": "string",
	}, "name")
	result := ValidateStructWithSchema(TestStruct{Name: "argus"}, schema)
	if !result.Success {
		t.Fatal("expected struct schema validation to pass")
	}
}

func TestNewSchemaBuilder(t *testing.T) {
	builder := NewSchemaBuilder()
	if builder == nil {
		t.Fatal("expected non-nil schema builder")
	}
}

func TestFormatSchemaError(t *testing.T) {
	schema := QuickSchema(map[string]string{
		"name": "string",
	}, "name")
	result := ValidateJSONSchema(map[string]any{}, schema)
	if result.Success {
		t.Fatal("expected schema validation to fail for missing required field")
	}
	msg := FormatSchemaError(result)
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestGetCompiledRegexInvalid(t *testing.T) {
	_, err := GetCompiledRegex(`[invalid`)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestQuickSchema(t *testing.T) {
	s := QuickSchema(map[string]string{"name": "string"}, "name")
	if s.Type != "object" {
		t.Fatal("expected object type")
	}
}
