/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 00:00:00
 * @FilePath: \go-argus\constants\compare.go
 * @Description: Central comparison operator names
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

// CompareOperator 操作符类型
type CompareOperator string

// String 返回操作符字符串表示
func (op CompareOperator) String() string {
	return string(op)
}

const (
	OpEqual                    CompareOperator = RuleEq
	OpNotEqual                 CompareOperator = RuleNe
	OpGreaterThan              CompareOperator = RuleGT
	OpGreaterThanOrEqual       CompareOperator = RuleGTE
	OpLessThan                 CompareOperator = RuleLT
	OpLessThanOrEqual          CompareOperator = RuleLTE
	OpContains                 CompareOperator = "contains"
	OpNotContains              CompareOperator = "not_contains"
	OpHasPrefix                CompareOperator = "has_prefix"
	OpHasSuffix                CompareOperator = "has_suffix"
	OpRegex                    CompareOperator = RuleRegex
	OpEmpty                    CompareOperator = "empty"
	OpNotEmpty                 CompareOperator = RuleNotEmpty
	OpSymbolEqual              CompareOperator = RuleSymbolEqual
	OpSymbolNotEqual           CompareOperator = RuleSymbolNotEqual
	OpSymbolGreaterThan        CompareOperator = RuleSymbolGreaterThan
	OpSymbolGreaterThanOrEqual CompareOperator = RuleSymbolGreaterThanOrEqual
	OpSymbolLessThan           CompareOperator = RuleSymbolLessThan
	OpSymbolLessThanOrEqual    CompareOperator = RuleSymbolLessThanOrEqual
)
