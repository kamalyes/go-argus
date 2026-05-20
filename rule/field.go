/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-19 13:58:18
 * @FilePath: \go-argus\rule\field.go
 * @Description: 字段路径解析模块，支持 Go 字段名、json 名、snake_case 和嵌套字段访问
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"strings"
	"sync"

	"github.com/kamalyes/go-argus/utils"
	"github.com/kamalyes/go-argus/validate"
)

var fieldLookupCache sync.Map

// FieldByPath 根据字段路径读取结构体字段
func FieldByPath(root reflect.Value, path string) (reflect.Value, bool) {
	root = validate.DerefReflect(root)
	if !root.IsValid() || path == "" {
		return reflect.Value{}, false
	}
	if !strings.ContainsRune(path, '.') {
		if !root.IsValid() || root.Kind() != reflect.Struct {
			return reflect.Value{}, false
		}
		return directField(root, path)
	}
	current := root
	for _, part := range strings.Split(path, ".") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		current = validate.DerefReflect(current)
		if !current.IsValid() || current.Kind() != reflect.Struct {
			return reflect.Value{}, false
		}
		next, ok := directField(current, part)
		if !ok {
			return reflect.Value{}, false
		}
		current = next
	}
	return current, true
}

// IsRequiredIf 判断 required_if 条件是否触发
func IsRequiredIf(parent reflect.Value, param string) bool {
	parts := strings.Fields(param)
	return isRequiredIfParts(parent, parts)
}

// IsRequiredIfFast 判断 required_if 条件是否触发（预拆分参数版本）
func IsRequiredIfFast(parent reflect.Value, parts []string) bool {
	return isRequiredIfParts(parent, parts)
}

func isRequiredIfParts(parent reflect.Value, parts []string) bool {
	if len(parts) < 2 || len(parts)%2 != 0 {
		return false
	}
	for i := 0; i < len(parts); i += 2 {
		value, ok := FieldByPath(parent, parts[i])
		if !ok || scalarString(value) != parts[i+1] {
			return false
		}
	}
	return true
}

// IsRequiredWith 判断 required_with 条件是否触发
func IsRequiredWith(parent reflect.Value, param string) bool {
	for _, field := range strings.Fields(param) {
		value, ok := FieldByPath(parent, field)
		if ok && !validate.IsEmptyValue(value) {
			return true
		}
	}
	return false
}

// CompareField 将当前字段和目标字段按操作符比较
func CompareField(current reflect.Value, parent reflect.Value, targetPath string, op string) bool {
	target, ok := FieldByPath(parent, targetPath)
	if !ok {
		return false
	}
	return CompareValue(current, target, op)
}

// CompareFieldDerefed 将已解引用的当前字段和目标字段按操作符比较
// 调用方已对 current 做 DerefReflect，避免重复解引用
func CompareFieldDerefed(current reflect.Value, parent reflect.Value, targetPath string, op string) bool {
	target, ok := FieldByPath(parent, targetPath)
	if !ok {
		return false
	}
	return CompareValue(current, target, op)
}

// CompareValue 按操作符比较两个值，优先支持时间，其次支持数值和字符串
func CompareValue(left reflect.Value, right reflect.Value, op string) bool {
	if lt, lok := TimeValue(left, ""); lok {
		rt, rok := TimeValue(right, "")
		return rok && compareFloat(float64(lt.UnixNano()), float64(rt.UnixNano()), op)
	}
	lf, lok := floatValue(left)
	rf, rok := floatValue(right)
	if lok && rok {
		return compareFloat(lf, rf, op)
	}
	// 快速路径：两边都是字符串时直接比较，避免 ScalarString 中的 fmt.Sprint 开销
	if left.Kind() == reflect.String && right.Kind() == reflect.String {
		ls, rs := left.String(), right.String()
		switch op {
		case "eq":
			return ls == rs
		case "ne":
			return ls != rs
		case "gt":
			return ls > rs
		case "gte":
			return ls >= rs
		case "lt":
			return ls < rs
		case "lte":
			return ls <= rs
		default:
			return false
		}
	}
	ls := scalarString(left)
	rs := scalarString(right)
	switch op {
	case "eq":
		return ls == rs
	case "ne":
		return ls != rs
	case "gt":
		return ls > rs
	case "gte":
		return ls >= rs
	case "lt":
		return ls < rs
	case "lte":
		return ls <= rs
	default:
		return false
	}
}

// OneOfFast 判断字段值是否在候选列表中（精确匹配）
func OneOfFast(field reflect.Value, parts []string) bool {
	actual, ok := validate.ScalarString(field)
	if !ok {
		return false
	}
	for _, item := range parts {
		if actual == item {
			return true
		}
	}
	return false
}

// OneOfCIFast 判断字段值是否在候选列表中（忽略大小写）
func OneOfCIFast(field reflect.Value, parts []string) bool {
	actual, ok := validate.ScalarString(field)
	if !ok {
		return false
	}
	for _, item := range parts {
		if strings.EqualFold(actual, item) {
			return true
		}
	}
	return false
}

// IsRequiredWithAll 判断 required_with_all 条件是否触发（所有指定字段均非空）
func IsRequiredWithAll(parent reflect.Value, parts []string) bool {
	if len(parts) == 0 {
		return false
	}
	for _, field := range parts {
		value, ok := FieldByPath(parent, field)
		if !ok || validate.IsEmptyValueWithStruct(value, true) {
			return false
		}
	}
	return true
}

// IsRequiredWithout 判断 required_without 条件是否触发（任一指定字段为空）
func IsRequiredWithout(parent reflect.Value, parts []string) bool {
	for _, field := range parts {
		value, ok := FieldByPath(parent, field)
		if !ok || validate.IsEmptyValueWithStruct(value, true) {
			return true
		}
	}
	return false
}

// IsRequiredWithoutAll 判断 required_without_all 条件是否触发（所有指定字段均为空）
func IsRequiredWithoutAll(parent reflect.Value, parts []string) bool {
	if len(parts) == 0 {
		return false
	}
	for _, field := range parts {
		value, ok := FieldByPath(parent, field)
		if ok && !validate.IsEmptyValueWithStruct(value, true) {
			return false
		}
	}
	return true
}

// Range 判断起始字段值是否小于结束字段值
func Range(parent reflect.Value, param string) bool {
	sep := ","
	if strings.Contains(param, "|") {
		sep = "|"
	}
	parts := strings.Split(param, sep)
	if len(parts) != 2 {
		return false
	}
	start, ok := FieldByPath(parent, strings.TrimSpace(parts[0]))
	if !ok {
		return false
	}
	end, ok := FieldByPath(parent, strings.TrimSpace(parts[1]))
	if !ok {
		return false
	}
	return CompareValue(start, end, "lt")
}

// FieldContains 判断当前字段字符串是否包含目标字段的值
func FieldContains(field reflect.Value, parent reflect.Value, param string) bool {
	other, ok := FieldByPath(parent, param)
	if !ok {
		return false
	}
	left, ok := validate.StringValueFromField(field)
	if !ok {
		return false
	}
	right, ok := validate.ScalarString(other)
	return ok && strings.Contains(left, right)
}

// directField 直接根据字段名获取字段值，支持嵌套字段访问
func directField(current reflect.Value, name string) (reflect.Value, bool) {
	if idx, ok := fieldLookup(current.Type())[name]; ok {
		return current.Field(idx), true
	}
	return reflect.Value{}, false
}

// fieldLookup 构建字段名到索引的映射，缓存结构体字段查找表
func fieldLookup(typ reflect.Type) map[string]int {
	if cached, ok := fieldLookupCache.Load(typ); ok {
		return cached.(map[string]int)
	}
	lookup := make(map[string]int, typ.NumField()*4)
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" {
			continue
		}
		addFieldLookupName(lookup, sf.Name, i)
		if jsonTag := sf.Tag.Get("json"); jsonTag != "" {
			if jsonName, _, _ := strings.Cut(jsonTag, ","); jsonName != "" && jsonName != "-" {
				addFieldLookupName(lookup, jsonName, i)
			}
		}
		addFieldLookupName(lookup, utils.LowerCamel(sf.Name), i)
		addFieldLookupName(lookup, utils.SnakeCase(sf.Name), i)
	}
	actual, _ := fieldLookupCache.LoadOrStore(typ, lookup)
	return actual.(map[string]int)
}

// addFieldLookupName 添加字段名到索引的映射，避免重复添加
func addFieldLookupName(lookup map[string]int, name string, index int) {
	if name == "" {
		return
	}
	if _, exists := lookup[name]; !exists {
		lookup[name] = index
	}
}

// fieldNames 生成字段名的候选列表，包括 Go 字段名、json 名、snake_case
func fieldNames(sf reflect.StructField) []string {
	names := []string{sf.Name, utils.LowerCamel(sf.Name), utils.SnakeCase(sf.Name)}
	if jsonName := strings.Split(sf.Tag.Get("json"), ",")[0]; jsonName != "" && jsonName != "-" {
		names = append(names, jsonName)
	}
	return names
}

// scalarString 获取字段值的字符串表示，支持指针和空值
func scalarString(v reflect.Value) string {
	v = validate.DerefReflect(v)
	if !v.IsValid() {
		return ""
	}
	return validate.StringValue(v)
}

// floatValue 获取字段值的浮点数表示，支持指针和空值
func floatValue(v reflect.Value) (float64, bool) {
	v = validate.DerefReflect(v)
	if !v.IsValid() {
		return 0, false
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	default:
		return 0, false
	}
}

// compareFloat 对比两个浮点数是否符合指定操作符
func compareFloat(left, right float64, op string) bool {
	switch op {
	case "eq":
		return left == right
	case "ne":
		return left != right
	case "gt":
		return left > right
	case "gte":
		return left >= right
	case "lt":
		return left < right
	case "lte":
		return left <= right
	default:
		return false
	}
}
