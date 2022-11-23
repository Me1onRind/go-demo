//go:generate stringer -type=IdType
package idpo

type IdType int32

const (
	UserIdType IdType = 1
)
