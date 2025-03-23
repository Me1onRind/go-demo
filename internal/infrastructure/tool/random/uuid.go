package random

import (
	"github.com/google/uuid"
)

var (
	uuidData = []byte("xxx")
)

func UUID() string {
	namespace := uuid.New()
	return uuid.NewSHA1(namespace, uuidData).String()
}
