/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-16 00:00:00
 * @FilePath: \go-argus\i18n\i18n.go
 * @Description: 国际化消息资源，提供英文、中文、日语、韩语、法语、德语、西班牙语、俄语、繁体中文的错误消息
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"strings"
	"sync"
)

var (
	mu     sync.RWMutex
	store  = map[string]map[string]string{}
	locale = "en"
)

func init() {
	store["en"] = EnMessages()
	store["zh"] = ZhMessages()
	store["zh-TW"] = ZhTWMessages()
	store["ja"] = JaMessages()
	store["ko"] = KoMessages()
	store["fr"] = FrMessages()
	store["de"] = DeMessages()
	store["es"] = EsMessages()
	store["ru"] = RuMessages()
}

// SetLocale 设置全局语言环境
func SetLocale(l string) {
	l = Normalize(l)
	mu.Lock()
	locale = l
	mu.Unlock()
}

// GetLocale 返回当前全局语言环境
func GetLocale() string {
	mu.RLock()
	defer mu.RUnlock()
	return locale
}

// Register 注册或覆盖某个语言下的单个翻译模板
func Register(l, key, template string) {
	l = Normalize(l)
	key = strings.TrimSpace(key)
	if l == "" || key == "" || strings.TrimSpace(template) == "" {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if store[l] == nil {
		store[l] = make(map[string]string)
	}
	store[l][key] = template
}

// RegisterMessages 批量注册某个语言下的翻译模板
func RegisterMessages(l string, messages map[string]string) {
	for key, template := range messages {
		Register(l, key, template)
	}
}

// Msg 根据 key 查找当前语言环境的消息模板，并用 args 替换占位符
func Msg(key string, args ...map[string]string) string {
	template := Lookup(GetLocale(), key)
	if template == "" {
		return key
	}
	if len(args) > 0 {
		for k, v := range args[0] {
			template = strings.ReplaceAll(template, "{"+k+"}", v)
		}
	}
	return template
}

// Lookup 按 locale → en 顺序查找消息模板
func Lookup(l, key string) string {
	l = Normalize(l)
	mu.RLock()
	defer mu.RUnlock()
	if items := store[l]; items != nil {
		if t := items[key]; t != "" {
			return t
		}
	}
	if l != "en" {
		if items := store["en"]; items != nil {
			return items[key]
		}
	}
	return ""
}

// Normalize 规范化语言标签
func Normalize(l string) string {
	l = strings.ToLower(strings.TrimSpace(l))
	l = strings.ReplaceAll(l, "_", "-")
	switch {
	case l == "", strings.HasPrefix(l, "en"):
		return "en"
	case l == "zh-tw", l == "zh-hant":
		return "zh-TW"
	case l == "zh", strings.HasPrefix(l, "zh-"):
		return "zh"
	case strings.HasPrefix(l, "ja"):
		return "ja"
	case strings.HasPrefix(l, "ko"):
		return "ko"
	case strings.HasPrefix(l, "fr"):
		return "fr"
	case strings.HasPrefix(l, "de"):
		return "de"
	case strings.HasPrefix(l, "es"):
		return "es"
	case strings.HasPrefix(l, "ru"):
		return "ru"
	default:
		return l
	}
}
