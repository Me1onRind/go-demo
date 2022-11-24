package code

const (
	Success int32 = 0

	Unexpect           = -10000000
	JsonEncodeFail     = -10000001
	JsonDecodeFail     = -10000002
	ProtocolDecodeFail = -10000003
	ReadDBFail         = -10000004
	WriteDBFail        = -10000005
	GenerateIdFail     = -10000006
)

func IsWarning(code int32) bool {
	return code > Success
}
