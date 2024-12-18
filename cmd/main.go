package main

import (
	"app/internal/db"
	"app/internal/handlers"
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
	// USERS
	// ----------------
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/register", handlers.Register)
	// mux.HandleFunc("/forgot-password", handlers.ForgotPassword)
	// mux.HandleFunc("/reset-password", handlers.ResetPassword)
	// mux.HandleFunc("/profile", handlers.Profile)

	// ----------------
	// SETTINGS
	// ----------------
	mux.HandleFunc("/settings", handlers.Settings)

	// ----------------
	// LEARN
	// ----------------

	mux.HandleFunc("/learn", handlers.Learn)

	// ----------------
	// CURRICULUM MODIFICATION
	// ----------------

	mux.HandleFunc("/curriculum-modification", handlers.CurriculumModification)

	mux.HandleFunc("/curriculum-modification/set-active-proposal-ID", handlers.SetActiveProposalID)
	mux.HandleFunc("PUT /set-open-unit", handlers.SetOpenUnit)

	// Proposal
	mux.HandleFunc("POST /curriculum-modification/proposal/create", handlers.CreateProposal)
	mux.HandleFunc("PUT /curriculum-modification/proposal/{id}/update", handlers.UpdateProposal)
	mux.HandleFunc("DELETE /curriculum-modification/proposal/{id}", handlers.DeleteProposal)
	mux.HandleFunc("PUT /curriculum-modification/proposal/{id}/submit", handlers.SubmitProposal)

	// Changes
	mux.HandleFunc("POST /curriculum-modification/proposal/{id}/unit_creation", handlers.CreateUnitCreation)
	mux.HandleFunc("PUT /curriculum-modification/proposal/{id}/unit_creation/{change_id}", handlers.UpdateUnitCreation)
	mux.HandleFunc("DELETE /curriculum-modification/proposal/{id}/unit_creation/{change_id}", handlers.DeleteUnitCreation)

	mux.HandleFunc("PUT /curriculum-modification/proposal/{id}/unit_deletion/{unit_id}", handlers.CreateUnitDeletion)
	mux.HandleFunc("DELETE /curriculum-modification/proposal/{id}/unit_deletion/{change_id}", handlers.DeleteUnitDeletion)

	mux.HandleFunc("PUT /curriculum-modification/proposal/{id}/unit_rename/{unit_id}", handlers.CreateUnitRename)
	mux.HandleFunc("DELETE /curriculum-modification/proposal/{id}/unit_rename/{change_id}", handlers.DeleteUnitRename)

	mux.HandleFunc("PUT /curriculum-modification/proposal/{id}/content_modification/", handlers.CreateContentModification)
	mux.HandleFunc("DELETE /curriculum-modification/proposal/{id}/content_modification/{change_id}", handlers.DeleteContentModification)

	mux.HandleFunc("POST /curriculum-modification/proposal/{id}/toggle_dependency", handlers.ToggleDependency)
	mux.HandleFunc("DELETE /curriculum-modification/proposal/{id}/dependency_creation/{change_id}", handlers.DeleteDependencyCreation)
	mux.HandleFunc("DELETE /curriculum-modification/proposal/{id}/dependency_deletion/{change_id}", handlers.DeleteDependencyDeletion)

	// Polls
	mux.HandleFunc("GET /curriculum-modification/polls", handlers.Polls)
	mux.HandleFunc("GET /curriculum-modification/poll/{id}", handlers.Poll)
	mux.HandleFunc("POST /curriculum-modification/poll/{id}/yes", handlers.VoteYes)
	mux.HandleFunc("POST /curriculum-modification/poll/{id}/no", handlers.VoteNo)

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
