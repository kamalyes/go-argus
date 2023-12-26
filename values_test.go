/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-29 00:00:00
 * @FilePath: \go-argus\values_test.go
 * @Description: values.go 测试，覆盖空值判断、解引用、过滤值归一化等导出函数
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"reflect"
	"testing"
	"time"
)

func TestIsEmptyValue(t *testing.T) {
	if !IsEmptyValue(reflect.ValueOf("")) {
		t.Fatal("expected empty string to be empty")
	}
	if IsEmptyValue(reflect.ValueOf("hello")) {
		t.Fatal("expected non-empty string to not be empty")
	}
}

func TestIsTimeEmpty(t *testing.T) {
	var zero time.Time
	if !IsTimeEmpty(&zero) {
		t.Fatal("expected zero time to be empty")
	}
	now := time.Now()
	if IsTimeEmpty(&now) {
		t.Fatal("expected current time to not be empty")
	}
	if !IsTimeEmpty(nil) {
		t.Fatal("expected nil time to be empty")
	}
}

func TestIsTimeValid(t *testing.T) {
	now := time.Now()
	if !IsTimeValid(now) {
		t.Fatal("expected time.Time to be valid")
	}
	if IsTimeValid(time.Time{}) {
		t.Fatal("expected zero time.Time to be invalid")
	}
	if !IsTimeValid("not-a-time") {
		t.Fatal("expected string to be valid (default case returns true)")
	}
	if IsTimeValid(nil) {
		t.Fatal("expected nil to not be valid time")
	}
}

func TestHasEmpty(t *testing.T) {
	has, idx := HasEmpty([]interface{}{"hello", "", "world"})
	if !has || idx != 1 {
		t.Fatalf("expected empty at index 1, got has=%v idx=%d", has, idx)
	}
	has, idx = HasEmpty([]interface{}{"a", "b"})
	if has {
		t.Fatal("expected no empty values")
	}
}

func TestIsAllEmpty(t *testing.T) {
	if !IsAllEmpty([]interface{}{"", nil, 0}) {
		t.Fatal("expected all empty")
	}
	if IsAllEmpty([]interface{}{"", "x"}) {
		t.Fatal("expected not all empty")
	}
}

func TestIsUndefined(t *testing.T) {
	if !IsUndefined("undefined") {
		t.Fatal("expected undefined")
	}
	if IsUndefined("defined") {
		t.Fatal("expected not undefined")
	}
}

func TestIsNull(t *testing.T) {
	if !IsNull("null") {
		t.Fatal("expected null")
	}
	if IsNull("notnull") {
		t.Fatal("expected not null")
	}
}

func TestIfNullOrUndefined(t *testing.T) {
	if !IfNullOrUndefined("null") {
		t.Fatal("expected null")
	}
	if !IfNullOrUndefined("undefined") {
		t.Fatal("expected undefined")
	}
	if IfNullOrUndefined("value") {
		t.Fatal("expected neither null nor undefined")
	}
}

func TestContainsChinese(t *testing.T) {
	if !ContainsChinese("你好世界") {
		t.Fatal("expected Chinese characters")
	}
	if ContainsChinese("hello world") {
		t.Fatal("expected no Chinese characters")
	}
}

func TestEmptyToDefault(t *testing.T) {
	if EmptyToDefault("", "default") != "default" {
		t.Fatal("expected default value for empty string")
	}
	if EmptyToDefault("value", "default") != "value" {
		t.Fatal("expected original value for non-empty string")
	}
}

func TestIsNil(t *testing.T) {
	if !IsNil(nil) {
		t.Fatal("expected nil to be nil")
	}
	var ch chan int
	if !IsNil(ch) {
		t.Fatal("expected nil channel to be nil")
	}
	var m map[string]int
	if !IsNil(m) {
		t.Fatal("expected nil map to be nil")
	}
	var s []int
	if !IsNil(s) {
		t.Fatal("expected nil slice to be nil")
	}
	if IsNil("hello") {
		t.Fatal("expected string to not be nil")
	}
}

func TestIsFuncType(t *testing.T) {
	if !IsFuncType[func()]() {
		t.Fatal("expected func type")
	}
	if IsFuncType[string]() {
		t.Fatal("expected string to not be func type")
	}
}

func TestIsCEmpty(t *testing.T) {
	if !IsCEmpty("") {
		t.Fatal("expected empty string to be empty")
	}
	if IsCEmpty("hello") {
		t.Fatal("expected non-empty string to not be empty")
	}
	if !IsCEmpty(0) {
		t.Fatal("expected zero to be empty")
	}
}

func TestDerefValue(t *testing.T) {
	val := 42
	ptr := &val
	result, ok := DerefValue(ptr)
	if !ok || result != 42 {
		t.Fatalf("expected 42, got %v ok=%v", result, ok)
	}
	result, ok = DerefValue(nil)
	if ok {
		t.Fatal("expected nil to not be deref-able")
	}
	result, ok = DerefValue("direct")
	if !ok || result != "direct" {
		t.Fatal("expected direct value")
	}
}

func TestIsSafeFieldName(t *testing.T) {
	if !IsSafeFieldName("valid_field") {
		t.Fatal("expected valid field name")
	}
	if IsSafeFieldName("invalid field!") {
		t.Fatal("expected invalid field name")
	}
}

func TestIsAllowedField(t *testing.T) {
	if !IsAllowedField("name", []string{"name", "age"}) {
		t.Fatal("expected allowed field")
	}
	if IsAllowedField("email", []string{"name", "age"}) {
		t.Fatal("expected disallowed field")
	}
	if !IsAllowedField("any") {
		t.Fatal("expected any field to be allowed with no whitelist")
	}
}

func TestUnwrapProtobufWrapper(t *testing.T) {
	result, ok := UnwrapProtobufWrapper("not-a-wrapper")
	if ok {
		t.Fatalf("expected non-wrapper to not unwrap, got ok=%v result=%v", ok, result)
	}
}

func TestUnwrapProtobufWrapperNil(t *testing.T) {
	result, ok := UnwrapProtobufWrapper(nil)
	if ok {
		t.Fatalf("expected nil to not unwrap, got ok=%v result=%v", ok, result)
	}
}

func TestIsEmptyAfterDeref(t *testing.T) {
	val, ok := IsEmptyAfterDeref("")
	if !ok {
		t.Fatal("expected empty string to be empty after deref")
	}
	if val != nil {
		t.Fatalf("expected nil for empty string, got %v", val)
	}
	val, ok = IsEmptyAfterDeref("hello")
	if ok {
		t.Fatal("expected non-empty string to not be empty after deref")
	}
	if val != "hello" {
		t.Fatalf("expected hello, got %v", val)
	}
}

func TestNormalizeFilterValue(t *testing.T) {
	if NormalizeFilterValue(nil) != nil {
		t.Fatal("expected nil to normalize to nil")
	}
	if NormalizeFilterValue("hello") != "hello" {
		t.Fatal("expected string to pass through")
	}
	result := NormalizeFilterValue([]interface{}{1, 2})
	arr, ok := result.([]interface{})
	if !ok || len(arr) != 2 {
		t.Fatalf("expected slice of length 2, got %v", result)
	}
}

func TestNormalizeFilterValueSlice(t *testing.T) {
	result := NormalizeFilterValueSlice([]interface{}{nil, "hello", 42})
	if len(result) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(result))
	}
}

func TestNormalizeFilterValueIfNotEmpty(t *testing.T) {
	val, ok := NormalizeFilterValueIfNotEmpty("")
	if !ok {
		t.Fatal("expected empty string to be empty (ok=true)")
	}
	if val != nil {
		t.Fatalf("expected nil for empty string, got %v", val)
	}
	val, ok = NormalizeFilterValueIfNotEmpty("hello")
	if ok {
		t.Fatal("expected non-empty string to not be empty (ok=false)")
	}
	if val != "hello" {
		t.Fatalf("expected hello, got %v", val)
	}
	val, ok = NormalizeFilterValueIfNotEmpty(nil)
	if !ok {
		t.Fatal("expected nil to be empty (ok=true)")
	}
}
