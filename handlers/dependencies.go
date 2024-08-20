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
		unit_id, _ := strconv.ParseInt(r.Form.Get("unit_id"), 10, 64)
		depends_on_id, _ := strconv.ParseInt(r.Form.Get("depends_on_id"), 10, 64)

		db.CreateDependency(unit_id, depends_on_id)
		http.Redirect(w, r, "/unit/"+strconv.FormatInt(unit_id, 10), http.StatusSeeOther)

	case "DELETE":
		unitID, _ := strconv.ParseInt(r.URL.Query().Get("unit_id"), 10, 64)
		dependsOnID, _ := strconv.ParseInt(r.URL.Query().Get("depends_on_id"), 10, 64)
		dependencyID := db.FindDependency(unitID, dependsOnID)
		db.DeleteDependency(dependencyID)
		http.Redirect(w, r, "/unit/"+strconv.FormatInt(unitID, 10), http.StatusSeeOther)
	}
}
