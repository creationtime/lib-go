package errors

import (
	"github.com/pquerna/ffjson/ffjson"
)

const (
	ErrTokenInvalid = 40401 // Token无效

	ErrTimeout        = 40501 // 服务超时
	ErrServiceInvalid = 40503 // 服务调用错误

	ErrInternal      = 40500 // 服务器内部错误
	ErrDatabase      = 40510 // 数据库处理错误
	ErrNotFound      = 40524 // 请求记录不存在
	ErrParamsInvalid = 40520 // 请求参数错误

	ErrMarshalError   = 40530 // 参数Marshal失败
	ErrUnMarshalError = 40531 // 参数UnMarshal失败

)

// Error implements the error interface.
type Error struct {
	Code   int32  `jsonpb:"code"`
	Detail string `jsonpb:"detail"`
}

func (e *Error) Error() string {
	b, _ := ffjson.Marshal(e)
	return string(b)
}

// New generates a custom error.
func New(code int32, message string) error {
	return &Error{
		Code:   code,
		Detail: message,
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	err2 := ffjson.Unmarshal([]byte(err), e)
	if err2 != nil {
		e.Code = ErrInternal
		e.Detail = err
	}
	return e
}

func Unauthorized(message string) error {
	return &Error{
		Code:   ErrTokenInvalid,
		Detail: message,
	}
}

func DataBaseError(message string) error {
	return &Error{
		Code:   ErrDatabase,
		Detail: message,
	}
}

func NotFound(message string) error {
	return &Error{
		Code:   ErrNotFound,
		Detail: message,
	}
}

func BadRequest(message string) error {
	return &Error{
		Code:   ErrParamsInvalid,
		Detail: message,
	}
}

func RequestTimeout(message string) error {
	return &Error{
		Code:   ErrTimeout,
		Detail: message,
	}
}

func InternalServerError(message string) error {
	return &Error{
		Code:   ErrInternal,
		Detail: message,
	}
}

func ServerUnavailable(message string) error {
	return &Error{
		Code:   ErrServiceInvalid,
		Detail: message,
	}
}
