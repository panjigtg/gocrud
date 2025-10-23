package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

)

type MongoAlumni struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OriginalID  int                `bson:"original_id" json:"original_id"`
	NIM         string             `bson:"nim" json:"nim"`
	Nama        string             `bson:"nama" json:"nama"`
	Jurusan     string             `bson:"jurusan" json:"jurusan"`
	Angkatan    int                `bson:"angkatan" json:"angkatan"`
	TahunLulus  int                `bson:"tahun_lulus" json:"tahun_lulus"`
	Email       string             `bson:"email" json:"email"`
	NoTelepon   *string            `bson:"no_telepon,omitempty" json:"no_telepon,omitempty"`
	Alamat      *string            `bson:"alamat,omitempty" json:"alamat,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	UserID      int                `bson:"user_id" json:"user_id"`
}