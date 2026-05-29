/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-19 08:57:16
 * @FilePath: \go-argus\reexport.go
 * @Description: 根包兼容入口，统一转发 validate 子包的所有导出能力
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"net"
	"reflect"
	"regexp"
	"time"

	"github.com/kamalyes/go-argus/schema"
	"github.com/kamalyes/go-argus/validate"
)

// ────────────────────────────────────────
// 比较校验（validate/compare）
// ────────────────────────────────────────

// CompareOperator 表示通用比较操作符
type CompareOperator = validate.CompareOperator

const (
	OpEqual                    = validate.OpEqual
	OpNotEqual                 = validate.OpNotEqual
	OpGreaterThan              = validate.OpGreaterThan
	OpGreaterThanOrEqual       = validate.OpGreaterThanOrEqual
	OpLessThan                 = validate.OpLessThan
	OpLessThanOrEqual          = validate.OpLessThanOrEqual
	OpContains                 = validate.OpContains
	OpNotContains              = validate.OpNotContains
	OpHasPrefix                = validate.OpHasPrefix
	OpHasSuffix                = validate.OpHasSuffix
	OpRegex                    = validate.OpRegex
	OpEmpty                    = validate.OpEmpty
	OpNotEmpty                 = validate.OpNotEmpty
	OpSymbolEqual              = validate.OpSymbolEqual
	OpSymbolNotEqual           = validate.OpSymbolNotEqual
	OpSymbolGreaterThan        = validate.OpSymbolGreaterThan
	OpSymbolGreaterThanOrEqual = validate.OpSymbolGreaterThanOrEqual
	OpSymbolLessThan           = validate.OpSymbolLessThan
	OpSymbolLessThanOrEqual    = validate.OpSymbolLessThanOrEqual
)

// CompareResult 表示一次比较校验结果
type CompareResult = validate.CompareResult

// Number 表示可比较数值类型集合
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// CompareNumbers 比较两个数值
func CompareNumbers[T Number](actual, expect T, op CompareOperator) CompareResult {
	return validate.CompareNumbers(actual, expect, op)
}

// CompareStrings 比较两个字符串
func CompareStrings(actual, expect string, op CompareOperator) CompareResult {
	return validate.CompareStrings(actual, expect, op)
}

// ValidateString 校验字符串关系
func ValidateString(actual, expect string, op CompareOperator) CompareResult {
	return validate.ValidateString(actual, expect, op)
}

// ValidateContains 校验字节内容是否包含子串
func ValidateContains(body []byte, substring string) CompareResult {
	return validate.ValidateContains(body, substring)
}

// ValidateNotContains 校验字节内容是否不包含子串
func ValidateNotContains(body []byte, substring string) CompareResult {
	return validate.ValidateNotContains(body, substring)
}

// ValidateStatusCode 比较 HTTP 状态码
func ValidateStatusCode(statusCode, expected int, op CompareOperator) CompareResult {
	return validate.ValidateStatusCode(statusCode, expected, op)
}

// ValidateStatusCodeRange 校验 HTTP 状态码是否在闭区间内
func ValidateStatusCodeRange(actual, min, max int) CompareResult {
	return validate.ValidateStatusCodeRange(actual, min, max)
}

// ValidateHeader 根据操作符比较 Header 值
func ValidateHeader(headers map[string]string, key, expected string, op CompareOperator) CompareResult {
	return validate.ValidateHeader(headers, key, expected, op)
}

// ValidateContentType 校验 Content-Type 是否包含期望类型
func ValidateContentType(headers map[string]string, expected string) CompareResult {
	return validate.ValidateContentType(headers, expected)
}

// ────────────────────────────────────────
// 格式校验（validate/format）
// ────────────────────────────────────────

// ValidateRegex 校验字节内容是否匹配正则表达式
func ValidateRegex(body []byte, pattern string) CompareResult {
	return validate.ValidateRegex(body, pattern)
}

// ValidateEmail 校验 Email 格式
func ValidateEmail(email string) CompareResult {
	return validate.ValidateEmail(email)
}

// ValidateIPAddress 校验 IP 地址格式
func ValidateIPAddress(ipStr string) CompareResult {
	return validate.ValidateIPAddress(ipStr)
}

// ValidateProtocol 校验 URL 协议是否在允许列表中
func ValidateProtocol(urlStr string, allowedProtocols ...string) CompareResult {
	return validate.ValidateProtocol(urlStr, allowedProtocols...)
}

// ValidateHTTP 校验 HTTP/HTTPS URL
func ValidateHTTP(urlStr string) CompareResult {
	return validate.ValidateHTTP(urlStr)
}

// ValidateWebSocket 校验 WebSocket URL
func ValidateWebSocket(urlStr string) CompareResult {
	return validate.ValidateWebSocket(urlStr)
}

// ValidateUUID 校验 UUID 格式
func ValidateUUID(uuidStr string) CompareResult {
	return validate.ValidateUUID(uuidStr)
}

// ValidateBase64 校验 Base64 字符串
func ValidateBase64(str string) CompareResult {
	return validate.ValidateBase64(str)
}

// IsEmail 判断字符串是否为有效 Email
func IsEmail(email string) bool {
	return validate.IsEmail(email)
}

// IsIP 判断字符串是否为有效 IP
func IsIP(ip string) bool {
	return validate.IsIP(ip)
}

// IsUUID 判断字符串是否为有效 UUID
func IsUUID(uuid string) bool {
	return validate.IsUUID(uuid)
}

// IsBase64 判断字符串是否为有效 Base64
func IsBase64(str string) bool {
	return validate.IsBase64(str)
}

// ────────────────────────────────────────
// 网络校验（validate/network）
// ────────────────────────────────────────

// IPBase 提供面向对象风格的 IP 校验入口
type IPBase = validate.IPBase

// IPSet 表示预编译 IP 规则集合
type IPSet = validate.IPSet

// CompileIPSet 将 IP 规则编译为可复用集合
func CompileIPSet(patterns []string) (*IPSet, error) {
	return validate.CompileIPSet(patterns)
}

// MustCompileIPSet 编译 IP 规则，失败时 panic
func MustCompileIPSet(patterns []string) *IPSet {
	return validate.MustCompileIPSet(patterns)
}

// MatchPathInList 判断路径是否命中任意路径前缀
func MatchPathInList(path string, patterns []string) bool {
	return validate.MatchPathInList(path, patterns)
}

// MatchPathGlob Glob 模式匹配路径（支持 * 和 ? 通配符）
func MatchPathGlob(path, pattern string) bool {
	return validate.MatchPathGlob(path, pattern)
}

// IsIPAllowed 判断 IP 是否在允许列表中
func IsIPAllowed(ip string, cidrList []string) bool {
	return validate.IsIPAllowed(ip, cidrList)
}

// IsIPBlocked 判断 IP 是否在黑名单中
func IsIPBlocked(ip string, blacklist []string) bool {
	return validate.IsIPBlocked(ip, blacklist)
}

// MatchIPPattern 判断 IP 是否命中单个规则
func MatchIPPattern(ip, pattern string) bool {
	return validate.MatchIPPattern(ip, pattern)
}

// MatchIPInList 判断 IP 是否命中任意规则
func MatchIPInList(ip string, ipList []string) bool {
	return validate.MatchIPInList(ip, ipList)
}

// IsIPInRange 判断 IPv4 地址是否落在闭区间内
func IsIPInRange(ip, start, end net.IP) bool {
	return validate.IsIPInRange(ip, start, end)
}

// MatchIPWithWildcard 使用星号通配符匹配 IPv4
func MatchIPWithWildcard(ip, pattern string) bool {
	return validate.MatchIPWithWildcard(ip, pattern)
}

// IsPrivateIP 判断 IP 是否属于常见私有或本地地址段
func IsPrivateIP(ip string) bool {
	return validate.IsPrivateIP(ip)
}

// ────────────────────────────────────────
// JSON 校验（validate/json）
// ────────────────────────────────────────

// ValidateJSON 校验 JSON 字节是否有效
func ValidateJSON(data []byte) error {
	return validate.ValidateJSON(data)
}

// IsJSONNull 判断 JSON 字节是否为 null
func IsJSONNull(data []byte) bool {
	return validate.IsJSONNull(data)
}

// IsJSONColumnType 判断数据库列类型是否属于 JSON 类型
func IsJSONColumnType(dbType string) bool {
	return validate.IsJSONColumnType(dbType)
}

// ValidateJSONWithData 校验 JSON 并返回反序列化数据
func ValidateJSONWithData(body []byte) (any, error) {
	return validate.ValidateJSONWithData(body)
}

// ValidateJSONField 校验 JSON 顶层字段是否等于期望值
func ValidateJSONField(body []byte, field string, expected any) CompareResult {
	return validate.ValidateJSONField(body, field, expected)
}

// ValidateJSONFields 批量校验 JSON 顶层字段
func ValidateJSONFields(body []byte, rules map[string]any) []CompareResult {
	return validate.ValidateJSONFields(body, rules)
}

// LookupJSONPath 按轻量路径读取 JSON 数据
func LookupJSONPath(data any, path string) (any, bool) {
	return validate.LookupJSONPath(data, path)
}

// ValidateJSONPath 校验 JSON 路径的值
func ValidateJSONPath(body []byte, jsonPath string, expected any, op CompareOperator) CompareResult {
	return validate.ValidateJSONPath(body, jsonPath, expected, op)
}

// ValidateJSONPathExists 校验 JSON 路径是否存在
func ValidateJSONPathExists(body []byte, jsonPath string) CompareResult {
	return validate.ValidateJSONPathExists(body, jsonPath)
}

// ────────────────────────────────────────
// 枚举校验（validate/enum）
// ────────────────────────────────────────

// NewEnumValidator 创建枚举校验器
func NewEnumValidator[T comparable](values ...T) *validate.EnumValidator[T] {
	return validate.NewEnumValidator(values...)
}

// ────────────────────────────────────────
// 空值与过滤值（validate/empty）
// ────────────────────────────────────────

// IsEmptyValue 判断 reflect.Value 是否为空值
func IsEmptyValue(v reflect.Value) bool {
	return validate.IsEmptyValue(v)
}

// IsTimeEmpty 判断时间是否为空或早于 Unix epoch
func IsTimeEmpty(t *time.Time) bool {
	return validate.IsTimeEmpty(t)
}

// IsTimeValid 判断时间值是否有效
func IsTimeValid(timeVal interface{}) bool {
	return validate.IsTimeValid(timeVal)
}

// HasEmpty 判断切片中是否存在空值
func HasEmpty(elems []interface{}) (bool, int) {
	return validate.HasEmpty(elems)
}

// IsAllEmpty 判断切片中所有元素是否都为空
func IsAllEmpty(elems []interface{}) bool {
	return validate.IsAllEmpty(elems)
}

// IsUndefined 判断字符串是否为 undefined
func IsUndefined(str string) bool {
	return validate.IsUndefined(str)
}

// IsNull 判断字符串是否为 null
func IsNull(str string) bool {
	return validate.IsNull(str)
}

// IfNullOrUndefined 判断字符串是否为 null 或 undefined
func IfNullOrUndefined(str string) bool {
	return validate.IfNullOrUndefined(str)
}

// ContainsChinese 判断字符串是否包含中文
func ContainsChinese(s string) bool {
	return validate.ContainsChinese(s)
}

// EmptyToDefault 在字符串为空时返回默认值
func EmptyToDefault(str string, defaultStr string) string {
	return validate.EmptyToDefault(str, defaultStr)
}

// IsNil 判断 interface 是否为 nil 或内部持有 nil
func IsNil(x interface{}) bool {
	return validate.IsNil(x)
}

// IsFuncType 判断泛型类型是否为函数
func IsFuncType[T any]() bool {
	return validate.IsFuncType[T]()
}

// IsCEmpty 判断可比较值是否为零值
func IsCEmpty[T comparable](v T) bool {
	return validate.IsCEmpty(v)
}

// DerefValue 解开 interface 中的指针值
func DerefValue(value interface{}) (interface{}, bool) {
	return validate.DerefValue(value)
}

// IsSafeFieldName 判断字段名是否安全
func IsSafeFieldName(field string) bool {
	return validate.IsSafeFieldName(field)
}

// IsAllowedField 判断字段是否在白名单中
func IsAllowedField(field string, allowedFields ...[]string) bool {
	return validate.IsAllowedField(field, allowedFields...)
}

// UnwrapProtobufWrapper 通过反射解开 protobuf wrapper
func UnwrapProtobufWrapper(value interface{}) (interface{}, bool) {
	return validate.UnwrapProtobufWrapper(value)
}

// IsEmptyAfterDeref 解引用后判断值是否为空
func IsEmptyAfterDeref(value interface{}) (interface{}, bool) {
	return validate.IsEmptyAfterDeref(value)
}

// NormalizeFilterValue 归一化过滤值
func NormalizeFilterValue(value interface{}) interface{} {
	return validate.NormalizeFilterValue(value)
}

// NormalizeFilterValueSlice 归一化过滤值切片
func NormalizeFilterValueSlice(values []interface{}) []interface{} {
	return validate.NormalizeFilterValueSlice(values)
}

// NormalizeFilterValueIfNotEmpty 过滤空值后返回归一化值
func NormalizeFilterValueIfNotEmpty(value interface{}) (interface{}, bool) {
	return validate.NormalizeFilterValueIfNotEmpty(value)
}

// ────────────────────────────────────────
// JSON 扫描（validate/json）
// ────────────────────────────────────────

// SkipJSONSpaces 跳过 JSON 字节流中的空白字符
func SkipJSONSpaces(data []byte, i int) int {
	return validate.SkipJSONSpaces(data, i)
}

// ScanJSONString 扫描 JSON 字符串，返回字符串结束后一位的位置
func ScanJSONString(data []byte, start int) (int, error) {
	return validate.ScanJSONString(data, start)
}

// ScanJSONValueEnd 扫描任意 JSON 值，返回值结束后一位的位置
func ScanJSONValueEnd(data []byte, start int) (int, error) {
	return validate.ScanJSONValueEnd(data, start)
}

// ────────────────────────────────────────
// 正则缓存（validate/format）
// ────────────────────────────────────────

// GetCompiledRegex 获取编译的正则（带缓存）
func GetCompiledRegex(pattern string) (*regexp.Regexp, error) {
	return validate.GetCompiledRegex(pattern)
}

// ClearRegexCache 清空正则缓存
func ClearRegexCache() {
	validate.ClearRegexCache()
}

// ValidateIP 校验 IP 地址格式（兼容旧 API）
func ValidateIP(ipStr string) CompareResult {
	return validate.ValidateIP(ipStr)
}

// ────────────────────────────────────────
// Schema 校验（schema 子包）
// ────────────────────────────────────────

// JSONSchema 描述 Argus 支持的 JSON Schema 子集
type JSONSchema = schema.JSONSchema

// SchemaBuilder 提供链式构建 JSONSchema 的能力
type SchemaBuilder = schema.SchemaBuilder

// ValidateJSONSchema 校验数据是否符合 schema
func ValidateJSONSchema(data interface{}, s interface{}) CompareResult {
	return schema.ValidateJSONSchema(data, s)
}

// ValidateStructWithSchema 校验结构体或 map 是否符合 schema
func ValidateStructWithSchema(structData interface{}, s interface{}) CompareResult {
	return schema.ValidateStructWithSchema(structData, s)
}

// NewSchemaBuilder 创建 SchemaBuilder
func NewSchemaBuilder() *SchemaBuilder {
	return schema.NewSchemaBuilder()
}

// QuickSchema 根据字段类型快速创建对象 schema
func QuickSchema(properties map[string]string, required ...string) JSONSchema {
	return schema.QuickSchema(properties, required...)
}

// FormatSchemaError 提取 schema 校验错误消息
func FormatSchemaError(result CompareResult) string {
	return schema.FormatSchemaError(result)
}
