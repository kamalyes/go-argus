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

func TestNumericValueInt(t *testing.T) {
	v, ok := NumericValue(reflect.ValueOf(42))
	if !ok || v != 42 {
		t.Fatalf("expected 42, got %f ok=%v", v, ok)
	}
}

func TestNumericValueUint(t *testing.T) {
	v, ok := NumericValue(reflect.ValueOf(uint(42)))
	if !ok || v != 42 {
		t.Fatalf("expected 42, got %f ok=%v", v, ok)
	}
}

func TestNumericValueFloat(t *testing.T) {
	v, ok := NumericValue(reflect.ValueOf(3.14))
	if !ok || v != 3.14 {
		t.Fatalf("expected 3.14, got %f ok=%v", v, ok)
	}
}

func TestNumericValueString(t *testing.T) {
	v, ok := NumericValue(reflect.ValueOf("42"))
	if !ok || v != 42 {
		t.Fatalf("expected 42, got %f ok=%v", v, ok)
	}
}

func TestNumericValueInvalid(t *testing.T) {
	_, ok := NumericValue(reflect.ValueOf([]int{1}))
	if ok {
		t.Fatal("expected slice to fail")
	}
}

func TestNumericValueInvalidReflect(t *testing.T) {
	_, ok := NumericValue(reflect.Value{})
	if ok {
		t.Fatal("expected invalid value to fail")
	}
}

func TestStringValueFromField(t *testing.T) {
	s, ok := StringValueFromField(reflect.ValueOf("hello"))
	if !ok || s != "hello" {
		t.Fatalf("expected hello, got %s ok=%v", s, ok)
	}
}

func TestStringValueFromFieldNonString(t *testing.T) {
	_, ok := StringValueFromField(reflect.ValueOf(42))
	if ok {
		t.Fatal("expected int to fail")
	}
}

func TestStringValueFromFieldInvalid(t *testing.T) {
	_, ok := StringValueFromField(reflect.Value{})
	if ok {
		t.Fatal("expected invalid value to fail")
	}
}

func TestBytesValue(t *testing.T) {
	b, ok := BytesValue(reflect.ValueOf([]byte("hello")))
	if !ok || string(b) != "hello" {
		t.Fatalf("expected hello, got %s ok=%v", string(b), ok)
	}
}

func TestBytesValueInvalid(t *testing.T) {
	_, ok := BytesValue(reflect.ValueOf("hello"))
	if ok {
		t.Fatal("expected string to fail BytesValue")
	}
}

func TestBytesValueInvalidReflect(t *testing.T) {
	_, ok := BytesValue(reflect.Value{})
	if ok {
		t.Fatal("expected invalid value to fail")
	}
}

func TestScalarString(t *testing.T) {
	s, ok := ScalarString(reflect.ValueOf("hello"))
	if !ok || s != "hello" {
		t.Fatalf("expected hello, got %s ok=%v", s, ok)
	}
}

func TestScalarStringNonString(t *testing.T) {
	s, ok := ScalarString(reflect.ValueOf(42))
	if !ok || s != "42" {
		t.Fatalf("expected 42, got %s ok=%v", s, ok)
	}
}

func TestScalarStringInvalid(t *testing.T) {
	_, ok := ScalarString(reflect.Value{})
	if ok {
		t.Fatal("expected invalid value to fail")
	}
}

func TestMatchStringRunes(t *testing.T) {
	if !MatchStringRunes(reflect.ValueOf("abc"), func(r rune) bool { return r >= 'a' && r <= 'z' }) {
		t.Fatal("expected match to pass")
	}
	if MatchStringRunes(reflect.ValueOf("abC"), func(r rune) bool { return r >= 'a' && r <= 'z' }) {
		t.Fatal("expected match to fail for uppercase")
	}
}

func TestMatchStringRunesEmpty(t *testing.T) {
	if MatchStringRunes(reflect.ValueOf(""), func(r rune) bool { return true }) {
		t.Fatal("expected empty string to fail")
	}
}

func TestMatchStringRunesNonString(t *testing.T) {
	if MatchStringRunes(reflect.ValueOf(42), func(r rune) bool { return true }) {
		t.Fatal("expected int to fail")
	}
}

func TestParseFloatDelegate(t *testing.T) {
	v, ok := ParseFloat("3.14")
	if !ok || v != 3.14 {
		t.Fatalf("expected 3.14, got %f ok=%v", v, ok)
	}
}
