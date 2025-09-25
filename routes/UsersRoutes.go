package routes

import (
	"strconv"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"crudprojectgo/app/services"
)

type UsersHandler struct {
	service *services.UsersService
}

func NewUsersHandler(service *services.UsersService) *UsersHandler {
	return &UsersHandler{service: service}
}

func (h *UsersHandler) SetupRoutes(app *fiber.App) {
	users := app.Group("/users") // base route

	users.Get("/", h.GetUsers)
}

// Handler function
func (h *UsersHandler) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	result, err := h.service.GetUsersService(page, limit, sortBy, order, search)
	if err != nil {
		fmt.Println("[USERS HANDLER ERROR]:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data users"})
	}

	return c.JSON(result)
}
