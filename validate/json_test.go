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

	"github.com/kamalyes/go-argus/constants"
)

// ==================== ValidateJSON ====================

func TestValidateJSON(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		if err := ValidateJSON([]byte(`{"key":"value"}`)); err != nil {
			t.Fatal("expected valid JSON")
		}
	})
	t.Run("invalid", func(t *testing.T) {
		if err := ValidateJSON([]byte(`{invalid}`)); err == nil {
			t.Fatal("expected invalid JSON to fail")
		}
	})
}

// ==================== IsJSONNull ====================

func TestIsJSONNull(t *testing.T) {
	if !IsJSONNull([]byte("null")) {
		t.Fatal("expected null")
	}
	if IsJSONNull([]byte(`{"a":1}`)) {
		t.Fatal("expected not null")
	}
}

// ==================== IsJSONColumnType ====================

func TestIsJSONColumnType(t *testing.T) {
	tests := []struct {
		name   string
		dbType string
		want   bool
	}{
		{"plain json", "json", true},
		{"plain jsonb", "jsonb", true},
		{"uppercase JSON", "JSON", true},
		{"uppercase JSONB", "JSONB", true},
		{"json with paren", "json(100)", true},
		{"jsonb with paren", "jsonb(200)", true},
		{"JSON with paren uppercase", "JSON(100)", true},
		{"json with spaces around paren", "  json(100)  ", true},
		{"text type", "text", false},
		{"varchar with paren", "varchar(255)", false},
		{"TEXT with paren uppercase", "TEXT(50)", false},
		{"empty string", "", false},
		{"integer type", "integer", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSONColumnType(tt.dbType); got != tt.want {
				t.Errorf("IsJSONColumnType(%q) = %v, want %v", tt.dbType, got, tt.want)
			}
		})
	}
}

// ==================== ValidateJSONWithData ====================

func TestValidateJSONWithData(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		data, err := ValidateJSONWithData([]byte(`{"key":"value"}`))
		if err != nil {
			t.Fatal("expected valid JSON")
		}
		m, ok := data.(map[string]interface{})
		if !ok || m["key"] != "value" {
			t.Fatal("expected key=value")
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ValidateJSONWithData([]byte(`{invalid}`))
		if err == nil {
			t.Fatal("expected invalid JSON to fail")
		}
	})
}

// ==================== ValidateJSONField ====================

func TestValidateJSONField(t *testing.T) {
	t.Run("field match", func(t *testing.T) {
		r := ValidateJSONField([]byte(`{"name":"argus"}`), "name", "argus")
		if !r.Success {
			t.Fatal("expected field match")
		}
	})
	t.Run("invalid JSON", func(t *testing.T) {
		r := ValidateJSONField([]byte(`{invalid}`), "name", "argus")
		if r.Success || r.Message == "" {
			t.Fatal("expected invalid JSON to fail")
		}
	})
	t.Run("root not object", func(t *testing.T) {
		r := ValidateJSONField([]byte(`[1,2,3]`), "name", "argus")
		if r.Success || r.Message == "" {
			t.Fatal("expected root not object to fail")
		}
	})
	t.Run("field not found", func(t *testing.T) {
		r := ValidateJSONField([]byte(`{"name":"argus"}`), "age", 25)
		if r.Success || r.Message == "" {
			t.Fatal("expected field not found to fail")
		}
	})
	t.Run("value mismatch", func(t *testing.T) {
		r := ValidateJSONField([]byte(`{"name":"argus"}`), "name", "other")
		if r.Success {
			t.Fatal("expected value mismatch")
		}
	})
}

// ==================== ValidateJSONFields ====================

func TestValidateJSONFields(t *testing.T) {
	results := ValidateJSONFields([]byte(`{"name":"argus","ver":1}`), map[string]any{
		"name": "argus",
		"ver":  1,
	})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

// ==================== LookupJSONPath ====================

func TestLookupJSONPath(t *testing.T) {
	t.Run("nested map", func(t *testing.T) {
		data := map[string]any{
			"user": map[string]any{
				"name": "argus",
			},
		}
		v, ok := LookupJSONPath(data, "$.user.name")
		if !ok || v != "argus" {
			t.Fatal("expected to find user.name=argus")
		}
	})
	t.Run("empty path returns data", func(t *testing.T) {
		data := "hello"
		v, ok := LookupJSONPath(data, "")
		if !ok || v != "hello" {
			t.Fatal("expected empty path to return data")
		}
	})
	t.Run("not found", func(t *testing.T) {
		data := map[string]any{"name": "argus"}
		_, ok := LookupJSONPath(data, "missing")
		if ok {
			t.Fatal("expected not found")
		}
	})
	t.Run("not map", func(t *testing.T) {
		data := "hello"
		_, ok := LookupJSONPath(data, "name")
		if ok {
			t.Fatal("expected not map to fail")
		}
	})
	t.Run("array index", func(t *testing.T) {
		data := map[string]any{
			"items": []any{"a", "b", "c"},
		}
		v, ok := LookupJSONPath(data, "items[1]")
		if !ok || v != "b" {
			t.Fatal("expected items[1]=b")
		}
	})
	t.Run("array out of range", func(t *testing.T) {
		data := map[string]any{
			"items": []any{"a"},
		}
		_, ok := LookupJSONPath(data, "items[5]")
		if ok {
			t.Fatal("expected out of range to fail")
		}
	})
	t.Run("not array", func(t *testing.T) {
		data := map[string]any{
			"items": "not-array",
		}
		_, ok := LookupJSONPath(data, "items[0]")
		if ok {
			t.Fatal("expected not array to fail")
		}
	})
	t.Run("name and index", func(t *testing.T) {
		data := map[string]any{
			"items": []any{map[string]any{"name": "first"}, map[string]any{"name": "second"}},
		}
		v, ok := LookupJSONPath(data, "items[1].name")
		if !ok || v != "second" {
			t.Fatal("expected items[1].name=second")
		}
	})
	t.Run("dollar prefix", func(t *testing.T) {
		data := map[string]any{"name": "argus"}
		v, ok := LookupJSONPath(data, "$name")
		if !ok || v != "argus" {
			t.Fatal("expected $name=argus")
		}
	})
}

// ==================== ValidateJSONPath ====================

func TestValidateJSONPath(t *testing.T) {
	t.Run("path match", func(t *testing.T) {
		r := ValidateJSONPath([]byte(`{"name":"argus"}`), "name", "argus", constants.OpEqual)
		if !r.Success {
			t.Fatal("expected path match")
		}
	})
	t.Run("invalid JSON", func(t *testing.T) {
		r := ValidateJSONPath([]byte(`{invalid}`), "name", "argus", constants.OpEqual)
		if r.Success || r.Message == "" {
			t.Fatal("expected invalid JSON to fail")
		}
	})
	t.Run("not found", func(t *testing.T) {
		r := ValidateJSONPath([]byte(`{"name":"argus"}`), "missing", "argus", constants.OpEqual)
		if r.Success {
			t.Fatal("expected path not found to fail")
		}
	})
}

// ==================== ValidateJSONPathExists ====================

func TestValidateJSONPathExists(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		r := ValidateJSONPathExists([]byte(`{"name":"argus"}`), "name")
		if !r.Success {
			t.Fatal("expected path to exist")
		}
	})
	t.Run("invalid JSON", func(t *testing.T) {
		r := ValidateJSONPathExists([]byte(`{invalid}`), "name")
		if r.Success || r.Message == "" {
			t.Fatal("expected invalid JSON to fail")
		}
	})
	t.Run("not found", func(t *testing.T) {
		r := ValidateJSONPathExists([]byte(`{"name":"argus"}`), "missing")
		if r.Success {
			t.Fatal("expected path not found")
		}
	})
}

// ==================== parsePathPart ====================

func TestParsePathPart(t *testing.T) {
	t.Run("no index", func(t *testing.T) {
		name, indexes := parsePathPart("name")
		if name != "name" || len(indexes) != 0 {
			t.Fatal("expected name without indexes")
		}
	})
	t.Run("invalid index", func(t *testing.T) {
		name, indexes := parsePathPart("items[abc]")
		if name != "items" || len(indexes) != 0 {
			t.Fatal("expected name with invalid index to have no indexes")
		}
	})
	t.Run("empty bracket", func(t *testing.T) {
		name, indexes := parsePathPart("items[]")
		if name != "items" || len(indexes) != 0 {
			t.Fatal("expected empty bracket to have no indexes")
		}
	})
	t.Run("single char bracket", func(t *testing.T) {
		name, indexes := parsePathPart("items[a]")
		if name != "items" || len(indexes) != 0 {
			t.Fatal("expected single char bracket to have no indexes")
		}
	})
}

// ==================== splitJSONPath ====================

func TestSplitJSONPath(t *testing.T) {
	t.Run("empty parts collapsed", func(t *testing.T) {
		parts := splitJSONPath("a..b")
		if len(parts) != 2 {
			t.Fatalf("expected 2 parts, got %d", len(parts))
		}
	})
}

// ==================== JSON Scanner Helpers ====================

func TestJSONScannerHelpers(t *testing.T) {
	data := []byte(" \t{\"name\":\"argus\",\"items\":[1,true,null]}\n")
	start := SkipJSONSpaces(data, 0)
	if start != 2 {
		t.Fatalf("expected first non-space at 2, got %d", start)
	}
	end, err := ScanJSONValueEnd(data, start)
	if err != nil {
		t.Fatalf("expected object scan to pass: %v", err)
	}
	if end <= start {
		t.Fatalf("expected scan to advance, got %d", end)
	}
	strEnd, err := ScanJSONString([]byte(`"a\"b"`), 0)
	if err != nil || strEnd != len(`"a\"b"`) {
		t.Fatalf("expected string scan to consume escaped string, end=%d err=%v", strEnd, err)
	}
	scalarEnd, err := ScanJSONValueEnd([]byte("true,rest"), 0)
	if err != nil || scalarEnd != len("true") {
		t.Fatalf("expected scalar scan to stop at comma, end=%d err=%v", scalarEnd, err)
	}
	if _, err := ScanJSONValueEnd([]byte(`{"bad":[1}`), 0); err == nil {
		t.Fatal("expected mismatched composite to fail")
	}
	if _, err := ScanJSONString([]byte(`bad`), 0); err == nil {
		t.Fatal("expected non-string scan to fail")
	}
}

// ==================== IsValidJSONBytes ====================

func TestIsValidJSONBytes(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{"valid object", `{"key":"value"}`, true},
		{"valid array", `[1,2,3]`, true},
		{"valid null", `null`, true},
		{"valid number", `42`, true},
		{"valid string", `"hello"`, true},
		{"valid with whitespace", `  {"a":1}  `, true},
		{"invalid object", `{invalid}`, false},
		{"invalid trailing comma", `{"a":1,}`, false},
		{"invalid unclosed", `{"a":`, false},
		{"empty bytes", ``, false},
		{"invalid bare text", `hello`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidJSONBytes([]byte(tt.data)); got != tt.want {
				t.Errorf("IsValidJSONBytes(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

// ==================== IsJSONSpace ====================

func TestIsJSONSpace(t *testing.T) {
	tests := []struct {
		name string
		c    byte
		want bool
	}{
		{"space", ' ', true},
		{"tab", '\t', true},
		{"newline", '\n', true},
		{"carriage return", '\r', true},
		{"letter", 'a', false},
		{"digit", '0', false},
		{"quote", '"', false},
		{"brace open", '{', false},
		{"bracket close", ']', false},
		{"comma", ',', false},
		{"colon", ':', false},
		{"null byte", 0, false},
		{"vertical tab (not json space)", '\v', false},
		{"form feed (not json space)", '\f', false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSONSpace(tt.c); got != tt.want {
				t.Errorf("IsJSONSpace(%q) = %v, want %v", tt.c, got, tt.want)
			}
		})
	}
}

// ==================== SkipJSONSpaces ====================

func TestSkipJSONSpaces(t *testing.T) {
	tests := []struct {
		name string
		data string
		i    int
		want int
	}{
		{"no leading space", `{"a":1}`, 0, 0},
		{"single space", ` {"a":1}`, 0, 1},
		{"mixed spaces", " \t\n\r{}", 0, 4},
		{"all spaces (loop exits via condition)", "   ", 0, 3},
		{"all mixed spaces to end", " \t\n ", 0, 4},
		{"start beyond length", `abc`, 5, 5},
		{"start at non-space mid", `a b c`, 2, 2},
		{"empty data", ``, 0, 0},
		{"start in middle of spaces", "   x", 1, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SkipJSONSpaces([]byte(tt.data), tt.i); got != tt.want {
				t.Errorf("SkipJSONSpaces(%q, %d) = %d, want %d", tt.data, tt.i, got, tt.want)
			}
		})
	}
}

// ==================== ScanJSONString ====================

func TestScanJSONString(t *testing.T) {
	t.Run("simple string", func(t *testing.T) {
		end, err := ScanJSONString([]byte(`"abc"`), 0)
		if err != nil || end != 5 {
			t.Fatalf("expected end=5 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("escaped string", func(t *testing.T) {
		end, err := ScanJSONString([]byte(`"a\"b"`), 0)
		if err != nil || end != len(`"a\"b"`) {
			t.Fatalf("expected end=%d err=nil, got end=%d err=%v", len(`"a\"b"`), end, err)
		}
	})
	t.Run("escaped quote inside string", func(t *testing.T) {
		// `"x\"y"` => 字符串内容为 x"y，结尾引号在最后
		end, err := ScanJSONString([]byte(`"x\"y"`), 0)
		if err != nil || end != 6 {
			t.Fatalf("expected end=6 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("non-string start", func(t *testing.T) {
		_, err := ScanJSONString([]byte(`abc`), 0)
		if err == nil {
			t.Fatal("expected error for non-string start")
		}
	})
	t.Run("start beyond length", func(t *testing.T) {
		_, err := ScanJSONString([]byte(`"abc"`), 10)
		if err == nil {
			t.Fatal("expected error for start beyond length")
		}
	})
	t.Run("unterminated string", func(t *testing.T) {
		// no closing quote -> loop ends without returning -> error
		_, err := ScanJSONString([]byte(`"abc`), 0)
		if err == nil {
			t.Fatal("expected error for unterminated string")
		}
	})
	t.Run("empty string quotes", func(t *testing.T) {
		end, err := ScanJSONString([]byte(`""`), 0)
		if err != nil || end != 2 {
			t.Fatalf("expected end=2 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("only opening quote", func(t *testing.T) {
		_, err := ScanJSONString([]byte(`"`), 0)
		if err == nil {
			t.Fatal("expected error for lone opening quote")
		}
	})
}

// ==================== ScanJSONValueEnd ====================

func TestScanJSONValueEnd(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`"abc",rest`), 0)
		if err != nil || end != 5 {
			t.Fatalf("expected end=5 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("object value", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`{"a":1}`), 0)
		if err != nil || end != 7 {
			t.Fatalf("expected end=7 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("array value", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`[1,2,3]`), 0)
		if err != nil || end != 7 {
			t.Fatalf("expected end=7 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar value", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`true,rest`), 0)
		if err != nil || end != 4 {
			t.Fatalf("expected end=4 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("start beyond length", func(t *testing.T) {
		_, err := ScanJSONValueEnd([]byte(`abc`), 5)
		if err == nil {
			t.Fatal("expected error for start beyond length")
		}
	})
	t.Run("empty data", func(t *testing.T) {
		_, err := ScanJSONValueEnd([]byte(``), 0)
		if err == nil {
			t.Fatal("expected error for empty data")
		}
	})
}

// ==================== scanJSONCompositeEnd ====================

func TestScanJSONCompositeEnd(t *testing.T) {
	t.Run("nested object", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`{"a":{"b":1}}`), 0)
		if err != nil || end != 13 {
			t.Fatalf("expected end=13 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("array inside object", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`{"a":[1,2]}`), 0)
		if err != nil || end != 11 {
			t.Fatalf("expected end=11 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("mismatched closing", func(t *testing.T) {
		_, err := ScanJSONValueEnd([]byte(`{"bad":1]`), 0)
		if err == nil {
			t.Fatal("expected error for mismatched closing bracket")
		}
	})
	t.Run("lone closing brace direct", func(t *testing.T) {
		// 直接调用 scanJSONCompositeEnd，覆盖 stack 为空时 last < 0 分支
		_, err := scanJSONCompositeEnd([]byte(`}`), 0)
		if err == nil {
			t.Fatal("expected error for lone closing brace")
		}
	})
	t.Run("lone closing bracket direct", func(t *testing.T) {
		_, err := scanJSONCompositeEnd([]byte(`]`), 0)
		if err == nil {
			t.Fatal("expected error for lone closing bracket")
		}
	})
	t.Run("unterminated composite", func(t *testing.T) {
		_, err := ScanJSONValueEnd([]byte(`{"a":1`), 0)
		if err == nil {
			t.Fatal("expected error for unterminated composite")
		}
	})
	t.Run("broken string inside composite", func(t *testing.T) {
		// ScanJSONString inside composite returns error -> propagate
		_, err := ScanJSONValueEnd([]byte(`{"a":"bad`), 0)
		if err == nil {
			t.Fatal("expected error for broken string inside composite")
		}
	})
	t.Run("string with escaped quote inside composite", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`{"a":"b\"c"}`), 0)
		if err != nil || end != 12 {
			t.Fatalf("expected end=12 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("unclosed outer with closed inner", func(t *testing.T) {
		// inner object closes, but outer remains on stack -> loop ends with non-empty stack -> error
		_, err := ScanJSONValueEnd([]byte(`{"a":{}`), 0)
		if err == nil {
			t.Fatal("expected error for unclosed outer composite")
		}
	})
}

// ==================== scanJSONScalarEnd ====================

func TestScanJSONScalarEnd(t *testing.T) {
	// Access via ScanJSONValueEnd (dispatched to scanJSONScalarEnd for non-string/object/array).
	t.Run("scalar then comma", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`true,rest`), 0)
		if err != nil || end != 4 {
			t.Fatalf("expected end=4 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar then closing brace", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`true}`), 0)
		if err != nil || end != 4 {
			t.Fatalf("expected end=4 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar then closing bracket", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`true]`), 0)
		if err != nil || end != 4 {
			t.Fatalf("expected end=4 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar then whitespace then comma", func(t *testing.T) {
		// "true ,x": whitespace at 4, end=4, inner finds ',' at 5 -> return end(4)
		end, err := ScanJSONValueEnd([]byte(`true ,x`), 0)
		if err != nil || end != 4 {
			t.Fatalf("expected end=4 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar then whitespace then default char", func(t *testing.T) {
		// "true x": whitespace at 4, end=4, inner finds 'x' at 5 -> return i(5)
		end, err := ScanJSONValueEnd([]byte(`true x`), 0)
		if err != nil || end != 5 {
			t.Fatalf("expected end=5 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar then trailing whitespace to end", func(t *testing.T) {
		// "true   ": whitespace at 4, end=4, inner consumes all spaces, i reaches len -> return end(4)
		end, err := ScanJSONValueEnd([]byte(`true   `), 0)
		if err != nil || end != 4 {
			t.Fatalf("expected end=4 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar with mixed whitespace then comma", func(t *testing.T) {
		// "true \t\n,": whitespace at 4, end=4, inner consumes \t\n then ',' -> return end(4)
		end, err := ScanJSONValueEnd([]byte("true \t\n,"), 0)
		if err != nil || end != 4 {
			t.Fatalf("expected end=4 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar with mixed whitespace then default", func(t *testing.T) {
		// "true \t\nx": whitespace at 4, end=4, inner consumes \t\n then 'x' -> return i(7)
		end, err := ScanJSONValueEnd([]byte("true \t\nx"), 0)
		if err != nil || end != 7 {
			t.Fatalf("expected end=7 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar no terminator to end", func(t *testing.T) {
		// "true": no terminator, loop ends, return len(4)
		end, err := ScanJSONValueEnd([]byte(`true`), 0)
		if err != nil || end != 4 {
			t.Fatalf("expected end=4 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("number no terminator to end", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`12345`), 0)
		if err != nil || end != 5 {
			t.Fatalf("expected end=5 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("null then closing brace", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`null}`), 0)
		if err != nil || end != 4 {
			t.Fatalf("expected end=4 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar immediately closing bracket", func(t *testing.T) {
		end, err := ScanJSONValueEnd([]byte(`1]`), 0)
		if err != nil || end != 1 {
			t.Fatalf("expected end=1 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("scalar immediately closing brace with leading space", func(t *testing.T) {
		// " }": starts with space -> whitespace branch, end=0, inner '}' -> return end(0)
		end, err := ScanJSONValueEnd([]byte(` }`), 0)
		if err != nil || end != 0 {
			t.Fatalf("expected end=0 err=nil, got end=%d err=%v", end, err)
		}
	})
	t.Run("only whitespace to end", func(t *testing.T) {
		// "   ": whitespace branch, end=0, inner consumes all -> return end(0)
		end, err := ScanJSONValueEnd([]byte(`   `), 0)
		if err != nil || end != 0 {
			t.Fatalf("expected end=0 err=nil, got end=%d err=%v", end, err)
		}
	})
}
