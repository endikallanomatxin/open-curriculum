package handlers

import (
	"app/db"
	"fmt"
	"net/http"
)

func Teach(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Units        []db.Unit
		Dependencies []db.Dependency
		UnitsByLevel map[int][]db.Unit
	}{
		Units:        db.GetUnits(),
		Dependencies: db.GetAllDependencies(),
	}

	RenderTemplate(w, r, "teach.html", data)
}

func Proposal(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d", &id)

	data := struct {
	}{}

	RenderTemplate(w, r, "proposal.html", data)
}

func NewProposal(w http.ResponseWriter, r *http.Request) {
	// If the request contains form data, parse it

	r.ParseForm()
	title := r.Form.Get("title")
	description := r.Form.Get("description")

	fmt.Println("something")

	p := db.Proposal{
		Title:       title,
		Description: description,
	}

	db.CreateProposal(p)

	http.Redirect(w, r, "/teach", http.StatusSeeOther)
}
