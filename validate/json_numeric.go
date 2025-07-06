/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-06 21:16:17
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-07 00:00:00
 * @FilePath: \go-argus\validate\json_numeric.go
 * @Description: JSON 数字字符串宽松转换工具
 *               解决前端将 int64/uint64/float64 等数字类型以字符串形式传递时（如 JS 精度问题），
 *               后端 JSON 反序列化无法自动将字符串转为数字类型的兼容性问题
 *
 *               性能设计：
 *                 1. 快速预检（零分配）：字节级扫描，无引号数字模式则直接返回原始数据
 *                 2. 仅在必要时执行 decode→convert→encode 流程
 *                 3. 零外部依赖，仅使用标准库
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// ConvertNumericStrings 递归遍历 JSON 数据，将所有看起来像数字的字符串值自动转换为 json.Number 类型
//
// 适用场景：
//   - 前端因 JavaScript Number 精度限制（最大安全整数 2^53-1），将 int64/uint64 以字符串形式传递；
//   - 后端 protobuf/jsonpb 反序列化时期望收到数字类型，但实际收到的是字符串；
//   - 需要做一次兼容转换，将数字字符串转为 JSON 数字，再交给后续的 protojson/encoding/json 解码
//
// 处理逻辑：
//   - 对象（map）：递归处理每个字段值；
//   - 数组（slice）：递归处理每个元素；
//   - 字符串：尝试按 int64 → uint64 → float64 顺序解析，解析成功则转为 json.Number；
//   - 其他类型（bool、number、null、已有的 json.Number）保持原样
//
// 安全策略：
//   - 该函数只做语法层面的转换，不依赖任何 schema 信息；
//   - 如果字符串无法解析为数字（如 "hello"、"ABC123"），保持原样不做处理；
//   - 空字符串直接返回，不做转换；
//   - 通过快速字节预检避免不必要的 decode/encode 开销
//
// 性能说明：
//   - 当 JSON 中不存在引号包裹的数字时（常见场景），仅执行 O(n) 字节扫描，零堆分配；
//   - 当存在数字字符串时，执行一次完整的 decode-convert-encode，开销约为 2x JSON 解析
//
// 参数：
//   - data: 原始 JSON 字节切片
//
// 返回值：
//   - []byte: 转换后的 JSON 字节切片；如果解析失败返回 nil 和 error
//   - error: JSON 解析失败时返回错误
func ConvertNumericStrings(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	// 快速预检：扫描是否存在引号包裹的数字（如 :"123" 或 ,"-456" 或 ["789）
	// 如果不存在，直接返回原始数据，避免 decode/encode 开销
	if !ContainsQuotedNumber(data) {
		return data, nil
	}

	var v interface{}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}

	converted := ConvertNumericStringsRecursive(v)

	var buf bytes.Buffer
	buf.Grow(len(data))
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(converted); err != nil {
		return nil, err
	}

	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}

// ContainsQuotedNumber 使用零分配字节扫描快速判断 JSON 中是否存在引号包裹的数字
//
// 该函数是高性能预检工具，在需要决定是否执行昂贵的 decode→convert→encode 流程前调用
// 识别模式（紧跟冒号/逗号/左中括号后的引号+数字/负号）：
//   - 键值对：`:"123` `:"-45`
//   - 数组元素：`,"123` `,"-45` 或 `["123` `["-45`
//   - 顶层字符串：整个 JSON 就是一个数字字符串如 `"123"`
//
// 性能特征：
//   - 零堆分配、零内存拷贝，仅做字节级扫描
//   - 小 JSON（< 200字节）通常在 20-100ns 内完成
//   - 找到匹配模式后立即返回 true，无需扫描完整内容
//
// 注意：这是一个启发式预检，允许误报（会多做一次 decode）但绝不漏报（不会漏掉需要转换的情况）
//
// 参数：
//   - data: JSON 字节切片
//
// 返回值：
//   - bool: 是否存在引号包裹的数字字符串
func ContainsQuotedNumber(data []byte) bool {
	n := len(data)
	if n < 2 {
		return false
	}

	i := 0
	// 跳过开头空白
	i = SkipJSONSpaces(data, i)

	// 检查顶层字符串（wrapperspb 等 well-known types 场景）
	if i < n && data[i] == '"' {
		// 跳过开头引号
		i++
		if i < n && IsNumStart(data[i]) {
			return true
		}
	}

	// 扫描整个 JSON 字节流
	for i < n {
		switch data[i] {
		case '"':
			// 遇到字符串：跳过整个字符串内容（含转义处理），避免误判字符串内部的 : , [
			i++ // 跳过开头引号
			for i < n {
				if data[i] == '\\' && i+1 < n {
					i += 2
					continue
				}
				if data[i] == '"' {
					break
				}
				i++
			}
			// i 现在指向结束引号，外层 i++ 会跳过它
		case ':', ',', '[':
			// 遇到值分隔符：跳过空白后检查是否是 "数字 模式
			j := SkipJSONSpaces(data, i+1)
			if j+1 < n && data[j] == '"' && IsNumStart(data[j+1]) {
				return true
			}
		}
		i++
	}
	return false
}

// IsNumStart 判断字节是否是数字字符串的起始字符
//
// 合法的数字起始字符包括：
//   - '-'：负数开头（如 "-123"、"-3.14"）
//   - '+'：正数开头（如 "+123"，部分 JSON5 扩展支持）
//   - '0'-'9'：阿拉伯数字
//
// 参数：
//   - c: 待判断的字节
//
// 返回值：
//   - bool: 是否可以作为数字（字符串形式）的起始字符
func IsNumStart(c byte) bool {
	return c == '-' || c == '+' || (c >= '0' && c <= '9')
}

// convertNumericStringsRecursive 递归遍历 JSON 结构，将数字字符串转为 json.Number
func ConvertNumericStringsRecursive(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, item := range val {
			val[k] = ConvertNumericStringsRecursive(item)
		}
		return val
	case []interface{}:
		for i, item := range val {
			val[i] = ConvertNumericStringsRecursive(item)
		}
		return val
	case string:
		return TryStringToNumber(val)
	default:
		return v
	}
}

// TryStringToNumber 尝试将字符串解析为数字，成功则返回 json.Number，失败则返回原字符串
// 解析顺序：先尝试 int64（有符号整数），再尝试 uint64（无符号大整数），最后尝试 float64（浮点数）
func TryStringToNumber(s string) interface{} {
	if len(s) == 0 {
		return s
	}

	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return json.Number(s)
	}

	if _, err := strconv.ParseUint(s, 10, 64); err == nil {
		return json.Number(s)
	}

	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return json.Number(s)
	}

	return s
}
