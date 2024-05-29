package handlers

import (
	"app/db"
	"app/models"
	"app/services"
	"fmt"
	"net/http"
)

func Teach(w http.ResponseWriter, r *http.Request) {

	active_proposal := GetActiveProposal(r)
	graph := services.GetProposedGraph(active_proposal.ID)
	positionedGraph := services.CalculatePositions(graph)

	fmt.Println("Graph", graph)

	data := struct {
		PositionedGraph models.PositionedGraph
		Proposals       []models.Proposal
		ActiveProposal  models.Proposal
	}{
		PositionedGraph: positionedGraph,
		Proposals:       db.GetProposals(),
		ActiveProposal:  GetActiveProposal(r),
	}

	RenderTemplate(w, r, "teach.html", data, nil)
}

func CreateProposal(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	description := r.Form.Get("description")

	p := models.Proposal{
		Title:       title,
		Description: description,
	}

	db.CreateProposal(p)

	// Render only the block proposals from the teach.html template
	data := struct {
		Proposals      []models.Proposal
		ActiveProposal models.Proposal
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
	p := models.Proposal{
		ID:          id,
		Title:       title,
		Description: description,
	}
	db.UpdateProposal(p)

	// Render the updated template
	data := struct {
		Proposals      []models.Proposal
		ActiveProposal models.Proposal
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
		Proposals      []models.Proposal
		ActiveProposal models.Proposal
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

	u := models.Unit{
		Name:    name,
		Content: "",
	}

	db.CreateUnit(u)
}
