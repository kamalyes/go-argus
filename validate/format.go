/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 01:53:23
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
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/kamalyes/go-argus/i18n"
)

var (
	uuidRegex    = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	semverRegex  = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	bicRegex     = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}(?:[A-Z0-9]{3})?$`)
	ethAddrRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	btcAddrRegex = regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$|^[bc1q][a-z0-9]{39,59}$`)
	bcp47Regex   = regexp.MustCompile(`^[a-zA-Z]{2,3}(?:-[a-zA-Z]{4})?(?:-(?:[a-zA-Z]{2}|\d{3}))?(?:-[a-zA-Z0-9]{5,8})*(?:-[a-zA-Z0-9]{1,8})*$`)
	datauriRegex = regexp.MustCompile(`^data:([^;,]*)?(;[^;,]*)*;?(base64,)?,`)
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
	if !IsEmail(email) {
		result.Message = i18n.Msg(MsgFormatEmailInvalid, map[string]string{"error": "invalid email format"})
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

func ValidateSemver(version string) CompareResult {
	result := CompareResult{Actual: version, Expect: "valid semantic version"}
	if !semverRegex.MatchString(version) {
		result.Message = i18n.Msg(MsgFormatSemverInvalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateISBN10(isbn string) CompareResult {
	result := CompareResult{Actual: isbn, Expect: "valid ISBN-10"}
	if !isISBN10(isbn) {
		result.Message = i18n.Msg(MsgFormatISBN10Invalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateISBN13(isbn string) CompareResult {
	result := CompareResult{Actual: isbn, Expect: "valid ISBN-13"}
	if !isISBN13(isbn) {
		result.Message = i18n.Msg(MsgFormatISBN13Invalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateISSN(issn string) CompareResult {
	result := CompareResult{Actual: issn, Expect: "valid ISSN"}
	if !isISSN(issn) {
		result.Message = i18n.Msg(MsgFormatISSNInvalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateBIC(bic string) CompareResult {
	result := CompareResult{Actual: bic, Expect: "valid BIC (SWIFT code)"}
	if !bicRegex.MatchString(bic) {
		result.Message = i18n.Msg(MsgFormatBICInvalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateCron(expr string) CompareResult {
	result := CompareResult{Actual: expr, Expect: "valid cron expression"}
	if !isCron(expr) {
		result.Message = i18n.Msg(MsgFormatCronInvalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateDataURI(uri string) CompareResult {
	result := CompareResult{Actual: uri, Expect: "valid Data URI"}
	if !strings.HasPrefix(uri, "data:") || !datauriRegex.MatchString(uri) {
		result.Message = i18n.Msg(MsgFormatDataURIInvalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateBCP47(tag string) CompareResult {
	result := CompareResult{Actual: tag, Expect: "valid BCP 47 language tag"}
	if !bcp47Regex.MatchString(tag) {
		result.Message = i18n.Msg(MsgFormatBCP47Invalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateEthAddr(addr string) CompareResult {
	result := CompareResult{Actual: addr, Expect: "valid Ethereum address"}
	if !ethAddrRegex.MatchString(addr) {
		result.Message = i18n.Msg(MsgFormatEthAddrInvalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateBtcAddr(addr string) CompareResult {
	result := CompareResult{Actual: addr, Expect: "valid Bitcoin address"}
	if !btcAddrRegex.MatchString(addr) {
		result.Message = i18n.Msg(MsgFormatBtcAddrInvalid)
		return result
	}
	result.Success = true
	return result
}

func IsSemver(version string) bool {
	return semverRegex.MatchString(version)
}

func IsISBN10(isbn string) bool {
	return isISBN10(isbn)
}

func IsISBN13(isbn string) bool {
	return isISBN13(isbn)
}

func IsISSN(issn string) bool {
	return isISSN(issn)
}

func IsBIC(bic string) bool {
	return bicRegex.MatchString(bic)
}

func IsCron(expr string) bool {
	return isCron(expr)
}

func IsDataURI(uri string) bool {
	return strings.HasPrefix(uri, "data:") && datauriRegex.MatchString(uri)
}

func IsBCP47(tag string) bool {
	return bcp47Regex.MatchString(tag)
}

func IsEthAddr(addr string) bool {
	return ethAddrRegex.MatchString(addr)
}

func IsBtcAddr(addr string) bool {
	return btcAddrRegex.MatchString(addr)
}

func isISBN10(s string) bool {
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, " ", "")
	if len(s) != 10 {
		return false
	}
	sum := 0
	for i := 0; i < 9; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
		sum += int(s[i]-'0') * (10 - i)
	}
	last := s[9]
	if last == 'X' || last == 'x' {
		sum += 10
	} else if last >= '0' && last <= '9' {
		sum += int(last - '0')
	} else {
		return false
	}
	return sum%11 == 0
}

func isISBN13(s string) bool {
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, " ", "")
	if len(s) != 13 {
		return false
	}
	sum := 0
	for i := 0; i < 12; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
		weight := 1
		if i%2 == 1 {
			weight = 3
		}
		sum += int(s[i]-'0') * weight
	}
	check := (10 - sum%10) % 10
	return s[12] >= '0' && s[12] <= '9' && int(s[12]-'0') == check
}

func isISSN(s string) bool {
	s = strings.ReplaceAll(s, "-", "")
	if len(s) != 8 {
		return false
	}
	sum := 0
	for i := 0; i < 7; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
		sum += int(s[i]-'0') * (8 - i)
	}
	last := s[7]
	var check int
	if last == 'X' || last == 'x' {
		check = 10
	} else if last >= '0' && last <= '9' {
		check = int(last - '0')
	} else {
		return false
	}
	return sum%11 == check
}

func isCron(expr string) bool {
	fields := strings.Fields(expr)
	if len(fields) != 5 && len(fields) != 6 {
		return false
	}
	for _, f := range fields {
		if !isValidCronField(f) {
			return false
		}
	}
	return true
}

func isValidCronField(field string) bool {
	parts := strings.Split(field, ",")
	for _, part := range parts {
		stepParts := strings.SplitN(part, "/", 2)
		base := stepParts[0]
		if len(stepParts) == 2 {
			step := stepParts[1]
			if step == "" {
				return false
			}
			for _, c := range step {
				if c != '*' && (c < '0' || c > '9') {
					return false
				}
			}
		}
		if base == "*" {
			continue
		}
		rangeParts := strings.SplitN(base, "-", 2)
		for _, rp := range rangeParts {
			if rp == "" {
				return false
			}
			for _, c := range rp {
				if c < '0' || c > '9' {
					return false
				}
			}
		}
	}
	return true
}
