package routes

import (
	"crudprojectgo/app/models"
	"crudprojectgo/app/services"
	"crudprojectgo/middleware"
	"github.com/gofiber/fiber/v2"
	"time"
	"crudprojectgo/utils"
)

type AuthHandler struct {
	service *services.AuthServices
}

func NewAuthHandler(service *services.AuthServices) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) SetupRoutes(app *fiber.App) {
	auth := app.Group("/auth")

	auth.Post("/login", h.Login)
	auth.Post("/register", h.Register)
	auth.Post("/refresh", h.RefreshToken)
	auth.Get("/logout", h.Logout)
	// route dilindungi token JWT
	auth.Get("/profile", middleware.AuthRequired(), h.Profile)
	auth.Get("/admin", middleware.AuthRequired(), middleware.AdminOnly(), h.AdminDashboard)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request salah"})
	}

	res, err := h.service.Login(req)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	refreshToken, err := utils.GenerateRefreshToken(res.User)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat refresh token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
	})

		return c.JSON(fiber.Map{
		"access_token": res.Token,
		"user": fiber.Map{
			"id":       res.User.ID,
			"username": res.User.Username,
			"email":    res.User.Email,
			"role":     res.User.Role,
		},
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request salah"})
	}

	user, err := h.service.Register(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	username := c.Locals("username")
	role := c.Locals("role")
	return c.JSON(fiber.Map{
		"username": username,
		"role":     role,
	})
}

func (h *AuthHandler) AdminDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Selamat datang admin!",
	})
}
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token tidak ditemukan"})
	}

	// Validasi dan generate access token baru
	newAccessToken, err := h.service.RefreshAccessToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"access_token": newAccessToken,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Hapus cookie refresh_token
	c.ClearCookie("refresh_token")

	// (opsional) blacklist access token kalau kamu pakai sistem itu

	return c.JSON(fiber.Map{
		"message": "Berhasil logout",
	})
}
