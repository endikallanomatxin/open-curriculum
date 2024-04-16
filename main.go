package main

import (
	"fmt"
	"net/http"

	"app/db"
	"app/handlers"
)

func main() {
	db.Init()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/manifest", handlers.Manifest)
	mux.HandleFunc("/foundation", handlers.Foundation)
	mux.HandleFunc("/units", handlers.Units)
	mux.HandleFunc("/unit/{id}", handlers.Unit)
	mux.HandleFunc("/dependencies", handlers.Dependencies)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
