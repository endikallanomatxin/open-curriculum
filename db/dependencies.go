package db

import (
	"log"
)

type Dependency struct {
	ID          int
	UnitID      int
	DependsOnID int
}

func DependenciesCreateTables() {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS dependencies (
		id SERIAL PRIMARY KEY,
		unit_id INT,
		depends_on_id INT,
		UNIQUE (unit_id, depends_on_id),
		FOREIGN KEY (unit_id) REFERENCES units(id),
		FOREIGN KEY (depends_on_id) REFERENCES units(id)
	);`)
	if err != nil {
		log.Fatalf("Error creating table: %q", err)
	}
}

func CheckDependency(unit_id int, depends_on_id int) bool {

	// Check for self-dependency
	if unit_id == depends_on_id {
		return false
	}

	// Check for circular dependency

	// Follow the dependencies until the end, if it finds the original unit, it's a circular dependency
	// TODO

	return true

}

func GetUnitDependencies(unit_id int) []Unit {
	rows, err := db.Query("SELECT units.id, units.name, units.description FROM dependencies JOIN units ON dependencies.depends_on_id = units.id WHERE dependencies.unit_id = $1", unit_id)
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
	rows, err := db.Query("SELECT id, unit_id, depends_on_id FROM dependencies")
	if err != nil {
		log.Fatalf("Error querying dependencies: %q", err)
	}
	defer rows.Close()

	dependencies := []Dependency{}
	for rows.Next() {
		var d Dependency
		err := rows.Scan(&d.ID, &d.UnitID, &d.DependsOnID)
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
