package config

import (
	"crudprojectgo/middleware"
	"crudprojectgo/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, s *ServiceContainer) {
	app.Use(middleware.SanitizeMiddleware)

	routes.AlumniRoutes(app, s.Alumni)
	routes.AuthRoutes(app, s.Auth)
	routes.UsersRoutes(app, s.Users)
	routes.PekerjaanRoutes(app,s.Pekerjaan)
	routes.PekerjaanMongoRoutes(app, s.PekerjaanMongo)
	routes.FileUploadRoutes(app, s.FileUpload)
}

func StartServer(app *fiber.App) {
	ensureUploadDirs([]string{
		"./uploads/foto",
		"./uploads/sertifikat",
	})

	const addr = "127.0.0.1:3000"
	log.Printf("Server berjalan di http://%s", addr)
	log.Fatal(app.Listen(addr))
}

func ensureUploadDirs(paths []string) {
	for _, p := range paths {
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			log.Fatalf("Gagal membuat direktori upload: %v", err)
		}
	}
}