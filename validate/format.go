/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 02:55:52
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
	if !IsUUID(uuidStr) {
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

func ParseSemverNum(s string, pos *int) bool {
	return parseSemverNumFmt(s, pos)
}

func ParseSemverPreRelease(s string, pos *int) bool {
	return parseSemverPreReleaseFmt(s, pos)
}

func ParseSemverBuildMeta(s string, pos *int) bool {
	return parseSemverBuildMetaFmt(s, pos)
}

func IsValidCronField(field string) bool {
	inRange := false
	for i := 0; i < len(field); i++ {
		if !validateCronChar(field, i, &inRange) {
			return false
		}
	}
	return true
}

func LuhnDouble(n int) int {
	n *= 2
	if n > 9 {
		return n - 9
	}
	return n
}

func IsISBN10CheckDigit(c byte, sum int) bool {
	if c == 'X' || c == 'x' {
		return (sum+10)%11 == 0
	}
	if c < '0' || c > '9' {
		return false
	}
	return (sum+int(c-'0'))%11 == 0
}

func IsLuhnChecksum(s string) bool {
	sum := 0
	double := false
	digits := 0
	for i := len(s) - 1; i >= 0; i-- {
		r := s[i]
		if r == ' ' || r == '-' {
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
		n := int(r - '0')
		if double {
			n = LuhnDouble(n)
		}
		sum += n
		double = !double
		digits++
	}
	return digits > 0 && sum%10 == 0
}

func isSemver(s string) bool {
	i := 0
	if i < len(s) && s[i] == 'v' {
		i++
	}
	if !parseSemverNumFmt(s, &i) || i >= len(s) || s[i] != '.' {
		return false
	}
	i++
	if !parseSemverNumFmt(s, &i) || i >= len(s) || s[i] != '.' {
		return false
	}
	i++
	if !parseSemverNumFmt(s, &i) {
		return false
	}
	if !parseSemverPreReleaseFmt(s, &i) {
		return false
	}
	if !parseSemverBuildMetaFmt(s, &i) {
		return false
	}
	return i == len(s)
}

func parseSemverPreReleaseFmt(s string, i *int) bool {
	if *i >= len(s) || s[*i] != '-' {
		return true
	}
	*i++
	if !parseSemverIdentFmt(s, i) {
		return false
	}
	for *i < len(s) && s[*i] == '.' {
		*i++
		if !parseSemverIdentFmt(s, i) {
			return false
		}
	}
	return true
}

func parseSemverBuildMetaFmt(s string, i *int) bool {
	if *i >= len(s) || s[*i] != '+' {
		return true
	}
	*i++
	if !parseSemverBuildFmt(s, i) {
		return false
	}
	for *i < len(s) && s[*i] == '.' {
		*i++
		if !parseSemverBuildFmt(s, i) {
			return false
		}
	}
	return true
}

func parseSemverNumFmt(s string, pos *int) bool {
	if *pos >= len(s) || s[*pos] < '0' || s[*pos] > '9' {
		return false
	}
	if s[*pos] == '0' {
		*pos++
		return true
	}
	for *pos < len(s) && s[*pos] >= '0' && s[*pos] <= '9' {
		*pos++
	}
	return true
}

func parseSemverIdentFmt(s string, pos *int) bool {
	start := *pos
	for *pos < len(s) && s[*pos] != '.' && s[*pos] != '+' {
		if !isSemverIdentCharFmt(s[*pos]) {
			return false
		}
		*pos++
	}
	if *pos == start {
		return false
	}
	return hasNonZeroAlphaNumFmt(s, start, *pos)
}

func isSemverIdentCharFmt(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-'
}

func hasNonZeroAlphaNumFmt(s string, start, end int) bool {
	for j := start; j < end; j++ {
		if s[j] != '-' && s[j] != '0' {
			return true
		}
	}
	return end-start <= 1
}

func parseSemverBuildFmt(s string, pos *int) bool {
	if *pos >= len(s) {
		return false
	}
	for *pos < len(s) && s[*pos] != '.' && s[*pos] != '+' {
		if !isSemverIdentCharFmt(s[*pos]) {
			return false
		}
		*pos++
	}
	return *pos > 0 && s[*pos-1] != '.' && s[*pos-1] != '-'
}

func isBIC(s string) bool {
	n := len(s)
	if n != 8 && n != 11 {
		return false
	}
	for i := 0; i < 6; i++ {
		if !isUpperCaseLetterFmt(s[i]) {
			return false
		}
	}
	for i := 6; i < n; i++ {
		if !isAlphanumFmt(s[i]) {
			return false
		}
	}
	return true
}

func isUpperCaseLetterFmt(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isDataURI(s string) bool {
	if len(s) < 6 || !hasDataPrefixFmt(s) {
		return false
	}
	i := 5
	i = skipDataURIMimeTypeFmt(s, i)
	i = skipDataURIParamsFmt(s, i)
	return i < len(s) && s[i] == ','
}

func hasDataPrefixFmt(s string) bool {
	return s[0] == 'd' && s[1] == 'a' && s[2] == 't' && s[3] == 'a' && s[4] == ':'
}

func skipDataURIMimeTypeFmt(s string, i int) int {
	for i < len(s) && s[i] != ';' && s[i] != ',' {
		if s[i] < ' ' || s[i] > '~' {
			return len(s)
		}
		i++
	}
	return i
}

func skipDataURIParamsFmt(s string, i int) int {
	for i < len(s) && s[i] == ';' {
		i++
		i = skipBase64IfPresentFmt(s, i)
		for i < len(s) && s[i] != ';' && s[i] != ',' {
			if s[i] < ' ' || s[i] > '~' {
				return len(s)
			}
			i++
		}
	}
	return i
}

func skipBase64IfPresentFmt(s string, i int) int {
	if i+6 <= len(s) && s[i] == 'b' && s[i+1] == 'a' && s[i+2] == 's' && s[i+3] == 'e' && s[i+4] == '6' && s[i+5] == '4' {
		return i + 6
	}
	return i
}

func isBCP47(s string) bool {
	if len(s) < 2 {
		return false
	}
	i := 0
	if !isAlphaFmt(s, &i, 2, 3) {
		return false
	}
	i = parseBCP47ExtLangFmt(s, i)
	i = parseBCP47ScriptFmt(s, i)
	i = parseBCP47RegionFmt(s, i)
	return parseBCP47VariantsFmt(s, i) == len(s)
}

func parseBCP47ExtLangFmt(s string, i int) int {
	if i >= len(s) || s[i] != '-' {
		return i
	}
	i++
	if i < len(s) && isAlphaAtFmt(s, i, 4) {
		i += 4
		if i < len(s) && s[i] == '-' {
			i++
		}
	}
	return i
}

func parseBCP47ScriptFmt(s string, i int) int {
	if i < len(s) && isAlphaAtFmt(s, i, 2) {
		return i + 2
	}
	return i
}

func parseBCP47RegionFmt(s string, i int) int {
	if i < len(s) && isDigitAtFmt(s, i, 3) {
		return i + 3
	}
	return i
}

func parseBCP47VariantsFmt(s string, i int) int {
	for i < len(s) && s[i] == '-' {
		i++
		start := i
		for i < len(s) && s[i] != '-' {
			if !isAlphanumFmt(s[i]) {
				return -1
			}
			i++
		}
		if i == start || i-start > 8 {
			return -1
		}
	}
	return i
}

func isAlphanumFmt(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func isAlphaFmt(s string, pos *int, minLen, maxLen int) bool {
	start := *pos
	for *pos < len(s) && *pos-start < maxLen {
		c := s[*pos]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			break
		}
		*pos++
	}
	return *pos-start >= minLen && *pos-start <= maxLen
}

func isAlphaAtFmt(s string, pos, length int) bool {
	if pos+length > len(s) {
		return false
	}
	for j := 0; j < length; j++ {
		c := s[pos+j]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}
	return true
}

func isDigitAtFmt(s string, pos, length int) bool {
	if pos+length > len(s) {
		return false
	}
	for j := 0; j < length; j++ {
		c := s[pos+j]
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isEthAddr(s string) bool {
	if len(s) != 42 {
		return false
	}
	if s[0] != '0' || (s[1] != 'x' && s[1] != 'X') {
		return false
	}
	for i := 2; i < 42; i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func isBtcAddr(s string) bool {
	n := len(s)
	if n < 26 || n > 62 {
		return false
	}
	if s[0] == '1' || s[0] == '3' {
		return isBtcLegacyAddrFmt(s, n)
	}
	return isBtcBech32AddrFmt(s, n)
}

func isBtcLegacyAddrFmt(s string, n int) bool {
	if n < 26 || n > 35 {
		return false
	}
	for i := 1; i < n; i++ {
		if !isBase58CharFmt(s[i]) {
			return false
		}
	}
	return true
}

func isBtcBech32AddrFmt(s string, n int) bool {
	if n < 42 || n > 62 || len(s) < 4 || s[0] != 'b' || s[1] != 'c' || s[2] != '1' || s[3] != 'q' {
		return false
	}
	for i := 4; i < n; i++ {
		if !((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= '0' && s[i] <= '9')) {
			return false
		}
	}
	return true
}

func isBase58CharFmt(c byte) bool {
	return (c >= '1' && c <= '9') ||
		(c >= 'A' && c <= 'H') ||
		(c >= 'J' && c <= 'N') ||
		(c >= 'P' && c <= 'Z') ||
		(c >= 'a' && c <= 'k') ||
		(c >= 'm' && c <= 'z')
}

func isEmailLocal(local string) bool {
	for i := 0; i < len(local); i++ {
		if !isValidEmailLocalChar(local[i]) {
			return false
		}
	}
	return local[0] != '.' && local[len(local)-1] != '.' && !containsConsecutiveDots(local)
}

func isValidEmailLocalChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
		c == '.' || c == '!' || c == '#' || c == '$' || c == '%' || c == '&' || c == '\'' ||
		c == '*' || c == '+' || c == '/' || c == '=' || c == '?' || c == '^' || c == '_' || c == '`' ||
		c == '{' || c == '|' || c == '}' || c == '~' || c == '-'
}

func containsConsecutiveDots(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '.' && s[i+1] == '.' {
			return true
		}
	}
	return false
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
	uuid = strings.TrimSpace(uuid)
	if len(uuid) != 36 {
		return false
	}
	for i := 0; i < len(uuid); i++ {
		switch i {
		case 8, 13, 18, 23:
			if uuid[i] != '-' {
				return false
			}
		default:
			if !IsHexChar(uuid[i]) {
				return false
			}
		}
	}
	return true
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
	if !isSemver(version) {
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
	if !isBIC(bic) {
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
	if !isDataURI(uri) {
		result.Message = i18n.Msg(MsgFormatDataURIInvalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateBCP47(tag string) CompareResult {
	result := CompareResult{Actual: tag, Expect: "valid BCP 47 language tag"}
	if !isBCP47(tag) {
		result.Message = i18n.Msg(MsgFormatBCP47Invalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateEthAddr(addr string) CompareResult {
	result := CompareResult{Actual: addr, Expect: "valid Ethereum address"}
	if !isEthAddr(addr) {
		result.Message = i18n.Msg(MsgFormatEthAddrInvalid)
		return result
	}
	result.Success = true
	return result
}

func ValidateBtcAddr(addr string) CompareResult {
	result := CompareResult{Actual: addr, Expect: "valid Bitcoin address"}
	if !isBtcAddr(addr) {
		result.Message = i18n.Msg(MsgFormatBtcAddrInvalid)
		return result
	}
	result.Success = true
	return result
}

func IsSemver(version string) bool {
	return isSemver(version)
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
	return isBIC(bic)
}

func IsCron(expr string) bool {
	return isCron(expr)
}

func IsDataURI(uri string) bool {
	return isDataURI(uri)
}

func IsBCP47(tag string) bool {
	return isBCP47(tag)
}

func IsEthAddr(addr string) bool {
	return isEthAddr(addr)
}

func IsBtcAddr(addr string) bool {
	return isBtcAddr(addr)
}

func isISBN10(s string) bool {
	digits := 0
	sum := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' || c == ' ' {
			continue
		}
		digits++
		if digits > 10 {
			return false
		}
		if digits == 10 {
			return isISBN10CheckDigitFmt(c, sum)
		}
		if c < '0' || c > '9' {
			return false
		}
		sum += int(c-'0') * (11 - digits)
	}
	return false
}

func isISBN10CheckDigitFmt(c byte, sum int) bool {
	if c == 'X' || c == 'x' {
		sum += 10
	} else if c >= '0' && c <= '9' {
		sum += int(c - '0')
	} else {
		return false
	}
	return sum%11 == 0
}

func isISBN13(s string) bool {
	digits := 0
	sum := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' || c == ' ' {
			continue
		}
		digits++
		if digits > 13 {
			return false
		}
		if c < '0' || c > '9' {
			return false
		}
		weight := 1
		if digits%2 == 0 {
			weight = 3
		}
		sum += int(c-'0') * weight
	}
	if digits != 13 {
		return false
	}
	return sum%10 == 0 || (10-sum%10)%10 == int(s[len(s)-1]-'0')
}

func isISSN(s string) bool {
	digits := 0
	sum := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' || c == ' ' {
			continue
		}
		digits++
		if digits > 8 {
			return false
		}
		if digits == 8 {
			if c == 'X' || c == 'x' {
				return digits == 8 && sum%11 == 10
			}
			if c < '0' || c > '9' {
				return false
			}
			return digits == 8 && sum%11 == int(c-'0')
		}
		if c < '0' || c > '9' {
			return false
		}
		sum += int(c-'0') * (9 - digits)
	}
	return false
}

func isCron(expr string) bool {
	count := 0
	start := 0
	for i := 0; i <= len(expr); i++ {
		if i == len(expr) || expr[i] == ' ' || expr[i] == '\t' {
			if i > start {
				count++
				if !isValidCronFieldZeroAlloc(expr[start:i]) {
					return false
				}
			}
			start = i + 1
		}
	}
	return count == 5 || count == 6
}

func isValidCronFieldZeroAlloc(field string) bool {
	inRange := false
	for i := 0; i < len(field); i++ {
		if !validateCronChar(field, i, &inRange) {
			return false
		}
	}
	return true
}

func validateCronChar(field string, i int, inRange *bool) bool {
	c := field[i]
	switch {
	case c == ',':
		if *inRange {
			return false
		}
	case c == '/':
		return i > 0 && i < len(field)-1 && validateCronStep(field, i+1)
	case c == '-':
		if *inRange {
			return false
		}
		*inRange = true
	case c == '*':
		return i == 0 || field[i-1] == ',' || field[i-1] == '/'
	case c >= '0' && c <= '9':
	default:
		return false
	}
	return true
}

func validateCronStep(field string, start int) bool {
	for j := start; j < len(field); j++ {
		d := field[j]
		if d == ',' {
			break
		}
		if d != '*' && (d < '0' || d > '9') {
			return false
		}
	}
	return true
}
