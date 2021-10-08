package main

import (
	"github.com/NekruzRakhimov/favran/routes"
	"github.com/NekruzRakhimov/favran/utils"
)

func main() {
	utils.ReadSettings()

	//db.StartDbConnection()

	routes.InitAllRoutes()
}
