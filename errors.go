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
	Value() interface{}
	Param() string
	Kind() reflect.Kind
	Type() reflect.Type
	Error() string
}

var _ FieldError = (*fieldError)(nil)
var _ error = (*fieldError)(nil)

type fieldError struct {
	tag         string
	actualTag   string
	ns          string
	structNs    string
	field       string
	structField string
	value       interface{}
	param       string
	kind        reflect.Kind
	typ         reflect.Type
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

// Value 返回字段原始值
func (fe *fieldError) Value() interface{} {
	return fe.value
}

// Param 返回规则参数
func (fe *fieldError) Param() string {
	return fe.param
}

// Kind 返回字段反射 Kind
func (fe *fieldError) Kind() reflect.Kind {
	return fe.kind
}

// Type 返回字段反射 Type
func (fe *fieldError) Type() reflect.Type {
	return fe.typ
}

// Error 返回兼容 go-playground 的错误字符串
func (fe *fieldError) Error() string {
	return fmt.Sprintf(fieldErrMsg, fe.ns, fe.field, fe.tag)
}
