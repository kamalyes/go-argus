/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\compare_test.go
 * @Description: compare.go 测试，覆盖数值比较、字符串比较、状态码和 Header 校验
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validate

import (
	"testing"
)

func TestCompareNumbersEqual(t *testing.T) {
	r := CompareNumbers(10, 10, OpEqual)
	if !r.Success {
		t.Fatal("expected equal")
	}
}

func TestCompareNumbersSymbolEqual(t *testing.T) {
	r := CompareNumbers(10, 10, OpSymbolEqual)
	if !r.Success {
		t.Fatal("expected symbol equal")
	}
}

func TestCompareNumbersNotEqual(t *testing.T) {
	r := CompareNumbers(10, 20, OpNotEqual)
	if !r.Success {
		t.Fatal("expected not equal")
	}
}

func TestCompareNumbersSymbolNotEqual(t *testing.T) {
	r := CompareNumbers(10, 20, OpSymbolNotEqual)
	if !r.Success {
		t.Fatal("expected symbol not equal")
	}
}

func TestCompareNumbersGreaterThan(t *testing.T) {
	r := CompareNumbers(20, 10, OpGreaterThan)
	if !r.Success {
		t.Fatal("expected greater than")
	}
}

func TestCompareNumbersSymbolGreaterThan(t *testing.T) {
	r := CompareNumbers(20, 10, OpSymbolGreaterThan)
	if !r.Success {
		t.Fatal("expected symbol greater than")
	}
}

func TestCompareNumbersGreaterThanOrEqual(t *testing.T) {
	r := CompareNumbers(10, 10, OpGreaterThanOrEqual)
	if !r.Success {
		t.Fatal("expected greater than or equal")
	}
}

func TestCompareNumbersSymbolGreaterThanOrEqual(t *testing.T) {
	r := CompareNumbers(10, 10, OpSymbolGreaterThanOrEqual)
	if !r.Success {
		t.Fatal("expected symbol greater than or equal")
	}
}

func TestCompareNumbersLessThan(t *testing.T) {
	r := CompareNumbers(5, 10, OpLessThan)
	if !r.Success {
		t.Fatal("expected less than")
	}
}

func TestCompareNumbersSymbolLessThan(t *testing.T) {
	r := CompareNumbers(5, 10, OpSymbolLessThan)
	if !r.Success {
		t.Fatal("expected symbol less than")
	}
}

func TestCompareNumbersLessThanOrEqual(t *testing.T) {
	r := CompareNumbers(10, 10, OpLessThanOrEqual)
	if !r.Success {
		t.Fatal("expected less than or equal")
	}
}

func TestCompareNumbersSymbolLessThanOrEqual(t *testing.T) {
	r := CompareNumbers(10, 10, OpSymbolLessThanOrEqual)
	if !r.Success {
		t.Fatal("expected symbol less than or equal")
	}
}

func TestCompareNumbersUnsupportedOp(t *testing.T) {
	r := CompareNumbers(10, 10, OpContains)
	if r.Success || r.Message == "" {
		t.Fatal("expected unsupported op to fail with message")
	}
}

func TestCompareNumbersFailed(t *testing.T) {
	r := CompareNumbers(5, 10, OpGreaterThan)
	if r.Success {
		t.Fatal("expected gt to fail")
	}
	if r.Message == "" {
		t.Fatal("expected failure message")
	}
}

func TestCompareStringsEqual(t *testing.T) {
	r := CompareStrings("hello", "hello", OpEqual)
	if !r.Success {
		t.Fatal("expected equal")
	}
}

func TestCompareStringsSymbolEqual(t *testing.T) {
	r := CompareStrings("hello", "hello", OpSymbolEqual)
	if !r.Success {
		t.Fatal("expected symbol equal")
	}
}

func TestCompareStringsNotEqual(t *testing.T) {
	r := CompareStrings("hello", "world", OpNotEqual)
	if !r.Success {
		t.Fatal("expected not equal")
	}
}

func TestCompareStringsSymbolNotEqual(t *testing.T) {
	r := CompareStrings("hello", "world", OpSymbolNotEqual)
	if !r.Success {
		t.Fatal("expected symbol not equal")
	}
}

func TestCompareStringsContains(t *testing.T) {
	r := CompareStrings("hello world", "world", OpContains)
	if !r.Success {
		t.Fatal("expected contains")
	}
}

func TestCompareStringsNotContains(t *testing.T) {
	r := CompareStrings("hello world", "xyz", OpNotContains)
	if !r.Success {
		t.Fatal("expected not contains")
	}
}

func TestCompareStringsHasPrefix(t *testing.T) {
	r := CompareStrings("hello world", "hello", OpHasPrefix)
	if !r.Success {
		t.Fatal("expected has prefix")
	}
}

func TestCompareStringsHasSuffix(t *testing.T) {
	r := CompareStrings("hello world", "world", OpHasSuffix)
	if !r.Success {
		t.Fatal("expected has suffix")
	}
}

func TestCompareStringsEmpty(t *testing.T) {
	r := CompareStrings("   ", "", OpEmpty)
	if !r.Success {
		t.Fatal("expected empty")
	}
}

func TestCompareStringsNotEmpty(t *testing.T) {
	r := CompareStrings("hello", "", OpNotEmpty)
	if !r.Success {
		t.Fatal("expected not empty")
	}
}

func TestCompareStringsRegex(t *testing.T) {
	r := CompareStrings("hello123", `\d+`, OpRegex)
	if !r.Success {
		t.Fatal("expected regex match")
	}
}

func TestCompareStringsRegexCompileError(t *testing.T) {
	r := CompareStrings("hello", `[`, OpRegex)
	if r.Success || r.Message == "" {
		t.Fatal("expected regex compile error")
	}
}

func TestCompareStringsUnsupportedOp(t *testing.T) {
	r := CompareStrings("hello", "world", OpGreaterThan)
	if r.Success || r.Message == "" {
		t.Fatal("expected unsupported op to fail with message")
	}
}

func TestCompareStringsFailed(t *testing.T) {
	r := CompareStrings("abc", "xyz", OpContains)
	if r.Success {
		t.Fatal("expected contains to fail")
	}
	if r.Message == "" {
		t.Fatal("expected failure message")
	}
}

func TestValidateString(t *testing.T) {
	r := ValidateString("hello", "hello", OpEqual)
	if !r.Success {
		t.Fatal("expected equal")
	}
}

func TestValidateContains(t *testing.T) {
	r := ValidateContains([]byte("hello world"), "world")
	if !r.Success {
		t.Fatal("expected contains")
	}
}

func TestValidateNotContains(t *testing.T) {
	r := ValidateNotContains([]byte("hello world"), "xyz")
	if !r.Success {
		t.Fatal("expected not contains")
	}
}

func TestValidateStatusCode(t *testing.T) {
	r := ValidateStatusCode(200, 200, OpEqual)
	if !r.Success {
		t.Fatal("expected status code equal")
	}
}

func TestValidateStatusCodeRange(t *testing.T) {
	r := ValidateStatusCodeRange(200, 200, 300)
	if !r.Success {
		t.Fatal("expected in range")
	}
}

func TestValidateStatusCodeRangeFail(t *testing.T) {
	r := ValidateStatusCodeRange(500, 200, 300)
	if r.Success {
		t.Fatal("expected out of range")
	}
	if r.Message == "" {
		t.Fatal("expected failure message")
	}
}

func TestValidateHeader(t *testing.T) {
	headers := map[string]string{"Content-Type": "application/json"}
	r := ValidateHeader(headers, "Content-Type", "application/json", OpEqual)
	if !r.Success {
		t.Fatal("expected header equal")
	}
}

func TestValidateContentType(t *testing.T) {
	headers := map[string]string{"Content-Type": "application/json; charset=utf-8"}
	r := ValidateContentType(headers, "application/json")
	if !r.Success {
		t.Fatal("expected content type match")
	}
}

func TestValidateContentTypeLowerKey(t *testing.T) {
	headers := map[string]string{"content-type": "text/html"}
	r := ValidateContentType(headers, "text/html")
	if !r.Success {
		t.Fatal("expected content type match with lower key")
	}
}

func TestCompareOperatorString(t *testing.T) {
	if OpEqual.String() != "eq" {
		t.Fatalf("expected eq, got %s", OpEqual.String())
	}
}

func TestFmtOp(t *testing.T) {
	if fmtOp(OpEqual) != "eq" {
		t.Fatalf("expected eq, got %s", fmtOp(OpEqual))
	}
}
