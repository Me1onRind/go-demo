//go:generate stringer -type=IdType
package idpo

type IdType uint32

const (
	UserIdType IdType = 1
)
