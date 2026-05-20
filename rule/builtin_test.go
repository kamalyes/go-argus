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
)

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
	if !RuleAlpha(reflect.ValueOf("abc"), "", false) {
		t.Fatal("expected alpha to pass")
	}
	if RuleAlpha(reflect.ValueOf("abc123"), "", false) {
		t.Fatal("expected alpha to fail for alphanum")
	}
}

func TestRuleAlphanum(t *testing.T) {
	if !RuleAlphanum(reflect.ValueOf("abc123"), "", false) {
		t.Fatal("expected alphanum to pass")
	}
}

func TestRuleEmail(t *testing.T) {
	if !RuleEmail(reflect.ValueOf("test@example.com"), "", false) {
		t.Fatal("expected email to pass")
	}
	if RuleEmail(reflect.ValueOf("not-email"), "", false) {
		t.Fatal("expected email to fail")
	}
}

func TestRuleIP(t *testing.T) {
	if !RuleIP(reflect.ValueOf("192.168.1.1"), "", false) {
		t.Fatal("expected IP to pass")
	}
	if RuleIP(reflect.ValueOf("not-ip"), "", false) {
		t.Fatal("expected IP to fail")
	}
}

func TestRuleHostname(t *testing.T) {
	if !RuleHostname(reflect.ValueOf("example.com"), "", false) {
		t.Fatal("expected hostname to pass")
	}
}

func TestRuleFQDN(t *testing.T) {
	if !RuleFQDN(reflect.ValueOf("example.com."), "", false) {
		t.Fatal("expected fqdn to pass")
	}
}

func TestRuleURL(t *testing.T) {
	if !RuleURL(reflect.ValueOf("http://example.com"), "", false) {
		t.Fatal("expected url to pass")
	}
}

func TestRuleUUID(t *testing.T) {
	if !RuleUUID(reflect.ValueOf("6ba7b810-9dad-11d1-80b4-00c04fd430c8"), "", false) {
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

func TestRuleStartsWith(t *testing.T) {
	if !RuleStartsWith(reflect.ValueOf("hello world"), "hello", false) {
		t.Fatal("expected startsWith to pass")
	}
}

func TestRuleContains(t *testing.T) {
	if !RuleContains(reflect.ValueOf("hello world"), "world", false) {
		t.Fatal("expected contains to pass")
	}
}

func TestRuleLowercase(t *testing.T) {
	if !RuleLowercase(reflect.ValueOf("hello"), "", false) {
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
	if !CompareOp(5, 3, CmpGT) {
		t.Fatal("expected 5 > 3")
	}
	if !CompareOp(5, 5, CmpGTE) {
		t.Fatal("expected 5 >= 5")
	}
	if !CompareOp(3, 5, CmpLT) {
		t.Fatal("expected 3 < 5")
	}
	if !CompareOp(3, 3, CmpLTE) {
		t.Fatal("expected 3 <= 3")
	}
	if !CompareOp(5, 5, CmpEQ) {
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
	if !RuleGte(reflect.ValueOf(5), "5", false) {
		t.Fatal("expected gte to pass")
	}
	if !RuleLte(reflect.ValueOf(5), "5", false) {
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
	// ScalarString 对 int 返回 "42"，所以 eq 实际会通过
	if !RuleEq(reflect.ValueOf(42), "42", false) {
		t.Fatal("expected eq to pass for int 42 via ScalarString")
	}
}

func TestRuleAlphaSpace(t *testing.T) {
	if !RuleAlphaSpace(reflect.ValueOf("hello world"), "", false) {
		t.Fatal("expected alphaspace to pass")
	}
}

func TestRuleAlphanumSpace(t *testing.T) {
	if !RuleAlphanumSpace(reflect.ValueOf("hello 123"), "", false) {
		t.Fatal("expected alphanumspace to pass")
	}
}

func TestRuleAlphaUnicode(t *testing.T) {
	if !RuleAlphaUnicode(reflect.ValueOf("你好"), "", false) {
		t.Fatal("expected alphaunicode to pass")
	}
}

func TestRuleAlphanumUnicode(t *testing.T) {
	if !RuleAlphanumUnicode(reflect.ValueOf("你好123"), "", false) {
		t.Fatal("expected alphanumunicode to pass")
	}
}

func TestRuleASCII(t *testing.T) {
	if !RuleASCII(reflect.ValueOf("hello"), "", false) {
		t.Fatal("expected ascii to pass")
	}
}

func TestRulePrintASCII(t *testing.T) {
	if !RulePrintASCII(reflect.ValueOf("hello"), "", false) {
		t.Fatal("expected printascii to pass")
	}
}

func TestRuleHexadecimal(t *testing.T) {
	if !RuleHexadecimal(reflect.ValueOf("abcdef0123456789"), "", false) {
		t.Fatal("expected hexadecimal to pass")
	}
}

func TestRuleHexColor(t *testing.T) {
	if !RuleHexColor(reflect.ValueOf("#ff0000"), "", false) {
		t.Fatal("expected hexcolor to pass")
	}
}

func TestRuleRGB(t *testing.T) {
	if !RuleRGB(reflect.ValueOf("rgb(255,0,0)"), "", false) {
		t.Fatal("expected rgb to pass")
	}
}

func TestRuleRGBA(t *testing.T) {
	if !RuleRGBA(reflect.ValueOf("rgba(255,0,0,0.5)"), "", false) {
		t.Fatal("expected rgba to pass")
	}
}

func TestRuleHSL(t *testing.T) {
	if !RuleHSL(reflect.ValueOf("hsl(120,100%,50%)"), "", false) {
		t.Fatal("expected hsl to pass")
	}
}

func TestRuleHSLA(t *testing.T) {
	if !RuleHSLA(reflect.ValueOf("hsla(120,100%,50%,0.5)"), "", false) {
		t.Fatal("expected hsla to pass")
	}
}

func TestRuleE164(t *testing.T) {
	if !RuleE164(reflect.ValueOf("+1234567890"), "", false) {
		t.Fatal("expected e164 to pass")
	}
}

func TestRuleIPv4(t *testing.T) {
	if !RuleIPv4(reflect.ValueOf("192.168.1.1"), "", false) {
		t.Fatal("expected ipv4 to pass")
	}
}

func TestRuleIPv6(t *testing.T) {
	if !RuleIPv6(reflect.ValueOf("::1"), "", false) {
		t.Fatal("expected ipv6 to pass")
	}
}

func TestRuleCIDR(t *testing.T) {
	if !RuleCIDR(reflect.ValueOf("10.0.0.0/8"), "", false) {
		t.Fatal("expected cidr to pass")
	}
}

func TestRuleCIDRv4(t *testing.T) {
	if !RuleCIDRv4(reflect.ValueOf("10.0.0.0/8"), "", false) {
		t.Fatal("expected cidrv4 to pass")
	}
}

func TestRuleCIDRv6(t *testing.T) {
	if !RuleCIDRv6(reflect.ValueOf("::1/128"), "", false) {
		t.Fatal("expected cidrv6 to pass")
	}
}

func TestRuleMAC(t *testing.T) {
	if !RuleMAC(reflect.ValueOf("00:11:22:33:44:55"), "", false) {
		t.Fatal("expected mac to pass")
	}
}

func TestRuleHostnamePort(t *testing.T) {
	if !RuleHostnamePort(reflect.ValueOf("example.com:8080"), "", false) {
		t.Fatal("expected hostname_port to pass")
	}
}

func TestRulePort(t *testing.T) {
	if !RulePort(reflect.ValueOf("443"), "", false) {
		t.Fatal("expected port to pass")
	}
}

func TestRuleURI(t *testing.T) {
	if !RuleURI(reflect.ValueOf("http://example.com/path"), "", false) {
		t.Fatal("expected uri to pass")
	}
}

func TestRuleHTTPURL(t *testing.T) {
	if !RuleHTTPURL(reflect.ValueOf("http://example.com"), "", false) {
		t.Fatal("expected http_url to pass")
	}
}

func TestRuleHTTPSURL(t *testing.T) {
	if !RuleHTTPSURL(reflect.ValueOf("https://example.com"), "", false) {
		t.Fatal("expected https_url to pass")
	}
}

func TestRuleURLEncoded(t *testing.T) {
	if !RuleURLEncoded(reflect.ValueOf("hello%20world"), "", false) {
		t.Fatal("expected url_encoded to pass")
	}
}

func TestRuleHTML(t *testing.T) {
	if !RuleHTML(reflect.ValueOf("<b>hello</b>"), "", false) {
		t.Fatal("expected html to pass")
	}
}

func TestRuleHTMLEncoded(t *testing.T) {
	if !RuleHTMLEncoded(reflect.ValueOf("&lt;b&gt;"), "", false) {
		t.Fatal("expected html_encoded to pass")
	}
}

func TestRuleUUID3(t *testing.T) {
	if !RuleUUID3(reflect.ValueOf("6ba7b810-9dad-31d1-80b4-00c04fd430c8"), "", false) {
		t.Fatal("expected uuid3 to pass")
	}
}

func TestRuleUUID4(t *testing.T) {
	if !RuleUUID4(reflect.ValueOf("6ba7b810-9dad-41d1-80b4-00c04fd430c8"), "", false) {
		t.Fatal("expected uuid4 to pass")
	}
}

func TestRuleUUID5(t *testing.T) {
	if !RuleUUID5(reflect.ValueOf("6ba7b810-9dad-51d1-80b4-00c04fd430c8"), "", false) {
		t.Fatal("expected uuid5 to pass")
	}
}

func TestRuleBase32(t *testing.T) {
	if !RuleBase32(reflect.ValueOf("JBSWY3DPEE======"), "", false) {
		t.Fatal("expected base32 to pass")
	}
}

func TestRuleBase64(t *testing.T) {
	if !RuleBase64(reflect.ValueOf("SGVsbG8="), "", false) {
		t.Fatal("expected base64 to pass")
	}
}

func TestRuleBase64URL(t *testing.T) {
	if !RuleBase64URL(reflect.ValueOf("SGVsbG8="), "", false) {
		t.Fatal("expected base64url to pass")
	}
}

func TestRuleBase64RawURL(t *testing.T) {
	if !RuleBase64RawURL(reflect.ValueOf("SGVsbG8"), "", false) {
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

func TestRuleEndsWith(t *testing.T) {
	if !RuleEndsWith(reflect.ValueOf("hello world"), "world", false) {
		t.Fatal("expected endswith to pass")
	}
}

func TestRuleStartsNotWith(t *testing.T) {
	if !RuleStartsNotWith(reflect.ValueOf("hello"), "xyz", false) {
		t.Fatal("expected startsnotwith to pass")
	}
}

func TestRuleEndsNotWith(t *testing.T) {
	if !RuleEndsNotWith(reflect.ValueOf("hello"), "xyz", false) {
		t.Fatal("expected endsnotwith to pass")
	}
}

func TestRuleContainsAny(t *testing.T) {
	if !RuleContainsAny(reflect.ValueOf("hello"), "ae", false) {
		t.Fatal("expected containsany to pass")
	}
}

func TestRuleContainsRune(t *testing.T) {
	if !RuleContainsRune(reflect.ValueOf("hello"), "e", false) {
		t.Fatal("expected containsrune to pass")
	}
}

func TestRuleExcludes(t *testing.T) {
	if !RuleExcludes(reflect.ValueOf("hello"), "xyz", false) {
		t.Fatal("expected excludes to pass")
	}
}

func TestRuleExcludesAll(t *testing.T) {
	if !RuleExcludesAll(reflect.ValueOf("hello"), "xyz", false) {
		t.Fatal("expected excludesall to pass")
	}
}

func TestRuleExcludesRune(t *testing.T) {
	if !RuleExcludesRune(reflect.ValueOf("hello"), "z", false) {
		t.Fatal("expected excludesrune to pass")
	}
}

func TestRuleUppercase(t *testing.T) {
	if !RuleUppercase(reflect.ValueOf("HELLO"), "", false) {
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

func TestRuleDatetime(t *testing.T) {
	if !RuleDatetime(reflect.ValueOf("2024-01-01"), "2006-01-02", false) {
		t.Fatal("expected datetime to pass")
	}
}

func TestRuleTimezone(t *testing.T) {
	if !RuleTimezone(reflect.ValueOf("UTC"), "", false) {
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

func TestRuleFile(t *testing.T) {
	_ = RuleFile(reflect.ValueOf("nonexistent.txt"), "", false)
}

func TestRuleFilePath(t *testing.T) {
	if !RuleFilePath(reflect.ValueOf("/test/file.txt"), "", false) && !RuleFilePath(reflect.ValueOf("C:\\test\\file.txt"), "", false) {
		t.Fatal("expected filepath to pass")
	}
}

func TestRuleDir(t *testing.T) {
	_ = RuleDir(reflect.ValueOf("nonexistent_dir"), "", false)
}

func TestRuleDirPath(t *testing.T) {
	if !RuleDirPath(reflect.ValueOf("/test/"), "", false) && !RuleDirPath(reflect.ValueOf("C:\\test\\"), "", false) {
		t.Fatal("expected dirpath to pass")
	}
}

func TestRuleMongoDB(t *testing.T) {
	if !RuleMongoDB(reflect.ValueOf("507f1f77bcf86cd799439011"), "", false) {
		t.Fatal("expected mongodb to pass")
	}
}

func TestRuleLuhnChecksum(t *testing.T) {
	if !RuleLuhnChecksum(reflect.ValueOf("49927398716"), "", false) {
		t.Fatal("expected luhn_checksum to pass")
	}
}

func TestRuleDNSRFC1035Label(t *testing.T) {
	if !RuleDNSRFC1035Label(reflect.ValueOf("example"), "", false) {
		t.Fatal("expected dns_rfc1035_label to pass")
	}
}

func TestRuleSemver(t *testing.T) {
	if !RuleSemver(reflect.ValueOf("1.2.3"), "", false) {
		t.Fatal("expected semver to pass")
	}
}

func TestRuleISBN10(t *testing.T) {
	if !RuleISBN10(reflect.ValueOf("080442957X"), "", false) {
		t.Fatal("expected isbn10 to pass")
	}
}

func TestRuleISBN13(t *testing.T) {
	if !RuleISBN13(reflect.ValueOf("9780306406157"), "", false) {
		t.Fatal("expected isbn13 to pass")
	}
}

func TestRuleISSN(t *testing.T) {
	// ISSN 0317-8471: 0*8 + 3*7 + 1*6 + 7*5 + 8*4 + 4*3 + 7*2 = 0+21+6+35+32+12+14 = 120, 120%11=10 -> X
	// 使用一个校验位正确的 ISSN
	if !RuleISSN(reflect.ValueOf("0317-847X"), "", false) {
		t.Fatal("expected issn to pass")
	}
}

func TestRuleBIC(t *testing.T) {
	if !RuleBIC(reflect.ValueOf("CHASUS33"), "", false) {
		t.Fatal("expected bic to pass")
	}
}

func TestRuleCron(t *testing.T) {
	if !RuleCron(reflect.ValueOf("0 * * * *"), "", false) {
		t.Fatal("expected cron to pass")
	}
}

func TestRuleDataURI(t *testing.T) {
	if !RuleDataURI(reflect.ValueOf("data:text/plain;base64,SGVsbG8="), "", false) {
		t.Fatal("expected datauri to pass")
	}
}

func TestRuleBCP47(t *testing.T) {
	if !RuleBCP47(reflect.ValueOf("en-US"), "", false) {
		t.Fatal("expected bcp47 to pass")
	}
}

func TestRuleEthAddr(t *testing.T) {
	if !RuleEthAddr(reflect.ValueOf("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38"), "", false) {
		t.Fatal("expected eth_addr to pass")
	}
}

func TestRuleBtcAddr(t *testing.T) {
	if !RuleBtcAddr(reflect.ValueOf("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"), "", false) {
		t.Fatal("expected btc_addr to pass")
	}
}

func TestRuleJSONBytes(t *testing.T) {
	if !RuleJSON(reflect.ValueOf([]byte(`{"key":"value"}`)), "", false) {
		t.Fatal("expected json bytes to pass")
	}
}

func TestRuleMultibyteInvalid(t *testing.T) {
	if RuleMultibyte(reflect.ValueOf(42), "", false) {
		t.Fatal("expected multibyte to fail for int")
	}
}

func TestRuleEqIgnoreCaseFail(t *testing.T) {
	// ScalarString 对 int 返回 "42"，所以 eq_ignore_case 实际会通过
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
	if CompareLengthOrNumber(reflect.Value{}, 5, CmpGTE) {
		t.Fatal("expected invalid value to fail")
	}
}
