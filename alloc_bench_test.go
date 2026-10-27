/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-17 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-27 15:02:26
 * @FilePath: \go-argus\alloc_bench_test.go
 * @Description: Argus 基本回归基准测试，性能对比测试请移步 go-argus-benchmark
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"strconv"
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

// ══════════════════════════════════════════════════════════════════════════════
// 性能基线基准（原 zz_perfprobe_bench_test.go，合并至此）
//
// 覆盖：禁用级别、字段遍历规模（padded3/10/20/40）、proto 宽窄结构体、
// 嵌套与 dive、跨字段比较、Var 单槽缓存、自定义规则并发。
// ══════════════════════════════════════════════════════════════════════════════

// ---------- 1. small struct, 3 fields, all valid ----------

type probeSmall struct {
	Name  string `validate:"required,min=2,max=32"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}

func BenchmarkProbeSmallAllValid(b *testing.B) {
	v := New()
	s := &probeSmall{Name: "kamal", Email: "k@example.com", Age: 30}
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

// ---------- 2. larger struct, 15 fields, mixed rules ----------

type probeLarge struct {
	F01 string  `validate:"required,min=2,max=32"`
	F02 string  `validate:"required,email"`
	F03 int     `validate:"gte=0,lte=130"`
	F04 string  `validate:"omitempty,url"`
	F05 string  `validate:"required,alphanum"`
	F06 string  `validate:"required,oneof=red green blue"`
	F07 int     `validate:"min=1"`
	F08 string  `validate:"required,uuid4"`
	F09 string  `validate:"omitempty,numeric"`
	F10 string  `validate:"required,len=8"`
	F11 float64 `validate:"gt=0"`
	F12 string  `validate:"required,ipv4"`
	F13 string  `validate:"omitempty,alpha"`
	F14 bool    `validate:"-"`
	F15 string  `validate:"required,max=64"`
}

func newProbeLarge() *probeLarge {
	return &probeLarge{
		F01: "kamal", F02: "k@example.com", F03: 30, F04: "https://example.com",
		F05: "abc123", F06: "green", F07: 5, F08: "6ba7b810-9dad-41d1-80b4-00c04fd430c8",
		F09: "12345", F10: "abcdefgh", F11: 1.5, F12: "192.168.1.1", F13: "abc",
		F14: true, F15: "hello",
	}
}

func BenchmarkProbeLargeMixedAllValid(b *testing.B) {
	v := New()
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

// ---------- 3./4. fail on FIRST field vs fail on LAST field ----------

func BenchmarkProbeLargeFailFirstField(b *testing.B) {
	v := New()
	s := newProbeLarge()
	s.F01 = "" // first field fails -> early exit
	if err := v.Struct(s); err == nil {
		b.Fatal("fixture must fail")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkProbeLargeFailLastField(b *testing.B) {
	v := New()
	s := newProbeLarge()
	s.F15 = "" // last field fails -> full walk
	if err := v.Struct(s); err == nil {
		b.Fatal("fixture must fail")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

// ---------- 5. nested structs ----------

type probeStreet struct {
	Street string `validate:"required,min=3"`
	Zip    string `validate:"required,len=5"`
}

type probeAddress struct {
	Line probeStreet
	City string `validate:"required,min=2"`
}

type probeUserNested struct {
	Name    string `validate:"required,min=2"`
	Address probeAddress
}

func BenchmarkProbeNestedAllValid(b *testing.B) {
	v := New()
	s := &probeUserNested{
		Name: "kamal",
		Address: probeAddress{
			Line: probeStreet{Street: "Main St", Zip: "12345"},
			City: "Beijing",
		},
	}
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkProbeNestedFailDeepest(b *testing.B) {
	v := New()
	s := &probeUserNested{
		Name: "kamal",
		Address: probeAddress{
			Line: probeStreet{Street: "Main St", Zip: "123"}, // deepest field fails
			City: "Beijing",
		},
	}
	if err := v.Struct(s); err == nil {
		b.Fatal("fixture must fail")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

// ---------- 6. slice field with per-element rules ----------

type probeSliceScalar struct {
	Tags []string `validate:"required,dive,required,min=2,max=16"`
}

type probeSliceStruct struct {
	Items []probeStreet `validate:"required,dive"`
}

func BenchmarkProbeSliceScalarDive10(b *testing.B) {
	v := New()
	s := &probeSliceScalar{Tags: make([]string, 0, 10)}
	for i := 0; i < 10; i++ {
		s.Tags = append(s.Tags, "tag"+strconv.Itoa(i))
	}
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkProbeSliceStructDive10(b *testing.B) {
	v := New()
	s := &probeSliceStruct{Items: make([]probeStreet, 0, 10)}
	for i := 0; i < 10; i++ {
		s.Items = append(s.Items, probeStreet{Street: "Main St", Zip: "12345"})
	}
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkProbeSliceScalarDive10FailLast(b *testing.B) {
	v := New()
	s := &probeSliceScalar{Tags: make([]string, 0, 10)}
	for i := 0; i < 9; i++ {
		s.Tags = append(s.Tags, "tag"+strconv.Itoa(i))
	}
	s.Tags = append(s.Tags, "x") // last element violates min=2
	if err := v.Struct(s); err == nil {
		b.Fatal("fixture must fail")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

// ---------- 7. concurrent ----------

func BenchmarkProbeSmallAllValidParallel(b *testing.B) {
	v := New()
	s := &probeSmall{Name: "kamal", Email: "k@example.com", Age: 30}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.Struct(s)
		}
	})
}

func BenchmarkProbeLargeMixedAllValidParallel(b *testing.B) {
	v := New()
	s := newProbeLarge()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.Struct(s)
		}
	})
}

func BenchmarkProbeLargeFailLastFieldParallel(b *testing.B) {
	v := New()
	s := newProbeLarge()
	s.F15 = ""
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.Struct(s)
		}
	})
}

// Custom-rule path: exercises the RWMutex-guarded v.validations lookup.
func BenchmarkProbeCustomRuleParallel(b *testing.B) {
	v := New()
	_ = v.RegisterValidation("probenoop", func(fl FieldLevel) bool { return true })
	type T struct {
		A string `validate:"probenoop"`
	}
	s := &T{A: "x"}
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.Struct(s)
		}
	})
}

// ---------- 8. proto/@inject_tag shape: many fields, few rules ----------

// probeProtoWide mimics a generated proto message: 40 fields, only 3 carry rules.
type probeProtoWide struct {
	A01 string `validate:"required"`
	A02 string
	A03 string
	A04 string
	A05 string
	A06 string
	A07 string
	A08 string
	A09 string
	A10 string
	A11 int32
	A12 int32
	A13 int32
	A14 int32
	A15 int32
	A16 int64
	A17 int64
	A18 int64
	A19 int64
	A20 int64 `validate:"gte=0"`
	A21 bool
	A22 bool
	A23 bool
	A24 bool
	A25 bool
	A26 float64
	A27 float64
	A28 float64
	A29 float64
	A30 float64
	A31 string
	A32 string
	A33 string
	A34 string
	A35 string
	A36 string
	A37 string
	A38 string
	A39 string
	A40 string `validate:"omitempty,max=64"`
}

// probeProtoNarrow has the same 3 rule-bearing fields and nothing else.
type probeProtoNarrow struct {
	A01 string `validate:"required"`
	A20 int64  `validate:"gte=0"`
	A40 string `validate:"omitempty,max=64"`
}

func BenchmarkProbeProtoWide40Fields3Rules(b *testing.B) {
	v := New()
	s := &probeProtoWide{A01: "x", A20: 1}
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkProbeProtoNarrow3Fields3Rules(b *testing.B) {
	v := New()
	s := &probeProtoNarrow{A01: "x", A20: 1}
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

// ---------- 9. error-consumption cost: does the caller pay for messages it discards? ----------

func BenchmarkProbeFailDiscardError(b *testing.B) {
	v := New()
	s := newProbeLarge()
	s.F15 = ""
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkProbeFailReadErrorMessage(b *testing.B) {
	v := New()
	s := newProbeLarge()
	s.F15 = ""
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := v.Struct(s); err != nil {
			_ = err.Error()
		}
	}
}

// Graduated field counts: identical rule work (3 rule-bearing fields),
// increasing numbers of rule-free padding fields. Isolates the per-ruleless-field cost.

type padded3 struct {
	A string `validate:"required"`
	B int64  `validate:"gte=0"`
	C string `validate:"omitempty,max=64"`
}

type padded10 struct {
	A                          string `validate:"required"`
	B                          int64  `validate:"gte=0"`
	C                          string `validate:"omitempty,max=64"`
	P1, P2, P3, P4, P5, P6, P7 string
}

type padded20 struct {
	A                                       string `validate:"required"`
	B                                       int64  `validate:"gte=0"`
	C                                       string `validate:"omitempty,max=64"`
	P1, P2, P3, P4, P5, P6, P7, P8, P9, P10 string
	P11, P12, P13, P14, P15, P16, P17       string
}

type padded40 struct {
	A                                       string `validate:"required"`
	B                                       int64  `validate:"gte=0"`
	C                                       string `validate:"omitempty,max=64"`
	P1, P2, P3, P4, P5, P6, P7, P8, P9, P10 string
	P11, P12, P13, P14, P15, P16, P17, P18  string
	P19, P20, P21, P22, P23, P24, P25, P26  string
	P27, P28, P29, P30, P31, P32, P33, P34  string
	P35, P36, P37                           string
}

func BenchmarkPad03Fields(b *testing.B) {
	v := New()
	s := &padded3{A: "x", B: 1}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkPad10Fields(b *testing.B) {
	v := New()
	s := &padded10{A: "x", B: 1}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkPad20Fields(b *testing.B) {
	v := New()
	s := &padded20{A: "x", B: 1}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkPad40Fields(b *testing.B) {
	v := New()
	s := &padded40{A: "x", B: 1}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

// Serial counterpart to BenchmarkProbeCustomRuleParallel, to compute lock scaling.
func BenchmarkCustomRuleSerial(b *testing.B) {
	v := New()
	_ = v.RegisterValidation("probenoop2", func(fl FieldLevel) bool { return true })
	type T struct {
		A string `validate:"probenoop2"`
	}
	s := &T{A: "x"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

// Builtin single-field struct, serial, as the contention-free control.
func BenchmarkBuiltinOneFieldSerial(b *testing.B) {
	v := New()
	type T struct {
		A string `validate:"required"`
	}
	s := &T{A: "x"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkBuiltinOneFieldParallel(b *testing.B) {
	v := New()
	type T struct {
		A string `validate:"required"`
	}
	s := &T{A: "x"}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.Struct(s)
		}
	})
}

// Cross-field (gtfield) on a PASSING validation: checks Interface() boxing allocs.
func BenchmarkCrossFieldIntSuccess(b *testing.B) {
	v := New()
	type T struct {
		Start int
		End   int `validate:"gtfield=Start"`
	}
	s := &T{Start: 1, End: 5}
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

func BenchmarkCrossFieldStringSuccess(b *testing.B) {
	v := New()
	type T struct {
		A string
		B string `validate:"eqfield=A"`
	}
	s := &T{A: "x", B: "x"}
	if err := v.Struct(s); err != nil {
		b.Fatalf("fixture must pass: %v", err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Struct(s)
	}
}

// The varSlot single-slot bypass cache (validator.go:362-376) stores a freshly
// allocated *varRuleSlot on every miss. These probe the hit vs rotating-miss case.

func BenchmarkVarSlotSingleTag(b *testing.B) {
	v := New()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("hello", "min=1,max=100")
	}
}

func BenchmarkVarSlotTwoTagsRotating(b *testing.B) {
	v := New()
	tags := [2]string{"min=1,max=100", "min=2,max=50"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("hello", tags[i&1])
	}
}

func BenchmarkVarSlotFourTagsRotating(b *testing.B) {
	v := New()
	tags := [4]string{"min=1,max=100", "min=2,max=50", "min=3,max=40", "min=4,max=30"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.VarString("hello", tags[i&3])
	}
}

func BenchmarkVarSlotSingleTagParallel(b *testing.B) {
	v := New()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.VarString("hello", "min=1,max=100")
		}
	})
}

func BenchmarkVarSlotFourTagsRotatingParallel(b *testing.B) {
	v := New()
	tags := [4]string{"min=1,max=100", "min=2,max=50", "min=3,max=40", "min=4,max=30"}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = v.VarString("hello", tags[i&3])
			i++
		}
	})
}

// i18n.Msg takes the global i18n RWMutex twice per message (GetLocale + Lookup)
// and runs strings.ReplaceAll per placeholder. These probe that under concurrency.

func BenchmarkErrMsgSerial(b *testing.B) {
	v := New()
	s := newProbeLarge()
	s.F15 = ""
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := v.Struct(s); err != nil {
			_ = err.Error()
		}
	}
}

func BenchmarkErrMsgParallel(b *testing.B) {
	v := New()
	s := newProbeLarge()
	s.F15 = ""
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := v.Struct(s); err != nil {
				_ = err.Error()
			}
		}
	})
}
