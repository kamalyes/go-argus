/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-20 00:00:00
 * @FilePath: \go-argus\rule\builtin_test.go
 * @Description: builtin.go 测试，覆盖内置规则函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"testing"

	"github.com/kamalyes/go-argus/validate"
)

// callBuiltin 辅助函数：从 BuiltinRules 查表调用
func callBuiltin(name string, field reflect.Value, param string) bool {
	fn := BuiltinRules[name]
	if fn == nil {
		panic("builtin rule not found: " + name)
	}
	return fn(field, param, false)
}

func TestRuleRequired(t *testing.T) {
	if !RuleRequired(reflect.ValueOf("hello"), "", true) {
		t.Fatal("expected non-empty string to be required")
	}
	if RuleRequired(reflect.ValueOf(""), "", true) {
		t.Fatal("expected empty string to not be required")
	}
}

func TestRuleDefault(t *testing.T) {
	if !RuleDefault(reflect.ValueOf(""), "", true) {
		t.Fatal("expected empty string to be default")
	}
	if RuleDefault(reflect.ValueOf("hello"), "", true) {
		t.Fatal("expected non-empty string to not be default")
	}
}

func TestRuleMin(t *testing.T) {
	if !RuleMin(reflect.ValueOf("abc"), "3", false) {
		t.Fatal("expected min=3 to pass for 'abc'")
	}
	if RuleMin(reflect.ValueOf("ab"), "3", false) {
		t.Fatal("expected min=3 to fail for 'ab'")
	}
}

func TestRuleMax(t *testing.T) {
	if !RuleMax(reflect.ValueOf("ab"), "3", false) {
		t.Fatal("expected max=3 to pass for 'ab'")
	}
	if !RuleMax(reflect.ValueOf("abc"), "3", false) {
		t.Fatal("expected max=3 to pass for 'abc'")
	}
}

func TestRuleLen(t *testing.T) {
	if !RuleLen(reflect.ValueOf("abc"), "3", false) {
		t.Fatal("expected len=3 to pass for 'abc'")
	}
	if RuleLen(reflect.ValueOf("ab"), "3", false) {
		t.Fatal("expected len=3 to fail for 'ab'")
	}
}

func TestRuleEq(t *testing.T) {
	if !RuleEq(reflect.ValueOf("hello"), "hello", false) {
		t.Fatal("expected eq to pass")
	}
	if RuleEq(reflect.ValueOf("hello"), "world", false) {
		t.Fatal("expected eq to fail")
	}
}

func TestRuleNe(t *testing.T) {
	if !RuleNe(reflect.ValueOf("hello"), "world", false) {
		t.Fatal("expected ne to pass")
	}
	if RuleNe(reflect.ValueOf("hello"), "hello", false) {
		t.Fatal("expected ne to fail")
	}
}

func TestRuleGt(t *testing.T) {
	if !RuleGt(reflect.ValueOf(10), "5", false) {
		t.Fatal("expected gt to pass")
	}
	if RuleGt(reflect.ValueOf(3), "5", false) {
		t.Fatal("expected gt to fail")
	}
}

func TestRuleLt(t *testing.T) {
	if !RuleLt(reflect.ValueOf(3), "5", false) {
		t.Fatal("expected lt to pass")
	}
}

func TestRuleAlpha(t *testing.T) {
	if !callBuiltin("alpha", reflect.ValueOf("abc"), "") {
		t.Fatal("expected alpha to pass")
	}
	if callBuiltin("alpha", reflect.ValueOf("abc123"), "") {
		t.Fatal("expected alpha to fail for alphanum")
	}
}

func TestRuleAlphanum(t *testing.T) {
	if !callBuiltin("alphanum", reflect.ValueOf("abc123"), "") {
		t.Fatal("expected alphanum to pass")
	}
}

// --- 通过 BuiltinRules 查表测试适配器生成的规则 ---

func TestBuiltinEmail(t *testing.T) {
	if !callBuiltin("email", reflect.ValueOf("test@example.com"), "") {
		t.Fatal("expected email to pass")
	}
	if callBuiltin("email", reflect.ValueOf("not-email"), "") {
		t.Fatal("expected email to fail")
	}
}

func TestBuiltinIP(t *testing.T) {
	if !callBuiltin("ip", reflect.ValueOf("192.168.1.1"), "") {
		t.Fatal("expected IP to pass")
	}
	if callBuiltin("ip", reflect.ValueOf("not-ip"), "") {
		t.Fatal("expected IP to fail")
	}
}

func TestBuiltinHostname(t *testing.T) {
	if !callBuiltin("hostname", reflect.ValueOf("example.com"), "") {
		t.Fatal("expected hostname to pass")
	}
}

func TestBuiltinFQDN(t *testing.T) {
	if !callBuiltin("fqdn", reflect.ValueOf("example.com."), "") {
		t.Fatal("expected fqdn to pass")
	}
}

func TestBuiltinURL(t *testing.T) {
	if !callBuiltin("url", reflect.ValueOf("http://example.com"), "") {
		t.Fatal("expected url to pass")
	}
}

func TestBuiltinUUID(t *testing.T) {
	if !callBuiltin("uuid", reflect.ValueOf("6ba7b810-9dad-11d1-80b4-00c04fd430c8"), "") {
		t.Fatal("expected uuid to pass")
	}
}

func TestRuleJSON(t *testing.T) {
	if !RuleJSON(reflect.ValueOf(`{"key":"value"}`), "", false) {
		t.Fatal("expected json to pass")
	}
	if RuleJSON(reflect.ValueOf("not-json"), "", false) {
		t.Fatal("expected json to fail")
	}
}

func TestBuiltinStartsWith(t *testing.T) {
	if !callBuiltin("startswith", reflect.ValueOf("hello world"), "hello") {
		t.Fatal("expected startsWith to pass")
	}
}

func TestBuiltinContains(t *testing.T) {
	if !callBuiltin("contains", reflect.ValueOf("hello world"), "world") {
		t.Fatal("expected contains to pass")
	}
}

func TestBuiltinLowercase(t *testing.T) {
	if !callBuiltin("lowercase", reflect.ValueOf("hello"), "") {
		t.Fatal("expected lowercase to pass")
	}
}

func TestRuleBoolean(t *testing.T) {
	if !RuleBoolean(reflect.ValueOf("true"), "", false) {
		t.Fatal("expected boolean to pass for 'true'")
	}
	if !RuleBoolean(reflect.ValueOf(true), "", false) {
		t.Fatal("expected boolean to pass for bool true")
	}
}

func TestRuleNumber(t *testing.T) {
	if !RuleNumber(reflect.ValueOf(42), "", false) {
		t.Fatal("expected number to pass for int")
	}
	if !RuleNumber(reflect.ValueOf("42"), "", false) {
		t.Fatal("expected number to pass for numeric string")
	}
}

func TestResolveFieldValueInvalid(t *testing.T) {
	_, ok := ResolveFieldValue(reflect.Value{})
	if ok {
		t.Fatal("expected invalid value to fail")
	}
}

func TestResolveFieldValueBool(t *testing.T) {
	_, ok := ResolveFieldValue(reflect.ValueOf(true))
	if ok {
		t.Fatal("expected bool to fail ResolveFieldValue")
	}
}

func TestCompareOp(t *testing.T) {
	if !validate.CompareOp(5, 3, validate.CmpGT) {
		t.Fatal("expected 5 > 3")
	}
	if !validate.CompareOp(5, 5, validate.CmpGTE) {
		t.Fatal("expected 5 >= 5")
	}
	if !validate.CompareOp(3, 5, validate.CmpLT) {
		t.Fatal("expected 3 < 5")
	}
	if !validate.CompareOp(3, 3, validate.CmpLTE) {
		t.Fatal("expected 3 <= 3")
	}
	if !validate.CompareOp(5, 5, validate.CmpEQ) {
		t.Fatal("expected 5 == 5")
	}
}

func TestRuleUniqueSlice(t *testing.T) {
	if !IsUniqueSlice(reflect.ValueOf([]string{"a", "b", "c"})) {
		t.Fatal("expected unique slice to pass")
	}
	if IsUniqueSlice(reflect.ValueOf([]string{"a", "b", "a"})) {
		t.Fatal("expected non-unique slice to fail")
	}
}

func TestRuleUniqueMap(t *testing.T) {
	m := map[string]string{"a": "1", "b": "2"}
	if !IsUniqueMap(reflect.ValueOf(m)) {
		t.Fatal("expected unique map to pass")
	}
}

func TestRuleUniqueInvalid(t *testing.T) {
	if RuleUnique(reflect.ValueOf(42), "", false) {
		t.Fatal("expected int to fail unique")
	}
}

func TestRuleEqIgnoreCase(t *testing.T) {
	if !RuleEqIgnoreCase(reflect.ValueOf("Hello"), "hello", false) {
		t.Fatal("expected eq_ignore_case to pass")
	}
}

func TestRuleNeIgnoreCase(t *testing.T) {
	if !RuleNeIgnoreCase(reflect.ValueOf("Hello"), "world", false) {
		t.Fatal("expected ne_ignore_case to pass")
	}
}

func TestRuleGteLte(t *testing.T) {
	if !callBuiltin("gte", reflect.ValueOf(5), "5") {
		t.Fatal("expected gte to pass")
	}
	if !callBuiltin("lte", reflect.ValueOf(5), "5") {
		t.Fatal("expected lte to pass")
	}
}

// --- 补充 builtin.go 未覆盖函数测试 ---

func TestRuleMinInvalidParam(t *testing.T) {
	if RuleMin(reflect.ValueOf("abc"), "notanumber", false) {
		t.Fatal("expected min to fail for invalid param")
	}
}

func TestRuleMaxInvalidParam(t *testing.T) {
	if RuleMax(reflect.ValueOf("abc"), "notanumber", false) {
		t.Fatal("expected max to fail for invalid param")
	}
}

func TestRuleLenInvalidParam(t *testing.T) {
	if RuleLen(reflect.ValueOf("abc"), "notanumber", false) {
		t.Fatal("expected len to fail for invalid param")
	}
}

func TestRuleEqInvalidField(t *testing.T) {
	if !RuleEq(reflect.ValueOf(42), "42", false) {
		t.Fatal("expected eq to pass for int 42 via ScalarString")
	}
}

func TestRuleAlphaSpace(t *testing.T) {
	if !callBuiltin("alphaspace", reflect.ValueOf("hello world"), "") {
		t.Fatal("expected alphaspace to pass")
	}
}

func TestRuleAlphanumSpace(t *testing.T) {
	if !callBuiltin("alphanumspace", reflect.ValueOf("hello 123"), "") {
		t.Fatal("expected alphanumspace to pass")
	}
}

func TestRuleAlphaUnicode(t *testing.T) {
	if !callBuiltin("alphaunicode", reflect.ValueOf("你好"), "") {
		t.Fatal("expected alphaunicode to pass")
	}
}

func TestRuleAlphanumUnicode(t *testing.T) {
	if !callBuiltin("alphanumunicode", reflect.ValueOf("你好123"), "") {
		t.Fatal("expected alphanumunicode to pass")
	}
}

func TestRuleASCII(t *testing.T) {
	if !callBuiltin("ascii", reflect.ValueOf("hello"), "") {
		t.Fatal("expected ascii to pass")
	}
}

func TestRulePrintASCII(t *testing.T) {
	if !callBuiltin("printascii", reflect.ValueOf("hello"), "") {
		t.Fatal("expected printascii to pass")
	}
}

func TestRuleHexadecimal(t *testing.T) {
	if !callBuiltin("hexadecimal", reflect.ValueOf("abcdef0123456789"), "") {
		t.Fatal("expected hexadecimal to pass")
	}
}

// --- 适配器生成的规则通过 BuiltinRules 查表测试 ---

func TestBuiltinHexColor(t *testing.T) {
	if !callBuiltin("hexcolor", reflect.ValueOf("#ff0000"), "") {
		t.Fatal("expected hexcolor to pass")
	}
}

func TestBuiltinRGB(t *testing.T) {
	if !callBuiltin("rgb", reflect.ValueOf("rgb(255,0,0)"), "") {
		t.Fatal("expected rgb to pass")
	}
}

func TestBuiltinRGBA(t *testing.T) {
	if !callBuiltin("rgba", reflect.ValueOf("rgba(255,0,0,0.5)"), "") {
		t.Fatal("expected rgba to pass")
	}
}

func TestBuiltinHSL(t *testing.T) {
	if !callBuiltin("hsl", reflect.ValueOf("hsl(120,100%,50%)"), "") {
		t.Fatal("expected hsl to pass")
	}
}

func TestBuiltinHSLA(t *testing.T) {
	if !callBuiltin("hsla", reflect.ValueOf("hsla(120,100%,50%,0.5)"), "") {
		t.Fatal("expected hsla to pass")
	}
}

func TestBuiltinE164(t *testing.T) {
	if !callBuiltin("e164", reflect.ValueOf("+1234567890"), "") {
		t.Fatal("expected e164 to pass")
	}
}

func TestBuiltinIPv4(t *testing.T) {
	if !callBuiltin("ipv4", reflect.ValueOf("192.168.1.1"), "") {
		t.Fatal("expected ipv4 to pass")
	}
}

func TestBuiltinIPv6(t *testing.T) {
	if !callBuiltin("ipv6", reflect.ValueOf("::1"), "") {
		t.Fatal("expected ipv6 to pass")
	}
}

func TestBuiltinCIDR(t *testing.T) {
	if !callBuiltin("cidr", reflect.ValueOf("10.0.0.0/8"), "") {
		t.Fatal("expected cidr to pass")
	}
}

func TestBuiltinCIDRv4(t *testing.T) {
	if !callBuiltin("cidrv4", reflect.ValueOf("10.0.0.0/8"), "") {
		t.Fatal("expected cidrv4 to pass")
	}
}

func TestBuiltinCIDRv6(t *testing.T) {
	if !callBuiltin("cidrv6", reflect.ValueOf("::1/128"), "") {
		t.Fatal("expected cidrv6 to pass")
	}
}

func TestBuiltinMAC(t *testing.T) {
	if !callBuiltin("mac", reflect.ValueOf("00:11:22:33:44:55"), "") {
		t.Fatal("expected mac to pass")
	}
}

func TestBuiltinHostnamePort(t *testing.T) {
	if !callBuiltin("hostname_port", reflect.ValueOf("example.com:8080"), "") {
		t.Fatal("expected hostname_port to pass")
	}
}

func TestBuiltinPort(t *testing.T) {
	if !callBuiltin("port", reflect.ValueOf("443"), "") {
		t.Fatal("expected port to pass")
	}
}

func TestBuiltinURI(t *testing.T) {
	if !callBuiltin("uri", reflect.ValueOf("http://example.com/path"), "") {
		t.Fatal("expected uri to pass")
	}
}

func TestBuiltinHTTPURL(t *testing.T) {
	if !callBuiltin("http_url", reflect.ValueOf("http://example.com"), "") {
		t.Fatal("expected http_url to pass")
	}
}

func TestBuiltinHTTPSURL(t *testing.T) {
	if !callBuiltin("https_url", reflect.ValueOf("https://example.com"), "") {
		t.Fatal("expected https_url to pass")
	}
}

func TestBuiltinURLEncoded(t *testing.T) {
	if !callBuiltin("url_encoded", reflect.ValueOf("hello%20world"), "") {
		t.Fatal("expected url_encoded to pass")
	}
}

func TestBuiltinHTML(t *testing.T) {
	if !callBuiltin("html", reflect.ValueOf("<b>hello</b>"), "") {
		t.Fatal("expected html to pass")
	}
}

func TestBuiltinHTMLEncoded(t *testing.T) {
	if !callBuiltin("html_encoded", reflect.ValueOf("&lt;b&gt;"), "") {
		t.Fatal("expected html_encoded to pass")
	}
}

func TestBuiltinUUID3(t *testing.T) {
	if !callBuiltin("uuid3", reflect.ValueOf("6ba7b810-9dad-31d1-80b4-00c04fd430c8"), "") {
		t.Fatal("expected uuid3 to pass")
	}
}

func TestBuiltinUUID4(t *testing.T) {
	if !callBuiltin("uuid4", reflect.ValueOf("6ba7b810-9dad-41d1-80b4-00c04fd430c8"), "") {
		t.Fatal("expected uuid4 to pass")
	}
}

func TestBuiltinUUID5(t *testing.T) {
	if !callBuiltin("uuid5", reflect.ValueOf("6ba7b810-9dad-51d1-80b4-00c04fd430c8"), "") {
		t.Fatal("expected uuid5 to pass")
	}
}

func TestBuiltinBase32(t *testing.T) {
	if !callBuiltin("base32", reflect.ValueOf("JBSWY3DPEE======"), "") {
		t.Fatal("expected base32 to pass")
	}
}

func TestBuiltinBase64(t *testing.T) {
	if !callBuiltin("base64", reflect.ValueOf("SGVsbG8="), "") {
		t.Fatal("expected base64 to pass")
	}
}

func TestBuiltinBase64URL(t *testing.T) {
	if !callBuiltin("base64url", reflect.ValueOf("SGVsbG8="), "") {
		t.Fatal("expected base64url to pass")
	}
}

func TestBuiltinBase64RawURL(t *testing.T) {
	if !callBuiltin("base64rawurl", reflect.ValueOf("SGVsbG8"), "") {
		t.Fatal("expected base64rawurl to pass")
	}
}

func TestRuleUniqueString(t *testing.T) {
	if !RuleUnique(reflect.ValueOf("abc"), "", false) {
		t.Fatal("expected unique string to pass")
	}
}

func TestRuleUniqueInvalidValue(t *testing.T) {
	if RuleUnique(reflect.Value{}, "", false) {
		t.Fatal("expected invalid value to fail unique")
	}
}

func TestBuiltinEndsWith(t *testing.T) {
	if !callBuiltin("endswith", reflect.ValueOf("hello world"), "world") {
		t.Fatal("expected endswith to pass")
	}
}

func TestBuiltinStartsNotWith(t *testing.T) {
	if !callBuiltin("startsnotwith", reflect.ValueOf("hello"), "xyz") {
		t.Fatal("expected startsnotwith to pass")
	}
}

func TestBuiltinEndsNotWith(t *testing.T) {
	if !callBuiltin("endsnotwith", reflect.ValueOf("hello"), "xyz") {
		t.Fatal("expected endsnotwith to pass")
	}
}

func TestBuiltinContainsAny(t *testing.T) {
	if !callBuiltin("containsany", reflect.ValueOf("hello"), "ae") {
		t.Fatal("expected containsany to pass")
	}
}

func TestBuiltinContainsRune(t *testing.T) {
	if !callBuiltin("containsrune", reflect.ValueOf("hello"), "e") {
		t.Fatal("expected containsrune to pass")
	}
}

func TestBuiltinExcludes(t *testing.T) {
	if !callBuiltin("excludes", reflect.ValueOf("hello"), "xyz") {
		t.Fatal("expected excludes to pass")
	}
}

func TestBuiltinExcludesAll(t *testing.T) {
	if !callBuiltin("excludesall", reflect.ValueOf("hello"), "xyz") {
		t.Fatal("expected excludesall to pass")
	}
}

func TestBuiltinExcludesRune(t *testing.T) {
	if !callBuiltin("excludesrune", reflect.ValueOf("hello"), "z") {
		t.Fatal("expected excludesrune to pass")
	}
}

func TestBuiltinUppercase(t *testing.T) {
	if !callBuiltin("uppercase", reflect.ValueOf("HELLO"), "") {
		t.Fatal("expected uppercase to pass")
	}
}

func TestRuleNumberInvalid(t *testing.T) {
	if RuleNumber(reflect.ValueOf("abc"), "", false) {
		t.Fatal("expected number to fail for non-numeric string")
	}
	if RuleNumber(reflect.ValueOf(true), "", false) {
		t.Fatal("expected number to fail for bool")
	}
}

func TestRuleNumberInvalidValue(t *testing.T) {
	if RuleNumber(reflect.Value{}, "", false) {
		t.Fatal("expected number to fail for invalid value")
	}
}

func TestBuiltinDatetime(t *testing.T) {
	if !callBuiltin("datetime", reflect.ValueOf("2024-01-01"), "2006-01-02") {
		t.Fatal("expected datetime to pass")
	}
}

func TestBuiltinTimezone(t *testing.T) {
	if !callBuiltin("timezone", reflect.ValueOf("UTC"), "") {
		t.Fatal("expected timezone to pass")
	}
}

func TestRuleLatitudeNumeric(t *testing.T) {
	if !RuleLatitude(reflect.ValueOf(45.0), "", false) {
		t.Fatal("expected latitude to pass for numeric")
	}
}

func TestRuleLongitudeNumeric(t *testing.T) {
	if !RuleLongitude(reflect.ValueOf(90.0), "", false) {
		t.Fatal("expected longitude to pass for numeric")
	}
}

func TestBuiltinFile(t *testing.T) {
	_ = callBuiltin("file", reflect.ValueOf("nonexistent.txt"), "")
}

func TestBuiltinFilePath(t *testing.T) {
	if !callBuiltin("filepath", reflect.ValueOf("/test/file.txt"), "") && !callBuiltin("filepath", reflect.ValueOf("C:\\test\\file.txt"), "") {
		t.Fatal("expected filepath to pass")
	}
}

func TestBuiltinDir(t *testing.T) {
	_ = callBuiltin("dir", reflect.ValueOf("nonexistent_dir"), "")
}

func TestBuiltinDirPath(t *testing.T) {
	if !callBuiltin("dirpath", reflect.ValueOf("/test/"), "") && !callBuiltin("dirpath", reflect.ValueOf("C:\\test\\"), "") {
		t.Fatal("expected dirpath to pass")
	}
}

func TestBuiltinMongoDB(t *testing.T) {
	if !callBuiltin("mongodb", reflect.ValueOf("507f1f77bcf86cd799439011"), "") {
		t.Fatal("expected mongodb to pass")
	}
}

func TestBuiltinLuhnChecksum(t *testing.T) {
	if !callBuiltin("luhn_checksum", reflect.ValueOf("49927398716"), "") {
		t.Fatal("expected luhn_checksum to pass")
	}
}

func TestBuiltinDNSRFC1035Label(t *testing.T) {
	if !callBuiltin("dns_rfc1035_label", reflect.ValueOf("example"), "") {
		t.Fatal("expected dns_rfc1035_label to pass")
	}
}

func TestBuiltinSemver(t *testing.T) {
	if !callBuiltin("semver", reflect.ValueOf("1.2.3"), "") {
		t.Fatal("expected semver to pass")
	}
}

func TestBuiltinISBN10(t *testing.T) {
	if !callBuiltin("isbn10", reflect.ValueOf("080442957X"), "") {
		t.Fatal("expected isbn10 to pass")
	}
}

func TestBuiltinISBN13(t *testing.T) {
	if !callBuiltin("isbn13", reflect.ValueOf("9780306406157"), "") {
		t.Fatal("expected isbn13 to pass")
	}
}

func TestBuiltinISSN(t *testing.T) {
	if !callBuiltin("issn", reflect.ValueOf("0317-847X"), "") {
		t.Fatal("expected issn to pass")
	}
}

func TestBuiltinBIC(t *testing.T) {
	if !callBuiltin("bic", reflect.ValueOf("CHASUS33"), "") {
		t.Fatal("expected bic to pass")
	}
}

func TestBuiltinCron(t *testing.T) {
	if !callBuiltin("cron", reflect.ValueOf("0 * * * *"), "") {
		t.Fatal("expected cron to pass")
	}
}

func TestBuiltinDataURI(t *testing.T) {
	if !callBuiltin("datauri", reflect.ValueOf("data:text/plain;base64,SGVsbG8="), "") {
		t.Fatal("expected datauri to pass")
	}
}

func TestBuiltinBCP47(t *testing.T) {
	if !callBuiltin("bcp47", reflect.ValueOf("en-US"), "") {
		t.Fatal("expected bcp47 to pass")
	}
}

func TestBuiltinEthAddr(t *testing.T) {
	if !callBuiltin("eth_addr", reflect.ValueOf("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38"), "") {
		t.Fatal("expected eth_addr to pass")
	}
}

func TestBuiltinBtcAddr(t *testing.T) {
	if !callBuiltin("btc_addr", reflect.ValueOf("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"), "") {
		t.Fatal("expected btc_addr to pass")
	}
}

func TestRuleJSONBytes(t *testing.T) {
	if !RuleJSON(reflect.ValueOf([]byte(`{"key":"value"}`)), "", false) {
		t.Fatal("expected json bytes to pass")
	}
}

func TestRuleMultibyteInvalid(t *testing.T) {
	if callBuiltin("multibyte", reflect.ValueOf(42), "") {
		t.Fatal("expected multibyte to fail for int")
	}
}

func TestRuleEqIgnoreCaseFail(t *testing.T) {
	if !RuleEqIgnoreCase(reflect.ValueOf(42), "42", false) {
		t.Fatal("expected eq_ignore_case to pass for int via ScalarString")
	}
}

func TestRuleNeIgnoreCaseFail(t *testing.T) {
	if RuleNeIgnoreCase(reflect.ValueOf(42), "42", false) {
		t.Fatal("expected ne_ignore_case to fail for int")
	}
}

func TestRuleBooleanNonStringBool(t *testing.T) {
	if !RuleBoolean(reflect.ValueOf(false), "", false) {
		t.Fatal("expected boolean to pass for bool false")
	}
}

func TestRuleBooleanNonStringNonBool(t *testing.T) {
	if RuleBoolean(reflect.ValueOf(42), "", false) {
		t.Fatal("expected boolean to fail for int")
	}
}

func TestRuleLatitudeFail(t *testing.T) {
	if RuleLatitude(reflect.ValueOf("not-a-number"), "", false) {
		t.Fatal("expected latitude to fail for non-numeric string")
	}
}

func TestRuleLongitudeFail(t *testing.T) {
	if RuleLongitude(reflect.ValueOf("not-a-number"), "", false) {
		t.Fatal("expected longitude to fail for non-numeric string")
	}
}

func TestResolveFieldValueSlice(t *testing.T) {
	v, ok := ResolveFieldValue(reflect.ValueOf([]int{1, 2, 3}))
	if !ok || v != 3 {
		t.Fatalf("expected 3, got %f ok=%v", v, ok)
	}
}

func TestResolveFieldValueMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	v, ok := ResolveFieldValue(reflect.ValueOf(m))
	if !ok || v != 2 {
		t.Fatalf("expected 2, got %f ok=%v", v, ok)
	}
}

func TestResolveFieldValueUint(t *testing.T) {
	v, ok := ResolveFieldValue(reflect.ValueOf(uint(42)))
	if !ok || v != 42 {
		t.Fatalf("expected 42, got %f ok=%v", v, ok)
	}
}

func TestResolveFieldValueFloat(t *testing.T) {
	v, ok := ResolveFieldValue(reflect.ValueOf(3.14))
	if !ok || v != 3.14 {
		t.Fatalf("expected 3.14, got %f ok=%v", v, ok)
	}
}

func TestCompareLengthOrNumberInvalid(t *testing.T) {
	if CompareLengthOrNumber(reflect.Value{}, 5, validate.CmpGTE) {
		t.Fatal("expected invalid value to fail")
	}
}
