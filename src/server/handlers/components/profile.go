package components

import (
	"encoding/json"
	"net/http"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/server/frontend"
	"github.com/enolgor/golang-webservice-template/server/session"
)

type getProfileData struct {
	LoggedIn    bool
	ProfileJson string
	UserJson    string
}

func Profile(app *application.App) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := &getProfileData{}
		session, err := session.GetAuthSession(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.LoggedIn = session.IsLoggedIn()
		if !data.LoggedIn {
			frontend.Components.Serve(w, "profile", data)
			return
		}
		profile := session.GetProfile()
		profileJson, err := json.MarshalIndent(session.GetOIDCProfile(), "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.ProfileJson = string(profileJson)
		user, err := app.GetUserFromProfile(req.Context(), &profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userJson, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.UserJson = string(userJson)
		frontend.Components.Serve(w, "profile", data)
	}
}
