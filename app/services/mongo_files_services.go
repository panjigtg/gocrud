package services

import (
	"crudprojectgo/app/models"
	mongoRepo "crudprojectgo/app/repository/mongo"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileService interface {
	UploadFoto(c *fiber.Ctx) error
	UploadSertifikat(c *fiber.Ctx) error
	GetAllFiles(c *fiber.Ctx) error
	GetFileByID(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type FileServiceMongo struct {
	Repo *mongoRepo.FileRepo
}

func NewFileServiceMongo(r *mongoRepo.FileRepo) *FileServiceMongo {
	return &FileServiceMongo{
		Repo: r,
	}
}

func (s *FileServiceMongo) UploadFoto(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("user_id"))
	role := fmt.Sprintf("%v", c.Locals("role"))
	paramUserID := c.Params("user_id")

	// Validasi role: user hanya bisa upload miliknya sendiri
	if role != "admin" && userID != paramUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Anda hanya dapat mengunggah foto untuk akun Anda sendiri",
		})
	}

	file, err := c.FormFile("foto")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "File foto tidak ditemukan"})
	}

	// Validasi ukuran
	if file.Size > 1*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": "Ukuran file maksimal 1MB"})
	}

	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return c.Status(400).JSON(fiber.Map{"error": "Format foto harus jpg, jpeg, atau png"})
	}

	// Simpan file
	saveDir := fmt.Sprintf("./uploads/foto/%s", paramUserID)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat direktori penyimpanan"})
	}

	savePath := filepath.Join(saveDir, file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan file"})
	}

	// Simpan metadata di Mongo
	record := models.File{
		ID:           primitive.NewObjectID(),
		FileName:     file.Filename,
		OriginalName: file.Filename,
		FilePath:     savePath,
		FileSize:     file.Size,
		FileType:     "foto",
		UploadedAt:   time.Now(),
		UpdatedAt:    time.Now(),
		Version:      1,
	}

	if err := s.Repo.Create(&record); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Gagal menyimpan metadata: %v", err)})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Foto berhasil diunggah",
		"data":    record,
	})
}


func (s *FileServiceMongo) UploadSertifikat(c *fiber.Ctx) error {
	userID := fmt.Sprintf("%v", c.Locals("user_id"))
	role := fmt.Sprintf("%v", c.Locals("role"))
	paramUserID := c.Params("user_id")

	if role != "admin" && userID != paramUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Anda hanya dapat mengunggah sertifikat untuk akun Anda sendiri",
		})
	}

	file, err := c.FormFile("sertifikat")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "File sertifikat tidak ditemukan"})
	}

	if file.Size > 2*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": "Ukuran file maksimal 2MB"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		return c.Status(400).JSON(fiber.Map{"error": "Format file harus PDF"})
	}

	saveDir := fmt.Sprintf("./uploads/sertifikat/%s", paramUserID)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat direktori penyimpanan"})
	}

	savePath := filepath.Join(saveDir, file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan file"})
	}

	record := models.File{
		ID:           primitive.NewObjectID(),
		FileName:     file.Filename,
		OriginalName: file.Filename,
		FilePath:     savePath,
		FileSize:     file.Size,
		FileType:     "sertifikat",
		UploadedAt:   time.Now(),
		UpdatedAt:    time.Now(),
		Version:      1,
	}

	if err := s.Repo.Create(&record); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Gagal menyimpan metadata: %v", err)})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Sertifikat berhasil diunggah",
		"data":    record,
	})
}


func (s *FileServiceMongo) GetAllFiles(c *fiber.Ctx) error {
	files, err := s.Repo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(files)
}


func (s *FileServiceMongo) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := s.Repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(file)
}


func (s *FileServiceMongo) DeleteFile(c *fiber.Ctx) error {
	role := fmt.Sprintf("%v", c.Locals("role"))
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Hanya admin yang dapat menghapus file",
		})
	}

	id := c.Params("id")
	file, err := s.Repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	if err := os.Remove(file.FilePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus file dari sistem"})
	}

	if err := s.Repo.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus metadata dari database"})
	}

	return c.JSON(fiber.Map{"message": "File berhasil dihapus"})
}
