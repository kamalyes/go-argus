/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-12-16 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-12-16 00:00:00
 * @FilePath: \go-argus\validate\constants.go
 * @Description: i18n 消息键常量，统一维护所有国际化消息的 key，避免硬编码字符串散落各处
 *
 * Copyright (c) 2023 by kamalyes, All Rights Reserved.
 */

package validate

// compare 系列消息键 —— 数值、字符串、状态码比较相关
const (
	MsgCompareUnsupportedNumberOp = "compare.unsupported_number_op"
	MsgCompareNumberFailed        = "compare.number_failed"
	MsgCompareUnsupportedStringOp = "compare.unsupported_string_op"
	MsgCompareStringFailed        = "compare.string_failed"
	MsgCompareRegexCompileFailed  = "compare.regex_compile_failed"
	MsgCompareStatusOutOfRange    = "compare.status_out_of_range"
)

// format 系列消息键 —— Email、IP、URL、UUID、Base64、正则格式校验相关
const (
	MsgFormatRegexCompileFailed     = "format.regex_compile_failed"
	MsgFormatRegexNotMatched        = "format.regex_not_matched"
	MsgFormatEmailEmpty             = "format.email_empty"
	MsgFormatEmailInvalid           = "format.email_invalid"
	MsgFormatEmailMalformed         = "format.email_malformed"
	MsgFormatIPInvalid              = "format.ip_invalid"
	MsgFormatURLMissingProtocol     = "format.url_missing_protocol"
	MsgFormatURLUnsupportedProtocol = "format.url_unsupported_protocol"
	MsgFormatUUIDInvalid            = "format.uuid_invalid"
	MsgFormatBase64Empty            = "format.base64_empty"
	MsgFormatBase64Invalid          = "format.base64_invalid"
)

// enum 系列消息键 —— 枚举校验相关
const (
	MsgEnumInvalidValue = "enum.invalid_value"
)

// json 系列消息键 —— JSON 校验相关
const (
	MsgJSONInvalid            = "json.invalid"
	MsgJSONRootNotObject      = "json.root_not_object"
	MsgJSONFieldNotFound      = "json.field_not_found"
	MsgJSONFieldValueMismatch = "json.field_value_mismatch"
	MsgJSONPathNotFound       = "json.path_not_found"
)

// network 系列消息键 —— 网络、IP、CIDR 校验相关
const (
	MsgNetworkIPInvalid      = "network.ip_invalid"
	MsgNetworkIPRangeInvalid = "network.ip_range_invalid"
	MsgNetworkIPRuleInvalid  = "network.ip_rule_invalid"
)

// schema 系列消息键 —— JSON Schema 校验相关
const (
	MsgSchemaEmpty              = "schema.empty"
	MsgSchemaTypeMismatch       = "schema.type_mismatch"
	MsgSchemaEnumMismatch       = "schema.enum_mismatch"
	MsgSchemaStringMinLength    = "schema.string_min_length"
	MsgSchemaStringMaxLength    = "schema.string_max_length"
	MsgSchemaNumberBelowMinimum = "schema.number_below_minimum"
	MsgSchemaNumberAboveMaximum = "schema.number_above_maximum"
	MsgSchemaFieldRequired      = "schema.field_required"
)
