/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\rule\field_test.go
 * @Description: field.go 测试，覆盖字段路径、跨字段规则和值比较
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"testing"
	"time"
)

type testStruct struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	private string
}

type nestedStruct struct {
	Inner testStruct `json:"inner"`
}

type multiNameStruct struct {
	UserName string `json:"user_name"`
}

func TestFieldByPathEmpty(t *testing.T) {
	_, ok := FieldByPath(reflect.Value{}, "")
	if ok {
		t.Fatal("expected false for empty path")
	}
}

func TestFieldByPathInvalidRoot(t *testing.T) {
	_, ok := FieldByPath(reflect.Value{}, "Name")
	if ok {
		t.Fatal("expected false for invalid root")
	}
}

func TestFieldByPathDirectField(t *testing.T) {
	s := testStruct{Name: "argus"}
	val, ok := FieldByPath(reflect.ValueOf(s), "Name")
	if !ok || val.String() != "argus" {
		t.Fatalf("expected to find Name=argus, ok=%v val=%v", ok, val)
	}
}

func TestFieldByPathJSONName(t *testing.T) {
	s := testStruct{Name: "argus"}
	val, ok := FieldByPath(reflect.ValueOf(s), "name")
	if !ok || val.String() != "argus" {
		t.Fatalf("expected to find name=argus via json tag, ok=%v", ok)
	}
}

func TestFieldByPathSnakeCase(t *testing.T) {
	s := multiNameStruct{UserName: "test"}
	val, ok := FieldByPath(reflect.ValueOf(s), "user_name")
	if !ok || val.String() != "test" {
		t.Fatalf("expected to find user_name via snake_case, ok=%v", ok)
	}
}

func TestFieldByPathLowerCamel(t *testing.T) {
	s := multiNameStruct{UserName: "test"}
	val, ok := FieldByPath(reflect.ValueOf(s), "userName")
	if !ok || val.String() != "test" {
		t.Fatalf("expected to find userName via lowerCamel, ok=%v", ok)
	}
}

func TestFieldByPathPrivateField(t *testing.T) {
	s := testStruct{private: "hidden"}
	_, ok := FieldByPath(reflect.ValueOf(s), "private")
	if ok {
		t.Fatal("expected private field to be inaccessible")
	}
}

func TestFieldByPathNotFound(t *testing.T) {
	s := testStruct{Name: "argus"}
	_, ok := FieldByPath(reflect.ValueOf(s), "Missing")
	if ok {
		t.Fatal("expected false for missing field")
	}
}

func TestFieldByPathNested(t *testing.T) {
	n := nestedStruct{Inner: testStruct{Name: "deep"}}
	val, ok := FieldByPath(reflect.ValueOf(n), "Inner.Name")
	if !ok || val.String() != "deep" {
		t.Fatalf("expected nested Name=deep, ok=%v", ok)
	}
}

func TestFieldByPathNestedNonStruct(t *testing.T) {
	s := testStruct{Name: "argus"}
	_, ok := FieldByPath(reflect.ValueOf(s), "Name.Length")
	if ok {
		t.Fatal("expected false for non-struct nested path")
	}
}

func TestFieldByPathPointer(t *testing.T) {
	s := &testStruct{Name: "ptr"}
	val, ok := FieldByPath(reflect.ValueOf(s), "Name")
	if !ok || val.String() != "ptr" {
		t.Fatalf("expected to deref pointer, ok=%v", ok)
	}
}

func TestFieldByPathEmptyPart(t *testing.T) {
	s := testStruct{Name: "argus"}
	val, ok := FieldByPath(reflect.ValueOf(s), ".Name")
	if !ok || val.String() != "argus" {
		t.Fatalf("expected to skip empty part, ok=%v", ok)
	}
}

func TestFieldByPathNilPointer(t *testing.T) {
	type withPtr struct {
		Name *string
	}
	s := withPtr{Name: nil}
	_, ok := FieldByPath(reflect.ValueOf(s), "Name")
	if !ok {
		t.Fatal("expected nil pointer field to be found")
	}
}

func TestIsRequiredIfTrue(t *testing.T) {
	s := testStruct{Name: "schedule", Age: 0}
	result := IsRequiredIf(reflect.ValueOf(s), "Name schedule")
	if !result {
		t.Fatal("expected required_if to trigger")
	}
}

func TestIsRequiredIfFalse(t *testing.T) {
	s := testStruct{Name: "instant", Age: 0}
	result := IsRequiredIf(reflect.ValueOf(s), "Name schedule")
	if result {
		t.Fatal("expected required_if not to trigger")
	}
}

func TestIsRequiredIfOddParts(t *testing.T) {
	s := testStruct{Name: "schedule"}
	result := IsRequiredIf(reflect.ValueOf(s), "Name")
	if result {
		t.Fatal("expected false for odd number of parts")
	}
}

func TestIsRequiredIfEmptyParam(t *testing.T) {
	s := testStruct{Name: "schedule"}
	result := IsRequiredIf(reflect.ValueOf(s), "")
	if result {
		t.Fatal("expected false for empty param")
	}
}

func TestIsRequiredWithTrue(t *testing.T) {
	s := testStruct{Name: "argus", Age: 10}
	result := IsRequiredWith(reflect.ValueOf(s), "Name")
	if !result {
		t.Fatal("expected required_with to trigger")
	}
}

func TestIsRequiredWithFalse(t *testing.T) {
	s := testStruct{Name: "", Age: 0}
	result := IsRequiredWith(reflect.ValueOf(s), "Name")
	if result {
		t.Fatal("expected required_with not to trigger for empty")
	}
}

func TestIsRequiredWithEmptyParam(t *testing.T) {
	s := testStruct{Name: "argus"}
	result := IsRequiredWith(reflect.ValueOf(s), "")
	if result {
		t.Fatal("expected empty param to not trigger")
	}
}

func TestIsRequiredWithMissingField(t *testing.T) {
	s := testStruct{Name: "argus"}
	result := IsRequiredWith(reflect.ValueOf(s), "NonExistent")
	if result {
		t.Fatal("expected missing field to not trigger")
	}
}

func TestCompareFieldSuccess(t *testing.T) {
	s := testStruct{Name: "hello", Age: 10}
	result := CompareField(reflect.ValueOf("hello"), reflect.ValueOf(s), "Name", "eq")
	if !result {
		t.Fatal("expected eq comparison to succeed")
	}
}

func TestCompareFieldNotFound(t *testing.T) {
	s := testStruct{Name: "hello"}
	result := CompareField(reflect.ValueOf("hello"), reflect.ValueOf(s), "Missing", "eq")
	if result {
		t.Fatal("expected false for missing field")
	}
}

func TestCompareValueNumbers(t *testing.T) {
	if !CompareValue(reflect.ValueOf(10), reflect.ValueOf(5), "gt") {
		t.Fatal("expected 10 > 5")
	}
	if !CompareValue(reflect.ValueOf(5), reflect.ValueOf(5), "eq") {
		t.Fatal("expected 5 == 5")
	}
	if CompareValue(reflect.ValueOf(5), reflect.ValueOf(10), "gt") {
		t.Fatal("expected 5 not > 10")
	}
}

func TestCompareValueStrings(t *testing.T) {
	if !CompareValue(reflect.ValueOf("abc"), reflect.ValueOf("abc"), "eq") {
		t.Fatal("expected string eq")
	}
	if !CompareValue(reflect.ValueOf("b"), reflect.ValueOf("a"), "gt") {
		t.Fatal("expected 'b' > 'a'")
	}
	if !CompareValue(reflect.ValueOf("a"), reflect.ValueOf("b"), "ne") {
		t.Fatal("expected 'a' != 'b'")
	}
}

func TestCompareValueUnknownOp(t *testing.T) {
	if CompareValue(reflect.ValueOf(1), reflect.ValueOf(2), "unknown") {
		t.Fatal("expected false for unknown op")
	}
}

func TestCompareValueTime(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Hour)
	if !CompareValue(reflect.ValueOf(later), reflect.ValueOf(now), "gt") {
		t.Fatal("expected later > now")
	}
}

func TestCompareValueAllStringOps(t *testing.T) {
	if !CompareValue(reflect.ValueOf("b"), reflect.ValueOf("a"), "gte") {
		t.Fatal("expected b >= a for strings")
	}
	if !CompareValue(reflect.ValueOf("a"), reflect.ValueOf("b"), "lte") {
		t.Fatal("expected a <= b for strings")
	}
	if !CompareValue(reflect.ValueOf("a"), reflect.ValueOf("b"), "lt") {
		t.Fatal("expected a < b for strings")
	}
	if !CompareValue(reflect.ValueOf("a"), reflect.ValueOf("a"), "lte") {
		t.Fatal("expected a <= a for strings")
	}
}

func TestCompareValueUintNumbers(t *testing.T) {
	if !CompareValue(reflect.ValueOf(uint(10)), reflect.ValueOf(uint(5)), "gt") {
		t.Fatal("expected uint 10 > 5")
	}
}

func TestCompareValueFloatNumbers(t *testing.T) {
	if !CompareValue(reflect.ValueOf(3.14), reflect.ValueOf(2.71), "gt") {
		t.Fatal("expected 3.14 > 2.71")
	}
}

func TestCompareValueNonComparableValues(t *testing.T) {
	if CompareValue(reflect.ValueOf([]int{1}), reflect.ValueOf([]int{2}), "eq") {
		t.Fatal("expected non-comparable values to fail")
	}
}

func TestCompareValueTimeVsNonTime(t *testing.T) {
	now := time.Now()
	result := CompareValue(reflect.ValueOf(now), reflect.ValueOf("not-a-time"), "gt")
	if result {
		t.Fatal("expected time vs non-time to fail")
	}
}

func TestCompareValueDefaultOp(t *testing.T) {
	if CompareValue(reflect.ValueOf("a"), reflect.ValueOf("b"), "invalid_op") {
		t.Fatal("expected invalid op to return false")
	}
}

func TestLowerCamel(t *testing.T) {
	if got := lowerCamel("UserName"); got != "userName" {
		t.Fatalf("expected userName, got %s", got)
	}
	if got := lowerCamel(""); got != "" {
		t.Fatalf("expected empty, got %s", got)
	}
}

func TestSnakeCase(t *testing.T) {
	if got := snakeCase("UserName"); got != "user_name" {
		t.Fatalf("expected user_name, got %s", got)
	}
	if got := snakeCase("ID"); got != "i_d" {
		t.Fatalf("expected i_d, got %s", got)
	}
}

func TestCompareFloatOps(t *testing.T) {
	if !compareFloat(1, 1, "eq") {
		t.Fatal("eq failed")
	}
	if compareFloat(1, 2, "eq") {
		t.Fatal("eq should be false")
	}
	if !compareFloat(1, 2, "ne") {
		t.Fatal("ne failed")
	}
	if !compareFloat(2, 1, "gt") {
		t.Fatal("gt failed")
	}
	if !compareFloat(2, 2, "gte") {
		t.Fatal("gte failed")
	}
	if !compareFloat(1, 2, "lt") {
		t.Fatal("lt failed")
	}
	if !compareFloat(2, 2, "lte") {
		t.Fatal("lte failed")
	}
	if compareFloat(1, 2, "unknown") {
		t.Fatal("unknown op should be false")
	}
}

func TestFloatValueInvalid(t *testing.T) {
	_, ok := floatValue(reflect.ValueOf("abc"))
	if ok {
		t.Fatal("expected string to fail floatValue")
	}
}

func TestFloatValueUint(t *testing.T) {
	v, ok := floatValue(reflect.ValueOf(uint(42)))
	if !ok || v != 42 {
		t.Fatalf("expected 42, got %f ok=%v", v, ok)
	}
}

func TestFloatValueFloat(t *testing.T) {
	v, ok := floatValue(reflect.ValueOf(3.14))
	if !ok || v != 3.14 {
		t.Fatalf("expected 3.14, got %f ok=%v", v, ok)
	}
}

func TestFloatValueInvalidKind(t *testing.T) {
	_, ok := floatValue(reflect.ValueOf([]int{1}))
	if ok {
		t.Fatal("expected slice to fail floatValue")
	}
}

func TestFloatValueInvalidValue(t *testing.T) {
	_, ok := floatValue(reflect.Value{})
	if ok {
		t.Fatal("expected invalid value to fail floatValue")
	}
}

func TestScalarStringInvalid(t *testing.T) {
	got := scalarString(reflect.Value{})
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
