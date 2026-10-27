/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-27 11:08:31
 * @FilePath: \go-argus\i18n\i18n.go
 * @Description: 国际化消息资源，提供英文、中文、日语、韩语、法语、德语、西班牙语、俄语、繁体中文的错误消息
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"strings"
	"sync"
	"sync/atomic"
)

// snapshot 消息库与语言环境的不可变快照
//
// 读路径（Msg/Lookup/GetLocale）一次原子加载即取得全部状态，不加锁。
// 原实现是 RWMutex + 裸 map，且 Msg 会连取两次 RLock
// （GetLocale 一次、Lookup 再一次）——RLock 对共享 readerCount 做原子 RMW，
// 多核在同一缓存行上乒乓，实测错误消息渲染在 4 核后不再提速、10 核反而退化
// 写路径（SetLocale/Register）在 writeMu 下复制整份快照再发布，属启动期/低频行为
type snapshot struct {
	store  map[string]map[string]string
	locale string // 已 Normalize，读路径无需再规范化
}

var (
	writeMu sync.Mutex
	state   atomic.Pointer[snapshot]
)

func init() {
	state.Store(&snapshot{
		store: map[string]map[string]string{
			"en":    EnMessages(),
			"zh":    ZhMessages(),
			"zh-TW": ZhTWMessages(),
			"ja":    JaMessages(),
			"ko":    KoMessages(),
			"fr":    FrMessages(),
			"de":    DeMessages(),
			"es":    EsMessages(),
			"ru":    RuMessages(),
		},
		locale: "en",
	})
}

// cloneLocked 复制当前快照，调用方必须持有 writeMu
func cloneLocked() *snapshot {
	cur := state.Load()
	next := &snapshot{
		store:  make(map[string]map[string]string, len(cur.store)+1),
		locale: cur.locale,
	}
	for l, items := range cur.store {
		// 语言表整体复用：Register 只会替换目标语言那一张表，其余共享无写入
		next.store[l] = items
	}
	return next
}

// SetLocale 设置全局语言环境
func SetLocale(l string) {
	l = Normalize(l)
	writeMu.Lock()
	defer writeMu.Unlock()
	next := cloneLocked()
	next.locale = l
	state.Store(next)
}

// GetLocale 返回当前全局语言环境
func GetLocale() string {
	return state.Load().locale
}

// Register 注册或覆盖某个语言下的单个翻译模板
func Register(l, key, template string) {
	l = Normalize(l)
	key = strings.TrimSpace(key)
	if l == "" || key == "" || strings.TrimSpace(template) == "" {
		return
	}
	writeMu.Lock()
	defer writeMu.Unlock()
	next := cloneLocked()
	// 目标语言表需单独深拷贝：快照共享其余语言表，直接改会让已发布的旧快照被并发读到半写状态
	items := make(map[string]string, len(next.store[l])+1)
	for k, v := range next.store[l] {
		items[k] = v
	}
	items[key] = template
	next.store[l] = items
	state.Store(next)
}

// RegisterMessages 批量注册某个语言下的翻译模板
func RegisterMessages(l string, messages map[string]string) {
	for key, template := range messages {
		Register(l, key, template)
	}
}

// Msg 根据 key 查找当前语言环境的消息模板，并用 args 替换占位符
//
// 单次快照加载即取得 locale 与消息库：原实现 GetLocale + Lookup 各取一次 RLock，
// 且 Lookup 会对已规范化的 locale 再跑一次 Normalize（内含 ToLower/ReplaceAll 分配）
func Msg(key string, args ...map[string]string) string {
	st := state.Load()
	// st.locale 恒为已规范化值，直接查表，跳过 Normalize
	template := lookupIn(st, st.locale, key)
	if template == "" {
		return key
	}
	if len(args) > 0 && len(args[0]) > 0 {
		return render(template, args[0])
	}
	return template
}

// render 单次扫描替换 {name} 占位符
//
// 原实现是 `for k, v := range args { template = strings.ReplaceAll(template, "{"+k+"}", v) }`，
// 每个占位符付两次分配（"{"+k+"}" 拼接，以及 ReplaceAll 产出的新串），
// 且每轮都重扫整个模板——时间 O(len×占位符数)、分配 O(占位符数)。
// 单次扫描改为 O(len) 时间、一次结果分配；键名取模板子串，不分配。
//
// 行为差异（更正了一处隐性不确定性）：原实现逐个 ReplaceAll 串行作用于
// **已替换**的中间结果，故当某个替换值自身含形如 {k} 的文本且 k 也在 args 中时，
// 它会被二次替换，且结果取决于 map 迭代顺序——同样输入可能产出不同输出。
// 单次扫描只替换模板自身的占位符，不再回头扫描替换值，结果恒定。
//
// 未在 args 中的占位符原样保留（与原实现一致）。
func render(template string, args map[string]string) string {
	// 无占位符的模板直接返回，零分配
	if strings.IndexByte(template, '{') < 0 {
		return template
	}

	var b strings.Builder
	b.Grow(len(template))

	for i := 0; i < len(template); {
		open := strings.IndexByte(template[i:], '{')
		if open < 0 {
			b.WriteString(template[i:])
			break
		}
		open += i
		b.WriteString(template[i:open])

		rest := template[open+1:]
		closeRel := strings.IndexByte(rest, '}')
		if closeRel < 0 {
			// 无闭合花括号：余下部分原样输出
			b.WriteString(template[open:])
			break
		}

		// 若在闭合 '}' 之前又出现 '{'，则当前 '{' 是字面量，不能把内层占位符一并吞掉。
		// 例：模板 "}反序{ {field}" 中第一个 '{' 后面紧跟的是另一个 '{'，
		// 若直接扫到末尾的 '}'，键名会变成 " {field" 而真正的 {field} 不被替换。
		if nextOpenRel := strings.IndexByte(rest, '{'); nextOpenRel >= 0 && nextOpenRel < closeRel {
			b.WriteByte('{')
			i = open + 1
			continue
		}

		closeIdx := open + 1 + closeRel
		// 键名为模板子串，不产生分配；map 查找亦不分配
		if v, ok := args[template[open+1:closeIdx]]; ok {
			b.WriteString(v)
		} else {
			b.WriteString(template[open : closeIdx+1])
		}
		i = closeIdx + 1
	}
	return b.String()
}

// Lookup 按 locale → en 顺序查找消息模板
//
// 公开入口，locale 来自外部输入故仍需 Normalize；内部热路径走 lookupIn
func Lookup(l, key string) string {
	return lookupIn(state.Load(), Normalize(l), key)
}

// lookupIn 在给定快照内查表，l 必须已规范化
func lookupIn(st *snapshot, l, key string) string {
	if items := st.store[l]; items != nil {
		if t := items[key]; t != "" {
			return t
		}
	}
	if l != "en" {
		if items := st.store["en"]; items != nil {
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
