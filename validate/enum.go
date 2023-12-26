/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\validate\enum.go
 * @Description: 泛型枚举校验器，用于业务状态、角色、类型等枚举值校验
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validate

import (
	"fmt"

	"github.com/kamalyes/go-argus/i18n"
)

// EnumValidator 表示一组可校验的枚举值
type EnumValidator[T comparable] struct {
	values map[T]struct{}
	order  []T
}

// NewEnumValidator 创建枚举校验器
func NewEnumValidator[T comparable](values ...T) *EnumValidator[T] {
	v := &EnumValidator[T]{values: make(map[T]struct{}, len(values))}
	v.Add(values...)
	return v
}

// IsValid 判断值是否属于枚举集合
func (v *EnumValidator[T]) IsValid(value T) bool {
	if v == nil {
		return false
	}
	_, ok := v.values[value]
	return ok
}

// MustBeValid 校验值是否属于枚举集合，失败时返回错误
func (v *EnumValidator[T]) MustBeValid(value T) error {
	if v.IsValid(value) {
		return nil
	}
	return fmt.Errorf(i18n.Msg(MsgEnumInvalidValue, map[string]string{"value": fmt.Sprint(value)}))
}

// GetValidValues 返回枚举值列表
func (v *EnumValidator[T]) GetValidValues() []T {
	if v == nil {
		return nil
	}
	out := make([]T, len(v.order))
	copy(out, v.order)
	return out
}

// GetValidValuesString 返回枚举值字符串列表
func (v *EnumValidator[T]) GetValidValuesString() []string {
	if v == nil {
		return nil
	}
	out := make([]string, 0, len(v.order))
	for _, item := range v.order {
		out = append(out, fmt.Sprint(item))
	}
	return out
}

// Count 返回枚举值数量
func (v *EnumValidator[T]) Count() int {
	if v == nil {
		return 0
	}
	return len(v.values)
}

// Contains 判断值是否存在于枚举集合
func (v *EnumValidator[T]) Contains(value T) bool {
	return v.IsValid(value)
}

// Add 增加枚举值
func (v *EnumValidator[T]) Add(values ...T) {
	if v.values == nil {
		v.values = make(map[T]struct{}, len(values))
	}
	for _, item := range values {
		if _, ok := v.values[item]; ok {
			continue
		}
		v.values[item] = struct{}{}
		v.order = append(v.order, item)
	}
}

// Remove 删除枚举值
func (v *EnumValidator[T]) Remove(values ...T) {
	if v == nil {
		return
	}
	for _, item := range values {
		if _, ok := v.values[item]; !ok {
			continue
		}
		delete(v.values, item)
		next := v.order[:0]
		for _, existing := range v.order {
			if existing != item {
				next = append(next, existing)
			}
		}
		v.order = next
	}
}

// Clear 清空枚举集合
func (v *EnumValidator[T]) Clear() {
	if v == nil {
		return
	}
	clear(v.values)
	v.order = nil
}

// Clone 复制枚举校验器
func (v *EnumValidator[T]) Clone() *EnumValidator[T] {
	if v == nil {
		return nil
	}
	return NewEnumValidator(v.order...)
}
