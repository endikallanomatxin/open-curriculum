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
	units := db.GetUnits()
	dependencies := db.GetAllDependencies()

	positionedUnits := services.PositionUnits(units, dependencies)

	data := struct {
		Dependencies    []models.Dependency
		PositionedUnits []models.PositionedUnit
	}{
		Dependencies:    dependencies,
		PositionedUnits: positionedUnits,
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
