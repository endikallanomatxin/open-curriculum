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
    http.HandleFunc("/units", unitsHandler)
    http.HandleFunc("/units/{id}", unitHandler)
    http.HandleFunc("/units/{id}/add_dependency", addDependencyHandler)
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "web/index.html")
    })

    fmt.Println("El servidor está escuchando en http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func unitsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Parsea la plantilla base y la específica
        tmpl, err := template.ParseFiles("web/base.html", "web/units.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Datos que se pasarán a la plantilla
        data := struct {
            Title string
            Units []db.Unit
        }{
            Title: "Página de Unidades",
            Units: db.GetUnits(),
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
        tmpl, err := template.ParseFiles("web/base.html", "web/unit.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Obtiene el ID de la URL
        id := 0
        fmt.Sscanf(r.URL.Path, "/units/%d", &id)

        // Datos que se pasarán a la plantilla
        data := struct {
            Title string
            Unit  db.Unit
            Dependencies []db.Unit
            Units []db.Unit
        }{
            Title: "Detalle de la Unidad",
            Unit:  db.GetUnit(id),
            Dependencies: db.GetDependencies(id),
            Units: db.GetUnits(),
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

        http.Redirect(w, r, "/units", http.StatusSeeOther)
    }
}

func addDependencyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "PUT" {
        id := 0
        fmt.Sscanf(r.URL.Path, "/units/%d/add_dependency", &id)

        r.ParseForm()
        depends_on_id, _ := strconv.Atoi(r.Form.Get("depends_on_id"))

        db.CreateDependency(id, depends_on_id)

        http.Redirect(w, r, fmt.Sprintf("/units/%d", id), http.StatusSeeOther)
    }
}
