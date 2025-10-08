package routes

import (
	"crudprojectgo/app/models"
	"crudprojectgo/app/services"
	"crudprojectgo/helper"
    "strconv"
    
    "github.com/gofiber/fiber/v2"
)

type MahasiswaHandler struct {
    service *services.MahasiswaService
}

func NewMahasiswaHandler(service *services.MahasiswaService) *MahasiswaHandler {
    return &MahasiswaHandler{service: service}
}

func (h *MahasiswaHandler) SetupRoutes(app *fiber.App) {
    m := app.Group("/alumni-management/mahasiswa")

    m.Get("/", h.GetAllMahasiswa)
    m.Get("/:id", h.GetMahasiswaByID)
    m.Post("/", h.CreateMahasiswa)
    m.Put("/:id", h.UpdateMahasiswa)
    m.Put("/softdelete/:id", h.SoftDeletes)
    m.Delete("/:id", h.DeleteMahasiswa)
}

func (h *MahasiswaHandler) GetAllMahasiswa(c *fiber.Ctx) error {
    data, err := h.service.GetAllMahasiswa()
    if err != nil {
        return helper.ErrorResponse(c, 500, "Gagal mengambil data mahasiswa")
    }
    return helper.SuccessResponse(c, "Data mahasiswa berhasil diambil", data)
}

func (h *MahasiswaHandler) GetMahasiswaByID(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    data, err := h.service.GetMahasiswaByID(id)
    if err != nil {
        return helper.ErrorResponse(c, 404, err.Error())
    }
    return helper.SuccessResponse(c, "Data mahasiswa berhasil diambil", data)
}

func (h *MahasiswaHandler) CreateMahasiswa(c *fiber.Ctx) error {
    var req models.CreateMahasiswaRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, 400, "Request body tidak valid")
    }
    data, err := h.service.CreateMahasiswa(&req)
    if err != nil {
        return helper.ErrorResponse(c, 400, err.Error())
    }
    return helper.CreatedResponse(c, "Mahasiswa berhasil ditambahkan", data)
}

func (h *MahasiswaHandler) UpdateMahasiswa(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    var req models.UpdateMahasiswaRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, 400, "Request body tidak valid")
    }
    data, err := h.service.UpdateMahasiswa(id, &req)
    if err != nil {
        return helper.ErrorResponse(c, 400, err.Error())
    }
    return helper.SuccessResponse(c, "Mahasiswa berhasil diupdate", data)
}

func (h *MahasiswaHandler) DeleteMahasiswa(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    if err := h.service.DeleteMahasiswa(id); err != nil {
        return helper.ErrorResponse(c, 404, err.Error())
    }
    return helper.SuccessResponse(c, "Mahasiswa berhasil dihapus", nil)
}

func (h *MahasiswaHandler) SoftDeletes(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    var req models.IsDeleted
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, 400, "Request body tidak valid")
    }
    data, err := h.service.SoftDeletes(id, &req)
    if err != nil {
        return helper.ErrorResponse(c, 400, err.Error())
    }
    return helper.SuccessResponse(c, "Mahasiswa berhasil dihapus secara soft delete", data)
}