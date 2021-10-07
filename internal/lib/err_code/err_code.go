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
	SUCCESS = NewError(Sucess, "Success", codes.OK)

	ServerInternalError = NewError(ServerInternal, "Server Internal Error", codes.Aborted)
	InvalidParamError   = NewError(InvalidParam, "Invalid Param", codes.Aborted)
	GRPCCallFailedError = NewError(GRPCCallFailed, "GRPC Call Failed", codes.Unavailable)
	JsonEncodeError     = NewError(JsonEncodeFailed, "Json Encode Failed", codes.Aborted)
	JsonDecodeError     = NewError(JsonDecodeFailed, "Json Decode Failed", codes.Aborted)
	AsyncTaskSendError  = NewError(AsyncTaskSendFailed, "Async Task Send Failed", codes.Unavailable)
	ReadDBError         = NewError(ReadDBFailed, "Read DB Failed", codes.Unavailable)
	WriteDBError        = NewError(WriteDBFailed, "Write DB Failed", codes.Unavailable)
	ReadRedisError      = NewError(ReadRedisFailed, "Read Redis Failed", codes.Unavailable)
	WriteRedisError     = NewError(WriteRedisFailed, "Write Redis Failed", codes.Unavailable)
	DBRecordExistError  = NewError(DBRecordExist, "DB Record Exist", codes.Aborted)
	GetLocalcacheError  = NewError(GetLocalcacheFailed, "Get Localcache Failed", codes.Aborted)

	GoDemoCommonFailedError     = NewError(GoDemoCommonFailed, "Call Go-demo Failed", codes.Aborted)
	KvConfigKeyNotExistError    = NewError(KvConfigKeyNotExist, "Config Key Not Exist", codes.Aborted)
	KvConfigValueWrongTypeError = NewError(KvConfigValueWrongType, "Config Value Type Is Wrong", codes.Aborted)
)
