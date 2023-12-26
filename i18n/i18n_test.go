/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-16 00:00:00
 * @FilePath: \go-argus\i18n\i18n_test.go
 * @Description: 国际化消息资源测试
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package i18n

import (
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
