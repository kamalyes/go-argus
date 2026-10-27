/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-27 11:08:31
 * @FilePath: \go-argus\i18n\i18n_test.go
 * @Description: 国际化消息资源测试
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package i18n

import (
	"strings"
	"sync"
	"testing"
)

func TestSetLocaleAndGetLocale(t *testing.T) {
	original := GetLocale()

	SetLocale("zh")
	if got := GetLocale(); got != "zh" {
		t.Errorf("GetLocale() = %q, want %q", got, "zh")
	}

	SetLocale("ja")
	if got := GetLocale(); got != "ja" {
		t.Errorf("GetLocale() = %q, want %q", got, "ja")
	}

	SetLocale(original)
}

func TestSetLocaleNormalizes(t *testing.T) {
	original := GetLocale()

	SetLocale("ZH")
	if got := GetLocale(); got != "zh" {
		t.Errorf("SetLocale(ZH): GetLocale() = %q, want %q", got, "zh")
	}

	SetLocale("JA-JP")
	if got := GetLocale(); got != "ja" {
		t.Errorf("SetLocale(JA-JP): GetLocale() = %q, want %q", got, "ja")
	}

	SetLocale(original)
}

func TestRegister(t *testing.T) {
	original := GetLocale()
	SetLocale("en")

	Register("en", "test.key", "hello {name}")

	got := Msg("test.key", map[string]string{"name": "world"})
	want := "hello world"
	if got != want {
		t.Errorf("Msg after Register = %q, want %q", got, want)
	}

	Register("en", "test.key", "overridden")
	got = Msg("test.key")
	if got != "overridden" {
		t.Errorf("Msg after override = %q, want %q", got, "overridden")
	}

	SetLocale(original)
}

func TestRegisterEmptyParams(t *testing.T) {
	Register("", "key", "value")
	Register("en", "", "value")
	Register("en", "key", "")
	Register("en", "key", "   ")
}

func TestRegisterNewLocale(t *testing.T) {
	Register("pt", "greeting", "olá")
	got := Lookup("pt", "greeting")
	if got != "olá" {
		t.Errorf("Lookup(pt, greeting) = %q, want %q", got, "olá")
	}
}

func TestRegisterMessages(t *testing.T) {
	original := GetLocale()
	SetLocale("en")

	msgs := map[string]string{
		"batch.key1": "batch value 1",
		"batch.key2": "batch value 2",
	}
	RegisterMessages("en", msgs)

	if got := Msg("batch.key1"); got != "batch value 1" {
		t.Errorf("Msg(batch.key1) = %q, want %q", got, "batch value 1")
	}
	if got := Msg("batch.key2"); got != "batch value 2" {
		t.Errorf("Msg(batch.key2) = %q, want %q", got, "batch value 2")
	}

	SetLocale(original)
}

func TestMsgWithArgs(t *testing.T) {
	original := GetLocale()
	SetLocale("en")

	Register("en", "greet", "hello {name}, age {age}")
	got := Msg("greet", map[string]string{"name": "alice", "age": "30"})
	want := "hello alice, age 30"
	if got != want {
		t.Errorf("Msg with args = %q, want %q", got, want)
	}

	SetLocale(original)
}

func TestMsgWithoutArgs(t *testing.T) {
	original := GetLocale()
	SetLocale("en")

	Register("en", "plain", "no placeholders")
	got := Msg("plain")
	if got != "no placeholders" {
		t.Errorf("Msg without args = %q, want %q", got, "no placeholders")
	}

	SetLocale(original)
}

func TestMsgMissingKey(t *testing.T) {
	original := GetLocale()
	SetLocale("en")

	got := Msg("nonexistent.key")
	if got != "nonexistent.key" {
		t.Errorf("Msg(missing) = %q, want %q", got, "nonexistent.key")
	}

	SetLocale(original)
}

func TestLookupFound(t *testing.T) {
	got := Lookup("en", "required")
	if got == "" {
		t.Error("Lookup(en, required) returned empty string")
	}
}

func TestLookupFallbackToEn(t *testing.T) {
	Register("xx", "unique.xx", "xx value")

	got := Lookup("xx", "required")
	if got == "" {
		t.Error("Lookup(xx, required) should fall back to en")
	}
}

func TestLookupEnMissingKey(t *testing.T) {
	got := Lookup("en", "absolutely.nonexistent.key")
	if got != "" {
		t.Errorf("Lookup(en, nonexistent) = %q, want empty", got)
	}
}

func TestLookupNonEnMissingKey(t *testing.T) {
	got := Lookup("xx", "absolutely.nonexistent.key")
	if got != "" {
		t.Errorf("Lookup(xx, nonexistent) = %q, want empty", got)
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "en"},
		{"en", "en"},
		{"EN", "en"},
		{"en-US", "en"},
		{"en_GB", "en"},
		{"zh", "zh"},
		{"ZH", "zh"},
		{"zh-CN", "zh"},
		{"zh-cn", "zh"},
		{"zh_CN", "zh"},
		{"zh-TW", "zh-TW"},
		{"zh-tw", "zh-TW"},
		{"zh-Hant", "zh-TW"},
		{"zh-hant", "zh-TW"},
		{"zh_Hant", "zh-TW"},
		{"ja", "ja"},
		{"JA", "ja"},
		{"ja-JP", "ja"},
		{"ko", "ko"},
		{"KO", "ko"},
		{"ko-KR", "ko"},
		{"fr", "fr"},
		{"FR", "fr"},
		{"fr-FR", "fr"},
		{"de", "de"},
		{"DE", "de"},
		{"de-DE", "de"},
		{"es", "es"},
		{"ES", "es"},
		{"es-ES", "es"},
		{"ru", "ru"},
		{"RU", "ru"},
		{"ru-RU", "ru"},
		{"pt", "pt"},
		{"it", "it"},
		{"  en  ", "en"},
		{"  zh  ", "zh"},
	}

	for _, tt := range tests {
		got := Normalize(tt.input)
		if got != tt.want {
			t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNormalizeUnderscore(t *testing.T) {
	got := Normalize("zh_CN")
	if got != "zh" {
		t.Errorf("Normalize(zh_CN) = %q, want %q", got, "zh")
	}
}

func TestInitStoresAllLanguages(t *testing.T) {
	locales := []string{"en", "zh", "zh-TW", "ja", "ko", "fr", "de", "es", "ru"}
	for _, l := range locales {
		got := Lookup(l, "required")
		if got == "" {
			t.Errorf("Lookup(%s, required) returned empty, expected a message", l)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			SetLocale("en")
		}()
		go func() {
			defer wg.Done()
			_ = GetLocale()
		}()
		go func() {
			defer wg.Done()
			_ = Msg("required")
		}()
	}
	wg.Wait()
}

func TestAllMessageFunctionsReturnNonEmpty(t *testing.T) {
	funcs := map[string]func() map[string]string{
		"en":    EnMessages,
		"zh":    ZhMessages,
		"zh-TW": ZhTWMessages,
		"ja":    JaMessages,
		"ko":    KoMessages,
		"fr":    FrMessages,
		"de":    DeMessages,
		"es":    EsMessages,
		"ru":    RuMessages,
	}

	for name, fn := range funcs {
		msgs := fn()
		if len(msgs) == 0 {
			t.Errorf("%sMessages() returned empty map", name)
		}
		for key, val := range msgs {
			if val == "" {
				t.Errorf("%sMessages()[%q] is empty", name, key)
			}
		}
	}
}

// TestAllLocalesHaveSameKeysAsEn 确保所有语言的键集合与 en 完全一致，防止遗漏翻译键
func TestAllLocalesHaveSameKeysAsEn(t *testing.T) {
	allFuncs := map[string]func() map[string]string{
		"zh":    ZhMessages,
		"zh-TW": ZhTWMessages,
		"ja":    JaMessages,
		"ko":    KoMessages,
		"fr":    FrMessages,
		"de":    DeMessages,
		"es":    EsMessages,
		"ru":    RuMessages,
	}

	enKeys := make(map[string]struct{}, len(EnMessages()))
	for k := range EnMessages() {
		enKeys[k] = struct{}{}
	}

	for locale, fn := range allFuncs {
		msgs := fn()
		for k := range enKeys {
			if _, ok := msgs[k]; !ok {
				t.Errorf("locale %q is missing key %q (present in en)", locale, k)
			}
		}
		for k := range msgs {
			if _, ok := enKeys[k]; !ok {
				t.Errorf("locale %q has extra key %q (not in en)", locale, k)
			}
		}
	}
}

func TestMsgWithEmptyArgs(t *testing.T) {
	original := GetLocale()
	SetLocale("en")

	Register("en", "test.empty.args", "value {placeholder}")
	got := Msg("test.empty.args")
	if got != "value {placeholder}" {
		t.Errorf("Msg with no args map = %q, want %q", got, "value {placeholder}")
	}

	got = Msg("test.empty.args", map[string]string{})
	if got != "value {placeholder}" {
		t.Errorf("Msg with empty args map = %q, want %q", got, "value {placeholder}")
	}

	SetLocale(original)
}

func TestLookupSameAsEn(t *testing.T) {
	got := Lookup("en", "nonexistent_same_as_en")
	if got != "" {
		t.Errorf("Lookup(en, nonexistent) = %q, want empty", got)
	}
}

func TestRegisterMessagesWithEmptyParams(t *testing.T) {
	RegisterMessages("", map[string]string{"key": "value"})
	RegisterMessages("en", map[string]string{"": "value"})
	RegisterMessages("en", map[string]string{"key": ""})
	RegisterMessages("en", map[string]string{"key": "   "})
}

// ══════════════════════════════════════════════════════════════════════════════
// 消息渲染：单次扫描替换（原 zz_render_probe_test.go，合并至此）
//
// 旧实现对每个占位符各做一次 strings.ReplaceAll，每次都整体复制一遍当前消息，
// 故 N 个占位符为 N 次全量拷贝。改单次扫描后降为一次缓冲写入。
//
// 更要紧的是旧实现**不确定**：链式 ReplaceAll 遍历 map，若某替换值自身含 {k}，
// 它会不会被二次替换取决于 map 迭代顺序——同一输入 200 次运行可得两种输出。
// TestRender_NoRecursiveSubstitution 即该不确定性的回归，故本次改动是正确性修复
// 而非单纯性能优化。TestRender_MatchesReplaceAll 守与旧实现的输出等价
// （已覆盖 "}反序{ {field}" 这类嵌套花括号的对抗输入）。
// ══════════════════════════════════════════════════════════════════════════════

// replaceAllRender 原实现，用于等价性对照与性能基线
func replaceAllRender(template string, args map[string]string) string {
	for k, v := range args {
		template = strings.ReplaceAll(template, "{"+k+"}", v)
	}
	return template
}

// TestRender_MatchesReplaceAll 单次扫描须与原 ReplaceAll 实现等价
//
// 仅覆盖「替换值本身不含占位符」的输入——原实现在那种输入下结果依赖
// map 迭代顺序（见 TestRender_NoRecursiveSubstitution），不构成可对照的基准。
func TestRender_MatchesReplaceAll(t *testing.T) {
	cases := []struct {
		tmpl string
		args map[string]string
	}{
		{"字段 {field} 不能为空", map[string]string{"field": "Name"}},
		{"{field} 须在 {min} 与 {max} 之间", map[string]string{"field": "Age", "min": "1", "max": "9"}},
		{"无占位符", map[string]string{"field": "X"}},
		{"{field}{field}{field}", map[string]string{"field": "A"}},
		{"缺失的 {nope} 保留", map[string]string{"field": "X"}},
		{"{field} 与未闭合 {brace", map[string]string{"field": "F"}},
		{"空值 {empty} 结束", map[string]string{"empty": ""}},
		{"{a}{b}", map[string]string{"a": "1", "b": "2"}},
		{"前缀 {field}", map[string]string{"field": "尾部"}},
		{"{field} 后缀", map[string]string{"field": "开头"}},
		{"}反序{ {field}", map[string]string{"field": "V"}},
		{"", map[string]string{"field": "X"}},
	}
	for _, c := range cases {
		want := replaceAllRender(c.tmpl, c.args)
		got := render(c.tmpl, c.args)
		if got != want {
			t.Errorf("render 与 ReplaceAll 不一致\n  tmpl=%q args=%v\n  want=%q\n  got =%q",
				c.tmpl, c.args, want, got)
		}
	}
}

// TestRender_NoRecursiveSubstitution 替换值内的占位符不得被二次替换
//
// 这是单次扫描相对原实现的**行为更正**，而非回归：
// 原实现把 ReplaceAll 串行作用于已替换的中间结果，故当替换值含 {b} 且 b 也在
// args 中时，是否被二次替换完全取决于 map 迭代顺序——同输入可得不同输出。
// 单次扫描只认模板自身的占位符，结果恒定。
func TestRender_NoRecursiveSubstitution(t *testing.T) {
	tmpl := "值={a}"
	args := map[string]string{"a": "{b}", "b": "BOOM"}

	got := render(tmpl, args)
	if got != "值={b}" {
		t.Errorf("替换值内的占位符被二次替换：期望 %q，实得 %q", "值={b}", got)
	}

	// 证明原实现在此输入下不确定：多轮运行可能得到两种结果之一
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		seen[replaceAllRender(tmpl, args)] = true
	}
	t.Logf("原 ReplaceAll 实现在该输入下的观测结果集: %v（>1 即证明依赖 map 迭代顺序）", keysOf(seen))
}

func keysOf(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// TestRender_Idempotent 同输入多次调用结果恒定
func TestRender_Idempotent(t *testing.T) {
	tmpl := "{field} 须在 {min} 与 {max} 之间"
	args := map[string]string{"field": "Age", "min": "1", "max": "9"}
	first := render(tmpl, args)
	for i := 0; i < 100; i++ {
		if got := render(tmpl, args); got != first {
			t.Fatalf("第 %d 次结果不同：%q vs %q", i, got, first)
		}
	}
}

// ── 性能对照 ──

var benchTmpl1 = "字段 {field} 不能为空"
var benchTmpl3 = "{field} 的取值 {actual} 不在 {expected} 允许范围内"
var benchArgs1 = map[string]string{"field": "Name"}
var benchArgs3 = map[string]string{"field": "Age", "actual": "120", "expected": "1-99"}

func BenchmarkProbeRender_ReplaceAll_1(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = replaceAllRender(benchTmpl1, benchArgs1)
	}
}

func BenchmarkProbeRender_SinglePass_1(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = render(benchTmpl1, benchArgs1)
	}
}

func BenchmarkProbeRender_ReplaceAll_3(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = replaceAllRender(benchTmpl3, benchArgs3)
	}
}

func BenchmarkProbeRender_SinglePass_3(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = render(benchTmpl3, benchArgs3)
	}
}

// BenchmarkProbeRender_NoPlaceholder 无占位符模板应零分配
func BenchmarkProbeRender_NoPlaceholder(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = render("没有任何占位符的消息模板", benchArgs3)
	}
}

// BenchmarkProbeMsg_End2End 经 Msg 的端到端（含查表）
func BenchmarkProbeMsg_End2End(b *testing.B) {
	Register("en", "probe.key", benchTmpl3)
	SetLocale("en")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Msg("probe.key", benchArgs3)
	}
}
