package models

import "time"

type Chall struct {
	ID 					int 	 	`json:"id"`
	Nama 				string 		`json:"nama"`
	Jurusan 			string 		`json:"jurusan"`
	Angkatan 			int 		`json:"angkatan"`
	BidangIndustri      string    	`json:"bidang_industri"`
	NamaPerusahaan      string    	`json:"nama_perusahaan"`
	PosisiJabatan       string     	`json:"posisi_jabatan"`
	TanggalMulaiKerja   time.Time   `json:"tanggal_mulai_kerja"`
	GajiRange           *string    	`json:"gaji_range"`
	TanggalSelesaiKerja *time.Time  `json:"tanggal_selesai_kerja"`
	CountPekerjaan 		int			`json:"count_pekerjaan"`
}



