package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	Code int

	message string
	cause   error
}

func NewError(code int, message string) *Error {
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

func (e *Error) Withf(msg string, a ...any) error {
	newE := *e
	newE.message = fmt.Sprintf("%s, %s", e.message, fmt.Sprintf(msg, a...))
	return &newE
}

func ExtractError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}
