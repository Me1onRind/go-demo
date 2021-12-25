package err_code

import (
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
	ReadDBFailed        = 10006
	WriteDBFailed       = 10007
	ReadRedisFailed     = 10008
	WriteRedisFailed    = 10009
	DBRecordExist       = 10010
	DBRecordNoExist     = 10011

	GetLocalcacheFailed = 10012

	GoDemoCommonFailed     = 20001
	KvConfigKeyNotExist    = 20002
	KvConfigValueWrongType = 20003
)

var (
	SUCCESS = newError(Sucess, "Success", codes.OK)

	ServerInternalError = newError(ServerInternal, "Server Internal Error", codes.Aborted)
	InvalidParamError   = newError(InvalidParam, "Invalid Param", codes.Aborted)
	GRPCCallFailedError = newError(GRPCCallFailed, "GRPC Call Failed", codes.Unavailable)
	JsonEncodeError     = newError(JsonEncodeFailed, "Json Encode Failed", codes.Aborted)
	JsonDecodeError     = newError(JsonDecodeFailed, "Json Decode Failed", codes.Aborted)
	AsyncTaskSendError  = newError(AsyncTaskSendFailed, "Async Task Send Failed", codes.Unavailable)
	ReadDBError         = newError(ReadDBFailed, "Read DB Failed", codes.Unavailable)
	WriteDBError        = newError(WriteDBFailed, "Write DB Failed", codes.Unavailable)
	ReadRedisError      = newError(ReadRedisFailed, "Read Redis Failed", codes.Unavailable)
	WriteRedisError     = newError(WriteRedisFailed, "Write Redis Failed", codes.Unavailable)
	DBRecordExistError  = newError(DBRecordExist, "DB Record Exist", codes.Aborted)
	GetLocalcacheError  = newError(GetLocalcacheFailed, "Get Localcache Failed", codes.Aborted)

	GoDemoCommonFailedError     = newError(GoDemoCommonFailed, "Call Go-demo Failed", codes.Aborted)
	KvConfigKeyNotExistError    = newError(KvConfigKeyNotExist, "Config Key Not Exist", codes.Aborted)
	KvConfigValueWrongTypeError = newError(KvConfigValueWrongType, "Config Value Type Is Wrong", codes.Aborted)
)
