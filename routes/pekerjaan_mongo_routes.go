package routes

import (
	"crudprojectgo/app/services"
	"crudprojectgo/middleware"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanMongoRoutes(app *fiber.App, service *services.PekerjaanServiceMongo) {
	r := app.Group("/api/v2/pekerjaan")

	r.Get("/", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetAll)
	r.Get("/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetByID)
	r.Get("/alumni/:alumni_id", middleware.AuthRequired(), middleware.RoleOnly("admin"), service.GetByAlumniID)
	r.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), service.Create)
	r.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.Update)
	r.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.Delete)
}
