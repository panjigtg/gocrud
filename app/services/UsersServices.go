package services

import (
	"strings"

	"crudprojectgo/app/models"
	"crudprojectgo/app/repository"
)

type UsersService struct {
	repo *repository.UsersRepository
}

func NewUsersService(repo *repository.UsersRepository) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s *UsersService) GetUsersService(page, limit int, sortBy, order, search string) (models.UserResponse, error) {
	offset := (page - 1) * limit

	// Validasi input sortBy dan order
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
		return models.UserResponse{}, err
	}

	total, err := s.repo.CountUsersRepo(search)
	if err != nil {
		return models.UserResponse{}, err
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

	return response, nil
}
