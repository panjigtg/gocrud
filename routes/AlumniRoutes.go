package routes

import (
   	"crudprojectgo/app/models"
	"crudprojectgo/app/services"
	"crudprojectgo/helper"
    "crudprojectgo/middleware"
    "strconv"
    "strings"
    
    "github.com/gofiber/fiber/v2"
)

type AlumniHandler struct {
    service *services.AlumniService
}

func NewAlumniHandler(service *services.AlumniService) *AlumniHandler {
    return &AlumniHandler{
        service: service,
    }
}

func (h *AlumniHandler) SetupRoutes(app *fiber.App) {
    alumni := app.Group("/alumni-management/alumni")

    alumni.Get("/", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.GetAllAlumniByFilter)
	alumni.Get("/", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.GetAllAlumni)
	alumni.Get("/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.GetAlumniByID)
	alumni.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), h.CreateAlumni)
	alumni.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), h.UpdateAlumni)
	alumni.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), h.DeleteAlumni)
}

func (h *AlumniHandler) GetAllAlumniByFilter(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "created_at")
    order := c.Query("order", "desc")
    search := c.Query("search", "")
    result, err := h.service.GetAllAlumniByFilter(page, limit, sortBy, order, search)
    if err != nil {
        if strings.Contains(err.Error(), "sortBy") {
        return helper.ErrorResponse(c, 404, "Kolom sortBy tidak ditemukan")
        }
        return helper.ErrorResponse(c, 500, "Gagal mengambil data alumni")
    }
    return helper.SuccessResponse(c, "Data alumni berhasil diambil", result)
}

func (h *AlumniHandler) GetAllAlumni(c *fiber.Ctx) error {
    alumni, err := h.service.GetAllAlumni()
    if err != nil {
        return helper.ErrorResponse(c, 500, "Gagal mengambil data alumni")
    }
    
    return helper.SuccessResponse(c, "Data alumni berhasil diambil", alumni)
}

func (h *AlumniHandler) GetAlumniByID(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    
    alumni, err := h.service.GetAlumniByID(id)
    if err != nil {
        return helper.ErrorResponse(c, 404, err.Error())
    }
    
    return helper.SuccessResponse(c, "Data alumni berhasil diambil", alumni)
}

func (h *AlumniHandler) CreateAlumni(c *fiber.Ctx) error {
    var req models.CreateAlumniRequest
    
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, 400, "Request body tidak valid")
    }
    
    alumni, err := h.service.CreateAlumni(&req)
    if err != nil {
        return helper.ErrorResponse(c, 400, err.Error())
    }
    
    return helper.CreatedResponse(c, "Alumni berhasil ditambahkan", alumni)
}

func (h *AlumniHandler) UpdateAlumni(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    
    var req models.UpdateAlumniRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, 400, "Request body tidak valid")
    }
    
    alumni, err := h.service.UpdateAlumni(id, &req)
    if err != nil {
        return helper.ErrorResponse(c, 400, err.Error())
    }
    
    return helper.SuccessResponse(c, "Alumni berhasil diupdate", alumni)
}

func (h *AlumniHandler) DeleteAlumni(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    
    err = h.service.DeleteAlumni(id)
    if err != nil {
        return helper.ErrorResponse(c, 404, err.Error())
    }
    
    return helper.SuccessResponse(c, "Alumni berhasil dihapus", nil)
}
