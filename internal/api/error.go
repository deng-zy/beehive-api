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
	// ErrTopicNotExists topic不存在
	ErrTopicNotExists = errors.New("topic not exists")
	// ErrTopicTooLarge topic太多
	ErrTopicTooLarge = errors.New("topic too large(max:100)")
	// ErrPayloadTooLarge 事件消息太长
	ErrPayloadTooLarge = errors.New("payload too large(max:2048)")
	// ErrTopicLengthTooLarge topic长度超过太长
	ErrTopicLengthTooLarge = errors.New("topic length too large(max:64)")
)
