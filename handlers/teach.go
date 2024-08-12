package handlers

import (
	"app/db"
	"app/models"
	"app/services"
	"fmt"
	"net/http"
	"strconv"
)

func renderTeachTemplate(w http.ResponseWriter, r *http.Request, activeProposalID int, openUnitTable string, openUnitID int) {

	activeProposal := models.Proposal{}

	if activeProposalID == 0 {
		activeProposal = models.Proposal{
			ID:          0,
			Title:       "No active proposal",
			Description: "There are no active proposals",
		}
	} else {
		activeProposal = db.GetProposal(activeProposalID)
		if activeProposal.Submitted {
			activeProposal = models.Proposal{
				ID:          0,
				Title:       "No active proposal",
				Description: "There are no active proposals",
			}
			activeProposalID = 0
			// En realidad se debería hacer un set active proposal a 0, pero así es más conveniente por ahora
		}
	}

	graph := services.GetProposedGraph(activeProposalID)
	positionedGraph := services.CalculatePositions(graph)

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

func UpdateUnitCreation(w http.ResponseWriter, r *http.Request) {
	proposal_id := 0
	unit_id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_creation/%d", &proposal_id, &unit_id)

	r.ParseForm()
	name := r.Form.Get("name")

	db.UpdateUnitCreation(unit_id, name)

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

func ToggleDependency(w http.ResponseWriter, r *http.Request) {
	// Get the proposal ID from the URL
	proposalID := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/toggle_dependency", &proposalID)

	// Get the unit_table and unit_id and depends_on_table and depends_on_id from the URI
	unitIsProposed, err := strconv.ParseBool(r.URL.Query().Get("unit_is_proposed"))
	if err != nil {
		http.Error(w, "Invalid unit_is_proposed", http.StatusBadRequest)
		fmt.Println("Invalid unit_is_proposed")
		return
	}
	unitID, err := strconv.Atoi(r.URL.Query().Get("unit_id"))
	if err != nil {
		http.Error(w, "Invalid unit ID", http.StatusBadRequest)
		fmt.Println("Invalid unit ID")
		return
	}

	dependsOnIsProposed, err := strconv.ParseBool(r.URL.Query().Get("depends_on_is_proposed"))
	if err != nil {
		http.Error(w, "Invalid depends_on_is_proposed", http.StatusBadRequest)
		fmt.Println("Invalid depends_on_is_proposed")
		return
	}
	dependsOnId, err := strconv.Atoi(r.URL.Query().Get("depends_on_id"))
	if err != nil {
		http.Error(w, "Invalid depends_on ID", http.StatusBadRequest)
		fmt.Println("Invalid depends_on ID")
		return
	}

	// Check if there exists any dependency change in the proposal
	proposalChanges := db.GetProposalChanges(proposalID)
	for _, change := range proposalChanges {
		// Change is a generic interface
		// Check what type it is
		switch change := change.(type) {
		case models.DependencyCreation:
			if change.UnitIsProposed == unitIsProposed && change.UnitID == unitID &&
				change.DependsOnIsProposed == dependsOnIsProposed && change.DependsOnID == dependsOnId {
				err := db.DeleteDependencyCreation(change.ID)
				if err != nil {
					fmt.Println(err)
				}
			}
		case models.DependencyDeletion:
			changedDependency := db.GetDependency(change.DependencyID)
			if changedDependency.UnitID == unitID && changedDependency.DependsOnID == dependsOnId {
				err := db.DeleteDependencyDeletion(change.ID)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	// If there isn't, then it has to be created

	dependencyID := db.FindDependency(unitID, dependsOnId)
	fmt.Println("Found dependency ID:", dependencyID)
	if dependencyID != 0 {
		err := db.CreateDependencyDeletion(proposalID, dependencyID)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err := db.CreateDependencyCreation(proposalID, unitIsProposed, unitID, dependsOnIsProposed, dependsOnId)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Toggling dependency between", unitID, "and", dependsOnId)
}

func DeleteDependencyCreation(w http.ResponseWriter, r *http.Request) {
	proposal_id := 0
	change_id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/dependency_creation/%d", &proposal_id, &change_id)

	db.DeleteDependencyCreation(change_id)

	activeProposalID := GetActiveProposalID(r)
	openUnitTable, openUnitID := GetOpenUnit(r)
	renderTeachTemplate(w, r, activeProposalID, openUnitTable, openUnitID)
}

func DeleteDependencyDeletion(w http.ResponseWriter, r *http.Request) {
	proposal_id := 0
	change_id := 0
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/dependency_deletion/%d", &proposal_id, &change_id)

	db.DeleteDependencyDeletion(change_id)

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
