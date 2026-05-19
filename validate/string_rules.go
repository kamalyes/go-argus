/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-19 00:00:00
 * @FilePath: \go-argus\validate\string_rules.go
 * @Description: 字符串规则校验实现，所有逻辑集中于此
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"encoding/base32"
	"encoding/base64"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// CmpOp 比较运算符类型
type CmpOp int

const (
	CmpEQ CmpOp = iota
	CmpLT
	CmpLTE
	CmpGT
	CmpGTE
	CmpNE
)

// CompareOp 比较运算符实现
func CompareOp(actual, expect float64, op CmpOp) bool {
	switch op {
	case CmpEQ:
		return actual == expect
	case CmpNE:
		return actual != expect
	case CmpLT:
		return actual < expect
	case CmpLTE:
		return actual <= expect
	case CmpGT:
		return actual > expect
	case CmpGTE:
		return actual >= expect
	default:
		return false
	}
}

var (
	ColorHexRegex = regexp.MustCompile(`^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	RGBRegex      = regexp.MustCompile(`^rgb\(\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*\)$`)
	RGBARegex     = regexp.MustCompile(`^rgba\(\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(0|1|0?\.\d+)\s*\)$`)
	HSLRegex      = regexp.MustCompile(`^hsl\(\s*(360|3[0-5]\d|[12]?\d?\d)\s*,\s*(100|[1-9]?\d)%\s*,\s*(100|[1-9]?\d)%\s*\)$`)
	HSLARegex     = regexp.MustCompile(`^hsla\(\s*(360|3[0-5]\d|[12]?\d?\d)\s*,\s*(100|[1-9]?\d)%\s*,\s*(100|[1-9]?\d)%\s*,\s*(0|1|0?\.\d+)\s*\)$`)
	E164Regex     = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	MongoIDRegex  = regexp.MustCompile(`^[0-9a-fA-F]{24}$`)
	DNSLabelRegex = regexp.MustCompile(`^[a-z]([-a-z0-9]*[a-z0-9])?$`)
)

// IsBlankString 检查字符串是否为空
func IsBlankString(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\t' && s[i] != '\n' && s[i] != '\r' {
			return false
		}
	}
	return true
}

// StringRequired 检查字符串是否为空
func StringRequired(s string) bool {
	return !IsBlankString(s)
}

func StringIsDefault(s string) bool {
	return IsBlankString(s)
}

func StringCompareLength(s string, param string, op CmpOp) bool {
	n, ok := ParseFloat(param)
	if !ok {
		return false
	}
	actual := float64(utf8.RuneCountInString(s))
	return CompareOp(actual, n, op)
}

func StringMin(s string, param string) bool {
	return StringCompareLength(s, param, CmpGTE)
}

func StringMax(s string, param string) bool {
	return StringCompareLength(s, param, CmpLTE)
}

func StringLen(s string, param string) bool {
	return StringCompareLength(s, param, CmpEQ)
}

func StringEq(s string, param string) bool {
	return s == param
}

func StringEqIgnoreCase(s string, param string) bool {
	return strings.EqualFold(s, param)
}

func StringNe(s string, param string) bool {
	return s != param
}

func StringNeIgnoreCase(s string, param string) bool {
	return !strings.EqualFold(s, param)
}

func StringGt(s string, param string) bool {
	return StringCompareLength(s, param, CmpGT)
}

func StringGte(s string, param string) bool {
	return StringCompareLength(s, param, CmpGTE)
}

func StringLt(s string, param string) bool {
	return StringCompareLength(s, param, CmpLT)
}

func StringLte(s string, param string) bool {
	return StringCompareLength(s, param, CmpLTE)
}

func StringMatchRunes(s string, fn func(rune) bool) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !fn(r) {
			return false
		}
	}
	return true
}

func StringAlpha(s string) bool {
	return StringMatchRunes(s, func(r rune) bool { return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') })
}

func StringAlphaSpace(s string) bool {
	return StringMatchRunes(s, func(r rune) bool { return r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') })
}

func StringAlphanum(s string) bool {
	return StringMatchRunes(s, func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	})
}

func StringAlphanumSpace(s string) bool {
	return StringMatchRunes(s, func(r rune) bool {
		return r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	})
}

func StringAlphaUnicode(s string) bool {
	return StringMatchRunes(s, unicode.IsLetter)
}

func StringAlphanumUnicode(s string) bool {
	return StringMatchRunes(s, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) })
}

func StringASCII(s string) bool {
	return StringMatchRunes(s, func(r rune) bool { return r <= unicode.MaxASCII })
}

func StringPrintASCII(s string) bool {
	return StringMatchRunes(s, func(r rune) bool { return r >= 0x20 && r <= 0x7e })
}

func StringMultibyte(s string) bool {
	return len(s) != utf8.RuneCountInString(s)
}

func StringHexadecimal(s string) bool {
	return StringMatchRunes(s, func(r rune) bool {
		return (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') || (r >= '0' && r <= '9')
	})
}

func StringHexColor(s string) bool {
	return ColorHexRegex.MatchString(s)
}

func StringRGB(s string) bool {
	return RGBRegex.MatchString(s)
}

func StringRGBA(s string) bool {
	return RGBARegex.MatchString(s)
}

func StringHSL(s string) bool {
	return HSLRegex.MatchString(s)
}

func StringHSLA(s string) bool {
	return HSLARegex.MatchString(s)
}

func StringE164(s string) bool {
	return E164Regex.MatchString(s)
}

func StringIP(s string) bool {
	return net.ParseIP(TrimSpaceIfNeeded(s)) != nil
}

func StringIPv4(s string) bool {
	ip := net.ParseIP(TrimSpaceIfNeeded(s))
	return ip != nil && ip.To4() != nil
}

func StringIPv6(s string) bool {
	ip := net.ParseIP(TrimSpaceIfNeeded(s))
	return ip != nil && ip.To4() == nil
}

func StringCIDR(s string) bool {
	_, _, err := net.ParseCIDR(strings.TrimSpace(s))
	return err == nil
}

func StringCIDRv4(s string) bool {
	ip, _, err := net.ParseCIDR(strings.TrimSpace(s))
	return err == nil && ip.To4() != nil
}

func StringCIDRv6(s string) bool {
	ip, _, err := net.ParseCIDR(strings.TrimSpace(s))
	return err == nil && ip.To4() == nil
}

func StringMAC(s string) bool {
	_, err := net.ParseMAC(strings.TrimSpace(s))
	return err == nil
}

func StringHostname(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return IsHostname(s)
}

func StringFQDN(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0 && s[len(s)-1] == '.' && IsHostname(s[:len(s)-1])
}

func StringHostnamePort(s string) bool {
	host, port, err := net.SplitHostPort(s)
	if err != nil || host == "" {
		return false
	}
	return StringPort(port)
}

func StringPort(s string) bool {
	n, err := strconv.Atoi(s)
	return err == nil && n >= 0 && n <= 65535
}

func StringURL(s string) bool {
	s = strings.TrimSpace(s)
	return HasSchemeAndHost(s)
}

func StringURI(s string) bool {
	s = strings.TrimSpace(s)
	colon := strings.Index(s, ":")
	if colon < 1 {
		return false
	}
	for i := 0; i < colon; i++ {
		c := s[i]
		if !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') && c != '+' && c != '-' && c != '.' {
			return false
		}
	}
	return len(s) > colon+1
}

func StringHTTPURL(s string) bool {
	s = TrimSpaceIfNeeded(s)
	if len(s) < 7 {
		return false
	}
	var colon int
	if s[0] == 'h' && s[1] == 't' && s[2] == 't' && s[3] == 'p' {
		if s[4] == ':' {
			colon = 4
		} else if s[4] == 's' && s[5] == ':' {
			colon = 5
		} else {
			return false
		}
	} else {
		return false
	}
	return HasHostAfterScheme(s, colon)
}

func StringHTTPSURL(s string) bool {
	s = TrimSpaceIfNeeded(s)
	if len(s) < 8 || s[0] != 'h' || s[1] != 't' || s[2] != 't' || s[3] != 'p' || s[4] != 's' || s[5] != ':' {
		return false
	}
	return HasHostAfterScheme(s, 5)
}

func IsHexChar(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func StringURLEncoded(s string) bool {
	hasPercent := false
	for i := 0; i < len(s); i++ {
		if s[i] == '%' {
			hasPercent = true
			if i+2 >= len(s) {
				return false
			}
			if !IsHexChar(s[i+1]) || !IsHexChar(s[i+2]) {
				return false
			}
			i += 2
		}
	}
	return hasPercent
}

func StringHTML(s string) bool {
	return strings.Contains(s, "<") && strings.Contains(s, ">")
}

func StringHTMLEncoded(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '&' {
			if i+2 < len(s) {
				if s[i+1] == '#' {
					return true
				}
				if s[i+1] == 'a' && i+3 < len(s) && s[i+2] == 'm' && s[i+3] == 'p' {
					return true
				}
				if s[i+1] == 'l' && i+3 < len(s) && s[i+2] == 't' && s[i+3] == ';' {
					return true
				}
				if s[i+1] == 'g' && i+3 < len(s) && s[i+2] == 't' && s[i+3] == ';' {
					return true
				}
				if s[i+1] == 'q' && i+4 < len(s) && s[i+2] == 'u' && s[i+3] == 'o' && s[i+4] == 't' {
					return true
				}
				if s[i+1] == 'a' && i+3 < len(s) && s[i+2] == 'p' && s[i+3] == 'o' && i+4 < len(s) && s[i+4] == 's' {
					return true
				}
				if s[i+1] == 'n' && i+4 < len(s) && s[i+2] == 'b' && s[i+3] == 's' && s[i+4] == 'p' {
					return true
				}
			}
		}
	}
	return false
}

func StringUUID(s string) bool {
	return IsUUID(s)
}

func StringUUID3(s string) bool {
	return len(s) == 36 && s[14] == '3' && IsUUID(s)
}

func StringUUID4(s string) bool {
	return len(s) == 36 && s[14] == '4' && IsUUID(s)
}

func StringUUID5(s string) bool {
	return len(s) == 36 && s[14] == '5' && IsUUID(s)
}

func StringBase32(s string) bool {
	ts := strings.TrimSpace(s)
	if ts == "" {
		return false
	}
	_, err := base32.StdEncoding.DecodeString(ts)
	return err == nil
}

func StringBase64(s string) bool {
	return IsBase64(s)
}

func StringBase64URL(s string) bool {
	ts := strings.TrimSpace(s)
	if ts == "" {
		return false
	}
	_, err := base64.URLEncoding.DecodeString(ts)
	return err == nil
}

func StringBase64RawURL(s string) bool {
	ts := strings.TrimSpace(s)
	if ts == "" {
		return false
	}
	_, err := base64.RawURLEncoding.DecodeString(ts)
	return err == nil
}

func StringJSON(s string) bool {
	return IsValidJSONString(s)
}

func IsValidJSONString(s string) bool {
	n := len(s)
	if n == 0 {
		return false
	}
	p := 0
	p = skipWS(s, p)
	if p >= n {
		return false
	}
	p = scanJSONVal(s, p)
	if p < 0 {
		return false
	}
	p = skipWS(s, p)
	return p == n
}

func skipWS(s string, i int) int {
	for i < len(s) {
		switch s[i] {
		case ' ', '\t', '\r', '\n':
			i++
		default:
			return i
		}
	}
	return i
}

func scanJSONVal(s string, i int) int {
	if i >= len(s) {
		return -1
	}
	switch s[i] {
	case '"':
		return scanJSONStr(s, i+1)
	case '{':
		return scanJSONObject(s, i+1)
	case '[':
		return scanJSONArray(s, i+1)
	case 't':
		if i+4 <= len(s) && s[i:i+4] == "true" {
			return i + 4
		}
		return -1
	case 'f':
		if i+5 <= len(s) && s[i:i+5] == "false" {
			return i + 5
		}
		return -1
	case 'n':
		if i+4 <= len(s) && s[i:i+4] == "null" {
			return i + 4
		}
		return -1
	default:
		if s[i] == '-' || (s[i] >= '0' && s[i] <= '9') {
			return scanJSONNum(s, i)
		}
		return -1
	}
}

func scanJSONStr(s string, i int) int {
	for i < len(s) {
		c := s[i]
		if c == '"' {
			return i + 1
		}
		if c == '\\' {
			i += 2
			continue
		}
		if c < 0x20 {
			return -1
		}
		i++
	}
	return -1
}

func scanJSONObject(s string, i int) int {
	i = skipWS(s, i)
	if i >= len(s) {
		return -1
	}
	if s[i] == '}' {
		return i + 1
	}
	for {
		i = skipWS(s, i)
		if i >= len(s) || s[i] != '"' {
			return -1
		}
		i = scanJSONStr(s, i+1)
		if i < 0 {
			return -1
		}
		i = skipWS(s, i)
		if i >= len(s) || s[i] != ':' {
			return -1
		}
		i = skipWS(s, i+1)
		i = scanJSONVal(s, i)
		if i < 0 {
			return -1
		}
		i = skipWS(s, i)
		if i >= len(s) {
			return -1
		}
		if s[i] == '}' {
			return i + 1
		}
		if s[i] != ',' {
			return -1
		}
		i++
	}
}

func scanJSONArray(s string, i int) int {
	i = skipWS(s, i)
	if i >= len(s) {
		return -1
	}
	if s[i] == ']' {
		return i + 1
	}
	for {
		i = skipWS(s, i)
		i = scanJSONVal(s, i)
		if i < 0 {
			return -1
		}
		i = skipWS(s, i)
		if i >= len(s) {
			return -1
		}
		if s[i] == ']' {
			return i + 1
		}
		if s[i] != ',' {
			return -1
		}
		i++
	}
}

func scanJSONNum(s string, i int) int {
	if i < len(s) && s[i] == '-' {
		i++
	}
	if i >= len(s) {
		return -1
	}
	if s[i] == '0' {
		i++
	} else if s[i] >= '1' && s[i] <= '9' {
		i++
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
	} else {
		return -1
	}
	if i < len(s) && s[i] == '.' {
		i++
		if i >= len(s) || s[i] < '0' || s[i] > '9' {
			return -1
		}
		i++
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
	}
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			i++
		}
		if i >= len(s) || s[i] < '0' || s[i] > '9' {
			return -1
		}
		i++
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
	}
	return i
}

func StringUnique(s string) bool {
	seen := make(map[rune]struct{}, utf8.RuneCountInString(s))
	for _, r := range s {
		if _, ok := seen[r]; ok {
			return false
		}
		seen[r] = struct{}{}
	}
	return true
}

func StringStartsWith(s string, param string) bool {
	return strings.HasPrefix(s, param)
}

func StringEndsWith(s string, param string) bool {
	return strings.HasSuffix(s, param)
}

func StringStartsNotWith(s string, param string) bool {
	return !strings.HasPrefix(s, param)
}

func StringEndsNotWith(s string, param string) bool {
	return !strings.HasSuffix(s, param)
}

func StringContains(s string, param string) bool {
	return strings.Contains(s, param)
}

func StringContainsAny(s string, param string) bool {
	return strings.ContainsAny(s, param)
}

func StringContainsRune(s string, param string) bool {
	r, _ := utf8.DecodeRuneInString(param)
	return r != utf8.RuneError && strings.ContainsRune(s, r)
}

func StringExcludes(s string, param string) bool {
	return !strings.Contains(s, param)
}

func StringExcludesAll(s string, param string) bool {
	return !strings.ContainsAny(s, param)
}

func StringExcludesRune(s string, param string) bool {
	r, _ := utf8.DecodeRuneInString(param)
	return r != utf8.RuneError && !strings.ContainsRune(s, r)
}

func StringLowercase(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			return false
		}
	}
	return true
}

func StringUppercase(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			return false
		}
	}
	return true
}

func StringBoolean(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

func StringNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func StringDatetime(s string, param string) bool {
	if param == "" {
		param = time.RFC3339
	}
	_, err := time.Parse(param, s)
	return err == nil
}

func StringTimezone(s string) bool {
	_, err := time.LoadLocation(s)
	return err == nil
}

func StringLatitude(s string) bool {
	n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return err == nil && n >= -90 && n <= 90
}

func StringLongitude(s string) bool {
	n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return err == nil && n >= -180 && n <= 180
}

func StringFile(s string) bool {
	info, err := os.Stat(s)
	return err == nil && !info.IsDir()
}

func StringFilePath(s string) bool {
	return strings.TrimSpace(s) != "" && filepath.Clean(s) != "."
}

func StringDir(s string) bool {
	info, err := os.Stat(s)
	return err == nil && info.IsDir()
}

func StringDirPath(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	cleaned := filepath.Clean(s)
	return cleaned != "." && !strings.Contains(filepath.Base(cleaned), ".")
}

func StringMongoDB(s string) bool {
	return MongoIDRegex.MatchString(s)
}

func StringDNSRFC1035Label(s string) bool {
	return len(s) <= 63 && DNSLabelRegex.MatchString(s)
}

// StringOneOf 判断字符串是否在候选列表中（精确匹配）
func StringOneOf(s string, parts []string) bool {
	for _, item := range parts {
		if s == item {
			return true
		}
	}
	return false
}

// StringOneOfCI 判断字符串是否在候选列表中（忽略大小写）
func StringOneOfCI(s string, parts []string) bool {
	for _, item := range parts {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
