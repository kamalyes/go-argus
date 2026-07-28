/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-20 22:34:24
 * @FilePath: \go-argus\validate\empty_test.go
 * @Description: empty.go 测试，覆盖空值判断、时间有效性、解引用、归一化和字段名安全校验
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validate

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

type customBool bool

func TestIsEmptyValueInvalid(t *testing.T) {
	if !IsEmptyValue(reflect.Value{}) {
		t.Fatal("expected invalid value to be empty")
	}
}

func TestIsEmptyValueNilPtr(t *testing.T) {
	var p *int
	if !IsEmptyValue(reflect.ValueOf(p)) {
		t.Fatal("expected nil ptr to be empty")
	}
}

func TestIsEmptyValueNilInterface(t *testing.T) {
	var i interface{}
	if !IsEmptyValue(reflect.ValueOf(i)) {
		t.Fatal("expected nil interface to be empty")
	}
}

func TestIsEmptyValueNonNilPtr(t *testing.T) {
	x := 42
	if IsEmptyValue(reflect.ValueOf(&x)) {
		t.Fatal("expected non-nil ptr to not be empty")
	}
}

func TestIsEmptyValueSlice(t *testing.T) {
	if !IsEmptyValue(reflect.ValueOf([]int{})) {
		t.Fatal("expected empty slice to be empty")
	}
	if IsEmptyValue(reflect.ValueOf([]int{1})) {
		t.Fatal("expected non-empty slice to not be empty")
	}
}

func TestIsEmptyValueMap(t *testing.T) {
	if !IsEmptyValue(reflect.ValueOf(map[string]int{})) {
		t.Fatal("expected empty map to be empty")
	}
}

func TestIsEmptyValueArray(t *testing.T) {
	arr := [0]int{}
	if !IsEmptyValue(reflect.ValueOf(arr)) {
		t.Fatal("expected empty array to be empty")
	}
}

func TestIsEmptyValueStringBlank(t *testing.T) {
	if !IsEmptyValue(reflect.ValueOf("   ")) {
		t.Fatal("expected blank string to be empty")
	}
}

func TestIsEmptyValueBool(t *testing.T) {
	if !IsEmptyValue(reflect.ValueOf(false)) {
		t.Fatal("expected false to be empty")
	}
}

func TestIsEmptyValueInt(t *testing.T) {
	if !IsEmptyValue(reflect.ValueOf(0)) {
		t.Fatal("expected zero int to be empty")
	}
}

func TestIsEmptyValueUint(t *testing.T) {
	if !IsEmptyValue(reflect.ValueOf(uint(0))) {
		t.Fatal("expected zero uint to be empty")
	}
}

func TestIsEmptyValueFloat(t *testing.T) {
	if !IsEmptyValue(reflect.ValueOf(0.0)) {
		t.Fatal("expected zero float to be empty")
	}
}

func TestIsEmptyValueTimeZero(t *testing.T) {
	if !IsEmptyValue(reflect.ValueOf(time.Time{})) {
		t.Fatal("expected zero time to be empty")
	}
}

func TestIsEmptyValueStructRequired(t *testing.T) {
	type s struct{ V int }
	if !IsEmptyValue(reflect.ValueOf(s{})) {
		t.Fatal("expected zero struct to be empty with requiredStructEnabled")
	}
}

func TestIsEmptyValueStructNotRequired(t *testing.T) {
	type s struct{ V int }
	if IsEmptyValueWithStruct(reflect.ValueOf(s{}), false) {
		t.Fatal("expected zero struct to not be empty without requiredStructEnabled")
	}
}

func TestIsEmptyValueFunc(t *testing.T) {
	var fn func()
	if !IsEmptyValue(reflect.ValueOf(fn)) {
		t.Fatal("expected nil func to be empty")
	}
}

func TestIsEmptyValueChan(t *testing.T) {
	var ch chan int
	if !IsEmptyValue(reflect.ValueOf(ch)) {
		t.Fatal("expected nil chan to be empty")
	}
}

func TestIsEmptyValueDefault(t *testing.T) {
	var up unsafe.Pointer
	if !IsEmptyValue(reflect.ValueOf(up)) {
		t.Fatal("expected nil unsafe pointer to be empty")
	}
}

func TestIsNilValue(t *testing.T) {
	var p *int
	if !IsNilValue(reflect.ValueOf(p)) {
		t.Fatal("expected nil ptr to be nil")
	}
	if IsNilValue(reflect.ValueOf(42)) {
		t.Fatal("expected int to not be nil")
	}
	if !IsNilValue(reflect.Value{}) {
		t.Fatal("expected invalid value to be nil")
	}
}

func TestDerefReflect(t *testing.T) {
	x := 42
	v := DerefReflect(reflect.ValueOf(&x))
	if v.Int() != 42 {
		t.Fatal("expected deref to 42")
	}
	var p *int
	v = DerefReflect(reflect.ValueOf(p))
	if v.IsValid() {
		t.Fatal("expected nil ptr deref to be invalid")
	}
}

func TestStringValue(t *testing.T) {
	x := 42
	if StringValue(reflect.ValueOf(&x)) != "42" {
		t.Fatal("expected string value 42")
	}
	if StringValue(reflect.Value{}) != "" {
		t.Fatal("expected invalid value to be empty string")
	}
}

func TestStringValueStringKind(t *testing.T) {
	if StringValue(reflect.ValueOf("abc")) != "abc" {
		t.Fatal("expected direct string branch")
	}
}

func TestStringValueCannotInterface(t *testing.T) {
	type secret struct{ v int }
	s := secret{v: 99}
	f := reflect.ValueOf(s).Field(0)
	result := StringValue(f)
	if result == "" {
		t.Fatal("expected non-empty string from unexported field")
	}
}

func TestIsTimeEmptyNil(t *testing.T) {
	if !IsTimeEmpty(nil) {
		t.Fatal("expected nil time to be empty")
	}
}

func TestIsTimeEmptyZero(t *testing.T) {
	if !IsTimeEmpty(&time.Time{}) {
		t.Fatal("expected zero time to be empty")
	}
}

func TestIsTimeEmptyBeforeEpoch(t *testing.T) {
	before := time.Unix(-1, 0)
	if !IsTimeEmpty(&before) {
		t.Fatal("expected before epoch to be empty")
	}
}

func TestIsTimeValid(t *testing.T) {
	now := time.Now()
	if !IsTimeValid(now) {
		t.Fatal("expected current time to be valid")
	}
	if IsTimeValid(nil) {
		t.Fatal("expected nil to be invalid")
	}
	if IsTimeValid(time.Time{}) {
		t.Fatal("expected zero time to be invalid")
	}
	p := &now
	if !IsTimeValid(p) {
		t.Fatal("expected time ptr to be valid")
	}
	if !IsTimeValid("not-time") {
		t.Fatal("expected non-time type to be valid (compat)")
	}
}

func TestHasEmpty(t *testing.T) {
	empty, count := HasEmpty([]interface{}{"", "hello"})
	if !empty || count != 1 {
		t.Fatalf("expected has empty, count=1, got empty=%v count=%d", empty, count)
	}
	empty, count = HasEmpty([]interface{}{})
	if !empty || count != 0 {
		t.Fatal("expected empty slice to have empty")
	}
}

func TestIsAllEmpty(t *testing.T) {
	if !IsAllEmpty([]interface{}{"", 0, nil}) {
		t.Fatal("expected all empty")
	}
	if IsAllEmpty([]interface{}{"hello", 0}) {
		t.Fatal("expected not all empty")
	}
}

func TestIsUndefined(t *testing.T) {
	if !IsUndefined("undefined") {
		t.Fatal("expected undefined")
	}
	if !IsUndefined("  UNDEFINED  ") {
		t.Fatal("expected undefined with spaces")
	}
}

func TestIsNull(t *testing.T) {
	if !IsNull("null") {
		t.Fatal("expected null")
	}
	if !IsNull("  NULL  ") {
		t.Fatal("expected null with spaces")
	}
}

func TestIfNullOrUndefined(t *testing.T) {
	if !IfNullOrUndefined("null") {
		t.Fatal("expected null or undefined")
	}
	if !IfNullOrUndefined("undefined") {
		t.Fatal("expected null or undefined")
	}
	if IfNullOrUndefined("hello") {
		t.Fatal("expected not null or undefined")
	}
}

func TestContainsChinese(t *testing.T) {
	if !ContainsChinese("你好") {
		t.Fatal("expected to contain chinese")
	}
	if ContainsChinese("hello") {
		t.Fatal("expected no chinese")
	}
}

func TestEmptyToDefault(t *testing.T) {
	if EmptyToDefault("", "default") != "default" {
		t.Fatal("expected default for empty")
	}
	if EmptyToDefault("hello", "default") != "hello" {
		t.Fatal("expected original value")
	}
}

func TestIsNil(t *testing.T) {
	if !IsNil(nil) {
		t.Fatal("expected nil to be nil")
	}
	var p *int
	if !IsNil(p) {
		t.Fatal("expected nil ptr to be nil")
	}
	if IsNil(42) {
		t.Fatal("expected int to not be nil")
	}
}

func TestIsFuncType(t *testing.T) {
	if !IsFuncType[func()]() {
		t.Fatal("expected func type")
	}
	if IsFuncType[int]() {
		t.Fatal("expected int not to be func type")
	}
}

func TestIsCEmpty(t *testing.T) {
	if !IsCEmpty(0) {
		t.Fatal("expected zero to be empty")
	}
	if IsCEmpty(42) {
		t.Fatal("expected 42 to not be empty")
	}
}

func TestDerefValue(t *testing.T) {
	x := 42
	v, ok := DerefValue(&x)
	if !ok || v.(int) != 42 {
		t.Fatal("expected deref to 42")
	}
	_, ok = DerefValue(nil)
	if ok {
		t.Fatal("expected nil to not deref")
	}
}

func TestIsSafeFieldName(t *testing.T) {
	if !IsSafeFieldName("hello_world.name") {
		t.Fatal("expected safe field name")
	}
	if IsSafeFieldName("") {
		t.Fatal("expected empty to be unsafe")
	}
	if IsSafeFieldName("hello world") {
		t.Fatal("expected space to be unsafe")
	}
}

func TestIsAllowedFieldWithWhitelist(t *testing.T) {
	if !IsAllowedField("name", []string{"name", "age"}) {
		t.Fatal("expected name in whitelist")
	}
	if IsAllowedField("email", []string{"name", "age"}) {
		t.Fatal("expected email not in whitelist")
	}
}

func TestIsAllowedFieldNoWhitelist(t *testing.T) {
	if !IsAllowedField("name_1") {
		t.Fatal("expected safe field name to be allowed")
	}
}

func TestIsEmptyValueNonNilFunc(t *testing.T) {
	fn := func() {}
	if IsEmptyValue(reflect.ValueOf(fn)) {
		t.Fatal("expected non-nil func to not be empty")
	}
}

func TestIsEmptyValueNonNilChan(t *testing.T) {
	ch := make(chan int)
	if IsEmptyValue(reflect.ValueOf(ch)) {
		t.Fatal("expected non-nil chan to not be empty")
	}
}

func TestDerefValueNilPtr(t *testing.T) {
	var p *int
	_, ok := DerefValue(p)
	if ok {
		t.Fatal("expected nil ptr to not deref")
	}
}

func TestUnwrapProtobufWrapperNilValue(t *testing.T) {
	_, ok := UnwrapProtobufWrapper(nil)
	if ok {
		t.Fatal("expected nil to not unwrap")
	}
}

type mockProtoMultiOut struct{}

func (m mockProtoMultiOut) GetValue() (string, error) { return "", nil }

func TestUnwrapProtobufWrapperMultiOut(t *testing.T) {
	_, ok := UnwrapProtobufWrapper(mockProtoMultiOut{})
	if ok {
		t.Fatal("expected multi-out method to fail unwrap")
	}
}

type mockProtoWithInput struct{}

func (m mockProtoWithInput) GetValue(_ int) string { return "" }

func TestUnwrapProtobufWrapperWithInput(t *testing.T) {
	result, ok := UnwrapProtobufWrapper(mockProtoWithInput{})
	t.Logf("result=%v, ok=%v", result, ok)
	if ok {
		t.Fatal("expected method with input to fail unwrap")
	}
	v := reflect.ValueOf(mockProtoWithInput{})
	m := v.MethodByName("GetValue")
	t.Logf("method valid=%v, numIn=%d, numOut=%d", m.IsValid(), m.Type().NumIn(), m.Type().NumOut())
}

type mockProtoCannotInterface struct{}

func (m mockProtoCannotInterface) GetValue() reflect.Value {
	return reflect.Value{}
}

func TestUnwrapProtobufWrapperEmptyReflectValue(t *testing.T) {
	result, ok := UnwrapProtobufWrapper(mockProtoCannotInterface{})
	if !ok {
		t.Fatal("expected unwrap to succeed with reflect.Value return")
	}
	rv, isRV := result.(reflect.Value)
	if !isRV || rv.IsValid() {
		t.Fatal("expected invalid reflect.Value result")
	}
}

func TestIsEmptyAfterDerefProtobufEmpty(t *testing.T) {
	w := mockProtoWrapper{value: ""}
	_, empty := IsEmptyAfterDeref(w)
	if !empty {
		t.Fatal("expected empty protobuf wrapper to be empty")
	}
}

func TestIsEmptyAfterDerefProtobufNonEmpty(t *testing.T) {
	w := mockProtoWrapper{value: "hello"}
	v, empty := IsEmptyAfterDeref(w)
	if empty {
		t.Fatal("expected non-empty protobuf wrapper to not be empty")
	}
	if v.(string) != "hello" {
		t.Fatal("expected hello")
	}
}

func TestIsEmptyAfterDerefProtobufNilPtrValue(t *testing.T) {
	_, empty := IsEmptyAfterDeref(mockProtoNilPtrValue{})
	if !empty {
		t.Fatal("expected protobuf wrapper nil pointer value to be empty")
	}
}

func TestIsEmptyAfterDerefProtobufScalarZero(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtoIntValue{})
	if empty {
		t.Fatal("expected protobuf scalar zero to be treated as non-empty")
	}
	if v.(int) != 0 {
		t.Fatal("expected protobuf scalar zero value")
	}
}

func TestNormalizeFilterValueInvalidReflect(t *testing.T) {
	result := NormalizeFilterValue(42)
	if result.(int) != 42 {
		t.Fatal("expected int to pass through")
	}
}

func TestNormalizeFilterValueArray(t *testing.T) {
	arr := [2]int{1, 2}
	result := NormalizeFilterValue(arr)
	arrResult, ok := result.([]interface{})
	if !ok || len(arrResult) != 2 {
		t.Fatal("expected array to normalize to slice")
	}
}

type mockProtoWrapper struct {
	value string
}

func (m mockProtoWrapper) GetValue() string { return m.value }

// --- protobuf enum mock 类型 ---
// 模拟 protoc-gen-go 生成的 enum：type XXX int32 + Number() + String()
// Number() 返回命名类型 mockEnumNumber（对应 protoreflect.EnumNumber）
// String() 返回枚举名（对应 LOGIN_RESULT_STATUS_FAILURE 等）

type mockEnumNumber int32

type mockProtobufEnum int32

const (
	mockProtobufEnumUnspecified mockProtobufEnum = 0
	mockProtobufEnumSuccess     mockProtobufEnum = 1
	mockProtobufEnumFailure     mockProtobufEnum = 2
)

func (e mockProtobufEnum) Number() mockEnumNumber { return mockEnumNumber(e) }

func (e mockProtobufEnum) String() string {
	switch e {
	case mockProtobufEnumSuccess:
		return "MOCK_ENUM_SUCCESS"
	case mockProtobufEnumFailure:
		return "MOCK_ENUM_FAILURE"
	default:
		return "MOCK_ENUM_UNSPECIFIED"
	}
}

// mockProtobufEnumNoMethod 没有 Number() 方法的类型，不应被解包
type mockProtobufEnumNoMethod int32

func (e mockProtobufEnumNoMethod) String() string { return "NO_METHOD" }

type mockProtoNilPtrValue struct{}

func (m mockProtoNilPtrValue) GetValue() *int { return nil }

type mockProtoIntValue struct{}

func (m mockProtoIntValue) GetValue() int { return 0 }

// --- 模拟 protobuf wrapperspb 类型 ---

type mockProtoInt32Value struct {
	val int32
}

func (m mockProtoInt32Value) GetValue() int32 { return m.val }

type mockProtoBoolValue struct {
	val bool
}

func (m mockProtoBoolValue) GetValue() bool { return m.val }

type mockProtoFloatValue struct {
	val float64
}

func (m mockProtoFloatValue) GetValue() float64 { return m.val }

type mockProtoNilPtrStringValue struct{}

func (m mockProtoNilPtrStringValue) GetValue() *string { return nil }

type mockProtoNonNilPtrStringValue struct {
	val string
}

func (m mockProtoNonNilPtrStringValue) GetValue() *string { return &m.val }

func TestUnwrapProtobufWrapper(t *testing.T) {
	w := mockProtoWrapper{value: "test"}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok || v.(string) != "test" {
		t.Fatal("expected unwrap to work")
	}
}

type mockProtoNoMethod struct{}

func TestUnwrapProtobufWrapperNoMethod(t *testing.T) {
	_, ok := UnwrapProtobufWrapper(mockProtoNoMethod{})
	if ok {
		t.Fatal("expected no method to fail unwrap")
	}
}

func TestUnwrapProtobufWrapperNoMethodPtr(t *testing.T) {
	ptr := &mockProtoNoMethod{}
	_, ok := UnwrapProtobufWrapper(ptr)
	if ok {
		t.Fatal("expected pointer without method to fail unwrap")
	}
}

type mockProtoAddr struct{ value string }

func (m *mockProtoAddr) GetValue() string { return m.value }

func TestUnwrapProtobufWrapperAddr(t *testing.T) {
	w := &mockProtoAddr{value: "addr"}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok || v.(string) != "addr" {
		t.Fatal("expected addr unwrap to work")
	}
}

func TestUnwrapProtobufWrapperValueWithPtrReceiver(t *testing.T) {
	w := mockProtoAddr{value: "value-receiver-via-pointer-method"}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok || v.(string) != "value-receiver-via-pointer-method" {
		t.Fatal("expected value type to unwrap via pointer receiver method")
	}
}

func TestIsEmptyAfterDerefBool(t *testing.T) {
	v, empty := IsEmptyAfterDeref(false)
	if empty {
		t.Fatal("expected bool false to not be empty")
	}
	if v.(bool) {
		t.Fatal("expected false bool")
	}
}

func TestIsEmptyAfterDerefCustomBool(t *testing.T) {
	v, empty := IsEmptyAfterDeref(customBool(false))
	if empty {
		t.Fatal("expected custom bool to not be empty")
	}
	if v.(customBool) != customBool(false) {
		t.Fatal("expected custom bool false")
	}
}

func TestIsEmptyAfterDerefNil(t *testing.T) {
	_, empty := IsEmptyAfterDeref(nil)
	if !empty {
		t.Fatal("expected nil to be empty")
	}
}

func TestIsEmptyAfterDerefEmptyString(t *testing.T) {
	_, empty := IsEmptyAfterDeref("")
	if !empty {
		t.Fatal("expected empty string to be empty")
	}
}

func TestIsEmptyAfterDerefNonEmpty(t *testing.T) {
	v, empty := IsEmptyAfterDeref("hello")
	if empty {
		t.Fatal("expected non-empty to not be empty")
	}
	if v.(string) != "hello" {
		t.Fatal("expected hello")
	}
}

func TestNormalizeFilterValueNilSlice(t *testing.T) {
	var s []interface{}
	result := NormalizeFilterValue(s)
	if result != nil {
		t.Fatal("expected nil slice to normalize to nil")
	}
}

func TestNormalizeFilterValueNilSliceNonEmpty(t *testing.T) {
	var s []string
	result := NormalizeFilterValue(s)
	if result != nil {
		t.Fatal("expected nil string slice to normalize to nil")
	}
}

func TestNormalizeFilterValueNilSliceViaInterface(t *testing.T) {
	var s []int
	var ifc interface{} = s
	result := NormalizeFilterValue(ifc)
	if result != nil {
		t.Fatal("expected nil slice via interface to normalize to nil")
	}
}

func TestNormalizeFilterValueSlice(t *testing.T) {
	s := []interface{}{"a", 1}
	result := NormalizeFilterValue(s)
	arr, ok := result.([]interface{})
	if !ok || len(arr) != 2 {
		t.Fatal("expected slice to normalize")
	}
}

func TestNormalizeFilterValueInvalid(t *testing.T) {
	result := NormalizeFilterValue(42)
	if result.(int) != 42 {
		t.Fatal("expected int to pass through")
	}
}

func TestNormalizeFilterValueNilInterface(t *testing.T) {
	var nilIfc interface{} = nil
	result := NormalizeFilterValue(nilIfc)
	if result != nil {
		t.Fatal("expected nil interface to return nil")
	}
}

func TestNormalizeFilterValueSliceNil(t *testing.T) {
	result := NormalizeFilterValueSlice(nil)
	if result != nil {
		t.Fatal("expected nil slice")
	}
}

func TestNormalizeFilterValueSliceValues(t *testing.T) {
	result := NormalizeFilterValueSlice([]interface{}{"a", "b"})
	if len(result) != 2 {
		t.Fatal("expected 2 elements")
	}
	if result[0] != "a" {
		t.Fatal("expected a")
	}
	if result[1] != "b" {
		t.Fatal("expected b")
	}
}

func TestNormalizeFilterValueIfNotEmpty(t *testing.T) {
	v, empty := NormalizeFilterValueIfNotEmpty("hello")
	if empty || v.(string) != "hello" {
		t.Fatal("expected non-empty value")
	}
	_, empty = NormalizeFilterValueIfNotEmpty("")
	if !empty {
		t.Fatal("expected empty value")
	}
}

func TestNormalizeFilterValueIfNotEmptyNumericZero(t *testing.T) {
	v, empty := NormalizeFilterValueIfNotEmpty(0)
	if !empty {
		t.Fatal("expected numeric zero to be treated as empty")
	}
	if v != nil {
		t.Fatal("expected normalized value to be nil when empty")
	}
}

func TestIsTimeType(t *testing.T) {
	if !IsTimeType(reflect.TypeOf(time.Time{})) {
		t.Fatal("expected time.Time type to match")
	}
	if IsTimeType(reflect.TypeOf("")) {
		t.Fatal("expected string type to not match time type")
	}
}

func TestIsFilterScalarKind(t *testing.T) {
	if !IsFilterScalarKind(reflect.Int) {
		t.Fatal("expected int kind to be scalar")
	}
	if !IsFilterScalarKind(reflect.Bool) {
		t.Fatal("expected bool kind to be scalar")
	}
	if IsFilterScalarKind(reflect.String) {
		t.Fatal("expected string kind to be non-scalar")
	}
}

func TestNormalizeFilterValueIfEmptySlice(t *testing.T) {
	s := []interface{}{}
	result := NormalizeFilterValueSlice(s)
	if len(result) != 0 {
		t.Fatal("expected empty slice to normalize to empty slice")
	}
}

// --- UnwrapProtobufWrapper 对各种 wrapper 类型的完整测试 ---

func TestUnwrapProtobufWrapperInt32(t *testing.T) {
	w := mockProtoInt32Value{val: 42}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok {
		t.Fatal("expected int32 wrapper to unwrap")
	}
	if v.(int32) != 42 {
		t.Fatalf("expected 42, got %v", v)
	}
}

func TestUnwrapProtobufWrapperInt32Zero(t *testing.T) {
	w := mockProtoInt32Value{val: 0}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok {
		t.Fatal("expected int32 zero wrapper to unwrap")
	}
	if v.(int32) != 0 {
		t.Fatalf("expected 0, got %v", v)
	}
}

func TestUnwrapProtobufWrapperBool(t *testing.T) {
	w := mockProtoBoolValue{val: true}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok {
		t.Fatal("expected bool wrapper to unwrap")
	}
	if v.(bool) != true {
		t.Fatalf("expected true, got %v", v)
	}
}

func TestUnwrapProtobufWrapperBoolFalse(t *testing.T) {
	w := mockProtoBoolValue{val: false}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok {
		t.Fatal("expected bool false wrapper to unwrap")
	}
	if v.(bool) != false {
		t.Fatalf("expected false, got %v", v)
	}
}

func TestUnwrapProtobufWrapperFloat(t *testing.T) {
	w := mockProtoFloatValue{val: 3.14}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok {
		t.Fatal("expected float wrapper to unwrap")
	}
	if v.(float64) != 3.14 {
		t.Fatalf("expected 3.14, got %v", v)
	}
}

func TestUnwrapProtobufWrapperNilPtrString(t *testing.T) {
	w := mockProtoNilPtrStringValue{}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok {
		t.Fatal("expected nil ptr string wrapper to unwrap")
	}
	// GetValue() 返回 *string 类型的 nil，Go typed nil != untyped nil
	ptr, isPtr := v.(*string)
	if !isPtr || ptr != nil {
		t.Fatalf("expected nil *string, got %v", v)
	}
}

func TestUnwrapProtobufWrapperNonNilPtrString(t *testing.T) {
	w := mockProtoNonNilPtrStringValue{val: "hello"}
	v, ok := UnwrapProtobufWrapper(w)
	if !ok {
		t.Fatal("expected non-nil ptr string wrapper to unwrap")
	}
	ptr, isPtr := v.(*string)
	if !isPtr || ptr == nil || *ptr != "hello" {
		t.Fatalf("expected *string with 'hello', got %v", v)
	}
}

// --- IsEmptyAfterDeref 对各种 wrapper 类型的完整链路测试 ---

func TestIsEmptyAfterDerefInt32NonZero(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtoInt32Value{val: 42})
	if empty {
		t.Fatal("expected int32 non-zero to not be empty")
	}
	if v.(int32) != 42 {
		t.Fatalf("expected 42, got %v", v)
	}
}

func TestIsEmptyAfterDerefInt32Zero(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtoInt32Value{val: 0})
	if empty {
		t.Fatal("expected int32 zero to not be empty (scalar kind)")
	}
	if v.(int32) != 0 {
		t.Fatalf("expected 0, got %v", v)
	}
}

func TestIsEmptyAfterDerefBoolTrue(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtoBoolValue{val: true})
	if empty {
		t.Fatal("expected bool true to not be empty")
	}
	if v.(bool) != true {
		t.Fatalf("expected true, got %v", v)
	}
}

func TestIsEmptyAfterDerefBoolFalse(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtoBoolValue{val: false})
	if empty {
		t.Fatal("expected bool false to not be empty (scalar kind)")
	}
	if v.(bool) != false {
		t.Fatalf("expected false, got %v", v)
	}
}

func TestIsEmptyAfterDerefFloatNonZero(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtoFloatValue{val: 1.5})
	if empty {
		t.Fatal("expected float non-zero to not be empty")
	}
	if v.(float64) != 1.5 {
		t.Fatalf("expected 1.5, got %v", v)
	}
}

func TestIsEmptyAfterDerefFloatZero(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtoFloatValue{val: 0.0})
	if empty {
		t.Fatal("expected float zero to not be empty (scalar kind)")
	}
	if v.(float64) != 0.0 {
		t.Fatalf("expected 0.0, got %v", v)
	}
}

func TestIsEmptyAfterDerefNilPtrString(t *testing.T) {
	_, empty := IsEmptyAfterDeref(mockProtoNilPtrStringValue{})
	if !empty {
		t.Fatal("expected nil ptr string wrapper to be empty")
	}
}

func TestIsEmptyAfterDerefNonNilPtrStringNonEmpty(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtoNonNilPtrStringValue{val: "hello"})
	if empty {
		t.Fatal("expected non-nil ptr string to not be empty")
	}
	s, ok := v.(string)
	if !ok || s != "hello" {
		t.Fatalf("expected dereferenced string 'hello', got %v", v)
	}
}

func TestIsEmptyAfterDerefNonNilPtrStringEmpty(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtoNonNilPtrStringValue{val: ""})
	if !empty {
		t.Fatal("expected non-nil ptr string with empty value to be empty")
	}
	_ = v
}

// --- NormalizeFilterValueIfNotEmpty 对各种 wrapper 类型的完整链路测试 ---

func TestNormalizeFilterValueIfNotEmptyInt32NonZero(t *testing.T) {
	v, empty := NormalizeFilterValueIfNotEmpty(mockProtoInt32Value{val: 42})
	if empty {
		t.Fatal("expected int32 non-zero to not be empty")
	}
	if v.(int32) != 42 {
		t.Fatalf("expected 42, got %v", v)
	}
}

func TestNormalizeFilterValueIfNotEmptyInt32Zero(t *testing.T) {
	v, empty := NormalizeFilterValueIfNotEmpty(mockProtoInt32Value{val: 0})
	if empty {
		t.Fatal("expected int32 zero to not be empty (scalar kind)")
	}
	if v.(int32) != 0 {
		t.Fatalf("expected 0, got %v", v)
	}
}

func TestNormalizeFilterValueIfNotEmptyBoolTrue(t *testing.T) {
	v, empty := NormalizeFilterValueIfNotEmpty(mockProtoBoolValue{val: true})
	if empty {
		t.Fatal("expected bool true to not be empty")
	}
	if v.(bool) != true {
		t.Fatalf("expected true, got %v", v)
	}
}

func TestNormalizeFilterValueIfNotEmptyBoolFalse(t *testing.T) {
	v, empty := NormalizeFilterValueIfNotEmpty(mockProtoBoolValue{val: false})
	if empty {
		t.Fatal("expected bool false to not be empty (scalar kind)")
	}
	if v.(bool) != false {
		t.Fatalf("expected false, got %v", v)
	}
}

func TestNormalizeFilterValueIfNotEmptyNilPtrString(t *testing.T) {
	_, empty := NormalizeFilterValueIfNotEmpty(mockProtoNilPtrStringValue{})
	if !empty {
		t.Fatal("expected nil ptr string to be empty")
	}
}

func TestNormalizeFilterValueIfNotEmptyNonNilPtrStringNonEmpty(t *testing.T) {
	v, empty := NormalizeFilterValueIfNotEmpty(mockProtoNonNilPtrStringValue{val: "hello"})
	if empty {
		t.Fatal("expected non-nil ptr string to not be empty")
	}
	s, ok := v.(string)
	if !ok || s != "hello" {
		t.Fatalf("expected dereferenced string 'hello', got %v", v)
	}
}

func TestNormalizeFilterValueIfNotEmptySliceWithWrappers(t *testing.T) {
	// 切片中包含 protobuf wrapper
	slice := []any{mockProtoInt32Value{val: 1}, mockProtoWrapper{value: "test"}}
	v, empty := NormalizeFilterValueIfNotEmpty(slice)
	if empty {
		t.Fatal("expected slice with wrappers to not be empty")
	}
	result, ok := v.([]any)
	if !ok || len(result) != 2 {
		t.Fatalf("expected slice of 2, got %v", v)
	}
	if result[0].(int32) != 1 {
		t.Fatalf("expected first element 1, got %v", result[0])
	}
	if result[1].(string) != "test" {
		t.Fatalf("expected second element 'test', got %v", result[1])
	}
}

// --- UnwrapProtobufEnum 测试 ---

func TestUnwrapProtobufEnumNonZero(t *testing.T) {
	v, ok := UnwrapProtobufEnum(mockProtobufEnumFailure)
	if !ok {
		t.Fatal("expected protobuf enum to be unwrapped")
	}
	// Number() 返回 mockEnumNumber(int32)，应转换为 int64
	n, isInt := v.(int64)
	if !isInt {
		t.Fatalf("expected int64, got %T", v)
	}
	if n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}
}

func TestUnwrapProtobufEnumZero(t *testing.T) {
	v, ok := UnwrapProtobufEnum(mockProtobufEnumUnspecified)
	if !ok {
		t.Fatal("expected protobuf enum zero to be unwrapped")
	}
	if v.(int64) != 0 {
		t.Fatalf("expected 0, got %v", v)
	}
}

func TestUnwrapProtobufEnumNil(t *testing.T) {
	_, ok := UnwrapProtobufEnum(nil)
	if ok {
		t.Fatal("expected nil to not be unwrapped")
	}
}

func TestUnwrapProtobufEnumNoNumberMethod(t *testing.T) {
	_, ok := UnwrapProtobufEnum(mockProtobufEnumNoMethod(2))
	if ok {
		t.Fatal("expected type without Number() to not be unwrapped")
	}
}

func TestUnwrapProtobufEnumPlainInt32(t *testing.T) {
	_, ok := UnwrapProtobufEnum(int32(2))
	if ok {
		t.Fatal("expected plain int32 to not be unwrapped")
	}
}

// --- IsEmptyAfterDeref 对 protobuf enum 的测试 ---

func TestIsEmptyAfterDerefEnumNonZero(t *testing.T) {
	v, empty := IsEmptyAfterDeref(mockProtobufEnumFailure)
	if empty {
		t.Fatal("expected protobuf enum non-zero to not be empty")
	}
	if v.(int64) != 2 {
		t.Fatalf("expected 2, got %v", v)
	}
}

func TestIsEmptyAfterDerefEnumZero(t *testing.T) {
	_, empty := IsEmptyAfterDeref(mockProtobufEnumUnspecified)
	if !empty {
		t.Fatal("expected protobuf enum zero (UNSPECIFIED) to be empty")
	}
}

// --- NormalizeFilterValue 对 protobuf enum 的测试 ---

func TestNormalizeFilterValueEnum(t *testing.T) {
	result := NormalizeFilterValue(mockProtobufEnumFailure)
	n, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T (%v)", result, result)
	}
	if n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}
}

func TestNormalizeFilterValueEnumZero(t *testing.T) {
	result := NormalizeFilterValue(mockProtobufEnumUnspecified)
	if result.(int64) != 0 {
		t.Fatalf("expected 0, got %v", result)
	}
}

// --- NormalizeFilterValueIfNotEmpty 对 protobuf enum 的完整链路测试 ---

func TestNormalizeFilterValueIfNotEmptyEnumNonZero(t *testing.T) {
	v, empty := NormalizeFilterValueIfNotEmpty(mockProtobufEnumFailure)
	if empty {
		t.Fatal("expected protobuf enum non-zero to not be empty")
	}
	n, ok := v.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T (%v)", v, v)
	}
	if n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}
}

func TestNormalizeFilterValueIfNotEmptyEnumZero(t *testing.T) {
	_, empty := NormalizeFilterValueIfNotEmpty(mockProtobufEnumUnspecified)
	if !empty {
		t.Fatal("expected protobuf enum zero (UNSPECIFIED) to be empty")
	}
}

// TestNormalizeFilterValueEnumResultHasNoStringMethod 验证解包后的值不会被驱动序列化为枚举名
// 这是 ClickHouse HTTP 驱动 TYPE_MISMATCH 错误的根因：原枚举值有 String() 方法返回枚举名
func TestNormalizeFilterValueEnumResultHasNoStringMethod(t *testing.T) {
	result := NormalizeFilterValue(mockProtobufEnumFailure)
	// int64 没有 String() 方法，fmt.Sprintf("%v") 会输出数字而非枚举名
	if fmt.Sprintf("%v", result) != "2" {
		t.Fatalf("expected '2', got '%s'", fmt.Sprintf("%v", result))
	}
}
