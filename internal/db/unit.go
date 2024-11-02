package db

import (
	"fmt"
	"log"

	"app/internal/models"
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

func CreateUnit(u models.Unit) int64 {
	var id int64
	err := db.QueryRow("INSERT INTO units (name, content) VALUES ($1, $2) RETURNING id", u.Name, u.Content).Scan(&id)
	if err != nil {
		log.Fatalf("Error creating unit: %q", err)
	}
	return id
}

func GetUnit(id int64) (models.Unit, error) {
	var u models.Unit
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

func ModifyUnitContent(id int64, content string) {
	_, err := db.Exec("UPDATE units SET content = $1 WHERE id = $2", content, id)
	if err != nil {
		log.Fatalf("Error modifying unit content: %q", err)
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
	fmt.Println(acceptedPolls)
	for _, poll := range acceptedPolls {
		// Poll is an interface
		// If it a SingleProposalPoll
		if poll, ok := poll.(models.SingleProposalPoll); ok {
			createdUnitMap := make(map[int64]int64)
			proposal := poll.Proposal
			for _, change := range proposal.Changes {
				switch change := change.(type) {
				case models.UnitCreation:
					CreatedUnitID := CreateUnit(models.Unit{Name: change.Name})
					createdUnitMap[change.ID] = CreatedUnitID
				case models.UnitDeletion:
					DeleteUnit(change.UnitID)
				case models.UnitRename:
					RenameUnit(change.UnitID, change.Name)
				case models.ContentModification:
					fmt.Println("ContentModification")
					fmt.Println(change)
					var UnitID int64
					if change.UnitIsProposed {
						UnitID = createdUnitMap[change.UnitID]
					} else {
						UnitID = change.UnitID
					}
					ModifyUnitContent(UnitID, change.Content)
				case models.DependencyCreation:
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
				case models.DependencyDeletion:
					DeleteDependency(change.DependencyID)
				}
			}
		}
	}
}
