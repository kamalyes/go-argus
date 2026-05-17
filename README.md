<div align="center">

# ⚡ Argus

**零依赖 · 高性能 · i18n 原生支持的 Go 结构体校验器**

[![Go Reference](https://pkg.go.dev/badge/github.com/kamalyes/go-argus.svg)](https://pkg.go.dev/github.com/kamalyes/go-argus)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamalyes/go-argus)](https://goreportcard.com/report/github.com/kamalyes/go-argus)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

[English](README_EN.md) · [中文](#)

</div>

---

## ✨ 特性

- 🚀 **零第三方依赖** — 仅依赖 Go 标准库，供应链安全无忧
- ⚡ **零反射 VarString 快速路径** — 字符串变量校验完全绕过 `reflect`，0 堆分配，比反射路径快 2~3 倍
- 🏷️ **97+ 内置字段规则** — required、min/max、email、IP、UUID、datetime、Luhn 校验、semver、ISBN、ISSN、BIC/SWIFT、cron、Data URI、BCP 47、以太坊/比特币地址等
- 🔗 **跨字段规则** — range（范围校验）、fieldcontains（字段包含）、requiredWithout 等
- 🌍 **i18n 原生支持** — 内置 9 种语言翻译（en/zh/zh-TW/ja/ko/fr/de/es/ru），一行代码切换，可扩展任意语言
- 🔄 **go-playground/validator 兼容** — struct tag 语法和 API 高度兼容，迁移成本极低
- 🧩 **JSON Schema 校验** — 轻量 JSON Schema 子集校验，适合 API 网关场景
- 🔒 **并发安全** — 校验器实例可复用，struct 编译结果自动缓存
- 🛠️ **自定义规则** — 支持 `RegisterValidation` 注册自定义校验函数，支持 context 透传
- 📊 **数组化错误输出** — `TranslateValidationErrors` 直接输出可序列化的 JSON 错误
- 🌐 **网关工具** — IP 黑白名单（CIDR/通配符）、HTTP 状态码、Header、Content-Type、JSON Path 校验
- 📎 **格式校验** — email、IP、UUID、base64、URL、URI（含 mailto/tel）、协议、WebSocket、semver、ISBN-10/13、ISSN、BIC/SWIFT、cron、Data URI、BCP 47 语言标签、以太坊/比特币地址
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
        V["Validate 实例<br/>Struct / Var / VarString 校验"]
        CACHE["编译缓存<br/>structPlan"]
        TAGS["内置规则<br/>87+ builtinRules"]
        STAGS["字符串规则<br/>stringRuleMap"]
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
    V --> STAGS
    V --> TRANS
    V --> OPTS
    V --> ERRORS
    V --> RPARSE
    V --> RFIELD
    V --> RTIME

    TAGS --> FORMAT
    TAGS --> COMPARE
    TAGS --> ENUM
    STAGS --> FORMAT

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
    style STAGS fill:#ffe0b2
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

## ⚡ VarString 零反射快速路径

对于字符串变量校验场景，`VarString` 提供完全绕过 `reflect` 的零分配快速路径：

```go
v := validator.New()

// 传统 Var 路径 — 通过 interface{} 装箱 + reflect
err := v.Var("user@example.com", "email")

// VarString 零反射路径 — 直接 string 参数，0 堆分配
err = v.VarString("user@example.com", "email")
```

**工作原理：**

- `VarString` 查找 `stringRuleMap`（所有字符串兼容规则的零反射实现），直接以 `string` 参数调用规则函数
- 不支持的规则（如跨字段规则 `eqfield`、`required_if`）自动降级到 reflect 路径，功能完全兼容
- 错误返回轻量级 `stringFieldError`，同样实现 `FieldError` 接口

**支持零反射的规则：**

`required` · `min` · `max` · `len` · `eq` · `ne` · `gt` · `gte` · `lt` · `lte` · `alpha` · `alphanum` · `email` · `url` · `uri` · `ip` · `ipv4` · `ipv6` · `uuid` · `uuid3/4/5` · `semver` · `isbn10/13` · `issn` · `bic` · `cron` · `base64` · `json` · `hostname` · `fqdn` · `mac` · `cidr` · `e164` · `lowercase` · `uppercase` · `boolean` · `number` · `datetime` · `latitude` · `longitude` · `eth_addr` · `btc_addr` · `bcp47` · `datauri` · `oneof` · `oneofci` · `contains` · `startswith` · `endswith` 等 70+ 规则

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
| 零反射字符串校验 | 不支持 | **VarString 0 allocs** |

---

## 🚀 性能基准测试

Argus 与 `go-playground/validator/v10` 的完整性能对比见 [go-argus-benchmark](https://github.com/kamalyes/go-argus-benchmark)。

### VarString 零反射路径 vs Var 反射路径

| 规则 | VarString (零反射) | Var (反射) | VarString 加速 |
|------|-------------------|-----------|---------------|
| `required` | **18 ns** / 0 B / 0 allocs | 49 ns / 16 B / 1 alloc | **2.7×** |
| `email` | **47 ns** / 0 B / 0 allocs | 81 ns / 16 B / 1 alloc | **1.7×** |
| `url` | **37 ns** / 0 B / 0 allocs | 64 ns / 16 B / 1 alloc | **1.7×** |
| `semver` | **28 ns** / 0 B / 0 allocs | 57 ns / 16 B / 1 alloc | **2.0×** |
| `isbn10` | **25 ns** / 0 B / 0 allocs | 60 ns / 16 B / 1 alloc | **2.4×** |
| `cron` | **44 ns** / 0 B / 0 allocs | 74 ns / 16 B / 1 alloc | **1.7×** |

### Argus vs go-playground/validator/v10

| 场景 | Argus | validator/v10 | 优势 |
|------|------:|--------------:|:----:|
| `Var_Email_Valid` | **87 ns** / 0 B / 0 allocs | 626 ns / 98 B / 5 allocs | 🚀 **7.2×** |
| `NestedWorkspace_Valid_Parallel` | **171 ns** / 192 B / 5 allocs | 768 ns / 1007 B / 33 allocs | 🚀 **4.5×** |
| `NestedWorkspace_Valid` | **1014 ns** / 192 B / 5 allocs | 3249 ns / 992 B / 33 allocs | 🚀 **3.2×** |
| `SimpleUser_Valid` | **341 ns** / 0 B / 0 allocs | 810 ns / 98 B / 5 allocs | 🚀 **2.4×** |

> 主要优化手段：零反射 VarString 快速路径、手写 email 解析器替代 `net/mail`、预编译规则分发表、`sync.Pool` 错误对象复用、零分配 `isEmptyValue`、零分配 lowercase/uppercase 字节检查、`json.NewDecoder` 替代 `json.Valid`、轻量 URL/URI 解析替代 `net/url` 等。详见 [go-argus-benchmark](https://github.com/kamalyes/go-argus-benchmark)。

---

## 🔗 生态项目

| 项目 | 说明 |
|------|------|
| [go-rpc-gateway](https://github.com/kamalyes/go-rpc-gateway) | 新一代企业级微服务网关框架，内置 Argus 作为 gRPC/HTTP 参数校验引擎 |
| [go-pbmo](https://github.com/kamalyes/go-pbmo) | 高性能 Protocol Buffer ↔ Model 双向转换库，集成 Argus 字段级参数校验 |
| [go-sqlbuilder](https://github.com/kamalyes/go-sqlbuilder) | 泛型 GORM 仓储层封装，使用 Argus 校验查询参数与模型字段 |
| [go-toolbox](https://github.com/kamalyes/go-toolbox) | 零依赖高性能 Go 工具库，集成 Argus 校验 HTTP 参数、字符串与数据结构 |

### go-rpc-gateway 集成

[go-rpc-gateway](https://github.com/kamalyes/go-rpc-gateway) 开箱即集成了 Argus，提供基于 struct tag 的 gRPC 拦截器，配合 `protoc-go-inject-tag` 在 pb 生成代码上注入 `validate:"..."` 标签，无需业务方手写参数校验：

```go
import "github.com/kamalyes/go-rpc-gateway/middleware"

// gRPC Unary 拦截器 — 自动校验 req 上的 validate 标签
unary := middleware.StructTagValidatorUnaryInterceptor()

// gRPC Stream 拦截器 — 自动校验每条流消息
stream := middleware.StructTagValidatorStreamInterceptor()
```

校验失败自动返回 `codes.InvalidArgument`，字段未注入标签则跳过，不产生误报。

---

## 📄 License

[MIT License](LICENSE)
