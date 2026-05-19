/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-18 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-18 00:00:00
 * @FilePath: \go-argus\utils\number.go
 * @Description: 数值解析工具，提供零依赖的浮点数解析能力
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package utils

import "fmt"

func ParseFloat(s string) (float64, bool) {
	n, err := ParseFloatStr(s)
	return n, err == nil
}

func ParseFloatStr(s string) (float64, error) {
	var n float64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			n = n*10 + float64(c-'0')
			continue
		}
		if c == '.' {
			frac := 0.0
			div := 1.0
			for j := i + 1; j < len(s); j++ {
				if s[j] < '0' || s[j] > '9' {
					return 0, fmt.Errorf("invalid float")
				}
				frac = frac*10 + float64(s[j]-'0')
				div *= 10
			}
			n += frac / div
			break
		}
		return 0, fmt.Errorf("invalid float")
	}
	return n, nil
}
