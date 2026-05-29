/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\network_test.go
 * @Description: network.go 测试，覆盖 IP 白名单、黑名单、CIDR、通配符、IP 范围和路径匹配
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validate

import (
	"net"
	"testing"
)

func netParseIP(s string) net.IP {
	return net.ParseIP(s)
}

func TestIPBaseValidateIP(t *testing.T) {
	b := &IPBase{}
	if err := b.ValidateIP("192.168.1.1"); err != nil {
		t.Fatal("expected valid IP")
	}
	if err := b.ValidateIP("invalid"); err == nil {
		t.Fatal("expected invalid IP to fail")
	}
}

func TestMatchPathInList(t *testing.T) {
	if !MatchPathInList("/api/users", []string{"/api/"}) {
		t.Fatal("expected prefix match")
	}
	if !MatchPathInList("/api/users", []string{"/api/users"}) {
		t.Fatal("expected exact match")
	}
	if !MatchPathInList("/v1/webhooks/WEBHOOK_TYPE_VERSION_BUILD", []string{"/v1/webhooks/*"}) {
		t.Fatal("expected wildcard match")
	}
	if !MatchPathInList("/api/v1/builds/callback", []string{"/api/*/builds/*"}) {
		t.Fatal("expected middle wildcard match")
	}
	if !MatchPathInList("/api/v1/users", []string{"/api/?1/users"}) {
		t.Fatal("expected question wildcard match")
	}
	if MatchPathInList("/v1/builds/callback", []string{"/v1/webhooks/*"}) {
		t.Fatal("expected wildcard to not match different prefix")
	}
	if MatchPathInList("/api/users", []string{""}) {
		t.Fatal("expected empty pattern to not match")
	}
	if MatchPathInList("/api/users", []string{}) {
		t.Fatal("expected empty list to not match")
	}
}

func TestCompileIPSetAllowAll(t *testing.T) {
	set, err := CompileIPSet([]string{"*"})
	if err != nil || !set.allowAll {
		t.Fatal("expected allow all")
	}
}

func TestCompileIPSetEmptyPattern(t *testing.T) {
	set, err := CompileIPSet([]string{""})
	if err != nil {
		t.Fatal(err)
	}
	if set.Contains("10.0.0.1") {
		t.Fatal("expected empty pattern to not match")
	}
}

func TestCompileIPSetExact(t *testing.T) {
	set, err := CompileIPSet([]string{"192.168.1.1"})
	if err != nil {
		t.Fatal(err)
	}
	if !set.Contains("192.168.1.1") {
		t.Fatal("expected exact match")
	}
	if set.Contains("192.168.1.2") {
		t.Fatal("expected no match for different IP")
	}
}

func TestCompileIPSetExactFromSplit(t *testing.T) {
	set, err := CompileIPSet([]string{"192.168.1.1 192.168.1.2"})
	if err != nil {
		t.Fatal(err)
	}
	if !set.Contains("192.168.1.1") {
		t.Fatal("expected first exact match")
	}
	if !set.Contains("192.168.1.2") {
		t.Fatal("expected second exact match")
	}
}

func TestCompileIPSetCIDR(t *testing.T) {
	set, err := CompileIPSet([]string{"10.0.0.0/8"})
	if err != nil {
		t.Fatal(err)
	}
	if !set.Contains("10.1.2.3") {
		t.Fatal("expected CIDR match")
	}
}

func TestCompileIPSetRange(t *testing.T) {
	set, err := CompileIPSet([]string{"192.168.1.1-192.168.1.100"})
	if err != nil {
		t.Fatal(err)
	}
	if !set.Contains("192.168.1.50") {
		t.Fatal("expected range match")
	}
}

func TestCompileIPSetWildcard(t *testing.T) {
	set, err := CompileIPSet([]string{"192.168.1.*"})
	if err != nil {
		t.Fatal(err)
	}
	if !set.Contains("192.168.1.50") {
		t.Fatal("expected wildcard match")
	}
}

func TestCompileIPSetInvalidRange(t *testing.T) {
	_, err := CompileIPSet([]string{"invalid-start-192.168.1.100"})
	if err == nil {
		t.Fatal("expected invalid range to fail")
	}
}

func TestCompileIPSetInvalidIP(t *testing.T) {
	_, err := CompileIPSet([]string{"abc"})
	if err == nil {
		t.Fatal("expected invalid IP to fail")
	}
}

func TestMustCompileIPSet(t *testing.T) {
	set := MustCompileIPSet([]string{"10.0.0.0/8"})
	if set == nil {
		t.Fatal("expected non-nil set")
	}
}

func TestMustCompileIPSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	MustCompileIPSet([]string{"not-an-ip"})
}

func TestIPSetContainsAllowAll(t *testing.T) {
	set, _ := CompileIPSet([]string{"*"})
	if !set.Contains("10.0.0.1") {
		t.Fatal("expected allow all to contain any IP")
	}
}

func TestMatchIPInListError(t *testing.T) {
	if MatchIPInList("10.0.0.1", []string{"not-an-ip"}) {
		t.Fatal("expected invalid pattern to not match")
	}
}

func TestIPSetContainsNil(t *testing.T) {
	var set *IPSet
	if set.Contains("10.0.0.1") {
		t.Fatal("expected nil set to not contain")
	}
}

func TestIPSetContainsInvalidIP(t *testing.T) {
	set, _ := CompileIPSet([]string{"10.0.0.0/8"})
	if set.Contains("invalid") {
		t.Fatal("expected invalid IP to not match")
	}
}

func TestIsIPAllowedEmpty(t *testing.T) {
	if !IsIPAllowed("10.0.0.1", nil) {
		t.Fatal("expected empty list to allow all")
	}
}

func TestIsIPAllowedWildcard(t *testing.T) {
	if !IsIPAllowed("10.0.0.1", []string{"*"}) {
		t.Fatal("expected wildcard to allow")
	}
}

func TestIsIPAllowedCIDR(t *testing.T) {
	if !IsIPAllowed("10.1.2.3", []string{"10.0.0.0/8"}) {
		t.Fatal("expected CIDR to allow")
	}
}

func TestIsIPBlockedEmpty(t *testing.T) {
	if IsIPBlocked("10.0.0.1", nil) {
		t.Fatal("expected empty list to not block")
	}
}

func TestIsIPBlockedMatch(t *testing.T) {
	if !IsIPBlocked("10.0.0.1", []string{"10.0.0.0/8"}) {
		t.Fatal("expected CIDR to block")
	}
}

func TestMatchIPPatternWildcard(t *testing.T) {
	if !MatchIPPattern("10.0.0.1", "*") {
		t.Fatal("expected wildcard to match")
	}
}

func TestMatchIPPatternExact(t *testing.T) {
	if !MatchIPPattern("10.0.0.1", "10.0.0.1") {
		t.Fatal("expected exact match")
	}
}

func TestMatchIPPatternCIDR(t *testing.T) {
	if !MatchIPPattern("10.1.2.3", "10.0.0.0/8") {
		t.Fatal("expected CIDR match")
	}
}

func TestIsIPInRange(t *testing.T) {
	if !IsIPInRange(netParseIP("192.168.1.50"), netParseIP("192.168.1.1"), netParseIP("192.168.1.100")) {
		t.Fatal("expected IP in range")
	}
	if IsIPInRange(netParseIP("192.168.2.1"), netParseIP("192.168.1.1"), netParseIP("192.168.1.100")) {
		t.Fatal("expected IP out of range")
	}
}

func TestIsIPInRangeIPv6(t *testing.T) {
	if IsIPInRange(netParseIP("::1"), netParseIP("192.168.1.1"), netParseIP("192.168.1.100")) {
		t.Fatal("expected IPv6 to not match IPv4 range")
	}
}

func TestMatchIPWithWildcard(t *testing.T) {
	if !MatchIPWithWildcard("192.168.1.50", "192.168.1.*") {
		t.Fatal("expected wildcard match")
	}
	if MatchIPWithWildcard("192.168.2.50", "192.168.1.*") {
		t.Fatal("expected wildcard not match")
	}
	if MatchIPWithWildcard("192.168.1.50", "192.168.*") {
		t.Fatal("expected different length to not match")
	}
}

func TestIsPrivateIP(t *testing.T) {
	if !IsPrivateIP("10.0.0.1") {
		t.Fatal("expected 10.x to be private")
	}
	if !IsPrivateIP("192.168.1.1") {
		t.Fatal("expected 192.168.x to be private")
	}
	if !IsPrivateIP("127.0.0.1") {
		t.Fatal("expected 127.x to be private")
	}
	if IsPrivateIP("8.8.8.8") {
		t.Fatal("expected 8.8.8.8 to be public")
	}
}

func TestSplitIPPatterns(t *testing.T) {
	patterns := splitIPPatterns("10.0.0.0/8\n172.16.0.0/12")
	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(patterns))
	}
}

func TestSplitIPPatternsComma(t *testing.T) {
	patterns := splitIPPatterns("10.0.0.0/8,172.16.0.0/12")
	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(patterns))
	}
}

func TestSplitIPPatternsRange(t *testing.T) {
	patterns := splitIPPatterns("192.168.1.1 - 192.168.1.100")
	if len(patterns) != 1 || patterns[0] != "192.168.1.1 - 192.168.1.100" {
		t.Fatalf("expected range pattern, got %v", patterns)
	}
}

func TestSplitIPPatternsSemicolon(t *testing.T) {
	patterns := splitIPPatterns("10.0.0.0/8;172.16.0.0/12")
	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(patterns))
	}
}

func TestSplitIPPatternsCRLF(t *testing.T) {
	patterns := splitIPPatterns("10.0.0.0/8\r\n172.16.0.0/12")
	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(patterns))
	}
}
