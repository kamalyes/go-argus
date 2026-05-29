/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 18:01:15
 * @FilePath: \go-argus\constants\cmp.go
 * @Description: 内部快速路径比较运算符，基于 int 的 switch 分发优于字符串
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

// CmpOp 整型比较运算符，用于内部快速 switch 分发
type CmpOp int

const (
	CmpEQ  CmpOp = iota // 等于
	CmpLT               // 小于
	CmpLTE              // 小于等于
	CmpGT               // 大于
	CmpGTE              // 大于等于
	CmpNE               // 不等于
)

// CmpOpFromStr 将字符串操作符转换为 CmpOp，未匹配返回 -1
func CmpOpFromStr(op string) CmpOp {
	switch op {
	case RuleEq:
		return CmpEQ
	case RuleNe:
		return CmpNE
	case RuleGT:
		return CmpGT
	case RuleGTE:
		return CmpGTE
	case RuleLT:
		return CmpLT
	case RuleLTE:
		return CmpLTE
	default:
		return -1
	}
}

// CmpOpForOperator 将 CompareOperator 转换为 CmpOp
func CmpOpForOperator(op CompareOperator) CmpOp {
	return CmpOpFromStr(op.String())
}
