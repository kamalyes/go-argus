/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-20 11:22:32
 * @FilePath: \go-argus\rule\string_rules_test.go
 * @Description: string_rules.go 测试，覆盖 StringRule* 委托函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rule

import "testing"

func TestStringRuleRequired(t *testing.T) {
	if !StringRuleRequired("hello", "") {
		t.Fatal("expected required to pass")
	}
	if StringRuleRequired("", "") {
		t.Fatal("expected required to fail for empty")
	}
}

func TestStringRuleIsDefault(t *testing.T) {
	if !StringRuleIsDefault("", "") {
		t.Fatal("expected isdefault to pass for empty")
	}
}

func TestStringRuleMinMaxLen(t *testing.T) {
	if !StringRuleMin("abc", "3") {
		t.Fatal("expected min to pass")
	}
	if !StringRuleMax("ab", "3") {
		t.Fatal("expected max to pass")
	}
	if !StringRuleLen("abc", "3") {
		t.Fatal("expected len to pass")
	}
}

func TestStringRuleEqNe(t *testing.T) {
	if !StringRuleEq("hello", "hello") {
		t.Fatal("expected eq to pass")
	}
	if !StringRuleNe("hello", "world") {
		t.Fatal("expected ne to pass")
	}
}

func TestStringRuleEqIgnoreCase(t *testing.T) {
	if !StringRuleEqIgnoreCase("Hello", "hello") {
		t.Fatal("expected eq_ignore_case to pass")
	}
}

func TestStringRuleNeIgnoreCase(t *testing.T) {
	if !StringRuleNeIgnoreCase("Hello", "world") {
		t.Fatal("expected ne_ignore_case to pass")
	}
}

func TestStringRuleGtGteLtLte(t *testing.T) {
	if !StringRuleGt("abcd", "3") {
		t.Fatal("expected gt to pass")
	}
	if !StringRuleGte("abc", "3") {
		t.Fatal("expected gte to pass")
	}
	if !StringRuleLt("ab", "3") {
		t.Fatal("expected lt to pass")
	}
	if !StringRuleLte("abc", "3") {
		t.Fatal("expected lte to pass")
	}
}

func TestStringRuleAlpha(t *testing.T) {
	if !StringRuleAlpha("abc", "") {
		t.Fatal("expected alpha to pass")
	}
}

func TestStringRuleAlphaSpace(t *testing.T) {
	if !StringRuleAlphaSpace("hello world", "") {
		t.Fatal("expected alphaspace to pass")
	}
}

func TestStringRuleAlphanum(t *testing.T) {
	if !StringRuleAlphanum("abc123", "") {
		t.Fatal("expected alphanum to pass")
	}
}

func TestStringRuleAlphanumSpace(t *testing.T) {
	if !StringRuleAlphanumSpace("abc 123", "") {
		t.Fatal("expected alphanumspace to pass")
	}
}

func TestStringRuleAlphaUnicode(t *testing.T) {
	if !StringRuleAlphaUnicode("你好", "") {
		t.Fatal("expected alphaunicode to pass")
	}
}

func TestStringRuleAlphanumUnicode(t *testing.T) {
	if !StringRuleAlphanumUnicode("你好123", "") {
		t.Fatal("expected alphanumunicode to pass")
	}
}

func TestStringRuleASCII(t *testing.T) {
	if !StringRuleASCII("hello", "") {
		t.Fatal("expected ascii to pass")
	}
}

func TestStringRulePrintASCII(t *testing.T) {
	if !StringRulePrintASCII("hello", "") {
		t.Fatal("expected printascii to pass")
	}
}

func TestStringRuleMultibyte(t *testing.T) {
	if !StringRuleMultibyte("你好", "") {
		t.Fatal("expected multibyte to pass")
	}
}

func TestStringRuleHexadecimal(t *testing.T) {
	if !StringRuleHexadecimal("abcdef0123456789", "") {
		t.Fatal("expected hexadecimal to pass")
	}
}

func TestStringRuleHexColor(t *testing.T) {
	if !StringRuleHexColor("#ff0000", "") {
		t.Fatal("expected hexcolor to pass")
	}
}

func TestStringRuleRGB(t *testing.T) {
	if !StringRuleRGB("rgb(255,0,0)", "") {
		t.Fatal("expected rgb to pass")
	}
}

func TestStringRuleRGBA(t *testing.T) {
	if !StringRuleRGBA("rgba(255,0,0,0.5)", "") {
		t.Fatal("expected rgba to pass")
	}
}

func TestStringRuleHSL(t *testing.T) {
	if !StringRuleHSL("hsl(120,100%,50%)", "") {
		t.Fatal("expected hsl to pass")
	}
}

func TestStringRuleHSLA(t *testing.T) {
	if !StringRuleHSLA("hsla(120,100%,50%,0.5)", "") {
		t.Fatal("expected hsla to pass")
	}
}

func TestStringRuleEmail(t *testing.T) {
	if !StringRuleEmail("test@example.com", "") {
		t.Fatal("expected email to pass")
	}
}

func TestStringRuleE164(t *testing.T) {
	if !StringRuleE164("+1234567890", "") {
		t.Fatal("expected e164 to pass")
	}
}

func TestStringRuleIP(t *testing.T) {
	if !StringRuleIP("192.168.1.1", "") {
		t.Fatal("expected ip to pass")
	}
}

func TestStringRuleIPv4(t *testing.T) {
	if !StringRuleIPv4("192.168.1.1", "") {
		t.Fatal("expected ipv4 to pass")
	}
}

func TestStringRuleIPv6(t *testing.T) {
	if !StringRuleIPv6("::1", "") {
		t.Fatal("expected ipv6 to pass")
	}
}

func TestStringRuleCIDR(t *testing.T) {
	if !StringRuleCIDR("10.0.0.0/8", "") {
		t.Fatal("expected cidr to pass")
	}
}

func TestStringRuleCIDRv4(t *testing.T) {
	if !StringRuleCIDRv4("10.0.0.0/8", "") {
		t.Fatal("expected cidrv4 to pass")
	}
}

func TestStringRuleCIDRv6(t *testing.T) {
	if !StringRuleCIDRv6("::1/128", "") {
		t.Fatal("expected cidrv6 to pass")
	}
}

func TestStringRuleMAC(t *testing.T) {
	if !StringRuleMAC("00:11:22:33:44:55", "") {
		t.Fatal("expected mac to pass")
	}
}

func TestStringRuleHostname(t *testing.T) {
	if !StringRuleHostname("example.com", "") {
		t.Fatal("expected hostname to pass")
	}
}

func TestStringRuleFQDN(t *testing.T) {
	if !StringRuleFQDN("example.com.", "") {
		t.Fatal("expected fqdn to pass")
	}
}

func TestStringRuleHostnamePort(t *testing.T) {
	if !StringRuleHostnamePort("example.com:8080", "") {
		t.Fatal("expected hostname_port to pass")
	}
}

func TestStringRulePort(t *testing.T) {
	if !StringRulePort("443", "") {
		t.Fatal("expected port to pass")
	}
}

func TestStringRuleURL(t *testing.T) {
	if !StringRuleURL("http://example.com", "") {
		t.Fatal("expected url to pass")
	}
}

func TestStringRuleURI(t *testing.T) {
	if !StringRuleURI("http://example.com/path", "") {
		t.Fatal("expected uri to pass")
	}
}

func TestStringRuleHTTPURL(t *testing.T) {
	if !StringRuleHTTPURL("http://example.com", "") {
		t.Fatal("expected http_url to pass")
	}
}

func TestStringRuleHTTPSURL(t *testing.T) {
	if !StringRuleHTTPSURL("https://example.com", "") {
		t.Fatal("expected https_url to pass")
	}
}

func TestStringRuleURLEncoded(t *testing.T) {
	if !StringRuleURLEncoded("hello%20world", "") {
		t.Fatal("expected url_encoded to pass")
	}
}

func TestStringRuleHTML(t *testing.T) {
	if !StringRuleHTML("<b>hello</b>", "") {
		t.Fatal("expected html to pass")
	}
}

func TestStringRuleHTMLEncoded(t *testing.T) {
	if !StringRuleHTMLEncoded("&lt;b&gt;", "") {
		t.Fatal("expected html_encoded to pass")
	}
}

func TestStringRuleUUID(t *testing.T) {
	if !StringRuleUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8", "") {
		t.Fatal("expected uuid to pass")
	}
}

func TestStringRuleUUID3(t *testing.T) {
	if !StringRuleUUID3("6ba7b810-9dad-31d1-80b4-00c04fd430c8", "") {
		t.Fatal("expected uuid3 to pass")
	}
}

func TestStringRuleUUID4(t *testing.T) {
	if !StringRuleUUID4("6ba7b810-9dad-41d1-80b4-00c04fd430c8", "") {
		t.Fatal("expected uuid4 to pass")
	}
}

func TestStringRuleUUID5(t *testing.T) {
	if !StringRuleUUID5("6ba7b810-9dad-51d1-80b4-00c04fd430c8", "") {
		t.Fatal("expected uuid5 to pass")
	}
}

func TestStringRuleBase32(t *testing.T) {
	if !StringRuleBase32("JBSWY3DPEE======", "") {
		t.Fatal("expected base32 to pass")
	}
}

func TestStringRuleBase64(t *testing.T) {
	if !StringRuleBase64("SGVsbG8=", "") {
		t.Fatal("expected base64 to pass")
	}
}

func TestStringRuleBase64URL(t *testing.T) {
	if !StringRuleBase64URL("SGVsbG8=", "") {
		t.Fatal("expected base64url to pass")
	}
}

func TestStringRuleBase64RawURL(t *testing.T) {
	if !StringRuleBase64RawURL("SGVsbG8", "") {
		t.Fatal("expected base64rawurl to pass")
	}
}

func TestStringRuleJSON(t *testing.T) {
	if !StringRuleJSON(`{"key":"value"}`, "") {
		t.Fatal("expected json to pass")
	}
}

func TestStringRuleUnique(t *testing.T) {
	if !StringRuleUnique("abc", "") {
		t.Fatal("expected unique to pass")
	}
}

func TestStringRuleStartsWith(t *testing.T) {
	if !StringRuleStartsWith("hello world", "hello") {
		t.Fatal("expected startswith to pass")
	}
}

func TestStringRuleEndsWith(t *testing.T) {
	if !StringRuleEndsWith("hello world", "world") {
		t.Fatal("expected endswith to pass")
	}
}

func TestStringRuleStartsNotWith(t *testing.T) {
	if !StringRuleStartsNotWith("hello", "xyz") {
		t.Fatal("expected startsnotwith to pass")
	}
}

func TestStringRuleEndsNotWith(t *testing.T) {
	if !StringRuleEndsNotWith("hello", "xyz") {
		t.Fatal("expected endsnotwith to pass")
	}
}

func TestStringRuleContains(t *testing.T) {
	if !StringRuleContains("hello world", "world") {
		t.Fatal("expected contains to pass")
	}
}

func TestStringRuleContainsAny(t *testing.T) {
	if !StringRuleContainsAny("hello", "ae") {
		t.Fatal("expected containsany to pass")
	}
}

func TestStringRuleContainsRune(t *testing.T) {
	if !StringRuleContainsRune("hello", "e") {
		t.Fatal("expected containsrune to pass")
	}
}

func TestStringRuleExcludes(t *testing.T) {
	if !StringRuleExcludes("hello", "xyz") {
		t.Fatal("expected excludes to pass")
	}
}

func TestStringRuleExcludesAll(t *testing.T) {
	if !StringRuleExcludesAll("hello", "xyz") {
		t.Fatal("expected excludesall to pass")
	}
}

func TestStringRuleExcludesRune(t *testing.T) {
	if !StringRuleExcludesRune("hello", "z") {
		t.Fatal("expected excludesrune to pass")
	}
}

func TestStringRuleLowercase(t *testing.T) {
	if !StringRuleLowercase("hello", "") {
		t.Fatal("expected lowercase to pass")
	}
}

func TestStringRuleUppercase(t *testing.T) {
	if !StringRuleUppercase("HELLO", "") {
		t.Fatal("expected uppercase to pass")
	}
}

func TestStringRuleBoolean(t *testing.T) {
	if !StringRuleBoolean("true", "") {
		t.Fatal("expected boolean to pass")
	}
}

func TestStringRuleNumber(t *testing.T) {
	if !StringRuleNumber("123", "") {
		t.Fatal("expected number to pass")
	}
}

func TestStringRuleDatetime(t *testing.T) {
	if !StringRuleDatetime("2024-01-01", "2006-01-02") {
		t.Fatal("expected datetime to pass")
	}
}

func TestStringRuleTimezone(t *testing.T) {
	if !StringRuleTimezone("UTC", "") {
		t.Fatal("expected timezone to pass")
	}
}

func TestStringRuleLatitude(t *testing.T) {
	if !StringRuleLatitude("45.0", "") {
		t.Fatal("expected latitude to pass")
	}
}

func TestStringRuleLongitude(t *testing.T) {
	if !StringRuleLongitude("90.0", "") {
		t.Fatal("expected longitude to pass")
	}
}

func TestStringRuleFile(t *testing.T) {
	_ = StringRuleFile("nonexistent.txt", "")
}

func TestStringRuleFilePath(t *testing.T) {
	if !StringRuleFilePath("/test/file.txt", "") && !StringRuleFilePath("C:\\test\\file.txt", "") {
		t.Fatal("expected filepath to pass")
	}
}

func TestStringRuleDir(t *testing.T) {
	_ = StringRuleDir("nonexistent_dir", "")
}

func TestStringRuleDirPath(t *testing.T) {
	if !StringRuleDirPath("/test/", "") && !StringRuleDirPath("C:\\test\\", "") {
		t.Fatal("expected dirpath to pass")
	}
}

func TestStringRuleMongoDB(t *testing.T) {
	if !StringRuleMongoDB("507f1f77bcf86cd799439011", "") {
		t.Fatal("expected mongodb to pass")
	}
}

func TestStringRuleLuhnChecksum(t *testing.T) {
	if !StringRuleLuhnChecksum("49927398716", "") {
		t.Fatal("expected luhn_checksum to pass")
	}
}

func TestStringRuleDNSRFC1035Label(t *testing.T) {
	if !StringRuleDNSRFC1035Label("example", "") {
		t.Fatal("expected dns_rfc1035_label to pass")
	}
}

func TestStringRuleSemver(t *testing.T) {
	if !StringRuleSemver("1.2.3", "") {
		t.Fatal("expected semver to pass")
	}
}

func TestStringRuleISBN10(t *testing.T) {
	if !StringRuleISBN10("080442957X", "") {
		t.Fatal("expected isbn10 to pass")
	}
}

func TestStringRuleISBN13(t *testing.T) {
	if !StringRuleISBN13("9780306406157", "") {
		t.Fatal("expected isbn13 to pass")
	}
}

func TestStringRuleISSN(t *testing.T) {
	if !StringRuleISSN("0317-847X", "") {
		t.Fatal("expected issn to pass")
	}
}

func TestStringRuleBIC(t *testing.T) {
	if !StringRuleBIC("CHASUS33", "") {
		t.Fatal("expected bic to pass")
	}
}

func TestStringRuleCron(t *testing.T) {
	if !StringRuleCron("0 * * * *", "") {
		t.Fatal("expected cron to pass")
	}
}

func TestStringRuleDataURI(t *testing.T) {
	if !StringRuleDataURI("data:text/plain;base64,SGVsbG8=", "") {
		t.Fatal("expected datauri to pass")
	}
}

func TestStringRuleBCP47(t *testing.T) {
	if !StringRuleBCP47("en-US", "") {
		t.Fatal("expected bcp47 to pass")
	}
}

func TestStringRuleEthAddr(t *testing.T) {
	if !StringRuleEthAddr("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38", "") {
		t.Fatal("expected eth_addr to pass")
	}
}

func TestStringRuleBtcAddr(t *testing.T) {
	if !StringRuleBtcAddr("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "") {
		t.Fatal("expected btc_addr to pass")
	}
}
