package httpclient

import (
	"context"
	"time"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type HTTPClient struct {
	fastClient fasthttp.Client
}

func NewHttpClient() *HTTPClient {
	h := &HTTPClient{
		fastClient: fasthttp.Client{},
	}
	return h
}

func (h *HTTPClient) SendJsonRequest(ctx context.Context, uri string, req, resp any, opts ...Option) error {
	var (
		reqBody []byte
		err     error
	)
	r := fasthttp.AcquireRequest()
	r.SetRequestURI(uri)
	r.Header.SetMethod(fasthttp.MethodPost)
	r.Header.SetContentType("application/json")
	rp := fasthttp.AcquireResponse()
	startTime := time.Now()
	defer func() {
		if err == nil {
			logger.CtxInfof(ctx, "uri:[%s],method:[POST],req:[%s],resp:[%s],duration:[%s]",
				uri, reqBody, rp.Body(), time.Since(startTime))
		} else {
			logger.CtxInfof(ctx, "uri:[%s],method:[POST],req:[%s],resp:[%s],duration:[%s],err:[%s]",
				uri, reqBody, rp.Body(), time.Since(startTime), err)
		}
		fasthttp.ReleaseRequest(r)
		fasthttp.ReleaseResponse(rp)
	}()

	if req != nil {
		reqBody, err = jsoniter.Marshal(req)
		if err != nil {
			return err
		}
		r.SetBody(reqBody)
	}

	for _, opt := range opts {
		opt(r)
	}

	err = h.fastClient.Do(r, rp)
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(rp.Body(), resp)
	if err != nil {
		return err
	}

	return nil
}
