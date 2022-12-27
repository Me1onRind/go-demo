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

func Test_Etract_Error(t *testing.T) {
	assert.Empty(t, ExtractError(io.EOF))

	e := NewError(200, "test")
	err := e.Wrap(io.EOF)
	extractE := ExtractError(err)
	assert.Equal(t, extractE.Code, 200)
}

func Test_Code_Equal(t *testing.T) {
	assert.Equal(t, false, CodeEqual(io.EOF, 200))
	e := NewError(200, "test")
	assert.Equal(t, true, CodeEqual(e, 200))
	assert.Equal(t, false, CodeEqual(e, 100))
}
