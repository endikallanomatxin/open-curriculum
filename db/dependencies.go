package db

import (
	"app/logic"
	"errors"
	"log"
)

func DependenciesCreateTables() {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS dependencies (
		id SERIAL PRIMARY KEY,
		unit_id BIGINT,
		depends_on_id BIGINT,
		UNIQUE (unit_id, depends_on_id),
		FOREIGN KEY (unit_id) REFERENCES units(id),
		FOREIGN KEY (depends_on_id) REFERENCES units(id)
	);`)
	if err != nil {
		log.Fatalf("Error creating table: %q", err)
	}
}

func CheckCircularDependency(unitID, dependsOnID int64) error {
	visited := make(map[int64]bool)
	path := make(map[int64]bool)

	// Verificar si hay un ciclo comenzando desde la unidad objetivo hacia la unidad fuente
	if hasCycle(dependsOnID, unitID, visited, path) {
		return errors.New("circular dependency detected")
	}

	return nil
}

func hasCycle(unitID, targetID int64, visited, path map[int64]bool) bool {
	// Marcar la unidad actual como visitada y agregarla al camino actual
	visited[unitID] = true
	path[unitID] = true

	// Obtener las dependencias de la unidad actual
	dependencies := GetUnitDependencies(unitID)

	// Iterar sobre las dependencias
	for _, dep := range dependencies {
		// Si la dependencia es la unidad objetivo, se detecta un ciclo
		if dep.ID == targetID {
			return true
		}
		// Si la dependencia no ha sido visitada aún, recursivamente verificar si hay ciclo
		if !visited[dep.ID] && hasCycle(dep.ID, targetID, visited, path) {
			return true
		} else if path[dep.ID] {
			// Si la dependencia está en el camino actual, se detecta un ciclo
			return true
		}
	}

	// Eliminar la unidad actual del camino actual
	path[unitID] = false
	return false
}

func GetUnitDependencies(unit_id int64) []logic.Unit {
	rows, err := db.Query("SELECT units.id, units.name, units.content FROM dependencies JOIN units ON dependencies.depends_on_id = units.id WHERE dependencies.unit_id = $1", unit_id)
	if err != nil {
		log.Fatalf("Error querying dependencies: %q", err)
	}
	defer rows.Close()

	units := []logic.Unit{}
	for rows.Next() {
		var u logic.Unit
		err := rows.Scan(&u.ID, &u.Name, &u.Content)
		if err != nil {
			log.Fatalf("Error scanning dependencies: %q", err)
		}
		units = append(units, u)
	}
	return units
}

func GetAllDependencies() []logic.Dependency {
	rows, err := db.Query("SELECT id, unit_id, depends_on_id FROM dependencies")
	if err != nil {
		log.Fatalf("Error querying dependencies: %q", err)
	}
	defer rows.Close()

	dependencies := []logic.Dependency{}
	for rows.Next() {
		var d logic.Dependency
		err := rows.Scan(&d.ID, &d.UnitID, &d.DependsOnID)
		if err != nil {
			log.Fatalf("Error scanning dependencies: %q", err)
		}
		dependencies = append(dependencies, d)
	}

	return dependencies
}

func GetDependency(id int64) logic.Dependency {
	var d logic.Dependency
	err := db.QueryRow("SELECT id, unit_id, depends_on_id FROM dependencies WHERE id = $1", id).Scan(&d.ID, &d.UnitID, &d.DependsOnID)
	if err != nil {
		log.Fatalf("Error querying dependency: %q", err)
	}
	return d
}

// If the dependency exists, return its ID.
// Otherwise, return 0
func FindDependency(unit_id int64, depends_on_id int64) int64 {
	rows, err := db.Query("SELECT id FROM dependencies WHERE unit_id = $1 AND depends_on_id = $2", unit_id, depends_on_id)
	if err != nil {
		log.Fatalf("Error querying dependencies: %q", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Fatalf("Error scanning dependencies: %q", err)
		}
		return id
	}

	return 0
}

func CreateDependency(unit_id int64, depends_on_id int64) error {
	// Check if the new dependency will create a circular dependency
	if err := CheckCircularDependency(unit_id, depends_on_id); err != nil {
		return err
	}

	// If no circular dependency is detected, proceed with creating the dependency
	_, err := db.Exec("INSERT INTO dependencies (unit_id, depends_on_id) VALUES ($1, $2)", unit_id, depends_on_id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteDependency(dependencyID int64) {
	_, err := db.Exec("DELETE FROM dependencies WHERE id = $1", dependencyID)
	if err != nil {
		log.Fatalf("Error deleting dependency: %q", err)
	}
}
