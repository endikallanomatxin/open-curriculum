package handlers

import (
	"app/internal/db"
	"app/internal/models"
	"fmt"
	"net/http"
)

func renderLearnTemplate(w http.ResponseWriter, r *http.Request) {

	_, openUnitID := GetOpenUnit(r)

	var graph models.Graph
	if openUnitID == 0 {
		graph = db.GetApprovedGraph()
	} else {
		graph = db.GetApprovedLocalGraph(openUnitID, 2, 2)
	}
	graph.SortAndPosition()

	var openUnit models.Unit
	var err error

	if openUnitID == 0 {
		openUnit = models.Unit{
			ID:      0,
			Name:    "No open unit",
			Content: "There are no open units",
		}
	} else {
		openUnit, err = db.GetUnit(openUnitID)
		if err != nil {
			http.Error(w, "Unit not found", http.StatusNotFound)
			return
		}
	}

	data := struct {
		Graph    models.Graph
		OpenUnit models.Unit
	}{
		Graph:    graph,
		OpenUnit: openUnit,
	}

	RenderTemplate(w, r, "learn.html.tmpl", data, nil)
}

func Learn(w http.ResponseWriter, r *http.Request) {
	renderLearnTemplate(w, r)
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
		Unit models.Unit
	}{
		Unit: unit,
	}

	RenderTemplate(w, r, "learn.html.tmpl", data, "unit_details")
}
