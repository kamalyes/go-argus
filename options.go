/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-25 17:35:19
 * @FilePath: \go-argus\options.go
 * @Description: 校验器配置项，提供兼容 go-playground validator 的构造选项
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validator

import "github.com/kamalyes/go-argus/i18n"

// Option 表示校验器初始化配置项
type Option func(*Validate)

// WithRequiredStructEnabled 允许 required 对非指针结构体零值生效
func WithRequiredStructEnabled() Option {
	return func(v *Validate) {
		v.requiredStructEnabled = true
	}
}

// WithPrivateFieldValidation 保留私有字段校验选项，用于平滑迁移 go-playground 调用
func WithPrivateFieldValidation() Option {
	return func(v *Validate) {
		v.privateFieldValidation = true
	}
}

// WithFailFast 首个错误即停止校验（默认关闭）
//
// 默认行为是收集全部字段错误：命中失败后仍继续遍历剩余字段，
// 因此"首字段失败"与"末字段失败"耗时相同（实测 15 字段结构体 482.9 vs 461.7 ns）。
// 请求入口的常见诉求是"拒绝这个请求"，调用方往往只需首个错误或 err != nil，
// 此时继续校验剩余字段是纯浪费。
//
// 开启后 ValidationErrors 只含首个错误。必须显式开启而非默认——
// 依赖"拿到全部错误"（如表单逐字段回显）的调用方不能被静默改掉语义。
func WithFailFast() Option {
	return func(v *Validate) {
		v.failFast = true
	}
}

// SetLocale 设置全局语言环境
func SetLocale(locale string) {
	i18n.SetLocale(locale)
}

// GetLocale 返回当前全局语言环境
func GetLocale() string {
	return i18n.GetLocale()
}

// RegisterI18n 注册或覆盖某个语言下的单个 i18n 消息模板
func RegisterI18n(locale, key, template string) {
	i18n.Register(locale, key, template)
}

// RegisterI18nMessages 批量注册某个语言下的 i18n 消息模板
func RegisterI18nMessages(locale string, messages map[string]string) {
	i18n.RegisterMessages(locale, messages)
}
