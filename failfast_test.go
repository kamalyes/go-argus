/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-25 17:35:19
 * @FilePath: \go-argus\failfast_test.go
 * @Description: WithFailFast 语义与性能测试
 *
 * 默认行为（收集全部错误）必须逐字节不变，failFast 仅在显式开启时生效。
 */

package validator

import (
	"testing"
)

type failFastSubject struct {
	A string `validate:"required"`
	B string `validate:"required"`
	C string `validate:"required"`
	D string `validate:"required"`
}

type failFastNested struct {
	Inner failFastSubject
	Tail  string `validate:"required"`
}

type failFastSlice struct {
	Items []string          `validate:"required,dive,required"`
	Pairs map[string]string `validate:"omitempty,dive,required"`
}

// errCount 取 ValidationErrors 长度，非该类型则失败
func errCount(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		return 0
	}
	ve, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("期望 ValidationErrors，实得 %T: %v", err, err)
	}
	return len(ve)
}

// TestFailFast_DefaultCollectsAllErrors 默认行为必须仍收集全部错误
//
// 这是向后兼容的红线：表单逐字段回显等调用方依赖拿到全部错误。
func TestFailFast_DefaultCollectsAllErrors(t *testing.T) {
	v := New()
	err := v.Struct(&failFastSubject{}) // 4 个字段全部为空
	if n := errCount(t, err); n != 4 {
		t.Fatalf("默认行为应收集全部 4 个错误，实得 %d: %v", n, err)
	}
}

// TestFailFast_StopsAtFirstError 开启后只返回首个错误
func TestFailFast_StopsAtFirstError(t *testing.T) {
	v := New(WithFailFast())
	err := v.Struct(&failFastSubject{})
	if n := errCount(t, err); n != 1 {
		t.Fatalf("failFast 应只返回 1 个错误，实得 %d: %v", n, err)
	}
}

// TestFailFast_FirstErrorIdentityMatchesDefault failFast 返回的首个错误
// 必须与默认模式的首个错误一致——只是提前收口，不改变判定顺序。
func TestFailFast_FirstErrorIdentityMatchesDefault(t *testing.T) {
	subject := &failFastSubject{A: "ok", B: ""} // 首错应落在 B

	full := New().Struct(subject)
	fast := New(WithFailFast()).Struct(subject)

	fullErrs, ok := full.(ValidationErrors)
	if !ok || len(fullErrs) == 0 {
		t.Fatalf("默认模式应返回 ValidationErrors，实得 %T: %v", full, full)
	}
	fastErrs, ok := fast.(ValidationErrors)
	if !ok || len(fastErrs) != 1 {
		t.Fatalf("failFast 应返回单个错误，实得 %T: %v", fast, fast)
	}

	if fullErrs[0].Field() != fastErrs[0].Field() {
		t.Fatalf("首错字段不一致: 默认 %q vs failFast %q",
			fullErrs[0].Field(), fastErrs[0].Field())
	}
	if fullErrs[0].Tag() != fastErrs[0].Tag() {
		t.Fatalf("首错规则不一致: 默认 %q vs failFast %q",
			fullErrs[0].Tag(), fastErrs[0].Tag())
	}
}

// TestFailFast_ValidPayloadUnaffected 合法入参在两种模式下都必须通过
func TestFailFast_ValidPayloadUnaffected(t *testing.T) {
	valid := &failFastSubject{A: "a", B: "b", C: "c", D: "d"}
	if err := New().Struct(valid); err != nil {
		t.Fatalf("默认模式误判合法入参: %v", err)
	}
	if err := New(WithFailFast()).Struct(valid); err != nil {
		t.Fatalf("failFast 模式误判合法入参: %v", err)
	}
}

// TestFailFast_NestedStructStops 嵌套结构体内部出错同样提前收口
func TestFailFast_NestedStructStops(t *testing.T) {
	subject := &failFastNested{} // Inner 4 个空 + Tail 空

	if n := errCount(t, New().Struct(subject)); n < 2 {
		t.Fatalf("默认模式应收集嵌套与外层的多个错误，实得 %d", n)
	}
	if n := errCount(t, New(WithFailFast()).Struct(subject)); n != 1 {
		t.Fatalf("failFast 在嵌套场景应只返回 1 个错误，实得 %d", n)
	}
}

// TestFailFast_DiveStops dive 元素出错提前收口（切片与 map 两条路径）
func TestFailFast_DiveStops(t *testing.T) {
	subject := &failFastSlice{
		Items: []string{"", "", ""},
		Pairs: map[string]string{"k1": "", "k2": ""},
	}

	fullN := errCount(t, New().Struct(subject))
	if fullN < 2 {
		t.Fatalf("默认模式应收集多个 dive 错误，实得 %d", fullN)
	}
	if n := errCount(t, New(WithFailFast()).Struct(subject)); n != 1 {
		t.Fatalf("failFast 在 dive 场景应只返回 1 个错误，实得 %d", n)
	}
}

// --- 性能 ---

func BenchmarkFailFast_LargeFailFirstField(b *testing.B) {
	v := New(WithFailFast())
	s := newProbeLarge()
	s.F01 = ""
	if err := v.Struct(s); err == nil {
		b.Fatal("fixture must fail")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkFailFast_LargeFailLastField(b *testing.B) {
	v := New(WithFailFast())
	s := newProbeLarge()
	s.F12 = "not-an-ip"
	if err := v.Struct(s); err == nil {
		b.Fatal("fixture must fail")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

// BenchmarkFailFast_LargeAllValid 全合法路径不得因开启 failFast 而变慢
func BenchmarkFailFast_LargeAllValid(b *testing.B) {
	v := New(WithFailFast())
	s := newProbeLarge()
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}
