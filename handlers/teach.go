package handlers

import (
	"app/db"
	"app/models"
	"app/services"
	"fmt"
	"net/http"
)

func renderTeachTemplate(w http.ResponseWriter, r *http.Request, activeProposalID int, openUnitTable string, openUnitID int) {
	graph := services.GetProposedGraph(activeProposalID)
	positionedGraph := services.CalculatePositions(graph)

	activeProposal := models.Proposal{}

	if activeProposalID == 0 {
		activeProposal = models.Proposal{
			ID:          0,
			Title:       "No active proposal",
			Description: "There are no active proposals",
		}
	} else {
		activeProposal = db.GetProposal(activeProposalID)
	}

	openUnit := models.Unit{}

	if openUnitID == 0 {
		openUnit = models.Unit{
			ID:      0,
			Name:    "No open unit",
			Content: "There are no open units",
		}
	} else {
		if openUnitTable == "units" {
			openUnit = db.GetUnit(openUnitID)
		}
		if openUnitTable == "unit_creations" {
			unitCreation, err := db.GetUnitCreation(openUnitID)
			if err != nil {
				fmt.Println(err)
			}
			openUnit = models.Unit{
				ID:      unitCreation.ID,
				Name:    unitCreation.Name,
				Content: "",
				Type:    "ProposedCreation",
			}
		}
	}

	data := struct {
		PositionedGraph models.PositionedGraph
		Proposals       []models.Proposal
		ActiveProposal  models.Proposal
		OpenUnit        models.Unit
	}{
		PositionedGraph: positionedGraph,
		Proposals:       db.GetUnsubmittedProposals(),
		ActiveProposal:  activeProposal,
		OpenUnit:        openUnit,
	}

	RenderTemplate(w, r, "teach.html", data, nil)
}

func Teach(w http.ResponseWriter, r *http.Request) {
	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func CreateProposal(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	description := r.Form.Get("description")

	p := models.Proposal{
		Title:       title,
		Description: description,
	}

	id := db.CreateProposal(p)

	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, id, openUnitTable, openUnitID)
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

	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func DeleteProposal(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d", &id)

	db.DeleteProposal(id)

	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func SubmitProposal(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/submit", &id)

	db.SubmitProposal(id)

	http.Redirect(w, r, "/teach", http.StatusFound)
}

func CreateUnitCreation(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_creation", &id)

	r.ParseForm()
	name := r.Form.Get("name")

	_, err := db.CreateUnitCreation(id, name)
	if err != nil {
		fmt.Println(err)
	}

	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func DeleteUnitCreation(w http.ResponseWriter, r *http.Request) {
	proposal_id := 0
	change_id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_creation/%d", &proposal_id, &change_id)

	db.DeleteUnitCreation(change_id)

	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func CreateUnitDeletion(w http.ResponseWriter, r *http.Request) {
	proposal_id := 0
	unit_id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_deletion/%d", &proposal_id, &unit_id)

	db.CreateUnitDeletion(proposal_id, unit_id)

	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func DeleteUnitDeletion(w http.ResponseWriter, r *http.Request) {
	proposal_id := 0
	change_id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_deletion/%d", &proposal_id, &change_id)

	fmt.Println("Deleting unit deletion", change_id)
	err := db.DeleteUnitDeletion(change_id)
	if err != nil {
		fmt.Println(err)
	}

	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func CreateUnitRename(w http.ResponseWriter, r *http.Request) {
	proposal_id := 0
	unit_id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_rename/%d", &proposal_id, &unit_id)

	r.ParseForm()
	name := r.Form.Get("name")

	db.CreateUnitRename(proposal_id, unit_id, name)

	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func DeleteUnitRename(w http.ResponseWriter, r *http.Request) {
	proposal_id := 0
	change_id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_rename/%d", &proposal_id, &change_id)

	db.DeleteUnitRename(change_id)

	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func Polls(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Polls []interface{}
	}{
		Polls: db.GetUnResolvedPolls(),
	}
	RenderTemplate(w, r, "polls.html", data, nil)
}

func Poll(w http.ResponseWriter, r *http.Request) {

	pollID := 0
	fmt.Sscanf(r.URL.Path, "/teach/poll/%d", &pollID)

	poll := db.GetPoll(pollID)

	data := struct {
		Poll models.SingleProposalPoll
	}{
		Poll: poll,
	}

	RenderTemplate(w, r, "poll.html", data, nil)
}

func VoteYes(w http.ResponseWriter, r *http.Request) {
	pollID := 0
	fmt.Sscanf(r.URL.Path, "/teach/poll/%d/yes", &pollID)

	db.VoteYes(pollID)

	redirectTo := fmt.Sprintf("/teach/poll/%d", pollID)
	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func VoteNo(w http.ResponseWriter, r *http.Request) {
	pollID := 0
	fmt.Sscanf(r.URL.Path, "/teach/poll/%d/no", &pollID)

	db.VoteNo(pollID)

	redirectTo := fmt.Sprintf("/teach/poll/%d", pollID)
	http.Redirect(w, r, redirectTo, http.StatusFound)
}
