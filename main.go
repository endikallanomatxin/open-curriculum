package main

import (
    "fmt"
    "net/http"

	"servidor/db"
)

func main() {
    db.InitializeDB()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "web/index.html")
    })

    fmt.Println("El servidor está escuchando en http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
