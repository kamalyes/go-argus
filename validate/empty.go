/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\empty.go
 * @Description: 空值、时间有效性、解引用和过滤值归一化能力
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// IsEmptyValue 判断 reflect.Value 是否为空值
func IsEmptyValue(v reflect.Value) bool {
	return IsEmptyValueWithStruct(v, true)
}

// IsEmptyValueWithStruct 判断 reflect.Value 是否为空值，并允许控制结构体零值语义
func IsEmptyValueWithStruct(v reflect.Value, requiredStructEnabled bool) bool {
	if !v.IsValid() {
		return true
	}
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.String:
		s := v.String()
		for i := 0; i < len(s); i++ {
			if s[i] != ' ' && s[i] != '\t' && s[i] != '\n' && s[i] != '\r' {
				return false
			}
		}
		return true
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			t, _ := v.Interface().(time.Time)
			return IsTimeEmpty(&t)
		}
		if !requiredStructEnabled {
			return false
		}
		return v.IsZero()
	case reflect.Func, reflect.Chan:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

// IsNilValue 判断 reflect.Value 是否为 nil
func IsNilValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// DerefReflect 解开指针和 interface，遇到 nil 返回无效 Value
func DerefReflect(v reflect.Value) reflect.Value {
	for v.IsValid() && (v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr) {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

// StringValue 将 reflect.Value 转为字符串
func StringValue(v reflect.Value) string {
	v = DerefReflect(v)
	if !v.IsValid() {
		return ""
	}
	if v.Kind() == reflect.String {
		return v.String()
	}
	if v.CanInterface() {
		return fmt.Sprint(v.Interface())
	}
	return fmt.Sprint(v)
}

// IsTimeEmpty 判断时间是否为空或早于 Unix epoch
func IsTimeEmpty(t *time.Time) bool {
	if t == nil {
		return true
	}
	return t.IsZero() || t.Unix() <= 0
}

// IsTimeValid 判断时间值是否有效，非时间类型视为有效以兼容过滤场景
func IsTimeValid(timeVal interface{}) bool {
	if timeVal == nil {
		return false
	}
	switch v := timeVal.(type) {
	case time.Time:
		return !v.IsZero() && v.After(time.Unix(0, 0))
	case *time.Time:
		return v != nil && !v.IsZero() && v.After(time.Unix(0, 0))
	default:
		return true
	}
}

// HasEmpty 判断切片中是否存在空值
func HasEmpty(elems []interface{}) (bool, int) {
	if len(elems) == 0 {
		return true, 0
	}
	count := 0
	for _, elem := range elems {
		if IsEmptyValue(reflect.ValueOf(elem)) {
			count++
		}
	}
	return count > 0, count
}

// IsAllEmpty 判断切片中所有元素是否都为空
func IsAllEmpty(elems []interface{}) bool {
	for _, elem := range elems {
		if !IsEmptyValue(reflect.ValueOf(elem)) {
			return false
		}
	}
	return true
}

// IsUndefined 判断字符串是否为 undefined
func IsUndefined(str string) bool {
	return strings.EqualFold(strings.TrimSpace(str), "undefined")
}

// IsNull 判断字符串是否为 null
func IsNull(str string) bool {
	return strings.EqualFold(strings.TrimSpace(str), "null")
}

// IfNullOrUndefined 判断字符串是否为 null 或 undefined
func IfNullOrUndefined(str string) bool {
	return IsNull(str) || IsUndefined(str)
}

// ContainsChinese 判断字符串是否包含中文
func ContainsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// EmptyToDefault 在字符串为空时返回默认值
func EmptyToDefault(str string, defaultStr string) string {
	if strings.TrimSpace(str) == "" {
		return defaultStr
	}
	return str
}

// IsNil 判断 interface 是否为 nil 或内部持有 nil
func IsNil(x interface{}) bool {
	if x == nil {
		return true
	}
	return IsNilValue(reflect.ValueOf(x))
}

// IsFuncType 判断泛型类型是否为函数
func IsFuncType[T any]() bool {
	var zero T
	t := reflect.TypeOf(zero)
	return t != nil && t.Kind() == reflect.Func
}

// IsCEmpty 判断可比较值是否为零值
func IsCEmpty[T comparable](v T) bool {
	var zero T
	return v == zero
}

// DerefValue 解开 interface 中的指针值
func DerefValue(value interface{}) (interface{}, bool) {
	if value == nil {
		return nil, false
	}
	v := DerefReflect(reflect.ValueOf(value))
	if !v.IsValid() || !v.CanInterface() {
		return nil, false
	}
	return v.Interface(), true
}

// IsSafeFieldName 判断字段名是否只包含字母、数字、下划线和点号
func IsSafeFieldName(field string) bool {
	if field == "" {
		return false
	}
	for _, ch := range field {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '_' || ch == '.') {
			return false
		}
	}
	return true
}

// IsAllowedField 判断字段是否在白名单中，未传白名单时退化为安全字段名检查
func IsAllowedField(field string, allowedFields ...[]string) bool {
	if len(allowedFields) > 0 && len(allowedFields[0]) > 0 {
		for _, allowed := range allowedFields[0] {
			if field == allowed {
				return true
			}
		}
		return false
	}
	return IsSafeFieldName(field)
}

// UnwrapProtobufWrapper 通过反射解开 protobuf wrapper，避免引入 protobuf 依赖
func UnwrapProtobufWrapper(value interface{}) (interface{}, bool) {
	if value == nil {
		return nil, false
	}
	v := reflect.ValueOf(value)
	if IsNilValue(v) {
		return nil, true
	}
	method := v.MethodByName("GetValue")
	if !method.IsValid() || method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
		return nil, false
	}
	out := method.Call(nil)
	return out[0].Interface(), true
}

// IsEmptyAfterDeref 解引用后判断值是否为空，适合 SQL/query 过滤条件
func IsEmptyAfterDeref(value interface{}) (interface{}, bool) {
	if unwrapped, ok := UnwrapProtobufWrapper(value); ok {
		if IsEmptyValue(reflect.ValueOf(unwrapped)) {
			return nil, true
		}
		return unwrapped, false
	}
	deref, ok := DerefValue(value)
	if !ok {
		return nil, true
	}
	if b, isBool := deref.(bool); isBool {
		return b, false
	}
	if IsEmptyValue(reflect.ValueOf(deref)) {
		return nil, true
	}
	return deref, false
}

// NormalizeFilterValue 归一化过滤值，支持 protobuf wrapper 和任意切片
func NormalizeFilterValue(value interface{}) interface{} {
	if normalized, ok := UnwrapProtobufWrapper(value); ok {
		return normalized
	}
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return value
	}
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		out := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			out[i] = NormalizeFilterValue(v.Index(i).Interface())
		}
		return out
	}
	return value
}

// NormalizeFilterValueSlice 归一化过滤值切片
func NormalizeFilterValueSlice(values []interface{}) []interface{} {
	if values == nil {
		return nil
	}
	out := make([]interface{}, len(values))
	for i, value := range values {
		out[i] = NormalizeFilterValue(value)
	}
	return out
}

// NormalizeFilterValueIfNotEmpty 过滤空值后返回归一化值
func NormalizeFilterValueIfNotEmpty(value interface{}) (interface{}, bool) {
	deref, isEmpty := IsEmptyAfterDeref(value)
	if isEmpty {
		return nil, true
	}
	return NormalizeFilterValue(deref), false
}
