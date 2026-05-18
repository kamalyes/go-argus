/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-17 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-18 12:26:16
 * @FilePath: \go-argus\alloc_bench_test.go
 * @Description: Argus 基本回归基准测试，性能对比测试请移步 go-argus-benchmark
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"
	"time"
)

func BenchmarkVarStringPath(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("hello", "required")
	}
}

func BenchmarkVarInterfacePath(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.Var("hello", "required")
	}
}

func BenchmarkVarStringPathParallel(b *testing.B) {
	v := New()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.VarString("hello", "required")
		}
	})
}

func BenchmarkVarInterfacePathParallel(b *testing.B) {
	v := New()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.Var("hello", "required")
		}
	})
}

func BenchmarkFieldSuccess(b *testing.B) {
	v := New()
	s := "1"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Var(&s, "len=1")
	}
}

func BenchmarkFieldFailure(b *testing.B) {
	v := New()
	s := "12"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Var(&s, "len=1")
	}
}

func BenchmarkFieldArrayDiveSuccess(b *testing.B) {
	v := New()
	m := []string{"val1", "val2", "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Var(m, "required,dive,required")
	}
}

func BenchmarkFieldMapDiveSuccess(b *testing.B) {
	v := New()
	m := map[string]string{"val1": "val1", "val2": "val2", "val3": "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Var(m, "required,dive,required")
	}
}

func BenchmarkFieldOrTagSuccess(b *testing.B) {
	v := New()
	s := "rgba(0,0,0,1)"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Var(s, "rgb|rgba")
	}
}

type benchSubTest struct {
	Test string `validate:"required"`
}

type benchTestString struct {
	Required  string `validate:"required"`
	Len       string `validate:"len=10"`
	Min       string `validate:"min=1"`
	Max       string `validate:"max=10"`
	MinMax    string `validate:"min=1,max=10"`
	Lt        string `validate:"lt=10"`
	Lte       string `validate:"lte=10"`
	Gt        string `validate:"gt=10"`
	Gte       string `validate:"gte=10"`
	OmitEmpty string `validate:"omitempty,min=1,max=10"`
	Sub       *benchSubTest
	SubIgnore *benchSubTest `validate:"-"`
	Anonymous struct {
		A string `validate:"required"`
	}
}

func BenchmarkStructSimpleSuccess(b *testing.B) {
	v := New()
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Struct(validFoo)
	}
}

func BenchmarkStructSimpleFailure(b *testing.B) {
	v := New()
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	invalidFoo := &Foo{StringValue: "Fo", IntValue: 3}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Struct(invalidFoo)
	}
}

func BenchmarkStructComplexSuccess(b *testing.B) {
	v := New()
	tSuccess := &benchTestString{
		Required:  "Required",
		Len:       "length==10",
		Min:       "min=1",
		Max:       "1234567890",
		MinMax:    "12345",
		Lt:        "012345678",
		Lte:       "0123456789",
		Gt:        "01234567890",
		Gte:       "0123456789",
		OmitEmpty: "",
		Sub: &benchSubTest{
			Test: "1",
		},
		SubIgnore: &benchSubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "1",
		},
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Struct(tSuccess)
	}
}

func BenchmarkStructComplexFailure(b *testing.B) {
	v := New()
	tFail := &benchTestString{
		Required:  "",
		Len:       "",
		Min:       "",
		Max:       "12345678901",
		MinMax:    "",
		Lt:        "0123456789",
		Lte:       "01234567890",
		Gt:        "1",
		Gte:       "1",
		OmitEmpty: "12345678901",
		Sub: &benchSubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "",
		},
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Struct(tFail)
	}
}

func BenchmarkStructSimpleCrossFieldSuccess(b *testing.B) {
	v := New()
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Struct(test)
	}
}

func BenchmarkStructSimpleCrossFieldFailure(b *testing.B) {
	v := New()
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * -5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Struct(test)
	}
}

type benchOneof struct {
	Color string `validate:"oneof=red green"`
}

func BenchmarkOneof(b *testing.B) {
	w := &benchOneof{Color: "green"}
	val := New()
	for i := 0; i < b.N; i++ {
		_ = val.Struct(w)
	}
}

func BenchmarkVarString_Email(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("test@example.com", "required,email")
	}
}

func BenchmarkVarString_URL(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("https://example.com/path", "required,url")
	}
}

func BenchmarkVarString_MultiRule(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("hello world", "required,min=1,max=100")
	}
}

func BenchmarkVarString_BIC(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("DEUTDEFF", "bic")
	}
}

func BenchmarkVarString_Cron(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("*/5 * * * *", "cron")
	}
}

func BenchmarkVarString_UUID(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("6ba7b810-9dad-11d1-80b4-00c04fd430c8", "uuid")
	}
}

func BenchmarkVarString_IPv4(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("192.168.1.1", "ipv4")
	}
}

func BenchmarkVarString_Semver(b *testing.B) {
	v := New()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("1.2.3-alpha.1+build.123", "semver")
	}
}
