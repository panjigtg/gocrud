package services

import (
	"crudprojectgo/app/models"
	"crudprojectgo/app/repository"
	"crudprojectgo/helper"
	"database/sql"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MahasiswaService struct {
	repo *repository.MahasiswaRepository
}

func NewMahasiswaService(repo *repository.MahasiswaRepository) *MahasiswaService {
	return &MahasiswaService{repo: repo}
}

func (s *MahasiswaService) GetAllMahasiswa(c *fiber.Ctx) error {
	data, err := s.repo.GetAll()
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengambil data mahasiswa")
	}
	return helper.SuccessResponse(c, "Data mahasiswa berhasil diambil", data)
}


func (s *MahasiswaService) GetMahasiswaByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	mhs, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 404, "Mahasiswa tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "Kesalahan server")
	}
	return helper.SuccessResponse(c, "Data mahasiswa berhasil diambil", mhs)
}


func (s *MahasiswaService) CreateMahasiswa(c *fiber.Ctx) error {
	var req models.CreateMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Request body tidak valid")
	}

	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return helper.ErrorResponse(c, 400, "NIM, nama, jurusan, dan email wajib diisi")
	}

	data, err := s.repo.Create(&req)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal menambahkan mahasiswa")
	}
	return helper.CreatedResponse(c, "Mahasiswa berhasil ditambahkan", data)
}

func (s *MahasiswaService) UpdateMahasiswa(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	var req models.UpdateMahasiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Request body tidak valid")
	}

	if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return helper.ErrorResponse(c, 400, "Nama, jurusan, dan email wajib diisi")
	}

	data, err := s.repo.Update(id, &req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 404, "Mahasiswa tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "Gagal mengupdate data mahasiswa")
	}
	return helper.SuccessResponse(c, "Mahasiswa berhasil diupdate", data)
}

func (s *MahasiswaService) DeleteMahasiswa(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	err = s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 404, "Mahasiswa tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "Gagal menghapus mahasiswa")
	}
	return helper.SuccessResponse(c, "Mahasiswa berhasil dihapus", nil)
}

func (s *MahasiswaService) SoftDeletes(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	var req models.IsDeleted
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Request body tidak valid")
	}

	_, err = s.repo.SoftDeletes(id, &req)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal melakukan soft delete")
	}

	return helper.SuccessResponse(c, "Mahasiswa berhasil dihapus secara soft delete", nil)
}
