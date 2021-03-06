package err_code

import (
	"errors"
	"fmt"

	"github.com/Me1onRind/go-demo/protobuf/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errcodes = map[int32]*Error{}

type Error struct {
	Code     int32
	Message  string
	Param    string
	grpcCode codes.Code
}

func newError(code int32, message string, grpcCode codes.Code) *Error {
	if _, ok := errcodes[code]; ok {
		panic(fmt.Sprintf("Error code:%d exist", code))
	}

	e := &Error{
		Code:     code,
		Message:  message,
		grpcCode: grpcCode,
	}
	errcodes[code] = e
	return e
}

func (e *Error) clone() *Error {
	return &Error{
		Code:     e.Code,
		Message:  e.Message,
		grpcCode: e.grpcCode,
	}
}

func (e *Error) String() string {
	if len(e.Param) == 0 {
		return e.Message
	}
	return fmt.Sprintf("%s:%s", e.Message, e.Param)
}

func (e *Error) With(param string) *Error {
	nE := e.clone()
	nE.Param = param
	return nE
}

func (e *Error) Withf(format string, v ...interface{}) *Error {
	return e.With(fmt.Sprintf(format, v...))
}

func (e *Error) WithErr(err error) *Error {
	return e.With(err.Error())
}

func (e *Error) GrpcErr() error {
	pbErr := &pb.Error{
		Code:    e.Code,
		Message: e.String(),
	}
	s, err := status.New(e.grpcCode, fmt.Sprintf("Code[%d],Msg[%s]", pbErr.Code, pbErr.Message)).WithDetails(pbErr)
	if err != nil {
		panic(err)
	}
	return s.Err()
}

func (e *Error) GenError() error {
	return errors.New(e.String())
}
