package api

import "errors"

var (
	// ErrInvalidParam 无效参数
	ErrInvalidParam = errors.New("invalid param")
	// ErrObjectNotFound 数据对象不存在
	ErrObjectNotFound = errors.New("object not found")
	// ErrAuthFail 认证失败
	ErrAuthFail = errors.New("auth failed")
	// ErrTokenIsEmpty token为空
	ErrTokenIsEmpty = errors.New("token is empty")
)
