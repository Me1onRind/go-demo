//go:generate stringer -type=JobBackendType
package async

type JobBackendType int

const (
	KafkaBackendJob JobBackendType = 1
	RedisBackendJob JobBackendType = 2
)
