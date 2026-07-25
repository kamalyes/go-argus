/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-21 01:26:05
 * @FilePath: \go-argus\validate\pattern_test.go
 * @Description: pattern.go 测试，覆盖数字/字符/密码模式校验函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"testing"
)

func TestIsIntOrFloat(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"整数", "123", true},
		{"两位小数", "123.45", true},
		{"超过两位小数", "123.456", false},
		{"非数字字符", "abc", false},
		{"空字符串", "", false},
		{"末尾小数点", "12.", true},
		{"仅有小数点", ".", false},
		{"整数部分为零", "0", true},
		{"一位小数", "1.5", true},
		{"无整数部分", ".5", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIntOrFloat(tt.s); got != tt.want {
				t.Errorf("IsIntOrFloat(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsDigits(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"纯数字", "123", true},
		{"空字符串返回true", "", true},
		{"包含字母", "12a", false},
		{"包含小数点", "12.3", false},
		{"单个数字", "0", true},
		{"前导零", "00123", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDigits(tt.s); got != tt.want {
				t.Errorf("IsDigits(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsDigitsLengthN(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want bool
	}{
		{"长度等于n", "123", 3, true},
		{"长度小于n", "12", 3, false},
		{"长度大于n", "1234", 3, false},
		{"包含非数字", "12a", 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDigitsLengthN(tt.s, tt.n); got != tt.want {
				t.Errorf("IsDigitsLengthN(%q, %d) = %v, want %v", tt.s, tt.n, got, tt.want)
			}
		})
	}
}

func TestIsDigitsLengthGeN(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want bool
	}{
		{"长度大于n", "123", 2, true},
		{"长度等于n", "12", 2, true},
		{"长度小于n", "1", 2, false},
		{"包含非数字", "1a", 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDigitsLengthGeN(tt.s, tt.n); got != tt.want {
				t.Errorf("IsDigitsLengthGeN(%q, %d) = %v, want %v", tt.s, tt.n, got, tt.want)
			}
		})
	}
}

func TestIsDigitsLengthMN(t *testing.T) {
	tests := []struct {
		name string
		s    string
		m    int
		n    int
		want bool
	}{
		{"长度在范围内", "123", 2, 4, true},
		{"长度等于下界", "12", 2, 4, true},
		{"长度等于上界", "1234", 2, 4, true},
		{"长度小于下界", "1", 2, 4, false},
		{"长度大于上界", "12345", 2, 4, false},
		{"包含非数字", "12a", 2, 4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDigitsLengthMN(tt.s, tt.m, tt.n); got != tt.want {
				t.Errorf("IsDigitsLengthMN(%q, %d, %d) = %v, want %v", tt.s, tt.m, tt.n, got, tt.want)
			}
		})
	}
}

func TestIsNonZeroLeadingDigits(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"非零开头", "123", true},
		{"零开头", "012", false},
		{"空字符串", "", false},
		{"单个零", "0", false},
		{"单个非零", "1", true},
		{"包含非数字", "12a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNonZeroLeadingDigits(tt.s); got != tt.want {
				t.Errorf("IsNonZeroLeadingDigits(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsPositiveNonZeroInt(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"带正号", "+123", true},
		{"不带正号", "123", true},
		{"负号", "-123", false},
		{"正号零", "+0", false},
		{"空字符串", "", false},
		{"零开头", "012", false},
		{"仅正号", "+", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPositiveNonZeroInt(tt.s); got != tt.want {
				t.Errorf("IsPositiveNonZeroInt(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsNegativeNonZeroInt(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"负整数", "-123", true},
		{"正整数", "123", false},
		{"负零", "-0", false},
		{"负零开头", "-012", false},
		{"仅负号", "-", false},
		{"空字符串", "", false},
		{"负数带字母", "-12a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNegativeNonZeroInt(tt.s); got != tt.want {
				t.Errorf("IsNegativeNonZeroInt(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsDecimalN(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want bool
	}{
		{"正好n位小数", "123.45", 2, true},
		{"无小数部分", "123", 2, true},
		{"少于n位小数", "123.4", 2, false},
		{"多于n位小数", "123.456", 2, false},
		{"空字符串", "", 2, false},
		{"多个小数点", "1.2.3", 1, false},
		{"零位小数", "123.", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDecimalN(tt.s, tt.n); got != tt.want {
				t.Errorf("IsDecimalN(%q, %d) = %v, want %v", tt.s, tt.n, got, tt.want)
			}
		})
	}
}

func TestIsDecimalMN(t *testing.T) {
	tests := []struct {
		name string
		s    string
		m    int
		n    int
		want bool
	}{
		{"小数位数在范围内", "123.45", 1, 3, true},
		{"小数位数超出上界", "123.4567", 1, 3, false},
		{"无小数部分且m为零", "123", 0, 3, true},
		{"无小数部分且m非零", "123", 1, 3, false},
		{"小数位数等于下界", "123.1", 1, 3, true},
		{"小数位数等于上界", "123.123", 1, 3, true},
		{"空字符串", "", 0, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDecimalMN(tt.s, tt.m, tt.n); got != tt.want {
				t.Errorf("IsDecimalMN(%q, %d, %d) = %v, want %v", tt.s, tt.m, tt.n, got, tt.want)
			}
		})
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"纯字母", "Hello", true},
		{"字母加数字", "Hello123", false},
		{"空字符串", "", false},
		{"中文字符", "你好", false},
		{"带空格", "Hello World", false},
		{"单个字母", "a", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlpha(tt.s); got != tt.want {
				t.Errorf("IsAlpha(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsUpperAlpha(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"全大写", "HELLO", true},
		{"混合大小写", "Hello", false},
		{"空字符串", "", false},
		{"全小写", "hello", false},
		{"带数字", "HELLO123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUpperAlpha(tt.s); got != tt.want {
				t.Errorf("IsUpperAlpha(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsLowerAlpha(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"全小写", "hello", true},
		{"混合大小写", "Hello", false},
		{"空字符串", "", false},
		{"全大写", "HELLO", false},
		{"带数字", "hello123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLowerAlpha(tt.s); got != tt.want {
				t.Errorf("IsLowerAlpha(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsAlphanumeric(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"字母加数字", "Hello123", true},
		{"带特殊字符", "Hello!", false},
		{"空字符串", "", false},
		{"纯字母", "Hello", true},
		{"纯数字", "12345", true},
		{"带空格", "Hello 123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlphanumeric(tt.s); got != tt.want {
				t.Errorf("IsAlphanumeric(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsWordChars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"字母数字下划线", "Hello_123", true},
		{"带特殊字符", "Hello!", false},
		{"空字符串", "", false},
		{"纯下划线", "___", true},
		{"带空格", "Hello World", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWordChars(tt.s); got != tt.want {
				t.Errorf("IsWordChars(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsLengthN(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want bool
	}{
		{"中文rune长度等于n", "你好吗", 3, true},
		{"中文rune长度不等于n", "你好", 3, false},
		{"英文长度等于n", "abc", 3, true},
		{"英文长度不等于n", "ab", 3, false},
		{"空字符串n为零", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLengthN(tt.s, tt.n); got != tt.want {
				t.Errorf("IsLengthN(%q, %d) = %v, want %v", tt.s, tt.n, got, tt.want)
			}
		})
	}
}

func TestIsHex(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"小写十六进制", "deadbeef", true},
		{"大写十六进制", "DEADBEEF", true},
		{"非十六进制字符", "xyz", false},
		{"空字符串", "", false},
		{"混合大小写", "DeadBeef", true},
		{"带0x前缀", "0xdead", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHex(tt.s); got != tt.want {
				t.Errorf("IsHex(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestHasSpecialChars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"包含感叹号", "hello!", true},
		{"无特殊字符", "hello", false},
		{"多个特殊字符", "a@b#", true},
		{"空字符串", "", false},
		{"仅特殊字符", "@#$", true},
		{"包含下划线", "hello_world", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasSpecialChars(tt.s); got != tt.want {
				t.Errorf("HasSpecialChars(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestHasDoubleByte(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"中文字符", "你好", true},
		{"纯ASCII", "hello", false},
		{"带重音字符", "héllo", true},
		{"空字符串", "", false},
		{"中英混合", "hello你好", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasDoubleByte(tt.s); got != tt.want {
				t.Errorf("HasDoubleByte(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsEmptyLine(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"空格", "   ", true},
		{"制表符换行", "\t\n", true},
		{"非空白字符", "a", false},
		{"空字符串", "", true},
		{"回车换行", "\r\n", true},
		{"带空格的文本", "  a  ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyLine(tt.s); got != tt.want {
				t.Errorf("IsEmptyLine(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsTimeFormat(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"日期", "2024-01-15", true},
		{"日期时间", "2024-01-15 10:30:00", true},
		{"非日期", "not-a-date", false},
		{"斜杠分隔", "2024/01/15", true},
		{"带毫秒", "2024-01-15 10:30:00.123", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTimeFormat(tt.s); got != tt.want {
				t.Errorf("IsTimeFormat(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsPasswordPattern(t *testing.T) {
	tests := []struct {
		name string
		s    string
		m    int
		n    int
		want bool
	}{
		{"合法密码", "abc123", 6, 10, true},
		{"数字开头", "1abc", 6, 10, false},
		{"带下划线", "abc_123", 6, 10, true},
		{"太短", "abc", 6, 10, false},
		{"太长", "abcdefghijklmnop", 6, 10, false},
		{"特殊字符开头", "_abc123", 6, 10, false},
		{"空字符串", "", 6, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPasswordPattern(tt.s, tt.m, tt.n); got != tt.want {
				t.Errorf("IsPasswordPattern(%q, %d, %d) = %v, want %v", tt.s, tt.m, tt.n, got, tt.want)
			}
		})
	}
}

func TestIsStrongPassword(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"强密码", "Abc123!@", true},
		{"无大写无特殊", "abc123", false},
		{"无小写", "ABC123!", false},
		{"太短", "Ab1!", false},
		{"无数字", "Abcdef!@", false},
		{"无特殊字符", "Abc12345", false},
		{"刚好8位", "Abc123!x", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStrongPassword(tt.s); got != tt.want {
				t.Errorf("IsStrongPassword(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsAllChinese(t *testing.T) {
	tests := []struct {
		name             string
		s                string
		wantIsAllChinese bool
		wantNonChinese   int
	}{
		{"全中文", "你好吗", true, 0},
		{"包含非中文", "你好a", false, 1},
		{"空字符串", "", false, 0},
		{"纯英文", "hello", false, 5},
		{"中文带标点", "你好。", false, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsAll, gotNon := IsAllChinese(tt.s)
			if gotIsAll != tt.wantIsAllChinese || gotNon != tt.wantNonChinese {
				t.Errorf("IsAllChinese(%q) = (%v, %d), want (%v, %d)", tt.s, gotIsAll, gotNon, tt.wantIsAllChinese, tt.wantNonChinese)
			}
		})
	}
}

func TestContainsChineseChars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"中英混合", "hello你好", true},
		{"纯英文", "hello", false},
		{"空字符串", "", false},
		{"纯中文", "你好", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsChineseChars(tt.s); got != tt.want {
				t.Errorf("ContainsChineseChars(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsTrueString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"小写true", "true", true},
		{"大写TRUE", "TRUE", true},
		{"数字1", "1", true},
		{"小写yes", "yes", true},
		{"大写YES", "YES", true},
		{"false", "false", false},
		{"空字符串", "", false},
		{"零", "0", false},
		{"no", "no", false},
		{"混合大小写True", "True", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTrueString(tt.s); got != tt.want {
				t.Errorf("IsTrueString(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
