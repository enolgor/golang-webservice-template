package session

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

var cookieStore *sessions.CookieStore

type Key int

func Initialize() {
	cookieStore = sessions.NewCookieStore([]byte("6rWcPjU2U#ss#rK*k7m@XZye34U3pLpJ"))
	gob.Register(Key(0))
	gob.Register(map[Key]any{})
	gob.Register(map[string]interface{}{})
}
