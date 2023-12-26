/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\enum_test.go
 * @Description: enum.go 测试，覆盖泛型枚举校验器全部方法
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validate

import (
	"testing"
)

func TestEnumValidatorIsValid(t *testing.T) {
	roles := NewEnumValidator("admin", "user")
	if !roles.IsValid("admin") {
		t.Fatal("expected admin to be valid")
	}
	if roles.IsValid("guest") {
		t.Fatal("expected guest to be invalid")
	}
}

func TestEnumValidatorNilIsValid(t *testing.T) {
	var roles *EnumValidator[string]
	if roles.IsValid("admin") {
		t.Fatal("expected nil validator to return false")
	}
}

func TestEnumValidatorMustBeValid(t *testing.T) {
	roles := NewEnumValidator("admin")
	if err := roles.MustBeValid("admin"); err != nil {
		t.Fatal("expected admin to be valid")
	}
	if err := roles.MustBeValid("guest"); err == nil {
		t.Fatal("expected guest to be invalid")
	}
}

func TestEnumValidatorGetValidValues(t *testing.T) {
	roles := NewEnumValidator("admin", "user")
	values := roles.GetValidValues()
	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}
	var nilRoles *EnumValidator[string]
	if nilRoles.GetValidValues() != nil {
		t.Fatal("expected nil validator to return nil")
	}
}

func TestEnumValidatorGetValidValuesString(t *testing.T) {
	roles := NewEnumValidator(1, 2, 3)
	strs := roles.GetValidValuesString()
	if len(strs) != 3 {
		t.Fatalf("expected 3 strings, got %d", len(strs))
	}
	var nilRoles *EnumValidator[int]
	if nilRoles.GetValidValuesString() != nil {
		t.Fatal("expected nil validator to return nil")
	}
}

func TestEnumValidatorCount(t *testing.T) {
	roles := NewEnumValidator("a", "b")
	if roles.Count() != 2 {
		t.Fatalf("expected 2, got %d", roles.Count())
	}
	var nilRoles *EnumValidator[string]
	if nilRoles.Count() != 0 {
		t.Fatal("expected nil validator to return 0")
	}
}

func TestEnumValidatorContains(t *testing.T) {
	roles := NewEnumValidator("admin")
	if !roles.Contains("admin") {
		t.Fatal("expected contains admin")
	}
}

func TestEnumValidatorAdd(t *testing.T) {
	roles := NewEnumValidator("admin")
	roles.Add("user", "admin")
	if roles.Count() != 2 {
		t.Fatalf("expected 2 after add, got %d", roles.Count())
	}
}

func TestEnumValidatorAddNilMap(t *testing.T) {
	roles := &EnumValidator[string]{}
	roles.Add("admin")
	if !roles.IsValid("admin") {
		t.Fatal("expected admin to be valid after add")
	}
}

func TestEnumValidatorRemove(t *testing.T) {
	roles := NewEnumValidator("admin", "user")
	roles.Remove("admin")
	if roles.IsValid("admin") {
		t.Fatal("expected admin to be removed")
	}
	if !roles.IsValid("user") {
		t.Fatal("expected user to still be valid")
	}
}

func TestEnumValidatorRemoveNil(t *testing.T) {
	var roles *EnumValidator[string]
	roles.Remove("admin")
}

func TestEnumValidatorRemoveNonExistent(t *testing.T) {
	roles := NewEnumValidator("admin")
	roles.Remove("guest")
	if !roles.IsValid("admin") {
		t.Fatal("expected admin to still be valid")
	}
}

func TestEnumValidatorClear(t *testing.T) {
	roles := NewEnumValidator("admin", "user")
	roles.Clear()
	if roles.Count() != 0 {
		t.Fatalf("expected 0 after clear, got %d", roles.Count())
	}
}

func TestEnumValidatorClearNil(t *testing.T) {
	var roles *EnumValidator[string]
	roles.Clear()
}

func TestEnumValidatorClone(t *testing.T) {
	roles := NewEnumValidator("admin", "user")
	cloned := roles.Clone()
	if !cloned.IsValid("admin") {
		t.Fatal("expected cloned to have admin")
	}
	cloned.Remove("admin")
	if !roles.IsValid("admin") {
		t.Fatal("expected original to still have admin")
	}
}

func TestEnumValidatorCloneNil(t *testing.T) {
	var roles *EnumValidator[string]
	if roles.Clone() != nil {
		t.Fatal("expected nil clone")
	}
}
