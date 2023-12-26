/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-20 00:00:00
 * @FilePath: \go-argus\cache_test.go
 * @Description: cache.go 测试，覆盖结构体编译缓存和字段名解析
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"reflect"
	"testing"
)

type cacheTestStruct struct {
	Name  string `json:"name" validate:"required"`
	Email string `validate:"required,email"`
}

type cacheTestPrivate struct {
	name string `validate:"required"`
}

type cacheTestSkip struct {
	Name  string `json:"name" validate:"required"`
	Skip  string `json:"-" validate:"-"`
	NoTag string
}

func TestCompileStructBasic(t *testing.T) {
	v := New()
	plan := v.compileStruct(reflect.TypeOf(cacheTestStruct{}))
	if plan.name != "cacheTestStruct" {
		t.Fatalf("expected struct name cacheTestStruct, got %s", plan.name)
	}
	if len(plan.fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(plan.fields))
	}
}

func TestCompileStructCaching(t *testing.T) {
	v := New()
	typ := reflect.TypeOf(cacheTestStruct{})
	plan1 := v.compileStruct(typ)
	plan2 := v.compileStruct(typ)
	if plan1 != plan2 {
		t.Fatal("expected same cached plan")
	}
}

func TestCompileStructPrivateField(t *testing.T) {
	v := New()
	plan := v.compileStruct(reflect.TypeOf(cacheTestPrivate{}))
	if len(plan.fields) != 0 {
		t.Fatalf("expected 0 fields without private validation, got %d", len(plan.fields))
	}
	vPrivate := New(WithPrivateFieldValidation())
	planPrivate := vPrivate.compileStruct(reflect.TypeOf(cacheTestPrivate{}))
	if len(planPrivate.fields) != 1 {
		t.Fatalf("expected 1 field with private validation, got %d", len(planPrivate.fields))
	}
}

func TestCompileStructSkipTag(t *testing.T) {
	v := New()
	plan := v.compileStruct(reflect.TypeOf(cacheTestSkip{}))
	for _, f := range plan.fields {
		if f.name == "Skip" {
			t.Fatal("expected Skip field to be skipped")
		}
	}
}

func TestResolveFieldNameWithTagFunc(t *testing.T) {
	v := New()
	v.RegisterTagNameFunc(func(sf reflect.StructField) string {
		return "custom_" + sf.Name
	})
	plan := v.compileStruct(reflect.TypeOf(cacheTestStruct{}))
	if plan.fields[0].altName != "custom_Name" {
		t.Fatalf("expected custom_Name, got %s", plan.fields[0].altName)
	}
}

func TestResolveFieldNameWithJSONTag(t *testing.T) {
	v := New()
	plan := v.compileStruct(reflect.TypeOf(cacheTestStruct{}))
	if plan.fields[0].altName != "name" {
		t.Fatalf("expected name from json tag, got %s", plan.fields[0].altName)
	}
	if plan.fields[1].altName != "Email" {
		t.Fatalf("expected Email (no json tag), got %s", plan.fields[1].altName)
	}
}

func TestParseRules(t *testing.T) {
	rules := parseRules("required,min=2,max=16")
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if rules[0].name != "required" {
		t.Fatalf("expected first rule required, got %s", rules[0].name)
	}
	if rules[1].name != "min" || rules[1].param != "2" {
		t.Fatalf("expected min=2, got %s=%s", rules[1].name, rules[1].param)
	}
}

func TestParseRulesEmpty(t *testing.T) {
	rules := parseRules("")
	if len(rules) != 0 {
		t.Fatalf("expected 0 rules for empty tag, got %d", len(rules))
	}
}
