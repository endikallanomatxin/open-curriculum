package main

import (
	"net/http"
	"fmt"
	"app/db"
	"app/handlers"
)

func main() {
	db.Init()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/manifest", handlers.Manifest)
	mux.HandleFunc("/manifest-es", handlers.ManifestES)
	mux.HandleFunc("/foundation", handlers.Foundation)
	mux.HandleFunc("/units", handlers.Units)
	mux.HandleFunc("/unit/{id}", handlers.Unit)
	mux.HandleFunc("/dependencies", handlers.Dependencies)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))


	err := http.ListenAndServeTLS(":443",
		"data/certbot/live/opencurriculum.eus/cert.pem",
		"data/certbot/live/opencurriculum.eus/privkey.pem", mux)
	if err != nil {
		fmt.Println("Error starting server")
		panic(err)
	}
}
