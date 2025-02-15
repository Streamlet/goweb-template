package core

type Error struct {
	Code    int
	Message string
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

const (
	// 通用错误码
	Error_Internal             = 1 // 内部错误
	Error_DbError              = 2 // 数据库错误
	Error_JsonEncodeError      = 3 // JSON 编码错误
	Error_JsonDecodeError      = 4 // JSON 解码错误
	Error_RequestTooFrequently = 5 // 请求过于频繁
	Error_InvalidArgument      = 6 // 参数错误
)
