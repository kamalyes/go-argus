/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-06 23:16:17
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-07 00:16:05
 * @FilePath: \go-argus\validate\json_numeric_test.go
 * @Description: json_numeric.go 测试，覆盖数字字符串转换和快速预检功能
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// ==================== ConvertNumericStrings 测试 ====================

func TestConvertNumericStrings_Int64String(t *testing.T) {
	// 真实 bug 场景：前端传入的大 int64 ID
	input := `{"items":[{"bank_bin_id":"1191749077254635521","enabled":false}]}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(converted)
	if !strings.Contains(s, `"bank_bin_id":1191749077254635521`) {
		t.Errorf("expected int64 string to be converted, got: %s", s)
	}
	if strings.Contains(s, `"bank_bin_id":"1191749077254635521"`) {
		t.Errorf("int64 string should not remain as string, got: %s", s)
	}
	// 验证转换结果是合法 JSON
	var v interface{}
	dec := json.NewDecoder(bytes.NewReader(converted))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("converted result should be valid JSON: %v", err)
	}
}

func TestConvertNumericStrings_FloatString(t *testing.T) {
	input := `{"price":"3.14","name":"test"}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(converted)
	if !strings.Contains(s, `"price":3.14`) {
		t.Errorf("expected float string to be converted, got: %s", s)
	}
	if !strings.Contains(s, `"name":"test"`) {
		t.Errorf("non-numeric string should remain, got: %s", s)
	}
}

func TestConvertNumericStrings_Uint64String(t *testing.T) {
	input := `{"id":"18446744073709551615"}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(converted), `"id":18446744073709551615`) {
		t.Errorf("expected uint64 conversion, got: %s", converted)
	}
}

func TestConvertNumericStrings_LargeInt64(t *testing.T) {
	// 最大 int64 值
	input := `{"id":"9223372036854775807"}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(converted), `"id":9223372036854775807`) {
		t.Errorf("expected max int64 conversion, got: %s", converted)
	}
}

func TestConvertNumericStrings_AlreadyNumber(t *testing.T) {
	input := `{"count":42,"price":3.14,"name":"hello"}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(converted)
	if !strings.Contains(s, `"count":42`) {
		t.Errorf("number field should remain, got: %s", s)
	}
	if !strings.Contains(s, `"name":"hello"`) {
		t.Errorf("string field should remain, got: %s", s)
	}
}

func TestConvertNumericStrings_NonNumericString(t *testing.T) {
	input := `{"name":"hello","code":"abc123","flag":"true","desc":"ABC123"}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(converted)
	if !strings.Contains(s, `"name":"hello"`) {
		t.Errorf("non-numeric string should remain, got: %s", s)
	}
	if !strings.Contains(s, `"code":"abc123"`) {
		t.Errorf("alphanumeric string should remain, got: %s", s)
	}
}

func TestConvertNumericStrings_EmptyString(t *testing.T) {
	input := `{"name":""}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(converted), `"name":""`) {
		t.Errorf("empty string should remain, got: %s", converted)
	}
}

func TestConvertNumericStrings_BooleanAndNull(t *testing.T) {
	input := `{"active":true,"deleted":false,"extra":null,"ref":null}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(converted)
	if !strings.Contains(s, `"active":true`) {
		t.Errorf("true should remain, got: %s", s)
	}
	if !strings.Contains(s, `"deleted":false`) {
		t.Errorf("false should remain, got: %s", s)
	}
	if !strings.Contains(s, `"extra":null`) {
		t.Errorf("null should remain, got: %s", s)
	}
}

func TestConvertNumericStrings_Array(t *testing.T) {
	input := `{"ids":["1","2","3"],"names":["a","b"]}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(converted)
	if !strings.Contains(s, `"ids":[1,2,3]`) {
		t.Errorf("array numeric strings should be converted, got: %s", s)
	}
	if !strings.Contains(s, `"names":["a","b"]`) {
		t.Errorf("non-numeric array elements should remain, got: %s", s)
	}
}

func TestConvertNumericStrings_Nested(t *testing.T) {
	input := `{"user":{"id":"123","name":"test","profile":{"age":"30","score":"98.5"}},"scores":["98","5.5"]}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(converted)
	if !strings.Contains(s, `"id":123`) {
		t.Errorf("nested id should be converted, got: %s", s)
	}
	if !strings.Contains(s, `"age":30`) {
		t.Errorf("nested age should be converted, got: %s", s)
	}
	if !strings.Contains(s, `"score":98.5`) {
		t.Errorf("nested score should be converted, got: %s", s)
	}
	if !strings.Contains(s, `"name":"test"`) {
		t.Errorf("string in nested object should remain, got: %s", s)
	}
}

func TestConvertNumericStrings_InvalidJSON(t *testing.T) {
	// 包含数字字符串但语法错误的 JSON 应该返回错误
	input := `{"id":"123"`
	_, err := ConvertNumericStrings([]byte(input))
	if err == nil {
		t.Fatal("expected error for invalid JSON with numeric strings")
	}
}

func TestConvertNumericStrings_EmptyInput(t *testing.T) {
	converted, err := ConvertNumericStrings([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(converted) != 0 {
		t.Errorf("expected empty result, got: %s", converted)
	}
}

func TestConvertNumericStrings_NegativeNumbers(t *testing.T) {
	input := `{"temp":"-5.5","offset":"-100"}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(converted)
	if !strings.Contains(s, `"temp":-5.5`) {
		t.Errorf("expected negative float, got: %s", s)
	}
	if !strings.Contains(s, `"offset":-100`) {
		t.Errorf("expected negative int, got: %s", s)
	}
}

func TestConvertNumericStrings_ZeroValues(t *testing.T) {
	input := `{"int":"0","float":"0.0","neg":"-1"}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(converted)
	if !strings.Contains(s, `"int":0`) {
		t.Errorf("zero int should be converted, got: %s", s)
	}
	if strings.Contains(s, `"0"`) {
		t.Errorf("zero strings should not remain as string, got: %s", s)
	}
	if !strings.Contains(s, `"neg":-1`) {
		t.Errorf("negative one should be converted, got: %s", s)
	}
}

func TestConvertNumericStrings_DecodesToNumber(t *testing.T) {
	input := `{"id":"1191749077254635521","name":"test"}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var v map[string]interface{}
	dec := json.NewDecoder(strings.NewReader(string(converted)))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("converted result should be valid JSON: %v", err)
	}
	if _, ok := v["id"].(json.Number); !ok {
		t.Errorf("expected id to be json.Number, got %T: %v", v["id"], v["id"])
	}
	if _, ok := v["name"].(string); !ok {
		t.Errorf("expected name to be string, got %T: %v", v["name"], v["name"])
	}
}

// ==================== 顶层字符串（wrapperspb 场景） ====================

func TestConvertNumericStrings_PlainIntString(t *testing.T) {
	// 顶层数字字符串（wrapperspb.Int64Value 等场景）
	input := `"1191749077254635521"`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(converted) != `1191749077254635521` {
		t.Errorf("expected top-level int string to be converted, got: %s", converted)
	}
}

func TestConvertNumericStrings_PlainFloatString(t *testing.T) {
	input := `"3.14"`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(converted) != `3.14` {
		t.Errorf("expected top-level float string to be converted, got: %s", converted)
	}
}

func TestConvertNumericStrings_PlainNonNumericString(t *testing.T) {
	input := `"hello"`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(converted) != `"hello"` {
		t.Errorf("top-level non-numeric string should remain, got: %s", converted)
	}
}

// ==================== ContainsQuotedNumber 预检函数测试 ====================

func TestContainsQuotedNumber(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"object with int64 string", `{"bank_bin_id":"1191749077254635521"}`, true},
		{"object with all numbers", `{"count":42,"price":3.14}`, false},
		{"object with string value", `{"name":"hello"}`, false},
		{"array with numeric strings", `["1","2","3"]`, true},
		{"array with mixed values", `[1,"hello",3]`, false},
		{"array starting with numeric string", `["1",2,3]`, true},
		{"nested object", `{"user":{"id":"123"}}`, true},
		{"top-level numeric string", `"1191749077254635521"`, true},
		{"top-level non-numeric string", `"hello"`, false},
		{"negative numeric string", `{"id":"-456"}`, true},
		{"plus sign numeric string", `{"id":"+789"}`, true},
		{"float string", `{"price":"3.14"}`, true},
		{"null and bool", `{"a":true,"b":null}`, false},
		{"string with colon inside", `{"desc":"key:value"}`, false},
		{"string with comma inside", `{"desc":"a,b,c"}`, false},
		{"empty object", `{}`, false},
		{"empty array", `[]`, false},
		{"whitespace before value", `{ "id" : "123" }`, true},
		{"no spaces, multiple fields", `{"a":"1","b":"hello","c":"2"}`, true},
		{"empty input", ``, false},
		{"single char", `"`, false},
		{"scientific notation", `{"val":"1e10"}`, true},
		{"zero string", `{"val":"0"}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsQuotedNumber([]byte(tt.input))
			if got != tt.want {
				t.Errorf("ContainsQuotedNumber(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// ==================== IsNumStart 测试 ====================

func TestIsNumStart(t *testing.T) {
	tests := []struct {
		c    byte
		want bool
	}{
		{'0', true}, {'1', true}, {'5', true}, {'9', true},
		{'-', true}, {'+', true},
		{'a', false}, {'A', false}, {'z', false},
		{'"', false}, {' ', false}, {':', false}, {'.', false},
		{'{', false}, {'[', false}, {'}', false}, {']', false},
	}
	for _, tt := range tests {
		got := IsNumStart(tt.c)
		if got != tt.want {
			t.Errorf("IsNumStart(%q) = %v, want %v", tt.c, got, tt.want)
		}
	}
}
