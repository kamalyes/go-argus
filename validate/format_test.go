/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\format_test.go
 * @Description: format.go 测试，覆盖 Email、IP、URL、UUID、Base64 和正则校验
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validate

import (
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
	r := ValidateISSN("0317847X")
	if !r.Success {
		t.Fatal("expected valid ISSN")
	}
	r = ValidateISSN("03178470")
	if r.Success || r.Message == "" {
		t.Fatal("expected wrong checksum to fail")
	}
}

func TestIsISSN(t *testing.T) {
	if !IsISSN("0317847X") {
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
