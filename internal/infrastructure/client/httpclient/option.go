package httpclient

import (
	"time"

	"github.com/valyala/fasthttp"
)

type Option func(r *fasthttp.Request)

func WithTimeout(t time.Duration) Option {
	return func(r *fasthttp.Request) {
		r.SetTimeout(t)
	}
}
