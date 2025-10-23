package config

import (
	"log"
	"crudprojectgo/database"
	"github.com/joho/godotenv"
	
	
)

func RunApp() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal load .env file")
	}

	db := database.KoneksiDB()
	mongo := database.MongoConnections()
	defer db.Close()

	repos := InitRepositories(db, mongo)
	services := InitServices(repos)
	app := SetupFiber()

	RegisterRoutes(app, services)
	StartServer(app)
}
