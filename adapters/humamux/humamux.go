package humagmux

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/queryparam"
	"github.com/gorilla/mux"
)

type gmuxContext struct {
	op *huma.Operation
	r  *http.Request
	w  http.ResponseWriter
}

func (ctx *gmuxContext) GetOperation() *huma.Operation {
	return ctx.op
}

func (ctx *gmuxContext) GetContext() context.Context {
	return ctx.r.Context()
}

func (ctx *gmuxContext) GetMethod() string {
	return ctx.r.Method
}

func (ctx *gmuxContext) GetURL() url.URL {
	return *ctx.r.URL
}

func (ctx *gmuxContext) GetParam(name string) string {
	return mux.Vars(ctx.r)[name]
}

func (ctx *gmuxContext) GetQuery(name string) string {
	return queryparam.Get(ctx.r.URL.RawQuery, name)
}

func (ctx *gmuxContext) GetHeader(name string) string {
	return ctx.r.Header.Get(name)
}

func (ctx *gmuxContext) EachHeader(cb func(name, value string)) {
	for name, values := range ctx.r.Header {
		for _, value := range values {
			cb(name, value)
		}
	}
}

func (ctx *gmuxContext) GetBodyReader() io.Reader {
	return ctx.r.Body
}

func (ctx *gmuxContext) SetReadDeadline(deadline time.Time) error {
	return huma.SetReadDeadline(ctx.w, deadline)
}

func (ctx *gmuxContext) WriteStatus(code int) {
	ctx.w.WriteHeader(code)
}

func (ctx *gmuxContext) AppendHeader(name string, value string) {
	ctx.w.Header().Add(name, value)
}

func (ctx *gmuxContext) WriteHeader(name string, value string) {
	ctx.w.Header().Set(name, value)
}

func (ctx *gmuxContext) BodyWriter() io.Writer {
	return ctx.w
}

type gMux struct {
	router *mux.Router
}

func (a *gMux) Handle(op *huma.Operation, handler func(huma.Context)) {
	m := op.Method
	a.router.HandleFunc(op.Path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == m {
			handler(&gmuxContext{op: op, r: r, w: w})
		}
	})
}

func (a *gMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func New(r *mux.Router, config huma.Config) huma.API {
	return huma.NewAPI(config, &gMux{router: r})
}
