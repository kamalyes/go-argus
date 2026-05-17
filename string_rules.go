/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-17 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 15:57:07
 * @FilePath: \go-argus\string_rules.go
 * @Description: Argus 字符串规则实现
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"html"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/kamalyes/go-argus/validate"
)

type stringRuleFunc func(s string, param string) bool

var stringRuleMap = map[string]stringRuleFunc{
	"required":          stringRequired,
	"isdefault":         stringIsDefault,
	"min":               stringMin,
	"max":               stringMax,
	"len":               stringLen,
	"eq":                stringEq,
	"eq_ignore_case":    stringEqIgnoreCase,
	"ne":                stringNe,
	"ne_ignore_case":    stringNeIgnoreCase,
	"gt":                stringGt,
	"gte":               stringGte,
	"lt":                stringLt,
	"lte":               stringLte,
	"alpha":             stringAlpha,
	"alphaspace":        stringAlphaSpace,
	"alphanum":          stringAlphanum,
	"alphanumspace":     stringAlphanumSpace,
	"alphaunicode":      stringAlphaUnicode,
	"alphanumunicode":   stringAlphanumUnicode,
	"ascii":             stringASCII,
	"printascii":        stringPrintASCII,
	"multibyte":         stringMultibyte,
	"hexadecimal":       stringHexadecimal,
	"hexcolor":          stringHexColor,
	"rgb":               stringRGB,
	"rgba":              stringRGBA,
	"hsl":               stringHSL,
	"hsla":              stringHSLA,
	"email":             stringEmail,
	"e164":              stringE164,
	"ip":                stringIP,
	"ip_addr":           stringIP,
	"ipv4":              stringIPv4,
	"ipv6":              stringIPv6,
	"cidr":              stringCIDR,
	"cidrv4":            stringCIDRv4,
	"cidrv6":            stringCIDRv6,
	"mac":               stringMAC,
	"hostname":          stringHostname,
	"hostname_rfc1123":  stringHostname,
	"fqdn":              stringFQDN,
	"hostname_port":     stringHostnamePort,
	"port":              stringPort,
	"url":               stringURL,
	"uri":               stringURI,
	"http_url":          stringHTTPURL,
	"https_url":         stringHTTPSURL,
	"url_encoded":       stringURLEncoded,
	"html":              stringHTML,
	"html_encoded":      stringHTMLEncoded,
	"uuid":              stringUUID,
	"uuid3":             stringUUID3,
	"uuid4":             stringUUID4,
	"uuid5":             stringUUID5,
	"uuid_rfc4122":      stringUUID,
	"uuid3_rfc4122":     stringUUID3,
	"uuid4_rfc4122":     stringUUID4,
	"uuid5_rfc4122":     stringUUID5,
	"base32":            stringBase32,
	"base64":            stringBase64,
	"base64url":         stringBase64URL,
	"base64rawurl":      stringBase64RawURL,
	"json":              stringJSON,
	"unique":            stringUnique,
	"startswith":        stringStartsWith,
	"endswith":          stringEndsWith,
	"startsnotwith":     stringStartsNotWith,
	"endsnotwith":       stringEndsNotWith,
	"contains":          stringContains,
	"containsany":       stringContainsAny,
	"containsrune":      stringContainsRune,
	"excludes":          stringExcludes,
	"excludesall":       stringExcludesAll,
	"excludesrune":      stringExcludesRune,
	"lowercase":         stringLowercase,
	"uppercase":         stringUppercase,
	"boolean":           stringBoolean,
	"number":            stringNumber,
	"numeric":           stringNumber,
	"datetime":          stringDatetime,
	"timezone":          stringTimezone,
	"latitude":          stringLatitude,
	"longitude":         stringLongitude,
	"file":              stringFile,
	"filepath":          stringFilePath,
	"dir":               stringDir,
	"dirpath":           stringDirPath,
	"mongodb":           stringMongoDB,
	"luhn_checksum":     stringLuhnChecksum,
	"credit_card":       stringLuhnChecksum,
	"dns_rfc1035_label": stringDNSRFC1035Label,
	"semver":            stringSemver,
	"isbn10":            stringISBN10,
	"isbn13":            stringISBN13,
	"issn":              stringISSN,
	"bic":               stringBIC,
	"cron":              stringCron,
	"datauri":           stringDataURI,
	"bcp47":             stringBCP47,
	"eth_addr":          stringEthAddr,
	"btc_addr":          stringBtcAddr,
}

func isBlankString(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\t' && s[i] != '\n' && s[i] != '\r' {
			return false
		}
	}
	return true
}

func stringRequired(s string, _ string) bool {
	return !isBlankString(s)
}

func stringIsDefault(s string, _ string) bool {
	return isBlankString(s)
}

func stringCompareLength(s string, param string, op cmpOp) bool {
	n, ok := parseFloat(param)
	if !ok {
		return false
	}
	actual := float64(utf8.RuneCountInString(s))
	return compareOp(actual, n, op)
}

func stringMin(s string, param string) bool {
	return stringCompareLength(s, param, cmpGTE)
}

func stringMax(s string, param string) bool {
	return stringCompareLength(s, param, cmpLTE)
}

func stringLen(s string, param string) bool {
	return stringCompareLength(s, param, cmpEQ)
}

func stringEq(s string, param string) bool {
	return s == param
}

func stringEqIgnoreCase(s string, param string) bool {
	return strings.EqualFold(s, param)
}

func stringNe(s string, param string) bool {
	return s != param
}

func stringNeIgnoreCase(s string, param string) bool {
	return !strings.EqualFold(s, param)
}

func stringGt(s string, param string) bool {
	return stringCompareLength(s, param, cmpGT)
}

func stringGte(s string, param string) bool {
	return stringCompareLength(s, param, cmpGTE)
}

func stringLt(s string, param string) bool {
	return stringCompareLength(s, param, cmpLT)
}

func stringLte(s string, param string) bool {
	return stringCompareLength(s, param, cmpLTE)
}

func stringMatchRunes(s string, fn func(rune) bool) bool {
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

func stringAlpha(s string, _ string) bool {
	return stringMatchRunes(s, func(r rune) bool { return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') })
}

func stringAlphaSpace(s string, _ string) bool {
	return stringMatchRunes(s, func(r rune) bool { return r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') })
}

func stringAlphanum(s string, _ string) bool {
	return stringMatchRunes(s, func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	})
}

func stringAlphanumSpace(s string, _ string) bool {
	return stringMatchRunes(s, func(r rune) bool {
		return r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	})
}

func stringAlphaUnicode(s string, _ string) bool {
	return stringMatchRunes(s, unicode.IsLetter)
}

func stringAlphanumUnicode(s string, _ string) bool {
	return stringMatchRunes(s, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) })
}

func stringASCII(s string, _ string) bool {
	return stringMatchRunes(s, func(r rune) bool { return r <= unicode.MaxASCII })
}

func stringPrintASCII(s string, _ string) bool {
	return stringMatchRunes(s, func(r rune) bool { return r >= 0x20 && r <= 0x7e })
}

func stringMultibyte(s string, _ string) bool {
	return len(s) != utf8.RuneCountInString(s)
}

func stringHexadecimal(s string, _ string) bool {
	return stringMatchRunes(s, func(r rune) bool {
		return (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') || (r >= '0' && r <= '9')
	})
}

func stringHexColor(s string, _ string) bool {
	return colorHexRegex.MatchString(s)
}

func stringRGB(s string, _ string) bool {
	return rgbRegex.MatchString(s)
}

func stringRGBA(s string, _ string) bool {
	return rgbaRegex.MatchString(s)
}

func stringHSL(s string, _ string) bool {
	return hslRegex.MatchString(s)
}

func stringHSLA(s string, _ string) bool {
	return hslaRegex.MatchString(s)
}

func stringEmail(s string, _ string) bool {
	return validate.IsEmail(s)
}

func stringE164(s string, _ string) bool {
	return e164Regex.MatchString(s)
}

func stringIP(s string, _ string) bool {
	return net.ParseIP(trimSpaceIfNeeded(s)) != nil
}

func stringIPv4(s string, _ string) bool {
	ip := net.ParseIP(trimSpaceIfNeeded(s))
	return ip != nil && ip.To4() != nil
}

func stringIPv6(s string, _ string) bool {
	ip := net.ParseIP(trimSpaceIfNeeded(s))
	return ip != nil && ip.To4() == nil
}

func stringCIDR(s string, _ string) bool {
	_, _, err := net.ParseCIDR(strings.TrimSpace(s))
	return err == nil
}

func stringCIDRv4(s string, _ string) bool {
	ip, _, err := net.ParseCIDR(strings.TrimSpace(s))
	return err == nil && ip.To4() != nil
}

func stringCIDRv6(s string, _ string) bool {
	ip, _, err := net.ParseCIDR(strings.TrimSpace(s))
	return err == nil && ip.To4() == nil
}

func stringMAC(s string, _ string) bool {
	_, err := net.ParseMAC(strings.TrimSpace(s))
	return err == nil
}

func stringHostname(s string, _ string) bool {
	return isHostname(strings.TrimSuffix(strings.TrimSpace(s), "."))
}

func stringFQDN(s string, _ string) bool {
	ts := strings.TrimSpace(s)
	return strings.HasSuffix(ts, ".") && isHostname(strings.TrimSuffix(ts, "."))
}

func stringHostnamePort(s string, _ string) bool {
	host, port, err := net.SplitHostPort(s)
	if err != nil || host == "" {
		return false
	}
	return stringPort(port, "")
}

func stringPort(s string, _ string) bool {
	n, err := strconv.Atoi(s)
	return err == nil && n >= 0 && n <= 65535
}

func stringURL(s string, _ string) bool {
	s = strings.TrimSpace(s)
	return hasSchemeAndHost(s)
}

func stringURI(s string, _ string) bool {
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

func stringHTTPURL(s string, _ string) bool {
	s = strings.TrimSpace(s)
	colon := strings.Index(s, ":")
	if colon < 0 {
		return false
	}
	scheme := s[:colon]
	if scheme != "http" && scheme != "https" {
		return false
	}
	return hasHostAfterScheme(s, colon)
}

func stringHTTPSURL(s string, _ string) bool {
	s = strings.TrimSpace(s)
	colon := strings.Index(s, ":")
	if colon < 0 {
		return false
	}
	if s[:colon] != "https" {
		return false
	}
	return hasHostAfterScheme(s, colon)
}

func stringURLEncoded(s string, _ string) bool {
	if !strings.Contains(s, "%") {
		return false
	}
	_, err := url.QueryUnescape(s)
	return err == nil
}

func stringHTML(s string, _ string) bool {
	return strings.Contains(s, "<") && strings.Contains(s, ">")
}

func stringHTMLEncoded(s string, _ string) bool {
	return html.UnescapeString(s) != s
}

func stringUUID(s string, _ string) bool {
	return validate.IsUUID(s)
}

func stringUUID3(s string, _ string) bool {
	return len(s) == 36 && s[14] == '3' && validate.IsUUID(s)
}

func stringUUID4(s string, _ string) bool {
	return len(s) == 36 && s[14] == '4' && validate.IsUUID(s)
}

func stringUUID5(s string, _ string) bool {
	return len(s) == 36 && s[14] == '5' && validate.IsUUID(s)
}

func stringBase32(s string, _ string) bool {
	ts := strings.TrimSpace(s)
	if ts == "" {
		return false
	}
	_, err := base32.StdEncoding.DecodeString(ts)
	return err == nil
}

func stringBase64(s string, _ string) bool {
	return validate.IsBase64(s)
}

func stringBase64URL(s string, _ string) bool {
	ts := strings.TrimSpace(s)
	if ts == "" {
		return false
	}
	_, err := base64.URLEncoding.DecodeString(ts)
	return err == nil
}

func stringBase64RawURL(s string, _ string) bool {
	ts := strings.TrimSpace(s)
	if ts == "" {
		return false
	}
	_, err := base64.RawURLEncoding.DecodeString(ts)
	return err == nil
}

func stringJSON(s string, _ string) bool {
	return json.NewDecoder(strings.NewReader(s)).Decode(new(interface{})) == nil
}

func stringUnique(s string, _ string) bool {
	seen := make(map[rune]struct{}, utf8.RuneCountInString(s))
	for _, r := range s {
		if _, ok := seen[r]; ok {
			return false
		}
		seen[r] = struct{}{}
	}
	return true
}

func stringStartsWith(s string, param string) bool {
	return strings.HasPrefix(s, param)
}

func stringEndsWith(s string, param string) bool {
	return strings.HasSuffix(s, param)
}

func stringStartsNotWith(s string, param string) bool {
	return !strings.HasPrefix(s, param)
}

func stringEndsNotWith(s string, param string) bool {
	return !strings.HasSuffix(s, param)
}

func stringContains(s string, param string) bool {
	return strings.Contains(s, param)
}

func stringContainsAny(s string, param string) bool {
	return strings.ContainsAny(s, param)
}

func stringContainsRune(s string, param string) bool {
	r, _ := utf8.DecodeRuneInString(param)
	return r != utf8.RuneError && strings.ContainsRune(s, r)
}

func stringExcludes(s string, param string) bool {
	return !strings.Contains(s, param)
}

func stringExcludesAll(s string, param string) bool {
	return !strings.ContainsAny(s, param)
}

func stringExcludesRune(s string, param string) bool {
	r, _ := utf8.DecodeRuneInString(param)
	return r != utf8.RuneError && !strings.ContainsRune(s, r)
}

func stringLowercase(s string, _ string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			return false
		}
	}
	return true
}

func stringUppercase(s string, _ string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			return false
		}
	}
	return true
}

func stringBoolean(s string, _ string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

func stringNumber(s string, _ string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func stringDatetime(s string, param string) bool {
	if param == "" {
		param = time.RFC3339
	}
	_, err := time.Parse(param, s)
	return err == nil
}

func stringTimezone(s string, _ string) bool {
	_, err := time.LoadLocation(s)
	return err == nil
}

func stringLatitude(s string, _ string) bool {
	n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return err == nil && n >= -90 && n <= 90
}

func stringLongitude(s string, _ string) bool {
	n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return err == nil && n >= -180 && n <= 180
}

func stringFile(s string, _ string) bool {
	info, err := os.Stat(s)
	return err == nil && !info.IsDir()
}

func stringFilePath(s string, _ string) bool {
	return strings.TrimSpace(s) != "" && filepath.Clean(s) != "."
}

func stringDir(s string, _ string) bool {
	info, err := os.Stat(s)
	return err == nil && info.IsDir()
}

func stringDirPath(s string, _ string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	cleaned := filepath.Clean(s)
	return cleaned != "." && !strings.Contains(filepath.Base(cleaned), ".")
}

func stringMongoDB(s string, _ string) bool {
	return mongoIDRegex.MatchString(s)
}

func stringLuhnChecksum(s string, _ string) bool {
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
			n = luhnDouble(n)
		}
		sum += n
		double = !double
		digits++
	}
	return digits > 0 && sum%10 == 0
}

func stringDNSRFC1035Label(s string, _ string) bool {
	return len(s) <= 63 && dnsLabelRegex.MatchString(s)
}

func stringSemver(s string, _ string) bool {
	i := 0
	if i < len(s) && s[i] == 'v' {
		i++
	}
	if !parseSemverNum(s, &i) || i >= len(s) || s[i] != '.' {
		return false
	}
	i++
	if !parseSemverNum(s, &i) || i >= len(s) || s[i] != '.' {
		return false
	}
	i++
	if !parseSemverNum(s, &i) {
		return false
	}
	if !parseSemverPreRelease(s, &i) {
		return false
	}
	if !parseSemverBuildMeta(s, &i) {
		return false
	}
	return i == len(s)
}

func stringISBN10(s string, _ string) bool {
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
			return isISBN10CheckDigit(c, sum)
		}
		if c < '0' || c > '9' {
			return false
		}
		sum += int(c-'0') * (11 - digits)
	}
	return false
}

func stringISBN13(s string, _ string) bool {
	digits := 0
	sum := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' || c == ' ' {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
		digits++
		if digits > 13 {
			return false
		}
		if digits < 13 {
			weight := 1
			if digits%2 == 0 {
				weight = 3
			}
			sum += int(c-'0') * weight
		}
	}
	if digits != 13 {
		return false
	}
	check := (10 - sum%10) % 10
	return int(s[len(s)-1]-'0') == check
}

func stringISSN(s string, _ string) bool {
	digits := 0
	sum := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' {
			continue
		}
		digits++
		if digits > 8 {
			return false
		}
		if digits == 8 {
			if c == 'X' || c == 'x' {
				return sum%11 == 10
			}
			if c < '0' || c > '9' {
				return false
			}
			return sum%11 == int(c-'0')
		}
		if c < '0' || c > '9' {
			return false
		}
		sum += int(c-'0') * (9 - digits)
	}
	return false
}

func stringBIC(s string, _ string) bool {
	n := len(s)
	if n != 8 && n != 11 {
		return false
	}
	for i := 0; i < 4; i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	for i := 4; i < 6; i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	for i := 6; i < n; i++ {
		c := s[i]
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}

func stringCron(s string, _ string) bool {
	count := 0
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ' ' || s[i] == '\t' {
			if i > start {
				count++
				if !isValidCronFieldZeroAlloc(s[start:i]) {
					return false
				}
			}
			start = i + 1
		}
	}
	return count == 5 || count == 6
}

func stringDataURI(s string, _ string) bool {
	if len(s) < 6 || !hasDataPrefix(s) {
		return false
	}
	i := 5
	i = skipDataURIMimeType(s, i)
	i = skipDataURIParams(s, i)
	return i < len(s) && s[i] == ','
}

func stringBCP47(s string, _ string) bool {
	if len(s) < 2 {
		return false
	}
	i := 0
	if !isAlpha(s, &i, 2, 3) {
		return false
	}
	i = parseBCP47ExtLang(s, i)
	i = parseBCP47Script(s, i)
	i = parseBCP47Region(s, i)
	return parseBCP47Variants(s, i) == len(s)
}

func stringEthAddr(s string, _ string) bool {
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

func stringBtcAddr(s string, _ string) bool {
	n := len(s)
	if n < 26 || n > 62 {
		return false
	}
	if s[0] == '1' || s[0] == '3' {
		return isBtcLegacyAddr(s, n)
	}
	return isBtcBech32Addr(s, n)
}
