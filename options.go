/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
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

// SetLocale 设置全局语言环境
func SetLocale(locale string) {
	i18n.SetLocale(locale)
}

// GetLocale 返回当前全局语言环境
func GetLocale() string {
	return i18n.GetLocale()
}

// RegisterI18n 注册或覆盖某个语言下的单个 i18n 消息模板
func RegisterI18n(locale string, key string, template string) {
	i18n.Register(locale, key, template)
}

// RegisterI18nMessages 批量注册某个语言下的 i18n 消息模板
func RegisterI18nMessages(locale string, messages map[string]string) {
	i18n.RegisterMessages(locale, messages)
}
