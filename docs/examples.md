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

