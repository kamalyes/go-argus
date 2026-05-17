/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-17 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 11:34:00
 * @FilePath: \go-argus\string_field_error_test.go
 * @Description: string_field_error.go 测试，覆盖 stringFieldError 所有方法
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"reflect"
	"testing"
)

func TestStringFieldErrorTag(t *testing.T) {
	e := &stringFieldError{tag: "required", param: "", value: ""}
	if e.Tag() != "required" {
		t.Fatalf("expected tag 'required', got '%s'", e.Tag())
	}
}

func TestStringFieldErrorActualTag(t *testing.T) {
	e := &stringFieldError{tag: "min", param: "3", value: "ab"}
	if e.ActualTag() != "min" {
		t.Fatalf("expected actual tag 'min', got '%s'", e.ActualTag())
	}
}

func TestStringFieldErrorNamespace(t *testing.T) {
	e := &stringFieldError{tag: "required", param: "", value: ""}
	if e.Namespace() != "" {
		t.Fatalf("expected empty namespace, got '%s'", e.Namespace())
	}
}

func TestStringFieldErrorStructNamespace(t *testing.T) {
	e := &stringFieldError{tag: "required", param: "", value: ""}
	if e.StructNamespace() != "" {
		t.Fatalf("expected empty struct namespace, got '%s'", e.StructNamespace())
	}
}

func TestStringFieldErrorField(t *testing.T) {
	e := &stringFieldError{tag: "required", param: "", value: ""}
	if e.Field() != "" {
		t.Fatalf("expected empty field, got '%s'", e.Field())
	}
}

func TestStringFieldErrorStructField(t *testing.T) {
	e := &stringFieldError{tag: "required", param: "", value: ""}
	if e.StructField() != "" {
		t.Fatalf("expected empty struct field, got '%s'", e.StructField())
	}
}

func TestStringFieldErrorValue(t *testing.T) {
	e := &stringFieldError{tag: "required", param: "", value: "hello"}
	if e.Value() != "hello" {
		t.Fatalf("expected value 'hello', got '%v'", e.Value())
	}
}

func TestStringFieldErrorParam(t *testing.T) {
	e := &stringFieldError{tag: "min", param: "3", value: "ab"}
	if e.Param() != "3" {
		t.Fatalf("expected param '3', got '%s'", e.Param())
	}
}

func TestStringFieldErrorKind(t *testing.T) {
	e := &stringFieldError{tag: "required", param: "", value: ""}
	if e.Kind() != reflect.String {
		t.Fatalf("expected kind reflect.String, got %v", e.Kind())
	}
}

func TestStringFieldType(t *testing.T) {
	e := &stringFieldError{tag: "required", param: "", value: ""}
	if e.Type() != reflect.TypeOf("") {
		t.Fatalf("expected type string, got %v", e.Type())
	}
}

func TestStringFieldErrorError(t *testing.T) {
	e := &stringFieldError{tag: "required", param: "", value: ""}
	errStr := e.Error()
	if errStr == "" {
		t.Fatal("expected non-empty error string")
	}
}

func TestStringFieldErrorImplementsFieldError(t *testing.T) {
	var _ FieldError = (*stringFieldError)(nil)
}
