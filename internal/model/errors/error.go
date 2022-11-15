package errors

import (
	"errors"
	"fmt"
)

var (
	duplicateCheck               = map[int32]*Error{}
	DuplicateRegisterError error = errors.New("Code has benn register")
)

type Error struct {
	Code int32

	message string
	cause   error
}

func NewUniqCodeError(code int32, message string) (*Error, error) {
	if value, ok := duplicateCheck[code]; ok {
		return nil, fmt.Errorf("Code:[%d] uniq new failed:[%w], exist value:[{message:%s}]", code, DuplicateRegisterError, value.message)
	}

	e := &Error{
		Code:    code,
		message: message,
	}
	duplicateCheck[code] = e
	return e, nil
}

func (e *Error) Error() string {
	if e.cause == nil {
		return e.message
	}

	return fmt.Sprintf("%s, cause:[%s]", e.message, e.cause.Error())
}

func (e *Error) Warp(err error) {
	e.cause = err
}

func (e *Error) Unwrap() error {
	return e.cause
}
