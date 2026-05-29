/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\json.go
 * @Description: JSON 校验能力，提供 JSON 有效性、字段读取和轻量路径匹配
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kamalyes/go-argus/constants"
	"github.com/kamalyes/go-argus/i18n"
)

// ValidateJSON 校验 JSON 字节是否有效
func ValidateJSON(data []byte) error {
	if json.Valid(data) {
		return nil
	}
	return fmt.Errorf(i18n.Msg(MsgJSONInvalid))
}

func IsValidJSONBytes(data []byte) bool {
	return json.Valid(data)
}

// IsJSONNull 判断 JSON 字节是否为 null
func IsJSONNull(data []byte) bool {
	return strings.TrimSpace(string(data)) == "null"
}

// IsJSONColumnType 判断数据库列类型是否属于 JSON 类型
func IsJSONColumnType(dbType string) bool {
	dbType = strings.TrimSpace(strings.ToLower(dbType))
	if idx := strings.IndexByte(dbType, '('); idx >= 0 {
		dbType = strings.TrimSpace(dbType[:idx])
	}
	return dbType == "json" || dbType == "jsonb"
}

// ValidateJSONWithData 校验 JSON 并返回反序列化数据
func ValidateJSONWithData(body []byte) (any, error) {
	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// ValidateJSONField 校验 JSON 顶层字段是否等于期望值
func ValidateJSONField(body []byte, field string, expected any) CompareResult {
	result := CompareResult{Expect: fmt.Sprint(expected)}
	data, err := ValidateJSONWithData(body)
	if err != nil {
		result.Message = err.Error()
		return result
	}
	obj, ok := data.(map[string]any)
	if !ok {
		result.Message = i18n.Msg(MsgJSONRootNotObject)
		return result
	}
	actual, ok := obj[field]
	if !ok {
		result.Message = i18n.Msg(MsgJSONFieldNotFound, map[string]string{"field": field})
		return result
	}
	result.Actual = fmt.Sprint(actual)
	result.Success = fmt.Sprint(actual) == fmt.Sprint(expected)
	if !result.Success {
		result.Message = i18n.Msg(MsgJSONFieldValueMismatch)
	}
	return result
}

// ValidateJSONFields 批量校验 JSON 顶层字段
func ValidateJSONFields(body []byte, rules map[string]any) []CompareResult {
	out := make([]CompareResult, 0, len(rules))
	for field, expected := range rules {
		out = append(out, ValidateJSONField(body, field, expected))
	}
	return out
}

// LookupJSONPath 按轻量路径读取 JSON 数据，支持 $.a.b、a.b、items[0].name
func LookupJSONPath(data any, path string) (any, bool) {
	path = strings.TrimSpace(path)
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")
	if path == "" {
		return data, true
	}
	current := data
	for _, part := range splitJSONPath(path) {
		name, indexes := parsePathPart(part)
		if name != "" {
			obj, ok := current.(map[string]any)
			if !ok {
				return nil, false
			}
			current, ok = obj[name]
			if !ok {
				return nil, false
			}
		}
		for _, idx := range indexes {
			arr, ok := current.([]any)
			if !ok || idx < 0 || idx >= len(arr) {
				return nil, false
			}
			current = arr[idx]
		}
	}
	return current, true
}

// ValidateJSONPath 校验 JSON 路径的值
func ValidateJSONPath(body []byte, jsonPath string, expected any, op constants.CompareOperator) CompareResult {
	result := CompareResult{Expect: fmt.Sprint(expected)}
	data, err := ValidateJSONWithData(body)
	if err != nil {
		result.Message = err.Error()
		return result
	}
	actual, ok := LookupJSONPath(data, jsonPath)
	if !ok {
		result.Message = i18n.Msg(MsgJSONPathNotFound, map[string]string{"path": jsonPath})
		return result
	}
	result.Actual = fmt.Sprint(actual)
	return CompareStrings(result.Actual, fmt.Sprint(expected), op)
}

// ValidateJSONPathExists 校验 JSON 路径是否存在
func ValidateJSONPathExists(body []byte, jsonPath string) CompareResult {
	result := CompareResult{Expect: "path exists"}
	data, err := ValidateJSONWithData(body)
	if err != nil {
		result.Message = err.Error()
		return result
	}
	_, ok := LookupJSONPath(data, jsonPath)
	result.Success = ok
	if !ok {
		result.Message = i18n.Msg(MsgJSONPathNotFound, map[string]string{"path": jsonPath})
	}
	return result
}

// SkipJSONSpaces 跳过 JSON 字节流中的空白字符，并返回下一个非空白位置
func SkipJSONSpaces(data []byte, i int) int {
	for i < len(data) {
		switch data[i] {
		case ' ', '\n', '\r', '\t':
			i++
		default:
			return i
		}
	}
	return i
}

// ScanJSONString 扫描 JSON 字符串，并返回字符串结束后一位的位置
func ScanJSONString(data []byte, start int) (int, error) {
	if start >= len(data) || data[start] != '"' {
		return 0, fmt.Errorf(i18n.Msg(MsgJSONInvalid))
	}
	for i := start + 1; i < len(data); i++ {
		switch data[i] {
		case '\\':
			i++
		case '"':
			return i + 1, nil
		}
	}
	return 0, fmt.Errorf(i18n.Msg(MsgJSONInvalid))
}

// ScanJSONValueEnd 扫描任意 JSON 值，并返回值结束后一位的位置
func ScanJSONValueEnd(data []byte, start int) (int, error) {
	if start >= len(data) {
		return 0, fmt.Errorf(i18n.Msg(MsgJSONInvalid))
	}
	switch data[start] {
	case '"':
		return ScanJSONString(data, start)
	case '{', '[':
		return scanJSONCompositeEnd(data, start)
	default:
		return scanJSONScalarEnd(data, start)
	}
}

func scanJSONCompositeEnd(data []byte, start int) (int, error) {
	stack := make([]byte, 0, 4)
	for i := start; i < len(data); i++ {
		switch data[i] {
		case '"':
			end, err := ScanJSONString(data, i)
			if err != nil {
				return 0, err
			}
			i = end - 1
		case '{':
			stack = append(stack, '}')
		case '[':
			stack = append(stack, ']')
		case '}', ']':
			last := len(stack) - 1
			if last < 0 || stack[last] != data[i] {
				return 0, fmt.Errorf(i18n.Msg(MsgJSONInvalid))
			}
			stack = stack[:last]
			if len(stack) == 0 {
				return i + 1, nil
			}
		}
	}
	return 0, fmt.Errorf(i18n.Msg(MsgJSONInvalid))
}

func scanJSONScalarEnd(data []byte, start int) (int, error) {
	for i := start; i < len(data); i++ {
		switch data[i] {
		case ',', '}', ']':
			return i, nil
		case ' ', '\n', '\r', '\t':
			end := i
			for i < len(data) {
				switch data[i] {
				case ' ', '\n', '\r', '\t':
					i++
				case ',', '}', ']':
					return end, nil
				default:
					return i, nil
				}
			}
			return end, nil
		}
	}
	return len(data), nil
}

func splitJSONPath(path string) []string {
	parts := strings.Split(path, ".")
	out := parts[:0]
	for _, part := range parts {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func parsePathPart(part string) (string, []int) {
	name := part
	var indexes []int
	if i := strings.Index(part, "["); i >= 0 {
		name = part[:i]
		rest := part[i:]
		for strings.HasPrefix(rest, "[") {
			end := strings.Index(rest, "]")
			if end <= 1 {
				break
			}
			idx, err := strconv.Atoi(rest[1:end])
			if err != nil {
				break
			}
			indexes = append(indexes, idx)
			rest = rest[end+1:]
		}
	}
	return name, indexes
}
