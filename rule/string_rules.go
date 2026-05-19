/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-19 00:00:00
 * @FilePath: \go-argus\rule\string_rules.go
 * @Description: 字符串规则映射，薄层委托 validate 子包
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"github.com/kamalyes/go-argus/validate"
)

type StringRuleFunc func(s string, param string) bool

var StringRuleMap = map[string]StringRuleFunc{
	"required":          StringRuleRequired,
	"isdefault":         StringRuleIsDefault,
	"min":               StringRuleMin,
	"max":               StringRuleMax,
	"len":               StringRuleLen,
	"eq":                StringRuleEq,
	"eq_ignore_case":    StringRuleEqIgnoreCase,
	"ne":                StringRuleNe,
	"ne_ignore_case":    StringRuleNeIgnoreCase,
	"gt":                StringRuleGt,
	"gte":               StringRuleGte,
	"lt":                StringRuleLt,
	"lte":               StringRuleLte,
	"alpha":             StringRuleAlpha,
	"alphaspace":        StringRuleAlphaSpace,
	"alphanum":          StringRuleAlphanum,
	"alphanumspace":     StringRuleAlphanumSpace,
	"alphaunicode":      StringRuleAlphaUnicode,
	"alphanumunicode":   StringRuleAlphanumUnicode,
	"ascii":             StringRuleASCII,
	"printascii":        StringRulePrintASCII,
	"multibyte":         StringRuleMultibyte,
	"hexadecimal":       StringRuleHexadecimal,
	"hexcolor":          StringRuleHexColor,
	"rgb":               StringRuleRGB,
	"rgba":              StringRuleRGBA,
	"hsl":               StringRuleHSL,
	"hsla":              StringRuleHSLA,
	"email":             StringRuleEmail,
	"e164":              StringRuleE164,
	"ip":                StringRuleIP,
	"ip_addr":           StringRuleIP,
	"ipv4":              StringRuleIPv4,
	"ipv6":              StringRuleIPv6,
	"cidr":              StringRuleCIDR,
	"cidrv4":            StringRuleCIDRv4,
	"cidrv6":            StringRuleCIDRv6,
	"mac":               StringRuleMAC,
	"hostname":          StringRuleHostname,
	"hostname_rfc1123":  StringRuleHostname,
	"fqdn":              StringRuleFQDN,
	"hostname_port":     StringRuleHostnamePort,
	"port":              StringRulePort,
	"url":               StringRuleURL,
	"uri":               StringRuleURI,
	"http_url":          StringRuleHTTPURL,
	"https_url":         StringRuleHTTPSURL,
	"url_encoded":       StringRuleURLEncoded,
	"html":              StringRuleHTML,
	"html_encoded":      StringRuleHTMLEncoded,
	"uuid":              StringRuleUUID,
	"uuid3":             StringRuleUUID3,
	"uuid4":             StringRuleUUID4,
	"uuid5":             StringRuleUUID5,
	"uuid_rfc4122":      StringRuleUUID,
	"uuid3_rfc4122":     StringRuleUUID3,
	"uuid4_rfc4122":     StringRuleUUID4,
	"uuid5_rfc4122":     StringRuleUUID5,
	"base32":            StringRuleBase32,
	"base64":            StringRuleBase64,
	"base64url":         StringRuleBase64URL,
	"base64rawurl":      StringRuleBase64RawURL,
	"json":              StringRuleJSON,
	"unique":            StringRuleUnique,
	"startswith":        StringRuleStartsWith,
	"endswith":          StringRuleEndsWith,
	"startsnotwith":     StringRuleStartsNotWith,
	"endsnotwith":       StringRuleEndsNotWith,
	"contains":          StringRuleContains,
	"containsany":       StringRuleContainsAny,
	"containsrune":      StringRuleContainsRune,
	"excludes":          StringRuleExcludes,
	"excludesall":       StringRuleExcludesAll,
	"excludesrune":      StringRuleExcludesRune,
	"lowercase":         StringRuleLowercase,
	"uppercase":         StringRuleUppercase,
	"boolean":           StringRuleBoolean,
	"number":            StringRuleNumber,
	"numeric":           StringRuleNumber,
	"datetime":          StringRuleDatetime,
	"timezone":          StringRuleTimezone,
	"latitude":          StringRuleLatitude,
	"longitude":         StringRuleLongitude,
	"file":              StringRuleFile,
	"filepath":          StringRuleFilePath,
	"dir":               StringRuleDir,
	"dirpath":           StringRuleDirPath,
	"mongodb":           StringRuleMongoDB,
	"luhn_checksum":     StringRuleLuhnChecksum,
	"credit_card":       StringRuleLuhnChecksum,
	"dns_rfc1035_label": StringRuleDNSRFC1035Label,
	"semver":            StringRuleSemver,
	"isbn10":            StringRuleISBN10,
	"isbn13":            StringRuleISBN13,
	"issn":              StringRuleISSN,
	"bic":               StringRuleBIC,
	"cron":              StringRuleCron,
	"datauri":           StringRuleDataURI,
	"bcp47":             StringRuleBCP47,
	"eth_addr":          StringRuleEthAddr,
	"btc_addr":          StringRuleBtcAddr,
}

func StringRuleRequired(s string, _ string) bool        { return validate.StringRequired(s) }
func StringRuleIsDefault(s string, _ string) bool       { return validate.StringIsDefault(s) }
func StringRuleMin(s string, param string) bool         { return validate.StringMin(s, param) }
func StringRuleMax(s string, param string) bool         { return validate.StringMax(s, param) }
func StringRuleLen(s string, param string) bool         { return validate.StringLen(s, param) }
func StringRuleEq(s string, param string) bool          { return validate.StringEq(s, param) }
func StringRuleEqIgnoreCase(s string, param string) bool { return validate.StringEqIgnoreCase(s, param) }
func StringRuleNe(s string, param string) bool          { return validate.StringNe(s, param) }
func StringRuleNeIgnoreCase(s string, param string) bool { return validate.StringNeIgnoreCase(s, param) }
func StringRuleGt(s string, param string) bool          { return validate.StringGt(s, param) }
func StringRuleGte(s string, param string) bool         { return validate.StringGte(s, param) }
func StringRuleLt(s string, param string) bool          { return validate.StringLt(s, param) }
func StringRuleLte(s string, param string) bool         { return validate.StringLte(s, param) }
func StringRuleAlpha(s string, _ string) bool           { return validate.StringAlpha(s) }
func StringRuleAlphaSpace(s string, _ string) bool      { return validate.StringAlphaSpace(s) }
func StringRuleAlphanum(s string, _ string) bool        { return validate.StringAlphanum(s) }
func StringRuleAlphanumSpace(s string, _ string) bool   { return validate.StringAlphanumSpace(s) }
func StringRuleAlphaUnicode(s string, _ string) bool    { return validate.StringAlphaUnicode(s) }
func StringRuleAlphanumUnicode(s string, _ string) bool { return validate.StringAlphanumUnicode(s) }
func StringRuleASCII(s string, _ string) bool           { return validate.StringASCII(s) }
func StringRulePrintASCII(s string, _ string) bool      { return validate.StringPrintASCII(s) }
func StringRuleMultibyte(s string, _ string) bool       { return validate.StringMultibyte(s) }
func StringRuleHexadecimal(s string, _ string) bool     { return validate.StringHexadecimal(s) }
func StringRuleHexColor(s string, _ string) bool        { return validate.StringHexColor(s) }
func StringRuleRGB(s string, _ string) bool             { return validate.StringRGB(s) }
func StringRuleRGBA(s string, _ string) bool            { return validate.StringRGBA(s) }
func StringRuleHSL(s string, _ string) bool             { return validate.StringHSL(s) }
func StringRuleHSLA(s string, _ string) bool            { return validate.StringHSLA(s) }
func StringRuleEmail(s string, _ string) bool           { return validate.IsEmail(s) }
func StringRuleE164(s string, _ string) bool            { return validate.StringE164(s) }
func StringRuleIP(s string, _ string) bool              { return validate.StringIP(s) }
func StringRuleIPv4(s string, _ string) bool            { return validate.StringIPv4(s) }
func StringRuleIPv6(s string, _ string) bool            { return validate.StringIPv6(s) }
func StringRuleCIDR(s string, _ string) bool            { return validate.StringCIDR(s) }
func StringRuleCIDRv4(s string, _ string) bool          { return validate.StringCIDRv4(s) }
func StringRuleCIDRv6(s string, _ string) bool          { return validate.StringCIDRv6(s) }
func StringRuleMAC(s string, _ string) bool             { return validate.StringMAC(s) }
func StringRuleHostname(s string, _ string) bool        { return validate.StringHostname(s) }
func StringRuleFQDN(s string, _ string) bool            { return validate.StringFQDN(s) }
func StringRuleHostnamePort(s string, _ string) bool    { return validate.StringHostnamePort(s) }
func StringRulePort(s string, _ string) bool            { return validate.StringPort(s) }
func StringRuleURL(s string, _ string) bool             { return validate.StringURL(s) }
func StringRuleURI(s string, _ string) bool             { return validate.StringURI(s) }
func StringRuleHTTPURL(s string, _ string) bool         { return validate.StringHTTPURL(s) }
func StringRuleHTTPSURL(s string, _ string) bool        { return validate.StringHTTPSURL(s) }
func StringRuleURLEncoded(s string, _ string) bool      { return validate.StringURLEncoded(s) }
func StringRuleHTML(s string, _ string) bool            { return validate.StringHTML(s) }
func StringRuleHTMLEncoded(s string, _ string) bool     { return validate.StringHTMLEncoded(s) }
func StringRuleUUID(s string, _ string) bool            { return validate.StringUUID(s) }
func StringRuleUUID3(s string, _ string) bool           { return validate.StringUUID3(s) }
func StringRuleUUID4(s string, _ string) bool           { return validate.StringUUID4(s) }
func StringRuleUUID5(s string, _ string) bool           { return validate.StringUUID5(s) }
func StringRuleBase32(s string, _ string) bool          { return validate.StringBase32(s) }
func StringRuleBase64(s string, _ string) bool          { return validate.StringBase64(s) }
func StringRuleBase64URL(s string, _ string) bool       { return validate.StringBase64URL(s) }
func StringRuleBase64RawURL(s string, _ string) bool    { return validate.StringBase64RawURL(s) }
func StringRuleJSON(s string, _ string) bool            { return validate.StringJSON(s) }
func StringRuleUnique(s string, _ string) bool          { return validate.StringUnique(s) }
func StringRuleStartsWith(s string, param string) bool  { return validate.StringStartsWith(s, param) }
func StringRuleEndsWith(s string, param string) bool    { return validate.StringEndsWith(s, param) }
func StringRuleStartsNotWith(s string, param string) bool { return validate.StringStartsNotWith(s, param) }
func StringRuleEndsNotWith(s string, param string) bool { return validate.StringEndsNotWith(s, param) }
func StringRuleContains(s string, param string) bool    { return validate.StringContains(s, param) }
func StringRuleContainsAny(s string, param string) bool { return validate.StringContainsAny(s, param) }
func StringRuleContainsRune(s string, param string) bool { return validate.StringContainsRune(s, param) }
func StringRuleExcludes(s string, param string) bool    { return validate.StringExcludes(s, param) }
func StringRuleExcludesAll(s string, param string) bool { return validate.StringExcludesAll(s, param) }
func StringRuleExcludesRune(s string, param string) bool { return validate.StringExcludesRune(s, param) }
func StringRuleLowercase(s string, _ string) bool       { return validate.StringLowercase(s) }
func StringRuleUppercase(s string, _ string) bool       { return validate.StringUppercase(s) }
func StringRuleBoolean(s string, _ string) bool         { return validate.StringBoolean(s) }
func StringRuleNumber(s string, _ string) bool          { return validate.StringNumber(s) }
func StringRuleDatetime(s string, param string) bool    { return validate.StringDatetime(s, param) }
func StringRuleTimezone(s string, _ string) bool        { return validate.StringTimezone(s) }
func StringRuleLatitude(s string, _ string) bool        { return validate.StringLatitude(s) }
func StringRuleLongitude(s string, _ string) bool       { return validate.StringLongitude(s) }
func StringRuleFile(s string, _ string) bool            { return validate.StringFile(s) }
func StringRuleFilePath(s string, _ string) bool        { return validate.StringFilePath(s) }
func StringRuleDir(s string, _ string) bool             { return validate.StringDir(s) }
func StringRuleDirPath(s string, _ string) bool         { return validate.StringDirPath(s) }
func StringRuleMongoDB(s string, _ string) bool         { return validate.StringMongoDB(s) }
func StringRuleLuhnChecksum(s string, _ string) bool    { return validate.IsLuhnChecksum(s) }
func StringRuleDNSRFC1035Label(s string, _ string) bool { return validate.StringDNSRFC1035Label(s) }
func StringRuleSemver(s string, _ string) bool          { return validate.IsSemver(s) }
func StringRuleISBN10(s string, _ string) bool          { return validate.IsISBN10(s) }
func StringRuleISBN13(s string, _ string) bool          { return validate.IsISBN13(s) }
func StringRuleISSN(s string, _ string) bool            { return validate.IsISSN(s) }
func StringRuleBIC(s string, _ string) bool             { return validate.IsBIC(s) }
func StringRuleCron(s string, _ string) bool            { return validate.IsCron(s) }
func StringRuleDataURI(s string, _ string) bool         { return validate.IsDataURI(s) }
func StringRuleBCP47(s string, _ string) bool           { return validate.IsBCP47(s) }
func StringRuleEthAddr(s string, _ string) bool         { return validate.IsEthAddr(s) }
func StringRuleBtcAddr(s string, _ string) bool         { return validate.IsBtcAddr(s) }
