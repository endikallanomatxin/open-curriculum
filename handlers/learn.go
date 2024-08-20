package handlers

import (
	"app/db"
	"app/logic"
	"fmt"
	"net/http"
)

func Learn(w http.ResponseWriter, r *http.Request) {
	// Unitl more logic is implemented, let's just get all
	activeProposalID := GetActiveProposalID(r)
	graph := db.GetProposedGraph(activeProposalID)
	positionedGraph := graph.Positioned()

	data := struct {
		PositionedGraph logic.PositionedGraph
		Proposals       []logic.Proposal
		ActiveProposal  logic.Proposal
	}{
		PositionedGraph: positionedGraph,
		Proposals:       db.GetUnsubmittedProposals(),
		ActiveProposal:  db.GetProposal(activeProposalID),
	}

	RenderTemplate(w, r, "learn.html", data, nil)
}

func GetUnitDetails(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/unit/%d/details", &id)

	unit, err := db.GetUnit(int64(id))
	if err != nil {
		http.Error(w, "Unit not found", http.StatusNotFound)
		return
	}

	data := struct {
		Unit logic.Unit
	}{
		Unit: unit,
	}

	RenderTemplate(w, r, "learn.html", data, "unit_details")
}
