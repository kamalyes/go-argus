/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-20 00:00:00
 * @FilePath: \go-argus\rule\plan_test.go
 * @Description: plan.go 测试，覆盖规则计划解析和 or 规则展开
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"testing"

	"github.com/kamalyes/go-argus/constants"
)

func TestParseRulesEmpty(t *testing.T) {
	if got := ParseRules(""); len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestParseRulesSingle(t *testing.T) {
	rules := ParseRules("required")
	if len(rules) != 1 || rules[0].Name != "required" {
		t.Fatalf("unexpected: %v", rules)
	}
}

func TestParseRulesWithParam(t *testing.T) {
	rules := ParseRules("min=3,max=10")
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Name != "min" || rules[0].Param != "3" {
		t.Fatalf("unexpected first: %v", rules[0])
	}
}

func TestParseRulePlanOrRules(t *testing.T) {
	item := Rule{Name: "oneof", Param: "a|b|c", Raw: "a|b|c"}
	rp := ParseRulePlan(item)
	if len(rp.OrRules) != 3 {
		t.Fatalf("expected 3 or rules, got %d", len(rp.OrRules))
	}
	if rp.OrRules[0].Name != "a" || rp.OrRules[1].Name != "b" || rp.OrRules[2].Name != "c" {
		t.Fatalf("unexpected or rules: %v", rp.OrRules)
	}
}

func TestParseRulePlanNoOr(t *testing.T) {
	item := Rule{Name: "required", Param: "", Raw: "required"}
	rp := ParseRulePlan(item)
	if len(rp.OrRules) != 0 {
		t.Fatalf("expected no or rules, got %d", len(rp.OrRules))
	}
	if rp.Name != "required" {
		t.Fatalf("expected name=required, got %s", rp.Name)
	}
}

func TestParseRulePlanEmptyRaw(t *testing.T) {
	item := Rule{Name: "required", Param: "", Raw: ""}
	rp := ParseRulePlan(item)
	if rp.Name != "required" {
		t.Fatalf("expected name=required, got %s", rp.Name)
	}
}

func TestParseSingleRulePlanWithParam(t *testing.T) {
	rp := ParseSingleRulePlan("min=3")
	if rp.Name != "min" || rp.Param != "3" {
		t.Fatalf("unexpected: %v", rp)
	}
}

func TestParseSingleRulePlanNoParam(t *testing.T) {
	rp := ParseSingleRulePlan("required")
	if rp.Name != "required" || rp.Param != "" {
		t.Fatalf("unexpected: %v", rp)
	}
}

func TestPrepareRulePlanOneOf(t *testing.T) {
	rp := PrepareRulePlan(RulePlan{Name: "oneof", Param: "a b c"})
	if len(rp.ParamParts) != 3 || rp.ParamParts[0] != "a" {
		t.Fatalf("unexpected param parts: %v", rp.ParamParts)
	}
}

func TestPrepareRulePlanRequiredWith(t *testing.T) {
	rp := PrepareRulePlan(RulePlan{Name: "required_with", Param: "A B"})
	if len(rp.ParamParts) != 2 {
		t.Fatalf("unexpected param parts: %v", rp.ParamParts)
	}
}

func TestPrepareRulePlanNoSplit(t *testing.T) {
	rp := PrepareRulePlan(RulePlan{Name: "min", Param: "3"})
	if len(rp.ParamParts) != 0 {
		t.Fatalf("expected no param parts for min, got %v", rp.ParamParts)
	}
}

func TestPrepareRulePlanCmpOp(t *testing.T) {
	rp := PrepareRulePlan(RulePlan{Name: "afterfield", Param: "StartTime"})
	if !rp.HasCmpOp || rp.CmpOp != constants.CmpGT {
		t.Fatalf("expected afterfield to precompute gt, got has=%v op=%v", rp.HasCmpOp, rp.CmpOp)
	}
}

func TestCmpOpForRule(t *testing.T) {
	cases := map[string]constants.CmpOp{
		"len":         constants.CmpEQ,
		"min":         constants.CmpGTE,
		"max":         constants.CmpLTE,
		"gt":          constants.CmpGT,
		"gte":         constants.CmpGTE,
		"lt":          constants.CmpLT,
		"lte":         constants.CmpLTE,
		"eqfield":     constants.CmpEQ,
		"nefield":     constants.CmpNE,
		"gtfield":     constants.CmpGT,
		"afterfield":  constants.CmpGT,
		"gtefield":    constants.CmpGTE,
		"ltfield":     constants.CmpLT,
		"beforefield": constants.CmpLT,
		"ltefield":    constants.CmpLTE,
		"eqcsfield":   constants.CmpEQ,
		"necsfield":   constants.CmpNE,
		"gtcsfield":   constants.CmpGT,
		"gtecsfield":  constants.CmpGTE,
		"ltcsfield":   constants.CmpLT,
		"ltecsfield":  constants.CmpLTE,
	}
	for name, want := range cases {
		if got := CmpOpForRule(name); got != want {
			t.Fatalf("expected %s to map to %v, got %v", name, want, got)
		}
	}
	if got := CmpOpForRule("unknown"); got != constants.CmpOp(-1) {
		t.Fatalf("expected unknown rule to map to -1, got %v", got)
	}
}

func TestPrepareRulePlanEmptyParam(t *testing.T) {
	rp := PrepareRulePlan(RulePlan{Name: "oneof", Param: ""})
	if len(rp.ParamParts) != 0 {
		t.Fatalf("expected no param parts for empty param, got %v", rp.ParamParts)
	}
}

func TestSplitRuleOrSingle(t *testing.T) {
	parts := SplitRuleOr("required")
	if parts != nil {
		t.Fatalf("expected nil for single rule, got %v", parts)
	}
}

func TestSplitRuleOrMultiple(t *testing.T) {
	parts := SplitRuleOr("a|b|c")
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d: %v", len(parts), parts)
	}
}

func TestSplitRuleOrWithEquals(t *testing.T) {
	parts := SplitRuleOr("a=x|b=y")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d: %v", len(parts), parts)
	}
}

func TestSplitRuleOrWithEqualsMissingInOne(t *testing.T) {
	// 如果第一个有 = 但后续没有，返回 nil
	parts := SplitRuleOr("a=x|b")
	if parts != nil {
		t.Fatalf("expected nil for mixed = pattern, got %v", parts)
	}
}

func TestSplitRuleOrEscaped(t *testing.T) {
	parts := SplitRuleOr(`a\|b`)
	if parts != nil {
		t.Fatalf("expected nil for escaped pipe (single part), got %v", parts)
	}
}

func TestFieldPlan(t *testing.T) {
	fp := FieldPlan{
		Index:       []int{0},
		Name:        "Name",
		AltName:     "name",
		Typ:         reflect.TypeOf(""),
		Rules:       ParseRules("required"),
		HasValidate: true,
		NsPrefix:    "Test.name",
	}
	if !fp.HasValidate {
		t.Fatal("expected HasValidate=true")
	}
	if len(fp.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(fp.Rules))
	}
}

func TestStructPlan(t *testing.T) {
	sp := StructPlan{
		Name:   "Test",
		Fields: []FieldPlan{{Name: "Name"}},
	}
	if sp.Name != "Test" || len(sp.Fields) != 1 {
		t.Fatalf("unexpected: %v", sp)
	}
}
