package components

import (
	"net/http"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/server/frontend"
	"github.com/enolgor/golang-webservice-template/server/session"
)

type getNavData struct {
	LoggedIn bool
}

func Nav(app *application.App) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := &getNavData{}
		session, err := session.GetAuthSession(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.LoggedIn = session.IsLoggedIn()
		frontend.Components.Serve(w, "nav", data)
	}
}
