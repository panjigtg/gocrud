package config

import (
	"log"
	"os"

	"crudprojectgo/database"
	docs "crudprojectgo/docs" // swagger docs
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

	// Swagger info (tanpa ubah struktur lama)
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	docs.SwaggerInfo.Host = host + ":" + port
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	RegisterRoutes(app, services)
	StartServer(app)
}
