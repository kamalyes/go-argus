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
