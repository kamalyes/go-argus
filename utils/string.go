/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-18 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-18 00:00:00
 * @FilePath: \go-argus\utils\string.go
 * @Description: 字符串命名转换工具，提供 lowerCamel 和 snake_case 转换
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package utils

import "strings"

func LowerCamel(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func SnakeCase(s string) string {
	var out strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out.WriteByte('_')
		}
		out.WriteRune(r)
	}
	return strings.ToLower(out.String())
}

// JoinNS 拼接命名空间路径，空值一侧直接返回另一侧
func JoinNS(parent string, child string) string {
	if parent == "" {
		return child
	}
	if child == "" {
		return parent
	}
	return parent + "." + child
}
