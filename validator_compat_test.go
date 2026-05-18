/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-28 00:00:00
 * @FilePath: \go-argus\validator_compat_test.go
 * @Description: 从 validator 兼容能力中抽取零依赖规则测试，覆盖常用字段、跨字段、集合和 i18n 输出。
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"reflect"
	"strings"
	"testing"
)

type compatProfile struct {
	Name       string            `json:"name" validate:"required,min=2,max=16,alphanumunicode"`
	Email      string            `json:"email" validate:"required,email"`
	Age        int               `json:"age" validate:"gte=18,lte=120"`
	Password   string            `json:"password" validate:"required,min=8"`
	Confirm    string            `json:"confirm" validate:"required,eqfield=password"`
	Website    string            `json:"website" validate:"omitempty,http_url"`
	Role       string            `json:"role" validate:"oneof=admin member guest"`
	Tags       []string          `json:"tags" validate:"min=1,dive,required,lowercase"`
	Meta       map[string]string `json:"meta" validate:"omitempty,dive,required"`
	TraceID    string            `json:"trace_id" validate:"omitempty,uuid4"`
	RemoteIP   string            `json:"remote_ip" validate:"omitempty,ip"`
	RemoteCIDR string            `json:"remote_cidr" validate:"omitempty,cidr"`
}

func TestValidatorCompatibilityCoreRules(t *testing.T) {
	v := New()
	err := v.Struct(compatProfile{
		Name:       "Argus用户",
		Email:      "argus@example.com",
		Age:        30,
		Password:   "secret123",
		Confirm:    "secret123",
		Website:    "https://example.com",
		Role:       "admin",
		Tags:       []string{"gateway", "validator"},
		Meta:       map[string]string{"env": "prod"},
		TraceID:    "550e8400-e29b-41d4-a716-446655440000",
		RemoteIP:   "10.0.0.1",
		RemoteCIDR: "10.0.0.0/8",
	})
	if err != nil {
		t.Fatalf("expected compat profile to pass: %v", err)
	}
}

func TestValidatorCompatibilityCoreFailures(t *testing.T) {
	v := New()
	err := v.Struct(compatProfile{
		Name:     "A",
		Email:    "bad-email",
		Age:      17,
		Password: "secret123",
		Confirm:  "different",
		Role:     "root",
		Tags:     []string{"Gateway"},
	})
	if err == nil {
		t.Fatal("expected invalid profile to fail")
	}
	validationErrors, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	if len(validationErrors) < 6 {
		t.Fatalf("expected multiple field errors, got %d: %v", len(validationErrors), validationErrors)
	}
}

func TestValidatorCompatibilityConditionalRules(t *testing.T) {
	type request struct {
		Mode        string `json:"mode" validate:"required,oneof=sync async"`
		CallbackURL string `json:"callback_url" validate:"required_if=mode async,http_url"`
		Token       string `json:"token" validate:"required_without=callback_url"`
		Debug       string `json:"debug" validate:"excluded_if=mode sync"`
	}

	v := New()
	if err := v.Struct(request{Mode: "async", CallbackURL: "https://example.com/callback"}); err != nil {
		t.Fatalf("expected async request to pass: %v", err)
	}
	err := v.Struct(request{Mode: "async"})
	if err == nil {
		t.Fatal("expected callback_url or token to be required")
	}
	err = v.Struct(request{Mode: "sync", Debug: "true"})
	if err == nil {
		t.Fatal("expected debug to be excluded in sync mode")
	}
}

func TestValidatorCompatibilityFormatRules(t *testing.T) {
	cases := []struct {
		name  string
		value interface{}
		tag   string
	}{
		{name: "base32", value: "JBSWY3DPEB3W64TMMQ======", tag: "base32"},
		{name: "base64url", value: "YXJndXM=", tag: "base64url"},
		{name: "hexcolor", value: "#12ffaa", tag: "hexcolor"},
		{name: "rgb", value: "rgb(12, 34, 255)", tag: "rgb"},
		{name: "e164", value: "+8613800138000", tag: "e164"},
		{name: "hostname", value: "api.example.com", tag: "hostname"},
		{name: "port", value: "443", tag: "port"},
		{name: "mongodb", value: "507f1f77bcf86cd799439011", tag: "mongodb"},
		{name: "credit_card", value: "4111111111111111", tag: "credit_card"},
	}
	v := New()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Var(tt.value, tt.tag); err != nil {
				t.Fatalf("expected %s to pass %s: %v", tt.value, tt.tag, err)
			}
		})
	}
}

func TestValidatorCompatibilityCustomValidation(t *testing.T) {
	v := New()
	err := v.RegisterValidation("notblank", func(fl FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})
	if err != nil {
		t.Fatalf("register validation failed: %v", err)
	}
	if err := v.Var(" argus ", "notblank"); err != nil {
		t.Fatalf("expected custom validation to pass: %v", err)
	}
	if err := v.Var("   ", "notblank"); err == nil {
		t.Fatal("expected custom validation to fail")
	}
}

func TestValidationErrorsTranslateArray(t *testing.T) {
	type request struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}
	v := New()
	err := v.Struct(request{})
	if err == nil {
		t.Fatal("expected missing fields")
	}
	messages := TranslateValidationErrors(err, "zh-CN")
	if len(messages) != 2 {
		t.Fatalf("expected two translated messages, got %#v", messages)
	}
	if messages[0].Field != "name" || messages[0].Message != "name 为必填字段" {
		t.Fatalf("unexpected first message: %#v", messages[0])
	}
	validationErrors := err.(ValidationErrors)
	if got := validationErrors.MissingFields(); !reflect.DeepEqual(got, []string{"name", "email"}) {
		t.Fatalf("unexpected missing fields: %#v", got)
	}
}
