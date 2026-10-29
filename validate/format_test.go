/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-27 19:58:37
 * @FilePath: \go-argus\validate\format_test.go
 * @Description: format.go 测试，覆盖 Email、IP、URL、UUID、Base64 和正则校验
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validate

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestValidateRegex(t *testing.T) {
	r := ValidateRegex([]byte("hello123"), `\d+`)
	if !r.Success {
		t.Fatal("expected regex match")
	}
}

func TestValidateRegexCompileError(t *testing.T) {
	r := ValidateRegex([]byte("hello"), `[`)
	if r.Success || r.Message == "" {
		t.Fatal("expected regex compile error")
	}
}

func TestValidateRegexNotMatched(t *testing.T) {
	r := ValidateRegex([]byte("hello"), `^\d+$`)
	if r.Success {
		t.Fatal("expected regex not matched")
	}
}

func TestValidateEmailValid(t *testing.T) {
	r := ValidateEmail("hello@example.com")
	if !r.Success {
		t.Fatal("expected valid email")
	}
}

func TestValidateEmailEmpty(t *testing.T) {
	r := ValidateEmail("")
	if r.Success || r.Message == "" {
		t.Fatal("expected empty email to fail")
	}
}

func TestValidateEmailInvalid(t *testing.T) {
	r := ValidateEmail("not-an-email")
	if r.Success || r.Message == "" {
		t.Fatal("expected invalid email to fail")
	}
}

func TestValidateEmailMalformed(t *testing.T) {
	r := ValidateEmail("hello@.com")
	if r.Success || r.Message == "" {
		t.Fatal("expected malformed email to fail")
	}
}

func TestValidateEmailNoLocalPart(t *testing.T) {
	r := ValidateEmail("@example.com")
	if r.Success || r.Message == "" {
		t.Fatal("expected no local part to fail")
	}
}

func TestValidateEmailNoDomain(t *testing.T) {
	r := ValidateEmail("hello@")
	if r.Success || r.Message == "" {
		t.Fatal("expected no domain to fail")
	}
}

func TestValidateEmailNoDotInDomain(t *testing.T) {
	r := ValidateEmail("hello@localhost")
	if r.Success || r.Message == "" {
		t.Fatal("expected no dot in domain to fail")
	}
}

func TestValidateIPAddressValid(t *testing.T) {
	r := ValidateIPAddress("192.168.1.1")
	if !r.Success {
		t.Fatal("expected valid IP")
	}
}

func TestValidateIPAddressInvalid(t *testing.T) {
	r := ValidateIPAddress("999.999.999.999")
	if r.Success {
		t.Fatal("expected invalid IP to fail")
	}
}

func TestValidateIPAlias(t *testing.T) {
	r := ValidateIP("127.0.0.1")
	if !r.Success {
		t.Fatal("expected ValidateIP alias to pass")
	}
}

func TestValidateProtocolValid(t *testing.T) {
	r := ValidateProtocol("https://example.com")
	if !r.Success {
		t.Fatal("expected valid protocol")
	}
}

func TestValidateProtocolCustom(t *testing.T) {
	r := ValidateProtocol("ftp://example.com", "ftp")
	if !r.Success {
		t.Fatal("expected custom protocol")
	}
}

func TestValidateProtocolMissingScheme(t *testing.T) {
	r := ValidateProtocol("example.com")
	if r.Success {
		t.Fatal("expected missing scheme to fail")
	}
}

func TestValidateProtocolUnsupported(t *testing.T) {
	r := ValidateProtocol("ssh://example.com")
	if r.Success {
		t.Fatal("expected unsupported protocol to fail")
	}
}

func TestValidateHTTP(t *testing.T) {
	r := ValidateHTTP("http://example.com")
	if !r.Success {
		t.Fatal("expected http to be valid")
	}
}

func TestValidateWebSocket(t *testing.T) {
	r := ValidateWebSocket("ws://example.com")
	if !r.Success {
		t.Fatal("expected ws to be valid")
	}
}

func TestValidateUUIDValid(t *testing.T) {
	r := ValidateUUID("550e8400-e29b-41d4-a716-446655440000")
	if !r.Success {
		t.Fatal("expected valid UUID")
	}
}

func TestValidateUUIDInvalid(t *testing.T) {
	r := ValidateUUID("not-a-uuid")
	if r.Success {
		t.Fatal("expected invalid UUID to fail")
	}
}

func TestValidateBase64Valid(t *testing.T) {
	r := ValidateBase64("SGVsbG8gV29ybGQ=")
	if !r.Success {
		t.Fatal("expected valid base64")
	}
}

func TestValidateBase64Empty(t *testing.T) {
	r := ValidateBase64("")
	if r.Success || r.Message == "" {
		t.Fatal("expected empty base64 to fail")
	}
}

func TestValidateBase64Invalid(t *testing.T) {
	r := ValidateBase64("!!!invalid!!!")
	if r.Success {
		t.Fatal("expected invalid base64 to fail")
	}
}

func TestIsEmail(t *testing.T) {
	if !IsEmail("hello@example.com") {
		t.Fatal("expected true")
	}
	if IsEmail("not-email") {
		t.Fatal("expected false")
	}
}

func TestIsIP(t *testing.T) {
	if !IsIP("192.168.1.1") {
		t.Fatal("expected true")
	}
	if IsIP("invalid") {
		t.Fatal("expected false")
	}
}

func TestIsUUID(t *testing.T) {
	if !IsUUID("550e8400-e29b-41d4-a716-446655440000") {
		t.Fatal("expected true")
	}
	if IsUUID("invalid") {
		t.Fatal("expected false")
	}
}

// uuidByScan 是 IsUUID 改成「先验分隔符、再分段扫十六进制」之前的逐字符实现
//
// 保留它作为等价性参照：分段版把 32 个数据字符上的 4 路位置判定省掉了，
// 该优化的唯一正确性依据就是两者判定完全一致
func uuidByScan(uuid string) bool {
	uuid = strings.TrimSpace(uuid)
	if len(uuid) != 36 {
		return false
	}
	for i := 0; i < len(uuid); i++ {
		switch i {
		case 8, 13, 18, 23:
			if uuid[i] != '-' {
				return false
			}
		default:
			if !IsHexChar(uuid[i]) {
				return false
			}
		}
	}
	return true
}

// TestIsUUID_MatchesCharScan 守住分段版与逐字符版的逐例等价
//
// 覆盖：有效 UUID 的每一位 × 字母表的单字符变异（分隔符位与数据位的全部边界）、
// 长度 30~41 的随机串，以及空串、空白、缺位、多位等退化输入
func TestIsUUID_MatchesCharScan(t *testing.T) {
	alphabet := []byte("0123456789abcdefABCDEF-gz_ .\t\n")
	cases := []string{
		"", " ", "\t", "550e8400-e29b-41d4-a716-446655440000",
		" 550e8400-e29b-41d4-a716-446655440000 ",
		"550e8400-e29b-41d4-a716-44665544000",
		"550e8400-e29b-41d4-a716-4466554400000",
		"550e8400e29b41d4a716446655440000",
	}
	base := []byte("550e8400-e29b-41d4-a716-446655440000")
	for pos := range base {
		for _, c := range alphabet {
			mutated := append([]byte(nil), base...)
			mutated[pos] = c
			cases = append(cases, string(mutated))
		}
	}
	rng := rand.New(rand.NewSource(1))
	for n := 0; n < 200000; n++ {
		b := make([]byte, 30+rng.Intn(12))
		for i := range b {
			b[i] = alphabet[rng.Intn(len(alphabet))]
		}
		cases = append(cases, string(b))
	}
	for _, s := range cases {
		if got, want := IsUUID(s), uuidByScan(s); got != want {
			t.Fatalf("IsUUID(%q) = %v，逐字符版 = %v", s, got, want)
		}
	}
}

func TestIsBase64(t *testing.T) {
	if !IsBase64("SGVsbG8=") {
		t.Fatal("expected true")
	}
	if IsBase64("!!!") {
		t.Fatal("expected false")
	}
}

func TestValidateSemver(t *testing.T) {
	r := ValidateSemver("1.2.3")
	if !r.Success {
		t.Fatal("expected valid semver")
	}
	r = ValidateSemver("v1.0.0-alpha.1+build.123")
	if !r.Success {
		t.Fatal("expected valid semver with pre-release and build")
	}
	r = ValidateSemver("1.2")
	if r.Success || r.Message == "" {
		t.Fatal("expected invalid semver to fail")
	}
	r = ValidateSemver("abc")
	if r.Success || r.Message == "" {
		t.Fatal("expected non-semver to fail")
	}
}

func TestIsSemver(t *testing.T) {
	if !IsSemver("1.2.3") {
		t.Fatal("expected true")
	}
	if IsSemver("1.2") {
		t.Fatal("expected false")
	}
}

func TestValidateISBN10(t *testing.T) {
	r := ValidateISBN10("0471958697")
	if !r.Success {
		t.Fatal("expected valid ISBN-10")
	}
	r = ValidateISBN10("0306406152")
	if !r.Success {
		t.Fatal("expected valid ISBN-10 with numeric check digit")
	}
	r = ValidateISBN10("0471958698")
	if r.Success || r.Message == "" {
		t.Fatal("expected wrong checksum to fail")
	}
}

func TestIsISBN10(t *testing.T) {
	if !IsISBN10("0471958697") {
		t.Fatal("expected true")
	}
	if IsISBN10("0471958698") {
		t.Fatal("expected false for wrong checksum")
	}
}

func TestValidateISBN13(t *testing.T) {
	r := ValidateISBN13("9780471117094")
	if !r.Success {
		t.Fatal("expected valid ISBN-13")
	}
	r = ValidateISBN13("9780471117095")
	if r.Success || r.Message == "" {
		t.Fatal("expected wrong checksum to fail")
	}
}

func TestIsISBN13(t *testing.T) {
	if !IsISBN13("9780471117094") {
		t.Fatal("expected true")
	}
	if IsISBN13("9780471117095") {
		t.Fatal("expected false for wrong checksum")
	}
}

func TestValidateISSN(t *testing.T) {
	r := ValidateISSN("03178471")
	if !r.Success {
		t.Fatal("expected valid ISSN")
	}
	r = ValidateISSN("03178470")
	if r.Success || r.Message == "" {
		t.Fatal("expected wrong checksum to fail")
	}
}

func TestIsISSN(t *testing.T) {
	if !IsISSN("03178471") {
		t.Fatal("expected true")
	}
	if IsISSN("03178470") {
		t.Fatal("expected false for wrong checksum")
	}
}

func TestValidateBIC(t *testing.T) {
	r := ValidateBIC("CHASUS33")
	if !r.Success {
		t.Fatal("expected valid BIC")
	}
	r = ValidateBIC("CHASUS33XXX")
	if !r.Success {
		t.Fatal("expected valid 11-char BIC")
	}
	r = ValidateBIC("INVALID")
	if r.Success || r.Message == "" {
		t.Fatal("expected invalid BIC to fail")
	}
}

func TestIsBIC(t *testing.T) {
	if !IsBIC("CHASUS33") {
		t.Fatal("expected true")
	}
	if IsBIC("INVALID") {
		t.Fatal("expected false")
	}
}

func TestValidateCron(t *testing.T) {
	r := ValidateCron("*/5 * * * *")
	if !r.Success {
		t.Fatal("expected valid cron")
	}
	r = ValidateCron("* * * *")
	if r.Success || r.Message == "" {
		t.Fatal("expected invalid cron (4 fields) to fail")
	}
}

func TestIsCron(t *testing.T) {
	if !IsCron("*/5 * * * *") {
		t.Fatal("expected true")
	}
	if IsCron("* * * *") {
		t.Fatal("expected false")
	}
}

func TestValidateDataURI(t *testing.T) {
	r := ValidateDataURI("data:text/plain;base64,SGVsbG8=")
	if !r.Success {
		t.Fatal("expected valid data URI")
	}
	r = ValidateDataURI("http://example.com")
	if r.Success || r.Message == "" {
		t.Fatal("expected non-data URI to fail")
	}
}

func TestIsDataURI(t *testing.T) {
	if !IsDataURI("data:text/plain;base64,SGVsbG8=") {
		t.Fatal("expected true")
	}
	if IsDataURI("http://example.com") {
		t.Fatal("expected false")
	}
}

func TestValidateBCP47(t *testing.T) {
	r := ValidateBCP47("zh-CN")
	if !r.Success {
		t.Fatal("expected valid BCP 47")
	}
	r = ValidateBCP47("1")
	if r.Success || r.Message == "" {
		t.Fatal("expected invalid BCP 47 to fail")
	}
}

func TestIsBCP47(t *testing.T) {
	if !IsBCP47("en-US") {
		t.Fatal("expected true")
	}
	if IsBCP47("1") {
		t.Fatal("expected false")
	}
}

func TestValidateEthAddr(t *testing.T) {
	r := ValidateEthAddr("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38")
	if !r.Success {
		t.Fatal("expected valid Ethereum address")
	}
	r = ValidateEthAddr("0x1234")
	if r.Success || r.Message == "" {
		t.Fatal("expected too-short address to fail")
	}
}

func TestIsEthAddr(t *testing.T) {
	if !IsEthAddr("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38") {
		t.Fatal("expected true")
	}
	if IsEthAddr("0x1234") {
		t.Fatal("expected false")
	}
}

func TestValidateBtcAddr(t *testing.T) {
	r := ValidateBtcAddr("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	if !r.Success {
		t.Fatal("expected valid Bitcoin address")
	}
	r = ValidateBtcAddr("not-a-btc-address")
	if r.Success || r.Message == "" {
		t.Fatal("expected invalid Bitcoin address to fail")
	}
}

func TestIsBtcAddr(t *testing.T) {
	if !IsBtcAddr("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa") {
		t.Fatal("expected true")
	}
	if IsBtcAddr("not-a-btc-address") {
		t.Fatal("expected false")
	}
}

// refISSN 按 ISO 3297 独立实现一遍：7 位主体加权 8..2，
// 校验位（X 记作 10）满足 (加权和 + 校验值) % 11 == 0
func refISSN(s string) bool {
	var ds []byte
	for i := 0; i < len(s); i++ {
		if c := s[i]; c != '-' && c != ' ' {
			ds = append(ds, c)
		}
	}
	if len(ds) != 8 {
		return false
	}
	sum := 0
	for i := 0; i < 7; i++ {
		if ds[i] < '0' || ds[i] > '9' {
			return false
		}
		sum += int(ds[i]-'0') * (8 - i)
	}
	last := ds[7]
	var check int
	switch {
	case last == 'X' || last == 'x':
		check = 10
	case last >= '0' && last <= '9':
		check = int(last - '0')
	default:
		return false
	}
	return (sum+check)%11 == 0
}

// refISBN13 按 ISO 2108 独立实现：13 位交替加权 1/3，总和 % 10 == 0
func refISBN13(s string) bool {
	var ds []byte
	for i := 0; i < len(s); i++ {
		if c := s[i]; c != '-' && c != ' ' {
			ds = append(ds, c)
		}
	}
	if len(ds) != 13 {
		return false
	}
	sum := 0
	for i := 0; i < 13; i++ {
		if ds[i] < '0' || ds[i] > '9' {
			return false
		}
		w := 1
		if i%2 == 1 {
			w = 3
		}
		sum += int(ds[i]-'0') * w
	}
	return sum%10 == 0
}

// TestISSNMatchesSpec ISSN 必须与规范实现逐例同结论
func TestISSNMatchesSpec(t *testing.T) {
	cases := []string{
		// 真实合法 ISSN（校验位分别落在 1 / 0 / 5 / X 上）
		"0317-8471", "2049-3630", "0378-5955", "1050-124X",
		"03178471", "20493630",
		// 非法
		"0317-847X", "0317-8470", "0317-8472", "1050-1240", "1050-1241",
		"", "0317-847", "0317-84711", "0317-847Y", "abcd-efgh",
	}
	r := rand.New(rand.NewSource(20261101))
	for i := 0; i < 20000; i++ {
		n := 6 + r.Intn(5)
		var b strings.Builder
		for j := 0; j < n; j++ {
			b.WriteByte("0123456789Xx- "[r.Intn(14)])
		}
		cases = append(cases, b.String())
	}
	// 结构化随机：7 位随机主体 × 全部 11 种校验位，逐一对照
	for i := 0; i < 2000; i++ {
		body := ""
		for j := 0; j < 7; j++ {
			body += strconv.Itoa(r.Intn(10))
		}
		for _, c := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "X"} {
			cases = append(cases, body[:4]+"-"+body[4:]+c)
		}
	}
	for _, s := range cases {
		if got, want := IsISSN(s), refISSN(s); got != want {
			t.Fatalf("IsISSN(%q) = %v, 规范实现 = %v", s, got, want)
		}
	}
}

// TestISBN13MatchesSpec ISBN-13 必须与规范实现逐例同结论
func TestISBN13MatchesSpec(t *testing.T) {
	cases := []string{
		"978-0-306-40615-7", "9780306406157", "978-3-16-148410-0",
		"978-0-471-95869-0", "9780306406150", "9780306406155",
		"", "97803064061", "97803064061577", "978030640615X",
	}
	r := rand.New(rand.NewSource(20261102))
	// 结构化随机：12 位随机主体 × 全部 10 种校验位，逐一对照
	for i := 0; i < 4000; i++ {
		body := ""
		for j := 0; j < 12; j++ {
			body += strconv.Itoa(r.Intn(10))
		}
		for d := 0; d < 10; d++ {
			cases = append(cases, body+strconv.Itoa(d))
		}
	}
	for _, s := range cases {
		if got, want := IsISBN13(s), refISBN13(s); got != want {
			t.Fatalf("IsISBN13(%q) = %v, 规范实现 = %v", s, got, want)
		}
	}
}
