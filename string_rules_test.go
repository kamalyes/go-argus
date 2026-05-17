/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-17 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 00:00:00
 * @FilePath: \go-argus\string_rules_test.go
 * @Description: string_rules.go 测试，覆盖所有零反射字符串规则函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"os"
	"testing"
)

func TestIsBlankString(t *testing.T) {
	if !isBlankString("") {
		t.Fatal("expected empty string to be blank")
	}
	if !isBlankString("   ") {
		t.Fatal("expected spaces to be blank")
	}
	if !isBlankString("\t\n\r") {
		t.Fatal("expected whitespace to be blank")
	}
	if isBlankString("a") {
		t.Fatal("expected non-whitespace to not be blank")
	}
}

func TestStringRequired(t *testing.T) {
	if stringRequired("", "") {
		t.Fatal("expected empty to fail required")
	}
	if stringRequired("   ", "") {
		t.Fatal("expected spaces to fail required")
	}
	if !stringRequired("hello", "") {
		t.Fatal("expected non-empty to pass required")
	}
}

func TestStringIsDefault(t *testing.T) {
	if !stringIsDefault("", "") {
		t.Fatal("expected empty to pass isdefault")
	}
	if !stringIsDefault("   ", "") {
		t.Fatal("expected spaces to pass isdefault")
	}
	if stringIsDefault("hello", "") {
		t.Fatal("expected non-empty to fail isdefault")
	}
}

func TestStringMin(t *testing.T) {
	if !stringMin("abc", "2") {
		t.Fatal("expected abc to pass min=2")
	}
	if stringMin("a", "2") {
		t.Fatal("expected a to fail min=2")
	}
	if stringMin("abc", "abc") {
		t.Fatal("expected min to fail for non-numeric param")
	}
}

func TestStringMax(t *testing.T) {
	if !stringMax("ab", "3") {
		t.Fatal("expected ab to pass max=3")
	}
	if stringMax("abcd", "3") {
		t.Fatal("expected abcd to fail max=3")
	}
	if stringMax("abc", "abc") {
		t.Fatal("expected max to fail for non-numeric param")
	}
}

func TestStringLen(t *testing.T) {
	if !stringLen("abc", "3") {
		t.Fatal("expected abc to pass len=3")
	}
	if stringLen("ab", "3") {
		t.Fatal("expected ab to fail len=3")
	}
	if stringLen("abc", "abc") {
		t.Fatal("expected len to fail for non-numeric param")
	}
}

func TestStringEq(t *testing.T) {
	if !stringEq("hello", "hello") {
		t.Fatal("expected equal strings to pass")
	}
	if stringEq("hello", "world") {
		t.Fatal("expected different strings to fail")
	}
}

func TestStringEqIgnoreCase(t *testing.T) {
	if !stringEqIgnoreCase("Hello", "hello") {
		t.Fatal("expected case-insensitive equal to pass")
	}
	if stringEqIgnoreCase("Hello", "world") {
		t.Fatal("expected different strings to fail")
	}
}

func TestStringNe(t *testing.T) {
	if !stringNe("hello", "world") {
		t.Fatal("expected different strings to pass ne")
	}
	if stringNe("hello", "hello") {
		t.Fatal("expected equal strings to fail ne")
	}
}

func TestStringNeIgnoreCase(t *testing.T) {
	if !stringNeIgnoreCase("Hello", "world") {
		t.Fatal("expected different strings to pass ne_ignore_case")
	}
	if stringNeIgnoreCase("Hello", "hello") {
		t.Fatal("expected equal strings to fail ne_ignore_case")
	}
}

func TestStringGtGteLtLte(t *testing.T) {
	if !stringGt("abc", "2") {
		t.Fatal("expected len(abc)>2 to pass")
	}
	if stringGt("abc", "3") {
		t.Fatal("expected len(abc)>3 to fail")
	}
	if !stringGte("abc", "3") {
		t.Fatal("expected len(abc)>=3 to pass")
	}
	if stringGte("ab", "3") {
		t.Fatal("expected len(ab)>=3 to fail")
	}
	if !stringLt("ab", "3") {
		t.Fatal("expected len(ab)<3 to pass")
	}
	if stringLt("abc", "3") {
		t.Fatal("expected len(abc)<3 to fail")
	}
	if !stringLte("abc", "3") {
		t.Fatal("expected len(abc)<=3 to pass")
	}
	if stringLte("abcd", "3") {
		t.Fatal("expected len(abcd)<=3 to fail")
	}
}

func TestStringAlpha(t *testing.T) {
	if !stringAlpha("abc", "") {
		t.Fatal("expected alpha to pass")
	}
	if stringAlpha("abc123", "") {
		t.Fatal("expected alpha to fail for alphanum")
	}
	if stringAlpha("", "") {
		t.Fatal("expected alpha to fail for empty")
	}
}

func TestStringAlphaSpace(t *testing.T) {
	if !stringAlphaSpace("hello world", "") {
		t.Fatal("expected alphaspace to pass")
	}
	if stringAlphaSpace("hello123", "") {
		t.Fatal("expected alphaspace to fail for digits")
	}
}

func TestStringAlphanum(t *testing.T) {
	if !stringAlphanum("abc123", "") {
		t.Fatal("expected alphanum to pass")
	}
	if stringAlphanum("abc-123", "") {
		t.Fatal("expected alphanum to fail for hyphen")
	}
}

func TestStringAlphanumSpace(t *testing.T) {
	if !stringAlphanumSpace("abc 123", "") {
		t.Fatal("expected alphanumspace to pass")
	}
}

func TestStringAlphaUnicode(t *testing.T) {
	if !stringAlphaUnicode("你好世界", "") {
		t.Fatal("expected alphaunicode to pass")
	}
	if stringAlphaUnicode("你好1", "") {
		t.Fatal("expected alphaunicode to fail for digits")
	}
}

func TestStringAlphanumUnicode(t *testing.T) {
	if !stringAlphanumUnicode("你好123", "") {
		t.Fatal("expected alphanumunicode to pass")
	}
}

func TestStringASCII(t *testing.T) {
	if !stringASCII("hello", "") {
		t.Fatal("expected ascii to pass")
	}
	if stringASCII("你好", "") {
		t.Fatal("expected ascii to fail for unicode")
	}
}

func TestStringPrintASCII(t *testing.T) {
	if !stringPrintASCII("hello world", "") {
		t.Fatal("expected printascii to pass")
	}
	if stringPrintASCII(string(rune(0x1f)), "") {
		t.Fatal("expected printascii to fail for control char")
	}
}

func TestStringMultibyte(t *testing.T) {
	if !stringMultibyte("你好", "") {
		t.Fatal("expected multibyte to pass")
	}
	if stringMultibyte("hello", "") {
		t.Fatal("expected multibyte to fail for ascii")
	}
}

func TestStringHexadecimal(t *testing.T) {
	if !stringHexadecimal("deadbeef", "") {
		t.Fatal("expected hexadecimal to pass")
	}
	if stringHexadecimal("xyz", "") {
		t.Fatal("expected hexadecimal to fail")
	}
}

func TestStringHexColor(t *testing.T) {
	if !stringHexColor("#12ffaa", "") {
		t.Fatal("expected hexcolor to pass")
	}
	if stringHexColor("invalid", "") {
		t.Fatal("expected hexcolor to fail")
	}
}

func TestStringRGB(t *testing.T) {
	if !stringRGB("rgb(12, 34, 255)", "") {
		t.Fatal("expected rgb to pass")
	}
	if stringRGB("rgb(300, 34, 255)", "") {
		t.Fatal("expected rgb to fail for out-of-range")
	}
}

func TestStringRGBA(t *testing.T) {
	if !stringRGBA("rgba(12, 34, 255, 0.5)", "") {
		t.Fatal("expected rgba to pass")
	}
}

func TestStringHSL(t *testing.T) {
	if !stringHSL("hsl(120, 50%, 75%)", "") {
		t.Fatal("expected hsl to pass")
	}
}

func TestStringHSLA(t *testing.T) {
	if !stringHSLA("hsla(120, 50%, 75%, 0.5)", "") {
		t.Fatal("expected hsla to pass")
	}
}

func TestStringEmail(t *testing.T) {
	if !stringEmail("user@example.com", "") {
		t.Fatal("expected email to pass")
	}
	if stringEmail("not-email", "") {
		t.Fatal("expected email to fail")
	}
}

func TestStringE164(t *testing.T) {
	if !stringE164("+8613800138000", "") {
		t.Fatal("expected e164 to pass")
	}
	if stringE164("12345", "") {
		t.Fatal("expected e164 to fail")
	}
}

func TestStringIP(t *testing.T) {
	if !stringIP("192.168.1.1", "") {
		t.Fatal("expected ip to pass")
	}
	if !stringIP("::1", "") {
		t.Fatal("expected ipv6 to pass ip")
	}
	if stringIP("not-an-ip", "") {
		t.Fatal("expected ip to fail")
	}
}

func TestStringIPv4(t *testing.T) {
	if !stringIPv4("192.168.1.1", "") {
		t.Fatal("expected ipv4 to pass")
	}
	if stringIPv4("::1", "") {
		t.Fatal("expected ipv6 to fail ipv4")
	}
}

func TestStringIPv6(t *testing.T) {
	if !stringIPv6("::1", "") {
		t.Fatal("expected ipv6 to pass")
	}
	if stringIPv6("192.168.1.1", "") {
		t.Fatal("expected ipv4 to fail ipv6")
	}
}

func TestStringCIDR(t *testing.T) {
	if !stringCIDR("192.168.1.0/24", "") {
		t.Fatal("expected cidr to pass")
	}
	if stringCIDR("invalid", "") {
		t.Fatal("expected cidr to fail")
	}
}

func TestStringCIDRv4(t *testing.T) {
	if !stringCIDRv4("192.168.1.0/24", "") {
		t.Fatal("expected cidrv4 to pass")
	}
	if stringCIDRv4("::1/128", "") {
		t.Fatal("expected ipv6 cidr to fail cidrv4")
	}
}

func TestStringCIDRv6(t *testing.T) {
	if !stringCIDRv6("::1/128", "") {
		t.Fatal("expected cidrv6 to pass")
	}
	if stringCIDRv6("192.168.1.0/24", "") {
		t.Fatal("expected ipv4 cidr to fail cidrv6")
	}
}

func TestStringMAC(t *testing.T) {
	if !stringMAC("00:11:22:33:44:55", "") {
		t.Fatal("expected mac to pass")
	}
	if stringMAC("invalid", "") {
		t.Fatal("expected mac to fail")
	}
}

func TestStringHostname(t *testing.T) {
	if !stringHostname("api.example.com", "") {
		t.Fatal("expected hostname to pass")
	}
	if stringHostname("-invalid.com", "") {
		t.Fatal("expected hostname to fail")
	}
	if !stringHostname("api.example.com.", "") {
		t.Fatal("expected hostname to pass with trailing dot")
	}
}

func TestStringFQDN(t *testing.T) {
	if !stringFQDN("api.example.com.", "") {
		t.Fatal("expected fqdn to pass")
	}
	if stringFQDN("api.example.com", "") {
		t.Fatal("expected fqdn to fail without trailing dot")
	}
}

func TestStringHostnamePort(t *testing.T) {
	if !stringHostnamePort("example.com:8080", "") {
		t.Fatal("expected hostname_port to pass")
	}
	if stringHostnamePort("example.com:99999", "") {
		t.Fatal("expected hostname_port to fail for invalid port")
	}
	if stringHostnamePort(":8080", "") {
		t.Fatal("expected hostname_port to fail for empty host")
	}
	if stringHostnamePort("noport", "") {
		t.Fatal("expected hostname_port to fail for missing port")
	}
}

func TestStringPort(t *testing.T) {
	if !stringPort("443", "") {
		t.Fatal("expected port to pass")
	}
	if stringPort("99999", "") {
		t.Fatal("expected port to fail for out-of-range")
	}
	if stringPort("-1", "") {
		t.Fatal("expected port to fail for negative")
	}
	if stringPort("abc", "") {
		t.Fatal("expected port to fail for non-numeric")
	}
}

func TestStringURL(t *testing.T) {
	if !stringURL("https://example.com/path", "") {
		t.Fatal("expected url to pass")
	}
	if stringURL("not-a-url", "") {
		t.Fatal("expected url to fail")
	}
}

func TestStringURI(t *testing.T) {
	if !stringURI("https://example.com/path", "") {
		t.Fatal("expected uri to pass for https")
	}
	if !stringURI("mailto:user@example.com", "") {
		t.Fatal("expected uri to pass for mailto")
	}
	if !stringURI("tel:+1-234-567-8900", "") {
		t.Fatal("expected uri to pass for tel")
	}
	if stringURI("no-scheme-here", "") {
		t.Fatal("expected uri to fail for no scheme")
	}
	if stringURI(":invalid", "") {
		t.Fatal("expected uri to fail for empty scheme")
	}
}

func TestStringHTTPURL(t *testing.T) {
	if !stringHTTPURL("https://example.com", "") {
		t.Fatal("expected http_url to pass for https")
	}
	if !stringHTTPURL("http://example.com", "") {
		t.Fatal("expected http_url to pass for http")
	}
	if stringHTTPURL("ftp://example.com", "") {
		t.Fatal("expected http_url to fail for ftp")
	}
	if stringHTTPURL("https://", "") {
		t.Fatal("expected http_url to fail for no host")
	}
	if stringHTTPURL("noscheme", "") {
		t.Fatal("expected http_url to fail for no scheme")
	}
}

func TestStringHTTPSURL(t *testing.T) {
	if !stringHTTPSURL("https://example.com", "") {
		t.Fatal("expected https_url to pass")
	}
	if stringHTTPSURL("http://example.com", "") {
		t.Fatal("expected https_url to fail for http")
	}
	if stringHTTPSURL("https://", "") {
		t.Fatal("expected https_url to fail for no host")
	}
	if stringHTTPSURL("noscheme", "") {
		t.Fatal("expected https_url to fail for no scheme")
	}
}

func TestStringURLEncoded(t *testing.T) {
	if !stringURLEncoded("hello%20world", "") {
		t.Fatal("expected url_encoded to pass")
	}
	if stringURLEncoded("hello%ZZ", "") {
		t.Fatal("expected url_encoded to fail for invalid encoding")
	}
	if stringURLEncoded("nopercent", "") {
		t.Fatal("expected url_encoded to fail without percent")
	}
}

func TestStringHTML(t *testing.T) {
	if !stringHTML("<b>bold</b>", "") {
		t.Fatal("expected html to pass")
	}
	if stringHTML("no html", "") {
		t.Fatal("expected html to fail")
	}
}

func TestStringHTMLEncoded(t *testing.T) {
	if !stringHTMLEncoded("&amp;", "") {
		t.Fatal("expected html_encoded to pass")
	}
	if stringHTMLEncoded("plain", "") {
		t.Fatal("expected html_encoded to fail for plain text")
	}
}

func TestStringUUID(t *testing.T) {
	if !stringUUID("550e8400-e29b-41d4-a716-446655440000", "") {
		t.Fatal("expected uuid to pass")
	}
	if stringUUID("not-a-uuid", "") {
		t.Fatal("expected uuid to fail")
	}
}

func TestStringUUID3(t *testing.T) {
	if !stringUUID3("550e8400-e29b-31d4-a716-446655440000", "") {
		t.Fatal("expected uuid3 to pass")
	}
	if stringUUID3("550e8400-e29b-41d4-a716-446655440000", "") {
		t.Fatal("expected uuid3 to fail for uuid4")
	}
}

func TestStringUUID4(t *testing.T) {
	if !stringUUID4("550e8400-e29b-41d4-a716-446655440000", "") {
		t.Fatal("expected uuid4 to pass")
	}
}

func TestStringUUID5(t *testing.T) {
	if !stringUUID5("550e8400-e29b-51d4-a716-446655440000", "") {
		t.Fatal("expected uuid5 to pass")
	}
}

func TestStringBase32(t *testing.T) {
	if !stringBase32("JBSWY3DPEB3W64TMMQ======", "") {
		t.Fatal("expected base32 to pass")
	}
	if stringBase32("!!!invalid!!!", "") {
		t.Fatal("expected base32 to fail")
	}
	if stringBase32("", "") {
		t.Fatal("expected base32 to fail for empty")
	}
}

func TestStringBase64(t *testing.T) {
	if !stringBase64("YXJndXM=", "") {
		t.Fatal("expected base64 to pass")
	}
	if stringBase64("!!!invalid!!!", "") {
		t.Fatal("expected base64 to fail")
	}
}

func TestStringBase64URL(t *testing.T) {
	if !stringBase64URL("YXJndXM=", "") {
		t.Fatal("expected base64url to pass")
	}
	if stringBase64URL("", "") {
		t.Fatal("expected base64url to fail for empty")
	}
}

func TestStringBase64RawURL(t *testing.T) {
	if !stringBase64RawURL("YXJndXM", "") {
		t.Fatal("expected base64rawurl to pass")
	}
	if stringBase64RawURL("", "") {
		t.Fatal("expected base64rawurl to fail for empty")
	}
}

func TestStringJSON(t *testing.T) {
	if !stringJSON(`{"key":"value"}`, "") {
		t.Fatal("expected json to pass")
	}
	if stringJSON("{invalid}", "") {
		t.Fatal("expected json to fail")
	}
}

func TestStringUnique(t *testing.T) {
	if !stringUnique("abcdef", "") {
		t.Fatal("expected unique string to pass")
	}
	if stringUnique("aabc", "") {
		t.Fatal("expected non-unique string to fail")
	}
}

func TestStringStartsWith(t *testing.T) {
	if !stringStartsWith("hello world", "hello") {
		t.Fatal("expected startswith to pass")
	}
	if stringStartsWith("hello world", "world") {
		t.Fatal("expected startswith to fail")
	}
}

func TestStringEndsWith(t *testing.T) {
	if !stringEndsWith("hello world", "world") {
		t.Fatal("expected endswith to pass")
	}
	if stringEndsWith("hello world", "hello") {
		t.Fatal("expected endswith to fail")
	}
}

func TestStringStartsNotWith(t *testing.T) {
	if !stringStartsNotWith("hello world", "world") {
		t.Fatal("expected startsnotwith to pass")
	}
	if stringStartsNotWith("hello world", "hello") {
		t.Fatal("expected startsnotwith to fail")
	}
}

func TestStringEndsNotWith(t *testing.T) {
	if !stringEndsNotWith("hello world", "hello") {
		t.Fatal("expected endsnotwith to pass")
	}
	if stringEndsNotWith("hello world", "world") {
		t.Fatal("expected endsnotwith to fail")
	}
}

func TestStringContains(t *testing.T) {
	if !stringContains("hello world", "world") {
		t.Fatal("expected contains to pass")
	}
	if stringContains("hello world", "xyz") {
		t.Fatal("expected contains to fail")
	}
}

func TestStringContainsAny(t *testing.T) {
	if !stringContainsAny("hello", "h") {
		t.Fatal("expected containsany to pass")
	}
	if stringContainsAny("hello", "xyz") {
		t.Fatal("expected containsany to fail")
	}
}

func TestStringContainsRune(t *testing.T) {
	if !stringContainsRune("hello", "h") {
		t.Fatal("expected containsrune to pass")
	}
	if stringContainsRune("hello", "z") {
		t.Fatal("expected containsrune to fail")
	}
	if stringContainsRune("hello", "") {
		t.Fatal("expected containsrune to fail for empty param")
	}
}

func TestStringExcludes(t *testing.T) {
	if !stringExcludes("hello", "world") {
		t.Fatal("expected excludes to pass")
	}
	if stringExcludes("hello world", "world") {
		t.Fatal("expected excludes to fail")
	}
}

func TestStringExcludesAll(t *testing.T) {
	if !stringExcludesAll("hello", "xyz") {
		t.Fatal("expected excludesall to pass")
	}
	if stringExcludesAll("hello", "h") {
		t.Fatal("expected excludesall to fail")
	}
}

func TestStringExcludesRune(t *testing.T) {
	if !stringExcludesRune("hello", "z") {
		t.Fatal("expected excludesrune to pass")
	}
	if stringExcludesRune("hello", "h") {
		t.Fatal("expected excludesrune to fail")
	}
	if stringExcludesRune("hello", "") {
		t.Fatal("expected excludesrune to fail for empty param")
	}
}

func TestStringLowercase(t *testing.T) {
	if !stringLowercase("hello", "") {
		t.Fatal("expected lowercase to pass")
	}
	if stringLowercase("Hello", "") {
		t.Fatal("expected lowercase to fail")
	}
	if !stringLowercase("", "") {
		t.Fatal("expected lowercase to pass for empty")
	}
	if !stringLowercase("hello123", "") {
		t.Fatal("expected lowercase to pass for lowercase+digits")
	}
	if stringLowercase("Hello123", "") {
		t.Fatal("expected lowercase to fail for mixed case+digits")
	}
}

func TestStringUppercase(t *testing.T) {
	if !stringUppercase("HELLO", "") {
		t.Fatal("expected uppercase to pass")
	}
	if stringUppercase("Hello", "") {
		t.Fatal("expected uppercase to fail")
	}
	if !stringUppercase("", "") {
		t.Fatal("expected uppercase to pass for empty")
	}
	if !stringUppercase("HELLO123", "") {
		t.Fatal("expected uppercase to pass for uppercase+digits")
	}
	if stringUppercase("Hello123", "") {
		t.Fatal("expected uppercase to fail for mixed case+digits")
	}
}

func TestStringBoolean(t *testing.T) {
	if !stringBoolean("true", "") {
		t.Fatal("expected boolean to pass for 'true'")
	}
	if !stringBoolean("1", "") {
		t.Fatal("expected boolean to pass for '1'")
	}
	if stringBoolean("maybe", "") {
		t.Fatal("expected boolean to fail for 'maybe'")
	}
}

func TestStringNumber(t *testing.T) {
	if !stringNumber("3.14", "") {
		t.Fatal("expected number to pass for numeric string")
	}
	if stringNumber("abc", "") {
		t.Fatal("expected number to fail for non-numeric string")
	}
}

func TestStringDatetime(t *testing.T) {
	if !stringDatetime("2023-12-06T00:00:00Z", "") {
		t.Fatal("expected datetime to pass for RFC3339")
	}
	if stringDatetime("invalid", "") {
		t.Fatal("expected datetime to fail")
	}
	if !stringDatetime("2023-12-06", "2006-01-02") {
		t.Fatal("expected datetime with custom layout to pass")
	}
}

func TestStringTimezone(t *testing.T) {
	if !stringTimezone("UTC", "") {
		t.Fatal("expected timezone to pass for UTC")
	}
	if stringTimezone("Invalid/Zone", "") {
		t.Fatal("expected timezone to fail")
	}
}

func TestStringLatitude(t *testing.T) {
	if !stringLatitude("45.0", "") {
		t.Fatal("expected latitude to pass")
	}
	if stringLatitude("91.0", "") {
		t.Fatal("expected latitude to fail for out-of-range")
	}
	if stringLatitude("abc", "") {
		t.Fatal("expected latitude to fail for non-numeric")
	}
}

func TestStringLongitude(t *testing.T) {
	if !stringLongitude("90.0", "") {
		t.Fatal("expected longitude to pass")
	}
	if stringLongitude("181.0", "") {
		t.Fatal("expected longitude to fail for out-of-range")
	}
	if stringLongitude("abc", "") {
		t.Fatal("expected longitude to fail for non-numeric")
	}
}

func TestStringFile(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Skip("cannot get executable path")
	}
	if !stringFile(exe, "") {
		t.Fatalf("expected file to pass for %s", exe)
	}
	if stringFile("/nonexistent/path/file.txt", "") {
		t.Fatal("expected file to fail for nonexistent")
	}
}

func TestStringFilePath(t *testing.T) {
	if !stringFilePath("/usr/local/bin/app", "") {
		t.Fatal("expected filepath to pass")
	}
	if stringFilePath("", "") {
		t.Fatal("expected filepath to fail for empty")
	}
}

func TestStringDir(t *testing.T) {
	if !stringDir(os.TempDir(), "") {
		t.Fatalf("expected dir to pass for %s", os.TempDir())
	}
	if stringDir("/nonexistent/dir", "") {
		t.Fatal("expected dir to fail for nonexistent")
	}
}

func TestStringDirPath(t *testing.T) {
	if !stringDirPath("/usr/local/bin", "") {
		t.Fatal("expected dirpath to pass")
	}
	if stringDirPath("", "") {
		t.Fatal("expected dirpath to fail for empty")
	}
	if stringDirPath(".", "") {
		t.Fatal("expected dirpath to fail for '.'")
	}
	if stringDirPath("/path/to/file.txt", "") {
		t.Fatal("expected dirpath to fail for path with dot in base")
	}
}

func TestStringMongoDB(t *testing.T) {
	if !stringMongoDB("507f1f77bcf86cd799439011", "") {
		t.Fatal("expected mongodb to pass")
	}
	if stringMongoDB("invalid", "") {
		t.Fatal("expected mongodb to fail")
	}
}

func TestStringLuhnChecksum(t *testing.T) {
	if !stringLuhnChecksum("4111111111111111", "") {
		t.Fatal("expected luhn_checksum to pass")
	}
	if stringLuhnChecksum("4111111111111112", "") {
		t.Fatal("expected luhn_checksum to fail for bad checksum")
	}
	if !stringLuhnChecksum("79927398713", "") {
		t.Fatal("expected luhn_checksum to pass for valid number")
	}
	if stringLuhnChecksum("abc", "") {
		t.Fatal("expected luhn_checksum to fail for non-digits")
	}
	if !stringLuhnChecksum("4111 1111 1111 1111", "") {
		t.Fatal("expected luhn_checksum to pass with spaces")
	}
	if !stringLuhnChecksum("4111-1111-1111-1111", "") {
		t.Fatal("expected luhn_checksum to pass with dashes")
	}
	if stringLuhnChecksum("  -  ", "") {
		t.Fatal("expected luhn_checksum to fail for no digits")
	}
}

func TestStringDNSRFC1035Label(t *testing.T) {
	if !stringDNSRFC1035Label("my-label", "") {
		t.Fatal("expected dns_rfc1035_label to pass")
	}
	if stringDNSRFC1035Label("Invalid", "") {
		t.Fatal("expected dns_rfc1035_label to fail for uppercase")
	}
	longLabel := ""
	for i := 0; i < 64; i++ {
		longLabel += "a"
	}
	if stringDNSRFC1035Label(longLabel, "") {
		t.Fatal("expected dns_rfc1035_label to fail for too long label")
	}
}

func TestStringSemver(t *testing.T) {
	if !stringSemver("1.2.3", "") {
		t.Fatal("expected semver to pass for 1.2.3")
	}
	if !stringSemver("v1.2.3", "") {
		t.Fatal("expected semver to pass for v1.2.3")
	}
	if !stringSemver("1.0.0-alpha.1", "") {
		t.Fatal("expected semver to pass for pre-release")
	}
	if !stringSemver("1.0.0+build.123", "") {
		t.Fatal("expected semver to pass with build metadata")
	}
	if stringSemver("1.2", "") {
		t.Fatal("expected semver to fail for incomplete version")
	}
	if stringSemver("abc", "") {
		t.Fatal("expected semver to fail for non-version")
	}
	if stringSemver("01.2.3", "") {
		t.Fatal("expected semver to fail for leading zero")
	}
}

func TestStringISBN10(t *testing.T) {
	if !stringISBN10("0471958697", "") {
		t.Fatal("expected isbn10 to pass for valid ISBN-10")
	}
	if !stringISBN10("0-471-95869-7", "") {
		t.Fatal("expected isbn10 to pass with dashes")
	}
	if !stringISBN10("0 471 95869 7", "") {
		t.Fatal("expected isbn10 to pass with spaces")
	}
	if stringISBN10("0471958698", "") {
		t.Fatal("expected isbn10 to fail for wrong checksum")
	}
	if stringISBN10("123456789", "") {
		t.Fatal("expected isbn10 to fail for too short")
	}
	if stringISBN10("12345678901", "") {
		t.Fatal("expected isbn10 to fail for too long")
	}
	if !stringISBN10("080442957X", "") {
		t.Fatal("expected isbn10 to pass for ISBN ending with X")
	}
}

func TestStringISBN13(t *testing.T) {
	if !stringISBN13("9780471117094", "") {
		t.Fatal("expected isbn13 to pass for valid ISBN-13")
	}
	if !stringISBN13("978-0-471-11709-4", "") {
		t.Fatal("expected isbn13 to pass with dashes")
	}
	if stringISBN13("9780471117095", "") {
		t.Fatal("expected isbn13 to fail for wrong checksum")
	}
	if stringISBN13("978047111709", "") {
		t.Fatal("expected isbn13 to fail for too short")
	}
}

func TestStringISSN(t *testing.T) {
	if !stringISSN("0317847X", "") {
		t.Fatal("expected issn to pass for valid ISSN")
	}
	if !stringISSN("0317-847X", "") {
		t.Fatal("expected issn to pass with dash")
	}
	if stringISSN("03178470", "") {
		t.Fatal("expected issn to fail for wrong checksum")
	}
	if stringISSN("1234567", "") {
		t.Fatal("expected issn to fail for too short")
	}
}

func TestStringBIC(t *testing.T) {
	if !stringBIC("CHASUS33", "") {
		t.Fatal("expected bic to pass for valid BIC")
	}
	if !stringBIC("CHASUS33XXX", "") {
		t.Fatal("expected bic to pass for 11-char BIC")
	}
	if stringBIC("CHASU", "") {
		t.Fatal("expected bic to fail for too short")
	}
	if stringBIC("12345678", "") {
		t.Fatal("expected bic to fail for numeric")
	}
}

func TestStringCron(t *testing.T) {
	if !stringCron("*/5 * * * *", "") {
		t.Fatal("expected cron to pass for standard 5-field")
	}
	if !stringCron("0 0 1 1 *", "") {
		t.Fatal("expected cron to pass for yearly")
	}
	if !stringCron("0 0 * * * *", "") {
		t.Fatal("expected cron to pass for 6-field")
	}
	if stringCron("* * * *", "") {
		t.Fatal("expected cron to fail for 4-field")
	}
	if stringCron("* * * * * * *", "") {
		t.Fatal("expected cron to fail for 7-field")
	}
}

func TestStringDataURI(t *testing.T) {
	if !stringDataURI("data:text/plain;base64,SGVsbG8=", "") {
		t.Fatal("expected datauri to pass for valid data URI")
	}
	if !stringDataURI("data:text/html,Hello", "") {
		t.Fatal("expected datauri to pass for simple data URI")
	}
	if stringDataURI("http://example.com", "") {
		t.Fatal("expected datauri to fail for non-data URI")
	}
	if stringDataURI("data:", "") {
		t.Fatal("expected datauri to fail for empty data URI")
	}
}

func TestStringBCP47(t *testing.T) {
	if !stringBCP47("en", "") {
		t.Fatal("expected bcp47 to pass for 'en'")
	}
	if !stringBCP47("zh-CN", "") {
		t.Fatal("expected bcp47 to pass for 'zh-CN'")
	}
	if !stringBCP47("en-US", "") {
		t.Fatal("expected bcp47 to pass for 'en-US'")
	}
	if !stringBCP47("sr-Latn-RS", "") {
		t.Fatal("expected bcp47 to pass for 'sr-Latn-RS'")
	}
	if stringBCP47("1", "") {
		t.Fatal("expected bcp47 to fail for single digit")
	}
}

func TestStringEthAddr(t *testing.T) {
	if !stringEthAddr("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38", "") {
		t.Fatal("expected eth_addr to pass for valid address")
	}
	if !stringEthAddr("0x0000000000000000000000000000000000000000", "") {
		t.Fatal("expected eth_addr to pass for zero address")
	}
	if stringEthAddr("742d35Cc6634C0532925a3b844Bc9e7595f2bD38", "") {
		t.Fatal("expected eth_addr to fail without 0x prefix")
	}
	if stringEthAddr("0x1234", "") {
		t.Fatal("expected eth_addr to fail for too short")
	}
}

func TestStringBtcAddr(t *testing.T) {
	if !stringBtcAddr("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "") {
		t.Fatal("expected btc_addr to pass for legacy address")
	}
	if stringBtcAddr("not-a-btc-address", "") {
		t.Fatal("expected btc_addr to fail for invalid address")
	}
	if stringBtcAddr("short", "") {
		t.Fatal("expected btc_addr to fail for too short")
	}
}

func TestStringRuleMapCompleteness(t *testing.T) {
	requiredRules := []string{
		"required", "isdefault", "min", "max", "len", "eq", "eq_ignore_case",
		"ne", "ne_ignore_case", "gt", "gte", "lt", "lte", "alpha", "alphaspace",
		"alphanum", "alphanumspace", "alphaunicode", "alphanumunicode", "ascii",
		"printascii", "multibyte", "hexadecimal", "hexcolor", "rgb", "rgba",
		"hsl", "hsla", "email", "e164", "ip", "ip_addr", "ipv4", "ipv6",
		"cidr", "cidrv4", "cidrv6", "mac", "hostname", "hostname_rfc1123",
		"fqdn", "hostname_port", "port", "url", "uri", "http_url", "https_url",
		"url_encoded", "html", "html_encoded", "uuid", "uuid3", "uuid4", "uuid5",
		"uuid_rfc4122", "uuid3_rfc4122", "uuid4_rfc4122", "uuid5_rfc4122",
		"base32", "base64", "base64url", "base64rawurl", "json", "unique",
		"startswith", "endswith", "startsnotwith", "endsnotwith", "contains",
		"containsany", "containsrune", "excludes", "excludesall", "excludesrune",
		"lowercase", "uppercase", "boolean", "number", "numeric", "datetime",
		"timezone", "latitude", "longitude", "file", "filepath", "dir", "dirpath",
		"mongodb", "luhn_checksum", "credit_card", "dns_rfc1035_label", "semver",
		"isbn10", "isbn13", "issn", "bic", "cron", "datauri", "bcp47",
		"eth_addr", "btc_addr",
	}
	for _, name := range requiredRules {
		if _, ok := stringRuleMap[name]; !ok {
			t.Errorf("stringRuleMap missing rule: %s", name)
		}
	}
}

func TestStringSemverEdgeCases(t *testing.T) {
	if !stringSemver("1.0.0-alpha", "") {
		t.Fatal("expected semver to pass for pre-release without dot")
	}
	if !stringSemver("1.0.0-alpha.beta", "") {
		t.Fatal("expected semver to pass for pre-release with dot")
	}
	if !stringSemver("1.0.0-0.3.7", "") {
		t.Fatal("expected semver to pass for numeric pre-release")
	}
	if !stringSemver("1.0.0-x.7.z.92", "") {
		t.Fatal("expected semver to pass for mixed pre-release")
	}
	if !stringSemver("1.0.0-alpha+001", "") {
		t.Fatal("expected semver to pass for pre-release with build")
	}
	if stringSemver("1.0.0-", "") {
		t.Fatal("expected semver to fail for empty pre-release")
	}
	if stringSemver("1.0.0-.", "") {
		t.Fatal("expected semver to fail for dot-only pre-release")
	}
	if stringSemver("1.0.0+", "") {
		t.Fatal("expected semver to fail for empty build metadata")
	}
	if stringSemver("1.0.0+.", "") {
		t.Fatal("expected semver to fail for dot-only build metadata")
	}
	if stringSemver("1.0.0+build.", "") {
		t.Fatal("expected semver to fail for trailing dot in build")
	}
	if stringSemver("1.0.0-alpha.", "") {
		t.Fatal("expected semver to fail for trailing dot in pre-release")
	}
	if stringSemver("1.0.0-alpha!bad", "") {
		t.Fatal("expected semver to fail for invalid char in pre-release")
	}
	if stringSemver("1.0.0+build!bad", "") {
		t.Fatal("expected semver to fail for invalid char in build")
	}
	if stringSemver("1.2.3.4", "") {
		t.Fatal("expected semver to fail for 4-part version")
	}
}

func TestStringISBN10EdgeCases(t *testing.T) {
	if stringISBN10("ABCDEFGHIJ", "") {
		t.Fatal("expected isbn10 to fail for all alpha")
	}
	if stringISBN10("", "") {
		t.Fatal("expected isbn10 to fail for empty")
	}
}

func TestStringISSNEdgeCases(t *testing.T) {
	if stringISSN("0317847", "") {
		t.Fatal("expected issn to fail for 7 chars without dash")
	}
	if stringISSN("0317-8470", "") {
		t.Fatal("expected issn to fail for wrong check digit")
	}
	if stringISSN("0317-847", "") {
		t.Fatal("expected issn to fail for short with dash")
	}
	if stringISSN("", "") {
		t.Fatal("expected issn to fail for empty")
	}
}

func TestStringBICEdgeCases(t *testing.T) {
	if stringBIC("12ASUS33", "") {
		t.Fatal("expected bic to fail for numeric bank code")
	}
	if stringBIC("CH12US33", "") {
		t.Fatal("expected bic to fail for numeric country code")
	}
	if stringBIC("CHA1US33", "") {
		t.Fatal("expected bic to fail for digit in bank code")
	}
	if stringBIC("CHASU1", "") {
		t.Fatal("expected bic to fail for digit in location code")
	}
}

func TestStringCronEdgeCases(t *testing.T) {
	if stringCron("1--2 * * * *", "") {
		t.Fatal("expected cron to fail for double dash")
	}
	if stringCron("1//2 * * * *", "") {
		t.Fatal("expected cron to fail for double slash")
	}
	if !stringCron("1-30 * * * *", "") {
		t.Fatal("expected cron to pass for valid range")
	}
	if !stringCron("1,15 * * * *", "") {
		t.Fatal("expected cron to pass for valid list")
	}
	if !stringCron("1-30/5 * * * *", "") {
		t.Fatal("expected cron to pass for valid step")
	}
	if stringCron("a * * * *", "") {
		t.Fatal("expected cron to fail for alpha char")
	}
	if stringCron("1*2 * * * *", "") {
		t.Fatal("expected cron to fail for invalid star position")
	}
}

func TestStringDataURIEdgeCases(t *testing.T) {
	if !stringDataURI("data:text/plain;charset=utf-8;base64,SGVsbG8=", "") {
		t.Fatal("expected datauri to pass with charset param")
	}
	if !stringDataURI("data:image/png;base64,iVBOR", "") {
		t.Fatal("expected datauri to pass for image")
	}
}

func TestStringBCP47EdgeCases(t *testing.T) {
	if !stringBCP47("en-GB", "") {
		t.Fatal("expected bcp47 to pass for en-GB")
	}
	if !stringBCP47("zh-Hans-CN", "") {
		t.Fatal("expected bcp47 to pass for zh-Hans-CN")
	}
	if stringBCP47("12", "") {
		t.Fatal("expected bcp47 to fail for numeric primary")
	}
}

func TestStringEthAddrEdgeCases(t *testing.T) {
	if stringEthAddr("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD3G", "") {
		t.Fatal("expected eth_addr to fail for invalid hex char")
	}
}

func TestStringBtcAddrEdgeCases(t *testing.T) {
	if !stringBtcAddr("3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy", "") {
		t.Fatal("expected btc_addr to pass for P2SH address")
	}
	if !stringBtcAddr("bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4", "") {
		t.Fatal("expected btc_addr to pass for bech32 address")
	}
	if stringBtcAddr("bc1q!nvalid", "") {
		t.Fatal("expected btc_addr to fail for invalid bech32")
	}
}

func TestStringURLEdgeCases(t *testing.T) {
	if stringURI("://no-scheme", "") {
		t.Fatal("expected uri to fail for no scheme")
	}
	if stringURI("ht tp://example.com", "") {
		t.Fatal("expected uri to fail for space in scheme")
	}
	if stringURI("http:", "") {
		t.Fatal("expected uri to fail for scheme only")
	}
}

func TestStringSemverParseFailures(t *testing.T) {
	if stringSemver("1.2", "") {
		t.Fatal("expected semver to fail for missing patch")
	}
	if stringSemver("1.0.0-alpha.", "") {
		t.Fatal("expected semver to fail for trailing dot in pre-release")
	}
	if stringSemver("1.0.0+build.", "") {
		t.Fatal("expected semver to fail for trailing dot in build")
	}
	if stringSemver("1.0.0-", "") {
		t.Fatal("expected semver to fail for empty pre-release")
	}
	if stringSemver("1.0.0+", "") {
		t.Fatal("expected semver to fail for empty build")
	}
	if stringSemver("1.0.0-alpha!bad", "") {
		t.Fatal("expected semver to fail for invalid pre-release char")
	}
	if stringSemver("1.0.0+build!bad", "") {
		t.Fatal("expected semver to fail for invalid build char")
	}
}

func TestStringISBN10EdgeCases2(t *testing.T) {
	if stringISBN10("ABCDEFGHIJ", "") {
		t.Fatal("expected isbn10 to fail for all alpha")
	}
}

func TestStringISBN13EdgeCases(t *testing.T) {
	if stringISBN13("978047111709A", "") {
		t.Fatal("expected isbn13 to fail for non-digit")
	}
	if stringISBN13("978047111709", "") {
		t.Fatal("expected isbn13 to fail for too short")
	}
	if stringISBN13("978-0-471-11709-4-extra", "") {
		t.Fatal("expected isbn13 to fail for too many chars")
	}
}

func TestStringISSNEdgeCases2(t *testing.T) {
	if stringISSN("03178471", "") {
		t.Fatal("expected issn to fail for wrong check digit")
	}
	if stringISSN("031X8471", "") {
		t.Fatal("expected issn to fail for non-digit in middle")
	}
	if stringISSN("0317847", "") {
		t.Fatal("expected issn to fail for too short")
	}
}

func TestStringBICEdgeCases2(t *testing.T) {
	if stringBIC("CH1SUS33", "") {
		t.Fatal("expected bic to fail for digit in bank code")
	}
	if stringBIC("CHASU133", "") {
		t.Fatal("expected bic to fail for digit in location code")
	}
	if stringBIC("CHASUS3!", "") {
		t.Fatal("expected bic to fail for invalid char in branch")
	}
}

func TestStringEthAddrEdgeCases2(t *testing.T) {
	if !stringEthAddr("0X742d35Cc6634C0532925a3b844Bc9e7595f2bD38", "") {
		t.Fatal("expected eth_addr to pass for 0X prefix")
	}
}

func TestStringSemverMinorFail(t *testing.T) {
	if stringSemver("1.", "") {
		t.Fatal("expected semver to fail for missing minor")
	}
}

func TestStringISBN13TooManyDigits(t *testing.T) {
	if stringISBN13("978-0-471-11709-4-extra", "") {
		t.Fatal("expected isbn13 to fail for too many digits")
	}
}

func TestStringISSNTooManyDigits(t *testing.T) {
	if stringISSN("031784712", "") {
		t.Fatal("expected issn to fail for too many digits")
	}
}

func TestStringISSNCheckDigitNonDigitNonX(t *testing.T) {
	if stringISSN("0317847A", "") {
		t.Fatal("expected issn to fail for non-digit non-X check digit")
	}
}

func TestStringEthAddrNo0xPrefix(t *testing.T) {
	if stringEthAddr("1x742d35Cc6634C0532925a3b844Bc9e7595f2bD38", "") {
		t.Fatal("expected eth_addr to fail without 0x prefix")
	}
}
