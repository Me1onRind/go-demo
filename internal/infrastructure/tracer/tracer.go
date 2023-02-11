package tracer

import (
	"context"
	opentracing "github.com/opentracing/opentracing-go"
)

type spanKey struct{}

func WithSpan(ctx context.Context, span opentracing.Span) context.Context {
	return context.WithValue(ctx, spanKey{}, span)
}
