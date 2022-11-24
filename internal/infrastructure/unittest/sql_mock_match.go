package unittest

import (
	"database/sql/driver"
	"github.com/Me1onRind/go-demo/internal/model/generic"
	"time"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type Greater[T generic.Integer] struct {
	compareValue T
}

func NewGreater[T generic.Integer](compareValue T) *Greater[T] {
	return &Greater[T]{
		compareValue: compareValue,
	}
}

func (g Greater[T]) Match(v driver.Value) bool {
	value, ok := v.(T)
	if !ok {
		return false
	}
	return value > g.compareValue
}

type Now struct{}

func (n Now) Match(v driver.Value) bool {
	now := time.Now().Unix()
	if value, ok := v.(int64); ok {
		return now-value < 5 || value-now < 5
	}
	return false
}
