/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 17:59:17
 * @FilePath: \go-argus\validate\string_rules.go
 * @Description: 字符串规则校验实现，所有逻辑集中于此
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"encoding/base32"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/kamalyes/go-argus/constants"
)

type CmpOp = constants.CmpOp

const (
	CmpEQ  = constants.CmpEQ
	CmpLT  = constants.CmpLT
	CmpLTE = constants.CmpLTE
	CmpGT  = constants.CmpGT
	CmpGTE = constants.CmpGTE
	CmpNE  = constants.CmpNE
)

func CmpOpFromStr(op string) CmpOp {
	return constants.CmpOpFromStr(op)
}

func CompareStringsOp(left, right string, op constants.CmpOp) bool {
	switch op {
	case constants.CmpEQ:
		return left == right
	case constants.CmpNE:
		return left != right
	case constants.CmpGT:
		return left > right
	case constants.CmpGTE:
		return left >= right
	case constants.CmpLT:
		return left < right
	case constants.CmpLTE:
		return left <= right
	default:
		return false
	}
}

func CompareOp(actual, expect float64, op constants.CmpOp) bool {
	switch op {
	case constants.CmpEQ:
		return actual == expect
	case constants.CmpNE:
		return actual != expect
	case constants.CmpLT:
		return actual < expect
	case constants.CmpLTE:
		return actual <= expect
	case constants.CmpGT:
		return actual > expect
	case constants.CmpGTE:
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

func StringCompareLength(s string, param string, op constants.CmpOp) bool {
	n, ok := ParseFloat(param)
	if !ok {
		return false
	}
	actual := float64(utf8.RuneCountInString(s))
	return CompareOp(actual, n, op)
}

func StringMin(s string, param string) bool {
	return StringCompareLength(s, param, constants.CmpGTE)
}

func StringMax(s string, param string) bool {
	return StringCompareLength(s, param, constants.CmpLTE)
}

func StringLen(s string, param string) bool {
	return StringCompareLength(s, param, constants.CmpEQ)
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
	return StringCompareLength(s, param, constants.CmpGT)
}

func StringGte(s string, param string) bool {
	return StringCompareLength(s, param, constants.CmpGTE)
}

func StringLt(s string, param string) bool {
	return StringCompareLength(s, param, constants.CmpLT)
}

func StringLte(s string, param string) bool {
	return StringCompareLength(s, param, constants.CmpLTE)
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
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	switch len(s) {
	case 3, 4, 6, 8:
	default:
		return false
	}
	for i := 0; i < len(s); i++ {
		if !IsHexChar(s[i]) {
			return false
		}
	}
	return true
}

func StringRGB(s string) bool {
	return parseRGBLike(s, "rgb(", false)
}

func StringRGBA(s string) bool {
	return parseRGBLike(s, "rgba(", true)
}

func StringHSL(s string) bool {
	return parseHSLLike(s, "hsl(", false)
}

func StringHSLA(s string) bool {
	return parseHSLLike(s, "hsla(", true)
}

func StringE164(s string) bool {
	if len(s) < 3 || len(s) > 16 || s[0] != '+' || s[1] < '1' || s[1] > '9' {
		return false
	}
	for i := 2; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func parseRGBLike(s string, prefix string, alpha bool) bool {
	if !strings.HasPrefix(s, prefix) || len(s) <= len(prefix) || s[len(s)-1] != ')' {
		return false
	}
	i := len(prefix)
	for n := 0; n < 3; n++ {
		value, ok := parseColorInt(s, &i)
		if !ok || value > 255 {
			return false
		}
		if !consumeColorSep(s, &i, n == 2 && alpha) {
			return false
		}
	}
	if alpha {
		if !parseAlpha(s, &i) {
			return false
		}
		skipASCIIWS(s, &i)
		return i == len(s)-1
	}
	skipASCIIWS(s, &i)
	return i == len(s)-1
}

func parseHSLLike(s string, prefix string, alpha bool) bool {
	if !strings.HasPrefix(s, prefix) || len(s) <= len(prefix) || s[len(s)-1] != ')' {
		return false
	}
	i := len(prefix)
	hue, ok := parseColorInt(s, &i)
	if !ok || hue > 360 || !consumeColorSep(s, &i, false) {
		return false
	}
	for n := 0; n < 2; n++ {
		pct, ok := parseColorInt(s, &i)
		if !ok || pct > 100 {
			return false
		}
		skipASCIIWS(s, &i)
		if i >= len(s) || s[i] != '%' {
			return false
		}
		i++
		if !consumeColorSep(s, &i, n == 1 && alpha) {
			return false
		}
	}
	if alpha {
		if !parseAlpha(s, &i) {
			return false
		}
		skipASCIIWS(s, &i)
		return i == len(s)-1
	}
	skipASCIIWS(s, &i)
	return i == len(s)-1
}

func parseColorInt(s string, i *int) (int, bool) {
	skipASCIIWS(s, i)
	start := *i
	value := 0
	for *i < len(s) && s[*i] >= '0' && s[*i] <= '9' {
		value = value*10 + int(s[*i]-'0')
		*i++
	}
	return value, *i > start
}

func consumeColorSep(s string, i *int, needComma bool) bool {
	skipASCIIWS(s, i)
	if needComma {
		if *i >= len(s) || s[*i] != ',' {
			return false
		}
		*i++
		return true
	}
	if *i < len(s) && s[*i] == ',' {
		*i++
		return true
	}
	return *i < len(s) && s[*i] == ')'
}

func parseAlpha(s string, i *int) bool {
	skipASCIIWS(s, i)
	if *i >= len(s) {
		return false
	}
	if s[*i] == '0' {
		*i++
		if *i < len(s) && s[*i] == '.' {
			*i++
			start := *i
			for *i < len(s) && s[*i] >= '0' && s[*i] <= '9' {
				*i++
			}
			return *i > start
		}
		return true
	}
	if s[*i] == '1' {
		*i++
		return true
	}
	if s[*i] == '.' {
		*i++
		start := *i
		for *i < len(s) && s[*i] >= '0' && s[*i] <= '9' {
			*i++
		}
		return *i > start
	}
	return false
}

func skipASCIIWS(s string, i *int) {
	for *i < len(s) {
		switch s[*i] {
		case ' ', '\t', '\n', '\r':
			*i++
		default:
			return
		}
	}
}

func StringIP(s string) bool {
	s = TrimSpaceIfNeeded(s)
	return parseIPv4Fast(s) || net.ParseIP(s) != nil
}

func StringIPv4(s string) bool {
	s = TrimSpaceIfNeeded(s)
	if parseIPv4Fast(s) {
		return true
	}
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() != nil
}

func StringIPv6(s string) bool {
	s = TrimSpaceIfNeeded(s)
	if s == "::1" || s == "::" {
		return true
	}
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() == nil
}

func StringCIDR(s string) bool {
	s = TrimSpaceIfNeeded(s)
	if parseIPv4CIDRFast(s) {
		return true
	}
	_, _, err := net.ParseCIDR(s)
	return err == nil
}

func StringCIDRv4(s string) bool {
	s = TrimSpaceIfNeeded(s)
	if parseIPv4CIDRFast(s) {
		return true
	}
	ip, _, err := net.ParseCIDR(s)
	return err == nil && ip.To4() != nil
}

func StringCIDRv6(s string) bool {
	s = TrimSpaceIfNeeded(s)
	ip, _, err := net.ParseCIDR(s)
	return err == nil && ip.To4() == nil
}

func StringMAC(s string) bool {
	s = TrimSpaceIfNeeded(s)
	if parseMACFast(s) {
		return true
	}
	_, err := net.ParseMAC(s)
	return err == nil
}

func parseIPv4Fast(s string) bool {
	part := 0
	value := 0
	digits := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			if digits == 0 && c == '0' && i+1 < len(s) && s[i+1] >= '0' && s[i+1] <= '9' {
				return false
			}
			value = value*10 + int(c-'0')
			digits++
			if digits > 3 || value > 255 {
				return false
			}
			continue
		}
		if c != '.' || digits == 0 || part == 3 {
			return false
		}
		part++
		value = 0
		digits = 0
	}
	return part == 3 && digits > 0
}

func parseIPv4CIDRFast(s string) bool {
	slash := strings.IndexByte(s, '/')
	if slash <= 0 || slash == len(s)-1 {
		return false
	}
	if !parseIPv4Fast(s[:slash]) {
		return false
	}
	bits := 0
	for i := slash + 1; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return false
		}
		bits = bits*10 + int(c-'0')
		if bits > 32 {
			return false
		}
	}
	return true
}

func parseMACFast(s string) bool {
	if len(s) == 12 {
		for i := 0; i < len(s); i++ {
			if !IsHexChar(s[i]) {
				return false
			}
		}
		return true
	}
	if len(s) != 17 {
		return false
	}
	sep := s[2]
	if sep != ':' && sep != '-' {
		return false
	}
	for i := 0; i < len(s); i++ {
		if (i+1)%3 == 0 {
			if s[i] != sep {
				return false
			}
			continue
		}
		if !IsHexChar(s[i]) {
			return false
		}
	}
	return true
}

func StringHostname(s string) bool {
	// 内联 TrimSpace：仅当首尾有空白时才分配
	if len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = strings.TrimSpace(s)
	}
	// hostname 允许尾部点号（FQDN 形式），去掉后再验证
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return IsHostname(s)
}

func StringFQDN(s string) bool {
	// 内联 TrimSpace
	if len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = strings.TrimSpace(s)
	}
	n := len(s)
	if n == 0 || s[n-1] != '.' {
		return false
	}
	return IsHostname(s[:n-1])
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
	return isUUIDVersion(s, '3')
}

func StringUUID4(s string) bool {
	return isUUIDVersion(s, '4')
}

func StringUUID5(s string) bool {
	return isUUIDVersion(s, '5')
}

func isUUIDVersion(s string, version byte) bool {
	if len(s) != 36 || s[14] != version {
		return false
	}
	return IsUUID(s)
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
	s = TrimSpaceIfNeeded(s)
	return isBase64Syntax(s, false, true) || isBase64Syntax(s, false, false)
}

func StringBase64URL(s string) bool {
	ts := strings.TrimSpace(s)
	if ts == "" {
		return false
	}
	return isBase64Syntax(ts, true, true)
}

func StringBase64RawURL(s string) bool {
	ts := strings.TrimSpace(s)
	if ts == "" {
		return false
	}
	return isBase64Syntax(ts, true, false)
}

func isBase64Syntax(s string, urlSafe bool, padded bool) bool {
	if s == "" {
		return false
	}
	padding := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '=' {
			padding++
			if !padded || padding > 2 || i < len(s)-2 {
				return false
			}
			continue
		}
		ok := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
		if urlSafe {
			ok = ok || c == '-' || c == '_'
		} else {
			ok = ok || c == '+' || c == '/'
		}
		if padding > 0 || !ok {
			return false
		}
	}
	if padded {
		return len(s)%4 == 0
	}
	return len(s)%4 != 1
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

// StringISSN 验证 ISSN 格式
func StringISSN(s string) bool { return IsISSN(s) }

// StringBIC 验证 BIC/SWIFT 代码
func StringBIC(s string) bool { return IsBIC(s) }

// StringCron 验证 cron 表达式
func StringCron(s string) bool { return IsCron(s) }

// StringDataURI 验证 Data URI 格式
func StringDataURI(s string) bool { return IsDataURI(s) }

// StringBCP47 验证 BCP47 语言标签
func StringBCP47(s string) bool { return IsBCP47(s) }

// StringEthAddr 验证以太坊地址
func StringEthAddr(s string) bool { return IsEthAddr(s) }

// StringBtcAddr 验证比特币地址
func StringBtcAddr(s string) bool { return IsBtcAddr(s) }

// StringRuleFunc 字符串规则函数签名
type StringRuleFunc func(s string, param string) bool

// noParamAdapter 将无参数的字符串校验函数适配为 StringRuleFunc
func noParamAdapter(fn func(string) bool) StringRuleFunc {
	return func(s string, _ string) bool { return fn(s) }
}

// StringRuleMap 字符串规则映射表，VarString 快速路径直接查表
var StringRuleMap = map[string]StringRuleFunc{
	constants.RuleRequired:        noParamAdapter(StringRequired),
	constants.RuleIsDefault:       noParamAdapter(StringIsDefault),
	constants.RuleMin:             StringMin,
	constants.RuleMax:             StringMax,
	constants.RuleLen:             StringLen,
	constants.RuleEq:              StringEq,
	constants.RuleEqIgnoreCase:    StringEqIgnoreCase,
	constants.RuleNe:              StringNe,
	constants.RuleNeIgnoreCase:    StringNeIgnoreCase,
	constants.RuleGT:              StringGt,
	constants.RuleGTE:             StringGte,
	constants.RuleLT:              StringLt,
	constants.RuleLTE:             StringLte,
	constants.RuleAlpha:           noParamAdapter(StringAlpha),
	constants.RuleAlphaSpace:      noParamAdapter(StringAlphaSpace),
	constants.RuleAlphanum:        noParamAdapter(StringAlphanum),
	constants.RuleAlphanumSpace:   noParamAdapter(StringAlphanumSpace),
	constants.RuleAlphaUnicode:    noParamAdapter(StringAlphaUnicode),
	constants.RuleAlphanumUnicode: noParamAdapter(StringAlphanumUnicode),
	constants.RuleASCII:           noParamAdapter(StringASCII),
	constants.RulePrintASCII:      noParamAdapter(StringPrintASCII),
	constants.RuleMultibyte:       noParamAdapter(StringMultibyte),
	constants.RuleHexadecimal:     noParamAdapter(StringHexadecimal),
	constants.RuleHexColor:        noParamAdapter(StringHexColor),
	constants.RuleRGB:             noParamAdapter(StringRGB),
	constants.RuleRGBA:            noParamAdapter(StringRGBA),
	constants.RuleHSL:             noParamAdapter(StringHSL),
	constants.RuleHSLA:            noParamAdapter(StringHSLA),
	constants.RuleEmail:           noParamAdapter(IsEmail),
	constants.RuleE164:            noParamAdapter(StringE164),
	constants.RuleIP:              noParamAdapter(StringIP),
	constants.RuleIPAddr:          noParamAdapter(StringIP),
	constants.RuleIPv4:            noParamAdapter(StringIPv4),
	constants.RuleIPv6:            noParamAdapter(StringIPv6),
	constants.RuleCIDR:            noParamAdapter(StringCIDR),
	constants.RuleCIDRv4:          noParamAdapter(StringCIDRv4),
	constants.RuleCIDRv6:          noParamAdapter(StringCIDRv6),
	constants.RuleMAC:             noParamAdapter(StringMAC),
	constants.RuleHostname:        noParamAdapter(StringHostname),
	constants.RuleHostnameRFC1123: noParamAdapter(StringHostname),
	constants.RuleFQDN:            noParamAdapter(StringFQDN),
	constants.RuleHostnamePort:    noParamAdapter(StringHostnamePort),
	constants.RulePort:            noParamAdapter(StringPort),
	constants.RuleURL:             noParamAdapter(StringURL),
	constants.RuleURI:             noParamAdapter(StringURI),
	constants.RuleHTTPURL:         noParamAdapter(StringHTTPURL),
	constants.RuleHTTPSURL:        noParamAdapter(StringHTTPSURL),
	constants.RuleURLEncoded:      noParamAdapter(StringURLEncoded),
	constants.RuleHTML:            noParamAdapter(StringHTML),
	constants.RuleHTMLEncoded:     noParamAdapter(StringHTMLEncoded),
	constants.RuleUUID:            noParamAdapter(StringUUID),
	constants.RuleUUID3:           noParamAdapter(StringUUID3),
	constants.RuleUUID4:           noParamAdapter(StringUUID4),
	constants.RuleUUID5:           noParamAdapter(StringUUID5),
	constants.RuleUUIDRFC4122:     noParamAdapter(StringUUID),
	constants.RuleUUID3RFC4122:    noParamAdapter(StringUUID3),
	constants.RuleUUID4RFC4122:    noParamAdapter(StringUUID4),
	constants.RuleUUID5RFC4122:    noParamAdapter(StringUUID5),
	constants.RuleBase32:          noParamAdapter(StringBase32),
	constants.RuleBase64:          noParamAdapter(StringBase64),
	constants.RuleBase64URL:       noParamAdapter(StringBase64URL),
	constants.RuleBase64RawURL:    noParamAdapter(StringBase64RawURL),
	constants.RuleJSON:            noParamAdapter(StringJSON),
	constants.RuleUnique:          noParamAdapter(StringUnique),
	constants.RuleStartsWith:      StringStartsWith,
	constants.RuleEndsWith:        StringEndsWith,
	constants.RuleStartsNotWith:   StringStartsNotWith,
	constants.RuleEndsNotWith:     StringEndsNotWith,
	constants.RuleContains:        StringContains,
	constants.RuleContainsAny:     StringContainsAny,
	constants.RuleContainsRune:    StringContainsRune,
	constants.RuleExcludes:        StringExcludes,
	constants.RuleExcludesAll:     StringExcludesAll,
	constants.RuleExcludesRune:    StringExcludesRune,
	constants.RuleLowercase:       noParamAdapter(StringLowercase),
	constants.RuleUppercase:       noParamAdapter(StringUppercase),
	constants.RuleBoolean:         noParamAdapter(StringBoolean),
	constants.RuleNumber:          noParamAdapter(StringNumber),
	constants.RuleNumeric:         noParamAdapter(StringNumber),
	constants.RuleDatetime:        StringDatetime,
	constants.RuleTimezone:        noParamAdapter(StringTimezone),
	constants.RuleLatitude:        noParamAdapter(StringLatitude),
	constants.RuleLongitude:       noParamAdapter(StringLongitude),
	constants.RuleFile:            noParamAdapter(StringFile),
	constants.RuleFilepath:        noParamAdapter(StringFilePath),
	constants.RuleDir:             noParamAdapter(StringDir),
	constants.RuleDirpath:         noParamAdapter(StringDirPath),
	constants.RuleMongoDB:         noParamAdapter(StringMongoDB),
	constants.RuleLuhnChecksum:    noParamAdapter(IsLuhnChecksum),
	constants.RuleCreditCard:      noParamAdapter(IsLuhnChecksum),
	constants.RuleDNSRFC1035Label: noParamAdapter(StringDNSRFC1035Label),
	constants.RuleSemver:          noParamAdapter(IsSemver),
	constants.RuleISBN10:          noParamAdapter(IsISBN10),
	constants.RuleISBN13:          noParamAdapter(IsISBN13),
	constants.RuleISSN:            noParamAdapter(StringISSN),
	constants.RuleBIC:             noParamAdapter(StringBIC),
	constants.RuleCron:            noParamAdapter(StringCron),
	constants.RuleDataURI:         noParamAdapter(StringDataURI),
	constants.RuleBCP47:           noParamAdapter(StringBCP47),
	constants.RuleEthAddr:         noParamAdapter(StringEthAddr),
	constants.RuleBtcAddr:         noParamAdapter(StringBtcAddr),
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
