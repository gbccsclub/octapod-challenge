package web

import (
	"html/template"
	"net/http"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w http.ResponseWriter, name string, data interface{}) {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		panic(err)
	}
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}
