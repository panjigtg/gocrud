package config

import (
	"database/sql"
	"crudprojectgo/app/repository"
	"crudprojectgo/app/services"
)

type RepositoryContainer struct {
	Mahasiswa *repository.MahasiswaRepository
	Alumni    *repository.AlumniRepository
	Pekerjaan *repository.PekerjaanRepository
	Auth      *repository.AuthRepository
	Users     *repository.UsersRepository
}

type ServiceContainer struct {
	Mahasiswa *services.MahasiswaService
	Alumni    *services.AlumniService
	Pekerjaan *services.PekerjaanService
	Auth      *services.AuthServices
	Users     *services.UsersService
}

func InitRepositories(db *sql.DB) *RepositoryContainer {
	return &RepositoryContainer{
		Mahasiswa: repository.NewMahasiswaRepository(db),
		Alumni:    repository.NewAlumniRepository(db),
		Pekerjaan: repository.NewPekerjaanRepository(db),
		Auth:      repository.NewAuthRepository(db),
		Users:     repository.NewUsersRepository(db),
	}
}

func InitServices(r *RepositoryContainer) *ServiceContainer {
	return &ServiceContainer{
		Mahasiswa: services.NewMahasiswaService(r.Mahasiswa),
		Alumni:    services.NewAlumniService(r.Alumni),
		Pekerjaan: services.NewPekerjaanService(r.Pekerjaan, r.Alumni),
		Auth:      services.NewAuthServices(r.Auth),
		Users:     services.NewUsersService(r.Users),
	}
}
