/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-20 00:00:00
 * @FilePath: \go-argus\utils\string_test.go
 * @Description: 字符串命名转换工具测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package utils

import "testing"

func TestLowerCamel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"UserName", "userName"},
		{"Name", "name"},
		{"A", "a"},
		{"", ""},
		{"ABC", "aBC"},
	}
	for _, tt := range tests {
		if got := LowerCamel(tt.input); got != tt.want {
			t.Errorf("LowerCamel(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"UserName", "user_name"},
		{"ID", "i_d"},
		{"HTTPServer", "h_t_t_p_server"},
		{"", ""},
		{"lower", "lower"},
	}
	for _, tt := range tests {
		if got := SnakeCase(tt.input); got != tt.want {
			t.Errorf("SnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestJoinNS(t *testing.T) {
	tests := []struct {
		parent string
		child  string
		want   string
	}{
		{"Parent", "Child", "Parent.Child"},
		{"", "Child", "Child"},
		{"Parent", "", "Parent"},
		{"", "", ""},
	}
	for _, tt := range tests {
		if got := JoinNS(tt.parent, tt.child); got != tt.want {
			t.Errorf("JoinNS(%q, %q) = %q, want %q", tt.parent, tt.child, got, tt.want)
		}
	}
}
