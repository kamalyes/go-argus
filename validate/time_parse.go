/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-16 23:18:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-25 23:25:06
 * @FilePath: \go-argus\validate\time_parse.go
 * @Description: 时间字段解析（星期、月份）
 *
 * 支持多种格式的星期和月份字符串解析，返回标准 time 类型
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseWeek 解析星期字段，返回 time.Weekday
//
// 支持格式：
//   - 数字：0-6（0=Sunday）
//   - 缩写：M/MON, T/TUE, W/WED, R/THU, F/FRI, S/SAT, U/SUN
//   - 全称：MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, SUNDAY
func ParseWeek(week string) (time.Weekday, error) {
	// 尝试数字解析
	if val, err := strconv.Atoi(week); err == nil && val >= 0 && val <= 6 {
		return time.Weekday(val), nil
	}

	// 转大写并去除空格
	upper := strings.ToUpper(strings.TrimSpace(week))

	switch upper {
	case "M", "MON", "MONDAY":
		return time.Monday, nil
	case "T", "TUE", "TUESDAY":
		return time.Tuesday, nil
	case "W", "WED", "WEDNESDAY":
		return time.Wednesday, nil
	case "R", "THU", "THURSDAY":
		return time.Thursday, nil
	case "F", "FRI", "FRIDAY":
		return time.Friday, nil
	case "S", "SAT", "SATURDAY":
		return time.Saturday, nil
	case "U", "SUN", "SUNDAY":
		return time.Sunday, nil
	}

	return 0, fmt.Errorf("invalid week: %s", week)
}

// ParseMonth 解析月份字段，返回 time.Month
//
// 支持格式：
//   - 数字：1-12
//   - 缩写：JAN, FEB, MAR, APR, MAY, JUN, JUL, AUG, SEP, OCT, NOV, DEC
//   - 全称：JANUARY, FEBRUARY, MARCH, APRIL, MAY, JUNE, JULY, AUGUST, SEPTEMBER, OCTOBER, NOVEMBER, DECEMBER
func ParseMonth(month string) (time.Month, error) {
	upper := strings.ToUpper(strings.TrimSpace(month))

	// 尝试数字解析
	if val, err := strconv.Atoi(upper); err == nil {
		if val >= 1 && val <= 12 {
			return time.Month(val), nil
		}
		return 0, fmt.Errorf("invalid month: %s", month)
	}

	switch upper {
	case "JAN", "JANUARY", "J":
		return time.January, nil
	case "FEB", "FEBRUARY", "F":
		return time.February, nil
	case "MAR", "MARCH", "M":
		return time.March, nil
	case "APR", "APRIL", "A":
		return time.April, nil
	case "MAY", "Y":
		return time.May, nil
	case "JUN", "JUNE", "N":
		return time.June, nil
	case "JUL", "JULY", "L":
		return time.July, nil
	case "AUG", "AUGUST", "G":
		return time.August, nil
	case "SEP", "SEPTEMBER", "S":
		return time.September, nil
	case "OCT", "OCTOBER", "T":
		return time.October, nil
	case "NOV", "NOVEMBER", "V":
		return time.November, nil
	case "DEC", "DECEMBER", "C":
		return time.December, nil
	}

	return 0, fmt.Errorf("invalid month: %s", month)
}
