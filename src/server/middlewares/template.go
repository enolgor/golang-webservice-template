package middlewares

import (
	"net/http"
	"strings"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/server/frontend"
)

func StaticTemplate(app *application.App) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			file := strings.TrimPrefix(req.URL.Path, "/")
			if frontend.Pages.Has(file) {
				frontend.Pages.Serve(w, file, nil)
				return
			}
			next.ServeHTTP(w, req)
		}
	}
}
