package routes

import (
	"crudprojectgo/app/services"
	"crudprojectgo/middleware"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoutes(app *fiber.App, service *services.PekerjaanService) {
	r := app.Group("/api/v1/pekerjaan")

	r.Get("/trash", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetTrash)
	r.Get("/filter", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetPekerjaanByFilter)
	r.Get("/", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetAllPekerjaan)
	r.Get("/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetPekerjaanByID)
	r.Get("/alumni/:alumni_id", middleware.AuthRequired(), middleware.RoleOnly("admin"), service.GetPekerjaanByAlumniID)

	r.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), service.CreatePekerjaan)
	r.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.UpdatePekerjaan)
	r.Put("/restore/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.RestorePekerjaan)
	r.Put("/softdel/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.SoftDeletePekerjaan)
	r.Delete("/harddel/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.HardDeletePekerjaan)
}
