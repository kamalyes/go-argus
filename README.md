<div align="center">

# ⚡ Argus

**零依赖 · 高性能 · i18n 原生支持的 Go 结构体校验器**

[![Go Reference](https://pkg.go.dev/badge/github.com/kamalyes/go-argus.svg)](https://pkg.go.dev/github.com/kamalyes/go-argus)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamalyes/go-argus)](https://goreportcard.com/report/github.com/kamalyes/go-argus)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

[English](#) · [中文](#)

</div>

---

## ✨ 特性

- 🚀 **零第三方依赖** — 仅依赖 Go 标准库，供应链安全无忧
- 🏷️ **87+ 内置字段规则** — required、min/max、email、IP、UUID、datetime、Luhn 校验等
- 🔗 **跨字段规则** — range（范围校验）、fieldcontains（字段包含）、requiredWithout 等
- 🌍 **i18n 原生支持** — 内置 9 种语言翻译（en/zh/zh-TW/ja/ko/fr/de/es/ru），一行代码切换，可扩展任意语言
- 🔄 **go-playground/validator 兼容** — struct tag 语法和 API 高度兼容，迁移成本极低
- 🧩 **JSON Schema 校验** — 轻量 JSON Schema 子集校验，适合 API 网关场景
- 🔒 **并发安全** — 校验器实例可复用，struct 编译结果自动缓存
- 🛠️ **自定义规则** — 支持 `RegisterValidation` 注册自定义校验函数，支持 context 透传
- 📊 **数组化错误输出** — `TranslateValidationErrors` 直接输出可序列化的 JSON 错误
- 🌐 **网关工具** — IP 黑白名单（CIDR/通配符）、HTTP 状态码、Header、Content-Type、JSON Path 校验
- 📎 **格式校验** — email、IP、UUID、base64、URL、协议、WebSocket
- 📦 **泛型枚举校验器** — `NewEnumValidator[T]` 类型安全的枚举值校验
- 🔀 **标签逗号转义** — `\,` 在参数中保留逗号，`|` 作为替代分隔符
- 🛑 **规则执行策略** — 单字段失败即短路，其他字段不受影响

---

## 🏗️ 架构

```mermaid
graph TB
    subgraph "用户层 User Layer"
        APP["应用代码"]
    end

    subgraph "根包 validator"
        V["Validate 实例<br/>Struct / Var 校验"]
        CACHE["编译缓存<br/>structPlan"]
        TAGS["内置规则<br/>87+ builtinRules"]
        TRANS["错误翻译<br/>translations.go → i18n.Lookup"]
        OPTS["配置选项<br/>Option / SetLocale"]
        ERRORS["错误模型<br/>ValidationErrors"]
    end

    subgraph "rule 包"
        RPARSE["标签解析<br/>ParseTag"]
        RFIELD["字段路径<br/>FieldByPath"]
        RTIME["时间规则<br/>TimeValue / ResolveTimeExpr"]
    end

    subgraph "validate 包"
        COMPARE["比较校验<br/>CompareNumbers / Strings"]
        FORMAT["格式校验<br/>Email / IP / URL / UUID / Base64"]
        ENUM["枚举校验<br/>EnumValidator"]
        JSON["JSON 校验<br/>ValidateJSON / JSONPath"]
        NETWORK["网络校验<br/>IPSet / CIDR / 通配符"]
        CONSTANTS["消息常量<br/>constants.go"]
    end

    subgraph "i18n 包（统一翻译存储）"
        I18N["i18n 核心<br/>SetLocale / Msg / Lookup / Register"]
        I18N_EN["en.go"]
        I18N_ZH["zh.go / zh_tw.go"]
        I18N_JA["ja.go / ko.go"]
        I18N_OTHER["fr.go / de.go / es.go / ru.go"]
    end

    subgraph "schema 包"
        SCHEMA["JSON Schema<br/>ValidateJSONSchema"]
    end

    APP --> V
    V --> CACHE
    V --> TAGS
    V --> TRANS
    V --> OPTS
    V --> ERRORS
    V --> RPARSE
    V --> RFIELD
    V --> RTIME

    TAGS --> FORMAT
    TAGS --> COMPARE
    TAGS --> ENUM

    APP --> COMPARE
    APP --> FORMAT
    APP --> NETWORK
    APP --> SCHEMA

    SCHEMA --> I18N
    COMPARE --> I18N
    FORMAT --> I18N
    ENUM --> I18N
    JSON --> I18N
    NETWORK --> I18N

    OPTS --> I18N
    TRANS --> I18N

    style APP fill:#e1f5fe
    style V fill:#fff3e0
    style I18N fill:#e8f5e9
    style SCHEMA fill:#fce4ec
```

## 📦 安装

```bash
go get github.com/kamalyes/go-argus
```

> 要求 Go 1.21+

## 🚀 快速开始

```go
package main

import (
    "fmt"
    "github.com/kamalyes/go-argus"
)

type User struct {
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"gte=0,lte=150"`
}

func main() {
    v := validator.New()
    err := v.Struct(User{Name: "A", Email: "bad", Age: -1})

    // 一行切换语言
    validator.SetLocale("zh")
    messages := validator.TranslateValidationErrors(err, "zh")
    for _, msg := range messages {
        fmt.Printf("%s: %s\n", msg.Field, msg.Message)
    }
    // 注册新语言（9 种内置语言：en/zh/zh-TW/ja/ko/fr/de/es/ru）
    validator.RegisterI18nMessages("pt", map[string]string{
        "required": "{field} é obrigatório",
    })
    // name: name 不能小于 2
    // email: email 必须是有效的 Email
    // age: age 必须大于或等于 0
}
```

## 📚 文档

| 文档 | 说明 |
|------|------|
| [docs/tags.md](docs/tags.md) | 所有校验标签完整参考 |
| [docs/i18n.md](docs/i18n.md) | 国际化使用指南 |
| [docs/examples.md](docs/examples.md) | 完整使用示例 |

---

## 🔄 从 go-playground/validator 迁移

Argus 的 struct tag 语法和核心 API 与 `go-playground/validator` 高度兼容：

```go
// go-playground/validator
import "github.com/go-playground/validator/v10"
v := validator.New()

// Argus — 只需改 import 路径
import "github.com/kamalyes/go-argus"
v := validator.New()
```

主要差异：

| 特性 | go-playground/validator | Argus |
|------|------------------------|-------|
| 第三方依赖 | 多个（如 utranslator） | **零依赖** |
| i18n | 需额外安装 translator | **内置 9 种语言** |
| JSON Schema | 不支持 | **内置** |
| IP/CIDR/网络 | 不支持 | **内置** |
| 比较校验 | 不支持 | **内置** |

---

## 📄 License

[MIT License](LICENSE)
