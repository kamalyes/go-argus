/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-18 10:53:49
 * @FilePath: \go-argus\rules.go
 * @Description: 根包内置字段规则，负责单字段格式、长度、数值和枚举校验
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
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
	s, ok := scalarString(field)
	return ok && s == param
}

func ruleEqIgnoreCase(field reflect.Value, param string, _ bool) bool {
	s, ok := scalarString(field)
	return ok && strings.EqualFold(s, param)
}

func ruleNe(field reflect.Value, param string, _ bool) bool {
	s, ok := scalarString(field)
	return ok && s != param
}

func ruleNeIgnoreCase(field reflect.Value, param string, _ bool) bool {
	s, ok := scalarString(field)
	return ok && !strings.EqualFold(s, param)
}

func ruleGt(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, cmpGT)
}

func ruleGte(field reflect.Value, param string, hasCtx bool) bool {
	return ruleMin(field, param, hasCtx)
}

func ruleLt(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, cmpLT)
}

func ruleLte(field reflect.Value, param string, hasCtx bool) bool {
	return ruleMax(field, param, hasCtx)
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
	return ok && stringMultibyte(s, "")
}

func ruleHexadecimal(field reflect.Value, _ string, _ bool) bool {
	return matchStringRunes(field, func(r rune) bool {
		return (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') || (r >= '0' && r <= '9')
	})
}

func ruleHexColor(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringHexColor(s, "")
}

func ruleRGB(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringRGB(s, "")
}

func ruleRGBA(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringRGBA(s, "")
}

func ruleHSL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringHSL(s, "")
}

func ruleHSLA(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringHSLA(s, "")
}

func ruleEmail(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringEmail(s, "")
}

func ruleE164(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringE164(s, "")
}

func ruleIP(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringIP(s, "")
}

func ruleIPv4(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringIPv4(s, "")
}

func ruleIPv6(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringIPv6(s, "")
}

func ruleCIDR(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringCIDR(s, "")
}

func ruleCIDRv4(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringCIDRv4(s, "")
}

func ruleCIDRv6(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringCIDRv6(s, "")
}

func ruleMAC(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringMAC(s, "")
}

func ruleHostname(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringHostname(s, "")
}

func ruleFQDN(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringFQDN(s, "")
}

func ruleHostnamePort(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringHostnamePort(s, "")
}

func rulePort(field reflect.Value, _ string, _ bool) bool {
	s, ok := scalarString(field)
	return ok && stringPort(s, "")
}

func ruleURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringURL(s, "")
}

func ruleURI(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringURI(s, "")
}

func ruleHTTPURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringHTTPURL(s, "")
}

func ruleHTTPSURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringHTTPSURL(s, "")
}

func ruleURLEncoded(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringURLEncoded(s, "")
}

func ruleHTML(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringHTML(s, "")
}

func ruleHTMLEncoded(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringHTMLEncoded(s, "")
}

func ruleUUID(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringUUID(s, "")
}

func ruleUUID3(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringUUID3(s, "")
}

func ruleUUID4(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringUUID4(s, "")
}

func ruleUUID5(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringUUID5(s, "")
}

func ruleBase32(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringBase32(s, "")
}

func ruleBase64(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringBase64(s, "")
}

func ruleBase64URL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringBase64URL(s, "")
}

func ruleBase64RawURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringBase64RawURL(s, "")
}

func ruleJSON(field reflect.Value, _ string, _ bool) bool {
	if b, ok := bytesValue(field); ok {
		return validate.IsValidJSONBytes(b)
	}
	s, ok := stringValue(field)
	return ok && stringJSON(s, "")
}

func ruleUnique(field reflect.Value, _ string, _ bool) bool {
	field = derefValue(field)
	if !field.IsValid() {
		return false
	}
	switch field.Kind() {
	case reflect.String:
		return stringUnique(field.String(), "")
	case reflect.Slice, reflect.Array:
		return isUniqueSlice(field)
	case reflect.Map:
		return isUniqueMap(field)
	default:
		return false
	}
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
	return ok && stringStartsWith(s, param)
}

func ruleEndsWith(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringEndsWith(s, param)
}

func ruleStartsNotWith(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringStartsNotWith(s, param)
}

func ruleEndsNotWith(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringEndsNotWith(s, param)
}

func ruleContains(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringContains(s, param)
}

func ruleContainsAny(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringContainsAny(s, param)
}

func ruleContainsRune(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringContainsRune(s, param)
}

func ruleExcludes(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringExcludes(s, param)
}

func ruleExcludesAll(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringExcludesAll(s, param)
}

func ruleExcludesRune(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringExcludesRune(s, param)
}

func ruleLowercase(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringLowercase(s, "")
}

func ruleUppercase(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringUppercase(s, "")
}

func ruleBoolean(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return derefValue(field).Kind() == reflect.Bool
	}
	return stringBoolean(s, "")
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
		return stringNumber(field.String(), "")
	default:
		return false
	}
}

func ruleDatetime(field reflect.Value, param string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringDatetime(s, param)
}

func ruleTimezone(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringTimezone(s, "")
}

func ruleLatitude(field reflect.Value, _ string, _ bool) bool {
	if n, ok := numericValue(field); ok {
		return n >= -90 && n <= 90
	}
	s, ok := stringValue(field)
	return ok && stringLatitude(s, "")
}

func ruleLongitude(field reflect.Value, _ string, _ bool) bool {
	if n, ok := numericValue(field); ok {
		return n >= -180 && n <= 180
	}
	s, ok := stringValue(field)
	return ok && stringLongitude(s, "")
}

func ruleFile(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringFile(s, "")
}

func ruleFilePath(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringFilePath(s, "")
}

func ruleDir(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringDir(s, "")
}

func ruleDirPath(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringDirPath(s, "")
}

func ruleMongoDB(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringMongoDB(s, "")
}

func ruleLuhnChecksum(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringLuhnChecksum(s, "")
}

func ruleDNSRFC1035Label(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringDNSRFC1035Label(s, "")
}

func ruleSemver(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringSemver(s, "")
}

func ruleISBN10(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringISBN10(s, "")
}

func ruleISBN13(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringISBN13(s, "")
}

func ruleISSN(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringISSN(s, "")
}

func ruleBIC(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringBIC(s, "")
}

func ruleCron(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringCron(s, "")
}

func ruleDataURI(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringDataURI(s, "")
}

func ruleBCP47(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringBCP47(s, "")
}

func ruleEthAddr(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringEthAddr(s, "")
}

func ruleBtcAddr(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && stringBtcAddr(s, "")
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
	default:
		return actual == expect
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
		return parseFloat(strings.TrimSpace(field.String()))
	default:
		return 0, false
	}
}

func parseFloat(s string) (float64, bool) {
	n, err := parseFloatStr(s)
	return n, err == nil
}

func parseFloatStr(s string) (float64, error) {
	var n float64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			n = n*10 + float64(c-'0')
			continue
		}
		if c == '.' {
			frac := 0.0
			div := 1.0
			for j := i + 1; j < len(s); j++ {
				if s[j] < '0' || s[j] > '9' {
					return 0, fmt.Errorf("invalid float")
				}
				frac = frac*10 + float64(s[j]-'0')
				div *= 10
			}
			n += frac / div
			break
		}
		return 0, fmt.Errorf("invalid float")
	}
	return n, nil
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
	if field.Kind() == reflect.String {
		return field.String(), true
	}
	if field.CanInterface() {
		return fmt.Sprint(field.Interface()), true
	}
	return "", false
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

func luhnDouble(n int) int {
	n *= 2
	if n > 9 {
		return n - 9
	}
	return n
}

func parseSemverNum(s string, i *int) bool {
	return validate.ParseSemverNum(s, i)
}

func parseSemverPreRelease(s string, i *int) bool {
	return validate.ParseSemverPreRelease(s, i)
}

func parseSemverBuildMeta(s string, i *int) bool {
	return validate.ParseSemverBuildMeta(s, i)
}

func isISBN10CheckDigit(c byte, sum int) bool {
	if c == 'X' || c == 'x' {
		return (sum+10)%11 == 0
	}
	if c < '0' || c > '9' {
		return false
	}
	return (sum+int(c-'0'))%11 == 0
}

func isValidCronFieldZeroAlloc(field string) bool {
	return validate.IsValidCronField(field)
}

func hasDataPrefix(s string) bool {
	return len(s) >= 5 && s[0] == 'd' && s[1] == 'a' && s[2] == 't' && s[3] == 'a' && s[4] == ':'
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
		if i+6 <= len(s) && s[i] == 'b' && s[i+1] == 'a' && s[i+2] == 's' && s[i+3] == 'e' && s[i+4] == '6' && s[i+5] == '4' {
			i += 6
		}
		for i < len(s) && s[i] != ';' && s[i] != ',' {
			if s[i] < ' ' || s[i] > '~' {
				return len(s)
			}
			i++
		}
	}
	return i
}

func isAlpha(s string, i *int, minLen, maxLen int) bool {
	start := *i
	for *i < len(s) && *i-start < maxLen {
		c := s[*i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			break
		}
		*i++
	}
	n := *i - start
	return n >= minLen && n <= maxLen
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

func isAlphanum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
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
	return (c >= '1' && c <= '9') ||
		(c >= 'A' && c <= 'H') ||
		(c >= 'J' && c <= 'N') ||
		(c >= 'P' && c <= 'Z') ||
		(c >= 'a' && c <= 'k') ||
		(c >= 'm' && c <= 'z')
}
