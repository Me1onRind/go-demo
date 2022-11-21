package errors

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Wrap_Success(t *testing.T) {
	message := "My Error"
	e := NewError(200, message)
	newE := e.Wrap(io.EOF)
	t.Log(newE)
	assert.Equal(t, true, errors.Is(newE, io.EOF))
	assert.Equal(t, "My Error, cause:[EOF]", newE.Error())
}

func Test_Errors_As(t *testing.T) {
	e := NewError(200, "test")
	err := e.Wrap(io.EOF)

	var customErr *Error
	assert.Equal(t, true, errors.As(err, &customErr))
	assert.EqualValues(t, 200, customErr.Code)
}
