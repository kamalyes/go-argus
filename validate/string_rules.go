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
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/kamalyes/go-argus/constants"
)

var timezoneCache sync.Map

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
	RGBRegex      = regexp.MustCompile(
		`^rgb\(\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*\)$`,
	)
	RGBARegex = regexp.MustCompile(
		`^rgba\(\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(0|1|0?\.\d+)\s*\)$`,
	)
	HSLRegex = regexp.MustCompile(
		`^hsl\(\s*(360|3[0-5]\d|[12]?\d?\d)\s*,\s*(100|[1-9]?\d)%\s*,\s*(100|[1-9]?\d)%\s*\)$`,
	)
	HSLARegex = regexp.MustCompile(
		`^hsla\(\s*(360|3[0-5]\d|[12]?\d?\d)\s*,\s*(100|[1-9]?\d)%\s*,\s*(100|[1-9]?\d)%\s*,\s*(0|1|0?\.\d+)\s*\)$`,
	)
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

func StringCompareLength(s, param string, op constants.CmpOp) bool {
	n, ok := ParseFloat(param)
	if !ok {
		return false
	}
	actual := float64(utf8.RuneCountInString(s))
	return CompareOp(actual, n, op)
}

func StringMin(s, param string) bool {
	return StringCompareLength(s, param, constants.CmpGTE)
}

func StringMax(s, param string) bool {
	return StringCompareLength(s, param, constants.CmpLTE)
}

func StringLen(s, param string) bool {
	return StringCompareLength(s, param, constants.CmpEQ)
}

func StringEq(s, param string) bool {
	return s == param
}

func StringEqIgnoreCase(s, param string) bool {
	return strings.EqualFold(s, param)
}

func StringNe(s, param string) bool {
	return s != param
}

func StringNeIgnoreCase(s, param string) bool {
	return !strings.EqualFold(s, param)
}

func StringGt(s, param string) bool {
	return StringCompareLength(s, param, constants.CmpGT)
}

func StringGte(s, param string) bool {
	return StringCompareLength(s, param, constants.CmpGTE)
}

func StringLt(s, param string) bool {
	return StringCompareLength(s, param, constants.CmpLT)
}

func StringLte(s, param string) bool {
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
	return StringMatchRunes(
		s,
		func(r rune) bool { return r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') },
	)
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

func parseRGBLike(s, prefix string, alpha bool) bool {
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

func parseHSLLike(s, prefix string, alpha bool) bool {
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

// parsePrefix 解析 CIDR，返回 (是否合法, 是否属于 IPv4 族)
//
// 用 netip.ParsePrefix 而非 net.ParseCIDR：后者要为 net.IP 与 net.IPNet 分配
// 字节切片，实测 2001:db8::/32 要 70.3ns / 96B / 4 allocs，而前者返回纯值类型，
// 同一输入 21.4ns / 0B / 0 allocs（3.3x）校验只要结论，不要解析产物
//
// IPv4 族的判定对齐原先的 ip.To4() != nil：它对 IPv4 和 IPv4-in-IPv6
// （::ffff:1.2.3.4）都返回非 nil，netip 侧对应 Is4() || Is4In6()
// 两者都拒绝带 zone 的地址（fe80::1%eth0/64）
func parsePrefix(s string) (ok, isV4 bool) {
	p, err := netip.ParsePrefix(s)
	if err != nil {
		// netip 不接受前缀位数带前导零（"2001:db8::/032"），net.ParseCIDR 接受。
		// go-argus 对齐 go-playground/validator 的宽松语义，故去零重试一次；
		// 只在解析已经失败后才走这条路，正常输入不付代价。
		// 该分歧由 TestCIDRMatchesNetParseCIDR 的差分用例抓出
		norm, changed := stripPrefixBitsLeadingZeros(s)
		if !changed {
			return false, false
		}
		if p, err = netip.ParsePrefix(norm); err != nil {
			return false, false
		}
	}
	addr := p.Addr()
	return true, addr.Is4() || addr.Is4In6()
}

// stripPrefixBitsLeadingZeros 去掉 "addr/bits" 中 bits 的前导零
//
// 不含 '/'、bits 为空、bits 无前导零、或 bits 含非数字时原样返回并报告未改动
func stripPrefixBitsLeadingZeros(s string) (string, bool) {
	i := strings.LastIndexByte(s, '/')
	if i < 0 || i == len(s)-1 {
		return s, false
	}
	bits := s[i+1:]
	if bits[0] != '0' || len(bits) == 1 {
		return s, false
	}
	j := 0
	for j < len(bits)-1 && bits[j] == '0' {
		j++
	}
	for k := j; k < len(bits); k++ {
		if bits[k] < '0' || bits[k] > '9' {
			return s, false
		}
	}
	return s[:i+1] + bits[j:], true
}

func StringCIDR(s string) bool {
	s = TrimSpaceIfNeeded(s)
	if parseIPv4CIDRFast(s) {
		return true
	}
	ok, _ := parsePrefix(s)
	return ok
}

func StringCIDRv4(s string) bool {
	s = TrimSpaceIfNeeded(s)
	if parseIPv4CIDRFast(s) {
		return true
	}
	ok, isV4 := parsePrefix(s)
	return ok && isV4
}

func StringCIDRv6(s string) bool {
	s = TrimSpaceIfNeeded(s)
	ok, isV4 := parsePrefix(s)
	return ok && !isV4
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
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '+' && c != '-' && c != '.' {
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
	return isBase32Syntax(ts)
}

// isBase32Syntax 校验 base32 标准编码的语法，不做实际解码
//
// 原实现调 base32.StdEncoding.DecodeString 并丢弃结果，只为取 err：解码缓冲区
// 按输入长度线性分配（16 字符输入即 1 alloc/16B，长串按比例放大），而校验只需要
// 字符集与填充形状。改为就地扫描后与 isBase64Syntax 的写法对齐，分配降为零
//
// 与标准库解码器的等价点：字母表为 A-Z2-7、换行符（\r \n）在任意位置被忽略、
// 去掉换行后的总长须为 8 的倍数、填充符只出现在末尾且个数为 1/3/4/6
// （对应末组 7/5/4/2 个数据字符）。等价性由 TestBase32SyntaxMatchesDecoder 的差分用例守住
func isBase32Syntax(s string) bool {
	if s == "" {
		return false
	}
	n := 0       // 去掉换行后的有效字符数
	padding := 0 // 末尾填充符个数
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\r' || c == '\n' {
			continue
		}
		n++
		if c == '=' {
			padding++
			continue
		}
		// 填充符之后不得再出现数据字符
		if padding > 0 {
			return false
		}
		if (c < 'A' || c > 'Z') && (c < '2' || c > '7') {
			return false
		}
	}
	if n == 0 || n%8 != 0 {
		return false
	}
	switch padding {
	case 0, 1, 3, 4, 6:
		return true
	}
	return false
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

func isBase64Syntax(s string, urlSafe, padded bool) bool {
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

func StringStartsWith(s, param string) bool {
	return strings.HasPrefix(s, param)
}

func StringEndsWith(s, param string) bool {
	return strings.HasSuffix(s, param)
}

func StringStartsNotWith(s, param string) bool {
	return !strings.HasPrefix(s, param)
}

func StringEndsNotWith(s, param string) bool {
	return !strings.HasSuffix(s, param)
}

func StringContains(s, param string) bool {
	return strings.Contains(s, param)
}

func StringContainsAny(s, param string) bool {
	return strings.ContainsAny(s, param)
}

func StringContainsRune(s, param string) bool {
	r, _ := utf8.DecodeRuneInString(param)
	return r != utf8.RuneError && strings.ContainsRune(s, r)
}

func StringExcludes(s, param string) bool {
	return !strings.Contains(s, param)
}

func StringExcludesAll(s, param string) bool {
	return !strings.ContainsAny(s, param)
}

func StringExcludesRune(s, param string) bool {
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

func StringDatetime(s, param string) bool {
	if param == "" {
		param = time.RFC3339
	}
	_, err := time.Parse(param, s)
	return err == nil
}

func StringTimezone(s string) bool {
	if cached, ok := timezoneCache.Load(s); ok {
		return cached.(bool)
	}
	_, err := time.LoadLocation(s)
	valid := err == nil
	timezoneCache.Store(s, valid)
	return valid
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

// StringMongoDB 校验 MongoDB ObjectID，等价于 `^[0-9a-fA-F]{24}$`
//
// 原实现走 MongoIDRegex.MatchString：定长 + 单一字符类的模式交给正则引擎，
// 单次匹配 119ns，而同形状的手写扫描（见 StringAlpha 一族）只要约 10ns
func StringMongoDB(s string) bool {
	if len(s) != 24 {
		return false
	}
	for i := 0; i < 24; i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}

// StringDNSRFC1035Label 校验 RFC 1035 标签，等价于 `len<=63 && ^[a-z]([-a-z0-9]*[a-z0-9])?$`
//
// 原实现走 DNSLabelRegex.MatchString，单次匹配 105ns；手写扫描约 10ns
func StringDNSRFC1035Label(s string) bool {
	if len(s) == 0 || len(s) > 63 {
		return false
	}
	// 首位必须是小写字母
	if s[0] < 'a' || s[0] > 'z' {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' {
			return false
		}
	}
	// 末位不可为连字符：正则的 ([-a-z0-9]*[a-z0-9])? 分组要求结尾落在 [a-z0-9]
	return s[len(s)-1] != '-'
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
type StringRuleFunc func(s, param string) bool

// noParamAdapter 将无参数的字符串校验函数适配为 StringRuleFunc
func noParamAdapter(fn func(string) bool) StringRuleFunc {
	return func(s, _ string) bool { return fn(s) }
}

// StringRuleMap 字符串规则映射表，VarString 快速路径直接查表
var StringRuleMap = map[string]StringRuleFunc{
	constants.RuleRequired:            noParamAdapter(StringRequired),
	constants.RuleIsDefault:           noParamAdapter(StringIsDefault),
	constants.RuleMin:                 StringMin,
	constants.RuleMax:                 StringMax,
	constants.RuleLen:                 StringLen,
	constants.RuleEq:                  StringEq,
	constants.RuleEqIgnoreCase:        StringEqIgnoreCase,
	constants.RuleNe:                  StringNe,
	constants.RuleNeIgnoreCase:        StringNeIgnoreCase,
	constants.RuleGT:                  StringGt,
	constants.RuleGTE:                 StringGte,
	constants.RuleLT:                  StringLt,
	constants.RuleLTE:                 StringLte,
	constants.RuleAlpha:               noParamAdapter(StringAlpha),
	constants.RuleAlphaSpace:          noParamAdapter(StringAlphaSpace),
	constants.RuleAlphanum:            noParamAdapter(StringAlphanum),
	constants.RuleAlphanumSpace:       noParamAdapter(StringAlphanumSpace),
	constants.RuleAlphaUnicode:        noParamAdapter(StringAlphaUnicode),
	constants.RuleAlphanumUnicode:     noParamAdapter(StringAlphanumUnicode),
	constants.RuleASCII:               noParamAdapter(StringASCII),
	constants.RulePrintASCII:          noParamAdapter(StringPrintASCII),
	constants.RuleMultibyte:           noParamAdapter(StringMultibyte),
	constants.RuleHexadecimal:         noParamAdapter(StringHexadecimal),
	constants.RuleHexColor:            noParamAdapter(StringHexColor),
	constants.RuleRGB:                 noParamAdapter(StringRGB),
	constants.RuleRGBA:                noParamAdapter(StringRGBA),
	constants.RuleHSL:                 noParamAdapter(StringHSL),
	constants.RuleHSLA:                noParamAdapter(StringHSLA),
	constants.RuleEmail:               noParamAdapter(IsEmail),
	constants.RuleE164:                noParamAdapter(StringE164),
	constants.RuleIP:                  noParamAdapter(StringIP),
	constants.RuleIPAddr:              noParamAdapter(StringIP),
	constants.RuleIPv4:                noParamAdapter(StringIPv4),
	constants.RuleIPv6:                noParamAdapter(StringIPv6),
	constants.RuleCIDR:                noParamAdapter(StringCIDR),
	constants.RuleCIDRv4:              noParamAdapter(StringCIDRv4),
	constants.RuleCIDRv6:              noParamAdapter(StringCIDRv6),
	constants.RuleMAC:                 noParamAdapter(StringMAC),
	constants.RuleHostname:            noParamAdapter(StringHostname),
	constants.RuleHostnameRFC1123:     noParamAdapter(StringHostname),
	constants.RuleFQDN:                noParamAdapter(StringFQDN),
	constants.RuleHostnamePort:        noParamAdapter(StringHostnamePort),
	constants.RulePort:                noParamAdapter(StringPort),
	constants.RuleURL:                 noParamAdapter(StringURL),
	constants.RuleURI:                 noParamAdapter(StringURI),
	constants.RuleHTTPURL:             noParamAdapter(StringHTTPURL),
	constants.RuleHTTPSURL:            noParamAdapter(StringHTTPSURL),
	constants.RuleURLEncoded:          noParamAdapter(StringURLEncoded),
	constants.RuleHTML:                noParamAdapter(StringHTML),
	constants.RuleHTMLEncoded:         noParamAdapter(StringHTMLEncoded),
	constants.RuleUUID:                noParamAdapter(StringUUID),
	constants.RuleUUID3:               noParamAdapter(StringUUID3),
	constants.RuleUUID4:               noParamAdapter(StringUUID4),
	constants.RuleUUID5:               noParamAdapter(StringUUID5),
	constants.RuleUUIDRFC4122:         noParamAdapter(StringUUID),
	constants.RuleUUID3RFC4122:        noParamAdapter(StringUUID3),
	constants.RuleUUID4RFC4122:        noParamAdapter(StringUUID4),
	constants.RuleUUID5RFC4122:        noParamAdapter(StringUUID5),
	constants.RuleBase32:              noParamAdapter(StringBase32),
	constants.RuleBase64:              noParamAdapter(StringBase64),
	constants.RuleBase64URL:           noParamAdapter(StringBase64URL),
	constants.RuleBase64RawURL:        noParamAdapter(StringBase64RawURL),
	constants.RuleJSON:                noParamAdapter(StringJSON),
	constants.RuleUnique:              noParamAdapter(StringUnique),
	constants.RuleStartsWith:          StringStartsWith,
	constants.RuleEndsWith:            StringEndsWith,
	constants.RuleStartsNotWith:       StringStartsNotWith,
	constants.RuleEndsNotWith:         StringEndsNotWith,
	constants.RuleContains:            StringContains,
	constants.RuleContainsAny:         StringContainsAny,
	constants.RuleContainsRune:        StringContainsRune,
	constants.RuleExcludes:            StringExcludes,
	constants.RuleExcludesAll:         StringExcludesAll,
	constants.RuleExcludesRune:        StringExcludesRune,
	constants.RuleLowercase:           noParamAdapter(StringLowercase),
	constants.RuleUppercase:           noParamAdapter(StringUppercase),
	constants.RuleBoolean:             noParamAdapter(StringBoolean),
	constants.RuleNumber:              noParamAdapter(StringNumber),
	constants.RuleNumeric:             noParamAdapter(StringNumber),
	constants.RuleDatetime:            StringDatetime,
	constants.RuleTimezone:            noParamAdapter(StringTimezone),
	constants.RuleLatitude:            noParamAdapter(StringLatitude),
	constants.RuleLongitude:           noParamAdapter(StringLongitude),
	constants.RuleFile:                noParamAdapter(StringFile),
	constants.RuleFilepath:            noParamAdapter(StringFilePath),
	constants.RuleDir:                 noParamAdapter(StringDir),
	constants.RuleDirpath:             noParamAdapter(StringDirPath),
	constants.RuleMongoDB:             noParamAdapter(StringMongoDB),
	constants.RuleLuhnChecksum:        noParamAdapter(IsLuhnChecksum),
	constants.RuleCreditCard:          noParamAdapter(IsLuhnChecksum),
	constants.RuleDNSRFC1035Label:     noParamAdapter(StringDNSRFC1035Label),
	constants.RuleSemver:              noParamAdapter(IsSemver),
	constants.RuleISBN10:              noParamAdapter(IsISBN10),
	constants.RuleISBN13:              noParamAdapter(IsISBN13),
	constants.RuleISSN:                noParamAdapter(StringISSN),
	constants.RuleBIC:                 noParamAdapter(StringBIC),
	constants.RuleCron:                noParamAdapter(StringCron),
	constants.RuleDataURI:             noParamAdapter(StringDataURI),
	constants.RuleBCP47:               noParamAdapter(StringBCP47),
	constants.RuleEthAddr:             noParamAdapter(StringEthAddr),
	constants.RuleBtcAddr:             noParamAdapter(StringBtcAddr),
	constants.RuleDigits:              noParamAdapter(IsDigits),
	constants.RuleIntOrFloat:          noParamAdapter(IsIntOrFloat),
	constants.RuleWordChars:           noParamAdapter(IsWordChars),
	constants.RuleUpperAlpha:          noParamAdapter(IsUpperAlpha),
	constants.RuleLowerAlpha:          noParamAdapter(IsLowerAlpha),
	constants.RuleStrongPassword:      noParamAdapter(IsStrongPassword),
	constants.RuleHasSpecial:          noParamAdapter(HasSpecialChars),
	constants.RuleHasDoubleByte:       noParamAdapter(HasDoubleByte),
	constants.RuleEmptyLine:           noParamAdapter(IsEmptyLine),
	constants.RuleChinesePhone:        noParamAdapter(IsChinesePhoneNumber),
	constants.RuleChineseIDCard:       noParamAdapter(IsChineseIDCard),
	constants.RuleChineseIDCardStrict: noParamAdapter(IsChineseIDCardWithChecksum),
}

// StringOneOf 判断字符串是否在候选列表中（精确匹配）
func StringOneOf(s string, parts []string) bool {
	return slices.Contains(parts, s)
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
