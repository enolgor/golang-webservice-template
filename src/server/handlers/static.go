package handlers

import (
	"net/http"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/server/frontend"
)

func Static(app *application.App) http.HandlerFunc {
	return http.FileServer(http.FS(frontend.Static)).ServeHTTP
}
