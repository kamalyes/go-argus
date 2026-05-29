/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-20 22:26:07
 * @FilePath: \go-argus\rule\field.go
 * @Description: 字段路径解析模块，支持 Go 字段名、json 名、snake_case 和嵌套字段访问
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"strconv"
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
func IsRequiredIf(parent reflect.Value, parts []string) bool {
	if len(parts) < 2 || len(parts)%2 != 0 {
		return false
	}
	for i := 0; i < len(parts); i += 2 {
		value, ok := FieldByPath(parent, parts[i])
		if !ok || !matchScalarString(value, parts[i+1]) {
			return false
		}
	}
	return true
}

// matchScalarString 判断字段值是否与目标字符串匹配
// 同时支持数字字符串（如 "2"）和枚举名称（如 "TRACKING_TYPE_ADJUST"），
// 以兼容 proto 枚举类型的 required_if 校验
func matchScalarString(v reflect.Value, target string) bool {
	s, _ := validate.ScalarString(v)
	if s == target {
		return true
	}
	v = validate.DerefReflect(v)
	if !v.IsValid() || !v.CanInterface() {
		return false
	}
	// 对整数类型，额外尝试数字字符串匹配（fmt.Sprint 对实现了 String() 的枚举会返回枚举名而非数字）
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if strconv.FormatInt(v.Int(), 10) == target {
			return true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if strconv.FormatUint(v.Uint(), 10) == target {
			return true
		}
	}
	// 额外尝试 String() 方法（proto 枚举实现了 String()）
	if stringer, ok := v.Interface().(interface{ String() string }); ok {
		if stringer.String() == target {
			return true
		}
	}
	return false
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

// CompareFieldDerefed 将已解引用的当前字段和目标字段按操作符比较
// 调用方已对 current 做 DerefReflect，避免重复解引用
func CompareFieldDerefed(current reflect.Value, parent reflect.Value, targetPath string, op string) bool {
	target, ok := FieldByPath(parent, targetPath)
	if !ok {
		return false
	}
	return CompareValue(current, target, op)
}

// FieldIndexByPath 预解析字段路径，供结构体计划复用
func FieldIndexByPath(root reflect.Type, path string) ([]int, bool) {
	if root == nil || path == "" {
		return nil, false
	}
	for root.Kind() == reflect.Ptr {
		root = root.Elem()
	}
	if root.Kind() != reflect.Struct {
		return nil, false
	}
	var index []int
	current := root
	for _, part := range strings.Split(path, ".") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx, ok := fieldLookup(current)[part]
		if !ok {
			return nil, false
		}
		sf := current.Field(idx)
		index = append(index, sf.Index...)
		current = sf.Type
		for current.Kind() == reflect.Ptr {
			current = current.Elem()
		}
	}
	return index, len(index) > 0
}

// CompareValue 按操作符比较两个值，优先支持时间，其次支持数值和字符串
func CompareValue(left reflect.Value, right reflect.Value, op string) bool {
	cmpOp := validate.CmpOpFromStr(op)
	return CompareValueOp(left, right, cmpOp)
}

// CompareValueOp 按预解析操作符比较两个值
func CompareValueOp(left reflect.Value, right reflect.Value, cmpOp validate.CmpOp) bool {
	if cmpOp < validate.CmpEQ || cmpOp > validate.CmpNE {
		return false
	}
	left = validate.DerefReflect(left)
	right = validate.DerefReflect(right)
	if !left.IsValid() || !right.IsValid() {
		return false
	}
	if left.Kind() == reflect.String && right.Kind() == reflect.String {
		switch cmpOp {
		case validate.CmpEQ, validate.CmpNE:
			return validate.CompareStringsOp(left.String(), right.String(), cmpOp)
		}
	}
	if lt, lok := TimeValue(left, ""); lok {
		rt, rok := TimeValue(right, "")
		return rok && validate.CompareOp(float64(lt.UnixNano()), float64(rt.UnixNano()), cmpOp)
	}
	lf, lok := validate.NumericValue(left)
	rf, rok := validate.NumericValue(right)
	if lok && rok {
		return validate.CompareOp(lf, rf, cmpOp)
	}
	// 快速路径：两边都是字符串时直接比较，避免 ScalarString 中的 fmt.Sprint 开销
	if left.Kind() == reflect.String && right.Kind() == reflect.String {
		return validate.CompareStringsOp(left.String(), right.String(), cmpOp)
	}
	ls, _ := validate.ScalarString(left)
	rs, _ := validate.ScalarString(right)
	return validate.CompareStringsOp(ls, rs, cmpOp)
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
