/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-06 16:08:21
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-06 17:36:05
 * @FilePath: \go-argus\rule\code.go
 * @Description: 业务错误码规格解析，把 code 标签与 validate 内联 code= 编译进规则计划
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"strconv"
	"strings"
)

// CodeSpec 字段级业务错误码规格，由 code 标签与 validate 内联 code= 解析得出
//
// 码在编译期写进 RulePlan.Code，运行期的 fieldError 只持一个 *RulePlan 指针
// （见 errors.go 的结构体注释），所以带码的错误与不带码的错误大小完全相同，
// 失败路径不多一次分配、GC 也不多一个待扫描指针字
type CodeSpec struct {
	Default int64            // 字段级默认码，标签里的裸数字项
	ByRule  map[string]int64 // 按规则名细分的码，标签里的 rule=code 项
}

// IsZero 规格是否为空，空规格不触碰任何规则计划
func (s CodeSpec) IsZero() bool {
	return s.Default == 0 && len(s.ByRule) == 0
}

// CodeFor 取规则名对应的码，按规则项优先于默认项；都没有返回 0
func (s CodeSpec) CodeFor(name string) int64 {
	if len(s.ByRule) > 0 {
		if c, ok := s.ByRule[name]; ok {
			return c
		}
	}
	return s.Default
}

// Merge 以 other 覆盖本规格并返回结果，用于「code 标签覆盖 validate 内联码」
//
// 覆盖只在同级之间发生：other 的默认项覆盖本默认项，other 的按规则项覆盖
// 同名按规则项。本规格的按规则项不会被 other 的默认项盖掉——按规则永远比
// 默认具体，无论它来自哪一处标签
func (s CodeSpec) Merge(other CodeSpec) CodeSpec {
	if other.Default != 0 {
		s.Default = other.Default
	}
	if len(other.ByRule) == 0 {
		return s
	}
	merged := make(map[string]int64, len(s.ByRule)+len(other.ByRule))
	for name, c := range s.ByRule {
		merged[name] = c
	}
	for name, c := range other.ByRule {
		merged[name] = c
	}
	s.ByRule = merged
	return s
}

// ParseCodeSpec 解析 code 标签
//
//	code:"110026"                                 字段级默认码
//	code:"required=110007,min=110008,max=110009"  按规则细分
//	code:"110000,min=110008"                      默认码 + 单规则覆盖
//
// 非十进制整数项与空项被静默忽略：标签是声明，解析期没有返回错误的出口，
// 宁可该规则不带码，也不能让一个写错的码把整个结构体的校验搞挂
func ParseCodeSpec(tag string) CodeSpec {
	var spec CodeSpec
	if strings.TrimSpace(tag) == "" {
		return spec
	}
	for _, part := range splitTag(tag) {
		spec.setEntry(part)
	}
	return spec
}

// setEntry 写入单个规格项，形如 110026 或 min=110008
func (s *CodeSpec) setEntry(entry string) {
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return
	}
	name, value, ok := strings.Cut(entry, "=")
	if !ok {
		if c, valid := parseCodeValue(entry); valid {
			s.Default = c
		}
		return
	}
	name = strings.TrimSpace(name)
	c, valid := parseCodeValue(value)
	if name == "" || !valid {
		return
	}
	if s.ByRule == nil {
		s.ByRule = make(map[string]int64, 4)
	}
	s.ByRule[name] = c
}

// parseCodeValue 解析十进制错误码，0 与非法值一律判为无效
func parseCodeValue(value string) (int64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	c, err := strconv.ParseInt(value, 10, 64)
	if err != nil || c == 0 {
		return 0, false
	}
	return c, true
}

// ApplyCodeSpec 把规格写入规则计划
//
// 控制规则（omitempty/dive 等）跳过：它们不产生错误，带码只是让计划数据多一份噪声
func ApplyCodeSpec(rules []RulePlan, spec CodeSpec) {
	if spec.IsZero() {
		return
	}
	for i := range rules {
		if rules[i].Kind&KindControl != 0 {
			continue
		}
		if c := spec.CodeFor(rules[i].Name); c != 0 {
			rules[i].Code = c
		}
	}
}
