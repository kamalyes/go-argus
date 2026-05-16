# 使用示例

## 安装

```bash
go get github.com/kamalyes/go-argus
```

***

## 1. 结构体标签校验

最核心的用法，通过 `validate` 标签声明规则：

```go
package main

import (
    "fmt"
    "github.com/kamalyes/go-argus"
)

type User struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"gte=0,lte=150"`
    Password string `json:"password" validate:"required,min=8"`
}

func main() {
    v := validator.New()

    user := User{
        Name:     "A",
        Email:    "invalid-email",
        Age:      -1,
        Password: "123",
    }

    err := v.Struct(user)
    if err != nil {
        for _, fe := range err.(validator.ValidationErrors) {
            fmt.Printf("field=%s tag=%s param=%s value=%v\n",
                fe.Field(), fe.Tag(), fe.Param(), fe.Value())
        }
    }
}
```

输出：

```
field=name tag=min param=2 value=A
field=email tag=email param= value=invalid-email
field=age tag=gte param=0 value=-1
field=password tag=min param=8 value=123
```

***

## 2. i18n 错误消息

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/kamalyes/go-argus"
)

type User struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
}

func main() {
    v := validator.New()
    err := v.Struct(User{})

    // 中文翻译
    messages := validator.TranslateValidationErrors(err, "zh")
    out, _ := json.MarshalIndent(messages, "", "  ")
    fmt.Println(string(out))
}
```

输出：

```json
[
  {
    "field": "name",
    "namespace": "User.name",
    "tag": "required",
    "message": "name 为必填字段"
  },
  {
    "field": "email",
    "namespace": "User.email",
    "tag": "required",
    "message": "email 为必填字段"
  }
]
```

***

## 3. 单变量校验

不需要定义结构体，直接校验单个值：

```go
v := validator.New()

err := v.Var("test@example.com", "required,email")
fmt.Println(err) // <nil>

err = v.Var("", "required,email")
fmt.Println(err) // Key: '' Error:Field validation for '' failed on the 'required' tag
```

***

## 4. 条件必填

```go
type Address struct {
    Country string `json:"country" validate:"required"`
    State   string `json:"state" validate:"required_if=Country US"`
    ZipCode string `json:"zip_code" validate:"required_with=State"`
}

v := validator.New()
err := v.Struct(Address{Country: "US"})
// State 和 ZipCode 都会报 required 错误
```

***

## 5. 跨字段比较

```go
type TimeRange struct {
    Start string `json:"start" validate:"required,datetime=2006-01-02"`
    End   string `json:"end" validate:"required,datetime=2006-01-02,gtfield=Start"`
}

type PasswordChange struct {
    NewPassword     string `json:"new_password" validate:"required,min=8"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
```

***

## 6. 切片元素逐一校验

使用 `dive` 标签对切片/Map 的每个元素进行校验：

```go
type Team struct {
    Members []string `json:"members" validate:"required,dive,required,email"`
}

v := validator.New()
err := v.Struct(Team{
    Members: []string{"ok@test.com", "bad-email", ""},
})
// "bad-email" 和 "" 都会报错
```

***

## 7. 自定义校验规则

```go
import "github.com/kamalyes/go-argus"

v := validator.New()

// 注册自定义规则
v.RegisterValidation("even", func(fl validator.FieldLevel) bool {
    n, ok := fl.Field().Interface().(int)
    return ok && n%2 == 0
})

type Config struct {
    Count int `json:"count" validate:"required,even"`
}

err := v.Struct(Config{Count: 3})
fmt.Println(err) // Count failed on 'even' tag
```

***

## 8. 使用 json tag 作为字段名

```go
v := validator.New()
v.RegisterTagNameFunc(func(sf reflect.StructField) string {
    name := strings.SplitN(sf.Tag.Get("json"), ",", 2)[0]
    if name == "" || name == "-" {
        return sf.Name
    }
    return name
})
```

***

## 9. 数值与字符串比较

```go
import "github.com/kamalyes/go-argus"

// 数值比较
result := validator.CompareNumbers(10, 5, validator.OpGreaterThan)
fmt.Println(result.Success) // true

// 字符串比较
result = validator.CompareStrings("hello", "world", validator.OpContains)
fmt.Println(result.Success) // false

// HTTP 状态码范围校验
result = validator.ValidateStatusCodeRange(200, 200, 299)
fmt.Println(result.Success) // true
```

***

## 10. IP 白名单 / 黑名单

```go
import "github.com/kamalyes/go-argus"

// 预编译 IP 规则（适合高频场景）
allowSet, _ := validator.CompileIPSet([]string{
    "10.0.0.0/8",
    "192.168.1.0/24",
    "172.16.*",
})

allowSet.Contains("10.1.2.3")     // true
allowSet.Contains("8.8.8.8")      // false

// 简单判断
validator.IsIPAllowed("10.1.2.3", []string{"10.0.0.0/8"})  // true
validator.IsIPBlocked("8.8.8.8", []string{"8.8.8.0/24"})    // true
validator.IsPrivateIP("192.168.1.1")                         // true
```

***

## 11. JSON Schema 校验

```go
import "github.com/kamalyes/go-argus/schema"

schemaDef := schema.JSONSchema{
    Type:     "object",
    Required: []string{"name", "age"},
    Properties: map[string]schema.JSONSchema{
        "name": {Type: "string", MinLength: intPtr(1), MaxLength: intPtr(100)},
        "age":  {Type: "integer", Minimum: floatPtr(0), Maximum: floatPtr(150)},
    },
}

data := map[string]interface{}{
    "name": "Argus",
    "age":  float64(200), // 超过最大值
}

result := schema.ValidateJSONSchema(data, schemaDef)
fmt.Println(result.Success) // false
fmt.Println(result.Message) // $.age above maximum 150

func intPtr(v int) *int       { return &v }
func floatPtr(v float64) *float64 { return &v }
```

***

## 12. 枚举校验

```go
import "github.com/kamalyes/go-argus"

roles := validator.NewEnumValidator("admin", "editor", "viewer")

err := roles.MustBeValid("admin")  // nil
err = roles.MustBeValid("hacker")  // invalid enum value: hacker
```

***

## 13. 全局 i18n 设置

```go
import "github.com/kamalyes/go-argus"

// 设置全局语言（同时影响根包和子包）
validator.SetLocale("zh")

// 注册新语言（9 种内置语言：en/zh/zh-TW/ja/ko/fr/de/es/ru）
validator.RegisterI18nMessages("pt", map[string]string{
    "required": "{field} é obrigatório",
    "email":    "{field} deve ser um e-mail válido",
})

// 切换到葡萄牙语
validator.SetLocale("pt")

// 校验器内部消息也会自动跟随语言
validator.SetLocale("ja")
result := validator.ValidateEmail("")
fmt.Println(result.Message) // メールアドレスが空です
```

***

## 14. gRPC 拦截器集成（配合 protoc-go-inject-tag）

在 gRPC 微服务中，通过 `protoc-go-inject-tag` 在 pb 生成代码字段上注入 `validate:"..."` 标签，配合 Argus 校验拦截器，无需在 service 层手写参数校验：

**Proto 文件声明 inject-tag：**

```protobuf
// announcement.proto
message CreateAnnouncementRequest {
  string announcement_id = 1;                       // @inject_tag: validate:"required"
  google.protobuf.Timestamp start_time = 2;         // @inject_tag: validate:"required"
  google.protobuf.Timestamp end_time = 3;           // @inject_tag: validate:"required"
  AnnouncementScope scope = 4;                      // @inject_tag: validate:"required"
  repeated AnnouncementContent contents = 7;        // @inject_tag: validate:"required,dive,required"
}
```

**生成的 Go 代码（自动注入标签）：**

```go
type CreateAnnouncementRequest struct {
    AnnouncementId string                 `protobuf:"..." json:"announcement_id,omitempty" validate:"required"`
    StartTime      *timestamppb.Timestamp `protobuf:"..." json:"start_time,omitempty" validate:"required"`
    EndTime        *timestamppb.Timestamp `protobuf:"..." json:"end_time,omitempty" validate:"required"`
    Scope          AnnouncementScope      `protobuf:"..." json:"scope,omitempty" validate:"required"`
    Contents       []*AnnouncementContent `protobuf:"..." json:"contents,omitempty" validate:"required,dive,required"`
}
```

**注册拦截器：**

```go
package middleware

import (
    "context"
    "strings"
    "sync"

    validator "github.com/kamalyes/go-argus"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

var (
    vOnce sync.Once
    v     *validator.Validate
)

func getValidator() *validator.Validate {
    vOnce.Do(func() {
        v = validator.New(validator.WithRequiredStructEnabled())
    })
    return v
}

func StructTagValidatorUnaryInterceptor() grpc.UnaryServerInterceptor {
    v := getValidator()
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        if req == nil {
            return handler(ctx, req)
        }
        if err := v.Struct(req); err != nil {
            return nil, status.Error(codes.InvalidArgument, formatError(err))
        }
        return handler(ctx, req)
    }
}

func formatError(err error) string {
    ve, ok := err.(validator.ValidationErrors)
    if !ok || len(ve) == 0 {
        return err.Error()
    }
    parts := make([]string, 0, len(ve))
    for _, fe := range ve {
        parts = append(parts, fe.Namespace()+": "+fe.Tag())
    }
    return "invalid argument: " + strings.Join(parts, ", ")
}
```

**在 gRPC Server 中使用：**

```go
s := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        middleware.StructTagValidatorUnaryInterceptor(),
    ),
)
```

***

## 15. Protobuf 多规则注入实战

微服务中常见的多规则组合注入案例：

```protobuf
message CreateWorkspaceRequest {
  string code = 1;                    // 工作空间编码 @inject_tag: validate:"required,alphanum,min=2,max=20"
  string display_name = 2;            // 显示名称 @inject_tag: validate:"required,min=2,max=100"
  string admin_email = 3;             // 管理员邮箱 @inject_tag: validate:"required,email"
  bytes encrypted_secret = 4;         // RSA加密密钥 @inject_tag: validate:"required,min=128"
  string infra_config_id = 9;         // 基础设施配置ID @inject_tag: validate:"required,uuid"
  string addon_plan_id = 10;          // 增值套餐ID @inject_tag: validate:"omitempty,uuid"
}

message UpdateEnvironmentRequest {
  string env_id = 1;                  // 环境ID @inject_tag: validate:"required,uuid"
  google.protobuf.StringValue display_name = 3;  // 环境名称 @inject_tag: validate:"omitnil,required,min=1,max=50"
  google.protobuf.StringValue region = 5;        // 区域 @inject_tag: validate:"omitnil,required,oneof=US EU AP SA"
  enums.EnvPhase phase = 3;           // 环境阶段 @inject_tag: validate:"required,oneof=1 2 3"
}

message ListReleaseRequest {
  common.Paging page_request = 1;     // 分页 @inject_tag: validate:"required"
  string release_tag = 2;             // 发布标签（可选） @inject_tag: validate:"omitempty,min=1,max=50"
}
```

***

## 16. 时间范围校验

公告、活动、预约等业务中常见的时间区间校验：

```go
type CreateEventRequest struct {
    StartTime time.Time `json:"start_time" validate:"required"`
    EndTime   time.Time `json:"end_time" validate:"required,afterfield=StartTime"`
}

type ScheduleEvent struct {
    Event CreateEventRequest `json:"event" validate:"range=StartTime|EndTime"`
}

v := validator.New(validator.WithRequiredStructEnabled())

err := v.Struct(ScheduleEvent{
    Event: CreateEventRequest{
        StartTime: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
        EndTime:   time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
    },
})
// EndTime 早于 StartTime → 校验失败
```

**使用时间表达式校验：**

```go
type Coupon struct {
    Code      string    `json:"code" validate:"required,alphanum,min=4,max=20"`
    ValidFrom time.Time `json:"valid_from" validate:"required,after=now-7d"`
    ValidTo   time.Time `json:"valid_to" validate:"required,after=now+30d"`
    MaxUsage  int       `json:"max_usage" validate:"required,gte=1,lte=100000"`
}

v := validator.New(validator.WithRequiredStructEnabled())
err := v.Struct(Coupon{
    Code:      "SUMMER2026",
    ValidFrom: time.Now().Add(-3 * 24 * time.Hour),
    ValidTo:   time.Now().Add(60 * 24 * time.Hour),
    MaxUsage:  5000,
})
```

***

## 17. 互斥字段与条件排除

管理后台常见场景：某些字段互斥，不能同时填写：

```go
type UpdateConfigRequest struct {
    ID    string `json:"id" validate:"required,uuid"`
    Value string `json:"value" validate:"excluded_with=File"`
    File  string `json:"file" validate:"excluded_with=Value,filepath"`
}

type NotificationSetting struct {
    Channel string `json:"channel" validate:"required,oneof=email sms push"`
    Email   string `json:"email" validate:"required_if=Channel email,email"`
    Phone   string `json:"phone" validate:"required_if=Channel sms,e164"`
    DeviceToken string `json:"device_token" validate:"required_if=Channel push,min=32,max=256"`
    QuietStart string `json:"quiet_start" validate:"omitempty,datetime=15:04"`
    QuietEnd   string `json:"quiet_end" validate:"omitempty,datetime=15:04,gtfield=QuietStart"`
}
```

***

## 18. 嵌套结构体 + dive 深度校验

表单提交、批量导入等复杂嵌套场景：

```go
type AnnouncementContent struct {
    Language string `json:"language" validate:"required,oneof=zh en ja ko"`
    Title    string `json:"title" validate:"required,min=1,max=200"`
    Body     string `json:"body" validate:"required,min=1,max=5000"`
    ImageURL string `json:"image_url" validate:"omitempty,http_url"`
}

type CreateAnnouncementRequest struct {
    ID               string               `json:"id" validate:"required,uuid"`
    StartTime        time.Time            `json:"start_time" validate:"required"`
    EndTime          time.Time            `json:"end_time" validate:"required,afterfield=StartTime"`
    Scope            string               `json:"scope" validate:"required,oneof=all ios android web"`
    DisplayPage      string               `json:"display_page" validate:"required,oneof=home lobby settings"`
    DisplayFrequency string               `json:"display_frequency" validate:"required,oneof=once daily always"`
    Contents         []AnnouncementContent `json:"contents" validate:"required,min=1,dive"`
    ClickAction      string               `json:"click_action" validate:"required,oneof=none url page"`
    ClickURL         string               `json:"click_url" validate:"required_if=ClickAction url,http_url"`
}

v := validator.New(validator.WithRequiredStructEnabled())
err := v.Struct(CreateAnnouncementRequest{
    ID:               "550e8400-e29b-41d4-a716-446655440000",
    StartTime:        time.Now(),
    EndTime:          time.Now().Add(24 * time.Hour),
    Scope:            "all",
    DisplayPage:      "home",
    DisplayFrequency: "once",
    Contents: []AnnouncementContent{
        {Language: "zh", Title: "系统维护通知", Body: "将于今晚进行系统维护", ImageURL: "https://cdn.example.com/maintenance.png"},
        {Language: "en", Title: "", Body: "System maintenance tonight"}, // Title 为空 → 校验失败
    },
    ClickAction: "url",
    ClickURL:    "not-a-url", // http_url 校验失败
})
```

***

## 19. 分页请求通用校验

几乎所有列表接口都需要分页参数校验：

```go
type Paging struct {
    Page     int32 `json:"page" validate:"required,min=1"`
    PageSize int32 `json:"page_size" validate:"required,min=1,max=100"`
}

type ListRequest struct {
    Paging     Paging  `json:"paging" validate:"required"`
    Keyword    string  `json:"keyword" validate:"omitempty,min=1,max=100"`
    SortBy     string  `json:"sort_by" validate:"omitempty,oneof=created_at updated_at name"`
    SortOrder  string  `json:"sort_order" validate:"omitempty,oneof=asc desc"`
    Status     int32   `json:"status" validate:"omitempty,oneof=0 1 2"`
    StartTime  string  `json:"start_time" validate:"omitempty,datetime=2006-01-02"`
    EndTime    string  `json:"end_time" validate:"omitempty,datetime=2006-01-02,gtfield=StartTime"`
}

v := validator.New(validator.WithRequiredStructEnabled())
err := v.Struct(ListRequest{
    Paging:    Paging{Page: 0, PageSize: 200}, // Page < 1, PageSize > 100 → 校验失败
    SortBy:    "invalid_field",
    SortOrder: "random",
})
```

***

## 20. 用户注册/登录表单校验

```go
type RegisterRequest struct {
    Username string `json:"username" validate:"required,alphanumunicode,min=3,max=30"`
    Email    string `json:"email" validate:"required,email"`
    Phone    string `json:"phone" validate:"required_without=Email,e164"`
    Password string `json:"password" validate:"required,min=8,max=128,containsany=!@#$%^&*"`
    Confirm  string `json:"confirm" validate:"required,eqfield=Password"`
    AgreeTOS bool   `json:"agree_tos" validate:"required,eq=true"`
}

type LoginRequest struct {
    Account  string `json:"account" validate:"required"` // 邮箱或手机号
    Password string `json:"password" validate:"required,min=8"`
    Captcha  string `json:"captcha" validate:"omitempty,len=6,numeric"`
}

type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" validate:"required,min=8"`
    NewPassword string `json:"new_password" validate:"required,min=8,nefield=OldPassword"`
    Confirm     string `json:"confirm" validate:"required,eqfield=NewPassword"`
}
```

***

## 21. 文件上传与资源校验

```go
type UploadRequest struct {
    Filename    string `json:"filename" validate:"required,excludesall=/\\,"`
    ContentType string `json:"content_type" validate:"required,oneof=image/jpeg image/png image/gif image/webp application/pdf"`
    MaxSizeMB   int    `json:"max_size_mb" validate:"required,gte=1,lte=50"`
    StoragePath string `json:"storage_path" validate:"required,dirpath"`
}

type CDNResource struct {
    URL        string `json:"url" validate:"required,https_url"`
    CDNPath    string `json:"cdn_path" validate:"required,startswith=/cdn/"`
    Expiration int64  `json:"expiration" validate:"required,gt=0"`
}
```

***

## 22. 网络与安全配置校验

```go
type NetworkConfig struct {
    ListenAddr string `json:"listen_addr" validate:"required,hostname_port"`
    AdminCIDR  string `json:"admin_cidr" validate:"required,cidrv4"`
    TrustedProxies []string `json:"trusted_proxies" validate:"required,dive,cidr"`
    DNS        string `json:"dns" validate:"required,ip"`
    Domain     string `json:"domain" validate:"required,fqdn"`
    SSLEnabled bool   `json:"ssl_enabled" validate:"omitempty"`
    CertPath   string `json:"cert_path" validate:"required_if=SSLEnabled true,file"`
    KeyPath    string `json:"key_path" validate:"required_if=SSLEnabled true,file"`
}

type SecurityPolicy struct {
    MaxLoginAttempts int      `json:"max_login_attempts" validate:"required,gte=1,lte=20"`
    LockoutDuration  int      `json:"lockout_duration" validate:"required,gte=1,lte=1440"` // 分钟
    AllowedIPs       []string `json:"allowed_ips" validate:"omitempty,dive,ip"`
    BlockedCountries []string `json:"blocked_countries" validate:"omitempty,dive,alpha,len=2"`
    TokenExpiry      int      `json:"token_expiry" validate:"required,gte=1,lte=8760"` // 小时
}
```

***

## 23. 电商/支付场景校验

```go
type OrderItem struct {
    ProductID string  `json:"product_id" validate:"required,uuid"`
    SKU       string  `json:"sku" validate:"required,alphanum"`
    Quantity  int     `json:"quantity" validate:"required,gte=1,lte=999"`
    UnitPrice float64 `json:"unit_price" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
    Items       []OrderItem `json:"items" validate:"required,min=1,dive"`
    CouponCode  string      `json:"coupon_code" validate:"omitempty,alphanum,min=3,max=20"`
    ShippingFee float64     `json:"shipping_fee" validate:"required,gte=0"`
    Currency    string      `json:"currency" validate:"required,alpha,len=3"`
    PaymentMethod string    `json:"payment_method" validate:"required,oneof=credit_card paypal bank_transfer crypto"`
    BillingAddress *Address `json:"billing_address" validate:"required"`
}

type Address struct {
    Country string `json:"country" validate:"required,alpha,len=2"`
    State   string `json:"state" validate:"required,min=1,max=100"`
    City    string `json:"city" validate:"required,min=1,max=100"`
    Street  string `json:"street" validate:"required,min=1,max=200"`
    ZipCode string `json:"zip_code" validate:"required,alphanum,min=3,max=10"`
}
```

***

## 24. SaaS 工作空间配置校验

多工作空间/多环境的 SaaS 平台常见校验场景：

```go
type CreateWorkspaceRequest struct {
    Code            string `json:"code" validate:"required,alphanum,min=2,max=20"`
    DisplayName     string `json:"display_name" validate:"required,min=2,max=100"`
    AdminEmail      string `json:"admin_email" validate:"required,email"`
    EncryptedSecret []byte `json:"encrypted_secret" validate:"required,min=128"`
    InfraConfigID   string `json:"infra_config_id" validate:"required,uuid"`
    AddonPlanID     string `json:"addon_plan_id" validate:"omitempty,uuid"`
}

type EnvironmentConfig struct {
    EnvID             string   `json:"env_id" validate:"required,uuid"`
    DisplayName       string   `json:"display_name" validate:"required,min=1,max=50"`
    Region            string   `json:"region" validate:"required,oneof=US EU AP SA ME AF"`
    SupportedLocales  []string `json:"supported_locales" validate:"required,min=1,dive,oneof=zh en ja ko es pt ru ar hi"`
    Theme             string   `json:"theme" validate:"required,oneof=dark light auto"`
    NamespacePrefix   string   `json:"namespace_prefix" validate:"required,dns_rfc1035_label"`
}

type ReleaseCreateRequest struct {
    Channel      string `json:"channel" validate:"required,oneof=stable beta canary"`
    ReleaseTag   string `json:"release_tag" validate:"required,startswith=v,alphanum"`
    ArtifactURL  string `json:"artifact_url" validate:"required,https_url"`
}
```

***

## 25. 批量操作与 ID 列表校验

```go
type BatchDeleteRequest struct {
    IDs []string `json:"ids" validate:"required,min=1,max=100,dive,uuid"`
}

type BatchUpdateStatusRequest struct {
    IDs     []string `json:"ids" validate:"required,min=1,max=100,dive,uuid"`
    Status  int32    `json:"status" validate:"required,oneof=1 2 3"`
    Reason  string   `json:"reason" validate:"required_unless=Status 1,min=5,max=500"`
}

type AssignRolesRequest struct {
    UserID string   `json:"user_id" validate:"required,uuid"`
    Roles  []string `json:"roles" validate:"required,min=1,max=10,dive,oneof=admin editor viewer operator auditor"`
}
```

***

## 26. 条件必填组合场景

多种条件必填规则组合使用，覆盖复杂业务逻辑：

```go
type ShippingRequest struct {
    Method      string  `json:"method" validate:"required,oneof=standard express pickup"`
    Address     string  `json:"address" validate:"required_unless=Method pickup,min=5,max=200"`
    StoreID     string  `json:"store_id" validate:"required_if=Method pickup,uuid"`
    ExpressFee  float64 `json:"express_fee" validate:"required_if=Method express,gt=0"`
    Insured     bool    `json:"insured"`
    InsuredAmount float64 `json:"insured_amount" validate:"required_if=Insured true,gt=0"`
}

type QueryRequest struct {
    Type       string `json:"type" validate:"required,oneof=id keyword date_range"`
    ID         string `json:"id" validate:"required_if=Type id,uuid"`
    Keyword    string `json:"keyword" validate:"required_if=Type keyword,min=1,max=100"`
    StartDate  string `json:"start_date" validate:"required_if=Type date_range,datetime=2006-01-02"`
    EndDate    string `json:"end_date" validate:"required_if=Type date_range,datetime=2006-01-02,gtfield=StartDate"`
}
```

***

## 27. 自定义校验器实战

### 幂等性 Token 校验

```go
v := validator.New()

v.RegisterValidation("idempotency_key", func(fl validator.FieldLevel) bool {
    s, ok := fl.Field().Interface().(string)
    if !ok {
        return false
    }
    return len(s) >= 16 && len(s) <= 128 && isHexString(s)
})

func isHexString(s string) bool {
    for _, c := range s {
        if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
            return false
        }
    }
    return true
}

type CreatePaymentRequest struct {
    IdempotencyKey string  `json:"idempotency_key" validate:"required,idempotency_key"`
    Amount         float64 `json:"amount" validate:"required,gt=0"`
    Currency       string  `json:"currency" validate:"required,alpha,len=3"`
}
```

### 密码强度校验

```go
v := validator.New()

v.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
    s, ok := fl.Field().Interface().(string)
    if !ok || len(s) < 8 {
        return false
    }
    var hasUpper, hasLower, hasDigit, hasSpecial bool
    for _, c := range s {
        switch {
        case unicode.IsUpper(c):
            hasUpper = true
        case unicode.IsLower(c):
            hasLower = true
        case unicode.IsDigit(c):
            hasDigit = true
        case unicode.IsPunct(c) || unicode.IsSymbol(c):
            hasSpecial = true
        }
    }
    return hasUpper && hasLower && hasDigit && hasSpecial
})

type SecureUser struct {
    Username string `json:"username" validate:"required,alphanum,min=3,max=30"`
    Password string `json:"password" validate:"required,strong_password"`
}
```

### 跨服务字段校验

```go
v := validator.New()

v.RegisterValidationCtx("exists_in_service", func(ctx context.Context, fl validator.FieldLevel) bool {
    svc, ok := ctx.Value("userService").(UserService)
    if !ok {
        return false
    }
    email := fl.Field().String()
    return svc.EmailExists(ctx, email)
})

type InviteRequest struct {
    Email string `json:"email" validate:"required,email,exists_in_service"`
    Role  string `json:"role" validate:"required,oneof=admin editor viewer"`
}

// 调用时注入 context
ctx := context.WithValue(context.Background(), "userService", userSvc)
err := v.StructCtx(ctx, InviteRequest{Email: "new@example.com", Role: "editor"})
```

***

## 28. 错误处理最佳实践

### 结构化错误响应

```go
func formatValidationError(err error) map[string]string {
    ve, ok := err.(validator.ValidationErrors)
    if !ok {
        return map[string]string{"error": err.Error()}
    }
    errs := make(map[string]string, len(ve))
    for _, fe := range ve {
        switch fe.Tag() {
        case "required":
            errs[fe.Field()] = fe.Field() + " 为必填字段"
        case "email":
            errs[fe.Field()] = fe.Field() + " 格式不正确"
        case "min":
            errs[fe.Field()] = fe.Field() + " 长度不能小于 " + fe.Param()
        case "max":
            errs[fe.Field()] = fe.Field() + " 长度不能大于 " + fe.Param()
        case "oneof":
            errs[fe.Field()] = fe.Field() + " 必须是 " + fe.Param() + " 之一"
        case "gtfield":
            errs[fe.Field()] = fe.Field() + " 必须大于 " + fe.Param()
        default:
            errs[fe.Field()] = fe.Field() + " 校验失败 (" + fe.Tag() + ")"
        }
    }
    return errs
}
```

### 配合 i18n 返回多语言错误

```go
func handleRequest(req interface{}, locale string) (interface{}, error) {
    v := validator.New(validator.WithRequiredStructEnabled())
    err := v.Struct(req)
    if err == nil {
        return nil, nil
    }

    messages := validator.TranslateValidationErrors(err, locale)

    type FieldError struct {
        Field     string `json:"field"`
        Namespace string `json:"namespace"`
        Tag       string `json:"tag"`
        Message   string `json:"message"`
    }

    var errors []FieldError
    for _, m := range messages {
        errors = append(errors, FieldError{
            Field:     m.Field,
            Namespace: m.Namespace,
            Tag:       m.Tag,
            Message:   m.Message,
        })
    }

    return errors, status.Error(codes.InvalidArgument, "validation failed")
}
```

***

## 29. WithRequiredStructEnabled 详解

默认情况下 `required` 对非指针结构体零值不生效，启用后空结构体也会被判定为零值：

```go
type Inner struct {
    Name string `validate:"required"`
}

type Outer struct {
    Inner Inner `validate:"required"`
}

v1 := validator.New()
v1.Struct(Outer{Inner: Inner{Name: ""}}) // Inner 零值但非指针 → required 不触发

v2 := validator.New(validator.WithRequiredStructEnabled())
v2.Struct(Outer{Inner: Inner{Name: ""}}) // Inner 零值 → required 触发

v3 := validator.New()
v3.Struct(Outer{Inner: Inner{Name: "ok"}}) // Inner 非零值 → 始终通过
```

> **推荐**：在 gRPC 拦截器中始终启用 `WithRequiredStructEnabled()`，确保 protobuf 嵌套消息的 `required` 标签对空结构体也能生效

***

## 30. omitempty / omitnil / omitzero 对比

```go
type OptionalFields struct {
    Nickname  string  `json:"nickname" validate:"omitempty,min=2,max=50"`   // 零值（""）时跳过
    Count     int     `json:"count" validate:"omitzero,gte=1"`              // 零值（0）时跳过
    Config    *Config `json:"config" validate:"omitnil,required"`           // nil 时跳过，非 nil 时 required 生效
    AvatarURL string  `json:"avatar_url" validate:"omitempty,http_url"`     // 空字符串跳过，非空时校验 URL
}
```

| 标签          | 跳过条件                          | 典型用途     |
| ----------- | ----------------------------- | -------- |
| `omitempty` | 字段为零值（`""`、`0`、`nil`、`false`） | 可选字符串/切片 |
| `omitzero`  | 与 `omitempty` 行为一致            | 可选数值字段   |
| `omitnil`   | 字段为 `nil`（仅指针/接口/slice/map）   | 可选指针结构体  |

