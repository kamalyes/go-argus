/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 18:00:58
 * @FilePath: \go-argus\constants\rule_compare.go
 * @Description: 标量比较与符号运算规则名常量
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

const (
	RuleLen = "len" // 长度等于
	RuleMin = "min" // 最小值/长度
	RuleMax = "max" // 最大值/长度
	RuleEq  = "eq"  // 等于
	RuleNe  = "ne"  // 不等于
	RuleGT  = "gt"  // 大于
	RuleGTE = "gte" // 大于等于
	RuleLT  = "lt"  // 小于
	RuleLTE = "lte" // 小于等于

	RuleSymbolEqual              = "="  // 符号：等于
	RuleSymbolNotEqual           = "!=" // 符号：不等于
	RuleSymbolGreaterThan        = ">"  // 符号：大于
	RuleSymbolGreaterThanOrEqual = ">=" // 符号：大于等于
	RuleSymbolLessThan           = "<"  // 符号：小于
	RuleSymbolLessThanOrEqual    = "<=" // 符号：小于等于

	RuleAfter    = "after"     // 时间晚于
	RuleBefore   = "before"    // 时间早于
	RuleRange    = "range"     // 范围校验
	RuleRegex    = "regex"     // 正则匹配
	RuleNotEmpty = "not_empty" // 非空
)
