package routes

import (
	"crudprojectgo/app/services"
	"crudprojectgo/middleware"

	"github.com/gofiber/fiber/v2"
)

func AlumniRoutes(app *fiber.App, service *services.AlumniService) {
	alumni := app.Group("/api/v1/alumni")

	alumni.Get("/", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetAllAlumniByFilter)
	alumni.Get("/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetAlumniByID)
	alumni.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), service.CreateAlumni)
	alumni.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.UpdateAlumni)
	alumni.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.DeleteAlumni)
}
