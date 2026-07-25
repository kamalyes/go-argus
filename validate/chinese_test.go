/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-21 01:25:56
 * @FilePath: \go-argus\validate\chinese_test.go
 * @Description: chinese.go 测试，覆盖中国大陆手机号、身份证号校验
 *
 * 校验和计算说明：
 *   - "11010119900307653" 的加权和为 226，226%11=6，checkMap[6]='6'
 *   - "11010119900307854" 的加权和为 244，244%11=2，checkMap[2]='X'
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"testing"
)

func TestIsChinesePhoneNumber(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"标准手机号", "13800138000", true},
		{"158开头", "15812345678", true},
		{"第二位为2", "12800138000", false},
		{"10位号码", "1380013800", false},
		{"12位号码", "138001380001", false},
		{"首位非1", "23800138000", false},
		{"包含字母", "1380013800a", false},
		{"199开头", "19912345678", true},
		{"170开头", "17012345678", true},
		{"空字符串", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChinesePhoneNumber(tt.s); got != tt.want {
				t.Errorf("IsChinesePhoneNumber(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsChineseIDCard(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"18位身份证", "110101199003076534", true},
		{"18位带X", "11010119900307653X", true},
		{"18位带小写x", "11010119900307653x", true},
		{"15位身份证", "110101900307001", true},
		{"19位", "1101011990030765345", false},
		{"过短", "12345", false},
		{"空字符串", "", false},
		{"18位含字母", "11010119900307653Y", false},
		{"15位含字母", "1101011990030a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChineseIDCard(tt.s); got != tt.want {
				t.Errorf("IsChineseIDCard(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestCalculateIDCardChecksum(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		// "11010119900307653" 加权和=226, 226%11=6, checkMap[6]='6'
		{"校验码为6", "11010119900307653", "6"},
		// "11010119900307854" 加权和=244, 244%11=2, checkMap[2]='X'
		{"校验码为X", "11010119900307854", "X"},
		// "11010119900307855" 加权和=246, 246%11=4, checkMap[4]='8'
		{"校验码为8", "11010119900307855", "8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateIDCardChecksum(tt.id); got != tt.want {
				t.Errorf("CalculateIDCardChecksum(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestIsChineseIDCardWithChecksum(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		// "11010119900307653" 的正确校验码为 '6'，故 "110101199003076536" 为有效身份证
		{"正确校验码", "110101199003076536", true},
		// 最后一位应为 6，此处为 4，校验失败
		{"错误校验码", "110101199003076534", false},
		// "11010119900307854" 的正确校验码为 'X'，故 "11010119900307854X" 为有效身份证
		{"正确校验码带X", "11010119900307854X", true},
		// 小写 x 同样接受
		{"正确校验码带小写x", "11010119900307854x", true},
		// "11010119900307854" 的正确校验码为 X，此处给 0，校验失败
		{"错误校验码带X位置", "110101199003078540", false},
		{"长度不为18", "11010119900307", false},
		{"前17位含字母", "1101011990030765a4", false},
		{"最后一位非法字符", "11010119900307653Y", false},
		{"空字符串", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChineseIDCardWithChecksum(tt.s); got != tt.want {
				t.Errorf("IsChineseIDCardWithChecksum(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
