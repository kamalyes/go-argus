/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-26 09:11:35
 * @FilePath: \go-argus\rule\time.go
 * @Description: 时间规则模块，提供 time.Time、protobuf Timestamp 反射识别和 now 表达式解析
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package rule

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/kamalyes/go-argus/constants"
	"github.com/kamalyes/go-argus/validate"
)

// TimeValue 将字段值识别为 time.Time，核心包不直接依赖 protobuf
func TimeValue(value reflect.Value, layout string) (time.Time, bool) {
	value = validate.DerefReflect(value)
	if !value.IsValid() {
		return time.Time{}, false
	}
	// 按 kind 分流，避免对不可能承载时间值的类型付装箱与反射探测代价。
	//
	// 原实现无条件走 value.Interface() 类型断言，再依次探测 callAsTime /
	// callSecondsNanos / readSecondsNanos。而 CompareValueOp（rule/field.go）在
	// NumericValue 之前先试 TimeValue，于是每次整数跨字段比较都要：
	// 装箱一次（逃逸到堆）+ 两次 MethodByName 反射探测（各带 defer/recover 闭包），
	// 全部注定失败——这是"校验通过"路径上仅剩的分配来源之一
	switch value.Kind() {
	case reflect.String:
		// 字符串直接取值，无需 Interface() 装箱
		s := value.String()
		if layout == "" && !looksLikeTimeString(s) {
			return time.Time{}, false
		}
		return parseTimeString(s, layout)

	case reflect.Struct, reflect.Interface:
		// time.Time 本身与 protobuf Timestamp 等载体均为结构体/接口，
		// 仅在此分支付装箱代价，随后按原顺序走三个 fallback
		if value.CanInterface() {
			if t, ok := value.Interface().(time.Time); ok {
				return t, true
			}
		}

	default:
		// 数值/布尔/切片/映射等不可能是时间值：
		// readSecondsNanos 本就要求 Kind()==Struct，callAsTime/callSecondsNanos
		// 依赖 AsTime()/Seconds+Nanos 载体（实践中均为结构体），故直接短路
		// 理论差异：带 AsTime() 方法的具名非结构体类型不再被识别——
		// 现实中无此形态（protobuf Timestamp 是结构体指针），换取的是所有数值比较路径零分配
		return time.Time{}, false
	}

	if t, ok := callAsTime(value); ok {
		return t, true
	}
	if t, ok := callSecondsNanos(value); ok {
		return t, true
	}
	if t, ok := readSecondsNanos(value); ok {
		return t, true
	}
	return time.Time{}, false
}

func looksLikeTimeString(value string) bool {
	value = strings.TrimSpace(value)
	return len(value) >= len("2006-01-02") &&
		isDigit(value[0]) && isDigit(value[1]) && isDigit(value[2]) && isDigit(value[3]) &&
		value[4] == '-' &&
		isDigit(value[5]) && isDigit(value[6]) &&
		value[7] == '-' &&
		isDigit(value[8]) && isDigit(value[9])
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// ResolveTimeExpr 解析 now、now+5m、now-30d 这类时间表达式
func ResolveTimeExpr(expr string, now time.Time) (time.Time, bool) {
	expr = strings.TrimSpace(expr)
	if expr == "" || expr == "now" {
		return now, expr == "now"
	}
	if !strings.HasPrefix(expr, "now+") && !strings.HasPrefix(expr, "now-") {
		return time.Time{}, false
	}
	sign := expr[3]
	d, ok := parseDuration(expr[4:])
	if !ok {
		return time.Time{}, false
	}
	if sign == '-' {
		d = -d
	}
	return now.Add(d), true
}

// CompareTimeExpr 将字段时间与表达式时间比较
func CompareTimeExpr(field reflect.Value, expr, op string, now time.Time) bool {
	left, ok := TimeValue(field, "")
	if !ok {
		return false
	}
	right, ok := ResolveTimeExpr(expr, now)
	return ok && validate.CompareOp(float64(left.UnixNano()), float64(right.UnixNano()), constants.CmpOpFromStr(op))
}

func parseTimeString(value, layout string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	if layout != "" {
		t, err := time.Parse(layout, value)
		return t, err == nil
	}
	for _, item := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(item, value); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func callAsTime(value reflect.Value) (_ time.Time, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	method := value.MethodByName("AsTime")
	if !method.IsValid() || method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
		return time.Time{}, false
	}
	out := method.Call(nil)
	t, ok := out[0].Interface().(time.Time)
	return t, ok
}

func callSecondsNanos(value reflect.Value) (_ time.Time, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	secondsMethod := value.MethodByName("GetSeconds")
	if !secondsMethod.IsValid() || secondsMethod.Type().NumIn() != 0 || secondsMethod.Type().NumOut() != 1 {
		return time.Time{}, false
	}
	secondsOut := secondsMethod.Call(nil)
	if secondsOut[0].Kind() != reflect.Int64 && secondsOut[0].Kind() != reflect.Int {
		return time.Time{}, false
	}
	seconds := secondsOut[0].Int()
	nanos := int64(0)
	nanosMethod := value.MethodByName("GetNanos")
	if nanosMethod.IsValid() && nanosMethod.Type().NumIn() == 0 && nanosMethod.Type().NumOut() == 1 {
		nanosOut := nanosMethod.Call(nil)
		if nanosOut[0].Kind() == reflect.Int64 || nanosOut[0].Kind() == reflect.Int {
			nanos = nanosOut[0].Int()
		}
	}
	return time.Unix(seconds, nanos), true
}

func readSecondsNanos(value reflect.Value) (time.Time, bool) {
	value = validate.DerefReflect(value)
	if !value.IsValid() || value.Kind() != reflect.Struct {
		return time.Time{}, false
	}
	secondsField := value.FieldByName("Seconds")
	if !secondsField.IsValid() {
		return time.Time{}, false
	}
	seconds := secondsField.Int()
	nanos := int64(0)
	if nanosField := value.FieldByName("Nanos"); nanosField.IsValid() {
		nanos = nanosField.Int()
	}
	return time.Unix(seconds, nanos), true
}

func parseDuration(value string) (time.Duration, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	if before, ok := strings.CutSuffix(value, "d"); ok {
		n, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return 0, false
		}
		return time.Duration(n) * 24 * time.Hour, true
	}
	d, err := time.ParseDuration(value)
	return d, err == nil
}
