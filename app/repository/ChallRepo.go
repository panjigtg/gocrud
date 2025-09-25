package repository

import (
    "crudprojectgo/app/models"
    "database/sql"
	"fmt"
)

type ChallRepository struct {
	db *sql.DB
}

func NewChallRepository(db *sql.DB) *ChallRepository {
    return &ChallRepository{
        db: db,
    }
}

func (r *ChallRepository) GetAll() ([]models.Chall, error){
	query :=`
		SELECT 
			a.id,
			a.nama,
			a.jurusan,
			a.angkatan,
			p.bidang_industri,
			p.nama_perusahaan,
			p.posisi_jabatan,
			p.tanggal_mulai_kerja,
			p.gaji_range,
			p.tanggal_selesai_kerja,
			COUNT(p.id) OVER (PARTITION BY a.id) AS count_pekerjaan
		FROM alumni a
		JOIN pekerjaan_alumni p ON a.id = p.alumni_id
		WHERE p.tanggal_mulai_kerja <= NOW() - INTERVAL '1 YEAR'
	`
	rows, err := r.db.Query(query)
    if err != nil {
		fmt.Println("ERROR SCAN:", err)
        return nil, err
    }
    defer rows.Close()

	var Chall []models.Chall
    for rows.Next() {
        var c models.Chall
        err := rows.Scan(
			&c.ID,
			&c.Nama,
			&c.Jurusan,
			&c.Angkatan,
			&c.BidangIndustri,
			&c.NamaPerusahaan,
			&c.PosisiJabatan,
			&c.TanggalMulaiKerja,
			&c.GajiRange,
			&c.TanggalSelesaiKerja,
			&c.CountPekerjaan,
		)
	if err != nil {
		fmt.Println("🔥 ERROR SCAN:", err)
		return nil, err
	}
	Chall = append(Chall, c)
	}
	 return Chall, nil
}

