/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-21 02:26:21
 * @FilePath: \go-argus\validate\time_parse_test.go
 * @Description: time_parse.go 测试，覆盖星期和月份字段解析
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"testing"
	"time"
)

func TestParseWeek(t *testing.T) {
	tests := []struct {
		name    string
		week    string
		want    time.Weekday
		wantErr bool
	}{
		{"数字0周日", "0", time.Sunday, false},
		{"数字1周一", "1", time.Monday, false},
		{"数字6周六", "6", time.Saturday, false},
		{"数字7越界", "7", 0, true},
		{"缩写MON", "MON", time.Monday, false},
		{"全称MONDAY", "MONDAY", time.Monday, false},
		{"缩写FRI", "FRI", time.Friday, false},
		{"全称FRIDAY", "FRIDAY", time.Friday, false},
		{"缩写SUN", "SUN", time.Sunday, false},
		{"全称SUNDAY", "SUNDAY", time.Sunday, false},
		{"缩写WED", "WED", time.Wednesday, false},
		{"全称WEDNESDAY", "WEDNESDAY", time.Wednesday, false},
		{"无效字符串", "invalid", 0, true},
		{"小写monday", "monday", time.Monday, false},
		{"带空格", "  MON  ", time.Monday, false},
		{"单字母M", "M", time.Monday, false},
		{"单字母F", "F", time.Friday, false},
		{"空字符串", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseWeek(tt.week)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWeek(%q) error = %v, wantErr %v", tt.week, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseWeek(%q) = %v, want %v", tt.week, got, tt.want)
			}
		})
	}
}

func TestParseMonth(t *testing.T) {
	tests := []struct {
		name    string
		month   string
		want    time.Month
		wantErr bool
	}{
		{"数字1", "1", time.January, false},
		{"数字12", "12", time.December, false},
		{"数字13越界", "13", 0, true},
		{"数字0越界", "0", 0, true},
		{"缩写JAN", "JAN", time.January, false},
		{"全称JANUARY", "JANUARY", time.January, false},
		{"缩写DEC", "DEC", time.December, false},
		{"全称DECEMBER", "DECEMBER", time.December, false},
		{"缩写JUN", "JUN", time.June, false},
		{"全称JUNE", "JUNE", time.June, false},
		{"无效字符串", "invalid", 0, true},
		{"小写january", "january", time.January, false},
		{"带空格", "  JAN  ", time.January, false},
		{"单字母J", "J", time.January, false},
		{"空字符串", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMonth(tt.month)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMonth(%q) error = %v, wantErr %v", tt.month, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseMonth(%q) = %v, want %v", tt.month, got, tt.want)
			}
		})
	}
}
