/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-27 00:00:00
 * @FilePath: \go-argus\options_test.go
 * @Description: options.go 测试，覆盖校验器配置项和 i18n 全局函数
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"
)

func TestWithRequiredStructEnabled(t *testing.T) {
	v := New(WithRequiredStructEnabled())
	if !v.requiredStructEnabled {
		t.Fatal("expected requiredStructEnabled to be true")
	}
}

func TestWithPrivateFieldValidation(t *testing.T) {
	v := New(WithPrivateFieldValidation())
	if !v.privateFieldValidation {
		t.Fatal("expected privateFieldValidation to be true")
	}
}

func TestSetLocaleAndGetLocale(t *testing.T) {
	original := GetLocale()
	SetLocale("zh-CN")
	if GetLocale() != "zh" {
		t.Fatalf("expected zh (normalized from zh-CN), got %s", GetLocale())
	}
	SetLocale("en-US")
	if GetLocale() != "en" {
		t.Fatalf("expected en (normalized from en-US), got %s", GetLocale())
	}
	SetLocale(original)
}

func TestRegisterI18n(t *testing.T) {
	RegisterI18n("en", "test_key", "test template {field}")
}

func TestRegisterI18nMessages(t *testing.T) {
	RegisterI18nMessages("en", map[string]string{
		"test_key1": "template1",
		"test_key2": "template2",
	})
}
