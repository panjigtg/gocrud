package routes

import (
	"crudprojectgo/app/services"

	"github.com/gofiber/fiber/v2"
)

func MahasiswaRoutes(app *fiber.App, service *services.MahasiswaService) {
	m := app.Group("/api/mahasiswa")

	m.Get("/", service.GetAllMahasiswa)
	m.Get("/:id", service.GetMahasiswaByID)
	m.Post("/", service.CreateMahasiswa)
	m.Put("/:id", service.UpdateMahasiswa)
	m.Delete("/:id", service.DeleteMahasiswa)
	m.Put("/softdelete/:id", service.SoftDeletes)
}
