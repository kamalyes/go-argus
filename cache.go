/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\cache.go
 * @Description: 结构体编译缓存，将字段和 validate 标签预编译为运行计划
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"reflect"
	"strings"

	"github.com/kamalyes/go-argus/rule"
)

type rulePlan struct {
	name       string
	param      string
	paramParts []string
	orRules    []rulePlan
}

type fieldPlan struct {
	index          []int
	name           string
	altName        string
	typ            reflect.Type
	rules          []rulePlan
	hasValidate    bool
	nsPrefix       string
	structNsPrefix string
}

type structPlan struct {
	name   string
	fields []fieldPlan
}

func (v *Validate) compileStruct(t reflect.Type) *structPlan {
	if cached, ok := v.structCache.Load(t); ok {
		return cached.(*structPlan)
	}

	typeName := t.Name()
	plan := &structPlan{name: typeName, fields: make([]fieldPlan, 0, t.NumField())}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		// if sf.Anonymous && sf.Type.Kind() == reflect.Struct {}
		if sf.PkgPath != "" && !v.privateFieldValidation {
			continue
		}

		tag := sf.Tag.Get(v.tagName)
		if tag == "-" {
			continue
		}

		altName := v.resolveFieldName(sf)
		fp := fieldPlan{
			index:          sf.Index,
			name:           sf.Name,
			altName:        altName,
			typ:            sf.Type,
			rules:          parseRules(tag),
			hasValidate:    tag != "",
			nsPrefix:       joinNS(typeName, altName),
			structNsPrefix: joinNS(typeName, sf.Name),
		}
		plan.fields = append(plan.fields, fp)
	}

	actual, _ := v.structCache.LoadOrStore(t, plan)
	return actual.(*structPlan)
}

func (v *Validate) resolveFieldName(sf reflect.StructField) string {
	if v.tagNameFunc != nil {
		if name := v.tagNameFunc(sf); name != "" {
			return name
		}
	}
	if jsonTag := sf.Tag.Get("json"); jsonTag != "" {
		name := strings.Split(jsonTag, ",")[0]
		if name != "" && name != "-" {
			return name
		}
	}
	return sf.Name
}

func parseRules(tag string) []rulePlan {
	parsed := rule.ParseTag(tag)
	rules := make([]rulePlan, 0, len(parsed))
	for _, item := range parsed {
		rules = append(rules, parseRulePlan(item))
	}
	return rules
}

func parseRulePlan(item rule.Rule) rulePlan {
	raw := strings.TrimSpace(item.Raw)
	if raw != "" {
		parts := splitRuleOr(raw)
		if len(parts) > 1 {
			rp := rulePlan{name: raw, orRules: make([]rulePlan, 0, len(parts))}
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				rp.orRules = append(rp.orRules, parseSingleRulePlan(part))
			}
			return rp
		}
	}
	return prepareRulePlan(rulePlan{name: item.Name, param: item.Param})
}

func parseSingleRulePlan(raw string) rulePlan {
	name, param, ok := strings.Cut(raw, "=")
	if !ok {
		name = raw
	}
	return prepareRulePlan(rulePlan{name: strings.TrimSpace(name), param: strings.TrimSpace(param)})
}

func prepareRulePlan(rp rulePlan) rulePlan {
	switch rp.name {
	case "oneof", "oneofci", "noneof", "noneofci",
		"required_with", "required_with_all", "required_without", "required_without_all",
		"excluded_with", "excluded_with_all", "excluded_without", "excluded_without_all",
		"required_if", "required_unless", "excluded_if", "excluded_unless":
		if rp.param != "" {
			rp.paramParts = strings.Fields(rp.param)
		}
	}
	return rp
}

func splitRuleOr(s string) []string {
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
