/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-16 23:18:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-25 23:19:56
 * @FilePath: \go-argus\validate\ip_util.go
 * @Description: IP 地址工具函数，补充 network.go 的 IP 校验能力
 *
 * 提供链路本地、唯一本地、全球单播、文档地址等 IP 属性判断
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import "net"

// HasLocalIP 判断 IP 字符串是否为本地或私有地址
// 包括：localhost、127.0.0.1、::1，以及 IPv4 私有地址段
func HasLocalIP(ip string) bool {
	// 快速检查硬编码的本地地址
	switch ip {
	case "localhost", "127.0.0.1", "::1":
		return true
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// IPv4 私有地址段
	if ip4 := parsedIP.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return true // 10.0.0.0/8
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true // 172.16.0.0/12
		case ip4[0] == 192 && ip4[1] == 168:
			return true // 192.168.0.0/16
		case ip4[0] == 127:
			return true // 127.0.0.0/8 loopback
		case ip4[0] == 169 && ip4[1] == 254:
			return true // 169.254.0.0/16 link-local
		}
		return false
	}

	// IPv6 回环地址
	return parsedIP.IsLoopback()
}

// IsLinkLocalIP 检查 IP 是否为链路本地地址
// IPv4: 169.254.0.0/16
// IPv6: fe80::/10
func IsLinkLocalIP(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 169 && ip4[1] == 254
	}
	if ip6 := ip.To16(); ip6 != nil {
		return ip6[0] == 0xfe && (ip6[1]&0xc0) == 0x80 // fe80::/10
	}
	return false
}

// IsUniqueLocalAddress 检查 IP 是否为 IPv6 唯一本地地址（ULA）
// ULA 范围：fc00::/7（包含 fc00::/8 和 fd00::/8）
func IsUniqueLocalAddress(ip net.IP) bool {
	if ip6 := ip.To16(); ip6 != nil {
		return ip6[0] == 0xFC || ip6[0] == 0xFD
	}
	return false
}

// IsGlobalUnicast 检查 IP 字符串是否为全球单播地址
// 排除：私有地址、回环、链路本地、文档地址等
func IsGlobalUnicast(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	if ip4 := parsedIP.To4(); ip4 != nil {
		// IPv4：排除私有、回环、链路本地
		switch {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		case ip4[0] == 127:
			return false
		case ip4[0] == 169 && ip4[1] == 254:
			return false
		}
		return true
	}

	if ip6 := parsedIP.To16(); ip6 != nil {
		return !parsedIP.IsUnspecified() &&
			!parsedIP.IsLoopback() &&
			!IsLinkLocalIP(parsedIP) &&
			!parsedIP.IsPrivate() &&
			!IsDocumentationAddress(ip6)
	}

	return false
}

// IsDocumentationAddress 检查 IP 是否为文档专用地址
// IPv6 文档地址前缀：2001:db8::/32
func IsDocumentationAddress(ip net.IP) bool {
	return len(ip) >= 4 &&
		ip[0] == 0x20 && ip[1] == 0x01 &&
		ip[2] == 0x0d && ip[3] == 0xb8
}
