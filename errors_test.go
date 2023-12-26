/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-26 00:00:00
 * @FilePath: \go-argus\errors_test.go
 * @Description: errors.go 测试，覆盖 InvalidValidationError、ValidationErrors 和 FieldError
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"reflect"
	"testing"
)

func TestInvalidValidationErrorNilType(t *testing.T) {
	err := &InvalidValidationError{}
	if err.Error() != "validator: (nil)" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestInvalidValidationErrorWithType(t *testing.T) {
	err := &InvalidValidationError{Type: reflect.TypeOf(0)}
	if err.Error() != "validator: (nil int)" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestValidationErrorsEmpty(t *testing.T) {
	var ve ValidationErrors
	if ve.Error() != "" {
		t.Fatalf("expected empty error string, got %q", ve.Error())
	}
}

func TestValidationErrorsMultiple(t *testing.T) {
	ve := ValidationErrors{
		&fieldError{tag: "required", ns: "Name", field: "Name"},
		&fieldError{tag: "email", ns: "Email", field: "Email"},
	}
	s := ve.Error()
	if s == "" {
		t.Fatal("expected non-empty error string")
	}
}

func TestFieldErrorMethods(t *testing.T) {
	fe := &fieldError{
		tag:         "min",
		actualTag:   "min",
		ns:          "User.Name",
		structNs:    "User.Name",
		field:       "name",
		structField: "Name",
		value:       "ab",
		param:       "3",
		kind:        reflect.String,
		typ:         reflect.TypeOf(""),
	}
	if fe.Tag() != "min" {
		t.Fatalf("expected tag min, got %s", fe.Tag())
	}
	if fe.ActualTag() != "min" {
		t.Fatalf("expected actualTag min, got %s", fe.ActualTag())
	}
	if fe.Namespace() != "User.Name" {
		t.Fatalf("expected namespace User.Name, got %s", fe.Namespace())
	}
	if fe.StructNamespace() != "User.Name" {
		t.Fatalf("expected structNs User.Name, got %s", fe.StructNamespace())
	}
	if fe.Field() != "name" {
		t.Fatalf("expected field name, got %s", fe.Field())
	}
	if fe.StructField() != "Name" {
		t.Fatalf("expected structField Name, got %s", fe.StructField())
	}
	if fe.Value() != "ab" {
		t.Fatalf("expected value ab, got %v", fe.Value())
	}
	if fe.Param() != "3" {
		t.Fatalf("expected param 3, got %s", fe.Param())
	}
	if fe.Kind() != reflect.String {
		t.Fatalf("expected kind String, got %v", fe.Kind())
	}
	if fe.Type() != reflect.TypeOf("") {
		t.Fatalf("expected type string, got %v", fe.Type())
	}
	expectedErr := "Key: 'User.Name' Error:Field validation for 'name' failed on the 'min' tag"
	if fe.Error() != expectedErr {
		t.Fatalf("unexpected error: %s", fe.Error())
	}
}

func TestFieldErrorNilValue(t *testing.T) {
	fe := &fieldError{
		tag:         "required",
		ns:          "Test.Field",
		field:       "field",
		structField: "Field",
		kind:        reflect.Invalid,
	}
	if fe.Value() != nil {
		t.Fatal("expected nil value")
	}
	if fe.Kind() != reflect.Invalid {
		t.Fatal("expected Invalid kind")
	}
	if fe.Type() != nil {
		t.Fatal("expected nil type")
	}
}
