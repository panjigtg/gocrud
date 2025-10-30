package routes

import (
	"crudprojectgo/app/services"
	"crudprojectgo/middleware"

	"github.com/gofiber/fiber/v2"
)

func FileUploadRoutes(app *fiber.App, service *services.FileServiceMongo) {
	r := app.Group("/api/v2/files")

	r.Post("/foto/:user_id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.UploadFoto)
	r.Post("/sertifikat/:user_id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.UploadSertifikat)
	r.Get("/", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetAllFiles)
	r.Get("/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), service.GetFileByID)
	r.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.DeleteFile)
}