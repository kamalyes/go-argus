/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 00:00:00
 * @FilePath: \go-argus\constants\rules.go
 * @Description: Central validation rule names and rule groups
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

const (
	RuleEmpty = ""

	RuleRequired = "required"

	RuleOmitEmpty = "omitempty"
	RuleOmitZero  = "omitzero"
	RuleOmitNil   = "omitnil"

	RuleDive          = "dive"
	RuleKeys          = "keys"
	RuleEndKeys       = "endkeys"
	RuleStructOnly    = "structonly"
	RuleNoStructLevel = "nostructlevel"

	RuleLen = "len"
	RuleMin = "min"
	RuleMax = "max"
	RuleEq  = "eq"
	RuleNe  = "ne"
	RuleGT  = "gt"
	RuleGTE = "gte"
	RuleLT  = "lt"
	RuleLTE = "lte"

	RuleSymbolEqual              = "="
	RuleSymbolNotEqual           = "!="
	RuleSymbolGreaterThan        = ">"
	RuleSymbolGreaterThanOrEqual = ">="
	RuleSymbolLessThan           = "<"
	RuleSymbolLessThanOrEqual    = "<="

	RuleOneOf    = "oneof"
	RuleOneOfCI  = "oneofci"
	RuleNoneOf   = "noneof"
	RuleNoneOfCI = "noneofci"

	RuleRequiredIf         = "required_if"
	RuleRequiredUnless     = "required_unless"
	RuleRequiredWith       = "required_with"
	RuleRequiredWithAll    = "required_with_all"
	RuleRequiredWithout    = "required_without"
	RuleRequiredWithoutAll = "required_without_all"

	RuleExcludedIf         = "excluded_if"
	RuleExcludedUnless     = "excluded_unless"
	RuleExcludedWith       = "excluded_with"
	RuleExcludedWithAll    = "excluded_with_all"
	RuleExcludedWithout    = "excluded_without"
	RuleExcludedWithoutAll = "excluded_without_all"

	RuleEqField     = "eqfield"
	RuleNeField     = "nefield"
	RuleGTField     = "gtfield"
	RuleGTEField    = "gtefield"
	RuleLTField     = "ltfield"
	RuleLTEField    = "ltefield"
	RuleAfterField  = "afterfield"
	RuleBeforeField = "beforefield"

	RuleEqCSField  = "eqcsfield"
	RuleNeCSField  = "necsfield"
	RuleGTCSField  = "gtcsfield"
	RuleGTECSField = "gtecsfield"
	RuleLTCSField  = "ltcsfield"
	RuleLTECSField = "ltecsfield"

	RuleAfter         = "after"
	RuleBefore        = "before"
	RuleRange         = "range"
	RuleFieldContains = "fieldcontains"
	RuleFieldExcludes = "fieldexcludes"
	RuleRegex         = "regex"
	RuleNotEmpty      = "not_empty"
)

// NeedsParamParts 是否需要参数部分
func NeedsParamParts(name string) bool {
	switch name {
	case RuleOneOf, RuleOneOfCI, RuleNoneOf, RuleNoneOfCI,
		RuleRequiredWith, RuleRequiredWithAll, RuleRequiredWithout, RuleRequiredWithoutAll,
		RuleExcludedWith, RuleExcludedWithAll, RuleExcludedWithout, RuleExcludedWithoutAll,
		RuleRequiredIf, RuleRequiredUnless, RuleExcludedIf, RuleExcludedUnless:
		return true
	default:
		return false
	}
}

// IsScalarCompareRule 是否为标量比较规则
func IsScalarCompareRule(name string) bool {
	switch name {
	case RuleMin, RuleMax, RuleLen, RuleGT, RuleGTE, RuleLT, RuleLTE:
		return true
	default:
		return false
	}
}

// IsFieldCompareRule 是否为跨字段比较规则
func IsFieldCompareRule(name string) bool {
	switch name {
	case RuleEqField, RuleNeField, RuleGTField, RuleAfterField, RuleGTEField, RuleLTField, RuleBeforeField, RuleLTEField,
		RuleEqCSField, RuleNeCSField, RuleGTCSField, RuleGTECSField, RuleLTCSField, RuleLTECSField:
		return true
	default:
		return false
	}
}

// IsLocalFieldCompareRule 是否为同结构体字段比较规则
func IsLocalFieldCompareRule(name string) bool {
	switch name {
	case RuleEqField, RuleNeField, RuleGTField, RuleAfterField, RuleGTEField, RuleLTField, RuleBeforeField, RuleLTEField:
		return true
	default:
		return false
	}
}

// IsCrossStructFieldCompareRule 是否为跨结构体字段比较规则
func IsCrossStructFieldCompareRule(name string) bool {
	switch name {
	case RuleEqCSField, RuleNeCSField, RuleGTCSField, RuleGTECSField, RuleLTCSField, RuleLTECSField:
		return true
	default:
		return false
	}
}

// IsOmitEmptyRule 是否为omitempty规则
func IsOmitEmptyRule(name string) bool {
	return name == RuleOmitEmpty || name == RuleOmitZero
}

// IsStructControlRule 是否为结构体控制规则
func IsStructControlRule(name string) bool {
	switch name {
	case RuleStructOnly, RuleNoStructLevel, RuleEmpty:
		return true
	default:
		return false
	}
}

// IsDiveControlRule 是否为递归控制规则
func IsDiveControlRule(name string) bool {
	switch name {
	case RuleKeys, RuleEndKeys, RuleStructOnly, RuleNoStructLevel, RuleEmpty:
		return true
	default:
		return false
	}
}

// StopsStructDive 是否停止结构体递归
func StopsStructDive(name string) bool {
	switch name {
	case RuleDive, RuleNoStructLevel, RuleStructOnly:
		return true
	default:
		return false
	}
}

// CompareOperatorForRule 返回规则对应的比较符
func CompareOperatorForRule(name string) string {
	switch name {
	case RuleLen, RuleEqField, RuleEqCSField:
		return RuleEq
	case RuleMin, RuleGTE, RuleGTEField, RuleGTECSField:
		return RuleGTE
	case RuleMax, RuleLTE, RuleLTEField, RuleLTECSField:
		return RuleLTE
	case RuleGT, RuleGTField, RuleGTCSField, RuleAfterField:
		return RuleGT
	case RuleLT, RuleLTField, RuleLTCSField, RuleBeforeField:
		return RuleLT
	case RuleNeField, RuleNeCSField:
		return RuleNe
	default:
		return RuleEmpty
	}
}
