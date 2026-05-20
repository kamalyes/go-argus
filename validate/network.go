/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\network.go
 * @Description: 网络校验能力，提供 IP 白名单、黑名单、CIDR、通配符、IP 范围和路径匹配
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/kamalyes/go-argus/i18n"
)

// IPBase 提供面向对象风格的 IP 校验入口，兼容 go-toolbox 旧 API
type IPBase struct{}

// ValidateIP 校验 IP 字符串是否有效
func (b *IPBase) ValidateIP(ip string) error {
	if net.ParseIP(strings.TrimSpace(ip)) == nil {
		return fmt.Errorf(i18n.Msg(MsgNetworkIPInvalid, map[string]string{"value": ip}))
	}
	return nil
}

// MatchPathInList 判断路径是否命中任意路径前缀
func MatchPathInList(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}
		if path == pattern || strings.HasPrefix(path, pattern) {
			return true
		}
	}
	return false
}

// IPSet 表示预编译 IP 规则集合，适合高频请求路径复用
type IPSet struct {
	allowAll  bool
	exact     map[string]struct{}
	cidrs     []*net.IPNet
	ranges    [][2]net.IP
	wildcards []string
}

// CompileIPSet 将 IP 规则编译为可复用集合
func CompileIPSet(patterns []string) (*IPSet, error) {
	set := &IPSet{exact: make(map[string]struct{})}
	for _, raw := range patterns {
		for _, pattern := range splitIPPatterns(raw) {
			if pattern == "*" {
				set.allowAll = true
				continue
			}
			if strings.Contains(pattern, "*") {
				set.wildcards = append(set.wildcards, pattern)
				continue
			}
			if strings.Contains(pattern, "-") {
				parts := strings.SplitN(pattern, "-", 2)
				start := net.ParseIP(strings.TrimSpace(parts[0]))
				end := net.ParseIP(strings.TrimSpace(parts[1]))
				if start == nil || end == nil {
					return nil, fmt.Errorf(i18n.Msg(MsgNetworkIPRangeInvalid, map[string]string{"value": pattern}))
				}
				set.ranges = append(set.ranges, [2]net.IP{start, end})
				continue
			}
			if _, ipNet, err := net.ParseCIDR(pattern); err == nil {
				set.cidrs = append(set.cidrs, ipNet)
				continue
			}
			ip := net.ParseIP(pattern)
			if ip == nil {
				return nil, fmt.Errorf(i18n.Msg(MsgNetworkIPRuleInvalid, map[string]string{"value": pattern}))
			}
			set.exact[pattern] = struct{}{}
		}
	}
	return set, nil
}

// MustCompileIPSet 编译 IP 规则，失败时 panic，适合初始化阶段使用
func MustCompileIPSet(patterns []string) *IPSet {
	set, err := CompileIPSet(patterns)
	if err != nil {
		panic(err)
	}
	return set
}

// Contains 判断 IP 是否命中规则集合
func (set *IPSet) Contains(ip string) bool {
	if set == nil {
		return false
	}
	rawIP := strings.TrimSpace(ip)
	parsedIP := net.ParseIP(rawIP)
	if parsedIP == nil {
		return false
	}
	if set.allowAll {
		return true
	}
	if _, ok := set.exact[rawIP]; ok {
		return true
	}
	for _, ipNet := range set.cidrs {
		if ipNet.Contains(parsedIP) {
			return true
		}
	}
	for _, item := range set.ranges {
		if IsIPInRange(parsedIP, item[0], item[1]) {
			return true
		}
	}
	for _, pattern := range set.wildcards {
		if MatchIPWithWildcard(rawIP, pattern) {
			return true
		}
	}
	return false
}

// IsIPAllowed 判断 IP 是否在允许列表中，空列表表示允许所有
func IsIPAllowed(ip string, cidrList []string) bool {
	if len(cidrList) == 0 {
		return true
	}
	for _, cidr := range cidrList {
		if strings.TrimSpace(cidr) == "*" {
			return true
		}
	}
	return MatchIPInList(ip, cidrList)
}

// IsIPBlocked 判断 IP 是否在黑名单中，空列表表示不阻止
func IsIPBlocked(ip string, blacklist []string) bool {
	if len(blacklist) == 0 {
		return false
	}
	return MatchIPInList(ip, blacklist)
}

// MatchIPPattern 判断 IP 是否命中单个规则
func MatchIPPattern(ip, pattern string) bool {
	if strings.TrimSpace(pattern) == "*" {
		return true
	}
	if ip == pattern {
		return true
	}
	return MatchIPInList(ip, []string{pattern})
}

// MatchIPInList 判断 IP 是否命中任意规则
func MatchIPInList(ip string, ipList []string) bool {
	set, err := CompileIPSet(ipList)
	if err != nil {
		return false
	}
	return set.Contains(ip)
}

// IsIPInRange 判断 IPv4 地址是否落在闭区间内
func IsIPInRange(ip, start, end net.IP) bool {
	ipv4 := ip.To4()
	startv4 := start.To4()
	endv4 := end.To4()
	if ipv4 == nil || startv4 == nil || endv4 == nil {
		return false
	}
	return bytes.Compare(ipv4, startv4) >= 0 && bytes.Compare(ipv4, endv4) <= 0
}

// MatchIPWithWildcard 使用星号通配符匹配 IPv4
func MatchIPWithWildcard(ip, pattern string) bool {
	ipParts := strings.Split(ip, ".")
	patternParts := strings.Split(pattern, ".")
	if len(ipParts) != len(patternParts) {
		return false
	}
	for i := 0; i < len(ipParts); i++ {
		if patternParts[i] != "*" && ipParts[i] != patternParts[i] {
			return false
		}
	}
	return true
}

// IsPrivateIP 判断 IP 是否属于常见私有或本地地址段
func IsPrivateIP(ip string) bool {
	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}
	return MatchIPInList(ip, privateBlocks)
}

func splitIPPatterns(pattern string) []string {
	pattern = strings.ReplaceAll(pattern, "\r\n", "\n")
	pattern = strings.NewReplacer(";", "\n", ",", "\n", "\t", "\n").Replace(pattern)
	lines := strings.Split(pattern, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(line, "-") {
			line = strings.TrimSpace(line)
			if line != "" {
				out = append(out, line)
			}
			continue
		}
		out = append(out, strings.Fields(line)...)
	}
	return out
}

// isHostnameLabel 判断单个 label 是否符合 hostname 规范（零正则、零分配）
// 规则：首尾必须是字母或数字，中间允许字母、数字、连字符，长度 1-63
func isHostnameLabel(label string) bool {
	n := len(label)
	if n == 0 || n > 63 {
		return false
	}
	// 首字符必须是字母或数字
	c := label[0]
	if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
		return false
	}
	// 尾字符必须是字母或数字
	c = label[n-1]
	if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
		return false
	}
	// 中间字符允许字母、数字、连字符
	for i := 1; i < n-1; i++ {
		c = label[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}

// IsHostname 判断字符串是否为合法 hostname（零正则、零分配）
func IsHostname(host string) bool {
	if host == "" || len(host) > 253 {
		return false
	}
	start := 0
	for i := 0; i <= len(host); i++ {
		if i == len(host) || host[i] == '.' {
			if !isHostnameLabel(host[start:i]) {
				return false
			}
			start = i + 1
		}
	}
	return true
}

func TrimSpaceIfNeeded(s string) string {
	if len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		return strings.TrimSpace(s)
	}
	return s
}

func HasSchemeAndHost(s string) bool {
	colon := strings.Index(s, ":")
	if colon < 1 {
		return false
	}
	return HasHostAfterScheme(s, colon)
}

func HasHostAfterScheme(s string, colonIdx int) bool {
	if len(s) <= colonIdx+2 || s[colonIdx+1] != '/' || s[colonIdx+2] != '/' {
		return false
	}
	hostStart := colonIdx + 3
	if hostStart >= len(s) {
		return false
	}
	for i := hostStart; i < len(s); i++ {
		if IsHostTerminator(s[i]) {
			break
		}
		if !IsValidHostChar(s, hostStart, i) {
			return false
		}
	}
	return true
}

func IsHostTerminator(c byte) bool {
	return c == '/' || c == '?' || c == '#'
}

func IsValidHostChar(s string, hostStart, i int) bool {
	c := s[i]
	if c == ':' || c == '@' {
		return true
	}
	return !(hostStart == i && (c == '.' || c == '-'))
}
