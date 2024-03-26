package main

import (
    "fmt"
    "html/template"
    "net/http"

	"app/db"
)

func main() {
    db.InitializeDB()
    http.HandleFunc("/units", unitsHandler)
    http.HandleFunc("/units/create", createUnitHandler)
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "web/index.html")
    })

    fmt.Println("El servidor está escuchando en http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func unitsHandler(w http.ResponseWriter, r *http.Request) {
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
}

func createUnitHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        unitsHandler(w, r)
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