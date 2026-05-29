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

	"github.com/kamalyes/go-argus/constants"
)

func TestCompareNumbersEqual(t *testing.T) {
	r := CompareNumbers(10, 10, constants.OpEqual)
	if !r.Success {
		t.Fatal("expected equal")
	}
}

func TestCompareNumbersSymbolEqual(t *testing.T) {
	r := CompareNumbers(10, 10, constants.OpSymbolEqual)
	if !r.Success {
		t.Fatal("expected symbol equal")
	}
}

func TestCompareNumbersNotEqual(t *testing.T) {
	r := CompareNumbers(10, 20, constants.OpNotEqual)
	if !r.Success {
		t.Fatal("expected not equal")
	}
}

func TestCompareNumbersSymbolNotEqual(t *testing.T) {
	r := CompareNumbers(10, 20, constants.OpSymbolNotEqual)
	if !r.Success {
		t.Fatal("expected symbol not equal")
	}
}

func TestCompareNumbersGreaterThan(t *testing.T) {
	r := CompareNumbers(20, 10, constants.OpGreaterThan)
	if !r.Success {
		t.Fatal("expected greater than")
	}
}

func TestCompareNumbersSymbolGreaterThan(t *testing.T) {
	r := CompareNumbers(20, 10, constants.OpSymbolGreaterThan)
	if !r.Success {
		t.Fatal("expected symbol greater than")
	}
}

func TestCompareNumbersGreaterThanOrEqual(t *testing.T) {
	r := CompareNumbers(10, 10, constants.OpGreaterThanOrEqual)
	if !r.Success {
		t.Fatal("expected greater than or equal")
	}
}

func TestCompareNumbersSymbolGreaterThanOrEqual(t *testing.T) {
	r := CompareNumbers(10, 10, constants.OpSymbolGreaterThanOrEqual)
	if !r.Success {
		t.Fatal("expected symbol greater than or equal")
	}
}

func TestCompareNumbersLessThan(t *testing.T) {
	r := CompareNumbers(5, 10, constants.OpLessThan)
	if !r.Success {
		t.Fatal("expected less than")
	}
}

func TestCompareNumbersSymbolLessThan(t *testing.T) {
	r := CompareNumbers(5, 10, constants.OpSymbolLessThan)
	if !r.Success {
		t.Fatal("expected symbol less than")
	}
}

func TestCompareNumbersLessThanOrEqual(t *testing.T) {
	r := CompareNumbers(10, 10, constants.OpLessThanOrEqual)
	if !r.Success {
		t.Fatal("expected less than or equal")
	}
}

func TestCompareNumbersSymbolLessThanOrEqual(t *testing.T) {
	r := CompareNumbers(10, 10, constants.OpSymbolLessThanOrEqual)
	if !r.Success {
		t.Fatal("expected symbol less than or equal")
	}
}

func TestCompareNumbersUnsupportedOp(t *testing.T) {
	r := CompareNumbers(10, 10, constants.OpContains)
	if r.Success || r.Message == "" {
		t.Fatal("expected unsupported op to fail with message")
	}
}

func TestCompareNumbersFailed(t *testing.T) {
	r := CompareNumbers(5, 10, constants.OpGreaterThan)
	if r.Success {
		t.Fatal("expected gt to fail")
	}
	if r.Message == "" {
		t.Fatal("expected failure message")
	}
}

func TestCompareStringsEqual(t *testing.T) {
	r := CompareStrings("hello", "hello", constants.OpEqual)
	if !r.Success {
		t.Fatal("expected equal")
	}
}

func TestCompareStringsSymbolEqual(t *testing.T) {
	r := CompareStrings("hello", "hello", constants.OpSymbolEqual)
	if !r.Success {
		t.Fatal("expected symbol equal")
	}
}

func TestCompareStringsNotEqual(t *testing.T) {
	r := CompareStrings("hello", "world", constants.OpNotEqual)
	if !r.Success {
		t.Fatal("expected not equal")
	}
}

func TestCompareStringsSymbolNotEqual(t *testing.T) {
	r := CompareStrings("hello", "world", constants.OpSymbolNotEqual)
	if !r.Success {
		t.Fatal("expected symbol not equal")
	}
}

func TestCompareStringsContains(t *testing.T) {
	r := CompareStrings("hello world", "world", constants.OpContains)
	if !r.Success {
		t.Fatal("expected contains")
	}
}

func TestCompareStringsNotContains(t *testing.T) {
	r := CompareStrings("hello world", "xyz", constants.OpNotContains)
	if !r.Success {
		t.Fatal("expected not contains")
	}
}

func TestCompareStringsHasPrefix(t *testing.T) {
	r := CompareStrings("hello world", "hello", constants.OpHasPrefix)
	if !r.Success {
		t.Fatal("expected has prefix")
	}
}

func TestCompareStringsHasSuffix(t *testing.T) {
	r := CompareStrings("hello world", "world", constants.OpHasSuffix)
	if !r.Success {
		t.Fatal("expected has suffix")
	}
}

func TestCompareStringsEmpty(t *testing.T) {
	r := CompareStrings("   ", "", constants.OpEmpty)
	if !r.Success {
		t.Fatal("expected empty")
	}
}

func TestCompareStringsNotEmpty(t *testing.T) {
	r := CompareStrings("hello", "", constants.OpNotEmpty)
	if !r.Success {
		t.Fatal("expected not empty")
	}
}

func TestCompareStringsRegex(t *testing.T) {
	r := CompareStrings("hello123", `\d+`, constants.OpRegex)
	if !r.Success {
		t.Fatal("expected regex match")
	}
}

func TestCompareStringsRegexCompileError(t *testing.T) {
	r := CompareStrings("hello", `[`, constants.OpRegex)
	if r.Success || r.Message == "" {
		t.Fatal("expected regex compile error")
	}
}

func TestCompareStringsUnsupportedOp(t *testing.T) {
	r := CompareStrings("hello", "world", constants.OpGreaterThan)
	if r.Success || r.Message == "" {
		t.Fatal("expected unsupported op to fail with message")
	}
}

func TestCompareStringsFailed(t *testing.T) {
	r := CompareStrings("abc", "xyz", constants.OpContains)
	if r.Success {
		t.Fatal("expected contains to fail")
	}
	if r.Message == "" {
		t.Fatal("expected failure message")
	}
}

func TestValidateString(t *testing.T) {
	r := ValidateString("hello", "hello", constants.OpEqual)
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
	r := ValidateStatusCode(200, 200, constants.OpEqual)
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
	r := ValidateHeader(headers, "Content-Type", "application/json", constants.OpEqual)
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
	if constants.OpEqual.String() != "eq" {
		t.Fatalf("expected eq, got %s", constants.OpEqual.String())
	}
}
