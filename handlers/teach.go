package handlers

import (
	"app/db"
	"fmt"
	"net/http"
	"strconv"
)

func GetActiveProposal(r *http.Request) db.Proposal {
	var active_proposal db.Proposal

	// If the request contains an active proposal id, change the active proposal
	if r.URL.Query().Get("active_proposal_id") != "" {
		active_proposa_id, err := strconv.Atoi(r.URL.Query().Get("active_proposal_id"))
		if err != nil {
			fmt.Println("Error converting active_proposal_id to int")
		}
		active_proposal = db.GetProposal(active_proposa_id)
	} else {
		active_proposal = db.Proposal{
			ID:          0,
			Title:       "No active proposal",
			Description: "There are no active proposals",
		}
	}

	return active_proposal
}

func Teach(w http.ResponseWriter, r *http.Request) {

	// Until more logic is implemented, let's just get all
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
		ActiveProposal: GetActiveProposal(r),
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
		ActiveProposal: GetActiveProposal(r),
	}

	RenderTemplate(w, r, "teach.html", data, "main")
}

func UpdateProposal(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/update", &id)

	// Parse the form data
	r.ParseForm()
	title := r.Form.Get("title")
	description := r.Form.Get("description")

	// Update the proposal in the database
	p := db.Proposal{
		ID:          id,
		Title:       title,
		Description: description,
	}
	db.UpdateProposal(p)

	// Render the updated template
	data := struct {
		Proposals      []db.Proposal
		ActiveProposal db.Proposal
	}{
		Proposals:      db.GetProposals(),
		ActiveProposal: db.GetProposal(id),
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
		ActiveProposal: GetActiveProposal(r),
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
