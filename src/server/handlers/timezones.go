package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/utils/timezone"
)

const resolvedTimezone string = `
window.resolveTimezone = () => { let tz = Intl.DateTimeFormat().resolvedOptions().timeZone; return window.timezones.includes(tz) ? tz : null; };
`

func TimeZonesJS(app *application.App) http.HandlerFunc {
	jsontz, err := json.MarshalIndent(timezone.TimeZones, "", "  ")
	if err != nil {
		panic(err)
	}
	js := fmt.Sprintf("window.timezones = %s;\n%s", jsontz, resolvedTimezone)
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, js)
	}
}
