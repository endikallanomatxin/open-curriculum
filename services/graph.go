package services

import (
	"app/db"
	"app/models"
)

func CalculatePositions(graph models.Graph) models.PositionedGraph {

	// First, get units by level

	unassignedUnits := graph.Units
	unitsByLevel := make(map[int][]models.Unit)

	// Iterate over the units and assign them to a level
	// Units array will be empty when all units are assigned to a level
	for level := 0; len(unassignedUnits) > 0; level++ {
		unitsByLevel[level] = []models.Unit{}

		for _, checkingU := range unassignedUnits {
			dependsOnUnassigned := false

			for _, d := range graph.Dependencies {
				for _, otherU := range unassignedUnits {
					if checkingU.ID == d.UnitID && otherU.ID == d.DependsOnID {
						dependsOnUnassigned = true
						break
					}
				}
			}

			if !dependsOnUnassigned {
				unitsByLevel[level] = append(unitsByLevel[level], checkingU)
			}
		}

		// Remove assigned units from the unassigned units array
		for _, u := range unitsByLevel[level] {
			for i, unassignedU := range unassignedUnits {
				if u.ID == unassignedU.ID {
					unassignedUnits = append(unassignedUnits[:i], unassignedUnits[i+1:]...)
					break
				}
			}
		}
	}

	// Ahora, calcular la posición horizontal de cada unidad
	positionedUnits := []models.PositionedUnit{}
	for _, units := range unitsByLevel {
		// Distribuir equitativamente las unidades en el nivel
		unitCount := len(units)
		for i, unit := range units {
			horizontalPosition := float64(i) / float64(unitCount-1) // Normalizar entre 0 y 1
			positionedUnits = append(positionedUnits, models.PositionedUnit{
				Unit:               unit,
				HorizontalPosition: horizontalPosition,
			})
		}
	}

	positionedGraph := models.PositionedGraph{
		PositionedUnits: positionedUnits,
		Dependencies:    graph.Dependencies,
	}

	return positionedGraph
}

func GetProposedGraph(proposalID int) models.Graph {
	units := db.GetUnits()
	dependencies := db.GetAllDependencies()

	if proposalID == 0 {
		return models.Graph{
			Units:        units,
			Dependencies: dependencies,
		}
	}

	proposal := db.GetProposal(proposalID)

	// Apply proposal changes
	for _, change := range proposal.Changes {
		switch change := change.(type) {
		case models.UnitCreation:
			units = append(units, models.Unit{
				ID:                change.ID,
				Name:              change.Name,
				IsAProposedChange: true,
			})
		}
	}

	return models.Graph{
		Units:        units,
		Dependencies: dependencies,
	}
}
