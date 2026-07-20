/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-25 23:17:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-25 23:19:09
 * @FilePath: \go-argus\validate\pattern.go
 * @Description: 数字/字符/密码模式校验，全部使用手动字节扫描，零正则零分配
 *
 * 设计原则：
 *   1. 手动 byte 遍历替代 regexp.MatchString，消除正则编译与匹配开销
 *   2. 早退：发现非法字节立即返回 false
 *   3. 函数式 API（无 struct 持有状态），参数化校验直接接受参数
 *   4. ASCII 算术（char - '0'）替代 strconv.Atoi
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"regexp"
	"strings"
	"unicode"
)

// ────────────────────────────────────────
// 数字校验（手动字节扫描）
// ────────────────────────────────────────

// isLowerHexByte 判断字节是否为小写十六进制字符
func isLowerHexByte(c byte) bool { return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') }

// isUpperHexByte 判断字节是否为大写十六进制字符
func isUpperHexByte(c byte) bool { return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') }

// allDigits 检查字符串是否全部由 ASCII 数字组成（空串返回 false）
func allDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// IsIntOrFloat 检查字符串是否为整数或最多2位小数
// 规则：至少1位整数 + 可选小数点 + 0~2位小数
func IsIntOrFloat(s string) bool {
	if len(s) == 0 {
		return false
	}
	i := 0
	// 整数部分：至少1位数字
	for ; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
	}
	if i == 0 {
		return false // 没有整数部分
	}
	// 可选小数部分
	if i < len(s) {
		if s[i] != '.' {
			return false
		}
		i++
		frac := 0
		for ; i < len(s) && s[i] >= '0' && s[i] <= '9' && frac < 2; i++ {
			frac++
		}
		if i != len(s) {
			return false // 小数超过2位或有非法字符
		}
	}
	return true
}

// IsDigits 检查字符串是否全部由数字组成（空串返回 true）
func IsDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// IsDigitsLengthN 检查字符串是否为长度等于 n 的纯数字
func IsDigitsLengthN(s string, n int) bool {
	if len(s) != n {
		return false
	}
	return allDigits(s)
}

// IsDigitsLengthGeN 检查字符串是否为长度不小于 n 的纯数字
func IsDigitsLengthGeN(s string, n int) bool {
	if len(s) < n {
		return false
	}
	return allDigits(s)
}

// IsDigitsLengthMN 检查字符串是否为长度在 m 到 n 之间的纯数字
func IsDigitsLengthMN(s string, m, n int) bool {
	if len(s) < m || len(s) > n {
		return false
	}
	return allDigits(s)
}

// IsNonZeroLeadingDigits 检查字符串是否为非零开头的纯数字
func IsNonZeroLeadingDigits(s string) bool {
	if len(s) == 0 || s[0] == '0' {
		return false
	}
	return allDigits(s)
}

// IsPositiveNonZeroInt 检查字符串是否为正整数（可选+前缀，首位非零）
// 规则：\+?[1-9][0-9]*
func IsPositiveNonZeroInt(s string) bool {
	if len(s) == 0 {
		return false
	}
	i := 0
	if s[0] == '+' {
		i = 1
	}
	if i >= len(s) || s[i] == '0' {
		return false
	}
	return allDigits(s[i:])
}

// IsNegativeNonZeroInt 检查字符串是否为负整数（首位非零）
// 规则：\-[1-9][0-9]*
func IsNegativeNonZeroInt(s string) bool {
	if len(s) < 2 || s[0] != '-' || s[1] == '0' {
		return false
	}
	return allDigits(s[1:])
}

// IsDecimalN 检查字符串是否为有 n 位小数的正实数
// 规则：整数部分 + . + n 位小数（或无小数部分）
func IsDecimalN(s string, n int) bool {
	if len(s) == 0 {
		return false
	}
	dotIdx := -1
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] >= '0' && s[i] <= '9':
			// OK
		case s[i] == '.':
			if dotIdx != -1 {
				return false // 多个小数点
			}
			dotIdx = i
		default:
			return false
		}
	}
	if dotIdx == -1 {
		return true // 无小数部分，合法
	}
	fracLen := len(s) - dotIdx - 1
	return fracLen == n
}

// IsDecimalMN 检查字符串是否为有 m 到 n 位小数的正实数
func IsDecimalMN(s string, m, n int) bool {
	if len(s) == 0 {
		return false
	}
	dotIdx := -1
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] >= '0' && s[i] <= '9':
		case s[i] == '.':
			if dotIdx != -1 {
				return false
			}
			dotIdx = i
		default:
			return false
		}
	}
	if dotIdx == -1 {
		return m == 0 // 无小数部分时，仅 m==0 时合法
	}
	fracLen := len(s) - dotIdx - 1
	return fracLen >= m && fracLen <= n
}

// ────────────────────────────────────────
// 字符串/字符校验（手动字节扫描）
// ────────────────────────────────────────

// IsAlpha 检查字符串是否全部由英文字母组成（大小写不限）
func IsAlpha(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}
	}
	return true
}

// IsUpperAlpha 检查字符串是否全部由大写英文字母组成
func IsUpperAlpha(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	return true
}

// IsLowerAlpha 检查字符串是否全部由小写英文字母组成
func IsLowerAlpha(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < 'a' || s[i] > 'z' {
			return false
		}
	}
	return true
}

// IsAlphanumeric 检查字符串是否由数字和英文字母组成
func IsAlphanumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}

// IsWordChars 检查字符串是否由字母、数字、下划线组成（\w+）
func IsWordChars(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

// IsLengthN 检查字符串的 rune 长度是否等于 n
// 注意：按 rune 计数，非 byte 长度
func IsLengthN(s string, n int) bool {
	return len([]rune(s)) == n
}

// IsHex 检查字符串是否为有效的十六进制数
func IsHex(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isLowerHexByte(s[i]) && !isUpperHexByte(s[i]) {
			return false
		}
	}
	return true
}

// HasSpecialChars 检查字符串是否包含特殊字符
// 特殊字符集合：!@#$%^&*()_+[]{}|;':",./<>?
func HasSpecialChars(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '!', '@', '#', '$', '%', '^', '&', '*', '(', ')',
			'_', '+', '[', ']', '{', '}', '|', ';', '\'', ':',
			'"', ',', '.', '/', '<', '>', '?':
			return true
		}
	}
	return false
}

// HasDoubleByte 检查字符串是否包含双字节字符（非 ASCII）
// 用于检测中文、日文、韩文等 CJK 字符
func HasDoubleByte(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return true
		}
	}
	return false
}

// IsEmptyLine 检查字符串是否为空白行（全部由空白字符组成）
func IsEmptyLine(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r':
		default:
			return false
		}
	}
	return true
}

// timeFormatRegex 预编译的时间格式正则（包级变量，仅编译一次）
var timeFormatRegex = regexp.MustCompile(`(\d{4}[-/\.]\d{1,2}[-/\.]\d{1,2})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`)

// IsTimeFormat 校验时间格式
func IsTimeFormat(s string) bool {
	return timeFormatRegex.MatchString(s)
}

// ────────────────────────────────────────
// 密码校验
// ────────────────────────────────────────

// IsPasswordPattern 检查密码是否符合规则：
// 以英文字母开头，由字母、数字、下划线组成，长度在 m 到 n 之间
func IsPasswordPattern(s string, m, n int) bool {
	runeLen := len([]rune(s))
	if runeLen < m || runeLen > n {
		return false
	}
	if len(s) == 0 {
		return false
	}
	// 首字符必须是英文字母
	first := s[0]
	if !((first >= 'A' && first <= 'Z') || (first >= 'a' && first <= 'z')) {
		return false
	}
	return IsWordChars(s)
}

// IsStrongPassword 检查密码强度：
//   - 长度至少 8 个字符
//   - 包含至少一个小写字母
//   - 包含至少一个大写字母
//   - 包含至少一个数字
//   - 包含至少一个特殊字符
//
// 使用手动字节扫描，一次遍历检查所有字符类别，零分配
func IsStrongPassword(s string) bool {
	if len(s) < 8 {
		return false
	}
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= '0' && c <= '9':
			hasDigit = true
		default:
			// 非字母数字视为特殊字符
			if isSpecialPasswordChar(c) {
				hasSpecial = true
			}
		}
		// 早退：所有条件都满足
		if hasLower && hasUpper && hasDigit && hasSpecial {
			return true
		}
	}
	return hasLower && hasUpper && hasDigit && hasSpecial
}

// isSpecialPasswordChar 判断是否为密码强度校验中的特殊字符
func isSpecialPasswordChar(c byte) bool {
	switch c {
	case '!', '@', '#', '~', '$', '%', '^', '&', '*', '(', ')',
		',', '.', '?', '"', ':', '{', '}', '|', '<', '>':
		return true
	}
	return false
}

// ────────────────────────────────────────
// 中文字符校验
// ────────────────────────────────────────

// IsAllChinese 检查字符串是否全部由汉字组成
// 返回值：isAllChinese 表示是否全为汉字，nonChineseCount 表示非汉字字符数
func IsAllChinese(str string) (isAllChinese bool, nonChineseCount int) {
	for _, v := range str {
		if !unicode.Is(unicode.Han, v) {
			nonChineseCount++
		}
	}
	if nonChineseCount == 0 && len(str) > 0 {
		isAllChinese = true
	}
	return
}

// ContainsChineseChars 检查字符串是否包含中文字符
// 注意：go-argus 已有 ContainsChinese 函数（string_rules.go），此函数为其补充
// 如果只需判断是否包含中文，直接使用 validate.ContainsChinese
func ContainsChineseChars(s string) bool {
	for _, v := range s {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

// ────────────────────────────────────────
// 字符串真值判断
// ────────────────────────────────────────

// IsTrueString 判断字符串是否表示 true（true/1/yes，大小写不敏感）
// 使用 strings.EqualFold 避免分配
func IsTrueString(s string) bool {
	return strings.EqualFold(s, "true") ||
		s == "1" ||
		strings.EqualFold(s, "yes")
}
