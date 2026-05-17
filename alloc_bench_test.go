/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-17 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 11:16:50
 * @FilePath: \go-argus\alloc_bench_test.go
 * @Description: Argus 内存分配测试基准
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"reflect"
	"testing"
)

func BenchmarkReflectValueOfInterface(b *testing.B) {
	s := "hello"
	for i := 0; i < b.N; i++ {
		_ = reflect.ValueOf(s)
	}
}

func BenchmarkReflectValueOfString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = reflect.ValueOf("hello")
	}
}

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

func BenchmarkVarString_Email(b *testing.B) {
	v := New()
	email := "test@example.com"
	for i := 0; i < b.N; i++ {
		_ = v.VarString(email, "required,email")
	}
}

func BenchmarkVar_Email(b *testing.B) {
	v := New()
	email := "test@example.com"
	for i := 0; i < b.N; i++ {
		_ = v.Var(email, "required,email")
	}
}

func BenchmarkVarString_URL(b *testing.B) {
	v := New()
	u := "https://example.com/path"
	for i := 0; i < b.N; i++ {
		_ = v.VarString(u, "required,url")
	}
}

func BenchmarkVar_URL(b *testing.B) {
	v := New()
	u := "https://example.com/path"
	for i := 0; i < b.N; i++ {
		_ = v.Var(u, "required,url")
	}
}

func BenchmarkVarString_Semver(b *testing.B) {
	v := New()
	sv := "1.2.3-alpha.1+build.123"
	for i := 0; i < b.N; i++ {
		_ = v.VarString(sv, "semver")
	}
}

func BenchmarkVar_Semver(b *testing.B) {
	v := New()
	sv := "1.2.3-alpha.1+build.123"
	for i := 0; i < b.N; i++ {
		_ = v.Var(sv, "semver")
	}
}

func BenchmarkVarString_IPv4(b *testing.B) {
	v := New()
	ip := "192.168.1.1"
	for i := 0; i < b.N; i++ {
		_ = v.VarString(ip, "ipv4")
	}
}

func BenchmarkVar_IPv4(b *testing.B) {
	v := New()
	ip := "192.168.1.1"
	for i := 0; i < b.N; i++ {
		_ = v.Var(ip, "ipv4")
	}
}

func BenchmarkVarString_MultiRule(b *testing.B) {
	v := New()
	s := "hello world"
	for i := 0; i < b.N; i++ {
		_ = v.VarString(s, "required,min=1,max=100")
	}
}

func BenchmarkVar_MultiRule(b *testing.B) {
	v := New()
	s := "hello world"
	for i := 0; i < b.N; i++ {
		_ = v.Var(s, "required,min=1,max=100")
	}
}
