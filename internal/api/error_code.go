package api

const (
	// InvalidParams 参数错误
	InvalidParams = 10000
	// ObjectNotFound 对象不存在
	ObjectNotFound = iota
	// AuthFail 验证失败
	AuthFail
	// SystemInternalError 系统内部错误
	SystemInternalError
)
