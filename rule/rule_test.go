/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\rule\rule_test.go
 * @Description: rule.go 测试，覆盖标签解析
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"testing"
)

func TestParseTagEmpty(t *testing.T) {
	if got := ParseTag(""); got != nil {
		t.Fatalf("expected nil for empty tag, got %v", got)
	}
}

func TestParseTagSingle(t *testing.T) {
	rules := ParseTag("required")
	if len(rules) != 1 || rules[0].Name != "required" {
		t.Fatalf("unexpected rules: %v", rules)
	}
}

func TestParseTagWithParam(t *testing.T) {
	rules := ParseTag("min=3,max=16")
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Name != "min" || rules[0].Param != "3" {
		t.Fatalf("unexpected first rule: %v", rules[0])
	}
	if rules[1].Name != "max" || rules[1].Param != "16" {
		t.Fatalf("unexpected second rule: %v", rules[1])
	}
}

func TestParseTagEscapedComma(t *testing.T) {
	rules := ParseTag(`contains=a\,b`)
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Name != "contains" || rules[0].Param != "a\\,b" {
		t.Fatalf("unexpected rule: %v", rules[0])
	}
}

func TestParseTagTrimsSpaces(t *testing.T) {
	rules := ParseTag(" required , min = 3 ")
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Name != "required" || rules[1].Name != "min" || rules[1].Param != "3" {
		t.Fatalf("unexpected rules: %v", rules)
	}
}

func TestParseTagSkipsEmpty(t *testing.T) {
	rules := ParseTag("required,,min=1")
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d: %v", len(rules), rules)
	}
}

func TestSplitTagEscapedBackslash(t *testing.T) {
	rules := ParseTag(`contains=a\\b`)
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
}
