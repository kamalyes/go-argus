/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 16:59:27
 * @FilePath: \go-argus\validate\compare.go
 * @Description: 通用比较校验能力，提供数值、字符串、HTTP 状态码和 Header 比较
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

// Package validate 承载从 go-toolbox/pkg/validator 迁移而来的通用校验能力
package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kamalyes/go-argus/constants"
	"github.com/kamalyes/go-argus/i18n"
)

// CompareResult 表示一次比较校验结果
type CompareResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Actual  string `json:"actual"`
	Expect  string `json:"expect"`
}

// Number 表示可比较数值类型集合
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// CompareNumbers 比较两个数值
func CompareNumbers[T Number](actual, expect T, op constants.CompareOperator) CompareResult {
	result := CompareResult{Actual: fmt.Sprint(actual), Expect: fmt.Sprint(expect)}
	switch op {
	case constants.OpEqual, constants.OpSymbolEqual:
		result.Success = actual == expect
	case constants.OpNotEqual, constants.OpSymbolNotEqual:
		result.Success = actual != expect
	case constants.OpGreaterThan, constants.OpSymbolGreaterThan:
		result.Success = actual > expect
	case constants.OpGreaterThanOrEqual, constants.OpSymbolGreaterThanOrEqual:
		result.Success = actual >= expect
	case constants.OpLessThan, constants.OpSymbolLessThan:
		result.Success = actual < expect
	case constants.OpLessThanOrEqual, constants.OpSymbolLessThanOrEqual:
		result.Success = actual <= expect
	default:
		result.Message = i18n.Msg(MsgCompareUnsupportedNumberOp)
	}
	if !result.Success && result.Message == "" {
		result.Message = i18n.Msg(MsgCompareNumberFailed, map[string]string{"actual": fmt.Sprint(actual), "op": op.String(), "expected": fmt.Sprint(expect)})
	}
	return result
}

// CompareStrings 比较两个字符串
func CompareStrings(actual, expect string, op constants.CompareOperator) CompareResult {
	result := CompareResult{Actual: actual, Expect: expect}
	switch op {
	case constants.OpEqual, constants.OpSymbolEqual:
		result.Success = actual == expect
	case constants.OpNotEqual, constants.OpSymbolNotEqual:
		result.Success = actual != expect
	case constants.OpContains:
		result.Success = strings.Contains(actual, expect)
	case constants.OpNotContains:
		result.Success = !strings.Contains(actual, expect)
	case constants.OpHasPrefix:
		result.Success = strings.HasPrefix(actual, expect)
	case constants.OpHasSuffix:
		result.Success = strings.HasSuffix(actual, expect)
	case constants.OpEmpty:
		result.Success = strings.TrimSpace(actual) == ""
		result.Expect = "empty string"
	case constants.OpNotEmpty:
		result.Success = strings.TrimSpace(actual) != ""
		result.Expect = "non-empty string"
	case constants.OpRegex:
		re, err := regexp.Compile(expect)
		if err != nil {
			result.Message = i18n.Msg(MsgCompareRegexCompileFailed, map[string]string{"error": err.Error()})
			return result
		}
		result.Success = re.MatchString(actual)
	default:
		result.Message = i18n.Msg(MsgCompareUnsupportedStringOp)
	}
	if !result.Success && result.Message == "" {
		result.Message = i18n.Msg(MsgCompareStringFailed, map[string]string{"actual": actual, "op": op.String(), "expected": expect})
	}
	return result
}

// ValidateString 校验字符串关系，保留 go-toolbox 旧函数名
func ValidateString(actual, expect string, op constants.CompareOperator) CompareResult {
	return CompareStrings(actual, expect, op)
}

// ValidateContains 校验字节内容是否包含子串
func ValidateContains(body []byte, substring string) CompareResult {
	return CompareStrings(string(body), substring, constants.OpContains)
}

// ValidateNotContains 校验字节内容是否不包含子串
func ValidateNotContains(body []byte, substring string) CompareResult {
	return CompareStrings(string(body), substring, constants.OpNotContains)
}

// ValidateStatusCode 比较 HTTP 状态码
func ValidateStatusCode(statusCode, expected int, op constants.CompareOperator) CompareResult {
	return CompareNumbers(statusCode, expected, op)
}

// ValidateStatusCodeRange 校验 HTTP 状态码是否在闭区间内
func ValidateStatusCodeRange(actual, min, max int) CompareResult {
	result := CompareResult{
		Actual: fmt.Sprint(actual),
		Expect: fmt.Sprintf("%d-%d", min, max),
	}
	result.Success = actual >= min && actual <= max
	if !result.Success {
		result.Message = i18n.Msg(MsgCompareStatusOutOfRange, map[string]string{"actual": fmt.Sprint(actual), "min": fmt.Sprint(min), "max": fmt.Sprint(max)})
	}
	return result
}

// ValidateHeader 根据操作符比较 Header 值
func ValidateHeader(headers map[string]string, key, expected string, op constants.CompareOperator) CompareResult {
	return CompareStrings(headers[key], expected, op)
}

// ValidateContentType 校验 Content-Type 是否包含期望类型
func ValidateContentType(headers map[string]string, expected string) CompareResult {
	actual := headers["Content-Type"]
	if actual == "" {
		actual = headers["content-type"]
	}
	return CompareStrings(actual, expected, constants.OpContains)
}
