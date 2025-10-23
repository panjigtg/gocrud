package config

import (
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"

	
	"crudprojectgo/app/repository/psql"
	"crudprojectgo/app/services"
)

type RepositoryContainer struct {
	Mahasiswa *psql.MahasiswaRepository
	Alumni    *psql.AlumniRepository
	Pekerjaan *psql.PekerjaanRepository
	Auth      *psql.AuthRepository
	Users     *psql.UsersRepository
}

type ServiceContainer struct {
	Mahasiswa *services.MahasiswaService
	Alumni    *services.AlumniService
	Pekerjaan *services.PekerjaanService
	Auth      *services.AuthServices
	Users     *services.UsersService
}

func InitRepositories(db *sql.DB,  mongo *mongo.Database) *RepositoryContainer {
	return &RepositoryContainer{
		Mahasiswa: psql.NewMahasiswaRepository(db),
		Alumni:    psql.NewAlumniRepository(db),
		Pekerjaan: psql.NewPekerjaanRepository(db),
		Auth:      psql.NewAuthRepository(db),
		Users:     psql.NewUsersRepository(db),
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
