/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 16:54:16
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
	"unicode/utf8"

	"github.com/kamalyes/go-argus/constants"
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
	switch val := field.(type) {
	case string:
		if tag == constants.RuleRequired {
			if validate.StringRequired(val) {
				return nil
			}
			return ValidationErrors{&stringFieldError{tag: constants.RuleRequired, value: val}}
		}
		if name, param, ok := strings.Cut(tag, "="); ok && name == "endswith" {
			if strings.HasSuffix(val, param) {
				return nil
			}
			return ValidationErrors{&stringFieldError{tag: name, param: param, value: val}}
		}
		rules := v.cachedVarRules(tag)
		return v.varStringRules(ctx, val, rules, true)
	case *string:
		if val == nil {
			rules := v.cachedVarRules(tag)
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
		if tag == constants.RuleRequired {
			if validate.StringRequired(*val) {
				return nil
			}
			return ValidationErrors{&stringFieldError{tag: constants.RuleRequired, value: *val}}
		}
		if name, param, ok := strings.Cut(tag, "="); ok && name == "endswith" {
			if strings.HasSuffix(*val, param) {
				return nil
			}
			return ValidationErrors{&stringFieldError{tag: name, param: param, value: *val}}
		}
		rules := v.cachedVarRules(tag)
		return v.varStringRules(ctx, *val, rules, true)
	}
	rules := v.cachedVarRules(tag)
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
		case constants.RuleOmitEmpty, constants.RuleOmitZero:
			if validate.IsBlankString(field) {
				return nil
			}
			continue
		case constants.RuleOmitNil:
			continue
		case constants.RuleStructOnly, constants.RuleNoStructLevel, constants.RuleEmpty:
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
	if constants.IsScalarCompareRule(r.Name) {
		return r.HasNumber && validate.CompareOp(float64(utf8.RuneCountInString(field)), r.Number, r.CmpOp), true
	}
	if fn, ok := rule.StringRuleMap[r.Name]; ok {
		return fn == nil || fn(field, r.Param), true
	}
	switch r.Name {
	case constants.RuleOneOf:
		return validate.StringOneOf(field, r.ParamParts), true
	case constants.RuleOneOfCI:
		return validate.StringOneOfCI(field, r.ParamParts), true
	case constants.RuleNoneOf:
		return !validate.StringOneOf(field, r.ParamParts), true
	case constants.RuleNoneOfCI:
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

		if fp.MayDiveStruct {
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
		case constants.RuleOmitEmpty, constants.RuleOmitZero:
			if validate.IsEmptyValueWithStruct(derefed, v.requiredStructEnabled) {
				return
			}
			continue
		case constants.RuleOmitNil:
			if validate.IsNilValue(field) {
				return
			}
			continue
		case constants.RuleDive:
			v.applyDive(ctx, top, parent, derefed, ns, structNs, fieldName, structFieldName, rules[i+1:], errs)
			return
		case constants.RuleStructOnly, constants.RuleNoStructLevel, constants.RuleEmpty:
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
	keyRules, valueRules := splitDiveRules(rules)
	switch field.Kind() {
	case reflect.Slice, reflect.Array:
		// 复用 []byte 缓冲区避免 strconv.Itoa 分配
		var idxBuf []byte
		for i := 0; i < field.Len(); i++ {
			item := field.Index(i)
			if v.applyRulesFast(ctx, top, parent, item, fieldName, structFieldName, valueRules) {
				continue
			}
			idxBuf = strconv.AppendInt(idxBuf[:0], int64(i), 10)
			childNS := ns + "[" + string(idxBuf) + "]"
			childStructNS := structNs + "[" + string(idxBuf) + "]"
			if len(valueRules) == 1 && valueRules[0].Name == constants.RuleRequired {
				*errs = append(*errs, newFieldError(item, childNS, childStructNS, fieldName, structFieldName, valueRules[0]))
				continue
			}
			v.applyRules(ctx, top, parent, item, childNS, childStructNS, fieldName, structFieldName, valueRules, errs)
		}
	case reflect.Map:
		for _, key := range field.MapKeys() {
			value := field.MapIndex(key)
			if v.applyRulesFast(ctx, top, parent, key, fieldName, structFieldName, keyRules) &&
				v.applyRulesFast(ctx, top, parent, value, fieldName, structFieldName, valueRules) {
				continue
			}
			keyText := validate.StringValue(key)
			childNS := ns + "[" + keyText + "]"
			childStructNS := structNs + "[" + keyText + "]"
			if len(keyRules) > 0 {
				v.applyRules(ctx, top, parent, key, childNS, childStructNS, fieldName, structFieldName, keyRules, errs)
			}
			v.applyRules(ctx, top, parent, value, childNS, childStructNS, fieldName, structFieldName, valueRules, errs)
		}
	}
}

func splitDiveRules(rules []rule.RulePlan) ([]rule.RulePlan, []rule.RulePlan) {
	if len(rules) == 0 || rules[0].Name != constants.RuleKeys {
		return nil, rules
	}
	for i := 1; i < len(rules); i++ {
		if rules[i].Name == constants.RuleEndKeys {
			return rules[1:i], rules[i+1:]
		}
	}
	return rules[1:], nil
}

// applyRulesFast 校验成功路径，不构造错误命名空间
func (v *Validate) applyRulesFast(ctx context.Context, top reflect.Value, parent reflect.Value, field reflect.Value, fieldName string, structFieldName string, rules []rule.RulePlan) bool {
	if len(rules) == 0 {
		return true
	}
	derefed := validate.DerefReflect(field)
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		switch rule.Name {
		case constants.RuleOmitEmpty, constants.RuleOmitZero:
			if validate.IsEmptyValueWithStruct(derefed, v.requiredStructEnabled) {
				return true
			}
			continue
		case constants.RuleOmitNil:
			if validate.IsNilValue(field) {
				return true
			}
			continue
		case constants.RuleDive:
			keyRules, valueRules := splitDiveRules(rules[i+1:])
			switch derefed.Kind() {
			case reflect.Slice, reflect.Array:
				for j := 0; j < derefed.Len(); j++ {
					if !v.applyRulesFast(ctx, top, parent, derefed.Index(j), fieldName, structFieldName, valueRules) {
						return false
					}
				}
			case reflect.Map:
				for _, key := range derefed.MapKeys() {
					if !v.applyRulesFast(ctx, top, parent, key, fieldName, structFieldName, keyRules) ||
						!v.applyRulesFast(ctx, top, parent, derefed.MapIndex(key), fieldName, structFieldName, valueRules) {
						return false
					}
				}
			}
			return true
		case constants.RuleKeys, constants.RuleEndKeys, constants.RuleStructOnly, constants.RuleNoStructLevel, constants.RuleEmpty:
			continue
		}

		ok := false
		if len(rule.OrRules) > 0 {
			for j := 0; j < len(rule.OrRules); j++ {
				orRule := rule.OrRules[j]
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
			return false
		}
	}
	return true
}

// evalRule 评估单条规则，优先查 dispatch 表，其次查 builtin 表，最后查自定义注册
func (v *Validate) evalRule(ctx context.Context, top reflect.Value, parent reflect.Value, field reflect.Value, fieldName string, structFieldName string, plan rule.RulePlan) bool {
	if constants.IsScalarCompareRule(plan.Name) {
		return plan.HasNumber && rule.CompareLengthOrNumber(field, plan.Number, plan.CmpOp)
	}
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
	if !rule.IsRequiredIf(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredUnless(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if rule.IsRequiredIf(parent, plan.ParamParts) {
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
	if !rule.IsRequiredIf(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedUnless(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	if rule.IsRequiredIf(parent, plan.ParamParts) {
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
	target := parent
	if constants.IsCrossStructFieldCompareRule(plan.Name) {
		target = top
	}
	derefed := validate.DerefReflect(field)
	var other reflect.Value
	if len(plan.FieldIndex) > 0 {
		target = validate.DerefReflect(target)
		if !target.IsValid() || target.Kind() != reflect.Struct {
			return false
		}
		other = target.FieldByIndex(plan.FieldIndex)
	} else {
		var ok bool
		other, ok = rule.FieldByPath(target, plan.Param)
		if !ok {
			return false
		}
	}
	if !plan.HasCmpOp {
		plan.CmpOp = rule.CmpOpForRule(plan.Name)
	}
	return rule.CompareValueOp(derefed, other, plan.CmpOp)
}

func (v *Validate) evalFieldContains(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return rule.FieldContains(field, parent, plan.Param)
}

func (v *Validate) evalFieldExcludes(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return !rule.FieldContains(field, parent, plan.Param)
}

func (v *Validate) evalAfter(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return rule.CompareTimeExpr(field, plan.Param, constants.RuleGT, time.Now())
}

func (v *Validate) evalBefore(top, parent, field reflect.Value, plan rule.RulePlan) bool {
	return rule.CompareTimeExpr(field, plan.Param, constants.RuleLT, time.Now())
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

// evalDispatchTable 需要跨字段访问的规则分派表
var evalDispatchTable = map[string]evalDispatchFn{
	constants.RuleRequiredIf:         (*Validate).evalRequiredIf,
	constants.RuleRequiredUnless:     (*Validate).evalRequiredUnless,
	constants.RuleRequiredWith:       (*Validate).evalRequiredWith,
	constants.RuleRequiredWithAll:    (*Validate).evalRequiredWithAll,
	constants.RuleRequiredWithout:    (*Validate).evalRequiredWithout,
	constants.RuleRequiredWithoutAll: (*Validate).evalRequiredWithoutAll,
	constants.RuleExcludedIf:         (*Validate).evalExcludedIf,
	constants.RuleExcludedUnless:     (*Validate).evalExcludedUnless,
	constants.RuleExcludedWith:       (*Validate).evalExcludedWith,
	constants.RuleExcludedWithAll:    (*Validate).evalExcludedWithAll,
	constants.RuleExcludedWithout:    (*Validate).evalExcludedWithout,
	constants.RuleExcludedWithoutAll: (*Validate).evalExcludedWithoutAll,
	constants.RuleEqField:            (*Validate).evalCmpField,
	constants.RuleNeField:            (*Validate).evalCmpField,
	constants.RuleGTField:            (*Validate).evalCmpField,
	constants.RuleAfterField:         (*Validate).evalCmpField,
	constants.RuleGTEField:           (*Validate).evalCmpField,
	constants.RuleLTField:            (*Validate).evalCmpField,
	constants.RuleBeforeField:        (*Validate).evalCmpField,
	constants.RuleLTEField:           (*Validate).evalCmpField,
	constants.RuleEqCSField:          (*Validate).evalCmpField,
	constants.RuleNeCSField:          (*Validate).evalCmpField,
	constants.RuleGTCSField:          (*Validate).evalCmpField,
	constants.RuleGTECSField:         (*Validate).evalCmpField,
	constants.RuleLTCSField:          (*Validate).evalCmpField,
	constants.RuleLTECSField:         (*Validate).evalCmpField,
	constants.RuleFieldContains:      (*Validate).evalFieldContains,
	constants.RuleFieldExcludes:      (*Validate).evalFieldExcludes,
	constants.RuleAfter:              (*Validate).evalAfter,
	constants.RuleBefore:             (*Validate).evalBefore,
	constants.RuleRange:              (*Validate).evalRange,
	constants.RuleOneOf:              (*Validate).evalOneOf,
	constants.RuleOneOfCI:            (*Validate).evalOneOfCI,
	constants.RuleNoneOf:             (*Validate).evalNoneOf,
	constants.RuleNoneOfCI:           (*Validate).evalNoneOfCI,
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

func mayDiveStructType(typ reflect.Type, rules []rule.RulePlan) bool {
	for _, rule := range rules {
		if constants.StopsStructDive(rule.Name) {
			return false
		}
	}
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.Kind() == reflect.Struct && !validate.IsTimeType(typ)
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
		rules := rule.ParseRules(tag)
		for i := range rules {
			if constants.IsLocalFieldCompareRule(rules[i].Name) {
				if index, ok := rule.FieldIndexByPath(t, rules[i].Param); ok {
					rules[i].FieldIndex = index
				}
			}
			for j := range rules[i].OrRules {
				if constants.IsLocalFieldCompareRule(rules[i].OrRules[j].Name) {
					if index, ok := rule.FieldIndexByPath(t, rules[i].OrRules[j].Param); ok {
						rules[i].OrRules[j].FieldIndex = index
					}
				}
			}
		}
		fp := rule.FieldPlan{
			Index:          sf.Index,
			Name:           sf.Name,
			AltName:        altName,
			Typ:            sf.Type,
			Rules:          rules,
			HasValidate:    tag != "",
			MayDiveStruct:  mayDiveStructType(sf.Type, rules),
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
