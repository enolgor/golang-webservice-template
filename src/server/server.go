package server

import (
	"net/http"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/server/frontend"
	"github.com/enolgor/golang-webservice-template/server/muxc"
	"github.com/enolgor/golang-webservice-template/server/session"
)

func NewMux(app *application.App) (*http.ServeMux, error) {
	session.Initialize()
	if err := frontend.Init(app); err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	muxc.ConfigureMux(mux, app)
	return mux, nil
}
