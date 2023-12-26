/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-06 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-06 00:00:00
 * @FilePath: \go-argus\schema\schema.go
 * @Description: JSON Schema 子集校验模块，提供零依赖结构化数据校验能力
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

// Package schema 提供轻量 JSON Schema 子集校验能力
package schema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kamalyes/go-argus/i18n"
	"github.com/kamalyes/go-argus/validate"
)

// JSONSchema 描述 Argus 支持的 JSON Schema 子集
type JSONSchema struct {
	Type                 string                `json:"type,omitempty"`
	Required             []string              `json:"required,omitempty"`
	Properties           map[string]JSONSchema `json:"properties,omitempty"`
	Items                *JSONSchema           `json:"items,omitempty"`
	Enum                 []interface{}         `json:"enum,omitempty"`
	MinLength            *int                  `json:"minLength,omitempty"`
	MaxLength            *int                  `json:"maxLength,omitempty"`
	Minimum              *float64              `json:"minimum,omitempty"`
	Maximum              *float64              `json:"maximum,omitempty"`
	AdditionalProperties bool                  `json:"additionalProperties,omitempty"`
}

// ValidateJSONSchema 校验数据是否符合 schema
func ValidateJSONSchema(data interface{}, schema interface{}) validate.CompareResult {
	result := validate.CompareResult{Actual: fmt.Sprint(data), Expect: "valid JSON schema"}
	compiled, err := normalizeSchema(schema)
	if err != nil {
		result.Message = err.Error()
		return result
	}
	if err := validateValue(data, compiled, "$"); err != nil {
		result.Message = err.Error()
		return result
	}
	result.Success = true
	return result
}

// ValidateStructWithSchema 校验结构体或 map 是否符合 schema
func ValidateStructWithSchema(structData interface{}, schema interface{}) validate.CompareResult {
	raw, err := json.Marshal(structData)
	if err != nil {
		return validate.CompareResult{Message: err.Error()}
	}
	var data interface{}
	json.Unmarshal(raw, &data)
	return ValidateJSONSchema(data, schema)
}

// SchemaBuilder 提供链式构建 JSONSchema 的能力
type SchemaBuilder struct {
	schema JSONSchema
}

// NewSchemaBuilder 创建 SchemaBuilder
func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{schema: JSONSchema{Properties: make(map[string]JSONSchema)}}
}

// Type 设置 schema 类型
func (b *SchemaBuilder) Type(t string) *SchemaBuilder {
	b.schema.Type = t
	return b
}

// Required 设置必填字段
func (b *SchemaBuilder) Required(fields ...string) *SchemaBuilder {
	b.schema.Required = append(b.schema.Required, fields...)
	return b
}

// Property 增加属性 schema
func (b *SchemaBuilder) Property(name string, propSchema JSONSchema) *SchemaBuilder {
	if b.schema.Properties == nil {
		b.schema.Properties = make(map[string]JSONSchema)
	}
	b.schema.Properties[name] = propSchema
	return b
}

// StringProperty 增加字符串属性
func (b *SchemaBuilder) StringProperty(name string, minLen, maxLen int) *SchemaBuilder {
	prop := JSONSchema{Type: "string"}
	if minLen >= 0 {
		prop.MinLength = &minLen
	}
	if maxLen >= 0 {
		prop.MaxLength = &maxLen
	}
	return b.Property(name, prop)
}

// NumberProperty 增加数值属性
func (b *SchemaBuilder) NumberProperty(name string, min, max *float64) *SchemaBuilder {
	return b.Property(name, JSONSchema{Type: "number", Minimum: min, Maximum: max})
}

// ArrayProperty 增加数组属性
func (b *SchemaBuilder) ArrayProperty(name string, items JSONSchema) *SchemaBuilder {
	return b.Property(name, JSONSchema{Type: "array", Items: &items})
}

// Enum 设置枚举值
func (b *SchemaBuilder) Enum(values ...interface{}) *SchemaBuilder {
	b.schema.Enum = append(b.schema.Enum, values...)
	return b
}

// Build 返回 JSONSchema
func (b *SchemaBuilder) Build() JSONSchema {
	return b.schema
}

// BuildJSON 返回 JSON 字符串
func (b *SchemaBuilder) BuildJSON() string {
	raw, _ := json.Marshal(b.schema)
	return string(raw)
}

// QuickSchema 根据字段类型快速创建对象 schema
func QuickSchema(properties map[string]string, required ...string) JSONSchema {
	s := JSONSchema{Type: "object", Required: required, Properties: make(map[string]JSONSchema, len(properties))}
	for name, typ := range properties {
		s.Properties[name] = JSONSchema{Type: typ}
	}
	return s
}

// FormatSchemaError 提取 schema 校验错误消息
func FormatSchemaError(result validate.CompareResult) string {
	if result.Success {
		return ""
	}
	return result.Message
}

func normalizeSchema(schema interface{}) (JSONSchema, error) {
	switch v := schema.(type) {
	case JSONSchema:
		return v, nil
	case *JSONSchema:
		if v == nil {
			return JSONSchema{}, fmt.Errorf(i18n.Msg(validate.MsgSchemaEmpty))
		}
		return *v, nil
	case []byte:
		var s JSONSchema
		return s, json.Unmarshal(v, &s)
	case string:
		var s JSONSchema
		return s, json.Unmarshal([]byte(v), &s)
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return JSONSchema{}, err
		}
		var s JSONSchema
		return s, json.Unmarshal(raw, &s)
	}
}

func validateValue(value interface{}, schema JSONSchema, path string) error {
	if schema.Type != "" && !matchesType(value, schema.Type) {
		return fmt.Errorf(i18n.Msg(validate.MsgSchemaTypeMismatch, map[string]string{"path": path, "type": schema.Type}))
	}
	if len(schema.Enum) > 0 && !containsEnum(value, schema.Enum) {
		return fmt.Errorf(i18n.Msg(validate.MsgSchemaEnumMismatch, map[string]string{"path": path}))
	}
	switch schema.Type {
	case "string":
		return validateString(value, schema, path)
	case "number", "integer":
		return validateNumber(value, schema, path)
	case "array":
		return validateArray(value, schema, path)
	case "object":
		return validateObject(value, schema, path)
	default:
		if schema.Properties != nil {
			return validateObject(value, schema, path)
		}
		return nil
	}
}

func validateString(value interface{}, schema JSONSchema, path string) error {
	s := value.(string)
	if schema.MinLength != nil && len([]rune(s)) < *schema.MinLength {
		return fmt.Errorf(i18n.Msg(validate.MsgSchemaStringMinLength, map[string]string{"path": path, "min": fmt.Sprint(*schema.MinLength)}))
	}
	if schema.MaxLength != nil && len([]rune(s)) > *schema.MaxLength {
		return fmt.Errorf(i18n.Msg(validate.MsgSchemaStringMaxLength, map[string]string{"path": path, "max": fmt.Sprint(*schema.MaxLength)}))
	}
	return nil
}

func validateNumber(value interface{}, schema JSONSchema, path string) error {
	n := value.(float64)
	if schema.Minimum != nil && n < *schema.Minimum {
		return fmt.Errorf(i18n.Msg(validate.MsgSchemaNumberBelowMinimum, map[string]string{"path": path, "min": fmt.Sprint(*schema.Minimum)}))
	}
	if schema.Maximum != nil && n > *schema.Maximum {
		return fmt.Errorf(i18n.Msg(validate.MsgSchemaNumberAboveMaximum, map[string]string{"path": path, "max": fmt.Sprint(*schema.Maximum)}))
	}
	return nil
}

func validateArray(value interface{}, schema JSONSchema, path string) error {
	arr, ok := value.([]interface{})
	if !ok || schema.Items == nil {
		return nil
	}
	for i, item := range arr {
		if err := validateValue(item, *schema.Items, fmt.Sprintf("%s[%d]", path, i)); err != nil {
			return err
		}
	}
	return nil
}

func validateObject(value interface{}, schema JSONSchema, path string) error {
	obj := value.(map[string]interface{})
	for _, field := range schema.Required {
		if _, ok := obj[field]; !ok {
			return fmt.Errorf(i18n.Msg(validate.MsgSchemaFieldRequired, map[string]string{"path": path, "field": field}))
		}
	}
	for name, prop := range schema.Properties {
		if item, ok := obj[name]; ok {
			if err := validateValue(item, prop, strings.TrimSuffix(path+".", ".")+"."+name); err != nil {
				return err
			}
		}
	}
	return nil
}

func matchesType(value interface{}, typ string) bool {
	switch typ {
	case "object":
		_, ok := value.(map[string]interface{})
		return ok
	case "array":
		_, ok := value.([]interface{})
		return ok
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		_, ok := value.(float64)
		return ok
	case "integer":
		n, ok := value.(float64)
		return ok && n == float64(int64(n))
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "null":
		return value == nil
	default:
		return true
	}
}

func containsEnum(value interface{}, values []interface{}) bool {
	actual := fmt.Sprint(value)
	for _, item := range values {
		if actual == fmt.Sprint(item) {
			return true
		}
	}
	return false
}
