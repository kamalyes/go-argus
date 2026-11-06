/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 16:30:52
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 17:58:36
 * @FilePath: \go-argus\rule\code_test.go
 * @Description: code.go 测试，覆盖错误码规格解析、合并优先级与落盘
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package rule

import "testing"

func TestParseCodeSpecDefault(t *testing.T) {
	spec := ParseCodeSpec("110026")
	if spec.Default != 110026 {
		t.Fatalf("default code = %d, want 110026", spec.Default)
	}
	if len(spec.ByRule) != 0 {
		t.Fatalf("ByRule = %v, want empty", spec.ByRule)
	}
	if spec.CodeFor("required") != 110026 {
		t.Fatalf("CodeFor(required) = %d, want 110026", spec.CodeFor("required"))
	}
}

func TestParseCodeSpecByRule(t *testing.T) {
	spec := ParseCodeSpec("required=110007,min=110008,max=110009")
	cases := map[string]int64{"required": 110007, "min": 110008, "max": 110009}
	for name, want := range cases {
		if got := spec.CodeFor(name); got != want {
			t.Fatalf("CodeFor(%s) = %d, want %d", name, got, want)
		}
	}
	if spec.Default != 0 {
		t.Fatalf("default code = %d, want 0", spec.Default)
	}
	// 未列出的规则在没有默认项时取不到码
	if got := spec.CodeFor("email"); got != 0 {
		t.Fatalf("CodeFor(email) = %d, want 0", got)
	}
}

func TestParseCodeSpecMixed(t *testing.T) {
	spec := ParseCodeSpec("110000, min=110008 ")
	if spec.Default != 110000 {
		t.Fatalf("default code = %d, want 110000", spec.Default)
	}
	if got := spec.CodeFor("min"); got != 110008 {
		t.Fatalf("CodeFor(min) = %d, want 110008", got)
	}
	if got := spec.CodeFor("required"); got != 110000 {
		t.Fatalf("CodeFor(required) = %d, want 110000", got)
	}
}

func TestParseCodeSpecInvalidEntriesIgnored(t *testing.T) {
	for _, tag := range []string{"", "   ", "abc", "min=abc", "=110008", "0", "min=0", "1.5"} {
		if spec := ParseCodeSpec(tag); !spec.IsZero() {
			t.Fatalf("ParseCodeSpec(%q) = %+v, want zero spec", tag, spec)
		}
	}
}

func TestCodeSpecMergeOtherWins(t *testing.T) {
	inline := ParseCodeSpec("110000,min=110008")
	tagged := ParseCodeSpec("110001,min=110088")
	merged := inline.Merge(tagged)
	if merged.Default != 110001 {
		t.Fatalf("merged default = %d, want 110001", merged.Default)
	}
	if got := merged.CodeFor("min"); got != 110088 {
		t.Fatalf("merged CodeFor(min) = %d, want 110088", got)
	}
}

// TestCodeSpecMergeByRuleBeatsDefault 按规则项不会被另一处的默认项盖掉
func TestCodeSpecMergeByRuleBeatsDefault(t *testing.T) {
	merged := ParseCodeSpec("min=110008").Merge(ParseCodeSpec("110000"))
	if got := merged.CodeFor("min"); got != 110008 {
		t.Fatalf("CodeFor(min) = %d, want 110008", got)
	}
	if got := merged.CodeFor("required"); got != 110000 {
		t.Fatalf("CodeFor(required) = %d, want 110000", got)
	}
}

func TestCodeSpecMergeIntoEmpty(t *testing.T) {
	var empty CodeSpec
	merged := empty.Merge(ParseCodeSpec("110026"))
	if merged.Default != 110026 {
		t.Fatalf("merged default = %d, want 110026", merged.Default)
	}
	if kept := ParseCodeSpec("110026").Merge(CodeSpec{}); kept.Default != 110026 {
		t.Fatalf("merge with zero spec changed default to %d", kept.Default)
	}
}

func TestApplyCodeSpecSkipsControlRules(t *testing.T) {
	rules := ParseRules("omitempty,min=3,dive,required")
	ApplyCodeSpec(rules, ParseCodeSpec("110008"))
	for _, r := range rules {
		switch r.Name {
		case "omitempty", "dive":
			if r.Code != 0 {
				t.Fatalf("control rule %s got code %d, want 0", r.Name, r.Code)
			}
		default:
			if r.Code != 110008 {
				t.Fatalf("rule %s got code %d, want 110008", r.Name, r.Code)
			}
		}
	}
}

func TestApplyCodeSpecZeroSpecNoop(t *testing.T) {
	rules := ParseRules("required")
	ApplyCodeSpec(rules, CodeSpec{})
	if rules[0].Code != 0 {
		t.Fatalf("code = %d, want 0", rules[0].Code)
	}
}

// TestParseRulesWithCodesStripsInline code= 必须在解析期被剥离，
// 否则 evalRule 对未知规则名返回 false，带码字段会恒定失败
func TestParseRulesWithCodesStripsInline(t *testing.T) {
	rules, spec := ParseRulesWithCodes("required,code=110026,min=3")
	if len(rules) != 2 {
		t.Fatalf("rules = %d, want 2: %+v", len(rules), rules)
	}
	if rules[0].Name != "required" || rules[1].Name != "min" {
		t.Fatalf("unexpected rules: %s, %s", rules[0].Name, rules[1].Name)
	}
	if spec.Default != 110026 {
		t.Fatalf("inline default = %d, want 110026", spec.Default)
	}
}

// TestParseRulesWithCodesStripsInvalidInline 写错的码也要剥离，不能漏成规则
func TestParseRulesWithCodesStripsInvalidInline(t *testing.T) {
	rules, spec := ParseRulesWithCodes("required,code=abc")
	if len(rules) != 1 || rules[0].Name != "required" {
		t.Fatalf("unexpected rules: %+v", rules)
	}
	if !spec.IsZero() {
		t.Fatalf("spec = %+v, want zero", spec)
	}
}

func TestParseRulesWithCodesInlineByRule(t *testing.T) {
	rules, spec := ParseRulesWithCodes("required,min=8,code=required=110007,code=min=110008")
	if len(rules) != 2 {
		t.Fatalf("rules = %d, want 2: %+v", len(rules), rules)
	}
	if got := spec.CodeFor("required"); got != 110007 {
		t.Fatalf("CodeFor(required) = %d, want 110007", got)
	}
	if got := spec.CodeFor("min"); got != 110008 {
		t.Fatalf("CodeFor(min) = %d, want 110008", got)
	}
}

// TestParseRulesAppliesInlineCode ParseRules 保持原签名，内联码就地落盘
func TestParseRulesAppliesInlineCode(t *testing.T) {
	rules := ParseRules("required,code=110026")
	if len(rules) != 1 {
		t.Fatalf("rules = %d, want 1", len(rules))
	}
	if rules[0].Code != 110026 {
		t.Fatalf("code = %d, want 110026", rules[0].Code)
	}
}

func TestParseRulesNoCodeStaysZero(t *testing.T) {
	for _, r := range ParseRules("required,min=3,email") {
		if r.Code != 0 {
			t.Fatalf("rule %s got code %d, want 0", r.Name, r.Code)
		}
	}
}
