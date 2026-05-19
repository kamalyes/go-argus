/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-19 00:00:00
 * @FilePath: \go-argus\rule\plan.go
 * @Description: 规则计划解析，将 validate 标签预编译为运行计划
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"strings"
)

// RulePlan 表示一条规则的运行计划
type RulePlan struct {
	Name       string     // 规则名称
	Param      string     // 规则参数原始值
	ParamParts []string   // 规则参数按空白拆分
	OrRules    []RulePlan // 或规则列表（对应 | 分隔符）
}

// FieldPlan 表示结构体字段的校验计划
type FieldPlan struct {
	Index          []int  // 字段索引路径
	Name           string // Go 字段名
	AltName        string // 展示名（json tag 或自定义）
	Typ            reflect.Type
	Rules          []RulePlan // 字段上的规则列表
	HasValidate    bool       // 是否包含校验标签
	NsPrefix       string     // 命名空间前缀
	StructNsPrefix string     // 结构体命名空间前缀
}

// StructPlan 表示结构体的编译计划
type StructPlan struct {
	Name   string      // 结构体名称
	Fields []FieldPlan // 字段计划列表
}

// ParseRules 将 validate 标签解析为规则计划列表
func ParseRules(tag string) []RulePlan {
	parsed := ParseTag(tag)
	rules := make([]RulePlan, 0, len(parsed))
	for _, item := range parsed {
		rules = append(rules, ParseRulePlan(item))
	}
	return rules
}

// ParseRulePlan 将单条规则解析为运行计划，支持或规则展开
func ParseRulePlan(item Rule) RulePlan {
	raw := strings.TrimSpace(item.Raw)
	if raw != "" {
		parts := SplitRuleOr(raw)
		if len(parts) > 1 {
			rp := RulePlan{Name: raw, OrRules: make([]RulePlan, 0, len(parts))}
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				rp.OrRules = append(rp.OrRules, ParseSingleRulePlan(part))
			}
			return rp
		}
	}
	return PrepareRulePlan(RulePlan{Name: item.Name, Param: item.Param})
}

// ParseSingleRulePlan 解析单条规则字符串为运行计划
func ParseSingleRulePlan(raw string) RulePlan {
	name, param, ok := strings.Cut(raw, "=")
	if !ok {
		name = raw
	}
	return PrepareRulePlan(RulePlan{Name: strings.TrimSpace(name), Param: strings.TrimSpace(param)})
}

// PrepareRulePlan 对规则计划做后处理，预拆分需要空白分词的规则参数
func PrepareRulePlan(rp RulePlan) RulePlan {
	switch rp.Name {
	case "oneof", "oneofci", "noneof", "noneofci",
		"required_with", "required_with_all", "required_without", "required_without_all",
		"excluded_with", "excluded_with_all", "excluded_without", "excluded_without_all",
		"required_if", "required_unless", "excluded_if", "excluded_unless":
		if rp.Param != "" {
			rp.ParamParts = strings.Fields(rp.Param)
		}
	}
	return rp
}

// SplitRuleOr 按 | 分隔或规则，支持转义符
func SplitRuleOr(s string) []string {
	start := 0
	escaped := false
	var parts []string
	for i := 0; i < len(s); i++ {
		if escaped {
			escaped = false
			continue
		}
		switch s[i] {
		case '\\':
			escaped = true
		case '|':
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	if len(parts) == 0 {
		return nil
	}
	parts = append(parts, s[start:])
	if strings.Contains(parts[0], "=") {
		for i := 1; i < len(parts); i++ {
			if !strings.Contains(parts[i], "=") {
				return nil
			}
		}
	}
	return parts
}
