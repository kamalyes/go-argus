/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-16 00:00:00
 * @FilePath: \go-argus\translations.go
 * @Description: 零依赖错误翻译模块，提供字段级 i18n 提示和数组化错误输出
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kamalyes/go-argus/i18n"
)

// ValidationMessage 表示一个可直接序列化给 HTTP/gRPC 网关的字段错误
type ValidationMessage struct {
	Field           string      `json:"field"`
	Namespace       string      `json:"namespace"`
	StructField     string      `json:"struct_field,omitempty"`
	StructNamespace string      `json:"struct_namespace,omitempty"`
	Tag             string      `json:"tag"`
	ActualTag       string      `json:"actual_tag"`
	Param           string      `json:"param,omitempty"`
	Value           interface{} `json:"value,omitempty"`
	Message         string      `json:"message"`
}

// RegisterTranslation 注册或覆盖某个语言下的单个规则翻译模板
func RegisterTranslation(locale string, tag string, template string) {
	i18n.Register(locale, tag, template)
}

// RegisterTranslations 批量注册某个语言下的规则翻译模板
func RegisterTranslations(locale string, items map[string]string) {
	i18n.RegisterMessages(locale, items)
}

// TranslateValidationErrors 将 error 转换为数组化的 i18n 字段错误
func TranslateValidationErrors(err error, locale string) []ValidationMessage {
	if err == nil {
		return nil
	}
	var validationErrors ValidationErrors
	if errors.As(err, &validationErrors) {
		return validationErrors.Translate(locale)
	}
	return []ValidationMessage{{Message: err.Error()}}
}

// Translate 将 ValidationErrors 转换为数组化的 i18n 字段错误
func (ve ValidationErrors) Translate(locale string) []ValidationMessage {
	if len(ve) == 0 {
		return nil
	}
	messages := make([]ValidationMessage, 0, len(ve))
	for _, fe := range ve {
		messages = append(messages, translateFieldError(fe, locale))
	}
	return messages
}

// RequiredMessages 只返回 required 系列规则产生的缺失字段错误
func (ve ValidationErrors) RequiredMessages(locale string) []ValidationMessage {
	if len(ve) == 0 {
		return nil
	}
	messages := make([]ValidationMessage, 0, len(ve))
	for _, fe := range ve {
		if isRequiredTag(fe.Tag()) {
			messages = append(messages, translateFieldError(fe, locale))
		}
	}
	return messages
}

// MissingFields 返回 required 系列规则中缺失的字段名，字段名优先使用 json tag
func (ve ValidationErrors) MissingFields() []string {
	if len(ve) == 0 {
		return nil
	}
	fields := make([]string, 0, len(ve))
	for _, fe := range ve {
		if isRequiredTag(fe.Tag()) {
			fields = append(fields, fe.Field())
		}
	}
	return fields
}

func translateFieldError(fe FieldError, locale string) ValidationMessage {
	msg := ValidationMessage{
		Field:           fe.Field(),
		Namespace:       fe.Namespace(),
		StructField:     fe.StructField(),
		StructNamespace: fe.StructNamespace(),
		Tag:             fe.Tag(),
		ActualTag:       fe.ActualTag(),
		Param:           fe.Param(),
		Value:           safeMessageValue(fe.Value()),
	}
	msg.Message = renderTranslation(locale, fe)
	return msg
}

func renderTranslation(locale string, fe FieldError) string {
	template := lookupTranslation(locale, fe.Tag())
	if template == "" {
		template = lookupTranslation(locale, "default")
	}
	if template == "" {
		return fe.Error()
	}
	replacements := map[string]string{
		"{field}":     fe.Field(),
		"{namespace}": fe.Namespace(),
		"{tag}":       fe.Tag(),
		"{param}":     fe.Param(),
		"{value}":     fmt.Sprint(fe.Value()),
	}
	for old, value := range replacements {
		template = strings.ReplaceAll(template, old, value)
	}
	return template
}

func lookupTranslation(locale string, tag string) string {
	return i18n.Lookup(locale, tag)
}

func isRequiredTag(tag string) bool {
	switch tag {
	case "required", "required_if", "required_unless", "required_with", "required_with_all", "required_without", "required_without_all":
		return true
	default:
		return false
	}
}

func safeMessageValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		if rv.IsNil() {
			return nil
		}
	}
	return value
}
