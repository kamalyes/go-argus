/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-20 00:00:00
 * @FilePath: \go-argus\compare_test.go
 * @Description: compare.go 测试，覆盖数值比较、字符串比较、状态码校验等包装函数
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"
)

func TestCompareNumbers(t *testing.T) {
	r := CompareNumbers(10, 5, OpGreaterThan)
	if !r.Success {
		t.Fatal("expected 10 > 5")
	}
	r = CompareNumbers(3, 5, OpLessThan)
	if !r.Success {
		t.Fatal("expected 3 < 5")
	}
	r = CompareNumbers(5, 5, OpEqual)
	if !r.Success {
		t.Fatal("expected 5 == 5")
	}
	r = CompareNumbers(5, 5, OpGreaterThanOrEqual)
	if !r.Success {
		t.Fatal("expected 5 >= 5")
	}
	r = CompareNumbers(5, 5, OpLessThanOrEqual)
	if !r.Success {
		t.Fatal("expected 5 <= 5")
	}
	r = CompareNumbers(3, 5, OpNotEqual)
	if !r.Success {
		t.Fatal("expected 3 != 5")
	}
}

func TestCompareStrings(t *testing.T) {
	r := CompareStrings("abc", "abc", OpEqual)
	if !r.Success {
		t.Fatal("expected strings equal")
	}
	r = CompareStrings("abc", "def", OpNotEqual)
	if !r.Success {
		t.Fatal("expected strings not equal")
	}
}

func TestValidateString(t *testing.T) {
	r := ValidateString("hello", "hello", OpEqual)
	if !r.Success {
		t.Fatal("expected string equal")
	}
	r = ValidateString("hello", "world", OpContains)
	if r.Success {
		t.Fatal("expected hello does not contain world")
	}
	r = ValidateString("hello world", "hello", OpContains)
	if !r.Success {
		t.Fatal("expected hello world contains hello")
	}
	r = ValidateString("hello", "he", OpHasPrefix)
	if !r.Success {
		t.Fatal("expected hello has prefix he")
	}
	r = ValidateString("hello", "lo", OpHasSuffix)
	if !r.Success {
		t.Fatal("expected hello has suffix lo")
	}
}

func TestValidateContains(t *testing.T) {
	r := ValidateContains([]byte("hello world"), "world")
	if !r.Success {
		t.Fatal("expected body to contain substring")
	}
	r = ValidateContains([]byte("hello"), "world")
	if r.Success {
		t.Fatal("expected body to not contain substring")
	}
}

func TestValidateNotContains(t *testing.T) {
	r := ValidateNotContains([]byte("hello"), "world")
	if !r.Success {
		t.Fatal("expected body to not contain substring")
	}
	r = ValidateNotContains([]byte("hello world"), "world")
	if r.Success {
		t.Fatal("expected body to contain substring")
	}
}

func TestValidateStatusCode(t *testing.T) {
	r := ValidateStatusCode(200, 200, OpEqual)
	if !r.Success {
		t.Fatal("expected status code equal")
	}
	r = ValidateStatusCode(404, 200, OpNotEqual)
	if !r.Success {
		t.Fatal("expected status code not equal")
	}
}

func TestValidateStatusCodeRange(t *testing.T) {
	r := ValidateStatusCodeRange(200, 200, 299)
	if !r.Success {
		t.Fatal("expected 200 in [200,299]")
	}
	r = ValidateStatusCodeRange(199, 200, 299)
	if r.Success {
		t.Fatal("expected 199 not in [200,299]")
	}
	r = ValidateStatusCodeRange(300, 200, 299)
	if r.Success {
		t.Fatal("expected 300 not in [200,299]")
	}
}

func TestValidateHeader(t *testing.T) {
	headers := map[string]string{"Content-Type": "application/json"}
	r := ValidateHeader(headers, "Content-Type", "application/json", OpEqual)
	if !r.Success {
		t.Fatal("expected header equal")
	}
	r = ValidateHeader(headers, "Content-Type", "text/html", OpNotEqual)
	if !r.Success {
		t.Fatal("expected header not equal")
	}
	r = ValidateHeader(headers, "X-Missing", "something", OpEqual)
	if r.Success {
		t.Fatal("expected missing header to fail")
	}
}

func TestValidateContentType(t *testing.T) {
	headers := map[string]string{"Content-Type": "application/json; charset=utf-8"}
	r := ValidateContentType(headers, "application/json")
	if !r.Success {
		t.Fatal("expected content type to match")
	}
	r = ValidateContentType(headers, "text/html")
	if r.Success {
		t.Fatal("expected content type to not match")
	}
	r = ValidateContentType(map[string]string{}, "application/json")
	if r.Success {
		t.Fatal("expected missing content type to fail")
	}
}
