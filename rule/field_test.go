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

	"github.com/kamalyes/go-argus/utils"
	"github.com/kamalyes/go-argus/validate"
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

func TestFieldIndexByPath(t *testing.T) {
	typ := reflect.TypeOf(nestedStruct{})
	index, ok := FieldIndexByPath(typ, "inner.name")
	if !ok {
		t.Fatal("expected nested json path to resolve")
	}
	value := reflect.ValueOf(nestedStruct{Inner: testStruct{Name: "indexed"}}).FieldByIndex(index)
	if value.String() != "indexed" {
		t.Fatalf("expected indexed value, got %s", value.String())
	}
}

func TestFieldIndexByPathInvalid(t *testing.T) {
	if _, ok := FieldIndexByPath(nil, "Name"); ok {
		t.Fatal("expected nil root to fail")
	}
	if _, ok := FieldIndexByPath(reflect.TypeOf(""), "Name"); ok {
		t.Fatal("expected non-struct root to fail")
	}
	if _, ok := FieldIndexByPath(reflect.TypeOf(testStruct{}), "Missing"); ok {
		t.Fatal("expected missing field to fail")
	}
	if _, ok := FieldIndexByPath(reflect.TypeOf(testStruct{}), ""); ok {
		t.Fatal("expected empty path to fail")
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
	result := IsRequiredIf(reflect.ValueOf(s), []string{"Name", "schedule"})
	if !result {
		t.Fatal("expected required_if to trigger")
	}
}

func TestIsRequiredIfFalse(t *testing.T) {
	s := testStruct{Name: "instant", Age: 0}
	result := IsRequiredIf(reflect.ValueOf(s), []string{"Name", "schedule"})
	if result {
		t.Fatal("expected required_if not to trigger")
	}
}

func TestIsRequiredIfOddParts(t *testing.T) {
	s := testStruct{Name: "schedule"}
	result := IsRequiredIf(reflect.ValueOf(s), []string{"Name"})
	if result {
		t.Fatal("expected false for odd number of parts")
	}
}

func TestIsRequiredIfEmptyParam(t *testing.T) {
	s := testStruct{Name: "schedule"}
	result := IsRequiredIf(reflect.ValueOf(s), []string{})
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

func TestCompareFieldDerefedSuccess(t *testing.T) {
	s := testStruct{Name: "hello", Age: 10}
	result := CompareFieldDerefed(reflect.ValueOf("hello"), reflect.ValueOf(s), "Name", "eq")
	if !result {
		t.Fatal("expected derefed eq comparison to succeed")
	}
}

func TestCompareFieldDerefedNotFound(t *testing.T) {
	s := testStruct{Name: "hello"}
	result := CompareFieldDerefed(reflect.ValueOf("hello"), reflect.ValueOf(s), "Missing", "eq")
	if result {
		t.Fatal("expected false for missing field in derefed")
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

func TestCompareValueOpInvalid(t *testing.T) {
	if CompareValueOp(reflect.ValueOf("a"), reflect.ValueOf("a"), validate.CmpOp(-1)) {
		t.Fatal("expected invalid cmp op to fail")
	}
}

func TestLowerCamel(t *testing.T) {
	if got := utils.LowerCamel("UserName"); got != "userName" {
		t.Fatalf("expected userName, got %s", got)
	}
	if got := utils.LowerCamel(""); got != "" {
		t.Fatalf("expected empty, got %s", got)
	}
}

func TestSnakeCase(t *testing.T) {
	if got := utils.SnakeCase("UserName"); got != "user_name" {
		t.Fatalf("expected user_name, got %s", got)
	}
	if got := utils.SnakeCase("ID"); got != "i_d" {
		t.Fatalf("expected i_d, got %s", got)
	}
}

func TestCompareFloatOps(t *testing.T) {
	cmp := func(left, right float64, op string) bool {
		return validate.CompareOp(left, right, validate.CmpOpFromStr(op))
	}
	if !cmp(1, 1, "eq") {
		t.Fatal("eq failed")
	}
	if cmp(1, 2, "eq") {
		t.Fatal("eq should be false")
	}
	if !cmp(1, 2, "ne") {
		t.Fatal("ne failed")
	}
	if !cmp(2, 1, "gt") {
		t.Fatal("gt failed")
	}
	if !cmp(2, 2, "gte") {
		t.Fatal("gte failed")
	}
	if !cmp(1, 2, "lt") {
		t.Fatal("lt failed")
	}
	if !cmp(2, 2, "lte") {
		t.Fatal("lte failed")
	}
	if cmp(1, 2, "unknown") {
		t.Fatal("unknown op should be false")
	}
}

func TestFloatValueInvalid(t *testing.T) {
	_, ok := validate.NumericValue(reflect.ValueOf("abc"))
	if ok {
		t.Fatal("expected non-numeric string to fail NumericValue")
	}
}

func TestFloatValueUint(t *testing.T) {
	v, ok := validate.NumericValue(reflect.ValueOf(uint(42)))
	if !ok || v != 42 {
		t.Fatalf("expected 42, got %f ok=%v", v, ok)
	}
}

func TestFloatValueFloat(t *testing.T) {
	v, ok := validate.NumericValue(reflect.ValueOf(3.14))
	if !ok || v != 3.14 {
		t.Fatalf("expected 3.14, got %f ok=%v", v, ok)
	}
}

func TestFloatValueInvalidKind(t *testing.T) {
	_, ok := validate.NumericValue(reflect.ValueOf([]int{1}))
	if ok {
		t.Fatal("expected slice to fail NumericValue")
	}
}

func TestFloatValueInvalidValue(t *testing.T) {
	_, ok := validate.NumericValue(reflect.Value{})
	if ok {
		t.Fatal("expected invalid value to fail NumericValue")
	}
}

func TestScalarStringInvalid(t *testing.T) {
	got, _ := validate.ScalarString(reflect.Value{})
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

// --- 补充 field.go 未覆盖函数测试 ---

func TestIsRequiredIfFast(t *testing.T) {
	s := testStruct{Name: "schedule", Age: 0}
	result := IsRequiredIf(reflect.ValueOf(s), []string{"Name", "schedule"})
	if !result {
		t.Fatal("expected required_if to trigger")
	}
}

func TestIsRequiredWithAllTrue(t *testing.T) {
	s := testStruct{Name: "argus", Age: 10}
	result := IsRequiredWithAll(reflect.ValueOf(s), []string{"Name", "Age"})
	if !result {
		t.Fatal("expected required_with_all to trigger")
	}
}

func TestIsRequiredWithAllFalse(t *testing.T) {
	s := testStruct{Name: "argus", Age: 0}
	result := IsRequiredWithAll(reflect.ValueOf(s), []string{"Name", "Age"})
	if result {
		t.Fatal("expected required_with_all to not trigger when one is zero")
	}
}

func TestIsRequiredWithAllEmpty(t *testing.T) {
	s := testStruct{Name: "argus"}
	result := IsRequiredWithAll(reflect.ValueOf(s), []string{})
	if result {
		t.Fatal("expected required_with_all to not trigger for empty parts")
	}
}

func TestIsRequiredWithoutTrue(t *testing.T) {
	s := testStruct{Name: "", Age: 0}
	result := IsRequiredWithout(reflect.ValueOf(s), []string{"Name", "Age"})
	if !result {
		t.Fatal("expected required_without to trigger when field is empty")
	}
}

func TestIsRequiredWithoutFalse(t *testing.T) {
	s := testStruct{Name: "argus", Age: 10}
	result := IsRequiredWithout(reflect.ValueOf(s), []string{"Name", "Age"})
	if result {
		t.Fatal("expected required_without to not trigger when all fields are non-empty")
	}
}

func TestIsRequiredWithoutAllTrue(t *testing.T) {
	s := testStruct{Name: "", Age: 0}
	result := IsRequiredWithoutAll(reflect.ValueOf(s), []string{"Name", "Age"})
	if !result {
		t.Fatal("expected required_without_all to trigger when all fields are empty")
	}
}

func TestIsRequiredWithoutAllFalse(t *testing.T) {
	s := testStruct{Name: "argus", Age: 0}
	result := IsRequiredWithoutAll(reflect.ValueOf(s), []string{"Name", "Age"})
	if result {
		t.Fatal("expected required_without_all to not trigger when some fields are non-empty")
	}
}

func TestIsRequiredWithoutAllEmpty(t *testing.T) {
	s := testStruct{Name: "argus"}
	result := IsRequiredWithoutAll(reflect.ValueOf(s), []string{})
	if result {
		t.Fatal("expected required_without_all to not trigger for empty parts")
	}
}

func TestRangeSuccess(t *testing.T) {
	type rangeStruct struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
	s := rangeStruct{Start: 1, End: 10}
	result := Range(reflect.ValueOf(s), "start,end")
	if !result {
		t.Fatal("expected range to pass when start < end")
	}
}

func TestRangePipeSeparator(t *testing.T) {
	type rangeStruct struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
	s := rangeStruct{Start: 1, End: 10}
	result := Range(reflect.ValueOf(s), "start|end")
	if !result {
		t.Fatal("expected range with pipe separator to pass")
	}
}

func TestRangeFail(t *testing.T) {
	type rangeStruct struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
	s := rangeStruct{Start: 10, End: 1}
	result := Range(reflect.ValueOf(s), "start,end")
	if result {
		t.Fatal("expected range to fail when start >= end")
	}
}

func TestRangeInvalidParts(t *testing.T) {
	s := testStruct{Name: "argus"}
	result := Range(reflect.ValueOf(s), "Name")
	if result {
		t.Fatal("expected range to fail for single part")
	}
}

func TestRangeMissingField(t *testing.T) {
	s := testStruct{Name: "argus"}
	result := Range(reflect.ValueOf(s), "Name,Missing")
	if result {
		t.Fatal("expected range to fail for missing field")
	}
}

func TestFieldContainsSuccess(t *testing.T) {
	s := testStruct{Name: "hello world"}
	result := FieldContains(reflect.ValueOf("hello world"), reflect.ValueOf(s), "Name")
	if !result {
		t.Fatal("expected field contains to pass")
	}
}

func TestFieldContainsFail(t *testing.T) {
	s := testStruct{Name: "hello"}
	result := FieldContains(reflect.ValueOf("hello"), reflect.ValueOf(s), "Missing")
	if result {
		t.Fatal("expected field contains to fail for missing field")
	}
}

func TestFieldContainsNonString(t *testing.T) {
	s := testStruct{Age: 10}
	result := FieldContains(reflect.ValueOf(10), reflect.ValueOf(s), "Age")
	if result {
		t.Fatal("expected field contains to fail for non-string field")
	}
}

func TestOneOfFast(t *testing.T) {
	if !OneOfFast(reflect.ValueOf("a"), []string{"a", "b", "c"}) {
		t.Fatal("expected oneof to pass")
	}
	if OneOfFast(reflect.ValueOf("d"), []string{"a", "b", "c"}) {
		t.Fatal("expected oneof to fail")
	}
}

func TestOneOfFastNonString(t *testing.T) {
	// ScalarString 对 int 返回 "42"，所以 OneOfFast 实际会通过
	if !OneOfFast(reflect.ValueOf(42), []string{"42"}) {
		t.Fatal("expected oneof to pass for int via ScalarString")
	}
}

func TestOneOfCIFast(t *testing.T) {
	if !OneOfCIFast(reflect.ValueOf("A"), []string{"a", "b", "c"}) {
		t.Fatal("expected oneofci to pass")
	}
	if OneOfCIFast(reflect.ValueOf("D"), []string{"a", "b", "c"}) {
		t.Fatal("expected oneofci to fail")
	}
}

func TestOneOfCIFastNonString(t *testing.T) {
	// ScalarString 对 int 返回 "42"，所以 OneOfCIFast 实际会通过
	if !OneOfCIFast(reflect.ValueOf(42), []string{"42"}) {
		t.Fatal("expected oneofci to pass for int via ScalarString")
	}
}

func TestCompareValueScalarStringPath(t *testing.T) {
	// 测试非字符串类型走 scalarString 路径
	if !CompareValue(reflect.ValueOf(10), reflect.ValueOf(10), "eq") {
		t.Fatal("expected scalar string eq to pass")
	}
	if !CompareValue(reflect.ValueOf(10), reflect.ValueOf(5), "gt") {
		t.Fatal("expected scalar string gt to pass")
	}
	if !CompareValue(reflect.ValueOf(5), reflect.ValueOf(10), "lt") {
		t.Fatal("expected scalar string lt to pass")
	}
	if !CompareValue(reflect.ValueOf(5), reflect.ValueOf(5), "gte") {
		t.Fatal("expected scalar string gte to pass")
	}
	if !CompareValue(reflect.ValueOf(5), reflect.ValueOf(5), "lte") {
		t.Fatal("expected scalar string lte to pass")
	}
	if !CompareValue(reflect.ValueOf(5), reflect.ValueOf(10), "ne") {
		t.Fatal("expected scalar string ne to pass")
	}
}

func TestFieldByPathNonStructRoot(t *testing.T) {
	_, ok := FieldByPath(reflect.ValueOf("hello"), "Name")
	if ok {
		t.Fatal("expected false for non-struct root")
	}
}

// --- 模拟 proto 枚举类型测试 ---

// protoEnum 模拟 protobuf 生成的枚举类型，实现了 String() 方法
type protoEnum int32

const (
	protoEnum_UNSPECIFIED protoEnum = 0
	protoEnum_SCHEDULED   protoEnum = 2
	protoEnum_ADJUST      protoEnum = 1
)

func (e protoEnum) String() string {
	switch e {
	case 0:
		return "UNSPECIFIED"
	case 1:
		return "ADJUST"
	case 2:
		return "SCHEDULED"
	default:
		return "UNKNOWN"
	}
}

type enumStruct struct {
	Timing protoEnum `json:"timing"`
	Name   string    `json:"name"`
}

func TestIsRequiredIfWithProtoEnumByNumber(t *testing.T) {
	// required_if=timing 2 → 当 timing == 2 (SCHEDULED) 时触发
	s := enumStruct{Timing: protoEnum_SCHEDULED, Name: "test"}
	result := IsRequiredIf(reflect.ValueOf(s), []string{"timing", "2"})
	if !result {
		t.Fatal("expected required_if to trigger with numeric value for proto enum")
	}
}

func TestIsRequiredIfWithProtoEnumByName(t *testing.T) {
	// required_if=timing SCHEDULED → 当 timing.String() == "SCHEDULED" 时触发
	s := enumStruct{Timing: protoEnum_SCHEDULED, Name: "test"}
	result := IsRequiredIf(reflect.ValueOf(s), []string{"timing", "SCHEDULED"})
	if !result {
		t.Fatal("expected required_if to trigger with enum name for proto enum")
	}
}

func TestIsRequiredIfWithProtoEnumNotMatch(t *testing.T) {
	// timing == 0 (UNSPECIFIED), 不匹配 SCHEDULED
	s := enumStruct{Timing: protoEnum_UNSPECIFIED, Name: "test"}
	result := IsRequiredIf(reflect.ValueOf(s), []string{"timing", "SCHEDULED"})
	if result {
		t.Fatal("expected required_if not to trigger when enum value doesn't match")
	}
}

func TestIsRequiredIfWithProtoEnumNumberNotMatch(t *testing.T) {
	// timing == 0, 不匹配 "2"
	s := enumStruct{Timing: protoEnum_UNSPECIFIED, Name: "test"}
	result := IsRequiredIf(reflect.ValueOf(s), []string{"timing", "2"})
	if result {
		t.Fatal("expected required_if not to trigger when numeric value doesn't match")
	}
}

func TestMatchScalarStringWithProtoEnum(t *testing.T) {
	// 测试 matchScalarString 同时支持数字和枚举名称
	v := reflect.ValueOf(protoEnum_SCHEDULED)
	if !matchScalarString(v, "2") {
		t.Fatal("expected matchScalarString to match numeric string '2'")
	}
	if !matchScalarString(v, "SCHEDULED") {
		t.Fatal("expected matchScalarString to match enum name 'SCHEDULED'")
	}
	if matchScalarString(v, "1") {
		t.Fatal("expected matchScalarString not to match '1'")
	}
	if matchScalarString(v, "ADJUST") {
		t.Fatal("expected matchScalarString not to match 'ADJUST'")
	}
}
