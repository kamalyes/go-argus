/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\field_level.go
 * @Description: 自定义字段校验上下文，提供当前字段、父结构和标签参数访问能力
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"context"
	"reflect"
)

// Func 表示自定义字段校验函数
type Func func(fl FieldLevel) bool

// FuncCtx 表示带 context 的自定义字段校验函数
type FuncCtx func(ctx context.Context, fl FieldLevel) bool

func wrapFunc(fn Func) FuncCtx {
	if fn == nil {
		return nil
	}
	return func(_ context.Context, fl FieldLevel) bool {
		return fn(fl)
	}
}

// FieldLevel 为自定义规则提供当前字段、父结构和标签参数
type FieldLevel interface {
	Top() reflect.Value
	Parent() reflect.Value
	Field() reflect.Value
	FieldName() string
	StructFieldName() string
	GetTag() string
	Param() string
}

type fieldLevel struct {
	top             reflect.Value
	parent          reflect.Value
	field           reflect.Value
	fieldName       string
	structFieldName string
	tag             string
	param           string
}

func (fl fieldLevel) Top() reflect.Value {
	return fl.top
}

func (fl fieldLevel) Parent() reflect.Value {
	return fl.parent
}

func (fl fieldLevel) Field() reflect.Value {
	return fl.field
}

func (fl fieldLevel) FieldName() string {
	return fl.fieldName
}

func (fl fieldLevel) StructFieldName() string {
	return fl.structFieldName
}

func (fl fieldLevel) GetTag() string {
	return fl.tag
}

func (fl fieldLevel) Param() string {
	return fl.param
}
