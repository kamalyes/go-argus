/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-28 00:00:00
 * @FilePath: \go-argus\translations_test.go
 * @Description: translations.go 测试，覆盖 i18n 翻译、RequiredMessages 和 MissingFields
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"errors"
	"testing"
)

func TestRegisterTranslation(t *testing.T) {
	RegisterTranslation("en", "custom_tag", "{field} failed custom validation")
}

func TestRegisterTranslations(t *testing.T) {
	RegisterTranslations("en", map[string]string{
		"custom1": "custom message 1",
		"custom2": "custom message 2",
	})
}

func TestTranslateValidationErrorsNil(t *testing.T) {
	result := TranslateValidationErrors(nil, "en")
	if result != nil {
		t.Fatal("expected nil for nil error")
	}
}

func TestTranslateValidationErrorsNonValidation(t *testing.T) {
	result := TranslateValidationErrors(errors.New("some error"), "en")
	if len(result) != 1 {
		t.Fatalf("expected 1 message, got %d", len(result))
	}
	if result[0].Message != "some error" {
		t.Fatalf("expected 'some error', got %s", result[0].Message)
	}
}

func TestTranslateValidationErrorsWithValidationErrors(t *testing.T) {
	type request struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}
	v := New()
	err := v.Struct(request{})
	messages := TranslateValidationErrors(err, "zh-CN")
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
	if messages[0].Field != "name" {
		t.Fatalf("expected field name, got %s", messages[0].Field)
	}
	if messages[0].Message == "" {
		t.Fatal("expected non-empty message")
	}
}

func TestValidationErrorsTranslateEmpty(t *testing.T) {
	var ve ValidationErrors
	result := ve.Translate("en")
	if result != nil {
		t.Fatal("expected nil for empty ValidationErrors")
	}
}

func TestValidationErrorsRequiredMessages(t *testing.T) {
	type request struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"email"`
	}
	v := New()
	err := v.Struct(request{Name: ""})
	ve := err.(ValidationErrors)
	messages := ve.RequiredMessages("zh-CN")
	if len(messages) != 1 {
		t.Fatalf("expected 1 required message, got %d", len(messages))
	}
	if messages[0].Tag != "required" {
		t.Fatalf("expected required tag, got %s", messages[0].Tag)
	}
}

func TestValidationErrorsRequiredMessagesEmpty(t *testing.T) {
	var ve ValidationErrors
	result := ve.RequiredMessages("en")
	if result != nil {
		t.Fatal("expected nil for empty ValidationErrors")
	}
}

func TestValidationErrorsMissingFields(t *testing.T) {
	type request struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required"`
	}
	v := New()
	err := v.Struct(request{})
	ve := err.(ValidationErrors)
	fields := ve.MissingFields()
	if len(fields) != 2 {
		t.Fatalf("expected 2 missing fields, got %d", len(fields))
	}
}

func TestValidationErrorsMissingFieldsEmpty(t *testing.T) {
	var ve ValidationErrors
	result := ve.MissingFields()
	if result != nil {
		t.Fatal("expected nil for empty ValidationErrors")
	}
}

func TestIsRequiredTag(t *testing.T) {
	tests := []struct {
		tag      string
		expected bool
	}{
		{"required", true},
		{"required_if", true},
		{"required_unless", true},
		{"required_with", true},
		{"required_with_all", true},
		{"required_without", true},
		{"required_without_all", true},
		{"email", false},
		{"min", false},
	}
	for _, tt := range tests {
		if isRequiredTag(tt.tag) != tt.expected {
			t.Fatalf("isRequiredTag(%q) = %v, expected %v", tt.tag, !tt.expected, tt.expected)
		}
	}
}

func TestSafeMessageValueNil(t *testing.T) {
	if safeMessageValue(nil) != nil {
		t.Fatal("expected nil for nil value")
	}
}

func TestSafeMessageValueNilChan(t *testing.T) {
	var ch chan int
	if safeMessageValue(ch) != nil {
		t.Fatal("expected nil for nil channel")
	}
}

func TestSafeMessageValueNilMap(t *testing.T) {
	var m map[string]int
	if safeMessageValue(m) != nil {
		t.Fatal("expected nil for nil map")
	}
}

func TestSafeMessageValueNilSlice(t *testing.T) {
	var s []int
	if safeMessageValue(s) != nil {
		t.Fatal("expected nil for nil slice")
	}
}

func TestSafeMessageValueNilPtr(t *testing.T) {
	var p *int
	if safeMessageValue(p) != nil {
		t.Fatal("expected nil for nil pointer")
	}
}

func TestSafeMessageValueNilFunc(t *testing.T) {
	var fn func()
	if safeMessageValue(fn) != nil {
		t.Fatal("expected nil for nil func")
	}
}

func TestSafeMessageValueNilInterface(t *testing.T) {
	var iface interface{}
	if safeMessageValue(iface) != nil {
		t.Fatal("expected nil for nil interface")
	}
}

func TestSafeMessageValueNonNil(t *testing.T) {
	if safeMessageValue("hello") != "hello" {
		t.Fatal("expected original value for non-nil")
	}
	if safeMessageValue(42) != 42 {
		t.Fatal("expected original int value")
	}
}

func TestRenderTranslationWithTemplate(t *testing.T) {
	RegisterTranslation("test", "test_tag", "{field} failed {tag} with param {param}")
	fe := &fieldError{
		tag:   "test_tag",
		ns:    "Test.Field",
		field: "myField",
		param: "5",
	}
	result := renderTranslation("test", fe)
	if result != "myField failed test_tag with param 5" {
		t.Fatalf("unexpected translation: %s", result)
	}
}

func TestRenderTranslationFallbackDefault(t *testing.T) {
	RegisterTranslation("fallback_test", "default", "{field} validation failed")
	fe := &fieldError{
		tag:   "unknown_tag",
		ns:    "Test.Field",
		field: "myField",
	}
	result := renderTranslation("fallback_test", fe)
	if result != "myField validation failed" {
		t.Fatalf("unexpected fallback: %s", result)
	}
}

func TestRenderTranslationNoTemplate(t *testing.T) {
	fe := &fieldError{
		tag:   "missing_tag",
		ns:    "Test.Field",
		field: "myField",
	}
	result := renderTranslation("nonexistent_locale", fe)
	if result == "" {
		t.Fatal("expected non-empty fallback to fe.Error()")
	}
}
