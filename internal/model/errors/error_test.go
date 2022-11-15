package errors

import (
	"errors"
	"testing"

	"io"

	"github.com/stretchr/testify/assert"
)

func Test_New_Success(t *testing.T) {
	e, err := NewUniqCodeError(100, "My Error")
	t.Log(e)
	if assert.Empty(t, err) {
		assert.Equal(t, "My Error", e.Error())
	}
}

func Test_New_Duplicate(t *testing.T) {
	_, err := NewUniqCodeError(100, "My Error")
	assert.Equal(t, true, errors.Is(err, DuplicateRegisterError))
	t.Log(err)
}

func Test_Wrap_Success(t *testing.T) {
	message := "My Error"
	e, _ := NewUniqCodeError(200, message)
	e.Warp(io.EOF)
	t.Log(e)
	assert.Equal(t, true, errors.Is(e, io.EOF))
	assert.Equal(t, "My Error, cause:[EOF]", e.Error())
}
