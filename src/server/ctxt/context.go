package ctxt

import (
	"context"
	"io"
	"time"

	"github.com/enolgor/golang-webservice-template/models"
)

type ContextKey int

const (
	ctxt_request_id_key ContextKey = iota
	ctxt_request_time_key
	ctxt_cookies_key
	ctxt_profile_key
	ctxt_intercepted_response_status_key
	ctxt_intercepted_response_content_key
)

func SetRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ctxt_request_id_key, requestID)
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(ctxt_request_id_key).(string); ok {
		return reqID
	}
	return ""
}

func SetRequestTime(ctx context.Context, t time.Time) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ctxt_request_time_key, t)
}

func GetRequestTime(ctx context.Context) time.Time {
	if ctx == nil {
		return time.Time{}
	}
	if t, ok := ctx.Value(ctxt_request_time_key).(time.Time); ok {
		return t
	}
	return time.Time{}
}

func SetUser(ctx context.Context, user *models.User) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ctxt_profile_key, user)
}

func GetUser(ctx context.Context) *models.User {
	if ctx == nil {
		return nil
	}
	if profile, ok := ctx.Value(ctxt_profile_key).(*models.User); ok {
		return profile
	}
	return nil
}

func SetInterceptedResponseStatus(ctx context.Context, status int) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ctxt_intercepted_response_status_key, status)
}

func GetInterceptedResponseStatus(ctx context.Context) int {
	if ctx == nil {
		return 0
	}
	if status, ok := ctx.Value(ctxt_intercepted_response_status_key).(int); ok {
		return status
	}
	return 0
}

func SetInterceptedResponseContent(ctx context.Context, content io.ReadCloser) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ctxt_intercepted_response_content_key, content)
}

func GetInterceptedResponseContent(ctx context.Context) io.ReadCloser {
	if ctx == nil {
		return nil
	}
	if content, ok := ctx.Value(ctxt_intercepted_response_content_key).(io.ReadCloser); ok {
		return content
	}
	return nil
}
