package handlers

import (
	"app/db"
	"fmt"
	"net/http"
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

	RenderTemplate(w, r, "learn.html", data, nil)
}

func GetUnitDetails(w http.ResponseWriter, r *http.Request) {
	id := 0
	fmt.Sscanf(r.URL.Path, "/unit/%d/details", &id)

	data := struct {
		Unit db.Unit
	}{
		Unit: db.GetUnit(id),
	}

	RenderTemplate(w, r, "learn.html", data, "unit_details")
}
