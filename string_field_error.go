/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-17 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-17 11:12:25
 * @FilePath: \go-argus\string_field_error.go
 * @Description: Argus 字符串字段错误实现
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"fmt"
	"reflect"
)

type stringFieldError struct {
	tag   string
	param string
	value string
}

var _ FieldError = (*stringFieldError)(nil)

func (e *stringFieldError) Tag() string             { return e.tag }
func (e *stringFieldError) ActualTag() string       { return e.tag }
func (e *stringFieldError) Namespace() string       { return "" }
func (e *stringFieldError) StructNamespace() string { return "" }
func (e *stringFieldError) Field() string           { return "" }
func (e *stringFieldError) StructField() string     { return "" }
func (e *stringFieldError) Value() interface{}      { return e.value }
func (e *stringFieldError) Param() string           { return e.param }
func (e *stringFieldError) Kind() reflect.Kind      { return reflect.String }
func (e *stringFieldError) Type() reflect.Type      { return reflect.TypeOf("") }
func (e *stringFieldError) Error() string {
	return fmt.Sprintf(fieldErrMsg, "", "", e.tag)
}
