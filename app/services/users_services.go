package services

import (
	"strings"
	"strconv"

	"crudprojectgo/app/models"
	"crudprojectgo/app/repository/psql"
	"crudprojectgo/helper"

	"github.com/gofiber/fiber/v2"
)

type UsersService struct {
	repo *psql.UsersRepository
}

func NewUsersService(repo *psql.UsersRepository) *UsersService {
	return &UsersService{repo: repo}
}

// HandleGetAllUsers godoc
// @Summary Dapatkan semua user
// @Description Mengambil seluruh data user dari database dengan pagination, sorting, dan pencarian
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Nomor halaman (default 1)"
// @Param limit query int false "Jumlah data per halaman (default 10)"
// @Param sortBy query string false "Kolom pengurutan" Enums(id,name,email,created_at)
// @Param order query string false "Arah pengurutan" Enums(asc,desc)
// @Param search query string false "Kata kunci pencarian"
// @Success 200 {object} models.UsersListOK
// @Failure 400 {object} models.ErrorPayload
// @Failure 500 {object} models.ErrorPayload
// @Router /users/ [get]
func (s *UsersService) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

	// validasi sortBy & order
	sortByWhitelist := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"created_at": true,
	}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	users, err := s.repo.GetUsersRepo(search, sortBy, order, page, limit, offset)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengambil data user")
	}

	total, err := s.repo.CountUsersRepo(search)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal menghitung jumlah user")
	}

	response := models.UserResponse{
		Data: users,
		Meta: models.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return helper.SuccessResponse(c, "Data user berhasil diambil", response)
}
