/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-25 23:16:16
 * @FilePath: \go-argus\validate\kind.go
 * @Description: 整数校验（与 RuleNumber 互补，不重复）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"math"
	"reflect"
)

// IsWholeNumber 判断 float64 是否为整数（whole number）
//
// 用于 RuleInteger 规则：float 类型字段必须为整数值
//
// 实现说明：
//   - 使用 math.Trunc 而非 int64 转换，避免大数溢出
//     int64 转换在 |f| > math.MaxInt64 时会返回错误结果（溢出）
//     math.Trunc 直接操作 float64 的尾数位，无溢出风险，支持全 float64 范围
//   - 显式排除 NaN 和 Inf
//     NaN：math.Trunc(NaN) = NaN，但 NaN != NaN，所以会被误判为非整数（巧合正确）
//     Inf：math.Trunc(+Inf) = +Inf，+Inf == +Inf 为 true，需显式排除
//
// 性能：2 次比较（IsNaN/IsInf）+ 1 次 Trunc + 1 次比较，无分配
// 相比 int64 转换版本：性能相当，但精度完整，无溢出风险
func IsWholeNumber(f float64) bool {
	// 显式排除 NaN 和 Inf（标准库实现为位运算，零分配，纳秒级）
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return false
	}
	// math.Trunc 截断小数部分，保留整数部分
	// 直接操作 float64 尾数位，不经过 int64，无溢出风险
	return math.Trunc(f) == f
}

// GetReflectKind 获取任意值的 reflect.Kind
func GetReflectKind(value interface{}) reflect.Kind {
	if value == nil {
		return reflect.Invalid
	}
	return reflect.ValueOf(value).Kind()
}

// StringInteger 判断字符串是否表示一个整数
//
// 支持格式：
//   - "123"    → true
//   - "-123"   → true
//   - "+123"   → true
//   - "0"      → true
//   - "12.3"   → false（含小数点）
//   - "12a"    → false（含非数字字符）
//   - ""       → false（空字符串）
//   - "+"      → false（只有符号）
//   - " 123"   → false（含空格，需调用方先 TrimSpace）
//
// 性能：手动字节扫描，零分配
// 相比 strconv.ParseFloat(s, 64) 节省：
//   - 1 次 float64 解析
//   - 1 次 error 对象分配
//   - 1 次 strings.Contains 等辅助调用
func StringInteger(s string) bool {
	n := len(s)
	if n == 0 {
		return false
	}

	i := 0
	// 处理正负号
	if s[0] == '+' || s[0] == '-' {
		i = 1
		if n == 1 {
			return false // 只有符号，不是有效整数
		}
	}

	// 扫描数字字符（手动字节遍历，零分配）
	for ; i < n; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}

	return true
}
