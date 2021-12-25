package gateway

import (
	"reflect"
)

func parserProtocol(paramType interface{}) interface{} {
	valueType := reflect.TypeOf(paramType)
	switch valueType.Kind() {
	case reflect.Ptr:
		return reflect.New(valueType.Elem()).Interface()
	case reflect.Struct:
		return reflect.New(valueType).Interface()
	default:
		return nil
	}
}
