package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"app/db"
)

func main() {
	db.InitializeDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("web/templates/base.html", "web/templates/index.html")
		_ = tmpl.ExecuteTemplate(w, "base.html", struct {
			Title string
		}{
			Title: "Página de Inicio",
		})
	})

	http.HandleFunc("/manifest", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("web/templates/base.html", "web/templates/manifest.html")
		_ = tmpl.ExecuteTemplate(w, "base.html", struct {
			Title string
		}{
			Title: "Manifest",
		})
	})

	http.HandleFunc("/foundation", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("web/templates/base.html", "web/templates/foundation.html")
		_ = tmpl.ExecuteTemplate(w, "base.html", struct {
			Title string
		}{
			Title: "Foundation",
		})
	})

	http.HandleFunc("/units", unitsHandler)
	http.HandleFunc("/units/{id}", unitHandler)
	http.HandleFunc("/dependencies", dependenciesHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	fmt.Println("El servidor está escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func unitsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Parsea la plantilla base y la específica
		tmpl, err := template.ParseFiles("web/templates/base.html", "web/templates/units.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Units by levels based on their dependencies
		// Level 0 means the unit has no dependencies
		// Level 1 means the unit depends on a unit with level 0
		// Level 2 means the unit depends on a unit with level 1, and so on
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
						if checkingU.ID == d.ID && otherU.ID == d.DependsOnID {
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

		// Ejecuta la plantilla, automáticamente "extiende" la base
		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
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

func unitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Parsea la plantilla base y la específica
		tmpl, err := template.ParseFiles("web/templates/base.html", "web/templates/unit.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Obtiene el ID de la URL
		id := 0
		fmt.Sscanf(r.URL.Path, "/units/%d", &id)

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

		// Ejecuta la plantilla, automáticamente "extiende" la base
		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "DELETE" {
		id := 0
		fmt.Sscanf(r.URL.Path, "/units/%d", &id)

		db.DeleteUnit(id)

		w.Header().Set("HX-Redirect", "/units")
		w.WriteHeader(http.StatusOK)
	}
}

func dependenciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		r.ParseForm()
		unit_id, _ := strconv.Atoi(r.Form.Get("unit_id"))
		depends_on_id, _ := strconv.Atoi(r.Form.Get("depends_on_id"))

		db.CreateDependency(unit_id, depends_on_id)
		http.Redirect(w, r, "/units/"+strconv.Itoa(unit_id), http.StatusSeeOther)

	} else if r.Method == "DELETE" {
		unit_id, _ := strconv.Atoi(r.URL.Query().Get("unit_id"))
		depends_on_id, _ := strconv.Atoi(r.URL.Query().Get("depends_on_id"))
		db.DeleteDependency(unit_id, depends_on_id)
		http.Redirect(w, r, "/units/"+strconv.Itoa(unit_id), http.StatusSeeOther)
	}
}
