/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 00:00:00
 * @FilePath: \go-argus\constants\rule_helpers.go
 * @Description: 规则分类辅助函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

// NeedsParamParts 是否需要将参数按空白拆分为多个部分
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

// IsScalarCompareRule 是否为标量比较规则（min/max/len/gt/gte/lt/lte）
func IsScalarCompareRule(name string) bool {
	switch name {
	case RuleMin, RuleMax, RuleLen, RuleGT, RuleGTE, RuleLT, RuleLTE:
		return true
	default:
		return false
	}
}

// IsFieldCompareRule 是否为跨字段比较规则（含同结构体和跨结构体）
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

// IsOmitEmptyRule 是否为省略空值规则（omitempty/omitzero）
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

// CompareOperatorForRule 返回规则对应的比较操作符字符串
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
