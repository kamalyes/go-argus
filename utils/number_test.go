/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-20 00:00:00
 * @FilePath: \go-argus\utils\number_test.go
 * @Description: 数值解析工具测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package utils

import "testing"

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input string
		want  float64
		ok    bool
	}{
		{"42", 42, true},
		{"3.14", 3.14, true},
		{"0", 0, true},
		{"", 0, true}, // 空字符串返回 0 无错误
		{"abc", 0, false},
		{"12.34.56", 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseFloat(tt.input)
		if ok != tt.ok || (tt.ok && got != tt.want) {
			t.Errorf("ParseFloat(%q) = (%f, %v), want (%f, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseFloatStr(t *testing.T) {
	n, err := ParseFloatStr("100")
	if err != nil || n != 100 {
		t.Fatalf("expected 100, got %f err=%v", n, err)
	}
	n, err = ParseFloatStr("1.5")
	if err != nil || n != 1.5 {
		t.Fatalf("expected 1.5, got %f err=%v", n, err)
	}
	_, err = ParseFloatStr("invalid")
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
	// 空字符串返回 0 无错误
	n, err = ParseFloatStr("")
	if err != nil || n != 0 {
		t.Fatalf("expected 0 for empty, got %f err=%v", n, err)
	}
}
