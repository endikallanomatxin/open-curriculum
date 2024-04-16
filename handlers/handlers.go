package handlers

import (
	"net/http"
	"text/template"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("web/templates/base.html", "web/templates/"+tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html", nil)
}

func Manifest(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "manifest.html", nil)
}

func Foundation(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "foundation.html", nil)
}
