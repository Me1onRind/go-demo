// Code generated by "stringer -type=IdType"; DO NOT EDIT.

package idpo

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UserIdType-1]
}

const _IdType_name = "UserIdType"

var _IdType_index = [...]uint8{0, 10}

func (i IdType) String() string {
	i -= 1
	if i < 0 || i >= IdType(len(_IdType_index)-1) {
		return "IdType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _IdType_name[_IdType_index[i]:_IdType_index[i+1]]
}
