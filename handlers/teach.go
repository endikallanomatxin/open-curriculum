package handlers

import (
	"app/db"
	"app/models"
	"app/services"
	"fmt"
	"net/http"
	"strconv"
)

func GetActiveProposal(r *http.Request) models.Proposal {
	var active_proposal models.Proposal

	// If the request contains an active proposal id, change the active proposal
	if r.URL.Query().Get("active_proposal_id") != "" {
		active_proposa_id, err := strconv.Atoi(r.URL.Query().Get("active_proposal_id"))
		if err != nil {
			fmt.Println("Error converting active_proposal_id to int")
		}
		active_proposal = db.GetProposal(active_proposa_id)
	} else {
		active_proposal = models.Proposal{
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

	positionedUnits := services.PositionUnits(units, dependencies)

	active_proposal := GetActiveProposal(r)
	graph := db.GetProposedGraph(active_proposal.ID)

	data := struct {
		Dependencies    []models.Dependency
		PositionedUnits []models.PositionedUnit
		Proposals       []models.Proposal
		ActiveProposal  models.Proposal
		Graph           models.Graph
	}{
		Dependencies:    dependencies,
		PositionedUnits: positionedUnits,
		Proposals:       db.GetProposals(),
		ActiveProposal:  GetActiveProposal(r),
		Graph:           graph,
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
