package main

import (
	"app/db"
	"app/handlers"
	"crypto/tls"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
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

	mux.HandleFunc("/set-language", handlers.SetLanguageCookie)

	if os.Getenv("ENV") == "dev" {
		http.ListenAndServe(":8080", mux)
	} else {
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("/certs"),
		}

		server := &http.Server{
			Addr:    ":443",
			Handler: mux,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		server.ListenAndServeTLS("", "")
	}
}
