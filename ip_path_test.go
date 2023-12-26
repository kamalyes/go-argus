/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-27 00:00:00
 * @FilePath: \go-argus\ip_path_test.go
 * @Description: ip_path.go 测试，覆盖 IP 校验、路径匹配和 CIDR 规则
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"net"
	"testing"
)

func TestCompileIPSet(t *testing.T) {
	set, err := CompileIPSet([]string{"192.168.1.*", "10.0.0.0/8"})
	if err != nil {
		t.Fatalf("expected valid IP set: %v", err)
	}
	if set == nil {
		t.Fatal("expected non-nil IP set")
	}
}

func TestMustCompileIPSet(t *testing.T) {
	set := MustCompileIPSet([]string{"192.168.1.*"})
	if set == nil {
		t.Fatal("expected non-nil IP set")
	}
}

func TestMustCompileIPSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid IP set")
		}
	}()
	MustCompileIPSet([]string{"not-an-ip"})
}

func TestMatchPathInList(t *testing.T) {
	if !MatchPathInList("/api/v1/users", []string{"/api/v1"}) {
		t.Fatal("expected path prefix match")
	}
	if MatchPathInList("/other/path", []string{"/api/v1"}) {
		t.Fatal("expected no path prefix match")
	}
	if !MatchPathInList("/api/v1", []string{"/api/v1"}) {
		t.Fatal("expected exact path match")
	}
}

func TestIsIPAllowed(t *testing.T) {
	if !IsIPAllowed("192.168.1.12", []string{"192.168.1.*"}) {
		t.Fatal("expected wildcard IP to be allowed")
	}
	if IsIPAllowed("10.0.0.1", []string{"192.168.1.*"}) {
		t.Fatal("expected non-matching IP to not be allowed")
	}
}

func TestIsIPBlocked(t *testing.T) {
	if !IsIPBlocked("10.0.0.1", []string{"10.0.0.0/8"}) {
		t.Fatal("expected CIDR-blocked IP")
	}
	if IsIPBlocked("192.168.1.1", []string{"10.0.0.0/8"}) {
		t.Fatal("expected non-blocked IP")
	}
}

func TestMatchIPPattern(t *testing.T) {
	if !MatchIPPattern("192.168.1.1", "192.168.1.*") {
		t.Fatal("expected wildcard match")
	}
	if !MatchIPPattern("10.0.0.1", "10.0.0.0/8") {
		t.Fatal("expected CIDR match")
	}
	if MatchIPPattern("192.168.1.1", "10.0.0.0/8") {
		t.Fatal("expected no CIDR match")
	}
}

func TestMatchIPInList(t *testing.T) {
	if !MatchIPInList("192.168.1.1", []string{"192.168.1.*"}) {
		t.Fatal("expected IP in list")
	}
	if MatchIPInList("10.0.0.1", []string{"192.168.1.*"}) {
		t.Fatal("expected IP not in list")
	}
}

func TestIsIPInRange(t *testing.T) {
	start := net.ParseIP("192.168.1.1")
	end := net.ParseIP("192.168.1.254")
	if !IsIPInRange(net.ParseIP("192.168.1.100"), start, end) {
		t.Fatal("expected IP in range")
	}
	if IsIPInRange(net.ParseIP("192.168.2.1"), start, end) {
		t.Fatal("expected IP out of range")
	}
}

func TestMatchIPWithWildcard(t *testing.T) {
	if !MatchIPWithWildcard("192.168.1.1", "192.168.1.*") {
		t.Fatal("expected wildcard match")
	}
	if MatchIPWithWildcard("192.168.2.1", "192.168.1.*") {
		t.Fatal("expected no wildcard match")
	}
}

func TestIsPrivateIP(t *testing.T) {
	if !IsPrivateIP("192.168.1.1") {
		t.Fatal("expected private IP")
	}
	if !IsPrivateIP("10.0.0.1") {
		t.Fatal("expected private IP")
	}
	if !IsPrivateIP("172.16.0.1") {
		t.Fatal("expected private IP")
	}
	if IsPrivateIP("8.8.8.8") {
		t.Fatal("expected public IP")
	}
}
