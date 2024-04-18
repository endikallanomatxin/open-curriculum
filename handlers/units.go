package handlers

import (
	"app/db"
	"fmt"
	"net/http"
)

func Units(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		// Units by levels based on their dependencies
		unassignedUnits := db.GetUnits()
		dependencies := db.GetAllDependencies()

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

		// Datos que se pasarán a la plantilla
		data := struct {
			Title        string
			Units        []db.Unit
			Dependencies []db.Dependency
			UnitsByLevel map[int][]db.Unit
		}{
			Title:        "Página de Unidades",
			Units:        db.GetUnits(),
			Dependencies: db.GetAllDependencies(),
			UnitsByLevel: unitsByLevel,
		}

		RenderTemplate(w, "units.html", data)

	case "POST":
		r.ParseForm()
		name := r.Form.Get("name")
		description := r.Form.Get("description")

		u := db.Unit{
			Name:        name,
			Description: description,
		}

		db.CreateUnit(u)

		http.Redirect(w, r, "/units", http.StatusSeeOther)
	}
}

func Unit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		// Obtiene el ID de la URL
		id := 0
		fmt.Sscanf(r.URL.Path, "/unit/%d", &id)

		// Datos que se pasarán a la plantilla
		data := struct {
			Title        string
			Unit         db.Unit
			Dependencies []db.Unit
			Units        []db.Unit
		}{
			Title:        "Detalle de la Unidad",
			Unit:         db.GetUnit(id),
			Dependencies: db.GetUnitDependencies(id),
			Units:        db.GetUnits(),
		}

		RenderTemplate(w, "unit.html", data)

	case "DELETE":
		id := 0
		fmt.Sscanf(r.URL.Path, "/unit/%d", &id)

		db.DeleteUnit(id)

		w.Header().Set("HX-Redirect", "/units")
		w.WriteHeader(http.StatusOK)
	}
}
