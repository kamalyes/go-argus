/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validator.go
 * @Description: Argus 根校验器，提供 struct tag 校验、变量校验、自定义规则和兼容入口
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kamalyes/go-argus/rule"
)

const defaultTagName = "validate"

// TagNameFunc 用于解析字段展示名，例如从 json tag 中取字段名
type TagNameFunc func(field reflect.StructField) string

// Validate 是可复用且并发安全的校验器实例
type Validate struct {
	tagName                string
	requiredStructEnabled  bool
	privateFieldValidation bool
	tagNameFunc            TagNameFunc

	mu          sync.RWMutex
	validations map[string]FuncCtx
	structCache sync.Map
}

// New 创建校验器实例
func New(options ...Option) *Validate {
	v := &Validate{
		tagName:     defaultTagName,
		validations: make(map[string]FuncCtx),
	}
	for _, opt := range options {
		if opt != nil {
			opt(v)
		}
	}
	return v
}

// SetTagName 设置校验标签名，默认使用 validate
func (v *Validate) SetTagName(name string) {
	if strings.TrimSpace(name) == "" {
		return
	}
	v.tagName = name
	v.structCache = sync.Map{}
}

// RegisterTagNameFunc 注册字段展示名解析函数
func (v *Validate) RegisterTagNameFunc(fn TagNameFunc) {
	v.tagNameFunc = fn
	v.structCache = sync.Map{}
}

// RegisterValidation 注册自定义字段校验函数
func (v *Validate) RegisterValidation(tag string, fn Func, _ ...bool) error {
	return v.RegisterValidationCtx(tag, wrapFunc(fn), false)
}

// RegisterValidationCtx 注册带 context 的自定义字段校验函数
func (v *Validate) RegisterValidationCtx(tag string, fn FuncCtx, _ ...bool) error {
	tag = strings.TrimSpace(tag)
	if tag == "" || fn == nil {
		return &InvalidValidationError{}
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	v.validations[tag] = fn
	return nil
}

// Struct 根据结构体字段上的 validate 标签执行校验
func (v *Validate) Struct(s interface{}) error {
	return v.StructCtx(context.Background(), s)
}

// StructCtx 根据结构体字段上的 validate 标签执行校验，并传递 context
func (v *Validate) StructCtx(ctx context.Context, s interface{}) error {
	current := reflect.ValueOf(s)
	if !current.IsValid() {
		return &InvalidValidationError{}
	}
	current = derefValue(current)
	if !current.IsValid() || current.Kind() != reflect.Struct {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	errs := make(ValidationErrors, 0)
	v.validateStruct(ctx, current, current, current.Type().Name(), current.Type().Name(), &errs)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Var 按标签表达式校验单个变量
func (v *Validate) Var(field interface{}, tag string) error {
	return v.VarCtx(context.Background(), field, tag)
}

// VarCtx 按标签表达式校验单个变量，并传递 context
func (v *Validate) VarCtx(ctx context.Context, field interface{}, tag string) error {
	rv := reflect.ValueOf(field)
	rules := parseRules(tag)
	errs := make(ValidationErrors, 0)
	v.applyRules(ctx, reflect.Value{}, rv, rv, "", "", "", "", rules, &errs)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (v *Validate) validateStruct(ctx context.Context, top reflect.Value, current reflect.Value, ns string, structNs string, errs *ValidationErrors) {
	current = derefValue(current)
	if !current.IsValid() || current.Kind() != reflect.Struct {
		return
	}

	plan := v.compileStruct(current.Type())
	for _, fp := range plan.fields {
		field := current.FieldByIndex(fp.index)
		fieldNS := joinNS(ns, fp.altName)
		fieldStructNS := joinNS(structNs, fp.name)

		before := len(*errs)
		v.applyRules(ctx, top, current, field, fieldNS, fieldStructNS, fp.altName, fp.name, fp.rules, errs)
		if len(*errs) != before {
			continue
		}

		if shouldDiveIntoStruct(field, fp.rules) {
			nested := derefValue(field)
			if nested.IsValid() && nested.Kind() == reflect.Struct {
				v.validateStruct(ctx, top, nested, fieldNS, fieldStructNS, errs)
			}
		}
	}
}

func (v *Validate) applyRules(ctx context.Context, top reflect.Value, parent reflect.Value, field reflect.Value, ns string, structNs string, fieldName string, structFieldName string, rules []rulePlan, errs *ValidationErrors) {
	if len(rules) == 0 {
		return
	}
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		switch rule.name {
		case "omitempty", "omitzero":
			if isEmptyValue(field, v.requiredStructEnabled) {
				return
			}
			continue
		case "omitnil":
			if isNilValue(field) {
				return
			}
			continue
		case "dive":
			v.applyDive(ctx, top, parent, field, ns, structNs, fieldName, structFieldName, rules[i+1:], errs)
			return
		case "structonly", "nostructlevel":
			continue
		case "":
			continue
		}

		ok := v.evalRule(ctx, top, parent, field, fieldName, structFieldName, rule)
		if !ok {
			*errs = append(*errs, newFieldError(field, ns, structNs, fieldName, structFieldName, rule))
			return
		}
	}
}

func (v *Validate) applyDive(ctx context.Context, top reflect.Value, parent reflect.Value, field reflect.Value, ns string, structNs string, fieldName string, structFieldName string, rules []rulePlan, errs *ValidationErrors) {
	field = derefValue(field)
	if !field.IsValid() {
		return
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < field.Len(); i++ {
			childNS := ns + "[" + strconv.Itoa(i) + "]"
			childStructNS := structNs + "[" + strconv.Itoa(i) + "]"
			v.applyRules(ctx, top, parent, field.Index(i), childNS, childStructNS, fieldName, structFieldName, rules, errs)
		}
	case reflect.Map:
		for _, key := range field.MapKeys() {
			keyText := toStringValue(key)
			childNS := ns + "[" + keyText + "]"
			childStructNS := structNs + "[" + keyText + "]"
			v.applyRules(ctx, top, parent, field.MapIndex(key), childNS, childStructNS, fieldName, structFieldName, rules, errs)
		}
	}
}

func (v *Validate) evalRule(ctx context.Context, top reflect.Value, parent reflect.Value, field reflect.Value, fieldName string, structFieldName string, plan rulePlan) bool {
	switch plan.name {
	case "required_if":
		if !rulepkgRequiredIf(parent, plan.param) {
			return true
		}
		return !isEmptyValue(field, v.requiredStructEnabled)
	case "required_unless":
		if rulepkgRequiredIf(parent, plan.param) {
			return true
		}
		return !isEmptyValue(field, v.requiredStructEnabled)
	case "required_with":
		if !rule.IsRequiredWith(parent, plan.param) {
			return true
		}
		return !isEmptyValue(field, v.requiredStructEnabled)
	case "required_with_all":
		if !ruleRequiredWithAll(parent, plan.param) {
			return true
		}
		return !isEmptyValue(field, v.requiredStructEnabled)
	case "required_without":
		if !ruleRequiredWithout(parent, plan.param) {
			return true
		}
		return !isEmptyValue(field, v.requiredStructEnabled)
	case "required_without_all":
		if !ruleRequiredWithoutAll(parent, plan.param) {
			return true
		}
		return !isEmptyValue(field, v.requiredStructEnabled)
	case "excluded_if":
		if !rulepkgRequiredIf(parent, plan.param) {
			return true
		}
		return isEmptyValue(field, v.requiredStructEnabled)
	case "excluded_unless":
		if rulepkgRequiredIf(parent, plan.param) {
			return true
		}
		return isEmptyValue(field, v.requiredStructEnabled)
	case "excluded_with":
		if !rule.IsRequiredWith(parent, plan.param) {
			return true
		}
		return isEmptyValue(field, v.requiredStructEnabled)
	case "excluded_with_all":
		if !ruleRequiredWithAll(parent, plan.param) {
			return true
		}
		return isEmptyValue(field, v.requiredStructEnabled)
	case "excluded_without":
		if !ruleRequiredWithout(parent, plan.param) {
			return true
		}
		return isEmptyValue(field, v.requiredStructEnabled)
	case "excluded_without_all":
		if !ruleRequiredWithoutAll(parent, plan.param) {
			return true
		}
		return isEmptyValue(field, v.requiredStructEnabled)
	case "eqfield":
		return rule.CompareField(field, parent, plan.param, "eq")
	case "nefield":
		return rule.CompareField(field, parent, plan.param, "ne")
	case "gtfield", "afterfield":
		return rule.CompareField(field, parent, plan.param, "gt")
	case "gtefield":
		return rule.CompareField(field, parent, plan.param, "gte")
	case "ltfield", "beforefield":
		return rule.CompareField(field, parent, plan.param, "lt")
	case "ltefield":
		return rule.CompareField(field, parent, plan.param, "lte")
	case "eqcsfield":
		return rule.CompareField(field, top, plan.param, "eq")
	case "necsfield":
		return rule.CompareField(field, top, plan.param, "ne")
	case "gtcsfield":
		return rule.CompareField(field, top, plan.param, "gt")
	case "gtecsfield":
		return rule.CompareField(field, top, plan.param, "gte")
	case "ltcsfield":
		return rule.CompareField(field, top, plan.param, "lt")
	case "ltecsfield":
		return rule.CompareField(field, top, plan.param, "lte")
	case "fieldcontains":
		return ruleFieldContains(field, parent, plan.param)
	case "fieldexcludes":
		return !ruleFieldContains(field, parent, plan.param)
	case "after":
		return rule.CompareTimeExpr(field, plan.param, "gt", time.Now())
	case "before":
		return rule.CompareTimeExpr(field, plan.param, "lt", time.Now())
	case "range":
		return ruleRange(parent, plan.param)
	}

	if fn, ok := builtinRules[plan.name]; ok {
		return fn(field, plan.param, v.requiredStructEnabled)
	}

	v.mu.RLock()
	fn := v.validations[plan.name]
	v.mu.RUnlock()
	if fn == nil {
		return false
	}

	return fn(ctx, fieldLevel{
		top:             top,
		parent:          parent,
		field:           derefValue(field),
		fieldName:       fieldName,
		structFieldName: structFieldName,
		tag:             plan.name,
		param:           plan.param,
	})
}

func rulepkgRequiredIf(parent reflect.Value, param string) bool {
	return rule.IsRequiredIf(parent, param)
}

func ruleRequiredWithAll(parent reflect.Value, param string) bool {
	for _, field := range strings.Fields(param) {
		value, ok := rule.FieldByPath(parent, field)
		if !ok || isEmptyValue(value, true) {
			return false
		}
	}
	return strings.TrimSpace(param) != ""
}

func ruleRequiredWithout(parent reflect.Value, param string) bool {
	for _, field := range strings.Fields(param) {
		value, ok := rule.FieldByPath(parent, field)
		if !ok || isEmptyValue(value, true) {
			return true
		}
	}
	return false
}

func ruleRequiredWithoutAll(parent reflect.Value, param string) bool {
	for _, field := range strings.Fields(param) {
		value, ok := rule.FieldByPath(parent, field)
		if ok && !isEmptyValue(value, true) {
			return false
		}
	}
	return strings.TrimSpace(param) != ""
}

func ruleRange(parent reflect.Value, param string) bool {
	sep := ","
	if strings.Contains(param, "|") {
		sep = "|"
	}
	parts := strings.Split(param, sep)
	if len(parts) != 2 {
		return false
	}
	start, ok := rule.FieldByPath(parent, strings.TrimSpace(parts[0]))
	if !ok {
		return false
	}
	end, ok := rule.FieldByPath(parent, strings.TrimSpace(parts[1]))
	if !ok {
		return false
	}
	return rule.CompareValue(start, end, "lt")
}

func ruleFieldContains(field reflect.Value, parent reflect.Value, param string) bool {
	other, ok := rule.FieldByPath(parent, param)
	if !ok {
		return false
	}
	left, ok := stringValue(field)
	if !ok {
		return false
	}
	right, ok := scalarString(other)
	return ok && strings.Contains(left, right)
}

func newFieldError(field reflect.Value, ns string, structNs string, fieldName string, structFieldName string, rule rulePlan) FieldError {
	value := interface{}(nil)
	current := derefValue(field)
	if current.IsValid() && current.CanInterface() {
		value = current.Interface()
	}
	kind := reflect.Invalid
	var typ reflect.Type
	if current.IsValid() {
		kind = current.Kind()
		typ = current.Type()
	}
	return &fieldError{
		tag:         rule.name,
		actualTag:   rule.name,
		ns:          ns,
		structNs:    structNs,
		field:       fieldName,
		structField: structFieldName,
		value:       value,
		param:       rule.param,
		kind:        kind,
		typ:         typ,
	}
}

func shouldDiveIntoStruct(field reflect.Value, rules []rulePlan) bool {
	for _, rule := range rules {
		switch rule.name {
		case "dive", "nostructlevel", "structonly":
			return false
		}
	}
	field = derefValue(field)
	return field.IsValid() && field.Kind() == reflect.Struct && !isTimeType(field.Type())
}

func joinNS(parent string, child string) string {
	if parent == "" {
		return child
	}
	if child == "" {
		return parent
	}
	return parent + "." + child
}
