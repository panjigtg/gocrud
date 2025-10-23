package config

import (
	"database/sql"

	psqlRepo "crudprojectgo/app/repository/psql"
	mongoRepo "crudprojectgo/app/repository/mongo"
	"crudprojectgo/app/services"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type RepositoryContainer struct {
	Mahasiswa      *psqlRepo.MahasiswaRepository
	Alumni         *psqlRepo.AlumniRepository
	Pekerjaan      *psqlRepo.PekerjaanRepository
	PekerjaanMongo *mongoRepo.PekerjaanRepo // ✅ dari folder kamu, bukan driver
	Auth           *psqlRepo.AuthRepository
	Users          *psqlRepo.UsersRepository
}

type ServiceContainer struct {
	Mahasiswa      *services.MahasiswaService
	Alumni         *services.AlumniService
	Pekerjaan      *services.PekerjaanService      // PostgreSQL version
	PekerjaanMongo *services.PekerjaanServiceMongo // MongoDB version
	Auth           *services.AuthServices
	Users          *services.UsersService
}

// ✅ sekarang parameternya pakai mongoDriver.Database
func InitRepositories(db *sql.DB, mongo *mongoDriver.Database) *RepositoryContainer {
	return &RepositoryContainer{
		Mahasiswa:      psqlRepo.NewMahasiswaRepository(db),
		Alumni:         psqlRepo.NewAlumniRepository(db),
		Pekerjaan:      psqlRepo.NewPekerjaanRepository(db),
		PekerjaanMongo: mongoRepo.NewPekerjaanRepo(mongo), 
		Auth:           psqlRepo.NewAuthRepository(db),
		Users:          psqlRepo.NewUsersRepository(db),
	}
}

func InitServices(r *RepositoryContainer) *ServiceContainer {
	return &ServiceContainer{
		Mahasiswa:      services.NewMahasiswaService(r.Mahasiswa),
		Alumni:         services.NewAlumniService(r.Alumni),
		Pekerjaan:      services.NewPekerjaanService(r.Pekerjaan, r.Alumni),
		PekerjaanMongo: services.NewPekerjaanServiceMongo(r.PekerjaanMongo),
		Auth:           services.NewAuthServices(r.Auth),
		Users:          services.NewUsersService(r.Users),
	}
}
