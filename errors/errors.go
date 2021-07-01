package errors

import (
	"github.com/pquerna/ffjson/ffjson"
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
		Code:   DatabaseCode,
		Detail: message,
	}
}

func NotFound(message string) error {
	return &Error{
		Code:   NotFoundCode,
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
