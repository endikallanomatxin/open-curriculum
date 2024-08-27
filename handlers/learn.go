package handlers

import (
	"app/db"
	"app/logic"
	"fmt"
	"net/http"
)

func renderLearnTemplate(w http.ResponseWriter, r *http.Request) {

	_, openUnitID := GetOpenUnit(r)

	graph := db.GetApprovedGraph()
	graph.SortAndPosition()

	var openUnit logic.Unit
	var err error

	if openUnitID == 0 {
		openUnit = logic.Unit{
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
		Graph    logic.Graph
		OpenUnit logic.Unit
	}{
		Graph:    graph,
		OpenUnit: openUnit,
	}

	RenderTemplate(w, r, "learn.html", data, nil)
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
		Unit logic.Unit
	}{
		Unit: unit,
	}

	RenderTemplate(w, r, "learn.html", data, "unit_details")
}
