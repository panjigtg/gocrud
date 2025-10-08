package models

import (
	"database/sql"
	// "time"
)

type Mahasiswa struct {
	ID int `json:"id"`
	NIM string `json:"nim"`
	Nama string `json:"nama"`
	Jurusan string `json:"jurusan"`
	Angkatan string `json:"angkatan"`
	Email string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IsDeleted sql.NullTime `json:"-"`
}

type CreateMahasiswaRequest struct {
	NIM string `json:"nim"`
	Nama string `json:"nama"`
	Jurusan string `json:"jurusan"`
	Angkatan string `json:"angkatan"`
	Email string `json:"email"`
}

type UpdateMahasiswaRequest struct {
	Nama string `json:"nama"`
	Jurusan string `json:"jurusan"`
	Angkatan string `json:"angkatan"`
	Email string `json:"email"`
}

type IsDeleted struct {
	ID 			int 		`json:"id"`
	IsDeleted 	sql.NullTime   `json:"-"`
}
