package routes

import (
   	// "crudprojectgo/app/models"
	"crudprojectgo/app/services"
	"crudprojectgo/helper"
    // "strconv"
    "fmt"
    
    "github.com/gofiber/fiber/v2"
)

type ChallHandler struct {
    service *services.ChallServices
}

func NewChallHandler(service *services.ChallServices) *ChallHandler {
    return &ChallHandler{
        service: service,
    }
}

func (h *ChallHandler) SetupRoutes(app *fiber.App) {
    alumni := app.Group("/alumni-management")
    
    alumni.Get("/chall", h.GetAllChall)
}

func (h *ChallHandler) GetAllChall(c *fiber.Ctx) error {
    alumni, err := h.service.GetAllChall()
    if err != nil {
        fmt.Println("ERROR:", err)
        return helper.ErrorResponse(c, 500, "Gagal mengambil data alumni")
    }
    
    return helper.SuccessResponse(c, "Data alumni berhasil diambil", alumni)
}