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

// ==================== ConvertNumericStrings ====================

func TestConvertNumericStrings(t *testing.T) {
	type convertCase struct {
		name           string
		input          string
		nilInput       bool // 传入 nil 切片而非空字符串
		wantErr        bool
		wantEqual      string // 结果需完全等于此值
		wantNil        bool   // 结果需为 nil
		wantEmpty      bool   // 结果需为空切片（len==0 但可能非 nil）
		wantContains   []string
		wantNotContain []string
		wantNoNewline  bool // 结果末尾不得有换行
	}

	cases := []convertCase{
		// --- 对象内数字字符串转数字 ---
		{
			name:           "int64 string in object",
			input:          `{"items":[{"bank_bin_id":"1191749077254635521","enabled":false}]}`,
			wantContains:   []string{`"bank_bin_id":1191749077254635521`},
			wantNotContain: []string{`"bank_bin_id":"1191749077254635521"`},
		},
		{
			name:         "float string in object",
			input:        `{"price":"3.14","name":"test"}`,
			wantContains: []string{`"price":3.14`, `"name":"test"`},
		},
		{
			name:         "uint64 string in object",
			input:        `{"id":"18446744073709551615"}`,
			wantContains: []string{`"id":18446744073709551615`},
		},
		{
			name:         "max int64 string",
			input:        `{"id":"9223372036854775807"}`,
			wantContains: []string{`"id":9223372036854775807`},
		},
		{
			name:         "already number stays number",
			input:        `{"count":42,"price":3.14,"name":"hello"}`,
			wantContains: []string{`"count":42`, `"name":"hello"`},
		},
		{
			name:         "non-numeric strings remain",
			input:        `{"name":"hello","code":"abc123","flag":"true","desc":"ABC123"}`,
			wantContains: []string{`"name":"hello"`, `"code":"abc123"`},
		},
		{
			name:         "empty string remains",
			input:        `{"name":""}`,
			wantContains: []string{`"name":""`},
		},
		{
			name:         "boolean and null remain",
			input:        `{"active":true,"deleted":false,"extra":null,"ref":null}`,
			wantContains: []string{`"active":true`, `"deleted":false`, `"extra":null`},
		},
		{
			name:         "array numeric strings converted",
			input:        `{"ids":["1","2","3"],"names":["a","b"]}`,
			wantContains: []string{`"ids":[1,2,3]`, `"names":["a","b"]`},
		},
		{
			name:         "nested object numeric strings",
			input:        `{"user":{"id":"123","name":"test","profile":{"age":"30","score":"98.5"}},"scores":["98","5.5"]}`,
			wantContains: []string{`"id":123`, `"age":30`, `"score":98.5`, `"name":"test"`},
		},
		{
			name:         "negative numbers",
			input:        `{"temp":"-5.5","offset":"-100"}`,
			wantContains: []string{`"temp":-5.5`, `"offset":-100`},
		},
		{
			name:           "zero values",
			input:          `{"int":"0","float":"0.0","neg":"-1"}`,
			wantContains:   []string{`"int":0`, `"neg":-1`},
			wantNotContain: []string{`"0"`},
		},
		{
			name:         "mixed types in object",
			input:        `{"id":"123","active":true,"extra":null,"count":42,"price":3.14}`,
			wantContains: []string{`"id":123`, `"active":true`, `"extra":null`, `"count":42`},
		},
		{
			name:         "float overflow string remains",
			input:        `{"val":"1e9999"}`,
			wantContains: []string{`"val":"1e9999"`},
		},
		// --- 顶层值 ---
		{
			name:      "top-level int string",
			input:     `"1191749077254635521"`,
			wantEqual: `1191749077254635521`,
		},
		{
			name:      "top-level float string",
			input:     `"3.14"`,
			wantEqual: `3.14`,
		},
		{
			name:      "top-level non-numeric string remains",
			input:     `"hello"`,
			wantEqual: `"hello"`,
		},
		{
			name:      "top-level number remains",
			input:     `12345`,
			wantEqual: `12345`,
		},
		{
			name:      "top-level bool remains",
			input:     `true`,
			wantEqual: `true`,
		},
		// --- 边界与错误场景 ---
		{
			name:    "invalid JSON returns error",
			input:   `{"id":"123"`,
			wantErr: true,
		},
		{
			name:      "empty input returns empty",
			input:     ``,
			wantEmpty: true,
		},
		{
			name:     "nil input returns nil",
			nilInput: true,
			wantNil:  true,
		},
		{
			name:          "trailing newline trimmed",
			input:         `{"id":"123"}`,
			wantNoNewline: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			if tt.nilInput {
				data = nil
			} else {
				data = []byte(tt.input)
			}
			converted, err := ConvertNumericStrings(data)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantNil {
				if converted != nil {
					t.Errorf("expected nil, got %v", converted)
				}
				return
			}
			if tt.wantEmpty {
				if len(converted) != 0 {
					t.Errorf("expected empty result, got: %s", converted)
				}
				return
			}
			s := string(converted)
			for _, sub := range tt.wantContains {
				if !strings.Contains(s, sub) {
					t.Errorf("expected result to contain %q, got: %s", sub, s)
				}
			}
			for _, sub := range tt.wantNotContain {
				if strings.Contains(s, sub) {
					t.Errorf("expected result to NOT contain %q, got: %s", sub, s)
				}
			}
			if tt.wantEqual != "" && s != tt.wantEqual {
				t.Errorf("expected %q, got: %s", tt.wantEqual, s)
			}
			if tt.wantNoNewline && len(converted) > 0 && converted[len(converted)-1] == '\n' {
				t.Errorf("expected trailing newline to be trimmed, got %q", converted)
			}
		})
	}
}

// TestConvertNumericStrings_DecodesToNumber 验证转换结果解码后为 json.Number 类型
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

// TestConvertNumericStrings_ValidJSON 验证转换结果为合法 JSON
func TestConvertNumericStrings_ValidJSON(t *testing.T) {
	input := `{"items":[{"bank_bin_id":"1191749077254635521","enabled":false}]}`
	converted, err := ConvertNumericStrings([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var v interface{}
	dec := json.NewDecoder(bytes.NewReader(converted))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("converted result should be valid JSON: %v", err)
	}
}

// ==================== ContainsQuotedNumber ====================

func TestContainsQuotedNumber(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// --- 基础场景 ---
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
		// --- 边界场景（覆盖各分支末尾条件） ---
		{"leading spaces then lone quote", `   "`, false},
		{"string ending with backslash", `{"a":"x\`, false},
		{"array bracket then quote at end", `["`, false},
		{"colon then quote at end", `:"`, false},
		{"comma then quote at end", `,"`, false},
		{"escaped quote then numeric string", `{"a":"b\"c","d":"123"}`, true},
		{"string with internal bracket", `{"desc":"a[b"}`, false},
		{"deeply nested numeric string", `{"a":{"b":{"c":"42"}}}`, true},
		{"array starts with numeric string", `["1"]`, true},
		{"array starts with non-numeric string", `["a"]`, false},
		{"top-level number", `123`, false},
		{"top-level bool", `true`, false},
		{"top-level null", `null`, false},
		{"only whitespace", `   `, false},
		{"single byte", `a`, false},
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

// ==================== IsNumStart ====================

func TestIsNumStart(t *testing.T) {
	tests := []struct {
		name string
		c    byte
		want bool
	}{
		{"digit 0", '0', true},
		{"digit 1", '1', true},
		{"digit 5", '5', true},
		{"digit 9", '9', true},
		{"minus sign", '-', true},
		{"plus sign", '+', true},
		{"slash before 0", '/', false},
		{"colon after 9", ':', false},
		{"lowercase letter", 'a', false},
		{"uppercase letter", 'A', false},
		{"letter z", 'z', false},
		{"quote", '"', false},
		{"space", ' ', false},
		{"colon", ':', false},
		{"dot", '.', false},
		{"open brace", '{', false},
		{"open bracket", '[', false},
		{"close brace", '}', false},
		{"close bracket", ']', false},
		{"null byte", 0, false},
		{"255 byte", 255, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumStart(tt.c); got != tt.want {
				t.Errorf("IsNumStart(%q) = %v, want %v", tt.c, got, tt.want)
			}
		})
	}
}

// ==================== TryStringToNumber ====================

func TestTryStringToNumber(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  interface{}
	}{
		// 空字符串 -> 原样返回（覆盖 len(s)==0 分支）
		{"empty string", "", ""},
		// int64 范围
		{"int64 positive", "123", json.Number("123")},
		{"int64 negative", "-456", json.Number("-456")},
		{"int64 zero", "0", json.Number("0")},
		{"int64 max", "9223372036854775807", json.Number("9223372036854775807")},
		{"int64 min", "-9223372036854775808", json.Number("-9223372036854775808")},
		// uint64 范围（超出 int64，需走 ParseUint 分支）
		{"uint64 max", "18446744073709551615", json.Number("18446744073709551615")},
		{"uint64 over int64", "9223372036854775808", json.Number("9223372036854775808")},
		// float64 范围（ParseInt/Uint 失败，ParseFloat 成功）
		{"float", "3.14", json.Number("3.14")},
		{"float negative", "-2.5", json.Number("-2.5")},
		{"float scientific", "1e10", json.Number("1e10")},
		{"float scientific negative", "-1.5e-3", json.Number("-1.5e-3")},
		{"float zero", "0.0", json.Number("0.0")},
		// 纯非数字字符串 -> 原样返回（覆盖 return s 末尾分支）
		{"non-numeric", "hello", "hello"},
		{"alphanumeric", "abc123", "abc123"},
		{"starts with letter", "a123", "a123"},
		// 指数过大 -> ParseFloat 返回 +Inf 与 ErrRange，err != nil -> 原样返回
		{"float overflow", "1e9999", "1e9999"},
		// 多个符号开头非数字 -> 原样返回
		{"plus sign then letters", "+abc", "+abc"},
		{"dash then letters", "-abc", "-abc"},
		// 带小数点：strconv.ParseFloat 接受 "1." 为 1.0
		{"dot only", ".", "."},
		{"trailing dot parses as float", "1.", json.Number("1.")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TryStringToNumber(tt.input)
			switch wantv := tt.want.(type) {
			case json.Number:
				gotNum, ok := got.(json.Number)
				if !ok {
					t.Errorf("TryStringToNumber(%q) type = %T, want json.Number", tt.input, got)
					return
				}
				if string(gotNum) != string(wantv) {
					t.Errorf("TryStringToNumber(%q) = %s, want %s", tt.input, gotNum, wantv)
				}
			case string:
				gotStr, ok := got.(string)
				if !ok {
					t.Errorf("TryStringToNumber(%q) type = %T, want string", tt.input, got)
					return
				}
				if gotStr != wantv {
					t.Errorf("TryStringToNumber(%q) = %q, want %q", tt.input, gotStr, wantv)
				}
			}
		})
	}
}

// ==================== ConvertNumericStringsRecursive ====================

func TestConvertNumericStringsRecursive(t *testing.T) {
	t.Run("nil value covers default case", func(t *testing.T) {
		got := ConvertNumericStringsRecursive(nil)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("bool value covers default case", func(t *testing.T) {
		got := ConvertNumericStringsRecursive(true)
		if got != true {
			t.Errorf("expected true, got %v", got)
		}
	})
	t.Run("json.Number value covers default case", func(t *testing.T) {
		n := json.Number("42")
		got := ConvertNumericStringsRecursive(n)
		gotNum, ok := got.(json.Number)
		if !ok || string(gotNum) != "42" {
			t.Errorf("expected json.Number 42, got %v", got)
		}
	})
	t.Run("int value covers default case", func(t *testing.T) {
		// 非标准 JSON 类型也应走 default 分支
		got := ConvertNumericStringsRecursive(42)
		if got != 42 {
			t.Errorf("expected 42, got %v", got)
		}
	})
	t.Run("nested map with mixed types", func(t *testing.T) {
		v := map[string]interface{}{
			"id":     "123",
			"name":   "test",
			"active": true,
			"extra":  nil,
			"count":  json.Number("42"),
		}
		got := ConvertNumericStringsRecursive(v)
		m, ok := got.(map[string]interface{})
		if !ok {
			t.Fatalf("expected map, got %T", got)
		}
		if _, ok := m["id"].(json.Number); !ok {
			t.Errorf("expected id to be json.Number, got %T", m["id"])
		}
		if m["name"] != "test" {
			t.Errorf("expected name to remain string, got %v", m["name"])
		}
		if m["active"] != true {
			t.Errorf("expected active to remain bool, got %v", m["active"])
		}
		if m["extra"] != nil {
			t.Errorf("expected extra to remain nil, got %v", m["extra"])
		}
	})
	t.Run("array with mixed types", func(t *testing.T) {
		v := []interface{}{"123", "hello", true, nil, json.Number("9")}
		got := ConvertNumericStringsRecursive(v)
		arr, ok := got.([]interface{})
		if !ok {
			t.Fatalf("expected slice, got %T", got)
		}
		if _, ok := arr[0].(json.Number); !ok {
			t.Errorf("expected arr[0] to be json.Number, got %T", arr[0])
		}
		if arr[1] != "hello" {
			t.Errorf("expected arr[1] to remain string, got %v", arr[1])
		}
		if arr[2] != true {
			t.Errorf("expected arr[2] to remain bool, got %v", arr[2])
		}
	})
	t.Run("string returns via TryStringToNumber", func(t *testing.T) {
		got := ConvertNumericStringsRecursive("100")
		if _, ok := got.(json.Number); !ok {
			t.Errorf("expected json.Number, got %T", got)
		}
	})
}
