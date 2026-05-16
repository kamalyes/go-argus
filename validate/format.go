/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\format.go
 * @Description: 格式校验能力，提供 Email、IP、URL、UUID、Base64 和正则校验
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/kamalyes/go-argus/i18n"
)

var (
	uuidRegex    = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	regexCache   = make(map[string]*regexp.Regexp)
	regexCacheMu sync.RWMutex
)

// GetCompiledRegex 获取编译的正则（带缓存）
func GetCompiledRegex(pattern string) (*regexp.Regexp, error) {
	regexCacheMu.RLock()
	cached, exists := regexCache[pattern]
	regexCacheMu.RUnlock()

	if exists {
		return cached, nil
	}

	regexCacheMu.Lock()
	defer regexCacheMu.Unlock()

	if cached, exists := regexCache[pattern]; exists {
		return cached, nil
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	regexCache[pattern] = compiled
	return compiled, nil
}

// ClearRegexCache 清空正则缓存
func ClearRegexCache() {
	regexCacheMu.Lock()
	defer regexCacheMu.Unlock()
	regexCache = make(map[string]*regexp.Regexp)
}

// ValidateRegex 校验字节内容是否匹配正则表达式
func ValidateRegex(body []byte, pattern string) CompareResult {
	result := CompareResult{Actual: string(body), Expect: pattern}
	re, err := regexp.Compile(pattern)
	if err != nil {
		result.Message = i18n.Msg(MsgFormatRegexCompileFailed, map[string]string{"error": err.Error()})
		return result
	}
	result.Success = re.Match(body)
	if !result.Success {
		result.Message = i18n.Msg(MsgFormatRegexNotMatched)
	}
	return result
}

// ValidateEmail 校验 Email 格式
func ValidateEmail(email string) CompareResult {
	result := CompareResult{Actual: email, Expect: "valid email format"}
	email = strings.TrimSpace(email)
	if email == "" {
		result.Message = i18n.Msg(MsgFormatEmailEmpty)
		return result
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		result.Message = i18n.Msg(MsgFormatEmailInvalid, map[string]string{"error": err.Error()})
		return result
	}
	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" || !strings.Contains(parts[1], ".") {
		result.Message = i18n.Msg(MsgFormatEmailMalformed)
		return result
	}
	result.Success = true
	return result
}

// ValidateIPAddress 校验 IP 地址格式
func ValidateIPAddress(ipStr string) CompareResult {
	result := CompareResult{Actual: ipStr, Expect: "valid IP address"}
	if net.ParseIP(strings.TrimSpace(ipStr)) == nil {
		result.Message = i18n.Msg(MsgFormatIPInvalid, map[string]string{"value": ipStr})
		return result
	}
	result.Success = true
	return result
}

// ValidateProtocol 校验 URL 协议是否在允许列表中
func ValidateProtocol(urlStr string, allowedProtocols ...string) CompareResult {
	result := CompareResult{Actual: urlStr}
	if len(allowedProtocols) == 0 {
		allowedProtocols = []string{"http", "https", "ws", "wss", "ftp", "ftps"}
	}
	result.Expect = fmt.Sprintf("valid URL protocol: %v", allowedProtocols)
	u, err := url.Parse(strings.TrimSpace(urlStr))
	if err != nil || u.Scheme == "" {
		result.Message = i18n.Msg(MsgFormatURLMissingProtocol)
		return result
	}
	for _, allowed := range allowedProtocols {
		if strings.EqualFold(u.Scheme, allowed) {
			result.Success = true
			return result
		}
	}
	result.Message = i18n.Msg(MsgFormatURLUnsupportedProtocol, map[string]string{"value": u.Scheme})
	return result
}

// ValidateHTTP 校验 HTTP/HTTPS URL
func ValidateHTTP(urlStr string) CompareResult {
	return ValidateProtocol(urlStr, "http", "https")
}

// ValidateWebSocket 校验 WebSocket URL
func ValidateWebSocket(urlStr string) CompareResult {
	return ValidateProtocol(urlStr, "ws", "wss")
}

// ValidateUUID 校验 UUID 格式
func ValidateUUID(uuidStr string) CompareResult {
	result := CompareResult{Actual: uuidStr, Expect: "valid UUID format"}
	if !uuidRegex.MatchString(strings.TrimSpace(uuidStr)) {
		result.Message = i18n.Msg(MsgFormatUUIDInvalid)
		return result
	}
	result.Success = true
	return result
}

// ValidateBase64 校验 Base64 字符串
func ValidateBase64(str string) CompareResult {
	result := CompareResult{Actual: str, Expect: "valid Base64 encoding"}
	str = strings.TrimSpace(str)
	if str == "" {
		result.Message = i18n.Msg(MsgFormatBase64Empty)
		return result
	}
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.RawURLEncoding} {
		if _, err := enc.DecodeString(str); err == nil {
			result.Success = true
			return result
		}
	}
	result.Message = i18n.Msg(MsgFormatBase64Invalid)
	return result
}

// IsEmail 判断字符串是否为有效 Email
func IsEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	n := len(email)
	if n > 254 {
		return false
	}
	at := strings.LastIndex(email, "@")
	if at <= 0 || at == n-1 {
		return false
	}
	local := email[:at]
	domain := email[at+1:]
	if local == "" || len(local) > 64 {
		return false
	}
	if !isEmailLocal(local) {
		return false
	}
	if !isEmailDomain(domain) {
		return false
	}
	return true
}

func isEmailLocal(local string) bool {
	for i := 0; i < len(local); i++ {
		c := local[i]
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
			continue
		}
		if c == '.' || c == '!' || c == '#' || c == '$' || c == '%' || c == '&' || c == '\'' {
			continue
		}
		if c == '*' || c == '+' || c == '/' || c == '=' || c == '?' || c == '^' || c == '_' || c == '`' {
			continue
		}
		if c == '{' || c == '|' || c == '}' || c == '~' || c == '-' {
			continue
		}
		return false
	}
	if local[0] == '.' || local[len(local)-1] == '.' {
		return false
	}
	for i := 0; i < len(local)-1; i++ {
		if local[i] == '.' && local[i+1] == '.' {
			return false
		}
	}
	return true
}

func isEmailDomain(domain string) bool {
	if domain == "" || len(domain) > 253 {
		return false
	}
	if domain[0] == '.' || domain[len(domain)-1] == '.' {
		return false
	}
	hasDot := false
	start := 0
	for i := 0; i <= len(domain); i++ {
		if i == len(domain) || domain[i] == '.' {
			if i == start {
				return false
			}
			if !isDomainLabel(domain[start:i]) {
				return false
			}
			if i < len(domain) {
				if i > 0 && domain[i-1] == '.' {
					return false
				}
				hasDot = true
			}
			start = i + 1
		}
	}
	return hasDot
}

func isDomainLabel(label string) bool {
	if label == "" || len(label) > 63 {
		return false
	}
	if label[0] == '-' || label[len(label)-1] == '-' {
		return false
	}
	for i := 0; i < len(label); i++ {
		c := label[i]
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' {
			continue
		}
		return false
	}
	return true
}

// IsIP 判断字符串是否为有效 IP
func IsIP(ip string) bool {
	return net.ParseIP(strings.TrimSpace(ip)) != nil
}

// IsUUID 判断字符串是否为有效 UUID
func IsUUID(uuid string) bool {
	return uuidRegex.MatchString(strings.TrimSpace(uuid))
}

// IsBase64 判断字符串是否为有效 Base64
func IsBase64(str string) bool {
	str = strings.TrimSpace(str)
	if str == "" {
		return false
	}
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.RawURLEncoding} {
		if _, err := enc.DecodeString(str); err == nil {
			return true
		}
	}
	return false
}

// ValidateIP 校验 IP 地址格式（ValidateIPAddress 的别名，兼容旧 API）
func ValidateIP(ipStr string) CompareResult {
	return ValidateIPAddress(ipStr)
}
