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

	createTables := `
    CREATE TABLE IF NOT EXISTS groups (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        group_id INT,
        FOREIGN KEY (group_id) REFERENCES groups(id)
    );

    CREATE TABLE IF NOT EXISTS units (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        group_id INT,
        FOREIGN KEY (group_id) REFERENCES groups(id)
    );

    CREATE TABLE IF NOT EXISTS dependencies (
        unit_id INT,
        depends_on_id INT,
        PRIMARY KEY (unit_id, depends_on_id),
        FOREIGN KEY (unit_id) REFERENCES units(id),
        FOREIGN KEY (depends_on_id) REFERENCES units(id)
    );`

	// Ejecutar los comandos SQL
	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatalf("Error creating tables: %q", err)
	}

	log.Println("Database initialized")
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Fatalf("Error closing database: %q", err)
	}
}

type Unit struct {
	ID          int
	Name        string
	Description string
}

func GetUnits() []Unit {

	rows, err := db.Query("SELECT id, name, description FROM units")
	if err != nil {
		log.Fatalf("Error querying units: %q", err)
	}
	defer rows.Close()

	units := []Unit{}
	for rows.Next() {
		var u Unit
		err := rows.Scan(&u.ID, &u.Name, &u.Description)
		if err != nil {
			log.Fatalf("Error scanning units: %q", err)
		}
		units = append(units, u)
	}

	return units
}

func CreateUnit(u Unit) {
	_, err := db.Exec("INSERT INTO units (name, description) VALUES ($1, $2)", u.Name, u.Description)
	if err != nil {
		log.Fatalf("Error creating unit: %q", err)
	}
}

func GetUnit(id int) Unit {
	var u Unit
	err := db.QueryRow("SELECT id, name, description FROM units WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Description)
	if err != nil {
		log.Fatalf("Error querying unit: %q", err)
	}
	return u
}

func DeleteUnit(id int) {
	// First delete the dependencies
	_, err := db.Exec("DELETE FROM dependencies WHERE unit_id = $1 OR depends_on_id = $1", id)
	if err != nil {
		log.Fatalf("Error deleting dependencies: %q", err)
	}
	_, err = db.Exec("DELETE FROM units WHERE id = $1", id)
	if err != nil {
		log.Fatalf("Error deleting unit: %q", err)
	}
}

type Dependency struct {
	ID          int
	DependsOnID int
}

func CheckDependency(id int, depends_on_id int) bool {

	// Check for self-dependency
	if id == depends_on_id {
		return false
	}

	// Check for circular dependency

	// Follow the dependencies until the end, if it finds the original unit, it's a circular dependency
	// TODO

	return true

}

func GetUnitDependencies(id int) []Unit {
	rows, err := db.Query("SELECT u.id, u.name, u.description FROM units u JOIN dependencies d ON u.id = d.depends_on_id WHERE d.unit_id = $1", id)
	if err != nil {
		log.Fatalf("Error querying dependencies: %q", err)
	}
	defer rows.Close()

	units := []Unit{}
	for rows.Next() {
		var u Unit
		err := rows.Scan(&u.ID, &u.Name, &u.Description)
		if err != nil {
			log.Fatalf("Error scanning dependencies: %q", err)
		}
		units = append(units, u)
	}

	return units
}

func GetAllDependencies() []Dependency {
	rows, err := db.Query("SELECT unit_id, depends_on_id FROM dependencies")
	if err != nil {
		log.Fatalf("Error querying dependencies: %q", err)
	}
	defer rows.Close()

	dependencies := []Dependency{}
	for rows.Next() {
		var d Dependency
		err := rows.Scan(&d.ID, &d.DependsOnID)
		if err != nil {
			log.Fatalf("Error scanning dependencies: %q", err)
		}
		dependencies = append(dependencies, d)
	}

	return dependencies
}

func CreateDependency(unit_id int, depends_on_id int) {
	if !CheckDependency(unit_id, depends_on_id) {
		log.Fatalf("Error creating dependency: Circular dependency")
		return
	}
	_, err := db.Exec("INSERT INTO dependencies (unit_id, depends_on_id) VALUES ($1, $2)", unit_id, depends_on_id)
	if err != nil {
		log.Fatalf("Error creating dependency: %q", err)
	}
}

func DeleteDependency(unit_id int, depends_on_id int) {
	_, err := db.Exec("DELETE FROM dependencies WHERE unit_id = $1 AND depends_on_id = $2", unit_id, depends_on_id)
	if err != nil {
		log.Fatalf("Error deleting dependency: %q", err)
	}
}
