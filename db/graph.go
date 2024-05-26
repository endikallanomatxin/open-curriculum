package db

import (
	models "app/models"
)

func GetProposedGraph(proposalID int) models.Graph {
	units := GetUnits()
	dependencies := GetAllDependencies()

	// TODO: Apply proposal changes

	return models.Graph{
		Units:        units,
		Dependencies: dependencies,
	}
}
