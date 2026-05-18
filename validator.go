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

var bgCtx = context.Background()

var errsPool = sync.Pool{
	New: func() interface{} {
		s := make(ValidationErrors, 0, 8)
		return &s
	},
}

var fieldLevelPool = sync.Pool{
	New: func() interface{} {
		return &fieldLevel{}
	},
}

func acquireErrors() *ValidationErrors {
	return errsPool.Get().(*ValidationErrors)
}

func releaseErrors(errs *ValidationErrors) {
	*errs = (*errs)[:0]
	errsPool.Put(errs)
}

func acquireFieldLevel() *fieldLevel {
	return fieldLevelPool.Get().(*fieldLevel)
}

func releaseFieldLevel(fl *fieldLevel) {
	fl.top = reflect.Value{}
	fl.parent = reflect.Value{}
	fl.field = reflect.Value{}
	fl.fieldName = ""
	fl.structFieldName = ""
	fl.tag = ""
	fl.param = ""
	fieldLevelPool.Put(fl)
}

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
	varCache    sync.Map
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
	return v.StructCtx(bgCtx, s)
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

	errs := acquireErrors()
	v.validateStruct(ctx, current, current, current.Type().Name(), current.Type().Name(), errs)
	if len(*errs) > 0 {
		result := make(ValidationErrors, len(*errs))
		copy(result, *errs)
		releaseErrors(errs)
		return result
	}
	releaseErrors(errs)
	return nil
}

// Var 按标签表达式校验单个变量
func (v *Validate) Var(field interface{}, tag string) error {
	return v.VarCtx(bgCtx, field, tag)
}

// VarString 按标签表达式校验字符串变量，零分配快速路径
func (v *Validate) VarString(field string, tag string) error {
	return v.VarStringCtx(bgCtx, field, tag)
}

// VarCtx 按标签表达式校验单个变量，并传递 context
func (v *Validate) VarCtx(ctx context.Context, field interface{}, tag string) error {
	rules := v.cachedVarRules(tag)
	if s, ok := field.(string); ok {
		return v.varStringRules(ctx, s, rules, true)
	}
	rv := reflect.ValueOf(field)
	errs := acquireErrors()
	v.applyRules(ctx, reflect.Value{}, rv, rv, "", "", "", "", rules, errs)
	if len(*errs) > 0 {
		result := make(ValidationErrors, len(*errs))
		copy(result, *errs)
		releaseErrors(errs)
		return result
	}
	releaseErrors(errs)
	return nil
}

// VarStringCtx 按标签表达式校验字符串变量，零反射快速路径
func (v *Validate) VarStringCtx(ctx context.Context, field string, tag string) error {
	return v.varStringRules(ctx, field, v.cachedVarRules(tag), false)
}

func (v *Validate) varStringRules(ctx context.Context, field string, rules []rulePlan, wrapError bool) error {
	for i := 0; i < len(rules); i++ {
		r := rules[i]
		switch r.name {
		case "omitempty", "omitzero":
			if isBlankString(field) {
				return nil
			}
			continue
		case "omitnil":
			continue
		case "structonly", "nostructlevel", "":
			continue
		}
		if len(r.orRules) > 0 {
			ok, handled := v.evalStringOr(ctx, field, r.orRules)
			if !handled {
				return v.varStringReflectPath(ctx, field, rules)
			}
			if !ok {
				return v.stringRuleError(field, r, wrapError)
			}
			continue
		}
		ok, handled := evalStringRule(field, r)
		if handled {
			if !ok {
				return v.stringRuleError(field, r, wrapError)
			}
			continue
		}
		return v.varStringReflectPath(ctx, field, rules)
	}
	return nil
}

func (v *Validate) evalStringOr(ctx context.Context, field string, rules []rulePlan) (bool, bool) {
	for i := 0; i < len(rules); i++ {
		ok, handled := evalStringRule(field, rules[i])
		if !handled {
			return false, false
		}
		if ok {
			return true, true
		}
	}
	return false, true
}

func evalStringRule(field string, r rulePlan) (bool, bool) {
	if fn, ok := stringRuleMap[r.name]; ok {
		return fn == nil || fn(field, r.param), true
	}
	switch r.name {
	case "oneof":
		return stringOneOf(field, r.paramParts), true
	case "oneofci":
		return stringOneOfCI(field, r.paramParts), true
	case "noneof":
		return !stringOneOf(field, r.paramParts), true
	case "noneofci":
		return !stringOneOfCI(field, r.paramParts), true
	default:
		return false, false
	}
}

func (v *Validate) stringRuleError(field string, rule rulePlan, wrap bool) error {
	fe := &stringFieldError{tag: rule.name, param: rule.param, value: field}
	if wrap {
		return ValidationErrors{fe}
	}
	return fe
}

func (v *Validate) varStringReflectPath(ctx context.Context, field string, rules []rulePlan) error {
	rv := reflect.ValueOf(field)
	errs := acquireErrors()
	v.applyRules(ctx, reflect.Value{}, rv, rv, "", "", "", "", rules, errs)
	if len(*errs) > 0 {
		result := make(ValidationErrors, len(*errs))
		copy(result, *errs)
		releaseErrors(errs)
		return result
	}
	releaseErrors(errs)
	return nil
}

func stringOneOf(s string, parts []string) bool {
	for _, item := range parts {
		if s == item {
			return true
		}
	}
	return false
}

func stringOneOfCI(s string, parts []string) bool {
	for _, item := range parts {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func (v *Validate) cachedVarRules(tag string) []rulePlan {
	if cached, ok := v.varCache.Load(tag); ok {
		return cached.([]rulePlan)
	}
	rules := parseRules(tag)
	actual, _ := v.varCache.LoadOrStore(tag, rules)
	return actual.([]rulePlan)
}

func (v *Validate) validateStruct(ctx context.Context, top reflect.Value, current reflect.Value, ns string, structNs string, errs *ValidationErrors) {
	current = derefValue(current)
	if !current.IsValid() || current.Kind() != reflect.Struct {
		return
	}

	plan := v.compileStruct(current.Type())
	for _, fp := range plan.fields {
		field := current.FieldByIndex(fp.index)
		fieldNS := fp.nsPrefix
		fieldStructNS := fp.structNsPrefix
		if fieldNS == "" {
			fieldNS = joinNS(ns, fp.altName)
		}
		if fieldStructNS == "" {
			fieldStructNS = joinNS(structNs, fp.name)
		}

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
	derefed := derefValue(field)
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		switch rule.name {
		case "omitempty", "omitzero":
			if isEmptyValue(derefed, v.requiredStructEnabled) {
				return
			}
			continue
		case "omitnil":
			if isNilValue(field) {
				return
			}
			continue
		case "dive":
			v.applyDive(ctx, top, parent, derefed, ns, structNs, fieldName, structFieldName, rules[i+1:], errs)
			return
		case "structonly", "nostructlevel":
			continue
		case "":
			continue
		}

		ok := false
		if len(rule.orRules) > 0 {
			for j := 0; j < len(rule.orRules); j++ {
				if v.evalRule(ctx, top, parent, derefed, fieldName, structFieldName, rule.orRules[j]) {
					ok = true
					break
				}
			}
		} else {
			ok = v.evalRule(ctx, top, parent, derefed, fieldName, structFieldName, rule)
		}
		if !ok {
			*errs = append(*errs, newFieldError(derefed, ns, structNs, fieldName, structFieldName, rule))
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
	if disp, ok := evalDispatchTable[plan.name]; ok {
		return disp(v, top, parent, field, plan)
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

	fl := acquireFieldLevel()
	fl.top = top
	fl.parent = parent
	fl.field = derefValue(field)
	fl.fieldName = fieldName
	fl.structFieldName = structFieldName
	fl.tag = plan.name
	fl.param = plan.param
	result := fn(ctx, fl)
	releaseFieldLevel(fl)
	return result
}

type evalDispatchFn func(v *Validate, top, parent, field reflect.Value, plan rulePlan) bool

func (v *Validate) evalRequiredIf(top, parent, field reflect.Value, plan rulePlan) bool {
	if !rule.IsRequiredIfFast(parent, plan.paramParts) {
		return true
	}
	return !isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredUnless(top, parent, field reflect.Value, plan rulePlan) bool {
	if rule.IsRequiredIfFast(parent, plan.paramParts) {
		return true
	}
	return !isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWith(top, parent, field reflect.Value, plan rulePlan) bool {
	if !rule.IsRequiredWith(parent, plan.param) {
		return true
	}
	return !isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWithAll(top, parent, field reflect.Value, plan rulePlan) bool {
	if !ruleRequiredWithAllFast(parent, plan.paramParts) {
		return true
	}
	return !isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWithout(top, parent, field reflect.Value, plan rulePlan) bool {
	if !ruleRequiredWithoutFast(parent, plan.paramParts) {
		return true
	}
	return !isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWithoutAll(top, parent, field reflect.Value, plan rulePlan) bool {
	if !ruleRequiredWithoutAllFast(parent, plan.paramParts) {
		return true
	}
	return !isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedIf(top, parent, field reflect.Value, plan rulePlan) bool {
	if !rule.IsRequiredIfFast(parent, plan.paramParts) {
		return true
	}
	return isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedUnless(top, parent, field reflect.Value, plan rulePlan) bool {
	if rule.IsRequiredIfFast(parent, plan.paramParts) {
		return true
	}
	return isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWith(top, parent, field reflect.Value, plan rulePlan) bool {
	if !rule.IsRequiredWith(parent, plan.param) {
		return true
	}
	return isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWithAll(top, parent, field reflect.Value, plan rulePlan) bool {
	if !ruleRequiredWithAllFast(parent, plan.paramParts) {
		return true
	}
	return isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWithout(top, parent, field reflect.Value, plan rulePlan) bool {
	if !ruleRequiredWithoutFast(parent, plan.paramParts) {
		return true
	}
	return isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWithoutAll(top, parent, field reflect.Value, plan rulePlan) bool {
	if !ruleRequiredWithoutAllFast(parent, plan.paramParts) {
		return true
	}
	return isEmptyValue(field, v.requiredStructEnabled)
}

func (v *Validate) evalCmpField(top, parent, field reflect.Value, plan rulePlan) bool {
	op := cmpFieldOps[plan.name]
	target := parent
	if strings.HasSuffix(plan.name, "csfield") {
		target = top
	}
	return rule.CompareField(field, target, plan.param, op)
}

func (v *Validate) evalFieldContains(top, parent, field reflect.Value, plan rulePlan) bool {
	return ruleFieldContains(field, parent, plan.param)
}

func (v *Validate) evalFieldExcludes(top, parent, field reflect.Value, plan rulePlan) bool {
	return !ruleFieldContains(field, parent, plan.param)
}

func (v *Validate) evalAfter(top, parent, field reflect.Value, plan rulePlan) bool {
	return rule.CompareTimeExpr(field, plan.param, "gt", time.Now())
}

func (v *Validate) evalBefore(top, parent, field reflect.Value, plan rulePlan) bool {
	return rule.CompareTimeExpr(field, plan.param, "lt", time.Now())
}

func (v *Validate) evalRange(top, parent, field reflect.Value, plan rulePlan) bool {
	return ruleRange(parent, plan.param)
}

func (v *Validate) evalOneOf(top, parent, field reflect.Value, plan rulePlan) bool {
	return ruleOneOfFast(field, plan.paramParts)
}

func (v *Validate) evalOneOfCI(top, parent, field reflect.Value, plan rulePlan) bool {
	return ruleOneOfCIFast(field, plan.paramParts)
}

func (v *Validate) evalNoneOf(top, parent, field reflect.Value, plan rulePlan) bool {
	return !ruleOneOfFast(field, plan.paramParts)
}

func (v *Validate) evalNoneOfCI(top, parent, field reflect.Value, plan rulePlan) bool {
	return !ruleOneOfCIFast(field, plan.paramParts)
}

var cmpFieldOps = map[string]string{
	"eqfield":     "eq",
	"nefield":     "ne",
	"gtfield":     "gt",
	"afterfield":  "gt",
	"gtefield":    "gte",
	"ltfield":     "lt",
	"beforefield": "lt",
	"ltefield":    "lte",
	"eqcsfield":   "eq",
	"necsfield":   "ne",
	"gtcsfield":   "gt",
	"gtecsfield":  "gte",
	"ltcsfield":   "lt",
	"ltecsfield":  "lte",
}

var evalDispatchTable = map[string]evalDispatchFn{
	"required_if":          (*Validate).evalRequiredIf,
	"required_unless":      (*Validate).evalRequiredUnless,
	"required_with":        (*Validate).evalRequiredWith,
	"required_with_all":    (*Validate).evalRequiredWithAll,
	"required_without":     (*Validate).evalRequiredWithout,
	"required_without_all": (*Validate).evalRequiredWithoutAll,
	"excluded_if":          (*Validate).evalExcludedIf,
	"excluded_unless":      (*Validate).evalExcludedUnless,
	"excluded_with":        (*Validate).evalExcludedWith,
	"excluded_with_all":    (*Validate).evalExcludedWithAll,
	"excluded_without":     (*Validate).evalExcludedWithout,
	"excluded_without_all": (*Validate).evalExcludedWithoutAll,
	"eqfield":              (*Validate).evalCmpField,
	"nefield":              (*Validate).evalCmpField,
	"gtfield":              (*Validate).evalCmpField,
	"afterfield":           (*Validate).evalCmpField,
	"gtefield":             (*Validate).evalCmpField,
	"ltfield":              (*Validate).evalCmpField,
	"beforefield":          (*Validate).evalCmpField,
	"ltefield":             (*Validate).evalCmpField,
	"eqcsfield":            (*Validate).evalCmpField,
	"necsfield":            (*Validate).evalCmpField,
	"gtcsfield":            (*Validate).evalCmpField,
	"gtecsfield":           (*Validate).evalCmpField,
	"ltcsfield":            (*Validate).evalCmpField,
	"ltecsfield":           (*Validate).evalCmpField,
	"fieldcontains":        (*Validate).evalFieldContains,
	"fieldexcludes":        (*Validate).evalFieldExcludes,
	"after":                (*Validate).evalAfter,
	"before":               (*Validate).evalBefore,
	"range":                (*Validate).evalRange,
	"oneof":                (*Validate).evalOneOf,
	"oneofci":              (*Validate).evalOneOfCI,
	"noneof":               (*Validate).evalNoneOf,
	"noneofci":             (*Validate).evalNoneOfCI,
}

func ruleOneOfFast(field reflect.Value, parts []string) bool {
	actual, ok := scalarString(field)
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

func ruleOneOfCIFast(field reflect.Value, parts []string) bool {
	actual, ok := scalarString(field)
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

func ruleRequiredWithAllFast(parent reflect.Value, parts []string) bool {
	if len(parts) == 0 {
		return false
	}
	for _, field := range parts {
		value, ok := rule.FieldByPath(parent, field)
		if !ok || isEmptyValue(value, true) {
			return false
		}
	}
	return true
}

func ruleRequiredWithoutFast(parent reflect.Value, parts []string) bool {
	for _, field := range parts {
		value, ok := rule.FieldByPath(parent, field)
		if !ok || isEmptyValue(value, true) {
			return true
		}
	}
	return false
}

func ruleRequiredWithoutAllFast(parent reflect.Value, parts []string) bool {
	if len(parts) == 0 {
		return false
	}
	for _, field := range parts {
		value, ok := rule.FieldByPath(parent, field)
		if ok && !isEmptyValue(value, true) {
			return false
		}
	}
	return true
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
