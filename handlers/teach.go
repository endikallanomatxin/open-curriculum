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

	data := struct {
		PositionedGraph models.PositionedGraph
		Proposals       []models.Proposal
		ActiveProposal  models.Proposal
	}{
		PositionedGraph: positionedGraph,
		Proposals:       db.GetUnsubmittedProposals(),
		ActiveProposal:  active_proposal,
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
		Proposals:      db.GetUnsubmittedProposals(),
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
		Proposals:      db.GetUnsubmittedProposals(),
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
		Proposals:      db.GetUnsubmittedProposals(),
		ActiveProposal: GetActiveProposal(r),
	}

	RenderTemplate(w, r, "teach.html", data, "main")
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

	active_proposal := GetActiveProposal(r)
	graph := services.GetProposedGraph(active_proposal.ID)
	positionedGraph := services.CalculatePositions(graph)

	data := struct {
		PositionedGraph models.PositionedGraph
		Proposals       []models.Proposal
		ActiveProposal  models.Proposal
	}{
		PositionedGraph: positionedGraph,
		Proposals:       db.GetUnsubmittedProposals(),
		ActiveProposal:  active_proposal,
	}

	// Just send ok
	RenderTemplate(w, r, "teach.html", data, nil)
}

func DeleteUnitCreation(w http.ResponseWriter, r *http.Request) {
	proposal_id := 0
	change_id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_creation/%d", &proposal_id, &change_id)

	db.DeleteUnitCreation(change_id)

	active_proposal := GetActiveProposal(r)
	graph := services.GetProposedGraph(active_proposal.ID)
	positionedGraph := services.CalculatePositions(graph)

	data := struct {
		PositionedGraph models.PositionedGraph
		Proposals       []models.Proposal
		ActiveProposal  models.Proposal
	}{
		PositionedGraph: positionedGraph,
		Proposals:       db.GetUnsubmittedProposals(),
		ActiveProposal:  active_proposal,
	}

	RenderTemplate(w, r, "teach.html", data, nil)
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
