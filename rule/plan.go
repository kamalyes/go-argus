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

type RulePlan struct {
	Name       string
	Param      string
	ParamParts []string
	OrRules    []RulePlan
}

type FieldPlan struct {
	Index          []int
	Name           string
	AltName        string
	Typ            reflect.Type
	Rules          []RulePlan
	HasValidate    bool
	NsPrefix       string
	StructNsPrefix string
}

type StructPlan struct {
	Name   string
	Fields []FieldPlan
}

func ParseRules(tag string) []RulePlan {
	parsed := ParseTag(tag)
	rules := make([]RulePlan, 0, len(parsed))
	for _, item := range parsed {
		rules = append(rules, ParseRulePlan(item))
	}
	return rules
}

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

func ParseSingleRulePlan(raw string) RulePlan {
	name, param, ok := strings.Cut(raw, "=")
	if !ok {
		name = raw
	}
	return PrepareRulePlan(RulePlan{Name: strings.TrimSpace(name), Param: strings.TrimSpace(param)})
}

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
