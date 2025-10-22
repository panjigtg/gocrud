package config

import (
	"log"
	"crudprojectgo/routes"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, s *ServiceContainer) {
	routes.AlumniRoutes(app, s.Alumni)
	routes.AuthRoutes(app, s.Auth)
	routes.UsersRoutes(app, s.Users)
	routes.PekerjaanRoutes(app,s.Pekerjaan)
}

func StartServer(app *fiber.App) {
	const addr = "127.0.0.1:3000"
	log.Printf("Server berjalan di http://%s", addr)
	log.Fatal(app.Listen(addr))
}
