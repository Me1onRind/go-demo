package coder

import (
	json "github.com/json-iterator/go"
)

func JsonMarshal(v interface{}) string {
	result, _ := json.MarshalToString(v)
	return result
}
