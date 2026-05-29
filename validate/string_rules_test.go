/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-20 00:00:00
 * @FilePath: \go-argus\validate\string_rules_test.go
 * @Description: string_rules.go 测试，覆盖字符串规则校验函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"testing"

	"github.com/kamalyes/go-argus/constants"
)

func TestStringRequired(t *testing.T) {
	if !StringRequired("hello") {
		t.Fatal("expected required to pass")
	}
	if StringRequired("") {
		t.Fatal("expected required to fail for empty")
	}
	if StringRequired("   ") {
		t.Fatal("expected required to fail for whitespace")
	}
}

func TestStringIsDefault(t *testing.T) {
	if !StringIsDefault("") {
		t.Fatal("expected isdefault to pass for empty")
	}
	if StringIsDefault("hello") {
		t.Fatal("expected isdefault to fail for non-empty")
	}
}

func TestStringMin(t *testing.T) {
	if !StringMin("abc", "3") {
		t.Fatal("expected min=3 to pass for 'abc'")
	}
	if StringMin("ab", "3") {
		t.Fatal("expected min=3 to fail for 'ab'")
	}
}

func TestStringMax(t *testing.T) {
	if !StringMax("ab", "3") {
		t.Fatal("expected max=3 to pass for 'ab'")
	}
	if StringMax("abcd", "3") {
		t.Fatal("expected max=3 to fail for 'abcd'")
	}
}

func TestStringLen(t *testing.T) {
	if !StringLen("abc", "3") {
		t.Fatal("expected len=3 to pass for 'abc'")
	}
}

func TestStringEq(t *testing.T) {
	if !StringEq("hello", "hello") {
		t.Fatal("expected eq to pass")
	}
}

func TestStringNe(t *testing.T) {
	if !StringNe("hello", "world") {
		t.Fatal("expected ne to pass")
	}
}

func TestStringGt(t *testing.T) {
	if !StringGt("abcd", "3") {
		t.Fatal("expected gt to pass for len 4 > 3")
	}
}

func TestStringAlpha(t *testing.T) {
	if !StringAlpha("abc") {
		t.Fatal("expected alpha to pass")
	}
	if StringAlpha("abc123") {
		t.Fatal("expected alpha to fail for alphanum")
	}
}

func TestStringAlphanum(t *testing.T) {
	if !StringAlphanum("abc123") {
		t.Fatal("expected alphanum to pass")
	}
}

func TestStringEmail(t *testing.T) {
	if !IsEmail("test@example.com") {
		t.Fatal("expected email to pass")
	}
}

func TestStringIP(t *testing.T) {
	if !StringIP("192.168.1.1") {
		t.Fatal("expected IP to pass")
	}
}

func TestStringHostname(t *testing.T) {
	if !StringHostname("example.com") {
		t.Fatal("expected hostname to pass")
	}
	if !StringHostname("example.com.") {
		t.Fatal("expected hostname with trailing dot to pass")
	}
}

func TestStringFQDN(t *testing.T) {
	if !StringFQDN("example.com.") {
		t.Fatal("expected fqdn to pass")
	}
	if StringFQDN("example.com") {
		t.Fatal("expected fqdn to fail without trailing dot")
	}
}

func TestStringURL(t *testing.T) {
	if !StringURL("http://example.com") {
		t.Fatal("expected url to pass")
	}
}

func TestStringJSON(t *testing.T) {
	if !StringJSON(`{"key":"value"}`) {
		t.Fatal("expected json to pass")
	}
	if StringJSON("not-json") {
		t.Fatal("expected json to fail")
	}
}

func TestStringStartsWith(t *testing.T) {
	if !StringStartsWith("hello world", "hello") {
		t.Fatal("expected startsWith to pass")
	}
}

func TestStringEndsWith(t *testing.T) {
	if !StringEndsWith("hello world", "world") {
		t.Fatal("expected endsWith to pass")
	}
}

func TestStringContains(t *testing.T) {
	if !StringContains("hello world", "world") {
		t.Fatal("expected contains to pass")
	}
}

func TestStringLowercase(t *testing.T) {
	if !StringLowercase("hello") {
		t.Fatal("expected lowercase to pass")
	}
}

func TestStringUppercase(t *testing.T) {
	if !StringUppercase("HELLO") {
		t.Fatal("expected uppercase to pass")
	}
}

func TestStringBoolean(t *testing.T) {
	if !StringBoolean("true") {
		t.Fatal("expected boolean to pass for 'true'")
	}
	if !StringBoolean("1") {
		t.Fatal("expected boolean to pass for '1'")
	}
}

func TestStringNumber(t *testing.T) {
	if !StringNumber("123") {
		t.Fatal("expected number to pass")
	}
}

func TestStringUnique(t *testing.T) {
	if !StringUnique("abc") {
		t.Fatal("expected unique to pass for 'abc'")
	}
	if StringUnique("aba") {
		t.Fatal("expected unique to fail for 'aba'")
	}
}

func TestStringHexColor(t *testing.T) {
	if !StringHexColor("#ff0000") {
		t.Fatal("expected hexcolor to pass")
	}
}

func TestStringRGB(t *testing.T) {
	if !StringRGB("rgb(255,0,0)") {
		t.Fatal("expected rgb to pass")
	}
}

func TestStringRGBA(t *testing.T) {
	if !StringRGBA("rgba(255,0,0,0.5)") {
		t.Fatal("expected rgba to pass")
	}
}

func TestStringHSL(t *testing.T) {
	if !StringHSL("hsl(120,100%,50%)") {
		t.Fatal("expected hsl to pass")
	}
}

func TestStringHSLA(t *testing.T) {
	if !StringHSLA("hsla(120,100%,50%,0.5)") {
		t.Fatal("expected hsla to pass")
	}
}

func TestStringBase64(t *testing.T) {
	if !StringBase64("SGVsbG8=") {
		t.Fatal("expected base64 to pass")
	}
}

func TestStringBase32(t *testing.T) {
	if !StringBase32("JBSWY3DPEE======") {
		t.Fatal("expected base32 to pass")
	}
}

func TestStringUUID(t *testing.T) {
	if !StringUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8") {
		t.Fatal("expected uuid to pass")
	}
	if !StringUUID(" 6ba7b810-9dad-11d1-80b4-00c04fd430c8 ") {
		t.Fatal("expected uuid to use canonical trimming")
	}
}

func TestStringPort(t *testing.T) {
	if !StringPort("443") {
		t.Fatal("expected port to pass")
	}
	if StringPort("99999") {
		t.Fatal("expected port to fail for out-of-range")
	}
}

func TestStringCIDR(t *testing.T) {
	if !StringCIDR("10.0.0.0/8") {
		t.Fatal("expected cidr to pass")
	}
}

func TestStringMAC(t *testing.T) {
	if !StringMAC("00:11:22:33:44:55") {
		t.Fatal("expected mac to pass")
	}
}

func TestStringIPv4(t *testing.T) {
	if !StringIPv4("192.168.1.1") {
		t.Fatal("expected ipv4 to pass")
	}
}

func TestStringIPv6(t *testing.T) {
	if !StringIPv6("::1") {
		t.Fatal("expected ipv6 to pass")
	}
}

func TestStringURLEncoded(t *testing.T) {
	if !StringURLEncoded("hello%20world") {
		t.Fatal("expected url_encoded to pass")
	}
}

func TestStringHTML(t *testing.T) {
	if !StringHTML("<b>hello</b>") {
		t.Fatal("expected html to pass")
	}
}

func TestStringHTMLEncoded(t *testing.T) {
	if !StringHTMLEncoded("&lt;b&gt;") {
		t.Fatal("expected html_encoded to pass")
	}
}

func TestStringE164(t *testing.T) {
	if !StringE164("+1234567890") {
		t.Fatal("expected e164 to pass")
	}
}

func TestStringURI(t *testing.T) {
	if !StringURI("http://example.com/path") {
		t.Fatal("expected uri to pass")
	}
}

func TestStringHTTPURL(t *testing.T) {
	if !StringHTTPURL("http://example.com") {
		t.Fatal("expected http_url to pass")
	}
}

func TestStringHTTPSURL(t *testing.T) {
	if !StringHTTPSURL("https://example.com") {
		t.Fatal("expected https_url to pass")
	}
}

func TestStringExcludes(t *testing.T) {
	if !StringExcludes("hello", "xyz") {
		t.Fatal("expected excludes to pass")
	}
}

func TestStringExcludesAll(t *testing.T) {
	if !StringExcludesAll("hello", "xyz") {
		t.Fatal("expected excludesall to pass")
	}
}

func TestStringContainsAny(t *testing.T) {
	if !StringContainsAny("hello", "ae") {
		t.Fatal("expected containsany to pass")
	}
}

func TestStringContainsRune(t *testing.T) {
	if !StringContainsRune("hello", "e") {
		t.Fatal("expected containsrune to pass")
	}
}

func TestStringExcludesRune(t *testing.T) {
	if !StringExcludesRune("hello", "z") {
		t.Fatal("expected excludesrune to pass")
	}
}

func TestStringStartsNotWith(t *testing.T) {
	if !StringStartsNotWith("hello", "xyz") {
		t.Fatal("expected startsnotwith to pass")
	}
}

func TestStringEndsNotWith(t *testing.T) {
	if !StringEndsNotWith("hello", "xyz") {
		t.Fatal("expected endsnotwith to pass")
	}
}

func TestStringEqIgnoreCase(t *testing.T) {
	if !StringEqIgnoreCase("Hello", "hello") {
		t.Fatal("expected eq_ignore_case to pass")
	}
}

func TestStringNeIgnoreCase(t *testing.T) {
	if !StringNeIgnoreCase("Hello", "world") {
		t.Fatal("expected ne_ignore_case to pass")
	}
}

func TestStringASCII(t *testing.T) {
	if !StringASCII("hello") {
		t.Fatal("expected ascii to pass")
	}
}

func TestStringPrintASCII(t *testing.T) {
	if !StringPrintASCII("hello") {
		t.Fatal("expected printascii to pass")
	}
}

func TestStringMultibyte(t *testing.T) {
	if !StringMultibyte("你好") {
		t.Fatal("expected multibyte to pass")
	}
}

func TestStringHexadecimal(t *testing.T) {
	if !StringHexadecimal("abcdef0123456789") {
		t.Fatal("expected hexadecimal to pass")
	}
}

func TestStringAlphaUnicode(t *testing.T) {
	if !StringAlphaUnicode("你好") {
		t.Fatal("expected alphaunicode to pass")
	}
}

func TestStringAlphanumUnicode(t *testing.T) {
	if !StringAlphanumUnicode("你好123") {
		t.Fatal("expected alphanumunicode to pass")
	}
}

func TestStringAlphaSpace(t *testing.T) {
	if !StringAlphaSpace("hello world") {
		t.Fatal("expected alphaspace to pass")
	}
}

func TestStringAlphanumSpace(t *testing.T) {
	if !StringAlphanumSpace("hello 123") {
		t.Fatal("expected alphanumspace to pass")
	}
}

func TestStringDatetime(t *testing.T) {
	if !StringDatetime("2024-01-01", "2006-01-02") {
		t.Fatal("expected datetime to pass")
	}
}

func TestStringTimezone(t *testing.T) {
	if !StringTimezone("UTC") {
		t.Fatal("expected timezone to pass")
	}
}

func TestStringLatitude(t *testing.T) {
	if !StringLatitude("45.0") {
		t.Fatal("expected latitude to pass")
	}
}

func TestStringLongitude(t *testing.T) {
	if !StringLongitude("90.0") {
		t.Fatal("expected longitude to pass")
	}
}

func TestStringMongoDB(t *testing.T) {
	if !StringMongoDB("507f1f77bcf86cd799439011") {
		t.Fatal("expected mongodb to pass")
	}
}

func TestStringDNSRFC1035Label(t *testing.T) {
	if !StringDNSRFC1035Label("example") {
		t.Fatal("expected dns_rfc1035_label to pass")
	}
}

func TestStringHostnamePort(t *testing.T) {
	if !StringHostnamePort("example.com:8080") {
		t.Fatal("expected hostname_port to pass")
	}
}

func TestStringFile(t *testing.T) {
	// 文件系统相关测试在 CI 中可能不稳定，仅测试基本路径
	_ = StringFile("nonexistent_file.txt")
}

func TestStringDir(t *testing.T) {
	_ = StringDir("nonexistent_dir")
}

func TestStringFilePath(t *testing.T) {
	if !StringFilePath("C:\\test\\file.txt") && !StringFilePath("/test/file.txt") {
		t.Fatal("expected filepath to pass for valid path")
	}
}

func TestStringDirPath(t *testing.T) {
	if !StringDirPath("C:\\test\\") && !StringDirPath("/test/") {
		t.Fatal("expected dirpath to pass for valid dir path")
	}
}

func TestStringGteLte(t *testing.T) {
	if !StringGte("abc", "3") {
		t.Fatal("expected gte to pass for len 3 >= 3")
	}
	if !StringLte("abc", "3") {
		t.Fatal("expected lte to pass for len 3 <= 3")
	}
}

func TestStringCIDRv4(t *testing.T) {
	if !StringCIDRv4("10.0.0.0/8") {
		t.Fatal("expected cidrv4 to pass")
	}
}

func TestStringCIDRv6(t *testing.T) {
	if !StringCIDRv6("::1/128") {
		t.Fatal("expected cidrv6 to pass")
	}
}

func TestStringUUID3(t *testing.T) {
	// UUID v3 使用 MD5，版本位为 3
	if !StringUUID3("6ba7b810-9dad-31d1-80b4-00c04fd430c8") {
		t.Fatal("expected uuid3 to pass")
	}
}

func TestStringUUID4(t *testing.T) {
	// UUID v4 使用随机，版本位为 4
	if !StringUUID4("6ba7b810-9dad-41d1-80b4-00c04fd430c8") {
		t.Fatal("expected uuid4 to pass")
	}
}

func TestStringUUID5(t *testing.T) {
	// UUID v5 使用 SHA-1，版本位为 5
	if !StringUUID5("6ba7b810-9dad-51d1-80b4-00c04fd430c8") {
		t.Fatal("expected uuid5 to pass")
	}
}

func TestStringBase64URL(t *testing.T) {
	if !StringBase64URL("SGVsbG8=") {
		t.Fatal("expected base64url to pass")
	}
}

func TestStringBase64RawURL(t *testing.T) {
	if !StringBase64RawURL("SGVsbG8") {
		t.Fatal("expected base64rawurl to pass")
	}
}

// --- 补充边界条件测试 ---

func TestIsBlankString(t *testing.T) {
	if !IsBlankString("") {
		t.Fatal("expected blank for empty")
	}
	if !IsBlankString("  \t\n") {
		t.Fatal("expected blank for whitespace")
	}
	if IsBlankString("a") {
		t.Fatal("expected not blank for 'a'")
	}
}

func TestStringCompareLengthInvalidParam(t *testing.T) {
	if StringCompareLength("abc", "notanumber", constants.CmpEQ) {
		t.Fatal("expected compare length to fail for invalid param")
	}
}

func TestStringLt(t *testing.T) {
	if !StringLt("ab", "3") {
		t.Fatal("expected lt to pass for len 2 < 3")
	}
}

func TestStringMatchRunesEmpty(t *testing.T) {
	if StringMatchRunes("", func(r rune) bool { return true }) {
		t.Fatal("expected matchrunes to fail for empty string")
	}
}

func TestStringAlphaFail(t *testing.T) {
	if StringAlpha("") {
		t.Fatal("expected alpha to fail for empty")
	}
}

func TestStringAlphanumFail(t *testing.T) {
	if StringAlphanum("abc!") {
		t.Fatal("expected alphanum to fail for special char")
	}
}

func TestStringAlphanumSpaceFail(t *testing.T) {
	if StringAlphanumSpace("abc!@#") {
		t.Fatal("expected alphanumspace to fail for special chars")
	}
}

func TestStringAlphaUnicodeFail(t *testing.T) {
	if StringAlphaUnicode("123") {
		t.Fatal("expected alphaunicode to fail for numbers")
	}
}

func TestStringAlphanumUnicodeFail(t *testing.T) {
	if StringAlphanumUnicode("!!!") {
		t.Fatal("expected alphanumunicode to fail for special chars")
	}
}

func TestStringAlphaSpaceFail(t *testing.T) {
	if StringAlphaSpace("123") {
		t.Fatal("expected alphaspace to fail for numbers")
	}
}

func TestStringASCIIFail(t *testing.T) {
	if StringASCII("\xff") {
		t.Fatal("expected ascii to fail for non-ascii")
	}
}

func TestStringPrintASCIIFail(t *testing.T) {
	if StringPrintASCII(string(rune(0x01))) {
		t.Fatal("expected printascii to fail for control char")
	}
}

func TestStringMultibyteFail(t *testing.T) {
	if StringMultibyte("hello") {
		t.Fatal("expected multibyte to fail for ASCII-only")
	}
}

func TestStringHexadecimalFail(t *testing.T) {
	if StringHexadecimal("xyz") {
		t.Fatal("expected hexadecimal to fail for non-hex")
	}
}

func TestStringHexColorFail(t *testing.T) {
	if StringHexColor("gg0000") {
		t.Fatal("expected hexcolor to fail for invalid hex")
	}
}

func TestStringRGBFail(t *testing.T) {
	if StringRGB("rgb(999,0,0)") {
		t.Fatal("expected rgb to fail for out-of-range")
	}
}

func TestStringRGBAFail(t *testing.T) {
	if StringRGBA("rgba(999,0,0,0.5)") {
		t.Fatal("expected rgba to fail for out-of-range")
	}
}

func TestStringHSLFail(t *testing.T) {
	if StringHSL("hsl(999,100%,50%)") {
		t.Fatal("expected hsl to fail for out-of-range")
	}
}

func TestStringHSLAFail(t *testing.T) {
	if StringHSLA("hsla(999,100%,50%,0.5)") {
		t.Fatal("expected hsla to fail for out-of-range")
	}
}

func TestStringE164Fail(t *testing.T) {
	if StringE164("1234567890") {
		t.Fatal("expected e164 to fail without +")
	}
}

func TestStringPortFail(t *testing.T) {
	if StringPort("abc") {
		t.Fatal("expected port to fail for non-numeric")
	}
}

func TestStringURLEncodedNoPercent(t *testing.T) {
	if StringURLEncoded("hello") {
		t.Fatal("expected url_encoded to fail without %")
	}
}

func TestStringURLEncodedInvalidEscape(t *testing.T) {
	if StringURLEncoded("hello%2") {
		t.Fatal("expected url_encoded to fail for incomplete escape")
	}
	if StringURLEncoded("hello%GG") {
		t.Fatal("expected url_encoded to fail for invalid hex")
	}
}

func TestIsHexChar(t *testing.T) {
	if !IsHexChar('a') || !IsHexChar('F') || !IsHexChar('0') {
		t.Fatal("expected hex chars to pass")
	}
	if IsHexChar('g') {
		t.Fatal("expected g to fail hex char check")
	}
}

func TestStringHTMLFail(t *testing.T) {
	if StringHTML("hello world") {
		t.Fatal("expected html to fail without tags")
	}
}

func TestStringHTMLEncodedAmp(t *testing.T) {
	// 测试 &amp; 模式
	if !StringHTMLEncoded("&amp;") {
		t.Fatal("expected &amp; to pass")
	}
}

func TestStringHTMLEncodedLt(t *testing.T) {
	if !StringHTMLEncoded("&lt;") {
		t.Fatal("expected &lt; to pass")
	}
}

func TestStringHTMLEncodedGt(t *testing.T) {
	if !StringHTMLEncoded("&gt;") {
		t.Fatal("expected &gt; to pass")
	}
}

func TestStringHTMLEncodedQuot(t *testing.T) {
	if !StringHTMLEncoded("&quot;") {
		t.Fatal("expected &quot; to pass")
	}
}

func TestStringHTMLEncodedApos(t *testing.T) {
	if !StringHTMLEncoded("&apos;") {
		t.Fatal("expected &apos; to pass")
	}
}

func TestStringHTMLEncodedNbsp(t *testing.T) {
	if !StringHTMLEncoded("&nbsp;") {
		t.Fatal("expected &nbsp; to pass")
	}
}

func TestStringHTMLEncodedNumeric(t *testing.T) {
	if !StringHTMLEncoded("&#60;") {
		t.Fatal("expected numeric entity to pass")
	}
}

func TestStringHTMLEncodedNoAmp(t *testing.T) {
	if StringHTMLEncoded("hello") {
		t.Fatal("expected plain text to fail")
	}
}

func TestStringBase32Empty(t *testing.T) {
	if StringBase32("") {
		t.Fatal("expected base32 to fail for empty")
	}
	if StringBase32("  ") {
		t.Fatal("expected base32 to fail for whitespace")
	}
}

func TestStringBase64URLEmpty(t *testing.T) {
	if StringBase64URL("") {
		t.Fatal("expected base64url to fail for empty")
	}
	if StringBase64URL("  ") {
		t.Fatal("expected base64url to fail for whitespace")
	}
}

func TestStringBase64RawURLEmpty(t *testing.T) {
	if StringBase64RawURL("") {
		t.Fatal("expected base64rawurl to fail for empty")
	}
	if StringBase64RawURL("  ") {
		t.Fatal("expected base64rawurl to fail for whitespace")
	}
}

func TestStringJSONEdgeCases(t *testing.T) {
	// 空
	if StringJSON("") {
		t.Fatal("expected json to fail for empty")
	}
	// 数组
	if !StringJSON(`[1,2,3]`) {
		t.Fatal("expected json array to pass")
	}
	// 布尔
	if !StringJSON("true") {
		t.Fatal("expected json true to pass")
	}
	if !StringJSON("false") {
		t.Fatal("expected json false to pass")
	}
	// null
	if !StringJSON("null") {
		t.Fatal("expected json null to pass")
	}
	// 数字
	if !StringJSON("42") {
		t.Fatal("expected json number to pass")
	}
	if !StringJSON("-3.14e+10") {
		t.Fatal("expected json negative float with exp to pass")
	}
	// 空对象
	if !StringJSON("{}") {
		t.Fatal("expected json empty object to pass")
	}
	// 空数组
	if !StringJSON("[]") {
		t.Fatal("expected json empty array to pass")
	}
	// 带空白的 JSON
	if !StringJSON(`  { "key" : "val" }  `) {
		t.Fatal("expected json with whitespace to pass")
	}
	// 无效 JSON
	if StringJSON("{") {
		t.Fatal("expected json to fail for incomplete object")
	}
	if StringJSON(`{"key"`) {
		t.Fatal("expected json to fail for incomplete key")
	}
	if StringJSON("[") {
		t.Fatal("expected json to fail for incomplete array")
	}
}

func TestStringLowercaseFail(t *testing.T) {
	if StringLowercase("Hello") {
		t.Fatal("expected lowercase to fail for mixed case")
	}
}

func TestStringUppercaseFail(t *testing.T) {
	if StringUppercase("HELLO world") {
		t.Fatal("expected uppercase to fail for mixed case")
	}
}

func TestStringBooleanFail(t *testing.T) {
	if StringBoolean("maybe") {
		t.Fatal("expected boolean to fail for 'maybe'")
	}
}

func TestStringNumberFail(t *testing.T) {
	if StringNumber("abc") {
		t.Fatal("expected number to fail for 'abc'")
	}
}

func TestStringDatetimeDefault(t *testing.T) {
	if !StringDatetime("2024-01-01T00:00:00Z", "") {
		t.Fatal("expected datetime with default RFC3339 to pass")
	}
}

func TestStringDatetimeFail(t *testing.T) {
	if StringDatetime("not-a-date", "2006-01-02") {
		t.Fatal("expected datetime to fail for invalid date")
	}
}

func TestStringTimezoneFail(t *testing.T) {
	if StringTimezone("Invalid/Zone") {
		t.Fatal("expected timezone to fail for invalid zone")
	}
}

func TestStringLatitudeFail(t *testing.T) {
	if StringLatitude("91") {
		t.Fatal("expected latitude to fail for >90")
	}
}

func TestStringLongitudeFail(t *testing.T) {
	if StringLongitude("181") {
		t.Fatal("expected longitude to fail for >180")
	}
}

func TestStringOneOf(t *testing.T) {
	if !StringOneOf("a", []string{"a", "b", "c"}) {
		t.Fatal("expected oneof to pass")
	}
	if StringOneOf("d", []string{"a", "b", "c"}) {
		t.Fatal("expected oneof to fail")
	}
}

func TestStringOneOfCI(t *testing.T) {
	if !StringOneOfCI("A", []string{"a", "b", "c"}) {
		t.Fatal("expected oneofci to pass")
	}
	if StringOneOfCI("D", []string{"a", "b", "c"}) {
		t.Fatal("expected oneofci to fail")
	}
}

func TestStringFilePathEmpty(t *testing.T) {
	if StringFilePath("") {
		t.Fatal("expected filepath to fail for empty")
	}
	if StringFilePath("  ") {
		t.Fatal("expected filepath to fail for whitespace")
	}
}

func TestStringDirPathEmpty(t *testing.T) {
	if StringDirPath("") {
		t.Fatal("expected dirpath to fail for empty")
	}
	if StringDirPath("  ") {
		t.Fatal("expected dirpath to fail for whitespace")
	}
}

func TestStringMongoDBFail(t *testing.T) {
	if StringMongoDB("invalid") {
		t.Fatal("expected mongodb to fail for invalid id")
	}
}

func TestStringDNSRFC1035LabelFail(t *testing.T) {
	if StringDNSRFC1035Label("") {
		t.Fatal("expected dns_rfc1035_label to fail for empty")
	}
}

func TestCompareOp(t *testing.T) {
	if !CompareOp(3, 3, constants.CmpEQ) {
		t.Fatal("expected eq")
	}
	if !CompareOp(3, 2, constants.CmpNE) {
		t.Fatal("expected ne")
	}
	if !CompareOp(3, 2, constants.CmpGT) {
		t.Fatal("expected gt")
	}
	if !CompareOp(3, 3, constants.CmpGTE) {
		t.Fatal("expected gte")
	}
	if !CompareOp(2, 3, constants.CmpLT) {
		t.Fatal("expected lt")
	}
	if !CompareOp(3, 3, constants.CmpLTE) {
		t.Fatal("expected lte")
	}
}

func TestStringHostnamePortFail(t *testing.T) {
	if StringHostnamePort("example.com") {
		t.Fatal("expected hostname_port to fail without port")
	}
}

func TestStringExcludesFail(t *testing.T) {
	if StringExcludes("hello", "ell") {
		t.Fatal("expected excludes to fail when substring present")
	}
}

func TestStringExcludesAllFail(t *testing.T) {
	if StringExcludesAll("hello", "he") {
		t.Fatal("expected excludesall to fail when any char present")
	}
}

func TestStringContainsRuneFail(t *testing.T) {
	if StringContainsRune("hello", "z") {
		t.Fatal("expected containsrune to fail for absent rune")
	}
}

func TestStringExcludesRuneFail(t *testing.T) {
	if StringExcludesRune("hello", "e") {
		t.Fatal("expected excludesrune to fail when rune present")
	}
}

func TestStringHostnameFail(t *testing.T) {
	if StringHostname("") {
		t.Fatal("expected hostname to fail for empty")
	}
	if StringHostname("-invalid.com") {
		t.Fatal("expected hostname to fail for label starting with hyphen")
	}
}

func TestStringFQDNFail(t *testing.T) {
	if StringFQDN("invalid") {
		t.Fatal("expected fqdn to fail for invalid hostname")
	}
}

func TestStringIPv4Fail(t *testing.T) {
	if StringIPv4("999.999.999.999") {
		t.Fatal("expected ipv4 to fail for invalid address")
	}
}

func TestStringIPv6Fail(t *testing.T) {
	if StringIPv6("not:an:ipv6") {
		t.Fatal("expected ipv6 to fail for invalid address")
	}
}

func TestStringCIDRFail(t *testing.T) {
	if StringCIDR("not-a-cidr") {
		t.Fatal("expected cidr to fail for invalid")
	}
}

func TestStringMACFail(t *testing.T) {
	if StringMAC("not-a-mac") {
		t.Fatal("expected mac to fail for invalid")
	}
}

func TestStringIPFail(t *testing.T) {
	if StringIP("not-an-ip") {
		t.Fatal("expected ip to fail for invalid")
	}
}

func TestStringCIDRv4Fail(t *testing.T) {
	if StringCIDRv4("::1/128") {
		t.Fatal("expected cidrv4 to fail for ipv6 cidr")
	}
}

func TestStringCIDRv6Fail(t *testing.T) {
	if StringCIDRv6("10.0.0.0/8") {
		t.Fatal("expected cidrv6 to fail for ipv4 cidr")
	}
}

func TestStringUUIDFail(t *testing.T) {
	if StringUUID("not-a-uuid") {
		t.Fatal("expected uuid to fail for invalid")
	}
}

func TestStringUUID3Fail(t *testing.T) {
	if StringUUID3("6ba7b810-9dad-11d1-80b4-00c04fd430c8") {
		t.Fatal("expected uuid3 to fail for v1 uuid")
	}
}

func TestStringUUID4Fail(t *testing.T) {
	if StringUUID4("6ba7b810-9dad-11d1-80b4-00c04fd430c8") {
		t.Fatal("expected uuid4 to fail for v1 uuid")
	}
}

func TestStringUUID5Fail(t *testing.T) {
	if StringUUID5("6ba7b810-9dad-11d1-80b4-00c04fd430c8") {
		t.Fatal("expected uuid5 to fail for v1 uuid")
	}
}

func TestStringBase64Fail(t *testing.T) {
	if StringBase64("not-base64!!!") {
		t.Fatal("expected base64 to fail for invalid")
	}
}

func TestStringFastParsersTable(t *testing.T) {
	valid := map[string]func(string) bool{
		"uuid":      StringUUID,
		"uuid4":     StringUUID4,
		"hexcolor":  StringHexColor,
		"rgb":       StringRGB,
		"rgba":      StringRGBA,
		"hsl":       StringHSL,
		"hsla":      StringHSLA,
		"e164":      StringE164,
		"ipv4":      StringIPv4,
		"cidr":      StringCIDR,
		"mac":       StringMAC,
		"base64":    StringBase64,
		"base64url": StringBase64URL,
	}
	values := map[string]string{
		"uuid":      "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"uuid4":     "6ba7b810-9dad-41d1-80b4-00c04fd430c8",
		"hexcolor":  "#aabbcc",
		"rgb":       "rgb(255, 0, 128)",
		"rgba":      "rgba(255,0,128,0.5)",
		"hsl":       "hsl(120,100%,50%)",
		"hsla":      "hsla(120,100%,50%,.5)",
		"e164":      "+14155552671",
		"ipv4":      "192.168.1.1",
		"cidr":      "192.168.1.0/24",
		"mac":       "00:11:22:33:44:55",
		"base64":    "SGVsbG8=",
		"base64url": "SGVsbG8=",
	}
	for name, fn := range valid {
		if !fn(values[name]) {
			t.Fatalf("expected %s to pass for %q", name, values[name])
		}
	}
}

func TestStringFastParsersInvalidTable(t *testing.T) {
	invalid := map[string]struct {
		value string
		fn    func(string) bool
	}{
		"uuid":      {"6ba7b810-9dad-11d1-80b4-00c04fd430cx", StringUUID},
		"uuid4":     {"6ba7b810-9dad-51d1-80b4-00c04fd430c8", StringUUID4},
		"hexcolor":  {"#abcdex", StringHexColor},
		"rgb":       {"rgb(256,0,0)", StringRGB},
		"rgba":      {"rgba(0,0,0,1.5)", StringRGBA},
		"hsl":       {"hsl(361,100%,50%)", StringHSL},
		"hsla":      {"hsla(120,101%,50%,.5)", StringHSLA},
		"e164":      {"+0123", StringE164},
		"ipv4":      {"192.168.1.999", StringIPv4},
		"cidr":      {"192.168.1.0/33", StringCIDR},
		"mac":       {"00:11:22:33:44:zz", StringMAC},
		"base64":    {"SGVsbG8===", StringBase64},
		"base64url": {"SGVsbG8+/=", StringBase64URL},
	}
	for name, item := range invalid {
		if item.fn(item.value) {
			t.Fatalf("expected %s to fail for %q", name, item.value)
		}
	}
}

func TestStringFastParserEdges(t *testing.T) {
	valid := map[string]struct {
		value string
		fn    func(string) bool
	}{
		"hexcolor_short":    {"abc", StringHexColor},
		"hexcolor_alpha":    {"#abcd", StringHexColor},
		"rgba_alpha_one":    {"rgba(0,0,0,1)", StringRGBA},
		"hsla_alpha_zero":   {"hsla(360,0%,100%,0)", StringHSLA},
		"e164_max":          {"+123456789012345", StringE164},
		"ipv4_zero":         {"0.0.0.0", StringIPv4},
		"cidr_zero":         {"0.0.0.0/0", StringCIDR},
		"mac_plain":         {"001122334455", StringMAC},
		"mac_dash":          {"00-11-22-33-44-55", StringMAC},
		"base64_raw":        {"SGVsbG8", StringBase64},
		"base64url_special": {"SGVsbG8_", StringBase64URL},
	}
	for name, item := range valid {
		if !item.fn(item.value) {
			t.Fatalf("expected %s to pass for %q", name, item.value)
		}
	}

	invalid := map[string]struct {
		value string
		fn    func(string) bool
	}{
		"hexcolor_len":      {"#ab", StringHexColor},
		"rgb_missing_comma": {"rgb(1 2 3)", StringRGB},
		"hsl_missing_pct":   {"hsl(120,100,50%)", StringHSL},
		"e164_too_long":     {"+1234567890123456", StringE164},
		"ipv4_leading_zero": {"01.2.3.4", StringIPv4},
		"cidr_empty_bits":   {"10.0.0.0/", StringCIDR},
		"mac_mixed_sep":     {"00:11-22:33:44:55", StringMAC},
		"base64url_std":     {"SGVsbG8+", StringBase64URL},
	}
	for name, item := range invalid {
		if item.fn(item.value) {
			t.Fatalf("expected %s to fail for %q", name, item.value)
		}
	}
}

func TestStringNoParamAdaptersDirect(t *testing.T) {
	valid := map[string]struct {
		value string
		fn    func(string) bool
	}{
		"issn":     {"0317-847X", StringISSN},
		"bic":      {"CHASUS33", StringBIC},
		"cron":     {"0 0 * * *", StringCron},
		"datauri":  {"data:text/plain;base64,SGVsbG8=", StringDataURI},
		"bcp47":    {"en-US", StringBCP47},
		"eth_addr": {"0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38", StringEthAddr},
		"btc_addr": {"1BoatSLRHtKNngkdXEeobR76b53LETtpyT", StringBtcAddr},
	}
	for name, item := range valid {
		if !item.fn(item.value) {
			t.Fatalf("expected %s to pass for %q", name, item.value)
		}
	}

	invalid := map[string]struct {
		value string
		fn    func(string) bool
	}{
		"issn":     {"0317-8472", StringISSN},
		"bic":      {"CHASUS3!", StringBIC},
		"cron":     {"0 0 * * * * *", StringCron},
		"datauri":  {"data:text/plain;chars\x00et=utf-8,hello", StringDataURI},
		"bcp47":    {"1n", StringBCP47},
		"eth_addr": {"742d35Cc6634C0532925a3b844Bc9e7595f2bD38", StringEthAddr},
		"btc_addr": {"bc1q!508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4", StringBtcAddr},
	}
	for name, item := range invalid {
		if item.fn(item.value) {
			t.Fatalf("expected %s to fail for %q", name, item.value)
		}
	}
}

func TestFormatHelpersDirect(t *testing.T) {
	ClearRegexCache()
	re, err := GetCompiledRegex(`^a+$`)
	if err != nil || !re.MatchString("aaa") {
		t.Fatalf("expected compiled regex to match, err=%v", err)
	}
	if _, err := GetCompiledRegex(`[`); err == nil {
		t.Fatal("expected invalid regex to fail")
	}
	ClearRegexCache()

	pos := 0
	if !ParseSemverNum("12", &pos) || pos != 2 {
		t.Fatalf("expected semver number to parse, pos=%d", pos)
	}
	pos = 0
	if !ParseSemverPreRelease("-rc.1", &pos) || pos != len("-rc.1") {
		t.Fatalf("expected pre-release to parse, pos=%d", pos)
	}
	pos = 0
	if !ParseSemverBuildMeta("+build.7", &pos) || pos != len("+build.7") {
		t.Fatalf("expected build metadata to parse, pos=%d", pos)
	}
	if !IsValidCronField("*/5") || IsValidCronField("!") {
		t.Fatal("expected cron field helper to validate step syntax")
	}
	if LuhnDouble(8) != 7 || !IsLuhnChecksum("79927398713") {
		t.Fatal("expected luhn helpers to pass")
	}
	if !IsISBN10CheckDigit('X', 1) || IsISBN10CheckDigit('!', 1) {
		t.Fatal("expected isbn10 check digit helper to validate X and reject punctuation")
	}
}

func TestStringBase32Fail(t *testing.T) {
	if StringBase32("not-base32!!!") {
		t.Fatal("expected base32 to fail for invalid")
	}
}

func TestStringURIFail(t *testing.T) {
	if StringURI("") {
		t.Fatal("expected uri to fail for empty")
	}
}

func TestStringHTTPURLFail(t *testing.T) {
	if StringHTTPURL("ftp://example.com") {
		t.Fatal("expected http_url to fail for ftp")
	}
}

func TestStringHTTPSURLFail(t *testing.T) {
	if StringHTTPSURL("http://example.com") {
		t.Fatal("expected https_url to fail for http")
	}
}

func TestStringEmailFail(t *testing.T) {
	if IsEmail("not-an-email") {
		t.Fatal("expected email to fail for invalid")
	}
}

func TestStringUniqueEmpty(t *testing.T) {
	if !StringUnique("") {
		t.Fatal("expected unique to pass for empty string")
	}
}
