/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-26 00:00:00
 * @FilePath: \go-argus\field_level_test.go
 * @Description: field_level.go 测试，覆盖 FieldLevel 接口方法和 wrapFunc
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"context"
	"reflect"
	"testing"
)

func TestWrapFunc(t *testing.T) {
	fn := func(fl FieldLevel) bool {
		return fl.Field().String() != ""
	}
	wrapped := wrapFunc(fn)
	if wrapped == nil {
		t.Fatal("expected non-nil wrapped function")
	}
}

func TestWrapFuncNil(t *testing.T) {
	wrapped := wrapFunc(nil)
	if wrapped != nil {
		t.Fatal("expected nil wrapped function")
	}
}

func TestFieldLevelMethods(t *testing.T) {
	top := reflect.ValueOf("top")
	parent := reflect.ValueOf("parent")
	field := reflect.ValueOf("field_value")
	fl := fieldLevel{
		top:             top,
		parent:          parent,
		field:           field,
		fieldName:       "fieldName",
		structFieldName: "structFieldName",
		tag:             "custom",
		param:           "paramValue",
	}
	if fl.Top().String() != "top" {
		t.Fatal("expected top value")
	}
	if fl.Parent().String() != "parent" {
		t.Fatal("expected parent value")
	}
	if fl.Field().String() != "field_value" {
		t.Fatal("expected field value")
	}
	if fl.FieldName() != "fieldName" {
		t.Fatal("expected fieldName")
	}
	if fl.StructFieldName() != "structFieldName" {
		t.Fatal("expected structFieldName")
	}
	if fl.GetTag() != "custom" {
		t.Fatal("expected tag custom")
	}
	if fl.Param() != "paramValue" {
		t.Fatal("expected param paramValue")
	}
}

func TestRegisterValidationWithFunc(t *testing.T) {
	v := New()
	fn := func(fl FieldLevel) bool {
		return len(fl.Field().String()) > 0
	}
	err := v.RegisterValidation("custom_rule", fn)
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}
	err = v.Var("hello", "custom_rule")
	if err != nil {
		t.Fatalf("expected custom rule to pass: %v", err)
	}
}

func TestRegisterValidationCtxWithFuncCtx(t *testing.T) {
	v := New()
	fn := func(ctx context.Context, fl FieldLevel) bool {
		return fl.Param() != ""
	}
	err := v.RegisterValidationCtx("ctx_rule", fn)
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}
}

func TestRegisterValidationCtxEmptyTag(t *testing.T) {
	v := New()
	err := v.RegisterValidationCtx("", func(ctx context.Context, fl FieldLevel) bool { return true })
	if err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestRegisterValidationCtxNilFunc(t *testing.T) {
	v := New()
	err := v.RegisterValidationCtx("tag", nil)
	if err == nil {
		t.Fatal("expected error for nil function")
	}
}
