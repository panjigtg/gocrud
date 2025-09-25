package services

import (
	"crudprojectgo/app/models"
	"crudprojectgo/app/repository"
	"database/sql"
	"errors"
    "strings"

	// "github.com/gofiber/fiber/v2/middleware/limiter"
)

type AlumniService struct {
    repo *repository.AlumniRepository
}


func NewAlumniService(repo *repository.AlumniRepository) *AlumniService {
    return &AlumniService{
        repo: repo,
    }
}

func (s *AlumniService) GetAllAlumni() ([]models.Alumni, error) {
    return s.repo.GetAll()
}

func (s *AlumniService) GetAlumniByID(id int) (*models.Alumni, error) {
    alumni, err := s.repo.GetByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("alumni tidak ditemukan")
        }
        return nil, err
    }
    return alumni, nil
}

func (s *AlumniService) CreateAlumni(req *models.CreateAlumniRequest) (*models.Alumni, error) {
    // Validasi input
    if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
        return nil, errors.New("NIM, nama, jurusan, dan email harus diisi")
    }
    
    if req.Angkatan <= 0 || req.TahunLulus <= 0 {
        return nil, errors.New("angkatan dan tahun lulus harus valid")
    }
    
    if req.TahunLulus < req.Angkatan {
        return nil, errors.New("tahun lulus tidak boleh lebih kecil dari angkatan")
    }
    
    return s.repo.Create(req)
}

func (s *AlumniService) UpdateAlumni(id int, req *models.UpdateAlumniRequest) (*models.Alumni, error) {
    // Validasi input
    if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
        return nil, errors.New("nama, jurusan, dan email harus diisi")
    }
    
    if req.Angkatan <= 0 || req.TahunLulus <= 0 {
        return nil, errors.New("angkatan dan tahun lulus harus valid")
    }
    
    if req.TahunLulus < req.Angkatan {
        return nil, errors.New("tahun lulus tidak boleh lebih kecil dari angkatan")
    }
    
    alumni, err := s.repo.Update(id, req)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("alumni tidak ditemukan")
        }
        return nil, err
    }
    
    return alumni, nil
}

func (s *AlumniService) DeleteAlumni(id int) error {
    err := s.repo.Delete(id)
    if err != nil {
        if err == sql.ErrNoRows {
            return errors.New("alumni tidak ditemukan")
        }
        return err
    }
    return nil
}

func (s *AlumniService) GetAllAlumniByFilter(page,limit int, sortBy, order, search string) (models.AlumniResponse, error) {
    offset := (page - 1) * limit

    sortByWhitelist := map[string]bool{
        "nama":       true,
        "jurusan":    true,
        "angkatan":   true,
        "tahun_lulus": true,
    }

    if !sortByWhitelist[sortBy] {
        return models.AlumniResponse{}, errors.New("kolom sortBy tidak valid")
    }
    
    if strings.ToLower(order) != "desc" {
        order = "asc"
    }

    alumni, err := s.repo.GetByFilter(search, sortBy, order, limit, offset)
    if err != nil {
        return models.AlumniResponse{}, err
    }

    total, err := s.repo.CountAlumniRepo(search)
    if err != nil {
        return models.AlumniResponse{}, err
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
    return response, nil
}