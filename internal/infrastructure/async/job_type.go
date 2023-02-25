//go:generate stringer -type=BackendType
package async

type BackendType uint8

const (
	KafkaBackend BackendType = 1
)
