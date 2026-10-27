/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-27 19:58:37
 * @FilePath: \go-argus\rule\plan.go
 * @Description: 规则计划解析，将 validate 标签预编译为运行计划
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"strings"

	"github.com/kamalyes/go-argus/constants"
	"github.com/kamalyes/go-argus/validate"
)

// RuleKind 规则类别位掩码，解析期一次性预分类，运行期以 int 位测试替代字符串比较与查表
type RuleKind uint8

const (
	// KindScalarCompare 标量比较规则（min/max/len/gt/gte/lt/lte），求值路径最优先短路
	KindScalarCompare RuleKind = 1 << iota
	// KindControl 流程控制规则（omitempty/dive/keys 等），不产生求值
	KindControl
	// KindHasStr StringRuleMap 命中，StrFn 为预解析函数
	KindHasStr
	// KindHasBuiltin BuiltinRules 命中，Builtin 为预解析函数
	KindHasBuiltin
)

// RulePlan 表示一条规则的运行计划
type RulePlan struct {
	Name       string     // 规则名称
	Param      string     // 规则参数原始值
	ParamParts []string   // 规则参数按空白拆分
	OrRules    []RulePlan // 或规则列表（对应 | 分隔符）
	CmpOp      constants.CmpOp
	HasCmpOp   bool
	FieldIndex []int
	Number     float64
	HasNumber  bool
	Kind       RuleKind       // 类别位掩码
	StrFn      StringRuleFunc // StringRuleMap 预解析函数
	Builtin    BuiltinRule    // BuiltinRules 预解析函数
}

// FieldNames 字段的命名空间与名称四元组，编译期一次填好，求值路径按指针直取
//
// 原先这四个串是摊开传的，applyRules 的入参恰好 24 字；arm64 ABIInternal 只有
// 16 个整数入参寄存器，超出部分每次调用都要落栈再取回，而「全部通过」的热路径
// 根本不读它们——只有错误构造、dive 子路径拼接和自定义规则上下文才需要。
// 收进本结构体后 applyRules 降到 13 字、evalRule 降到 7 字，帧也相应变小
//
// 注：dive 子元素的命名空间仍由 applyRules/applyDive 以 ns/structNs 形参单独传，
// 不另建一个 FieldNames 局部变量——applyRules 与 applyDive 互递归，
// 局部变量取址后会流进这个强连通分量里的形参，逃逸分析收敛不了，
// 会把它按出错元素逐个搬上堆（实测 dive 失败路径 6→7 次分配）
type FieldNames struct {
	NS          string // 展示命名空间，形如 Type.Field
	StructNS    string // 结构体命名空间，形如 Type.GoField
	Field       string // 展示名（json tag 或自定义）
	StructField string // Go 字段名
}

// FieldPlan 表示结构体字段的校验计划
type FieldPlan struct {
	Index          []int  // 字段索引路径
	Name           string // Go 字段名
	AltName        string // 展示名（json tag 或自定义）
	Typ            reflect.Type
	Rules          []RulePlan // 字段上的规则列表
	HasValidate    bool       // 是否包含校验标签
	MayDiveStruct  bool       // 是否可能递归进入结构体字段
	NsPrefix       string     // 命名空间前缀
	StructNsPrefix string     // 结构体命名空间前缀
	// Names 编译期预填的四元组，供求值路径按指针直取；
	// 与上面四个字段同源，仅为省掉每字段一次装配而冗余存一份
	Names FieldNames
}

// StructPlan 表示结构体的编译计划
type StructPlan struct {
	Name   string      // 结构体名称
	Fields []FieldPlan // 字段计划列表
}

// ParseRules 将 validate 标签解析为规则计划列表
func ParseRules(tag string) []RulePlan {
	parsed := ParseTag(tag)
	rules := make([]RulePlan, 0, len(parsed))
	for _, item := range parsed {
		rules = append(rules, ParseRulePlan(item))
	}
	return rules
}

// ParseRulePlan 将单条规则解析为运行计划，支持或规则展开
func ParseRulePlan(item Rule) RulePlan {
	raw := strings.TrimSpace(item.Raw)
	if raw != "" {
		parts := SplitRuleOr(raw)
		if len(parts) > 1 {
			rp := RulePlan{Name: raw, OrRules: make([]RulePlan, 0, len(parts))}
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				rp.OrRules = append(rp.OrRules, ParseSingleRulePlan(part))
			}
			return rp
		}
	}
	return PrepareRulePlan(RulePlan{Name: item.Name, Param: item.Param})
}

// ParseSingleRulePlan 解析单条规则字符串为运行计划
func ParseSingleRulePlan(raw string) RulePlan {
	name, param, ok := strings.Cut(raw, "=")
	if !ok {
		name = raw
	}
	return PrepareRulePlan(RulePlan{Name: strings.TrimSpace(name), Param: strings.TrimSpace(param)})
}

// PrepareRulePlan 对规则计划做后处理，预拆分需要空白分词的规则参数
// 同时完成类别预分类与函数预解析，类别判定的优先级与求值路径保持一致（标量比较最优先短路）
func PrepareRulePlan(rp RulePlan) RulePlan {
	if constants.NeedsParamParts(rp.Name) {
		if rp.Param != "" {
			rp.ParamParts = strings.Fields(rp.Param)
		}
	}
	if constants.IsScalarCompareRule(rp.Name) {
		rp.Kind |= KindScalarCompare
		rp.CmpOp = CmpOpForRule(rp.Name)
		rp.HasCmpOp = true
		if rp.Param != "" {
			rp.Number, rp.HasNumber = validate.ParseFloat(rp.Param)
		}
		return rp
	}
	if IsControlRule(rp.Name) {
		rp.Kind |= KindControl
		return rp
	}
	if fn, ok := StringRuleMap[rp.Name]; ok {
		rp.Kind |= KindHasStr
		rp.StrFn = fn
	}
	if fn, ok := BuiltinRules[rp.Name]; ok {
		rp.Kind |= KindHasBuiltin
		rp.Builtin = fn
	}
	if constants.IsFieldCompareRule(rp.Name) {
		rp.CmpOp = CmpOpForRule(rp.Name)
		rp.HasCmpOp = true
	}
	return rp
}

// IsControlRule 是否为流程控制规则，不参与规则求值
func IsControlRule(name string) bool {
	switch name {
	case constants.RuleOmitEmpty, constants.RuleOmitZero, constants.RuleOmitNil,
		constants.RuleDive, constants.RuleKeys, constants.RuleEndKeys,
		constants.RuleStructOnly, constants.RuleNoStructLevel, constants.RuleEmpty:
		return true
	}
	return false
}

// CmpOpForRule maps validation rules to comparison operators.
func CmpOpForRule(name string) constants.CmpOp {
	op := constants.CompareOperatorForRule(name)
	if op == constants.RuleEmpty {
		return -1
	}
	return constants.CmpOpFromStr(op)
}

// SplitRuleOr 按 | 分隔或规则，支持转义符
func SplitRuleOr(s string) []string {
	start := 0
	escaped := false
	var parts []string
	for i := 0; i < len(s); i++ {
		if escaped {
			escaped = false
			continue
		}
		switch s[i] {
		case '\\':
			escaped = true
		case '|':
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	if len(parts) == 0 {
		return nil
	}
	parts = append(parts, s[start:])
	if strings.Contains(parts[0], "=") {
		for i := 1; i < len(parts); i++ {
			if !strings.Contains(parts[i], "=") {
				return nil
			}
		}
	}
	return parts
}
