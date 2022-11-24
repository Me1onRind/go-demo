package gerror

import (
	"github.com/Me1onRind/go-demo/internal/model/errors"
	"github.com/Me1onRind/go-demo/protocol/code"
)

var (
	ReadDBError     = errors.NewError(code.ReadDBFail, "Read DB Fail")
	WriteDBError    = errors.NewError(code.WriteDBFail, "Write DB Fail")
	GenerateIdError = errors.NewError(code.GenerateIdFail, "Generate Id Fail")
)
