package handlers

import (
	"app/db"
	"fmt"
	"net/http"
	"text/template"
)

func Learn(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Units        []db.Unit
		Dependencies []db.Dependency
		UnitsByLevel map[int][]db.Unit
	}{
		Units:        db.GetUnits(),
		Dependencies: db.GetAllDependencies(),
		UnitsByLevel: db.GetUnitsByLevel(),
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
