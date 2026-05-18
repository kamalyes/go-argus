/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
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

func directField(current reflect.Value, name string) (reflect.Value, bool) {
	if idx, ok := fieldLookup(current.Type())[name]; ok {
		return current.Field(idx), true
	}
	return reflect.Value{}, false
}

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
		addFieldLookupName(lookup, lowerCamel(sf.Name), i)
		addFieldLookupName(lookup, snakeCase(sf.Name), i)
	}
	actual, _ := fieldLookupCache.LoadOrStore(typ, lookup)
	return actual.(map[string]int)
}

func addFieldLookupName(lookup map[string]int, name string, index int) {
	if name == "" {
		return
	}
	if _, exists := lookup[name]; !exists {
		lookup[name] = index
	}
}

func fieldNames(sf reflect.StructField) []string {
	names := []string{sf.Name, lowerCamel(sf.Name), snakeCase(sf.Name)}
	if jsonName := strings.Split(sf.Tag.Get("json"), ",")[0]; jsonName != "" && jsonName != "-" {
		names = append(names, jsonName)
	}
	return names
}

func scalarString(v reflect.Value) string {
	v = validate.DerefReflect(v)
	if !v.IsValid() {
		return ""
	}
	return validate.StringValue(v)
}

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

func lowerCamel(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func snakeCase(s string) string {
	var out strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out.WriteByte('_')
		}
		out.WriteRune(r)
	}
	return strings.ToLower(out.String())
}
