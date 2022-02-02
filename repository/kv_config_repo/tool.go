package kv_config_dao

import (
	json "github.com/json-iterator/go"
	"strconv"
	"strings"
)

func boolean(value string) (realValue bool, ok bool) {
	value = strings.ToLower(value)
	if value == "true" {
		return true, true
	}

	if value == "false" {
		return false, true
	}

	return false, false
}

func integer(value string) (int64, bool) {
	realValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return realValue, true
}

func float(value string) (float64, bool) {
	realValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return realValue, true
}

func dict(value string) (map[string]interface{}, bool) {
	realValue := map[string]interface{}{}
	if err := json.Unmarshal([]byte(value), &realValue); err != nil {
		return nil, false
	}
	return realValue, true
}

func list(value string) ([]interface{}, bool) {
	realValue := []interface{}{}
	if err := json.Unmarshal([]byte(value), &realValue); err != nil {
		return nil, false
	}
	return realValue, true
}
