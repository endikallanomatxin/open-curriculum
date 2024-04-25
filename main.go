package main

import (
	"app/db"
	"app/handlers"
	"fmt"
	"net/http"
	"os"
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

	if os.Getenv("ENV") == "dev" {
		http.ListenAndServe(":8080", mux)
	} else if os.Getenv("ENV") == "prod" {
		
		err := http.ListenAndServeTLS(":443",
			"/etc/letsencrypt/live/opencurriculum.eus/cert.pem",
			"/etc/letsencrypt/live/opencurriculum.eus/privkey.pem", mux)
		if err != nil {
			fmt.Println("Error starting server")
			panic(err)
		}
	}
}
