package middlewares

import (
	"net/http"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/models"
	"github.com/enolgor/golang-webservice-template/server/ctxt"
	"github.com/enolgor/golang-webservice-template/server/session"
)

func Authenticated(app *application.App, allowed ...models.Role) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			session, err := session.GetAuthSession(req)
			if err != nil {
				http.Error(w, "Unable to retrieve session", http.StatusInternalServerError)
				return
			}
			profile := session.GetProfile()
			if profile.IsEmpty() {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			user, err := app.GetUserFromProfile(req.Context(), &profile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(allowed) > 0 && !user.Roles.AnyOf(allowed...) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := ctxt.SetUser(req.Context(), user)
			next.ServeHTTP(w, req.WithContext(ctx))
		}
	}
}
