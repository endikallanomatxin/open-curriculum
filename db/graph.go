package db

import (
	"app/logic"
)

func GetApprovedGraph() logic.Graph {
	units := GetUnits()
	dependencies := GetAllDependencies()

	// Mark units as existing
	for i := range units {
		units[i].Type = "Existing"
	}

	// Mark dependencies as existing
	for i := range dependencies {
		dependencies[i].Type = "Existing"
	}

	for i := range units {
		units[i].Relevance = 1
	}

	return logic.Graph{
		Units:        units,
		Dependencies: dependencies,
	}
}

func GetApprovedLocalGraph(
	unitID int64,
	antecessorDepth int,
	successorDepth int) logic.Graph {

	unit, err := GetUnit(unitID)
	if err != nil {
		return logic.Graph{}
	}

	units := []logic.Unit{unit}

	// Get all antecessors
	previousDepthAntecessors := []logic.Unit{unit}
	var thisDepthAntecessors []logic.Unit
	for i := 0; i < antecessorDepth; i++ {
		thisDepthAntecessors = []logic.Unit{}
		for _, u := range previousDepthAntecessors {
			antecessors := GetUnitDirectAntecessors(u.ID)
			if len(antecessors) > 0 {
				thisDepthAntecessors = append(thisDepthAntecessors, antecessors...)
				// Add only new units
				for _, antecessor := range antecessors {
					alreadyExists := false
					for _, existingUnit := range units {
						if existingUnit.ID == antecessor.ID {
							alreadyExists = true
							break
						}
					}
					if !alreadyExists {
						units = append(units, antecessor)
					}
				}
			}
		}
		previousDepthAntecessors = thisDepthAntecessors
	}

	// Get all successors
	previousDepthSuccessors := []logic.Unit{unit}
	var thisDepthSuccessors []logic.Unit
	for i := 0; i < successorDepth; i++ {
		thisDepthSuccessors = []logic.Unit{}
		for _, u := range previousDepthSuccessors {
			successors := GetUnitDirectSuccessors(u.ID)
			if len(successors) > 0 {
				thisDepthSuccessors = append(thisDepthSuccessors, successors...)
				// Add only new units
				for _, successor := range successors {
					alreadyExists := false
					for _, existingUnit := range units {
						if existingUnit.ID == successor.ID {
							alreadyExists = true
							break
						}
					}
					if !alreadyExists {
						units = append(units, successor)
					}
				}
			}
		}
		previousDepthSuccessors = thisDepthSuccessors
	}

	// Find all dependencies between the units
	dependencies := []logic.Dependency{}
	for _, oneUnit := range units {
		for _, anotherUnit := range units {
			if oneUnit.ID == anotherUnit.ID {
				continue
			}
			dependencyID := FindDependency(oneUnit.ID, anotherUnit.ID)
			if dependencyID != 0 {
				dependencies = append(dependencies, GetDependency(dependencyID))
			}
		}
	}

	// Mark units as existing
	for i := range units {
		units[i].Type = "Existing"
	}

	// Mark dependencies as existing
	for i := range dependencies {
		dependencies[i].Type = "Existing"
	}

	for i := range units {
		units[i].Relevance = 1
	}

	return logic.Graph{
		Units:        units,
		Dependencies: dependencies,
	}
}

func GetProposedGraph(proposalID int64) logic.Graph {

	graph := GetApprovedGraph()

	if proposalID == 0 {
		return graph
	}

	proposal := GetProposal(proposalID)

	// Apply proposal changes
	for _, change := range proposal.Changes {
		switch change := change.(type) {
		case logic.UnitDeletion:
			for _, unit := range graph.Units {
				if unit.ID == change.UnitID {
					unit.Type = "ProposedDeletion"
					unit.ChangeID = change.ID
				}
			}
		case logic.UnitCreation:
			graph.Units = append(graph.Units, logic.Unit{
				ChangeID:  change.ID,
				Name:      change.Name,
				Type:      "ProposedCreation",
				Relevance: 1,
			})
		case logic.UnitRename:
			for _, unit := range graph.Units {
				if unit.ID == change.UnitID {
					unit.Type = "ProposedRename"
					unit.ChangeID = change.ID
					unit.Name = change.Name
				}
			}
		case logic.ContentModification:
			for i, unit := range graph.Units {
				if (!change.UnitIsProposed && unit.Type == "Existing" && unit.ID == change.UnitID) ||
					(change.UnitIsProposed && unit.Type == "ProposedCreation" && unit.ChangeID == change.UnitID) {
					graph.Units[i].Content = change.Content
				}
			}
		case logic.DependencyCreation:
			graph.Dependencies = append(graph.Dependencies, logic.Dependency{
				ID:                  change.ID,
				Type:                "ProposedCreation",
				UnitIsProposed:      change.UnitIsProposed,
				UnitID:              change.UnitID,
				DependsOnIsProposed: change.DependsOnIsProposed,
				DependsOnID:         change.DependsOnID,
			})
		case logic.DependencyDeletion:
			changedDependency := GetDependency(change.DependencyID)
			for _, dependency := range graph.Dependencies {
				if dependency.UnitID == changedDependency.UnitID &&
					dependency.DependsOnID == changedDependency.DependsOnID {
					dependency.Type = "ProposedDeletion"
					dependency.ChangeID = change.ID
				}
			}
		}
	}

	return graph
}
