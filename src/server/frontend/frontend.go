package frontend

import (
	"errors"
	"io/fs"

	"github.com/enolgor/golang-webservice-template/application"
)

var Templates fs.FS
var Static fs.FS

var Pages *HtmlTemplateStore
var Components *HtmlTemplateStore

func Init(app *application.App) error {
	var errs, err error
	Templates, err = fs.Sub(templates, "templates")
	errs = errors.Join(errs, err)
	Pages, err = newHtmlTemplateStore(app.Log, "pages", "base.html")
	errs = errors.Join(errs, err)
	Components, err = newHtmlTemplateStore(app.Log, "components", "")
	errs = errors.Join(errs, err)
	errs = errors.Join(errs, InitStatic(app.Log))
	return errs
}
