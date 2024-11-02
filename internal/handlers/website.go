package handlers

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index.html.tmpl", nil, nil)
}

func Manifest(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "manifest.html.tmpl", nil, nil)
}

func Foundation(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "foundation.html.tmpl", nil, nil)
}
