package mongo

import (
	"context"
	"crudprojectgo/app/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PekerjaanRepo struct {
	Col *mongo.Collection
}

func NewPekerjaanRepo(db *mongo.Database) *PekerjaanRepo {
	return &PekerjaanRepo{
		Col: db.Collection("pekerjaan_alumni"),
	}
}

func (r *PekerjaanRepo) GetAll() ([]models.PekerjaanAlumniMongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.Col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.PekerjaanAlumniMongo
	for cursor.Next(ctx) {
		var p models.PekerjaanAlumniMongo
		cursor.Decode(&p)
		results = append(results, p)
	}
	return results, nil
}

func (r *PekerjaanRepo) GetByID(id string) (models.PekerjaanAlumniMongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	var result models.PekerjaanAlumniMongo
	err := r.Col.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	return result, err
}

func (r *PekerjaanRepo) Create(p *models.PekerjaanAlumniMongo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Col.InsertOne(ctx, p)
	return err
}

func (r *PekerjaanRepo) Update(id string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := r.Col.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	return err
}

func (r *PekerjaanRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := r.Col.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}