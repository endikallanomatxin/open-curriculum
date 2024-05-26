package db

import (
	"log"

	models "app/models"
)

func UnitsCreateTables() {
	createTables := `
	CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		group_id INT,
		FOREIGN KEY (group_id) REFERENCES groups(id)
	);

	CREATE TABLE IF NOT EXISTS units (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		content TEXT,
		group_id INT,
		FOREIGN KEY (group_id) REFERENCES groups(id)
	);`

	// Ejecutar los comandos SQL
	_, err := db.Exec(createTables)
	if err != nil {
		log.Fatalf("Error creating tables: %q", err)
	}
}

func GetUnits() []models.Unit {

	rows, err := db.Query("SELECT id, name, content FROM units")
	if err != nil {
		log.Fatalf("Error querying units: %q", err)
	}
	defer rows.Close()

	units := []models.Unit{}
	for rows.Next() {
		var u models.Unit
		err := rows.Scan(&u.ID, &u.Name, &u.Content)
		if err != nil {
			log.Fatalf("Error scanning units: %q", err)
		}
		units = append(units, u)
	}

	return units
}

func CreateUnit(u models.Unit) {
	_, err := db.Exec("INSERT INTO units (name, content) VALUES ($1, $2)", u.Name, u.Content)
	if err != nil {
		log.Fatalf("Error creating unit: %q", err)
	}
}

func GetUnit(id int) models.Unit {
	var u models.Unit
	err := db.QueryRow("SELECT id, name, content FROM units WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Content)
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
