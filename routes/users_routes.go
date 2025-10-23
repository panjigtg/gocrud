package routes

import (
	"crudprojectgo/app/services"
	"github.com/gofiber/fiber/v2"
)

func UsersRoutes(app *fiber.App, service *services.UsersService) {
	r := app.Group("/users")
	r.Get("/", service.GetUsers)
}
