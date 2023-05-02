//go:generate stringer -type=JobBackendType
package async

type JobBackendType int

const (
	KafkaBackendJob JobBackendType = 1
)
