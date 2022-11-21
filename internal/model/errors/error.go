package errors

import (
	"fmt"
)

type Error struct {
	Code int32

	message string
	cause   error
}

func NewError(code int32, message string) *Error {
	e := &Error{
		Code:    code,
		message: message,
	}
	return e
}

func (e *Error) Error() string {
	if e.cause == nil {
		return e.message
	}

	return fmt.Sprintf("%s, cause:[%s]", e.message, e.cause.Error())
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Wrap(err error) error {
	newE := *e
	newE.cause = err
	return &newE
}
