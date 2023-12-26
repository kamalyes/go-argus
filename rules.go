/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
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
	"oneof":             ruleOneOf,
	"oneofci":           ruleOneOfCI,
	"noneof":            ruleNoneOf,
	"noneofci":          ruleNoneOfCI,
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
	return ok && compareLengthOrNumber(field, n, func(actual, expect float64) bool { return actual >= expect })
}

func ruleMax(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, func(actual, expect float64) bool { return actual <= expect })
}

func ruleLen(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, func(actual, expect float64) bool { return actual == expect })
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
	return ok && compareLengthOrNumber(field, n, func(actual, expect float64) bool { return actual > expect })
}

func ruleGte(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, func(actual, expect float64) bool { return actual >= expect })
}

func ruleLt(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, func(actual, expect float64) bool { return actual < expect })
}

func ruleLte(field reflect.Value, param string, _ bool) bool {
	n, ok := parseFloat(param)
	return ok && compareLengthOrNumber(field, n, func(actual, expect float64) bool { return actual <= expect })
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
	return ok && net.ParseIP(strings.TrimSpace(s)) != nil
}

func ruleIPv4(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	ip := net.ParseIP(strings.TrimSpace(s))
	return ok && ip != nil && ip.To4() != nil
}

func ruleIPv6(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	ip := net.ParseIP(strings.TrimSpace(s))
	return ok && ip != nil && ip.To4() == nil
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
	u, err := url.ParseRequestURI(strings.TrimSpace(s))
	return err == nil && u.Scheme != ""
}

func ruleURI(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	_, err := url.ParseRequestURI(strings.TrimSpace(s))
	return err == nil
}

func ruleHTTPURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	u, err := url.Parse(strings.TrimSpace(s))
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func ruleHTTPSURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	if !ok {
		return false
	}
	u, err := url.Parse(strings.TrimSpace(s))
	return err == nil && u.Scheme == "https" && u.Host != ""
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
	if bytes, ok := bytesValue(field); ok {
		return json.Valid(bytes)
	}
	if s, ok := stringValue(field); ok {
		return json.Valid([]byte(s))
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
	return ok && s == strings.ToLower(s)
}

func ruleUppercase(field reflect.Value, _ string, _ bool) bool {
	s, ok := stringValue(field)
	return ok && s == strings.ToUpper(s)
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

func compareLengthOrNumber(field reflect.Value, expect float64, cmp func(float64, float64) bool) bool {
	field = derefValue(field)
	if !field.IsValid() {
		return false
	}
	switch field.Kind() {
	case reflect.String:
		return cmp(float64(utf8.RuneCountInString(field.String())), expect)
	case reflect.Slice, reflect.Array, reflect.Map:
		return cmp(float64(field.Len()), expect)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cmp(float64(field.Int()), expect)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return cmp(float64(field.Uint()), expect)
	case reflect.Float32, reflect.Float64:
		return cmp(field.Float(), expect)
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
