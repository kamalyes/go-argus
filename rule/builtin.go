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
	"unicode/utf8"

	"github.com/kamalyes/go-argus/validate"
)

// BuiltinRule 内置规则函数签名，接收反射值、参数和 requiredStructEnabled 标志
type BuiltinRule func(field reflect.Value, param string, requiredStructEnabled bool) bool

// stringBuiltinAdapter 将 StringRuleFunc 适配为 BuiltinRule
// 自动从字段提取字符串值，委托给 StringRuleFunc 执行
func stringBuiltinAdapter(fn validate.StringRuleFunc) BuiltinRule {
	return func(field reflect.Value, param string, _ bool) bool {
		s, ok := validate.StringValueFromField(field)
		return ok && fn(s, param)
	}
}

// noParamBuiltinAdapter 将无参数的字符串校验函数直接适配为 BuiltinRule
func noParamBuiltinAdapter(fn func(string) bool) BuiltinRule {
	return func(field reflect.Value, _ string, _ bool) bool {
		s, ok := validate.StringValueFromField(field)
		return ok && fn(s)
	}
}

// noParamScalarBuiltinAdapter 将无参数的字符串校验函数适配为 BuiltinRule（使用 ScalarString 提取值，支持 int 等类型）
func noParamScalarBuiltinAdapter(fn func(string) bool) BuiltinRule {
	return func(field reflect.Value, _ string, _ bool) bool {
		s, ok := validate.ScalarString(field)
		return ok && fn(s)
	}
}

// --- 需要特殊逻辑的规则，无法用适配器自动生成 ---

func RuleRequired(field reflect.Value, _ string, requiredStructEnabled bool) bool {
	return !validate.IsEmptyValueWithStruct(field, requiredStructEnabled)
}

func RuleDefault(field reflect.Value, _ string, requiredStructEnabled bool) bool {
	return validate.IsEmptyValueWithStruct(field, requiredStructEnabled)
}

func RuleMin(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, validate.CmpGTE)
}

func RuleMax(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, validate.CmpLTE)
}

func RuleLen(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, validate.CmpEQ)
}

func RuleGt(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, validate.CmpGT)
}

func RuleLt(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, validate.CmpLT)
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

func RuleJSON(field reflect.Value, _ string, _ bool) bool {
	if b, ok := validate.BytesValue(field); ok {
		return validate.IsValidJSONBytes(b)
	}
	s, ok := validate.StringValueFromField(field)
	return ok && validate.StringJSON(s)
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

// --- 数值/长度比较辅助函数 ---

// CompareLengthOrNumber 对字段值进行长度或数值比较
func CompareLengthOrNumber(field reflect.Value, expect float64, op validate.CmpOp) bool {
	actual, ok := ResolveFieldValue(field)
	if !ok {
		return false
	}
	return validate.CompareOp(actual, expect, op)
}

// ResolveFieldValue 解析字段的数值或长度
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

// IsUniqueSlice 判断切片元素是否唯一
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

// IsUniqueMap 判断 Map 值是否唯一
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

// BuiltinRules 内置规则映射表，evalTable 在 init 中合并此表和 dispatch 表
var BuiltinRules = map[string]BuiltinRule{
	// 需要特殊逻辑的规则
	"required":        RuleRequired,
	"isdefault":       RuleDefault,
	"min":             RuleMin,
	"max":             RuleMax,
	"len":             RuleLen,
	"eq":              RuleEq,
	"eq_ignore_case":  RuleEqIgnoreCase,
	"ne":              RuleNe,
	"ne_ignore_case":  RuleNeIgnoreCase,
	"gt":              RuleGt,
	"gte":             RuleMin,
	"lt":              RuleLt,
	"lte":             RuleMax,
	"alpha":           noParamBuiltinAdapter(validate.StringAlpha),
	"alphaspace":      noParamBuiltinAdapter(validate.StringAlphaSpace),
	"alphanum":        noParamBuiltinAdapter(validate.StringAlphanum),
	"alphanumspace":   noParamBuiltinAdapter(validate.StringAlphanumSpace),
	"alphaunicode":    noParamBuiltinAdapter(validate.StringAlphaUnicode),
	"alphanumunicode": noParamBuiltinAdapter(validate.StringAlphanumUnicode),
	"ascii":           noParamBuiltinAdapter(validate.StringASCII),
	"printascii":      noParamBuiltinAdapter(validate.StringPrintASCII),
	"multibyte":       noParamBuiltinAdapter(validate.StringMultibyte),
	"hexadecimal":     noParamBuiltinAdapter(validate.StringHexadecimal),
	"unique":          RuleUnique,
	"boolean":         RuleBoolean,
	"number":          RuleNumber,
	"numeric":         RuleNumber,
	"json":            RuleJSON,
	"latitude":        RuleLatitude,
	"longitude":       RuleLongitude,

	// 以下规则可通过适配器自动生成
	"hexcolor":          noParamBuiltinAdapter(validate.StringHexColor),
	"rgb":               noParamBuiltinAdapter(validate.StringRGB),
	"rgba":              noParamBuiltinAdapter(validate.StringRGBA),
	"hsl":               noParamBuiltinAdapter(validate.StringHSL),
	"hsla":              noParamBuiltinAdapter(validate.StringHSLA),
	"email":             noParamBuiltinAdapter(validate.IsEmail),
	"e164":              noParamBuiltinAdapter(validate.StringE164),
	"ip":                noParamBuiltinAdapter(validate.StringIP),
	"ip_addr":           noParamBuiltinAdapter(validate.StringIP),
	"ipv4":              noParamBuiltinAdapter(validate.StringIPv4),
	"ipv6":              noParamBuiltinAdapter(validate.StringIPv6),
	"cidr":              noParamBuiltinAdapter(validate.StringCIDR),
	"cidrv4":            noParamBuiltinAdapter(validate.StringCIDRv4),
	"cidrv6":            noParamBuiltinAdapter(validate.StringCIDRv6),
	"mac":               noParamBuiltinAdapter(validate.StringMAC),
	"hostname":          noParamBuiltinAdapter(validate.StringHostname),
	"hostname_rfc1123":  noParamBuiltinAdapter(validate.StringHostname),
	"fqdn":              noParamBuiltinAdapter(validate.StringFQDN),
	"hostname_port":     noParamBuiltinAdapter(validate.StringHostnamePort),
	"port":              noParamScalarBuiltinAdapter(validate.StringPort),
	"url":               noParamBuiltinAdapter(validate.StringURL),
	"uri":               noParamBuiltinAdapter(validate.StringURI),
	"http_url":          noParamBuiltinAdapter(validate.StringHTTPURL),
	"https_url":         noParamBuiltinAdapter(validate.StringHTTPSURL),
	"url_encoded":       noParamBuiltinAdapter(validate.StringURLEncoded),
	"html":              noParamBuiltinAdapter(validate.StringHTML),
	"html_encoded":      noParamBuiltinAdapter(validate.StringHTMLEncoded),
	"uuid":              noParamBuiltinAdapter(validate.StringUUID),
	"uuid3":             noParamBuiltinAdapter(validate.StringUUID3),
	"uuid4":             noParamBuiltinAdapter(validate.StringUUID4),
	"uuid5":             noParamBuiltinAdapter(validate.StringUUID5),
	"uuid_rfc4122":      noParamBuiltinAdapter(validate.StringUUID),
	"uuid3_rfc4122":     noParamBuiltinAdapter(validate.StringUUID3),
	"uuid4_rfc4122":     noParamBuiltinAdapter(validate.StringUUID4),
	"uuid5_rfc4122":     noParamBuiltinAdapter(validate.StringUUID5),
	"base32":            noParamBuiltinAdapter(validate.StringBase32),
	"base64":            noParamBuiltinAdapter(validate.StringBase64),
	"base64url":         noParamBuiltinAdapter(validate.StringBase64URL),
	"base64rawurl":      noParamBuiltinAdapter(validate.StringBase64RawURL),
	"startswith":        stringBuiltinAdapter(validate.StringStartsWith),
	"endswith":          stringBuiltinAdapter(validate.StringEndsWith),
	"startsnotwith":     stringBuiltinAdapter(validate.StringStartsNotWith),
	"endsnotwith":       stringBuiltinAdapter(validate.StringEndsNotWith),
	"contains":          stringBuiltinAdapter(validate.StringContains),
	"containsany":       stringBuiltinAdapter(validate.StringContainsAny),
	"containsrune":      stringBuiltinAdapter(validate.StringContainsRune),
	"excludes":          stringBuiltinAdapter(validate.StringExcludes),
	"excludesall":       stringBuiltinAdapter(validate.StringExcludesAll),
	"excludesrune":      stringBuiltinAdapter(validate.StringExcludesRune),
	"lowercase":         noParamBuiltinAdapter(validate.StringLowercase),
	"uppercase":         noParamBuiltinAdapter(validate.StringUppercase),
	"datetime":          stringBuiltinAdapter(validate.StringDatetime),
	"timezone":          noParamBuiltinAdapter(validate.StringTimezone),
	"file":              noParamBuiltinAdapter(validate.StringFile),
	"filepath":          noParamBuiltinAdapter(validate.StringFilePath),
	"dir":               noParamBuiltinAdapter(validate.StringDir),
	"dirpath":           noParamBuiltinAdapter(validate.StringDirPath),
	"mongodb":           noParamBuiltinAdapter(validate.StringMongoDB),
	"luhn_checksum":     noParamBuiltinAdapter(validate.IsLuhnChecksum),
	"credit_card":       noParamBuiltinAdapter(validate.IsLuhnChecksum),
	"dns_rfc1035_label": noParamBuiltinAdapter(validate.StringDNSRFC1035Label),
	"semver":            noParamBuiltinAdapter(validate.IsSemver),
	"isbn10":            noParamBuiltinAdapter(validate.IsISBN10),
	"isbn13":            noParamBuiltinAdapter(validate.IsISBN13),
	"issn":              noParamBuiltinAdapter(validate.StringISSN),
	"bic":               noParamBuiltinAdapter(validate.StringBIC),
	"cron":              noParamBuiltinAdapter(validate.StringCron),
	"datauri":           noParamBuiltinAdapter(validate.StringDataURI),
	"bcp47":             noParamBuiltinAdapter(validate.StringBCP47),
	"eth_addr":          noParamBuiltinAdapter(validate.StringEthAddr),
	"btc_addr":          noParamBuiltinAdapter(validate.StringBtcAddr),
}
