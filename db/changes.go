package db

import (
	"app/logic"
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
		CREATE TABLE IF NOT EXISTS document_modifications (
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			unit_id BIGINT,
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
			id BIGSERIAL PRIMARY KEY,
			proposal_id BIGINT,
			unit_id BIGINT
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

func GetProposalChanges(proposalId int64) []logic.Change {
	// Change is a generic interface
	// All change types tables have to be queried

	changes := []logic.Change{}

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
		var c logic.UnitCreation
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
		var c logic.UnitDeletion
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
		var c logic.UnitRename
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
		var c logic.DependencyCreation
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
		var c logic.DependencyDeletion
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
		var c logic.DocumentModification
		err := rows.Scan(&c.ID, &c.UnitID, &c.FromLine, &c.ToLine, &c.Content)
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

func GetUnitCreation(changeId int64) (logic.UnitCreation, error) {
	var c logic.UnitCreation
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

func GetProposedGraph(proposalID int64) logic.Graph {
	units := GetUnits()
	dependencies := GetAllDependencies()

	// Mark units as existing
	for i := range units {
		units[i].Type = "Existing"
	}

	// Mark dependencies as existing
	for i := range dependencies {
		dependencies[i].Type = "Existing"
	}

	if proposalID == 0 {
		return logic.Graph{
			Units:        units,
			Dependencies: dependencies,
		}
	}

	proposal := GetProposal(proposalID)

	// Apply proposal changes
	for _, change := range proposal.Changes {
		switch change := change.(type) {
		case logic.UnitDeletion:
			for i, unit := range units {
				if unit.ID == change.UnitID {
					unit.Type = "ProposedDeletion"
					unit.ChangeID = change.ID
					units[i] = unit
				}
			}
		case logic.UnitCreation:
			units = append(units, logic.Unit{
				ChangeID: change.ID,
				Name:     change.Name,
				Type:     "ProposedCreation",
			})
		case logic.UnitRename:
			for i, unit := range units {
				if unit.ID == change.UnitID {
					unit.Type = "ProposedRename"
					unit.ChangeID = change.ID
					unit.Name = change.Name
					units[i] = unit
				}
			}
		case logic.DependencyCreation:
			dependencies = append(dependencies, logic.Dependency{
				ID:                  change.ID,
				Type:                "ProposedCreation",
				UnitIsProposed:      change.UnitIsProposed,
				UnitID:              change.UnitID,
				DependsOnIsProposed: change.DependsOnIsProposed,
				DependsOnID:         change.DependsOnID,
			})
		case logic.DependencyDeletion:
			changedDependency := GetDependency(change.DependencyID)
			for i, dependency := range dependencies {
				if dependency.UnitID == changedDependency.UnitID &&
					dependency.DependsOnID == changedDependency.DependsOnID {
					dependency.Type = "ProposedDeletion"
					dependency.ChangeID = change.ID
					dependencies[i] = dependency
				}
			}
		}
	}

	return logic.Graph{
		Units:        units,
		Dependencies: dependencies,
	}
}
