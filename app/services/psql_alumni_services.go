package services

import (
	"crudprojectgo/app/models"
	"crudprojectgo/app/repository/psql"
	"crudprojectgo/helper"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	repo *psql.AlumniRepository
}

func NewAlumniService(repo *psql.AlumniRepository) *AlumniService {
	return &AlumniService{repo: repo}
}

// GetAllAlumniByFilter - GET /alumni-management/alumni
func (s *AlumniService) GetAllAlumniByFilter(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "desc")
	search := c.Query("search", "")

	offset := (page - 1) * limit
	sortByWhitelist := map[string]bool{
		"nama":        true,
		"jurusan":     true,
		"angkatan":    true,
		"tahun_lulus": true,
	}

	if !sortByWhitelist[sortBy] {
		return helper.ErrorResponse(c, 400, "kolom sortBy tidak valid")
	}

	order = strings.ToLower(order)
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	alumni, err := s.repo.GetByFilter(search, sortBy, order, limit, offset)
	if err != nil {
		return helper.ErrorResponse(c, 500, "gagal mengambil data alumni dari repository")
	}

	total, err := s.repo.CountAlumniRepo(search)
	if err != nil {
		return helper.ErrorResponse(c, 500, "gagal menghitung total data alumni")
	}

	response := models.AlumniResponse{
		Data: alumni,
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

	return helper.SuccessResponse(c, "Data alumni berhasil diambil", response)
}

// GetAlumniByID - GET /alumni-management/alumni/:id
func (s *AlumniService) GetAlumniByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	alumni, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 404, "alumni tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "gagal mengambil data alumni")
	}

	return helper.SuccessResponse(c, "Data alumni berhasil diambil", alumni)
}


// CreateAlumni - POST /alumni-management/alumni
func (s *AlumniService) CreateAlumni(c *fiber.Ctx) error {
	var req models.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Request body tidak valid")
	}

	if strings.TrimSpace(req.NIM) == "" ||
		strings.TrimSpace(req.Nama) == "" ||
		strings.TrimSpace(req.Jurusan) == "" ||
		strings.TrimSpace(req.Email) == "" {
		return helper.ErrorResponse(c, 400, "semua field wajib diisi")
	}

	if req.Angkatan <= 0 || req.TahunLulus <= 0 {
		return helper.ErrorResponse(c, 400, "angkatan dan tahun lulus harus valid")
	}

	if req.TahunLulus < req.Angkatan {
		return helper.ErrorResponse(c, 400, "tahun lulus tidak boleh lebih kecil dari angkatan")
	}

	alumni, err := s.repo.Create(&req)
	if err != nil {
		return helper.ErrorResponse(c, 500, "gagal menambahkan alumni")
	}

	return helper.CreatedResponse(c, "Alumni berhasil ditambahkan", alumni)
}

// UpdateAlumni - PUT /alumni-management/alumni/:id
func (s *AlumniService) UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	var req models.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Request body tidak valid")
	}

	if strings.TrimSpace(req.Nama) == "" ||
		strings.TrimSpace(req.Jurusan) == "" ||
		strings.TrimSpace(req.Email) == "" {
		return helper.ErrorResponse(c, 400, "semua field wajib diisi")
	}

	if req.Angkatan <= 0 || req.TahunLulus <= 0 {
		return helper.ErrorResponse(c, 400, "angkatan dan tahun lulus harus valid")
	}

	if req.TahunLulus < req.Angkatan {
		return helper.ErrorResponse(c, 400, "tahun lulus tidak boleh lebih kecil dari angkatan")
	}

	alumni, err := s.repo.Update(id, &req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 404, "alumni tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "gagal memperbarui alumni")
	}

	return helper.SuccessResponse(c, "Alumni berhasil diupdate", alumni)
}


// DeleteAlumni - DELETE /alumni-management/alumni/:id
func (s *AlumniService) DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	err = s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 404, "alumni tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "gagal menghapus data alumni")
	}

	return helper.SuccessResponse(c, "Alumni berhasil dihapus", nil)
}
