package err_code

import (
	"github.com/Me1onRind/go-demo/internal/core/common"
	"google.golang.org/grpc/codes"
)

const (
	Sucess = 0

	ServerInternal      = 10000
	InvalidParam        = 10001
	GRPCCallFailed      = 10002
	JsonEncodeFailed    = 10003
	JsonDecodeFailed    = 10004
	AsyncTaskSendFailed = 10005

	GoDemoCommonFailed = 20001
)

var (
	SUCCESS = common.NewError(Sucess, "Success", codes.OK)

	ServerInternalError = common.NewError(ServerInternal, "Server Internal Error", codes.Aborted)
	InvalidParamError   = common.NewError(InvalidParam, "Invalid Param", codes.Aborted)
	GRPCCallFailedError = common.NewError(GRPCCallFailed, "GRPC Call Failed", codes.Unavailable)
	JsonEncodeError     = common.NewError(JsonEncodeFailed, "Json Encode Failed", codes.Aborted)
	JsonDecodeError     = common.NewError(JsonDecodeFailed, "Json Decode Failed", codes.Aborted)
	AsyncTaskSendError  = common.NewError(AsyncTaskSendFailed, "Async Task Send Failed", codes.Unavailable)

	GoDemoCommonFailedError = common.NewError(GoDemoCommonFailed, "Call Go-demo Failed", codes.Aborted)
)
