/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\errors.go
 * @Description: 校验错误模型，提供字段错误、错误集合和兼容错误格式
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

const fieldErrMsg = "Key: '%s' Error:Field validation for '%s' failed on the '%s' tag"

// InvalidValidationError 表示传入校验器的参数类型不合法
type InvalidValidationError struct {
	Type reflect.Type
}

func (e *InvalidValidationError) Error() string {
	if e.Type == nil {
		return "validator: (nil)"
	}
	return "validator: (nil " + e.Type.String() + ")"
}

// ValidationErrors 表示字段校验失败集合
type ValidationErrors []FieldError

func (ve ValidationErrors) Error() string {
	var buff bytes.Buffer
	for i := 0; i < len(ve); i++ {
		if i > 0 {
			buff.WriteByte('\n')
		}
		buff.WriteString(ve[i].Error())
	}
	return strings.TrimSpace(buff.String())
}

// FieldError 表示单个字段校验失败详情
type FieldError interface {
	Tag() string
	ActualTag() string
	Namespace() string
	StructNamespace() string
	Field() string
	StructField() string
	Value() any
	Param() string
	Kind() reflect.Kind
	Type() reflect.Type
	Error() string
}

var (
	_ FieldError = (*fieldError)(nil)
	_ error      = (*fieldError)(nil)
)

type fieldError struct {
	tag         string
	actualTag   string
	ns          string
	structNs    string
	field       string
	structField string
	fieldVal    reflect.Value // 惰性求值源，Value/Kind/Type 访问时才解引用
	value       any           // Value() 首次访问装箱缓存
	param       string
}

// Tag 返回失败的规则名
func (fe *fieldError) Tag() string {
	return fe.tag
}

// ActualTag 返回实际失败的规则名，兼容别名规则场景
func (fe *fieldError) ActualTag() string {
	return fe.actualTag
}

// Namespace 返回字段展示路径
func (fe *fieldError) Namespace() string {
	return fe.ns
}

// StructNamespace 返回字段结构体路径
func (fe *fieldError) StructNamespace() string {
	return fe.structNs
}

// Field 返回字段展示名
func (fe *fieldError) Field() string {
	return fe.field
}

// StructField 返回结构体字段名
func (fe *fieldError) StructField() string {
	return fe.structField
}

// Value 返回字段原始值，首次访问时从 reflect.Value 装箱并缓存
func (fe *fieldError) Value() any {
	if fe.value == nil && fe.fieldVal.IsValid() && fe.fieldVal.CanInterface() {
		fe.value = fe.fieldVal.Interface()
	}
	return fe.value
}

// Param 返回规则参数
func (fe *fieldError) Param() string {
	return fe.param
}

// Kind 返回字段反射 Kind，无效值返回 reflect.Invalid
func (fe *fieldError) Kind() reflect.Kind {
	if !fe.fieldVal.IsValid() {
		return reflect.Invalid
	}
	return fe.fieldVal.Kind()
}

// Type 返回字段反射 Type，无效值返回 nil
func (fe *fieldError) Type() reflect.Type {
	if !fe.fieldVal.IsValid() {
		return nil
	}
	return fe.fieldVal.Type()
}

// Error 返回兼容 go-playground 的错误字符串
func (fe *fieldError) Error() string {
	return fmt.Sprintf(fieldErrMsg, fe.ns, fe.field, fe.tag)
}

var _ FieldError = (*stringFieldError)(nil)

type stringFieldError struct {
	tag   string
	param string
	value string
}

func (e *stringFieldError) Tag() string             { return e.tag }
func (e *stringFieldError) ActualTag() string       { return e.tag }
func (e *stringFieldError) Namespace() string       { return "" }
func (e *stringFieldError) StructNamespace() string { return "" }
func (e *stringFieldError) Field() string           { return "" }
func (e *stringFieldError) StructField() string     { return "" }
func (e *stringFieldError) Value() any              { return e.value }
func (e *stringFieldError) Param() string           { return e.param }
func (e *stringFieldError) Kind() reflect.Kind      { return reflect.String }
func (e *stringFieldError) Type() reflect.Type      { return reflect.TypeOf("") }
func (e *stringFieldError) Error() string {
	return fmt.Sprintf(fieldErrMsg, "", "", e.tag)
}
