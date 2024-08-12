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
	mux.HandleFunc("/foundation", handlers.Foundation)

	mux.HandleFunc("/set-language", handlers.SetLanguageCookie)

	mux.HandleFunc("GET /unit/{id}/details", handlers.GetUnitDetails)

	// ----------------

	// LEARN
	mux.HandleFunc("/learn", handlers.Learn)

	// ----------------

	// TEACH
	mux.HandleFunc("/teach", handlers.Teach)
	mux.HandleFunc("/teach/set-active-proposal-ID", handlers.SetActiveProposalID)
	mux.HandleFunc("PUT /teach/set-open-unit", handlers.SetOpenUnit)

	// Proposal
	mux.HandleFunc("POST /teach/proposal/create", handlers.CreateProposal)
	mux.HandleFunc("PUT /teach/proposal/{id}/update", handlers.UpdateProposal)
	mux.HandleFunc("DELETE /teach/proposal/{id}", handlers.DeleteProposal)
	mux.HandleFunc("PUT /teach/proposal/{id}/submit", handlers.SubmitProposal)

	// Changes
	mux.HandleFunc("POST /teach/proposal/{id}/unit_creation", handlers.CreateUnitCreation)
	mux.HandleFunc("PUT /teach/proposal/{id}/unit_creation/{change_id}", handlers.UpdateUnitCreation)
	mux.HandleFunc("DELETE /teach/proposal/{id}/unit_creation/{change_id}", handlers.DeleteUnitCreation)

	mux.HandleFunc("PUT /teach/proposal/{id}/unit_deletion/{unit_id}", handlers.CreateUnitDeletion)
	mux.HandleFunc("DELETE /teach/proposal/{id}/unit_deletion/{change_id}", handlers.DeleteUnitDeletion)

	mux.HandleFunc("PUT /teach/proposal/{id}/unit_rename/{unit_id}", handlers.CreateUnitRename)
	mux.HandleFunc("DELETE /teach/proposal/{id}/unit_rename/{change_id}", handlers.DeleteUnitRename)

	mux.HandleFunc("POST /teach/proposal/{id}/toggle_dependency", handlers.ToggleDependency)
	mux.HandleFunc("DELETE /teach/proposal/{id}/dependency_creation/{change_id}", handlers.DeleteDependencyCreation)
	mux.HandleFunc("DELETE /teach/proposal/{id}/dependency_deletion/{change_id}", handlers.DeleteDependencyDeletion)

	// Polls
	mux.HandleFunc("GET /teach/polls", handlers.Polls)
	mux.HandleFunc("GET /teach/poll/{id}", handlers.Poll)
	mux.HandleFunc("POST /teach/poll/{id}/yes", handlers.VoteYes)
	mux.HandleFunc("POST /teach/poll/{id}/no", handlers.VoteNo)

	mux.HandleFunc("/units", handlers.Units)
	mux.HandleFunc("/unit/{id}", handlers.Unit)
	mux.HandleFunc("/dependencies", handlers.Dependencies)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

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
