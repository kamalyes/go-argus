/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 00:00:00
 * @FilePath: \go-argus\constants\rule_string.go
 * @Description: 字符串校验规则名常量
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

const (
	RuleAlpha           = "alpha"           // 纯字母
	RuleAlphaSpace      = "alphaspace"      // 字母+空格
	RuleAlphanum        = "alphanum"        // 字母数字
	RuleAlphanumSpace   = "alphanumspace"   // 字母数字+空格
	RuleAlphaUnicode    = "alphaunicode"    // Unicode 字母
	RuleAlphanumUnicode = "alphanumunicode" // Unicode 字母数字

	RuleASCII       = "ascii"       // ASCII 字符
	RulePrintASCII  = "printascii"  // 可打印 ASCII
	RuleMultibyte   = "multibyte"   // 多字节字符
	RuleHexadecimal = "hexadecimal" // 十六进制字符串

	RuleStartsWith    = "startswith"    // 前缀匹配
	RuleEndsWith      = "endswith"      // 后缀匹配
	RuleStartsNotWith = "startsnotwith" // 前缀不匹配
	RuleEndsNotWith   = "endsnotwith"   // 后缀不匹配

	RuleContains     = "contains"     // 包含子串
	RuleContainsAny  = "containsany"  // 包含任一字符
	RuleContainsRune = "containsrune" // 包含指定 rune

	RuleExcludes     = "excludes"     // 不包含子串
	RuleExcludesAll  = "excludesall"  // 不包含任一字符
	RuleExcludesRune = "excludesrune" // 不包含指定 rune

	RuleLowercase = "lowercase" // 全小写
	RuleUppercase = "uppercase" // 全大写

	RuleEqIgnoreCase = "eq_ignore_case" // 忽略大小写等于
	RuleNeIgnoreCase = "ne_ignore_case" // 忽略大小写不等于
)
