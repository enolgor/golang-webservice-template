package interceptors

import (
	"io"
	"net/http"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/server/ctxt"
	"github.com/enolgor/golang-webservice-template/server/frontend"
	"github.com/enolgor/golang-webservice-template/server/middlewares"
)

var overrideErrorMessages map[int]string = map[int]string{
	http.StatusUnauthorized: "Sorry, you are not authorized to view this page",
}

func ErrorInterceptor(app *application.App) middlewares.Interceptor {
	return &errorInterceptor{app}
}

type errorTemplateData struct {
	StatusCode int
	StatusText string
	Content    string
}

type errorInterceptor struct {
	app *application.App
}

func (ei *errorInterceptor) ShouldIntercept(status int) bool {
	return status >= 400
}

func (ei *errorInterceptor) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	status := ctxt.GetInterceptedResponseStatus(req.Context())
	text := http.StatusText(status)
	if text == "" {
		text = "Unknown"
	}
	contentText := ""
	if txt, ok := overrideErrorMessages[status]; ok {
		contentText = txt
	} else {
		contentReader := ctxt.GetInterceptedResponseContent(req.Context())
		if contentReader != nil {
			content, err := io.ReadAll(contentReader)
			if err == nil {
				contentText = string(content)
			}
		}
	}

	w.WriteHeader(status)
	frontend.Pages.Serve(w, "error", errorTemplateData{
		StatusCode: status,
		StatusText: text,
		Content:    contentText,
	})
}
