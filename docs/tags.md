# 校验标签（Tags）完整参考

Argus 通过 `validate` 结构体标签声明字段校验规则，语法与 `go-playground/validator` 高度兼容

## 语法格式

```go
type User struct {
    Name  string `validate:"required,min=2,max=50"`
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=0,lte=150"`
}
```

- 多条规则用逗号 `,` 分隔，按声明顺序依次校验
- 带参数的规则使用 `rule=param` 格式
- 逗号可用 `\` 转义以避免被当作规则分隔符，但反斜杠本身会保留在参数值中（例如 `contains=a\,b` 的参数为 `a\,b`，而非 `a,b`）。对于需要内部逗号分隔参数的规则（如 `range`），请使用 `|` 作为替代分隔符：`range=Start|End`
- 单个字段的规则校验失败时，立即停止该字段的后续规则校验并返回错误；其他字段的校验不受影响

***

## 必填与排除

### `required`

**场景**：确保字段不为零值（空字符串、0、nil、空切片等）

```go
type Req struct {
    Name string `validate:"required"`
}
```

**校验逻辑**：字段值经过指针解引用后，判断是否为零值。对于布尔类型 `false` 不算零值；对于结构体，当启用 `WithRequiredStructEnabled()` 时，空结构体视为零值

***

### `required_if`

**场景**：当另一个字段等于特定值时，当前字段必填

```go
type Req struct {
    Type   string `validate:"required,oneof=email sms"`
    Email  string `validate:"required_if=Type email"`
    Phone  string `validate:"required_if=Type sms"`
}
```

**参数格式**：`Field1 Value1 [Field2 Value2]...`，字段名与值成对出现，多对之间为 AND 关系（所有条件都满足时才触发）

**校验逻辑**：通过 `FieldByPath` 查找同结构体中的字段，将字段值转为字符串后与参数值比较。字段名支持 Go 字段名、json tag 名、lowerCamel 和 snake\_case 形式

***

### `required_unless`

**场景**：除非另一个字段等于特定值，否则当前字段必填（与 `required_if` 逻辑相反）

```go
type Req struct {
    Role   string `validate:"required,oneof=admin user"`
    Secret string `validate:"required_unless=Role user"`
}
```

**参数格式**：与 `required_if` 相同

**校验逻辑**：当所有条件对都满足时，跳过必填检查；否则要求当前字段非空

***

### `required_with`

**场景**：当任一指定字段非空时，当前字段必填

```go
type Req struct {
    Phone string `validate:"required_with=Email"`
    Email string `validate:"required_with=Phone"`
}
```

**参数格式**：`Field1 Field2 ...`，空格分隔多个字段名

**校验逻辑**：遍历所有指定字段，只要有一个字段非空（非零值），则当前字段必填

***

### `required_with_all`

**场景**：当所有指定字段都非空时，当前字段必填

```go
type Req struct {
    Street string `validate:"required_with_all=City Country"`
    City   string
    Country string
}
```

**参数格式**：与 `required_with` 相同

**校验逻辑**：所有指定字段都非空时，当前字段才必填

***

### `required_without`

**场景**：当任一指定字段为空时，当前字段必填

```go
type Req struct {
    Email string `validate:"required_without=Phone"`
    Phone string `validate:"required_without=Email"`
}
```

**参数格式**：与 `required_with` 相同

**校验逻辑**：只要有一个指定字段为空（零值），则当前字段必填

***

### `required_without_all`

**场景**：当所有指定字段都为空时，当前字段必填

```go
type Req struct {
    BackupEmail string `validate:"required_without_all=Email Phone"`
    Email       string
    Phone       string
}
```

**参数格式**：与 `required_with` 相同

**校验逻辑**：所有指定字段都为空时，当前字段才必填

***

### `isdefault`

**场景**：确保字段为零值（与 `required` 相反），常用于条件性排除

```go
type Req struct {
    AutoGen string `validate:"isdefault"`
}
```

**校验逻辑**：字段值经过指针解引用后，判断是否为零值

***

### `omitempty`

**场景**：字段为零值时跳过后续所有校验，常用于可选字段

```go
type Req struct {
    Nickname string `validate:"omitempty,min=2,max=50"`
}
```

**校验逻辑**：字段为零值时直接返回（跳过后续所有规则）；非零值时继续校验后续规则

> **注意**：`omitempty` 必须放在规则列表前面，否则前面的规则仍会执行

***

### `omitzero`

**场景**：与 `omitempty` 行为一致，字段为零值时跳过后续校验

```go
type Req struct {
    Count int `validate:"omitzero,gte=1"`
}
```

***

### `omitnil`

**场景**：字段为 nil 时跳过后续校验，专用于指针/接口/slice/map/channel 类型

```go
type Req struct {
    Config *Config `validate:"omitnil,required"`
}
```

**校验逻辑**：仅当字段为 nil（Kind 为 Ptr/Interface/Slice/Map/Chan/Func 且 `IsNil()` 为 true）时跳过，零值但不为 nil 的字段不会被跳过

***

### `excluded_if`

**场景**：当条件满足时，当前字段必须为空（与 `required_if` 互斥）

```go
type Req struct {
    Mode  string `validate:"required,oneof=basic advanced"`
    Debug string `validate:"excluded_if=Mode basic"`
}
```

**参数格式**：与 `required_if` 相同（`Field1 Value1 [Field2 Value2]...`）

**校验逻辑**：当所有条件对都满足时，当前字段必须为空值

***

### `excluded_unless`

**场景**：当条件不满足时，当前字段必须为空

```go
type Req struct {
    Role    string `validate:"required,oneof=admin user"`
    AdminOp string `validate:"excluded_unless=Role admin"`
}
```

**参数格式**：与 `required_if` 相同

**校验逻辑**：当所有条件对都满足时，允许当前字段非空；否则当前字段必须为空

***

### `excluded_with`

**场景**：当任一指定字段非空时，当前字段必须为空

```go
type Req struct {
    Email string `validate:"excluded_with=Phone"`
    Phone string
}
```

**参数格式**：`Field1 Field2 ...`

**校验逻辑**：任一指定字段非空时，当前字段必须为空值

***

### `excluded_with_all`

**场景**：当所有指定字段都非空时，当前字段必须为空

```go
type Req struct {
    LogDetail string `validate:"excluded_with_all=Verbose Debug"`
    Verbose   bool
    Debug     bool
}
```

**校验逻辑**：所有指定字段都非空时，当前字段必须为空值

***

### `excluded_without`

**场景**：当任一指定字段为空时，当前字段必须为空

```go
type Req struct {
    AltEmail string `validate:"excluded_without=Email"`
    Email    string `validate:"required"`
}
```

**校验逻辑**：任一指定字段为空时，当前字段必须为空值

***

### `excluded_without_all`

**场景**：当所有指定字段都为空时，当前字段必须为空

```go
type Req struct {
    Recovery string `validate:"excluded_without_all=Email Phone"`
    Email    string
    Phone    string
}
```

**校验逻辑**：所有指定字段都为空时，当前字段必须为空值

***

## 数值与长度比较

> **核心规则**：对于字符串类型，`min`/`max`/`len`/`gt`/`gte`/`lt`/`lte` 比较的是 **字符数**（`rune` 长度，即 Unicode 字符数）；对于切片/数组/Map，比较的是元素个数；对于数值类型，比较的是数值本身

### `eq`

**场景**：字符串精确匹配

```go
type Req struct {
    Status string `validate:"eq=active"`
}
```

**校验逻辑**：将字段值转为字符串后与参数比较，区分大小写

***

### `eq_ignore_case`

**场景**：字符串匹配，不区分大小写

```go
type Req struct {
    Status string `validate:"eq_ignore_case=ACTIVE"`
}
```

**校验逻辑**：使用 `strings.EqualFold` 比较

***

### `ne`

**场景**：字符串不等于指定值

```go
type Req struct {
    Status string `validate:"ne=deleted"`
}
```

***

### `ne_ignore_case`

**场景**：字符串不等于指定值，不区分大小写

```go
type Req struct {
    Status string `validate:"ne_ignore_case=DRAFT"`
}
```

***

### `gt`

**场景**：大于指定值

```go
type Req struct {
    Age     int    `validate:"gt=0"`
    Name    string `validate:"gt=3"`     // 字符数 > 3
    Scores  []int  `validate:"gt=0"`     // 元素个数 > 0
}
```

**校验逻辑**：字符串比较 rune 长度，数值比较数值大小，切片/Map 比较元素个数

***

### `gte`

**场景**：大于等于指定值

```go
type Req struct {
    Age int `validate:"gte=18"`
}
```

***

### `lt`

**场景**：小于指定值

```go
type Req struct {
    Age     int    `validate:"lt=150"`
    Name    string `validate:"lt=100"`   // 字符数 < 100
}
```

***

### `lte`

**场景**：小于等于指定值

```go
type Req struct {
    Discount float64 `validate:"lte=100"`
}
```

***

### `min`

**场景**：最小值/最小长度约束

```go
type Req struct {
    Name    string `validate:"min=2"`    // 至少 2 个字符
    Age     int    `validate:"min=0"`
    Tags    []string `validate:"min=1"`  // 至少 1 个元素
}
```

**校验逻辑**：等价于 `gte`

***

### `max`

**场景**：最大值/最大长度约束

```go
type Req struct {
    Name    string `validate:"max=50"`   // 最多 50 个字符
    Age     int    `validate:"max=150"`
    Tags    []string `validate:"max=10"` // 最多 10 个元素
}
```

**校验逻辑**：等价于 `lte`

***

### `len`

**场景**：精确长度/值约束

```go
type Req struct {
    Code    string `validate:"len=6"`    // 必须恰好 6 个字符
    Pin     int    `validate:"len=6"`    // 数值必须等于 6
    Items   []string `validate:"len=3"`  // 必须恰好 3 个元素
}
```

**校验逻辑**：字符串比较 rune 长度，数值比较数值大小，切片/Map 比较元素个数

***

## 跨字段比较

> **字段查找规则**：所有跨字段规则通过 `FieldByPath` 查找同结构体（`*field` 系列）或顶层结构体（`*csfield` 系列）中的字段。字段名支持以下形式：
>
> - Go 导出字段名：`FirstName`
> - JSON tag 名：`first_name`
> - lowerCamel：`firstName`
> - snake\_case：`first_name`
>
> 支持嵌套路径，如 `Address.City`

### `eqfield`

**场景**：当前字段必须等于同结构体中另一字段的值

```go
type PasswordForm struct {
    Password        string `validate:"required,min=8"`
    ConfirmPassword string `validate:"required,eqfield=Password"`
}
```

**校验逻辑**：优先尝试时间比较，其次数值比较，最后字符串比较

***

### `nefield`

**场景**：当前字段不能等于同结构体中另一字段的值

```go
type Req struct {
    NewPassword string `validate:"required,min=8"`
    OldPassword string `validate:"required,nefield=NewPassword"`
}
```

***

### `gtfield` / `gtefield` / `ltfield` / `ltefield`

**场景**：与同结构体字段进行大小比较

```go
type DateRange struct {
    StartDate time.Time `validate:"required"`
    EndDate   time.Time `validate:"required,gtfield=StartDate"`
}

type PriceRange struct {
    MinPrice float64 `validate:"required,gte=0"`
    MaxPrice float64 `validate:"required,gtfield=MinPrice"`
}
```

**校验逻辑**：优先尝试时间比较（支持 `time.Time`、时间字符串、protobuf Timestamp），其次数值比较，最后字符串字典序比较

***

### `eqcsfield` / `necsfield` / `gtcsfield` / `gtecsfield` / `ltcsfield` / `ltecsfield`

**场景**：与**顶层**结构体中的字段比较（跨嵌套层级）

```go
type Inner struct {
    Value int `validate:"gtcsfield=Top.Min"`
}
type Top struct {
    Min   int
    Inner Inner
}
```

**校验逻辑**：与 `*field` 系列相同，区别在于从顶层结构体开始查找字段

***

### `fieldcontains`

**场景**：当前字段必须包含指定字段的值

```go
type Req struct {
    Keyword string
    Content string `validate:"fieldcontains=Keyword"`
}
```

**校验逻辑**：将当前字段和目标字段都转为字符串，检查 `strings.Contains(当前值, 目标值)`

***

### `fieldexcludes`

**场景**：当前字段不能包含指定字段的值

```go
type Req struct {
    Forbidden string
    Content   string `validate:"fieldexcludes=Forbidden"`
}
```

**校验逻辑**：与 `fieldcontains` 相反，检查当前字符串不包含目标字段的值

***

## 时间比较

> **时间识别**：Argus 支持以下时间类型：
>
> - `time.Time` 结构体
> - 时间字符串（自动尝试 RFC3339Nano、RFC3339、`2006-01-02 15:04:05`、`2006-01-02` 格式）
> - Protobuf Timestamp（通过 `AsTime()` 或 `GetSeconds()/GetNanos()` 方法）
> - 包含 `Seconds` 和 `Nanos` 字段的结构体

### `datetime`

**场景**：验证字符串是否符合指定时间格式

```go
type Req struct {
    CreatedAt string `validate:"datetime=2006-01-02"`
    UpdatedAt string `validate:"datetime"`            // 默认 RFC3339
}
```

**参数**：Go 时间布局字符串，省略时默认 `time.RFC3339`

**校验逻辑**：使用 `time.Parse(layout, value)` 解析，解析成功即通过

***

### `after`

**场景**：时间必须晚于指定表达式

```go
type Req struct {
    StartTime time.Time `validate:"after=now"`
    ExpiredAt time.Time `validate:"after=now+30d"`
    BookedAt  string    `validate:"after=now-7d"`
}
```

**参数**：时间表达式，支持：

- `now` — 当前时间
- `now+5m` — 当前时间 + 5 分钟
- `now-30d` — 当前时间 - 30 天
- `now+2h30m` — 当前时间 + 2 小时 30 分钟
- 支持 Go `time.Duration` 格式（`h`/`m`/`s`/`ms`/`µs`/`ns`）以及 `d`（天）

**校验逻辑**：将字段值识别为时间，将表达式解析为时间，比较 `字段时间 > 表达式时间`

***

### `before`

**场景**：时间必须早于指定表达式

```go
type Req struct {
    BirthDate time.Time `validate:"before=now"`
    EndTime   string    `validate:"before=now+1h"`
}
```

**校验逻辑**：与 `after` 相反，比较 `字段时间 < 表达式时间`

***

### `afterfield` / `beforefield`

**场景**：与同结构体中另一字段时间比较

```go
type Event struct {
    StartAt time.Time `validate:"required"`
    EndAt   time.Time `validate:"required,afterfield=StartAt"`
}
```

**校验逻辑**：`afterfield` 等价于 `gtfield`，`beforefield` 等价于 `ltfield`，均优先使用时间比较

***

### `range`

**场景**：验证起始字段值小于结束字段值（常用于时间范围校验）

```go
type Req struct {
    StartAt time.Time `validate:"required"`
    EndAt   time.Time `validate:"required"`
}
// 在父结构体上使用，用 | 分隔两个字段名
type Schedule struct {
    Req Req `validate:"range=StartAt|EndAt"`
}
```

**参数格式**：`StartField|EndField`（使用 `|` 分隔两个字段名）

> **注意**：由于 validate 标签中逗号是规则分隔符，且 `\` 转义后反斜杠会保留在参数值中，`range` 规则内部使用 `|` 作为字段分隔符。当参数中包含 `|` 时，`ruleRange` 会自动切换为 `|` 分隔模式

**校验逻辑**：通过 `FieldByPath` 查找两个字段的值，使用 `CompareValue` 比较 `StartField < EndField`

***

### `timezone`

**场景**：验证字符串是否为有效的 IANA 时区名

```go
type Req struct {
    TZ string `validate:"timezone"`
}
```

**校验逻辑**：使用 `time.LoadLocation(value)` 验证，如 `Asia/Shanghai`、`UTC`、`America/New_York`

***

## 字符串内容

### `alpha`

**场景**：仅允许 ASCII 字母（a-z, A-Z）

```go
type Req struct {
    Code string `validate:"alpha"`
}
```

**通过**：`"Hello"`、`"abc"`  **不通过**：`"Hello123"`、`"你好"`

***

### `alphaspace`

**场景**：仅允许 ASCII 字母和空格

```go
type Req struct {
    Name string `validate:"alphaspace"`
}
```

**通过**：`"Hello World"`  **不通过**：`"Hello123"`

***

### `alphanum`

**场景**：仅允许 ASCII 字母和数字

```go
type Req struct {
    Username string `validate:"alphanum"`
}
```

**通过**：`"User123"`  **不通过**：`"User 123"`、`"用户123"`

***

### `alphanumspace`

**场景**：仅允许 ASCII 字母、数字和空格

```go
type Req struct {
    Title string `validate:"alphanumspace"`
}
```

***

### `alphaunicode`

**场景**：仅允许 Unicode 字母（支持中文、日文等）

```go
type Req struct {
    Name string `validate:"alphaunicode"`
}
```

**通过**：`"你好"`、`"Hello"`、`"こんにちは"`  **不通过**：`"Hello123"`

***

### `alphanumunicode`

**场景**：仅允许 Unicode 字母和数字

```go
type Req struct {
    Name string `validate:"alphanumunicode"`
}
```

**通过**：`"你好123"`、`"Hello42"`  **不通过**：`"Hello!"`

***

### `ascii`

**场景**：仅允许 ASCII 字符（0x00 - 0x7F）

```go
type Req struct {
    Data string `validate:"ascii"`
}
```

***

### `printascii`

**场景**：仅允许可打印 ASCII 字符（0x20 - 0x7E）

```go
type Req struct {
    Text string `validate:"printascii"`
}
```

***

### `multibyte`

**场景**：必须包含至少一个多字节字符（如中文、日文、Emoji）

```go
type Req struct {
    Content string `validate:"multibyte"`
}
```

**校验逻辑**：`len(s) != utf8.RuneCountInString(s)`，即字节长度不等于字符数

***

### `lowercase`

**场景**：字符串必须全部小写

```go
type Req struct {
    Code string `validate:"lowercase"`
}
```

**校验逻辑**：`s == strings.ToLower(s)`

***

### `uppercase`

**场景**：字符串必须全部大写

```go
type Req struct {
    Code string `validate:"uppercase"`
}
```

***

### `boolean`

**场景**：验证布尔值或可解析为布尔的字符串

```go
type Req struct {
    Active string `validate:"boolean"`  // "true"/"false"/"1"/"0"
    Flag   bool   `validate:"boolean"`
}
```

**校验逻辑**：字符串类型使用 `strconv.ParseBool` 解析（接受 `1/0/t/f/T/F/true/false/TRUE/FALSE`）；布尔类型直接通过

***

### `number` / `numeric`

**场景**：验证数值或可解析为数值的字符串

```go
type Req struct {
    Count string `validate:"numeric"`  // "42"、"3.14"
    Age   int    `validate:"numeric"`  // 数值类型直接通过
}
```

**校验逻辑**：数值类型（int/uint/float 系列）直接通过；字符串类型使用 `strconv.ParseFloat` 解析

***

## 字符串匹配

### `contains`

**场景**：字符串必须包含指定子串

```go
type Req struct {
    Email string `validate:"contains=@"`
}
```

***

### `containsany`

**场景**：字符串必须包含指定字符集中的任一字符

```go
type Req struct {
    Password string `validate:"containsany=!@#$%"`
}
```

**校验逻辑**：使用 `strings.ContainsAny(s, chars)`

***

### `containsrune`

**场景**：字符串必须包含指定 rune

```go
type Req struct {
    Text string `validate:"containsrune=中"`
}
```

**校验逻辑**：参数的第一个 rune，使用 `strings.ContainsRune` 检查

***

### `excludes`

**场景**：字符串不能包含指定子串

```go
type Req struct {
    Username string `validate:"excludes= "`
}
```

***

### `excludesall`

**场景**：字符串不能包含指定字符集中的任一字符

```go
type Req struct {
    Filename string `validate:"excludesall=/\\"`
}
```

***

### `excludesrune`

**场景**：字符串不能包含指定 rune

```go
type Req struct {
    Text string `validate:"excludesrune=<>"`
}
```

***

### `startswith`

**场景**：字符串必须以指定前缀开头

```go
type Req struct {
    URL string `validate:"startswith=https://"`
}
```

***

### `endswith`

**场景**：字符串必须以指定后缀结尾

```go
type Req struct {
    File string `validate:"endswith=.go"`
}
```

***

### `startsnotwith`

**场景**：字符串不能以指定前缀开头

```go
type Req struct {
    Code string `validate:"startsnotwith=0"`
}
```

***

### `endsnotwith`

**场景**：字符串不能以指定后缀结尾

```go
type Req struct {
    Path string `validate:"endsnotwith=/"`
}
```

***

## 枚举与唯一

### `oneof`

**场景**：值必须是指定列表中的某一项

```go
type Req struct {
    Role   string `validate:"oneof=admin user guest"`
    Status int    `validate:"oneof=1 2 3"`
}
```

**参数格式**：空格分隔的值列表

**校验逻辑**：将字段值转为字符串后，逐一与列表项比较（区分大小写）

***

### `oneofci`

**场景**：同 `oneof`，但不区分大小写

```go
type Req struct {
    Color string `validate:"oneofci=Red Green Blue"`
}
```

***

### `noneof`

**场景**：值不能是指定列表中的任一项

```go
type Req struct {
    Username string `validate:"noneof=admin root system"`
}
```

***

### `noneofci`

**场景**：同 `noneof`，但不区分大小写

***

### `unique`

**场景**：确保元素唯一性

```go
type Req struct {
    Tags    string   `validate:"unique"`     // 字符串中每个 rune 唯一
    Emails  []string `validate:"unique"`     // 切片元素唯一
    Scores  map[string]int `validate:"unique"` // Map 值唯一
}
```

**校验逻辑**：

- **字符串**：每个 rune 必须唯一
- **切片/数组**：每个元素转为字符串后必须唯一
- **Map**：每个值转为字符串后必须唯一

***

## 网络地址

### `ip` / `ip_addr`

**场景**：验证 IP 地址（IPv4 或 IPv6）

```go
type Req struct {
    Host string `validate:"ip"`
}
```

**校验逻辑**：使用 `net.ParseIP` 解析，支持 IPv4 和 IPv6

***

### `ipv4`

**场景**：验证 IPv4 地址

```go
type Req struct {
    Host string `validate:"ipv4"`
}
```

**校验逻辑**：`net.ParseIP` 解析后检查 `ip.To4() != nil`

***

### `ipv6`

**场景**：验证 IPv6 地址

```go
type Req struct {
    Host string `validate:"ipv6"`
}
```

**校验逻辑**：`net.ParseIP` 解析后检查 `ip.To4() == nil`

***

### `cidr`

**场景**：验证 CIDR 表示法

```go
type Req struct {
    Network string `validate:"cidr"`
}
```

**通过**：`"192.168.1.0/24"`、`"2001:db8::/32"`

***

### `cidrv4`

**场景**：验证 IPv4 CIDR

```go
type Req struct {
    Network string `validate:"cidrv4"`
}
```

**通过**：`"192.168.1.0/24"`  **不通过**：`"2001:db8::/32"`

***

### `cidrv6`

**场景**：验证 IPv6 CIDR

***

### `mac`

**场景**：验证 MAC 地址

```go
type Req struct {
    MAC string `validate:"mac"`
}
```

**通过**：`"01:23:45:67:89:ab"`、`"01-23-45-67-89-ab"`

***

### `hostname` / `hostname_rfc1123`

**场景**：验证 RFC 1123 主机名

```go
type Req struct {
    Host string `validate:"hostname"`
}
```

**校验逻辑**：每个标签必须以字母或数字开头和结尾，中间可含连字符，长度 1-63，允许末尾有点号

***

### `fqdn`

**场景**：验证完全限定域名（FQDN）

```go
type Req struct {
    Domain string `validate:"fqdn"`
}
```

**校验逻辑**：必须以点号结尾，且点号前的部分符合主机名规则

**通过**：`"example.com."`  **不通过**：`"example.com"`

***

### `hostname_port`

**场景**：验证 `hostname:port` 格式

```go
type Req struct {
    Addr string `validate:"hostname_port"`
}
```

**通过**：`"example.com:8080"`、`"localhost:3000"`

**校验逻辑**：使用 `net.SplitHostPort` 分割后分别验证主机名和端口

***

### `port`

**场景**：验证端口号

```go
type Req struct {
    Port int    `validate:"port"`
    PortStr string `validate:"port"`
}
```

**校验逻辑**：转为整数后检查 0-65535 范围

***

### `url`

**场景**：验证 URL（必须包含 scheme）

```go
type Req struct {
    Webhook string `validate:"url"`
}
```

**通过**：`"https://example.com"`、`"ftp://files.example.com"`

**校验逻辑**：使用 `url.ParseRequestURI` 解析，且 scheme 不为空

***

### `uri`

**场景**：验证 URI

```go
type Req struct {
    Resource string `validate:"uri"`
}
```

**校验逻辑**：使用 `url.ParseRequestURI` 解析，不要求 scheme

***

### `http_url`

**场景**：验证 HTTP/HTTPS URL

```go
type Req struct {
    API string `validate:"http_url"`
}
```

**通过**：`"http://localhost:8080/api"`、`"https://example.com"`

**校验逻辑**：scheme 必须是 `http` 或 `https`，且 host 不为空

***

### `https_url`

**场景**：验证 HTTPS URL

```go
type Req struct {
    SecureAPI string `validate:"https_url"`
}
```

**校验逻辑**：scheme 必须是 `https`，且 host 不为空

***

### `url_encoded`

**场景**：验证 URL 编码字符串

```go
type Req struct {
    Query string `validate:"url_encoded"`
}
```

**校验逻辑**：字符串必须包含 `%`，且 `url.QueryUnescape` 解析成功

***

### `email`

**场景**：验证 Email 地址

```go
type Req struct {
    Email string `validate:"email"`
}
```

**校验逻辑**：使用 `validator.IsEmail` 验证

***

### `e164`

**场景**：验证 E.164 国际电话号码格式

```go
type Req struct {
    Phone string `validate:"e164"`
}
```

**通过**：`"+8613800138000"`、`"+14155552671"`

**校验逻辑**：正则 `^\+[1-9]\d{1,14}$`，以 `+` 开头，后跟 1-15 位数字

***

## 编码与格式

### `json`

**场景**：验证 JSON 字符串或 `[]byte`

```go
type Req struct {
    Payload  string `validate:"json"`
    RawBody  []byte `validate:"json"`
}
```

**校验逻辑**：使用 `json.Valid` 验证

***

### `base32`

**场景**：验证 Base32 编码字符串

```go
type Req struct {
    Data string `validate:"base32"`
}
```

**校验逻辑**：使用 `base32.StdEncoding.DecodeString` 解码

***

### `base64`

**场景**：验证 Base64 编码字符串

```go
type Req struct {
    Data string `validate:"base64"`
}
```

**校验逻辑**：使用 `validator.IsBase64` 验证

***

### `base64url`

**场景**：验证 Base64 URL 编码字符串

```go
type Req struct {
    Token string `validate:"base64url"`
}
```

**校验逻辑**：使用 `base64.URLEncoding.DecodeString` 解码

***

### `base64rawurl`

**场景**：验证 Base64 Raw URL 编码字符串（无 padding）

```go
type Req struct {
    Token string `validate:"base64rawurl"`
}
```

**校验逻辑**：使用 `base64.RawURLEncoding.DecodeString` 解码

***

### `uuid` / `uuid_rfc4122`

**场景**：验证 UUID 格式

```go
type Req struct {
    ID string `validate:"uuid"`
}
```

**校验逻辑**：使用 `validator.IsUUID` 验证，接受任何版本的 UUID

***

### `uuid3` / `uuid3_rfc4122`

**场景**：验证 UUID v3

```go
type Req struct {
    ID string `validate:"uuid3"`
}
```

**校验逻辑**：验证 UUID 格式且版本号为 3

***

### `uuid4` / `uuid4_rfc4122`

**场景**：验证 UUID v4

```go
type Req struct {
    ID string `validate:"uuid4"`
}
```

***

### `uuid5` / `uuid5_rfc4122`

**场景**：验证 UUID v5

```go
type Req struct {
    ID string `validate:"uuid5"`
}
```

***

### `hexadecimal`

**场景**：验证十六进制字符串

```go
type Req struct {
    Hex string `validate:"hexadecimal"`
}
```

**通过**：`"deadbeef"`、`"CAFEBABE"`  **不通过**：`"xyz"`

***

### `hexcolor`

**场景**：验证十六进制颜色值

```go
type Req struct {
    Color string `validate:"hexcolor"`
}
```

**通过**：`"#fff"`、`"#000000"`、`"#FFFFFFFF"`  **不通过**：`"red"`

**校验逻辑**：正则匹配，支持 3/4/6/8 位十六进制，`#` 前缀可选

***

### `html`

**场景**：字符串必须包含 HTML 标签

```go
type Req struct {
    Content string `validate:"html"`
}
```

**校验逻辑**：同时包含 `<` 和 `>` 字符

***

### `html_encoded`

**场景**：验证 HTML 编码字符串

```go
type Req struct {
    Content string `validate:"html_encoded"`
}
```

**校验逻辑**：`html.UnescapeString(s) != s`，即解码后与原串不同

***

## 颜色

### `rgb`

**场景**：验证 RGB 颜色值

```go
type Req struct {
    Color string `validate:"rgb"`
}
```

**通过**：`"rgb(255, 128, 0)"`  **不通过**：`"rgb(300, 0, 0)"`

**校验逻辑**：正则匹配，每个分量 0-255

***

### `rgba`

**场景**：验证 RGBA 颜色值

```go
type Req struct {
    Color string `validate:"rgba"`
}
```

**通过**：`"rgba(255, 128, 0, 0.5)"`  **不通过**：`"rgba(255, 128, 0, 2)"`

**校验逻辑**：正则匹配，RGB 分量 0-255，Alpha 分量 0-1

***

### `hsl`

**场景**：验证 HSL 颜色值

```go
type Req struct {
    Color string `validate:"hsl"`
}
```

**通过**：`"hsl(180, 50%, 50%)"`

**校验逻辑**：色相 0-360，饱和度 0-100%，亮度 0-100%

***

### `hsla`

**场景**：验证 HSLA 颜色值

```go
type Req struct {
    Color string `validate:"hsla"`
}
```

**通过**：`"hsla(180, 50%, 50%, 0.5)"`

***

## 地理与文件

### `latitude`

**场景**：验证纬度值

```go
type Req struct {
    Lat float64 `validate:"latitude"`
}
```

**校验逻辑**：数值范围 -90 \~ 90，支持 int/uint/float/string 类型

***

### `longitude`

**场景**：验证经度值

```go
type Req struct {
    Lng float64 `validate:"longitude"`
}
```

**校验逻辑**：数值范围 -180 \~ 180，支持 int/uint/float/string 类型

***

### `file`

**场景**：验证路径指向一个存在的文件

```go
type Req struct {
    Path string `validate:"file"`
}
```

**校验逻辑**：`os.Stat(path)` 成功且 `!info.IsDir()`

***

### `filepath`

**场景**：验证有效的文件路径格式（不检查文件是否存在）

```go
type Req struct {
    Path string `validate:"filepath"`
}
```

**校验逻辑**：非空且 `filepath.Clean(path) != "."`

***

### `dir`

**场景**：验证路径指向一个存在的目录

```go
type Req struct {
    Path string `validate:"dir"`
}
```

**校验逻辑**：`os.Stat(path)` 成功且 `info.IsDir()`

***

### `dirpath`

**场景**：验证有效的目录路径格式（不检查目录是否存在）

```go
type Req struct {
    Path string `validate:"dirpath"`
}
```

**校验逻辑**：非空、`filepath.Clean(path) != "."`，且路径最后一段不含 `.`

***

## 专业格式校验

### `semver`

**场景**：验证语义化版本号（Semantic Versioning）

```go
type Req struct {
    Version string `validate:"semver"`
}
```

**通过**：`"1.2.3"`、`"v1.0.0-alpha.1"`、`"2.0.0+build.123"`  **不通过**：`"1.2"`、`"1"`、`"abc"`

**校验逻辑**：正则匹配 `v?主版本.次版本.修订版[-预发布][+构建元数据]`，预发布和构建元数据可选

***

### `isbn10`

**场景**：验证 ISBN-10 国际标准书号

```go
type Req struct {
    BookCode string `validate:"isbn10"`
}
```

**通过**：`"0471958697"`、`"080442957X"`  **不通过**：`"0471958698"`（校验位错误）、`"123456789"`

**校验逻辑**：10 位数字（最后一位可为 X 表示 10），加权求和 `∑(dᵢ × (11-i))` 必须能被 11 整除。自动去除横线和空格

***

### `isbn13`

**场景**：验证 ISBN-13 国际标准书号

```go
type Req struct {
    BookCode string `validate:"isbn13"`
}
```

**通过**：`"9780471117094"`  **不通过**：`"9780471117095"`（校验位错误）

**校验逻辑**：13 位数字，奇数位权重 1、偶数位权重 3，加权求和必须能被 10 整除。自动去除横线和空格

***

### `issn`

**场景**：验证 ISSN 国际标准连续出版物号

```go
type Req struct {
    SerialCode string `validate:"issn"`
}
```

**通过**：`"0317847X"`、`"0317-847X"`  **不通过**：`"03178470"`（校验位错误）

**校验逻辑**：8 位字符（最后一位可为 X 表示 10），加权求和 `∑(dᵢ × (9-i))` 必须能被 11 整除。自动去除横线

***

### `bic`

**场景**：验证 BIC/SWIFT 银行识别码

```go
type Req struct {
    SwiftCode string `validate:"bic"`
}
```

**通过**：`"CHASUS33"`、`"CHASUS33XXX"`  **不通过**：`"INVALID"`、`"AB"`

**校验逻辑**：正则匹配 8 或 11 位，前 4 位字母（银行代码）、2 位字母（国家代码）、2 位字母数字（地区代码）、可选 3 位字母数字（分行代码）

***

### `cron`

**场景**：验证 cron 表达式格式

```go
type Req struct {
    Schedule string `validate:"cron"`
}
```

**通过**：`"*/5 * * * *"`、`"0 6 * * 1-5"`  **不通过**：`"* * * *"`（4 字段）、`"* * * * * * *"`（7 字段）

**校验逻辑**：必须为 5 或 6 个空白分隔的字段，每个字段为非空字符串

***

### `datauri`

**场景**：验证 Data URI（RFC 2397）

```go
type Req struct {
    Avatar string `validate:"datauri"`
}
```

**通过**：`"data:text/plain;base64,SGVsbG8="`、`"data:text/html,Hello"`  **不通过**：`"http://example.com"`、`"data:,"`

**校验逻辑**：正则匹配 `data:[<mediatype>][;base64],<data>` 格式

***

### `bcp47`

**场景**：验证 BCP 47 语言标签

```go
type Req struct {
    Locale string `validate:"bcp47"`
}
```

**通过**：`"en"`、`"zh-CN"`、`"sr-Latn-RS"`  **不通过**：`"1"`、`"a"`

**校验逻辑**：正则匹配 `语言[-脚本][-地区][-扩展]` 格式，语言 2-3 字母，脚本 4 字母，地区 2 字母或 3 数字

***

### `eth_addr`

**场景**：验证以太坊地址

```go
type Req struct {
    Wallet string `validate:"eth_addr"`
}
```

**通过**：`"0x742d35Cc6634C0532925a3b844Bc9e7595f2bD38"`  **不通过**：`"0x1234"`、`"742d35Cc6634C0532925a3b844Bc9e7595f2bD38"`

**校验逻辑**：以 `0x` 开头，后跟 40 位十六进制字符，总长度 42

***

### `btc_addr`

**场景**：验证比特币地址

```go
type Req struct {
    Wallet string `validate:"btc_addr"`
}
```

**通过**：`"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"`（Legacy）、`"3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy"`（P2SH）  **不通过**：`"not-a-btc-address"`

**校验逻辑**：正则匹配 Legacy（1 开头 25-34 位）、P2SH（3 开头 25-34 位）或 Bech32（bc1q 开头 39-59 位）格式

***

## 特殊标识

### `mongodb`

**场景**：验证 MongoDB ObjectId

```go
type Req struct {
    ID string `validate:"mongodb"`
}
```

**校验逻辑**：正则 `^[0-9a-fA-F]{24}$`，24 位十六进制字符串

***

### `luhn_checksum` / `credit_card`

**场景**：验证 Luhn 校验和（信用卡号等）

```go
type Req struct {
    CardNumber string `validate:"credit_card"`
}
```

**校验逻辑**：标准 Luhn 算法，忽略空格和连字符，至少 1 位数字，校验和能被 10 整除

***

### `dns_rfc1035_label`

**场景**：验证 DNS RFC 1035 标签

```go
type Req struct {
    Label string `validate:"dns_rfc1035_label"`
}
```

**校验逻辑**：以小写字母开头，仅含小写字母/数字/连字符，长度 ≤ 63

***

### `dive`

**场景**：对切片/Map 的每个元素逐一校验

```go
type Req struct {
    Emails []string `validate:"dive,email"`
    Scores map[string]int `validate:"dive,gte=0,lte=100"`
}
```

**校验逻辑**：`dive` 之后的规则作用于每个元素。对于切片，命名空间为 `Field[0]`、`Field[1]`...；对于 Map，为 `Field[key1]`、`Field[key2]`...

> **注意**：`dive` 必须是最后一个修饰符，其后的所有规则都作用于元素级别

***

### `structonly`

**场景**：仅校验结构体本身字段上的标签，不递归校验嵌套结构体

```go
type Inner struct {
    Name string `validate:"required"`
}
type Outer struct {
    Inner Inner `validate:"structonly"`
}
```

**效果**：只校验 `Outer.Inner` 字段本身的规则，不会递归进入 `Inner` 校验 `Name` 字段

***

### `nostructlevel`

**场景**：跳过当前字段的结构体递归校验（与 `structonly` 行为一致）

```go
type Inner struct {
    Name string `validate:"required"`
}
type Outer struct {
    Inner Inner `validate:"nostructlevel"`
}
```

**效果**：跳过对 `Inner` 嵌套结构体的递归校验，`Inner.Name` 的 `required` 规则不会被检查

> **注意**：当前实现中 `nostructlevel` 与 `structonly` 行为一致，均阻止递归进入嵌套结构体。这与 `go-playground/validator` 中 `nostructlevel` 的语义（跳过结构体级别校验器但继续递归字段校验）有所不同

***

## 字段查找规则

跨字段规则（`required_if`、`eqfield`、`fieldcontains` 等）通过 `FieldByPath` 查找目标字段，支持以下名称形式：

| 形式          | 示例             | 说明                      |
| ----------- | -------------- | ----------------------- |
| Go 字段名      | `FirstName`    | 结构体导出字段名                |
| JSON tag    | `first_name`   | `json:"first_name"` 标签值 |
| lowerCamel  | `firstName`    | 首字母小写的字段名               |
| snake\_case | `first_name`   | 下划线分隔的小写字段名             |
| 嵌套路径        | `Address.City` | 点号分隔的多级字段               |

> **注意**：`*field` 系列在同结构体中查找，`*csfield` 系列在顶层结构体中查找

## 零值判断规则

`required`、`omitempty`、`excluded_*` 等规则依赖零值判断，零值定义如下：

| 类型                      | 零值                                      |
| ----------------------- | --------------------------------------- |
| 字符串                     | `""`                                    |
| 数值（int/uint/float）      | `0`                                     |
| 布尔                      | `false`                                 |
| 指针/接口/slice/map/channel | `nil`                                   |
| 结构体                     | 启用 `WithRequiredStructEnabled()` 时为空结构体 |
| `time.Time`             | `time.Time{}`（零值时间）                     |

> 指针类型会先解引用再判断，`*string(nil)` 为零值，`*string("")` 为非零值

