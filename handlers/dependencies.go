package handlers

import (
	"app/db"
	"net/http"
	"strconv"
)

func Dependencies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		r.ParseForm()
		unit_id, _ := strconv.Atoi(r.Form.Get("unit_id"))
		depends_on_id, _ := strconv.Atoi(r.Form.Get("depends_on_id"))

		db.CreateDependency(unit_id, depends_on_id)
		http.Redirect(w, r, "/unit/"+strconv.Itoa(unit_id), http.StatusSeeOther)

	case "DELETE":
		unitID, _ := strconv.Atoi(r.URL.Query().Get("unit_id"))
		dependsOnID, _ := strconv.Atoi(r.URL.Query().Get("depends_on_id"))
		dependencyID := db.FindDependency(unitID, dependsOnID)
		db.DeleteDependency(dependencyID)
		http.Redirect(w, r, "/unit/"+strconv.Itoa(unitID), http.StatusSeeOther)
	}
}
