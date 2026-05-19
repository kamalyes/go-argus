/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-19 00:00:00
 * @FilePath: \go-argus\rule\builtin.go
 * @Description: 内置字段规则，负责单字段格式、长度、数值和枚举校验
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/kamalyes/go-argus/validate"
)

type BuiltinRule func(field reflect.Value, param string, requiredStructEnabled bool) bool

var BuiltinRules = map[string]BuiltinRule{
	"required":          RuleRequired,
	"isdefault":         RuleDefault,
	"min":               RuleMin,
	"max":               RuleMax,
	"len":               RuleLen,
	"eq":                RuleEq,
	"eq_ignore_case":    RuleEqIgnoreCase,
	"ne":                RuleNe,
	"ne_ignore_case":    RuleNeIgnoreCase,
	"gt":                RuleGt,
	"gte":               RuleGte,
	"lt":                RuleLt,
	"lte":               RuleLte,
	"alpha":             RuleAlpha,
	"alphaspace":        RuleAlphaSpace,
	"alphanum":          RuleAlphanum,
	"alphanumspace":     RuleAlphanumSpace,
	"alphaunicode":      RuleAlphaUnicode,
	"alphanumunicode":   RuleAlphanumUnicode,
	"ascii":             RuleASCII,
	"printascii":        RulePrintASCII,
	"multibyte":         RuleMultibyte,
	"hexadecimal":       RuleHexadecimal,
	"hexcolor":          RuleHexColor,
	"rgb":               RuleRGB,
	"rgba":              RuleRGBA,
	"hsl":               RuleHSL,
	"hsla":              RuleHSLA,
	"email":             RuleEmail,
	"e164":              RuleE164,
	"ip":                RuleIP,
	"ip_addr":           RuleIP,
	"ipv4":              RuleIPv4,
	"ipv6":              RuleIPv6,
	"cidr":              RuleCIDR,
	"cidrv4":            RuleCIDRv4,
	"cidrv6":            RuleCIDRv6,
	"mac":               RuleMAC,
	"hostname":          RuleHostname,
	"hostname_rfc1123":  RuleHostname,
	"fqdn":              RuleFQDN,
	"hostname_port":     RuleHostnamePort,
	"port":              RulePort,
	"url":               RuleURL,
	"uri":               RuleURI,
	"http_url":          RuleHTTPURL,
	"https_url":         RuleHTTPSURL,
	"url_encoded":       RuleURLEncoded,
	"html":              RuleHTML,
	"html_encoded":      RuleHTMLEncoded,
	"uuid":              RuleUUID,
	"uuid3":             RuleUUID3,
	"uuid4":             RuleUUID4,
	"uuid5":             RuleUUID5,
	"uuid_rfc4122":      RuleUUID,
	"uuid3_rfc4122":     RuleUUID3,
	"uuid4_rfc4122":     RuleUUID4,
	"uuid5_rfc4122":     RuleUUID5,
	"base32":            RuleBase32,
	"base64":            RuleBase64,
	"base64url":         RuleBase64URL,
	"base64rawurl":      RuleBase64RawURL,
	"json":              RuleJSON,
	"unique":            RuleUnique,
	"startswith":        RuleStartsWith,
	"endswith":          RuleEndsWith,
	"startsnotwith":     RuleStartsNotWith,
	"endsnotwith":       RuleEndsNotWith,
	"contains":          RuleContains,
	"containsany":       RuleContainsAny,
	"containsrune":      RuleContainsRune,
	"excludes":          RuleExcludes,
	"excludesall":       RuleExcludesAll,
	"excludesrune":      RuleExcludesRune,
	"lowercase":         RuleLowercase,
	"uppercase":         RuleUppercase,
	"boolean":           RuleBoolean,
	"number":            RuleNumber,
	"numeric":           RuleNumber,
	"datetime":          RuleDatetime,
	"timezone":          RuleTimezone,
	"latitude":          RuleLatitude,
	"longitude":         RuleLongitude,
	"file":              RuleFile,
	"filepath":          RuleFilePath,
	"dir":               RuleDir,
	"dirpath":           RuleDirPath,
	"mongodb":           RuleMongoDB,
	"luhn_checksum":     RuleLuhnChecksum,
	"credit_card":       RuleLuhnChecksum,
	"dns_rfc1035_label": RuleDNSRFC1035Label,
	"semver":            RuleSemver,
	"isbn10":            RuleISBN10,
	"isbn13":            RuleISBN13,
	"issn":              RuleISSN,
	"bic":               RuleBIC,
	"cron":              RuleCron,
	"datauri":           RuleDataURI,
	"bcp47":             RuleBCP47,
	"eth_addr":          RuleEthAddr,
	"btc_addr":          RuleBtcAddr,
}

type CmpOp int

const (
	CmpGTE CmpOp = iota
	CmpLTE
	CmpGT
	CmpLT
	CmpEQ
)

func CompareOp(actual, expect float64, op CmpOp) bool {
	switch op {
	case CmpGTE:
		return actual >= expect
	case CmpLTE:
		return actual <= expect
	case CmpGT:
		return actual > expect
	case CmpLT:
		return actual < expect
	default:
		return actual == expect
	}
}

func CompareLengthOrNumber(field reflect.Value, expect float64, op CmpOp) bool {
	actual, ok := ResolveFieldValue(field)
	if !ok {
		return false
	}
	return CompareOp(actual, expect, op)
}

func ResolveFieldValue(field reflect.Value) (float64, bool) {
	field = validate.DerefReflect(field)
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

func RuleRequired(field reflect.Value, _ string, requiredStructEnabled bool) bool {
	return !validate.IsEmptyValueWithStruct(field, requiredStructEnabled)
}

func RuleDefault(field reflect.Value, _ string, requiredStructEnabled bool) bool {
	return validate.IsEmptyValueWithStruct(field, requiredStructEnabled)
}

func RuleMin(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, CmpGTE)
}

func RuleMax(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, CmpLTE)
}

func RuleLen(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, CmpEQ)
}

func RuleEq(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.ScalarString(field)
	return ok && s == param
}

func RuleEqIgnoreCase(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.ScalarString(field)
	return ok && strings.EqualFold(s, param)
}

func RuleNe(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.ScalarString(field)
	return ok && s != param
}

func RuleNeIgnoreCase(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.ScalarString(field)
	return ok && !strings.EqualFold(s, param)
}

func RuleGt(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, CmpGT)
}

func RuleGte(field reflect.Value, param string, hasCtx bool) bool {
	return RuleMin(field, param, hasCtx)
}

func RuleLt(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, CmpLT)
}

func RuleLte(field reflect.Value, param string, hasCtx bool) bool {
	return RuleMax(field, param, hasCtx)
}

func RuleAlpha(field reflect.Value, _ string, _ bool) bool {
	return validate.MatchStringRunes(field, func(r rune) bool { return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') })
}

func RuleAlphaSpace(field reflect.Value, _ string, _ bool) bool {
	return validate.MatchStringRunes(field, func(r rune) bool { return r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') })
}

func RuleAlphanum(field reflect.Value, _ string, _ bool) bool {
	return validate.MatchStringRunes(field, func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	})
}

func RuleAlphanumSpace(field reflect.Value, _ string, _ bool) bool {
	return validate.MatchStringRunes(field, func(r rune) bool {
		return r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	})
}

func RuleAlphaUnicode(field reflect.Value, _ string, _ bool) bool {
	return validate.MatchStringRunes(field, unicode.IsLetter)
}

func RuleAlphanumUnicode(field reflect.Value, _ string, _ bool) bool {
	return validate.MatchStringRunes(field, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) })
}

func RuleASCII(field reflect.Value, _ string, _ bool) bool {
	return validate.MatchStringRunes(field, func(r rune) bool { return r <= unicode.MaxASCII })
}

func RulePrintASCII(field reflect.Value, _ string, _ bool) bool {
	return validate.MatchStringRunes(field, func(r rune) bool { return r >= 0x20 && r <= 0x7e })
}

func RuleMultibyte(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringMultibyte(s)
}

func RuleHexadecimal(field reflect.Value, _ string, _ bool) bool {
	return validate.MatchStringRunes(field, func(r rune) bool {
		return (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') || (r >= '0' && r <= '9')
	})
}

func RuleHexColor(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringHexColor(s)
}

func RuleRGB(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringRGB(s)
}

func RuleRGBA(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringRGBA(s)
}

func RuleHSL(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringHSL(s)
}

func RuleHSLA(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringHSLA(s)
}

func RuleEmail(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsEmail(s)
}

func RuleE164(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringE164(s)
}

func RuleIP(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringIP(s)
}

func RuleIPv4(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringIPv4(s)
}

func RuleIPv6(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringIPv6(s)
}

func RuleCIDR(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringCIDR(s)
}

func RuleCIDRv4(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringCIDRv4(s)
}

func RuleCIDRv6(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringCIDRv6(s)
}

func RuleMAC(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringMAC(s)
}

func RuleHostname(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringHostname(s)
}

func RuleFQDN(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringFQDN(s)
}

func RuleHostnamePort(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringHostnamePort(s)
}

func RulePort(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.ScalarString(field)
	return ok && validate.StringPort(s)
}

func RuleURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringURL(s)
}

func RuleURI(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringURI(s)
}

func RuleHTTPURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringHTTPURL(s)
}

func RuleHTTPSURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringHTTPSURL(s)
}

func RuleURLEncoded(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringURLEncoded(s)
}

func RuleHTML(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringHTML(s)
}

func RuleHTMLEncoded(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringHTMLEncoded(s)
}

func RuleUUID(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringUUID(s)
}

func RuleUUID3(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringUUID3(s)
}

func RuleUUID4(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringUUID4(s)
}

func RuleUUID5(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringUUID5(s)
}

func RuleBase32(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringBase32(s)
}

func RuleBase64(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringBase64(s)
}

func RuleBase64URL(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringBase64URL(s)
}

func RuleBase64RawURL(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringBase64RawURL(s)
}

func RuleJSON(field reflect.Value, _ string, _ bool) bool {
	if b, ok := validate.BytesValue(field); ok {
		return validate.IsValidJSONBytes(b)
	}
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringJSON(s)
}

func RuleUnique(field reflect.Value, _ string, _ bool) bool {
	field = validate.DerefReflect(field)
	if !field.IsValid() {
		return false
	}
	switch field.Kind() {
	case reflect.String:
		return validate.StringUnique(field.String())
	case reflect.Slice, reflect.Array:
		return IsUniqueSlice(field)
	case reflect.Map:
		return IsUniqueMap(field)
	default:
		return false
	}
}

func IsUniqueSlice(field reflect.Value) bool {
	seen := make(map[string]struct{}, field.Len())
	for i := 0; i < field.Len(); i++ {
		key := validate.StringValue(validate.DerefReflect(field.Index(i)))
		if _, ok := seen[key]; ok {
			return false
		}
		seen[key] = struct{}{}
	}
	return true
}

func IsUniqueMap(field reflect.Value) bool {
	seen := make(map[string]struct{}, field.Len())
	for _, key := range field.MapKeys() {
		valueKey := validate.StringValue(validate.DerefReflect(field.MapIndex(key)))
		if _, ok := seen[valueKey]; ok {
			return false
		}
		seen[valueKey] = struct{}{}
	}
	return true
}

func RuleStartsWith(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringStartsWith(s, param)
}

func RuleEndsWith(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringEndsWith(s, param)
}

func RuleStartsNotWith(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringStartsNotWith(s, param)
}

func RuleEndsNotWith(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringEndsNotWith(s, param)
}

func RuleContains(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringContains(s, param)
}

func RuleContainsAny(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringContainsAny(s, param)
}

func RuleContainsRune(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringContainsRune(s, param)
}

func RuleExcludes(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringExcludes(s, param)
}

func RuleExcludesAll(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringExcludesAll(s, param)
}

func RuleExcludesRune(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringExcludesRune(s, param)
}

func RuleLowercase(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringLowercase(s)
}

func RuleUppercase(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringUppercase(s)
}

func RuleBoolean(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	if !ok {
		return validate.DerefReflect(field).Kind() == reflect.Bool
	}
	return validate.StringBoolean(s)
}

func RuleNumber(field reflect.Value, _ string, _ bool) bool {
	field = validate.DerefReflect(field)
	if !field.IsValid() {
		return false
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return validate.StringNumber(field.String())
	default:
		return false
	}
}

func RuleDatetime(field reflect.Value, param string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringDatetime(s, param)
}

func RuleTimezone(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringTimezone(s)
}

func RuleLatitude(field reflect.Value, _ string, _ bool) bool {
	if n, ok := validate.NumericValue(field); ok {
		return n >= -90 && n <= 90
	}
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringLatitude(s)
}

func RuleLongitude(field reflect.Value, _ string, _ bool) bool {
	if n, ok := validate.NumericValue(field); ok {
		return n >= -180 && n <= 180
	}
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringLongitude(s)
}

func RuleFile(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringFile(s)
}

func RuleFilePath(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringFilePath(s)
}

func RuleDir(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringDir(s)
}

func RuleDirPath(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringDirPath(s)
}

func RuleMongoDB(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringMongoDB(s)
}

func RuleLuhnChecksum(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsLuhnChecksum(s)
}

func RuleDNSRFC1035Label(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringDNSRFC1035Label(s)
}

func RuleSemver(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsSemver(s)
}

func RuleISBN10(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsISBN10(s)
}

func RuleISBN13(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsISBN13(s)
}

func RuleISSN(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsISSN(s)
}

func RuleBIC(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsBIC(s)
}

func RuleCron(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsCron(s)
}

func RuleDataURI(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsDataURI(s)
}

func RuleBCP47(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsBCP47(s)
}

func RuleEthAddr(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsEthAddr(s)
}

func RuleBtcAddr(field reflect.Value, _ string, _ bool) bool {
	s, ok := validate.StringValueFromField(field)
	return ok && validate.IsBtcAddr(s)
}
