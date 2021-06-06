package common

import (
	"fmt"

	"github.com/Me1onRind/go-demo/protobuf/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errcodes = map[int32]struct{}{}

type Error struct {
	Code    int32
	Message string
	Param   string
}

func NewError(code int32, message string) *Error {
	if _, ok := errcodes[code]; ok {
		panic(fmt.Sprintf("Error code:%d exist", code))
	}

	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) String() string {
	if len(e.Param) == 0 {
		return e.Message
	}
	return fmt.Sprintf("%s:%s", e.Message, e.Param)
}

func (e *Error) With(param string) {
	e.Param = param
}

func (e *Error) Withf(format string, v ...interface{}) {
	e.Param = fmt.Sprintf(format, v...)
}

func (e *Error) GrpcErr() error {
	pbErr := &pb.Error{
		Code:    e.Code,
		Message: e.String(),
	}
	s, err := status.New(toGRPCCode(e.Code), pbErr.Message).WithDetails(pbErr)
	if err != nil {
		panic(err)
	}
	return s.Err()
}

func toGRPCCode(code int32) codes.Code {
	var statusCode codes.Code
	switch code {
	default:
		statusCode = codes.Unknown
	}
	return statusCode
}
