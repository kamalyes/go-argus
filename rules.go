/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 02:28:05
 * @FilePath: \go-argus\rules.go
 * @Description: 根包内置字段规则，负责单字段格式、长度、数值和枚举校验
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
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
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/kamalyes/go-argus/validate"
)

type builtinRule func(field reflect.Value, param string, requiredStructEnabled bool) bool

var builtinRules = map[string]builtinRule{
	"required":          ruleRequired,
	"isdefault":         ruleDefault,
	"min":               ruleMin,
	"max":               ruleMax,
	"len":               ruleLen,
	"eq":                ruleEq,
	"eq_ignore_case":    ruleEqIgnoreCase,
	"ne":                ruleNe,
	"ne_ignore_case":    ruleNeIgnoreCase,
	"gt":                ruleGt,
	"gte":               ruleGte,
	"lt":                ruleLt,
	"lte":               ruleLte,
	"alpha":             ruleAlpha,
	"alphaspace":        ruleAlphaSpace,
	"alphanum":          ruleAlphanum,
	"alphanumspace":     ruleAlphanumSpace,
	"alphaunicode":      ruleAlphaUnicode,
	"alphanumunicode":   ruleAlphanumUnicode,
	"ascii":             ruleASCII,
	"printascii":        rulePrintASCII,
	"multibyte":         ruleMultibyte,
	"hexadecimal":       ruleHexadecimal,
	"hexcolor":          ruleHexColor,
	"rgb":               ruleRGB,
	"rgba":              ruleRGBA,
	"hsl":               ruleHSL,
	"hsla":              ruleHSLA,
	"email":             ruleEmail,
	"e164":              ruleE164,
	"ip":                ruleIP,
	"ip_addr":           ruleIP,
	"ipv4":              ruleIPv4,
	"ipv6":              ruleIPv6,
	"cidr":              ruleCIDR,
	"cidrv4":            ruleCIDRv4,
	"cidrv6":            ruleCIDRv6,
	"mac":               ruleMAC,
	"hostname":          ruleHostname,
	"hostname_rfc1123":  ruleHostname,
	"fqdn":              ruleFQDN,
	"hostname_port":     ruleHostnamePort,
	"port":              rulePort,
	"url":               ruleURL,
	"uri":               ruleURI,
	"http_url":          ruleHTTPURL,
	"https_url":         ruleHTTPSURL,
	"url_encoded":       ruleURLEncoded,
	"html":              ruleHTML,
	"html_encoded":      ruleHTMLEncoded,
	"uuid":              ruleUUID,
	"uuid3":             ruleUUID3,
	"uuid4":             ruleUUID4,
	"uuid5":             ruleUUID5,
	"uuid_rfc4122":      ruleUUID,
	"uuid3_rfc4122":     ruleUUID3,
	"uuid4_rfc4122":     ruleUUID4,
	"uuid5_rfc4122":     ruleUUID5,
	"base32":            ruleBase32,
	"base64":            ruleBase64,
	"base64url":         ruleBase64URL,
	"base64rawurl":      ruleBase64RawURL,
	"json":              ruleJSON,
	"unique":            ruleUnique,
	"startswith":        ruleStartsWith,
	"endswith":          ruleEndsWith,
	"startsnotwith":     ruleStartsNotWith,
	"endsnotwith":       ruleEndsNotWith,
	"contains":          ruleContains,
	"containsany":       ruleContainsAny,
	"containsrune":      ruleContainsRune,
	"excludes":          ruleExcludes,
	"excludesall":       ruleExcludesAll,
	"excludesrune":      ruleExcludesRune,
	"lowercase":         ruleLowercase,
	"uppercase":         ruleUppercase,
	"boolean":           ruleBoolean,
	"number":            ruleNumber,
	"numeric":           ruleNumber,
	"datetime":          ruleDatetime,
	"timezone":          ruleTimezone,
	"latitude":          ruleLatitude,
	"longitude":         ruleLongitude,
	"file":              ruleFile,
	"filepath":          ruleFilePath,
	"dir":               ruleDir,
	"dirpath":           ruleDirPath,
	"mongodb":           ruleMongoDB,
	"luhn_checksum":     ruleLuhnChecksum,
	"credit_card":       ruleLuhnChecksum,
	"dns_rfc1035_label": ruleDNSRFC1035Label,
	"semver":            ruleSemver,
	"isbn10":            ruleISBN10,
	"isbn13":            ruleISBN13,
	"issn":              ruleISSN,
	"bic":               ruleBIC,
	"cron":              ruleCron,
	"datauri":           ruleDataURI,
	"bcp47":             ruleBCP47,
	"eth_addr":          ruleEthAddr,
	"btc_addr":          ruleBtcAddr,
}

var (
	colorHexRegex      = regexp.MustCompile(`^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	rgbRegex           = regexp.MustCompile(`^rgb\(\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*\)$`)
	rgbaRegex          = regexp.MustCompile(`^rgba\(\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(25[0-5]|2[0-4]\d|1?\d?\d)\s*,\s*(0|1|0?\.\d+)\s*\)$`)
	hslRegex           = regexp.MustCompile(`^hsl\(\s*(360|3[0-5]\d|[12]?\d?\d)\s*,\s*(100|[1-9]?\d)%\s*,\s*(100|[1-9]?\d)%\s*\)$`)
	hslaRegex          = regexp.MustCompile(`^hsla\(\s*(360|3[0-5]\d|[12]?\d?\d)\s*,\s*(100|[1-9]?\d)%\s*,\s*(100|[1-9]?\d)%\s*,\s*(0|1|0?\.\d+)\s*\)$`)
	e164Regex          = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	hostnameLabelRegex = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?$`)
	mongoIDRegex       = regexp.MustCompile(`^[0-9a-fA-F]{24}$`)
	dnsLabelRegex      = regexp.MustCompile(`^[a-z]([-a-z0-9]*[a-z0-9])?$`)
)

func ruleRequired(field reflect.Value, _ string, requiredStructEnabled bool) bool {
	return !isEmptyValue(field, requiredStructEnabled)
}

func ruleDefault(field reflect.Value, _ string, requiredStructEnabled bool) bool {
	return isEmptyValue(field, requiredStructEnabled)
}

func ruleMin(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, cmpGTE)
}

func ruleMax(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, cmpLTE)
}

func ruleLen(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, cmpEQ)
}

func ruleEq(field reflect.Value, param string, _ bool) bool {
	actual, ok := scalarString(field)
	return ok && actual == param
}

func ruleEqIgnoreCase(field reflect.Value, param string, _ bool) bool {
	actual, ok := scalarString(field)
	return ok && strings.EqualFold(actual, param)
}

func ruleNe(field reflect.Value, param string, _ bool) bool {
	actual, ok := scalarString(field)
	return ok && actual != param
}

func ruleNeIgnoreCase(field reflect.Value, param string, _ bool) bool {
	actual, ok := scalarString(field)
	return ok && !strings.EqualFold(actual, param)
}

func ruleGt(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, cmpGT)
}

func ruleGte(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, cmpGTE)
}

func ruleLt(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, cmpLT)
}

func ruleLte(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, cmpLTE)
}

func ruleAlpha(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, func(r rune) bool { return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') })
}

func ruleAlphaSpace(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, func(r rune) bool { return r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') })
}

func ruleAlphanum(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	})
}

func ruleAlphanumSpace(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, func(r rune) bool {
		return r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	})
}

func ruleAlphaUnicode(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, unicode.IsLetter)
}

func ruleAlphanumUnicode(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) })
}

func ruleASCII(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, func(r rune) bool { return r <= unicode.MaxASCII })
}

func rulePrintASCII(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, func(r rune) bool { return r >= 0x20 && r <= 0x7e })
}

func ruleMultibyte(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && len(s) != utf8.RuneCountInString(s)
}

func ruleHexadecimal(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, func(r rune) bool {
		return (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') || (r >= '0' && r <= '9')
	})
}

func ruleHexColor(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && colorHexRegex.MatchString(s)
}

func ruleRGB(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && rgbRegex.MatchString(s)
}

func ruleRGBA(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && rgbaRegex.MatchString(s)
}

func ruleHSL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && hslRegex.MatchString(s)
}

func ruleHSLA(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && hslaRegex.MatchString(s)
}

func ruleEmail(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && validate.IsEmail(s)
}

func ruleE164(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && e164Regex.MatchString(s)
}

func ruleIP(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && net.ParseIP(trimSpaceIfNeeded(s)) != nil
}

func ruleIPv4(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	ip := net.ParseIP(trimSpaceIfNeeded(s))
	return ip != nil && ip.To4() != nil
}

func ruleIPv6(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	ip := net.ParseIP(trimSpaceIfNeeded(s))
	return ip != nil && ip.To4() == nil
}

func ruleCIDR(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	_, _, err := net.ParseCIDR(strings.TrimSpace(s))
	return err == nil
}

func ruleCIDRv4(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	ip, _, err := net.ParseCIDR(strings.TrimSpace(s))
	return err == nil && ip.To4() != nil
}

func ruleCIDRv6(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	ip, _, err := net.ParseCIDR(strings.TrimSpace(s))
	return err == nil && ip.To4() == nil
}

func ruleMAC(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	_, err := net.ParseMAC(strings.TrimSpace(s))
	return err == nil
}

func ruleHostname(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	return isHostname(strings.TrimSuffix(strings.TrimSpace(s), "."))
}

func ruleFQDN(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && strings.HasSuffix(strings.TrimSpace(s), ".") && isHostname(strings.TrimSuffix(strings.TrimSpace(s), "."))
}

func ruleHostnamePort(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	host, port, err := net.SplitHostPort(s)
	return err == nil && host != "" && rulePort(reflect.ValueOf(port), "", false)
}

func rulePort(field reflect.Value, _ string, _ bool) bool {
	s, ok := scalarString(field)
	if !ok {
		return false
	}
	n, err := strconv.Atoi(s)
	return err == nil && n >= 0 && n <= 65535
}

func ruleURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	return hasSchemeAndHost(s)
}

func ruleURI(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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
	if len(s) <= colon+1 {
		return false
	}
	return true
}

func ruleHTTPURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func ruleHTTPSURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	colon := strings.Index(s, ":")
	if colon < 0 {
		return false
	}
	scheme := s[:colon]
	if scheme != "https" {
		return false
	}
	return hasHostAfterScheme(s, colon)
}

func hasSchemeAndHost(s string) bool {
	colon := strings.Index(s, ":")
	if colon < 1 {
		return false
	}
	return hasHostAfterScheme(s, colon)
}

func hasHostAfterScheme(s string, colonIdx int) bool {
	if len(s) <= colonIdx+2 || s[colonIdx+1] != '/' || s[colonIdx+2] != '/' {
		return false
	}
	hostStart := colonIdx + 3
	if hostStart >= len(s) {
		return false
	}
	for i := hostStart; i < len(s); i++ {
		if isHostTerminator(s[i]) {
			break
		}
		if !isValidHostChar(s, hostStart, i) {
			return false
		}
	}
	return true
}

func isHostTerminator(c byte) bool {
	return c == '/' || c == '?' || c == '#'
}

func isValidHostChar(s string, hostStart, i int) bool {
	c := s[i]
	if c == ':' || c == '@' {
		return true
	}
	return !(hostStart == i && (c == '.' || c == '-'))
}

func ruleURLEncoded(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok || !strings.Contains(s, "%") {
		return false
	}
	_, err := url.QueryUnescape(s)
	return err == nil
}

func ruleHTML(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && strings.Contains(s, "<") && strings.Contains(s, ">")
}

func ruleHTMLEncoded(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && html.UnescapeString(s) != s
}

func ruleUUID(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && validate.IsUUID(s)
}

func ruleUUID3(field reflect.Value, _ string, _ bool) bool {
	return uuidVersion(field, '3')
}

func ruleUUID4(field reflect.Value, _ string, _ bool) bool {
	return uuidVersion(field, '4')
}

func ruleUUID5(field reflect.Value, _ string, _ bool) bool {
	return uuidVersion(field, '5')
}

func ruleBase32(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok || strings.TrimSpace(s) == "" {
		return false
	}
	_, err := base32.StdEncoding.DecodeString(strings.TrimSpace(s))
	return err == nil
}

func ruleBase64(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && validate.IsBase64(s)
}

func ruleBase64URL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok || strings.TrimSpace(s) == "" {
		return false
	}
	_, err := base64.URLEncoding.DecodeString(strings.TrimSpace(s))
	return err == nil
}

func ruleBase64RawURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok || strings.TrimSpace(s) == "" {
		return false
	}
	_, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(s))
	return err == nil
}

func ruleJSON(field reflect.Value, _ string, _ bool) bool {
	if b, ok := bytesValue(field); ok {
		return json.Valid(b)
	}
	if s, ok := stringValue(field); ok {
		return json.NewDecoder(strings.NewReader(s)).Decode(new(interface{})) == nil
	}
	return false
}

func ruleUnique(field reflect.Value, _ string, _ bool) bool {
	field = derefValue(field)
	if !field.IsValid() {
		return false
	}
	switch field.Kind() {
	case reflect.String:
		return isUniqueRunes(field.String())
	case reflect.Slice, reflect.Array:
		return isUniqueSlice(field)
	case reflect.Map:
		return isUniqueMap(field)
	default:
		return false
	}
}

func isUniqueRunes(s string) bool {
	seen := make(map[rune]struct{}, utf8.RuneCountInString(s))
	for _, r := range s {
		if _, ok := seen[r]; ok {
			return false
		}
		seen[r] = struct{}{}
	}
	return true
}

func isUniqueSlice(field reflect.Value) bool {
	seen := make(map[string]struct{}, field.Len())
	for i := 0; i < field.Len(); i++ {
		key := toStringValue(derefValue(field.Index(i)))
		if _, ok := seen[key]; ok {
			return false
		}
		seen[key] = struct{}{}
	}
	return true
}

func isUniqueMap(field reflect.Value) bool {
	seen := make(map[string]struct{}, field.Len())
	for _, key := range field.MapKeys() {
		valueKey := toStringValue(derefValue(field.MapIndex(key)))
		if _, ok := seen[valueKey]; ok {
			return false
		}
		seen[valueKey] = struct{}{}
	}
	return true
}

func ruleStartsWith(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && strings.HasPrefix(s, param)
}

func ruleEndsWith(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && strings.HasSuffix(s, param)
}

func ruleStartsNotWith(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && !strings.HasPrefix(s, param)
}

func ruleEndsNotWith(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && !strings.HasSuffix(s, param)
}

func ruleContains(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && strings.Contains(s, param)
}

func ruleContainsAny(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && strings.ContainsAny(s, param)
}

func ruleContainsRune(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	r, _ := utf8.DecodeRuneInString(param)
	return r != utf8.RuneError && strings.ContainsRune(s, r)
}

func ruleExcludes(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && !strings.Contains(s, param)
}

func ruleExcludesAll(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && !strings.ContainsAny(s, param)
}

func ruleExcludesRune(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	r, _ := utf8.DecodeRuneInString(param)
	return r != utf8.RuneError && !strings.ContainsRune(s, r)
}

func ruleLowercase(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

func ruleUppercase(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			return false
		}
	}
	return true
}

func ruleBoolean(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return derefValue(field).Kind() == reflect.Bool
	}
	_, err := strconv.ParseBool(s)
	return err == nil
}

func ruleNumber(field reflect.Value, _ string, _ bool) bool {
	field = derefValue(field)
	if !field.IsValid() {
		return false
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		_, err := strconv.ParseFloat(field.String(), 64)
		return err == nil
	default:
		return false
	}
}

func ruleDatetime(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	if param == "" {
		param = time.RFC3339
	}
	_, err := time.Parse(param, s)
	return err == nil
}

func ruleTimezone(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	_, err := time.LoadLocation(s)
	return err == nil
}

func ruleLatitude(field reflect.Value, _ string, _ bool) bool {
	n, ok := numericValue(field)
	return ok && n >= -90 && n <= 90
}

func ruleLongitude(field reflect.Value, _ string, _ bool) bool {
	n, ok := numericValue(field)
	return ok && n >= -180 && n <= 180
}

func ruleFile(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	info, err := os.Stat(s)
	return err == nil && !info.IsDir()
}

func ruleFilePath(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && strings.TrimSpace(s) != "" && filepath.Clean(s) != "."
}

func ruleDir(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	info, err := os.Stat(s)
	return err == nil && info.IsDir()
}

func ruleDirPath(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok || strings.TrimSpace(s) == "" {
		return false
	}
	cleaned := filepath.Clean(s)
	return cleaned != "." && !strings.Contains(filepath.Base(cleaned), ".")
}

func ruleMongoDB(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && mongoIDRegex.MatchString(s)
}

func ruleLuhnChecksum(field reflect.Value, _ string, _ bool) bool {
	s, ok := scalarString(field)
	if !ok {
		return false
	}
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

func luhnDouble(n int) int {
	n *= 2
	if n > 9 {
		n -= 9
	}
	return n
}

func ruleDNSRFC1035Label(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && len(s) <= 63 && dnsLabelRegex.MatchString(s)
}

type cmpOp int

const (
	cmpGTE cmpOp = iota
	cmpLTE
	cmpGT
	cmpLT
	cmpEQ
)

func compareLengthOrNumber(field reflect.Value, expect float64, op cmpOp) bool {
	actual, ok := resolveFieldValue(field)
	if !ok {
		return false
	}
	return compareOp(actual, expect, op)
}

func resolveFieldValue(field reflect.Value) (float64, bool) {
	field = derefValue(field)
	if !field.IsValid() {
		return 0, false
	}
	switch field.Kind() {
	case reflect.String:
		return float64(utf8.RuneCountInString(field.String())), true
	case reflect.Slice, reflect.Array, reflect.Map:
		return float64(field.Len()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(field.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(field.Uint()), true
	case reflect.Float32, reflect.Float64:
		return field.Float(), true
	default:
		return 0, false
	}
}

func compareOp(actual, expect float64, op cmpOp) bool {
	switch op {
	case cmpGTE:
		return actual >= expect
	case cmpLTE:
		return actual <= expect
	case cmpGT:
		return actual > expect
	case cmpLT:
		return actual < expect
	case cmpEQ:
		return actual == expect
	default:
		return false
	}
}

func numericValue(field reflect.Value) (float64, bool) {
	field = derefValue(field)
	if !field.IsValid() {
		return 0, false
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(field.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(field.Uint()), true
	case reflect.Float32, reflect.Float64:
		return field.Float(), true
	case reflect.String:
		n, err := strconv.ParseFloat(strings.TrimSpace(field.String()), 64)
		return n, err == nil
	default:
		return 0, false
	}
}

func parseFloat(s string) (float64, bool) {
	n, err := strconv.ParseFloat(s, 64)
	return n, err == nil
}

func stringValue(field reflect.Value) (string, bool) {
	field = derefValue(field)
	if !field.IsValid() || field.Kind() != reflect.String {
		return "", false
	}
	return field.String(), true
}

func bytesValue(field reflect.Value) ([]byte, bool) {
	field = derefValue(field)
	if !field.IsValid() {
		return nil, false
	}
	if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Uint8 {
		return field.Bytes(), true
	}
	return nil, false
}

func scalarString(field reflect.Value) (string, bool) {
	field = derefValue(field)
	if !field.IsValid() {
		return "", false
	}
	switch field.Kind() {
	case reflect.String:
		return field.String(), true
	case reflect.Bool:
		return strconv.FormatBool(field.Bool()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(field.Uint(), 10), true
	case reflect.Float32:
		return strconv.FormatFloat(field.Float(), 'f', -1, 32), true
	case reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', -1, 64), true
	default:
		return "", false
	}
}

func matchStringRunes(field reflect.Value, fn func(rune) bool) bool {
	s, ok := stringValue(field)
	if !ok || s == "" {
		return false
	}
	for _, r := range s {
		if !fn(r) {
			return false
		}
	}
	return true
}

func uuidVersion(field reflect.Value, version byte) bool {
	s, ok := stringValue(field)
	return ok && len(s) == 36 && s[14] == version && validate.IsUUID(s)
}

func isHostname(host string) bool {
	if host == "" || len(host) > 253 {
		return false
	}
	labels := strings.Split(host, ".")
	for _, label := range labels {
		if !hostnameLabelRegex.MatchString(label) {
			return false
		}
	}
	return true
}

func trimSpaceIfNeeded(s string) string {
	if len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		return strings.TrimSpace(s)
	}
	return s
}

func ruleSemver(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func parseSemverPreRelease(s string, i *int) bool {
	if *i >= len(s) || s[*i] != '-' {
		return true
	}
	*i++
	if !parseSemverIdent(s, i) {
		return false
	}
	for *i < len(s) && s[*i] == '.' {
		*i++
		if !parseSemverIdent(s, i) {
			return false
		}
	}
	return true
}

func parseSemverBuildMeta(s string, i *int) bool {
	if *i >= len(s) || s[*i] != '+' {
		return true
	}
	*i++
	if !parseSemverBuild(s, i) {
		return false
	}
	for *i < len(s) && s[*i] == '.' {
		*i++
		if !parseSemverBuild(s, i) {
			return false
		}
	}
	return true
}

func parseSemverNum(s string, pos *int) bool {
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

func parseSemverIdent(s string, pos *int) bool {
	start := *pos
	for *pos < len(s) && s[*pos] != '.' && s[*pos] != '+' {
		if !isSemverIdentChar(s[*pos]) {
			return false
		}
		*pos++
	}
	if *pos == start {
		return false
	}
	return hasNonZeroAlphaNum(s, start, *pos)
}

func isSemverIdentChar(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-'
}

func hasNonZeroAlphaNum(s string, start, end int) bool {
	for j := start; j < end; j++ {
		if s[j] != '-' && s[j] != '0' {
			return true
		}
	}
	return end-start <= 1
}

func parseSemverBuild(s string, pos *int) bool {
	if *pos >= len(s) {
		return false
	}
	for *pos < len(s) && s[*pos] != '.' && s[*pos] != '+' {
		if !isSemverIdentChar(s[*pos]) {
			return false
		}
		*pos++
	}
	return *pos > 0 && s[*pos-1] != '.' && s[*pos-1] != '-'
}

func ruleISBN10(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func isISBN10CheckDigit(c byte, sum int) bool {
	if c == 'X' || c == 'x' {
		sum += 10
	} else if c >= '0' && c <= '9' {
		sum += int(c - '0')
	} else {
		return false
	}
	return sum%11 == 0
}

func ruleISBN13(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func ruleISSN(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func ruleBIC(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func ruleCron(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func isValidCronFieldZeroAlloc(field string) bool {
	inRange := false
	inStep := false
	for i := 0; i < len(field); i++ {
		c := field[i]
		switch {
		case c == ',':
			inRange = false
			inStep = false
		case c == '/':
			if inStep {
				return false
			}
			inStep = true
		case c == '-':
			if inRange {
				return false
			}
			inRange = true
		case c == '*':
			if i > 0 && field[i-1] != ',' && field[i-1] != '/' {
				return false
			}
		case c >= '0' && c <= '9':
		default:
			return false
		}
	}
	return true
}

func ruleDataURI(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok || len(s) < 6 || !hasDataPrefix(s) {
		return false
	}
	i := 5
	i = skipDataURIMimeType(s, i)
	i = skipDataURIParams(s, i)
	return i < len(s) && s[i] == ','
}

func hasDataPrefix(s string) bool {
	return s[0] == 'd' && s[1] == 'a' && s[2] == 't' && s[3] == 'a' && s[4] == ':'
}

func skipDataURIMimeType(s string, i int) int {
	for i < len(s) && s[i] != ';' && s[i] != ',' {
		if s[i] < ' ' || s[i] > '~' {
			return len(s)
		}
		i++
	}
	return i
}

func skipDataURIParams(s string, i int) int {
	for i < len(s) && s[i] == ';' {
		i++
		i = skipBase64IfPresent(s, i)
		for i < len(s) && s[i] != ';' && s[i] != ',' {
			if s[i] < ' ' || s[i] > '~' {
				return len(s)
			}
			i++
		}
	}
	return i
}

func skipBase64IfPresent(s string, i int) int {
	if i+6 <= len(s) && s[i] == 'b' && s[i+1] == 'a' && s[i+2] == 's' && s[i+3] == 'e' && s[i+4] == '6' && s[i+5] == '4' {
		return i + 6
	}
	return i
}

func ruleBCP47(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok || len(s) < 2 {
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

func parseBCP47ExtLang(s string, i int) int {
	if i >= len(s) || s[i] != '-' {
		return i
	}
	i++
	if i < len(s) && isAlphaAt(s, i, 4) {
		i += 4
		if i < len(s) && s[i] == '-' {
			i++
		}
	}
	return i
}

func parseBCP47Script(s string, i int) int {
	if i < len(s) && isAlphaAt(s, i, 2) {
		return i + 2
	}
	return i
}

func parseBCP47Region(s string, i int) int {
	if i < len(s) && isDigitAt(s, i, 3) {
		return i + 3
	}
	return i
}

func parseBCP47Variants(s string, i int) int {
	for i < len(s) && s[i] == '-' {
		i++
		start := i
		for i < len(s) && s[i] != '-' {
			if !isAlphanum(s[i]) {
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

func isAlphanum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func isAlpha(s string, pos *int, minLen, maxLen int) bool {
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

func isAlphaAt(s string, pos, length int) bool {
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

func isDigitAt(s string, pos, length int) bool {
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

func ruleEthAddr(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok || len(s) != 42 {
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

func ruleBtcAddr(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	n := len(s)
	if n < 26 || n > 62 {
		return false
	}
	if s[0] == '1' || s[0] == '3' {
		return isBtcLegacyAddr(s, n)
	}
	return isBtcBech32Addr(s, n)
}

func isBtcLegacyAddr(s string, n int) bool {
	if n < 26 || n > 35 {
		return false
	}
	for i := 1; i < n; i++ {
		if !isBase58Char(s[i]) {
			return false
		}
	}
	return true
}

func isBtcBech32Addr(s string, n int) bool {
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

func isBase58Char(c byte) bool {
	if c >= '1' && c <= '9' {
		return true
	}
	if c >= 'A' && c <= 'H' {
		return true
	}
	if c >= 'J' && c <= 'N' {
		return true
	}
	if c >= 'P' && c <= 'Z' {
		return true
	}
	if c >= 'a' && c <= 'k' {
		return true
	}
	if c >= 'm' && c <= 'z' {
		return true
	}
	return false
}
