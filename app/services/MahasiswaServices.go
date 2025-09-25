package services

import (
    "crudprojectgo/app/models"
    "crudprojectgo/app/repository"
    "database/sql"
    "errors"
)

type MahasiswaService struct {
    repo *repository.MahasiswaRepository
}

func NewMahasiswaService(repo *repository.MahasiswaRepository) *MahasiswaService {
    return &MahasiswaService{repo: repo}
}

func (s *MahasiswaService) GetAllMahasiswa() ([]models.Mahasiswa, error) {
    return s.repo.GetAll()
}

func (s *MahasiswaService) GetMahasiswaByID(id int) (*models.Mahasiswa, error) {
    mhs, err := s.repo.GetByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("mahasiswa tidak ditemukan")
        }
        return nil, err
    }
    return mhs, nil
}

func (s *MahasiswaService) CreateMahasiswa(req *models.CreateMahasiswaRequest) (*models.Mahasiswa, error) {
    if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
        return nil, errors.New("NIM, nama, jurusan, dan email wajib diisi")
    }
    return s.repo.Create(req)
}

func (s *MahasiswaService) UpdateMahasiswa(id int, req *models.UpdateMahasiswaRequest) (*models.Mahasiswa, error) {
    if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
        return nil, errors.New("nama, jurusan, dan email wajib diisi")
    }
    return s.repo.Update(id, req)
}

func (s *MahasiswaService) DeleteMahasiswa(id int) error {
    return s.repo.Delete(id)
}