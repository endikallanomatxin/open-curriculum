package handlers

import (
	"app/db"
	"app/models"
	"fmt"
	"net/http"
)

func Units(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		// Datos que se pasarán a la plantilla
		data := struct {
			Title        string
			Units        []models.Unit
			Dependencies []models.Dependency
		}{
			Title:        "Página de Unidades",
			Units:        db.GetUnits(),
			Dependencies: db.GetAllDependencies(),
		}

		RenderTemplate(w, r, "units.html", data, nil)

	case "POST":
		r.ParseForm()
		name := r.Form.Get("name")
		content := r.Form.Get("content")

		u := models.Unit{
			Name:    name,
			Content: content,
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

		unit, err := db.GetUnit(id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unidad no encontrada", http.StatusNotFound)
			return
		}

		// Datos que se pasarán a la plantilla
		data := struct {
			Title        string
			Unit         models.Unit
			Dependencies []models.Unit
			Units        []models.Unit
		}{
			Title:        "Detalle de la Unidad",
			Unit:         unit,
			Dependencies: db.GetUnitDependencies(id),
			Units:        db.GetUnits(),
		}

		RenderTemplate(w, r, "unit.html", data, nil)

	case "DELETE":
		id := 0
		fmt.Sscanf(r.URL.Path, "/unit/%d", &id)

		db.DeleteUnit(id)

		w.Header().Set("HX-Redirect", "/units")
		w.WriteHeader(http.StatusOK)
	}
}
