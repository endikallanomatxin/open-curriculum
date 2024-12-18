package db

import (
	"app/internal/models"
	"fmt"
)

// ChangesCreateTables creates the necessary tables for storing changes
func ChangesCreateTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS unit_creations (
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			name VARCHAR(255)
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS unit_deletions (
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			unit_id BIGINT
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS unit_renames (
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			unit_id BIGINT,
			name VARCHAR(255)
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS dependency_creations (
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			unit_is_proposed BOOLEAN,
			unit_id BIGINT,
			depends_on_is_proposed BOOLEAN,
			depends_on_id BIGINT
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS dependency_deletions (
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			dependency_id BIGINT
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS content_modifications (
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			unit_is_proposed BOOLEAN,
			unit_id BIGINT,
			content TEXT
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS content_file_uploads (
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			unit_id BIGINT,
			file_name VARCHAR(255),
			file_path VARCHAR(255)
		)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS video_modifications (
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			unit_id BIGINT,
			from_time INTEGER,
			to_time INTEGER,
			content TEXT
		)
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func GetProposalChanges(proposalId int64) []models.Change {
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
		SELECT id, unit_is_proposed, unit_id, content
		FROM content_modifications
		WHERE proposal_id = $1
	`, proposalId)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.ContentModification
		err := rows.Scan(&c.ID, &c.UnitIsProposed, &c.UnitID, &c.Content)
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
		var c models.ContentModification
		err := rows.Scan(&c.ID, &c.UnitID, &c.Content)
		if err != nil {
			fmt.Println(err)
		}
		changes = append(changes, c)
	}

	return changes
}

func CreateUnitCreation(proposalId int64, name string) (int64, error) {
	var id int64
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

func UpdateUnitCreation(changeId int64, name string) error {
	_, err := db.Exec(`
		UPDATE unit_creations
		SET name = $1
		WHERE id = $2
	`, name, changeId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUnitCreation(changeId int64) error {
	_, err := db.Exec(`
		DELETE FROM unit_creations
		WHERE id = $1
	`, changeId)
	if err != nil {
		return err
	}
	return nil
}

func GetUnitCreation(changeId int64) (models.UnitCreation, error) {
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

func CreateUnitDeletion(proposalId int64, unitId int64) (int64, error) {
	var id int64
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

func DeleteUnitDeletion(changeId int64) error {
	_, err := db.Exec(`
		DELETE FROM unit_deletions
		WHERE id = $1
	`, changeId)
	if err != nil {
		return err
	}
	return nil
}

func CreateUnitRename(proposalId int64, unitId int64, name string) (int64, error) {
	var id int64
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

func DeleteUnitRename(changeId int64) error {
	_, err := db.Exec(`
		DELETE FROM unit_renames
		WHERE id = $1
	`, changeId)
	if err != nil {
		return err
	}
	return nil
}

func CreateContentModification(proposalId int64, unitIsProposed bool, unitId int64, content string) error {
	_, err := db.Exec(`
		INSERT INTO content_modifications (proposal_id, unit_is_proposed, unit_id, content)
		VALUES ($1, $2, $3, $4)
	`, proposalId, unitIsProposed, unitId, content)
	if err != nil {
		return err
	}
	return nil
}

func DeleteContentModification(changeId int64) error {
	_, err := db.Exec(`
		DELETE FROM content_modifications
		WHERE id = $1
	`, changeId)
	if err != nil {
		return err
	}
	return nil
}

func FindContentModification(proposalId int64, unitIsProposed bool, unitId int64) int64 {
	var id int64
	err := db.QueryRow(`
		SELECT id
		FROM content_modifications
		WHERE proposal_id = $1 AND unit_is_proposed = $2 AND unit_id = $3
	`, proposalId, unitIsProposed, unitId).Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func CreateDependencyCreation(proposalID int64, unitIsProposed bool, unitID int64, dependsOnIsProposed bool, dependsOnID int64) (int64, error) {
	var id int64
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

func DeleteDependencyCreation(changeID int64) error {
	_, err := db.Exec(`
		DELETE FROM dependency_creations
		WHERE id = $1
	`, changeID)
	if err != nil {
		return err
	}
	return nil
}

func CreateDependencyDeletion(proposalID int64, dependencyID int64) error {
	_, err := db.Exec(`
		INSERT INTO dependency_deletions (proposal_id, dependency_id)
		VALUES ($1, $2)
	`, proposalID, dependencyID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDependencyDeletion(changeID int64) error {
	_, err := db.Exec(`
		DELETE FROM dependency_deletions
		WHERE id = $1
	`, changeID)
	if err != nil {
		return err
	}
	return nil
}
