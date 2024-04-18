package db

import "log"

type Unit struct {
	ID          int
	Name        string
	Description string
}

type Group struct {
	ID          int
	Name        string
	Description string
	GroupID     int
}

func UnitsCreateTables() {
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
	);`

	// Ejecutar los comandos SQL
	_, err := db.Exec(createTables)
	if err != nil {
		log.Fatalf("Error creating tables: %q", err)
	}
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
