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

type PekerjaanService struct {
	repo       *psql.PekerjaanRepository
	alumniRepo *psql.AlumniRepository
}

func NewPekerjaanService(repo *psql.PekerjaanRepository, alumniRepo *psql.AlumniRepository) *PekerjaanService {
	return &PekerjaanService{repo: repo, alumniRepo: alumniRepo}
}

func (s *PekerjaanService) GetAllPekerjaan(c *fiber.Ctx) error {
	data, err := s.repo.GetAll()
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengambil data pekerjaan")
	}
	return helper.SuccessResponse(c, "Data pekerjaan berhasil diambil", data)
}

func (s *PekerjaanService) GetPekerjaanByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}
	data, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 404, "Pekerjaan tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "Kesalahan server")
	}
	return helper.SuccessResponse(c, "Data pekerjaan berhasil diambil", data)
}

func (s *PekerjaanService) GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "Alumni ID tidak valid")
	}

	// pastikan alumni ada
	_, err = s.alumniRepo.GetByID(alumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.ErrorResponse(c, 404, "Alumni tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "Kesalahan saat mengambil alumni")
	}

	data, err := s.repo.GetByAlumniID(alumniID)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengambil data pekerjaan alumni")
	}
	if len(data) == 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Data pekerjaan alumni tidak ditemukan",
			"data":    []models.PekerjaanAlumni{},
		})
	}
	return helper.SuccessResponse(c, "Data pekerjaan alumni berhasil diambil", data)
}

func (s *PekerjaanService) CreatePekerjaan(c *fiber.Ctx) error {
	var req models.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Request body tidak valid")
	}

	if req.AlumniID <= 0 {
		return helper.ErrorResponse(c, 400, "Alumni ID harus valid")
	}
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" {
		return helper.ErrorResponse(c, 400, "Nama perusahaan, posisi jabatan, bidang industri, dan lokasi kerja wajib diisi")
	}

	validStatus := []string{"aktif", "selesai", "resigned"}
	isValid := false
	for _, v := range validStatus {
		if strings.EqualFold(req.StatusPekerjaan, v) {
			isValid = true
			break
		}
	}
	if !isValid {
		return helper.ErrorResponse(c, 400, "Status pekerjaan harus salah satu dari: aktif, selesai, resigned")
	}

	// cek alumni
	_, err := s.alumniRepo.GetByID(req.AlumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.ErrorResponse(c, 404, "Alumni tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "Kesalahan saat validasi alumni")
	}

	data, err := s.repo.Create(&req)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal menambahkan pekerjaan")
	}
	return helper.CreatedResponse(c, "Pekerjaan berhasil ditambahkan", data)
}

func (s *PekerjaanService) UpdatePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	var req models.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Request body tidak valid")
	}

	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" {
		return helper.ErrorResponse(c, 400, "Nama perusahaan, posisi jabatan, bidang industri wajib diisi")
	}

	data, err := s.repo.Update(id, &req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 404, "Pekerjaan tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "Gagal memperbarui pekerjaan")
	}
	return helper.SuccessResponse(c, "Pekerjaan berhasil diperbarui", data)
}

func (s *PekerjaanService) DeletePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}
	err = s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, 404, "Pekerjaan tidak ditemukan")
		}
		return helper.ErrorResponse(c, 500, "Gagal menghapus pekerjaan")
	}
	return helper.SuccessResponse(c, "Pekerjaan berhasil dihapus permanen", nil)
}


func (s *PekerjaanService) GetPekerjaanByFilter(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	sortMap := map[string]string{
		"nama alumni":      "a.nama",
		"nama perusahaan":  "p.nama_perusahaan",
		"posisi jabatan":   "p.posisi_jabatan",
		"bidang industri":  "p.bidang_industri",
		"lokasi kerja":     "p.lokasi_kerja",
		"gaji range":       "p.gaji_range",
		"status pekerjaan": "p.status_pekerjaan",
	}
	col, ok := sortMap[strings.ToLower(sortBy)]
	if !ok {
		return helper.ErrorResponse(c, 400, "Kolom sortBy tidak valid")
	}

	data, err := s.repo.GetByFilter(search, col, order, page, limit, (page-1)*limit)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengambil data pekerjaan")
	}

	total, _ := s.repo.CountPekerjaanRepo(search)
	meta := models.MetaInfo{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: (total + limit - 1) / limit,
	}
	return helper.SuccessResponse(c, "Data pekerjaan berhasil diambil", fiber.Map{"meta": meta, "data": data})
}


func (s *PekerjaanService) SoftDeletePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	err = s.repo.SoftDeletes(id)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal melakukan soft delete")
	}

	return helper.SuccessResponse(c, "Pekerjaan berhasil dihapus (soft delete)", nil)
}

func (s *PekerjaanService) GetTrash(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	data, err := s.repo.GetTrash(role, userID)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengambil data trash")
	}

	if len(data) == 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Tidak ada data di trash",
			"data":    []models.Trash{},
		})
	}

	return helper.SuccessResponse(c, "Data trash berhasil diambil", data)
}

func (s *PekerjaanService) RestorePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	err = s.repo.RestoreByID(id)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal merestore pekerjaan")
	}

	return helper.SuccessResponse(c, "Pekerjaan berhasil direstore", nil)
}

func (s *PekerjaanService) HardDeletePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	isDeleted, err := s.repo.IsInTrash(id)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal memeriksa status trash")
	}
	if !isDeleted {
		return helper.ErrorResponse(c, 400, "Data belum dihapus (soft delete) sehingga tidak bisa hard delete")
	}

	err = s.repo.DeletePermanent(id)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal menghapus data permanen")
	}

	return helper.SuccessResponse(c, "Pekerjaan berhasil dihapus permanen", nil)
}
