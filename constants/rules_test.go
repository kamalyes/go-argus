package constants

import "testing"

func TestCompareOperatorString(t *testing.T) {
	if OpGreaterThan.String() != RuleGT {
		t.Fatalf("expected %q, got %q", RuleGT, OpGreaterThan.String())
	}
}

func TestCmpOpFromStr(t *testing.T) {
	cases := map[string]CmpOp{
		RuleEq:  CmpEQ,
		RuleNe:  CmpNE,
		RuleGT:  CmpGT,
		RuleGTE: CmpGTE,
		RuleLT:  CmpLT,
		RuleLTE: CmpLTE,
	}
	for op, want := range cases {
		if got := CmpOpFromStr(op); got != want {
			t.Fatalf("expected %s to map to %v, got %v", op, want, got)
		}
	}
	if got := CmpOpFromStr("unknown"); got != CmpOp(-1) {
		t.Fatalf("expected unknown op to map to -1, got %v", got)
	}
	if got := CmpOpForOperator(OpGreaterThanOrEqual); got != CmpGTE {
		t.Fatalf("expected operator to map to %v, got %v", CmpGTE, got)
	}
}

func TestRuleGroups(t *testing.T) {
	if !NeedsParamParts(RuleOneOf) || NeedsParamParts(RuleMin) {
		t.Fatal("unexpected param parts rule classification")
	}
	if !IsScalarCompareRule(RuleLTE) || IsScalarCompareRule(RuleEq) {
		t.Fatal("unexpected scalar compare rule classification")
	}
	if !IsFieldCompareRule(RuleEqCSField) || IsFieldCompareRule(RuleRequired) {
		t.Fatal("unexpected field compare rule classification")
	}
	if !IsLocalFieldCompareRule(RuleEqField) || IsLocalFieldCompareRule(RuleEqCSField) {
		t.Fatal("unexpected local field compare rule classification")
	}
	if !IsCrossStructFieldCompareRule(RuleEqCSField) || IsCrossStructFieldCompareRule(RuleEqField) {
		t.Fatal("unexpected cross-struct field compare rule classification")
	}
	if !IsOmitEmptyRule(RuleOmitZero) || IsOmitEmptyRule(RuleRequired) {
		t.Fatal("unexpected omit-empty rule classification")
	}
	if !IsStructControlRule(RuleStructOnly) || IsStructControlRule(RuleDive) {
		t.Fatal("unexpected struct control rule classification")
	}
	if !IsDiveControlRule(RuleEndKeys) || IsDiveControlRule(RuleRequired) {
		t.Fatal("unexpected dive control rule classification")
	}
	if !StopsStructDive(RuleDive) || StopsStructDive(RuleRequired) {
		t.Fatal("unexpected struct dive stop classification")
	}
}

func TestCompareOperatorForRule(t *testing.T) {
	cases := map[string]string{
		RuleLen:         RuleEq,
		RuleEqField:     RuleEq,
		RuleMin:         RuleGTE,
		RuleGTECSField:  RuleGTE,
		RuleMax:         RuleLTE,
		RuleLTECSField:  RuleLTE,
		RuleGT:          RuleGT,
		RuleAfterField:  RuleGT,
		RuleLT:          RuleLT,
		RuleBeforeField: RuleLT,
		RuleNeField:     RuleNe,
		RuleNeCSField:   RuleNe,
	}
	for name, want := range cases {
		if got := CompareOperatorForRule(name); got != want {
			t.Fatalf("expected %s to map to %s, got %s", name, want, got)
		}
	}
	if got := CompareOperatorForRule(RuleRequired); got != RuleEmpty {
		t.Fatalf("expected unknown compare rule to map to empty, got %s", got)
	}
}
