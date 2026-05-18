/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\rule\rule.go
 * @Description: 规则模块入口，提供 validate 标签解析和规则基础结构
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

// Package rule 提供 Argus 的标签规则、跨字段规则和时间表达式能力
package rule

import "strings"

// Rule 表示一个从 validate 标签中解析出来的规则
type Rule struct {
	Name  string
	Param string
	Raw   string
}

// ParseTag 解析 validate 标签，支持逗号分隔和 name=value 参数
func ParseTag(tag string) []Rule {
	if tag == "" {
		return nil
	}
	parts := splitTag(tag)
	rules := make([]Rule, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		name, param, ok := strings.Cut(part, "=")
		if !ok {
			name = part
		}
		rules = append(rules, Rule{
			Name:  strings.TrimSpace(name),
			Param: strings.TrimSpace(param),
			Raw:   part,
		})
	}
	return rules
}

func splitTag(tag string) []string {
	var parts []string
	start := 0
	escaped := false
	for i := 0; i < len(tag); i++ {
		if escaped {
			escaped = false
			continue
		}
		switch tag[i] {
		case '\\':
			escaped = true
		case ',':
			parts = append(parts, tag[start:i])
			start = i + 1
		}
	}
	parts = append(parts, tag[start:])
	return parts
}
