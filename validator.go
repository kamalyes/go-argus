/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-19 13:16:11
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
	"github.com/kamalyes/go-argus/utils"
	"github.com/kamalyes/go-argus/validate"
)

// defaultTagName 默认校验标签名
const defaultTagName = "validate"

// bgCtx 默认后台上下文
var bgCtx = context.Background()

// errsPool 校验错误切片对象池
var errsPool = sync.Pool{
	New: func() interface{} {
		s := make(ValidationErrors, 0, 8)
		return &s
	},
}

// fieldLevelPool 字段级别校验上下文对象池
var fieldLevelPool = sync.Pool{
	New: func() interface{} {
		return &fieldLevel{}
	},
}

// acquireErrors 从对象池获取错误切片
func acquireErrors() *ValidationErrors {
	return errsPool.Get().(*ValidationErrors)
}

// releaseErrors 清空并归还错误切片到对象池
func releaseErrors(errs *ValidationErrors) {
	*errs = (*errs)[:0]
	errsPool.Put(errs)
}

// acquireFieldLevel 从对象池获取字段级别校验上下文
func acquireFieldLevel() *fieldLevel {
	return fieldLevelPool.Get().(*fieldLevel)
}

// releaseFieldLevel 重置并归还字段级别校验上下文到对象池
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
	current = validate.DerefReflect(current)
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
	switch val := field.(type) {
	case string:
		return v.varStringRules(ctx, val, rules, true)
	case *string:
		if val == nil {
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
		return v.varStringRules(ctx, *val, rules, true)
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

// varStringRules 零反射快速路径，按规则列表校验字符串
func (v *Validate) varStringRules(ctx context.Context, field string, rules []rule.RulePlan, wrapError bool) error {
	for i := 0; i < len(rules); i++ {
		r := rules[i]
		switch r.Name {
		case "omitempty", "omitzero":
			if validate.IsBlankString(field) {
				return nil
			}
			continue
		case "omitnil":
			continue
		case "structonly", "nostructlevel", "":
			continue
		}
		if len(r.OrRules) > 0 {
			ok, handled := v.evalStringOr(ctx, field, r.OrRules)
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

// evalStringOr 评估字符串或规则，任一通过即返回 true
func (v *Validate) evalStringOr(ctx context.Context, field string, rules []rule.RulePlan) (bool, bool) {
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

// evalStringRule 评估单条字符串规则，返回 (是否通过, 是否已处理)
func evalStringRule(field string, r rule.RulePlan) (bool, bool) {
	if fn, ok := rule.StringRuleMap[r.Name]; ok {
		return fn == nil || fn(field, r.Param), true
	}
	switch r.Name {
	case "oneof":
		return validate.StringOneOf(field, r.ParamParts), true
	case "oneofci":
		return validate.StringOneOfCI(field, r.ParamParts), true
	case "noneof":
		return !validate.StringOneOf(field, r.ParamParts), true
	case "noneofci":
		return !validate.StringOneOfCI(field, r.ParamParts), true
	default:
		return false, false
	}
}

// stringRuleError 构造字符串规则校验失败错误
func (v *Validate) stringRuleError(field string, rule rule.RulePlan, wrap bool) error {
	fe := &stringFieldError{tag: rule.Name, param: rule.Param, value: field}
	if wrap {
		return ValidationErrors{fe}
	}
	return fe
}

// varStringReflectPath 字符串回退到反射路径
func (v *Validate) varStringReflectPath(ctx context.Context, field string, rules []rule.RulePlan) error {
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

// cachedVarRules 缓存变量规则解析结果
func (v *Validate) cachedVarRules(tag string) []rule.RulePlan {
	if cached, ok := v.varCache.Load(tag); ok {
		return cached.([]rule.RulePlan)
	}
	rules := rule.ParseRules(tag)
	actual, _ := v.varCache.LoadOrStore(tag, rules)
	return actual.([]rule.RulePlan)
}

// validateStruct 递归校验结构体字段
func (v *Validate) validateStruct(ctx context.Context, top reflect.Value, current reflect.Value, ns string, structNs string, errs *ValidationErrors) {
	current = validate.DerefReflect(current)
	if !current.IsValid() || current.Kind() != reflect.Struct {
		return
	}

	plan := v.compileStruct(current.Type())
	for _, fp := range plan.Fields {
		field := current.FieldByIndex(fp.Index)
		fieldNS := fp.NsPrefix
		fieldStructNS := fp.StructNsPrefix
		if fieldNS == "" {
			fieldNS = utils.JoinNS(ns, fp.AltName)
		}
		if fieldStructNS == "" {
			fieldStructNS = utils.JoinNS(structNs, fp.Name)
		}

		before := len(*errs)
		v.applyRules(ctx, top, current, field, fieldNS, fieldStructNS, fp.AltName, fp.Name, fp.Rules, errs)
		if len(*errs) != before {
			continue
		}

		if shouldDiveIntoStruct(field, fp.Rules) {
			nested := validate.DerefReflect(field)
			if nested.IsValid() && nested.Kind() == reflect.Struct {
				v.validateStruct(ctx, top, nested, fieldNS, fieldStructNS, errs)
			}
		}
	}
}

// applyRules 按规则列表校验单个字段
func (v *Validate) applyRules(ctx context.Context, top reflect.Value, parent reflect.Value, field reflect.Value, ns string, structNs string, fieldName string, structFieldName string, rules []rule.RulePlan, errs *ValidationErrors) {
	if len(rules) == 0 {
		return
	}
	derefed := validate.DerefReflect(field)
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		switch rule.Name {
		case "omitempty", "omitzero":
			if validate.IsEmptyValueWithStruct(derefed, v.requiredStructEnabled) {
				return
			}
			continue
		case "omitnil":
			if validate.IsNilValue(field) {
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
		if len(rule.OrRules) > 0 {
			for j := 0; j < len(rule.OrRules); j++ {
				orRule := rule.OrRules[j]
				// 快速路径：or 规则通常是 builtin，直接查表避免完整 evalRule 开销
				if action, found := evalTable[orRule.Name]; found {
					if action.dispatch != nil {
						if action.dispatch(v, top, parent, derefed, orRule) {
							ok = true
							break
						}
					} else if action.builtin(derefed, orRule.Param, v.requiredStructEnabled) {
						ok = true
						break
					}
				} else if v.evalRule(ctx, top, parent, derefed, fieldName, structFieldName, orRule) {
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

// applyDive 递归校验切片、数组和 map 元素
func (v *Validate) applyDive(ctx context.Context, top reflect.Value, parent reflect.Value, field reflect.Value, ns string, structNs string, fieldName string, structFieldName string, rules []rule.RulePlan, errs *ValidationErrors) {
	field = validate.DerefReflect(field)
	if !field.IsValid() {
		return
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Array:
		// 复用 []byte 缓冲区避免 strconv.Itoa 分配
		var idxBuf []byte
		for i := 0; i < field.Len(); i++ {
			idxBuf = strconv.AppendInt(idxBuf[:0], int64(i), 10)
			childNS := ns + "[" + string(idxBuf) + "]"
			childStructNS := structNs + "[" + string(idxBuf) + "]"
			v.applyRules(ctx, top, parent, field.Index(i), childNS, childStructNS, fieldName, structFieldName, rules, errs)
		}
	case reflect.Map:
		for _, key := range field.MapKeys() {
			keyText := validate.StringValue(key)
			childNS := ns + "[" + keyText + "]"
			childStructNS := structNs + "[" + keyText + "]"
			v.applyRules(ctx, top, parent, field.MapIndex(key), childNS, childStructNS, fieldName, structFieldName, rules, errs)
		}
	}
}

// evalRule 评估单条规则，优先查 dispatch 表，其次查 builtin 表，最后查自定义注册
func (v *Validate) evalRule(ctx context.Context, top reflect.Value, parent reflect.Value, field reflect.Value, fieldName string, structFieldName string, plan rule.RulePlan) bool {
	if action, ok := evalTable[plan.Name]; ok {
		if action.dispatch != nil {
			return action.dispatch(v, top, parent, field, plan)
		}
		return action.builtin(field, plan.Param, v.requiredStructEnabled)
	}

	// 仅在 builtin/dispatch 未命中时才加读锁查自定义规则
	fn := v.getCustomValidation(plan.Name)
	if fn == nil {
		return false
	}

	fl := acquireFieldLevel()
	fl.top = top
	fl.parent = parent
	fl.field = validate.DerefReflect(field)
	fl.fieldName = fieldName
	fl.structFieldName = structFieldName
	fl.tag = plan.Name
	fl.param = plan.Param
	result := fn(ctx, fl)
	releaseFieldLevel(fl)
	return result
}

// getCustomValidation 读取自定义规则，使用 RLock 保护
func (v *Validate) getCustomValidation(name string) FuncCtx {
	v.mu.RLock()
	fn := v.validations[name]
	v.mu.RUnlock()
	return fn
}

// evalDispatchFn 规则分派函数签名
type evalDispatchFn func(v *Validate, top, parent, field reflect.Value, plan rule.RulePlan) bool

// evalAction 规则执行动作，dispatch 和 builtin 二选一
type evalAction struct {
	dispatch evalDispatchFn
	builtin  rule.BuiltinRule
}

// evalTable 规则执行表，init 中合并 dispatch 表和 builtin 表
var evalTable map[string]evalAction

func init() {
	evalTable = make(map[string]evalAction, len(evalDispatchTable)+len(rule.BuiltinRules))
	for name, fn := range evalDispatchTable {
		evalTable[name] = evalAction{dispatch: fn}
	}
	for name, fn := range rule.BuiltinRules {
		if _, exists := evalTable[name]; !exists {
			evalTable[name] = evalAction{builtin: fn}
		}
	}
}

func (v *Validate) evalRequiredIf(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredIfFast(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredUnless(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if rule.IsRequiredIfFast(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWith(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredWith(parent, plan.Param) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWithAll(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredWithAll(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWithout(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredWithout(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWithoutAll(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredWithoutAll(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedIf(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredIfFast(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedUnless(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if rule.IsRequiredIfFast(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWith(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredWith(parent, plan.Param) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWithAll(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredWithAll(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWithout(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredWithout(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWithoutAll(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if !rule.IsRequiredWithoutAll(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalCmpField(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	op := cmpFieldOps[plan.Name]
	target := parent
	if strings.HasSuffix(plan.Name, "csfield") {
		target = top
	}
	// 预解引用，避免 CompareValue 内部重复 DerefReflect
	derefed := validate.DerefReflect(field)
	return rule.CompareFieldDerefed(derefed, target, plan.Param, op)
}

func (v *Validate) evalFieldContains(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return rule.FieldContains(field, parent, plan.Param)
}

func (v *Validate) evalFieldExcludes(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return !rule.FieldContains(field, parent, plan.Param)
}

func (v *Validate) evalAfter(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return rule.CompareTimeExpr(field, plan.Param, "gt", time.Now())
}

func (v *Validate) evalBefore(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return rule.CompareTimeExpr(field, plan.Param, "lt", time.Now())
}

func (v *Validate) evalRange(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return rule.Range(parent, plan.Param)
}

func (v *Validate) evalOneOf(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return rule.OneOfFast(field, plan.ParamParts)
}

func (v *Validate) evalOneOfCI(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return rule.OneOfCIFast(field, plan.ParamParts)
}

func (v *Validate) evalNoneOf(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return !rule.OneOfFast(field, plan.ParamParts)
}

func (v *Validate) evalNoneOfCI(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return !rule.OneOfCIFast(field, plan.ParamParts)
}

// cmpFieldOps 跨字段比较规则到操作符的映射
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

// evalDispatchTable 需要跨字段访问的规则分派表
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

func newFieldError(field reflect.Value, ns string, structNs string, fieldName string, structFieldName string, rule rule.RulePlan) FieldError {
	value := interface{}(nil)
	current := validate.DerefReflect(field)
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
		tag:         rule.Name,
		actualTag:   rule.Name,
		ns:          ns,
		structNs:    structNs,
		field:       fieldName,
		structField: structFieldName,
		value:       value,
		param:       rule.Param,
		kind:        kind,
		typ:         typ,
	}
}

// shouldDiveIntoStruct 判断字段是否需要递归进入结构体校验
func shouldDiveIntoStruct(field reflect.Value, rules []rule.RulePlan) bool {
	for _, rule := range rules {
		switch rule.Name {
		case "dive", "nostructlevel", "structonly":
			return false
		}
	}
	field = validate.DerefReflect(field)
	return field.IsValid() && field.Kind() == reflect.Struct && !validate.IsTimeType(field.Type())
}

func (v *Validate) compileStruct(t reflect.Type) *rule.StructPlan {
	if cached, ok := v.structCache.Load(t); ok {
		return cached.(*rule.StructPlan)
	}

	typeName := t.Name()
	plan := &rule.StructPlan{Name: typeName, Fields: make([]rule.FieldPlan, 0, t.NumField())}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" && !v.privateFieldValidation {
			continue
		}

		tag := sf.Tag.Get(v.tagName)
		if tag == "-" {
			continue
		}

		altName := v.resolveFieldName(sf)
		fp := rule.FieldPlan{
			Index:          sf.Index,
			Name:           sf.Name,
			AltName:        altName,
			Typ:            sf.Type,
			Rules:          rule.ParseRules(tag),
			HasValidate:    tag != "",
			NsPrefix:       utils.JoinNS(typeName, altName),
			StructNsPrefix: utils.JoinNS(typeName, sf.Name),
		}
		plan.Fields = append(plan.Fields, fp)
	}

	actual, _ := v.structCache.LoadOrStore(t, plan)
	return actual.(*rule.StructPlan)
}

func (v *Validate) resolveFieldName(sf reflect.StructField) string {
	if v.tagNameFunc != nil {
		if name := v.tagNameFunc(sf); name != "" {
			return name
		}
	}
	if jsonTag := sf.Tag.Get("json"); jsonTag != "" {
		name := strings.Split(jsonTag, ",")[0]
		if name != "" && name != "-" {
			return name
		}
	}
	return sf.Name
}
