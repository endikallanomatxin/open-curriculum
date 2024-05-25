package handlers

import (
	"app/db"
)

type GraphedUnit struct {
	Unit               db.Unit
	HorizontalPosition float64
}

func MakeGraph(units []db.Unit, dependencies []db.Dependency) []GraphedUnit {

	// First, get units by level

	unassignedUnits := units
	unitsByLevel := make(map[int][]db.Unit)

	// Iterate over the units and assign them to a level
	// Units array will be empty when all units are assigned to a level
	for level := 0; len(unassignedUnits) > 0; level++ {
		unitsByLevel[level] = []db.Unit{}

		for _, checkingU := range unassignedUnits {
			dependsOnUnassigned := false

			for _, d := range dependencies {
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
	graphedUnits := []GraphedUnit{}
	for _, units := range unitsByLevel {
		// Distribuir equitativamente las unidades en el nivel
		unitCount := len(units)
		for i, unit := range units {
			horizontalPosition := float64(i) / float64(unitCount-1) // Normalizar entre 0 y 1
			graphedUnits = append(graphedUnits, GraphedUnit{
				Unit:               unit,
				HorizontalPosition: horizontalPosition,
			})
		}
	}

	return graphedUnits
}
