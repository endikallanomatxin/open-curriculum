package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	UnitsCreateTables()
	DependenciesCreateTables()
	ChangesCreateTables()
	ProposalsCreateTables()
	PollsCreateTables()

	log.Println("Database initialized")
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Fatalf("Error closing database: %q", err)
	}
}
