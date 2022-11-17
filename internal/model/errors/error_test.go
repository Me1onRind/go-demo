package errors

import (
	"errors"
	"testing"

	"io"

	"github.com/stretchr/testify/assert"
)

func Test_Wrap_Success(t *testing.T) {
	message := "My Error"
	e := NewError(200, message)
	e.Warp(io.EOF)
	t.Log(e)
	assert.Equal(t, true, errors.Is(e, io.EOF))
	assert.Equal(t, "My Error, cause:[EOF]", e.Error())
}
