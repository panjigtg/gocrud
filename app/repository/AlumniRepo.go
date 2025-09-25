package repository

import (
    "crudprojectgo/app/models"
    "database/sql"
    "time"
    "strings"   
    "fmt"
    "log"
)

type AlumniRepository struct {
    db *sql.DB
}

func NewAlumniRepository(db *sql.DB) *AlumniRepository {
    return &AlumniRepository{
        db: db,
    }
}

func (r *AlumniRepository) GetAll() ([]models.Alumni, error) {
    query := `
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
               no_telepon, alamat, created_at, updated_at 
        FROM alumni 
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var alumni []models.Alumni
    for rows.Next() {
        var a models.Alumni
        err := rows.Scan(
            &a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
            &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
            &a.CreatedAt, &a.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        alumni = append(alumni, a)
    }
    
    return alumni, nil
}

func (r *AlumniRepository) GetByID(id int) (*models.Alumni, error) {
    query := `
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
               no_telepon, alamat, created_at, updated_at 
        FROM alumni 
        WHERE id = $1
    `
    
    var a models.Alumni
    err := r.db.QueryRow(query, id).Scan(
        &a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
        &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
        &a.CreatedAt, &a.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &a, nil
}

func (r *AlumniRepository) Create(req *models.CreateAlumniRequest) (*models.Alumni, error) {
    query := `
        INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, updated_at
    `
    
    now := time.Now()
    var a models.Alumni
    
    err := r.db.QueryRow(
        query, req.NIM, req.Nama, req.Jurusan, req.Angkatan,
        req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, now, now,
    ).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
    
    if err != nil {
        return nil, err
    }
    
    // Set data dari request
    a.NIM = req.NIM
    a.Nama = req.Nama
    a.Jurusan = req.Jurusan
    a.Angkatan = req.Angkatan
    a.TahunLulus = req.TahunLulus
    a.Email = req.Email
    a.NoTelepon = req.NoTelepon
    a.Alamat = req.Alamat
    
    return &a, nil
}

func (r *AlumniRepository) Update(id int, req *models.UpdateAlumniRequest) (*models.Alumni, error) {
    query := `
        UPDATE alumni 
        SET nama = $1, jurusan = $2, angkatan = $3, tahun_lulus = $4, 
            email = $5, no_telepon = $6, alamat = $7, updated_at = $8
        WHERE id = $9
    `
    
    now := time.Now()
    result, err := r.db.Exec(
        query, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
        req.Email, req.NoTelepon, req.Alamat, now, id,
    )
    
    if err != nil {
        return nil, err
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return nil, sql.ErrNoRows
    }
    
    // Ambil data yang sudah diupdate
    return r.GetByID(id)
}

func (r *AlumniRepository) Delete(id int) error {
    query := "DELETE FROM alumni WHERE id = $1"
    
    result, err := r.db.Exec(query, id)
    if err != nil {
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}

func (r *AlumniRepository) GetByFilter(search, sortBy, order string, limit, offset int) ([]models.FilterAlumni, error){
    validSort := map[string]string{
        "id": "id", "nama": "nama", "jurusan": "jurusan", "angkatan": "angkatan", "tahun_lulus": "tahun_lulus", "created_at": "created_at",
    }
    sortCol, ok := validSort[sortBy]
    if !ok {
    sortCol = "id"
	}

    ord := "ASC"
	if strings.ToUpper(order) == "DESC" {
		ord = "DESC"
	}

    query := fmt.Sprintf(`
        SELECT a.nama,
            a.jurusan,
            a.angkatan,
            a.tahun_lulus,
            pa.nama_perusahaan,
            pa.status_pekerjaan
        FROM alumni a
        LEFT JOIN pekerjaan_alumni pa ON pa.alumni_id = a.id
        WHERE a.nama ILIKE $1 
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
     `, sortCol, ord)

    rows, err := r.db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("query Error:", err)
		return nil, err
	}

	defer rows.Close()
    
    var filterAlumni []models.FilterAlumni
    for rows.Next() {
        var f models.FilterAlumni
        if err := rows.Scan(
            &f.Nama, 
            &f.Jurusan, 
            &f.Angkatan, 
            &f.TahunLulus,
            &f.NamaPerusahaan,
            &f.StatusPekerjaan,);
        err != nil {
            return nil, err
        }
        filterAlumni = append(filterAlumni, f)
    }
    return filterAlumni, nil
    
}

func (r *AlumniRepository) CountAlumniRepo(search string) (int,error){
    var total int
    countQuery :=`
    Select Count(*) From alumni Where nama ILIKE $1 or jurusan ILIKE $1 or CAST(angkatan AS TEXT) ILIKE $1 or CAST(tahun_lulus AS TEXT) ILIKE $1
    `
    err := r.db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
    if err != nil && err != sql.ErrNoRows {
        return 0, err
    }
    return total, nil
}