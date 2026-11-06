/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 16:52:13
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 18:21:07
 * @FilePath: \go-argus\code_test.go
 * @Description: 业务错误码测试，覆盖 code 标签、内联码、规则级兜底与取码入口
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"fmt"
	"testing"
)

// claimsCodes 对齐真实业务里的「Token 字段缺失」码段，纯标签声明，无需手写 if 链
type claimsCodes struct {
	OrgID       string `json:"org_id" validate:"required" code:"110026"`
	WorkspaceID string `json:"workspace_id" validate:"required" code:"110027"`
	TeamID      string `json:"team_id" validate:"required" code:"110028"`
	AccountID   string `json:"account_id" validate:"required" code:"110029"`
	SessionID   string `json:"session_id" validate:"required" code:"110003"`
	Version     int    `json:"version" validate:"gt=0" code:"110003"`
}

func fullClaims() claimsCodes {
	return claimsCodes{
		OrgID:       "o_1",
		WorkspaceID: "ws_1",
		TeamID:      "tm_1",
		AccountID:   "a_1",
		SessionID:   "sess_1",
		Version:     1,
	}
}

func TestStructCodeTagPerField(t *testing.T) {
	v := New()
	cases := []struct {
		name   string
		mutate func(*claimsCodes)
		want   int64
	}{
		{"org missing", func(c *claimsCodes) { c.OrgID = " " }, 110026},
		{"workspace missing", func(c *claimsCodes) { c.WorkspaceID = "" }, 110027},
		{"team missing", func(c *claimsCodes) { c.TeamID = "" }, 110028},
		{"account missing", func(c *claimsCodes) { c.AccountID = "" }, 110029},
		{"session missing", func(c *claimsCodes) { c.SessionID = "" }, 110003},
		{"version not positive", func(c *claimsCodes) { c.Version = 0 }, 110003},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			claims := fullClaims()
			tc.mutate(&claims)
			err := v.Struct(claims)
			if err == nil {
				t.Fatal("expected validation error")
			}
			got, ok := CodeOf(err)
			if !ok || got != tc.want {
				t.Fatalf("CodeOf = (%d, %v), want (%d, true)", got, ok, tc.want)
			}
		})
	}
	if err := v.Struct(fullClaims()); err != nil {
		t.Fatalf("full claims should pass: %v", err)
	}
}

func TestStructCodeTagPerRule(t *testing.T) {
	type password struct {
		Value string `validate:"required,min=8,max=16" code:"required=110007,min=110008,max=110009"`
	}
	v := New()
	cases := []struct {
		value string
		want  int64
	}{
		{"", 110007},
		{"short", 110008},
		{"0123456789abcdefg", 110009},
	}
	for _, tc := range cases {
		got, ok := CodeOf(v.Struct(password{Value: tc.value}))
		if !ok || got != tc.want {
			t.Fatalf("value %q: CodeOf = (%d, %v), want (%d, true)", tc.value, got, ok, tc.want)
		}
	}
	if err := v.Struct(password{Value: "0123456789"}); err != nil {
		t.Fatalf("valid password should pass: %v", err)
	}
}

// TestStructCodeTagDefaultWithOverride 裸数字项当默认码，按规则项只覆盖自己那条
func TestStructCodeTagDefaultWithOverride(t *testing.T) {
	type form struct {
		Name string `validate:"required,min=2" code:"100002,min=110008"`
	}
	v := New()
	if got, _ := CodeOf(v.Struct(form{Name: ""})); got != 100002 {
		t.Fatalf("required code = %d, want 100002", got)
	}
	if got, _ := CodeOf(v.Struct(form{Name: "x"})); got != 110008 {
		t.Fatalf("min code = %d, want 110008", got)
	}
}

func TestStructInlineCode(t *testing.T) {
	type form struct {
		OrgID string `validate:"required,code=110026"`
	}
	got, ok := CodeOf(New().Struct(form{}))
	if !ok || got != 110026 {
		t.Fatalf("CodeOf = (%d, %v), want (110026, true)", got, ok)
	}
}

func TestStructInlineCodePerRule(t *testing.T) {
	type form struct {
		Pass string `validate:"required,min=8,code=required=110007,code=min=110008"`
	}
	v := New()
	if got, _ := CodeOf(v.Struct(form{Pass: ""})); got != 110007 {
		t.Fatalf("required code = %d, want 110007", got)
	}
	if got, _ := CodeOf(v.Struct(form{Pass: "short"})); got != 110008 {
		t.Fatalf("min code = %d, want 110008", got)
	}
}

// TestStructCodeTagBeatsInline 同级比较：code 标签覆盖 validate 内联码
func TestStructCodeTagBeatsInline(t *testing.T) {
	type form struct {
		OrgID string `validate:"required,code=100002" code:"110026"`
	}
	if got, _ := CodeOf(New().Struct(form{})); got != 110026 {
		t.Fatalf("code = %d, want 110026 (code tag wins)", got)
	}
}

// TestInlineCodeNotEvaluatedAsRule code= 是声明不是规则，不能让字段恒定失败
func TestInlineCodeNotEvaluatedAsRule(t *testing.T) {
	type form struct {
		OrgID   string `validate:"required,code=110026"`
		Garbage string `validate:"required,code=abc"`
	}
	if err := New().Struct(form{OrgID: "o_1", Garbage: "g"}); err != nil {
		t.Fatalf("valid struct should pass, got %v", err)
	}
}

func TestCodeTagInvalidValueIgnored(t *testing.T) {
	type form struct {
		OrgID string `validate:"required" code:"not-a-code"`
	}
	err := New().Struct(form{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if got, ok := CodeOf(err); ok || got != 0 {
		t.Fatalf("CodeOf = (%d, %v), want (0, false)", got, ok)
	}
}

func TestCodeSkipsControlRules(t *testing.T) {
	type form struct {
		Nick string `validate:"omitempty,min=3" code:"110008"`
	}
	v := New()
	if err := v.Struct(form{Nick: ""}); err != nil {
		t.Fatalf("omitempty should skip blank value: %v", err)
	}
	if got, _ := CodeOf(v.Struct(form{Nick: "ab"})); got != 110008 {
		t.Fatalf("min code = %d, want 110008", got)
	}
}

func TestCodeOnDiveElements(t *testing.T) {
	type form struct {
		Tags []string `validate:"required,dive,required" code:"110026"`
	}
	v := New()
	if got, _ := CodeOf(v.Struct(form{Tags: nil})); got != 110026 {
		t.Fatalf("slice required code = %d, want 110026", got)
	}
	if got, _ := CodeOf(v.Struct(form{Tags: []string{"a", ""}})); got != 110026 {
		t.Fatalf("element required code = %d, want 110026", got)
	}
}

func TestVarStringInlineCode(t *testing.T) {
	v := New()
	err := v.VarString("", "required,code=110026")
	if err == nil {
		t.Fatal("expected validation error")
	}
	// VarString 返回裸 FieldError，不包 ValidationErrors
	got, ok := CodeOf(err)
	if !ok || got != 110026 {
		t.Fatalf("CodeOf = (%d, %v), want (110026, true)", got, ok)
	}
	if err := v.VarString("x", "required,code=110026"); err != nil {
		t.Fatalf("valid value should pass: %v", err)
	}
}

func TestVarInlineCode(t *testing.T) {
	v := New()
	if got, _ := CodeOf(v.Var("", "required,code=110026")); got != 110026 {
		t.Fatalf("string path code = %d, want 110026", got)
	}
	if got, _ := CodeOf(v.Var((*string)(nil), "required,code=110026")); got != 110026 {
		t.Fatalf("nil pointer path code = %d, want 110026", got)
	}
	if got, _ := CodeOf(v.Var(0, "gt=0,code=110026")); got != 110026 {
		t.Fatalf("reflect path code = %d, want 110026", got)
	}
}

func TestRegisterRuleCodeFallback(t *testing.T) {
	type form struct {
		Name  string `validate:"required"`
		OrgID string `validate:"required" code:"110026"`
	}
	v := New()
	if err := v.RegisterRuleCode("required", 100002); err != nil {
		t.Fatalf("RegisterRuleCode: %v", err)
	}
	// 无标签码的字段取兜底码
	if got, _ := CodeOf(v.Struct(form{OrgID: "o_1"})); got != 100002 {
		t.Fatalf("fallback code = %d, want 100002", got)
	}
	// 标签码优先于兜底码
	if got, _ := CodeOf(v.Struct(form{Name: "n"})); got != 110026 {
		t.Fatalf("tag code = %d, want 110026", got)
	}
}

// TestRegisterRuleCodeInvalidatesPlanCache 码烘在编译期计划里，注册必须让旧计划失效
func TestRegisterRuleCodeInvalidatesPlanCache(t *testing.T) {
	type form struct {
		Name string `validate:"required"`
	}
	v := New()
	if got, ok := CodeOf(v.Struct(form{})); ok || got != 0 {
		t.Fatalf("CodeOf before register = (%d, %v), want (0, false)", got, ok)
	}
	if err := v.RegisterRuleCode("required", 100002); err != nil {
		t.Fatalf("RegisterRuleCode: %v", err)
	}
	if got, _ := CodeOf(v.Struct(form{})); got != 100002 {
		t.Fatalf("CodeOf after register = %d, want 100002", got)
	}

	// Var 规则缓存与分槽旁路缓存同样要失效
	if got, ok := CodeOf(v.Var(0, "gt=0")); ok || got != 0 {
		t.Fatalf("var CodeOf before register = (%d, %v), want (0, false)", got, ok)
	}
	if err := v.RegisterRuleCodes(map[string]int64{"gt": 100003}); err != nil {
		t.Fatalf("RegisterRuleCodes: %v", err)
	}
	if got, _ := CodeOf(v.Var(0, "gt=0")); got != 100003 {
		t.Fatalf("var CodeOf after register = %d, want 100003", got)
	}
}

// TestRegisterRuleCodeOnStringFastPath tag 恰为 required 的零分配快速路径也要带兜底码
func TestRegisterRuleCodeOnStringFastPath(t *testing.T) {
	v := New()
	if err := v.RegisterRuleCodes(map[string]int64{"required": 100002, "endswith": 100004}); err != nil {
		t.Fatalf("RegisterRuleCodes: %v", err)
	}
	if got, _ := CodeOf(v.VarString("", "required")); got != 100002 {
		t.Fatalf("VarString required code = %d, want 100002", got)
	}
	if got, _ := CodeOf(v.Var("", "required")); got != 100002 {
		t.Fatalf("Var required code = %d, want 100002", got)
	}
	s := ""
	if got, _ := CodeOf(v.Var(&s, "required")); got != 100002 {
		t.Fatalf("Var *string required code = %d, want 100002", got)
	}
	if got, _ := CodeOf(v.Var("abc", "endswith=z")); got != 100004 {
		t.Fatalf("Var endswith code = %d, want 100004", got)
	}
	v2 := "abc"
	if got, _ := CodeOf(v.Var(&v2, "endswith=z")); got != 100004 {
		t.Fatalf("Var *string endswith code = %d, want 100004", got)
	}
}

func TestRegisterRuleCodeRejectsEmpty(t *testing.T) {
	v := New()
	if err := v.RegisterRuleCode("", 100002); err == nil {
		t.Fatal("empty tag should be rejected")
	}
	if err := v.RegisterRuleCode("required", 0); err == nil {
		t.Fatal("zero code should be rejected")
	}
	if err := v.RegisterRuleCodes(nil); err == nil {
		t.Fatal("empty map should be rejected")
	}
	// 非法项跳过后整表仍可用
	if err := v.RegisterRuleCodes(map[string]int64{"  ": 1, "required": 100002}); err != nil {
		t.Fatalf("RegisterRuleCodes: %v", err)
	}
	if got, _ := CodeOf(v.VarString("", "required")); got != 100002 {
		t.Fatalf("code = %d, want 100002", got)
	}
}

func TestSetCodeTagName(t *testing.T) {
	type form struct {
		OrgID string `validate:"required" errcode:"110026"`
	}
	v := New()
	// 默认标签名读不到 errcode
	if got, ok := CodeOf(v.Struct(form{})); ok || got != 0 {
		t.Fatalf("CodeOf = (%d, %v), want (0, false)", got, ok)
	}
	v.SetCodeTagName("errcode")
	if got, _ := CodeOf(v.Struct(form{})); got != 110026 {
		t.Fatalf("CodeOf = %d, want 110026", got)
	}
	// 空名不改配置
	v.SetCodeTagName("  ")
	if got, _ := CodeOf(v.Struct(form{})); got != 110026 {
		t.Fatalf("CodeOf = %d, want 110026 after blank SetCodeTagName", got)
	}
}

// TestValidationErrorsCodeFirstNonZero 一次响应只能带一个码，取最先失败的那条
func TestValidationErrorsCodeFirstNonZero(t *testing.T) {
	type form struct {
		A string `validate:"required"`
		B string `validate:"required" code:"110027"`
		C string `validate:"required" code:"110028"`
	}
	err := New().Struct(form{})
	ve, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("err type = %T, want ValidationErrors", err)
	}
	if len(ve) != 3 {
		t.Fatalf("errors = %d, want 3", len(ve))
	}
	if ve[0].Code() != 0 {
		t.Fatalf("A code = %d, want 0", ve[0].Code())
	}
	if got := ve.Code(); got != 110027 {
		t.Fatalf("ValidationErrors.Code = %d, want 110027", got)
	}
}

func TestValidationErrorsCodeNoneDeclared(t *testing.T) {
	type form struct {
		A string `validate:"required"`
	}
	err := New().Struct(form{})
	ve := err.(ValidationErrors)
	if got := ve.Code(); got != 0 {
		t.Fatalf("Code = %d, want 0", got)
	}
	if got, ok := CodeOf(err); ok || got != 0 {
		t.Fatalf("CodeOf = (%d, %v), want (0, false)", got, ok)
	}
	if got := (ValidationErrors{}).Code(); got != 0 {
		t.Fatalf("empty Code = %d, want 0", got)
	}
}

func TestCodeOfNonValidationError(t *testing.T) {
	if got, ok := CodeOf(nil); ok || got != 0 {
		t.Fatalf("CodeOf(nil) = (%d, %v), want (0, false)", got, ok)
	}
	if got, ok := CodeOf(&InvalidValidationError{}); ok || got != 0 {
		t.Fatalf("CodeOf(InvalidValidationError) = (%d, %v), want (0, false)", got, ok)
	}
}

func TestTranslateCarriesCode(t *testing.T) {
	type form struct {
		OrgID string `json:"org_id" validate:"required" code:"110026"`
		Name  string `json:"name" validate:"required"`
	}
	err := New().Struct(form{})
	msgs := TranslateValidationErrors(err, "zh")
	if len(msgs) != 2 {
		t.Fatalf("messages = %d, want 2", len(msgs))
	}
	if msgs[0].Code != 110026 {
		t.Fatalf("org_id code = %d, want 110026", msgs[0].Code)
	}
	if msgs[1].Code != 0 {
		t.Fatalf("name code = %d, want 0", msgs[1].Code)
	}
	if msgs[0].Message == "" {
		t.Fatal("message should still be translated")
	}
}

// TestZeroValueFieldErrorCode 零值构造的 fieldError 取码不 panic
func TestZeroValueFieldErrorCode(t *testing.T) {
	fe := &fieldError{}
	if got := fe.Code(); got != 0 {
		t.Fatalf("Code = %d, want 0", got)
	}
}

// TestCodeOnOrRule or 规则的错误以整条表达式为规则名，按规则细分时键要写成同一条表达式
func TestCodeOnOrRule(t *testing.T) {
	type form struct {
		Color string `validate:"rgb|rgba" code:"rgb|rgba=110002"`
		Theme string `validate:"rgb|rgba" code:"110003"`
	}
	v := New()
	if got, _ := CodeOf(v.Struct(form{Color: "nope", Theme: "rgb(1,2,3)"})); got != 110002 {
		t.Fatalf("or rule by-name code = %d, want 110002", got)
	}
	if got, _ := CodeOf(v.Struct(form{Color: "rgb(1,2,3)", Theme: "nope"})); got != 110003 {
		t.Fatalf("or rule default code = %d, want 110003", got)
	}
}

// TestCodeOfWrappedError 调用方自己包装过的错误也要能取到码
func TestCodeOfWrappedError(t *testing.T) {
	type form struct {
		OrgID string `validate:"required" code:"110026"`
	}
	v := New()

	wrapped := fmt.Errorf("validate claims: %w", v.Struct(form{}))
	if got, ok := CodeOf(wrapped); !ok || got != 110026 {
		t.Fatalf("wrapped ValidationErrors: CodeOf = (%d, %v), want (110026, true)", got, ok)
	}

	// VarString 返回裸 FieldError，包装后走 errors.As 的接口分支
	wrappedField := fmt.Errorf("validate org: %w", v.VarString("", "required,code=110026"))
	if got, ok := CodeOf(wrappedField); !ok || got != 110026 {
		t.Fatalf("wrapped FieldError: CodeOf = (%d, %v), want (110026, true)", got, ok)
	}

	// 包装了不带码的校验错误
	type plain struct {
		Name string `validate:"required"`
	}
	if got, ok := CodeOf(fmt.Errorf("x: %w", v.Struct(plain{}))); ok || got != 0 {
		t.Fatalf("wrapped without code: CodeOf = (%d, %v), want (0, false)", got, ok)
	}
	if got, ok := CodeOf(fmt.Errorf("plain error")); ok || got != 0 {
		t.Fatalf("plain error: CodeOf = (%d, %v), want (0, false)", got, ok)
	}
}
