/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-21 00:00:00
 * @FilePath: \go-argus\validate\ip_util_test.go
 * @Description: ip_util.go 测试，覆盖本地IP、链路本地、唯一本地、全球单播、文档地址判断
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"net"
	"testing"
)

func TestHasLocalIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"回环地址", "127.0.0.1", true},
		{"localhost", "localhost", true},
		{"10私有段", "10.0.0.1", true},
		{"192.168私有段", "192.168.1.1", true},
		{"公网地址", "8.8.8.8", false},
		{"无效地址", "invalid", false},
		{"IPv6回环", "::1", true},
		{"172.16私有段", "172.16.0.1", true},
		{"172.32非私有段", "172.32.0.1", false},
		{"169.254链路本地", "169.254.1.1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasLocalIP(tt.ip); got != tt.want {
				t.Errorf("HasLocalIP(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIsLinkLocalIP(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want bool
	}{
		{"IPv4链路本地", net.ParseIP("169.254.0.1"), true},
		{"IPv4非链路本地", net.ParseIP("10.0.0.1"), false},
		{"IPv6链路本地", net.ParseIP("fe80::1"), true},
		{"IPv6非链路本地", net.ParseIP("2001:db8::1"), false},
		{"IPv4回环", net.ParseIP("127.0.0.1"), false},
		{"IPv4公网", net.ParseIP("8.8.8.8"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLinkLocalIP(tt.ip); got != tt.want {
				t.Errorf("IsLinkLocalIP(%v) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIsUniqueLocalAddress(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want bool
	}{
		{"fc00 ULA", net.ParseIP("fc00::1"), true},
		{"fd00 ULA", net.ParseIP("fd00::1"), true},
		{"非ULA IPv6", net.ParseIP("2001:db8::1"), false},
		{"IPv4地址", net.ParseIP("8.8.8.8"), false},
		{"IPv6回环", net.ParseIP("::1"), false},
		{"IPv6链路本地", net.ParseIP("fe80::1"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUniqueLocalAddress(tt.ip); got != tt.want {
				t.Errorf("IsUniqueLocalAddress(%v) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIsGlobalUnicast(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"IPv4公网", "8.8.8.8", true},
		{"10私有段", "10.0.0.1", false},
		{"回环地址", "127.0.0.1", false},
		{"链路本地", "169.254.1.1", false},
		{"192.168私有段", "192.168.1.1", false},
		{"无效地址", "invalid", false},
		{"IPv6公网", "2001:4860:4860::8888", true},
		{"IPv6文档地址", "2001:db8::1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGlobalUnicast(tt.ip); got != tt.want {
				t.Errorf("IsGlobalUnicast(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIsDocumentationAddress(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want bool
	}{
		{"IPv6文档地址", net.ParseIP("2001:db8::1"), true},
		{"IPv4公网", net.ParseIP("8.8.8.8"), false},
		{"IPv6链路本地", net.ParseIP("fe80::1"), false},
		{"IPv6回环", net.ParseIP("::1"), false},
		{"IPv6 ULA", net.ParseIP("fc00::1"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDocumentationAddress(tt.ip); got != tt.want {
				t.Errorf("IsDocumentationAddress(%v) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}
