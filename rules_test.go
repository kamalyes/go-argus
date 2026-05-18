/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 20:10:55
 * @FilePath: \go-argus\rules_test.go
 * @Description: rules.go 测试，覆盖所有内置字段规则和辅助函数
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"os"
	"reflect"
	"testing"
)

func TestRuleRequired(t *testing.T) {
	v := New(WithRequiredStructEnabled())
	if err := v.Var("", "required"); err == nil {
		t.Fatal("expected empty string to fail required")
	}
	if err := v.Var("hello", "required"); err != nil {
		t.Fatal("expected non-empty string to pass required")
	}
}

func TestRuleIsDefault(t *testing.T) {
	v := New()
	if err := v.Var("", "isdefault"); err != nil {
		t.Fatal("expected empty string to pass isdefault")
	}
	if err := v.Var("hello", "isdefault"); err == nil {
		t.Fatal("expected non-empty string to fail isdefault")
	}
}

func TestRuleMinMaxLen(t *testing.T) {
	v := New()
	if err := v.Var("abc", "min=2,max=5"); err != nil {
		t.Fatal("expected abc to pass min=2,max=5")
	}
	if err := v.Var("a", "min=2"); err == nil {
		t.Fatal("expected 'a' to fail min=2")
	}
	if err := v.Var("abcdef", "max=5"); err == nil {
		t.Fatal("expected 'abcdef' to fail max=5")
	}
	if err := v.Var("abc", "len=3"); err != nil {
		t.Fatal("expected 'abc' to pass len=3")
	}
	if err := v.Var("ab", "len=3"); err == nil {
		t.Fatal("expected 'ab' to fail len=3")
	}
}

func TestRuleMinMaxLenSlice(t *testing.T) {
	v := New()
	if err := v.Var([]string{"a", "b"}, "min=1,max=3"); err != nil {
		t.Fatal("expected slice to pass min=1,max=3")
	}
	if err := v.Var([]string{}, "min=1"); err == nil {
		t.Fatal("expected empty slice to fail min=1")
	}
}

func TestRuleMinMaxNumeric(t *testing.T) {
	v := New()
	if err := v.Var(5, "min=1,max=10"); err != nil {
		t.Fatal("expected 5 to pass min=1,max=10")
	}
	if err := v.Var(0, "min=1"); err == nil {
		t.Fatal("expected 0 to fail min=1")
	}
	if err := v.Var(11, "max=10"); err == nil {
		t.Fatal("expected 11 to fail max=10")
	}
}

func TestRuleEqNe(t *testing.T) {
	v := New()
	if err := v.Var("hello", "eq=hello"); err != nil {
		t.Fatal("expected eq=hello to pass")
	}
	if err := v.Var("hello", "eq=world"); err == nil {
		t.Fatal("expected eq=world to fail")
	}
	if err := v.Var("hello", "ne=world"); err != nil {
		t.Fatal("expected ne=world to pass")
	}
	if err := v.Var("hello", "ne=hello"); err == nil {
		t.Fatal("expected ne=hello to fail")
	}
}

func TestRuleEqIgnoreCaseNeIgnoreCase(t *testing.T) {
	v := New()
	if err := v.Var("Hello", "eq_ignore_case=hello"); err != nil {
		t.Fatal("expected eq_ignore_case to pass")
	}
	if err := v.Var("Hello", "ne_ignore_case=world"); err != nil {
		t.Fatal("expected ne_ignore_case to pass")
	}
	if err := v.Var("Hello", "ne_ignore_case=hello"); err == nil {
		t.Fatal("expected ne_ignore_case=hello to fail")
	}
}

func TestRuleGtGteLtLte(t *testing.T) {
	v := New()
	if err := v.Var(10, "gt=5"); err != nil {
		t.Fatal("expected 10 > 5")
	}
	if err := v.Var(5, "gt=5"); err == nil {
		t.Fatal("expected 5 not > 5")
	}
	if err := v.Var(5, "gte=5"); err != nil {
		t.Fatal("expected 5 >= 5")
	}
	if err := v.Var(4, "gte=5"); err == nil {
		t.Fatal("expected 4 not >= 5")
	}
	if err := v.Var(3, "lt=5"); err != nil {
		t.Fatal("expected 3 < 5")
	}
	if err := v.Var(5, "lt=5"); err == nil {
		t.Fatal("expected 5 not < 5")
	}
	if err := v.Var(5, "lte=5"); err != nil {
		t.Fatal("expected 5 <= 5")
	}
	if err := v.Var(6, "lte=5"); err == nil {
		t.Fatal("expected 6 not <= 5")
	}
}

func TestRuleAlpha(t *testing.T) {
	v := New()
	if err := v.Var("abc", "alpha"); err != nil {
		t.Fatal("expected alpha to pass")
	}
	if err := v.Var("abc123", "alpha"); err == nil {
		t.Fatal("expected alpha to fail for alphanum")
	}
}

func TestRuleAlphaSpace(t *testing.T) {
	v := New()
	if err := v.Var("hello world", "alphaspace"); err != nil {
		t.Fatal("expected alphaspace to pass")
	}
	if err := v.Var("hello123", "alphaspace"); err == nil {
		t.Fatal("expected alphaspace to fail for digits")
	}
}

func TestRuleAlphanum(t *testing.T) {
	v := New()
	if err := v.Var("abc123", "alphanum"); err != nil {
		t.Fatal("expected alphanum to pass")
	}
	if err := v.Var("abc-123", "alphanum"); err == nil {
		t.Fatal("expected alphanum to fail for hyphen")
	}
}

func TestRuleAlphanumSpace(t *testing.T) {
	v := New()
	if err := v.Var("abc 123", "alphanumspace"); err != nil {
		t.Fatal("expected alphanumspace to pass")
	}
}

func TestRuleAlphaUnicode(t *testing.T) {
	v := New()
	if err := v.Var("你好世界", "alphaunicode"); err != nil {
		t.Fatal("expected alphaunicode to pass")
	}
}

func TestRuleAlphanumUnicode(t *testing.T) {
	v := New()
	if err := v.Var("你好123", "alphanumunicode"); err != nil {
		t.Fatal("expected alphanumunicode to pass")
	}
}

func TestRuleASCII(t *testing.T) {
	v := New()
	if err := v.Var("hello", "ascii"); err != nil {
		t.Fatal("expected ascii to pass")
	}
	if err := v.Var("你好", "ascii"); err == nil {
		t.Fatal("expected ascii to fail for unicode")
	}
}

func TestRulePrintASCII(t *testing.T) {
	v := New()
	if err := v.Var("hello world", "printascii"); err != nil {
		t.Fatal("expected printascii to pass")
	}
}

func TestRuleMultibyte(t *testing.T) {
	v := New()
	if err := v.Var("你好", "multibyte"); err != nil {
		t.Fatal("expected multibyte to pass")
	}
	if err := v.Var("hello", "multibyte"); err == nil {
		t.Fatal("expected multibyte to fail for ascii")
	}
}

func TestRuleHexadecimal(t *testing.T) {
	v := New()
	if err := v.Var("deadbeef", "hexadecimal"); err != nil {
		t.Fatal("expected hexadecimal to pass")
	}
	if err := v.Var("xyz", "hexadecimal"); err == nil {
		t.Fatal("expected hexadecimal to fail")
	}
}

func TestRuleHexColor(t *testing.T) {
	v := New()
	if err := v.Var("#12ffaa", "hexcolor"); err != nil {
		t.Fatal("expected hexcolor to pass")
	}
	if err := v.Var("12ffaa", "hexcolor"); err != nil {
		t.Fatal("expected hexcolor without # to pass")
	}
}

func TestRuleRGB(t *testing.T) {
	v := New()
	if err := v.Var("rgb(12, 34, 255)", "rgb"); err != nil {
		t.Fatal("expected rgb to pass")
	}
	if err := v.Var("rgb(300, 34, 255)", "rgb"); err == nil {
		t.Fatal("expected rgb to fail for out-of-range")
	}
}

func TestRuleRGBA(t *testing.T) {
	v := New()
	if err := v.Var("rgba(12, 34, 255, 0.5)", "rgba"); err != nil {
		t.Fatal("expected rgba to pass")
	}
}

func TestRuleHSL(t *testing.T) {
	v := New()
	if err := v.Var("hsl(120, 50%, 75%)", "hsl"); err != nil {
		t.Fatal("expected hsl to pass")
	}
}

func TestRuleHSLA(t *testing.T) {
	v := New()
	if err := v.Var("hsla(120, 50%, 75%, 0.5)", "hsla"); err != nil {
		t.Fatal("expected hsla to pass")
	}
}

func TestRuleEmail(t *testing.T) {
	v := New()
	if err := v.Var("user@example.com", "email"); err != nil {
		t.Fatal("expected email to pass")
	}
	if err := v.Var("not-email", "email"); err == nil {
		t.Fatal("expected email to fail")
	}
}

func TestRuleE164(t *testing.T) {
	v := New()
	if err := v.Var("+8613800138000", "e164"); err != nil {
		t.Fatal("expected e164 to pass")
	}
	if err := v.Var("12345", "e164"); err == nil {
		t.Fatal("expected e164 to fail")
	}
}

func TestRuleIP(t *testing.T) {
	v := New()
	if err := v.Var("192.168.1.1", "ip"); err != nil {
		t.Fatal("expected ip to pass")
	}
	if err := v.Var("::1", "ip"); err != nil {
		t.Fatal("expected ipv6 to pass ip")
	}
	if err := v.Var("not-an-ip", "ip"); err == nil {
		t.Fatal("expected ip to fail")
	}
}

func TestRuleIPv4(t *testing.T) {
	v := New()
	if err := v.Var("192.168.1.1", "ipv4"); err != nil {
		t.Fatal("expected ipv4 to pass")
	}
	if err := v.Var("::1", "ipv4"); err == nil {
		t.Fatal("expected ipv6 to fail ipv4")
	}
}

func TestRuleIPv6(t *testing.T) {
	v := New()
	if err := v.Var("::1", "ipv6"); err != nil {
		t.Fatal("expected ipv6 to pass")
	}
	if err := v.Var("192.168.1.1", "ipv6"); err == nil {
		t.Fatal("expected ipv4 to fail ipv6")
	}
}

func TestRuleCIDR(t *testing.T) {
	v := New()
	if err := v.Var("192.168.1.0/24", "cidr"); err != nil {
		t.Fatal("expected cidr to pass")
	}
	if err := v.Var("invalid", "cidr"); err == nil {
		t.Fatal("expected cidr to fail")
	}
}

func TestRuleCIDRv4(t *testing.T) {
	v := New()
	if err := v.Var("192.168.1.0/24", "cidrv4"); err != nil {
		t.Fatal("expected cidrv4 to pass")
	}
	if err := v.Var("::1/128", "cidrv4"); err == nil {
		t.Fatal("expected ipv6 cidr to fail cidrv4")
	}
}

func TestRuleCIDRv6(t *testing.T) {
	v := New()
	if err := v.Var("::1/128", "cidrv6"); err != nil {
		t.Fatal("expected cidrv6 to pass")
	}
	if err := v.Var("192.168.1.0/24", "cidrv6"); err == nil {
		t.Fatal("expected ipv4 cidr to fail cidrv6")
	}
}

func TestRuleMAC(t *testing.T) {
	v := New()
	if err := v.Var("00:11:22:33:44:55", "mac"); err != nil {
		t.Fatal("expected mac to pass")
	}
	if err := v.Var("invalid", "mac"); err == nil {
		t.Fatal("expected mac to fail")
	}
}

func TestRuleHostname(t *testing.T) {
	v := New()
	if err := v.Var("api.example.com", "hostname"); err != nil {
		t.Fatal("expected hostname to pass")
	}
	if err := v.Var("-invalid.com", "hostname"); err == nil {
		t.Fatal("expected hostname to fail")
	}
}

func TestRuleFQDN(t *testing.T) {
	v := New()
	if err := v.Var("api.example.com.", "fqdn"); err != nil {
		t.Fatal("expected fqdn to pass")
	}
	if err := v.Var("api.example.com", "fqdn"); err == nil {
		t.Fatal("expected fqdn to fail without trailing dot")
	}
}

func TestRuleHostnamePort(t *testing.T) {
	v := New()
	if err := v.Var("example.com:8080", "hostname_port"); err != nil {
		t.Fatal("expected hostname_port to pass")
	}
	if err := v.Var("example.com:99999", "hostname_port"); err == nil {
		t.Fatal("expected hostname_port to fail for invalid port")
	}
}

func TestRulePort(t *testing.T) {
	v := New()
	if err := v.Var("443", "port"); err != nil {
		t.Fatal("expected port to pass")
	}
	if err := v.Var("99999", "port"); err == nil {
		t.Fatal("expected port to fail for out-of-range")
	}
}

func TestRuleURL(t *testing.T) {
	v := New()
	if err := v.Var("https://example.com/path", "url"); err != nil {
		t.Fatal("expected url to pass")
	}
	if err := v.Var("not-a-url", "url"); err == nil {
		t.Fatal("expected url to fail")
	}
}

func TestRuleURI(t *testing.T) {
	v := New()
	if err := v.Var("https://example.com/path", "uri"); err != nil {
		t.Fatal("expected uri to pass")
	}
}

func TestRuleHTTPURL(t *testing.T) {
	v := New()
	if err := v.Var("https://example.com", "http_url"); err != nil {
		t.Fatal("expected http_url to pass for https")
	}
	if err := v.Var("http://example.com", "http_url"); err != nil {
		t.Fatal("expected http_url to pass for http")
	}
	if err := v.Var("ftp://example.com", "http_url"); err == nil {
		t.Fatal("expected http_url to fail for ftp")
	}
}

func TestRuleHTTPSURL(t *testing.T) {
	v := New()
	if err := v.Var("https://example.com", "https_url"); err != nil {
		t.Fatal("expected https_url to pass")
	}
	if err := v.Var("http://example.com", "https_url"); err == nil {
		t.Fatal("expected https_url to fail for http")
	}
}

func TestRuleURLEncoded(t *testing.T) {
	v := New()
	if err := v.Var("hello%20world", "url_encoded"); err != nil {
		t.Fatal("expected url_encoded to pass")
	}
	if err := v.Var("hello%ZZ", "url_encoded"); err == nil {
		t.Fatal("expected url_encoded to fail for invalid encoding")
	}
	if err := v.Var("nopercent", "url_encoded"); err == nil {
		t.Fatal("expected url_encoded to fail without percent")
	}
}

func TestRuleHTML(t *testing.T) {
	v := New()
	if err := v.Var("<b>bold</b>", "html"); err != nil {
		t.Fatal("expected html to pass")
	}
	if err := v.Var("no html", "html"); err == nil {
		t.Fatal("expected html to fail")
	}
}

func TestRuleHTMLEncoded(t *testing.T) {
	v := New()
	if err := v.Var("&amp;", "html_encoded"); err != nil {
		t.Fatal("expected html_encoded to pass")
	}
	if err := v.Var("plain", "html_encoded"); err == nil {
		t.Fatal("expected html_encoded to fail for plain text")
	}
}

func TestRuleUUID(t *testing.T) {
	v := New()
	if err := v.Var("550e8400-e29b-41d4-a716-446655440000", "uuid"); err != nil {
		t.Fatal("expected uuid to pass")
	}
	if err := v.Var("not-a-uuid", "uuid"); err == nil {
		t.Fatal("expected uuid to fail")
	}
}

func TestRuleUUID3(t *testing.T) {
	v := New()
	if err := v.Var("550e8400-e29b-31d4-a716-446655440000", "uuid3"); err != nil {
		t.Fatal("expected uuid3 to pass")
	}
	if err := v.Var("550e8400-e29b-41d4-a716-446655440000", "uuid3"); err == nil {
		t.Fatal("expected uuid3 to fail for uuid4")
	}
}

func TestRuleUUID4(t *testing.T) {
	v := New()
	if err := v.Var("550e8400-e29b-41d4-a716-446655440000", "uuid4"); err != nil {
		t.Fatal("expected uuid4 to pass")
	}
}

func TestRuleUUID5(t *testing.T) {
	v := New()
	if err := v.Var("550e8400-e29b-51d4-a716-446655440000", "uuid5"); err != nil {
		t.Fatal("expected uuid5 to pass")
	}
}

func TestRuleBase32(t *testing.T) {
	v := New()
	if err := v.Var("JBSWY3DPEB3W64TMMQ======", "base32"); err != nil {
		t.Fatal("expected base32 to pass")
	}
	if err := v.Var("!!!invalid!!!", "base32"); err == nil {
		t.Fatal("expected base32 to fail")
	}
	if err := v.Var("", "base32"); err == nil {
		t.Fatal("expected base32 to fail for empty")
	}
}

func TestRuleBase64(t *testing.T) {
	v := New()
	if err := v.Var("YXJndXM=", "base64"); err != nil {
		t.Fatal("expected base64 to pass")
	}
	if err := v.Var("!!!invalid!!!", "base64"); err == nil {
		t.Fatal("expected base64 to fail")
	}
}

func TestRuleBase64URL(t *testing.T) {
	v := New()
	if err := v.Var("YXJndXM=", "base64url"); err != nil {
		t.Fatal("expected base64url to pass")
	}
	if err := v.Var("", "base64url"); err == nil {
		t.Fatal("expected base64url to fail for empty")
	}
}

func TestRuleBase64RawURL(t *testing.T) {
	v := New()
	if err := v.Var("YXJndXM", "base64rawurl"); err != nil {
		t.Fatal("expected base64rawurl to pass")
	}
	if err := v.Var("", "base64rawurl"); err == nil {
		t.Fatal("expected base64rawurl to fail for empty")
	}
}

func TestRuleJSON(t *testing.T) {
	v := New()
	if err := v.Var(`{"key":"value"}`, "json"); err != nil {
		t.Fatal("expected json to pass")
	}
	if err := v.Var("{invalid}", "json"); err == nil {
		t.Fatal("expected json to fail")
	}
}

func TestRuleJSONBytes(t *testing.T) {
	v := New()
	if err := v.Var([]byte(`{"key":"value"}`), "json"); err != nil {
		t.Fatal("expected json bytes to pass")
	}
}

func TestRuleOneOf(t *testing.T) {
	v := New()
	if err := v.Var("admin", "oneof=admin member guest"); err != nil {
		t.Fatal("expected oneof to pass")
	}
	if err := v.Var("root", "oneof=admin member guest"); err == nil {
		t.Fatal("expected oneof to fail")
	}
}

func TestRuleOneOfCI(t *testing.T) {
	v := New()
	if err := v.Var("Admin", "oneofci=admin member guest"); err != nil {
		t.Fatal("expected oneofci to pass")
	}
}

func TestRuleNoneOf(t *testing.T) {
	v := New()
	if err := v.Var("root", "noneof=admin member guest"); err != nil {
		t.Fatal("expected noneof to pass")
	}
	if err := v.Var("admin", "noneof=admin member guest"); err == nil {
		t.Fatal("expected noneof to fail")
	}
}

func TestRuleNoneOfCI(t *testing.T) {
	v := New()
	if err := v.Var("Admin", "noneofci=admin member guest"); err == nil {
		t.Fatal("expected noneofci to fail")
	}
}

func TestRuleUnique(t *testing.T) {
	v := New()
	if err := v.Var("abcdef", "unique"); err != nil {
		t.Fatal("expected unique string to pass")
	}
	if err := v.Var("aabc", "unique"); err == nil {
		t.Fatal("expected non-unique string to fail")
	}
	if err := v.Var([]string{"a", "b", "c"}, "unique"); err != nil {
		t.Fatal("expected unique slice to pass")
	}
	if err := v.Var([]string{"a", "a"}, "unique"); err == nil {
		t.Fatal("expected non-unique slice to fail")
	}
}

func TestRuleUniqueMapType(t *testing.T) {
	v := New()
	m := map[string]string{"a": "1", "b": "1"}
	if err := v.Var(m, "unique"); err == nil {
		t.Fatal("expected unique to fail for map with duplicate values")
	}
}

func TestRuleStartsWith(t *testing.T) {
	v := New()
	if err := v.Var("hello world", "startswith=hello"); err != nil {
		t.Fatal("expected startswith to pass")
	}
	if err := v.Var("hello world", "startswith=world"); err == nil {
		t.Fatal("expected startswith to fail")
	}
}

func TestRuleEndsWith(t *testing.T) {
	v := New()
	if err := v.Var("hello world", "endswith=world"); err != nil {
		t.Fatal("expected endswith to pass")
	}
	if err := v.Var("hello world", "endswith=hello"); err == nil {
		t.Fatal("expected endswith to fail")
	}
}

func TestRuleStartsNotWith(t *testing.T) {
	v := New()
	if err := v.Var("hello world", "startsnotwith=world"); err != nil {
		t.Fatal("expected startsnotwith to pass")
	}
	if err := v.Var("hello world", "startsnotwith=hello"); err == nil {
		t.Fatal("expected startsnotwith to fail")
	}
}

func TestRuleEndsNotWith(t *testing.T) {
	v := New()
	if err := v.Var("hello world", "endsnotwith=hello"); err != nil {
		t.Fatal("expected endsnotwith to pass")
	}
	if err := v.Var("hello world", "endsnotwith=world"); err == nil {
		t.Fatal("expected endsnotwith to fail")
	}
}

func TestRuleContains(t *testing.T) {
	v := New()
	if err := v.Var("hello world", "contains=world"); err != nil {
		t.Fatal("expected contains to pass")
	}
	if err := v.Var("hello world", "contains=xyz"); err == nil {
		t.Fatal("expected contains to fail")
	}
}

func TestRuleContainsAny(t *testing.T) {
	v := New()
	if err := v.Var("hello", "containsany=xyz"); err == nil {
		t.Fatal("expected containsany to fail")
	}
	if err := v.Var("hello", "containsany=h"); err != nil {
		t.Fatal("expected containsany to pass")
	}
}

func TestRuleContainsRune(t *testing.T) {
	v := New()
	if err := v.Var("hello", "containsrune=h"); err != nil {
		t.Fatal("expected containsrune to pass")
	}
	if err := v.Var("hello", "containsrune=z"); err == nil {
		t.Fatal("expected containsrune to fail")
	}
}

func TestRuleExcludes(t *testing.T) {
	v := New()
	if err := v.Var("hello", "excludes=world"); err != nil {
		t.Fatal("expected excludes to pass")
	}
	if err := v.Var("hello world", "excludes=world"); err == nil {
		t.Fatal("expected excludes to fail")
	}
}

func TestRuleExcludesAll(t *testing.T) {
	v := New()
	if err := v.Var("hello", "excludesall=xyz"); err != nil {
		t.Fatal("expected excludesall to pass")
	}
	if err := v.Var("hello", "excludesall=h"); err == nil {
		t.Fatal("expected excludesall to fail")
	}
}

func TestRuleExcludesRune(t *testing.T) {
	v := New()
	if err := v.Var("hello", "excludesrune=z"); err != nil {
		t.Fatal("expected excludesrune to pass")
	}
	if err := v.Var("hello", "excludesrune=h"); err == nil {
		t.Fatal("expected excludesrune to fail")
	}
}

func TestRuleLowercase(t *testing.T) {
	v := New()
	if err := v.Var("hello", "lowercase"); err != nil {
		t.Fatal("expected lowercase to pass")
	}
	if err := v.Var("Hello", "lowercase"); err == nil {
		t.Fatal("expected lowercase to fail")
	}
}

func TestRuleUppercase(t *testing.T) {
	v := New()
	if err := v.Var("HELLO", "uppercase"); err != nil {
		t.Fatal("expected uppercase to pass")
	}
	if err := v.Var("Hello", "uppercase"); err == nil {
		t.Fatal("expected uppercase to fail")
	}
}

func TestRuleBoolean(t *testing.T) {
	v := New()
	if err := v.Var("true", "boolean"); err != nil {
		t.Fatal("expected boolean to pass for 'true'")
	}
	if err := v.Var("1", "boolean"); err != nil {
		t.Fatal("expected boolean to pass for '1'")
	}
	if err := v.Var("maybe", "boolean"); err == nil {
		t.Fatal("expected boolean to fail for 'maybe'")
	}
}

func TestRuleBooleanType(t *testing.T) {
	v := New()
	if err := v.Var(true, "boolean"); err != nil {
		t.Fatal("expected boolean to pass for bool type")
	}
}

func TestRuleNumber(t *testing.T) {
	v := New()
	if err := v.Var(42, "number"); err != nil {
		t.Fatal("expected number to pass for int")
	}
	if err := v.Var("3.14", "number"); err != nil {
		t.Fatal("expected number to pass for numeric string")
	}
	if err := v.Var("abc", "number"); err == nil {
		t.Fatal("expected number to fail for non-numeric string")
	}
}

func TestRuleDatetime(t *testing.T) {
	v := New()
	if err := v.Var("2023-12-06T00:00:00Z", "datetime"); err != nil {
		t.Fatal("expected datetime to pass for RFC3339")
	}
	if err := v.Var("invalid", "datetime"); err == nil {
		t.Fatal("expected datetime to fail")
	}
	if err := v.Var("2023-12-06", "datetime=2006-01-02"); err != nil {
		t.Fatal("expected datetime with custom layout to pass")
	}
}

func TestRuleTimezone(t *testing.T) {
	v := New()
	if err := v.Var("UTC", "timezone"); err != nil {
		t.Fatal("expected timezone to pass for UTC")
	}
	if err := v.Var("Invalid/Zone", "timezone"); err == nil {
		t.Fatal("expected timezone to fail")
	}
}

func TestRuleLatitude(t *testing.T) {
	v := New()
	if err := v.Var(45.0, "latitude"); err != nil {
		t.Fatal("expected latitude to pass")
	}
	if err := v.Var(91.0, "latitude"); err == nil {
		t.Fatal("expected latitude to fail for out-of-range")
	}
}

func TestRuleLongitude(t *testing.T) {
	v := New()
	if err := v.Var(90.0, "longitude"); err != nil {
		t.Fatal("expected longitude to pass")
	}
	if err := v.Var(181.0, "longitude"); err == nil {
		t.Fatal("expected longitude to fail for out-of-range")
	}
}

func TestRuleFile(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Skip("cannot get executable path")
	}
	v := New()
	if err := v.Var(exe, "file"); err != nil {
		t.Fatalf("expected file to pass for %s: %v", exe, err)
	}
	if err := v.Var("/nonexistent/path/file.txt", "file"); err == nil {
		t.Fatal("expected file to fail for nonexistent")
	}
}

func TestRuleFilePath(t *testing.T) {
	v := New()
	if err := v.Var("/usr/local/bin/app", "filepath"); err != nil {
		t.Fatal("expected filepath to pass")
	}
	if err := v.Var("", "filepath"); err == nil {
		t.Fatal("expected filepath to fail for empty")
	}
}

func TestRuleDir(t *testing.T) {
	v := New()
	if err := v.Var(os.TempDir(), "dir"); err != nil {
		t.Fatalf("expected dir to pass for %s", os.TempDir())
	}
	if err := v.Var("/nonexistent/dir", "dir"); err == nil {
		t.Fatal("expected dir to fail for nonexistent")
	}
}

func TestRuleDirPath(t *testing.T) {
	v := New()
	if err := v.Var("/usr/local/bin", "dirpath"); err != nil {
		t.Fatal("expected dirpath to pass")
	}
	if err := v.Var("", "dirpath"); err == nil {
		t.Fatal("expected dirpath to fail for empty")
	}
}

func TestRuleMongoDB(t *testing.T) {
	v := New()
	if err := v.Var("507f1f77bcf86cd799439011", "mongodb"); err != nil {
		t.Fatal("expected mongodb to pass")
	}
	if err := v.Var("invalid", "mongodb"); err == nil {
		t.Fatal("expected mongodb to fail")
	}
}

func TestRuleLuhnChecksum(t *testing.T) {
	v := New()
	if err := v.Var("4111111111111111", "credit_card"); err != nil {
		t.Fatal("expected credit_card to pass")
	}
	if err := v.Var("4111111111111112", "credit_card"); err == nil {
		t.Fatal("expected credit_card to fail for bad checksum")
	}
	if err := v.Var("79927398713", "luhn_checksum"); err != nil {
		t.Fatal("expected luhn_checksum to pass")
	}
	if err := v.Var("abc", "luhn_checksum"); err == nil {
		t.Fatal("expected luhn_checksum to fail for non-digits")
	}
}

func TestRuleDNSRFC1035Label(t *testing.T) {
	v := New()
	if err := v.Var("my-label", "dns_rfc1035_label"); err != nil {
		t.Fatal("expected dns_rfc1035_label to pass")
	}
	if err := v.Var("Invalid", "dns_rfc1035_label"); err == nil {
		t.Fatal("expected dns_rfc1035_label to fail for uppercase")
	}
}

func TestCompareLengthOrNumberInvalidField(t *testing.T) {
	v := New()
	if err := v.Var(true, "min=1"); err == nil {
		t.Fatal("expected min to fail for bool type")
	}
}

func TestNumericValueInvalid(t *testing.T) {
	v := New()
	if err := v.Var(true, "latitude"); err == nil {
		t.Fatal("expected latitude to fail for bool type")
	}
}

func TestScalarStringFromBool(t *testing.T) {
	v := New()
	if err := v.Var(true, "eq=true"); err != nil {
		t.Fatal("expected eq=true to pass for bool")
	}
}

func TestScalarStringFromInt(t *testing.T) {
	v := New()
	if err := v.Var(42, "eq=42"); err != nil {
		t.Fatal("expected eq=42 to pass for int")
	}
}

func TestScalarStringFromUint(t *testing.T) {
	v := New()
	if err := v.Var(uint(42), "eq=42"); err != nil {
		t.Fatal("expected eq=42 to pass for uint")
	}
}

func TestScalarStringFromFloat32(t *testing.T) {
	v := New()
	if err := v.Var(float32(3.14), "eq=3.14"); err != nil {
		t.Fatal("expected eq=3.14 to pass for float32")
	}
}

func TestScalarStringFromFloat64(t *testing.T) {
	v := New()
	if err := v.Var(3.14, "eq=3.14"); err != nil {
		t.Fatal("expected eq=3.14 to pass for float64")
	}
}

func TestMatchStringRunesEmpty(t *testing.T) {
	v := New()
	if err := v.Var("", "alpha"); err == nil {
		t.Fatal("expected alpha to fail for empty string")
	}
}

func TestIsHostnameTooLong(t *testing.T) {
	longLabel := ""
	for i := 0; i < 254; i++ {
		longLabel += "a"
	}
	v := New()
	if err := v.Var(longLabel+".com", "hostname"); err == nil {
		t.Fatal("expected hostname to fail for too long host")
	}
}

func TestIsHostnameEmptyLabel(t *testing.T) {
	v := New()
	if err := v.Var("api..example.com", "hostname"); err == nil {
		t.Fatal("expected hostname to fail for empty label")
	}
}

func TestRuleDatetimeNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "datetime"); err == nil {
		t.Fatal("expected datetime to fail for non-string")
	}
}

func TestRuleTimezoneNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "timezone"); err == nil {
		t.Fatal("expected timezone to fail for non-string")
	}
}

func TestRuleNumberInvalidType(t *testing.T) {
	v := New()
	if err := v.Var([]int{1, 2}, "number"); err == nil {
		t.Fatal("expected number to fail for slice type")
	}
}

func TestRuleUniqueInvalidType(t *testing.T) {
	v := New()
	if err := v.Var(42, "unique"); err == nil {
		t.Fatal("expected unique to fail for int type")
	}
}

func TestRuleUniqueInvalidValue(t *testing.T) {
	v := New()
	if err := v.Var(nil, "unique"); err == nil {
		t.Fatal("expected unique to fail for nil")
	}
}

func TestRuleContainsRuneInvalidParam(t *testing.T) {
	v := New()
	if err := v.Var("hello", "containsrune="); err == nil {
		t.Fatal("expected containsrune to fail for empty param")
	}
}

func TestRuleExcludesRuneInvalidParam(t *testing.T) {
	v := New()
	if err := v.Var("hello", "excludesrune="); err == nil {
		t.Fatal("expected excludesrune to fail for empty param")
	}
}

func TestRuleHostnameNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "hostname"); err == nil {
		t.Fatal("expected hostname to fail for non-string")
	}
}

func TestRuleFQDNNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "fqdn"); err == nil {
		t.Fatal("expected fqdn to fail for non-string")
	}
}

func TestRuleHostnamePortNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "hostname_port"); err == nil {
		t.Fatal("expected hostname_port to fail for non-string")
	}
}

func TestRuleHostnamePortInvalidFormat(t *testing.T) {
	v := New()
	if err := v.Var("noport", "hostname_port"); err == nil {
		t.Fatal("expected hostname_port to fail for missing port")
	}
}

func TestRulePortNonString(t *testing.T) {
	v := New()
	if err := v.Var(443, "port"); err != nil {
		t.Fatalf("expected port to pass for int via scalarString: %v", err)
	}
}

func TestRuleURLNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "url"); err == nil {
		t.Fatal("expected url to fail for non-string")
	}
}

func TestRuleHTTPURLNoHost(t *testing.T) {
	v := New()
	if err := v.Var("https://", "http_url"); err == nil {
		t.Fatal("expected http_url to fail for no host")
	}
}

func TestRuleHTTPSURLNoHost(t *testing.T) {
	v := New()
	if err := v.Var("https://", "https_url"); err == nil {
		t.Fatal("expected https_url to fail for no host")
	}
}

func TestRuleDirPathWithDot(t *testing.T) {
	v := New()
	if err := v.Var(".", "dirpath"); err == nil {
		t.Fatal("expected dirpath to fail for '.'")
	}
}

func TestRuleDirPathWithDotInBase(t *testing.T) {
	v := New()
	if err := v.Var("/path/to/file.txt", "dirpath"); err == nil {
		t.Fatal("expected dirpath to fail for path with dot in base")
	}
}

func TestRuleDirNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "dir"); err == nil {
		t.Fatal("expected dir to fail for non-string")
	}
}

func TestRuleFileNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "file"); err == nil {
		t.Fatal("expected file to fail for non-string")
	}
}

func TestRuleFilePathNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "filepath"); err == nil {
		t.Fatal("expected filepath to fail for non-string")
	}
}

func TestRuleLuhnChecksumEmpty(t *testing.T) {
	v := New()
	if err := v.Var(true, "luhn_checksum"); err == nil {
		t.Fatal("expected luhn_checksum to fail for bool")
	}
}

func TestRuleDNSRFC1035LabelTooLong(t *testing.T) {
	longLabel := ""
	for i := 0; i < 64; i++ {
		longLabel += "a"
	}
	v := New()
	if err := v.Var(longLabel, "dns_rfc1035_label"); err == nil {
		t.Fatal("expected dns_rfc1035_label to fail for too long label")
	}
}

func TestRuleDNSRFC1035LabelNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "dns_rfc1035_label"); err == nil {
		t.Fatal("expected dns_rfc1035_label to fail for non-string")
	}
}

func TestRuleIPNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "ip"); err == nil {
		t.Fatal("expected ip to fail for non-string")
	}
}

func TestRuleIPv4NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "ipv4"); err == nil {
		t.Fatal("expected ipv4 to fail for non-string")
	}
}

func TestRuleIPv6NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "ipv6"); err == nil {
		t.Fatal("expected ipv6 to fail for non-string")
	}
}

func TestRuleCIDRNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "cidr"); err == nil {
		t.Fatal("expected cidr to fail for non-string")
	}
}

func TestRuleMACNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "mac"); err == nil {
		t.Fatal("expected mac to fail for non-string")
	}
}

func TestRuleOneOfNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "oneof=a b c"); err == nil {
		t.Fatal("expected oneof to fail for non-string")
	}
}

func TestRuleOneOfCINonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "oneofci=a b c"); err == nil {
		t.Fatal("expected oneofci to fail for non-string")
	}
}

func TestRuleE164NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "e164"); err == nil {
		t.Fatal("expected e164 to fail for non-string")
	}
}

func TestRuleURLEncodedNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "url_encoded"); err == nil {
		t.Fatal("expected url_encoded to fail for non-string")
	}
}

func TestRuleHTMLNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "html"); err == nil {
		t.Fatal("expected html to fail for non-string")
	}
}

func TestRuleHTMLEncodedNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "html_encoded"); err == nil {
		t.Fatal("expected html_encoded to fail for non-string")
	}
}

func TestRuleUUIDNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "uuid"); err == nil {
		t.Fatal("expected uuid to fail for non-string")
	}
}

func TestRuleUUID3NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "uuid3"); err == nil {
		t.Fatal("expected uuid3 to fail for non-string")
	}
}

func TestRuleBase32NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "base32"); err == nil {
		t.Fatal("expected base32 to fail for non-string")
	}
}

func TestRuleBase64NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "base64"); err == nil {
		t.Fatal("expected base64 to fail for non-string")
	}
}

func TestRuleBase64URLNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "base64url"); err == nil {
		t.Fatal("expected base64url to fail for non-string")
	}
}

func TestRuleBase64RawURLNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "base64rawurl"); err == nil {
		t.Fatal("expected base64rawurl to fail for non-string")
	}
}

func TestRuleMultibyteNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "multibyte"); err == nil {
		t.Fatal("expected multibyte to fail for non-string")
	}
}

func TestRuleHexColorNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "hexcolor"); err == nil {
		t.Fatal("expected hexcolor to fail for non-string")
	}
}

func TestRuleRGBNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "rgb"); err == nil {
		t.Fatal("expected rgb to fail for non-string")
	}
}

func TestRuleRGBANonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "rgba"); err == nil {
		t.Fatal("expected rgba to fail for non-string")
	}
}

func TestRuleHSLNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "hsl"); err == nil {
		t.Fatal("expected hsl to fail for non-string")
	}
}

func TestRuleHSLANonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "hsla"); err == nil {
		t.Fatal("expected hsla to fail for non-string")
	}
}

func TestRuleEmailNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "email"); err == nil {
		t.Fatal("expected email to fail for non-string")
	}
}

func TestRuleMongoDBNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "mongodb"); err == nil {
		t.Fatal("expected mongodb to fail for non-string")
	}
}

func TestRuleContainsNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "contains=x"); err == nil {
		t.Fatal("expected contains to fail for non-string")
	}
}

func TestRuleContainsAnyNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "containsany=x"); err == nil {
		t.Fatal("expected containsany to fail for non-string")
	}
}

func TestRuleContainsRuneNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "containsrune=x"); err == nil {
		t.Fatal("expected containsrune to fail for non-string")
	}
}

func TestRuleExcludesNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "excludes=x"); err == nil {
		t.Fatal("expected excludes to fail for non-string")
	}
}

func TestRuleExcludesAllNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "excludesall=x"); err == nil {
		t.Fatal("expected excludesall to fail for non-string")
	}
}

func TestRuleExcludesRuneNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "excludesrune=x"); err == nil {
		t.Fatal("expected excludesrune to fail for non-string")
	}
}

func TestRuleLowercaseNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "lowercase"); err == nil {
		t.Fatal("expected lowercase to fail for non-string")
	}
}

func TestRuleUppercaseNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "uppercase"); err == nil {
		t.Fatal("expected uppercase to fail for non-string")
	}
}

func TestRuleBooleanNonStringNonBool(t *testing.T) {
	v := New()
	if err := v.Var([]int{1}, "boolean"); err == nil {
		t.Fatal("expected boolean to fail for slice")
	}
}

func TestRuleLatitudeNonNumeric(t *testing.T) {
	v := New()
	if err := v.Var("abc", "latitude"); err == nil {
		t.Fatal("expected latitude to fail for non-numeric string")
	}
}

func TestRuleLongitudeNonNumeric(t *testing.T) {
	v := New()
	if err := v.Var("abc", "longitude"); err == nil {
		t.Fatal("expected longitude to fail for non-numeric string")
	}
}

func TestRuleStartsWithNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "startswith=x"); err == nil {
		t.Fatal("expected startswith to fail for non-string")
	}
}

func TestRuleEndsWithNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "endswith=x"); err == nil {
		t.Fatal("expected endswith to fail for non-string")
	}
}

func TestRuleStartsNotWithNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "startsnotwith=x"); err == nil {
		t.Fatal("expected startsnotwith to fail for non-string")
	}
}

func TestRuleEndsNotWithNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "endsnotwith=x"); err == nil {
		t.Fatal("expected endsnotwith to fail for non-string")
	}
}

func TestRuleDatetimeNonStringField(t *testing.T) {
	v := New()
	if err := v.Var(123, "datetime=2006-01-02"); err == nil {
		t.Fatal("expected datetime to fail for non-string")
	}
}

func TestRuleHostnameTrailingDot(t *testing.T) {
	v := New()
	if err := v.Var("api.example.com.", "hostname"); err != nil {
		t.Fatal("expected hostname to pass with trailing dot")
	}
}

func TestRuleCIDRv4NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "cidrv4"); err == nil {
		t.Fatal("expected cidrv4 to fail for non-string")
	}
}

func TestRuleCIDRv6NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "cidrv6"); err == nil {
		t.Fatal("expected cidrv6 to fail for non-string")
	}
}

func TestRuleHostnamePortEmptyHost(t *testing.T) {
	v := New()
	if err := v.Var(":8080", "hostname_port"); err == nil {
		t.Fatal("expected hostname_port to fail for empty host")
	}
}

func TestRuleHTTPURLNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "http_url"); err == nil {
		t.Fatal("expected http_url to fail for non-string")
	}
}

func TestRuleHTTPSURLNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "https_url"); err == nil {
		t.Fatal("expected https_url to fail for non-string")
	}
}

func TestRuleURINonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "uri"); err == nil {
		t.Fatal("expected uri to fail for non-string")
	}
}

func TestRuleMinInvalidParam(t *testing.T) {
	v := New()
	if err := v.Var("abc", "min=abc"); err == nil {
		t.Fatal("expected min to fail for non-numeric param")
	}
}

func TestRuleMaxInvalidParam(t *testing.T) {
	v := New()
	if err := v.Var("abc", "max=abc"); err == nil {
		t.Fatal("expected max to fail for non-numeric param")
	}
}

func TestRuleLenInvalidParam(t *testing.T) {
	v := New()
	if err := v.Var("abc", "len=abc"); err == nil {
		t.Fatal("expected len to fail for non-numeric param")
	}
}

func TestRuleGtInvalidParam(t *testing.T) {
	v := New()
	if err := v.Var(5, "gt=abc"); err == nil {
		t.Fatal("expected gt to fail for non-numeric param")
	}
}

func TestRuleGteInvalidParam(t *testing.T) {
	v := New()
	if err := v.Var(5, "gte=abc"); err == nil {
		t.Fatal("expected gte to fail for non-numeric param")
	}
}

func TestRuleLtInvalidParam(t *testing.T) {
	v := New()
	if err := v.Var(5, "lt=abc"); err == nil {
		t.Fatal("expected lt to fail for non-numeric param")
	}
}

func TestRuleLteInvalidParam(t *testing.T) {
	v := New()
	if err := v.Var(5, "lte=abc"); err == nil {
		t.Fatal("expected lte to fail for non-numeric param")
	}
}

func TestRuleMinInvalidValue(t *testing.T) {
	v := New()
	if err := v.Var(nil, "min=1"); err == nil {
		t.Fatal("expected min to fail for nil value")
	}
}

func TestRuleEqIgnoreCaseNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "eq_ignore_case=hello"); err == nil {
		t.Fatal("expected eq_ignore_case to fail for non-string")
	}
}

func TestRuleNeIgnoreCaseNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "ne_ignore_case=hello"); err != nil {
		t.Fatalf("expected ne_ignore_case to pass for int via scalarString: %v", err)
	}
}

func TestRuleNeNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "ne=hello"); err != nil {
		t.Fatalf("expected ne to pass for int via scalarString: %v", err)
	}
}

func TestRuleEqNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "eq=hello"); err == nil {
		t.Fatal("expected eq to fail for non-string")
	}
}

func TestRuleLuhnChecksumWithSpaces(t *testing.T) {
	v := New()
	if err := v.Var("4111 1111 1111 1111", "luhn_checksum"); err != nil {
		t.Fatal("expected luhn_checksum to pass with spaces")
	}
}

func TestRuleLuhnChecksumWithDashes(t *testing.T) {
	v := New()
	if err := v.Var("4111-1111-1111-1111", "luhn_checksum"); err != nil {
		t.Fatal("expected luhn_checksum to pass with dashes")
	}
}

func TestRuleLuhnChecksumEmptyDigits(t *testing.T) {
	v := New()
	if err := v.Var("  -  ", "luhn_checksum"); err == nil {
		t.Fatal("expected luhn_checksum to fail for no digits")
	}
}

func TestRuleNumberNilValue(t *testing.T) {
	v := New()
	if err := v.Var(nil, "number"); err == nil {
		t.Fatal("expected number to fail for nil")
	}
}

func TestRuleJSONInvalidType(t *testing.T) {
	v := New()
	if err := v.Var(123, "json"); err == nil {
		t.Fatal("expected json to fail for int")
	}
}

func TestRuleUniqueSliceInvalidElement(t *testing.T) {
	v := New()
	if err := v.Var([]interface{}{nil, nil}, "unique"); err == nil {
		t.Fatal("expected unique to fail for nil elements")
	}
}

func TestRulePortOutOfRange(t *testing.T) {
	v := New()
	if err := v.Var("70000", "port"); err == nil {
		t.Fatal("expected port to fail for out of range")
	}
	if err := v.Var("-1", "port"); err == nil {
		t.Fatal("expected port to fail for negative")
	}
}

func TestRuleOneOfNoMatch(t *testing.T) {
	v := New()
	if err := v.Var("d", "oneof=a b c"); err == nil {
		t.Fatal("expected oneof to fail for no match")
	}
}

func TestRuleOneOfCINoMatch(t *testing.T) {
	v := New()
	if err := v.Var("d", "oneofci=a b c"); err == nil {
		t.Fatal("expected oneofci to fail for no match")
	}
}

func TestCompareLengthOrNumberDefault(t *testing.T) {
	v := New()
	if err := v.Var(complex(1, 2), "min=1"); err == nil {
		t.Fatal("expected min to fail for unsupported type")
	}
}

func TestNumericValueString(t *testing.T) {
	v := New()
	if err := v.Var("3.14", "gt=3"); err != nil {
		t.Fatalf("expected numeric string to pass gt: %v", err)
	}
	if err := v.Var("abc", "gt=3"); err == nil {
		t.Fatal("expected non-numeric string to fail gt")
	}
}

func TestScalarStringDefault(t *testing.T) {
	v := New()
	if err := v.Var([]int{1, 2}, "oneof=1 2"); err == nil {
		t.Fatal("expected oneof to fail for slice")
	}
}

func TestBytesValueInvalid(t *testing.T) {
	v := New()
	if err := v.Var("hello", "json"); err == nil {
		t.Fatal("expected json to fail for non-bytes string")
	}
}

func TestRuleLuhnChecksumNonDigit(t *testing.T) {
	v := New()
	if err := v.Var("4111a11111111111", "luhn_checksum"); err == nil {
		t.Fatal("expected luhn_checksum to fail for non-digit char")
	}
}

func TestRuleUniqueMapInvalid(t *testing.T) {
	v := New()
	if err := v.Var(complex(1, 2), "unique"); err == nil {
		t.Fatal("expected unique to fail for complex type")
	}
}

func TestNumericValueDefault(t *testing.T) {
	v := New()
	if err := v.Var(complex(1, 2), "gt=0"); err == nil {
		t.Fatal("expected gt to fail for complex type")
	}
}

func TestBytesValueNonByteSlice(t *testing.T) {
	v := New()
	if err := v.Var([]int{1, 2, 3}, "json"); err == nil {
		t.Fatal("expected json to fail for non-byte slice")
	}
}

func TestScalarStringInvalidField(t *testing.T) {
	v := New()
	if err := v.Var([]int{1}, "port"); err == nil {
		t.Fatal("expected port to fail for slice")
	}
}

func TestRulePortNonNumeric(t *testing.T) {
	v := New()
	if err := v.Var("abc", "port"); err == nil {
		t.Fatal("expected port to fail for non-numeric string")
	}
}

func TestRuleUniqueMapDuplicate(t *testing.T) {
	v := New()
	if err := v.Var(map[string]int{"a": 1, "b": 1}, "unique"); err == nil {
		t.Fatal("expected unique to fail for map with duplicate values")
	}
}

func TestRuleLuhnChecksumDoubleGT9(t *testing.T) {
	v := New()
	if err := v.Var("79927398713", "luhn_checksum"); err != nil {
		t.Fatalf("expected luhn_checksum to pass for valid card: %v", err)
	}
}

func TestCompareLengthOrNumberInvalid(t *testing.T) {
	v := New()
	if err := v.Var(nil, "min=1"); err == nil {
		t.Fatal("expected min to fail for nil")
	}
}

func TestNumericValueUint(t *testing.T) {
	v := New()
	if err := v.Var(uint(5), "latitude"); err != nil {
		t.Fatalf("expected latitude to pass for uint in range: %v", err)
	}
}

func TestNumericValueInt(t *testing.T) {
	v := New()
	if err := v.Var(5, "latitude"); err != nil {
		t.Fatalf("expected latitude to pass for int in range: %v", err)
	}
}

func TestNumericValueFloat(t *testing.T) {
	v := New()
	if err := v.Var(45.0, "latitude"); err != nil {
		t.Fatalf("expected latitude to pass for float: %v", err)
	}
}

func TestNumericValueStringValid(t *testing.T) {
	v := New()
	if err := v.Var("45.0", "latitude"); err != nil {
		t.Fatalf("expected latitude to pass for numeric string: %v", err)
	}
}

func TestNumericValueInvalidDefault(t *testing.T) {
	v := New()
	if err := v.Var(complex(1, 2), "latitude"); err == nil {
		t.Fatal("expected latitude to fail for complex type")
	}
}

func TestBytesValueInvalidField(t *testing.T) {
	v := New()
	if err := v.Var(123, "json"); err == nil {
		t.Fatal("expected json to fail for int")
	}
}

func TestScalarStringDefaultBranch(t *testing.T) {
	v := New()
	if err := v.Var(complex(1, 2), "oneof=a b c"); err == nil {
		t.Fatal("expected oneof to fail for complex type")
	}
}

func TestRuleOneOfCIMatch(t *testing.T) {
	v := New()
	if err := v.Var("A", "oneofci=a b c"); err != nil {
		t.Fatalf("expected oneofci to pass for case-insensitive match: %v", err)
	}
}

func TestRuleUniqueMapPass(t *testing.T) {
	v := New()
	m := map[string]string{"a": "1", "b": "2"}
	if err := v.Var(m, "unique"); err != nil {
		t.Fatalf("expected unique map with distinct values to pass: %v", err)
	}
}

func TestRuleLuhnChecksumNoDigits(t *testing.T) {
	v := New()
	if err := v.Var("  -  ", "luhn_checksum"); err == nil {
		t.Fatal("expected luhn_checksum to fail for no digits")
	}
}

func TestCompareLengthOrNumberUint(t *testing.T) {
	v := New()
	if err := v.Var(uint(5), "min=3"); err != nil {
		t.Fatalf("expected min to pass for uint: %v", err)
	}
}

func TestCompareLengthOrNumberFloat(t *testing.T) {
	v := New()
	if err := v.Var(5.0, "min=3"); err != nil {
		t.Fatalf("expected min to pass for float: %v", err)
	}
}

func TestRuleOneOfCINonScalar(t *testing.T) {
	v := New()
	if err := v.Var([]int{1}, "oneofci=a b c"); err == nil {
		t.Fatal("expected oneofci to fail for slice")
	}
}

func TestRuleLuhnChecksumNonScalar(t *testing.T) {
	v := New()
	if err := v.Var([]int{1}, "luhn_checksum"); err == nil {
		t.Fatal("expected luhn_checksum to fail for slice")
	}
}

func TestNumericValueInvalidField(t *testing.T) {
	v := New()
	if err := v.Var(nil, "latitude"); err == nil {
		t.Fatal("expected latitude to fail for nil")
	}
}

func TestNumericValueDefaultType(t *testing.T) {
	v := New()
	if err := v.Var(complex(1, 2), "longitude"); err == nil {
		t.Fatal("expected longitude to fail for complex type")
	}
}

func TestBytesValueNilField(t *testing.T) {
	v := New()
	if err := v.Var(nil, "json"); err == nil {
		t.Fatal("expected json to fail for nil")
	}
}

func TestScalarStringNilField(t *testing.T) {
	v := New()
	if err := v.Var(nil, "port"); err == nil {
		t.Fatal("expected port to fail for nil")
	}
}

func TestRuleSemver(t *testing.T) {
	v := New()
	if err := v.Var("1.2.3", "semver"); err != nil {
		t.Fatal("expected semver to pass for 1.2.3")
	}
	if err := v.Var("v1.2.3", "semver"); err != nil {
		t.Fatal("expected semver to pass for v1.2.3")
	}
	if err := v.Var("1.0.0-alpha.1", "semver"); err != nil {
		t.Fatal("expected semver to pass for pre-release")
	}
	if err := v.Var("1.0.0+build.123", "semver"); err != nil {
		t.Fatal("expected semver to pass with build metadata")
	}
	if err := v.Var("1.2", "semver"); err == nil {
		t.Fatal("expected semver to fail for incomplete version")
	}
	if err := v.Var("abc", "semver"); err == nil {
		t.Fatal("expected semver to fail for non-version")
	}
	if err := v.Var("01.2.3", "semver"); err == nil {
		t.Fatal("expected semver to fail for leading zero")
	}
}

func TestRuleISBN10(t *testing.T) {
	v := New()
	if err := v.Var("0471958697", "isbn10"); err != nil {
		t.Fatal("expected isbn10 to pass for valid ISBN-10")
	}
	if err := v.Var("0-471-95869-7", "isbn10"); err != nil {
		t.Fatal("expected isbn10 to pass with dashes")
	}
	if err := v.Var("0 471 95869 7", "isbn10"); err != nil {
		t.Fatal("expected isbn10 to pass with spaces")
	}
	if err := v.Var("0471958698", "isbn10"); err == nil {
		t.Fatal("expected isbn10 to fail for wrong checksum")
	}
	if err := v.Var("123456789", "isbn10"); err == nil {
		t.Fatal("expected isbn10 to fail for too short")
	}
	if err := v.Var("12345678901", "isbn10"); err == nil {
		t.Fatal("expected isbn10 to fail for too long")
	}
}

func TestRuleISBN10WithX(t *testing.T) {
	v := New()
	if err := v.Var("080442957X", "isbn10"); err != nil {
		t.Fatal("expected isbn10 to pass for ISBN ending with X")
	}
}

func TestRuleISBN13(t *testing.T) {
	v := New()
	if err := v.Var("9780471117094", "isbn13"); err != nil {
		t.Fatal("expected isbn13 to pass for valid ISBN-13")
	}
	if err := v.Var("978-0-471-11709-4", "isbn13"); err != nil {
		t.Fatal("expected isbn13 to pass with dashes")
	}
	if err := v.Var("9780471117095", "isbn13"); err == nil {
		t.Fatal("expected isbn13 to fail for wrong checksum")
	}
	if err := v.Var("978047111709", "isbn13"); err == nil {
		t.Fatal("expected isbn13 to fail for too short")
	}
}

func TestRuleISSN(t *testing.T) {
	v := New()
	if err := v.Var("0317847X", "issn"); err != nil {
		t.Fatal("expected issn to pass for valid ISSN")
	}
	if err := v.Var("0317-847X", "issn"); err != nil {
		t.Fatal("expected issn to pass with dash")
	}
	if err := v.Var("03178470", "issn"); err == nil {
		t.Fatal("expected issn to fail for wrong checksum")
	}
	if err := v.Var("1234567", "issn"); err == nil {
		t.Fatal("expected issn to fail for too short")
	}
}

func TestRuleBIC(t *testing.T) {
	v := New()
	if err := v.Var("CHASUS33", "bic"); err != nil {
		t.Fatal("expected bic to pass for valid BIC")
	}
	if err := v.Var("CHASUS33XXX", "bic"); err != nil {
		t.Fatal("expected bic to pass for 11-char BIC")
	}
	if err := v.Var("CHASU", "bic"); err == nil {
		t.Fatal("expected bic to fail for too short")
	}
	if err := v.Var("12345678", "bic"); err == nil {
		t.Fatal("expected bic to fail for numeric")
	}
}

func TestRuleCron(t *testing.T) {
	v := New()
	if err := v.Var("*/5 * * * *", "cron"); err != nil {
		t.Fatal("expected cron to pass for standard 5-field")
	}
	if err := v.Var("0 0 1 1 *", "cron"); err != nil {
		t.Fatal("expected cron to pass for yearly")
	}
	if err := v.Var("0 0 * * * *", "cron"); err != nil {
		t.Fatal("expected cron to pass for 6-field")
	}
	if err := v.Var("* * * *", "cron"); err == nil {
		t.Fatal("expected cron to fail for 4-field")
	}
	if err := v.Var("* * * * * * *", "cron"); err == nil {
		t.Fatal("expected cron to fail for 7-field")
	}
}

func TestRuleDataURI(t *testing.T) {
	v := New()
	if err := v.Var("data:text/plain;base64,SGVsbG8=", "datauri"); err != nil {
		t.Fatal("expected datauri to pass for valid data URI")
	}
	if err := v.Var("data:text/html,Hello", "datauri"); err != nil {
		t.Fatal("expected datauri to pass for simple data URI")
	}
	if err := v.Var("http://example.com", "datauri"); err == nil {
		t.Fatal("expected datauri to fail for non-data URI")
	}
	if err := v.Var("data:", "datauri"); err == nil {
		t.Fatal("expected datauri to fail for empty data URI")
	}
}

func TestRuleBCP47(t *testing.T) {
	v := New()
	if err := v.Var("en", "bcp47"); err != nil {
		t.Fatal("expected bcp47 to pass for 'en'")
	}
	if err := v.Var("zh-CN", "bcp47"); err != nil {
		t.Fatal("expected bcp47 to pass for 'zh-CN'")
	}
	if err := v.Var("en-US", "bcp47"); err != nil {
		t.Fatal("expected bcp47 to pass for 'en-US'")
	}
	if err := v.Var("sr-Latn-RS", "bcp47"); err != nil {
		t.Fatal("expected bcp47 to pass for 'sr-Latn-RS'")
	}
	if err := v.Var("1", "bcp47"); err == nil {
		t.Fatal("expected bcp47 to fail for single digit")
	}
}

func TestRuleEthAddr(t *testing.T) {
	v := New()
	if err := v.Var("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38", "eth_addr"); err != nil {
		t.Fatal("expected eth_addr to pass for valid address")
	}
	if err := v.Var("0x0000000000000000000000000000000000000000", "eth_addr"); err != nil {
		t.Fatal("expected eth_addr to pass for zero address")
	}
	if err := v.Var("742d35Cc6634C0532925a3b844Bc9e7595f2bD38", "eth_addr"); err == nil {
		t.Fatal("expected eth_addr to fail without 0x prefix")
	}
	if err := v.Var("0x1234", "eth_addr"); err == nil {
		t.Fatal("expected eth_addr to fail for too short")
	}
}

func TestRuleBtcAddr(t *testing.T) {
	v := New()
	if err := v.Var("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "btc_addr"); err != nil {
		t.Fatal("expected btc_addr to pass for legacy address")
	}
	if err := v.Var("not-a-btc-address", "btc_addr"); err == nil {
		t.Fatal("expected btc_addr to fail for invalid address")
	}
}

func TestRuleSemverNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "semver"); err == nil {
		t.Fatal("expected semver to fail for non-string")
	}
}

func TestRuleISBN10NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "isbn10"); err == nil {
		t.Fatal("expected isbn10 to fail for non-string")
	}
}

func TestRuleISBN13NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "isbn13"); err == nil {
		t.Fatal("expected isbn13 to fail for non-string")
	}
}

func TestRuleISSNNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "issn"); err == nil {
		t.Fatal("expected issn to fail for non-string")
	}
}

func TestRuleBICNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "bic"); err == nil {
		t.Fatal("expected bic to fail for non-string")
	}
}

func TestRuleCronNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "cron"); err == nil {
		t.Fatal("expected cron to fail for non-string")
	}
}

func TestRuleDataURINonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "datauri"); err == nil {
		t.Fatal("expected datauri to fail for non-string")
	}
}

func TestRuleBCP47NonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "bcp47"); err == nil {
		t.Fatal("expected bcp47 to fail for non-string")
	}
}

func TestRuleEthAddrNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "eth_addr"); err == nil {
		t.Fatal("expected eth_addr to fail for non-string")
	}
}

func TestRuleBtcAddrNonString(t *testing.T) {
	v := New()
	if err := v.Var(123, "btc_addr"); err == nil {
		t.Fatal("expected btc_addr to fail for non-string")
	}
}

func TestRuleURIMailto(t *testing.T) {
	v := New()
	if err := v.Var("mailto:user@example.com", "uri"); err != nil {
		t.Fatal("expected uri to pass for mailto scheme")
	}
}

func TestRuleURITel(t *testing.T) {
	v := New()
	if err := v.Var("tel:+1-234-567-8900", "uri"); err != nil {
		t.Fatal("expected uri to pass for tel scheme")
	}
}

func TestRuleURINoScheme(t *testing.T) {
	v := New()
	if err := v.Var("no-scheme-here", "uri"); err == nil {
		t.Fatal("expected uri to fail for no scheme")
	}
}

func TestRuleURLEmptyHost(t *testing.T) {
	v := New()
	if err := v.Var("https://", "url"); err == nil {
		t.Fatal("expected url to fail for empty host")
	}
}

func TestRuleLowercaseEmpty(t *testing.T) {
	v := New()
	if err := v.Var("", "lowercase"); err != nil {
		t.Fatal("expected lowercase to pass for empty string")
	}
}

func TestRuleUppercaseEmpty(t *testing.T) {
	v := New()
	if err := v.Var("", "uppercase"); err != nil {
		t.Fatal("expected uppercase to pass for empty string")
	}
}

func TestRuleLowercaseMixedDigits(t *testing.T) {
	v := New()
	if err := v.Var("hello123", "lowercase"); err != nil {
		t.Fatal("expected lowercase to pass for lowercase+digits")
	}
	if err := v.Var("Hello123", "lowercase"); err == nil {
		t.Fatal("expected lowercase to fail for mixed case+digits")
	}
}

func TestRuleUppercaseMixedDigits(t *testing.T) {
	v := New()
	if err := v.Var("HELLO123", "uppercase"); err != nil {
		t.Fatal("expected uppercase to pass for uppercase+digits")
	}
	if err := v.Var("Hello123", "uppercase"); err == nil {
		t.Fatal("expected uppercase to fail for mixed case+digits")
	}
}

func TestRuleOneOfNonScalar(t *testing.T) {
	v := New()
	if err := v.Var([]int{1}, "oneof=a b c"); err == nil {
		t.Fatal("expected oneof to fail for non-scalar slice")
	}
}

func TestRuleNoneOfNonScalar(t *testing.T) {
	v := New()
	if err := v.Var([]int{1}, "noneof=a b c"); err != nil {
		t.Fatalf("expected noneof to pass for non-scalar slice (scalarString fails, ruleOneOfFast returns false, !false=true): %v", err)
	}
}

func TestRuleBtcAddrBech32(t *testing.T) {
	v := New()
	if err := v.Var("bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4", "btc_addr"); err != nil {
		t.Fatalf("expected btc_addr to pass for bech32 address: %v", err)
	}
	if err := v.Var("bc1qinvalid!@#", "btc_addr"); err == nil {
		t.Fatal("expected btc_addr to fail for invalid bech32")
	}
}

func TestRuleBtcAddrTooShort(t *testing.T) {
	v := New()
	if err := v.Var("short", "btc_addr"); err == nil {
		t.Fatal("expected btc_addr to fail for too short")
	}
}

func TestRuleBtcAddrTooLong(t *testing.T) {
	longAddr := "1"
	for i := 0; i < 70; i++ {
		longAddr += "A"
	}
	v := New()
	if err := v.Var(longAddr, "btc_addr"); err == nil {
		t.Fatal("expected btc_addr to fail for too long")
	}
}

func TestRuleURISchemeInvalidChar(t *testing.T) {
	v := New()
	if err := v.Var("ht tp://example.com", "uri"); err == nil {
		t.Fatal("expected uri to fail for space in scheme")
	}
}

func TestRuleURISchemeOnly(t *testing.T) {
	v := New()
	if err := v.Var("http:", "uri"); err == nil {
		t.Fatal("expected uri to fail for scheme only")
	}
}

func TestRuleHTTPSURLNonHTTPS(t *testing.T) {
	v := New()
	if err := v.Var("http://example.com", "https_url"); err == nil {
		t.Fatal("expected https_url to fail for http scheme")
	}
}

func TestRuleHTTPURLHostStartsWithDot(t *testing.T) {
	v := New()
	if err := v.Var("http://.example.com", "http_url"); err == nil {
		t.Fatal("expected http_url to fail for host starting with dot")
	}
}

func TestRuleHTTPURLHostStartsWithDash(t *testing.T) {
	v := New()
	if err := v.Var("http://-example.com", "http_url"); err == nil {
		t.Fatal("expected http_url to fail for host starting with dash")
	}
}

func TestRuleSemverPreReleaseInvalid(t *testing.T) {
	v := New()
	if err := v.Var("1.0.0-", "semver"); err == nil {
		t.Fatal("expected semver to fail for empty pre-release")
	}
}

func TestRuleSemverBuildMetaInvalid(t *testing.T) {
	v := New()
	if err := v.Var("1.0.0+", "semver"); err == nil {
		t.Fatal("expected semver to fail for empty build metadata")
	}
}

func TestRuleSemverPreReleaseTrailingDot(t *testing.T) {
	v := New()
	if err := v.Var("1.0.0-alpha.", "semver"); err == nil {
		t.Fatal("expected semver to fail for trailing dot in pre-release")
	}
}

func TestRuleSemverBuildMetaTrailingDot(t *testing.T) {
	v := New()
	if err := v.Var("1.0.0+build.", "semver"); err == nil {
		t.Fatal("expected semver to fail for trailing dot in build")
	}
}

func TestRuleSemverPreReleaseInvalidChar(t *testing.T) {
	v := New()
	if err := v.Var("1.0.0-alpha!bad", "semver"); err == nil {
		t.Fatal("expected semver to fail for invalid char in pre-release")
	}
}

func TestRuleSemverBuildMetaInvalidChar(t *testing.T) {
	v := New()
	if err := v.Var("1.0.0+build!bad", "semver"); err == nil {
		t.Fatal("expected semver to fail for invalid char in build")
	}
}

func TestRuleSemverPreReleaseOnlyZero(t *testing.T) {
	v := New()
	if err := v.Var("1.0.0-0", "semver"); err != nil {
		t.Fatalf("expected semver to pass for numeric pre-release: %v", err)
	}
}

func TestRuleSemverPreReleaseAllDashes(t *testing.T) {
	v := New()
	if err := v.Var("1.0.0---", "semver"); err == nil {
		t.Fatal("expected semver to fail for all-dash pre-release")
	}
}

func TestRuleSemverBuildMetaTrailingDash(t *testing.T) {
	v := New()
	if err := v.Var("1.0.0+build-", "semver"); err == nil {
		t.Fatal("expected semver to fail for trailing dash in build")
	}
}

func TestRuleISBN10NonDigitMiddle(t *testing.T) {
	v := New()
	if err := v.Var("0471X58697", "isbn10"); err == nil {
		t.Fatal("expected isbn10 to fail for non-digit in middle")
	}
}

func TestRuleISBN10TooFewDigits(t *testing.T) {
	v := New()
	if err := v.Var("047195869", "isbn10"); err == nil {
		t.Fatal("expected isbn10 to fail for too few digits")
	}
}

func TestRuleISBN13NonDigit(t *testing.T) {
	v := New()
	if err := v.Var("978047111709A", "isbn13"); err == nil {
		t.Fatal("expected isbn13 to fail for non-digit")
	}
}

func TestRuleISBN13TooFewDigits(t *testing.T) {
	v := New()
	if err := v.Var("978047111709", "isbn13"); err == nil {
		t.Fatal("expected isbn13 to fail for too few digits")
	}
}

func TestRuleISSNNonDigitMiddle(t *testing.T) {
	v := New()
	if err := v.Var("031X8471", "issn"); err == nil {
		t.Fatal("expected issn to fail for non-digit in middle")
	}
}

func TestRuleISSNTooFewDigits(t *testing.T) {
	v := New()
	if err := v.Var("0317847", "issn"); err == nil {
		t.Fatal("expected issn to fail for too few digits")
	}
}

func TestRuleBICBankCodeDigit(t *testing.T) {
	v := New()
	if err := v.Var("CH1SUS33", "bic"); err == nil {
		t.Fatal("expected bic to fail for digit in bank code")
	}
}

func TestRuleBICCountryCodeDigit(t *testing.T) {
	v := New()
	if err := v.Var("C1ASUS33", "bic"); err == nil {
		t.Fatal("expected bic to fail for digit in country code")
	}
}

func TestRuleBICLocationCodeDigit(t *testing.T) {
	v := New()
	if err := v.Var("CHASU133", "bic"); err == nil {
		t.Fatal("expected bic to fail for digit in location code")
	}
}

func TestRuleBICBranchCodeInvalid(t *testing.T) {
	v := New()
	if err := v.Var("CHASUS3!1", "bic"); err == nil {
		t.Fatal("expected bic to fail for invalid branch code")
	}
}

func TestRuleCronDoubleDash(t *testing.T) {
	v := New()
	if err := v.Var("1--2 * * * *", "cron"); err == nil {
		t.Fatal("expected cron to fail for double dash")
	}
}

func TestRuleCronDoubleSlash(t *testing.T) {
	v := New()
	if err := v.Var("1//2 * * * *", "cron"); err == nil {
		t.Fatal("expected cron to fail for double slash")
	}
}

func TestRuleCronInvalidChar(t *testing.T) {
	v := New()
	if err := v.Var("a * * * *", "cron"); err == nil {
		t.Fatal("expected cron to fail for alpha char")
	}
}

func TestRuleCronInvalidStarPos(t *testing.T) {
	v := New()
	if err := v.Var("1*2 * * * *", "cron"); err == nil {
		t.Fatal("expected cron to fail for invalid star position")
	}
}

func TestRuleDataURINoComma(t *testing.T) {
	v := New()
	if err := v.Var("data:text/plain;base64", "datauri"); err == nil {
		t.Fatal("expected datauri to fail for no comma")
	}
}

func TestRuleDataURIInvalidMimeType(t *testing.T) {
	v := New()
	if err := v.Var("data:;base64,SGVsbG8=", "datauri"); err != nil {
		t.Fatalf("expected datauri to pass for empty mime type with base64: %v", err)
	}
}

func TestRuleBCP47WithVariant(t *testing.T) {
	v := New()
	if err := v.Var("en-US-variant", "bcp47"); err != nil {
		t.Fatalf("expected bcp47 to pass for variant subtag: %v", err)
	}
}

func TestRuleBCP47NumericRegion(t *testing.T) {
	v := New()
	if err := v.Var("en-123", "bcp47"); err != nil {
		t.Fatalf("expected bcp47 to pass for numeric region: %v", err)
	}
}

func TestRuleBCP47ExtLang(t *testing.T) {
	v := New()
	if err := v.Var("zh-yue", "bcp47"); err == nil {
		t.Fatal("expected bcp47 to fail for extlang (4-char subtag)")
	}
}

func TestRuleEthAddrInvalidHex(t *testing.T) {
	v := New()
	if err := v.Var("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD3G", "eth_addr"); err == nil {
		t.Fatal("expected eth_addr to fail for invalid hex char")
	}
}

func TestRuleBtcAddrInvalidLegacy(t *testing.T) {
	v := New()
	if err := v.Var("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfN!", "btc_addr"); err == nil {
		t.Fatal("expected btc_addr to fail for invalid char")
	}
}

func TestRuleBtcAddrInvalidBech32(t *testing.T) {
	v := New()
	if err := v.Var("bc1q!nvalidaddress", "btc_addr"); err == nil {
		t.Fatal("expected btc_addr to fail for invalid bech32")
	}
}

func TestRuleHTTPURLNoDoubleSlash(t *testing.T) {
	v := New()
	if err := v.Var("http:example.com", "http_url"); err == nil {
		t.Fatal("expected http_url to fail for no double slash")
	}
}

func TestRuleHTTPURLSingleSlash(t *testing.T) {
	v := New()
	if err := v.Var("http:/example.com", "http_url"); err == nil {
		t.Fatal("expected http_url to fail for single slash")
	}
}

func TestRuleCompareOpInvalidField(t *testing.T) {
	v := New()
	if err := v.Var(nil, "gt=1"); err == nil {
		t.Fatal("expected gt to fail for nil")
	}
}

func TestTrimSpaceIfNeeded(t *testing.T) {
	v := New()
	if err := v.Var("  hello  ", "min=3"); err != nil {
		t.Fatalf("expected min to pass for trimmed string: %v", err)
	}
}

func TestRuleURIMailtoNoHost(t *testing.T) {
	v := New()
	if err := v.Var("mailto:", "uri"); err == nil {
		t.Fatal("expected uri to fail for mailto with no content after colon")
	}
}

func TestBuiltinRuleISBN13NonDigit(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"978-0-471-A1709-4", false},
		{"978047111709A", false},
		{"978-0-471-11709-4X", false},
		{"9780123456786", true}, // 有效的 ISBN13
	}
	for _, tc := range testCases {
		result := ruleISBN13(reflect.ValueOf(tc.input), "", false)
		if result != tc.expected {
			t.Errorf("ISBN13(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestBuiltinRuleISBN13WithNonDigitInMiddle(t *testing.T) {
	input := "9780A12345678"
	v := reflect.ValueOf(input)
	s, ok := stringValue(v)
	t.Logf("stringValue: %q, ok: %v", s, ok)

	for i, c := range s {
		t.Logf("s[%d] = %q (byte=%d)", i, c, s[i])
		if s[i] < '0' || s[i] > '9' {
			t.Logf("Found non-digit at index %d: %q", i, s[i])
		}
	}

	result := ruleISBN13(v, "", false)
	t.Logf("ruleISBN13 result: %v", result)
	if result {
		t.Fatal("expected isbn13 to fail for non-digit char in middle")
	}
}

func TestBuiltinRuleISSNTooManyDigits(t *testing.T) {
	if ruleISSN(reflect.ValueOf("0-317-8471-2"), "", false) {
		t.Fatal("expected issn to fail for too many digits")
	}
}

func TestBuiltinRuleBICBranchInvalidChar(t *testing.T) {
	if ruleBIC(reflect.ValueOf("CHASUS3!"), "", false) {
		t.Fatal("expected bic to fail for invalid branch code char")
	}
}

func TestBuiltinRuleDataURIParamsNonASCII(t *testing.T) {
	if ruleDataURI(reflect.ValueOf("data:text/plain;chars\x00et=utf-8,hello"), "", false) {
		t.Fatal("expected datauri to fail for non-ASCII in params")
	}
}

func TestBuiltinRuleBCP47VariantNonAlphanum(t *testing.T) {
	if ruleBCP47(reflect.ValueOf("en-US-variant!"), "", false) {
		t.Fatal("expected bcp47 to fail for non-alphanum in variant")
	}
}

func TestBuiltinRuleBtcAddrLegacyTooShort(t *testing.T) {
	if ruleBtcAddr(reflect.ValueOf("1A"), "", false) {
		t.Fatal("expected btc_addr to fail for too short legacy address")
	}
}

func TestBuiltinRuleSemverPatchFail(t *testing.T) {
	if ruleSemver(reflect.ValueOf("1.2."), "", false) {
		t.Fatal("expected semver to fail for empty patch")
	}
}
