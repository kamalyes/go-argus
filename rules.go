/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 01:58:16
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
		c := s[i]
		if c == '/' || c == '?' || c == '#' {
			break
		}
		if c == ':' || c == '@' {
			continue
		}
		if hostStart == i && (c == '.' || c == '-') {
			return false
		}
	}
	return true
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

func ruleOneOf(field reflect.Value, param string, _ bool) bool {
	actual, ok := scalarString(field)
	if !ok {
		return false
	}
	for _, item := range strings.Fields(param) {
		if actual == item {
			return true
		}
	}
	return false
}

func ruleOneOfCI(field reflect.Value, param string, _ bool) bool {
	actual, ok := scalarString(field)
	if !ok {
		return false
	}
	for _, item := range strings.Fields(param) {
		if strings.EqualFold(actual, item) {
			return true
		}
	}
	return false
}

func ruleNoneOf(field reflect.Value, param string, _ bool) bool {
	return !ruleOneOf(field, param, false)
}

func ruleNoneOfCI(field reflect.Value, param string, _ bool) bool {
	return !ruleOneOfCI(field, param, false)
}

func ruleUnique(field reflect.Value, _ string, _ bool) bool {
	field = derefValue(field)
	if !field.IsValid() {
		return false
	}
	switch field.Kind() {
	case reflect.String:
		seen := make(map[rune]struct{}, utf8.RuneCountInString(field.String()))
		for _, r := range field.String() {
			if _, ok := seen[r]; ok {
				return false
			}
			seen[r] = struct{}{}
		}
		return true
	case reflect.Slice, reflect.Array:
		seen := make(map[string]struct{}, field.Len())
		for i := 0; i < field.Len(); i++ {
			key := toStringValue(derefValue(field.Index(i)))
			if _, ok := seen[key]; ok {
				return false
			}
			seen[key] = struct{}{}
		}
		return true
	case reflect.Map:
		seen := make(map[string]struct{}, field.Len())
		for _, key := range field.MapKeys() {
			valueKey := toStringValue(derefValue(field.MapIndex(key)))
			if _, ok := seen[valueKey]; ok {
				return false
			}
			seen[valueKey] = struct{}{}
		}
		return true
	default:
		return false
	}
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
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		double = !double
		digits++
	}
	return digits > 0 && sum%10 == 0
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
	field = derefValue(field)
	if !field.IsValid() {
		return false
	}
	var actual float64
	switch field.Kind() {
	case reflect.String:
		actual = float64(utf8.RuneCountInString(field.String()))
	case reflect.Slice, reflect.Array, reflect.Map:
		actual = float64(field.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		actual = float64(field.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		actual = float64(field.Uint())
	case reflect.Float32, reflect.Float64:
		actual = field.Float()
	default:
		return false
	}
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

var (
	semverRegex  = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	bicRegex     = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}(?:[A-Z0-9]{3})?$`)
	ethAddrRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	btcAddrRegex = regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$|^[bc1q][a-z0-9]{39,59}$`)
	bcp47Regex   = regexp.MustCompile(`^[a-zA-Z]{2,3}(?:-[a-zA-Z]{4})?(?:-(?:[a-zA-Z]{2}|\d{3}))?(?:-[a-zA-Z0-9]{5,8})*(?:-[a-zA-Z0-9]{1,8})*$`)
	cronFieldRe  = regexp.MustCompile(`^\S+$`)
	datauriRegex = regexp.MustCompile(`^data:([^;,]*)?(;[^;,]*)*;?(base64,)?,`)
)

func ruleSemver(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && semverRegex.MatchString(s)
}

func ruleISBN10(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func ruleISBN13(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func ruleISSN(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
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

func ruleBIC(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && bicRegex.MatchString(s)
}

func ruleCron(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	fields := strings.Fields(s)
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

func ruleDataURI(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok || !strings.HasPrefix(s, "data:") {
		return false
	}
	return datauriRegex.MatchString(s)
}

func ruleBCP47(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && bcp47Regex.MatchString(s)
}

func ruleEthAddr(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && ethAddrRegex.MatchString(s)
}

func ruleBtcAddr(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && btcAddrRegex.MatchString(s)
}
