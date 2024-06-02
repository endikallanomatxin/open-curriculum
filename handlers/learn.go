package handlers

import (
	"app/db"
	"app/models"
	"app/services"
	"fmt"
	"net/http"
)

func Learn(w http.ResponseWriter, r *http.Request) {
	// Unitl more logic is implemented, let's just get all
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
		ActiveProposal:  GetActiveProposal(r),
	}

	RenderTemplate(w, r, "learn.html", data, nil)
}

func GetUnitDetails(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/unit/%d/details", &id)

	data := struct {
		Unit models.Unit
	}{
		Unit: db.GetUnit(id),
	}

	RenderTemplate(w, r, "learn.html", data, "unit_details")
}
