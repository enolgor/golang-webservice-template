package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/config"
	"github.com/enolgor/golang-webservice-template/server/session"
)

func Login(app *application.App) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		state, err := generateRandomState()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session, err := session.GetAuthSession(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.SetState(state)
		if err := session.Save(req, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, app.Authenticator.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func Auth0Callback(app *application.App) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, err := session.GetAuthSession(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if req.URL.Query().Get("state") != session.GetState() {
			http.Error(w, "Invalid state parameter.", http.StatusBadRequest)
			return
		}

		token, err := app.Authenticator.Exchange(req.Context(), req.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange an authorization code for a token.", http.StatusUnauthorized)
			return
		}

		idToken, err := app.Authenticator.VerifyIDToken(req.Context(), token)
		if err != nil {
			http.Error(w, "Failed to verify ID Token.", http.StatusInternalServerError)
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.SetOIDCProfile(profile)
		session.SetAccessToken(token.AccessToken)
		if err := session.Save(req, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	}
}

func Logout(app *application.App) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, err := session.GetAuthSession(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Clear()
		if err := session.Save(req, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logoutUrl, err := url.Parse("https://" + config.AUTH0_DOMAIN + "/v2/logout")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		scheme := "http"
		if req.TLS != nil {
			scheme = "https"
		}

		returnTo, err := url.Parse(scheme + "://" + req.Host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		parameters := url.Values{}
		parameters.Add("returnTo", returnTo.String())
		parameters.Add("client_id", config.AUTH0_CLIENT_ID)
		logoutUrl.RawQuery = parameters.Encode()

		http.Redirect(w, req, logoutUrl.String(), http.StatusTemporaryRedirect)
	}
}
