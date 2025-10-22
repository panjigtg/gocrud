package services

import (
	"strings"
	"strconv"

	"crudprojectgo/app/models"
	"crudprojectgo/app/repository"
	"crudprojectgo/helper"

	"github.com/gofiber/fiber/v2"
)

type UsersService struct {
	repo *repository.UsersRepository
}

func NewUsersService(repo *repository.UsersRepository) *UsersService {
	return &UsersService{repo: repo}
}

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
