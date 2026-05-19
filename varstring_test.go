/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-17 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 00:00:00
 * @FilePath: \go-argus\varstring_test.go
 * @Description: VarString 零反射快速路径测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"context"
	"testing"

	"github.com/kamalyes/go-argus/validate"
)

func TestVarStringRequired(t *testing.T) {
	v := New()
	if err := v.VarString("", "required"); err == nil {
		t.Fatal("expected empty string to fail required")
	}
	if err := v.VarString("hello", "required"); err != nil {
		t.Fatal("expected non-empty string to pass required")
	}
}

func TestVarStringEmail(t *testing.T) {
	v := New()
	if err := v.VarString("user@example.com", "email"); err != nil {
		t.Fatal("expected valid email to pass")
	}
	if err := v.VarString("not-email", "email"); err == nil {
		t.Fatal("expected invalid email to fail")
	}
}

func TestVarStringURL(t *testing.T) {
	v := New()
	if err := v.VarString("https://example.com/path", "url"); err != nil {
		t.Fatal("expected valid url to pass")
	}
	if err := v.VarString("not-a-url", "url"); err == nil {
		t.Fatal("expected invalid url to fail")
	}
}

func TestVarStringOmitEmpty(t *testing.T) {
	v := New()
	if err := v.VarString("", "omitempty,email"); err != nil {
		t.Fatalf("expected empty string with omitempty to pass: %v", err)
	}
	if err := v.VarString("not-email", "omitempty,email"); err == nil {
		t.Fatal("expected invalid email to fail with omitempty")
	}
	if err := v.VarString("user@example.com", "omitempty,email"); err != nil {
		t.Fatal("expected valid email to pass with omitempty")
	}
}

func TestVarStringOmitZero(t *testing.T) {
	v := New()
	if err := v.VarString("", "omitzero,email"); err != nil {
		t.Fatalf("expected empty string with omitzero to pass: %v", err)
	}
	if err := v.VarString("not-email", "omitzero,email"); err == nil {
		t.Fatal("expected invalid email to fail with omitzero")
	}
}

func TestVarStringOmitNil(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "omitnil,required"); err != nil {
		t.Fatal("expected omitnil to be skipped and required to pass")
	}
}

func TestVarStringStructOnly(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "structonly,required"); err != nil {
		t.Fatal("expected structonly to be skipped and required to pass")
	}
}

func TestVarStringNoStructLevel(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "nostructlevel,required"); err != nil {
		t.Fatal("expected nostructlevel to be skipped and required to pass")
	}
}

func TestVarStringOneOf(t *testing.T) {
	v := New()
	if err := v.VarString("admin", "oneof=admin member guest"); err != nil {
		t.Fatal("expected oneof to pass")
	}
	if err := v.VarString("root", "oneof=admin member guest"); err == nil {
		t.Fatal("expected oneof to fail")
	}
}

func TestVarStringOneOfCI(t *testing.T) {
	v := New()
	if err := v.VarString("Admin", "oneofci=admin member guest"); err != nil {
		t.Fatal("expected oneofci to pass")
	}
	if err := v.VarString("Root", "oneofci=admin member guest"); err == nil {
		t.Fatal("expected oneofci to fail")
	}
}

func TestVarStringNoneOf(t *testing.T) {
	v := New()
	if err := v.VarString("root", "noneof=admin member guest"); err != nil {
		t.Fatal("expected noneof to pass")
	}
	if err := v.VarString("admin", "noneof=admin member guest"); err == nil {
		t.Fatal("expected noneof to fail")
	}
}

func TestVarStringNoneOfCI(t *testing.T) {
	v := New()
	if err := v.VarString("Root", "noneofci=admin member guest"); err != nil {
		t.Fatal("expected noneofci to pass")
	}
	if err := v.VarString("Admin", "noneofci=admin member guest"); err == nil {
		t.Fatal("expected noneofci to fail")
	}
}

func TestVarStringMultiRule(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "required,min=1,max=100"); err != nil {
		t.Fatalf("expected multi-rule to pass: %v", err)
	}
	if err := v.VarString("", "required,min=1,max=100"); err == nil {
		t.Fatal("expected empty to fail required in multi-rule")
	}
	if err := v.VarString("a", "min=3,max=100"); err == nil {
		t.Fatal("expected short string to fail min in multi-rule")
	}
}

func TestVarStringReflectPathFallback(t *testing.T) {
	v := New()
	err := v.VarString("hello", "required_with=nonexistent_field")
	if err != nil {
		t.Fatalf("expected reflect path fallback to pass: %v", err)
	}
}

func TestVarStringCtx(t *testing.T) {
	v := New()
	ctx := context.Background()
	if err := v.VarStringCtx(ctx, "hello", "required"); err != nil {
		t.Fatal("expected VarStringCtx to pass")
	}
	if err := v.VarStringCtx(ctx, "", "required"); err == nil {
		t.Fatal("expected VarStringCtx to fail for empty")
	}
}

func TestVarStringEmptyTag(t *testing.T) {
	v := New()
	if err := v.VarString("hello", ""); err != nil {
		t.Fatalf("expected empty tag to pass: %v", err)
	}
}

func TestStringOneOfFunc(t *testing.T) {
	parts := []string{"admin", "member", "guest"}
	if !validate.StringOneOf("admin", parts) {
		t.Fatal("expected StringOneOf to find 'admin'")
	}
	if validate.StringOneOf("root", parts) {
		t.Fatal("expected StringOneOf to not find 'root'")
	}
}

func TestStringOneOfCIFunc(t *testing.T) {
	parts := []string{"admin", "member", "guest"}
	if !validate.StringOneOfCI("Admin", parts) {
		t.Fatal("expected StringOneOfCI to find 'Admin' case-insensitively")
	}
	if validate.StringOneOfCI("Root", parts) {
		t.Fatal("expected StringOneOfCI to not find 'Root'")
	}
}

func TestVarStringIPv4(t *testing.T) {
	v := New()
	if err := v.VarString("192.168.1.1", "ipv4"); err != nil {
		t.Fatal("expected ipv4 to pass")
	}
	if err := v.VarString("::1", "ipv4"); err == nil {
		t.Fatal("expected ipv6 to fail ipv4")
	}
}

func TestVarStringUUID(t *testing.T) {
	v := New()
	if err := v.VarString("550e8400-e29b-41d4-a716-446655440000", "uuid"); err != nil {
		t.Fatal("expected uuid to pass")
	}
	if err := v.VarString("not-a-uuid", "uuid"); err == nil {
		t.Fatal("expected uuid to fail")
	}
}

func TestVarStringBase64(t *testing.T) {
	v := New()
	if err := v.VarString("YXJndXM=", "base64"); err != nil {
		t.Fatal("expected base64 to pass")
	}
	if err := v.VarString("!!!invalid!!!", "base64"); err == nil {
		t.Fatal("expected base64 to fail")
	}
}

func TestVarStringJSON(t *testing.T) {
	v := New()
	if err := v.VarString(`{"key":"value"}`, "json"); err != nil {
		t.Fatal("expected json to pass")
	}
	if err := v.VarString("{invalid}", "json"); err == nil {
		t.Fatal("expected json to fail")
	}
}

func TestVarStringSemver(t *testing.T) {
	v := New()
	if err := v.VarString("1.2.3", "semver"); err != nil {
		t.Fatal("expected semver to pass")
	}
	if err := v.VarString("abc", "semver"); err == nil {
		t.Fatal("expected semver to fail")
	}
}

func TestVarStringCron(t *testing.T) {
	v := New()
	if err := v.VarString("*/5 * * * *", "cron"); err != nil {
		t.Fatal("expected cron to pass")
	}
	if err := v.VarString("* * * *", "cron"); err == nil {
		t.Fatal("expected cron to fail for 4-field")
	}
}

func TestVarStringContains(t *testing.T) {
	v := New()
	if err := v.VarString("hello world", "contains=world"); err != nil {
		t.Fatal("expected contains to pass")
	}
	if err := v.VarString("hello world", "contains=xyz"); err == nil {
		t.Fatal("expected contains to fail")
	}
}

func TestVarStringStartsWith(t *testing.T) {
	v := New()
	if err := v.VarString("hello world", "startswith=hello"); err != nil {
		t.Fatal("expected startswith to pass")
	}
	if err := v.VarString("hello world", "startswith=world"); err == nil {
		t.Fatal("expected startswith to fail")
	}
}

func TestVarStringDatetime(t *testing.T) {
	v := New()
	if err := v.VarString("2023-12-06T00:00:00Z", "datetime"); err != nil {
		t.Fatal("expected datetime to pass")
	}
	if err := v.VarString("invalid", "datetime"); err == nil {
		t.Fatal("expected datetime to fail")
	}
}

func TestVarStringNumber(t *testing.T) {
	v := New()
	if err := v.VarString("3.14", "number"); err != nil {
		t.Fatal("expected number to pass")
	}
	if err := v.VarString("abc", "number"); err == nil {
		t.Fatal("expected number to fail")
	}
}

func TestVarStringBoolean(t *testing.T) {
	v := New()
	if err := v.VarString("true", "boolean"); err != nil {
		t.Fatal("expected boolean to pass for 'true'")
	}
	if err := v.VarString("maybe", "boolean"); err == nil {
		t.Fatal("expected boolean to fail for 'maybe'")
	}
}

func TestVarStringLowercase(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "lowercase"); err != nil {
		t.Fatal("expected lowercase to pass")
	}
	if err := v.VarString("Hello", "lowercase"); err == nil {
		t.Fatal("expected lowercase to fail")
	}
}

func TestVarStringUppercase(t *testing.T) {
	v := New()
	if err := v.VarString("HELLO", "uppercase"); err != nil {
		t.Fatal("expected uppercase to pass")
	}
	if err := v.VarString("Hello", "uppercase"); err == nil {
		t.Fatal("expected uppercase to fail")
	}
}

func TestVarStringUnique(t *testing.T) {
	v := New()
	if err := v.VarString("abcdef", "unique"); err != nil {
		t.Fatal("expected unique to pass")
	}
	if err := v.VarString("aabc", "unique"); err == nil {
		t.Fatal("expected unique to fail")
	}
}

func TestVarStringLuhnChecksum(t *testing.T) {
	v := New()
	if err := v.VarString("4111111111111111", "credit_card"); err != nil {
		t.Fatal("expected credit_card to pass")
	}
	if err := v.VarString("4111111111111112", "credit_card"); err == nil {
		t.Fatal("expected credit_card to fail for bad checksum")
	}
}

func TestVarStringEthAddr(t *testing.T) {
	v := New()
	if err := v.VarString("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38", "eth_addr"); err != nil {
		t.Fatal("expected eth_addr to pass")
	}
	if err := v.VarString("0x1234", "eth_addr"); err == nil {
		t.Fatal("expected eth_addr to fail for too short")
	}
}

func TestVarStringBtcAddr(t *testing.T) {
	v := New()
	if err := v.VarString("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "btc_addr"); err != nil {
		t.Fatal("expected btc_addr to pass")
	}
	if err := v.VarString("not-a-btc-address", "btc_addr"); err == nil {
		t.Fatal("expected btc_addr to fail")
	}
}

func TestVarStringBCP47(t *testing.T) {
	v := New()
	if err := v.VarString("en-US", "bcp47"); err != nil {
		t.Fatal("expected bcp47 to pass")
	}
	if err := v.VarString("1", "bcp47"); err == nil {
		t.Fatal("expected bcp47 to fail")
	}
}

func TestVarStringDataURI(t *testing.T) {
	v := New()
	if err := v.VarString("data:text/plain;base64,SGVsbG8=", "datauri"); err != nil {
		t.Fatal("expected datauri to pass")
	}
	if err := v.VarString("http://example.com", "datauri"); err == nil {
		t.Fatal("expected datauri to fail")
	}
}

func TestVarStringISBN10(t *testing.T) {
	v := New()
	if err := v.VarString("0471958697", "isbn10"); err != nil {
		t.Fatal("expected isbn10 to pass")
	}
	if err := v.VarString("0471958698", "isbn10"); err == nil {
		t.Fatal("expected isbn10 to fail for wrong checksum")
	}
}

func TestVarStringISBN13(t *testing.T) {
	v := New()
	if err := v.VarString("9780471117094", "isbn13"); err != nil {
		t.Fatal("expected isbn13 to pass")
	}
	if err := v.VarString("9780471117095", "isbn13"); err == nil {
		t.Fatal("expected isbn13 to fail for wrong checksum")
	}
}

func TestVarStringISSN(t *testing.T) {
	v := New()
	if err := v.VarString("0317847X", "issn"); err != nil {
		t.Fatal("expected issn to pass")
	}
	if err := v.VarString("03178470", "issn"); err == nil {
		t.Fatal("expected issn to fail for wrong checksum")
	}
}

func TestVarStringBIC(t *testing.T) {
	v := New()
	if err := v.VarString("CHASUS33", "bic"); err != nil {
		t.Fatal("expected bic to pass")
	}
	if err := v.VarString("CHASU", "bic"); err == nil {
		t.Fatal("expected bic to fail for too short")
	}
}

func TestVarStringDNSRFC1035Label(t *testing.T) {
	v := New()
	if err := v.VarString("my-label", "dns_rfc1035_label"); err != nil {
		t.Fatal("expected dns_rfc1035_label to pass")
	}
	if err := v.VarString("Invalid", "dns_rfc1035_label"); err == nil {
		t.Fatal("expected dns_rfc1035_label to fail for uppercase")
	}
}

func TestVarStringIP(t *testing.T) {
	v := New()
	if err := v.VarString("192.168.1.1", "ip"); err != nil {
		t.Fatal("expected ip to pass")
	}
	if err := v.VarString("not-an-ip", "ip"); err == nil {
		t.Fatal("expected ip to fail")
	}
}

func TestVarStringCIDR(t *testing.T) {
	v := New()
	if err := v.VarString("192.168.1.0/24", "cidr"); err != nil {
		t.Fatal("expected cidr to pass")
	}
	if err := v.VarString("invalid", "cidr"); err == nil {
		t.Fatal("expected cidr to fail")
	}
}

func TestVarStringMAC(t *testing.T) {
	v := New()
	if err := v.VarString("00:11:22:33:44:55", "mac"); err != nil {
		t.Fatal("expected mac to pass")
	}
	if err := v.VarString("invalid", "mac"); err == nil {
		t.Fatal("expected mac to fail")
	}
}

func TestVarStringHostname(t *testing.T) {
	v := New()
	if err := v.VarString("api.example.com", "hostname"); err != nil {
		t.Fatal("expected hostname to pass")
	}
	if err := v.VarString("-invalid.com", "hostname"); err == nil {
		t.Fatal("expected hostname to fail")
	}
}

func TestVarStringFQDN(t *testing.T) {
	v := New()
	if err := v.VarString("api.example.com.", "fqdn"); err != nil {
		t.Fatal("expected fqdn to pass")
	}
	if err := v.VarString("api.example.com", "fqdn"); err == nil {
		t.Fatal("expected fqdn to fail without trailing dot")
	}
}

func TestVarStringPort(t *testing.T) {
	v := New()
	if err := v.VarString("443", "port"); err != nil {
		t.Fatal("expected port to pass")
	}
	if err := v.VarString("99999", "port"); err == nil {
		t.Fatal("expected port to fail for out-of-range")
	}
}

func TestVarStringE164(t *testing.T) {
	v := New()
	if err := v.VarString("+8613800138000", "e164"); err != nil {
		t.Fatal("expected e164 to pass")
	}
	if err := v.VarString("12345", "e164"); err == nil {
		t.Fatal("expected e164 to fail")
	}
}

func TestVarStringHexColor(t *testing.T) {
	v := New()
	if err := v.VarString("#12ffaa", "hexcolor"); err != nil {
		t.Fatal("expected hexcolor to pass")
	}
}

func TestVarStringRGB(t *testing.T) {
	v := New()
	if err := v.VarString("rgb(12, 34, 255)", "rgb"); err != nil {
		t.Fatal("expected rgb to pass")
	}
}

func TestVarStringRGBA(t *testing.T) {
	v := New()
	if err := v.VarString("rgba(12, 34, 255, 0.5)", "rgba"); err != nil {
		t.Fatal("expected rgba to pass")
	}
}

func TestVarStringHSL(t *testing.T) {
	v := New()
	if err := v.VarString("hsl(120, 50%, 75%)", "hsl"); err != nil {
		t.Fatal("expected hsl to pass")
	}
}

func TestVarStringHSLA(t *testing.T) {
	v := New()
	if err := v.VarString("hsla(120, 50%, 75%, 0.5)", "hsla"); err != nil {
		t.Fatal("expected hsla to pass")
	}
}

func TestVarStringAlpha(t *testing.T) {
	v := New()
	if err := v.VarString("abc", "alpha"); err != nil {
		t.Fatal("expected alpha to pass")
	}
	if err := v.VarString("abc123", "alpha"); err == nil {
		t.Fatal("expected alpha to fail for alphanum")
	}
}

func TestVarStringAlphanum(t *testing.T) {
	v := New()
	if err := v.VarString("abc123", "alphanum"); err != nil {
		t.Fatal("expected alphanum to pass")
	}
	if err := v.VarString("abc-123", "alphanum"); err == nil {
		t.Fatal("expected alphanum to fail for hyphen")
	}
}

func TestVarStringMongoDB(t *testing.T) {
	v := New()
	if err := v.VarString("507f1f77bcf86cd799439011", "mongodb"); err != nil {
		t.Fatal("expected mongodb to pass")
	}
	if err := v.VarString("invalid", "mongodb"); err == nil {
		t.Fatal("expected mongodb to fail")
	}
}

func TestVarStringTimezone(t *testing.T) {
	v := New()
	if err := v.VarString("UTC", "timezone"); err != nil {
		t.Fatal("expected timezone to pass")
	}
	if err := v.VarString("Invalid/Zone", "timezone"); err == nil {
		t.Fatal("expected timezone to fail")
	}
}

func TestVarStringLatitude(t *testing.T) {
	v := New()
	if err := v.VarString("45.0", "latitude"); err != nil {
		t.Fatal("expected latitude to pass")
	}
	if err := v.VarString("91.0", "latitude"); err == nil {
		t.Fatal("expected latitude to fail for out-of-range")
	}
}

func TestVarStringLongitude(t *testing.T) {
	v := New()
	if err := v.VarString("90.0", "longitude"); err != nil {
		t.Fatal("expected longitude to pass")
	}
	if err := v.VarString("181.0", "longitude"); err == nil {
		t.Fatal("expected longitude to fail for out-of-range")
	}
}

func TestVarStringExcludes(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "excludes=world"); err != nil {
		t.Fatal("expected excludes to pass")
	}
	if err := v.VarString("hello world", "excludes=world"); err == nil {
		t.Fatal("expected excludes to fail")
	}
}

func TestVarStringHTMLEncoded(t *testing.T) {
	v := New()
	if err := v.VarString("&amp;", "html_encoded"); err != nil {
		t.Fatal("expected html_encoded to pass")
	}
	if err := v.VarString("plain", "html_encoded"); err == nil {
		t.Fatal("expected html_encoded to fail for plain text")
	}
}

func TestVarStringURLEncoded(t *testing.T) {
	v := New()
	if err := v.VarString("hello%20world", "url_encoded"); err != nil {
		t.Fatal("expected url_encoded to pass")
	}
	if err := v.VarString("nopercent", "url_encoded"); err == nil {
		t.Fatal("expected url_encoded to fail without percent")
	}
}

func TestVarStringHTML(t *testing.T) {
	v := New()
	if err := v.VarString("<b>bold</b>", "html"); err != nil {
		t.Fatal("expected html to pass")
	}
	if err := v.VarString("no html", "html"); err == nil {
		t.Fatal("expected html to fail")
	}
}

func TestVarStringURI(t *testing.T) {
	v := New()
	if err := v.VarString("https://example.com/path", "uri"); err != nil {
		t.Fatal("expected uri to pass")
	}
	if err := v.VarString("no-scheme-here", "uri"); err == nil {
		t.Fatal("expected uri to fail for no scheme")
	}
}

func TestVarStringHTTPURL(t *testing.T) {
	v := New()
	if err := v.VarString("https://example.com", "http_url"); err != nil {
		t.Fatal("expected http_url to pass")
	}
	if err := v.VarString("ftp://example.com", "http_url"); err == nil {
		t.Fatal("expected http_url to fail for ftp")
	}
}

func TestVarStringHTTPSURL(t *testing.T) {
	v := New()
	if err := v.VarString("https://example.com", "https_url"); err != nil {
		t.Fatal("expected https_url to pass")
	}
	if err := v.VarString("http://example.com", "https_url"); err == nil {
		t.Fatal("expected https_url to fail for http")
	}
}

func TestVarStringHostnamePort(t *testing.T) {
	v := New()
	if err := v.VarString("example.com:8080", "hostname_port"); err != nil {
		t.Fatal("expected hostname_port to pass")
	}
	if err := v.VarString("example.com:99999", "hostname_port"); err == nil {
		t.Fatal("expected hostname_port to fail for invalid port")
	}
}

func TestVarStringBase32(t *testing.T) {
	v := New()
	if err := v.VarString("JBSWY3DPEB3W64TMMQ======", "base32"); err != nil {
		t.Fatal("expected base32 to pass")
	}
	if err := v.VarString("!!!invalid!!!", "base32"); err == nil {
		t.Fatal("expected base32 to fail")
	}
}

func TestVarStringBase64URL(t *testing.T) {
	v := New()
	if err := v.VarString("YXJndXM=", "base64url"); err != nil {
		t.Fatal("expected base64url to pass")
	}
}

func TestVarStringBase64RawURL(t *testing.T) {
	v := New()
	if err := v.VarString("YXJndXM", "base64rawurl"); err != nil {
		t.Fatal("expected base64rawurl to pass")
	}
}

func TestVarStringUUID3(t *testing.T) {
	v := New()
	if err := v.VarString("550e8400-e29b-31d4-a716-446655440000", "uuid3"); err != nil {
		t.Fatal("expected uuid3 to pass")
	}
}

func TestVarStringUUID4(t *testing.T) {
	v := New()
	if err := v.VarString("550e8400-e29b-41d4-a716-446655440000", "uuid4"); err != nil {
		t.Fatal("expected uuid4 to pass")
	}
}

func TestVarStringUUID5(t *testing.T) {
	v := New()
	if err := v.VarString("550e8400-e29b-51d4-a716-446655440000", "uuid5"); err != nil {
		t.Fatal("expected uuid5 to pass")
	}
}

func TestVarStringIsDefault(t *testing.T) {
	v := New()
	if err := v.VarString("", "isdefault"); err != nil {
		t.Fatal("expected empty string to pass isdefault")
	}
	if err := v.VarString("hello", "isdefault"); err == nil {
		t.Fatal("expected non-empty string to fail isdefault")
	}
}

func TestVarStringMinMaxLen(t *testing.T) {
	v := New()
	if err := v.VarString("abc", "min=2,max=5"); err != nil {
		t.Fatal("expected abc to pass min=2,max=5")
	}
	if err := v.VarString("a", "min=2"); err == nil {
		t.Fatal("expected 'a' to fail min=2")
	}
	if err := v.VarString("abcdef", "max=5"); err == nil {
		t.Fatal("expected 'abcdef' to fail max=5")
	}
	if err := v.VarString("abc", "len=3"); err != nil {
		t.Fatal("expected 'abc' to pass len=3")
	}
}

func TestVarStringEqNe(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "eq=hello"); err != nil {
		t.Fatal("expected eq=hello to pass")
	}
	if err := v.VarString("hello", "eq=world"); err == nil {
		t.Fatal("expected eq=world to fail")
	}
	if err := v.VarString("hello", "ne=world"); err != nil {
		t.Fatal("expected ne=world to pass")
	}
	if err := v.VarString("hello", "ne=hello"); err == nil {
		t.Fatal("expected ne=hello to fail")
	}
}

func TestVarStringEqIgnoreCaseNeIgnoreCase(t *testing.T) {
	v := New()
	if err := v.VarString("Hello", "eq_ignore_case=hello"); err != nil {
		t.Fatal("expected eq_ignore_case to pass")
	}
	if err := v.VarString("Hello", "ne_ignore_case=world"); err != nil {
		t.Fatal("expected ne_ignore_case to pass")
	}
	if err := v.VarString("Hello", "ne_ignore_case=hello"); err == nil {
		t.Fatal("expected ne_ignore_case=hello to fail")
	}
}

func TestVarStringGtGteLtLte(t *testing.T) {
	v := New()
	if err := v.VarString("abc", "gt=2"); err != nil {
		t.Fatal("expected gt=2 to pass for len=3")
	}
	if err := v.VarString("abc", "gte=3"); err != nil {
		t.Fatal("expected gte=3 to pass for len=3")
	}
	if err := v.VarString("ab", "lt=3"); err != nil {
		t.Fatal("expected lt=3 to pass for len=2")
	}
	if err := v.VarString("abc", "lte=3"); err != nil {
		t.Fatal("expected lte=3 to pass for len=3")
	}
}

func TestVarStringASCII(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "ascii"); err != nil {
		t.Fatal("expected ascii to pass")
	}
	if err := v.VarString("你好", "ascii"); err == nil {
		t.Fatal("expected ascii to fail for unicode")
	}
}

func TestVarStringMultibyte(t *testing.T) {
	v := New()
	if err := v.VarString("你好", "multibyte"); err != nil {
		t.Fatal("expected multibyte to pass")
	}
	if err := v.VarString("hello", "multibyte"); err == nil {
		t.Fatal("expected multibyte to fail for ascii")
	}
}

func TestVarStringHexadecimal(t *testing.T) {
	v := New()
	if err := v.VarString("deadbeef", "hexadecimal"); err != nil {
		t.Fatal("expected hexadecimal to pass")
	}
	if err := v.VarString("xyz", "hexadecimal"); err == nil {
		t.Fatal("expected hexadecimal to fail")
	}
}

func TestVarStringContainsAny(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "containsany=h"); err != nil {
		t.Fatal("expected containsany to pass")
	}
	if err := v.VarString("hello", "containsany=xyz"); err == nil {
		t.Fatal("expected containsany to fail")
	}
}

func TestVarStringContainsRune(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "containsrune=h"); err != nil {
		t.Fatal("expected containsrune to pass")
	}
	if err := v.VarString("hello", "containsrune=z"); err == nil {
		t.Fatal("expected containsrune to fail")
	}
}

func TestVarStringExcludesAll(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "excludesall=xyz"); err != nil {
		t.Fatal("expected excludesall to pass")
	}
	if err := v.VarString("hello", "excludesall=h"); err == nil {
		t.Fatal("expected excludesall to fail")
	}
}

func TestVarStringExcludesRune(t *testing.T) {
	v := New()
	if err := v.VarString("hello", "excludesrune=z"); err != nil {
		t.Fatal("expected excludesrune to pass")
	}
	if err := v.VarString("hello", "excludesrune=h"); err == nil {
		t.Fatal("expected excludesrune to fail")
	}
}

func TestVarStringEndsWith(t *testing.T) {
	v := New()
	if err := v.VarString("hello world", "endswith=world"); err != nil {
		t.Fatal("expected endswith to pass")
	}
	if err := v.VarString("hello world", "endswith=hello"); err == nil {
		t.Fatal("expected endswith to fail")
	}
}

func TestVarStringStartsNotWith(t *testing.T) {
	v := New()
	if err := v.VarString("hello world", "startsnotwith=world"); err != nil {
		t.Fatal("expected startsnotwith to pass")
	}
	if err := v.VarString("hello world", "startsnotwith=hello"); err == nil {
		t.Fatal("expected startsnotwith to fail")
	}
}

func TestVarStringEndsNotWith(t *testing.T) {
	v := New()
	if err := v.VarString("hello world", "endsnotwith=hello"); err != nil {
		t.Fatal("expected endsnotwith to pass")
	}
	if err := v.VarString("hello world", "endsnotwith=world"); err == nil {
		t.Fatal("expected endsnotwith to fail")
	}
}

func TestVarStringAlphaUnicode(t *testing.T) {
	v := New()
	if err := v.VarString("你好世界", "alphaunicode"); err != nil {
		t.Fatal("expected alphaunicode to pass")
	}
}

func TestVarStringAlphanumUnicode(t *testing.T) {
	v := New()
	if err := v.VarString("你好123", "alphanumunicode"); err != nil {
		t.Fatal("expected alphanumunicode to pass")
	}
}

func TestVarStringPrintASCII(t *testing.T) {
	v := New()
	if err := v.VarString("hello world", "printascii"); err != nil {
		t.Fatal("expected printascii to pass")
	}
}

func TestVarStringAlphaSpace(t *testing.T) {
	v := New()
	if err := v.VarString("hello world", "alphaspace"); err != nil {
		t.Fatal("expected alphaspace to pass")
	}
}

func TestVarStringAlphanumSpace(t *testing.T) {
	v := New()
	if err := v.VarString("abc 123", "alphanumspace"); err != nil {
		t.Fatal("expected alphanumspace to pass")
	}
}

func TestVarStringCIDRv4(t *testing.T) {
	v := New()
	if err := v.VarString("192.168.1.0/24", "cidrv4"); err != nil {
		t.Fatal("expected cidrv4 to pass")
	}
	if err := v.VarString("::1/128", "cidrv4"); err == nil {
		t.Fatal("expected ipv6 cidr to fail cidrv4")
	}
}

func TestVarStringCIDRv6(t *testing.T) {
	v := New()
	if err := v.VarString("::1/128", "cidrv6"); err != nil {
		t.Fatal("expected cidrv6 to pass")
	}
	if err := v.VarString("192.168.1.0/24", "cidrv6"); err == nil {
		t.Fatal("expected ipv4 cidr to fail cidrv6")
	}
}

func TestVarStringIPv6(t *testing.T) {
	v := New()
	if err := v.VarString("::1", "ipv6"); err != nil {
		t.Fatal("expected ipv6 to pass")
	}
	if err := v.VarString("192.168.1.1", "ipv6"); err == nil {
		t.Fatal("expected ipv4 to fail ipv6")
	}
}

func TestVarStringHostnameRFC1123(t *testing.T) {
	v := New()
	if err := v.VarString("api.example.com", "hostname_rfc1123"); err != nil {
		t.Fatal("expected hostname_rfc1123 to pass")
	}
}

func TestVarStringUUIDRFC4122(t *testing.T) {
	v := New()
	if err := v.VarString("550e8400-e29b-41d4-a716-446655440000", "uuid_rfc4122"); err != nil {
		t.Fatal("expected uuid_rfc4122 to pass")
	}
}

func TestVarStringUUID3RFC4122(t *testing.T) {
	v := New()
	if err := v.VarString("550e8400-e29b-31d4-a716-446655440000", "uuid3_rfc4122"); err != nil {
		t.Fatal("expected uuid3_rfc4122 to pass")
	}
}

func TestVarStringUUID4RFC4122(t *testing.T) {
	v := New()
	if err := v.VarString("550e8400-e29b-41d4-a716-446655440000", "uuid4_rfc4122"); err != nil {
		t.Fatal("expected uuid4_rfc4122 to pass")
	}
}

func TestVarStringUUID5RFC4122(t *testing.T) {
	v := New()
	if err := v.VarString("550e8400-e29b-51d4-a716-446655440000", "uuid5_rfc4122"); err != nil {
		t.Fatal("expected uuid5_rfc4122 to pass")
	}
}

func TestVarStringNumeric(t *testing.T) {
	v := New()
	if err := v.VarString("3.14", "numeric"); err != nil {
		t.Fatal("expected numeric to pass")
	}
	if err := v.VarString("abc", "numeric"); err == nil {
		t.Fatal("expected numeric to fail")
	}
}

func TestVarStringLuhnChecksumTag(t *testing.T) {
	v := New()
	if err := v.VarString("79927398713", "luhn_checksum"); err != nil {
		t.Fatal("expected luhn_checksum to pass")
	}
	if err := v.VarString("abc", "luhn_checksum"); err == nil {
		t.Fatal("expected luhn_checksum to fail for non-digits")
	}
}

func TestVarStringDatetimeCustomLayout(t *testing.T) {
	v := New()
	if err := v.VarString("2023-12-06", "datetime=2006-01-02"); err != nil {
		t.Fatal("expected datetime with custom layout to pass")
	}
}

func TestVarStringIPAddr(t *testing.T) {
	v := New()
	if err := v.VarString("192.168.1.1", "ip_addr"); err != nil {
		t.Fatal("expected ip_addr to pass")
	}
	if err := v.VarString("not-an-ip", "ip_addr"); err == nil {
		t.Fatal("expected ip_addr to fail")
	}
}
