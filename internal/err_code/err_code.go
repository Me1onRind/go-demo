package err_code

import "github.com/Me1onRind/go-demo/internal/core/common"

const (
	Sucess = 0

	ServerInternal = 10001
)

var (
	ServerInternalError = common.NewError(ServerInternal, "Server Internal Error")
)
