/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 00:00:00
 * @FilePath: \go-argus\validate\compare_alias.go
 * @Description: Compare operator compatibility aliases
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import "github.com/kamalyes/go-argus/constants"

// CompareOperator 表示通用比较操作符
type CompareOperator = constants.CompareOperator

const (
	// OpEqual 表示相等
	OpEqual = constants.OpEqual
	// OpNotEqual 表示不等
	OpNotEqual = constants.OpNotEqual
	// OpGreaterThan 表示大于
	OpGreaterThan = constants.OpGreaterThan
	// OpGreaterThanOrEqual 表示大于等于
	OpGreaterThanOrEqual = constants.OpGreaterThanOrEqual
	// OpLessThan 表示小于
	OpLessThan = constants.OpLessThan
	// OpLessThanOrEqual 表示小于等于
	OpLessThanOrEqual = constants.OpLessThanOrEqual
	// OpContains 表示包含
	OpContains = constants.OpContains
	// OpNotContains 表示不包含
	OpNotContains = constants.OpNotContains
	// OpHasPrefix 表示前缀匹配
	OpHasPrefix = constants.OpHasPrefix
	// OpHasSuffix 表示后缀匹配
	OpHasSuffix = constants.OpHasSuffix
	// OpRegex 表示正则匹配
	OpRegex = constants.OpRegex
	// OpEmpty 表示空字符串
	OpEmpty = constants.OpEmpty
	// OpNotEmpty 表示非空字符串
	OpNotEmpty = constants.OpNotEmpty

	// OpSymbolEqual 是相等操作符别名
	OpSymbolEqual = constants.OpSymbolEqual
	// OpSymbolNotEqual 是不等操作符别名
	OpSymbolNotEqual = constants.OpSymbolNotEqual
	// OpSymbolGreaterThan 是大于操作符别名
	OpSymbolGreaterThan = constants.OpSymbolGreaterThan
	// OpSymbolGreaterThanOrEqual 是大于等于操作符别名
	OpSymbolGreaterThanOrEqual = constants.OpSymbolGreaterThanOrEqual
	// OpSymbolLessThan 是小于操作符别名
	OpSymbolLessThan = constants.OpSymbolLessThan
	// OpSymbolLessThanOrEqual 是小于等于操作符别名
	OpSymbolLessThanOrEqual = constants.OpSymbolLessThanOrEqual
)
