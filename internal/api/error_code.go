package api

const (
	// InvalidParams 参数错误
	InvalidParams = 10000
	// ObjectNotFound 对象不存在
	ObjectNotFound = 10001
	// AuthFail 验证失败
	AuthFail = 10002
	// SystemInternalError 系统内部错误
	SystemInternalError = 10003
	// TopicNotExists topic不存在
	TopicNotExists = 10004
	// TopicTooLarge topic太多
	TopicTooLarge = 10005
	// PayloadTooLarge payload 长度超过限制
	PayloadTooLarge = 10006
	// TopicLengthTooLarge topic长度太长
	TopicLengthTooLarge = 10007
)
