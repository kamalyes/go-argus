/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-28 15:52:08
 * @FilePath: \go-argus\errors_test.go
 * @Description: errors.go 测试，覆盖 InvalidValidationError、ValidationErrors 和 FieldError
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestInvalidValidationErrorNilType(t *testing.T) {
	err := &InvalidValidationError{}
	if err.Error() != "validator: (nil)" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestInvalidValidationErrorWithType(t *testing.T) {
	err := &InvalidValidationError{Type: reflect.TypeOf(0)}
	if err.Error() != "validator: (nil int)" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestValidationErrorsEmpty(t *testing.T) {
	var ve ValidationErrors
	if ve.Error() != "" {
		t.Fatalf("expected empty error string, got %q", ve.Error())
	}
}

func TestValidationErrorsMultiple(t *testing.T) {
	ve := ValidationErrors{
		&fieldError{tag: "required", ns: "Name", field: "Name"},
		&fieldError{tag: "email", ns: "Email", field: "Email"},
	}
	// 多条拼接路径与逐条 Sprintf 换行连接逐字节一致（守护部件直写的不变量）
	want := fmt.Sprintf(fieldErrMsg, "Name", "Name", "required") + "\n" +
		fmt.Sprintf(fieldErrMsg, "Email", "Email", "email")
	if got := ve.Error(); got != want {
		t.Fatalf("join mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestFieldErrorMethods(t *testing.T) {
	fe := &fieldError{
		tag:         "min",
		ns:          "User.Name",
		structNs:    "User.Name",
		field:       "name",
		structField: "Name",
		fieldVal:    reflect.ValueOf("ab"),
		param:       "3",
	}
	if fe.Tag() != "min" {
		t.Fatalf("expected tag min, got %s", fe.Tag())
	}
	if fe.ActualTag() != "min" {
		t.Fatalf("expected actualTag min, got %s", fe.ActualTag())
	}
	if fe.Namespace() != "User.Name" {
		t.Fatalf("expected namespace User.Name, got %s", fe.Namespace())
	}
	if fe.StructNamespace() != "User.Name" {
		t.Fatalf("expected structNs User.Name, got %s", fe.StructNamespace())
	}
	if fe.Field() != "name" {
		t.Fatalf("expected field name, got %s", fe.Field())
	}
	if fe.StructField() != "Name" {
		t.Fatalf("expected structField Name, got %s", fe.StructField())
	}
	if fe.Value() != "ab" {
		t.Fatalf("expected value ab, got %v", fe.Value())
	}
	if fe.Param() != "3" {
		t.Fatalf("expected param 3, got %s", fe.Param())
	}
	if fe.Kind() != reflect.String {
		t.Fatalf("expected kind String, got %v", fe.Kind())
	}
	if fe.Type() != reflect.TypeOf("") {
		t.Fatalf("expected type string, got %v", fe.Type())
	}
	expectedErr := "Key: 'User.Name' Error:Field validation for 'name' failed on the 'min' tag"
	if fe.Error() != expectedErr {
		t.Fatalf("unexpected error: %s", fe.Error())
	}
}

func TestFieldErrorNilValue(t *testing.T) {
	// 零值 fieldVal 即 invalid，Value/Kind/Type 均返回零值
	fe := &fieldError{
		tag:         "required",
		ns:          "Test.Field",
		field:       "field",
		structField: "Field",
	}
	if fe.Value() != nil {
		t.Fatal("expected nil value")
	}
	if fe.Kind() != reflect.Invalid {
		t.Fatal("expected Invalid kind")
	}
	if fe.Type() != nil {
		t.Fatal("expected nil type")
	}
}

// ══════════════════════════════════════════════════════════════════════════════
// 错误消息渲染（原 zz_errmsg_probe_test.go，合并至此）
//
// bytes.Buffer.String() 会整体复制，strings.Builder.String() 零拷贝；
// 三处纯 %s 的 fmt.Sprintf 改手写拼接。TestFieldErrorMsg_MatchesSprintf
// 守输出逐字节不变（含引号、CJK、300 字符命名空间）。
// ══════════════════════════════════════════════════════════════════════════════

// TestFieldErrorMsg_MatchesSprintf 手工拼接必须与 fmt.Sprintf 逐字节一致
//
// fieldError.Error 改为按 fieldErrPrefix/Mid1/Mid2/Suffix 手工拼接以省去
// 格式串解析与 []any 装箱。该分段必须始终等价于 fieldErrMsg 的 %s 代换，
// 否则错误文案会与 go-playground 不兼容。本用例是那条不变量的闸门。
func TestFieldErrorMsg_MatchesSprintf(t *testing.T) {
	cases := [][3]string{
		{"User.Name", "Name", "required"},
		{"", "", ""},
		{"A.B.C[0].D", "D", "gte"},
		{"带中文.字段", "字段", "oneof"},
		{"q'uote", "fi'eld", "ta'g"},
		{strings.Repeat("x", 300), "f", "t"},
	}
	for _, c := range cases {
		want := fmt.Sprintf(fieldErrMsg, c[0], c[1], c[2])
		got := (&fieldError{ns: c[0], field: c[1], tag: c[2]}).Error()
		if got != want {
			t.Errorf("拼接与 Sprintf 不一致\n  ns=%q field=%q tag=%q\n  want=%q\n  got =%q",
				c[0], c[1], c[2], want, got)
		}
	}
}

// TestValidationErrorsMsg_Multi 多条错误的换行拼接语义不变
func TestValidationErrorsMsg_Multi(t *testing.T) {
	mk := func(ns, f, tag string) FieldError {
		return &fieldError{ns: ns, field: f, tag: tag}
	}
	// 空集合
	if got := (ValidationErrors{}).Error(); got != "" {
		t.Errorf("空集合应返回空串，实得 %q", got)
	}
	// 单条：须与该条自身一致
	one := ValidationErrors{mk("A", "a", "required")}
	if got, want := one.Error(), strings.TrimSpace(one[0].Error()); got != want {
		t.Errorf("单条不一致\n  want=%q\n  got =%q", got, want)
	}
	// 多条：换行分隔且顺序保持
	three := ValidationErrors{
		mk("A", "a", "required"),
		mk("B", "b", "gte"),
		mk("C", "c", "oneof"),
	}
	got := three.Error()
	parts := strings.Split(got, "\n")
	if len(parts) != 3 {
		t.Fatalf("期望 3 行，实得 %d 行: %q", len(parts), got)
	}
	for i, fe := range three {
		if parts[i] != fe.Error() {
			t.Errorf("第 %d 行不一致\n  want=%q\n  got =%q", i, fe.Error(), parts[i])
		}
	}
}

func BenchmarkProbeFieldErrorMsg(b *testing.B) {
	fe := &fieldError{ns: "ProbeLarge.F15", field: "F15", tag: "required"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fe.Error()
	}
}

func BenchmarkProbeValidationErrorsMsg(b *testing.B) {
	for _, n := range []int{1, 3, 10} {
		ve := make(ValidationErrors, n)
		for i := range ve {
			ve[i] = &fieldError{ns: "ProbeLarge.F15", field: "F15", tag: "required"}
		}
		b.Run(fmt.Sprintf("errors=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ve.Error()
			}
		})
	}
}
