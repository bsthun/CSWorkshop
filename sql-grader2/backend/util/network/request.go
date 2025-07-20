package network

import (
	"bytes"
	"fmt"
	"github.com/bsthun/gut"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"net/url"
)

func ConvertRequest(ctx *fasthttp.RequestCtx) *http.Request {
	defer func() {
		if err := recover(); err != nil {
			// TODO: use sentry logger
			gut.Debug(err)
		}
	}()

	r := new(http.Request)

	r.Method = string(ctx.Method())
	uri := ctx.URI()
	// * ignore error
	r.URL, _ = url.Parse(fmt.Sprintf("%s://%s%s", uri.Scheme(), uri.Host(), uri.Path()))

	// * headers
	r.Header = make(http.Header)
	r.Header.Add("Host", string(ctx.Host()))
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		r.Header.Add(string(key), string(value))
	})
	r.Host = string(ctx.Host())

	// * cookies
	ctx.Request.Header.VisitAllCookie(func(key, value []byte) {
		r.AddCookie(&http.Cookie{Name: string(key), Value: string(value)})
	})

	// * env
	r.RemoteAddr = ctx.RemoteAddr().String()

	// * query string
	r.URL.RawQuery = string(ctx.URI().QueryString())

	// * body
	r.Body = io.NopCloser(bytes.NewReader(ctx.Request.Body()))

	return r
}
