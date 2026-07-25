/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-16 23:18:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-25 23:19:56
 * @FilePath: \go-argus\validate\chinese.go
 * @Description: 中国大陆特定格式校验（手机号、身份证号）
 *
 * 设计原则：
 *   1. 手动 byte 遍历，零正则零分配
 *   2. 身份证校验和用 ASCII 算术（char - '0'）计算
 *   3. 早退：长度/前缀不符立即返回
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

// IsChinesePhoneNumber 检查字符串是否为有效的大陆手机号
// 规则：长度11位 + 首位为1 + 第二位为3-9 + 其余为数字
func IsChinesePhoneNumber(s string) bool {
	if len(s) != 11 {
		return false
	}
	// 首位必须为 1
	if s[0] != '1' {
		return false
	}
	// 第二位必须为 3-9
	c := s[1]
	if c < '3' || c > '9' {
		return false
	}
	// 其余9位必须为数字
	return allDigits(s[2:])
}

// IsChineseIDCard 检查字符串是否为有效的大陆身份证号（15位或18位）
// 18位身份证包含校验位，15位身份证不包含校验位
func IsChineseIDCard(s string) bool {
	switch len(s) {
	case 15:
		// 15位身份证：全数字
		return allDigits(s)
	case 18:
		// 18位身份证：前17位数字 + 最后一位数字或X
		if !allDigits(s[:17]) {
			return false
		}
		last := s[17]
		return (last >= '0' && last <= '9') || last == 'X' || last == 'x'
	default:
		return false
	}
}

// IsChineseIDCardWithChecksum 检查18位身份证号（含校验位验证）
// 比仅检查格式更严格：会计算校验和并验证最后一位
func IsChineseIDCardWithChecksum(s string) bool {
	if len(s) != 18 {
		return false
	}
	if !allDigits(s[:17]) {
		return false
	}
	// 最后一位可以是数字或X
	last := s[17]
	if !((last >= '0' && last <= '9') || last == 'X' || last == 'x') {
		return false
	}
	// 计算校验和
	expected := CalculateIDCardChecksum(s[:17])
	// 比较（X 大小写不敏感）
	return last == expected[0] || (last == 'x' && expected[0] == 'X')
}

// CalculateIDCardChecksum 计算给定17位身份证号的校验和字符
//
// 算法：
//  1. 前17位数字分别乘以权重系数 [7,9,10,5,8,4,2,1,6,3,7,9,10,5,8,4,2]
//  2. 求和后对11取模
//  3. 根据模值查表得到校验码：[1,0,X,9,8,7,6,5,4,3,2]
//
// 性能：使用 ASCII 算术（char - '0'），按 byte 遍历，零分配
func CalculateIDCardChecksum(id string) string {
	weights := [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 索引 sum%11 对应的校验码
	checkMap := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

	sum := 0
	// 按 byte 遍历，避免 rune 切片分配；身份证号均为 ASCII 数字
	for i := 0; i < len(id) && i < len(weights); i++ {
		sum += int(id[i]-'0') * weights[i]
	}

	return string(checkMap[sum%11])
}
