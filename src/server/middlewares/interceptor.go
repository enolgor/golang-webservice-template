package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/server/ctxt"
)

type Interceptor interface {
	http.Handler
	ShouldIntercept(status int) bool
}

func Intercept(app *application.App, interceptors ...Interceptor) func(http.HandlerFunc) http.HandlerFunc {
	shouldIntercept := func(status int) bool {
		for i := range interceptors {
			if interceptors[i].ShouldIntercept(status) {
				return true
			}
		}
		return false
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			ww := newResponseInterceptor(w, shouldIntercept)
			next.ServeHTTP(ww, req)
			status := ww.Status()
			for i := range interceptors {
				if interceptors[i].ShouldIntercept(status) {
					ww.NoSkip()
					clearHeaders(ww)
					ctx := ctxt.SetInterceptedResponseStatus(req.Context(), ww.Status())
					ctx = ctxt.SetInterceptedResponseContent(ctx, ww.InterceptedContent())
					interceptors[i].ServeHTTP(ww, req.WithContext(ctx))
					break
				}
			}
		}
	}
}

func clearHeaders(ww *responseInterceptor) {
	for k := range ww.Header() {
		ww.Header().Del(k)
	}
}

func newResponseInterceptor(rw http.ResponseWriter, skip func(int) bool) *responseInterceptor {
	if skip == nil {
		skip = func(status int) bool { return false }
	}
	nrw := &responseInterceptor{
		ResponseWriter:     rw,
		skip:               skip,
		noSkip:             false,
		interceptedContent: &bytes.Buffer{},
	}
	return nrw
}

type responseInterceptor struct {
	http.ResponseWriter
	skip               func(int) bool
	status             int
	wroteHeader        bool
	noSkip             bool
	interceptedContent *bytes.Buffer
}

func (rw *responseInterceptor) WriteHeader(s int) {
	if !rw.noSkip && rw.skip(s) {
		rw.status = s
		return
	}
	if rw.wroteHeader {
		return
	}
	rw.ResponseWriter.WriteHeader(s)
	rw.status = s
	rw.wroteHeader = true
}

func (rw *responseInterceptor) Write(b []byte) (int, error) {
	if !rw.noSkip && rw.skip(rw.status) {
		_, err := rw.interceptedContent.Write(b)
		return 0, err
	}
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *responseInterceptor) Status() int {
	return rw.status
}

func (rw *responseInterceptor) NoSkip() {
	rw.noSkip = true
}

func (rw *responseInterceptor) InterceptedContent() io.ReadCloser {
	return io.NopCloser(rw.interceptedContent)
}
