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
	name  string
	param string
}

type fieldPlan struct {
	index       []int
	name        string
	altName     string
	typ         reflect.Type
	rules       []rulePlan
	hasValidate bool
}

type structPlan struct {
	name   string
	fields []fieldPlan
}

func (v *Validate) compileStruct(t reflect.Type) *structPlan {
	if cached, ok := v.structCache.Load(t); ok {
		return cached.(*structPlan)
	}

	plan := &structPlan{name: t.Name(), fields: make([]fieldPlan, 0, t.NumField())}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.Anonymous && sf.Type.Kind() == reflect.Struct {
			// Embedded structs are handled as normal fields by reflection below.
		}
		if sf.PkgPath != "" && !v.privateFieldValidation {
			continue
		}

		tag := sf.Tag.Get(v.tagName)
		if tag == "-" {
			continue
		}

		fp := fieldPlan{
			index:       sf.Index,
			name:        sf.Name,
			altName:     v.resolveFieldName(sf),
			typ:         sf.Type,
			rules:       parseRules(tag),
			hasValidate: tag != "",
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
		rules = append(rules, rulePlan{name: item.Name, param: item.Param})
	}
	return rules
}
