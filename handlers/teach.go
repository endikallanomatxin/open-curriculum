package handlers

import (
	"app/db"
	"fmt"
	"net/http"
)

func Teach(w http.ResponseWriter, r *http.Request) {

	// Unitl more logic is implemented, let's just get all
	units := db.GetUnits()
	dependencies := db.GetAllDependencies()

	graphedUnits := MakeGraph(units, dependencies)

	data := struct {
		Dependencies   []db.Dependency
		GraphedUnits   []GraphedUnit
		Proposals      []db.Proposal
		ActiveProposal db.Proposal
	}{
		Dependencies:   dependencies,
		GraphedUnits:   graphedUnits,
		Proposals:      db.GetProposals(),
		ActiveProposal: db.GetActiveProposal(),
	}

	RenderTemplate(w, r, "teach.html", data, nil)
}

func CreateProposal(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	description := r.Form.Get("description")

	p := db.Proposal{
		Title:       title,
		Description: description,
	}

	db.CreateProposal(p)

	// Render only the block proposals from the teach.html template
	data := struct {
		Proposals      []db.Proposal
		ActiveProposal db.Proposal
	}{
		Proposals:      db.GetProposals(),
		ActiveProposal: db.GetActiveProposal(),
	}

	RenderTemplate(w, r, "teach.html", data, "main")
}

func DeleteProposal(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d", &id)

	db.DeleteProposal(id)

	data := struct {
		Proposals      []db.Proposal
		ActiveProposal db.Proposal
	}{
		Proposals:      db.GetProposals(),
		ActiveProposal: db.GetActiveProposal(),
	}

	RenderTemplate(w, r, "teach.html", data, "main")
}

func AddUnitCreation(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/add_change/unit_creation", &id)

	r.ParseForm()
	name := r.Form.Get("name")

	u := db.Unit{
		Name:    name,
		Content: "",
	}

	db.CreateUnit(u)
}

func Proposal(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d", &id)

	data := struct {
	}{}

	RenderTemplate(w, r, "proposal.html", data, nil)
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
