package mongo

import (
	"crudprojectgo/app/models"
	"context"
	"errors"
	"time"
	"sync"
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)	

type FilesUploads interface {
	Create(file *models.File) error
	FindAll() ([]models.File, error)
	FindByID(id string) (*models.File, error)
	Update(id string, update models.FileUpdate) error
	Delete(id string) error
}

type FileRepo struct {
	 Col 		*mongo.Collection
	 validate   *validator.Validate
	 mu   		 sync.Mutex
}	

var _ FilesUploads = (*FileRepo)(nil)


func NewFileRepository(db *mongo.Database) *FileRepo {
	return &FileRepo{
		Col: db.Collection("files_uploads"),
		validate:   validator.New(),
	}
}


func (r *FileRepo) Create(file *models.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.validate.Struct(file); err != nil {
		return fmt.Errorf("invalid file data: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	filter := bson.M{"file_name": file.FileName}
	count, err := r.Col.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("file with the same name already exists")
	}

	file.ID = primitive.NewObjectID()
	file.UploadedAt = time.Now()
	file.UpdatedAt = time.Now()

	_, err = r.Col.InsertOne(ctx, file)
	return err
}

func (r *FileRepo) FindAll() ([]models.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.Col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []models.File
	for cursor.Next(ctx) {
		var f models.File
		if err := cursor.Decode(&f); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (r *FileRepo) FindByID(id string) (*models.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !primitive.IsValidObjectID(id) {
		return nil, errors.New("invalid ID format")
	}
	objectID, _ := primitive.ObjectIDFromHex(id)

	var file models.File
	err := r.Col.FindOne(ctx, bson.M{"_id": objectID}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("file not found")
		}
		return nil, err
	}
	return &file, nil
}

func (r *FileRepo) Update(id string, update models.FileUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	updateFields := bson.M{}
	if update.Filename != "" {
		updateFields["filename"] = update.Filename
	}
	if update.Description != "" {
		updateFields["description"] = update.Description
	}
	if update.Size > 0 {
		updateFields["size"] = update.Size
	}

	if len(updateFields) == 0 {
		return errors.New("no valid fields to update")
	}

	updateFields["updated_at"] = time.Now()

	filter := bson.M{"_id": objectID, "version": update.Version}
	updateCmd := bson.M{
		"$set": updateFields,
		"$inc": bson.M{"version": 1},
	}

	res, err := r.Col.UpdateOne(ctx, filter, updateCmd)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("update conflict or file not found")
	}
	return nil
}

func (r *FileRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !primitive.IsValidObjectID(id) {
		return errors.New("invalid ID format")
	}
	objectID, _ := primitive.ObjectIDFromHex(id)

	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.Col.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("file not found")
	}
	return nil
}