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


func GetUnitsByLevel() map[int][]Unit {

	// Units by levels based on their dependencies
	unassignedUnits := GetUnits()
	dependencies := GetAllDependencies()

	unitsByLevel := make(map[int][]Unit)

	// Iterate over the units and assign them to a level
	// Units array will be empty when all units are assigned to a level
	for level := 0; len(unassignedUnits) > 0; level++ {
		unitsByLevel[level] = []Unit{}

		for _, checkingU := range unassignedUnits {
			dependsOnUnassigned := false

			for _, d := range dependencies {
				for _, otherU := range unassignedUnits {
					if checkingU.ID == d.UnitID && otherU.ID == d.DependsOnID {
						dependsOnUnassigned = true
						break
					}
				}
			}

			if !dependsOnUnassigned {
				unitsByLevel[level] = append(unitsByLevel[level], checkingU)
			}
		}

		// Remove assigned units from the unassigned units array
		for _, u := range unitsByLevel[level] {
			for i, unassignedU := range unassignedUnits {
				if u.ID == unassignedU.ID {
					unassignedUnits = append(unassignedUnits[:i], unassignedUnits[i+1:]...)
					break
				}
			}
		}
	}

	return unitsByLevel
}
