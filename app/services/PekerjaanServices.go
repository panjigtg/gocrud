package services

import (
	"crudprojectgo/app/models"
	"crudprojectgo/app/repository"
	"database/sql"
	"errors"
	"strings"
)

type PekerjaanService struct {
	repo       *repository.PekerjaanRepository
	alumniRepo *repository.AlumniRepository
}

func NewPekerjaanService(repo *repository.PekerjaanRepository, alumniRepo *repository.AlumniRepository) *PekerjaanService {
	return &PekerjaanService{
		repo:       repo,
		alumniRepo: alumniRepo,
	}
}

func (s *PekerjaanService) GetAllPekerjaan() ([]models.PekerjaanAlumni, error) {
	return s.repo.GetAll()
}

func (s *PekerjaanService) GetPekerjaanByID(id int) (*models.PekerjaanAlumni, error) {
	pekerjaan, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("pekerjaan tidak ditemukan")
		}
		return nil, err
	}
	return pekerjaan, nil
}

func (s *PekerjaanService) GetPekerjaanByAlumniID(alumniID int) ([]models.PekerjaanAlumni, error) {
	// Cek apakah alumni ada
	_, err := s.alumniRepo.GetByID(alumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("alumni tidak ditemukan")
		}
		return nil, err
	}

	pekerjaanList, err := s.repo.GetByAlumniID(alumniID)
	if err != nil {
		return nil, err
	}

	return pekerjaanList, nil
}

func (s *PekerjaanService) CreatePekerjaan(req *models.CreatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
	// Validasi input
	if req.AlumniID <= 0 {
		return nil, errors.New("alumni ID harus valid")
	}

	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" {
		return nil, errors.New("nama perusahaan, posisi jabatan, bidang industri, dan lokasi kerja harus diisi")
	}

	if req.TanggalMulaiKerja == "" {
		return nil, errors.New("tanggal mulai kerja harus diisi")
	}

	if req.StatusPekerjaan == "" {
		req.StatusPekerjaan = "aktif"
	}

	// Validasi status pekerjaan
	validStatus := []string{"aktif", "selesai", "resigned"}
	isValidStatus := false
	for _, status := range validStatus {
		if req.StatusPekerjaan == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return nil, errors.New("status pekerjaan harus salah satu dari: aktif, selesai, resigned")
	}

	// Cek apakah alumni ada
	_, err := s.alumniRepo.GetByID(req.AlumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("alumni tidak ditemukan")
		}
		return nil, err
	}

	return s.repo.Create(req)
}

func (s *PekerjaanService) UpdatePekerjaan(id int, req *models.UpdatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
	// Validasi input
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || req.BidangIndustri == "" || req.LokasiKerja == "" {
		return nil, errors.New("nama perusahaan, posisi jabatan, bidang industri, dan lokasi kerja harus diisi")
	}

	if req.TanggalMulaiKerja == "" {
		return nil, errors.New("tanggal mulai kerja harus diisi")
	}

	// Validasi status pekerjaan
	validStatus := []string{"aktif", "selesai", "resigned"}
	isValidStatus := false
	for _, status := range validStatus {
		if req.StatusPekerjaan == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return nil, errors.New("status pekerjaan harus salah satu dari: aktif, selesai, resigned")
	}

	pekerjaan, err := s.repo.Update(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("pekerjaan tidak ditemukan")
		}
		return nil, err
	}

	return pekerjaan, nil
}

func (s *PekerjaanService) DeletePekerjaan(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("pekerjaan tidak ditemukan")
		}
		return err
	}
	return nil
}

func (s *PekerjaanService) GetPekerjaanByFilter(search, sortBy, order string, page, limit int) (models.PekerjaanResponse, error) {
	offset := (page - 1) * limit

	SortByAlias := map[string]string{
		"nama alumni":      "a.nama",
		"nama perusahaan":  "p.nama_perusahaan",
		"posisi jabatan":   "p.posisi_jabatan",
		"bidang industri":  "p.bidang_industri",
		"lokasi kerja":     "p.lokasi_kerja",
		"gaji range":       "p.gaji_range",
		"status pekerjaan": "p.status_pekerjaan",
		"tahun mulai":      "COALESCE(EXTRACT(YEAR FROM p.tanggal_mulai_kerja), 0)",
		"tahun selesai":    "COALESCE(EXTRACT(YEAR FROM p.tanggal_selesai_kerja), 0)",
	}

	sortCol, ok := SortByAlias[strings.ToLower(sortBy)]
	if !ok {
		return models.PekerjaanResponse{}, errors.New("kolom sortBy tidak valid")
	}

	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	pekerjaan, err := s.repo.GetByFilter(search, sortCol, order, page, limit, offset)
	if err != nil {
		return models.PekerjaanResponse{}, err
	}

	total, err := s.repo.CountPekerjaanRepo(search)
	if err != nil {
		return models.PekerjaanResponse{}, err
	}

	response := models.PekerjaanResponse{
		Data: pekerjaan,
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

func (s *PekerjaanService) SoftDeletePekerjaan(pekerjaanID int, requesterID int, requesterRole string) error {
	// Ambil data pekerjaan
	pekerjaan, err := s.repo.GetByID(pekerjaanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data pekerjaan tidak ditemukan")
		}
		return err
	}

	// Ambil data alumni (validasi user_id)
	alumni, err := s.alumniRepo.GetByID(pekerjaan.AlumniID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("alumni terkait tidak ditemukan")
		}
		return err
	}

	// Validasi akses
	if requesterRole != "admin" && alumni.UserID != requesterID {
		return errors.New("anda tidak diizinkan menghapus pekerjaan milik orang lain")
	}

	// Jalankan soft delete
	return s.repo.SoftDeletes(pekerjaanID)
}


func (s *PekerjaanService) GetTrash(role string, userID int) ([]models.PekerjaanAlumni, error) {
	return s.repo.GetTrash(role, userID)
}


func (s *PekerjaanService) RestorePekerjaan(id int, requesterID int, requesterRole string) error {
	// Admin bisa langsung restore
	if requesterRole == "admin" {
		return s.repo.RestoreByID(id)
	}

	// Untuk user biasa → cek apakah data miliknya
	isOwner, err := s.repo.CheckOwnership(id, requesterID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("anda tidak memiliki akses untuk restore data ini")
	}

	return s.repo.RestoreByID(id)
}

func (s *PekerjaanService) HardDeletePekerjaan(id int, requesterID int, requesterRole string) error {
	if requesterRole == "admin" {
		return s.repo.DeletePermanent(id)
	}

	// Untuk user biasa, cek apakah miliknya sendiri
	isOwner, err := s.repo.CheckOwnership(id, requesterID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("anda tidak memiliki akses untuk menghapus data ini")
	}

	// Pastikan data memang sudah dihapus (is_deleted != NULL)
	isDeleted, err := s.repo.IsInTrash(id)
	if err != nil {
		return err
	}
	if !isDeleted {
		return errors.New("data belum dihapus (soft delete) sehingga tidak bisa hard delete")
	}

	return s.repo.DeletePermanent(id)
}
