package handlers

import (
	"net/http"
)

func Settings(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "settings.html.tmpl", nil, nil)
}
