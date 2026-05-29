/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 00:00:00
 * @FilePath: \go-argus\constants\rule_conditional.go
 * @Description: 条件与集合规则名常量
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

const (
	RuleOneOf    = "oneof"    // 枚举值（精确匹配）
	RuleOneOfCI  = "oneofci"  // 枚举值（忽略大小写）
	RuleNoneOf   = "noneof"   // 排除值（精确匹配）
	RuleNoneOfCI = "noneofci" // 排除值（忽略大小写）

	RuleRequiredIf         = "required_if"          // 条件必填：指定字段等于目标值时必填
	RuleRequiredUnless     = "required_unless"      // 条件必填：指定字段不等于目标值时必填
	RuleRequiredWith       = "required_with"        // 条件必填：任一指定字段非空时必填
	RuleRequiredWithAll    = "required_with_all"    // 条件必填：所有指定字段非空时必填
	RuleRequiredWithout    = "required_without"     // 条件必填：任一指定字段为空时必填
	RuleRequiredWithoutAll = "required_without_all" // 条件必填：所有指定字段为空时必填

	RuleExcludedIf         = "excluded_if"          // 条件排除：指定字段等于目标值时排除
	RuleExcludedUnless     = "excluded_unless"      // 条件排除：指定字段不等于目标值时排除
	RuleExcludedWith       = "excluded_with"        // 条件排除：任一指定字段非空时排除
	RuleExcludedWithAll    = "excluded_with_all"    // 条件排除：所有指定字段非空时排除
	RuleExcludedWithout    = "excluded_without"     // 条件排除：任一指定字段为空时排除
	RuleExcludedWithoutAll = "excluded_without_all" // 条件排除：所有指定字段为空时排除
)
