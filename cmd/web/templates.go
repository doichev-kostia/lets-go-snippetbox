package main

import (
	"github.com/justinas/nosurf"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"snippetbox.doichevkostia.dev/internal/models"
	"snippetbox.doichevkostia.dev/ui"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Toast           string
	IsAuthenticated bool
	CSRFToken       string
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Toast:           app.sessionManager.PopString(r.Context(), "toast"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.gohtml",
			"html/partials/*.gohtml",
			page,
		}

		tmpl, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}
