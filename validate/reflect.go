/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-19 00:00:00
 * @FilePath: \go-argus\validate\reflect.go
 * @Description: 反射值提取工具，提供字符串、字节、数值和 rune 匹配等零分配辅助函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"fmt"
	"reflect"

	"github.com/kamalyes/go-argus/utils"
)

func ParseFloat(s string) (float64, bool) {
	return utils.ParseFloat(s)
}

func NumericValue(field reflect.Value) (float64, bool) {
	field = DerefReflect(field)
	if !field.IsValid() {
		return 0, false
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(field.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(field.Uint()), true
	case reflect.Float32, reflect.Float64:
		return field.Float(), true
	case reflect.String:
		return utils.ParseFloat(field.String())
	default:
		return 0, false
	}
}

func StringValueFromField(field reflect.Value) (string, bool) {
	field = DerefReflect(field)
	if !field.IsValid() || field.Kind() != reflect.String {
		return "", false
	}
	return field.String(), true
}

func BytesValue(field reflect.Value) ([]byte, bool) {
	field = DerefReflect(field)
	if !field.IsValid() {
		return nil, false
	}
	if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Uint8 {
		return field.Bytes(), true
	}
	return nil, false
}

func ScalarString(field reflect.Value) (string, bool) {
	field = DerefReflect(field)
	if !field.IsValid() {
		return "", false
	}
	if field.Kind() == reflect.String {
		return field.String(), true
	}
	if field.CanInterface() {
		return fmt.Sprint(field.Interface()), true
	}
	return "", false
}

func MatchStringRunes(field reflect.Value, fn func(rune) bool) bool {
	s, ok := StringValueFromField(field)
	if !ok || s == "" {
		return false
	}
	for _, r := range s {
		if !fn(r) {
			return false
		}
	}
	return true
}
