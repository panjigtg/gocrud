package psql

import (
    "crudprojectgo/app/models"
    "database/sql"
    "time"
)

type MahasiswaRepository struct {
	db *sql.DB
}

func NewMahasiswaRepository(db *sql.DB) *MahasiswaRepository {
	return &MahasiswaRepository{
		db: db,
	}
}

func (r *MahasiswaRepository) GetAll() ([]models.Mahasiswa, error) {
    query := `
        SELECT id, nim, nama, jurusan, angkatan, email, created_at, updated_at 
        FROM mahasiswa
        WHERE is_deleted IS NULL
        ORDER BY created_at DESC
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var mahasiswas []models.Mahasiswa
    for rows.Next() {
        var m models.Mahasiswa
        err := rows.Scan(
            &m.ID, &m.NIM, &m.Nama, &m.Jurusan, &m.Angkatan,
            &m.Email, &m.CreatedAt, &m.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        mahasiswas = append(mahasiswas, m)
    }

    return mahasiswas, nil
}

func (r *MahasiswaRepository) GetByID(id int) (*models.Mahasiswa, error) {
    query := `
        SELECT id, nim, nama, jurusan, angkatan, email, created_at, updated_at
        FROM mahasiswa WHERE id = $1 AND is_deleted IS NULL
    `
    var m models.Mahasiswa
    err := r.db.QueryRow(query, id).Scan(
        &m.ID, &m.NIM, &m.Nama, &m.Jurusan, &m.Angkatan,
        &m.Email, &m.CreatedAt, &m.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &m, nil
}

func (r *MahasiswaRepository) Create(req *models.CreateMahasiswaRequest) (*models.Mahasiswa, error) {
    query := `
        INSERT INTO mahasiswa (nim, nama, jurusan, angkatan, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    now := time.Now()
    result, err := r.db.Exec(query,
        req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.Email, now, now,
    )
    if err != nil {
        return nil, err
    }

    id, _ := result.LastInsertId()
    return r.GetByID(int(id))
}

func (r *MahasiswaRepository) Update(id int, req *models.UpdateMahasiswaRequest) (*models.Mahasiswa, error) {
    query := `
        UPDATE mahasiswa 
        SET nama=$1, jurusan=$2, angkatan=$3, email=$4, updated_at=$5
        WHERE id=$6
    `
    now := time.Now()
    _, err := r.db.Exec(query, req.Nama, req.Jurusan, req.Angkatan, req.Email, now, id)
    if err != nil {
        return nil, err
    }

    return r.GetByID(id)
}

func (r *MahasiswaRepository) Delete(id int) error {
    query := "DELETE FROM mahasiswa WHERE id=$1"
    _, err := r.db.Exec(query, id)
    return err
}

func (r *MahasiswaRepository) SoftDeletes(id int, req *models.IsDeleted) (*models.Mahasiswa, error) {
    query := "UPDATE mahasiswa SET is_deleted = $2 WHERE id=$1"
    now := time.Now()
    _, err := r.db.Exec(query, id, now)
    return nil, err
}

