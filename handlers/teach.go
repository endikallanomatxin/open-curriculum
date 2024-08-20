package handlers

import (
	"app/db"
	"app/logic"
	"fmt"
	"net/http"
	"strconv"
)

func renderTeachTemplate(w http.ResponseWriter, r *http.Request) {

	activeProposalID := GetActiveProposalID(r)
	openUnitIsProposed, openUnitID := GetOpenUnit(r)

	activeProposal := logic.Proposal{}

	if activeProposalID == 0 {
		activeProposal = logic.Proposal{
			ID:          0,
			Title:       "No active proposal",
			Description: "There are no active proposals",
		}
	} else {
		activeProposal = db.GetProposal(activeProposalID)
		if activeProposal.Submitted {
			activeProposal = logic.Proposal{
				ID:          0,
				Title:       "No active proposal",
				Description: "There are no active proposals",
			}
			activeProposalID = 0
			// En realidad se debería hacer un set active proposal a 0, pero así es más conveniente por ahora
		}
	}

	graph := db.GetProposedGraph(activeProposalID)
	positionedGraph := graph.Positioned()

	openUnit := logic.Unit{}
	err := error(nil)

	if openUnitID == 0 {
		openUnit = logic.Unit{
			ID:      0,
			Name:    "No open unit",
			Content: "There are no open units",
		}
	} else {
		if !openUnitIsProposed {
			openUnit, err = db.GetUnit(openUnitID)
			if err != nil {
				fmt.Println(err)
				openUnit = logic.Unit{
					ID:      0,
					Name:    "No open unit",
					Content: "There are no open units",
				}
			}
		} else {
			unitCreation, err := db.GetUnitCreation(openUnitID)
			if err != nil {
				fmt.Println(err)
				openUnit = logic.Unit{
					ID:      0,
					Name:    "No open unit",
					Content: "There are no open units",
				}
			}
			openUnit = logic.Unit{
				ID:      unitCreation.ID,
				Name:    unitCreation.Name,
				Content: "",
				Type:    "ProposedCreation",
			}
		}
	}

	data := struct {
		PositionedGraph logic.PositionedGraph
		Proposals       []logic.Proposal
		ActiveProposal  logic.Proposal
		OpenUnit        logic.Unit
	}{
		PositionedGraph: positionedGraph,
		Proposals:       db.GetUnsubmittedProposals(),
		ActiveProposal:  activeProposal,
		OpenUnit:        openUnit,
	}

	RenderTemplate(w, r, "teach.html", data, nil)
}

func Teach(w http.ResponseWriter, r *http.Request) {
	renderTeachTemplate(w, r)
}

func CreateProposal(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	description := r.Form.Get("description")

	p := logic.Proposal{
		Title:       title,
		Description: description,
	}

	_ = db.CreateProposal(p)

	// TODO: Set active proposal to the id

	renderTeachTemplate(w, r)
}

func UpdateProposal(w http.ResponseWriter, r *http.Request) {
	var id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/update", &id)

	// Parse the form data
	r.ParseForm()
	title := r.Form.Get("title")
	description := r.Form.Get("description")

	// Update the proposal in the database
	p := logic.Proposal{
		ID:          id,
		Title:       title,
		Description: description,
	}
	db.UpdateProposal(p)

	renderTeachTemplate(w, r)
}

func DeleteProposal(w http.ResponseWriter, r *http.Request) {
	var id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d", &id)

	db.DeleteProposal(id)

	renderTeachTemplate(w, r)
}

func SubmitProposal(w http.ResponseWriter, r *http.Request) {
	var id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/submit", &id)

	db.SubmitProposal(id)

	http.Redirect(w, r, "/teach", http.StatusFound)
}

func CreateUnitCreation(w http.ResponseWriter, r *http.Request) {
	var id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_creation", &id)

	r.ParseForm()
	name := r.Form.Get("name")

	_, err := db.CreateUnitCreation(id, name)
	if err != nil {
		fmt.Println(err)
	}

	renderTeachTemplate(w, r)
}

func UpdateUnitCreation(w http.ResponseWriter, r *http.Request) {
	var proposal_id int64
	var unit_id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_creation/%d", &proposal_id, &unit_id)

	r.ParseForm()
	name := r.Form.Get("name")

	db.UpdateUnitCreation(unit_id, name)

	renderTeachTemplate(w, r)
}

func DeleteUnitCreation(w http.ResponseWriter, r *http.Request) {
	var proposal_id int64
	var change_id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_creation/%d", &proposal_id, &change_id)

	db.DeleteUnitCreation(change_id)

	renderTeachTemplate(w, r)
}

func CreateUnitDeletion(w http.ResponseWriter, r *http.Request) {
	var proposal_id int64
	var unit_id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_deletion/%d", &proposal_id, &unit_id)

	db.CreateUnitDeletion(proposal_id, unit_id)

	renderTeachTemplate(w, r)
}

func DeleteUnitDeletion(w http.ResponseWriter, r *http.Request) {
	var proposal_id int64
	var change_id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_deletion/%d", &proposal_id, &change_id)

	fmt.Println("Deleting unit deletion", change_id)
	err := db.DeleteUnitDeletion(change_id)
	if err != nil {
		fmt.Println(err)
	}

	renderTeachTemplate(w, r)
}

func CreateUnitRename(w http.ResponseWriter, r *http.Request) {
	var proposal_id int64
	var unit_id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_rename/%d", &proposal_id, &unit_id)

	r.ParseForm()
	name := r.Form.Get("name")

	db.CreateUnitRename(proposal_id, unit_id, name)

	renderTeachTemplate(w, r)
}

func DeleteUnitRename(w http.ResponseWriter, r *http.Request) {
	var proposal_id int64
	var change_id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/unit_rename/%d", &proposal_id, &change_id)

	db.DeleteUnitRename(change_id)

	renderTeachTemplate(w, r)
}

func ToggleDependency(w http.ResponseWriter, r *http.Request) {
	// Get the proposal ID from the URL
	var proposalID int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/toggle_dependency", &proposalID)

	// Get the unit_table and unit_id and depends_on_table and depends_on_id from the URI
	unitIsProposed, err := strconv.ParseBool(r.URL.Query().Get("unit_is_proposed"))
	if err != nil {
		http.Error(w, "Invalid unit_is_proposed", http.StatusBadRequest)
		fmt.Println("Invalid unit_is_proposed")
		return
	}
	unitID, err := strconv.ParseInt(r.URL.Query().Get("unit_id"), 10, 64)
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
	dependsOnId, err := strconv.ParseInt(r.URL.Query().Get("depends_on_id"), 10, 64)
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
		case logic.DependencyCreation:
			if change.UnitIsProposed == unitIsProposed && change.UnitID == unitID &&
				change.DependsOnIsProposed == dependsOnIsProposed && change.DependsOnID == dependsOnId {
				err := db.DeleteDependencyCreation(change.ID)
				if err != nil {
					fmt.Println(err)
				}
			}
		case logic.DependencyDeletion:
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
	var proposal_id int64
	var change_id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/dependency_creation/%d", &proposal_id, &change_id)

	db.DeleteDependencyCreation(change_id)

	renderTeachTemplate(w, r)
}

func DeleteDependencyDeletion(w http.ResponseWriter, r *http.Request) {
	var proposal_id int64
	var change_id int64
	fmt.Sscanf(r.URL.Path, "/teach/proposal/%d/dependency_deletion/%d", &proposal_id, &change_id)

	db.DeleteDependencyDeletion(change_id)

	renderTeachTemplate(w, r)
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

	var pollID int64
	fmt.Sscanf(r.URL.Path, "/teach/poll/%d", &pollID)

	poll := db.GetPoll(pollID)

	data := struct {
		Poll logic.SingleProposalPoll
	}{
		Poll: poll,
	}

	RenderTemplate(w, r, "poll.html", data, nil)
}

func VoteYes(w http.ResponseWriter, r *http.Request) {
	var pollID int64
	fmt.Sscanf(r.URL.Path, "/teach/poll/%d/yes", &pollID)

	db.VoteYes(pollID)

	redirectTo := fmt.Sprintf("/teach/poll/%d", pollID)
	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func VoteNo(w http.ResponseWriter, r *http.Request) {
	var pollID int64
	fmt.Sscanf(r.URL.Path, "/teach/poll/%d/no", &pollID)

	db.VoteNo(pollID)

	redirectTo := fmt.Sprintf("/teach/poll/%d", pollID)
	http.Redirect(w, r, redirectTo, http.StatusFound)
}
