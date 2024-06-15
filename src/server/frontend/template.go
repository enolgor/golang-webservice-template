package frontend

import (
	"embed"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"text/template"
)

//go:embed templates
var templates embed.FS

var htmlMime string = mime.TypeByExtension(".html")

type HtmlTemplateStore map[string]*template.Template

func newHtmlTemplateStore(log *slog.Logger, dir, baseTemplate string) (*HtmlTemplateStore, error) {
	tmplCache := make(map[string]*template.Template)
	err := fs.WalkDir(Templates, dir, func(fullPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name, err := filepath.Rel(dir, fullPath[:len(fullPath)-len(filepath.Ext(fullPath))])
		if err != nil {
			return err
		}
		var tmpl *template.Template
		if baseTemplate == "" {
			tmpl, err = template.ParseFS(Templates, fullPath)
		} else {
			tmpl, err = template.ParseFS(Templates, baseTemplate, fullPath)
		}
		log.Debug("template parsed", "template", fullPath, "base-template", baseTemplate)
		tmplCache[name] = tmpl
		return err
	})
	if err != nil {
		return nil, err
	}
	store := HtmlTemplateStore(tmplCache)
	return &store, nil
}

type Template template.Template

func (ts *HtmlTemplateStore) Has(filepath string) (ok bool) {
	_, ok = (*ts)[filepath]
	if !ok {
		_, ok = (*ts)[path.Join(filepath, "index")]
	}
	return
}

func (ts *HtmlTemplateStore) Serve(w http.ResponseWriter, filepath string, data any) {
	tmpl := (*ts)[filepath]
	if tmpl == nil {
		tmpl = (*ts)[path.Join(filepath, "index")]
	}
	if tmpl == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", htmlMime)
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
