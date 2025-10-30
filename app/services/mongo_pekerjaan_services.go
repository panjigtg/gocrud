package services

import (
	"context"
	"crudprojectgo/app/models"
	mongoRepo "crudprojectgo/app/repository/mongo"
	"crudprojectgo/helper"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struct utama service Mongo
type PekerjaanServiceMongo struct {
	Repo *mongoRepo.PekerjaanRepo
}

// Konstruktor agar bisa dipanggil di container.go
func NewPekerjaanServiceMongo(repo *mongoRepo.PekerjaanRepo) *PekerjaanServiceMongo {
	return &PekerjaanServiceMongo{
		Repo: repo,
	}
}


func (s *PekerjaanServiceMongo) GetAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.Repo.Col.Find(ctx, bson.M{})
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengambil data")
	}
	defer cursor.Close(ctx)

	var data []models.PekerjaanAlumniMongo
	for cursor.Next(ctx) {
		var p models.PekerjaanAlumniMongo
		cursor.Decode(&p)
		data = append(data, p)
	}

	return helper.SuccessResponse(c, "Data berhasil diambil", data)
}

func (s *PekerjaanServiceMongo) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	var data models.PekerjaanAlumniMongo
	err = s.Repo.Col.FindOne(context.TODO(), bson.M{"original_id": id}).Decode(&data)
	if err != nil {
		return helper.ErrorResponse(c, 404, "Data tidak ditemukan")
	}

	return helper.SuccessResponse(c, "Data berhasil diambil", data)
}

func (s *PekerjaanServiceMongo) GetByAlumniID(c *fiber.Ctx) error {
	alumniID, _ := c.ParamsInt("alumni_id")

	cursor, err := s.Repo.Col.Find(context.TODO(), bson.M{"alumni_id": alumniID})
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal mengambil data alumni")
	}
	defer cursor.Close(context.TODO())

	var data []models.PekerjaanAlumniMongo
	for cursor.Next(context.TODO()) {
		var p models.PekerjaanAlumniMongo
		cursor.Decode(&p)
		data = append(data, p)
	}

	return helper.SuccessResponse(c, "Data pekerjaan alumni berhasil diambil", data)
}

func (s *PekerjaanServiceMongo) Create(c *fiber.Ctx) error {
	var req models.PekerjaanAlumniMongo
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Data tidak valid")
	}

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	_, err := s.Repo.Col.InsertOne(context.TODO(), req)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal menambah data")
	}

	return helper.CreatedResponse(c, "Data berhasil ditambahkan", req)
}

func (s *PekerjaanServiceMongo) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	var req models.PekerjaanAlumniMongo
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, 400, "Data tidak valid")
	}

	update := bson.M{
		"$set": bson.M{
			"nama_perusahaan":   req.NamaPerusahaan,
			"posisi_jabatan":    req.PosisiJabatan,
			"bidang_industri":   req.BidangIndustri,
			"lokasi_kerja":      req.LokasiKerja,
			"status_pekerjaan":  req.StatusPekerjaan,
			"gaji_range":        req.GajiRange,
			"deskripsi_pekerjaan": req.DeskripsiPekerjaan,
			"updated_at":        time.Now(),
		},
	}

	result, err := s.Repo.Col.UpdateOne(
		context.TODO(),
		bson.M{"original_id": id},
		update,
	)
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal memperbarui data")
	}

	if result.MatchedCount == 0 {
		return helper.ErrorResponse(c, 404, "Data tidak ditemukan")
	}

	return helper.SuccessResponse(c, "Data berhasil diperbarui", nil)
}

func (s *PekerjaanServiceMongo) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return helper.ErrorResponse(c, 400, "ID tidak valid")
	}

	result, err := s.Repo.Col.DeleteOne(context.TODO(), bson.M{"original_id": id})
	if err != nil {
		return helper.ErrorResponse(c, 500, "Gagal menghapus data")
	}

	if result.DeletedCount == 0 {
		return helper.ErrorResponse(c, 404, "Data tidak ditemukan")
	}

	return helper.SuccessResponse(c, "Data berhasil dihapus", nil)
}
