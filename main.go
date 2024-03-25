package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "¡Hola, localhost!")
    })

    fmt.Println("El servidor está escuchando en http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
