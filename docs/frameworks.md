# HTTP 框架集成

Argus 与具体 Web 框架完全解耦：校验入口只是一个 `validator.New()` 创建的 `*Validate` 实例。主流框架都预留了校验器注入点，替换后既有绑定流程自动触发 Argus 校验，业务代码零侵入。

| 框架 | 注入点 | 生效范围 |
|------|--------|----------|
| gin | `binding.Validator` | `ShouldBindJSON` / `ShouldBindQuery` / `ShouldBindUri` 等全部 `ShouldBind*` |
| Echo | `e.Validator` | `c.Bind` 之后调用 `c.Validate` 触发 |
| Fiber | `fiber.Config.StructValidator`（v3+） | `c.Bind().Body / Query / URI` 解析后自动调用 |
| go-zero | `httpx.SetValidator`（v1.4.2+） | `httpx.Parse` 解析路径参数、query、header、JSON body 后自动调用 |
| chi 等轻量路由 | 无内置校验器 | handler 内手动 `v.Struct`，适用于任何框架 |

> 以下示例均在 gin v1.12、Echo v4.15、Fiber v3.5、go-zero v1.10 上实测通过，可直接复制运行。

***

## 1. gin 集成

### 1.1 接入原理

gin 的参数绑定统一经过 `binding.Validator` 全局接口，只要实现 `binding.StructValidator` 并替换默认实例，所有 `ShouldBind*` 就会自动走 Argus，存量 `validate` 标签语法兼容、无需改动：

```go
type StructValidator interface {
    ValidateStruct(interface{}) error
    Engine() interface{}
}
```

### 1.2 完整示例

```go
package main

import (
    "net/http"
    "reflect"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    validator "github.com/kamalyes/go-argus"
)

// LoginReq 登录入参，validate 标签语法与 go-playground 完全一致
type LoginReq struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// LoginResp 登录出参
type LoginResp struct {
    Token string `json:"token"`
}

// ErrorResponse 统一错误响应体
type ErrorResponse struct {
    Code    int                           `json:"code"`
    Message string                        `json:"message"`
    Errors  []validator.ValidationMessage `json:"errors"`
}

// argusValidator 适配 gin 的 binding.StructValidator 接口
type argusValidator struct {
    v *validator.Validate
}

// ValidateStruct 所有 ShouldBind* 绑定完成后自动调用
func (a *argusValidator) ValidateStruct(obj any) error {
    return a.v.Struct(obj)
}

// Engine 返回底层校验器实例
func (a *argusValidator) Engine() any {
    return a.v
}

// tagName 统一从 json/form/path/uri/header 标签提取字段名，
// 校验错误的 field 直接对齐 API 出入参命名
func tagName(sf reflect.StructField) string {
    for _, key := range []string{"json", "form", "path", "uri", "header"} {
        name := strings.SplitN(sf.Tag.Get(key), ",", 2)[0]
        if name != "" && name != "-" {
            return name
        }
    }
    return sf.Name
}

func init() {
    v := validator.New(validator.WithRequiredStructEnabled())
    v.RegisterTagNameFunc(tagName)
    validator.SetLocale("zh")
    // 全局替换 gin 默认校验器，须在服务启动前完成
    binding.Validator = &argusValidator{v: v}
}

func login(c *gin.Context) {
    var req LoginReq
    // 绑定 + 校验一步完成，校验失败时 err 即 validator.ValidationErrors
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "参数校验失败",
            Errors:  validator.TranslateValidationErrors(err, "zh"),
        })
        return
    }
    c.JSON(http.StatusOK, LoginResp{Token: "example-token"})
}

func main() {
    r := gin.Default()
    r.POST("/api/login", login)
    _ = r.Run(":8080")
}
```

query 参数同样生效，例如 `GET /api/users?page=0&page_size=200`：

```go
// ListUsersReq 列表查询入参
type ListUsersReq struct {
    Page     int    `form:"page" validate:"gte=1"`
    PageSize int    `form:"page_size" validate:"gte=1,lte=100"`
    Keyword  string `form:"keyword" validate:"max=64"`
}

func listUsers(c *gin.Context) {
    var req ListUsersReq
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "参数校验失败",
            Errors:  validator.TranslateValidationErrors(err, "zh"),
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": []string{}})
}
```

### 1.3 效果演示

```bash
curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"bad","password":"123"}'
```

```json
{
  "code": 400,
  "message": "参数校验失败",
  "errors": [
    {
      "field": "email",
      "namespace": "LoginReq.email",
      "struct_field": "Email",
      "struct_namespace": "LoginReq.Email",
      "tag": "email",
      "actual_tag": "email",
      "value": "bad",
      "message": "email 必须是有效的 Email"
    },
    {
      "field": "password",
      "namespace": "LoginReq.password",
      "struct_field": "Password",
      "struct_namespace": "LoginReq.Password",
      "tag": "min",
      "actual_tag": "min",
      "param": "8",
      "value": "123",
      "message": "password 不能小于 8"
    }
  ]
}
```

***

## 2. Echo 集成

### 2.1 接入原理

Echo 通过 `e.Validator` 注入校验器（`echo.Validator` 接口只有一个 `Validate(i any) error` 方法）。注意 Echo v4.15 起 `c.Bind` 只负责绑定，需在其后显式调用 `c.Validate` 触发校验：

```go
type Validator interface {
    Validate(i interface{}) error
}
```

### 2.2 完整示例

```go
package main

import (
    "net/http"
    "reflect"
    "strings"

    "github.com/labstack/echo/v4"
    validator "github.com/kamalyes/go-argus"
)

// LoginReq 登录入参
type LoginReq struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// LoginResp 登录出参
type LoginResp struct {
    Token string `json:"token"`
}

// ErrorResponse 统一错误响应体
type ErrorResponse struct {
    Code    int                           `json:"code"`
    Message string                        `json:"message"`
    Errors  []validator.ValidationMessage `json:"errors"`
}

// argusValidator 适配 echo 的 echo.Validator 接口
type argusValidator struct {
    v *validator.Validate
}

// Validate 由 c.Validate 触发
func (a *argusValidator) Validate(i any) error {
    return a.v.Struct(i)
}

// tagName 统一从 json/form/query/param 标签提取字段名
func tagName(sf reflect.StructField) string {
    for _, key := range []string{"json", "form", "query", "param"} {
        name := strings.SplitN(sf.Tag.Get(key), ",", 2)[0]
        if name != "" && name != "-" {
            return name
        }
    }
    return sf.Name
}

func main() {
    e := echo.New()

    v := validator.New(validator.WithRequiredStructEnabled())
    v.RegisterTagNameFunc(tagName)
    validator.SetLocale("zh")
    e.Validator = &argusValidator{v: v}

    e.POST("/api/login", func(c echo.Context) error {
        var req LoginReq
        if err := c.Bind(&req); err != nil {
            return c.JSON(http.StatusBadRequest, ErrorResponse{
                Code:    http.StatusBadRequest,
                Message: "参数校验失败",
                Errors:  validator.TranslateValidationErrors(err, "zh"),
            })
        }
        // 绑定后显式触发校验，内部调用 e.Validator.Validate
        if err := c.Validate(&req); err != nil {
            return c.JSON(http.StatusBadRequest, ErrorResponse{
                Code:    http.StatusBadRequest,
                Message: "参数校验失败",
                Errors:  validator.TranslateValidationErrors(err, "zh"),
            })
        }
        return c.JSON(http.StatusOK, LoginResp{Token: "example-token"})
    })

    e.Logger.Fatal(e.Start(":8080"))
}
```

### 2.3 效果演示

与 gin 章节的响应结构完全一致（同一个 `ErrorResponse`），略。

***

## 3. Fiber 集成

### 3.1 接入原理

Fiber v3 在 `fiber.Config` 中提供 `StructValidator` 字段，注入后 `c.Bind().Body / Query / URI` 解析完请求会自动调用校验器（v2 无此能力，需升级 v3）：

```go
type StructValidator interface {
    Validate(out any) error
}
```

### 3.2 完整示例

```go
package main

import (
    "reflect"
    "strings"

    "github.com/gofiber/fiber/v3"
    validator "github.com/kamalyes/go-argus"
)

// LoginReq 登录入参
type LoginReq struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// LoginResp 登录出参
type LoginResp struct {
    Token string `json:"token"`
}

// ErrorResponse 统一错误响应体
type ErrorResponse struct {
    Code    int                           `json:"code"`
    Message string                        `json:"message"`
    Errors  []validator.ValidationMessage `json:"errors"`
}

// argusValidator 适配 fiber 的 StructValidator 接口
type argusValidator struct {
    v *validator.Validate
}

// Validate 在 Body/Query/URI Parser 解析完请求后自动调用
func (a *argusValidator) Validate(out any) error {
    return a.v.Struct(out)
}

// tagName 统一从 json/form/query/params 标签提取字段名
func tagName(sf reflect.StructField) string {
    for _, key := range []string{"json", "form", "query", "params"} {
        name := strings.SplitN(sf.Tag.Get(key), ",", 2)[0]
        if name != "" && name != "-" {
            return name
        }
    }
    return sf.Name
}

func main() {
    v := validator.New(validator.WithRequiredStructEnabled())
    v.RegisterTagNameFunc(tagName)
    validator.SetLocale("zh")

    app := fiber.New(fiber.Config{
        // 注入 Argus，Body/Query/URI Parser 解析完请求后自动调用
        StructValidator: &argusValidator{v: v},
    })

    app.Post("/api/login", func(c fiber.Ctx) error {
        var req LoginReq
        // 绑定 + 校验一步完成，校验失败时 err 即 validator.ValidationErrors
        if err := c.Bind().Body(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
                Code:    fiber.StatusBadRequest,
                Message: "参数校验失败",
                Errors:  validator.TranslateValidationErrors(err, "zh"),
            })
        }
        return c.Status(fiber.StatusOK).JSON(LoginResp{Token: "example-token"})
    })

    _ = app.Listen(":8080")
}
```

### 3.3 效果演示

与 gin 章节的响应结构完全一致（同一个 `ErrorResponse`），略。

***

## 4. go-zero 集成

### 4.1 接入原理

go-zero（v1.4.2+）在 `httpx` 包预留了校验器接口，goctl 生成的 handler 在 `httpx.Parse` 解析完路径参数、query、header 和 JSON body 之后自动调用注入的校验器：

```go
type Validator interface {
    Validate(r *http.Request, data any) error
}

httpx.SetValidator(v) // 全局注入
```

### 4.2 声明校验标签

`validate` 标签直接写在 api 文件的字段上，goctl 会原样透传到生成的 `internal/types/types.go`，重新生成代码不会丢失。

注意 json tag 要统一加 `,optional`：go-zero 的 mapping 会把缺失的非 optional 字段直接报 `field "email" is not set` 而跳过校验器，把必填语义完全交给 `validate` 标签即可获得结构化错误输出：

```api
syntax = "v1"

type (
    RegisterReq {
        Username string `json:"username,optional" validate:"required,min=2,max=32"`
        Email    string `json:"email,optional" validate:"required,email"`
        Password string `json:"password,optional" validate:"required,min=8"`
    }

    RegisterResp {
        Id       int64  `json:"id"`
        Username string `json:"username"`
    }
)

@server (
    prefix: /api
)
service user-api {
    @handler register
    post /user/register (RegisterReq) returns (RegisterResp)
}
```

### 4.3 实现校验器

```go
package middleware

import (
    "net/http"
    "reflect"
    "strings"

    validator "github.com/kamalyes/go-argus"
)

// ArgusValidator 实现 go-zero 的 httpx.Validator 接口
type ArgusValidator struct {
    v *validator.Validate
}

// NewArgusValidator 创建全局可复用的校验器实例
func NewArgusValidator() *ArgusValidator {
    v := validator.New(validator.WithRequiredStructEnabled())
    v.RegisterTagNameFunc(tagName)
    return &ArgusValidator{v: v}
}

// tagName 统一从 json/path/form/header 标签提取字段名
func tagName(sf reflect.StructField) string {
    for _, key := range []string{"json", "path", "form", "header"} {
        name := strings.SplitN(sf.Tag.Get(key), ",", 2)[0]
        if name != "" && name != "-" {
            return name
        }
    }
    return sf.Name
}

// Validate 在 httpx.Parse 解析完请求后自动调用
func (a *ArgusValidator) Validate(r *http.Request, data any) error {
    if data == nil {
        return nil
    }
    return a.v.Struct(data)
}
```

### 4.4 注入与统一错误响应

goctl 生成的 handler 无需任何改动，校验失败会走到 `httpx.ErrorCtx`，配合 `SetErrorHandlerCtx` 统一输出结构化 JSON：

```go
// internal/handler/registerhandler.go — goctl 生成，无需修改
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.RegisterReq
        if err := httpx.Parse(r, &req); err != nil {
            httpx.ErrorCtx(r.Context(), w, err) // 校验失败在这里被统一转换
            return
        }
        // ...
    }
}
```

在 `main.go` 中注入：

```go
package main

import (
    "context"
    "errors"
    "flag"
    "net/http"

    "github.com/zeromicro/go-zero/conf"
    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/rest/httpx"

    validator "github.com/kamalyes/go-argus"

    "user/internal/config"
    "user/internal/handler"
    "user/internal/middleware"
    "user/internal/svc"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

// ErrorResponse 统一错误响应体
type ErrorResponse struct {
    Code    int                           `json:"code"`
    Message string                        `json:"message"`
    Errors  []validator.ValidationMessage `json:"errors"`
}

func main() {
    flag.Parse()

    var c config.Config
    conf.MustLoad(*configFile, &c)

    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()

    ctx := svc.NewServiceContext(c)
    handler.RegisterHandlers(server, ctx)

    // 注入 Argus 校验器，httpx.Parse 解析完请求后自动调用
    httpx.SetValidator(middleware.NewArgusValidator())

    // 校验失败统一返回结构化 JSON，其他错误保持 400 + 原始信息
    httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
        var ve validator.ValidationErrors
        if errors.As(err, &ve) {
            return http.StatusBadRequest, ErrorResponse{
                Code:    http.StatusBadRequest,
                Message: "参数校验失败",
                Errors:  validator.TranslateValidationErrors(err, "zh"),
            }
        }
        return http.StatusBadRequest, ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: err.Error(),
        }
    })

    server.Start()
}
```

### 4.5 效果演示

```bash
curl -s -X POST http://localhost:8888/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"k"}'
```

```json
{
  "code": 400,
  "message": "参数校验失败",
  "errors": [
    {
      "field": "username",
      "namespace": "RegisterReq.username",
      "struct_field": "Username",
      "struct_namespace": "RegisterReq.Username",
      "tag": "min",
      "actual_tag": "min",
      "param": "2",
      "value": "k",
      "message": "username 不能小于 2"
    },
    {
      "field": "email",
      "namespace": "RegisterReq.email",
      "struct_field": "Email",
      "struct_namespace": "RegisterReq.Email",
      "tag": "required",
      "actual_tag": "required",
      "value": "",
      "message": "email 为必填字段"
    },
    {
      "field": "password",
      "namespace": "RegisterReq.password",
      "struct_field": "Password",
      "struct_namespace": "RegisterReq.Password",
      "tag": "required",
      "actual_tag": "required",
      "value": "",
      "message": "password 为必填字段"
    }
  ]
}
```

***

## 5. chi 与通用模式

chi、标准库 `net/http` 等路由本身不带绑定与校验，解码完成后手动调用 `v.Struct` 即可。这也是接入任何未在本文档列出的框架的通用模式：

```go
package main

import (
    "encoding/json"
    "net/http"
    "reflect"
    "strings"

    "github.com/go-chi/chi/v5"
    validator "github.com/kamalyes/go-argus"
)

// LoginReq 登录入参
type LoginReq struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// LoginResp 登录出参
type LoginResp struct {
    Token string `json:"token"`
}

// ErrorResponse 统一错误响应体
type ErrorResponse struct {
    Code    int                           `json:"code"`
    Message string                        `json:"message"`
    Errors  []validator.ValidationMessage `json:"errors"`
}

var v = newValidator()

// newValidator 路由层没有注入点时，自行创建并全局复用
func newValidator() *validator.Validate {
    v := validator.New(validator.WithRequiredStructEnabled())
    v.RegisterTagNameFunc(tagName)
    validator.SetLocale("zh")
    return v
}

// tagName 统一从 json 标签提取字段名
func tagName(sf reflect.StructField) string {
    for _, key := range []string{"json"} {
        name := strings.SplitN(sf.Tag.Get(key), ",", 2)[0]
        if name != "" && name != "-" {
            return name
        }
    }
    return sf.Name
}

func writeJSON(w http.ResponseWriter, status int, body any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(body)
}

func login(w http.ResponseWriter, r *http.Request) {
    var req LoginReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "参数校验失败",
            Errors:  validator.TranslateValidationErrors(err, "zh"),
        })
        return
    }
    // 解码完成后手动触发校验
    if err := v.Struct(req); err != nil {
        writeJSON(w, http.StatusBadRequest, ErrorResponse{
            Code:    http.StatusBadRequest,
            Message: "参数校验失败",
            Errors:  validator.TranslateValidationErrors(err, "zh"),
        })
        return
    }
    writeJSON(w, http.StatusOK, LoginResp{Token: "example-token"})
}

func main() {
    r := chi.NewRouter()
    r.Post("/api/login", login)
    _ = http.ListenAndServe(":8080", r)
}
```

***

## 6. 进阶用法

### 6.1 按请求语言动态翻译

gin 中可在 handler 直接读取 `Accept-Language`：

```go
// localeFromHeader 从 Accept-Language 提取语言标签，9 种内置语言见 docs/i18n.md
func localeFromHeader(c *gin.Context) string {
    locale := c.GetHeader("Accept-Language")
    if len(locale) >= 2 {
        locale = locale[:2]
    }
    switch locale {
    case "en", "zh", "ja", "ko", "fr", "de", "es", "ru":
        return locale
    default:
        return "zh"
    }
}

c.JSON(http.StatusBadRequest, ErrorResponse{
    Code:    http.StatusBadRequest,
    Message: "参数校验失败",
    Errors:  validator.TranslateValidationErrors(err, localeFromHeader(c)),
})
```

```bash
curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d '{"email":"bad","password":"123"}'
```

```json
{
  "errors": [
    { "field": "email", "tag": "email", "value": "bad", "message": "email must be a valid email address" },
    { "field": "password", "tag": "min", "param": "8", "value": "123", "message": "password must be at least 8" }
  ]
}
```

go-zero 的错误处理回调只拿到 `ctx`，若需按请求语言翻译，可在 `Validate` 内翻译完成后返回携带消息的自定义错误，错误处理回调中直接输出：

```go
// TranslatedError 携带已翻译消息的错误，供统一错误处理直接输出
type TranslatedError struct {
    Messages []validator.ValidationMessage
}

func (e *TranslatedError) Error() string { return "参数校验失败" }

func (a *ArgusValidator) Validate(r *http.Request, data any) error {
    if data == nil {
        return nil
    }
    if err := a.v.Struct(data); err != nil {
        return &TranslatedError{Messages: validator.TranslateValidationErrors(err, localeFromHeader(r))}
    }
    return nil
}
```

### 6.2 提取缺失字段列表

区分"缺参"与"格式错误"是常见需求，例如缺参直接打点上报、格式错误提示用户修改。`ValidationErrors` 提供了两个开箱即用的方法：

```go
if ve, ok := err.(validator.ValidationErrors); ok {
    // 只返回 required 系列规则失败的错误
    missing := ve.RequiredMessages("zh")

    // 只返回 required 系列规则失败的字段名
    fields := ve.MissingFields() // ["email", "password"]
}
```

***

## 7. 常见问题

| 问题 | 说明 |
|------|------|
| gin 还需要 go-playground/validator 吗 | gin 模块自带该依赖，但替换 `binding.Validator` 后所有校验都走 Argus，`validate` 标签语法兼容、无需改动 |
| 替换校验器的时机 | gin 在 `main`/`init` 中设置 `binding.Validator`，Echo 在 `e.Validator` 赋值，Fiber 在 `fiber.New` 的 Config 中传入，go-zero 在 `rest.MustNewServer` 之后、`server.Start` 之前调用 `httpx.SetValidator`，都必须在服务启动前完成 |
| Echo 的 `c.Bind` 没触发校验 | Echo v4.15 起 `c.Bind` 只负责绑定，需在其后显式调用 `c.Validate(&req)` |
| go-zero 缺字段报 `field "email" is not set` | go-zero 的 mapping 把缺失的非 optional json 字段直接报错而跳过校验器，api 文件的 json tag 统一加 `,optional`，必填语义交给 `validate` 标签 |
| go-zero 版本要求 | `httpx.SetValidator` 需 go-zero v1.4.2+ |
| Fiber 版本要求 | `fiber.Config.StructValidator` 为 v3 特性，v2 需升级 |
| 重新生成代码会丢 validate 标签吗 | 不会，标签声明在 api 文件字段上，goctl 原样透传到生成的 types |
| 指针字段零值是否触发 required | 指针 `nil` 天然触发；非指针结构体零值需开启 `WithRequiredStructEnabled` |
| 非 ValidationErrors 错误怎么处理 | `TranslateValidationErrors` 会把 JSON 解码错误等包成单条 `message` 输出，无需额外分支 |
