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

	"github.com/kamalyes/go-argus/constants"
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
	return ok && CompareLengthOrNumber(field, n, constants.CmpGTE)
}

func RuleMax(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, constants.CmpLTE)
}

func RuleLen(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, constants.CmpEQ)
}

func RuleGt(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, constants.CmpGT)
}

func RuleLt(field reflect.Value, param string, _ bool) bool {
	n, ok := validate.ParseFloat(param)
	return ok && CompareLengthOrNumber(field, n, constants.CmpLT)
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
func CompareLengthOrNumber(field reflect.Value, expect float64, op constants.CmpOp) bool {
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
	constants.RuleRequired:        RuleRequired,
	constants.RuleIsDefault:       RuleDefault,
	constants.RuleMin:             RuleMin,
	constants.RuleMax:             RuleMax,
	constants.RuleLen:             RuleLen,
	constants.RuleEq:              RuleEq,
	constants.RuleEqIgnoreCase:    RuleEqIgnoreCase,
	constants.RuleNe:              RuleNe,
	constants.RuleNeIgnoreCase:    RuleNeIgnoreCase,
	constants.RuleGT:              RuleGt,
	constants.RuleGTE:             RuleMin,
	constants.RuleLT:              RuleLt,
	constants.RuleLTE:             RuleMax,
	constants.RuleAlpha:           noParamBuiltinAdapter(validate.StringAlpha),
	constants.RuleAlphaSpace:      noParamBuiltinAdapter(validate.StringAlphaSpace),
	constants.RuleAlphanum:        noParamBuiltinAdapter(validate.StringAlphanum),
	constants.RuleAlphanumSpace:   noParamBuiltinAdapter(validate.StringAlphanumSpace),
	constants.RuleAlphaUnicode:    noParamBuiltinAdapter(validate.StringAlphaUnicode),
	constants.RuleAlphanumUnicode: noParamBuiltinAdapter(validate.StringAlphanumUnicode),
	constants.RuleASCII:           noParamBuiltinAdapter(validate.StringASCII),
	constants.RulePrintASCII:      noParamBuiltinAdapter(validate.StringPrintASCII),
	constants.RuleMultibyte:       noParamBuiltinAdapter(validate.StringMultibyte),
	constants.RuleHexadecimal:     noParamBuiltinAdapter(validate.StringHexadecimal),
	constants.RuleUnique:          RuleUnique,
	constants.RuleBoolean:         RuleBoolean,
	constants.RuleNumber:          RuleNumber,
	constants.RuleNumeric:         RuleNumber,
	constants.RuleJSON:            RuleJSON,
	constants.RuleLatitude:        RuleLatitude,
	constants.RuleLongitude:       RuleLongitude,

	constants.RuleHexColor:        noParamBuiltinAdapter(validate.StringHexColor),
	constants.RuleRGB:             noParamBuiltinAdapter(validate.StringRGB),
	constants.RuleRGBA:            noParamBuiltinAdapter(validate.StringRGBA),
	constants.RuleHSL:             noParamBuiltinAdapter(validate.StringHSL),
	constants.RuleHSLA:            noParamBuiltinAdapter(validate.StringHSLA),
	constants.RuleEmail:           noParamBuiltinAdapter(validate.IsEmail),
	constants.RuleE164:            noParamBuiltinAdapter(validate.StringE164),
	constants.RuleIP:              noParamBuiltinAdapter(validate.StringIP),
	constants.RuleIPAddr:          noParamBuiltinAdapter(validate.StringIP),
	constants.RuleIPv4:            noParamBuiltinAdapter(validate.StringIPv4),
	constants.RuleIPv6:            noParamBuiltinAdapter(validate.StringIPv6),
	constants.RuleCIDR:            noParamBuiltinAdapter(validate.StringCIDR),
	constants.RuleCIDRv4:          noParamBuiltinAdapter(validate.StringCIDRv4),
	constants.RuleCIDRv6:          noParamBuiltinAdapter(validate.StringCIDRv6),
	constants.RuleMAC:             noParamBuiltinAdapter(validate.StringMAC),
	constants.RuleHostname:        noParamBuiltinAdapter(validate.StringHostname),
	constants.RuleHostnameRFC1123: noParamBuiltinAdapter(validate.StringHostname),
	constants.RuleFQDN:            noParamBuiltinAdapter(validate.StringFQDN),
	constants.RuleHostnamePort:    noParamBuiltinAdapter(validate.StringHostnamePort),
	constants.RulePort:            noParamScalarBuiltinAdapter(validate.StringPort),
	constants.RuleURL:             noParamBuiltinAdapter(validate.StringURL),
	constants.RuleURI:             noParamBuiltinAdapter(validate.StringURI),
	constants.RuleHTTPURL:         noParamBuiltinAdapter(validate.StringHTTPURL),
	constants.RuleHTTPSURL:        noParamBuiltinAdapter(validate.StringHTTPSURL),
	constants.RuleURLEncoded:      noParamBuiltinAdapter(validate.StringURLEncoded),
	constants.RuleHTML:            noParamBuiltinAdapter(validate.StringHTML),
	constants.RuleHTMLEncoded:     noParamBuiltinAdapter(validate.StringHTMLEncoded),
	constants.RuleUUID:            noParamBuiltinAdapter(validate.StringUUID),
	constants.RuleUUID3:           noParamBuiltinAdapter(validate.StringUUID3),
	constants.RuleUUID4:           noParamBuiltinAdapter(validate.StringUUID4),
	constants.RuleUUID5:           noParamBuiltinAdapter(validate.StringUUID5),
	constants.RuleUUIDRFC4122:     noParamBuiltinAdapter(validate.StringUUID),
	constants.RuleUUID3RFC4122:    noParamBuiltinAdapter(validate.StringUUID3),
	constants.RuleUUID4RFC4122:    noParamBuiltinAdapter(validate.StringUUID4),
	constants.RuleUUID5RFC4122:    noParamBuiltinAdapter(validate.StringUUID5),
	constants.RuleBase32:          noParamBuiltinAdapter(validate.StringBase32),
	constants.RuleBase64:          noParamBuiltinAdapter(validate.StringBase64),
	constants.RuleBase64URL:       noParamBuiltinAdapter(validate.StringBase64URL),
	constants.RuleBase64RawURL:    noParamBuiltinAdapter(validate.StringBase64RawURL),
	constants.RuleStartsWith:      stringBuiltinAdapter(validate.StringStartsWith),
	constants.RuleEndsWith:        stringBuiltinAdapter(validate.StringEndsWith),
	constants.RuleStartsNotWith:   stringBuiltinAdapter(validate.StringStartsNotWith),
	constants.RuleEndsNotWith:     stringBuiltinAdapter(validate.StringEndsNotWith),
	constants.RuleContains:        stringBuiltinAdapter(validate.StringContains),
	constants.RuleContainsAny:     stringBuiltinAdapter(validate.StringContainsAny),
	constants.RuleContainsRune:    stringBuiltinAdapter(validate.StringContainsRune),
	constants.RuleExcludes:        stringBuiltinAdapter(validate.StringExcludes),
	constants.RuleExcludesAll:     stringBuiltinAdapter(validate.StringExcludesAll),
	constants.RuleExcludesRune:    stringBuiltinAdapter(validate.StringExcludesRune),
	constants.RuleLowercase:       noParamBuiltinAdapter(validate.StringLowercase),
	constants.RuleUppercase:       noParamBuiltinAdapter(validate.StringUppercase),
	constants.RuleDatetime:        stringBuiltinAdapter(validate.StringDatetime),
	constants.RuleTimezone:        noParamBuiltinAdapter(validate.StringTimezone),
	constants.RuleFile:            noParamBuiltinAdapter(validate.StringFile),
	constants.RuleFilepath:        noParamBuiltinAdapter(validate.StringFilePath),
	constants.RuleDir:             noParamBuiltinAdapter(validate.StringDir),
	constants.RuleDirpath:         noParamBuiltinAdapter(validate.StringDirPath),
	constants.RuleMongoDB:         noParamBuiltinAdapter(validate.StringMongoDB),
	constants.RuleLuhnChecksum:    noParamBuiltinAdapter(validate.IsLuhnChecksum),
	constants.RuleCreditCard:      noParamBuiltinAdapter(validate.IsLuhnChecksum),
	constants.RuleDNSRFC1035Label: noParamBuiltinAdapter(validate.StringDNSRFC1035Label),
	constants.RuleSemver:          noParamBuiltinAdapter(validate.IsSemver),
	constants.RuleISBN10:          noParamBuiltinAdapter(validate.IsISBN10),
	constants.RuleISBN13:          noParamBuiltinAdapter(validate.IsISBN13),
	constants.RuleISSN:            noParamBuiltinAdapter(validate.StringISSN),
	constants.RuleBIC:             noParamBuiltinAdapter(validate.StringBIC),
	constants.RuleCron:            noParamBuiltinAdapter(validate.StringCron),
	constants.RuleDataURI:         noParamBuiltinAdapter(validate.StringDataURI),
	constants.RuleBCP47:           noParamBuiltinAdapter(validate.StringBCP47),
	constants.RuleEthAddr:         noParamBuiltinAdapter(validate.StringEthAddr),
	constants.RuleBtcAddr:         noParamBuiltinAdapter(validate.StringBtcAddr),
}
