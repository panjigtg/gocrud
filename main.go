package main

import (
    "crudprojectgo/helper"
    "crudprojectgo/app/repository"
    "crudprojectgo/routes"
    "crudprojectgo/app/services"
    "crudprojectgo/middleware"
    "log"
    
    "github.com/joho/godotenv"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Gagal load .env file")
    }

    // Koneksi database
    db := helper.KoneksiDB()
    defer db.Close()
    
    // Inisialisasi repository
    mahasiswaRepo := repository.NewMahasiswaRepository(db)
    alumniRepo := repository.NewAlumniRepository(db)
    pekerjaanRepo := repository.NewPekerjaanRepository(db)
    challRepo := repository.NewChallRepository(db)
    authRepo := repository.NewAuthRepository(db)
    usersRepo := repository.NewUsersRepository(db)
    
    // Inisialisasi service
    mahasiswaService := services.NewMahasiswaService(mahasiswaRepo)
    alumniService := services.NewAlumniService(alumniRepo)
    pekerjaanService := services.NewPekerjaanService(pekerjaanRepo, alumniRepo)
    challService := services.NewChallServices(challRepo)
    authService := services.NewAuthServices(authRepo)
    usersService := services.NewUsersService(usersRepo)

    // Inisialisasi handler
    mahasiswaHandler := routes.NewMahasiswaHandler(mahasiswaService)
    alumniHandler := routes.NewAlumniHandler(alumniService)
    pekerjaanHandler := routes.NewPekerjaanHandler(pekerjaanService)
    challHandler := routes.NewChallHandler(challService)
    authHandler := routes.NewAuthHandler(authService)
	usersHandler := routes.NewUsersHandler(usersService)

    
    // Inisialisasi Fiber
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(500).JSON(fiber.Map{
                "success": false,
                "error":   err.Error(),
            })
        },
    })
    
    // Middleware
    app.Use(logger.New())
    app.Use(cors.New())
    app.Use(middleware.RequestLogger())
    
    // Health check endpoint
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Alumni Management API is running!",
            "version": "1.0.0",
        })
    })
    
    // Setup routes
    mahasiswaHandler.SetupRoutes(app)
    alumniHandler.SetupRoutes(app)
    pekerjaanHandler.SetupRoutes(app)
    challHandler.SetupRoutes(app)
    authHandler.SetupRoutes(app)
    usersHandler.SetupRoutes(app)

    
    // Start server
    log.Println("Server starting on port 3000...")
    log.Fatal(app.Listen(":3000"))
}
