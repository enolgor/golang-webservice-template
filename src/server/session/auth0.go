package session

import (
	"net/http"

	"github.com/enolgor/golang-webservice-template/models"
	"github.com/gorilla/sessions"
)

const (
	authState Key = iota + 1
	authAccessToken
	authProfile
)

const authSessionName string = "auth"

type AuthSession interface {
	Save(req *http.Request, w http.ResponseWriter) error
	SetState(state string)
	GetState() string
	SetAccessToken(accessToken string)
	SetOIDCProfile(profile map[string]interface{})
	GetOIDCProfile() any
	GetProfile() models.Profile
	Clear()
	IsLoggedIn() bool
}

type authSession sessions.Session

func GetAuthSession(req *http.Request) (AuthSession, error) {
	session, err := cookieStore.Get(req, authSessionName)
	as := authSession(*session)
	return &as, err
}

func (as *authSession) Save(req *http.Request, w http.ResponseWriter) error {
	session := sessions.Session(*as)
	return session.Save(req, w)
}

func (as *authSession) SetState(state string) {
	as.Values[authState] = state
}

func (as *authSession) SetAccessToken(accessToken string) {
	as.Values[authAccessToken] = accessToken
}

func (as *authSession) SetOIDCProfile(profile map[string]interface{}) {
	as.Values[authProfile] = profile
}

func (as *authSession) GetOIDCProfile() any {
	return as.Values[authProfile]
}

func (as *authSession) GetProfile() (userProfile models.Profile) {
	oidcProfile := as.Values[authProfile]
	if oidcProfile == nil {
		return
	}
	if profileMap, ok := oidcProfile.(map[string]any); !ok {
		return
	} else {
		userProfile.UserID, _ = profileMap["sub"].(string)
		userProfile.Email, _ = profileMap["email"].(string)
	}
	return
}

func (as *authSession) GetState() string {
	if v, ok := as.Values[authState]; ok {
		if session, ok := v.(string); ok {
			return session
		}
	}
	return ""
}

func (as *authSession) Clear() {
	as.Options.MaxAge = -1
}

func (as *authSession) IsLoggedIn() bool {
	return !as.GetProfile().IsEmpty()
}
