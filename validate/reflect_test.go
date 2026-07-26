/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-20 00:00:00
 * @FilePath: \go-argus\validate\reflect_test.go
 * @Description: reflect.go 测试，覆盖反射值提取工具函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"reflect"
	"testing"
)

// coverageHolder 含一个未导出字段，用于触发 CanInterface() == false 分支
type coverageHolder struct {
	hidden  int
	Visible string
}

// ==================== NumericValue ====================

func TestNumericValue(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(42))
		if !ok || v != 42 {
			t.Fatalf("expected 42, got %f ok=%v", v, ok)
		}
	})
	t.Run("uint", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(uint(42)))
		if !ok || v != 42 {
			t.Fatalf("expected 42, got %f ok=%v", v, ok)
		}
	})
	t.Run("float64", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(3.14))
		if !ok || v != 3.14 {
			t.Fatalf("expected 3.14, got %f ok=%v", v, ok)
		}
	})
	t.Run("string numeric", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf("42"))
		if !ok || v != 42 {
			t.Fatalf("expected 42, got %f ok=%v", v, ok)
		}
	})
	t.Run("int8", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(int8(8)))
		if !ok || v != 8 {
			t.Fatalf("expected 8, got %f ok=%v", v, ok)
		}
	})
	t.Run("int16", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(int16(16)))
		if !ok || v != 16 {
			t.Fatalf("expected 16, got %f ok=%v", v, ok)
		}
	})
	t.Run("int32", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(int32(32)))
		if !ok || v != 32 {
			t.Fatalf("expected 32, got %f ok=%v", v, ok)
		}
	})
	t.Run("int64", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(int64(64)))
		if !ok || v != 64 {
			t.Fatalf("expected 64, got %f ok=%v", v, ok)
		}
	})
	t.Run("uint8", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(uint8(8)))
		if !ok || v != 8 {
			t.Fatalf("expected 8, got %f ok=%v", v, ok)
		}
	})
	t.Run("uint16", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(uint16(16)))
		if !ok || v != 16 {
			t.Fatalf("expected 16, got %f ok=%v", v, ok)
		}
	})
	t.Run("uint32", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(uint32(32)))
		if !ok || v != 32 {
			t.Fatalf("expected 32, got %f ok=%v", v, ok)
		}
	})
	t.Run("uint64", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(uint64(64)))
		if !ok || v != 64 {
			t.Fatalf("expected 64, got %f ok=%v", v, ok)
		}
	})
	t.Run("uintptr", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(uintptr(100)))
		if !ok || v != 100 {
			t.Fatalf("expected 100, got %f ok=%v", v, ok)
		}
	})
	t.Run("float32", func(t *testing.T) {
		v, ok := NumericValue(reflect.ValueOf(float32(1.5)))
		if !ok || v != 1.5 {
			t.Fatalf("expected 1.5, got %f ok=%v", v, ok)
		}
	})
	t.Run("pointer deref", func(t *testing.T) {
		n := 42
		v, ok := NumericValue(reflect.ValueOf(&n))
		if !ok || v != 42 {
			t.Fatalf("expected 42, got %f ok=%v", v, ok)
		}
	})
	t.Run("nil pointer", func(t *testing.T) {
		var n *int
		_, ok := NumericValue(reflect.ValueOf(n))
		if ok {
			t.Fatal("expected nil pointer to fail")
		}
	})
	t.Run("slice invalid", func(t *testing.T) {
		_, ok := NumericValue(reflect.ValueOf([]int{1}))
		if ok {
			t.Fatal("expected slice to fail")
		}
	})
	t.Run("invalid value", func(t *testing.T) {
		_, ok := NumericValue(reflect.Value{})
		if ok {
			t.Fatal("expected invalid value to fail")
		}
	})
}

// ==================== StringValueFromField ====================

func TestStringValueFromField(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		s, ok := StringValueFromField(reflect.ValueOf("hello"))
		if !ok || s != "hello" {
			t.Fatalf("expected hello, got %s ok=%v", s, ok)
		}
	})
	t.Run("non-string int", func(t *testing.T) {
		_, ok := StringValueFromField(reflect.ValueOf(42))
		if ok {
			t.Fatal("expected int to fail")
		}
	})
	t.Run("invalid value", func(t *testing.T) {
		_, ok := StringValueFromField(reflect.Value{})
		if ok {
			t.Fatal("expected invalid value to fail")
		}
	})
	t.Run("pointer to string", func(t *testing.T) {
		s := "hello"
		got, ok := StringValueFromField(reflect.ValueOf(&s))
		if !ok || got != "hello" {
			t.Fatalf("expected hello, got %s ok=%v", got, ok)
		}
	})
	t.Run("nil string pointer", func(t *testing.T) {
		var s *string
		_, ok := StringValueFromField(reflect.ValueOf(s))
		if ok {
			t.Fatal("expected nil pointer to fail")
		}
	})
}

// ==================== BytesValue ====================

func TestBytesValue(t *testing.T) {
	t.Run("valid bytes", func(t *testing.T) {
		b, ok := BytesValue(reflect.ValueOf([]byte("hello")))
		if !ok || string(b) != "hello" {
			t.Fatalf("expected hello, got %s ok=%v", string(b), ok)
		}
	})
	t.Run("string fails", func(t *testing.T) {
		_, ok := BytesValue(reflect.ValueOf("hello"))
		if ok {
			t.Fatal("expected string to fail BytesValue")
		}
	})
	t.Run("invalid value", func(t *testing.T) {
		_, ok := BytesValue(reflect.Value{})
		if ok {
			t.Fatal("expected invalid value to fail")
		}
	})
	t.Run("non-byte slice", func(t *testing.T) {
		_, ok := BytesValue(reflect.ValueOf([]int{1, 2}))
		if ok {
			t.Fatal("expected []int to fail")
		}
	})
	t.Run("pointer to bytes", func(t *testing.T) {
		b := []byte("ptr")
		got, ok := BytesValue(reflect.ValueOf(&b))
		if !ok || string(got) != "ptr" {
			t.Fatalf("expected ptr, got %s ok=%v", string(got), ok)
		}
	})
}

// ==================== ScalarString ====================

func TestScalarString(t *testing.T) {
	t.Run("string kind", func(t *testing.T) {
		s, ok := ScalarString(reflect.ValueOf("hello"))
		if !ok || s != "hello" {
			t.Fatalf("expected hello, got %s ok=%v", s, ok)
		}
	})
	t.Run("non-string int", func(t *testing.T) {
		s, ok := ScalarString(reflect.ValueOf(42))
		if !ok || s != "42" {
			t.Fatalf("expected 42, got %s ok=%v", s, ok)
		}
	})
	t.Run("non-string bool", func(t *testing.T) {
		s, ok := ScalarString(reflect.ValueOf(true))
		if !ok || s != "true" {
			t.Fatalf("expected true, got %s ok=%v", s, ok)
		}
	})
	t.Run("non-string float", func(t *testing.T) {
		s, ok := ScalarString(reflect.ValueOf(3.14))
		if !ok || s != "3.14" {
			t.Fatalf("expected 3.14, got %s ok=%v", s, ok)
		}
	})
	t.Run("non-string slice", func(t *testing.T) {
		s, ok := ScalarString(reflect.ValueOf([]int{1, 2, 3}))
		if !ok || s == "" {
			t.Fatalf("expected non-empty Sprint, got %s ok=%v", s, ok)
		}
	})
	t.Run("invalid value", func(t *testing.T) {
		_, ok := ScalarString(reflect.Value{})
		if ok {
			t.Fatal("expected invalid value to fail")
		}
	})
	t.Run("unexported field (CanInterface false)", func(t *testing.T) {
		h := coverageHolder{hidden: 42, Visible: "x"}
		v := reflect.ValueOf(h).Field(0) // hidden 字段，未导出
		if v.CanInterface() {
			t.Fatal("expected unexported field to have CanInterface()==false")
		}
		s, ok := ScalarString(v)
		if ok {
			t.Fatalf("expected ok=false for unexported field, got ok=true s=%q", s)
		}
		if s != "" {
			t.Fatalf("expected empty string for unexported field, got %q", s)
		}
	})
	t.Run("exported field", func(t *testing.T) {
		h := coverageHolder{hidden: 42, Visible: "x"}
		v := reflect.ValueOf(h).Field(1) // Visible 字段，已导出
		s, ok := ScalarString(v)
		if !ok || s != "x" {
			t.Fatalf("expected x, got %s ok=%v", s, ok)
		}
	})
}

// ==================== MatchStringRunes ====================

func TestMatchStringRunes(t *testing.T) {
	t.Run("all match", func(t *testing.T) {
		if !MatchStringRunes(reflect.ValueOf("abc"), func(r rune) bool { return r >= 'a' && r <= 'z' }) {
			t.Fatal("expected match to pass")
		}
	})
	t.Run("fail uppercase", func(t *testing.T) {
		if MatchStringRunes(reflect.ValueOf("abC"), func(r rune) bool { return r >= 'a' && r <= 'z' }) {
			t.Fatal("expected match to fail for uppercase")
		}
	})
	t.Run("empty string", func(t *testing.T) {
		if MatchStringRunes(reflect.ValueOf(""), func(r rune) bool { return true }) {
			t.Fatal("expected empty string to fail")
		}
	})
	t.Run("non-string", func(t *testing.T) {
		if MatchStringRunes(reflect.ValueOf(42), func(r rune) bool { return true }) {
			t.Fatal("expected int to fail")
		}
	})
	t.Run("digits match", func(t *testing.T) {
		if !MatchStringRunes(reflect.ValueOf("123"), func(r rune) bool { return r >= '0' && r <= '9' }) {
			t.Fatal("expected digits to match")
		}
	})
	t.Run("nil pointer", func(t *testing.T) {
		var s *string
		if MatchStringRunes(reflect.ValueOf(s), func(r rune) bool { return true }) {
			t.Fatal("expected nil pointer to fail")
		}
	})
}

// ==================== ParseFloat ====================

func TestParseFloat(t *testing.T) {
	// utils.ParseFloat 是零依赖实现：不处理负号、科学计数法；空串返回 (0, true)
	tests := []struct {
		name string
		in   string
		want float64
		ok   bool
	}{
		{"integer", "42", 42, true},
		{"float", "3.14", 3.14, true},
		{"zero", "0", 0, true},
		{"negative unsupported", "-5", 0, false},
		{"empty returns zero ok", "", 0, true},
		{"non-numeric", "abc", 0, false},
		{"trailing dot parses", "1.", 1, true},
		{"letter prefix", "a1", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, ok := ParseFloat(tt.in)
			if ok != tt.ok {
				t.Fatalf("expected ok=%v, got ok=%v (v=%f)", tt.ok, ok, v)
			}
			if ok && v != tt.want {
				t.Errorf("expected %f, got %f", tt.want, v)
			}
		})
	}
}
