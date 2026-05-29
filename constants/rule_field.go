/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 00:00:00
 * @FilePath: \go-argus\constants\rule_field.go
 * @Description: 跨字段比较规则名常量
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

const (
	RuleEqField     = "eqfield"     // 等于指定字段
	RuleNeField     = "nefield"     // 不等于指定字段
	RuleGTField     = "gtfield"     // 大于指定字段
	RuleGTEField    = "gtefield"    // 大于等于指定字段
	RuleLTField     = "ltfield"     // 小于指定字段
	RuleLTEField    = "ltefield"    // 小于等于指定字段
	RuleAfterField  = "afterfield"  // 时间晚于指定字段
	RuleBeforeField = "beforefield" // 时间早于指定字段

	RuleEqCSField  = "eqcsfield"  // 等于跨结构体字段
	RuleNeCSField  = "necsfield"  // 不等于跨结构体字段
	RuleGTCSField  = "gtcsfield"  // 大于跨结构体字段
	RuleGTECSField = "gtecsfield" // 大于等于跨结构体字段
	RuleLTCSField  = "ltcsfield"  // 小于跨结构体字段
	RuleLTECSField = "ltecsfield" // 小于等于跨结构体字段

	RuleFieldContains = "fieldcontains" // 字段包含指定字段值
	RuleFieldExcludes = "fieldexcludes" // 字段不包含指定字段值
)
