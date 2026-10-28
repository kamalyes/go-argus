/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-28 15:26:53
 * @FilePath: \go-argus\errors.go
 * @Description: 校验错误模型，提供字段错误、错误集合和兼容错误格式
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const fieldErrMsg = "Key: '%s' Error:Field validation for '%s' failed on the '%s' tag"

// fieldErrMsg 的字面分段，供手工拼接使用（见 fieldError.Error）
//
// 拼接结果必须与 fmt.Sprintf(fieldErrMsg, ns, field, tag) 逐字节一致，
// 该不变量由 TestFieldErrorMsg_MatchesSprintf 守护；改动此处须同步 fieldErrMsg。
const (
	fieldErrPrefix = "Key: '"
	fieldErrMid1   = "' Error:Field validation for '"
	fieldErrMid2   = "' failed on the '"
	fieldErrSuffix = "' tag"
)

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
	switch len(ve) {
	case 0:
		return ""
	case 1:
		// 单条错误是最常见情形：直接返回，省掉缓冲区与换行拼接
		return strings.TrimSpace(ve[0].Error())
	}

	// strings.Builder 而非 bytes.Buffer：后者的 String() 会把累积内容整体复制一份
	// （实测占本路径分配对象的 11.15%），Builder.String() 则零拷贝移交底层数组。
	// 不预先汇总长度：那需要把每条 Error() 先物化到切片（多一次分配），
	// 或对每条调用两次 Error()（成本翻倍）——Builder 自身的倍增增长已足够。
	var b strings.Builder
	for i := 0; i < len(ve); i++ {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(ve[i].Error())
	}
	return strings.TrimSpace(b.String())
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
//
// 无独立字段：别名在编译期已解析为规范名，生产构造路径（newFieldError）
// 恒有 actualTag == tag，省去每错误 16B 与一个 GC 指针扫描槽
func (fe *fieldError) ActualTag() string {
	return fe.tag
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
//
// 手工拼接而非 fmt.Sprintf：fieldErrMsg 全为 %s 代换，无需格式串解析，
// 也省掉可变参数的 []any 装箱。按最终长度精确 Grow，一次分配拼出结果，
// 且 Builder.String() 零拷贝。输出与 Sprintf 逐字节一致。
func (fe *fieldError) Error() string {
	var b strings.Builder
	b.Grow(len(fieldErrPrefix) + len(fe.ns) + len(fieldErrMid1) + len(fe.field) +
		len(fieldErrMid2) + len(fe.tag) + len(fieldErrSuffix))
	b.WriteString(fieldErrPrefix)
	b.WriteString(fe.ns)
	b.WriteString(fieldErrMid1)
	b.WriteString(fe.field)
	b.WriteString(fieldErrMid2)
	b.WriteString(fe.tag)
	b.WriteString(fieldErrSuffix)
	return b.String()
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
