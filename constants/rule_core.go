/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-29 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 00:00:00
 * @FilePath: \go-argus\constants\rule_core.go
 * @Description: 核心与结构体控制规则名常量
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

const (
	RuleEmpty = "" // 空规则

	RuleRequired  = "required"  // 必填
	RuleIsDefault = "isdefault" // 默认值

	RuleOmitEmpty = "omitempty" // 空值省略
	RuleOmitZero  = "omitzero"  // 零值省略
	RuleOmitNil   = "omitnil"   // nil 省略

	RuleDive          = "dive"          // 递归进入切片/映射
	RuleKeys          = "keys"          // 映射键规则开始
	RuleEndKeys       = "endkeys"       // 映射键规则结束
	RuleStructOnly    = "structonly"    // 仅校验结构体本身
	RuleNoStructLevel = "nostructlevel" // 跳过结构体级校验
)
