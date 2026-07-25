/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-21 00:00:00
 * @FilePath: \go-argus\constants\rule_pattern.go
 * @Description: 模式与中国格式校验规则名常量
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

const (
	RuleDigits              = "digits"           // 纯数字字符串
	RuleIntOrFloat          = "int_or_float"     // 整数或最多2位小数
	RuleWordChars           = "word_chars"       // 字母数字下划线
	RuleUpperAlpha          = "upper_alpha"      // 纯大写字母
	RuleLowerAlpha          = "lower_alpha"      // 纯小写字母
	RuleStrongPassword      = "strong_password"  // 强密码（8位+大小写+数字+特殊字符）
	RuleHasSpecial          = "has_special"      // 包含特殊字符
	RuleHasDoubleByte       = "has_double_byte"  // 包含双字节字符
	RuleEmptyLine           = "empty_line"       // 空白行
	RuleChinesePhone        = "cn_phone"         // 中国大陆手机号
	RuleChineseIDCard       = "cn_idcard"        // 中国大陆身份证号（15或18位）
	RuleChineseIDCardStrict = "cn_idcard_strict" // 身份证号含校验和验证
)
