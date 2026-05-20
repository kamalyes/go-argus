/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-19 00:00:00
 * @FilePath: \go-argus\rule\string_rules.go
 * @Description: 字符串规则映射，委托 validate 子包
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"github.com/kamalyes/go-argus/validate"
)

// StringRuleFunc 字符串规则函数签名
type StringRuleFunc = validate.StringRuleFunc

// StringRuleMap 字符串规则映射表，VarString 快速路径直接查表
var StringRuleMap = validate.StringRuleMap
