package base

import "context"

type Context struct {
	context.Context
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
	}
}
