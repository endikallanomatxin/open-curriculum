package db

import (
	models "app/models"
	"fmt"
)

func ChangesCreateTables() {
	// use models so that it doesn't get removed
	unit := models.Unit{}
	fmt.Print(unit)
}
