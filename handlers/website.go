package handlers

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index.html", nil)
}

func Manifest(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "manifest.html", nil)
}

func ManifestES(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "manifest-es.html", nil)
}

func Foundation(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "foundation.html", nil)
}
