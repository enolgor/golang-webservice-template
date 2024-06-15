package middlewares

import (
	"net/http"
	"time"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/server/ctxt"
	"github.com/google/uuid"
)

func Logger(app *application.App) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			requestId := uuid.NewString()
			now := time.Now()
			ctx := ctxt.SetRequestID(req.Context(), requestId)
			ctx = ctxt.SetRequestTime(ctx, now)
			ww := newResponseLogger(w)
			defer func() {
				duration := time.Since(now)
				app.Log.Info("request", "requestId", requestId, "method", req.Method, "path", req.URL.Path, "query", req.URL.RawQuery, "addr", req.RemoteAddr, "status", ww.Status(), "size", ww.Size(), "duration", duration)
			}()
			next.ServeHTTP(ww, req.WithContext(ctx))
		}
	}
}

func newResponseLogger(rw http.ResponseWriter) *responseLogger {
	nrw := &responseLogger{
		ResponseWriter: rw,
	}
	return nrw
}

type responseLogger struct {
	http.ResponseWriter
	pendingStatus int
	status        int
	size          int
}

func (rw *responseLogger) WriteHeader(s int) {
	rw.pendingStatus = s
	if rw.Written() {
		return
	}
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseLogger) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseLogger) Status() int {
	if rw.Written() {
		return rw.status
	}

	return rw.pendingStatus
}

func (rw *responseLogger) Size() int {
	return rw.size
}

func (rw *responseLogger) Written() bool {
	return rw.status != 0
}
