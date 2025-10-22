package routes

import (
	"crudprojectgo/app/services"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, service *services.AuthServices) {
	auth := app.Group("/auth")

	auth.Post("/login", service.Login)
	auth.Post("/register", service.Register)
	auth.Post("/logout", service.Logout)
}
