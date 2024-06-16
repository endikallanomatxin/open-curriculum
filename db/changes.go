package db

import (
	"app/models"
	"fmt"
)

// ChangesCreateTables creates the necessary tables for storing changes
func ChangesCreateTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS unit_creations (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			name VARCHAR(255)
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS unit_deletions (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			unit_id INTEGER
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS unit_renames (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			unit_id INTEGER,
			name VARCHAR(255)
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS dependency_creations (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			unit_is_proposed BOOLEAN,
			unit_id INTEGER,
			depends_on_is_proposed BOOLEAN,
			depends_on_id INTEGER
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS dependency_deletions (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			dependency_id INTEGER
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS document_modifications (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			unit_id INTEGER,
			from_line INTEGER,
			to_line INTEGER,
			content TEXT
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS document_file_uploads (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			unit_id INTEGER
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS video_modifications (
			id SERIAL PRIMARY KEY,
			proposal_id INTEGER,
			unit_id INTEGER,
			from_time INTEGER,
			to_time INTEGER,
			content TEXT
		)
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func GetProposalChanges(proposalId int) []models.Change {
	// Change is a generic interface
	// All change types tables have to be queried

	changes := []models.Change{}

	rows, err := db.Query(`
		SELECT id, name
		FROM unit_creations
		WHERE proposal_id = $1
	`, proposalId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.UnitCreation
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			fmt.Println(err)
		}
		changes = append(changes, c)
	}

	rows, err = db.Query(`
		SELECT id, unit_id
		FROM unit_deletions
		WHERE proposal_id = $1
	`, proposalId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.UnitDeletion
		err := rows.Scan(&c.ID, &c.UnitID)
		if err != nil {
			fmt.Println(err)
		}
		changes = append(changes, c)
	}

	rows, err = db.Query(`
		SELECT id, unit_id, name
		FROM unit_renames
		WHERE proposal_id = $1
	`, proposalId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.UnitRename
		err := rows.Scan(&c.ID, &c.UnitID, &c.Name)
		if err != nil {
			fmt.Println(err)
		}
		changes = append(changes, c)
	}

	rows, err = db.Query(`
		SELECT id, unit_is_proposed, unit_id, depends_on_is_proposed, depends_on_id
		FROM dependency_creations
		WHERE proposal_id = $1
	`, proposalId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.DependencyCreation
		err := rows.Scan(&c.ID, &c.UnitIsProposed, &c.UnitID, &c.DependsOnIsProposed, &c.DependsOnID)
		if err != nil {
			fmt.Println(err)
		}
		changes = append(changes, c)
	}

	rows, err = db.Query(`
		SELECT id, dependency_id
		FROM dependency_deletions
		WHERE proposal_id = $1
	`, proposalId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.DependencyDeletion
		err := rows.Scan(&c.ID, &c.DependencyID)
		if err != nil {
			fmt.Println(err)
		}
		changes = append(changes, c)
	}

	rows, err = db.Query(`
		SELECT id, unit_id, from_line, to_line, content
		FROM document_modifications
		WHERE proposal_id = $1
	`, proposalId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.DocumentModification
		err := rows.Scan(&c.ID, &c.UnitID, &c.FromLine, &c.ToLine, &c.Content)
		if err != nil {
			fmt.Println(err)
		}
		changes = append(changes, c)
	}

	return changes
}

func CreateUnitCreation(proposalId int, name string) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO unit_creations (proposal_id, name)
		VALUES ($1, $2)
		RETURNING id
	`, proposalId, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func DeleteUnitCreation(changeId int) error {
	_, err := db.Exec(`
		DELETE FROM unit_creations
		WHERE id = $1
	`, changeId)
	if err != nil {
		return err
	}
	return nil
}

func GetUnitCreation(changeId int) (models.UnitCreation, error) {
	var c models.UnitCreation
	err := db.QueryRow(`
		SELECT id, name
		FROM unit_creations
		WHERE id = $1
	`, changeId).Scan(&c.ID, &c.Name)
	if err != nil {
		return c, err
	}
	return c, nil
}

func CreateUnitDeletion(proposalId int, unitId int) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO unit_deletions (proposal_id, unit_id)
		VALUES ($1, $2)
		RETURNING id
	`, proposalId, unitId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func DeleteUnitDeletion(changeId int) error {
	_, err := db.Exec(`
		DELETE FROM unit_deletions
		WHERE id = $1
	`, changeId)
	if err != nil {
		return err
	}
	return nil
}

func CreateUnitRename(proposalId int, unitId int, name string) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO unit_renames (proposal_id, unit_id, name)
		VALUES ($1, $2, $3)
		RETURNING id
	`, proposalId, unitId, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func DeleteUnitRename(changeId int) error {
	_, err := db.Exec(`
		DELETE FROM unit_renames
		WHERE id = $1
	`, changeId)
	if err != nil {
		return err
	}
	return nil
}

func CreateDependencyCreation(proposalID int, unitIsProposed bool, unitID int, dependsOnIsProposed bool, dependsOnID int) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO dependency_creations (proposal_id, unit_is_proposed, unit_id, depends_on_is_proposed, depends_on_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, proposalID, unitIsProposed, unitID, dependsOnIsProposed, dependsOnID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func DeleteDependencyCreation(changeID int) error {
	_, err := db.Exec(`
		DELETE FROM dependency_creations
		WHERE id = $1
	`, changeID)
	if err != nil {
		return err
	}
	return nil
}

func CreateDependencyDeletion(proposalID int, dependencyID int) error {
	_, err := db.Exec(`
		INSERT INTO dependency_deletions (proposal_id, dependency_id)
		VALUES ($1, $2)
	`, proposalID, dependencyID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDependencyDeletion(changeID int) error {
	_, err := db.Exec(`
		DELETE FROM dependency_deletions
		WHERE id = $1
	`, changeID)
	if err != nil {
		return err
	}
	return nil
}
