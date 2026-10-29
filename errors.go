/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-28 15:52:08
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

	"github.com/kamalyes/go-argus/rule"
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

// appendFieldErrMsg 按 fieldErrMsg 格式把单条字段错误消息写入 b；
// fieldError.Error 与 ValidationErrors.Error 的多条拼接共用，格式单一来源。
// 拼接结果必须与 fmt.Sprintf(fieldErrMsg, ns, field, tag) 逐字节一致，
// 该不变量由 TestFieldErrorMsg_MatchesSprintf 与 TestValidationErrorsMultiple 守护。
func appendFieldErrMsg(b *strings.Builder, ns, field, tag string) {
	b.WriteString(fieldErrPrefix)
	b.WriteString(ns)
	b.WriteString(fieldErrMid1)
	b.WriteString(field)
	b.WriteString(fieldErrMid2)
	b.WriteString(tag)
	b.WriteString(fieldErrSuffix)
}

// fieldErrMsgLen 是 appendFieldErrMsg 输出的精确长度，供 Grow 一次分配使用
func fieldErrMsgLen(nsLen, fieldLen, tagLen int) int {
	return len(fieldErrPrefix) + nsLen + len(fieldErrMid1) + fieldLen +
		len(fieldErrMid2) + tagLen + len(fieldErrSuffix)
}

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
	// 总长可从字段部件直接求出（fieldErrMsgLen），无需把每条 Error() 先物化成
	// 临时串再统计——那会按条数线性分配。具体类型走部件直写共享同一个 Builder，
	// 非 fieldError 的实现回退 Error() 物化（罕见路径）。
	total := len(ve) - 1 // 换行符
	for i := 0; i < len(ve); i++ {
		if fe, ok := ve[i].(*fieldError); ok {
			total += fieldErrMsgLen(len(fe.ns), len(fe.fieldName()), len(fe.tagName()))
		} else {
			total += len(ve[i].Error())
		}
	}
	var b strings.Builder
	b.Grow(total)
	for i := 0; i < len(ve); i++ {
		if i > 0 {
			b.WriteByte('\n')
		}
		if fe, ok := ve[i].(*fieldError); ok {
			appendFieldErrMsg(&b, fe.ns, fe.fieldName(), fe.tagName())
		} else {
			b.WriteString(ve[i].Error())
		}
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

// fieldError 是 FieldError 的生产实现
//
// tag/param 与 field/structField 不自存，改为各持一个指向编译期数据的指针：
// RulePlan 来自 StructPlan.Fields[i].Rules（结构体计划）或 cachedVarRules
// （Var 规则缓存），FieldNames 来自 FieldPlan.Names，两者在编译后只读且跨
// goroutine 共享，寿命长于任何错误对象，故按指针引用是安全的。
//
// 这么做是因为失败路径的成本不在校验逻辑上：StructComplexFailure 的 profile
// 里库代码只占 applyRules 8.2% + newFieldError 5.9%，其余是 mallocgc、
// writeHeapBitsSmall、typePointers.next、madvise 这些分配器与 GC 工作。
// 四个串摊开存要 64 字节、4 个指针字，换成两个指针后结构体从 136 字节
// （144B 规格）降到 88 字节（96B 规格），每个错误少 48 字节、少 2 个待扫描指针字。
//
// ns/structNs 必须自存：dive 子元素的命名空间带 [i]/[key] 后缀，由运行期拼出，
// 不在编译期数据里（见 applyDive）。
type fieldError struct {
	names    *rule.FieldNames // 编译期常量，提供 field / structField
	plan     *rule.RulePlan   // 编译期常量，提供 tag / param
	ns       string           // 展示命名空间，dive 会覆盖
	structNs string           // 结构体命名空间，dive 会覆盖
	fieldVal reflect.Value    // 惰性求值源，Value/Kind/Type 访问时才解引用
	value    any              // Value() 首次访问装箱缓存
}

// tagName 取规则名；plan 为 nil（零值构造）时返回空串，不 panic
func (fe *fieldError) tagName() string {
	if fe.plan == nil {
		return ""
	}
	return fe.plan.Name
}

// fieldName 取展示名；names 为 nil（零值构造）时返回空串，不 panic
func (fe *fieldError) fieldName() string {
	if fe.names == nil {
		return ""
	}
	return fe.names.Field
}

// Tag 返回失败的规则名
func (fe *fieldError) Tag() string {
	return fe.tagName()
}

// ActualTag 返回实际失败的规则名，兼容别名规则场景
//
// 无独立字段：别名在编译期已解析为规范名，生产构造路径（newFieldError）
// 恒有 actualTag == tag，省去每错误 16B 与一个 GC 指针扫描槽
func (fe *fieldError) ActualTag() string {
	return fe.tagName()
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
	return fe.fieldName()
}

// StructField 返回结构体字段名
func (fe *fieldError) StructField() string {
	if fe.names == nil {
		return ""
	}
	return fe.names.StructField
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
	if fe.plan == nil {
		return ""
	}
	return fe.plan.Param
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
// 拼接本体收敛到 appendFieldErrMsg，与 ValidationErrors.Error 多条拼接共用。
func (fe *fieldError) Error() string {
	var b strings.Builder
	b.Grow(fieldErrMsgLen(len(fe.ns), len(fe.fieldName()), len(fe.tagName())))
	appendFieldErrMsg(&b, fe.ns, fe.fieldName(), fe.tagName())
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
