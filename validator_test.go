/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-28 00:00:00
 * @FilePath: \go-argus\validator_test.go
 * @Description: validator.go 测试，覆盖校验器核心逻辑、Struct/Var 校验和 evalRule 分支
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type scheduleRequest struct {
	Type      string    `json:"type" validate:"required,oneof=instant schedule"`
	StartTime time.Time `json:"start_time" validate:"required_if=type schedule,after=now,beforefield=end_time"`
	EndTime   time.Time `json:"end_time" validate:"required_with=start_time,afterfield=start_time"`
}

func TestStructTimeRangeRules(t *testing.T) {
	v := New(WithRequiredStructEnabled())
	now := time.Now()

	req := scheduleRequest{
		Type:      "schedule",
		StartTime: now.Add(time.Hour),
		EndTime:   now.Add(2 * time.Hour),
	}
	if err := v.Struct(req); err != nil {
		t.Fatalf("expected valid request, got %v", err)
	}
}

func TestStructTimeRangeRulesRejectPastStart(t *testing.T) {
	v := New(WithRequiredStructEnabled())
	now := time.Now()

	req := scheduleRequest{
		Type:      "schedule",
		StartTime: now.Add(-time.Hour),
		EndTime:   now.Add(2 * time.Hour),
	}
	err := v.Struct(req)
	if err == nil {
		t.Fatal("expected past start_time to fail")
	}
}

func TestStructTimeRangeRulesRejectEndBeforeStart(t *testing.T) {
	v := New(WithRequiredStructEnabled())
	now := time.Now()

	req := scheduleRequest{
		Type:      "schedule",
		StartTime: now.Add(2 * time.Hour),
		EndTime:   now.Add(time.Hour),
	}
	err := v.Struct(req)
	if err == nil {
		t.Fatal("expected end_time before start_time to fail")
	}
}

func TestRequiredIfUsesJSONFieldName(t *testing.T) {
	v := New(WithRequiredStructEnabled())

	req := scheduleRequest{Type: "schedule"}
	err := v.Struct(req)
	if err == nil {
		t.Fatal("expected required_if=type schedule to require start_time")
	}
}

func TestSetTagName(t *testing.T) {
	v := New()
	v.SetTagName("custom")
	if v.tagName != "custom" {
		t.Fatal("expected tag name to be custom")
	}
}

func TestSetTagNameEmpty(t *testing.T) {
	v := New()
	original := v.tagName
	v.SetTagName("  ")
	if v.tagName != original {
		t.Fatal("expected tag name to remain unchanged for blank input")
	}
}

func TestRegisterTagNameFunc(t *testing.T) {
	v := New()
	v.RegisterTagNameFunc(func(sf reflect.StructField) string {
		return "custom_" + sf.Name
	})
	if v.tagNameFunc == nil {
		t.Fatal("expected tagNameFunc to be set")
	}
}

func TestStructInvalidInput(t *testing.T) {
	v := New()
	err := v.Struct(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
	err = v.Struct("not-a-struct")
	if err == nil {
		t.Fatal("expected error for non-struct input")
	}
	err = v.Struct(42)
	if err == nil {
		t.Fatal("expected error for int input")
	}
}

func TestStructCtxInvalidInput(t *testing.T) {
	v := New()
	err := v.StructCtx(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestVar(t *testing.T) {
	v := New()
	if err := v.Var("test@example.com", "required,email"); err != nil {
		t.Fatalf("expected valid var: %v", err)
	}
}

func TestVarInvalid(t *testing.T) {
	v := New()
	if err := v.Var("bad", "email"); err == nil {
		t.Fatal("expected invalid var to fail")
	}
}

func TestVarCtx(t *testing.T) {
	v := New()
	if err := v.VarCtx(context.Background(), "test@example.com", "required,email"); err != nil {
		t.Fatalf("expected valid var: %v", err)
	}
}

func TestOmitEmpty(t *testing.T) {
	v := New()
	if err := v.Var("", "omitempty,email"); err != nil {
		t.Fatalf("expected omitempty to skip empty: %v", err)
	}
}

func TestOmitZero(t *testing.T) {
	v := New()
	if err := v.Var(0, "omitzero,min=1"); err != nil {
		t.Fatalf("expected omitzero to skip zero: %v", err)
	}
}

func TestOmitNil(t *testing.T) {
	v := New()
	var ptr *int
	if err := v.Var(ptr, "omitnil,min=1"); err != nil {
		t.Fatalf("expected omitnil to skip nil: %v", err)
	}
}

func TestDiveSlice(t *testing.T) {
	v := New()
	if err := v.Var([]string{"a", "b"}, "dive,required"); err != nil {
		t.Fatalf("expected dive to pass: %v", err)
	}
	if err := v.Var([]string{"a", ""}, "dive,required"); err == nil {
		t.Fatal("expected dive to fail for empty element")
	}
}

func TestDiveMap(t *testing.T) {
	v := New()
	m := map[string]string{"key1": "a", "key2": "b"}
	if err := v.Var(m, "dive,required"); err != nil {
		t.Fatalf("expected dive map to pass: %v", err)
	}
}

func TestDiveInvalidField(t *testing.T) {
	v := New()
	if err := v.Var(nil, "dive,required"); err != nil {
		t.Fatalf("expected dive to skip nil: %v", err)
	}
}

func TestRequiredWithAll(t *testing.T) {
	type req struct {
		A string `validate:"required_with_all=b c"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{A: "", B: "val", C: "val"}); err == nil {
		t.Fatal("expected required_with_all to fail when all fields present and A empty")
	}
	if err := v.Struct(req{A: "x", B: "val", C: "val"}); err != nil {
		t.Fatalf("expected required_with_all to pass: %v", err)
	}
	if err := v.Struct(req{A: "", B: "val", C: ""}); err != nil {
		t.Fatalf("expected required_with_all to skip when not all present: %v", err)
	}
}

func TestRequiredWithAllEmptyParam(t *testing.T) {
	type req struct {
		A string `validate:"required_with_all="`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{}); err != nil {
		t.Fatalf("expected required_with_all with empty param to skip: %v", err)
	}
}

func TestRequiredWithoutAll(t *testing.T) {
	type req struct {
		A string `validate:"required_without_all=b c"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{A: "", B: "", C: ""}); err == nil {
		t.Fatal("expected required_without_all to fail when all absent and A empty")
	}
	if err := v.Struct(req{A: "x", B: "", C: ""}); err != nil {
		t.Fatalf("expected required_without_all to pass: %v", err)
	}
	if err := v.Struct(req{A: "", B: "val", C: ""}); err != nil {
		t.Fatalf("expected required_without_all to skip when some present: %v", err)
	}
}

func TestRequiredWithoutAllEmptyParam(t *testing.T) {
	type req struct {
		A string `validate:"required_without_all="`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{}); err != nil {
		t.Fatalf("expected required_without_all with empty param to skip: %v", err)
	}
}

func TestExcludedIf(t *testing.T) {
	type req struct {
		Mode  string `validate:"required,oneof=sync async"`
		Debug string `validate:"excluded_if=mode sync"`
	}
	v := New()
	if err := v.Struct(req{Mode: "sync", Debug: ""}); err != nil {
		t.Fatalf("expected excluded_if to pass when field empty: %v", err)
	}
}

func TestExcludedUnless(t *testing.T) {
	type req struct {
		Mode  string `validate:"required,oneof=sync async"`
		Debug string `validate:"excluded_unless=mode debug"`
	}
	v := New()
	if err := v.Struct(req{Mode: "sync", Debug: ""}); err != nil {
		t.Fatalf("expected excluded_unless to pass when field empty: %v", err)
	}
}

func TestExcludedWith(t *testing.T) {
	type req struct {
		A string `validate:"excluded_with=b"`
		B string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: "val"}); err != nil {
		t.Fatalf("expected excluded_with to pass when field empty: %v", err)
	}
}

func TestExcludedWithAll(t *testing.T) {
	type req struct {
		A string `validate:"excluded_with_all=b c"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: "val", C: "val"}); err != nil {
		t.Fatalf("expected excluded_with_all to pass when field empty: %v", err)
	}
}

func TestExcludedWithout(t *testing.T) {
	type req struct {
		A string `validate:"excluded_without=b"`
		B string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: ""}); err != nil {
		t.Fatalf("expected excluded_without to pass when field empty: %v", err)
	}
}

func TestExcludedWithoutAll(t *testing.T) {
	type req struct {
		A string `validate:"excluded_without_all=b c"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: "", C: ""}); err != nil {
		t.Fatalf("expected excluded_without_all to pass when field empty: %v", err)
	}
}

func TestFieldContains(t *testing.T) {
	type req struct {
		Name    string `validate:"fieldcontains=Pattern"`
		Pattern string `validate:""`
	}
	v := New()
	if err := v.Struct(req{Name: "hello world", Pattern: "hello"}); err != nil {
		t.Fatalf("expected fieldcontains to pass: %v", err)
	}
	if err := v.Struct(req{Name: "goodbye", Pattern: "hello"}); err == nil {
		t.Fatal("expected fieldcontains to fail")
	}
}

func TestFieldExcludes(t *testing.T) {
	type req struct {
		Name    string `validate:"fieldexcludes=Pattern"`
		Pattern string `validate:""`
	}
	v := New()
	if err := v.Struct(req{Name: "goodbye", Pattern: "hello"}); err != nil {
		t.Fatalf("expected fieldexcludes to pass: %v", err)
	}
	if err := v.Struct(req{Name: "hello world", Pattern: "hello"}); err == nil {
		t.Fatal("expected fieldexcludes to fail")
	}
}

func TestRangeRule(t *testing.T) {
	type req struct {
		Min int `validate:""`
		Max int `validate:"range=Min|Max"`
	}
	v := New()
	if err := v.Struct(req{Min: 1, Max: 10}); err != nil {
		t.Fatalf("expected range to pass: %v", err)
	}
	if err := v.Struct(req{Min: 10, Max: 1}); err == nil {
		t.Fatal("expected range to fail for min > max")
	}
}

func TestRangeInvalidParts(t *testing.T) {
	type req struct {
		Start int `validate:"range=Onlyone"`
	}
	v := New()
	if err := v.Struct(req{Start: 1}); err == nil {
		t.Fatal("expected range to fail for single field")
	}
}

func TestRangeInvalidField(t *testing.T) {
	type req struct {
		Start int `validate:"range=Missing|End"`
		End   int `validate:""`
	}
	v := New()
	if err := v.Struct(req{Start: 1, End: 10}); err == nil {
		t.Fatal("expected range to fail for missing field")
	}
}

func TestEqField(t *testing.T) {
	type req struct {
		Password string `validate:"required"`
		Confirm  string `validate:"required,eqfield=password"`
	}
	v := New()
	if err := v.Struct(req{Password: "abc", Confirm: "abc"}); err != nil {
		t.Fatalf("expected eqfield to pass: %v", err)
	}
	if err := v.Struct(req{Password: "abc", Confirm: "xyz"}); err == nil {
		t.Fatal("expected eqfield to fail")
	}
}

func TestNeField(t *testing.T) {
	type req struct {
		Old string `validate:""`
		New string `validate:"nefield=old"`
	}
	v := New()
	if err := v.Struct(req{Old: "abc", New: "xyz"}); err != nil {
		t.Fatalf("expected nefield to pass: %v", err)
	}
}

func TestGtField(t *testing.T) {
	type req struct {
		Min int `validate:""`
		Max int `validate:"gtfield=min"`
	}
	v := New()
	if err := v.Struct(req{Min: 5, Max: 10}); err != nil {
		t.Fatalf("expected gtfield to pass: %v", err)
	}
}

func TestGteField(t *testing.T) {
	type req struct {
		Min int `validate:""`
		Max int `validate:"gtefield=min"`
	}
	v := New()
	if err := v.Struct(req{Min: 5, Max: 5}); err != nil {
		t.Fatalf("expected gtefield to pass: %v", err)
	}
}

func TestLtField(t *testing.T) {
	type req struct {
		Min int `validate:"ltfield=max"`
		Max int `validate:""`
	}
	v := New()
	if err := v.Struct(req{Min: 5, Max: 10}); err != nil {
		t.Fatalf("expected ltfield to pass: %v", err)
	}
}

func TestLteField(t *testing.T) {
	type req struct {
		Min int `validate:"ltefield=max"`
		Max int `validate:""`
	}
	v := New()
	if err := v.Struct(req{Min: 5, Max: 5}); err != nil {
		t.Fatalf("expected ltefield to pass: %v", err)
	}
}

func TestEqCsField(t *testing.T) {
	type inner struct {
		Value string `validate:"eqcsfield=Name"`
	}
	type outer struct {
		Name  string `validate:""`
		Inner inner  `validate:""`
	}
	v := New()
	if err := v.Struct(outer{Name: "test", Inner: inner{Value: "test"}}); err != nil {
		t.Fatalf("expected eqcsfield to pass: %v", err)
	}
}

func TestNeCsField(t *testing.T) {
	type inner struct {
		Value string `validate:"necsfield=Name"`
	}
	type outer struct {
		Name  string `validate:""`
		Inner inner  `validate:""`
	}
	v := New()
	if err := v.Struct(outer{Name: "test", Inner: inner{Value: "other"}}); err != nil {
		t.Fatalf("expected necsfield to pass: %v", err)
	}
}

func TestGtCsField(t *testing.T) {
	type inner struct {
		Value int `validate:"gtcsfield=Min"`
	}
	type outer struct {
		Min   int   `validate:""`
		Inner inner `validate:""`
	}
	v := New()
	if err := v.Struct(outer{Min: 5, Inner: inner{Value: 10}}); err != nil {
		t.Fatalf("expected gtcsfield to pass: %v", err)
	}
}

func TestGteCsField(t *testing.T) {
	type inner struct {
		Value int `validate:"gtecsfield=Min"`
	}
	type outer struct {
		Min   int   `validate:""`
		Inner inner `validate:""`
	}
	v := New()
	if err := v.Struct(outer{Min: 5, Inner: inner{Value: 5}}); err != nil {
		t.Fatalf("expected gtecsfield to pass: %v", err)
	}
}

func TestLtCsField(t *testing.T) {
	type inner struct {
		Value int `validate:"ltcsfield=Max"`
	}
	type outer struct {
		Max   int   `validate:""`
		Inner inner `validate:""`
	}
	v := New()
	if err := v.Struct(outer{Max: 10, Inner: inner{Value: 5}}); err != nil {
		t.Fatalf("expected ltcsfield to pass: %v", err)
	}
}

func TestLteCsField(t *testing.T) {
	type inner struct {
		Value int `validate:"ltecsfield=Max"`
	}
	type outer struct {
		Max   int   `validate:""`
		Inner inner `validate:""`
	}
	v := New()
	if err := v.Struct(outer{Max: 10, Inner: inner{Value: 10}}); err != nil {
		t.Fatalf("expected ltecsfield to pass: %v", err)
	}
}

func TestAfterBefore(t *testing.T) {
	type req struct {
		Start time.Time `validate:"after=now"`
		End   time.Time `validate:"before=now"`
	}
	v := New()
	now := time.Now()
	if err := v.Struct(req{Start: now.Add(time.Hour), End: now.Add(-time.Hour)}); err != nil {
		t.Fatalf("expected after/before to pass: %v", err)
	}
}

func TestJoinNS(t *testing.T) {
	if joinNS("", "child") != "child" {
		t.Fatal("expected child for empty parent")
	}
	if joinNS("parent", "") != "parent" {
		t.Fatal("expected parent for empty child")
	}
	if joinNS("parent", "child") != "parent.child" {
		t.Fatal("expected parent.child")
	}
}

func TestStructNoValidateTag(t *testing.T) {
	type noTag struct {
		Name string
	}
	v := New()
	if err := v.Struct(noTag{Name: "test"}); err != nil {
		t.Fatalf("expected struct without tags to pass: %v", err)
	}
}

func TestStructDiveIntoNestedStruct(t *testing.T) {
	type address struct {
		City string `validate:"required"`
	}
	type person struct {
		Name    string  `validate:"required"`
		Address address `validate:""`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(person{Name: "John", Address: address{City: "NYC"}}); err != nil {
		t.Fatalf("expected nested struct to pass: %v", err)
	}
	if err := v.Struct(person{Name: "John", Address: address{City: ""}}); err == nil {
		t.Fatal("expected nested struct to fail for empty city")
	}
}

func TestStructOnlyNoDive(t *testing.T) {
	type inner struct {
		Value string `validate:"required"`
	}
	type outer struct {
		Name  string `validate:"required"`
		Inner inner  `validate:"structonly"`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(outer{Name: "test", Inner: inner{Value: ""}}); err != nil {
		t.Fatalf("expected structonly to skip inner validation: %v", err)
	}
}

func TestNoStructLevel(t *testing.T) {
	type inner struct {
		Value string `validate:"required"`
	}
	type outer struct {
		Name  string `validate:"required"`
		Inner inner  `validate:"nostructlevel"`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(outer{Name: "test", Inner: inner{Value: ""}}); err != nil {
		t.Fatalf("expected nostructlevel to skip inner validation: %v", err)
	}
}

func TestUnknownRule(t *testing.T) {
	v := New()
	if err := v.Var("test", "unknown_rule"); err == nil {
		t.Fatal("expected unknown rule to fail")
	}
}

func TestEmptyRule(t *testing.T) {
	v := New()
	if err := v.Var("test", ""); err != nil {
		t.Fatalf("expected empty rule to pass: %v", err)
	}
}

func TestCustomValidationCtx(t *testing.T) {
	v := New()
	err := v.RegisterValidationCtx("ctx_check", func(ctx context.Context, fl FieldLevel) bool {
		return fl.Param() != ""
	})
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}
	if err := v.Var("test", "ctx_check=hello"); err != nil {
		t.Fatalf("expected ctx_check to pass: %v", err)
	}
}

func TestFieldContainsNonString(t *testing.T) {
	type req struct {
		Count   int    `validate:"fieldcontains=X"`
		Pattern string `validate:""`
	}
	v := New()
	if err := v.Struct(req{Count: 1, Pattern: "x"}); err == nil {
		t.Fatal("expected fieldcontains to fail for non-string field")
	}
}

func TestFieldContainsMissingField(t *testing.T) {
	type req struct {
		Name string `validate:"fieldcontains=Missing"`
	}
	v := New()
	if err := v.Struct(req{Name: "hello"}); err == nil {
		t.Fatal("expected fieldcontains to fail for missing field")
	}
}

func TestFieldContainsNonStringOther(t *testing.T) {
	type req struct {
		Name  string `validate:"fieldcontains=Count"`
		Count int    `validate:""`
	}
	v := New()
	if err := v.Struct(req{Name: "1", Count: 1}); err != nil {
		t.Fatalf("expected fieldcontains to pass for int via scalarString: %v", err)
	}
}

func TestRequiredWithout(t *testing.T) {
	type req struct {
		A string `validate:"required_without=b"`
		B string `validate:""`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{A: "", B: ""}); err == nil {
		t.Fatal("expected required_without to fail when both empty")
	}
	if err := v.Struct(req{A: "x", B: ""}); err != nil {
		t.Fatalf("expected required_without to pass: %v", err)
	}
}

func TestRequiredWithoutNoneEmpty(t *testing.T) {
	type req struct {
		A string `validate:"required_without=b"`
		B string `validate:""`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{A: "", B: "val"}); err != nil {
		t.Fatalf("expected required_without to skip when B present: %v", err)
	}
}

func TestExcludedIfConditionMet(t *testing.T) {
	type req struct {
		Mode  string `validate:"required"`
		Debug string `validate:"excluded_if=mode sync"`
	}
	v := New()
	err := v.Struct(req{Mode: "sync", Debug: "true"})
	if err == nil {
		t.Fatal("expected excluded_if to fail when condition met and field not empty")
	}
}

func TestExcludedUnlessConditionMet(t *testing.T) {
	type req struct {
		Mode  string `validate:"required"`
		Debug string `validate:"excluded_unless=mode debug"`
	}
	v := New()
	err := v.Struct(req{Mode: "sync", Debug: "true"})
	if err == nil {
		t.Fatal("expected excluded_unless to fail when condition not met and field not empty")
	}
}

func TestExcludedWithFieldPresentValue(t *testing.T) {
	type req struct {
		A string `validate:"excluded_with=b"`
		B string `validate:""`
	}
	v := New()
	err := v.Struct(req{A: "val", B: "val"})
	if err == nil {
		t.Fatal("expected excluded_with to fail when B present and A not empty")
	}
}

func TestExcludedWithAllAllPresent(t *testing.T) {
	type req struct {
		A string `validate:"excluded_with_all=b c"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New()
	err := v.Struct(req{A: "val", B: "val", C: "val"})
	if err == nil {
		t.Fatal("expected excluded_with_all to fail when all present and A not empty")
	}
}

func TestExcludedWithoutFieldAbsent(t *testing.T) {
	type req struct {
		A string `validate:"excluded_without=b"`
		B string `validate:""`
	}
	v := New()
	err := v.Struct(req{A: "val", B: ""})
	if err == nil {
		t.Fatal("expected excluded_without to fail when B absent and A not empty")
	}
}

func TestExcludedWithoutAllAllAbsent(t *testing.T) {
	type req struct {
		A string `validate:"excluded_without_all=b c"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New()
	err := v.Struct(req{A: "val", B: "", C: ""})
	if err == nil {
		t.Fatal("expected excluded_without_all to fail when all absent and A not empty")
	}
}

func TestApplyRulesEmpty(t *testing.T) {
	v := New()
	errs := make(ValidationErrors, 0)
	v.applyRules(context.Background(), reflect.Value{}, reflect.Value{}, reflect.ValueOf("test"), "ns", "sns", "f", "sf", nil, &errs)
	if len(errs) != 0 {
		t.Fatal("expected no errors for empty rules")
	}
}

func TestStructPtrInput(t *testing.T) {
	type req struct {
		Name string `validate:"required"`
	}
	v := New(WithRequiredStructEnabled())
	r := &req{Name: "test"}
	if err := v.Struct(r); err != nil {
		t.Fatalf("expected pointer to struct to pass: %v", err)
	}
}

func TestStructNilPtrInput(t *testing.T) {
	type req struct {
		Name string `validate:"required"`
	}
	v := New()
	var r *req
	err := v.Struct(r)
	if err == nil {
		t.Fatal("expected nil pointer to struct to fail")
	}
}

func TestBeforeField(t *testing.T) {
	type req struct {
		Start time.Time `validate:"beforefield=end"`
		End   time.Time `validate:""`
	}
	v := New()
	now := time.Now()
	if err := v.Struct(req{Start: now, End: now.Add(time.Hour)}); err != nil {
		t.Fatalf("expected beforefield to pass: %v", err)
	}
}

func TestAfterField(t *testing.T) {
	type req struct {
		End   time.Time `validate:"afterfield=start"`
		Start time.Time `validate:""`
	}
	v := New()
	now := time.Now()
	if err := v.Struct(req{Start: now, End: now.Add(time.Hour)}); err != nil {
		t.Fatalf("expected afterfield to pass: %v", err)
	}
}

func TestRequiredUnless(t *testing.T) {
	type req struct {
		Name string `validate:"required_unless=mode guest"`
		Mode string `validate:""`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{Name: "", Mode: "guest"}); err != nil {
		t.Fatalf("expected required_unless to skip when condition met: %v", err)
	}
	if err := v.Struct(req{Name: "", Mode: "admin"}); err == nil {
		t.Fatal("expected required_unless to fail when condition not met")
	}
}

func TestRequiredWith(t *testing.T) {
	type req struct {
		A string `validate:"required_with=b"`
		B string `validate:""`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{A: "", B: "val"}); err == nil {
		t.Fatal("expected required_with to fail when B present and A empty")
	}
	if err := v.Struct(req{A: "x", B: "val"}); err != nil {
		t.Fatalf("expected required_with to pass: %v", err)
	}
	if err := v.Struct(req{A: "", B: ""}); err != nil {
		t.Fatalf("expected required_with to skip when B absent: %v", err)
	}
}

func TestDiveSliceWithIndex(t *testing.T) {
	type req struct {
		Tags []string `validate:"dive,required"`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{Tags: []string{"a", "b"}}); err != nil {
		t.Fatalf("expected dive to pass: %v", err)
	}
	err := v.Struct(req{Tags: []string{"a", ""}})
	if err == nil {
		t.Fatal("expected dive to fail for empty element")
	}
	ve := err.(ValidationErrors)
	if !containsString(ve[0].Namespace(), "[1]") {
		t.Fatalf("expected index in namespace, got %s", ve[0].Namespace())
	}
}

func TestDiveMapWithKey(t *testing.T) {
	type req struct {
		Meta map[string]string `validate:"dive,required"`
	}
	v := New(WithRequiredStructEnabled())
	err := v.Struct(req{Meta: map[string]string{"key": ""}})
	if err == nil {
		t.Fatal("expected dive map to fail for empty value")
	}
	ve := err.(ValidationErrors)
	if !containsString(ve[0].Namespace(), "[key]") {
		t.Fatalf("expected key in namespace, got %s", ve[0].Namespace())
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestRangeCommaSeparator(t *testing.T) {
	type req struct {
		Min int `validate:""`
		Max int `validate:"range=Min|Max"`
	}
	v := New()
	if err := v.Struct(req{Min: 1, Max: 10}); err != nil {
		t.Fatalf("expected range with pipe separator to pass: %v", err)
	}
}

func TestValidateStructNonStruct(t *testing.T) {
	v := New()
	if err := v.Struct("hello"); err == nil {
		t.Fatal("expected Struct to fail for non-struct")
	}
}

func TestApplyRulesOmitEmptyEmpty(t *testing.T) {
	type req struct {
		Name string `validate:"omitempty,min=5"`
	}
	v := New()
	if err := v.Struct(req{Name: ""}); err != nil {
		t.Fatalf("expected omitempty to skip empty: %v", err)
	}
}

func TestApplyRulesOmitZeroZero(t *testing.T) {
	type req struct {
		Count int `validate:"omitzero,min=5"`
	}
	v := New()
	if err := v.Struct(req{Count: 0}); err != nil {
		t.Fatalf("expected omitzero to skip zero: %v", err)
	}
}

func TestApplyRulesOmitNilNil(t *testing.T) {
	type req struct {
		Ptr *string `validate:"omitnil,min=1"`
	}
	v := New()
	if err := v.Struct(req{Ptr: nil}); err != nil {
		t.Fatalf("expected omitnil to skip nil: %v", err)
	}
}

func TestEvalRuleCustomValidation(t *testing.T) {
	v := New()
	v.RegisterValidationCtx("custom", func(ctx context.Context, fl FieldLevel) bool {
		return fl.Field().String() == "valid"
	})
	if err := v.Var("invalid", "custom"); err == nil {
		t.Fatal("expected custom validation to fail")
	}
	if err := v.Var("valid", "custom"); err != nil {
		t.Fatalf("expected custom validation to pass: %v", err)
	}
}

func TestRenderTranslationDefaultFallback(t *testing.T) {
	SetLocale("en")
	type req struct {
		Name string `validate:"required"`
	}
	v := New(WithRequiredStructEnabled())
	err := v.Struct(req{})
	if err == nil {
		t.Fatal("expected required to fail")
	}
	msg := TranslateValidationErrors(err, "en")
	if len(msg) == 0 || msg[0].Message == "" {
		t.Fatal("expected non-empty translated message")
	}
}

func TestFieldContainsNonStringOtherScalar(t *testing.T) {
	type req struct {
		Name  string `validate:"fieldcontains=Count"`
		Count int    `validate:""`
	}
	v := New()
	if err := v.Struct(req{Name: "5", Count: 5}); err != nil {
		t.Fatalf("expected fieldcontains to pass with int field: %v", err)
	}
}

func TestRuleRangeFieldNotFound(t *testing.T) {
	type req struct {
		Start int `validate:"range=NonExistent|End"`
		End   int `validate:""`
	}
	v := New()
	if err := v.Struct(req{Start: 1, End: 10}); err == nil {
		t.Fatal("expected range to fail for non-existent field")
	}
}

func TestRuleRangeEndFieldNotFound(t *testing.T) {
	type req struct {
		Start int `validate:"range=End|NonExistent"`
		End   int `validate:""`
	}
	v := New()
	if err := v.Struct(req{Start: 1, End: 10}); err == nil {
		t.Fatal("expected range to fail for non-existent end field")
	}
}

func TestValidateStructNestedNoError(t *testing.T) {
	type inner struct {
		Name string `validate:"required"`
	}
	type outer struct {
		Inner inner `validate:"structonly"`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(outer{Inner: inner{Name: "ok"}}); err != nil {
		t.Fatalf("expected struct to pass: %v", err)
	}
}

func TestDiveInvalidKind(t *testing.T) {
	type req struct {
		Name string `validate:"dive,required"`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(req{Name: "hello"}); err != nil {
		t.Fatalf("expected dive on non-slice/map to silently pass: %v", err)
	}
}

func TestRequiredWithFieldEmpty(t *testing.T) {
	type req struct {
		A string `validate:"required_with=B"`
		B string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: "present"}); err == nil {
		t.Fatal("expected required_with to fail when A empty and B present")
	}
}

func TestExcludedWithFieldPresent2(t *testing.T) {
	type req struct {
		A string `validate:"excluded_with=B"`
		B string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "value", B: "present"}); err == nil {
		t.Fatal("expected excluded_with to fail when A has value and B present")
	}
}

func TestExcludedWithoutFieldAbsent2(t *testing.T) {
	type req struct {
		A string `validate:"excluded_without=B"`
		B string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "value", B: ""}); err == nil {
		t.Fatal("expected excluded_without to fail when A has value and B absent")
	}
}

func TestExcludedWithoutAllNonePresent2(t *testing.T) {
	type req struct {
		A string `validate:"excluded_without_all=B C"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "value", B: "", C: ""}); err == nil {
		t.Fatal("expected excluded_without_all to fail when A has value and B,C absent")
	}
}

func TestRequiredUnlessConditionTrue(t *testing.T) {
	type req struct {
		A string `validate:"required_unless=Mode off"`
		B string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: ""}); err == nil {
		t.Fatal("expected required_unless to fail when condition not met")
	}
}

func TestExcludedUnlessConditionTrue(t *testing.T) {
	type req struct {
		A string `validate:"excluded_unless=Mode on"`
		B string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "value", B: ""}); err == nil {
		t.Fatal("expected excluded_unless to fail when condition not met and A has value")
	}
}

func TestExcludedWithAllAllPresent2(t *testing.T) {
	type req struct {
		A string `validate:"excluded_with_all=B C"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "value", B: "b", C: "c"}); err == nil {
		t.Fatal("expected excluded_with_all to fail when A has value and B,C present")
	}
}

func TestValidateStructNestedNonStruct(t *testing.T) {
	type inner struct {
		Name string `validate:"required"`
	}
	type outer struct {
		Inner *inner `validate:""`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(outer{Inner: &inner{Name: "ok"}}); err != nil {
		t.Fatalf("expected nested struct ptr to pass: %v", err)
	}
}

func TestRenderTranslationUnknownTag(t *testing.T) {
	v := New()
	v.RegisterValidation("unknowntag", func(fl FieldLevel) bool {
		return false
	})
	err := v.Var("test", "unknowntag")
	if err == nil {
		t.Fatal("expected unknowntag to fail")
	}
	ve := err.(ValidationErrors)
	msgs := ve.Translate("xx")
	if len(msgs) == 0 {
		t.Fatal("expected translation message")
	}
	if msgs[0].Message == "" {
		t.Fatal("expected non-empty message from fe.Error() fallback")
	}
}

func TestRenderTranslationNilError(t *testing.T) {
	msg := TranslateValidationErrors(nil, "en")
	if msg != nil {
		t.Fatal("expected nil for nil error")
	}
}

func TestRenderTranslationNonValidationError(t *testing.T) {
	msg := TranslateValidationErrors(fmt.Errorf("some error"), "en")
	if len(msg) == 0 || msg[0].Message != "some error" {
		t.Fatal("expected error message as fallback")
	}
}

func TestApplyRulesOmitNil(t *testing.T) {
	type req struct {
		Ptr *string `validate:"omitnil,required"`
	}
	v := New(WithRequiredStructEnabled())
	s := "hello"
	if err := v.Struct(req{Ptr: &s}); err != nil {
		t.Fatalf("expected omitnil to continue when field not nil: %v", err)
	}
	if err := v.Struct(req{Ptr: nil}); err != nil {
		t.Fatalf("expected omitnil to skip nil pointer: %v", err)
	}
}

func TestApplyRulesStructOnly(t *testing.T) {
	type inner struct {
		Name string `validate:"required"`
	}
	type outer struct {
		Inner inner `validate:"structonly"`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(outer{Inner: inner{Name: ""}}); err != nil {
		t.Fatalf("expected structonly to skip nested validation: %v", err)
	}
}

func TestApplyRulesNoStructLevel(t *testing.T) {
	type inner struct {
		Name string `validate:"required"`
	}
	type outer struct {
		Inner inner `validate:"nostructlevel"`
	}
	v := New(WithRequiredStructEnabled())
	if err := v.Struct(outer{Inner: inner{Name: ""}}); err != nil {
		t.Fatalf("expected nostructlevel to skip nested validation: %v", err)
	}
}

func TestEvalRuleExcludedWith(t *testing.T) {
	type req struct {
		A string `validate:"excluded_with=B"`
		B string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: "present"}); err != nil {
		t.Fatalf("expected excluded_with to pass when A empty and B present: %v", err)
	}
}

func TestEvalRuleExcludedWithAll(t *testing.T) {
	type req struct {
		A string `validate:"excluded_with_all=B C"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: "b", C: "c"}); err != nil {
		t.Fatalf("expected excluded_with_all to pass when A empty: %v", err)
	}
}

func TestEvalRuleExcludedWithout(t *testing.T) {
	type req struct {
		A string `validate:"excluded_without=B"`
		B string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: "present"}); err != nil {
		t.Fatalf("expected excluded_without to pass when A empty and B present: %v", err)
	}
}

func TestEvalRuleExcludedWithoutAll(t *testing.T) {
	type req struct {
		A string `validate:"excluded_without_all=B C"`
		B string `validate:""`
		C string `validate:""`
	}
	v := New()
	if err := v.Struct(req{A: "", B: "b", C: "c"}); err != nil {
		t.Fatalf("expected excluded_without_all to pass when A empty and B,C present: %v", err)
	}
}

func TestEvalRuleEqCsField(t *testing.T) {
	type inner struct {
		Value int `validate:"eqcsfield=Value"`
	}
	type outer struct {
		Value int   `validate:""`
		Inner inner `validate:""`
	}
	v := New()
	if err := v.Struct(outer{Value: 5, Inner: inner{Value: 5}}); err != nil {
		t.Fatalf("expected eqcsfield to pass: %v", err)
	}
}

func TestFieldContainsNonScalarOther(t *testing.T) {
	type req struct {
		Name  string `validate:"fieldcontains=Items"`
		Items []int  `validate:""`
	}
	v := New()
	if err := v.Struct(req{Name: "hello", Items: []int{1, 2}}); err == nil {
		t.Fatal("expected fieldcontains to fail when other is non-scalar")
	}
}

func TestValidateStructNonStructValue(t *testing.T) {
	v := New()
	if err := v.Struct(123); err == nil {
		t.Fatal("expected Struct to fail for int")
	}
}

func TestExcludedUnlessConditionNotMet(t *testing.T) {
	type req struct {
		Mode  string `validate:"required,oneof=on off"`
		Debug string `validate:"excluded_unless=Mode on"`
	}
	v := New()
	if err := v.Struct(req{Mode: "off", Debug: ""}); err != nil {
		t.Fatalf("expected excluded_unless to pass when field empty and condition not met: %v", err)
	}
	if err := v.Struct(req{Mode: "off", Debug: "value"}); err == nil {
		t.Fatal("expected excluded_unless to fail when field non-empty and condition not met")
	}
}

func TestExcludedWithConditionMet(t *testing.T) {
	type req struct {
		B     string `validate:""`
		Debug string `validate:"excluded_with=B"`
	}
	v := New()
	if err := v.Struct(req{B: "present", Debug: ""}); err != nil {
		t.Fatalf("expected excluded_with to pass when field empty: %v", err)
	}
	if err := v.Struct(req{B: "present", Debug: "value"}); err == nil {
		t.Fatal("expected excluded_with to fail when field non-empty and B present")
	}
}

func TestExcludedWithAllConditionMet(t *testing.T) {
	type req struct {
		B     string `validate:""`
		C     string `validate:""`
		Debug string `validate:"excluded_with_all=B C"`
	}
	v := New()
	if err := v.Struct(req{B: "b", C: "c", Debug: ""}); err != nil {
		t.Fatalf("expected excluded_with_all to pass when field empty: %v", err)
	}
	if err := v.Struct(req{B: "b", C: "c", Debug: "value"}); err == nil {
		t.Fatal("expected excluded_with_all to fail when field non-empty and B,C present")
	}
}

func TestExcludedWithoutConditionMet(t *testing.T) {
	type req struct {
		B     string `validate:""`
		Debug string `validate:"excluded_without=B"`
	}
	v := New()
	if err := v.Struct(req{B: "", Debug: ""}); err != nil {
		t.Fatalf("expected excluded_without to pass when field empty: %v", err)
	}
	if err := v.Struct(req{B: "", Debug: "value"}); err == nil {
		t.Fatal("expected excluded_without to fail when field non-empty and B absent")
	}
}
