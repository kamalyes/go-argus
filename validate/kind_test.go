/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-25 23:15:16
 * @FilePath: \go-argus\validate\kind_test.go
 * @Description: 整数校验函数测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"math"
	"reflect"
	"testing"
)

func TestIsWholeNumber(t *testing.T) {
	// 有效整数（含大数，验证不溢出）
	wholeNumbers := []float64{
		0, 1, -1, 100, -100, 42, -42,
		1e10, -1e10,
		1e15, -1e15, // 超过 int64 范围的大数
		1e20, -1e20, // 远超 int64 范围
		math.MaxFloat64,      // float64 最大值（整数）
		-math.MaxFloat64,     // float64 最小值（整数）
		9223372036854775807,  // int64 max（边界）
		-9223372036854775808, // int64 min（边界）
	}
	for _, f := range wholeNumbers {
		if !IsWholeNumber(f) {
			t.Errorf("IsWholeNumber(%v) = false; want true", f)
		}
	}

	// 非整数
	nonWholeNumbers := []float64{0.5, -0.5, 3.14, -3.14, 0.1, -0.1, 1.5, -1.5}
	for _, f := range nonWholeNumbers {
		if IsWholeNumber(f) {
			t.Errorf("IsWholeNumber(%v) = true; want false", f)
		}
	}

	// 零值
	if !IsWholeNumber(0) {
		t.Error("IsWholeNumber(0) = false; want true")
	}
	if !IsWholeNumber(math.Copysign(0, -1)) {
		t.Error("IsWholeNumber(-0.0) = false; want true")
	}

	// 特殊值：NaN 和 Inf 必须返回 false
	if IsWholeNumber(math.NaN()) {
		t.Error("IsWholeNumber(NaN) = true; want false")
	}
	if IsWholeNumber(math.Inf(1)) {
		t.Error("IsWholeNumber(+Inf) = true; want false")
	}
	if IsWholeNumber(math.Inf(-1)) {
		t.Error("IsWholeNumber(-Inf) = true; want false")
	}

	// 大数带小数部分（float64 精度范围内的非整数）
	if IsWholeNumber(1e15 + 0.5) {
		t.Error("IsWholeNumber(1e15+0.5) = true; want false")
	}
	if IsWholeNumber(1e10 + 0.25) {
		t.Error("IsWholeNumber(1e10+0.25) = true; want false")
	}

	// 注意：1e20 + 0.5 在 float64 中等于 1e20（精度丢失），
	// 所以 IsWholeNumber(1e20+0.5) 返回 true 是正确的（float64 层面它就是整数）
	// 这不是 bug，而是 IEEE 754 double 的固有限制
}

func TestGetReflectKind(t *testing.T) {
	// nil 返回 Invalid
	if got := GetReflectKind(nil); got != reflect.Invalid {
		t.Errorf("GetReflectKind(nil) = %v; want Invalid", got)
	}

	cases := []struct {
		v    interface{}
		want reflect.Kind
	}{
		{42, reflect.Int},
		{"hello", reflect.String},
		{true, reflect.Bool},
		{3.14, reflect.Float64},
		{[]int{1, 2, 3}, reflect.Slice},
		{map[string]int{"a": 1}, reflect.Map},
	}
	for _, c := range cases {
		if got := GetReflectKind(c.v); got != c.want {
			t.Errorf("GetReflectKind(%v) = %v; want %v", c.v, got, c.want)
		}
	}

	// struct
	type TestStruct struct{ X int }
	if got := GetReflectKind(TestStruct{X: 1}); got != reflect.Struct {
		t.Errorf("GetReflectKind(struct) = %v; want Struct", got)
	}

	// ptr
	x := 42
	if got := GetReflectKind(&x); got != reflect.Ptr {
		t.Errorf("GetReflectKind(ptr) = %v; want Ptr", got)
	}
}

func TestStringInteger(t *testing.T) {
	// 有效整数
	validCases := []string{
		"0", "1", "-1", "123", "-123", "+123",
		"999999999999", "-999999999999",
		"9223372036854775807",  // int64 max
		"-9223372036854775808", // int64 min
	}
	for _, s := range validCases {
		if !StringInteger(s) {
			t.Errorf("StringInteger(%q) = false; want true", s)
		}
	}

	// 无效整数
	invalidCases := []string{
		"",      // 空字符串
		"+",     // 只有正号
		"-",     // 只有负号
		"12.3",  // 含小数点
		"12a",   // 含非数字字符
		" 123",  // 含前导空格
		"123 ",  // 含后置空格
		"1 2 3", // 含中间空格
		"0x1A",  // 十六进制
		"1e5",   // 科学计数法
		"++123", // 双正号
		"--123", // 双负号
		"+-123", // 正负号混用
		"12.0",  // 浮点数
		".5",    // 省略前导零
		"5.",    // 省略末尾零
	}
	for _, s := range invalidCases {
		if StringInteger(s) {
			t.Errorf("StringInteger(%q) = true; want false", s)
		}
	}

	// 边界情况
	if !StringInteger("0") {
		t.Error("StringInteger(\"0\") = false; want true")
	}
	if StringInteger("a") {
		t.Error("StringInteger(\"a\") = true; want false")
	}
	if StringInteger(".") {
		t.Error("StringInteger(\".\") = true; want false")
	}
}

// BenchmarkStringInteger 字符串整数校验性能基准
// 零分配手动字节扫描
func BenchmarkStringInteger(b *testing.B) {
	cases := []string{"12345", "-12345", "+12345", "12.345", "12abc"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, s := range cases {
			_ = StringInteger(s)
		}
	}
}

// BenchmarkIsWholeNumber 整数判断性能基准
func BenchmarkIsWholeNumber(b *testing.B) {
	cases := []float64{0, 1, -1, 3.14, -3.14, 100, 100.5}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, f := range cases {
			_ = IsWholeNumber(f)
		}
	}
}
