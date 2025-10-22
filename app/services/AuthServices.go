package services

import (
	"crudprojectgo/app/models"
	"crudprojectgo/app/repository/psql"
	"crudprojectgo/helper"
	"crudprojectgo/utils"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type AuthServices struct {
	repo *psql.AuthRepository
}

func NewAuthServices(repo *psql.AuthRepository) *AuthServices {
	return &AuthServices{repo: repo}
}

// --------------------
// LOGIN
// --------------------
func (s *AuthServices) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Request body tidak valid")
	}

	user, passwordHash, err := s.repo.GetUserByUsernameOrEmail(req.UsernameOrEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 401, "Email atau password salah")
		}
		return helper.ErrorResponse(c, 500, "Kesalahan server")
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		return helper.ErrorResponse(c, 401, "Password salah")
	}

	token, err := utils.GenerateAccessToken(user)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal generate token")
	}

	return helper.SuccessResponse(c, "Login berhasil", fiber.Map{
		"user":  user,
		"token": token,
	})
}

// --------------------
// REGISTER
// --------------------
func (s *AuthServices) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Request body tidak valid")
	}

	exist, _ := s.repo.CheckUserExist(req.Username, req.Email)
	if exist {
		return helper.ErrorResponse(c, 400, "Username atau email sudah digunakan")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengenkripsi password")
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	newUser, err := s.repo.CreateUserWithRole(req.Username, req.Email, hashedPassword, role)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal membuat user")
	}

	return helper.CreatedResponse(c, "Akun berhasil dibuat", newUser)
}

// --------------------
// LOGOUT
// --------------------
func (s *AuthServices) Logout(c *fiber.Ctx) error {
	// Untuk JWT stateless, logout = hapus token di sisi client.
	// Tapi kalau kamu pakai token blacklist, bisa disimpan ke DB atau Redis di sini.

	token := c.Get("Authorization")
	if token == "" {
		return helper.ErrorResponse(c, 400, "Token tidak ditemukan")
	}

	// Contoh pseudo-implementasi blacklist (opsional):
	// s.repo.BlacklistToken(token)

	return helper.SuccessResponse(c, "Logout berhasil, token tidak lagi berlaku", nil)
}
