package routes

import (
    "crudprojectgo/helper"
    "crudprojectgo/app/models"
    "crudprojectgo/app/services"
    "crudprojectgo/middleware"
    "strconv"
    "strings"
    
    "github.com/gofiber/fiber/v2"
)

type PekerjaanHandler struct {
    service *services.PekerjaanService
}

func NewPekerjaanHandler(service *services.PekerjaanService) *PekerjaanHandler {
    return &PekerjaanHandler{
        service: service,
    }
}

func (h *PekerjaanHandler) SetupRoutes(app *fiber.App) {
   pekerjaan := app.Group("/alumni-management/pekerjaan")

    pekerjaan.Get("/trash", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.GetTrash)
    pekerjaan.Get("/filter", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.GetPekerjaanByFilter)
	pekerjaan.Get("/", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.GetAllPekerjaan)
	pekerjaan.Get("/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.GetPekerjaanByID)
	pekerjaan.Get("/alumni/:alumni_id", middleware.AuthRequired(), middleware.RoleOnly("admin"), h.GetPekerjaanByAlumniID)

	pekerjaan.Post("/", middleware.AuthRequired(), middleware.AdminOnly(), h.CreatePekerjaan)
    pekerjaan.Put("/restore/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.RestorePekerjaan)
	pekerjaan.Put("/:id", middleware.AuthRequired(), middleware.AdminOnly(), h.UpdatePekerjaan)
    pekerjaan.Put("/softdel/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.SoftDeletePekerjaan)
    pekerjaan.Delete("/harddel/:id", middleware.AuthRequired(), middleware.RoleOnly("admin", "user"), h.HardDeletePekerjaan)
	pekerjaan.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), h.DeletePekerjaan)
}

func (h *PekerjaanHandler) GetAllPekerjaan(c *fiber.Ctx) error {
    pekerjaan, err := h.service.GetAllPekerjaan()
    if err != nil {
        return helper.ErrorResponse(c, 500, "Gagal mengambil data pekerjaan")
    }
    
    return helper.SuccessResponse(c, "Data pekerjaan berhasil diambil", pekerjaan)
}

func (h *PekerjaanHandler) GetPekerjaanByID(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    
    pekerjaan, err := h.service.GetPekerjaanByID(id)
    if err != nil {
        return helper.ErrorResponse(c, 404, err.Error())
    }
    
    return helper.SuccessResponse(c, "Data pekerjaan berhasil diambil", pekerjaan)
}

func (h *PekerjaanHandler) GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "Alumni ID tidak valid")
	}

	pekerjaan, err := h.service.GetPekerjaanByAlumniID(alumniID)
	if err != nil {
		if err.Error() == "alumni tidak ditemukan" {
			return helper.ErrorResponse(c, 404, err.Error())
		}
		return helper.ErrorResponse(c, 500, "Terjadi kesalahan saat mengambil data pekerjaan")
	}

	// Balas tetap 200 meskipun data kosong
	if len(pekerjaan) == 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Data pekerjaan alumni tidak ditemukan",
			"data":    []models.PekerjaanAlumni{},
		})
	}

	return helper.SuccessResponse(c, "Data pekerjaan alumni berhasil diambil", pekerjaan)
}


func (h *PekerjaanHandler) CreatePekerjaan(c *fiber.Ctx) error {
    var req models.CreatePekerjaanRequest
    
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, 400, "Request body tidak valid")
    }
    
    pekerjaan, err := h.service.CreatePekerjaan(&req)
    if err != nil {
        return helper.ErrorResponse(c, 400, err.Error())
    }
    
    return helper.CreatedResponse(c, "Pekerjaan berhasil ditambahkan", pekerjaan)
}

func (h *PekerjaanHandler) UpdatePekerjaan(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    
    var req models.UpdatePekerjaanRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, 400, "Request body tidak valid")
    }
    
    pekerjaan, err := h.service.UpdatePekerjaan(id, &req)
    if err != nil {
        return helper.ErrorResponse(c, 400, err.Error())
    }
    
    return helper.SuccessResponse(c, "Pekerjaan berhasil diupdate", pekerjaan)
}

func (h *PekerjaanHandler) DeletePekerjaan(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, 400, "ID tidak valid")
    }
    
    err = h.service.DeletePekerjaan(id)
    if err != nil {
        return helper.ErrorResponse(c, 404, err.Error())
    }
    
    return helper.SuccessResponse(c, "Pekerjaan berhasil dihapus", nil)
}

func (h *PekerjaanHandler) GetPekerjaanByFilter(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "id")
    order := c.Query("order", "asc")
    search := c.Query("search", "")
    result, err := h.service.GetPekerjaanByFilter(search, sortBy, order, page, limit)
    if err != nil {
        if strings.Contains(strings.ToLower(err.Error()), "sortby") {
        return helper.ErrorResponse(c, 404, "Kolom sortBy tidak ditemukan")
        }
        return helper.ErrorResponse(c, 500, "Gagal mengambil data pekerjaan")
    }
    return helper.SuccessResponse(c, "Data pekerjaan berhasil diambil", result)
}

func (h *PekerjaanHandler) SoftDeletePekerjaan(c *fiber.Ctx) error {
	idParam := c.Params("id")
	pekerjaanID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	requesterID := c.Locals("user_id").(int)
	requesterRole := c.Locals("role").(string)

	err = h.service.SoftDeletePekerjaan(pekerjaanID, requesterID, requesterRole)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus (soft delete)"})
}

func (h *PekerjaanHandler) GetTrash(c *fiber.Ctx) error {
	
	requesterID := c.Locals("user_id").(int)
	requesterRole := c.Locals("role").(string)

	data, err := h.service.GetTrash(requesterRole, requesterID)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengambil data trash")
	}

	if len(data) == 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Tidak ada data di trash",
			"data":    []models.PekerjaanAlumni{},
		})
	}

	return helper.SuccessResponse(c, "Data trash berhasil diambil", data)
}

func (h *PekerjaanHandler) RestorePekerjaan(c *fiber.Ctx) error {
	idParam := c.Params("id")
	pekerjaanID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	requesterID := c.Locals("user_id").(int)
	requesterRole := c.Locals("role").(string)

	err = h.service.RestorePekerjaan(pekerjaanID, requesterID, requesterRole)
	if err != nil {
		return helper.ErrorResponse(c, 403, err.Error())
	}

	return helper.SuccessResponse(c, "Pekerjaan berhasil direstore", nil)
}

func (h *PekerjaanHandler) HardDeletePekerjaan(c *fiber.Ctx) error {
	idParam := c.Params("id")
	pekerjaanID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	requesterID := c.Locals("user_id").(int)
	requesterRole := c.Locals("role").(string)

	err = h.service.HardDeletePekerjaan(pekerjaanID, requesterID, requesterRole)
	if err != nil {
		return helper.ErrorResponse(c, 403, err.Error())
	}

	return helper.SuccessResponse(c, "Data pekerjaan berhasil dihapus permanen", nil)
}
