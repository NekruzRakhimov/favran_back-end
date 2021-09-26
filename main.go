package main

import (
	"favran/db"
	"favran/routes"
	"favran/utils"
)

func main() {
	utils.ReadSettings()

	db.StartDbConnection()

	routes.InitAllRoutes()
}
