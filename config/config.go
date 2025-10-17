package config

import (
	"log"

	"crudprojectgo/app/repository"
	"crudprojectgo/app/services"
	"crudprojectgo/database"
	"crudprojectgo/middleware"
	"crudprojectgo/routes"

	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// RunApp menjalankan seluruh setup dan server
func RunApp() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal load .env file")
	}

	db := database.KoneksiDB()
	defer db.Close()

	repos := initRepositories(db)
	svcs := initServices(repos)

	app := setupFiber()
	registerRoutes(app, svcs)
	startServer(app)
}

func initRepositories(db *sql.DB) *repositoryContainer {
	return &repositoryContainer{
		Mahasiswa: repository.NewMahasiswaRepository(db),
		Alumni:    repository.NewAlumniRepository(db),
		Pekerjaan: repository.NewPekerjaanRepository(db),
		Chall:     repository.NewChallRepository(db),
		Auth:      repository.NewAuthRepository(db),
		Users:     repository.NewUsersRepository(db),
	}
}


func initServices(r *repositoryContainer) *serviceContainer {
	return &serviceContainer{
		Mahasiswa: services.NewMahasiswaService(r.Mahasiswa),
		Alumni:    services.NewAlumniService(r.Alumni),
		Pekerjaan: services.NewPekerjaanService(r.Pekerjaan, r.Alumni),
		Chall:     services.NewChallServices(r.Chall),
		Auth:      services.NewAuthServices(r.Auth),
		Users:     services.NewUsersService(r.Users),
	}
}

func setupFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(middleware.RequestLogger())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Alumni Management API is running",
			"version": "1.0.0",
		})
	})

	return app
}

func registerRoutes(app *fiber.App, s *serviceContainer) {
	routes.NewMahasiswaHandler(s.Mahasiswa).SetupRoutes(app)
	routes.NewAlumniHandler(s.Alumni).SetupRoutes(app)
	routes.NewPekerjaanHandler(s.Pekerjaan).SetupRoutes(app)
	routes.NewChallHandler(s.Chall).SetupRoutes(app)
	routes.NewAuthHandler(s.Auth).SetupRoutes(app)
	routes.NewUsersHandler(s.Users).SetupRoutes(app)
}

func startServer(app *fiber.App) {
	const addr = "127.0.0.1:3000"
	log.Printf("Server berjalan di http://%s", addr)
	log.Fatal(app.Listen(addr))
}


type repositoryContainer struct {
	Mahasiswa *repository.MahasiswaRepository
	Alumni    *repository.AlumniRepository
	Pekerjaan *repository.PekerjaanRepository
	Chall     *repository.ChallRepository
	Auth      *repository.AuthRepository
	Users     *repository.UsersRepository
}

type serviceContainer struct {
	Mahasiswa *services.MahasiswaService
	Alumni    *services.AlumniService
	Pekerjaan *services.PekerjaanService
	Chall     *services.ChallServices
	Auth      *services.AuthServices
	Users     *services.UsersService
}
