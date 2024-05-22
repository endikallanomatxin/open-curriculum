package handlers

import (
	"app/db"
	"fmt"
	"net/http"
	"text/template"
)

func Learn(w http.ResponseWriter, r *http.Request) {
	// Unitl more logic is implemented, let's just get all
	units := db.GetUnits()
	dependencies := db.GetAllDependencies()

	graphedUnits := MakeGraph(units, dependencies)

	data := struct {
		Dependencies []db.Dependency
		GraphedUnits []GraphedUnit
	}{
		Dependencies: dependencies,
		GraphedUnits: graphedUnits,
	}

	RenderTemplate(w, r, "learn.html", data)
}

func GetUnitDetails(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/unit/%d/details", &id)

	data := struct {
		Unit db.Unit
	}{
		Unit: db.GetUnit(id),
	}

	t, err := template.ParseFiles("web/templates/base.html", "web/templates/learn.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t.ExecuteTemplate(w, "unit_details", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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

	// Now, calculate the horizontal position of each unit
	// For now, let's just put them in a line
	graphedUnits := []GraphedUnit{}
	for _, units := range unitsByLevel {
		for i, u := range units {
			graphedUnits = append(graphedUnits, GraphedUnit{
				Unit:               u,
				HorizontalPosition: float64(i),
			})
		}
	}
	return graphedUnits
}
