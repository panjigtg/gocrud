package repository

import (
	"crudprojectgo/app/models"
	"database/sql"
	// "errors"
	"time"
    "strings"
    "fmt"
    "log"
)

type PekerjaanRepository struct {
    db *sql.DB
}

func NewPekerjaanRepository(db *sql.DB) *PekerjaanRepository {
    return &PekerjaanRepository{
        db: db,
    }
}

func (r *PekerjaanRepository) GetAll() ([]models.PekerjaanAlumni, error) {
    query := `
        SELECT p.id, p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, 
               p.bidang_industri, p.lokasi_kerja, p.gaji_range, 
               p.tanggal_mulai_kerja, p.tanggal_selesai_kerja, 
               p.status_pekerjaan, p.deskripsi_pekerjaan, 
               p.created_at, p.updated_at, a.nama as alumni_nama
        FROM pekerjaan_alumni p
        LEFT JOIN alumni a ON p.alumni_id = a.id
        ORDER BY p.created_at DESC
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var pekerjaan []models.PekerjaanAlumni
    for rows.Next() {
        var p models.PekerjaanAlumni
        err := rows.Scan(
            &p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
            &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
            &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
            &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
            &p.CreatedAt, &p.UpdatedAt, &p.AlumniNama,
        )
        if err != nil {
            return nil, err
        }
        pekerjaan = append(pekerjaan, p)
    }
    
    return pekerjaan, nil
}

func (r *PekerjaanRepository) GetByID(id int) (*models.PekerjaanAlumni, error) {
    query := `
        SELECT p.id, p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, 
               p.bidang_industri, p.lokasi_kerja, p.gaji_range, 
               p.tanggal_mulai_kerja, p.tanggal_selesai_kerja, 
               p.status_pekerjaan, p.deskripsi_pekerjaan, 
               p.created_at, p.updated_at, a.nama as alumni_nama
        FROM pekerjaan_alumni p
        LEFT JOIN alumni a ON p.alumni_id = a.id
        WHERE p.id = $1
    `
    
    var p models.PekerjaanAlumni
    err := r.db.QueryRow(query, id).Scan(
        &p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
        &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
        &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
        &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
        &p.CreatedAt, &p.UpdatedAt, &p.AlumniNama,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &p, nil
}

func (r *PekerjaanRepository) GetByAlumniID(alumniID int) ([]models.PekerjaanAlumni, error) {
    query := `
        SELECT p.id, p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, 
               p.bidang_industri, p.lokasi_kerja, p.gaji_range, 
               p.tanggal_mulai_kerja, p.tanggal_selesai_kerja, 
               p.status_pekerjaan, p.deskripsi_pekerjaan, 
               p.created_at, p.updated_at, a.nama as alumni_nama
        FROM pekerjaan_alumni p
        LEFT JOIN alumni a ON p.alumni_id = a.id
        WHERE p.alumni_id = $1
        ORDER BY p.tanggal_mulai_kerja DESC
    `
    
    rows, err := r.db.Query(query, alumniID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var pekerjaan []models.PekerjaanAlumni
    for rows.Next() {
        var p models.PekerjaanAlumni
        err := rows.Scan(
            &p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
            &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
            &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
            &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
            &p.CreatedAt, &p.UpdatedAt, &p.AlumniNama,
        )
        if err != nil {
            return nil, err
        }
        pekerjaan = append(pekerjaan, p)
    }
    
    return pekerjaan, nil
}

func (r *PekerjaanRepository) Create(req *models.CreatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
    // Parse tanggal
    tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
    if err != nil {
        return nil, err
    }
    
    var tanggalSelesai *time.Time
    if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
        t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
        if err != nil {
            return nil, err
        }
        tanggalSelesai = &t
    }
    
    query := `
        INSERT INTO pekerjaan_alumni 
        (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
         lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
         status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id, created_at, updated_at
    `
    
    now := time.Now()
    var p models.PekerjaanAlumni
    
    err = r.db.QueryRow(
        query, req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan,
        req.BidangIndustri, req.LokasiKerja, req.GajiRange,
        tanggalMulai, tanggalSelesai, req.StatusPekerjaan,
        req.DeskripsiPekerjaan, now, now,
    ).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
    
    if err != nil {
        return nil, err
    }
    
    // Set data dari request
    p.AlumniID = req.AlumniID
    p.NamaPerusahaan = req.NamaPerusahaan
    p.PosisiJabatan = req.PosisiJabatan
    p.BidangIndustri = req.BidangIndustri
    p.LokasiKerja = req.LokasiKerja
    p.GajiRange = req.GajiRange
    p.TanggalMulaiKerja = tanggalMulai
    p.TanggalSelesaiKerja = tanggalSelesai
    p.StatusPekerjaan = req.StatusPekerjaan
    p.DeskripsiPekerjaan = req.DeskripsiPekerjaan
    
    return &p, nil
}

func (r *PekerjaanRepository) Update(id int, req *models.UpdatePekerjaanRequest) (*models.PekerjaanAlumni, error) {
    // Parse tanggal
    tanggalMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
    if err != nil {
        return nil, err
    }
    
    var tanggalSelesai *time.Time
    if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
        t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
        if err != nil {
            return nil, err
        }
        tanggalSelesai = &t
    }
    
    query := `
        UPDATE pekerjaan_alumni 
        SET nama_perusahaan = $1, posisi_jabatan = $2, bidang_industri = $3,
            lokasi_kerja = $4, gaji_range = $5, tanggal_mulai_kerja = $6,
            tanggal_selesai_kerja = $7, status_pekerjaan = $8, 
            deskripsi_pekerjaan = $9, updated_at = $10
        WHERE id = $11
    `
    
    now := time.Now()
    result, err := r.db.Exec(
        query, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri,
        req.LokasiKerja, req.GajiRange, tanggalMulai, tanggalSelesai,
        req.StatusPekerjaan, req.DeskripsiPekerjaan, now, id,
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

func (r *PekerjaanRepository) Delete(id int) error {
    query := "DELETE FROM pekerjaan_alumni WHERE id = $1"
    
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

func (r *PekerjaanRepository) GetByFilter(search string, sortCol, order string, page, limit, offset int) ([]models.FilterPekerjaan, error) {
    // Validasi dan pemilihan ASC/DESC
    ord := "ASC"
    if strings.ToUpper(order) == "DESC" {
        ord = "DESC"
    }

    // Query final (sortCol sudah divalidasi di service)
    query := fmt.Sprintf(`
        SELECT 
            a.nama AS nama_alumni,
            p.nama_perusahaan,
            p.posisi_jabatan,
            p.bidang_industri,
            p.lokasi_kerja,
            p.gaji_range,
            p.status_pekerjaan,   
            EXTRACT(YEAR FROM p.tanggal_mulai_kerja) AS tahun_mulai,
            EXTRACT(YEAR FROM p.tanggal_selesai_kerja) AS tahun_selesai
        FROM pekerjaan_alumni p
        LEFT JOIN alumni a ON p.alumni_id = a.id
        WHERE p.nama_perusahaan ILIKE $1 
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortCol, ord)

    log.Println("📌 FINAL QUERY:", query)

    rows, err := r.db.Query(query, "%"+search+"%", limit, offset)
    if err != nil {
        log.Println("❌ query error:", err)
        return nil, err
    }
    defer rows.Close()

    var filterPekerjaan []models.FilterPekerjaan
    for rows.Next() {
        var f models.FilterPekerjaan
        var tahunMulai sql.NullFloat64
        var tahunSelesai sql.NullFloat64

        if err := rows.Scan(
            &f.Nama,
            &f.NamaPerusahaan,
            &f.PosisiJabatan,
            &f.BidangIndustri,
            &f.LokasiKerja,
            &f.GajiRange,
            &f.StatusPekerjaan,
            &tahunMulai,
            &tahunSelesai,
        ); err != nil {
            log.Println("❌ scan error:", err)
            return nil, err
        }

        if tahunMulai.Valid {
            tahun := int(tahunMulai.Float64)
            f.TahunMulaiKerja = &tahun
        }

        if tahunSelesai.Valid {
            tahun := int(tahunSelesai.Float64)
            f.TahunSelesaiKerja = &tahun
        }

        filterPekerjaan = append(filterPekerjaan, f)
    }

    if err := rows.Err(); err != nil {
        log.Println("❌ rows.Err error:", err)
        return nil, err
    }

    log.Println("✅ Total ditemukan:", len(filterPekerjaan))
    return filterPekerjaan, nil
}


func (r *PekerjaanRepository) CountPekerjaanRepo(search string) (int, error) {
    var total int
    query := `
        SELECT COUNT(*)
        FROM pekerjaan_alumni p
        WHERE p.nama_perusahaan ILIKE $1
        OR p.posisi_jabatan ILIKE $1
        OR p.bidang_industri ILIKE $1
        OR p.lokasi_kerja ILIKE $1
        OR p.status_pekerjaan ILIKE $1
    `
    err := r.db.QueryRow(query, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
         log.Println("Count error:", err)
	return 0, err
	}
	return total, nil
}