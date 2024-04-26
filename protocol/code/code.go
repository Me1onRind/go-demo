package code

const (
	Success   int = 0
	Duplicate int = 1

	Unexpect           = -100000
	JsonEncodeFail     = -100001
	JsonDecodeFail     = -100002
	ProtocolDecodeFail = -100003
	ReadDBFail         = -100004
	WriteDBFail        = -100005
	GenerateIdFail     = -100006
	InvalidJobProtocol = -100007
	SendKafkaFail      = -100008

	RecordExisted = -200000
)

func IsWarning(code int) bool {
	return code > Success
}
