package db

import (
	"app/logic"
	"log"
)

func UnitsCreateTables() {
	createTables := `
	CREATE TABLE IF NOT EXISTS groups (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		group_id BIGINT,
		FOREIGN KEY (group_id) REFERENCES groups(id)
	);

	CREATE TABLE IF NOT EXISTS units (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		content TEXT,
		group_id BIGINT,
		FOREIGN KEY (group_id) REFERENCES groups(id)
	);`

	// Ejecutar los comandos SQL
	_, err := db.Exec(createTables)
	if err != nil {
		log.Fatalf("Error creating tables: %q", err)
	}
}

func GetUnits() []logic.Unit {

	rows, err := db.Query("SELECT id, name, content FROM units")
	if err != nil {
		log.Fatalf("Error querying units: %q", err)
	}
	defer rows.Close()

	units := []logic.Unit{}
	for rows.Next() {
		var u logic.Unit
		err := rows.Scan(&u.ID, &u.Name, &u.Content)
		if err != nil {
			log.Fatalf("Error scanning units: %q", err)
		}
		units = append(units, u)
	}

	return units
}

func CreateUnit(u logic.Unit) int64 {
	var id int64
	err := db.QueryRow("INSERT INTO units (name, content) VALUES ($1, $2) RETURNING id", u.Name, u.Content).Scan(&id)
	if err != nil {
		log.Fatalf("Error creating unit: %q", err)
	}
	return id
}

func GetUnit(id int64) (logic.Unit, error) {
	var u logic.Unit
	err := db.QueryRow("SELECT id, name, content FROM units WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Content)
	return u, err
}

func DeleteUnit(id int64) {
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

func RenameUnit(id int64, name string) {
	_, err := db.Exec("UPDATE units SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		log.Fatalf("Error renaming unit: %q", err)
	}
}

func UpdateGraph() {
	// This goes over all accepted polls and updates the graph with its proposals

	// Delete all units and dependencies
	_, err := db.Exec("DELETE FROM dependencies")
	if err != nil {
		log.Fatalf("Error deleting dependencies: %q", err)
	}
	_, err = db.Exec("DELETE FROM units")
	if err != nil {
		log.Fatalf("Error deleting units: %q", err)
	}

	// Restart DB autoincrement
	_, err = db.Exec("ALTER SEQUENCE units_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatalf("Error restarting units_id_seq: %q", err)
	}
	_, err = db.Exec("ALTER SEQUENCE dependencies_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatalf("Error restarting dependencies_id_seq: %q", err)
	}

	// Get all accepted polls
	acceptedPolls := GetAcceptedPolls()
	for _, poll := range acceptedPolls {
		// Poll is an interface
		// If it a SingleProposalPoll
		if poll, ok := poll.(logic.SingleProposalPoll); ok {
			createdUnitMap := make(map[int64]int64)
			proposal := poll.Proposal
			for _, change := range proposal.Changes {
				switch change := change.(type) {
				case logic.UnitCreation:
					CreatedUnitID := CreateUnit(logic.Unit{Name: change.Name})
					createdUnitMap[change.ID] = CreatedUnitID
				case logic.UnitDeletion:
					DeleteUnit(change.UnitID)
				case logic.UnitRename:
					RenameUnit(change.UnitID, change.Name)
				case logic.DependencyCreation:
					var UnitID, DependsOnID int64
					if change.UnitIsProposed {
						UnitID = createdUnitMap[change.UnitID]
					} else {
						UnitID = change.UnitID
					}
					if change.DependsOnIsProposed {
						DependsOnID = createdUnitMap[change.DependsOnID]
					} else {
						DependsOnID = change.DependsOnID
					}
					CreateDependency(UnitID, DependsOnID)
				case logic.DependencyDeletion:
					DeleteDependency(change.DependencyID)
				}
			}
		}
	}
}
