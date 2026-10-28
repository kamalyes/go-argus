/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-28 15:58:31
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
	"sync/atomic"
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
	New: func() any {
		s := make(ValidationErrors, 0, 8)
		return &s
	},
}

// fieldLevelPool 字段级别校验上下文对象池
var fieldLevelPool = sync.Pool{
	New: func() any {
		return &fieldLevel{}
	},
}

// stringStringMapType map[string]string 的反射类型，map dive 零反射快路径的形状判定
var stringStringMapType = reflect.TypeOf(map[string]string{})

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
	failFast               bool
	tagNameFunc            TagNameFunc

	// regMu 仅序列化注册（写）路径，校验（读）路径不参与
	regMu sync.Mutex
	// validations 自定义规则表，写时复制
	//
	// 原实现是 mu.RWMutex + 裸 map，校验路径每次自定义规则求值都要 RLock——
	// RLock 对共享 readerCount 做原子 RMW，多核在同一缓存行上乒乓，
	// 实测 10 核比单核慢 1.9×（62.20→119.3 ns），而无锁的内建规则路径是 5.08× 正向扩展。
	// 改为写时复制后读路径退化为一次无竞争原子加载；注册是启动期行为，克隆成本无关紧要。
	// 与同结构体内 varSlots 的 atomic.Pointer 用法保持一致。
	validations atomic.Pointer[map[string]FuncCtx]
	structCache sync.Map
	varCache    sync.Map
	// varSlots Var/VarString 规则的分槽旁路缓存（见 varSlotCount）
	varSlots [varSlotCount]atomic.Pointer[varRuleSlot]
}

// New 创建校验器实例
func New(options ...Option) *Validate {
	v := &Validate{
		tagName: defaultTagName,
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
	// 写时复制：克隆当前表、写入新规则、整体发布。
	// 注册期间并发校验读到的仍是旧快照（合法：该规则此刻尚未注册完成），
	// 发布后续读立即可见——与原 RWMutex 的可见性语义等价。
	v.regMu.Lock()
	defer v.regMu.Unlock()

	var next map[string]FuncCtx
	if cur := v.validations.Load(); cur != nil {
		next = make(map[string]FuncCtx, len(*cur)+1)
		for k, f := range *cur {
			next[k] = f
		}
	} else {
		next = make(map[string]FuncCtx, 1)
	}
	next[tag] = fn
	v.validations.Store(&next)
	return nil
}

// Struct 根据结构体字段上的 validate 标签执行校验
func (v *Validate) Struct(s any) error {
	return v.StructCtx(bgCtx, s)
}

// StructCtx 根据结构体字段上的 validate 标签执行校验，并传递 context
func (v *Validate) StructCtx(ctx context.Context, s any) error {
	current := reflect.ValueOf(s)
	if !current.IsValid() {
		return &InvalidValidationError{}
	}
	current = validate.DerefReflect(current)
	if !current.IsValid() || current.Kind() != reflect.Struct {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	errs := acquireErrors()
	// 类型名直接取编译计划里已存的 Name，不再每次校验调一次 reflect.Type.Name()
	// ——后者内部要走 rtype.String() 再切掉包名前缀，profile 里占 260ms/9.07s（2.9%）
	// 顺带把计划下传，避免 walkStruct 再查一次 structCache
	plan := v.compileStruct(current.Type())
	v.walkStruct(ctx, current, current, plan, plan.Name, plan.Name, errs)
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
func (v *Validate) Var(field any, tag string) error {
	return v.VarCtx(bgCtx, field, tag)
}

// VarString 按标签表达式校验字符串变量，零分配快速路径
func (v *Validate) VarString(field, tag string) error {
	return v.VarStringCtx(bgCtx, field, tag)
}

// VarCtx 按标签表达式校验单个变量，并传递 context
func (v *Validate) VarCtx(ctx context.Context, field any, tag string) error {
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
			sc := evalScope{ctx: ctx, parent: rv, errs: errs}
			v.applyRules(&sc, &emptyFieldNames, "", "", rv, rules)
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
	sc := evalScope{ctx: ctx, parent: rv, errs: errs}
	v.applyRules(&sc, &emptyFieldNames, "", "", rv, rules)
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
func (v *Validate) VarStringCtx(ctx context.Context, field, tag string) error {
	if tag == constants.RuleRequired {
		if validate.StringRequired(field) {
			return nil
		}
		return &stringFieldError{tag: constants.RuleRequired, value: field}
	}
	return v.varStringRules(ctx, field, v.cachedVarRules(tag), false)
}

// varStringRules 零反射快速路径，按规则列表校验字符串
func (v *Validate) varStringRules(ctx context.Context, field string, rules []rule.RulePlan, wrapError bool) error {
	for i := 0; i < len(rules); i++ {
		r := &rules[i]
		if r.Kind&rule.KindControl != 0 {
			switch r.Name {
			case constants.RuleOmitEmpty, constants.RuleOmitZero:
				if validate.IsBlankString(field) {
					return nil
				}
			case constants.RuleOmitNil, constants.RuleStructOnly,
				constants.RuleNoStructLevel, constants.RuleEmpty:
			default:
				// dive/keys 等控制规则在字符串快速路径无意义，回退反射路径
				return v.varStringReflectPath(ctx, field, rules)
			}
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
		ok, handled := evalStringRule(field, &rules[i])
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
// 类别与函数已在解析期预解析，运行期零字符串比较、零查表
func evalStringRule(field string, r *rule.RulePlan) (bool, bool) {
	if r.Kind&rule.KindScalarCompare != 0 {
		return r.HasNumber && validate.CompareOp(float64(utf8.RuneCountInString(field)), r.Number, r.CmpOp), true
	}
	if r.Kind&rule.KindHasStr != 0 {
		return r.StrFn == nil || r.StrFn(field, r.Param), true
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
func (v *Validate) stringRuleError(field string, r *rule.RulePlan, wrap bool) error {
	fe := &stringFieldError{tag: r.Name, param: r.Param, value: field}
	if wrap {
		return ValidationErrors{fe}
	}
	return fe
}

// varStringReflectPath 字符串回退到反射路径
func (v *Validate) varStringReflectPath(ctx context.Context, field string, rules []rule.RulePlan) error {
	rv := reflect.ValueOf(field)
	errs := acquireErrors()
	sc := evalScope{ctx: ctx, parent: rv, errs: errs}
	v.applyRules(&sc, &emptyFieldNames, "", "", rv, rules)
	if len(*errs) > 0 {
		result := make(ValidationErrors, len(*errs))
		copy(result, *errs)
		releaseErrors(errs)
		return result
	}
	releaseErrors(errs)
	return nil
}

// varRuleSlot 旁路缓存条目，热点 tag 免走 sync.Map 遍历
type varRuleSlot struct {
	tag   string
	rules []rule.RulePlan
}

// varSlotCount 旁路缓存槽位数，必须为 2 的幂（下标用掩码取）
//
// 原实现是单槽：两个以上热点 tag 轮转时互相驱逐，每次都未命中，
// 既要新分配一个 varRuleSlot（48B），又要向同一个共享 atomic.Pointer 写入——
// 实测单 tag 11.79ns/0 分配，二 tag 退化到 33.18ns/48B，
// 四 tag 并行更是 49.61ns，比单 tag 并行（1.82ns）慢 27×，且慢于自身串行（34.62ns）：
// 典型的缓存行争用负扩展。
// 分槽后不同 tag 落入不同槽位，各自稳定为零分配，写入也不再挤同一缓存行。
const varSlotCount = 16

// varSlotIndex 按 tag 计算槽位下标
//
// 采样长度 + 首字节 + 两个内部字节（1/3、2/3 处）+ 末字节，常数时间，不遍历整串。
//
// 必须采样内部字节：规则 tag 高度同构，普遍共享前后缀——
// "min=2,max=50" 与 "min=3,max=40" 长度、首字节、末字节全同，差异只在中间。
// 仅用首尾会让这类 tag 全部撞进同一槽位，分槽退化为单槽。
func varSlotIndex(tag string) uint {
	n := len(tag)
	if n == 0 {
		return 0
	}
	h := uint(n)*31 + uint(tag[0])*7 + uint(tag[n-1])*13
	h += uint(tag[n/3]) * 131
	h += uint(tag[2*n/3]) * 17
	return h & (varSlotCount - 1)
}

// cachedVarRules 缓存变量规则解析结果，分槽旁路优先，未命中再走 sync.Map
func (v *Validate) cachedVarRules(tag string) []rule.RulePlan {
	idx := varSlotIndex(tag)
	if slot := v.varSlots[idx].Load(); slot != nil && slot.tag == tag {
		return slot.rules
	}
	var rules []rule.RulePlan
	if cached, ok := v.varCache.Load(tag); ok {
		rules = cached.([]rule.RulePlan)
	} else {
		parsed := rule.ParseRules(tag)
		actual, _ := v.varCache.LoadOrStore(tag, parsed)
		rules = actual.([]rule.RulePlan)
	}
	v.varSlots[idx].Store(&varRuleSlot{tag: tag, rules: rules})
	return rules
}

// evalScope 一次校验调用内跨字段、跨规则都不变的求值上下文
//
// 原先 applyRules 按值接收 ctx+top+parent+field+4 个命名空间字符串+rules+errs 共 24 字，
// evalRule 接收 17 字；而 arm64 ABIInternal 只有 16 个整数入参寄存器，
// 超出的部分每次调用都要由调用方落栈、被调方再取回。CPU profile 里
// evalRule 的函数序言独占 740ms/9.07s（8.2%），applyRules 的序言另占 270ms。
// 把 ctx/top/parent/errs 收进一个栈上结构体按指针传递后（再叠加 FieldNames 收口），
// applyRules 降到 13 字、evalRule 降到 7 字，帧从 $352/$208 缩到 $224/$160
//
// 本结构体只读不写：求值路径只从中取值，不回写任何字段，
// 故逃逸分析判定 leaking param content 而非 leaking param，
// 指针本身留在栈上，热路径零分配不变
type evalScope struct {
	ctx    context.Context
	top    reflect.Value
	parent reflect.Value
	errs   *ValidationErrors
}

// emptyFieldNames Var 系列入口的空命名空间四元组
//
// Var 校验没有结构体上下文，四个名字恒为空串；求值路径只读不写它，
// 故共享一份包级零值即可，不必每次调用在栈上装配
var emptyFieldNames rule.FieldNames

// validateStruct 递归校验结构体字段
func (v *Validate) validateStruct(
	ctx context.Context,
	top, current reflect.Value,
	ns, structNs string,
	errs *ValidationErrors,
) {
	current = validate.DerefReflect(current)
	if !current.IsValid() || current.Kind() != reflect.Struct {
		return
	}
	v.walkStruct(ctx, top, current, v.compileStruct(current.Type()), ns, structNs, errs)
}

// walkStruct 按已编译计划遍历结构体字段，要求 current 已是有效的 struct 值
func (v *Validate) walkStruct(
	ctx context.Context,
	top, current reflect.Value,
	plan *rule.StructPlan,
	ns, structNs string,
	errs *ValidationErrors,
) {
	// parent 对本结构体的全部字段恒为 current，故 scope 每层结构体只构造一次
	sc := evalScope{ctx: ctx, top: top, parent: current, errs: errs}
	// 按下标取指针而非 range 取值：FieldPlan 有 15 字，range 的隐式拷贝
	// 在 profile 里占 270ms/9.07s（3%），而循环体只读该计划
	for i := range plan.Fields {
		fp := &plan.Fields[i]
		// compileStruct 只枚举 t.NumField() 的直接字段，Index 恒为单元素，
		// 直通 Field 省去 FieldByIndex 这一层非内联调用（profile 自身占 270ms）；
		// 多级路径保留原调用，不假设计划一定来自 compileStruct
		var field reflect.Value
		if len(fp.Index) == 1 {
			field = current.Field(fp.Index[0])
		} else {
			field = current.FieldByIndex(fp.Index)
		}
		// 命名空间四元组通常在编译期已填好，直接取计划里的指针，每字段零装配；
		// 仅匿名结构体等前缀为空的情形才就地按父命名空间现算一份局部副本
		names := &fp.Names
		if fp.NsPrefix == "" || fp.StructNsPrefix == "" {
			local := fp.Names
			if local.NS == "" {
				local.NS = utils.JoinNS(ns, fp.AltName)
			}
			if local.StructNS == "" {
				local.StructNS = utils.JoinNS(structNs, fp.Name)
			}
			names = &local
		}

		before := len(*errs)
		v.applyRules(&sc, names, names.NS, names.StructNS, field, fp.Rules)
		if len(*errs) != before {
			// failFast：首个错误即收口，不再遍历剩余字段
			if v.failFast {
				return
			}
			continue
		}

		if fp.MayDiveStruct {
			nested := validate.DerefReflect(field)
			if nested.IsValid() && nested.Kind() == reflect.Struct {
				v.validateStruct(ctx, top, nested, names.NS, names.StructNS, errs)
				// failFast：嵌套结构体内部出错同样立即收口
				if v.failFast && len(*errs) > 0 {
					return
				}
			}
		}
	}
}

// applyRules 按规则列表校验单个字段
//
// ns/structNs 覆盖 names 里的同名字段：dive 子元素要带 [i]/[key] 后缀，
// 而 names 指向编译期计划（跨 goroutine 共享的只读数据）不可改写，
// 单独传两个串也比另建一个 FieldNames 局部变量便宜（见 rule.FieldNames 注释）
func (v *Validate) applyRules(
	sc *evalScope,
	names *rule.FieldNames,
	ns, structNs string,
	field reflect.Value,
	rules []rule.RulePlan,
) {
	if len(rules) == 0 {
		return
	}
	derefed := validate.DerefReflect(field)
	for i := 0; i < len(rules); i++ {
		r := &rules[i]
		if r.Kind&rule.KindControl != 0 {
			switch r.Name {
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
				v.applyDive(sc, names, ns, structNs, derefed, rules[i+1:])
				return
			case constants.RuleStructOnly, constants.RuleNoStructLevel, constants.RuleEmpty:
				continue
			}
			// 未消费的 keys/endkeys 等控制规则保持原有求值兜底路径
		}

		ok := false
		if len(r.OrRules) > 0 {
			for j := 0; j < len(r.OrRules); j++ {
				orRule := &r.OrRules[j]
				// 快速路径：or 规则已预解析函数直接调用，其余查表或走完整 evalRule
				if orRule.Kind&rule.KindHasBuiltin != 0 {
					if orRule.Builtin(derefed, orRule.Param, v.requiredStructEnabled) {
						ok = true
						break
					}
				} else if action, found := evalTable[orRule.Name]; found {
					if action.dispatch != nil {
						if action.dispatch(v, sc.top, sc.parent, derefed, orRule) {
							ok = true
							break
						}
					} else if action.builtin(derefed, orRule.Param, v.requiredStructEnabled) {
						ok = true
						break
					}
				} else if v.evalRule(sc, names, derefed, orRule) {
					ok = true
					break
				}
			}
		} else {
			ok = v.evalRule(sc, names, derefed, r)
		}
		if !ok {
			*sc.errs = append(*sc.errs, newFieldError(derefed, names, ns, structNs, r))
			return
		}
	}
}

// applyDive 递归校验切片、数组和 map 元素
func (v *Validate) applyDive(
	sc *evalScope,
	names *rule.FieldNames,
	ns, structNs string,
	field reflect.Value,
	rules []rule.RulePlan,
) {
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
			if v.applyRulesFast(sc, names, item, valueRules) {
				continue
			}
			idxBuf = strconv.AppendInt(idxBuf[:0], int64(i), 10)
			childNS := ns + "[" + string(idxBuf) + "]"
			childStructNS := structNs + "[" + string(idxBuf) + "]"
			if len(valueRules) == 1 && valueRules[0].Name == constants.RuleRequired {
				*sc.errs = append(
					*sc.errs,
					newFieldError(item, names, childNS, childStructNS, &valueRules[0]),
				)
				if v.failFast {
					return
				}
				continue
			}
			v.applyRules(sc, names, childNS, childStructNS, item, valueRules)
			// failFast：元素出错即停，不再遍历剩余元素
			if v.failFast && len(*sc.errs) > 0 {
				return
			}
		}
	case reflect.Map:
		// map[string]string 零反射快路径：值规则恰为单条 required、键规则为空或
		// 单条 required 时直接 range 具体类型，跳过 MapIter 的 SetIterValue 逐条目
		// 反射搬运（含 assignTo 类型核对，entries=100 profile 占 24%）与
		// applyRulesFast→evalRule→RuleRequired 的三层分派（占 27%）。
		// map 值为指针整形，Interface() 装箱零分配；承载体与 MapIter 的每次 dive
		// 常数分配（1 alloc/16B）也一并省去。Name 判等与切片分支的 required 特化
		// 一致："required" 在 BuiltinRules 唯一映射 RuleRequired，or 表达式与别名
		// 的 Name 均不等于 "required"，自然落回通用路径。
		if field.Type() == stringStringMapType && field.CanInterface() &&
			len(valueRules) == 1 && valueRules[0].Name == constants.RuleRequired &&
			(len(keyRules) == 0 || (len(keyRules) == 1 && keyRules[0].Name == constants.RuleRequired)) {
			m := field.Interface().(map[string]string)
			for k, s := range m {
				valBad := validate.IsBlankString(s)
				keyBad := len(keyRules) > 0 && validate.IsBlankString(k)
				if !keyBad && !valBad {
					continue
				}
				// 错误构造与通用路径逐字段对齐：childNS 用键原文；键失败后仍
				// 评估值规则，同一条目可报两条，顺序先键后值
				childNS := ns + "[" + k + "]"
				childStructNS := structNs + "[" + k + "]"
				if keyBad {
					*sc.errs = append(*sc.errs,
						newFieldError(reflect.ValueOf(k), names, childNS, childStructNS, &keyRules[0]))
				}
				if valBad {
					*sc.errs = append(*sc.errs,
						newFieldError(reflect.ValueOf(s), names, childNS, childStructNS, &valueRules[0]))
				}
				if v.failFast {
					return
				}
			}
			return
		}
		// MapRange 迭代器就地推进，不像 MapKeys() 那样先物化整个 []reflect.Value；
		// 校验通过的 map dive 原本要为此付一次按条目数增长的分配
		iter := field.MapRange()

		// 复用键/值承载体：iter.Key() / iter.Value() 每次调用都经 reflect.copyVal
		// → reflect.unsafe_New 分配一个新的值存储（map[string]string 即每条目两次
		// 16B 字符串头分配，实测占本路径全部分配的 100%）。SetIterKey/SetIterValue
		// 就地写入已有的可寻址值，故分配从「2×条目数」降为「每次 map dive 常数次」。
		//
		// 键承载体仅在存在 keys 规则时才需要：applyRulesFast 对空规则立即返回 true，
		// 错误路径另有 iter.Key() 兜底，故无 keys 规则时连承载体都不必分配。
		needKey := len(keyRules) > 0
		var keyHolder reflect.Value
		if needKey {
			keyHolder = reflect.New(field.Type().Key()).Elem()
		}
		valHolder := reflect.New(field.Type().Elem()).Elem()

		for iter.Next() {
			if needKey {
				keyHolder.SetIterKey(iter)
			}
			valHolder.SetIterValue(iter)
			if (!needKey || v.applyRulesFast(sc, names, keyHolder, keyRules)) &&
				v.applyRulesFast(sc, names, valHolder, valueRules) {
				continue
			}

			// 错误路径取 iter 的独立副本，不可复用承载体：newFieldError 会持有该
			// reflect.Value（Kind/Type/装箱均延迟到访问器求值），若传入承载体，
			// 后续条目的 SetIter* 会就地覆写它，已生成的错误将报告错误的值。
			key, value := iter.Key(), iter.Value()
			keyText := validate.StringValue(key)
			childNS := ns + "[" + keyText + "]"
			childStructNS := structNs + "[" + keyText + "]"
			if len(keyRules) > 0 {
				v.applyRules(sc, names, childNS, childStructNS, key, keyRules)
			}
			v.applyRules(sc, names, childNS, childStructNS, value, valueRules)
			// failFast：键或值出错即停，不再遍历剩余条目
			if v.failFast && len(*sc.errs) > 0 {
				return
			}
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
func (v *Validate) applyRulesFast(
	sc *evalScope,
	names *rule.FieldNames,
	field reflect.Value,
	rules []rule.RulePlan,
) bool {
	if len(rules) == 0 {
		return true
	}
	derefed := validate.DerefReflect(field)
	for i := 0; i < len(rules); i++ {
		r := &rules[i]
		if r.Kind&rule.KindControl != 0 {
			switch r.Name {
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
						if !v.applyRulesFast(sc, names, derefed.Index(j), valueRules) {
							return false
						}
					}
				case reflect.Map:
					// 同 applyDive：MapRange 迭代器就地推进，避免 MapKeys() 物化键切片；
					// 键/值承载体复用——iter.Key()/iter.Value() 每次调用都经
					// reflect.copyVal → unsafe_New 分配新值存储（每条目 2 次），
					// SetIterKey/SetIterValue 就地写入已有可寻址值，分配降为每次
					// dive 常数次。本函数是纯成功路径：失败即返回 false，错误由
					// 调用方经慢路径以 iter 的独立副本重建，承载体不会被错误对象
					// 持有，可安全复用。承载体按需建：对空规则列表 applyRulesFast
					// 不读取字段即返回 true，keys-only 或 values-only 的嵌套 dive
					// 连对应的承载体都不必分配。
					iter := derefed.MapRange()
					needKey := len(keyRules) > 0
					needVal := len(valueRules) > 0
					var keyHolder, valHolder reflect.Value
					if needKey {
						keyHolder = reflect.New(derefed.Type().Key()).Elem()
					}
					if needVal {
						valHolder = reflect.New(derefed.Type().Elem()).Elem()
					}
					for iter.Next() {
						if needKey {
							keyHolder.SetIterKey(iter)
							if !v.applyRulesFast(sc, names, keyHolder, keyRules) {
								return false
							}
						}
						if needVal {
							valHolder.SetIterValue(iter)
							if !v.applyRulesFast(sc, names, valHolder, valueRules) {
								return false
							}
						}
					}
				}
				return true
			case constants.RuleKeys,
				constants.RuleEndKeys,
				constants.RuleStructOnly,
				constants.RuleNoStructLevel,
				constants.RuleEmpty:
				continue
			}
		}

		ok := false
		if len(r.OrRules) > 0 {
			for j := 0; j < len(r.OrRules); j++ {
				orRule := &r.OrRules[j]
				if orRule.Kind&rule.KindHasBuiltin != 0 {
					if orRule.Builtin(derefed, orRule.Param, v.requiredStructEnabled) {
						ok = true
						break
					}
				} else if action, found := evalTable[orRule.Name]; found {
					if action.dispatch != nil {
						if action.dispatch(v, sc.top, sc.parent, derefed, orRule) {
							ok = true
							break
						}
					} else if action.builtin(derefed, orRule.Param, v.requiredStructEnabled) {
						ok = true
						break
					}
				} else if v.evalRule(sc, names, derefed, orRule) {
					ok = true
					break
				}
			}
		} else {
			ok = v.evalRule(sc, names, derefed, r)
		}
		if !ok {
			return false
		}
	}
	return true
}

// evalRule 评估单条规则，类别与函数已在解析期预解析，仅 dispatch 与自定义规则需查表
func (v *Validate) evalRule(
	sc *evalScope,
	names *rule.FieldNames,
	field reflect.Value,
	plan *rule.RulePlan,
) bool {
	if plan.Kind&rule.KindScalarCompare != 0 {
		return plan.HasNumber && rule.CompareLengthOrNumber(field, plan.Number, plan.CmpOp)
	}
	if plan.Kind&rule.KindHasBuiltin != 0 {
		return plan.Builtin(field, plan.Param, v.requiredStructEnabled)
	}
	if action, ok := evalTable[plan.Name]; ok {
		if action.dispatch != nil {
			return action.dispatch(v, sc.top, sc.parent, field, plan)
		}
		return action.builtin(field, plan.Param, v.requiredStructEnabled)
	}

	// 仅在 builtin/dispatch 未命中时才加读锁查自定义规则
	fn := v.getCustomValidation(plan.Name)
	if fn == nil {
		return false
	}

	fl := acquireFieldLevel()
	fl.top = sc.top
	fl.parent = sc.parent
	fl.field = validate.DerefReflect(field)
	fl.fieldName = names.Field
	fl.structFieldName = names.StructField
	fl.tag = plan.Name
	fl.param = plan.Param
	result := fn(sc.ctx, fl)
	releaseFieldLevel(fl)
	return result
}

// getCustomValidation 取自定义规则函数（校验热路径，一次无竞争原子加载）
//
// 表未初始化（从未注册过自定义规则）时返回 nil，调用方按"无此规则"处理——
// 与原裸 map 零值查找语义一致。
func (v *Validate) getCustomValidation(name string) FuncCtx {
	m := v.validations.Load()
	if m == nil {
		return nil
	}
	return (*m)[name]
}

// evalDispatchFn 规则分派函数签名
type evalDispatchFn func(v *Validate, top, parent, field reflect.Value, plan *rule.RulePlan) bool

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

func (v *Validate) evalRequiredIf(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredIf(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredUnless(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if rule.IsRequiredIf(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWith(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredWith(parent, plan.Param) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWithAll(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredWithAll(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWithout(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredWithout(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalRequiredWithoutAll(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredWithoutAll(parent, plan.ParamParts) {
		return true
	}
	return !validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedIf(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredIf(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedUnless(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if rule.IsRequiredIf(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWith(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredWith(parent, plan.Param) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWithAll(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredWithAll(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWithout(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredWithout(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalExcludedWithoutAll(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if !rule.IsRequiredWithoutAll(parent, plan.ParamParts) {
		return true
	}
	return validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled)
}

func (v *Validate) evalCmpField(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
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
	// 兜底推导只落局部变量，绝不回写 plan：plan 指向 structCache 共享的 StructPlan，
	// 校验路径是并发只读路径，回写即对共享结构体的无同步写（竞争）。
	// 正常情况下 PrepareRulePlan（rule/plan.go）已在编译期置好 CmpOp/HasCmpOp，
	// 此分支仅当某规则进了 evalDispatchTable 却漏登记 constants.IsFieldCompareRule 时触发——
	// 届时若回写便是真实 data race，故此处只读不写，编译期才是唯一正确的赋值位置。
	cmpOp := plan.CmpOp
	if !plan.HasCmpOp {
		cmpOp = rule.CmpOpForRule(plan.Name)
	}
	return rule.CompareValueOp(derefed, other, cmpOp)
}

// evalCmpFieldLen 跨字段长度比较，参数格式: FieldName 或 op:FieldName（op 可选 eq/ne/gt/gte/lt/lte，默认 eq）
// 支持嵌套字段路径如 Inner.Field，从顶层结构体开始查找
func (v *Validate) evalCmpFieldLen(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	opStr := constants.RuleEq
	targetPath := plan.Param
	if idx := strings.Index(plan.Param, ":"); idx >= 0 {
		candidate := plan.Param[:idx]
		switch candidate {
		case constants.RuleEq, constants.RuleNe, constants.RuleGT, constants.RuleGTE, constants.RuleLT, constants.RuleLTE:
			opStr = candidate
			targetPath = plan.Param[idx+1:]
		}
	}
	cmpOp := constants.CmpOpFromStr(opStr)

	currentLen := fieldLen(field)
	// 优先从 parent 查找，找不到则从 top（根结构体）查找，支持跨层级访问
	other, ok := rule.FieldByPath(parent, targetPath)
	if !ok {
		other, ok = rule.FieldByPath(top, targetPath)
	}
	if !ok {
		return false
	}
	otherLen := fieldLen(other)

	return validate.CompareOp(float64(currentLen), float64(otherLen), cmpOp)
}

// fieldLen 获取字段的长度（支持 slice/array/map/string/chan，其他类型返回 0）
func fieldLen(v reflect.Value) int {
	v = validate.DerefReflect(v)
	if !v.IsValid() {
		return 0
	}
	switch v.Kind() {
	case reflect.String:
		return utf8.RuneCountInString(v.String())
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return v.Len()
	default:
		return 0
	}
}

func (v *Validate) evalFieldContains(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	return rule.FieldContains(field, parent, plan.Param)
}

func (v *Validate) evalFieldExcludes(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	return !rule.FieldContains(field, parent, plan.Param)
}

func (v *Validate) evalAfter(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled) {
		return true
	}
	return rule.CompareTimeExpr(field, plan.Param, constants.RuleGT, time.Now())
}

func (v *Validate) evalBefore(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	if validate.IsEmptyValueWithStruct(field, v.requiredStructEnabled) {
		return true
	}
	return rule.CompareTimeExpr(field, plan.Param, constants.RuleLT, time.Now())
}

func (v *Validate) evalRange(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	return rule.Range(parent, plan.Param)
}

func (v *Validate) evalOneOf(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	return rule.OneOfFast(field, plan.ParamParts)
}

func (v *Validate) evalOneOfCI(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	return rule.OneOfCIFast(field, plan.ParamParts)
}

func (v *Validate) evalNoneOf(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
	return !rule.OneOfFast(field, plan.ParamParts)
}

func (v *Validate) evalNoneOfCI(top, parent, field reflect.Value, plan *rule.RulePlan) bool {
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
	constants.RuleFieldLen:           (*Validate).evalCmpFieldLen,
	constants.RuleAfter:              (*Validate).evalAfter,
	constants.RuleBefore:             (*Validate).evalBefore,
	constants.RuleRange:              (*Validate).evalRange,
	constants.RuleOneOf:              (*Validate).evalOneOf,
	constants.RuleOneOfCI:            (*Validate).evalOneOfCI,
	constants.RuleNoneOf:             (*Validate).evalNoneOf,
	constants.RuleNoneOfCI:           (*Validate).evalNoneOfCI,
}

func newFieldError(
	field reflect.Value,
	names *rule.FieldNames,
	ns, structNs string,
	r *rule.RulePlan,
) FieldError {
	// 惰性求值：构造仅存 deref 后的 reflect.Value，装箱与 Kind/Type 推迟到访问器
	// 失败路径每错误省去一次 Interface 装箱分配与两次反射调用
	current := validate.DerefReflect(field)
	return &fieldError{
		tag:         r.Name,
		ns:          ns,
		structNs:    structNs,
		field:       names.Field,
		structField: names.StructField,
		fieldVal:    current,
		param:       r.Param,
	}
}

func mayDiveStructType(typ reflect.Type, rules []rule.RulePlan) bool {
	for _, rule := range rules {
		if constants.StopsStructDive(rule.Name) {
			return false
		}
	}
	for typ.Kind() == reflect.Pointer {
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
		mayDive := mayDiveStructType(sf.Type, rules)

		// 既无规则又不可下潜的字段不入计划：它在 validateStruct 单轮迭代里可证明为 no-op
		// （applyRules 对空规则立即返回、MayDiveStruct 为 false 不下潜、FieldByIndex 结果无消费者），
		// 而 proto/@inject_tag 生成的宽结构体中这类字段占绝大多数，逐字段 FieldByIndex
		// 与命名空间计算使耗时正比于「总字段数」而非「带规则字段数」
		//
		// ⚠️ 判定必须是 len(rules)==0 && !mayDive，不可用 HasValidate：
		// 无 tag 的嵌套结构体字段 HasValidate 为 false，但必须保留以供下潜校验
		if len(rules) == 0 && !mayDive {
			continue
		}

		fp := rule.FieldPlan{
			Index:          sf.Index,
			Name:           sf.Name,
			AltName:        altName,
			Typ:            sf.Type,
			Rules:          rules,
			HasValidate:    tag != "",
			MayDiveStruct:  mayDive,
			NsPrefix:       utils.JoinNS(typeName, altName),
			StructNsPrefix: utils.JoinNS(typeName, sf.Name),
		}
		fp.Names = rule.FieldNames{
			NS:          fp.NsPrefix,
			StructNS:    fp.StructNsPrefix,
			Field:       altName,
			StructField: sf.Name,
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
